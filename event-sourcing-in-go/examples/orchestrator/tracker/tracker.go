package main

import (
	"fmt"

	"github.com/oklog/ulid"

	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"

	"github.com/Shopify/sarama"
)

const ConsumerPartition = 0

type Tracker struct {
	Logger       log.Logger
	Client       sarama.Client
	TrackerTopic string
	TrackerGroup string
	WorkerTopic  string

	consumer               sarama.Consumer
	producer               sarama.SyncProducer
	partitionOffsetManager sarama.PartitionOffsetManager
	offset                 int64

	state map[batchID]*batch
}

func (t *Tracker) Run() error {
	var err error

	// Init

	t.consumer, err = sarama.NewConsumerFromClient(t.Client)
	if err != nil {
		return fmt.Errorf("could not create consumer: %s", err)
	}
	defer t.consumer.Close()

	t.producer, err = sarama.NewSyncProducerFromClient(t.Client)
	if err != nil {
		return fmt.Errorf("could not create producer: %s", err)
	}
	defer t.producer.Close()

	om, err := sarama.NewOffsetManagerFromClient(t.TrackerGroup, t.Client)
	if err != nil {
		return fmt.Errorf("could not create offset manager: %s", err)
	}
	defer om.Close()

	if t.partitionOffsetManager, err = om.ManagePartition(t.TrackerTopic, ConsumerPartition); err != nil {
		return fmt.Errorf("could not create partition offset manager: %s", err)
	}
	defer t.partitionOffsetManager.Close()

	t.state = make(map[batchID]*batch)

	startOffset, _ := t.partitionOffsetManager.NextOffset()

	partitionConsumer, err := t.consumer.ConsumePartition(t.TrackerTopic, ConsumerPartition, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("could not consumer partition: %s", err)
	}
	defer partitionConsumer.Close()

	// Main loop

	for {
		inMsg, open := <-partitionConsumer.Messages()

		if !open {
			return fmt.Errorf("partition closed unexpectedly")
		}

		t.offset = inMsg.Offset

		id, event, err := unmarshalMessage(inMsg)
		if err != nil {
			t.Logger.Log("offset", inMsg.Offset, "msg", "invalid message", "err", err)
			continue
		}

		outMsgs, err := t.updateState(id, event)
		if err != nil {
			return fmt.Errorf("could not process event message: %s", err)
		}

		// Has this message been processed by a previous tracker instance?
		if inMsg.Offset >= startOffset {
			// No, this message has not been processed before -> Do follow-up actions

			if err := t.producer.SendMessages(outMsgs); err != nil {
				return fmt.Errorf("could not send event message: %s", err)
			}

			// Mark after send to ensure least-once guarantee
			t.partitionOffsetManager.MarkOffset(inMsg.Offset, "")
		}
	}
}

func unmarshalMessage(inMsg *sarama.ConsumerMessage) (batchID, orchestrator.Event, error) {
	id := ulid.ULID{}
	event := orchestrator.Event{}

	if err := proto.Unmarshal(inMsg.Value, &event); err != nil {
		return batchID(id), event, fmt.Errorf("could not unmarshal event: %s", err)
	}

	if err := id.UnmarshalBinary(event.GetBatchId()); err != nil {
		return batchID(id), event, fmt.Errorf("could not unmarshal batch ID: %s", err)
	}

	return batchID(id), event, nil
}

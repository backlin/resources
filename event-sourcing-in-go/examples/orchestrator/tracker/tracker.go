package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/go-kit/kit/log"
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

	state map[orchestrator.BatchID]*batch
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

	t.state = make(map[orchestrator.BatchID]*batch)

	startOffset, _ := t.partitionOffsetManager.NextOffset()

	partitionConsumer, err := t.consumer.ConsumePartition(t.TrackerTopic, ConsumerPartition, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("could not consumer partition: %s", err)
	}
	defer partitionConsumer.Close()

	// Main loop

	stateRecovered := false
	for {
		inMsg, open := <-partitionConsumer.Messages()

		if !open {
			return fmt.Errorf("partition closed unexpectedly")
		}

		t.offset = inMsg.Offset

		id, event, err := orchestrator.UnmarshalEvent(inMsg)
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

			if !stateRecovered {
				stateRecovered = true
				t.Logger.Log("offset", inMsg.Offset, "msg", "state recovered")
			}

			if len(outMsgs) > 0 {
				if err := t.producer.SendMessages(outMsgs); err != nil {
					return fmt.Errorf("could not send event message: %s", err)
				}
			}

			t.partitionOffsetManager.MarkOffset(inMsg.Offset, "")
			// Mark after send, not before, to ensure least-once guarantee
		}
	}
}

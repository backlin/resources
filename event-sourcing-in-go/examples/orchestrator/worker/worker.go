package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/bsm/sarama-cluster"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
)

type Worker struct {
	Logger       log.Logger
	Consumer     *cluster.Consumer
	Producer     sarama.SyncProducer
	TrackerTopic string
}

func (w Worker) Run() error {
	for {

		inMsg, open := <-w.Consumer.Messages()

		if !open {
			w.Logger.Log("consumer closed, shutting down")
			return nil
		}

		w.Logger.Log("offset", inMsg.Offset, "msg", "processing")

		outMsg, err := processMessage(inMsg, w.TrackerTopic)
		if err != nil {
			w.Logger.Log("offset", inMsg.Offset, "err", err, "msg", "failed processing work message")
			continue
		}

		// Send with simple retry logic
		for {
			_, _, err := w.Producer.SendMessage(outMsg)
			if err == nil {
				break
			}

			w.Logger.Log("offset", inMsg.Offset, "err", err, "msg", "failed sending event message")
			time.Sleep(10 * time.Second)
		}

		w.Consumer.MarkOffset(inMsg, "")

	}
}

func processMessage(inMsg *sarama.ConsumerMessage, trackerTopic string) (*sarama.ProducerMessage, error) {
	work := &orchestrator.Work{}
	if err := proto.Unmarshal(inMsg.Value, work); err != nil {
		return nil, fmt.Errorf("could not unmarshal message: %s", err)
		// You might want to consider failing softly here instead
	}

	time.Sleep(time.Duration(work.Duration)) // The requested work to be done

	value := orchestrator.Event{
		BatchId:     work.BatchId,
		JobId:       work.JobId,
		StatusLevel: orchestrator.Event_JOB,
		Status:      orchestrator.Event_SUCCESS,
	}

	b, err := value.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal %s job event: %s", value.Status, err)
		// The application is broken, always fail hard
	}

	outMsg := &sarama.ProducerMessage{
		Topic: trackerTopic,
		Value: sarama.ByteEncoder(b),
	}

	return outMsg, nil
}
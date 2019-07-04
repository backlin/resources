package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/bsm/sarama-cluster"
	"github.com/go-kit/kit/log"
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

		id, work, err := orchestrator.UnmarshalWork(inMsg)
		if err != nil {
			w.Logger.Log("offset", inMsg.Offset, "msg", "invalid message", "err", err)
			continue
		}

		w.Logger.Log("offset", inMsg.Offset, "batch_id", id, "job_id", work.JobId, "msg", "processing")

		outMsg, err := processWork(work, w.TrackerTopic)
		if err != nil {
			w.Logger.Log("offset", inMsg.Offset, "batch_id", id, "job_id", work.JobId, "err", err, "msg", "failed processing work message")
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

		w.Logger.Log("offset", inMsg.Offset, "batch_id", id, "job_id", work.JobId, "msg", "SUCCESS")
	}
}

func processWork(work orchestrator.Work, trackerTopic string) (*sarama.ProducerMessage, error) {
	// Do the actual work
	time.Sleep(time.Duration(work.Duration) * time.Millisecond)

	outStatus := orchestrator.Event_SUCCESS
	if rand.Float32() < work.FailureRate {
		outStatus = orchestrator.Event_FAILURE
	}

	// Package the output
	value := orchestrator.Event{
		BatchId:     work.BatchId,
		JobId:       work.JobId,
		StatusLevel: orchestrator.Event_JOB,
		Status:      outStatus,
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

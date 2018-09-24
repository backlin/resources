package main

import (
	"fmt"
	"math/rand"

	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
)

type batch struct {
	jobCount     int32
	successCount int32
	failureCount int32
}

func (t *Tracker) updateState(id orchestrator.BatchID, event orchestrator.Event) ([]*sarama.ProducerMessage, error) {
	if event.GetStatusLevel() == orchestrator.Event_JOB {
		outMsg, err := t.updateJob(id, event)
		if err != nil {
			return nil, fmt.Errorf("failed updating job: %s", err)
		}

		if outMsg == nil {
			return nil, nil
		}

		return []*sarama.ProducerMessage{outMsg}, nil
	}

	// Batch event
	outMsgs, err := t.updateBatch(id, event)
	if err != nil {
		return nil, fmt.Errorf("failed updating batch: %s", err)
	}

	return outMsgs, nil
}

func (t *Tracker) updateBatch(id orchestrator.BatchID, event orchestrator.Event) ([]*sarama.ProducerMessage, error) {
	t.Logger.Log("offset", t.offset, "batch_id", id, "status", event.GetStatus())

	switch event.GetStatus() {
	case orchestrator.Event_PENDING:
		return t.runBatch(id, event.GetJobCount())

	case orchestrator.Event_SUCCESS:
		// Run downstream dependencies

	case orchestrator.Event_FAILURE:
		// Cancel downstream dependencies
	}

	return nil, nil
}

func (t *Tracker) runBatch(id orchestrator.BatchID, jobCount int32) ([]*sarama.ProducerMessage, error) {
	if jobCount == 0 {
		t.Logger.Log("offset", t.offset, "batch_id", id, "err", "batch has no jobs, no point running it")
		return nil, nil
	}

	t.state[id] = &batch{
		jobCount: jobCount,
	}

	outMsgs := make([]*sarama.ProducerMessage, jobCount+1)

	// Create work messages
	for i := int32(0); i < jobCount; i++ {
		work := &orchestrator.Work{
			BatchId:  id.MustMarshalBinary(),
			JobId:    i,
			Duration: rand.Int63n(4000) + 2000,
		}

		b, err := work.Marshal()
		if err != nil {
			return nil, fmt.Errorf("could not marshal work message: %s", err)
		}
		outMsgs[i] = &sarama.ProducerMessage{
			Topic: t.WorkerTopic,
			Value: sarama.ByteEncoder(b),
		}
	}

	// Create updated event message
	newEvent := orchestrator.Event{
		BatchId:     id.MustMarshalBinary(),
		StatusLevel: orchestrator.Event_BATCH,
		Status:      orchestrator.Event_RUNNING,
	}

	b, err := newEvent.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal event message: %s", err)
	}

	outMsgs[jobCount] = &sarama.ProducerMessage{
		Topic: t.TrackerTopic,
		Value: sarama.ByteEncoder(b),
	}

	return outMsgs, nil
}

func (t *Tracker) updateJob(id orchestrator.BatchID, event orchestrator.Event) (*sarama.ProducerMessage, error) {
	b, ok := t.state[id]
	if !ok {
		return nil, fmt.Errorf("previously unseen batch")
	}

	newBatchStatus := b.updateJob(event.GetStatus())

	if newBatchStatus == orchestrator.Event_RUNNING {
		return nil, nil
	}

	// Batch is terminated (SUCCESS or FAILURE)

	outEvent := &orchestrator.Event{
		BatchId:     id.MustMarshalBinary(),
		StatusLevel: orchestrator.Event_BATCH,
		Status:      newBatchStatus,
	}

	bEvent, err := outEvent.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed marshaling event message: %s", err)
	}

	outMsg := &sarama.ProducerMessage{
		Topic: t.TrackerTopic,
		Value: sarama.ByteEncoder(bEvent),
	}

	return outMsg, nil
}

func (b *batch) updateJob(status orchestrator.Event_Status) orchestrator.Event_Status {
	switch status {

	case orchestrator.Event_SUCCESS:
		b.successCount++

		if b.successCount == b.jobCount {
			return orchestrator.Event_SUCCESS
		}

	case orchestrator.Event_FAILURE:
		b.failureCount++

		return orchestrator.Event_FAILURE
	}

	return orchestrator.Event_RUNNING
}

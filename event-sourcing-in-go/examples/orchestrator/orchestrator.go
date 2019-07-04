package orchestrator

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	"github.com/oklog/ulid"
)

type BatchID ulid.ULID

func (id BatchID) String() string {
	return ulid.ULID(id).String()
}

func (id BatchID) MustMarshalBinary() sarama.ByteEncoder {
	b, err := ulid.ULID(id).MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("failed to marshal ULID: %s", err))
	}

	return sarama.ByteEncoder(b)
}

func UnmarshalEvent(inMsg *sarama.ConsumerMessage) (BatchID, Event, error) {
	id := ulid.ULID{}
	event := Event{}

	if err := proto.Unmarshal(inMsg.Value, &event); err != nil {
		return BatchID(id), event, fmt.Errorf("could not unmarshal event: %s", err)
	}

	if err := id.UnmarshalBinary(event.GetBatchId()); err != nil {
		return BatchID(id), event, fmt.Errorf("could not unmarshal batch ID: %s", err)
	}

	return BatchID(id), event, nil
}

func UnmarshalWork(inMsg *sarama.ConsumerMessage) (BatchID, Work, error) {
	id := ulid.ULID{}
	work := Work{}

	if err := proto.Unmarshal(inMsg.Value, &work); err != nil {
		return BatchID(id), work, fmt.Errorf("could not unmarshal work: %s", err)
	}

	if err := id.UnmarshalBinary(work.GetBatchId()); err != nil {
		return BatchID(id), work, fmt.Errorf("could not unmarshal batch ID: %s", err)
	}

	return BatchID(id), work, nil
}

func DefaultConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	return config
}

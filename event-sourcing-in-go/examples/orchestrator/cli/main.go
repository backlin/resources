package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/stevenroose/gonfig"

	"github.com/Shopify/sarama"

	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/oklog/ulid"
)

var config = struct {
	BrokerAddresses []string
	TrackerTopic    string
	Create          bool  `id:"create",validate:"required"`
	JobCount        int32 `id:"job-count",validate:"gt=0"`
}{
	BrokerAddresses: []string{"localhost:9092"},
	TrackerTopic:    "tracker",
}

func main() {
	if err := gonfig.Load(&config, gonfig.Conf{}); err != nil {
		fmt.Printf("error reading arguments: %s", err)
		os.Exit(1)
	}

	if !config.Create {
		fmt.Println("You must use argument --create (future arguments to be added)")
		os.Exit(1)
	}

	if config.Create && config.JobCount <= 0 {
		fmt.Printf("Batch must have > 0 jobs, %d requested\n", config.JobCount)
		os.Exit(1)
	}

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	bID, err := id.MarshalBinary()
	if err != nil {
		fmt.Printf("internal error: %s\n", err)
		os.Exit(1)
	}

	event := &orchestrator.Event{
		BatchId:     bID,
		StatusLevel: orchestrator.Event_BATCH,
		Status:      orchestrator.Event_PENDING,
		JobCount:    config.JobCount,
	}

	bEvent, err := event.Marshal()
	if err != nil {
		fmt.Printf("internal error: %s\n", err)
		os.Exit(1)
	}

	outMsg := &sarama.ProducerMessage{
		Topic: config.TrackerTopic,
		Value: sarama.ByteEncoder(bEvent),
	}

	client, err := sarama.NewClient(config.BrokerAddresses, orchestrator.DefaultConfig())
	if err != nil {
		fmt.Printf("could not create client: %s\n", err)
		os.Exit(1)
	}
	defer client.Close()

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		fmt.Printf("could not create producer: %s\n", err)
		os.Exit(1)
	}
	defer producer.Close()

	if _, _, err := producer.SendMessage(outMsg); err != nil && false {
		fmt.Printf("internal error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Produced %s batch %s with %d jobs.\n", orchestrator.Event_PENDING, id, config.JobCount)
}

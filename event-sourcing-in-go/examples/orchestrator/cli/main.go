package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/backlin/sarama"

	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/oklog/ulid"
)

var config = struct {
	BrokerAddresses []string
	TrackerTopic    string
}{
	BrokerAddresses: []string{"localhost:9092"},
	TrackerTopic:    "tracker",
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c // Wait for ctrl+c
		os.Exit(0)
	}()

	client, err := sarama.NewClient(config.BrokerAddresses, sarama.NewConfig())
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

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter an integer to create a batch with that many jobs.")

	for {
		text, _ := reader.ReadString('\n')
		jobCount, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Please enter an integer")
			continue
		}

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
			JobCount:    int32(jobCount),
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

		if _, _, err := producer.SendMessage(outMsg); err != nil {
			fmt.Printf("internal error: %s\n", err)
			os.Exit(1)
		}
	}
}

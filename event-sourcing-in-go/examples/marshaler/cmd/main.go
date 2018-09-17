package main

import (
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/marshaler"
	"github.com/go-kit/kit/log"
)

var config = struct {
	BrokerAddresses []string
	InputTopic      string
	OutputTopic     string
	ConsumerGroup   string
}{
	BrokerAddresses: []string{"localhost:9092"},
	InputTopic:      "data_flatfile",
	OutputTopic:     "data_proto",
	ConsumerGroup:   "marshaler",
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	saramaConfig := sarama.NewConfig()

	client, err := sarama.NewClient(config.BrokerAddresses, saramaConfig)
	if err != nil {
		logger.Log("msg", "could not create client", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	offsetManager, err := sarama.NewOffsetManagerFromClient(config.ConsumerGroup, client)
	if err != nil {
		logger.Log("msg", "could not create offset manager", "err", err)
		os.Exit(1)
	}
	defer offsetManager.Close()

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Log("msg", "could not create producer", "err", err)
		os.Exit(1)
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		logger.Log("msg", "could not create consumer", "err", err)
		os.Exit(1)
	}

	/*
		Set up program termination on receiving SIGINT (ctrl+c).
		Note that this only works if you build the program and run the binary,
		but not if you run it by calling `go run`.
	*/
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c // Wait for ctrl+c
		consumer.Close()
	}()

	var filestorage marshaler.FileRetriever
	// TODO: Set up file retreival here

	marshaler := marshaler.Marshaler{
		Consumer:      consumer,
		Producer:      producer,
		OffsetManager: offsetManager,
		FileStorage:   filestorage,
		Logger:        logger,
	}

	marshaler.Run(config.InputTopic, config.OutputTopic)
}

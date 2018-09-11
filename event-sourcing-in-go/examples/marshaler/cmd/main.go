package main

import (
	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/marshaler"
	"github.com/go-kit/kit/log"
	"os"
)
var config = struct {
	BrokerAddresses []string
	InputTopic string
	OutputTopic string
}{
	BrokerAddresses: []string{"localhost:9092"},
	InputTopic: "data_flatfile",
	OutputTopic: "data_proto",
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	saramaConfig := sarama.NewConfig()

	client, err := sarama.NewClient(config.BrokerAddresses, saramaConfig)
	if err != nil {
		logger.Log("msg", "could not create client", "err", err)
		os.Exit(1)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		logger.Log("msg", "could not create consumer", "err", err)
		os.Exit(1)
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Log("msg", "could not create producer", "err", err)
		os.Exit(1)
	}

	var filestorage marshaler.FileRetriever
	// TODO: Set up file retreival here

	marshaler := marshaler.Marshaler{
		Consumer:    consumer,
		Producer:    producer,
		FileStorage: filestorage,
		Logger:      logger,
	}

	stop := make(chan struct{})
	go marshaler.Run(config.InputTopic, config.OutputTopic, sarama.OffsetOldest, stop)

	// Optionally, implement `stop <- struct{}` when receiving external request
}

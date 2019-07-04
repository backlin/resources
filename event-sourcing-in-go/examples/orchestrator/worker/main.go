package main

import (
	"os"

	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/go-kit/kit/log"
)

var config = struct {
	BrokerAddresses []string `id:"broker-addresses" validate:"required"`
	TrackerTopic    string   `id:"tracker-topic" validate:"required"`
	WorkerTopic     string   `id:"worker-topic" validate:"required"`
	WorkerGroup     string   `id:"worker-group" validate:"required"`
}{
	BrokerAddresses: []string{"localhost:9092"},
	TrackerTopic:    "tracker",
	WorkerTopic:     "worker",
	WorkerGroup:     "worker",
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	consumer, err := cluster.NewConsumer(config.BrokerAddresses, config.WorkerGroup, []string{config.WorkerTopic}, cluster.NewConfig())

	client, err := sarama.NewClient(config.BrokerAddresses, orchestrator.DefaultConfig())
	if err != nil {
		logger.Log("msg", "could not create client", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Log("msg", "could not create producer", "err", err)
		os.Exit(1)
	}
	defer producer.Close()

	worker := Worker{
		Logger:       logger,
		TrackerTopic: config.TrackerTopic,
		Consumer:     consumer,
		Producer:     producer,
	}

	logger.Log("msg", "running worker")
	if err := worker.Run(); err != nil {
		logger.Log("msg", "critical failure, shutting down tracker", "err", err)
		os.Exit(1)
	}
}

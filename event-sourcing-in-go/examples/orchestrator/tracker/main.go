package main

import (
	"os"

	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/log"
)

var config = struct {
	BrokerAddresses []string `id:"broker-addresses",validate:"required"`
	TrackerTopic    string   `id:"tracker-topic",validate:"required"`
	TrackerGroup    string   `id:"tracker-group",validate:"required"`
	WorkerTopic     string   `id:"worker-topic",validate:"required"`
}{
	BrokerAddresses: []string{"localhost:9092"},
	TrackerTopic:    "tracker",
	TrackerGroup:    "tracker",
	WorkerTopic:     "worker",
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

	tr := Tracker{
		Logger:       logger,
		Client:       client,
		TrackerTopic: config.TrackerTopic,
		TrackerGroup: config.TrackerGroup,
		WorkerTopic:  config.WorkerTopic,
	}

	if err := tr.Run(); err != nil {
		logger.Log("msg", "critical failure, shutting down tracker", "err", err)
		os.Exit(1)
	}
}

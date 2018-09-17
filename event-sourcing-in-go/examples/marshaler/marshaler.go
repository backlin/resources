package marshaler

import (
	"fmt"
	"io"

	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type Marshaler struct {
	Consumer      sarama.Consumer
	Producer      sarama.SyncProducer
	OffsetManager sarama.OffsetManager
	FileStorage   FileRetriever
	Logger        log.Logger
}

type FileRetriever interface {
	Get(url string) (io.Reader, error)
}

func (m Marshaler) Run(inputTopic, outputTopic string) error {
	const partition = 0

	offset := sarama.OffsetOldest
	var (
		pom sarama.PartitionOffsetManager
		err error
	)
	if m.OffsetManager != nil {
		pom, err = m.OffsetManager.ManagePartition(inputTopic, partition)
		if err != nil {
			return fmt.Errorf("could not manage offset of %q, partition %d", inputTopic, partition)
		}
		defer pom.Close()

		offset, _ = pom.NextOffset()
	}

	pc, err := m.Consumer.ConsumePartition(inputTopic, partition, offset)
	if err != nil {
		return fmt.Errorf("could not consume topic %q, partition %d, offset %d", inputTopic, partition, offset)
	}
	defer pc.Close()

	for {
		inMsg, open := <-pc.Messages()

		if !open {
			m.Logger.Log("partition consumer closed, shutting down")
			return nil
		}

		m.Logger.Log("msg", fmt.Sprintf("processing message (offset %d)", inMsg.Offset))

		outMsgs, err := m.processMessage(inMsg, outputTopic)
		if err != nil {
			return fmt.Errorf("could not process message: %s", err)
		}

		if err := m.Producer.SendMessages(outMsgs); err != nil {
			return fmt.Errorf("could not send message: %s", err)
		}

		if pom != nil {
			pom.MarkOffset(inMsg.Offset+1, "")
		}
	}
}

func (m Marshaler) processMessage(inMsg *sarama.ConsumerMessage, outputTopic string) ([]*sarama.ProducerMessage, error) {
	rawItems := &RawItemSet{}
	if err := proto.Unmarshal(inMsg.Value, rawItems); err != nil {
		return nil, fmt.Errorf("could not unmarshal message: %s", err)
	}

	rawFile, err := m.FileStorage.Get(rawItems.Url)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve raw file %q: %s", rawItems.Url, err)
	}

	items := &ItemSet{}
	if err := jsonpb.Unmarshal(rawFile, items); err != nil {
		return nil, fmt.Errorf("could not unmarshal raw file: %s", err)
	}

	outMsgs := make([]*sarama.ProducerMessage, len(items.Items))

	for i, item := range items.Items {
		b, err := item.Marshal()
		if err != nil {
			return nil, fmt.Errorf("could not marshal output: %s", err)
		}

		outMsgs[i] = &sarama.ProducerMessage{
			Topic: outputTopic,
			Value: sarama.ByteEncoder(b),
		}
	}

	return outMsgs, nil
}

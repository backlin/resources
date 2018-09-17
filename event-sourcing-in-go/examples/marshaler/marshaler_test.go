package marshaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

const (
	testInputTopic  = "test_in"
	testOutputTopic = "test_out"
)

func TestMarshaler(t *testing.T) {
	wantItem := &Item{
		Id:   1032,
		Name: "t-shirt",
		Attributes: []*Attribute{
			{Name: "size", Value: "small"},
			{Name: "color", Value: "blue"},
		},
	}

	// Make raw file
	testItems := &ItemSet{
		Items: []*Item{wantItem},
	}
	body, err := json.Marshal(testItems)
	if err != nil {
		t.Fatalf("failed to marshal raw file (test is broken): %s", err)
	}
	filestorage := mockFileStorage{body: body}

	// Make message that points to the raw file
	rawItems := &RawItemSet{
		Url: "s3://mock-bucket/path/body.json",
	}
	b, err := rawItems.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal raw item (test is broken): %s", err)
	}
	testMessage := &sarama.ConsumerMessage{
		Value: sarama.ByteEncoder(b),
	}

	// Mock Kafka
	config := sarama.NewConfig()
	consumer := mocks.NewConsumer(t, config)

	partitionConsumer := consumer.ExpectConsumePartition(testInputTopic, 0, 0)
	partitionConsumer.YieldMessage(testMessage)

	valueChecker := func(val []byte) error {
		defer consumer.Close()

		gotItem := &Item{}
		if err := proto.Unmarshal(val, gotItem); err != nil {
			return fmt.Errorf("could not unmarshal value written to Kafka: %s", err)
		}

		if !reflect.DeepEqual(wantItem, gotItem) {
			return fmt.Errorf("incorrect output: %s", cmp.Diff(wantItem, gotItem))
		}

		return nil
	}

	producer := mocks.NewSyncProducer(t, config)
	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(valueChecker)

	marshaler := Marshaler{
		Consumer:    consumer,
		Producer:    producer,
		FileStorage: filestorage,
		Logger:      testLogger{t: t},
	}

	if err := marshaler.Run(testInputTopic, testOutputTopic, 0); err != nil {
		t.Fatal(err)
	}
}

type mockFileStorage struct {
	body []byte
}

var _ FileRetriever = &mockFileStorage{}

func (r mockFileStorage) Get(url string) (io.Reader, error) {
	return bytes.NewReader(r.body), nil
}

type testLogger struct {
	t *testing.T
}

var _ log.Logger = &testLogger{}

func (l testLogger) Log(keyvals ...interface{}) error {
	l.t.Log(keyvals)
	return nil
}

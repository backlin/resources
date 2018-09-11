package marshaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/go-kit/kit/log"
	"github.com/google/go-cmp/cmp"
	"io"
	"reflect"
	"testing"
)

const (
	testInputTopic = "test_in"
	testOutputTopic = "test_out"
)

func TestMarshaler(t *testing.T) {
	wantItems := []*Item{
		{
			ID: 1032,
			Name: "t-shirt",
			Attributes: []*Attribute{
				{Name: "size", Value: &Attribute_StringValue{StringValue: "small"}},
				{Name: "color", Value: &Attribute_StringValue{StringValue: "blue"}},
			},
		},
	}

	b, err := json.Marshal(wantItems)
	if err != nil {
		t.Fatalf("test is broken: %s", err)
	}
	filestorage := mockFileStorage{body: b}

	testMessage := &sarama.ConsumerMessage{
		Value: sarama.ByteEncoder(b),
	}

	config := sarama.NewConfig()
	consumer := mocks.NewConsumer(t, config)

	partitionConsumer := consumer.ExpectConsumePartition(testInputTopic, 0, 0)
	partitionConsumer.YieldMessage(testMessage)

	stop := make(chan struct{})
	valueChecker := func(val []byte) error {
		defer func(){
			stop <- struct{}{}
		}()

		if !reflect.DeepEqual(val, wantItems) {
			return fmt.Errorf("incorrect output: %s", cmp.Diff(val, wantItems))
		}
		return nil
	}

	producer := mocks.NewSyncProducer(t, config)
	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(valueChecker)

	marshaler := Marshaler{
		Consumer: consumer,
		Producer: producer,
		Filestorage: filestorage,
		Logger: testLogger{t: t},
	}

	marshaler.Run(testInputTopic, testOutputTopic, 0, stop)
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

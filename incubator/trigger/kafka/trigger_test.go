package kafka

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"io/ioutil"
	"testing"
	"time"
)

var jsonMetadata = getJsonMetadata()
var ranAction = make(chan bool, 1)

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "name": "kafka",
  "settings": {
    "server": "127.0.0.1:9092",
    "configid": "test-flogo-trigger"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "topic": "test",
				"partition": "0"
      }
    },
		{
      "actionId": "local://testFlow2",
      "settings": {
        "topic": "flogo",
				"partition": "0"
      }
    }
  ]
}`

var kafkaAddrs = []string{"127.0.0.1:9092"}

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Infof("Ran Action: %v", uri)
	ranAction <- true
	return 0, nil, nil
}

/*func TestInit(t *testing.T) {
	log.Info("Testing Init")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	f := &kafkaFactory{}
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
} */

func TestEndpoint(t *testing.T) {
	log.Info("Testing Endpoint")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &kafkaFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()

	conf := kafka.NewBrokerConf("test-client")
	conf.AllowTopicCreation = true

	// connect to kafka cluster
	broker, err := kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		t.Error("Error connecting to kafka cluster")
		t.Fail()
	}
	defer broker.Close()

	producer := broker.Producer(kafka.NewProducerConf())
	message := "Test message from Flogo"
	msg := &proto.Message{Value: []byte(message)}
	log.Info("Publishing test message to topic 'test' on Kafka")
	if _, err := producer.Produce("test", 0, msg); err != nil {
		t.Error("Error sending message to topic test")
		t.Fail()
	}
	//broker.Close()
	log.Info("Waiting 5 seconds for the message to be handled by the trigger...")
	select {
		case <-ranAction:
			log.Debug("Message was handled OK by the trigger")
		case <-time.After(5 * time.Second):
			t.Error("No action called by trigger based on message")
			t.Fail()
			return
	}
	log.Info("Publishing test message to topic 'flogo' on Kafka")
	if _, err := producer.Produce("flogo", 0, msg); err != nil {
		t.Error("Error sending message to topic flogo")
		t.Fail()
	}
	//broker.Close()
	log.Info("Waiting 5 seconds for the message to be handled by the trigger...")
	select {
		case <-ranAction:
			log.Debug("Message was handled OK by the trigger")
		case <-time.After(5 * time.Second):
			t.Error("No action called by trigger based on message")
			t.Fail()
			return
	}

}

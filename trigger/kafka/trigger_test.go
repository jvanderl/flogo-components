package kafka

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"testing"
)

const testConfig string = `{
  "name": "kafka",
  "settings": {
    "server": "10.10.1.50:9092",
    "configid": "test-flogo-trigger",
    "topic": "test",
    "partition": "0"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "topic": "test"
      }
    }
  ]
}`

var kafkaAddrs = []string{"localhost:9092"}

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Infof("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestInit(t *testing.T) {
	log.Info("Testing Init")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	f := &kafkaFactory{}
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
}

func TestEndpoint(t *testing.T) {
	log.Info("Testing Endpoint")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &kafkaFactory{}
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
		log.Errorf("cannot connect to kafka cluster: %s", err)
	}
	defer broker.Close()

	producer := broker.Producer(kafka.NewProducerConf())
	message := "Test message from Flogo"
	msg := &proto.Message{Value: []byte(message)}
	log.Info("---- doing publish ----")
	if _, err := producer.Produce("test", 0, msg); err != nil {
		log.Errorf("cannot produce message to %s:%d: %s", "test", 0, err)
	}

	broker.Close()
	log.Info("Sample Publisher Disconnected")
}

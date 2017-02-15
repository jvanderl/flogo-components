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
    "server": "192.168.178.41:9092",
    "configid": "test-flogo-trigger",
    "topic": "test",
    "partition": "0"
  },
  "endpoints": [
    {
      "actionType": "flow",
      "actionURI": "local://testFlow",
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
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("kafka")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("kafka")

	runner := &TestRunner{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(config, runner)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("kafka")

	tgr.Start()
	defer tgr.Stop()

	conf := kafka.NewBrokerConf("test-client")
	conf.AllowTopicCreation = true

	// connect to kafka cluster
	broker, err := kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		log.Fatalf("cannot connect to kafka cluster: %s", err)
	}
	defer broker.Close()

	producer := broker.Producer(kafka.NewProducerConf())
	message := "Test message from Flogo"
	msg := &proto.Message{Value: []byte(message)}
	log.Debug("---- doing publish ----")
	if _, err := producer.Produce("test", 0, msg); err != nil {
		log.Fatalf("cannot produce message to %s:%d: %s", "test", 0, err)
	}

	broker.Close()
	log.Debug("Sample Publisher Disconnected")
}

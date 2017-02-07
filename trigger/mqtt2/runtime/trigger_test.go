package mqtt2

import (
	"context"
	"encoding/json"
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/eclipse/paho.mqtt.golang"
)

const testConfig string = `{
  "name": "mqtt2",
  "settings": {
    "topic": "flogo/#",
    "broker": "tcp://127.0.0.1:1883",
    "id": "flogoEngine",
    "user": "",
    "password": "",
    "store": "",
    "qos": "0",
    "cleansess": "false"
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "topic": "test_start"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("mqtt2")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("mqtt2")

	runner := &TestRunner{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(config, runner)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("mqtt2")

	tgr.Start()
	defer tgr.Stop()

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("flogoEngine")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Debug("---- doing publish ----")
	token := client.Publish("flogo/test_start", 0, false, "Test message payload!")
	token.Wait()

	client.Disconnect(250)
	log.Debug("Sample Publisher Disconnected")
}


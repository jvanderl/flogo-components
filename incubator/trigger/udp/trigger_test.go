package udp

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

// Run Once, Start Immediately
const testConfig1 string = `{
  "name": "udp",
  "settings": {
		"port": 22600,
		"multicast_group": "224.192.32.19"
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "handler_setting": "xxx"
      }
    }
  ]
}`

// Listen for F1-2017 data
const testConfig2 string = `{
  "name": "udp",
  "settings": {
		"port": 20777,
		"multicast_group": ""
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "handler_setting": "xxx"
      }
    }
  ]
}`

//192.168.1.19
type TestRunner struct {
}

var Test action.Runner

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Infof("Ran Action (Run): %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Infof("Ran Action (RunAction): %v", act)
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	log.Infof("Ran Action (Execute): %v", act)
	return nil, nil
}

func TestTimer(t *testing.T) {
	log.Info("Testing UDP")
	config := trigger.Config{}

	//  Owl PV monitor test
	json.Unmarshal([]byte(testConfig1), &config)

	// F1-2017 Telemtery
	//json.Unmarshal([]byte(testConfig2), &config)

	f := &udpTriggerFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
	log.Infof("Press CTRL-C to quit")
	for {
	}
}

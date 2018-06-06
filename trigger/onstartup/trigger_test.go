package onstartup

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
  "name": "onstartup",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow1",
      "settings": {
      }
    }
  ]
}`

// Multiple flows configurations
const testConfig2 string = `{
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow1",
      "settings": {
      }
    },
		{
      "actionId": "local://testFlow2",
      "settings": {
      }
    },
    {
      "actionId": "local://testFlow3",
      "settings": {
      }
    }
  ]
}`

type TestRunner struct {
}

var Test action.Runner

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action (Run): %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action (RunAction): %v", act)
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action (Execute): %v", act)
	return nil, nil
}

func TestSingle(t *testing.T) {
	log.Info("Testing OnStartup")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig1), &config)
	f := &OnStartupFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	//runner := &TestRunner{}
	//tgr.Initialize()
	//tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
	log.Infof("Press CTRL-C to quit")
	for {
	}
}

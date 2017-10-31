package tcm

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"io/ioutil"
	"testing"
)

var jsonMetadata = getJSONMetadata()

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "name": "tcm",
  "settings": {
    "url": "<Your TCM URL Here>",
		"authkey": "<Your TCM Auht Key Here>",
		"clientid": "flogo-testsubscriber",
    "certificate": ""
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "destinationname": "demo_tcm",
				"destinationmatch": "*",
				"messagename": "demo_tcm",
				"durable": "true",
				"durablename": "flogo_demo_tcm"
      }
    }
  ]
}`


type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Infof("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestEndpoint(t *testing.T) {
	log.Info("Testing Endpoint")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &tcmFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	// just loop
	for {}
}

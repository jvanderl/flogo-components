package gpio

import (
	"context"
	"encoding/json"
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"io/ioutil"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string{
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil{
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

// Run Once, Start Immediately
const testConfig1 string = `{
  "name": "gpio",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "gpiopin": "7",
				"state": "1"
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

func TestGpio(t *testing.T) {
	log.Info("Testing GPIO")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig1), &config)
	f := &GPIOFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
	log.Infof("Press CTRL-C to quit")
  for {}
}

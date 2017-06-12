package blemaster

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"context"
	"testing"
	"encoding/json"
	"io/ioutil"
	"time"
)

var jsonMetadata = getJSONMetadata()
var ranAction = make(chan bool, 1)

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

//		"deviceid": "00:15:83:00:97:DF", //CC41
// 		"deviceid": "A4:D5:78:6D:57:6C", //HM10

const testConfig string = `{
  "name": "blemaster",
  "settings": {
		"devicename": "IOTDEVICE",
		"deviceid": "00:15:83:00:97:DF",
		"autodisconnect": "true",
		"autoreconnect": "true",
		"reconnectinterval": "5",
		"intervaltype": "seconds"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "service": "ffe0",
				"characteristic": "ffe1"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Infof("Ran Action: %v", uri)
	ranAction <- true
	return 0, nil, nil
}

/*
func TestRegistered(t *testing.T) {
	act := trigger.Get("blemaster")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}
*/

func TestEndpoint(t *testing.T) {
	log.Info("Testing Endpoint")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &BleMasterFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	//generate some bt data here

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

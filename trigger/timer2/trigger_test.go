package timer2

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

const testConfig3 string = `{
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "false",
				"startImmediate": "true"
      }
    }
  ]
}`


const testConfig string = `{
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "false",
				"startDate" : "2017-06-14T01:41:00Z02:00"
      }
    }
  ]
}`

const testConfig2 string = `{
"name": "timer2",
"settings": {
},
"handlers": [
	{
		"actionId": "local://testFlow2",
		"settings": {
			"startImmediate": "true",
			"repeating": "true",
			"seconds": "5",
			"minutes": "0",
			"hours": "0"
		}
	}
]
}`

const testConfig4 string = `{
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "false",
        "startDate" : "2017-06-14T2:28:00Z02:00"
      }
    },
    {
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "true",
        "startDate" : "2017-06-14T2:28:00Z02:00",
				"seconds": "0",
				"minutes": "0",
        "hours": "24"
      }
    },
    {
      "actionId": "local://testFlow3",
      "settings": {
        "repeating": "true",
        "startDate" : "2017-06-14T2:28:00Z02:00",
				"seconds": "0",
        "minutes": "60",
				"hours": "0"
      }
    },
    {
      "actionId": "local://testFlow4",
      "settings": {
        "repeating": "true",
        "startDate" : "2017-06-14T2:28:00Z02:00",
        "seconds": "30",
				"minutes": "0",
				"hours": "0"
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

func TestTimer(t *testing.T) {
	log.Info("Testing Timer")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig2), &config)
	// New  factory
	f := &Timer2Factory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
  for {}
	log.Infof("Test timer done")
}

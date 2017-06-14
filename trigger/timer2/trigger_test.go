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


// Run Once, Start Immediately
const testConfig string = `{
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

//Run Every 5 seconds, start Immediately
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

// Run Once, Start Delayed at 2017-06-14T15:52:00Z02:00
const testConfig3 string = `{
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "false",
				"startImmediate": "false",
				"startDate" : "2017-06-14T16:56:00Z02:00"
      }
    }
  ]
}`


//Run Every 5 seconds, start Delayed
const testConfig4 string = `{
"name": "timer2",
"settings": {
},
"handlers": [
	{
		"actionId": "local://testFlow2",
		"settings": {
		"repeating": "true",
			"startImmediate": "false",
			"startDate" : "2017-06-14T17:18:00Z02:00",
			"seconds": "5",
			"minutes": "0",
			"hours": "0"
		}
	}
]
}`


// Multiple timer configurations
const testConfig5 string = `{
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
        "startDate" : "2017-06-14T15:52:00Z02:00",
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
	json.Unmarshal([]byte(testConfig4), &config)
	f := &Timer2Factory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
	log.Infof("Press CTRL-C to quit")
  for {}
}

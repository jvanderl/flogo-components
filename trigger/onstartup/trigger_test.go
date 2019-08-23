package onstartup

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "onstartup",
	"ref": "github.com/jvanderl/flogo-components/trigger/onstartup",
	"handlers": [
	  {
		"settings":{
		},
		"action":{
			"id":"dummy"
		}
	  }
	]
  }
  `

func TestInitOk(t *testing.T) {
	f := &Factory{}
	tgr, err := f.New(nil)
	assert.Nil(t, err)
	assert.NotNil(t, tgr)
}

func TestSingle(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	trg.Start()
	trg.Stop()
}

/*
func TestTimerTrigger_Initialize(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	err = trg.Start()
	assert.Nil(t, err)
	err = trg.Stop()
	assert.Nil(t, err)

}
*/
/*
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
  "name": "onstartup",
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
	//tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
	log.Infof("Press CTRL-C to quit")
	for {
	}
}
*/

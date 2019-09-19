package stomp

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "stomp",
	"ref": "github.com/jvanderl/flogo-components/trigger/stomp",
	"settings": {
        "address": "3dexp.18xfd05.ds:61613"
    },
	"handlers": [
	  {
		"settings": {
		  "source": "/queue/test"
		},
		"action" : {
		  "id": "dummy"
		}
	  }
	]
  }`

const testConfig2 string = `{
	"id": "stomp",
	"ref": "github.com/jvanderl/flogo-components/trigger/stomp",
	"settings": {
        "address": "3dexp.18xfd05.ds:61613"
    },
	"handlers": [
	  {
		"settings": {
		  "source": "/queue/test"
		},
		"action" : {
		  "id": "dummy"
		}
	  },
	  {
		"settings": {
		  "source": "/queue/test2"
		},
		"action" : {
		  "id": "dummy2"
		}
	  }
	]
  }`

func TestTrigger_Init(t *testing.T) {

	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
	})}

	trg, err := test.InitTrigger(f, config, actions)

	assert.Nil(t, err)
	assert.NotNil(t, trg)

}

func TestTrigger_Single(t *testing.T) {
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

func TestTrigger_Multi(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig2), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	}), "dummy2": test.NewDummyAction(func() {
		//stiil do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	trg.Start()
	trg.Stop()

}

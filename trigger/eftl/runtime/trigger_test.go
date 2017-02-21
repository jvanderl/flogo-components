package eftl

import (
	"context"
	"testing"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

const testConfig string = `{
  "name": "eftl",
  "settings": {
    "server": "192.168.178.41:9191",
    "channel": "/channel",
    "destination": "sample",
    "user": "user",
    "password": "password"
  },
  "endpoints": [
    {
      "actionType": "flow",
      "actionURI": "local://testFlow",
      "settings": {
        "destination": "flogo"
      }
    },
    {
      "actionType": "flow",
      "actionURI": "local://testFlow2",
      "settings": {
        "destination": "sample"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("eftl")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("eftl")

	runner := &TestRunner{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(config, runner)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("eftl")

	tgr.Start()
	defer tgr.Stop()

}

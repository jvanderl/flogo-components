package mqtt3

import (
	"context"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("mqtt3")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

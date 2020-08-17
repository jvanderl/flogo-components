package mqtt

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestPlain(t *testing.T) {

	settings := &Settings{Broker: "tcp://192.168.8.127:1883", Id: "flogo-tester"}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("topic", "bla")
	tc.SetInput("qos", 0)
	tc.SetInput("message", "blabla")
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

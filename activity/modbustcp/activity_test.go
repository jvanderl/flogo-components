package modbustcp

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

func TestReadCoils(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "ReadCoils")
	tc.SetInput("address", 300)
	tc.SetInput("numElements", 2)
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestReadDiscreteInputs(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "ReadDiscreteInputs")
	tc.SetInput("address", 200)
	tc.SetInput("numElements", 1)
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestReadInputRegisters(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "ReadInputRegisters")
	tc.SetInput("address", 100)
	tc.SetInput("numElements", 1)
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestReadHoldingRegisters(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "ReadHoldingRegisters")
	tc.SetInput("address", 400)
	tc.SetInput("numElements", 2)
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestWriteSingleCoil(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "WriteSingleCoil")
	tc.SetInput("address", 300)
	tc.SetInput("data", 0xFF00)
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestWriteMultipleCoils(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "WriteMultipleCoils")
	tc.SetInput("address", 300)
	tc.SetInput("numElements", 2)
	tc.SetInput("data", []byte{3})
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestWriteSingleRegister(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "WriteSingleRegister")
	tc.SetInput("address", 400)
	tc.SetInput("data", 43)
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

func TestWriteMultipleRegisters(t *testing.T) {

	settings := &Settings{Server: "192.168.8.141:502", Timeout: 10, SlaveID: 0xFF}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("operation", "WriteMultipleRegisters")
	tc.SetInput("address", 400)
	tc.SetInput("numElements", 2)
	tc.SetInput("data", []byte{0, 44, 1, 45})
	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, "OK", output.Result)
}

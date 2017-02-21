package eftl

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("eftl")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &MyActivity{metadata: md}

	tc := test.NewTestActivityContext(md)
	//setup attrs

	fmt.Println("Publishing a flogo test message to destination 'sample' on channel '/channel' on eFTL Server '192.168.178.41:9191'")

	tc.SetInput("server", "192.168.178.41:9191")
	tc.SetInput("channel", "/channel")
	tc.SetInput("destination", "flogo")
	tc.SetInput("user", "user")
	tc.SetInput("password", "password")
	tc.SetInput("message", "{\"deviceID\":\"5CCF7F942BCB\",\"distance\":9,\"distState\":\"Safe\"}")

	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	tc.SetInput("destination", "sample")

	act.Eval(tc)

	//check result attr
	result = tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

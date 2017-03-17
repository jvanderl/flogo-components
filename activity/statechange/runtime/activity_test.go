package statechange

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("statechange")

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

	fmt.Println("Setting up inital state for 'UID123123'")
	tc.SetInput("datasource", "UID123123")
	tc.SetInput("input1", "5")
	fmt.Println("First run should have changed is 'true', result is 'initial value', flags is '1'")

	act.Eval(tc)

	changed := tc.GetOutput("changed")
	result := tc.GetOutput("result")
	flags := tc.GetOutput("flags")

	fmt.Println("changed: ", changed)
	fmt.Println("result: ", result)
	fmt.Println("flags: ", flags)

	fmt.Println("Second run we offer 5 again. Should not detect change")

	act.Eval(tc)

	changed = tc.GetOutput("changed")
	result = tc.GetOutput("result")
	flags = tc.GetOutput("flags")

	fmt.Println("changed: ", changed)
	fmt.Println("result: ", result)
	fmt.Println("flags: ", flags)

	fmt.Println("Third run we offer a value of '20'. Should detect change again")

	tc.SetInput("input1", "20")

	act.Eval(tc)

	changed = tc.GetOutput("changed")
	result = tc.GetOutput("result")
	flags = tc.GetOutput("flags")

	fmt.Println("changed: ", changed)
	fmt.Println("result: ", result)
	fmt.Println("flags: ", flags)

	//check result attr

	if changed == nil {
		t.Fail()
	}
}

package getjson

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"fmt"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("getjson")

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

	//setup attrs
	tc.SetInput("input", "{\"distance\":150, \"status\":\"optimal\"}")
	tc.SetInput("name1", "distance")
	tc.SetInput("name2", "status")

	act.Eval(tc)

	result := tc.GetOutput("result")
	value1 := tc.GetOutput("value1")
	value2 := tc.GetOutput("value2")
	
	fmt.Println("result: ", result)
	fmt.Println("value1: ", value1)
	fmt.Println("value2: ", value2)

	if result == nil {
		t.Fail()
	}
}

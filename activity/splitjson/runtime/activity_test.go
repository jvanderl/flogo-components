package splitjson

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("splitjson")

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

	act.Eval(tc)

	result := tc.GetOutput("result")
	name1 := tc.GetOutput("name1")
	value1 := tc.GetOutput("value1")
	name2 := tc.GetOutput("name2")
	value2 := tc.GetOutput("value2")

	fmt.Println("result: ", result)
	fmt.Println("name1: ", name1)
	fmt.Println("value1: ", value1)
	fmt.Println("name2: ", name2)
	fmt.Println("value2: ", value2)

	if result == nil {
		t.Fail()
	}
}

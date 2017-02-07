package mqtt

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
	"fmt"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("mqtt")

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

	fmt.Println("Publishing a flogo test message to topic 'flogo' on broker 'localhost:1883'")

	tc.SetInput("broker", "tcp://127.0.0.1:1883")
	tc.SetInput("id", "flogo_tester")
	tc.SetInput("topic", "flogo")
	tc.SetInput("qos", 0)
	tc.SetInput("message", "This is a test message from flogo")

	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

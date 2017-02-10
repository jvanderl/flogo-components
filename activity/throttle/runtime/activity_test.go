package throttle

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
	"time"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("throttle")

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

	fmt.Println("Setting up Throttle for 'UID123123' with interval 5 seconds")
	tc.SetInput("datasource", "UID123123")
	tc.SetInput("interval", 5)
	tc.SetInput("intervaltype", "seconds")

	fmt.Println("First run should pass")

	act.Eval(tc)

	pass := tc.GetOutput("pass")
	reason := tc.GetOutput("reason")
	lasttimepassed := tc.GetOutput("lasttimepassed")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)
	fmt.Println("lasttimepassed: ", lasttimepassed)

	fmt.Println("Wait 2 seconds")
	time.Sleep(2000 * time.Millisecond)

	fmt.Println("Second run should not pass")

	act.Eval(tc)

	pass = tc.GetOutput("pass")
	reason = tc.GetOutput("reason")
	lasttimepassed = tc.GetOutput("lasttimepassed")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)
	fmt.Println("lasttimepassed: ", lasttimepassed)

	fmt.Println("Wait 3 more seconds")
	time.Sleep(3000 * time.Millisecond)

	fmt.Println("Third run should pass again")

	act.Eval(tc)

	pass = tc.GetOutput("pass")
	reason = tc.GetOutput("reason")
	lasttimepassed = tc.GetOutput("lasttimepassed")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)
	fmt.Println("lasttimepassed: ", lasttimepassed)
	//check result attr

	if pass == nil {
		t.Fail()
	}
}

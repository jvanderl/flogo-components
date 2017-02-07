package filter

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"fmt"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("filter")

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

	fmt.Println("Setting input to 150, minvalue: 100, maxvalue: 200. Should pass")

	tc.SetInput("input", "150")
	tc.SetInput("datatype", "int")
	tc.SetInput("minvalue", "100")
	tc.SetInput("maxvalue", "200")
	tc.SetInput("inverse", false)

	act.Eval(tc)

	//check result attr
	pass := tc.GetOutput("pass")
	reason := tc.GetOutput("reason")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)

////////////////////////

	fmt.Println("Setting input to 150, maxvalue: 100. Should not pass, value too high")

	tc.SetInput("input", "150")
	tc.SetInput("datatype", "int")
	tc.SetInput("maxvalue", "100")
	tc.SetInput("inverse", false)

	act.Eval(tc)

	//check result attr
	pass = tc.GetOutput("pass")
	reason = tc.GetOutput("reason")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)

////////////////////////

	fmt.Println("Setting input to 150, minvalue: 200. Should not pass, value too low")

	tc.SetInput("input", "150")
	tc.SetInput("datatype", "int")
	tc.SetInput("minvalue", "200")
	tc.SetInput("inverse", false)

	act.Eval(tc)

	//check result attr
	pass = tc.GetOutput("pass")
	reason = tc.GetOutput("reason")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)

////////////////////////
	fmt.Println("Setting input to 150, minvalue: 100, maxvalue: 200, inverse: true. Should not pass, filter out mid section")

	tc.SetInput("input", "150")
	tc.SetInput("datatype", "int")
	tc.SetInput("minvalue", "100")
	tc.SetInput("maxvalue", "200")
	tc.SetInput("inverse", true)

	act.Eval(tc)

	//check result attr
	pass = tc.GetOutput("pass")
	reason = tc.GetOutput("reason")

	fmt.Println("pass: ", pass)
	fmt.Println("reason: ", reason)

////////////////////////

	if reason == nil {
		t.Fail()
	}
}

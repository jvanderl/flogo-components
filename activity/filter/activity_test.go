package filter

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"io/ioutil"
	"testing"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}
		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}
	return activityMetadata
}

func TestCreate(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Created")
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

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

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

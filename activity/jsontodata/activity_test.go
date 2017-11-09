package jsontodata

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
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

	fmt.Println("splitting json string into separate name value pairs: '{\"distance\":150, \"status\":\"optimal\"}'")

	//setup attrs
	tc.SetInput("input", "{\"distance\":150, \"status\":\"optimal\"}")

	act.Eval(tc)

	result := tc.GetOutput("result")
	//var data map[string]string
	data := tc.GetOutput("data")

	fmt.Println("distance: ", data)
	//fmt.Println("status: ", data.status)

	if result == nil {
		t.Fail()
	}
}

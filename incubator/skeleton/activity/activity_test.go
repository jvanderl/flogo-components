package [package]

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

	fmt.Println("[functionality]")

	//set inputs
	tc.SetInput("[input1_name]", "[input1_value]")
	tc.SetInput("[input2_name]", "[input2_value]")

	//evaluate activity
	act.Eval(tc)

	//get outputs
	[output1_name] := tc.GetOutput("[output1_name]")
	[output2_name] := tc.GetOutput("[output2_name]")

	//print outputs
	fmt.Println("[output1_name]: ", output1_name)
	fmt.Println("[output2_name]: ", output2_name)

	if result == nil {
		t.Fail()
	}
}

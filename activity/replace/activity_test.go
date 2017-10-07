package replace

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

	fmt.Println("Stripping colons from mac addres 5C:CF:7F:94:2B:CB")

	//setup attrs
	tc.SetInput("find", ":")
	tc.SetInput("replace", "")
	tc.SetInput("input1", "5C:CF:7F:94:2B:CB")

	act.Eval(tc)

	result := tc.GetOutput("result")
	output1 := tc.GetOutput("output1")
	output2 := tc.GetOutput("output2")

	fmt.Println("result: ", result)
	fmt.Println("output1: ", output1)
	fmt.Println("output2: ", output2)

	if result == nil {
		t.Fail()
	}
}

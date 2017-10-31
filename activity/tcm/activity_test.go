package tcm

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
	//setup attrs

	fmt.Println("Publishing a flogo test message to TIBCO Cloud Messaging")

	tc.SetInput("url", "<Your TCM URL Here>")
	tc.SetInput("authkey", "Your TCM Auth Key Here")
	tc.SetInput("clientid", "flogo_test")
	tc.SetInput("messagename", "demo_tcm")
	tc.SetInput("messagevalue", "Hello TCM from FLOGO")
	tc.SetInput("certificate", "")

	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

package tcmpub

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"io/ioutil"
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

func TestActivityRegistration(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}


func TestEval(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())
	//setup attrs
	tc.SetInput("url", "<Your TCM URL Here>")
	tc.SetInput("authkey", "Your TCM Auth Key Here")
	tc.SetInput("clientid", "flogo_test")
	tc.SetInput("messagename", "demo_tcm")
	tc.SetInput("messagevalue", "Hello TCM from FLOGO")

	done,err := act.Eval(tc)
	assert.Nil(t, err)
	result := tc.GetOutput("result")
	assert.Equal(t, result, "ERR_CONNECT_HOST")
}

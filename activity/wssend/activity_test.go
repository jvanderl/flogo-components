package wssend

import (
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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

func TestOK(t *testing.T) {
	fmt.Println("Setting up to succeed.")

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("server", "localhost:8080")
	tc.SetInput("channel", "/echo")
	tc.SetInput("message", "haha!")

	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")
	fmt.Printf("result: %v\n", result)
}

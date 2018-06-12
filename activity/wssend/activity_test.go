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
	tc.SetInput("server", "localhost:9191")
	tc.SetInput("channel", "/test")
	tc.SetInput("message", "haha!")
	tc.SetInput("waitforresponse", "true")
	tc.SetInput("timeout", "5")

	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")
	fmt.Printf("result: %v\n", result)
	response := tc.GetOutput("response")
	fmt.Printf("response: %v\n", response)
}

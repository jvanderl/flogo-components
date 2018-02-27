package matchresponse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

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
	searchdataJSON := `[
		{"find": "thank you", "resp": "Don't mention it."},
		{"find": "bye", "resp": "See you later!"}
		]`
	var searchdata interface{}
	err := json.Unmarshal([]byte(searchdataJSON), &searchdata)
	if err != nil {
		panic("Unable to parse search data: " + err.Error())
	}

	fmt.Println("finding a match...")

	tc.SetInput("input", "Thank you for your help")
	tc.SetInput("searchdata", searchdata)

	act.Eval(tc)

	outmatch := tc.GetOutput("match")
	response := tc.GetOutput("response")

	fmt.Println("match: ", outmatch)
	fmt.Println("response: ", response)

	fmt.Println("finding another match...")
	tc.SetInput("input", "Bye for now")

	act.Eval(tc)

	outmatch = tc.GetOutput("match")
	response = tc.GetOutput("response")

	fmt.Println("match: ", outmatch)
	fmt.Println("response: ", response)

}

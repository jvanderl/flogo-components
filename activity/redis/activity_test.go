package redis

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

	fmt.Println("Pinging redis server on 'localhost:6379', expecting result 'PONG'")

	tc.SetInput("server", "127.0.0.1:6379")
	tc.SetInput("password", "")
	tc.SetInput("operation", "PING")
	tc.SetInput("database", 0)

	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")
	fmt.Println("result: ", result)



	//setup attrs

	fmt.Println("Setting key 'flogo_test' to 'Test123', expecting result 'OK'")

	tc.SetInput("operation", "SET")
	tc.SetInput("key", "flogo_test")
	tc.SetInput("value", "Test123")

	act.Eval(tc)

	//check result attr
	result = tc.GetOutput("result")
	fmt.Println("result: ", result)

	//setup attrs
	fmt.Println("Getting value for key 'flogo_test', expecting result 'Test123'")
	tc.SetInput("operation", "GET")

	act.Eval(tc)

	//check result attr
	result = tc.GetOutput("result")
	fmt.Println("result: ", result)

	//setup attrs
	fmt.Println("Deleting key 'flogo_test', expecting result 'OK'")
	tc.SetInput("operation", "DEL")

	act.Eval(tc)

	//check result attr
	result = tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}






}

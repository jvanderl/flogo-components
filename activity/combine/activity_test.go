package combine

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
			if err != nil{
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

	md := getActivityMetadata()
	act := NewActivity(getActivityMetadata())

	tc := test.NewTestActivityContext(getActivityMetadata())

	///////////////////

	fmt.Println("Combining 3 parts into path using delimiter, no prefix, no suffix")
	//setup attrs
	tc.SetInput("delimiter", "/")
	tc.SetInput("part1", "device")
	tc.SetInput("part2", "Dev12345")
	tc.SetInput("part3", "distance")

	act.Eval(tc)

	result := tc.GetOutput("result")

	fmt.Println("result: ", result)

	///////////////////

	fmt.Println("Combining 3 parts into path using delimiter, no prefix, no suffix, intented to start with delimiter")
	//setup attrs
	tc.SetInput("delimiter", "/")
	tc.SetInput("part1", "/device")
	tc.SetInput("part2", "Dev12345")
	tc.SetInput("part3", "distance")

	act.Eval(tc)

	result = tc.GetOutput("result")

	fmt.Println("result: ", result)

	///////////////////

	fmt.Println("Combining 3 parts into path using delimiter, with prefix, no suffix")
	//setup attrs
	tc.SetInput("delimiter", "/")
	tc.SetInput("prefix", "flogo")
	tc.SetInput("part1", "device")
	tc.SetInput("part2", "Dev12345")
	tc.SetInput("part3", "distance")

	act.Eval(tc)

	result = tc.GetOutput("result")

	fmt.Println("result: ", result)

	///////////////////

	fmt.Println("Combining 8 parts into path using delimiter, prefix and suffix")
	//setup attrs
	tc.SetInput("delimiter", "/")
	tc.SetInput("prefix", "/flogo")
	tc.SetInput("suffix", "parameter")
	tc.SetInput("part1", "device")
	tc.SetInput("part2", "Dev12345")
	tc.SetInput("part3", "distance")
	tc.SetInput("part4", "part4")
	tc.SetInput("part5", "part5")
	tc.SetInput("part6", "part6")
	tc.SetInput("part7", "part7")
	tc.SetInput("part8", "part8")

	act.Eval(tc)

	result = tc.GetOutput("result")

	fmt.Println("result: ", result)

	///////////////////

	fmt.Println("Simple concatenation")
	//setup attrs
	tc.SetInput("delimiter", nil)
	tc.SetInput("prefix", nil)
	tc.SetInput("suffix", nil)
	tc.SetInput("part1", "This")
	tc.SetInput("part2", "Will")
	tc.SetInput("part3", "Become")
	tc.SetInput("part4", "One")
	tc.SetInput("part5", "Word")
	tc.SetInput("part6", nil)
	tc.SetInput("part7", nil)
	tc.SetInput("part8", nil)

	act.Eval(tc)

	result = tc.GetOutput("result")

	fmt.Println("result: ", result)

	///////////////////

	fmt.Println("Create Simple JSON")
	//setup attrs
	tc.SetInput("delimiter", nil)
	tc.SetInput("prefix", "{")
	tc.SetInput("suffix", "}")
	tc.SetInput("part1", "\"distance\":")
	tc.SetInput("part2", "25")
	tc.SetInput("part3", ",")
	tc.SetInput("part4", "\"distState\":\"")
	tc.SetInput("part5", "Safe")
	tc.SetInput("part6", "\"")
	tc.SetInput("part7", nil)
	tc.SetInput("part8", nil)

	act.Eval(tc)

	result = tc.GetOutput("result")

	fmt.Println("result: ", result)

	///////////////////

	if result == nil {
		t.Fail()
	}
}

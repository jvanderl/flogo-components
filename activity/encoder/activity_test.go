package encoder

import (
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

func TestBase32Encode(t *testing.T) {

	fmt.Print("Testing Base32 Encoder\ninput: This is a test\n")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("input", "This is a test")
	tc.SetInput("action", "ENCODE")
	tc.SetInput("encoder", "BASE32")

	//eval
	act.Eval(tc)
	result := tc.GetOutput("result")
	status := tc.GetOutput("status")

	fmt.Printf("result: %v\n", result)
	fmt.Printf("status: %v\n", status)
}

func TestBase32Decode(t *testing.T) {

	fmt.Print("Testing Base32 Decoder\ninput: KRUGS4ZANFZSAYJAORSXG5A=\n")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("input", "KRUGS4ZANFZSAYJAORSXG5A=")
	tc.SetInput("action", "DECODE")
	tc.SetInput("encoder", "BASE32")

	//eval
	act.Eval(tc)
	result := tc.GetOutput("result")
	status := tc.GetOutput("status")

	fmt.Printf("result: %v\n", result)
	fmt.Printf("status: %v\n", status)
}

func TestBase64Encode(t *testing.T) {

	fmt.Print("Testing Base64 Encoder\ninput: This is a test\n")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("input", "This is a test")
	tc.SetInput("action", "ENCODE")
	tc.SetInput("encoder", "BASE64")

	//eval
	act.Eval(tc)
	result := tc.GetOutput("result")
	status := tc.GetOutput("status")

	fmt.Printf("result: %v\n", result)
	fmt.Printf("status: %v\n", status)
}

func TestBase64Decode(t *testing.T) {

	fmt.Print("Testing Base64 Decoder\ninput: VGhpcyBpcyBhIHRlc3Q=\n")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("input", "VGhpcyBpcyBhIHRlc3Q=")
	tc.SetInput("action", "DECODE")
	tc.SetInput("encoder", "BASE64")

	//eval
	act.Eval(tc)
	result := tc.GetOutput("result")
	status := tc.GetOutput("status")

	fmt.Printf("result: %v\n", result)
	fmt.Printf("status: %v\n", status)
}

func TestHexEncode(t *testing.T) {

	fmt.Print("Testing Hex Encoder\ninput: This is a test\n")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("input", "This is a test")
	tc.SetInput("action", "ENCODE")
	tc.SetInput("encoder", "HEX")

	//eval
	act.Eval(tc)
	result := tc.GetOutput("result")
	status := tc.GetOutput("status")

	fmt.Printf("result: %v\n", result)
	fmt.Printf("status: %v\n", status)
}

func TestHexDecode(t *testing.T) {

	fmt.Print("Testing Base64 Decoder\ninput: 5468697320697320612074657374\n")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("input", "5468697320697320612074657374")
	tc.SetInput("action", "DECODE")
	tc.SetInput("encoder", "HEX")

	//eval
	act.Eval(tc)
	result := tc.GetOutput("result")
	status := tc.GetOutput("status")

	fmt.Printf("result: %v\n", result)
	fmt.Printf("status: %v\n", status)
}

package checkiban

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

	fmt.Println("Checkinb IBAN NL40ABNA0517552264")

	tc.SetInput("iban", "NL40ABNA0517552264")

	act.Eval(tc)

	//check result attr

	result := tc.GetOutput("result")
	code := tc.GetOutput("code")
	printcode := tc.GetOutput("printcode")
	countrycode := tc.GetOutput("countrycode")
	checkdigits := tc.GetOutput("checkdigits")
	bban := tc.GetOutput("bban")
	ibanobj := tc.GetOutput("ibanobj")

	fmt.Println("result: ", result)
	fmt.Println("code: ", code)
	fmt.Println("printcode: ", printcode)
	fmt.Println("countrycode: ", countrycode)
	fmt.Println("checkdigits: ", checkdigits)
	fmt.Println("bban: ", bban)
	fmt.Println("ibanobj: ", ibanobj)

	if result == nil {
		t.Fail()
	}

}

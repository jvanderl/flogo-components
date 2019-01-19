package dslogin

import (
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

const reqPostStr string = `{
  "name": "my pet"
}
`

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

func TestDSLogin(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("3DPassportURL", "https://3dexp.18xfd05.ds/3DPassport")
	tc.SetInput("3DServiceURL", "https://3dexp.18xfd05.ds/3DSpace")
	tc.SetInput("userName", "demoleader")
	tc.SetInput("serviceName", "flogo")
	tc.SetInput("serviceSecret", "9500ca71-1012-4f81-a739-263d3d60a400")
	tc.SetInput("skipSsl", "true")

	//eval
	act.Eval(tc)

	val := tc.GetOutput("serviceAccessToken")
	fmt.Printf("serviceAccessToken: %v\n", val)
	val = tc.GetOutput("serviceRedirectURL")
	fmt.Printf("serviceRedirectURL: %v\n", val)
}

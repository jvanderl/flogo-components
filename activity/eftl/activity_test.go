package eftl

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

	fmt.Println("Publishing a flogo test message to destination 'flogo' on channel '/channel' on eFTL Server")

	tc.SetInput("server", "localhost:9191")
	tc.SetInput("clientid", "5CCF7F942BCB")
	tc.SetInput("channel", "/channel")
  tc.SetInput("destination", "flogo")
	tc.SetInput("user", "user")
	tc.SetInput("password", "password")
	tc.SetInput("message", "{\"deviceID\":\"5CCF7F942BCB\",\"distance\":9,\"distState\":\"Safe\"}")
	tc.SetInput("secure", false)
	tc.SetInput("certificate", "")

	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

/*
func TestEvalSecure(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())
	//setup attrs

	fmt.Println("Publishing a flogo test message to destination 'sample' on channel '/channel' on eFTL Server '10.10.1.50:9291'")

	tc.SetInput("server", "10.10.1.50:9291")
	tc.SetInput("channel", "/channel")
	tc.SetInput("destination", "flogo")
//	tc.SetInput("destination", "device")
	tc.SetInput("user", "user")
	tc.SetInput("password", "password")
	tc.SetInput("message", "{\"deviceID\":\"5C:CF:7F:94:2B:CB\",\"distance\":9,\"distState\":\"Safe\"}")
//	tc.SetInput("message", "{\"deviceID\":\"5C:CF:7F:94:2B:CB\",\"distFactor\":1.0}")
	tc.SetInput("secure", true)
	tc.SetInput("certificate", "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUV6ekNDQTdlZ0F3SUJBZ0lKQUlKU2RCd2QzZjVUTUEwR0NTcUdTSWIzRFFFQkJRVUFNSUdhTVFzd0NRWUQKVlFRR0V3Sk9UREVMTUFrR0ExVUVDQk1DV2tneEVqQVFCZ05WQkFjVENWSnZkSFJsY21SaGJURVhNQlVHQTFVRQpDaE1PVkVsQ1EwOGdVMjltZEhkaGNtVXhEVEFMQmdOVkJBc1RCRk5EVGt3eEh6QWRCZ05WQkFNVEZrcGhibk10ClRXRmpRbTl2YXkxUWNtOHViRzlqWVd3eElUQWZCZ2txaGtpRzl3MEJDUUVXRW1wMllXNWtaWEpzUUhScFltTnYKTG1OdmJUQWVGdzB4TnpBMU1Ea3lNREUzTWpkYUZ3MHhPREExTURreU1ERTNNamRhTUlHYU1Rc3dDUVlEVlFRRwpFd0pPVERFTE1Ba0dBMVVFQ0JNQ1drZ3hFakFRQmdOVkJBY1RDVkp2ZEhSbGNtUmhiVEVYTUJVR0ExVUVDaE1PClZFbENRMDhnVTI5bWRIZGhjbVV4RFRBTEJnTlZCQXNUQkZORFRrd3hIekFkQmdOVkJBTVRGa3BoYm5NdFRXRmoKUW05dmF5MVFjbTh1Ykc5allXd3hJVEFmQmdrcWhraUc5dzBCQ1FFV0VtcDJZVzVrWlhKc1FIUnBZbU52TG1OdgpiVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFOM1lKa1lWY3ppNlJ3T2piZmt3CjNmOVNxT3pYKys1MGRGTjcyTFU4ZHpiTGRoM29BVFkrY3pHZ1RlbkF4akJsNm9zM09aS1ZYaE85OHlDSzd5NzEKeWpEWE5zYzJHZ2ZnbGwvUVJLb3VXcnRCazAvV3BUVUo1MnZZZzdjeXgxeUFyWmZwWE5EVS9TSnhJYlpxODRSSgpXTDhlUnJsYlE1ZEFMZW1NSFZDM1BYWkZuUUFCTU1ON3JVaHk5UFJSVVNYUFp2TTZWT0ZGOHc3MUJlSXZXZW1XCmV1QTJnRSsvdGdRTE9JckJBZlJnUFkxMUp0ZjBqY0NoMDZ1VGJrYWpHd0hkc3hNQmwzbkhyRktDdFhrSFJ4NWEKd2pCbnYwMlhOT2lYa1hpcU5pUmNpSUNEemlwR09kMCtJc01TU29CWnZEMkNxMnExRUQ5Q0NlQVBianN3bXM5RQpibWNDQXdFQUFhT0NBUlF3Z2dFUU1CMEdBMVVkRGdRV0JCVGlqaVFucFF5NHN2RElFRUNMM0JTMlk5cnZ2RENCCnp3WURWUjBqQklISE1JSEVnQlRpamlRbnBReTRzdkRJRUVDTDNCUzJZOXJ2dktHQm9LU0JuVENCbWpFTE1Ba0cKQTFVRUJoTUNUa3d4Q3pBSkJnTlZCQWdUQWxwSU1SSXdFQVlEVlFRSEV3bFNiM1IwWlhKa1lXMHhGekFWQmdOVgpCQW9URGxSSlFrTlBJRk52Wm5SM1lYSmxNUTB3Q3dZRFZRUUxFd1JUUTA1TU1SOHdIUVlEVlFRREV4WktZVzV6CkxVMWhZMEp2YjJzdFVISnZMbXh2WTJGc01TRXdId1lKS29aSWh2Y05BUWtCRmhKcWRtRnVaR1Z5YkVCMGFXSmoKYnk1amIyMkNDUUNDVW5RY0hkMytVekFNQmdOVkhSTUVCVEFEQVFIL01BOEdBMVVkRVFRSU1BYUhCQW9LQVRJdwpEUVlKS29aSWh2Y05BUUVGQlFBRGdnRUJBTW5GZUN2djhwUnN6RUZyS295N1VRdEZCTHJlb09qMFptdFVlT1ZNCllxcGxhWGdnZE9rbEtHY0ZQM29jS3N3S3RWejNCZ3BxdjEwNm44Z3RDT1RuY3JZcG1aNDFQN3VBaXF4dTVnRGsKaWF4aHh4NjdxU1I5eXJ6R29aUjhNaVZEamF3ejNtMHZDbzNDQjljcHV0WVpYU055NjZIMlozb2pXQkJZTnhrVApkc2dTV1pIbzhZRVkxelErWXJ6M3lmNHJrQXJXREV1dlUxRk9ZaEc2M0oyRWN6RHptRXp1RHozaCtrZGtQTEhrCnBFNFhlY29tUXBhbEpGd3VSYndYUnUzNnpiWGVxUjRYOGNRRUs0ZmlDTWxjczROczZGbzJhUDkwTmFzelowSFoKUkIraVlON0p4Y3hsRTN5Y0Z0T3BWMDFRMU9zUzM1Q1V6cWVERzRRdVhWQnAxMGc9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K")

	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

} */

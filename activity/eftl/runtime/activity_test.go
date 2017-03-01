package eftl

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("eftl")

	if act == nil {
		t.Error("Activity Not Registered")
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

	md := activity.NewMetadata(jsonMetadata)
	act := &MyActivity{metadata: md}

	tc := test.NewTestActivityContext(md)
	//setup attrs

	fmt.Println("Publishing a flogo test message to destination 'sample' on channel '/channel' on eFTL Server 'tibco4.demo.local:9191'")

	tc.SetInput("server", "tibco4.demo.local:9191")
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

func TestEvalSecure(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &MyActivity{metadata: md}

	tc := test.NewTestActivityContext(md)
	//setup attrs

	fmt.Println("Publishing a flogo test message to destination 'sample' on channel '/channel' on eFTL Server 'tibco4.demo.local:9291'")

	tc.SetInput("server", "tibco4.demo.local:9291")
	tc.SetInput("channel", "/channel")
	tc.SetInput("destination", "flogo")
	tc.SetInput("user", "user")
	tc.SetInput("password", "password")
	tc.SetInput("message", "{\"deviceID\":\"5CCF7F942BCB\",\"distance\":9,\"distState\":\"Safe\"}")
	tc.SetInput("secure", true)
	tc.SetInput("certificate", "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUQrRENDQXVDZ0F3SUJBZ0lKQUtIRGtEUk5rdUdlTUEwR0NTcUdTSWIzRFFFQkN3VUFNSUdRTVFzd0NRWUQKVlFRR0V3Sk9UREVMTUFrR0ExVUVDQXdDV2tneEVqQVFCZ05WQkFjTUNWSnZkSFJsY21SaGJURVhNQlVHQTFVRQpDZ3dPVkVsQ1EwOGdVMjltZEhkaGNtVXhEVEFMQmdOVkJBc01CRk5EVGt3eEZUQVRCZ05WQkFNTURDb3VaR1Z0CmJ5NXNiMk5oYkRFaE1COEdDU3FHU0liM0RRRUpBUllTYW5aaGJtUmxjbXhBZEdsaVkyOHVZMjl0TUI0WERURTMKTURJeU56RXdORFV6TUZvWERURTNNRE15T1RFd05EVXpNRm93Z1pBeEN6QUpCZ05WQkFZVEFrNU1NUXN3Q1FZRApWUVFJREFKYVNERVNNQkFHQTFVRUJ3d0pVbTkwZEdWeVpHRnRNUmN3RlFZRFZRUUtEQTVVU1VKRFR5QlRiMlowCmQyRnlaVEVOTUFzR0ExVUVDd3dFVTBOT1RERVZNQk1HQTFVRUF3d01LaTVrWlcxdkxteHZZMkZzTVNFd0h3WUoKS29aSWh2Y05BUWtCRmhKcWRtRnVaR1Z5YkVCMGFXSmpieTVqYjIwd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQQpBNElCRHdBd2dnRUtBb0lCQVFDa3JyZnEyZVpyeDZVdEwxMVFjL3c0eVczQnVGZlVOQzFSRVAwZjBoYkh5aTlLCldJYzI5cENiNUs0TTVMSEpqK29lOTR4Y2JOVFVzMnlYZ3VwOExIRkFUUzN4WnhyUUJZR2lTaEZYdTRKY3hMaEoKZG1TbDBLUWQwT0x0SWtzSzQ5SDNScXhaOGpHMkorQURNS1loT0ZrM3dLS0FLRDduMXNJazVOK1JrNDVMYy9sVQo3UnZ3S01sTjE3eFdGdTdCSnVnYWtkSldTRUk1R1drclhzbHpHeXMxTDN6b1haM1dMdWI3VHZ3SFd6aG1UT0FtCmxmcDV5bUpFU0hXeXh1UDRDZ2JiRDJld25sTzIwMjc0cHdyYnNCL1hFQmtRamZmY3dtVENNNW5yVnlSKzhFVm4KaU0yZEFKdFdyN1V5eE5yK2o2NmhlZXhLWHFTUXhjc3Y0QTdxRG5MNUFnTUJBQUdqVXpCUk1CMEdBMVVkRGdRVwpCQlFYUWNpcy9VVk5LMURXVGRSTnVGbkNScHkwQ2pBZkJnTlZIU01FR0RBV2dCUVhRY2lzL1VWTksxRFdUZFJOCnVGbkNScHkwQ2pBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFCVXhsczkKemI3RkpSeVV3KzYrZG1sQWhEQ3I3K1YzN1ZGV3VoRmN1NlRwdXRIYWNTNFlvWmltR1d4Z1E2TU5lZEp2ajJoMQp4bWZudmVEc3Z6cHYweE9oU3kyRUVTOUd4emZhQUdCMXdzS0lzajdQVDh0TzlxMlVLOFkveE85OGdqWnFmZzNPCkpUVTA1UFZCSjIvS3VZNWlqUFI0d1V3WEVydXJYYTk0K1p2R0hBUlNSS0FnaHA2S0diTDRCZ1g3ZHhaTXRjK0kKQWhTVUNtSUtWbWh1VG94TFlETWJTS2hiR2VhaEl4SzN0RzNwMVk4NGxEbXRsbzN1Z3RuZXlhVWdOVjVGMWh1bgpVQVRWb2NDUy96b3Uzb3prSTlnOGhJUVV3bit4ZHlYOGZiYlFaVmRxUThSaXc2b0l6ZUdCejV1K1ptdnR3YTYzClh5NHFuYkhyTnNQVnJ2Q0wKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==")
	
	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

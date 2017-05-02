package eftl

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
	//md := getActivityMetadata()
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())
	//setup attrs

	fmt.Println("Publishing a flogo test message to destination 'sample' on channel '/channel' on eFTL Server '10.10.1.50:19191'")

	tc.SetInput("server", "localhost:19191")
	tc.SetInput("channel", "/channel")
	tc.SetInput("destination", "flogo")
	tc.SetInput("user", "user")
	tc.SetInput("password", "password")
//	tc.SetInput("message", "{\"deviceID\":\"5CCF7F942BCB\",\"distance\":9,\"distState\":\"Safe\"}")
	tc.SetInput("message", "Simple Message")
	tc.SetInput("secure", false)
	tc.SetInput("certificate", "DummyCert")

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

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())
	//setup attrs

	fmt.Println("Publishing a flogo test message to destination 'sample' on channel '/channel' on eFTL Server '10.10.1.50:9291'")

	tc.SetInput("server", "10.10.1.50:9291")
	tc.SetInput("channel", "/channel")
	tc.SetInput("destination", "flogo")
	tc.SetInput("user", "user")
	tc.SetInput("password", "password")
	tc.SetInput("message", "{\"deviceID\":\"5CCF7F942BCB\",\"distance\":9,\"distState\":\"Safe\"}")
	tc.SetInput("secure", true)
	tc.SetInput("certificate", "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVvakNDQTRxZ0F3SUJBZ0lKQUpsc1QyTTRkL1htTUEwR0NTcUdTSWIzRFFFQkJRVUFNSUdMTVFzd0NRWUQKVlFRR0V3Sk9UREVMTUFrR0ExVUVDQk1DV2tneEVqQVFCZ05WQkFjVENWSnZkSFJsY21SaGJURVhNQlVHQTFVRQpDaE1PVkVsQ1EwOGdVMjltZEhkaGNtVXhEVEFMQmdOVkJBc1RCRk5EVGt3eEVEQU9CZ05WQkFNVUJ5b3ViRzlqCllXd3hJVEFmQmdrcWhraUc5dzBCQ1FFV0VtcDJZVzVrWlhKc1FIUnBZbU52TG1OdmJUQWVGdzB4TnpBek1EY3gKTXpNME1EVmFGdzB4TnpBME1EWXhNek0wTURWYU1JR0xNUXN3Q1FZRFZRUUdFd0pPVERFTE1Ba0dBMVVFQ0JNQwpXa2d4RWpBUUJnTlZCQWNUQ1ZKdmRIUmxjbVJoYlRFWE1CVUdBMVVFQ2hNT1ZFbENRMDhnVTI5bWRIZGhjbVV4CkRUQUxCZ05WQkFzVEJGTkRUa3d4RURBT0JnTlZCQU1VQnlvdWJHOWpZV3d4SVRBZkJna3Foa2lHOXcwQkNRRVcKRW1wMllXNWtaWEpzUUhScFltTnZMbU52YlRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQwpnZ0VCQUpNWCtlcTlteTI4VnhUMzgrRzlFcWJLamFES2RDYVBLcVdWTHM0cFVYdk11SFY1SmowV2hYTnFONXIyClJJcnFFUzFuUVh4REErUmMyaHZ1QzZ0dkJtVFBqOXd6TWVNUGxBL3d2UzNGQUx4S1hLNnJQN1lSUHV5RVJITnkKaUNYNDV6QmRlVkJZZXV3TkR2T2grK250OG5Wa0tBMHBGVGRXdmhNbUJRcEFhazNDeVlueDZXcVExTHVnaXR1SQpHU2JjMFZwUzFnZ1RZSmNiL1N6dk85djJvTkhlcDg4TzBhWlZVVVVMOWVycHJIUE9JdWhYdXhMK1JidWxOVm5NCjZINDlBM2luak9KU3VzUE9vWUxCRVJYcmFPZG9PREtEcUtCTGQ2Z0hyWDdGMXRNL2dibUVCcWNmUlJSbzd5QzkKZUpLRWNyN3M4cnh3S01QNEVMYlJxWWUzQkFjQ0F3RUFBYU9DQVFVd2dnRUJNQjBHQTFVZERnUVdCQlRSNi8zTgpWN2Vmd0d4UExoUXhXRVo0a2VNTXV6Q0J3QVlEVlIwakJJRzRNSUcxZ0JUUjYvM05WN2Vmd0d4UExoUXhXRVo0CmtlTU11NkdCa2FTQmpqQ0JpekVMTUFrR0ExVUVCaE1DVGt3eEN6QUpCZ05WQkFnVEFscElNUkl3RUFZRFZRUUgKRXdsU2IzUjBaWEprWVcweEZ6QVZCZ05WQkFvVERsUkpRa05QSUZOdlpuUjNZWEpsTVEwd0N3WURWUVFMRXdSVApRMDVNTVJBd0RnWURWUVFERkFjcUxteHZZMkZzTVNFd0h3WUpLb1pJaHZjTkFRa0JGaEpxZG1GdVpHVnliRUIwCmFXSmpieTVqYjIyQ0NRQ1piRTlqT0hmMTVqQU1CZ05WSFJNRUJUQURBUUgvTUE4R0ExVWRFUVFJTUFhSEJBb0sKQVRJd0RRWUpLb1pJaHZjTkFRRUZCUUFEZ2dFQkFERjRYWUd3aVRoQTFZak9zZk5UZHgxTzBrdFJGQng4amNNMwpsUE4zMGl5OWhHdkZKTFZOaXlSa1QxOFZsU3BZYmcxclRpd3poR2RrRnJmSEdOdDloOFZtQTVVb2g1NURCUU1BCmFhdWREUi9OQmNuYlVmNnhDa0tqNHVtb3FSRXkyTXFsT0lZTE1wZzBsQ3FZOXFoRTlXWnZuVjk2amxySGVSbWUKdGhnSXhsbjZ6cC9hMGNTVzhQRDY2bnlGdnlGTFlaZjU1N2xqMExCQ29kb29KT2hSdlVoUDNwMnlKRG51NUIzOQp3MFF4ZVp4cGdGSzF1K0t1VUFsdmJrZUp5UGUxQzE4MzdDU2FKajBQM3ByTCt4NHRQVGVQZ1dGWVQ5UGxkRkpYCjVleUp2bUMxMklHZEt3ZS91YzhCelh5bm5Zelk1c0RhYWQwR2x3UDNSalNmQXFmdkYzWT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=")

	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

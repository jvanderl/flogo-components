package systeminfo

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
	"fmt"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("systeminfo")

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

	fmt.Println("Retrieving Sytem Information")

	act.Eval(tc)

	//check result attr

	hostname := tc.GetOutput("hostname")
	ipaddress := tc.GetOutput("ipaddress")

	fmt.Println("hostname: ", hostname)
	fmt.Println("ipaddress: ", ipaddress)

	if hostname == nil {
		t.Fail()
	}

}

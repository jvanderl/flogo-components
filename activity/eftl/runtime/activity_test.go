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
	tc.SetInput("certificate", `-----BEGIN CERTIFICATE-----
MIID+DCCAuCgAwIBAgIJAKHDkDRNkuGeMA0GCSqGSIb3DQEBCwUAMIGQMQswCQYD
VQQGEwJOTDELMAkGA1UECAwCWkgxEjAQBgNVBAcMCVJvdHRlcmRhbTEXMBUGA1UE
CgwOVElCQ08gU29mdHdhcmUxDTALBgNVBAsMBFNDTkwxFTATBgNVBAMMDCouZGVt
by5sb2NhbDEhMB8GCSqGSIb3DQEJARYSanZhbmRlcmxAdGliY28uY29tMB4XDTE3
MDIyNzEwNDUzMFoXDTE3MDMyOTEwNDUzMFowgZAxCzAJBgNVBAYTAk5MMQswCQYD
VQQIDAJaSDESMBAGA1UEBwwJUm90dGVyZGFtMRcwFQYDVQQKDA5USUJDTyBTb2Z0
d2FyZTENMAsGA1UECwwEU0NOTDEVMBMGA1UEAwwMKi5kZW1vLmxvY2FsMSEwHwYJ
KoZIhvcNAQkBFhJqdmFuZGVybEB0aWJjby5jb20wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQCkrrfq2eZrx6UtL11Qc/w4yW3BuFfUNC1REP0f0hbHyi9K
WIc29pCb5K4M5LHJj+oe94xcbNTUs2yXgup8LHFATS3xZxrQBYGiShFXu4JcxLhJ
dmSl0KQd0OLtIksK49H3RqxZ8jG2J+ADMKYhOFk3wKKAKD7n1sIk5N+Rk45Lc/lU
7RvwKMlN17xWFu7BJugakdJWSEI5GWkrXslzGys1L3zoXZ3WLub7TvwHWzhmTOAm
lfp5ymJESHWyxuP4CgbbD2ewnlO20274pwrbsB/XEBkQjffcwmTCM5nrVyR+8EVn
iM2dAJtWr7UyxNr+j66heexKXqSQxcsv4A7qDnL5AgMBAAGjUzBRMB0GA1UdDgQW
BBQXQcis/UVNK1DWTdRNuFnCRpy0CjAfBgNVHSMEGDAWgBQXQcis/UVNK1DWTdRN
uFnCRpy0CjAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBUxls9
zb7FJRyUw+6+dmlAhDCr7+V37VFWuhFcu6TputHacS4YoZimGWxgQ6MNedJvj2h1
xmfnveDsvzpv0xOhSy2EES9GxzfaAGB1wsKIsj7PT8tO9q2UK8Y/xO98gjZqfg3O
JTU05PVBJ2/KuY5ijPR4wUwXErurXa94+ZvGHARSRKAghp6KGbL4BgX7dxZMtc+I
AhSUCmIKVmhuToxLYDMbSKhbGeahIxK3tG3p1Y84lDmtlo3ugtneyaUgNV5F1hun
UATVocCS/zou3ozkI9g8hIQUwn+xdyX8fbbQZVdqQ8Riw6oIzeGBz5u+Zmvtwa63
Xy4qnbHrNsPVrvCL
-----END CERTIFICATE-----`)
	
	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}

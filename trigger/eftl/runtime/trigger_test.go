package eftl

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"testing"
)

const testConfig string = `{
  "name": "eftl",
  "settings": {
    "server": "tibco4.demo.local:9191",
    "channel": "/channel",
    "user": "user",
    "password": "password",
    "secure" : "false",
    "certificate" : ""
  },
  "endpoints": [
    {
      "actionType": "flow",
      "actionURI": "local://testFlow",
      "settings": {
        "destination": "flogo"
      }
    }
  ]
}`

const testConfigSecure string = `{
  "name": "eftl",
  "settings": {
    "server": "tibco4.demo.local:9291",
    "channel": "/channel",
    "user": "user",
    "password": "password",
    "secure" : "true",
    "certificate" : ""
  },
  "endpoints": [
    {
      "actionType": "flow",
      "actionURI": "local://testFlow",
      "settings": {
        "destination": "flogo"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("eftl")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("eftl")
	runner := &TestRunner{}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(config, runner)
}

func TestInitSecure(t *testing.T) {
	tgr := trigger.Get("eftl")
	runner := &TestRunner{}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfigSecure), config)
	config.Settings["certificate"] = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`
	tgr.Init(config, runner)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("eftl")

	tgr.Start()
	defer tgr.Stop()

}

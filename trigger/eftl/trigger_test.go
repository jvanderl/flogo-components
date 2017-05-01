package eftl

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"io/ioutil"
	"testing"
)

var jsonMetadata = getJSONMetadata()

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
	"name": "jvanderl-eftl",
  "settings": {
    "server": "10.10.1.50:19191",
    "channel": "/channel",
    "user": "user",
    "password": "password",
    "secure" : "false",
    "certificate" : ""
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "destination": "flogo"
      }
    }
  ]
}`

/*
const testConfigSecure string = `{
  "name": "eftl",
  "settings": {
    "server": "10.10.1.50:9291",
    "channel": "/channel",
    "user": "user",
    "password": "password",
    "secure" : "true",
    "certificate" : "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVvakNDQTRxZ0F3SUJBZ0lKQUpsc1QyTTRkL1htTUEwR0NTcUdTSWIzRFFFQkJRVUFNSUdMTVFzd0NRWUQKVlFRR0V3Sk9UREVMTUFrR0ExVUVDQk1DV2tneEVqQVFCZ05WQkFjVENWSnZkSFJsY21SaGJURVhNQlVHQTFVRQpDaE1PVkVsQ1EwOGdVMjltZEhkaGNtVXhEVEFMQmdOVkJBc1RCRk5EVGt3eEVEQU9CZ05WQkFNVUJ5b3ViRzlqCllXd3hJVEFmQmdrcWhraUc5dzBCQ1FFV0VtcDJZVzVrWlhKc1FIUnBZbU52TG1OdmJUQWVGdzB4TnpBek1EY3gKTXpNME1EVmFGdzB4TnpBME1EWXhNek0wTURWYU1JR0xNUXN3Q1FZRFZRUUdFd0pPVERFTE1Ba0dBMVVFQ0JNQwpXa2d4RWpBUUJnTlZCQWNUQ1ZKdmRIUmxjbVJoYlRFWE1CVUdBMVVFQ2hNT1ZFbENRMDhnVTI5bWRIZGhjbVV4CkRUQUxCZ05WQkFzVEJGTkRUa3d4RURBT0JnTlZCQU1VQnlvdWJHOWpZV3d4SVRBZkJna3Foa2lHOXcwQkNRRVcKRW1wMllXNWtaWEpzUUhScFltTnZMbU52YlRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQwpnZ0VCQUpNWCtlcTlteTI4VnhUMzgrRzlFcWJLamFES2RDYVBLcVdWTHM0cFVYdk11SFY1SmowV2hYTnFONXIyClJJcnFFUzFuUVh4REErUmMyaHZ1QzZ0dkJtVFBqOXd6TWVNUGxBL3d2UzNGQUx4S1hLNnJQN1lSUHV5RVJITnkKaUNYNDV6QmRlVkJZZXV3TkR2T2grK250OG5Wa0tBMHBGVGRXdmhNbUJRcEFhazNDeVlueDZXcVExTHVnaXR1SQpHU2JjMFZwUzFnZ1RZSmNiL1N6dk85djJvTkhlcDg4TzBhWlZVVVVMOWVycHJIUE9JdWhYdXhMK1JidWxOVm5NCjZINDlBM2luak9KU3VzUE9vWUxCRVJYcmFPZG9PREtEcUtCTGQ2Z0hyWDdGMXRNL2dibUVCcWNmUlJSbzd5QzkKZUpLRWNyN3M4cnh3S01QNEVMYlJxWWUzQkFjQ0F3RUFBYU9DQVFVd2dnRUJNQjBHQTFVZERnUVdCQlRSNi8zTgpWN2Vmd0d4UExoUXhXRVo0a2VNTXV6Q0J3QVlEVlIwakJJRzRNSUcxZ0JUUjYvM05WN2Vmd0d4UExoUXhXRVo0CmtlTU11NkdCa2FTQmpqQ0JpekVMTUFrR0ExVUVCaE1DVGt3eEN6QUpCZ05WQkFnVEFscElNUkl3RUFZRFZRUUgKRXdsU2IzUjBaWEprWVcweEZ6QVZCZ05WQkFvVERsUkpRa05QSUZOdlpuUjNZWEpsTVEwd0N3WURWUVFMRXdSVApRMDVNTVJBd0RnWURWUVFERkFjcUxteHZZMkZzTVNFd0h3WUpLb1pJaHZjTkFRa0JGaEpxZG1GdVpHVnliRUIwCmFXSmpieTVqYjIyQ0NRQ1piRTlqT0hmMTVqQU1CZ05WSFJNRUJUQURBUUgvTUE4R0ExVWRFUVFJTUFhSEJBb0sKQVRJd0RRWUpLb1pJaHZjTkFRRUZCUUFEZ2dFQkFERjRYWUd3aVRoQTFZak9zZk5UZHgxTzBrdFJGQng4amNNMwpsUE4zMGl5OWhHdkZKTFZOaXlSa1QxOFZsU3BZYmcxclRpd3poR2RrRnJmSEdOdDloOFZtQTVVb2g1NURCUU1BCmFhdWREUi9OQmNuYlVmNnhDa0tqNHVtb3FSRXkyTXFsT0lZTE1wZzBsQ3FZOXFoRTlXWnZuVjk2amxySGVSbWUKdGhnSXhsbjZ6cC9hMGNTVzhQRDY2bnlGdnlGTFlaZjU1N2xqMExCQ29kb29KT2hSdlVoUDNwMnlKRG51NUIzOQp3MFF4ZVp4cGdGSzF1K0t1VUFsdmJrZUp5UGUxQzE4MzdDU2FKajBQM3ByTCt4NHRQVGVQZ1dGWVQ5UGxkRkpYCjVleUp2bUMxMklHZEt3ZS91YzhCelh5bm5Zelk1c0RhYWQwR2x3UDNSalNmQXFmdkYzWT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
  },
  "handlers": [
    {
      "actionType": "flow",
      "actionURI": "local://testFlow",
      "settings": {
        "destination": "flogo"
      }
    }
  ]
}`
*/
type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestInit(t *testing.T) {
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)

	// New  factory
	f := &MyFactory{}
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)
}

/*func TestInitSecure(t *testing.T) {
	tgr := trigger.Get("eftl")
	runner := &TestRunner{}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfigSecure), config)
	tgr.Init(config, runner)
}
*/
func TestEndpoint(t *testing.T) {
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &MyFactory{}
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

}

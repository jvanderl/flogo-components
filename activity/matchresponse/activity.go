package matchresponse

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	input      = "input"
	searchdata = "searchdata"
	match      = "match"
	repsonse   = "repsonse"
)

var ifInput = ""

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-matchresponse")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	ifInput := context.GetInput(input).(string)
	log.Debugf("Got input: '%v'", ifInput)
	ifSearchData := context.GetInput(searchdata).([]interface{})
	log.Debugf("Got searchdata: '%v'", ifSearchData)
	response := "not found"
	match := "not found"
	for _, elem := range ifSearchData {
		element := elem.(map[string]interface{})
		find := element["find"].(string)
		log.Debugf("Checking input against : %v", find)
		if strings.Contains(strings.ToLower(ifInput), find) {
			log.Debug("Got a match!")
			match = find
			response = element["resp"].(string)
		}
	}
	context.SetOutput("match", match)
	context.SetOutput("response", response)

	return true, nil
}

package jsontodata

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	//"strings"
)

const (
	input  = "input"
	result = "result"
	data  = "data"
)


// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-jsontodata")

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
	cxinput := context.GetInput(input)
	input := cxinput.(string)

	log.Debug("Got input: '", input, "'")

	byt := []byte(input)

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		context.SetOutput("result", "ERROR_JSON_DECODE")
		return true, nil
	}
	log.Debug("Umarchalled: ", dat)

	/*data := make(map[string]string, len(dat))

	for key, value := range dat {
		log.Debug("got key: ", key, ", value: ", value)
		data[key] = value.(string)
		}
	context.SetOutput("data", data) */
	context.SetOutput("data", dat)
	context.SetOutput("result", "OK")

	return true, nil
}

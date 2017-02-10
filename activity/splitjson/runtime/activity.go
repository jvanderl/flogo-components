package splitjson

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
	//	"fmt"
)

const (
	input  = "input"
	result = "result"
	name1  = "name1"
	value1 = "value1"
	name2  = "name2"
	value2 = "value2"
	name3  = "name3"
	value3 = "value3"
	name4  = "name4"
	value4 = "value4"
	name5  = "name5"
	value5 = "value5"
	name6  = "name6"
	value6 = "value6"
	name7  = "name7"
	value7 = "value7"
	name8  = "name8"
	value8 = "value8"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-rest")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&MyActivity{metadata: md})
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

	context.SetOutput("name1", "")

	log.Debug("Umarchalled: ", dat)

	index := 1
	for key, value := range dat {
		log.Debug("got key: ", key, ", value: ", value)
		switch index {
		case 1:
			{
				context.SetOutput("name1", key)
				context.SetOutput("value1", value)
			}
		case 2:
			{
				context.SetOutput("name2", key)
				context.SetOutput("value2", value)
			}
		case 3:
			{
				context.SetOutput("name3", key)
				context.SetOutput("value3", value)
			}
		case 4:
			{
				context.SetOutput("name4", key)
				context.SetOutput("value4", value)
			}
		case 5:
			{
				context.SetOutput("name5", key)
				context.SetOutput("value5", value)
			}
		case 6:
			{
				context.SetOutput("name6", key)
				context.SetOutput("value6", value)
			}
		case 7:
			{
				context.SetOutput("name7", key)
				context.SetOutput("value7", value)
			}
		case 8:
			{
				context.SetOutput("name8", key)
				context.SetOutput("value8", value)
			}
		}
		index++
	}
	context.SetOutput("result", "OK")

	return true, nil
}

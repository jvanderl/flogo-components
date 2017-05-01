package getjson

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
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

var ifName1 = ""
var ifName2 = ""
var ifName3 = ""
var ifName4 = ""
var ifName5 = ""
var ifName6 = ""
var ifName7 = ""
var ifName8 = ""

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-getjson")

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
	log.Debug("Got input: '", ifInput, "'")
	ifName1, _ := context.GetInput(name1).(string)
	ifName2, _ := context.GetInput(name2).(string)
	ifName3, _ := context.GetInput(name3).(string)
	ifName4, _ := context.GetInput(name4).(string)
	ifName5, _ := context.GetInput(name5).(string)
	ifName6, _ := context.GetInput(name6).(string)
	ifName7, _ := context.GetInput(name7).(string)
	ifName8, _ := context.GetInput(name8).(string)

	byt := []byte(ifInput)

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		context.SetOutput("result", "ERROR_JSON_DECODE")
		return true, nil
	}

	log.Debug("Umarchalled: ", dat)

	context.SetOutput(value1, dat[ifName1])
	context.SetOutput(value2, dat[ifName2])
	context.SetOutput(value3, dat[ifName3])
	context.SetOutput(value4, dat[ifName4])
	context.SetOutput(value5, dat[ifName5])
	context.SetOutput(value6, dat[ifName6])
	context.SetOutput(value7, dat[ifName7])
	context.SetOutput(value8, dat[ifName8])

	/*

	       for key, value := range dat {
	   	    log.Debug("got key: ", key, ", value: ", value)
	       	switch key {
	       		case ifName1 : {
	       			context.SetOutput("value1", value)
	       			}
	       		case ifName2 : {
	       			context.SetOutput("value2", value)
	       			}
	       		case ifName3 : {
	       			context.SetOutput("value3", value)
	       			}
	       		case ifName4 : {
	       			context.SetOutput("value4", value)
	       			}
	       		case ifName5 : {
	       			context.SetOutput("value5", value)
	       			}
	       		case ifName6 : {
	       			context.SetOutput("value6", value)
	       			}
	       		case ifName7 : {
	       			context.SetOutput("value7", value)
	   	    			}
	       		case ifName8 : {
	       			context.SetOutput("value8", value)
	       			}
	       		}
	       }
	*/
	context.SetOutput("result", "OK")

	return true, nil
}

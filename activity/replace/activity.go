package replace

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strings"
)

const (
	find  = "find"
	replace = "replace"
	result = "result"
	input1  = "input1"
	output1 = "output1"
	input2  = "input2"
	output2 = "output2"
	input3  = "input3"
	output3 = "output3"
	input4  = "input4"
	output4 = "output4"
	input5  = "input5"
	output5 = "output5"
	input6  = "input6"
	output6 = "output6"
	input7  = "input7"
	output7 = "output7"
	input8  = "input8"
	output8 = "output8"
)

var ifInput1 = ""
var ifInput2 = ""
var ifInput3 = ""
var ifInput4 = ""
var ifInput5 = ""
var ifInput6 = ""
var ifInput7 = ""
var ifInput8 = ""

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-trim")

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
	ifFind := context.GetInput(find).(string)
	log.Debug("Got find: '", ifFind, "'")
	ifReplace := context.GetInput(replace).(string)
	log.Debug("Got replace: '", ifReplace, "'")
	ifInput1, _ := context.GetInput(input1).(string)
	ifInput2, _ := context.GetInput(input2).(string)
	ifInput3, _ := context.GetInput(input3).(string)
	ifInput4, _ := context.GetInput(input4).(string)
	ifInput5, _ := context.GetInput(input5).(string)
	ifInput6, _ := context.GetInput(input6).(string)
	ifInput7, _ := context.GetInput(input7).(string)
	ifInput8, _ := context.GetInput(input8).(string)

	context.SetOutput("output1", strings.Replace(ifInput1, ifFind, ifReplace, -1))
	context.SetOutput("output2", strings.Replace(ifInput2, ifFind, ifReplace, -1))
	context.SetOutput("output3", strings.Replace(ifInput3, ifFind, ifReplace, -1))
	context.SetOutput("output4", strings.Replace(ifInput4, ifFind, ifReplace, -1))
	context.SetOutput("output5", strings.Replace(ifInput5, ifFind, ifReplace, -1))
	context.SetOutput("output6", strings.Replace(ifInput6, ifFind, ifReplace, -1))
	context.SetOutput("output7", strings.Replace(ifInput7, ifFind, ifReplace, -1))
	context.SetOutput("output8", strings.Replace(ifInput8, ifFind, ifReplace, -1))
	context.SetOutput("result", "OK")

	return true, nil
}

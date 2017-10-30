package [package]

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"[import1]"
	"[import2]"
)

const (
	[input1_name] = "[input1_name]"
	[input2_name] = "[input2_name]"
	[output1_name] = "[output1_name]"
	[output2_name] = "[output2_name]"
)

//initialize input parameters to avoid "input nil" exceptions
var if[input1_name] [input1_type] = [input1_default]
var if[input2_name] [input2_type] = [input2_default]

//initialize output parameters
var if[output1_name] [output1_type] = [output1_default]
var if[output2_name] [output2_type] = [output2_default]

// log is the default package logger
var log = logger.GetLogger("activity-[git_user]-[package]")

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

	// read input parameters
	if[input1_name] := context.GetInput([input1_name]).([input1_type])
	log.Debug("Got [input1_name]: '", if[input1_name], "'")
	if[input2_name] := context.GetInput([input2_name]).([input2_type])
	log.Debug("Got [input2_name]: '", if[input2_name], "'")

	// **** YOUR LOGIC HERE ****


	//set output parameters
	context.SetOutput([output1_name], if[output1_name])
	context.SetOutput([output2_name], if[output2_name])

	return true, nil
}

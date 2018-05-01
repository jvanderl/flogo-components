package statechange

import (
	"strconv"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	datasource = "datasource"
	input      = "input"
	changed    = "changed"
	flags      = "flags"
	result     = "result"
	laststate  = "laststate"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-statechange")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	sync.Mutex
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

	//first check if unique source is set
	ifDataSource := context.GetInput(datasource).(string)
	if ifDataSource == "" {
		//data source is not set
		context.SetOutput(result, "DATASOURCE_NOT_SET")
		return true, nil
	}

	gotAnyInputs := false
	stateChanged := false
	resultText := ""
	ifFlags := 0

	// check each input
	for i := 1; i < 9; i++ {
		inputName := input + strconv.Itoa(i)
		//		log.Debugf("Checking input: %s", inputName)
		ifInput, ok := context.GetInput(inputName).(string)
		if ok && ifInput != "" {
			//input is not empty.
			//			log.Debugf("Input value: %s", ifInput)
			gotAnyInputs = true
			prevStateLocation := ifDataSource + "." + laststate + "." + inputName
			log.Debugf("Checking previous State Data at : %s", prevStateLocation)

			//data source is set, now try get history from global data
			var attribInput interface{} = ifInput
			log.Debugf("input attribute is %v", attribInput)
			var prevstatedat *data.Attribute
			prevstatedat, ok := data.GetGlobalScope().GetAttr(prevStateLocation)
			if ok {
				//got previous state
				prevState := prevstatedat.Value().(string)
				log.Debugf("previous state value is %v", prevState)
				if prevState != ifInput {
					//state has changed
					stateChanged = true
					log.Debug("State has changed")
					// set flag
					ifFlags |= (1 << uint(i-1))
					//					log.Debugf ("Flag is now set at: %s", ifFlags)
					// set result text
					if resultText != "" {
						resultText = resultText + ", "
					}
					resultText = resultText + inputName + " changed from " + prevState + " to " + ifInput
					// update previous state
					//					data.GetGlobalScope().SetAttrValue(prevStateLocation, ifInput)
					data.GetGlobalScope().SetAttrValue(prevStateLocation, attribInput)

				} else {
					//state has not changed
					//					log.Debug("State has not changed")
				}
			} else {
				//got no previous state
				//state has changed
				stateChanged = true
				log.Debug("State has changed because no previous state was found")
				// set flag
				ifFlags |= (1 << uint(i-1))
				//				log.Debugf ("Flag is now set at: %s", ifFlags)
				// set result text
				if resultText != "" {
					resultText = resultText + ", "
				}
				resultText = resultText + inputName + " is initally set to " + ifInput
				// add current input data to prevState
				dt, ok := data.ToTypeEnum("string")
				if ok {
					//					data.GetGlobalScope().AddAttr(prevStateLocation, dt, ifInput)
					data.GetGlobalScope().AddAttr(prevStateLocation, dt, attribInput)
				}
			}
		}
	}
	if !gotAnyInputs {
		context.SetOutput(result, "ERR_NO_INPUTS")
		return true, nil
	}

	if !stateChanged {
		resultText = "NO_CHANGE"
	}

	// When inputs were set, set the outputs
	context.SetOutput(changed, stateChanged)
	context.SetOutput(flags, ifFlags)
	context.SetOutput(result, resultText)

	return true, nil
}

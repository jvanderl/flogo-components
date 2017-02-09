package filter

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
	"strconv"
	"fmt"
)

const (
	input 		= "input"
	datatype 	= "datatype"
	minvalue 	= "minvalue"
	maxvalue 	= "maxvalue"
	inverse 	= "inverse"
	pass 		= "pass"
	reason 		= "reason"
)

var ifInput interface{}
var ifMinValue interface {}
var ifMaxValue interface {}
var ifInverse = false
var valueTooLow = false
var valueTooHigh = false
var minimumSet = false
var maximumSet = false

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
func (a *MyActivity) Eval(context activity.Context) (done bool, err error)  {

	// check input value data type
	datatypeInput := context.GetInput(datatype)
	ivdatatype, ok := datatypeInput.(string)
	if !ok {
		context.SetOutput(pass, false)
		context.SetOutput(reason, "DATATYPE_NOT_SET")
		return true, fmt.Errorf("Data type not set.")
	}

	// check the value matches the data type
	ifInput = validateValue(context, input, ivdatatype)
	if ifInput == nil {
		context.SetOutput(pass, false)
		context.SetOutput(reason, "INPUT_INVALID")
		return true, fmt.Errorf("Invalid input data.")
	}

	// check if inverse is set and applu
	cxInverse := context.GetInput(inverse)
	ivInverse, ok := cxInverse.(bool)
	if ok {
		// overwrite default inverse with input from context
		ifInverse = ivInverse
	}
	
	// check if minimum value is set and apply filter
	cxMinValue := context.GetInput(minvalue)
	if cxMinValue != "" {
		// there is a minvalue assigned, now check for validity
		ifMinValue = validateValue(context, minvalue, ivdatatype)
				if ifMinValue == nil {
				context.SetOutput(pass, false)
				context.SetOutput(reason, "MIN_VALUE_INVALID")
				return true, fmt.Errorf("Invalid minimum value.")
		} else {
			minimumSet = true
			if !ifInverse {
				// normal filter
				switch ivdatatype {
					case "int" : if ifInput.(int) < ifMinValue.(int) {valueTooLow = true}
					case "uint" : if ifInput.(uint) < ifMinValue.(uint) {valueTooLow = true}
					case "float32" : if ifInput.(float32) < ifMinValue.(float32) {valueTooLow = true}
				}
			} else {
				// inverse filter
				switch ivdatatype {
					case "int" : if ifInput.(int) > ifMinValue.(int) {valueTooHigh = true}
					case "uint" : if ifInput.(uint) > ifMinValue.(uint) {valueTooHigh = true}
					case "float32" : if ifInput.(float32) > ifMinValue.(float32) {valueTooHigh = true}
				}
			}
		}		
	}

	// check if maximum value is set and apply filter
	cxMaxValue := context.GetInput(maxvalue)
	if cxMaxValue != "" {
		// there is a minvalue assigned, now check for validity
		ifMaxValue = validateValue(context, maxvalue, ivdatatype)
			if ifMaxValue == nil {
			context.SetOutput(pass, false)
			context.SetOutput(reason, "MAX_VALUE_INVALID")
			return true, fmt.Errorf("Invalid maximum value.")
		} else {
			maximumSet = true
			if !ifInverse {
				// normal filter
				switch ivdatatype {
					case "int" : if ifInput.(int) > ifMaxValue.(int) {valueTooHigh = true}
					case "uint" : if ifInput.(uint) > ifMaxValue.(uint) {valueTooHigh = true}
					case "float32" : if ifInput.(float32) > ifMaxValue.(float32) {valueTooHigh = true}
				}
			} else {
				// inverse filter
				switch ivdatatype {
					case "int" : if ifInput.(int) < ifMaxValue.(int) {valueTooLow = true}
					case "uint" : if ifInput.(uint) < ifMaxValue.(uint) {valueTooLow = true}
					case "float32" : if ifInput.(float32) < ifMaxValue.(float32) {valueTooLow = true}
				}
			}
		}		
	}

	if (!ifInverse && valueTooLow) || (ifInverse && !minimumSet && valueTooLow) {
		// normal filter and value is too low or inverse filter, no minumum set and value too low
		context.SetOutput(pass, false)
		context.SetOutput(reason, "VALUE_TOO_LOW")
		return true, nil
	}
	if (!ifInverse && valueTooHigh) || (ifInverse && !maximumSet && valueTooHigh)  {
		// normal filter and value is too high or inverse filter, no maximum set and value too high
		context.SetOutput(pass, false)
		context.SetOutput(reason, "VALUE_TOO_HIGH")
		return true, nil
	}
	if (ifInverse && minimumSet && valueTooHigh && maximumSet && valueTooLow) {
		// normal filter and value is too high or inverse filter, no maximum set and value too high
		context.SetOutput(pass, false)
		context.SetOutput(reason, "MID_SECTION_FILTERED")
		return true, nil
	}
	// When not filtered out, put the input data in output
	context.SetOutput(pass, true)
	context.SetOutput(reason, "FILTER_PASSED")
	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils
////////////////////////////////////////////////////////////////////////////////////////

func validateValue(context activity.Context, element string, datatype string) interface{}  {

	cxDataInput := context.GetInput(element)
	dataInput := cxDataInput.(string)

	switch datatype {
		case "int": dataOutput, err := strconv.ParseInt(dataInput, 10, strconv.IntSize)
			if err == nil {
				return int(dataOutput)
			}
		case "uint": dataOutput, err := strconv.ParseUint(dataInput, 10, strconv.IntSize)
			if err == nil {
				return uint(dataOutput)
			}
		case "float32": dataOutput, err := strconv.ParseFloat(dataInput, 32)
			if err == nil {
				return float32(dataOutput)
			}
		default:
			return nil
	}
	return nil
}
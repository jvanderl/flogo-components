package filter

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/op/go-logging"
	"strconv"
	"time"
	"fmt"
	"sync"
)

const (
	input   	 = "input"
	datasource   = "datasource"
	datatype     = "datatype"
	minvalue 	 = "minvalue"
	maxvalue     = "maxvalue"
	interval     = "interval"
	intervaltype = "intervaltype"
	timelayout   = time.RFC3339Nano
)

var ifInput interface{}
var ifMinValue interface {}
var ifMaxValue interface {}
var ifLastTimePassed = ""
var valueTooLow = false
var valueTooHigh = false
var actualInterval = 0
var intervalTooShort = false

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-rest")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	sync.Mutex
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

	//filter out by default
	context.SetOutput("usevalue", false)
	context.SetOutput("result", "")

	// check input value data type
	datatypeInput := context.GetInput(datatype)
	ivdatatype, ok := datatypeInput.(string)
	if !ok {
		context.SetOutput("reason", "DATATYPE_NOT_SET")
		return true, fmt.Errorf("Data type not set.")
	}

	// check the value matches the data type
	ifInput = validateValue(context, input, ivdatatype)
	if ifInput == nil {
		context.SetOutput("reason", "INPUT_INVALID")
		return true, fmt.Errorf("Invalid input data.")
	}
	
	// check if minimum value is set and apply filter
	cxMinValue := context.GetInput(minvalue)
	if cxMinValue != "" {
		// there is a minvalue assigned, now check for validity
		ifMinValue = validateValue(context, minvalue, ivdatatype)
				if ifMinValue == nil {
				context.SetOutput("reason", "MIN_VALUE_INVALID")
				return true, fmt.Errorf("Invalid minimum value.")
		} else {
			switch ivdatatype {
				case "int" : if ifInput.(int) < ifMinValue.(int) {valueTooLow = true}
				case "uint" : if ifInput.(uint) < ifMinValue.(uint) {valueTooLow = true}
				case "float32" : if ifInput.(float32) < ifMinValue.(float32) {valueTooLow = true}
			}
			if valueTooLow {
					context.SetOutput("reason", "VALUE_TOO_LOW")
					return true, nil
			}
		}		
	}

	// check if maximum value is set and apply filter
	cxMaxValue := context.GetInput(maxvalue)
	if cxMaxValue != "" {
		// there is a minvalue assigned, now check for validity
		ifMaxValue = validateValue(context, maxvalue, ivdatatype)
				if ifMinValue == nil {
				context.SetOutput("reason", "MAX_VALUE_INVALID")
				return true, fmt.Errorf("Invalid maximum value.")
		} else {
			switch ivdatatype {
				case "int" : if ifInput.(int) > ifMaxValue.(int) {valueTooHigh = true}
				case "uint" : if ifInput.(uint) > ifMaxValue.(uint) {valueTooHigh = true}
				case "float32" : if ifInput.(float32) > ifMaxValue.(float32) {valueTooHigh = true}
			}
			if valueTooHigh {
					context.SetOutput("reason", "VALUE_TOO_HIGH")
					return true, nil
			}
		}		
	}

	//check against interval
//	log.Debug("***** About to test Interval *****")
	ifInterval, ok := context.GetInput(interval).(int)
	if ok && ifInput != 0 {
		//interval is set, now check if interval type is set
//		log.Debug("***** Interval is set at ", ifInterval, " *****")
		ifIntervalType := context.GetInput(intervaltype).(string)
		if ifIntervalType == "" {
//			log.Debug("***** Interval is empty *****")
			context.SetOutput("reason", "INTERVAL_TYPE_NOT_SET")
			return true, nil
		} else {
			//now check if unique source is set
//			log.Debug("***** Checking data source *****")
			ifDataSource := context.GetInput(datasource).(string)
			if ifDataSource == "" {
				//data source is not set
//				log.Debug("***** Data Source is not set *****")
				context.SetOutput("reason", "DATASOURCE_NOT_SET")
				return true, nil
			} else {
				//data source is set, now try get history from global data
//				log.Debug("***** DataSource is set at ", ifDataSource, ", getting last time from Globalvar *****")
				lasttimedat, ok := data.GetGlobalScope().GetAttr(ifDataSource)
				if ok {
					// previous data found, go ahead and analyse
//					log.Debug("***** found last time in global: ", lasttimedat.Value.(string), "*****")
					ifLastTimePassed := lasttimedat.Value.(string)
//					log.Debug("***** converting last time string to time *****")
					lasttime , err := time.Parse(timelayout, ifLastTimePassed)
					if err != nil {
						//invalid previous timestamp format
//						log.Debug ("***** unable to convert string to time *****")
						return true, fmt.Errorf("Invalid time format in history")
					} else {
						// valid previous timestamp, now check duration against interval
//						log.Debug ("***** Time converted OK, now calculating the interval duration *****")
						duration := time.Since(lasttime)
//						log.Debug ("***** duration is ", duration, " *****")
						switch ifIntervalType {
							case "hours"   : intervalTooShort = (int(duration.Hours()) < ifInterval)
							case "minutes" : intervalTooShort = (int(duration.Minutes()) < ifInterval) 
							case "seconds" : intervalTooShort = (int(duration.Seconds()) < ifInterval) 
							case "milliseconds" : intervalTooShort = (int(duration.Nanoseconds()/1e6) < ifInterval)
							default : {
//								log.Debug("***** Invalid Interval Type: ", ifIntervalType, " *****")
								context.SetOutput("reason", "INVALID_INTERVAL_TYPE")
								return true, nil
								}
						}
						if intervalTooShort {
							// will filter this one out
//							log.Debug("***** Interval was too short *****")
							context.SetOutput("reason", "TROTTLED_BY_INTERVAL")
							return true, nil
						} else {
							// as this is the final filter, value will be used. Update last time in global
//							log.Debug("***** Interval was OK, updating Gobal var with current timestamp")
							ifLastTimePassed = string(time.Now().Format(timelayout))
							data.GetGlobalScope().SetAttrValue(ifDataSource, ifLastTimePassed )

						}
					}
				} else {
					// did not get data from global var earlier, so we'll go ahead and update Global var 
//					log.Debug("***** Initial run, updating Gobal var with current timestamp")
					ifLastTimePassed = string(time.Now().Format(timelayout))
					dt, ok := data.ToTypeEnum("string")
					if ok {
						data.GetGlobalScope().AddAttr(ifDataSource, dt, ifLastTimePassed)
					}
				}
			}
		}
	}

	// When not filtered out, put the input data in output
	context.SetOutput("result", ifInput)
	context.SetOutput("usevalue", true)
	context.SetOutput("lasttimepassed", ifLastTimePassed)
	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils
////////////////////////////////////////////////////////////////////////////////////////

func validateValue(context activity.Context, element string, datatype string) interface{}  {

	dataInput := context.GetInput(element)

	switch datatype {
		case "int": dataOutput, err := strconv.ParseInt(dataInput.(string), 10, strconv.IntSize)
			if err == nil {
				return int(dataOutput)
			}
		case "uint": dataOutput, err := strconv.ParseUint(dataInput.(string), 10, strconv.IntSize)
			if err == nil {
				return uint(dataOutput)
			}
		case "float32": dataOutput, err := strconv.ParseFloat(dataInput.(string), 32)
			if err == nil {
				return float32(dataOutput)
			}
		default:
			return nil
	}
	return nil
}
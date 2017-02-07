package throttle

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/op/go-logging"
	"time"
	"fmt"
	"sync"
)

const (
	datasource   = "datasource"
	interval     = "interval"
	intervaltype = "intervaltype"
	pass 		 = "pass"
	reason		 = "reason"
	lasttimepassed = "lasttimepassed"
	timelayout   = time.RFC3339Nano
)

var ifLastTimePassed = ""
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

// trhottle by default
	context.SetOutput(pass, false)

	//check against interval
	log.Debug("Checking against interval")
	ifInterval, ok := context.GetInput(interval).(int)
	if ok && ifInterval != 0 {
		//interval is set, now check if interval type is set
		ifIntervalType := context.GetInput(intervaltype).(string)
		if ifIntervalType == "" {
			context.SetOutput(reason, "INTERVAL_TYPE_NOT_SET")
			return true, nil
		} else {
			//now check if unique source is set
			ifDataSource := context.GetInput(datasource).(string)
			if ifDataSource == "" {
				//data source is not set
				context.SetOutput(reason, "DATASOURCE_NOT_SET")
				return true, nil
			} else {
				//data source is set, now try get history from global data
				lasttimedat, ok := data.GetGlobalScope().GetAttr(ifDataSource)
				if ok {
					// previous data found, go ahead and analyse
					ifLastTimePassed := lasttimedat.Value.(string)
					lasttime , err := time.Parse(timelayout, ifLastTimePassed)
					if err != nil {
						//invalid previous timestamp format
						return true, fmt.Errorf("Invalid time format in history")
					} else {
						// valid previous timestamp, now check duration against interval
						log.Debug("Retrieved valid previous timestamp")
						duration := time.Since(lasttime)
						switch ifIntervalType {
							case "hours"   : intervalTooShort = (int(duration.Hours()) < ifInterval)
							case "minutes" : intervalTooShort = (int(duration.Minutes()) < ifInterval) 
							case "seconds" : intervalTooShort = (int(duration.Seconds()) < ifInterval) 
							case "milliseconds" : intervalTooShort = (int(duration.Nanoseconds()/1e6) < ifInterval)
							default : {
								// Invalid Interval Type
								context.SetOutput(reason, "INVALID_INTERVAL_TYPE")
								return true, nil
								}
						}
						if intervalTooShort {
							// will filter this one out
							log.Debug("Time since last timestamp is shorter than interval, will not pass")
							context.SetOutput(lasttimepassed, ifLastTimePassed)
							context.SetOutput(reason, "TROTTLED_BY_INTERVAL")
							return true, nil
						} else {
							// as this is the final filter, value will be used. Update last time in global
							log.Debug("Time since last timestamp is longer than interval, will pass")
							ifLastTimePassed = string(time.Now().Format(timelayout))
							context.SetOutput(lasttimepassed, ifLastTimePassed)
							data.GetGlobalScope().SetAttrValue(ifDataSource, ifLastTimePassed )

						}
					}
				} else {
					// did not get data from global var earlier, so we'll go ahead and update Global var 
					log.Debug("Retrieved no previous timestamp, will pass")
					ifLastTimePassed = string(time.Now().Format(timelayout))
					context.SetOutput(lasttimepassed, ifLastTimePassed)
					dt, ok := data.ToTypeEnum("string")
					if ok {
						data.GetGlobalScope().AddAttr(ifDataSource, dt, ifLastTimePassed)
					}
				}
			}
		}
	}

	// When not filtered out, put the input data in output
	context.SetOutput(pass, true)
	context.SetOutput(reason, "THROTTLE_PASSED")

	return true, nil
}

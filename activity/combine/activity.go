package combine

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strings"
)

const (
	delimiter = "delimiter"
	prefix    = "prefix"
	suffix    = "suffix"
	part1     = "part1"
	part2     = "part2"
	part3     = "part3"
	part4     = "part4"
	part5     = "part5"
	part6     = "part6"
	part7     = "part7"
	part8     = "part8"
	result    = "result"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-combine")

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
	ivDelimiter := getInputParameter(context, delimiter)
	ivPrefix := getInputParameter(context, prefix)
	ivSuffix := getInputParameter(context, suffix)
	ivPart1 := getInputParameter(context, part1)
	ivPart2 := getInputParameter(context, part2)
	ivPart3 := getInputParameter(context, part3)
	ivPart4 := getInputParameter(context, part4)
	ivPart5 := getInputParameter(context, part5)
	ivPart6 := getInputParameter(context, part6)
	ivPart7 := getInputParameter(context, part7)
	ivPart8 := getInputParameter(context, part8)
	/*	log.Debugf("delimiter: [%s]", ivDelimiter)
		log.Debugf("prefix: [%s]", ivPrefix)
		log.Debugf("suffix: [%s]", ivSuffix)
		log.Debugf("part1: [%s]", ivPart1)
		log.Debugf("part2: [%s]", ivPart2)
		log.Debugf("part3: [%s]", ivPart3)
		log.Debugf("part4: [%s]", ivPart4)
		log.Debugf("part5: [%s]", ivPart5)
		log.Debugf("part6: [%s]", ivPart6)
		log.Debugf("part7: [%s]", ivPart7)
		log.Debugf("part8: [%s]", ivPart8)
	*/
	var ivResult = ""

	//	log.Debug("Adding parts to result")

	ivResult = addPart(ivResult, ivPart1, ivDelimiter)
	ivResult = addPart(ivResult, ivPart2, ivDelimiter)
	ivResult = addPart(ivResult, ivPart3, ivDelimiter)
	ivResult = addPart(ivResult, ivPart4, ivDelimiter)
	ivResult = addPart(ivResult, ivPart5, ivDelimiter)
	ivResult = addPart(ivResult, ivPart6, ivDelimiter)
	ivResult = addPart(ivResult, ivPart7, ivDelimiter)
	ivResult = addPart(ivResult, ivPart8, ivDelimiter)

	//	log.Debugf("Result is now: [%s]", ivResult)

	//	log.Debug("Processing prefix")
	if ivPrefix != "" { // we have a prefix, result starts with this
		if ivDelimiter != "" { // add delimiter at start if present
			ivResult = ivDelimiter + ivResult
		}
		ivResult = ivPrefix + ivResult
	} else { // we don't have a prefix, delete delimiter if present
		if ivDelimiter != "" {
			strings.TrimPrefix(ivResult, ivDelimiter)
		}
	}

	//	log.Debugf("Result is now: [%s]", ivResult)

	//	log.Debug("Processing suffix")

	if ivSuffix != "" { // we have a suffix, should be added to result
		if ivDelimiter != "" { // when we have a delimiter, check if suffix starts with it
			if strings.HasPrefix(ivSuffix, ivDelimiter) { // already starts with delimiter, just add it
				ivResult = ivResult + ivSuffix
			} else {
				// suffix does not start with dilimiter, add it as well
				ivResult = ivResult + ivDelimiter + ivSuffix
			}
		} else { // there's no delimiter, just add the suffix
			ivResult = ivResult + ivSuffix
		}
	}

	//	log.Debugf("Result is now: [%s]. This is what we'll return.", ivResult)

	context.SetOutput(result, ivResult)

	return true, nil
}

func getInputParameter(context activity.Context, parameter string) string {
	output, ok := context.GetInput(parameter).(string)
	if !ok {
		return ""
	}
	return output
}

func addPart(input string, part string, delimiter string) string {
	if part == "" { // no part means no change to input
		return input
	}
	if delimiter == "" || (delimiter != "" && input == "") { // no delimiter or first entry without having prefix, just add part to input
		//	if delimiter == "" { // no delimiter, just add part to input
		return input + part
	}
	// part and delimiter, add both to input
	return input + delimiter + part
}

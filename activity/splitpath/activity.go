package splitpath

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strings"
)

const (
	input     = "input"
	delimiter = "delimiter"
	fixedpath = "fixedpath"
	result    = "result"
	part1     = "part1"
	part2     = "part2"
	part3     = "part3"
	part4     = "part4"
	part5     = "part5"
	part6     = "part6"
	part7     = "part7"
	part8     = "part8"
	remainder = "remainder"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-splitpath")

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
	ivInput := context.GetInput(input).(string)
	ivDelimiter := context.GetInput(delimiter).(string)
	ivFixedPath := context.GetInput(fixedpath).(string)
	//	log.Debug("input: '", ivInput, "'")
	//	log.Debug("delimiter: '", ivDelimiter, "'")
	//	log.Debug("fixedpath: '", ivFixedPath, "'")

	log.Debug("Proecssing fixed path")
	if ivFixedPath != "" {
		//forward the fixedpath to output
		context.SetOutput(fixedpath, ivFixedPath)
		//		log.Debug("Checking if start of input matches fixed path")
		//check if input actually starts with fixedpath
		if strings.HasPrefix(ivInput, ivFixedPath) == false {
			context.SetOutput(result, "PREFIX_MISMATCH")
			return true, nil
		}
		// strip fixed path from input
		//			log.Debug("fixed path matches, stripping it from the input")
		ivInput = strings.TrimPrefix(ivInput, ivFixedPath)
		//			log.Debug("input to split is now: ", ivInput)
	}

	// check if input starts with delimiter, remove it
	//	log.Debug("Checking if input starts with delimiter")
	if strings.HasPrefix(ivInput, ivDelimiter) {
		//		log.Debug("Input starts with delimiter, removing delimiter at start")
		ivInput = strings.TrimPrefix(ivInput, ivDelimiter)
		//		log.Debug("input to split is now: ", ivInput)
	}

	log.Debug("Splitting the input")
	parts := strings.Split(ivInput, ivDelimiter)
	numparts := len(parts)
	log.Debug("Number of parts: ", numparts)

	// check if there are more parts than 8
	log.Debug("Populating parts")
	if numparts > 8 {
		//		log.Debug("There are more than 8 parts")

		// first fill the individual parts
		for i := 0; i < 8; i++ {
			//			log.Debug("In for loop filling part: ", i)
			position := numparts - i - 1
			switch i {
			case 0:
				context.SetOutput(part1, parts[position])
			case 1:
				context.SetOutput(part2, parts[position])
			case 2:
				context.SetOutput(part3, parts[position])
			case 3:
				context.SetOutput(part4, parts[position])
			case 4:
				context.SetOutput(part5, parts[position])
			case 5:
				context.SetOutput(part6, parts[position])
			case 6:
				context.SetOutput(part7, parts[position])
			case 7:
				context.SetOutput(part8, parts[position])
			}
		}
		// now create the reamainder string
		log.Debug("Populating remainder")
		remainparts := parts[0 : numparts-8]
		ovRemainder := strings.Join(remainparts, ivDelimiter)
		context.SetOutput(remainder, ovRemainder)
	} else {
		//		log.Debug("There are no more than 8 parts")
		//all can go into parts, start with highest part
		for i := 0; i < numparts; i++ {
			//			log.Debug("In for loop filling part: ", i)
			position := numparts - i
			switch position {
			case 1:
				context.SetOutput(part1, parts[i])
			case 2:
				context.SetOutput(part2, parts[i])
			case 3:
				context.SetOutput(part3, parts[i])
			case 4:
				context.SetOutput(part4, parts[i])
			case 5:
				context.SetOutput(part5, parts[i])
			case 6:
				context.SetOutput(part6, parts[i])
			case 7:
				context.SetOutput(part7, parts[i])
			case 8:
				context.SetOutput(part8, parts[i])
			}
		}
	}

	context.SetOutput(result, "OK")
	return true, nil
}

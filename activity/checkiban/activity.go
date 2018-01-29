package checkiban

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/almerlucke/go-iban/iban"
)

var log = logger.GetLogger("activity-jvanderl-checkiban")

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
	ivIban := context.GetInput("iban").(string)

	code := ""
	printcode := ""
	result := ""

	log.Debugf("Checking IBAN: %v", ivIban)

	iban, err := iban.NewIBAN(ivIban)

	if err != nil {
		result = err.Error()
		log.Debugf("IBAN check failed: %v", result)
	} else {
		result = "OK"
		printcode = iban.PrintCode
		code = iban.Code
		log.Debugf("IBAN check OK: %v", printcode)
	}

	context.SetOutput("result", result)
	context.SetOutput("code", code)
	context.SetOutput("printcode", printcode)

	return true, nil
}

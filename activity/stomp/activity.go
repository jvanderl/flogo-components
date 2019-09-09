package stomp

import (
	"github.com/go-stomp/stomp"
	"github.com/project-flogo/core/activity"
)

func init() {
	_ = activity.Register(&Activity{})
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity is an activity that is used to invoke a REST Operation
// input    : {address, destination, message}
// outputs  : {result}
type Activity struct {
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval - Sends a message to Slack
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	ctx.Logger().Infof("Sending message on destination: %s", input.Address)

	output := &Output{}

	//Validate the inputs
	if input.Address == "" {
		output.Result = "ERR_ADDRESS_NOT_DEFINED"
		return false, nil
	}

	if input.Destination == "" {
		output.Result = "ERR_DESTINATION_NOT_DEFINED"
		return false, nil
	}

	//Connect to Stomp
	conn, err := stomp.Dial("tcp", input.Address)
	if err != nil {
		output.Result = "ERR_STOMP_CONNECT"
		return false, nil
	}

	defer conn.Disconnect()

	//Send the message
	err = conn.Send(
		input.Destination,     // destination
		"text/plain",          // content-type
		[]byte(input.Message)) // body

	if err != nil {
		output.Result = "ERR_STOMP_SEND"
		return false, nil
	}

	output.Result = "OK"
	return true, nil
}

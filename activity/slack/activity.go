package slack

import (
	"net/http"

	"github.com/nlopes/slack"
	"github.com/project-flogo/core/activity"
)

func init() {
	_ = activity.Register(&Activity{})
}

var savedCookies []http.Cookie

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity is an activity that is used to invoke a REST Operation
// input    : {token, channel, message}
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

	ctx.Logger().Infof("Sending message on channel: %s", input.Channel)

	output := &Output{}

	//Validate the inputs
	if input.Token == "" {
		output.Result = "ERR_TOKEN_NOT_DEFINED"
		return false, nil
	}

	if input.Channel == "" {
		output.Result = "ERR_CHANNEL_NOT_DEFINED"
		return false, nil
	}

	//Send the message
	api := slack.New(input.Token)
	channelID, timestamp, err := api.PostMessage(input.Channel, slack.MsgOptionText(input.Message, false))
	if err != nil {
		ctx.Logger().Errorf("Error sending slack message: %v", err)
		output.Result = err.Error()
		return false, err
	}
	ctx.Logger().Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
	output.Result = "OK"
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

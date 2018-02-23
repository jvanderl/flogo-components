package slack

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/nlopes/slack"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-slack")

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

// Eval implements activity.Activity.Eval - Sends a message to TIBCO eFTL
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the activity data from the context
	wsToken, _ := context.GetInput("token").(string)
	wsChannel, _ := context.GetInput("channel").(string)
	wsMessage, _ := context.GetInput("message").(string)

	//Validate the inputs
	if wsToken == "" {
		context.SetOutput("result", "ERR_TOKEN_NOT_DEFINED")
		return false, nil
	}

	if wsChannel == "" {
		context.SetOutput("result", "ERR_CHANNEL_NOT_DEFINED")
		return false, nil
	}

	// Connect to Slack
	api := slack.New(wsToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			log.Debugf("Event Received. Type: %v", msg.Type)
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				log.Debugf("Connected succesfully, count %v", ev.ConnectionCount)
				// Connected, now find the channel ID
				channel := findChannelID(api, wsChannel)
				if channel == "" {
					context.SetOutput("result", "ERR_CHANNEL_NOT_FOUND")
					return false, nil
				}
				log.Debugf("About to send message to channel: %v", channel)
				rtm.SendMessage(rtm.NewOutgoingMessage(wsMessage, channel))
				//break Loop
			case *slack.AckMessage:
				break Loop
			case *slack.MessageEvent:
				log.Debugf("Message Received: %v\n", ev)

			case *slack.RTMError:
				log.Debugf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Debug("Invalid credentials")
				context.SetOutput("result", "ERR_LOGIN_FAILED")
				return false, nil
			default:
				//Take no action
			}
		}

	}
	context.SetOutput("result", "OK")
	return true, nil
}

func findChannelID(api *slack.Client, name string) string {
	log.Debugf("Looking up channel %v", name)
	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	for _, channel := range channels {
		log.Debugf("Checking Channelname: %v", channel.Name)
		if channel.Name == name {
			log.Debug("Found what I'm looking for!")
			return channel.ID
		}
	}
	return ""
}

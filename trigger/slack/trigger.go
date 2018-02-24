package slack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/nlopes/slack"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-slack")

// slackTrigger is a stub for your Trigger implementation
type slackTrigger struct {
	metadata          *trigger.Metadata
	runner            action.Runner
	config            *trigger.Config
	channelToActionID map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &slackFactory{metadata: md}
}

// slackFactory Trigger factory
type slackFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *slackFactory) New(config *trigger.Config) trigger.Trigger {
	slackTrigger := &slackTrigger{metadata: t.metadata, config: config}
	return slackTrigger
}

// Metadata implements trigger.Trigger.Metadata
func (t *slackTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *slackTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements trigger.Trigger.Start
func (t *slackTrigger) Start() error {

	// start the trigger
	wsToken := t.config.GetSetting("token")

	// Read Actions from trigger endpoints
	t.channelToActionID = make(map[string]string)

	for _, handlerCfg := range t.config.Handlers {
		log.Debugf("handlers: [%s]", handlerCfg.ActionId)
		epdestination := handlerCfg.GetSetting("channel") + "_" + handlerCfg.GetSetting("matchtext")
		log.Debugf("destination: [%s]", epdestination)
		t.channelToActionID[epdestination] = handlerCfg.ActionId
		nobots, err := strconv.ParseBool(handlerCfg.GetSetting("nobots"))
		log.Debugf("nobots: [%v]", nobots)
		if err != nil {
			return err
		}

	}

	// connect to Slack
	api := slack.New(wsToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				//connected
				log.Debugf("Connection counter: %v", ev.ConnectionCount)
			case *slack.MessageEvent:
				//incoming message
				log.Debugf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)
				log.Debugf("info: %v", info)
				log.Debugf("prefix: %v", prefix)
				message := ev.Text
				channel := ev.Channel
				username := ev.User
				isBot := false
				channelInfo, err := api.GetChannelInfo(ev.Channel)
				if err == nil {
					channel = channelInfo.NameNormalized
				}
				userInfo, err := api.GetUserInfo(ev.User)
				if err == nil {
					username = userInfo.Name
					isBot = userInfo.IsBot
					log.Debugf("Is this a bot? : %v", isBot)
				}
				for _, handler := range t.config.Handlers {
					destChannel := handler.GetSetting("channel")
					destMatch := handler.GetSetting("matchtext")
					if destChannel == "*" || destChannel != "*" && channel == destChannel {
						//Channel matches, now check text
						if destMatch == "*" || destMatch != "*" && strings.Contains(strings.ToUpper(message), strings.ToUpper(destMatch)) {
							//Text matches, Run Action
							destination := destChannel + "_" + destMatch
							actionId, found := t.channelToActionID[destination]
							if found {
								//now check if we need to skip bots
								nobots, _ := strconv.ParseBool(handler.GetSetting("nobots"))
								if isBot && nobots {
									log.Debugf("Skipping Bot Message")
								} else {
									log.Debugf("About to run action for Id [%s]", actionId)
									response := t.RunAction(actionId, handler, message, channel, username)
									if response != "" {
										rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
									}
								}
							} else {
								log.Debug("actionId not found")
							}
						}
					}

				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}

	return nil
}

// Stop implements trigger.Trigger.Start
func (t *slackTrigger) Stop() error {
	// stop the trigger
	return nil
}

// RunAction starts a new Process Instance
func (t *slackTrigger) RunAction(actionID string, handlerCfg *trigger.HandlerConfig, message string, channel string, username string) string {
	log.Debug("Starting new Process Instance")
	log.Debugf("Message: %s", message)
	log.Debugf("Channel: %s", channel)
	log.Debugf("Username: %s", username)
	log.Debugf("ActionID: %v", actionID)
	log.Debugf("handlerCfg: %v", handlerCfg)

	req := t.constructStartRequest(message, channel, username)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	act := action.Get(actionID)

	ctx := trigger.NewInitialContext(startAttrs, handlerCfg)

	results, err := t.runner.RunAction(ctx, act, nil)

	var replyData interface{}
	var replyCode int

	if len(results) != 0 {
		dataAttr, ok := results["data"]
		if ok {
			replyData = dataAttr.Value()
			log.Debugf("Got Reply Data: %v", replyData)
		}
		codeAttr, ok := results["code"]
		if ok {
			replyCode, _ = data.CoerceToInteger(codeAttr.Value())
			log.Debugf("Got Reply Code: %v", replyCode)
		}
	}

	if err != nil {
		log.Debugf("Slack Trigger Error: %s", err.Error())
		return ""
	}

	msgText := ""
	if replyData != nil {
		for s, a := range replyData.(map[string]interface{}) {
			if s == "text" {
				msgText = a.(string)
			}
		}
		return msgText
	}

	return ""
}

func (t *slackTrigger) constructStartRequest(message string, channel string, username string) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	data["channel"] = channel
	data["user"] = username
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}

func convert(b []byte) string {
	n := len(b)
	return string(b[:n])
}

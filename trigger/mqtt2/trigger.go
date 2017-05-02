package mqtt2

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-mqtt2")

// todo: switch to use endpoint registration

// Mqtt2Trigger is simple MQTT trigger
type Mqtt2Trigger struct {
	metadata          *trigger.Metadata
	runner            action.Runner
	client            mqtt.Client
	config            *trigger.Config
	topicToActionId  map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MQTT2Factory{metadata: md}
}

// MQTT2Factory Trigger factory
type MQTT2Factory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *MQTT2Factory) New(config *trigger.Config) trigger.Trigger {
	return &Mqtt2Trigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *Mqtt2Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *Mqtt2Trigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *Mqtt2Trigger) Start() error {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(t.config.GetSetting("broker"))
	opts.SetClientID(t.config.GetSetting("id"))
	opts.SetUsername(t.config.GetSetting("user"))
	opts.SetPassword(t.config.GetSetting("password"))
	b, err := strconv.ParseBool(t.config.GetSetting("cleansess"))
	if err != nil {
		log.Error("Error converting \"cleansess\" to a boolean ", err.Error())
		return err
	}
	opts.SetCleanSession(b)
	if storeType := t.config.Settings["store"]; storeType != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(t.config.GetSetting("store")))
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {

		topic := msg.Topic()
		//TODO we should handle other types, since mqtt message format are data-agnostic
		payload := string(msg.Payload())
		log.Infof("Received msg: %s", payload)
		log.Infof("Actual topic: %s", topic)

		// try topic without wildcards
		actionId, found := t.topicToActionId[topic]

		if found {
			t.RunAction(actionId, payload, topic)
		} else {
			// search for wildcards

			for _, handlerCfg := range t.config.Handlers {
				eptopic := handlerCfg.GetSetting("topic")
				if strings.HasSuffix(eptopic, "/#") {
					// is wildcard, now check actual topic starts with wildcard
					if strings.HasPrefix(topic, strings.TrimSuffix(eptopic, "/#")) {
						// Got a match, now get the action for the wildcard topic
						//actionType, found := t.topicToActionType[eptopic]
						actionId, found := t.topicToActionId[eptopic]
						if found {
							t.RunAction(actionId, payload, topic)
						}
					}
				}
			}
		}

	})

	client := mqtt.NewClient(opts)
	t.client = client
	log.Infof("Connecting to broker [%s]", t.config.GetSetting("broker"))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Info( "Connected to broker")

	i, err := strconv.Atoi(t.config.GetSetting("qos"))
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return err
	}

//	t.topicToActionType = make(map[string]string)
	t.topicToActionId = make(map[string]string)

	for _, handlerCfg := range t.config.Handlers {
		if token := t.client.Subscribe(handlerCfg.GetSetting("topic"), byte(i), nil); token.Wait() && token.Error() != nil {
			log.Errorf("Error subscribing to topic %s: %s", handlerCfg.Settings["topic"], token.Error())
			panic(token.Error())
		} else {
			log.Infof("Subscribed to topic %s for action %s", handlerCfg.GetSetting("topic"), handlerCfg.ActionId)
			t.topicToActionId[handlerCfg.GetSetting("topic")] = handlerCfg.ActionId
		}
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *Mqtt2Trigger) Stop() error {
	//unsubscribe from topics
	for _, handlerCfg := range t.config.Handlers {
		log.Infof("Unsubcribing from topic: %s ", handlerCfg.GetSetting("topic"))
		if token := t.client.Unsubscribe(handlerCfg.GetSetting("topic")); token.Wait() && token.Error() != nil {
			log.Errorf("Error unsubscribing from topic %s: %s", handlerCfg.GetSetting("topic"), token.Error())
		}
	}

	t.client.Disconnect(250)

	return nil
}

// RunAction starts a new Process Instance
func (t *Mqtt2Trigger) RunAction(actionId string, payload string, topic string) {

	log.Info("Starting new Process Instance")
	log.Infof("Action Id: %s", actionId)
	log.Infof("Payload: %s", payload)
	log.Infof("Actual Topic: %s", topic)

	req := t.constructStartRequest(payload, topic)
	//err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	//if err != nil {
	//	//http.Error(w, err.Error(), http.StatusBadRequest)
	//	log.Error("Error Starting action ", err.Error())
	//	return
	//}

	//todo handle error
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)
	action := action.Get(actionId)
	context := trigger.NewContext(context.Background(), startAttrs)

	//todo handle error
	_, replyData, err := t.runner.Run(context, action, actionId, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s]", actionId)

	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			log.Error(err)
		} else {
			t.publishMessage(req.ReplyTo, string(data))
		}
	}
}

func (t *Mqtt2Trigger) publishMessage(topic string, message string) {

	log.Debug("ReplyTo topic: ", topic)
	log.Debug("Publishing message: ", message)

	qos, err := strconv.Atoi(t.config.GetSetting("qos"))
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}
	token := t.client.Publish(topic, byte(qos), false, message)
	token.Wait()
}

func (t *Mqtt2Trigger) constructStartRequest(message string, topic string) *StartRequest {

	log.Debug("Received contstruct start request")

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	data["actualtopic"] = topic
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

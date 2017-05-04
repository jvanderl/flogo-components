package eftl

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/go-eftl"
	"strconv"
	//	"encoding/json"
	//	"strings"
)

//var dat map[string]interface{}

// log is the default package logger

var log = logger.GetLogger("trigger-jvanderl-eftl")

// eftlTrigger is a stub for your Trigger implementation
type eftlTrigger struct {
	metadata              *trigger.Metadata
	runner                action.Runner
	config                *trigger.Config
	destinationToActionId map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &eftlFactory{metadata: md}
}

// eftlFactory Trigger factory
type eftlFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *eftlFactory) New(config *trigger.Config) trigger.Trigger {
	return &eftlTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *eftlTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *eftlTrigger) Init(runner action.Runner) {
	t.runner = runner
}

/*//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &eftlFactory{metadata: md}
}

// eftlFactory Trigger factory
type eftlFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *eftlFactory) New(config *trigger.Config) trigger.Trigger {
	return &eftlTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *eftlTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements trigger.Trigger.Init
func (t *eftlTrigger) Init(runner action.Runner) {
	t.runner = runner
}
*/
// Start implements trigger.Trigger.Start
func (t *eftlTrigger) Start() error {

	// start the trigger
	wsHost := t.config.GetSetting("server")
	wsChannel := t.config.GetSetting("channel")
	wsUser := t.config.GetSetting("user")
	wsPassword := t.config.GetSetting("password")
	wsSecure, err := strconv.ParseBool(t.config.GetSetting("secure"))
	if err != nil {
		return err
	}
	wsCert := "DummyCert"
	if wsSecure {
		wsCert = t.config.GetSetting("certificate")
	}

	// Read Actions from trigger endpoints
	t.destinationToActionId = make(map[string]string)

	for _, handlerCfg := range t.config.Handlers {
		log.Debugf("handlers: [%s]", handlerCfg.ActionId)
		epdestination := handlerCfg.GetSetting("destination")
		log.Debugf("destination: [%s]", epdestination)
		t.destinationToActionId[epdestination] = handlerCfg.ActionId
	}

	// Connect to eFTL server
	log.Infof("Connecting to eFTL server: [%s]", wsHost)
	eftlConn, err := eftl.Connect(wsHost, wsChannel, wsSecure, wsCert, "")
	if err != nil {
		log.Debugf("Error while connecting to wsHost: [%s]", err)
		return err
	}

	// Login to eFTL
	err = eftlConn.Login(wsUser, wsPassword)
	if err != nil {
		log.Debugf("Error while Loggin in: [%s]", err)
	}
	log.Debugf("Login succesful. client_id: [%s], id_token: [%s]", eftlConn.ClientID, eftlConn.ReconnectToken)

	//Subscribe to destination in endpoints
	for _, handler := range t.config.Handlers {
		log.Infof("Subscribing to destination: [%s]", handler.GetSetting("destination"))
		destination := "{\"_dest\":\"" + handler.GetSetting("destination") + "\"}"
		wsSubscriptionID, err := eftlConn.Subscribe(destination, "")
		if err != nil {
			log.Debugf("Error while subscribing in: [%s]", err)
		}
		log.Debugf("Subscribe succesful. subscription_id: [%s]", wsSubscriptionID)
	}

	for {
		message, destination, err := eftlConn.ReceiveMessage()
		if err != nil {
			return err
		}
		//actionType, found := t.destinationToActionType[destination]
		actionId, found := t.destinationToActionId[destination]
		if found {
			log.Debugf("About to run action for Id [%s]", actionId)
			t.RunAction(actionId, message, destination)
		} else {
			log.Debug("actionId not found")
		}
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *eftlTrigger) Stop() error {
	// stop the trigger
	return nil
}

// RunAction starts a new Process Instance
func (t *eftlTrigger) RunAction(actionId string, payload string, destination string) {

	log.Debug("Starting new Process Instance")
	log.Debugf("Action Id: ", actionId)
	log.Debugf("Payload: ", payload)
	log.Debugf("Destination: ", destination)

	req := t.constructStartRequest(payload)
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)
	action := action.Get(actionId)
	context := trigger.NewContext(context.Background(), startAttrs)
	_, replyData, err := t.runner.Run(context, action, actionId, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debug("Reply data: ", replyData)

	/*	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			log.Error(err)
		} else {
			t.publishMessage(req.ReplyTo, partition, string(data))
		}
	}*/
}

//func (t *eftlTrigger) constructStartRequest(message string, destination string) *StartRequest {
func (t *eftlTrigger) constructStartRequest(message string) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	//	data["destination"] = destination
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

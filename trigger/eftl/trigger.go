package eftl

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/go-eftl"
	"strconv"
)

var dat map[string]interface{}

// log is the default package logger

var log = logger.GetLogger("trigger-jvanderl-eftl")

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata                *trigger.Metadata
	runner                  action.Runner
	config                  *trigger.Config
	destinationToActionURI  map[string]string
	destinationToActionType map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyFactory{metadata: md}
}

// MyFactory Trigger factory
type MyFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *MyFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {

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
	t.destinationToActionType = make(map[string]string)
	t.destinationToActionURI = make(map[string]string)


	for _, handler := range t.config.Handlers {
		log.Infof("handlers: [%s]", handler.ActionId)
		t.destinationToActionURI[handler.GetSetting("destination")] = handler.ActionId
		t.destinationToActionType[handler.GetSetting("destination")] = "flow"
	}

	// Connect to eFTL server
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
		actionType, found := t.destinationToActionType[destination]
		actionURI, _ := t.destinationToActionURI[destination]
		if found {
			log.Infof ("About to run action for type [%s], uri [%s]", actionType, actionURI)
			t.RunAction(actionType, actionURI, message, destination)
		} else {
			log.Debug("actionType and URI not found")
		}
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}

// RunAction starts a new Process Instance
func (t *MyTrigger) RunAction(actionType string, actionURI string, payload string, destination string) {

	log.Info("Starting new Process Instance")
	log.Info("Action Type: ", actionType)
	log.Info("Action URI: ", actionURI)
	log.Info("Payload: ", payload)
	log.Info("Destination: ", destination)

	log.Info("Construct Request")
	req := t.constructStartRequest(payload, destination)
	log.Infof("Request data: [%s]", req.Data)

	log.Info("Set Start Attributes")
//	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	log.Info("Get Action to perform")
	action := action.Get(actionURI)

	log.Info("Set Trigger context")
	context := trigger.NewContext(context.Background(), startAttrs)

	log.Info("Call the runner")
	_, replyData, err := t.runner.Run(context, action, actionURI, nil)
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

func (t *MyTrigger) constructStartRequest(message string, destination string) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
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

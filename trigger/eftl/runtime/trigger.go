package eftl

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/jvanderl/go-eftl"
	"github.com/op/go-logging"
	"strconv"
)

var dat map[string]interface{}

// log is the default package logger

var log = logging.MustGetLogger("trigger-eftl")

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata                *trigger.Metadata
	runner                  action.Runner
	settings                map[string]string
	config                  *trigger.Config
	destinationToActionURI  map[string]string
	destinationToActionType map[string]string
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MyTrigger{metadata: md})
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(config *trigger.Config, runner action.Runner) {
	t.config = config
	t.settings = config.Settings
	t.runner = runner
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	// start the trigger
	wsHost := t.settings["server"]
	wsChannel := t.settings["channel"]
	wsUser := t.settings["user"]
	wsPassword := t.settings["password"]
	wsSecure, err := strconv.ParseBool(t.settings["secure"])
		if err != nil {
			return err
		}
	wsCert := t.settings["certificate"]

	// Read Actions from trigger endpoints
	t.destinationToActionType = make(map[string]string)
	t.destinationToActionURI = make(map[string]string)

	for _, endpoint := range t.config.Endpoints {
		t.destinationToActionURI[endpoint.Settings["destination"]] = endpoint.ActionURI
		t.destinationToActionType[endpoint.Settings["destination"]] = endpoint.ActionType
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
	for _, endpoint := range t.config.Endpoints {
		destination := "{\"_dest\":\"" + endpoint.Settings["destination"] + "\"}"
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

	log.Debug("Starting new Process Instance")
	log.Debug("Action Type: ", actionType)
	log.Debug("Action URI: ", actionURI)
	log.Debug("Payload: ", payload)
	log.Debug("Destination: ", destination)

	req := t.constructStartRequest(payload, destination)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionType)

	context := trigger.NewContext(context.Background(), startAttrs)

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

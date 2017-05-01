package eftl

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/go-eftl"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-eftl")

type eftlLoginMessage struct {
	Operator      int               `json:"op"`
	ClientType    string            `json:"client_type"`
	ClientVersion string            `json:"client_version"`
	User          string            `json:"user"`
	Password      string            `json:"password"`
	LoginOptions  map[string]string `json:"login_options"`
}

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
	wsHost := context.GetInput("server").(string)
	wsChannel := context.GetInput("channel").(string)
	wsDestination := context.GetInput("destination").(string)
	wsMessage := context.GetInput("message").(string)
	wsUser := context.GetInput("user").(string)
	wsPassword := context.GetInput("password").(string)
	wsSecure := context.GetInput("secure").(bool)
	wsCert := context.GetInput("certificate").(string)

	// Connect to eFTL server
	eftlConn, err := eftl.Connect(wsHost, wsChannel, wsSecure, wsCert, "")
	if err != nil {
		log.Debugf("Error while connecting to wsHost: [%s]", err)
		context.SetOutput("result", "ERR_CONNECT_HOST")
		return false, err
	}

	// Login to eFTL
	err = eftlConn.Login(wsUser, wsPassword)
	if err != nil {
		log.Debugf("Error while Loggin in: [%s]", err)
		context.SetOutput("result", "ERR_EFTL_LOGIN")
		return false, err
	}
	log.Debugf("Login succesful. client_id: [%s], id_token: [%s]", eftlConn.ClientID, eftlConn.ReconnectToken)

	// Send the message
	err = eftlConn.SendMessage(wsMessage, wsDestination)
	if err != nil {
		log.Debugf("Error while sending message to wsHost: [%s]", err)
		context.SetOutput("result", "ERR_SEND_MESSAGE")
		return false, err
	}

	context.SetOutput("result", "OK")

	return true, nil
}

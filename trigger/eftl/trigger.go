package eftl

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/tib-eftl"
	"strconv"
	"crypto/x509"
	"crypto/tls"
	"encoding/base64"
	"net/url"
	"fmt"
)

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
	eftlTrigger := &eftlTrigger{metadata: t.metadata, config: config}
	return eftlTrigger
}

// Metadata implements trigger.Trigger.Metadata
func (t *eftlTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *eftlTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements trigger.Trigger.Start
func (t *eftlTrigger) Start() error {

	// start the trigger
	wsHost := t.config.GetSetting("server")
	wsClientID := t.config.GetSetting("clientid")
	wsChannel := t.config.GetSetting("channel")
	wsUser := t.config.GetSetting("user")
	wsPassword := t.config.GetSetting("password")
	wsSecure, err := strconv.ParseBool(t.config.GetSetting("secure"))
	if err != nil {
		return err
	}
	wsCert := ""
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

	wsURL := url.URL{}
	if wsSecure {
		wsURL = url.URL{Scheme: "wss", Host: wsHost, Path: wsChannel}
	} else {
		wsURL = url.URL{Scheme: "ws", Host: wsHost, Path: wsChannel}
	}
	wsConn := wsURL.String()

	var tlsConfig *tls.Config

	if wsCert != "" {
		// TLS configuration uses CA certificate from a PEM file to
		// authenticate the server certificate when using wss:// for
		// a secure connection
		caCert, err := base64.StdEncoding.DecodeString(wsCert)
		if err != nil {
			log.Errorf("unable to decode certificate: %s", err)
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	} else {
		// TLS configuration accepts all server certificates
		// when using wss:// for a secure connection
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// channel for receiving connection errors
	errChan := make(chan error, 1)

	// set connection options
	opts := &eftl.Options{
		ClientID:  wsClientID,
		Username:  wsUser,
		Password:  wsPassword,
		TLSConfig: tlsConfig,
	}

	// connect to the server
	conn, err := eftl.Connect(wsConn, opts, errChan)
	if err != nil {
		log.Errorf("Error connecing to eFTL server: [%s]", err)
		return err
	}

	// close the connection when done
	defer conn.Disconnect()

	// channel for receiving subscription response
	subChan := make(chan *eftl.Subscription, 1)

	// channel for receiving published messages
	msgChan := make(chan eftl.Message, 1000)


	//Subscribe to destination in endpoints
	for _, handler := range t.config.Handlers {
		log.Infof("Subscribing to destination: [%s]", handler.GetSetting("destination"))
		// create the message content matcher
		matcher := fmt.Sprintf("{\"_dest\":\"%s\"}", handler.GetSetting("destination"))
		durablename := ""
		durable, err := strconv.ParseBool(handler.GetSetting("durable"))
		if err != nil {
			return err
		}
		if durable {
			durablename = handler.GetSetting("durablename")
		}
		conn.SubscribeAsync(matcher, durablename, msgChan, subChan)
	}

	for {
		select {
		case sub := <-subChan:
			if sub.Error != nil {
				log.Infof("subscribe operation failed: %s", sub.Error)
				return sub.Error
			}
			log.Infof("subscribed with matcher %s", sub.Matcher)
		case msg := <-msgChan:
			log.Infof("received message: %s", msg)
			destination := msg["_dest"].(string)
			message := msg["text"].(string)
			actionId, found := t.destinationToActionId[destination]
			if found {
				log.Debugf("About to run action for Id [%s]", actionId)
				t.RunAction(actionId, message, destination)
			} else {
				log.Debug("actionId not found")
			}

		case err := <-errChan:
			log.Infof("connection error: %s", err)
			return err
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
	log.Debugf("Action Id: %s", actionId)
	log.Debugf("Payload: %s", payload)

	req := t.constructStartRequest(payload)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionId)

	context := trigger.NewContext(context.Background(), startAttrs)

	_, replyData, err := t.runner.Run(context, action, actionId, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s]", actionId)
	log.Debugf("Reply data: [%v]", replyData)

}

func (t *eftlTrigger) constructStartRequest(message string) *StartRequest {

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
	ReplyTo     string                 `json:"replyTo"`
}

func convert(b []byte) string {
	n := len(b)
	return string(b[:n])
}

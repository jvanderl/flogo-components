package eftl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/tib-eftl"
	//	"strings"
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

	// Create array of channels for all handlers
	msgChans := make([]chan eftl.Message, len(t.config.Handlers))

	// Error channel for receiving connection errors
	errChan := make(chan error, 1)

	//Create the subsription channel [1]
	subChan := make(chan *eftl.Subscription, 1)

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

	//Subscribe to destination in endpoints
	for i, handler := range t.config.Handlers {
		msgChans[i] = make(chan eftl.Message)
		log.Infof("Subscribing to messages [%v]: [%s]", i, handler.GetSetting("matcher"))
		// create the message content matcher
		//complex matcher format like '{"_dest":"subject"}' can be used directly
		matcher := handler.GetSetting("matcher")
		if string(matcher[0:1]) != "{" {
			// simple destination, will need to form matcher
			matcher = fmt.Sprintf("{\"_dest\":\"%s\"}", handler.GetSetting("matcher"))
		}

		durablename := ""
		durable, err := strconv.ParseBool(handler.GetSetting("durable"))
		if err != nil {
			return err
		}
		if durable {
			durablename = handler.GetSetting("durablename")
		}
		log.Infof("created matcher: %v", matcher)
		conn.SubscribeAsync(matcher, durablename, msgChans[i], subChan)
	}

	for {
		select {
		case sub := <-subChan:
			if sub.Error != nil {
				log.Infof("subscribe operation failed: %s", sub.Error)
				return sub.Error
			}
			log.Infof("subscribed with matcher %s", sub.Matcher)

		case err := <-errChan:
			log.Infof("connection error: %s", err)
			return err
		}
		cases := make([]reflect.SelectCase, len(msgChans))
		for i, ch := range msgChans {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}
		remaining := len(cases)
		for remaining > 0 {
			chosen, value, ok := reflect.Select(cases)
			fmt.Printf("Read from channel %#v and received %s\n", chosen, value)
			log.Infof("received message: %s", value)
			if !ok {
				// The chosen channel has been closed, so zero out the channel to disable the case
				cases[chosen].Chan = reflect.ValueOf(nil)
				remaining--
				continue
			}
			//Get eFTL message from value
			msg, ok := value.Interface().(eftl.Message)
			if !ok {
				log.Error("Error casting regular message type")
				continue
			}
			/*			message := msg["text"].(string)
						log.Infof("Message Payload: %v", message)
						destination := msg["_dest"].(string)
						log.Infof("Message Destination: %v", destination)
						subject := msg["_subj"].(string)
						log.Infof("Message Subject: %v", subject)
						//actionId := t.config.Handlers[chosen-2].ActionId */
			actionId := t.config.Handlers[chosen].ActionId
			log.Debugf("About to run action for Id [%s]", actionId)
			//t.RunAction(actionId, message, destination, subject)
			t.RunAction(actionId, msg)
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
func (t *eftlTrigger) RunAction(actionId string, msg eftl.Message) {
	log.Debug("Starting new Process Instance")
	log.Debugf("Action Id: %s", actionId)

	req := t.constructStartRequest(msg)

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

func (t *eftlTrigger) constructStartRequest(message eftl.Message) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI string                 `json:"flowUri"`
	Data       map[string]interface{} `json:"data"`
	ReplyTo    string                 `json:"replyTo"`
}

func convert(b []byte) string {
	n := len(b)
	return string(b[:n])
}

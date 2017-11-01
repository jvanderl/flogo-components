package tcm

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
	"fmt"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-tcm")

// tcmTrigger is a stub for your Trigger implementation
type tcmTrigger struct {
	metadata              *trigger.Metadata
	runner                action.Runner
	config                *trigger.Config
	destinationToActionId map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &tcmFactory{metadata: md}
}

// tcmFactory Trigger factory
type tcmFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *tcmFactory) New(config *trigger.Config) trigger.Trigger {
	tcmTrigger := &tcmTrigger{metadata: t.metadata, config: config}
	return tcmTrigger
}

// Metadata implements trigger.Trigger.Metadata
func (t *tcmTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *tcmTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements trigger.Trigger.Start
func (t *tcmTrigger) Start() error {

	wsCert := ""

	// start the trigger
	wsURL := t.config.GetSetting("url")
	wsAuthKey := t.config.GetSetting("authkey")
	wsClientID := t.config.GetSetting("clientid")
  wsCert = t.config.GetSetting("certificate")

	// Read Actions from trigger endpoints
	t.destinationToActionId = make(map[string]string)

	for _, handlerCfg := range t.config.Handlers {
		log.Debugf("handlers: [%s]", handlerCfg.ActionId)
		epdestination := handlerCfg.GetSetting("destinationname") + "_" + handlerCfg.GetSetting("destinationmatch")
		log.Debugf("destination: [%s]", epdestination)
		t.destinationToActionId[epdestination] = handlerCfg.ActionId
	}

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
		Username:  "",
		Password:  wsAuthKey,
		TLSConfig: tlsConfig,
	}

	// connect to the server
	conn, err := eftl.Connect(wsURL, opts, errChan)
	if err != nil {
		log.Errorf("Error connecing to TCM server: [%s]", err)
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
		destName := handler.GetSetting("destinationname")
		destMatch := handler.GetSetting("destinationmatch")
		// create the message content matcher
		matcher := ""
		if (destMatch == "*") {
			matcher = fmt.Sprintf("{\"%s\":true}", destName)
		} else {
			matcher = fmt.Sprintf("{\"%s\":\"%s\"}", destName, destMatch)
		}
		durable, err := strconv.ParseBool(handler.GetSetting("durable"))
		if err != nil {
			return err
		}
		if durable {
			durablename := handler.GetSetting("durablename")
			log.Infof("Subscribing to destination: %s:%s, durable name:%s", destName, destMatch, durablename)
			conn.SubscribeAsync(matcher, durablename, msgChan, subChan)
		} else {
			log.Infof("Subscribing to destination: %s:%s", destName, destMatch)
			conn.SubscribeAsync(matcher, "", msgChan, subChan)
		}
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
			// see if we can find a matching handler
			for _, handler := range t.config.Handlers {
				destName := handler.GetSetting("destinationname")
				destMatch := handler.GetSetting("destinationmatch")
				msgName := handler.GetSetting("messagename")
				if ((msg[destName].(string) == destMatch) || (msg[destName] != nil && destMatch == "*")) {
					destination := destName + "_" + destMatch
					message := msg[msgName].(string)
					actionId, found := t.destinationToActionId[destination]
					if found {
						log.Debugf("About to run action for Id [%s]", actionId)
						t.RunAction(actionId, message, destination)
					} else {
						log.Debug("actionId not found")
					}
				}
			}

		case err := <-errChan:
			log.Infof("connection error: %s", err)
			return err
		}
	}
	return nil
}


// Stop implements trigger.Trigger.Start
func (t *tcmTrigger) Stop() error {
	// stop the trigger
	return nil
}

// RunAction starts a new Process Instance
func (t *tcmTrigger) RunAction(actionId string, payload string, destination string) {
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

func (t *tcmTrigger) constructStartRequest(message string) *StartRequest {

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

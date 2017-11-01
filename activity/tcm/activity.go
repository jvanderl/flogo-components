package tcm

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/tib-eftl"
	"crypto/x509"
	"crypto/tls"
	"encoding/base64"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-tcm")

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
	wsURL, _ := context.GetInput("url").(string)
	wsAuthKey, _ := context.GetInput("authkey").(string)
	wsClientID, _ := context.GetInput("clientid").(string)
	wsDestinationName, _ := context.GetInput("destinationname").(string)
	wsDestinationValue, _ := context.GetInput("destinationvalue").(string)
	wsMessageName, _ := context.GetInput("messagename").(string)
	wsMessageValue, _ := context.GetInput("messagevalue").(string)
//	wsCert, _ := context.GetInput("certificate").(string)
	wsCert := ""

	var tlsConfig *tls.Config

	if wsCert != "" {
		// TLS configuration uses CA certificate from a PEM file to
		// authenticate the server certificate when using wss:// for
		// a secure connection
		caCert, err := base64.StdEncoding.DecodeString(wsCert)
		if err != nil {
			log.Errorf("unable to decode certificate: %s", err)
			return false, err
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
		context.SetOutput("result", "ERR_CONNECT_HOST")
		return false, err
	}

	// close the connection when done
	defer conn.Disconnect()

	// channel for receiving publish completions
	compChan := make(chan *eftl.Completion, 1000)

	// publish the message
	if (wsDestinationValue != ""){
		conn.PublishAsync(eftl.Message{
			wsDestinationName: wsDestinationValue,
			wsMessageName: wsMessageValue,
			}, compChan)
		} else {
			conn.PublishAsync(eftl.Message{
				wsMessageName: wsMessageValue,
				}, compChan)
		}

	for {
		select {
		case comp := <-compChan:
			if comp.Error != nil {
				log.Errorf("Error while sending message to wsHost: [%s]", comp.Error)
				context.SetOutput("result", "ERR_SEND_MESSAGE")
				return false, comp.Error
			}
			log.Debugf("published message: %s", comp.Message)
			context.SetOutput("result", "OK")
			return true, nil
		case err := <-errChan:
			log.Errorf("connection error: %s", err)
			context.SetOutput("result", "ERR_CONNECT_HOST")
			return false, err
		}
	}

}

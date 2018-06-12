package wssend

import (
	"net/url"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/gorilla/websocket"
)

// log is the default package logger
var log = logger.GetLogger("activity-sendWSMessage")

// WsMsgActivity is a stub for your Activity implementation
type WsMsgActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new WsMsgActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &WsMsgActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *WsMsgActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval - Sends a message to a WebSocket enabled server like TIBCO eFTL
func (a *WsMsgActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the activity data from the context
	wsHost := context.GetInput("server").(string)
	wsChannel := context.GetInput("channel").(string)
	wsMessage := context.GetInput("message").(string)
	wsWaitResp := context.GetInput("waitforresponse").(bool)
	wsTimeout, _ := context.GetInput("timeout").(int)

	wsURL := url.URL{Scheme: "ws", Host: wsHost, Path: wsChannel}
	log.Debugf("connecting to %s", wsURL.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Debugf("Error while dialing to wsHost: ", err)
		context.SetOutput("response", err)
		context.SetOutput("result", "ERR_CONN_HOST")
		return false, err
	}

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(wsMessage))
	if err != nil {
		log.Debugf("Error while sending message to wsHost: [%s]", err)
		context.SetOutput("response", err)
		context.SetOutput("result", "ERR_SEND_MESSAGE")
		return false, err
	}

	if wsWaitResp == true {
		wsConn.SetReadDeadline(time.Now().Add(time.Duration(wsTimeout) * time.Second))
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			context.SetOutput("response", err)
			context.SetOutput("result", "ERR_WAITING_RESPONSE")
		} else {
			message := string(msg)
			log.Debugf("message: %v", message)
			context.SetOutput("response", message)
			context.SetOutput("result", "OK")
		}
	} else {
		context.SetOutput("result", "OK")
	}
	wsConn.Close()

	return true, nil
}

package wssend

import (
	"net/url"

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

	wsURL := url.URL{Scheme: "ws", Host: wsHost, Path: wsChannel}
	log.Debugf("connecting to %s", wsURL.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Debugf("Error while dialing to wsHost: ", err)
	}

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(wsMessage))
	if err != nil {
		log.Debugf("Error while sending message to wsHost: [%s]", err)
		return
	}
	context.SetOutput("result", "OK")
	wsConn.Close()

	return true, nil
}

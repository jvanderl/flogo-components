package eftl

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"encoding/json"
	"net/url"
	"context"
)

var dat map[string]interface{}

// log is the default package logger

var log = logging.MustGetLogger("trigger-eftl")

// MyTrigger is a stub for your Trigger implementation

type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	settings          map[string]string
	config            *trigger.Config
	destinationToActionURI  map[string]string
	destinationToActionType map[string]string
}

type eftlGeneric struct {
    element map[string]string
}
// {"op": 1, "client_type": "js", "client_version": "3.1.0   V7", "user":"user", "password":"password", "login_options": {"_qos": "true"}}


//type eftlLoginOptions struct {
//	_qos *bool
//}
type eftlLoginMessage struct {
	Operator int  `json:"op"`
	ClientType string `json:"client_type"`
	client_version string `json:"client_version"`
	User string `json:"user"`
	Password string `json:"password"`
	LoginOptions map[string]string `json:"login_options"`
}

// {"op":2,"client_id":"68404BBF-8831-42CD-8744-43385AEF3590","id_token":"ZkGGc24IVFacA9jIGOBgb2bDc2g=","timeout":600,"heartbeat":240,"max_size":8192,"_qos":"true"}

type eftlLoginResponse struct {
	Operator int  `json:"op"`
	ClientID string `json:"client_id"`	
	IDToken string `json:"id_token"`
	Timeout int `json:"timeout"`
	Heartbeat int `json:"heartbeat"`
	MaxSize int `json:"max_size"`
	QoS string `json:"_qos"`
}

// {"op":3,"id":"68404BBF-8831-42CD-8744-43385AEF3590.s.1","matcher":"{\"_dest\":\"sample\"}"}

type eftlSubscription struct {
	Operator int  `json:"op"`
	ClientID string `json:"id"`	
	Matcher string `json:"matcher"`
}

// {"op":4,"id":"68404BBF-8831-42CD-8744-43385AEF3590.s.1"}

type eftlSubscriptionResponse struct {
	Operator int  `json:"op"`
	SubscriptionID string `json:"id"`	
}


// {"op":7,"to":"68404BBF-8831-42CD-8744-43385AEF3590.s.1","seq":1,"body":{"_dest":"sample","text":"This is a sample eFTL message","number":1}}

type eftlBody struct {
	Destination string `json:"_dest"`
	Text string `json:"text"`
	Number int `json:"number"`
}

type eftlMessage struct {
	Operator int  `json:"op"`
	To string `json:"to"`
	Sequence int `json:"seq"`
	Body eftlBody `json:"body"`
} 

// {"op":9,"seq":1}

type eftlSequenceMsg struct {
	Operator int  `json:"op"`
	Sequence int `json:"seq"`
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
//    wsClientID := t.settings["clientid"]
    wsDestination := t.settings["destination"]
	wsUser := t.settings["user"]
	wsPassword := t.settings["password"]
	wsClientID := ""
	wsIDToken := ""
	wsSubscriptionID := ""
	
	log.Debug("server:", wsHost)
	log.Debug("channel:", wsChannel)
	log.Debug("clientid:", wsClientID)
	log.Debug("destination:", wsDestination)
	log.Debug("user:", wsUser)
	log.Debug("password:", wsPassword)

// Read Actions from trigger endpoints 
	t.destinationToActionType = make(map[string]string)
	t.destinationToActionURI = make(map[string]string)

	for _, endpoint := range t.config.Endpoints {
		t.destinationToActionURI[endpoint.Settings["destination"]] = endpoint.ActionURI
		t.destinationToActionType[endpoint.Settings["destination"]] = endpoint.ActionType
	}
	
	wsURL := url.URL{Scheme: "ws", Host: wsHost, Path: wsChannel}
	log.Debugf("connecting to %s", wsURL.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Debugf("Error while dialing to wsHost: ", err)
	}

	loginMessage := eftlLoginMessage{1, "js", "3.1.0   V7", wsUser, wsPassword, map[string]string{"_qos": "true"}}

	loginb, err := json.Marshal(loginMessage)
	if err != nil {
		log.Debugf("Error while marshalling login message: [%s]", err)
		return err
	}
	
	log.Debug("Sending login message")

	err = wsConn.WriteMessage(websocket.TextMessage, loginb)
	if err != nil {
		log.Debugf("Error while sending login message to wsHost: [%s]", err)
		return err
	}

	for {
		messageType, p, err := wsConn.ReadMessage()
		if err == nil {
			switch messageType {
		    	case websocket.TextMessage : {
		    		if err := json.Unmarshal(p, &dat); err != nil {
     				   panic(err)
    				}
    				eftlOp := dat["op"].(float64)
					switch eftlOp {
						case 0 : { //Heartbeat
							// {"op":0}
							log.Debug("Heartbeat Received")
						}
						case 2 : { // login response
							// {"op":2,"client_id":"68404BBF-8831-42CD-8744-43385AEF3590","id_token":"ZkGGc24IVFacA9jIGOBgb2bDc2g=","timeout":600,"heartbeat":240,"max_size":8192,"_qos":"true"}
				    		res := new(eftlLoginResponse)
						    if err := json.Unmarshal(p, &res); err != nil {
		        				return err
		    				}
//							log.Debug("Login Response Received: [%s]", convert(p))
							wsClientID = res.ClientID
							wsIDToken = res.IDToken
							log.Debugf("Login Succesful. client_id: [%s], id_token: [%s]", wsClientID, wsIDToken)
//							log.Debug("client_id:", wsClientID)
//							log.Debug("id_token:", wsIDToken)
							// Got ClientID, now start subscribing
							// {"op":3,"id":"68404BBF-8831-42CD-8744-43385AEF3590.s.1","matcher":"{\"_dest\":\"sample\"}"}
							subscriptionMessage := eftlSubscription{3, wsClientID + ".s.1", "{\"_dest\":\"" + wsDestination + "\"}"}

							subscrb, err := json.Marshal(subscriptionMessage)
							if err != nil {
								log.Debugf("Error while marshalling subscription message: [%s]", err)
								return err
							}
							
//							log.Debugf("Preparing to send subscription message: [%s]", convert(subscrb))
							log.Debugf("Subscribing to destination: [%s]", wsDestination)
							err = wsConn.WriteMessage(websocket.TextMessage, subscrb)
							if err != nil {
								log.Debugf("Error while sending subscription message to wsHost: [%s]", err)
								return err
							} 
//							log.Debug("Subscription sent succesfully")
						}
						case 4 : { // subscription response
							// {"op":4,"id":"68404BBF-8831-42CD-8744-43385AEF3590.s.1"}
				    		res := eftlSubscriptionResponse{}
						    if err := json.Unmarshal(p, &res); err != nil {
		        				return err
		    				}
//							log.Debug("Subscription Response Received [%s]", convert(p))
							wsSubscriptionID = res.SubscriptionID
							log.Debugf("Subscription successful. Subscription ID: [%s]", wsSubscriptionID)
						}
						case 7 : { // Regular Message
							// {"op":7,"to":"68404BBF-8831-42CD-8744-43385AEF3590.s.1","seq":1,"body":{"_dest":"sample","text":"This is a sample eFTL message","number":1}}
				    		res := eftlMessage{}
						    if err := json.Unmarshal(p, &res); err != nil {
		        				return err
		    				}
							message := res.Body.Text
							log.Debugf("Regular Message Received: [%s]", message )
							log.Debugf("Destination: [%s]", res.Body.Destination)
							actionType, found := t.destinationToActionType[wsDestination]
							actionURI, _ := t.destinationToActionURI[wsDestination]
							if found {
								log.Debugf("Found actionType [%s]", actionType)
								log.Debugf("Found actionURI [%s]", actionURI)
								t.RunAction(actionType, actionURI, message, wsDestination)
							} else {
								log.Debug("actionType and URI not found")
							}
						}
						case 9 : { //Sequence Message
							// {"op":9,"seq":1}
				    		res := eftlSequenceMsg{}
						    if err := json.Unmarshal(p, &res); err != nil {
		        				return err
		    				}
							log.Debugf("Seqence Received: [%s]", res.Sequence)
						}
						default : {
							log.Debugf("Other message Received: [%s]", convert(p))
						}
					}
		    	}
		    	case websocket.BinaryMessage : {
		    		log.Debug("Received Binary message", p)
		    	}
		    	case websocket.CloseMessage : {
		    		log.Debug("Received Close message", p)
		    		return nil
		    	}
		    	case websocket.PingMessage : {
		    		log.Debug("Received Ping message", p)
		    	}
		    	case websocket.PongMessage : {
		    		log.Debug("Received Pong message", p)
		    	}
		    }
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

	log.Debugf("Ran action: [%s-%s]", actionType, actionURI)
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

	log.Debug("Received contstruct start request")

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

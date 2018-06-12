package wsserver

import (
	"context"
	"flag"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/gorilla/websocket"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-wsserver")

var upgrader = websocket.Upgrader{} // use default options

//WsServerTrigger is th main structure for this trigger
type WsServerTrigger struct {
	metadata *trigger.Metadata
	//runner   action.Runner
	config   *trigger.Config
	handlers []*trigger.Handler
	//timers map[string]*scheduler.Job
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &WsServerFactory{metadata: md}
}

// WsServerFactory Trigger factory
type WsServerFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *WsServerFactory) New(config *trigger.Config) trigger.Trigger {
	return &WsServerTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *WsServerTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *WsServerTrigger) Init(runner action.Runner) {
	log.Debug("Trigger Init called")
}

// Initialize implements ext.Trigger.Initialize
func (t *WsServerTrigger) Initialize(ctx trigger.InitContext) error {
	log.Debug("Trigger Initialize called")

	t.handlers = ctx.GetHandlers()

	return nil
}

func (t *WsServerTrigger) handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade: %v", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		log.Debugf("RequestURI: %v", r.RequestURI)
		if err != nil {
			log.Errorf("read: %v", err)
			break
		}
		log.Debugf("RequestURI: %s", r.RequestURI)
		log.Debugf("recv: %s", message)

		for _, handler := range t.handlers {
			channel := handler.GetStringSetting("channel")
			log.Debugf("Checking URI aginst handler channel: %s", channel)
			if r.RequestURI == channel {
				log.Debug("Found matching handler, starting flow")
				response := t.Execute(handler, string(message), channel)
				if response != nil {
					log.Debugf("Got response message: %v", response)
					responseMsg := []byte(response.(string))
					err = c.WriteMessage(mt, responseMsg)
					if err != nil {
						log.Errorf("Error sending response: %v", err)
					}
				}
			}
		}

	}
}

// Start implements ext.Trigger.Start
func (t *WsServerTrigger) Start() error {

	log.Debug("Trigger Start Called")

	port := t.config.GetSetting("port")
	var addr = flag.String("addr", "localhost:"+port, "http service address")

	flag.Parse()

	http.HandleFunc("/", t.handleWS)

	http.ListenAndServe(*addr, nil)

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *WsServerTrigger) Stop() error {

	return nil
}

// Execute executes any handlers defined immediately on startup
func (t *WsServerTrigger) Execute(handler *trigger.Handler, wsmessage string, wschannel string) interface{} {
	log.Debug("Starting process")

	triggerData := map[string]interface{}{
		"message": wsmessage,
		"channel": wschannel,
	}

	response, err := handler.Handle(context.Background(), triggerData)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
		return nil
	}
	log.Debugf("Action call successful: %v", response)
	return response["response"].Value()
}

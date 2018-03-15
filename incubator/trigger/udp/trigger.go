package udp

import (
	"context"
	"net"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-udp")

// udpTriggerFactory My Trigger factory
type udpTriggerFactory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &udpTriggerFactory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *udpTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &udpTrigger{metadata: t.metadata, config: config}
}

// udpTrigger is a stub for your Trigger implementation
type udpTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	//handlers []*trigger.Handler
}

// Init implements trigger.Trigger.Init
func (t *udpTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Metadata implements trigger.Trigger.Metadata
func (t *udpTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

const (
	maxDatagramSize = 8192
)

// Start implements trigger.Trigger.Start
func (t *udpTrigger) Start() error {
	// start the trigger
	log.Debug("Start")
	// Get parms
	wsPort := t.config.GetSetting("port")
	wsGroup := t.config.GetSetting("multicast_group")
	var l *net.UDPConn

	//  for multicastgroup
	log.Debug("Resolve")

	addr, err := net.ResolveUDPAddr("udp", wsGroup+":"+wsPort)

	if err != nil {
		// Handle error
		log.Errorf("Error resolving address %v", err)
	}

	log.Debug("Resolved Addr O/P %v", addr)

	if wsGroup == "" {
		/* Now listen at selected port */
		log.Debug("ListenUDP")
		l, err = net.ListenUDP("udp", addr)
	} else {
		log.Debug("ListenMulticastUDP")
		l, err = net.ListenMulticastUDP("udp", nil, addr)
	}
	if err != nil {
		log.Errorf("ListenUDP failed: %v", err)
	}

	// common
	log.Debug("SetRead")
	l.SetReadBuffer(maxDatagramSize)

	for {
		buf := make([]byte, maxDatagramSize)

		log.Debug("BeforeREAD")
		n, addr, err := l.ReadFromUDP(buf)

		// Read ok ?
		if err != nil {
			log.Errorf("ReadFromUDP failed: %v", err)
		}

		log.Debug("afterRead ")
		payload := string(buf[0:n])

		log.Infof("Received %v from %v", payload, addr)

		handlers := t.config.Handlers

		log.Debug("Processing handlers")
		for _, handler := range handlers {
			t.RunAction(handler, payload)
		}

	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *udpTrigger) Stop() error {
	// stop the trigger
	return nil
}

func (t *udpTrigger) RunAction(handlerCfg *trigger.HandlerConfig, payload string) {

	log.Debug("Starting Payload data action")

	req := t.constructStartRequest(handlerCfg, payload)
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)
	action := action.Get(handlerCfg.ActionId)
	context := trigger.NewContext(context.Background(), startAttrs)

	log.Debugf("ActionID: '%s'", handlerCfg.ActionId)

	// Run the triggered code
	_, _, err := t.runner.Run(context, action, handlerCfg.ActionId, nil)
	if err != nil {
		log.Error("Error starting action: ", err.Error())
	}
}

func (t *udpTrigger) constructStartRequest(handlerCfg *trigger.HandlerConfig, payload string) *StartRequest {

	log.Debug("Received contstructStartRequest")

	req := &StartRequest{}

	data := make(map[string]interface{})
	data["payload"] = payload
	req.Data = data

	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI string                 `json:"flowUri"`
	Data       map[string]interface{} `json:"data"`
	ReplyTo    string                 `json:"replyTo"`
}

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
	handlers []*trigger.Handler
}

// Metadata implements trigger.Trigger.Metadata
func (t *udpTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

const (
	maxDatagramSize = 8192
)

func (t *udpTrigger) Initialize(ctx trigger.InitContext) error {
	log.Debug("Initialize")
	t.handlers = ctx.GetHandlers()
	return nil

}

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
		return nil
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
		log.Errorf("UDP Listen failed: %v", err)
	} else {

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
			payloadB := buf[0:n]

			log.Debugf("Received %v from %v", payload, addr)

			//handlers := t.config.Handlers
			trgData := make(map[string]interface{})
			trgData["payload"] = payload
			trgData["buffer"] = payloadB

			log.Debug("Processing handlers")
			for _, handler := range t.handlers {
				results, err := handler.Handle(context.Background(), trgData)
				if err != nil {
					log.Error("Error starting action: ", err.Error())
				}
				log.Debugf("Ran Handler: [%s]", handler)
				log.Debugf("Results: [%v]", results)
			}
		}
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *udpTrigger) Stop() error {
	// stop the trigger
	return nil
}

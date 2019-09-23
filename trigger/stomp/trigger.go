package stomp

import (
	"context"
	"fmt"

	stomp "github.com/go-stomp/stomp"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
}

type Trigger struct {
	logger        log.Logger
	settings      *Settings
	conn          *stomp.Conn
	stompHandlers []*StompHandler
}

// StompHandler is a Stomp topic handler
type StompHandler struct {
	logger       log.Logger
	handler      trigger.Handler
	subscription *stomp.Subscription
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s}, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Init implements trigger.Init
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	//Stomp Connection Options

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(t.settings.UserName, t.settings.Password),
		stomp.ConnOpt.Host(t.settings.Address),
		stomp.ConnOpt.HeartBeat(0, 0),
	}

	t.logger = ctx.Logger()

	t.logger.Infof("Connecting to Stomp server: %v", t.settings.Address)

	//conn, err := connectStomp(t.settings.Address)
	conn, err := stomp.Dial("tcp", t.settings.Address, options...)
	if err != nil {
		t.logger.Errorf("Error connecting to Stomp: %v", err)
		return err
	}
	t.conn = conn
	for _, handler := range ctx.GetHandlers() {
		stompHandler, err := t.NewStompHandler(handler)
		if err != nil {
			return err
		}
		t.stompHandlers = append(t.stompHandlers, stompHandler)
	}
	return nil
}

// NewStompHandler creates a new stomp handler to handle a topic
func (t *Trigger) NewStompHandler(handler trigger.Handler) (*StompHandler, error) {

	stompHandler := &StompHandler{logger: t.logger, handler: handler}

	handlerSetting := &HandlerSettings{}
	err := metadata.MapToStruct(handler.Settings(), handlerSetting, true)
	if err != nil {
		return nil, err
	}

	if handlerSetting.Source == "" {
		return nil, fmt.Errorf("source string was not provided for handler: [%s]", handler)
	}

	t.logger.Debugf("Initializing subscription to source [%s]", handlerSetting.Source)

	return stompHandler, nil
}

// registerDummyEventHandler is used for dummy event handler registration, this should be replaced
// with the appropriate event handling mechanism for the trigger.  Some form of a discriminator
// should be used for dispatching to different handlers.  For example a REST based trigger might
// dispatch based on the method and path.
func registerDummyEventHandler(discriminator string, onEvent dummyOnEvent) {
	//ignore
}

// dummyOnEvent is a dummy event handler for our dummy event source
type dummyOnEvent func(interface{})

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	t.logger.Debug("Trigger Starting")

	for i, stompHandler := range t.stompHandlers {
		// start subscribing

		handlerSetting := &HandlerSettings{}
		err := metadata.MapToStruct(stompHandler.handler.Settings(), handlerSetting, true)
		if err != nil {
			t.logger.Info("Error reading metadata for handler")
			return err
		}
		destination := handlerSetting.Source
		t.logger.Infof("Subscribing to destination: '%v'", destination)

		stompHandler.subscription, err = t.conn.Subscribe(handlerSetting.Source, stomp.AckAuto, stomp.SubscribeOpt.Header("id", string(i)))
		if err != nil {
			t.logger.Info("Error subscribing to destination")
			return err
		}

	}

	//start listening
	for {
		for _, stompHandler := range t.stompHandlers {
			m := <-stompHandler.subscription.C
			newActionHandler(t, stompHandler.handler, m.Body, m.Destination)
		}
	}

	//return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	//stop servers/services if necessary
	t.logger.Debug("Trigger Stopping")
	for _, handler := range t.stompHandlers {
		err := handler.subscription.Unsubscribe()
		if err != nil {
			t.logger.Errorf("Error unsubscribing: %v", err)
		}
	}

	return nil
}

// Execute executes any handlers defined immediately on startup
func newActionHandler(t *Trigger, handler trigger.Handler, msg interface{}, destination string) {
	t.logger.Debugf("Inside ActionHandler")

	triggerData := map[string]interface{}{
		"message":        msg,
		"originalSource": destination,
	}

	response, err := handler.Handle(context.Background(), triggerData)

	if err != nil {
		t.logger.Error("Error starting action: ", err.Error())
	} else {
		t.logger.Debugf("Action call successful: %v", response)
	}
}

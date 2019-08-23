package onstartup

import (
	"context"
	"time"

	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}
type Trigger struct {
	handlers []trigger.Handler
	logger   log.Logger
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	return &Trigger{}, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
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

// Init implements trigger.Init
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.handlers = ctx.GetHandlers()
	t.logger = ctx.Logger()
	t.logger.Info("Trigger Init Called")

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	t.logger.Info("Trigger Start Called")
	for _, handler := range t.handlers {
		t.logger.Infof("Executing Handler %v", handler)
		t.Execute(handler)
	}
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	//stop servers/services if necessary
	return nil
}

// Execute executes any handlers defined immediately on startup
func (t *Trigger) Execute(handler trigger.Handler) {
	t.logger.Debug("Starting process")

	triggerData := map[string]interface{}{
		"triggerTime": time.Now().String(),
	}

	response, err := handler.Handle(context.Background(), triggerData)

	if err != nil {
		t.logger.Error("Error starting action: ", err.Error())
	} else {
		t.logger.Debugf("Action call successful: %v", response)
	}
}

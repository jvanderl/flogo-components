package onstartup

import (
	"context"
	"time"
	//"math"
	//"strconv"
	//"strings"
	//"time"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	//"github.com/carlescere/scheduler"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-onstartup")

//OnStartupTrigger is th main structure for this trigger
type OnStartupTrigger struct {
	metadata *trigger.Metadata
	//runner   action.Runner
	config   *trigger.Config
	handlers []*trigger.Handler
	//timers map[string]*scheduler.Job
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &OnStartupFactory{metadata: md}
}

// OnStartupFactory Trigger factory
type OnStartupFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *OnStartupFactory) New(config *trigger.Config) trigger.Trigger {
	return &OnStartupTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *OnStartupTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *OnStartupTrigger) Init(runner action.Runner) {
	log.Debug("Trigger Init called")
	//	t.runner = runner
	//	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Initialize implements ext.Trigger.Initialize
func (t *OnStartupTrigger) Initialize(ctx trigger.InitContext) error {
	log.Debug("Trigger Initialize called")

	t.handlers = ctx.GetHandlers()

	return nil
}

// Start implements ext.Trigger.Start
func (t *OnStartupTrigger) Start() error {

	log.Debug("Trigger Start Called")

	for _, handler := range t.handlers {
		t.Execute(handler)
	}
	return nil
}

// Stop implements ext.Trigger.Stop
func (t *OnStartupTrigger) Stop() error {

	return nil
}

// Execute executes any handlers defined immediately on startup
func (t *OnStartupTrigger) Execute(handler *trigger.Handler) {
	log.Debug("Starting process")

	triggerData := map[string]interface{}{
		"triggerTime": time.Now().String(),
	}

	response, err := handler.Handle(context.Background(), triggerData)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
	} else {
		log.Debugf("Action call successful: %v", response)
	}
}

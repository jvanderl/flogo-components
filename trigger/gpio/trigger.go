package gpio

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/stianeikeland/go-rpio"
	"strconv"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-gpio")

type GPIOTrigger struct {
	metadata   *trigger.Metadata
	runner     action.Runner
	config     *trigger.Config
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &GPIOFactory{metadata:md}
}

// TimerFactory Timer Trigger factory
type GPIOFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *GPIOFactory) New(config *trigger.Config) trigger.Trigger {
	return &GPIOTrigger{metadata: t.metadata, config:config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *GPIOTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *GPIOTrigger) Init(runner action.Runner) {
	t.runner = runner
	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)

/*	//Open rpio
	log.Debug("Opening RPIO")
	err := rpio.Open()
	if (err != nil) {
		log.Errorf("Error opening RPIO: %s", err)
		return
	}

	log.Debug("Setting handler pins to input")
	//set the pins in trigger to input
	handlers := t.config.Handlers
	for _, handler := range handlers {
		gpiopin, err := strconv.ParseInt(handler.Settings["gpiopin"].(string), 10, 64)
		if (err != nil){
			log.Errorf("Error converting GPIO pin setting to int: %v", err)
		}
		pin := rpio.Pin(gpiopin)
		log.Debugf("Setting pin %v to input")
		pin.Input()
	}
*/
}

// Start implements ext.Trigger.Start
func (t *GPIOTrigger) Start() error {

	log.Debug("Start")
	//Open rpio
	log.Debug("Opening RPIO")
	err := rpio.Open()
	if (err != nil) {
		log.Errorf("Error opening RPIO: %s", err)
		return err
	}
	handlers := t.config.Handlers
	log.Debug("Processing handlers")

//loop
for {
	log.Debug("Inside for loop")

	for _, handler := range handlers {
		log.Debug("Inside handler loop")
		gpiopin, err := strconv.ParseInt(handler.Settings["gpiopin"].(string), 10, 64)
		log.Debugf("Checking Pin: %v", gpiopin)
		if (err != nil){
			log.Errorf("Error converting GPIO pin setting to int: %v", err)
		}
		stateSetting := handler.Settings["state"].(string)
		log.Debugf("Looking for state: %v", stateSetting)
		// assign rpi pin
		pin := rpio.Pin(gpiopin)
		log.Debug("Setting pin to input")
		pin.Input()
		//check what state to read and pull opposite first
		if (stateSetting == "1") {
			log.Debug("Pulling pin down first")
			pin.PullDown()
			res := pin.Read()
			log.Debugf("Got reading: %v", res)
			if res == rpio.High {
				log.Debugf("calling runaction: %v", handler.ActionId)
				t.RunAction(handler.ActionId, stateSetting)
			}
		} else {
			log.Debug("Pulling pin up first")
			pin.PullUp()
			res := pin.Read()
			log.Debugf("Got reading: %v", res)
			if res == rpio.Low {
				log.Debugf("calling runaction: %v", handler.ActionId)
				t.RunAction(handler.ActionId, stateSetting)
			}
		}
	}

}
	return nil
}

// Stop implements ext.Trigger.Stop
func (t *GPIOTrigger) Stop() error {

	log.Debug("Stopping GPIO Trigger")
	rpio.Close()
	return nil
}

func (t *GPIOTrigger) RunAction(actionID string, state string) {
	log.Debug("Starting Immediate \"Once\" process")
	req := t.constructStartRequest(state)
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)
	action := action.Get(actionID)
	context := trigger.NewContext(context.Background(), startAttrs)
	log.Debugf("Found action: '%+x'", action)
	//log.Debugf("ActionID: '%s'", handlerCfg.ActionId)
	_, _, err := t.runner.Run(context, action, actionID, nil)
	if err != nil {
		log.Error("Error starting action: ", err.Error())
	}
}

func (t *GPIOTrigger) constructStartRequest(gpiostate string) *StartRequest {

	log.Debug("Received contstruct start request")

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["state"] = gpiostate
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	ReplyTo     string                 `json:"replyTo"`
}

func GetSettingSafe (handlerCfg *trigger.HandlerConfig, setting string, defaultValue string) string {
var retString string
	defer func() {
		if r := recover(); r != nil {
			retString = defaultValue
		}
	}()
 retString = handlerCfg.GetSetting(setting)
 return retString
}

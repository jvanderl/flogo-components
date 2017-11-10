package gpio

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/stianeikeland/go-rpio"
	"strconv"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-gpio")
type pinConfig struct {
	pin		rpio.Pin
	desiredstate rpio.State
	laststate	rpio.State
	pull	bool
}

type GPIOTrigger struct {
	metadata   *trigger.Metadata
	runner     action.Runner
	config     *trigger.Config
	pindata		 map[string]*pinConfig
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
	t.pindata = make(map[string]*pinConfig)
	handlers := t.config.Handlers
	log.Debug("Processing handlers")

	interval,err := strconv.ParseInt(t.config.GetSetting("interval"), 10, 64)
	log.Debugf("Interval: %v", interval)

	//init the test structure

	for _, handler := range handlers {
		var pinJob pinConfig
		tmpUint, err := strconv.ParseUint(handler.Settings["gpiopin"].(string), 10, 8)
		pinJob.pin = rpio.Pin(tmpUint)
		log.Debugf("Checking Pin: %v", pinJob.pin)
		if (err != nil){
			log.Errorf("Error converting GPIO pin setting to int: %v", err)
			return err
		}
		stateSetting := handler.Settings["state"].(string)
		log.Debugf("Looking for state: %v", stateSetting)
		pinJob.pull, err = strconv.ParseBool(handler.Settings["pull"].(string))
		if (err != nil){
			log.Errorf("Error converting pull setting to bool: %v", err)
			return err
		}
		// assign rpi pin
		//pinJob.pin = rpio.Pin(gpiopin)
		pinJob.desiredstate = rpio.Low
		if (stateSetting == "1") {
			pinJob.desiredstate = rpio.High
		}
		t.pindata[handler.ActionId] = &pinJob
	}

	//loop

  tickChan := time.NewTicker(time.Millisecond * time.Duration(interval)).C

	for {
			select {
			case <- tickChan:
				for _, handler := range handlers {
					pinConf := t.pindata[handler.ActionId]
					pin := pinConf.pin
					pin.Input()
					//check what state to read and pull opposite first
					if (pinConf.desiredstate == rpio.High) {
						if (pinConf.pull) {
							log.Debug("Pulling pin down first")
							pin.PullDown()
						}
						res := pin.Read()
						log.Debugf("Got reading: %v", res)
						if (res == rpio.High && pinConf.laststate != res) {
							pinConf.laststate = res
							log.Debugf("calling runaction: %v", handler.ActionId)
							t.RunAction(handler.ActionId, "1")
						}
					} else {

						if (pinConf.pull) {
							log.Debug("Pulling pin up first")
							pin.PullUp()
						}
						res := pin.Read()
						log.Debugf("Got reading: %v", res)
						if (res == rpio.Low && pinConf.laststate != res) {
							pinConf.laststate = res
							log.Debugf("calling runaction: %v", handler.ActionId)
							t.RunAction(handler.ActionId, "0")
						}
					}
					t.pindata[handler.ActionId] = pinConf
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

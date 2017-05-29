package blemaster

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"strings"
	"strconv"
	"context"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-blemaster")

var done = make(chan struct{})

type BleService struct {
	serviceID string
	characteristic string
	actionID	string
}
type  BleTarget struct {
	devicename string
	deviceid string
	localname string
	bleservices []BleService
}

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	bletarget BleTarget
//	destinationToactionID map[string]string
}


//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &BleMasterFactory{metadata: md}
}

// BleMaster Trigger factory
type BleMasterFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *BleMasterFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *MyTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	// start the trigger
	t.bletarget.devicename = t.config.GetSetting("devicename")
	for _, handlerCfg := range t.config.Handlers {
		log.Infof("Adding BLE Service: [%s.%s]", strings.ToUpper(handlerCfg.GetSetting("service")),strings.ToUpper(handlerCfg.GetSetting("characteristic")))
		t.bletarget.bleservices = append(t.bletarget.bleservices, BleService{strings.ToUpper(handlerCfg.GetSetting("service")), strings.ToUpper(handlerCfg.GetSetting("characteristic")), handlerCfg.ActionId})
	}

	//attempt to open own local bt device
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Errorf("Failed to open device, err: %s\n", err)
		return err
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(t.onPeriphDiscovered),
		gatt.PeripheralConnected(t.onPeriphConnected),
		gatt.PeripheralDisconnected(t.onPeriphDisconnected),
	)

	d.Init(onStateChanged)
	<-done
	log.Info ("*** At end of MyTrigger Start ***")
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	log.Info ("In MyTrigger Stop")
	return nil
}

func onStateChanged(d gatt.Device, s gatt.State) {
	log.Infof("State: [%s]", s)
	switch s {
	case gatt.StatePoweredOn:
		log.Info("Scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func (t *MyTrigger) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	var found bool = false
	log.Info ("In onPeriphDiscovered")

	//check against localname
	if t.bletarget.devicename != "" {
		if strings.ToUpper(a.LocalName) != strings.ToUpper(t.bletarget.devicename) {
			return
		}
	}


	// check if service is specified in any handler
	for _, svc := range a.Services {
		for _, target := range t.bletarget.bleservices {
			log.Infof ("checking svc [%s]", svc.String())
	    if strings.ToUpper(svc.String()) == target.serviceID {
					log.Infof ("Found service [%s]", target.serviceID)
					found = true
	        continue
	    }
		}
	}
if !found {
	return
	}
	t.bletarget.deviceid = p.ID()
	t.bletarget.localname = a.LocalName
	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()
	// Now connect to the device
	p.Device().Connect(p)
}

func (t *MyTrigger) onPeriphConnected(p gatt.Peripheral, err error) {
	log.Info("Device Connected")
	defer p.Device().CancelConnection(p)

/*	if err := p.SetMTU(500); err != nil {
		log.Infof("Failed to set MTU, err: %s\n", err)
	}
*/
	// Discovery services
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		log.Errorf("Failed to discover services, err: %s\n", err)
		return
	}

	for _, s := range ss {
		for _, target := range t.bletarget.bleservices {
			if strings.ToUpper(s.UUID().String()) == target.serviceID {
				log.Infof ("Discovering characteristics for service [%s] ", target.serviceID)
				// Discovery characteristics
				cs, err := p.DiscoverCharacteristics(nil, s)
				if err != nil {
					log.Infof("Failed to discover characteristics, err: %s\n", err)
					continue
				}
				for _, c := range cs {

					if strings.ToUpper(c.UUID().String()) == target.characteristic {
						log.Infof ("Found characteristic [%s]", target.characteristic)

						// Subscribe to the characteristic
						log.Info("Subscribing to characteristic...")
						if (c.Properties() & (gatt.CharNotify | gatt.CharIndicate)) != 0 {
							f := func(c *gatt.Characteristic, b []byte, err error) {
								log.Infof("notified on %s.%s: % X | %q\n", target.serviceID, target.characteristic, b, b)
								t.RunAction (target.actionID, string(b), target.serviceID, target.characteristic)
							}
							if err := p.SetNotifyValue(c, f); err != nil {
								log.Infof("Failed to subscribe characteristic, err: %s\n", err)
								continue
							}
						}
					}
				}
			}
		}
	}
}

func (t *MyTrigger) onPeriphDisconnected(p gatt.Peripheral, err error) {
	log.Info("Device Disconnected")
//	close(done)
	reconnect, err := strconv.ParseBool(t.config.GetSetting("autoreconnect"))
	if err != nil {
		// Invalid Auto Reconnect Switch
		log.Errorf ("Invalid Auto Reconnect Switch [%s]", t.config.GetSetting("autoreconnect"))
		return
		}
  if reconnect {
		i, err := strconv.ParseInt(t.config.GetSetting("reconnectinterval"), 10, 64)
		if err != nil {
			{
				// Invalid Reconnect Interval
				log.Errorf ("Invalid Reconnect Interval [%s]", t.config.GetSetting("reconnectinterval"))
				return
			}
		}
		intervalDuration := time.Duration(i)
		log.Infof("Waiting to reconnect for %d %s", i, t.config.GetSetting("intervaltype"))
		switch t.config.GetSetting("intervaltype") {
		case "hours":
			time.Sleep(intervalDuration * time.Hour)
		case "minutes":
			time.Sleep(intervalDuration * time.Minute)
		case "seconds":
			time.Sleep(intervalDuration * time.Second)
		case "milliseconds":
			time.Sleep(intervalDuration * time.Millisecond)
		default:
			{
				// Invalid Interval Type
				log.Errorf ("Invalid Interval Type [%s]", t.config.GetSetting("intervaltype"))
				return
			}
		}
		p.Device().Connect(p)
	}
}

// RunAction starts a new Process Instance
func (t *MyTrigger) RunAction(actionID string, payload string, serviceID string, characteristic string) {

	log.Debug("Starting new Process Instance")
	log.Debugf("Action ID: ", actionID)
	log.Debugf("Payload: ", payload)
	log.Debugf("Service ID: ", serviceID)
	log.Debugf("Characteristic ", characteristic)

	req := t.constructStartRequest(payload)
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)
	action := action.Get(actionID)
	context := trigger.NewContext(context.Background(), startAttrs)
	_, replyData, err := t.runner.Run(context, action, actionID, nil)
	if err != nil {
		log.Error(err)
	}

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

func (t *MyTrigger) constructStartRequest(payload string) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["notification"] = payload
	data["deviceid"] = t.bletarget.deviceid
	data["localname"] = t.bletarget.localname
	//	data["destination"] = destination
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

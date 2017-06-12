package blemaster

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/go-gatt"
	"github.com/jvanderl/go-gatt/examples/option"
//	"github.com/paypal/gatt"
//	"github.com/paypal/gatt/examples/option"
	"strings"
	"strconv"
	"context"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-blemaster")

var done = make(chan struct{})

type BleService struct {
	serviceID gatt.UUID
	characteristic gatt.UUID
	actionID	string
	gotData bool
}
type  BleTarget struct {
	devicename string
	devicecheck bool
	autodisconnect bool
	autoreconnect bool
	sleeptime time.Duration
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
	t.bletarget.deviceid = t.config.GetSetting("deviceid")
	autodisconnect, err := strconv.ParseBool(t.config.GetSetting("autodisconnect"))

	if err != nil {
		// Invalid Auto Reconnect Switch
		log.Errorf ("Invalid Auto Reconnect Switch [%s]", t.config.GetSetting("autodisconnect"))
		return err
	}
	t.bletarget.autodisconnect = autodisconnect
	log.Infof ("Auto Disconnect [%v]", t.bletarget.autodisconnect)

	autoreconnect, err := strconv.ParseBool(t.config.GetSetting("autoreconnect"))
	if err != nil {
		// Invalid Auto Reconnect Switch
		log.Errorf ("Invalid Auto Reconnect Switch [%s]", t.config.GetSetting("autoreconnect"))
		return err
	}
	t.bletarget.autoreconnect = autoreconnect
	log.Infof ("Auto Reconnect [%v]", t.bletarget.autoreconnect)

	i, err := strconv.ParseInt(t.config.GetSetting("reconnectinterval"), 10, 64)
	if err != nil {
		{
			// Invalid Reconnect Interval
			log.Errorf ("Invalid Reconnect Interval [%s]", t.config.GetSetting("reconnectinterval"))
			return err
		}
	}
	intervalDuration := time.Duration(i)
	log.Infof("Waiting to reconnect for %d %s", i, t.config.GetSetting("intervaltype"))
	switch t.config.GetSetting("intervaltype") {
	case "hours":
		t.bletarget.sleeptime = intervalDuration * time.Hour
	case "minutes":
		t.bletarget.sleeptime = intervalDuration * time.Minute
	case "seconds":
		t.bletarget.sleeptime = intervalDuration * time.Second
	case "milliseconds":
		t.bletarget.sleeptime = intervalDuration * time.Millisecond
	default:
		{
			// Invalid Interval Type
			log.Errorf ("Invalid Interval Type [%s]", t.config.GetSetting("intervaltype"))
			return err
		}
	}

	log.Infof ("Auto Reconnect Interval [%v] [%s]", t.bletarget.autoreconnect, t.config.GetSetting("intervaltype"))

	for _, handlerCfg := range t.config.Handlers {
		log.Infof("Adding BLE Service: [%s.%s]", handlerCfg.GetSetting("service"),handlerCfg.GetSetting("characteristic"))
		t.bletarget.bleservices = append(t.bletarget.bleservices, BleService{gatt.MustParseUUID(handlerCfg.GetSetting("service")), gatt.MustParseUUID(handlerCfg.GetSetting("characteristic")), handlerCfg.ActionId, false})
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
	log.Info ("*** At end of MyTrigger Start ***")
	<-done
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
	//var found bool = false
	log.Info ("In onPeriphDiscovered")
	log.Infof ("a.LocalName [%s]", a.LocalName)
	log.Infof ("p.Name [%v]", p.Name())
	log.Infof ("p.ID [%v]", p.ID())

	//check against localname
	if strings.ToUpper(a.LocalName) == strings.ToUpper(t.bletarget.devicename) || strings.ToUpper(p.ID()) == strings.ToUpper(t.bletarget.deviceid) {
		if t.bletarget.deviceid == "" {
			t.bletarget.deviceid = p.ID()
		}
		if t.bletarget.localname == "" {
			t.bletarget.localname = a.LocalName
		}
		// Stop scanning once we've got the peripheral we're looking for.
		p.Device().StopScanning()
		// Now connect to the device
		p.Device().Connect(p)
	}
}

func (t *MyTrigger) onPeriphConnected(p gatt.Peripheral, err error) {
	log.Info("Device Connected")
	//defer p.Device().CancelConnection(p)
	time.Sleep(100 * time.Millisecond)
	// Discovery services
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		log.Errorf("Failed to discover services, err: %s\n", err)
		return
	}

	for _, s := range ss {
		for _, target := range t.bletarget.bleservices {
			if (s.UUID().Equal(target.serviceID)) {
				log.Infof ("Discovering characteristics for service [%s] ", target.serviceID)
				// Discovery characteristics
				cs, err := p.DiscoverCharacteristics(nil, s)
				if err != nil {
					log.Infof("Failed to discover characteristics, err: %s\n", err)
					continue
				}
				for _, c := range cs {
					if (c.UUID().Equal(target.characteristic)) {
						log.Infof ("Found characteristic [%s]", target.characteristic)
						// Discovery descriptors
						p.DiscoverDescriptors(nil, c)
						// Subscribe to the characteristic
						log.Info("Subscribing to characteristic...")
						if (c.Properties() & (gatt.CharNotify | gatt.CharIndicate)) != 0 {
							f := func(c *gatt.Characteristic, b []byte, err error) {
								log.Infof("notified on %s.%s: %q", target.serviceID, target.characteristic, b)
								t.RunAction (p, target, string(b))
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
  if t.bletarget.autoreconnect {
		time.Sleep(t.bletarget.sleeptime)
		p.Device().Connect(p)
	}
}

	// RunAction starts a new Process Instance
func (t *MyTrigger) RunAction(p gatt.Peripheral, bleService BleService, payload string) {

	log.Debug("Starting new Process Instance")
	log.Debugf("Action ID: ", bleService.actionID)
	log.Debugf("Payload: ", payload)
	log.Debugf("Service ID: ", bleService.serviceID)
	log.Debugf("Characteristic ", bleService.characteristic)

	if t.bletarget.autodisconnect == true {
		log.Info ("Checking if all data was received")
		var gotItAll bool = true
		for _, target := range t.bletarget.bleservices {
			if (target.serviceID.Equal(bleService.serviceID) && target.characteristic.Equal(bleService.characteristic)) {
				target.gotData = true
			}
			if !target.gotData {
				gotItAll = false
			}
		}

		if gotItAll {
			log.Info ("All data was recieved, disconnecting")
			p.Device().CancelConnection(p)
		}
	}

	req := t.constructStartRequest(payload)
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)
	action := action.Get(bleService.actionID)
	context := trigger.NewContext(context.Background(), startAttrs)
	_, replyData, err := t.runner.Run(context, action, bleService.actionID, nil)
	if err != nil {
		log.Error(err)
	}

	log.Info("Reply data: ", replyData)

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

func (t *MyTrigger) GotAllData() bool {
	for _, target := range t.bletarget.bleservices {
		log.Infof ("Checking target [%s] data was received: [%v]", target.serviceID, target.gotData)
		if !target.gotData {
			log.Infof ("Target [%s], gotData [%v]", target.serviceID, target.gotData)
			return false
		}
	}
	return true
}

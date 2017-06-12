package blewrite

import (
//	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/jvanderl/go-gatt"
	"github.com/jvanderl/go-gatt/examples/option"
	"strings"
	"strconv"
//	"context"
	"time"
	"errors"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-blewrite")

type response struct {
	result string
	err error
}
var resp response
var allDone = make(chan response)

const (
	devicename   		= "devicename"
	deviceid				= "deviceid"
	serviceid				= "serviceid"
	characteristic 	= "characteristic"
	bledata 				=	"bledata"
	disconnectdelay			= "disconnectdelay"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
	DeviceName string
	DeviceID string
	ServiceID gatt.UUID
	Characteristic gatt.UUID
	BLEdata string
	DisconnectDelay time.Duration
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	a.DeviceName = context.GetInput(devicename).(string)
	log.Infof ("DeviceName: [%s]", a.DeviceName)

	a.DeviceID = context.GetInput(deviceid).(string)
	log.Infof ("DeviceID: [%s]", a.DeviceID)

	a.ServiceID = gatt.MustParseUUID(context.GetInput(serviceid).(string))
	log.Infof ("ServiceID: [%s]", a.ServiceID)

	a.Characteristic = gatt.MustParseUUID(context.GetInput(characteristic).(string))
	log.Infof ("Characteristic: [%s]", a.Characteristic)

	a.BLEdata = context.GetInput(bledata).(string)
	log.Infof ("BLEdata: [%s]", a.BLEdata)

	var tmp int64
	tmp, err = strconv.ParseInt(context.GetInput(disconnectdelay).(string),10,32)
	if err != nil {
		log.Errorf ("Error converting AddNewLine to input to integer")
		a.DisconnectDelay = time.Duration(400)
	} else {
		a.DisconnectDelay = time.Duration(tmp)
	}

	//attempt to open own local bt device
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Errorf("Failed to open device, err: %s\n", err)
		return true, err
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(a.onPeriphDiscovered),
		gatt.PeripheralConnected(a.onPeriphConnected),
		gatt.PeripheralDisconnected(a.onPeriphDisconnected),
	)

	d.Init(a.onStateChanged)
	<-allDone
	context.SetOutput("result", resp.result)
	return true, resp.err
}

func (a *MyActivity) onPeriphDiscovered(p gatt.Peripheral, ad *gatt.Advertisement, rssi int) {
//	var found bool = false
	log.Info ("In onPeriphDiscovered")
	log.Infof ("ad.LocalName [%s]", ad.LocalName)
	log.Infof ("p.Name [%v]", p.Name())
	log.Infof ("p.ID [%v]", p.ID())
	log.Infof ("a.DeviceName [%v]", a.DeviceName)
	log.Infof ("a.DeviceID [%v]", a.DeviceID)

	//check against localname or peripheral ID
	if strings.ToUpper(ad.LocalName) == strings.ToUpper(a.DeviceName) || strings.ToUpper(p.ID()) == strings.ToUpper(a.DeviceID) {
		log.Infof ("Found matching LocalName or DeviceID")
		// Stop scanning once we've got the peripheral we're looking for.
		p.Device().StopScanning()
		// Now connect to the device
		p.Device().Connect(p)
	}
}


func (a *MyActivity) onPeriphConnected(p gatt.Peripheral, err error) {

	log.Info("Device Connected")
	defer p.Device().CancelConnection(p)
	time.Sleep(100 * time.Millisecond)

	// Discover services
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		resp.result = "ERR_DISCOVER_SERVICES"
		resp.err = err
		log.Errorf("Failed to discover services, err: %s\n", err)
		return
	}

	for _, s := range ss {
		if (s.UUID().Equal(a.ServiceID)) {
			log.Infof ("Discovering characteristics for service [%s] ", a.ServiceID)
			// Discovery characteristics
			cs, err := p.DiscoverCharacteristics(nil, s)
			if err != nil {
				log.Infof("Failed to discover characteristics, err: %s\n", err)
				continue
			}

			for _, c := range cs {
				if (c.UUID().Equal(a.Characteristic)) {
					log.Infof ("Found characteristic [%s]", a.Characteristic)

					p.DiscoverDescriptors(nil, c)

					if (c.Properties() & gatt.CharWriteNR) == 0  {
						resp.result = "ERR_WRITE_DISABLED"
						resp.err = errors.New("Charateristic is not WriteNR enabled")
						log.Errorf("Charateristic is not WriteNR enabled [%v]", c.Properties())
						return
					}
				}
				var tmpdata = []byte(a.BLEdata)
				log.Infof ("Writing BLE data [% X], [\"%s\"]...", tmpdata, tmpdata)
			  err = p.WriteCharacteristic(c, tmpdata, true)
				if err != nil {
					resp.result = "ERR_WRITE_CHARACTERISTIC"
					resp.err = err
					log.Errorf("Failed to write characteristic, err: %s", err)
					return
				}
				resp.result = "OK"
				resp.err = nil
				log.Info ("BLE data write successful")
				time.Sleep(a.DisconnectDelay * time.Millisecond)
				return
			}
		}
	}
}

func (a *MyActivity) onStateChanged(d gatt.Device, s gatt.State) {
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

func (a *MyActivity) onPeriphDisconnected(p gatt.Peripheral, err error) {
	log.Info("Device Disconnected")
	allDone <- resp
}

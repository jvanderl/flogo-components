package mqtt

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-rest")

const (
	broker   = "broker"
	topic    = "topic"
	qos      = "qos"
	payload  = "message"
	id       = "id"
	user     = "user"
	password = "password"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&MyActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	brokerInput := context.GetInput(broker)

	ivbroker, ok := brokerInput.(string)
	if !ok {
		context.SetOutput("result", "BROKER_NOT_SET")
		return true, fmt.Errorf("Broker not set.")
	}

	topicInput := context.GetInput(topic)

	ivtopic, ok := topicInput.(string)
	if !ok {
		context.SetOutput("result", "TOPIC_NOT_SET")
		return true, fmt.Errorf("Topic not set.")
	}

	payloadInput := context.GetInput(payload)

	ivpayload, ok := payloadInput.(string)
	if !ok {
		context.SetOutput("result", "PAYLOAD_NOT_SET")
		return true, fmt.Errorf("Payload not set.")
	}

	ivqos, ok := context.GetInput(qos).(int)

	if !ok {
		context.SetOutput("result", "QOS_NOT_SET")
		return true, fmt.Errorf("QoS not set.")
	}

	idInput := context.GetInput(id)

	ivID, ok := idInput.(string)
	if !ok {
		context.SetOutput("result", "CLIENTID_NOT_SET")
		return true, fmt.Errorf("Client ID not set.")
	}

	userInput := context.GetInput(user)

	ivUser, ok := userInput.(string)
	if !ok {
		//User not set, use default
		ivUser = ""
	}

	passwordInput := context.GetInput(password)

	ivPassword, ok := passwordInput.(string)
	if !ok {
		//Password not set, use default
		ivPassword = ""
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(ivbroker)
	opts.SetClientID(ivID)
	opts.SetUsername(ivUser)
	opts.SetPassword(ivPassword)
	client := mqtt.NewClient(opts)

	log.Debugf("MQTT Publisher connecting")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Debugf("MQTT Publisher connected, sending message")
	token := client.Publish(ivtopic, byte(ivqos), false, ivpayload)
	token.Wait()

	client.Disconnect(250)
	log.Debugf("MQTT Publisher disconnected")
	context.SetOutput("result", "OK")

	return true, nil
}

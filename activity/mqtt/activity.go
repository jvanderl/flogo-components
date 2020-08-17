package mqtt

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("Settings Broker: %s", settings.Broker)
	ctx.Logger().Debugf("Settings Id	: %s", settings.Id)
	ctx.Logger().Debugf("Settings User	: %s", settings.User)

	conn, err := getMqttConnection(ctx.Logger(), settings)
	if err != nil {
		return nil, err
	}

	act := &Activity{conn: conn}

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	conn *MqttConnection
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	ctx.Logger().Infof("Topic: %s", input.Topic)
	ctx.Logger().Infof("Qos: %v", input.Qos)
	ctx.Logger().Infof("Message: %s", input.Message)

	if act.conn.client.IsConnected() {
		ctx.Logger().Infof("Client is already connected")
	} else {
		ctx.Logger().Infof("MQTT Publisher connecting")
		if token := act.conn.client.Connect(); token.Wait() && token.Error() != nil {
			output := &Output{Result: "ERROR"}
			err = ctx.SetOutputObject(output)
			return false, token.Error()
		}
	}

	ctx.Logger().Infof("MQTT Publisher connected, sending message")
	token := act.conn.client.Publish(input.Topic, byte(input.Qos), false, input.Message)
	token.Wait()
	if token.Error() != nil {
		output := &Output{Result: "ERROR"}
		err = ctx.SetOutputObject(output)
		return false, token.Error()
	}

	output := &Output{Result: "OK"}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

/*
// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	brokerInput := context.GetInput(broker)

	ivbroker, ok := brokerInput.(string)
	if !ok {
		context.SetOutput("result", "BROKER_NOT_SET")
		return true, fmt.Errorf("broker not set")
	}

	topicInput := context.GetInput(topic)

	ivtopic, ok := topicInput.(string)
	if !ok {
		context.SetOutput("result", "TOPIC_NOT_SET")
		return true, fmt.Errorf("topic not set")
	}

	payloadInput := context.GetInput(payload)

	ivpayload, ok := payloadInput.(string)
	if !ok {
		context.SetOutput("result", "PAYLOAD_NOT_SET")
		return true, fmt.Errorf("payload not set")
	}

	ivqos, ok := context.GetInput(qos).(int)

	if !ok {
		context.SetOutput("result", "QOS_NOT_SET")
		return true, fmt.Errorf("qos not set")
	}

	idInput := context.GetInput(id)

	ivID, ok := idInput.(string)
	if !ok {
		context.SetOutput("result", "CLIENTID_NOT_SET")
		return true, fmt.Errorf("client id not set")
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
*/

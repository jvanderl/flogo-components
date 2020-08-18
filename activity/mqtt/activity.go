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

	ctx.Logger().Debugf("Topic: %s", input.Topic)
	ctx.Logger().Debugf("Qos: %v", input.Qos)
	ctx.Logger().Debugf("Message: %s", input.Message)

	if act.conn.client.IsConnected() {
		ctx.Logger().Debugf("Client is already connected")
	} else {
		ctx.Logger().Debugf("MQTT Publisher connecting")
		if token := act.conn.client.Connect(); token.Wait() && token.Error() != nil {
			output := &Output{Result: "ERROR"}
			err = ctx.SetOutputObject(output)
			return false, token.Error()
		}
	}

	ctx.Logger().Debugf("MQTT Publisher connected, sending message")
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

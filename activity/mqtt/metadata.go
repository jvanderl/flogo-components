package mqtt

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Broker   string `md:"broker,required"` // The MQTT Broker URI (tcp://[hostname]:[port])
	Id       string `md:"id,required"`     // The MQTT Connection Client ID
	User     string `md:"user"`            // The UserID used when connecting to the MQTT broker
	Password string `md:"password"`        // The Password used when connecting to the MQTT broker
}

type Input struct {
	Topic   string `md:"topic,required"`   // Topic on which the message is published
	Qos     int    `md:"qos,required"`     // MQTT Quality of Service. 0 = At most once, 1 = At least once, 2 = Exactly once.
	Message string `md:"message,required"` // The message payload
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"topic":   i.Topic,
		"qos":     i.Qos,
		"message": i.Message,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Topic, err = coerce.ToString(values["topic"])
	i.Qos, err = coerce.ToInt(values["qos"])
	if err != nil {
		return err
	}
	i.Message, err = coerce.ToString(values["message"])
	return err
}

type Output struct {
	Result string `md:"result"` // The result of message sending
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Result, _ = coerce.ToString(values["result"])

	return nil
}

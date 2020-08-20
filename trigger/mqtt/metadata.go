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

type HandlerSettings struct {
	Topic string `md:"topic,required"` // The topic to subscribe to
	Qos   int    `md:"qos,required"`   // MQTT Quality of Service. 0 = At most once, 1 = At least once, 2 = Exactly once.
}

type Output struct {
	Message     string `md:"message"`     // The message that was received
	ActualTopic string `md:"actualtopic"` // The message that was received
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}
	o.ActualTopic, err = coerce.ToString(values["actualtopic"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message":     o.Message,
		"actualtopic": o.ActualTopic,
	}
}

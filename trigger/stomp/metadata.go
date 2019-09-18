package stomp

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Address string `md:"address"` // The address of the stomp server. Example: localhost:61613
}

type HandlerSettings struct {
	Source string `md:"source,required"` // The topic or queue to subscribe to
}

type Output struct {
	Message        interface{} `md:"message"`        // The message that was received
	OriginalSource string      `md:"originalSource"` // The topic or queue name where the message was received on
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message":        o.Message,
		"originalSource": o.OriginalSource,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Message = values["message"]
	o.OriginalSource, err = coerce.ToString(values["originalSource"])
	if err != nil {
		return err
	}
	return nil
}

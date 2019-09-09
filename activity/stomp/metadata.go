package stomp

import "github.com/project-flogo/core/data/coerce"

type Input struct {
	Address     string `md:"address,required"`     //The Stomp Address example: localhost:61613
	Destination string `md:"destination,required"` //The queue to post message on
	Message     string `md:"message,required"`     //The message to post
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Address, err = coerce.ToString(values["address"])
	if err != nil {
		return err
	}
	i.Destination, err = coerce.ToString(values["destination"])
	if err != nil {
		return err
	}
	i.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"address":     i.Address,
		"destination": i.Destination,
		"message":     i.Message,
	}
}

type Output struct {
	Result string `md:"result"` // The Result of posting message
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Result, err = coerce.ToString(values["result"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
	}
}

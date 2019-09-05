package slack

import "github.com/project-flogo/core/data/coerce"

type Input struct {
	Token   string `md:"token,required"`   //The Slack App Token
	Channel string `md:"channel,required"` //The Channel to post message on
	Message string `md:"message,required"` //The message to post
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Token, err = coerce.ToString(values["token"])
	if err != nil {
		return err
	}
	i.Channel, err = coerce.ToString(values["channel"])
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
		"token":   i.Token,
		"channel": i.Channel,
		"message": i.Message,
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

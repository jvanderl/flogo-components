package onstartup

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
}

type Output struct {
	TriggerTime string `md:"triggerTime"` // Time the trigger was fired
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"triggerTime": o.TriggerTime,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.TriggerTime, err = coerce.ToString(values["triggerTime"])
	if err != nil {
		return err
	}
	return nil
}

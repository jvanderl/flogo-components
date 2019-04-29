package amqp

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestPublishQueue(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("server", "localhost")
	tc.SetInput("port", "5672")
	tc.SetInput("userID", "test")
	tc.SetInput("password", "test")
	tc.SetInput("exchange", "")
	tc.SetInput("routingKey", "hello")
	tc.SetInput("routingType", "queue")
	tc.SetInput("message", "Hello World, from Flogo!")
	tc.SetInput("durable", "false")
	tc.SetInput("autoDelete", "false")
	tc.SetInput("exclusive", "false")
	tc.SetInput("noWait", "false")

	act.Eval(tc)

	//check result attr
	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestPublishTopic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("server", "localhost")
	tc.SetInput("port", "5672")
	tc.SetInput("userID", "test")
	tc.SetInput("password", "test")
	tc.SetInput("exchange", "logs_topic")
	tc.SetInput("routingKey", "kern.critical")
	tc.SetInput("routingType", "topic")
	tc.SetInput("message", "A critical kernel error from Flogo")
	tc.SetInput("durable", "true")
	tc.SetInput("autoDelete", "false")
	tc.SetInput("exclusive", "false")
	tc.SetInput("noWait", "false")

	act.Eval(tc)

	//check result attr
	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

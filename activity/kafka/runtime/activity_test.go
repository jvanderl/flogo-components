package kafka

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("kafka")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &KafkaActivity{metadata: md}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(server, "10.10.1.50:9092")
	tc.SetInput(configid, "flogo-test")
	tc.SetInput(topic, "test")
	tc.SetInput(partition, "0")
	tc.SetInput(message, "Kafka test message")

	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")

	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

	// TODO how to do some checks if the activity has no Output?
}

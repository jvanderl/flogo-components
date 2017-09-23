package kafka

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"strconv"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-kafka")

const (
	server    = "server"
	configid  = "configid"
	topic     = "topic"
	message   = "message"
	partition = "partition"
	result    = "result"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	ifServers := []string{context.GetInput(server).(string)}
	log.Debug("ifServers: ", ifServers)

	ifConfigID := context.GetInput(configid).(string)
	ifTopic := context.GetInput(topic).(string)
	cxPartition := context.GetInput(partition)
	ifPartition, ok := cxPartition.(int32)
	if !ok {
		ifPartition = 0
	}

	ifMessage := context.GetInput(message).(string)

	conf := kafka.NewBrokerConf(ifConfigID)
	conf.AllowTopicCreation = true

	// connect to kafka cluster
	log.Debug("Connecting to Kafka server")
	broker, err := kafka.Dial(ifServers, conf)
	if err != nil {
		log.Errorf("cannot connect to kafka cluster: %s", err)
		context.SetOutput("result", "ERROR_KAFKA_CONNECT")
		return true, nil
	}
	defer broker.Close()
	log.Debug("Connected to Kafka server")

	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: []byte(ifMessage)}

	log.Debug("Sending message to Kafka server")
	resp, err := producer.Produce(ifTopic, ifPartition, msg)

	if err != nil {
		log.Error("Error sending message to Kafka broker:", err)
		context.SetOutput("result", "ERROR_KAFKA_SEND")
		return true, nil
	}

	log.Debug("Response:", resp)
	log.Debug("Message sent succesfully")

	context.SetOutput("result", strconv.FormatInt(resp, 10))
	return true, nil
}

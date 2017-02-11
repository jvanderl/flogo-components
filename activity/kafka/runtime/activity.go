package kafka

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-kafka")

const (
	server	    = "server"
	configid	= "configid"
	topic       = "topic"
	message     = "message"
	partition   = "partition"
	result      = "result"
)

// KafkaActivity is a Kafka Activity implementation
type KafkaActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&KafkaActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *KafkaActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *KafkaActivity) Eval(context activity.Context) (done bool, err error) {

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
		log.Fatalf("cannot connect to kafka cluster: %s", err)
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

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Response:", resp)
	}
	log.Debug("Message sent succesfully")

	context.SetOutput("result", resp)
	return true, nil
}

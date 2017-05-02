package kafka

import (
	"context"
	//	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"strconv"
	//	"bytes"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-kafka")
var broker kafka.Broker

// KafkaTrigger is simple Kafka trigger
type kafkaTrigger struct {
	metadata        *trigger.Metadata
	runner          action.Runner
	config          *trigger.Config
	topicToActionId map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &kafkaFactory{metadata: md}
}

// kafkaFactory Trigger factory
type kafkaFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *kafkaFactory) New(config *trigger.Config) trigger.Trigger {
	return &kafkaTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *kafkaTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *kafkaTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *kafkaTrigger) Start() error {

	ifServers := []string{t.config.GetSetting("server")}

	ifConfigID := t.config.GetSetting("configid")
	ifTopic := t.config.GetSetting("topic")
	cxPartition, err := strconv.ParseInt(t.config.GetSetting("partition"), 10, 32)
	ifPartition := int32(0)
	if err == nil {
		ifPartition = int32(cxPartition)
	}

	conf := kafka.NewBrokerConf(ifConfigID)
	conf.AllowTopicCreation = true

	// connect to kafka cluster
	log.Infof("Connecting to Kafka server(s): %s", ifServers)
	broker, err := kafka.Dial(ifServers, conf)
	if err != nil {
		log.Errorf("cannot connect to kafka cluster: %s", err)
		return err
	}
	defer broker.Close()
	log.Info("Connected to Kafka server")

	conf2 := kafka.NewConsumerConf(ifTopic, ifPartition)
	conf2.StartOffset = kafka.StartOffsetNewest
	log.Infof("subscribing to topic [%s]", ifTopic)
	consumer, err := broker.Consumer(conf2)
	if err != nil {
		log.Errorf("cannot create kafka consumer for %s:%d: %s", ifTopic, ifPartition, err)
	}
	log.Info("Subscription successful")

	t.topicToActionId = make(map[string]string)
	for _, handlerCfg := range t.config.Handlers {
		log.Infof("Regestering Action [%s] for topic [%s]", handlerCfg.ActionId, handlerCfg.GetSetting("topic"))
		t.topicToActionId[handlerCfg.GetSetting("topic")] = handlerCfg.ActionId
	}

	for {
		msg, err := consumer.Consume()
		if err != nil {
			if err != kafka.ErrNoData {
				log.Infof("cannot consume %q topic message: %s", ifTopic, err)
			}
			break
		}
		message := convert(msg.Value)
		log.Infof("Received message: %d - %s", msg.Offset, message)
		actionId, found := t.topicToActionId[ifTopic]
		if found {
			log.Infof("Found actionId: %s", actionId)
			t.RunAction(actionId, message, ifTopic, ifPartition)
		} else {
			log.Info("actionId not found")
		}
	}
	log.Info("consumer quit")

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *kafkaTrigger) Stop() error {

	broker.Close()

	return nil
}

// RunAction starts a new Process Instance
func (t *kafkaTrigger) RunAction(actionId string, payload string, topic string, partition int32) {

	log.Info("Starting new Process Instance")
	log.Infof("Action Id: %s", actionId)
	log.Infof("Payload: %s", payload)
	log.Infof("Actual Topic: %s ", topic)

	req := t.constructStartRequest(payload, topic)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionId)

	context := trigger.NewContext(context.Background(), startAttrs)

	_, replyData, err := t.runner.Run(context, action, actionId, nil)
	if err != nil {
		log.Error(err)
	}

	log.Infof("Ran action: [%s]", actionId)
	log.Infof("Reply data: [%s]", replyData)

}

func (t *kafkaTrigger) publishMessage(topic string, partition int32, message string) {

	log.Infof("ReplyTo topic: %s", topic)
	log.Infof("Publishing message: %s", message)

	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: []byte(message)}

	log.Info("Sending message to Kafka server")
	resp, err := producer.Produce(topic, partition, msg)

	if err != nil {
		log.Errorf("Error sending message to Kafka broker: %s", err)
	}

	log.Infof("Response: %s", resp)
	log.Info("Message sent succesfully")
}

func (t *kafkaTrigger) constructStartRequest(message string, topic string) *StartRequest {

	log.Info("Received contstruct start request")

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}

func convert(b []byte) string {
	n := len(b)
	return string(b[:n])
}

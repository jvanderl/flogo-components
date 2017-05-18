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

// Subscription scructure
type kafkaSubscription struct {
	topic			string
	partition int32
	consumer	kafka.Consumer
}

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
	log.Debug("Connected to Kafka server")

	//subscribe to topics
	var Subscriptions []kafkaSubscription
	t.topicToActionId = make(map[string]string)
	for _, handlerCfg := range t.config.Handlers {
		tTopic := handlerCfg.GetSetting("topic")
		//Set default partition to 0
		tPartition := int32(0)
		sPartition, err := strconv.ParseInt(handlerCfg.GetSetting("partition"), 10, 32)
		if err == nil {
			tPartition = int32(sPartition)
		}
		log.Infof("Regestering Action [%s] for topic [%s], partition [%d]", handlerCfg.ActionId, tTopic, tPartition)
		t.topicToActionId[tTopic] = handlerCfg.ActionId
		conf2 := kafka.NewConsumerConf(tTopic, tPartition)
		conf2.StartOffset = kafka.StartOffsetNewest
		log.Infof("subscribing to topic [%s]", tTopic)
		Consumer, err := broker.Consumer(conf2)
		if err != nil {
			log.Errorf("cannot create kafka consumer for %s:%d: %s", tTopic, tPartition, err)
		}
		subscription := kafkaSubscription{tTopic, tPartition, Consumer}
		Subscriptions = append(Subscriptions, subscription)
	}
	// run the message receiver
	go RunReceiver(t, Subscriptions )

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *kafkaTrigger) Stop() error {

	broker.Close()

	return nil
}

// RunAction starts a new Process Instance
func (t *kafkaTrigger) RunAction(actionId string, payload string, topic string, partition int32) {

	log.Debug("Starting new Process Instance")
	log.Debugf("Action Id: %s", actionId)
	log.Debugf("Payload: %s", payload)
	log.Debugf("Actual Topic: %s ", topic)

	req := t.constructStartRequest(payload, topic)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionId)

	context := trigger.NewContext(context.Background(), startAttrs)

	_, replyData, err := t.runner.Run(context, action, actionId, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s]", actionId)
	log.Debugf("Reply data: [%s]", replyData)

}

func (t *kafkaTrigger) publishMessage(topic string, partition int32, message string) {

	log.Debugf("ReplyTo topic: %s", topic)
	log.Debugf("Publishing message: %s", message)

	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: []byte(message)}

	log.Info("Sending message to Kafka server")
	resp, err := producer.Produce(topic, partition, msg)

	if err != nil {
		log.Errorf("Error sending message to Kafka broker: %s", err)
	}

	log.Debugf("Response: %s", resp)
	log.Debug("Message sent succesfully")
}

func (t *kafkaTrigger) constructStartRequest(message string, topic string) *StartRequest {

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

func RunReceiver(t *kafkaTrigger, Subscriptions []kafkaSubscription) error {
	for {
		for _, subscription := range Subscriptions {
			msg, err := subscription.consumer.Consume()
			if err != nil {
				if err != kafka.ErrNoData {
					log.Errorf("cannot consume %q topic message: %s", subscription.topic, err)
				}
				break
			}
			message := convert(msg.Value)
			log.Infof("Received message: %d on topic '%s', partition %d [%s]", msg.Offset, subscription.topic, subscription.partition, message)
			actionId, found := t.topicToActionId[subscription.topic]
			if found {
				log.Debugf("Found actionId: %s", actionId)
				t.RunAction(actionId, message, subscription.topic, subscription.partition)
			} else {
				log.Debug("actionId not found")
			}
		}
	}
	return nil
}

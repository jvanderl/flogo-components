package kafka

import (
	"context"
	//	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/op/go-logging"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"strconv"
	//	"bytes"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-kafka")
var broker kafka.Broker

// todo: switch to use endpoint registration

// KafkaTrigger is simple Kafka trigger
type KafkaTrigger struct {
	metadata          *trigger.Metadata
	runner            action.Runner
	settings          map[string]string
	config            *trigger.Config
	topicToActionURI  map[string]string
	topicToActionType map[string]string
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&KafkaTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *KafkaTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *KafkaTrigger) Init(config *trigger.Config, runner action.Runner) {

	t.config = config
	t.settings = config.Settings
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *KafkaTrigger) Start() error {

	ifServers := []string{t.settings["server"]}
	log.Debug("ifServers: ", ifServers)

	ifConfigID := t.settings["configid"]
	ifTopic := t.settings["topic"]
	cxPartition, err := strconv.ParseInt(t.settings["partition"], 10, 32)
	ifPartition := int32(0)
	if err == nil {
		ifPartition = int32(cxPartition)
	}

	conf := kafka.NewBrokerConf(ifConfigID)
	conf.AllowTopicCreation = true

	// connect to kafka cluster
	log.Debug("Connecting to Kafka server")
	broker, err := kafka.Dial(ifServers, conf)
	if err != nil {
		log.Fatalf("cannot connect to kafka cluster: %s", err)
		return err
	}
	defer broker.Close()
	log.Debug("Connected to Kafka server")

	conf2 := kafka.NewConsumerConf(ifTopic, ifPartition)
	conf2.StartOffset = kafka.StartOffsetNewest
	log.Debug("subscribing to topic ", ifTopic)
	consumer, err := broker.Consumer(conf2)
	if err != nil {
		log.Fatalf("cannot create kafka consumer for %s:%d: %s", ifTopic, ifPartition, err)
	}
	log.Debug("Subscription successful", ifTopic)

	t.topicToActionType = make(map[string]string)
	t.topicToActionURI = make(map[string]string)

	for _, endpoint := range t.config.Endpoints {
		t.topicToActionURI[endpoint.Settings["topic"]] = endpoint.ActionURI
		t.topicToActionType[endpoint.Settings["topic"]] = endpoint.ActionType
	}

	for {
		msg, err := consumer.Consume()
		if err != nil {
			if err != kafka.ErrNoData {
				log.Debug("cannot consume %q topic message: %s", ifTopic, err)
			}
			break
		}
		message := convert(msg.Value)
		log.Debug("Received message", msg.Offset, message)
		actionType, found := t.topicToActionType[ifTopic]
		actionURI, _ := t.topicToActionURI[ifTopic]
		if found {
			log.Debug("Found actionType", actionType)
			log.Debug("Found actionURI", actionURI)
			t.RunAction(actionType, actionURI, message, ifTopic, ifPartition)
		} else {
			log.Debug("actionType and URI not found")
		}
	}
	log.Debug("consumer quit")

	/*opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {

		topic := msg.Topic()
		//TODO we should handle other types, since mqtt message format are data-agnostic
		payload := string(msg.Payload())
		log.Debug("Received msg:", payload)
		log.Debug("Actual topic: ", topic)

		// try topic without wildcards
		actionType, found := t.topicToActionType[topic]
		actionURI, _ := t.topicToActionURI[topic]

		if found {
			t.RunAction(actionType, actionURI, payload, topic)
		} else {
			// search for wildcards

			for _, endpoint := range t.config.Endpoints {
				eptopic := endpoint.Settings["topic"]
				if strings.HasSuffix(eptopic, "/#") {
					// is wildcard, now check actual topic starts with wildcard
					if strings.HasPrefix(topic, strings.TrimSuffix(eptopic, "/#")) {
						// Got a match, now get the action for the wildcard topic
						actionType, found := t.topicToActionType[eptopic]
						actionURI, _ := t.topicToActionURI[eptopic]
						if found {
							t.RunAction(actionType, actionURI, payload, topic)
						}
					}
				}
			}
		}

	}
	)*/

	/*	client := mqtt.NewClient(opts)
		t.client = client
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		i, err := strconv.Atoi(t.settings["qos"])
		if err != nil {
			log.Error("Error converting \"qos\" to an integer ", err.Error())
			return err
		}

		t.topicToActionType = make(map[string]string)
		t.topicToActionURI = make(map[string]string)

		for _, endpoint := range t.config.Endpoints {
			if token := t.client.Subscribe(endpoint.Settings["topic"], byte(i), nil); token.Wait() && token.Error() != nil {
				log.Errorf("Error subscribing to topic %s: %s", endpoint.Settings["topic"], token.Error())
				panic(token.Error())
			} else {
				t.topicToActionURI[endpoint.Settings["topic"]] = endpoint.ActionURI
				t.topicToActionType[endpoint.Settings["topic"]] = endpoint.ActionType
			}
		} */

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *KafkaTrigger) Stop() error {

	broker.Close()

	return nil
}

// RunAction starts a new Process Instance
func (t *KafkaTrigger) RunAction(actionType string, actionURI string, payload string, topic string, partition int32) {

	log.Debug("Starting new Process Instance")
	log.Debug("Action Type: ", actionType)
	log.Debug("Action URI: ", actionURI)
	log.Debug("Payload: ", payload)
	log.Debug("Actual Topic: ", topic)

	req := t.constructStartRequest(payload, topic)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionType)

	context := trigger.NewContext(context.Background(), startAttrs)

	_, replyData, err := t.runner.Run(context, action, actionURI, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s-%s]", actionType, actionURI)
	log.Debug("Reply data: ", replyData)

	/*	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			log.Error(err)
		} else {
			t.publishMessage(req.ReplyTo, partition, string(data))
		}
	}*/
}

func (t *KafkaTrigger) publishMessage(topic string, partition int32, message string) {

	log.Debug("ReplyTo topic: ", topic)
	log.Debug("Publishing message: ", message)

	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: []byte(message)}

	log.Debug("Sending message to Kafka server")
	resp, err := producer.Produce(topic, partition, msg)

	if err != nil {
		log.Error("Error sending message to Kafka broker:", err)
	}

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Response:", resp)
	}
	log.Debug("Message sent succesfully")
}

func (t *KafkaTrigger) constructStartRequest(message string, topic string) *StartRequest {

	log.Debug("Received contstruct start request")

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

package mqtt

import (
	"context"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

var logger log.Logger

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Trigger struct {
	settings *Settings
	id       string
	conn     *MqttConnection
	handlers []*Handler
}

type Handler struct {
	handler trigger.Handler
	topic   string
}

type Factory struct {
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, settings: s}, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	logger = ctx.Logger()

	var err error
	t.conn, err = getMqttConnection(ctx.Logger(), t.settings)
	if t.conn.client.IsConnected() {
		logger.Debugf("Client is already connected")
	} else {
		logger.Debugf("MQTT Publisher connecting")
		if token := t.conn.client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	for _, handler := range ctx.GetHandlers() {
		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		err = t.subscribeTopic(s.Topic, s.Qos, t.conn.client)
		if err != nil {
			return err
		}
		t.handlers = append(t.handlers, &Handler{handler: handler, topic: s.Topic})
	}

	return err
}

// Start starts the mqtt trigger
func (t *Trigger) Start() error {

	go t.readMessages()

	return nil
}

// Stop stops the mqtt trigger
func (t *Trigger) Stop() error {

	for _, handler := range t.handlers {
		logger.Debugf("Unsubscribing from topic [%s]", handler.topic)
		if token := t.conn.client.Unsubscribe(handler.topic); token.Wait() && token.Error() != nil {
			logger.Errorf("Error unsubscribing from topic %s: %s", handler.topic, token.Error())
		} else {
			logger.Infof("Unsubscribed from topic %s", handler.topic)
		}
	}
	_ = t.conn.Stop
	return nil
}

func (t *Trigger) readMessages() {
	for {
		incoming := <-newMsg
		topic := incoming[0]
		message := incoming[1]
		logger.Infof("Received topic: %s, Message: %s\n", topic, message)

		for _, val := range t.handlers {
			if strings.HasSuffix(val.topic, "/#") {
				// is wildcard, now check actual topic starts with wildcard
				if strings.HasPrefix(topic, strings.TrimSuffix(val.topic, "/#")) {
					output := &Output{}
					output.Message = message
					output.ActualTopic = topic
					_, err := val.handler.Handle(context.Background(), output.ToMap())
					if err != nil {
						logger.Errorf("Error calling trigger action: %v", err)
					}
				}
			} else {
				// no wildcard, chech exact topic match
				if topic == val.topic {
					output := &Output{}
					output.Message = message
					output.ActualTopic = topic
					_, err := val.handler.Handle(context.Background(), output.ToMap())
					if err != nil {
						logger.Errorf("Error calling trigger action: %v", err)
					}
				}
			}
		}

	}
}

func (t *Trigger) subscribeTopic(topic string, qos int, client mqtt.Client) error {

	logger.Debugf("Subscribing to topic [%s]", topic)

	if token := client.Subscribe(topic, byte(qos), nil); token.Wait() && token.Error() != nil {
		logger.Errorf("Error subscribing to topic %s: %s", topic, token.Error())
		return (token.Error())
	} else {
		logger.Infof("Subscribed to topic %s", topic)
	}

	return nil
}

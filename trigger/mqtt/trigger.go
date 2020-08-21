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
		//t.handlers = append(t.handlers, handler)
		//err = t.subscribeTopic(ctx.Logger(), handler, t.conn.client)
		err = t.subscribeTopic(s.Topic, s.Qos, t.conn.client)
		if err != nil {
			return err
		}
		t.handlers = append(t.handlers, &Handler{handler: handler, topic: s.Topic})
		//t.mqttEvents = append(t.mqttEvents, registerMqttEventHandler(s.Topic, newActionHandler(handler)))
	}

	return err
}

// Start starts the mqtt trigger
func (t *Trigger) Start() error {

	go t.readMessages()

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *Trigger) Stop() error {

	for _, handler := range t.handlers {
		if token := t.conn.client.Unsubscribe(handler.topic); token.Wait() && token.Error() != nil {
			logger.Errorf("Error unsubscribing from topic %s: %s", handler.topic, token.Error())
		}
	}
	_ = t.conn.Stop
	return nil
}

/*
func (t *Trigger) Stop() error {

	for _, handler := range t.handlers {
		handlerSetting := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), handlerSetting, true)
		if err == nil {
			if token := t.conn.client.Unsubscribe(handlerSetting.Topic); token.Wait() && token.Error() != nil {
				//log.Errorf("Error unsubscribing from topic %s: %s", handlerCfg.GetSetting("topic"), token.Error())
			}
		}
	}
	_ = t.conn.Stop
	return nil
}
*/
/*
func (t *Trigger) readMessages() {
	for {
		incoming := <-newMsg
		topic := incoming[0]
		message := incoming[1]
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", topic, message)

		for _, val := range t.mqttEvents {
			go val(topic, message)
		}

	}
}
*/
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

/*
func (t *Trigger) subscribeTopic(logger log.Logger, handler trigger.Handler, client mqtt.Client) error {

	handlerSetting := &HandlerSettings{}
	err := metadata.MapToStruct(handler.Settings(), handlerSetting, true)
	if err != nil {
		return err
	}

	if handlerSetting.Topic == "" {
		return fmt.Errorf("topic string was not provided for handler: [%s]", handler)
	}

	logger.Debugf("Subscribing to topic [%s]", handlerSetting.Topic)

	if token := client.Subscribe(handlerSetting.Topic, byte(handlerSetting.Qos), nil); token.Wait() && token.Error() != nil {
		logger.Errorf("Error subscribing to topic %s: %s", handlerSetting.Topic, token.Error())
		return (token.Error())
	} else {
		logger.Infof("Subscribed to topic %s", handlerSetting.Topic)
		//t.topicToHandler[handlerSetting.Topic] = handler
	}

	return nil
}
*/
/*
// registerMqttEventHandler is used for mqtt event handler registration
func registerMqttEventHandler(topic string, onEvent mqttOnEvent) mqttOnEvent {
	//ignore
	return onEvent
}

// mqttOnEvent is an mqtt event handler for our mqtt event source
type mqttOnEvent func(actualTopic string, message string)

func newActionHandler(handler trigger.Handler) mqttOnEvent {

	return func(actualTopic string, message string) {
		handlerSetting := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), handlerSetting, true)
		if err == nil {
			if strings.HasSuffix(handlerSetting.Topic, "/#") {
				// is wildcard, now check actual topic starts with wildcard
				if strings.HasPrefix(actualTopic, strings.TrimSuffix(handlerSetting.Topic, "/#")) {
					output := &Output{}
					output.Message = message
					output.ActualTopic = actualTopic
					_, err := handler.Handle(context.Background(), output.ToMap())
					if err != nil {
						//handle error
					}
				}
			}
		}
	}
}
*/

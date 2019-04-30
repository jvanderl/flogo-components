package amqp

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/streadway/amqp"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-amqp")

const (
	ivServer   = "server"
	ivPort     = "port"
	ivUserID   = "userID"
	ivPassword = "password"
	//	ivExchange    = "exchange"
	//	ivRoutingKey  = "routingKey"
	//	ivRoutingType = "routingType"
	//	ivMessage     = "message"
	//	ivDurable     = "durable"
	//	ivAutoDelete  = "autoDelete"
	//	ivExclusive   = "exclusive"
	//	ivNoWait      = "noWait"

	ovMessage = "message"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	handler *trigger.Handler
}

var Consumers []Consumer

// AMQPTriggerFactory AMQP Trigger factory
type AMQPTriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &AMQPTriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *AMQPTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &AMQPTrigger{metadata: t.metadata, config: config}
}

// AMQPTrigger is a stub for your Trigger implementation
type AMQPTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	handlers []*trigger.Handler
}

// Init implements ext.Trigger.Init
func (t *AMQPTrigger) Init(runner action.Runner) {
	log.Debug("Trigger Init called")
	t.runner = runner
	//	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Initialize implements trigger.Init.Initialize
func (t *AMQPTrigger) Initialize(ctx trigger.InitContext) error {
	log.Debug("Trigger Initialize called")
	t.handlers = ctx.GetHandlers()

	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *AMQPTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *AMQPTrigger) Start() error {
	// start the trigger
	serverName := t.config.GetSetting(ivServer)
	serverPort := t.config.GetSetting(ivPort)
	userID := t.config.GetSetting(ivUserID)
	password := t.config.GetSetting(ivPassword)

	uri := "amqp://" + userID + ":" + password + "@" + serverName + ":" + serverPort + "/"
	log.Info("Adding Listeners")
	for i, handler := range t.config.Handlers {
		log.Info("Creating new Consumer...")
		exchange := handler.GetSetting("exchange")
		exchangeType := handler.GetSetting("exchangeType")
		queueName := handler.GetSetting("queueName")
		bindingKey := handler.GetSetting("bindingKey")
		consumerTag := handler.GetSetting("consumerTag")

		consumer, err := t.NewConsumer(uri, exchange, exchangeType, queueName, bindingKey, consumerTag, t.handlers[i])
		if err != nil {
			log.Error("Error Creating Consumer")
			return err
		}

		Consumers = append(Consumers, *consumer)
	}

	select {}

	return nil

}

// Stop implements trigger.Trigger.Start
func (t *AMQPTrigger) Stop() error {
	// stop the trigger

	log.Infof("shutting down")
	for _, c := range Consumers {
		if err := c.Shutdown(); err != nil {
			log.Errorf("error during shutdown: %s", err)
		}
	}

	return nil
}

func (t *AMQPTrigger) NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string, handler *trigger.Handler) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
		handler: handler,
	}

	var err error

	log.Infof("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Infof("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Infof("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Infof("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go t.handle(deliveries, c.done, handler)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Infof("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func (t *AMQPTrigger) handle(deliveries <-chan amqp.Delivery, done chan error, handler *trigger.Handler) {
	for d := range deliveries {
		log.Infof(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		log.Infof("Content Type: %s ", d.ContentType)
		log.Infof("Routing Key: %s ", d.RoutingKey)
		t.Execute(handler, d.Body, d.ContentType, d.RoutingKey)
		//		t.RunAction(actionID, "tesmessage", d.ContentType, d.RoutingKey)
		//		t.RunAction(actionID, string(d.Body), d.ContentType, d.RoutingKey)
		d.Ack(false)
	}
	//**** TODO add actual response runaction here ****

	log.Infof("handle: deliveries channel closed")
	done <- nil
}

// Execute executes any handlers defined immediately on startup
func (t *AMQPTrigger) Execute(handler *trigger.Handler, payload []byte, contentType string, routingKey string) {
	log.Debug("Starting process")

	triggerData := map[string]interface{}{
		"message":     payload,
		"contentType": contentType,
		"routingKey":  routingKey,
	}

	response, err := handler.Handle(context.Background(), triggerData)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
	} else {
		log.Debugf("Action call successful: %v", response)
	}
}

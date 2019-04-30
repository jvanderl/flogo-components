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
	conn     *amqp.Connection
	channel  *amqp.Channel
	tag      string
	done     chan error
	actionID string
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
	for _, handler := range t.config.Handlers {
		log.Info("Creating new Consumer...")
		exchange := handler.GetSetting("exchange")
		exchangeType := handler.GetSetting("exchangeType")
		queueName := handler.GetSetting("queueName")
		bindingKey := handler.GetSetting("bindingKey")
		consumerTag := handler.GetSetting("consumerTag")
		actionID := handler.ActionId

		consumer, err := t.NewConsumer(uri, exchange, exchangeType, queueName, bindingKey, consumerTag, actionID)
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

func (t *AMQPTrigger) NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string, actionID string) (*Consumer, error) {
	c := &Consumer{
		conn:     nil,
		channel:  nil,
		tag:      ctag,
		done:     make(chan error),
		actionID: actionID,
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

	go t.handle(deliveries, c.done, actionID)

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

func (t *AMQPTrigger) handle(deliveries <-chan amqp.Delivery, done chan error, actionID string) {
	for d := range deliveries {
		log.Infof(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		t.RunAction(actionID, d.Body, d.ContentType, d.RoutingKey)
		d.Ack(false)
	}
	//**** TODO add actual response runaction here ****

	log.Infof("handle: deliveries channel closed")
	done <- nil
}

// RunAction starts a new Process Instance
func (t *AMQPTrigger) RunAction(actionID string, payload []byte, contentType string, routingKey string) {

	log.Debug("Starting new Process Instance")
	log.Debugf("Action Id: %s", actionID)
	log.Debugf("Payload: %s", payload)
	log.Debugf("Content Type: %s ", contentType)
	log.Debugf("Routing Key: %s ", routingKey)

	req := t.constructStartRequest(payload, contentType, routingKey)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionID)

	context := trigger.NewContext(context.Background(), startAttrs)

	_, replyData, err := t.runner.Run(context, action, actionID, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s]", actionID)
	log.Debugf("Reply data: [%s]", replyData)

}

func (t *AMQPTrigger) constructStartRequest(message []byte, contentType string, routingKey string) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	data["contentType"] = contentType
	data["routingKey"] = routingKey
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI string                 `json:"flowUri"`
	Data       map[string]interface{} `json:"data"`
	ReplyTo    string                 `json:"replyTo"`
}

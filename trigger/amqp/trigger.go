package amqp

import (
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

type AMQPConsumer struct {
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
	tag     string
	done    chan error
}

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
	config   *trigger.Config
	handlers []*trigger.Handler
}

// Init implements ext.Trigger.Init
func (t *AMQPTrigger) Init(runner action.Runner) {
	log.Debug("Trigger Init called")
	//	t.runner = runner
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

	amqpURI := "amqp://" + userID + ":" + password + "@" + serverName + ":" + serverPort + "/"
	log.Infof("Connecting to AMQP server %s...", amqpURI)
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Error("Error connecting to AMQP server")
		return err
	}
	log.Info("Connected")
	defer conn.Close()

	// Create array of channels for all handlers

	var Consumers []AMQPConsumer

	log.Info("Adding Listeners")
	for i, handler := range t.config.Handlers {
		log.Info("Creating new Consumer...")
		tag := "tag" + string(i)
		consumer := AMQPConsumer{nil, nil, tag, make(chan error)}
		log.Info("Creating new Channel...")
		consumer.channel, err = conn.Channel()
		if err != nil {
			log.Error("Error Creating Channel")
			return err
		}
		defer consumer.channel.Close()
		queueName := handler.GetSetting("queueName")
		log.Infof("Declaring queue %s...", queueName)
		q, err := consumer.channel.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Error("Failed to declare a queue")
			return err
		}
		log.Info("Creating consumer for queue")
		consumer.msgs, err = consumer.channel.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Error("Failed to register a consumer")
			return err
		}
		Consumers = append(Consumers, consumer)
	}
	log.Info("Done. Starting listener.")

	forever := make(chan bool)

	go func() {
		for i := range Consumers {
			for d := range Consumers[i].msgs {
				log.Infof("Received AMQP Message: %s", d.Body)
				//**** TODO add actual response runaction here ****
			}
		}
	}()

	<-forever

	return nil

}

// Stop implements trigger.Trigger.Start
func (t *AMQPTrigger) Stop() error {
	// stop the trigger
	return nil
}

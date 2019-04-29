package amqp

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/streadway/amqp"
)

// log is the default package logger
var log = logger.GetLogger("activity-jl-amqp")

const (
	ivServer      = "server"
	ivPort        = "port"
	ivUserID      = "userID"
	ivPassword    = "password"
	ivExchange    = "exchange"
	ivRoutingKey  = "routingKey"
	ivRoutingType = "routingType"
	ivMessage     = "message"
	ivDurable     = "durable"
	ivAutoDelete  = "autoDelete"
	ivExclusive   = "exclusive"
	ivNoWait      = "noWait"

	ovResult = "result"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	serverName := context.GetInput(ivServer).(string)
	serverPort := context.GetInput(ivPort).(string)
	userID := context.GetInput(ivUserID).(string)
	password := context.GetInput(ivPassword).(string)
	message := context.GetInput(ivMessage).(string)
	exchange := context.GetInput(ivExchange).(string)
	routingKey := context.GetInput(ivRoutingKey).(string)
	routingType := context.GetInput(ivRoutingType).(string)
	durable := context.GetInput(ivDurable).(bool)
	autoDelete := context.GetInput(ivAutoDelete).(bool)
	exclusive := context.GetInput(ivExclusive).(bool)
	noWait := context.GetInput(ivNoWait).(bool)

	amqpURI := "amqp://" + userID + ":" + password + "@" + serverName + ":" + serverPort + "/"
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		context.SetOutput(ovResult, "ERR_CONNECT_AMQP")
		return false, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		context.SetOutput(ovResult, "ERR_OPENING_CHANNEL")
		return false, err
	}
	defer ch.Close()

	if routingType == "queue" {
		//expect publish to queue
		_, err := ch.QueueDeclare(
			routingKey, // name
			durable,    // durable
			autoDelete, // delete when unused
			exclusive,  // exclusive
			noWait,     // no-wait
			nil,        // arguments
		)
		if err != nil {
			context.SetOutput(ovResult, "ERR_DECLARE_QUEUE")
			return false, err
		}
	} else {
		//expect publish to topic
		err = ch.ExchangeDeclare(
			exchange,   // name
			"topic",    // type
			durable,    // durable
			autoDelete, // auto-deleted
			exclusive,  // internal
			noWait,     // no-wait
			nil,        // arguments
		)
		if err != nil {
			context.SetOutput(ovResult, "ERR_DECLARE_EXCHANGE")
			return false, err
		}
	}

	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	log.Debugf(" [x] Sent %s", message)
	if err != nil {
		context.SetOutput(ovResult, "ERR_PUBLISH_MSG")
		return false, err
	}
	context.SetOutput(ovResult, "OK")
	return true, nil
}

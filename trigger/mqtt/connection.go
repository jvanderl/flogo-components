package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/project-flogo/core/support/log"
)

var connections = make(map[string]*MqttConnection)

var newMsg = make(chan [2]string)

type MqttConnection struct {
	client mqtt.Client
}

func (c *MqttConnection) Connection() mqtt.Client {
	return c.client
}

func (c *MqttConnection) Stop() error {
	c.client.Disconnect(100)
	return nil
}

func getConnectionKey(settings *Settings) string {

	var connKey string

	connKey += settings.Broker
	if settings.User != "" {
		connKey += settings.User
	}

	return connKey
}

func getMqttConnection(logger log.Logger, settings *Settings) (*MqttConnection, error) {

	connKey := getConnectionKey(settings)

	if conn, ok := connections[connKey]; ok {
		logger.Debugf("Reusing cached MQTT connection [%s]", connKey)
		return conn, nil
	}

	newConn := &MqttConnection{}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(settings.Broker)
	opts.SetClientID(settings.Id)
	opts.SetUsername(settings.User)
	opts.SetPassword(settings.Password)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		newMsg <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(opts)
	newConn.client = client

	connections[connKey] = newConn
	logger.Debugf("Caching MQTT connection [%s]", connKey)

	return newConn, nil
}

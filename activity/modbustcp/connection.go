package modbustcp

import (
	"time"

	"github.com/goburrow/modbus"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
)

var connections = make(map[string]*ModbusTcpConnection)

type ModbusTcpConnection struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

func (c *ModbusTcpConnection) Connection() modbus.Client {
	return c.client
}

func (c *ModbusTcpConnection) Stop() error {
	c.handler.Close()
	return nil
}

func getConnectionKey(settings *Settings) string {

	var connKey string

	connKey += settings.Server
	connKey += string(settings.SlaveID)

	return connKey
}

func getModbusTcpConnection(logger log.Logger, settings *Settings) (*ModbusTcpConnection, error) {

	connKey := getConnectionKey(settings)

	if conn, ok := connections[connKey]; ok {
		logger.Debugf("Reusing cached Modbus TCP connection [%s]", connKey)
		return conn, nil
	}

	newConn := &ModbusTcpConnection{}

	mbhandler := modbus.NewTCPClientHandler(settings.Server)
	mbhandler.Timeout = time.Duration(settings.Timeout) * time.Second
	tmpInt, err := coerce.ToInt(settings.SlaveID)
	if err != nil {
		return nil, err
	}
	mbhandler.SlaveId = byte(tmpInt)
	//handler.Logger = logger

	err = mbhandler.Connect()
	if err != nil {
		logger.Debugf("Error connecting to Modbus Server: %v\n", err)
		return nil, err
	} else {
		logger.Debugf("Connected to Modbus Server\n")
		mbclient := modbus.NewClient(mbhandler)
		newConn.handler = mbhandler
		newConn.client = mbclient
	}

	connections[connKey] = newConn
	logger.Debugf("Caching Modbus TCP connection [%s]", connKey)

	return newConn, nil
}

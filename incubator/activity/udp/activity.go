package udp

import (
	"net"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-jl-udp")

const (
	ivMessage        = "message"
	ivMulticastGroup = "multicast_group"
	ivPort           = "port"

	ovResult = "result"
)

// MyActivity is an Activity that is used to send a UDP Message
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new MyActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Sends UDP Message
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	mGroup := context.GetInput(ivMulticastGroup).(string)
	log.Infof("multicast_group: %v", mGroup)
	mPort := context.GetInput(ivPort).(int)
	log.Infof("port: %v", strconv.Itoa(mPort))
	message := context.GetInput(ivMessage).(string)
	serverStr := mGroup + ":" + strconv.Itoa(mPort)
	log.Infof("ServerStr: %v", serverStr)
	localStr := mGroup + ":0"
	log.Infof("LocalStr: %v", localStr)

	serverAddr, err := net.ResolveUDPAddr("udp", serverStr)
	if err != nil {
		log.Errorf("Error Resolving Server Address: %v", err)
		context.SetOutput(ovResult, "ERR_RESOLVE_SERVER_ADDRESS")
		return false, nil
	}
	localAddr, err := net.ResolveUDPAddr("udp", localStr)
	if err != nil {
		log.Errorf("Error Resolving Local Address: %v", err)
		context.SetOutput(ovResult, "ERR_RESOLVE_LOCAL_ADDRESS")
		return false, nil
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Errorf("Error Connecting to UDP: %v", err)
		context.SetOutput(ovResult, "ERR_UDP_CONNECT")
		return false, nil
	}

	defer conn.Close()
	buf := []byte(message)
	_, err = conn.Write(buf)
	if err != nil {
		log.Errorf("Error Sending UDP Message: %v", err)
		context.SetOutput(ovResult, "ERR_UDP_SEND")
		return false, nil
	}

	context.SetOutput(ovResult, "OK")

	return true, nil
}

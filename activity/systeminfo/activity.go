package systeminfo

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"net"
	"os"
)

var log = logger.GetLogger("activity-jvanderl-eftl")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	ipaddrs := ""
	host, err := os.Hostname()
	log.Debug("Getting hostname")

	if err == nil {
		context.SetOutput("hostname", host)
		log.Debugf("Hostname is: [%s]", host)
		log.Debug("Getting IP address")
		addrs, _ := net.LookupIP(host)
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				ipaddrs = addr.String()
			}
		}
		context.SetOutput("ipaddress", ipaddrs)
		log.Debugf("IP address is: [%s]", ipaddrs)
	} else {
		context.SetOutput("hostname", "Unknown")
		context.SetOutput("ipaddress", "Unknown")
	}

	return true, nil
}

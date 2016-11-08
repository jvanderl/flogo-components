package systeminfo

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"os"
	"net"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&MyActivity{metadata: md})
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
	if err == nil {
		context.SetOutput("hostname", host)
		addrs, _ := net.LookupIP(host)
		for _, addr := range addrs {
    		if ipv4 := addr.To4(); ipv4 != nil {
        		 ipaddrs = addr.String()
    		}   
		}
		context.SetOutput("ipaddress", ipaddrs)
	} else {
		context.SetOutput("hostname", "Unknown")
		context.SetOutput("ipaddress", "Unknown")
	}

	return true, nil
}

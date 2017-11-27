package systeminfo

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"net"
	"os"
	"strings"
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
  includeNetmask := context.GetInput("includenetmask").(bool )

	ip4addrs := ""
	ip6addrs := ""
	macaddrs := ""
	gotIt := false
	host, err := os.Hostname()
	log.Debug("Getting hostname")

	if err == nil {
		context.SetOutput("hostname", host)
		log.Debugf("Hostname is: [%s]", host)
		log.Debug("Getting IP address")
/*		addrs, _ := net.LookupIP(host)
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				ipaddrs = addr.String()
			}
		}
*/
		interfaces, _ := net.Interfaces()
		for _, inter := range interfaces {
			log.Debugf("Found interface: %v, hwAdress: %v",inter.Name, inter.HardwareAddr)
			if addrs, err := inter.Addrs(); err == nil {
				macaddrs = inter.HardwareAddr.String()
				for _, addr := range addrs {
					log.Debugf("Interface: %v, Found IpAddress: %v", inter.Name, addr)
//					if ipv4 := addr.To4(); ipv4 != nil {
					if (strings.Contains(addr.String(), ":")) {
						ip6addrs = addr.String()
					} else {
						ip4addrs = addr.String()
						if (macaddrs!="") {
							gotIt = true
						}
					}
				}
				if (gotIt) {break}
			}
		}
		if (gotIt && includeNetmask == false) {
			ip4addrs = before(ip4addrs, "/")
			ip6addrs = before(ip6addrs, "/")
		}

		context.SetOutput("ipaddress", ip4addrs)
		context.SetOutput("ip6address", ip6addrs)
		context.SetOutput("macaddress", macaddrs)
		log.Debugf("IP v4 address is: [%s]", ip4addrs)
		log.Debugf("IP v6 address is: [%s]", ip6addrs)
		log.Debugf("Mac address is: [%s]", macaddrs)
	} else {
		context.SetOutput("hostname", "Unknown")
		context.SetOutput("ipaddress", "Unknown")
		context.SetOutput("ip6address", "Unknown")
		context.SetOutput("macaddress", "Unknown")
	}

	return true, nil
}

func before(value string, a string) string {
    // Get substring before a string.
    pos := strings.Index(value, a)
    if pos == -1 {
        return ""
    }
    return value[0:pos]
}

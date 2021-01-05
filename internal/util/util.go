package util

import (
	"io/ioutil"
	"net"
	"strings"
)

// DeadlineExceeded is the error returned if a timeout is specified.
var DeadlineExceeded error = deadlineExceededError{}

var hostHostnamePath = "/etc/host_hostname"

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string { return "deadline exceeded" }

// GetInterfaces returns the list of non-loopback interfaces
func GetInterfaces() (ips []net.IP, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	var addrs []net.Addr
	var ip net.IP

	for _, i := range ifaces {
		addrs, err = i.Addrs()
		if err != nil {
			return
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !ip.IsLoopback() {
				ips = append(ips, ip)
			}
		}
	}
	return
}

// HostHostname returns hostname from file defined in constant hostHostnamePath
// or the defaultValue in case the file does not exits.
func HostHostname(defaultValue string) string {
	if content, err := ioutil.ReadFile(hostHostnamePath); err == nil {
		return strings.Trim(string(content), "\n")
	}
	return defaultValue
}

// PanicOnError panics if the error is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

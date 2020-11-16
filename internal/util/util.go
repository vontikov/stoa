package util

import (
	"net"
)

// DeadlineExceeded is the error returned if a timeout is specified.
var DeadlineExceeded error = deadlineExceededError{}

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

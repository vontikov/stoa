package util

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strings"
	"time"
)

// ErrDeadlineExceeded is the error returned if a timeout is specified.
var ErrDeadlineExceeded = errors.New("deadline exceeded")

var hostHostnamePath = "/etc/host_hostname"

// DDNSBootstrap returns bootstrap string build with DNS lookup.
// XXX experimental, very basic
func DNSBootstrap(expectedClusterSize int, timeout time.Duration, sep string) (peers string, err error) {
	if expectedClusterSize < 3 || expectedClusterSize%2 == 0 {
		err = errors.New("cluster size must be odd number greater or equal to 3")
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	ifaces, err := GetInterfaces()
	if err != nil {
		return
	}

	t := time.After(timeout)
	for {
		select {
		case <-t:
			err = ErrDeadlineExceeded
			return
		default:
			var ips []net.IP
			ips, err = net.LookupIP(hostname)
			if err != nil {
				return
			}
			if len(ips) != expectedClusterSize {
				break
			}

			// the lowest will be the leader
			sort.Slice(ips, func(i, j int) bool { return bytes.Compare(ips[i], ips[j]) < 0 })
			for _, iface := range ifaces {
				for i, e := range ips {
					if bytes.Equal(e.To4(), iface.To4()) {
						if i == 0 {
							b := strings.Builder{}
							for _, ip := range ips {
								b.WriteString(ip.To4().String())
								b.WriteString(sep)
							}
							return strings.TrimSuffix(b.String(), sep), nil
						}
						return iface.To4().String(), nil
					}
				}
			}
		}
	}
}

// GetInterfaces returns the list of non-loopback interfaces
var GetInterfaces = func() (ips []net.IP, err error) {
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

package util

import (
	"net"
	"time"

	"github.com/vontikov/stoa/internal/cluster"
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

// WaitLeaderStatus waits until the peer p becomes the leader or timeout
// specified by the duration d happens.
// Returns nil on success, otherwise returns Deadline.
func WaitLeaderStatus(p *cluster.Peer, d time.Duration) error {
	const sleepDuration = 100 * time.Millisecond
	timeout := time.After(d)
	for {
		select {
		case <-timeout:
			return DeadlineExceeded
		default:
			if p.IsLeader() {
				return nil
			}
			time.Sleep(sleepDuration)
		}
	}
}

// WaitForSingleLeader waits until only one of the peers becomes a leader or
// timeout specified by the durationd happens.
// Returns nil on success, otherwise returns Deadline.
func WaitForSingleLeader(peers []*cluster.Peer, d time.Duration) error {
	const sleepDuration = 1000 * time.Millisecond
	timeout := time.After(d)
	a := 0
	for {
		select {
		case <-timeout:
			return DeadlineExceeded
		default:
			n := 0
			for _, p := range peers {
				if p.IsLeader() {
					n++
				}
			}
			if n == 1 {
				a++
				if a == 5 {
					return nil
				}
			}
			time.Sleep(sleepDuration)
		}
	}
}

// LeaderIndex scans the peers and returns the index of the first leader if
// found. Otherwise returns -1.
func LeaderIndex(peers []*cluster.Peer) int {
	for i, p := range peers {
		if p.IsLeader() {
			return i
		}
	}
	return -1
}

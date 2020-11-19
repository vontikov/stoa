package client

import (
	"time"
)

const (
	DefaultDialTimeout      = 30000 * time.Millisecond
	DefaultDiscoveryEnabled = true
	DefaultDiscoveryIP      = "224.5.1.1"
	DefaultDiscoveryPort    = 3500
	DefaultKeepAlivePeriod  = 1000 * time.Millisecond
	DefaultRetryTimeout     = 15000 * time.Millisecond
)

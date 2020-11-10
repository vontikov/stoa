package discovery

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/vontikov/stoa/internal/common"
	"github.com/vontikov/stoa/internal/logging"
)

// MessageSupplier produces discovery message on demand
type MessageSupplier func() ([]byte, error)

// Handler handles incoming discovery message received from src
type Handler func(src *net.UDPAddr, message []byte) error

type sender struct {
	conn   *net.UDPConn
	ms     MessageSupplier
	logger logging.Logger
}

type receiver struct {
	conn    *net.UDPConn
	handler Handler
	logger  logging.Logger
}

const (
	network = "udp"
)

var (
	maxDatagramSize = 1024
	duration        = 1000 * time.Millisecond
)

// Duration sets the duration between discovery ping messages
func Duration(d time.Duration) { duration = d }

func NewReceiver(ip string, port int, h Handler) (common.Runnable, error) {
	addr, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP(network, nil, addr)
	if err != nil {
		return nil, err
	}
	err = conn.SetReadBuffer(maxDatagramSize)
	if err != nil {
		return nil, err
	}
	return &receiver{
		conn:    conn,
		handler: h,
		logger:  logging.NewLogger("receiver"),
	}, nil
}

func (r *receiver) Run(ctx context.Context) error {
	buf := make([]byte, maxDatagramSize)
loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			r.conn.SetReadDeadline(time.Now().Add(duration))
			n, src, err := r.conn.ReadFromUDP(buf)
			if err != nil {
				if e, ok := err.(net.Error); !ok || !e.Timeout() {
					r.logger.Warn("Read error: ", err)
				}
				continue loop
			}
			if r.handler != nil {
				if err := r.handler(src, buf[:n]); err != nil {
					r.logger.Warn("Hadler error: ", err)
				}
			}
		}
	}
}

func NewSender(ip string, port int, ms MessageSupplier) (common.Runnable, error) {
	addr, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(network, nil, addr)
	if err != nil {
		return nil, err
	}

	return &sender{
		conn:   conn,
		ms:     ms,
		logger: logging.NewLogger("sender"),
	}, nil
}

func (s *sender) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := s.ms()
			if err != nil {
				s.logger.Warn("Message supplier error: ", err)
			}
			if _, err := s.conn.Write(msg); err != nil {
				s.logger.Warn("Write error: ", err)
			}
			time.Sleep(duration)
		}
	}
}

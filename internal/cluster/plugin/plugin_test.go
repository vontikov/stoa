package plugin

import (
	"testing"

	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"

	"github.com/vontikov/stoa/pkg/pb"
)

type dnsHandler struct {
	domains map[string]string
}

func newDNSHandler(service, host string) *dnsHandler {
	h := &dnsHandler{
		domains: make(map[string]string),
	}
	h.domains[service+"."] = host
	return h
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	if r.Question[0].Qtype == dns.TypeA {
		msg.Authoritative = true
		domain := msg.Question[0].Name
		if address, ok := h.domains[domain]; ok {
			msg.Answer = append(msg.Answer,
				&dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    60,
					},
					A: net.ParseIP(address),
				},
			)
		}
	}
	w.WriteMsg(&msg)
}

var errFsm = errors.New("test fsm error")

type testEvictionServer struct {
	pb.UnimplementedEvictionServer
}

func newTestEvictionServer() *testEvictionServer {
	return &testEvictionServer{}
}

func (*testEvictionServer) Dictionary(ctx context.Context, e *pb.DictionaryEntry) (*pb.Result, error) {
	return &pb.Result{Ok: true}, nil
}

type testFSM struct {
	c, max int
}

func (f *testFSM) DictionaryScan() chan *pb.DictionaryEntry {
	ch := make(chan *pb.DictionaryEntry)
	go func() {
		for i := 0; i < f.max; i++ {
			ch <- &pb.DictionaryEntry{}
		}
		close(ch)
	}()
	return ch
}

func (f *testFSM) DictionaryRemove(e *pb.DictionaryEntry) error {
	f.c++
	if f.c < f.max {
		return nil
	}
	return errFsm
}

func newTestFSM(max int) *testFSM {
	return &testFSM{max: max}
}

func TestPluginIntegration(t *testing.T) {
	const (
		dnsHost     = "127.0.0.1"
		dnsPort     = 1053
		serviceHost = "127.0.0.1"
		servicePort = 3701
		service     = "foo.bar"

		max = 10
	)

	TickerPeriod(200 * time.Millisecond)
	ErrorHandler(func(err error) bool { return false })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dnsServer := &dns.Server{
		Addr:    fmt.Sprintf(":%d", dnsPort),
		Net:     "udp",
		Handler: newDNSHandler(service, serviceHost),
	}

	grpcAddr := fmt.Sprintf("%s:%d", serviceHost, servicePort)
	lis, err := net.Listen("tcp", grpcAddr)
	assert.Nil(t, err)

	grpcServer := grpc.NewServer()
	pb.RegisterEvictionServer(grpcServer, newTestEvictionServer())

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error { return grpcServer.Serve(lis) })
	g.Go(dnsServer.ListenAndServe)

	target := fmt.Sprintf("dns://%s:%d/%s:%d", dnsHost, dnsPort, service, servicePort)
	fsm := newTestFSM(max)
	p := New(target, fsm)
	g.Go(func() error { return p.Run(ctx) })

	time.Sleep(1 * time.Second)
	dnsServer.Shutdown()
	grpcServer.GracefulStop()

	err = g.Wait()
	assert.Equal(t, max, fsm.c)
	assert.Equal(t, errFsm, err)
}

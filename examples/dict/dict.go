package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	stoa "github.com/vontikov/stoa/pkg/client"
)

const (
	bootstrap = "localhost:3001,localhost:3002,localhost:3003"
	dictName  = "dict"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go put(ctx)
	go get(ctx)

	<-signals
	cancel()
}

func put(ctx context.Context) {
	client, err := stoa.New(ctx, stoa.WithBootstrap(bootstrap))
	if err != nil {
		log.Fatal("client error: ", err)
	}

	d := client.Dictionary(dictName)

	t := time.NewTicker(1000 * time.Millisecond)
	defer t.Stop()

	n := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			var k, v bytes.Buffer
			fmt.Fprintf(&k, "key-%d", n)
			fmt.Fprintf(&v, "val-%d", n)

			if _, err := d.Put(ctx, k.Bytes(), v.Bytes()); err != nil {
				log.Fatal("put: ", err)
			}
			log.Println("put: ", string(k.Bytes()), string(v.Bytes()))
			n++
		}
	}
}

func get(ctx context.Context) {
	client, err := stoa.New(ctx, stoa.WithBootstrap(bootstrap))
	if err != nil {
		log.Fatal("client error: ", err)
	}

	d := client.Dictionary(dictName)

	t := time.NewTicker(1000 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			sz, _ := d.Size(ctx)
			log.Printf("size: %d\n", sz)

			kv, _ := d.Range(ctx)
			for e := range kv {
				k := e[0]
				v := e[1]
				log.Printf("get: %s %s\n", string(k), string(v))
			}
		}
	}
}

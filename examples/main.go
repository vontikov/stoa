package main

import (
	"context"
	"fmt"

	"github.com/vontikov/stoa/pkg/client"
)

func main() {
	const dictName = "dict"

	client, err := client.New()
	panicOnError(err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dict := client.Dictionary(dictName)
	sz, err := dict.Size(ctx)
	panicOnError(err)
	fmt.Println("Size: ", sz)

	err = dict.Remove(ctx, []byte("key"))
	panicOnError(err)

	r, err := dict.PutIfAbsent(ctx, []byte("key"), []byte("value"))
	panicOnError(err)
	fmt.Println("Result: ", r)

	r, err = dict.PutIfAbsent(ctx, []byte("key"), []byte("value"))
	panicOnError(err)
	fmt.Println("Result: ", r)

	sz, err = dict.Size(ctx)
	panicOnError(err)
	fmt.Println("Size: ", sz)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

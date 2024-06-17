package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/nirdosh17/tracer"
)

// go run cmd/main.go -hops 20 -timeout 120 google.com
func main() {
	hops := flag.Int("hops", 64, "max hops(TTL) for the packet, default: 64")
	timeout := flag.Int("timeout", 100, "timeout(ms) for ICMP response, default: 100")
	flag.Parse()

	host := flag.Arg(0)
	if len(host) == 0 {
		fmt.Println("no host provided!")
		os.Exit(1)
	}

	// consume live hops from channel
	var wg sync.WaitGroup
	wg.Add(1)

	c := make(chan tracer.Hop)
	go func() {
		for {
			hop, ok := <-c
			if !ok {
				// channel closed, so exiting
				wg.Done()
				return
			}
			fmt.Println("received hop", hop)
		}
	}()

	config := tracer.NewConfig().WithHops(*hops).WithTimeout(*timeout)
	t := tracer.NewTracer(config)
	err := t.Run(host, c)
	if err != nil {
		fmt.Println("Error from Trace runner: ", err)
	}
	wg.Wait()
}

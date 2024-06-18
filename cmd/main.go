package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/nirdosh17/tracer"
)

// go run cmd/main.go -hops 5 google.com
func main() {
	hops := flag.Int("hops", 64, "max hops(TTL) for the packet, default: 64")
	timeout := flag.Int("timeout", 5, "timeout(seconds) for ICMP response, default: 5")
	retries := flag.Int("retries", 2, "timeout(seconds) for ICMP response, default: 2")

	flag.Parse()

	host := flag.Arg(0)
	if len(host) == 0 {
		fmt.Println("no host provided!")
		os.Exit(1)
	}

	addr, _ := net.ResolveIPAddr("ip", host)
	fmt.Printf("tracing %v (%v), %v hops max, max retries: %v\n", host, addr.String(), *hops, *retries)

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
			printHop(hop)
		}
	}()

	config := tracer.NewConfig().WithHops(*hops).WithTimeout(*timeout)
	t := tracer.NewTracer(config)
	_, err := t.Run(host, c)
	wg.Wait()

	if err != nil {
		fmt.Println("Error from tracer: ", err)
	}
}

func printHop(hop tracer.Hop) {
	et := hop.ElapsedTime.String()
	if hop.ElapsedTime == 0 {
		et = "*"
	}
	fmt.Printf("%v.   %v    %v    %v\n", hop.TTL, hop.Addr, hop.Location, et)
}

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
	flag.Parse()

	host := flag.Arg(0)
	if len(host) == 0 {
		fmt.Println("no host provided!")
		os.Exit(1)
	}

	addr, _ := net.ResolveIPAddr("ip", host)
	fmt.Printf("tracing %v (%v), %v hops max\n", host, addr.String(), *hops)

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
	trace, err := t.Run(host, c)
	wg.Wait()

	if err != nil {
		fmt.Println("Error from tracer: ", err)
	} else {
		fmt.Println("Total Round Trip Time:", trace.RoundTripTime)
	}
}

func printHop(hop tracer.Hop) {
	fmt.Printf("%v.   %v    %v    %v\n", hop.TTL, hop.Addr, hop.Location, hop.ElapsedTime.String())
}

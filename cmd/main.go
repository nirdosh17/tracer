package main

import (
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/nirdosh17/tracer"
)

var usage = `
Usage:

	It requires privileged access to work with raw ICMP packets.
	Make sure to run the command as administrator.

	sudo trace [-hops] [-timeout] host

Examples:

	# trace with default settings, max hops: ` + fmt.Sprintf("%v", tracer.DEFAULT_HOPS) + `, timeout(seconds): ` + fmt.Sprintf("%v", tracer.DEFAULT_TIMEOUT_SECONDS) + `, retries: ` + fmt.Sprintf("%v", tracer.DEFAULT_MAX_RETRIES) + `
	trace google.com

	# trace 'n' number of hops
	trace -hops 10 example.com

	# if you are receiving blank response
	# try increasing the ICMP response timeout(-t) and retries(-r)
	trace -t 10 -r 5 example.com
`

func main() {
	hops := flag.Int("hops", tracer.DEFAULT_HOPS, "")
	timeout := flag.Int("t", tracer.DEFAULT_TIMEOUT_SECONDS, "")
	retries := flag.Int("r", tracer.DEFAULT_MAX_RETRIES, "")
	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	host := flag.Arg(0)
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

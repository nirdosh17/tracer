package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"sync"

	"github.com/nirdosh17/tracer"
)

var colorReset = "\033[0m"
var colorGreen = "\033[32m"
var usage = `
Usage:

	1. Tracing route requires privileged access to analyze raw ICMP packets.
	   Make sure to run the command as administrator.

		` + colorGreen + `sudo gotrace route [-hops] [-timeout] host` + colorReset + `

	2. Get your public ip
		` + colorGreen + `gotrace myip` + colorReset + `

Examples:

	# trace with default settings, max hops: ` + fmt.Sprintf("%v", tracer.DEFAULT_HOPS) + `, timeout(seconds): ` + fmt.Sprintf("%v", tracer.DEFAULT_TIMEOUT_SECONDS) + `, retries: ` + fmt.Sprintf("%v", tracer.DEFAULT_MAX_RETRIES) + `
	` + colorGreen + `sudo gotrace route google.com` + colorReset + `

	# trace 'n' number of hops
	` + colorGreen + `sudo gotrace route -hops 10 example.com` + colorReset + `

	# if you are receiving blank response
	# try increasing the ICMP response timeout(-t) and retries(-r)
	` + colorGreen + `sudo gotrace route -t 10 -r 5 example.com` + colorReset + `

	# get your public ipv4 address with ipv6 and location
	` + colorGreen + `gotrace myip` + colorReset + `

	# get ipv4 address only
	` + colorGreen + `gotrace myip --ipv4` + colorReset + `
`

func main() {
	routeCmd := flag.NewFlagSet("route", flag.ExitOnError)
	hops := routeCmd.Int("hops", tracer.DEFAULT_HOPS, "")
	timeout := routeCmd.Int("t", tracer.DEFAULT_TIMEOUT_SECONDS, "")
	retries := routeCmd.Int("r", tracer.DEFAULT_MAX_RETRIES, "")

	myipCmd := flag.NewFlagSet("myip", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Println(usage)
	}

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	// first arg is always binary name e.g. /tmp/go-build3122800919/b001/exe/main
	switch os.Args[1] {
	case "route":
		if !hasPrivilegedAccess() {
			fmt.Println("This command requires privileged access! Try with 'sudo'.")
			os.Exit(1)
		}

		routeCmd.Parse(os.Args[2:])
		if routeCmd.NArg() == 0 {
			flag.Usage()
			return
		}

		host := routeCmd.Arg(0)
		traceRoute(host, hops, retries, timeout)

	case "myip":
		ipv4 := myipCmd.Bool("ipv4", false, "print ipv4 address only")
		myipCmd.Parse(os.Args[2:])

		pubIPv4, pubIPv6, loc := tracer.PublicIP()

		if *ipv4 {
			fmt.Println(pubIPv4)
			return
		}

		fmt.Println("--- Your Public IP ---")
		fmt.Println("IPv4 - ", pubIPv4)
		fmt.Println("IPv6 - ", pubIPv6)
		fmt.Println(" Org - ", loc)

	case "-h", "-help":
		flag.Usage()

	default:
		fmt.Printf("unknown command '%v', run 'gotrace -help' for command usage\n", os.Args[1])
		os.Exit(1)
	}

}

func hasPrivilegedAccess() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	// UID '0' means root user
	return currentUser.Uid == "0"
}

func traceRoute(host string, hops, retries, timeout *int) {
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Printf("failed to resolve host '%v'\n", host)
		os.Exit(1)
	}

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
			hop.Print()
		}
	}()

	config := tracer.NewConfig().WithHops(*hops).WithTimeout(*timeout)
	t := tracer.NewTracer(config)
	_, err = t.Run(host, c)
	wg.Wait()

	if err != nil {
		fmt.Printf("failed to trace route for %v: %v\n", host, err)
		os.Exit(1)
	}
}

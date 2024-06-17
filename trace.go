package tracer

import (
	"fmt"
	"time"
)

type Tracer struct {
	Config *TracerConfig
}

type Hop struct {
	SourceIP   string
	SourcePort string
	// ip from which icmp packet was received
	DestIP   string
	DestPort string
	// current ttl(hops) of the packet
	TTL int
	// total time taken for this hop
	ElapsedTime time.Duration
}

func NewTracer(c *TracerConfig) *Tracer {
	return &Tracer{Config: c}
}

// Run sends packets to the specified host in loop recording each network hop until it reaches the destination or max hops is reached.
// It also collects traces in the given channel.
//
// e.g. domain = google.com
func (t Tracer) Run(domain string, traces chan Hop) error {
	fmt.Printf("configs: %+v\n", *t.Config)

	traces <- Hop{DestIP: "254.178.123.100"}

	close(traces)
	return nil
}

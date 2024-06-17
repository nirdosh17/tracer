# Tracer
A packet tracer in Go similer to Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop until it finally reaches the destination or max TTL is reached.

# Tracer
A packet tracer in Go similer to Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop until it finally reaches the destination or max TTL is reached.


### Test
```bash
go run cmd/main.go -help

# Yes, previleged access is needed while creatig raw sockets for ICMP
sudo go run cmd/main.go example.com

sudo go run cmd/main.go -hops 5 -timeout 2 example.com
```

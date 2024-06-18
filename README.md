# Tracer
Network diagnostic tool in Go inspired by Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop(router address) until it finally reaches the destination or max TTL is reached.

### Test
```bash
go run cmd/main.go -help

# Yes, privileged access is needed while creating raw sockets for ICMP
sudo go run cmd/main.go example.com

sudo go run cmd/main.go -hops 5 -timeout 2 example.com
```
**Sample output:**
```bash
tracing example.com (93.184.215.14), 64 hops max, max retries: 2
1.   192.168.101.1    private range    12.023833ms
2.   62.115.42.118    Arelion Sweden AB (France)    178.632ms
3.   62.115.122.159    TELIANET (United States)    281.884917ms
4.   62.115.123.125    TELIANET (United States)    305.271958ms
5.   62.115.175.71    Arelion Sweden AB (United States)    277.827958ms
6.   152.195.64.153    Edgecast Inc. (United States)    276.270542ms
7.   93.184.215.14    Edgecast Inc. (United Kingdom)    305.162792ms
```

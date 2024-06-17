# Tracer
A packet tracer in Go similar to Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop until it finally reaches the destination or max TTL is reached.


### Test
```bash
go run cmd/main.go -help

# Yes, privileged access is needed while creating raw sockets for ICMP
sudo go run cmd/main.go example.com

sudo go run cmd/main.go -hops 5 -timeout 2 example.com
```
**Sample output:**
```bash
tracing example.com (93.184.215.14), 64 hops max
1.   10.80.70.255     ()    1.154591ms
2.   138.197.249.110    Canada (DigitalOcean, LLC)    1.62412ms
3.   143.244.192.42    Canada (DigitalOcean, LLC)    867.255µs
4.   143.244.224.142    Canada (DigitalOcean, LLC)    923.144µs
5.   143.244.224.147    Canada (DigitalOcean, LLC)    1.125888ms
6.   213.248.88.244    United Kingdom (Arelion)    1.369646ms
7.   62.115.127.6    United Kingdom (Telia Company AB)    2.111305ms
8.   62.115.113.20    United States (Telia Company AB)    76.752301ms
9.   62.115.112.242    United States (Telia Company AB)    88.230903ms
10.   62.115.147.199    United States (Arelion Sweden AB)    78.454842ms
11.   152.195.68.135    United States (Verizon Business)    76.296614ms
12.   62.115.141.245    United States (Telia Company AB)    85.042073ms
13.   152.195.68.133    United States (Verizon Business)    122.720549ms
14.   93.184.215.14    United Kingdom ()    78.305133ms
Total Round Trip Time: 3.722480376s
```

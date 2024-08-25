# Tracer
Network diagnostic tool in Go inspired by Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop(router address) until it finally reaches the destination or max TTL is reached.

## [Usage](https://pkg.go.dev/github.com/nirdosh17/tracer)

### 1. CLI
  ```bash
  go install github.com/nirdosh17/tracer/cmd/gotrace@latest
  ```
  _You can find the binaries [HERE](https://github.com/nirdosh17/tracer/releases/latest)._

  **Commands:**
  ```bash
  # trace route to a host
  gotrace example.com

  # with options (max hops, timeout, retries)
  gotrace route -hops 5 -t 5 -r 5 example.com

  # get your public ip
  gotrace myip
  ```

  **Get your public IP:**

  ```bash
  $ gotrace myip
  +----------+----------------------------------------+
  | IPv4     | 101.129.138.66                         |
  | IPv6     | 2404:7c00:41:50ce:755f:69d3:c890:604b  |
  | Location | TELIANET (United States)               |
  +----------+----------------------------------------+

  # just get ipv4
  $ gotrace myip --ipv4
  101.129.138.66
  ```

  **Trace route:**

```bash
$ gotrace route example.com
tracing example.com (93.184.215.14), 64 hops max, max retries: 2
1.   192.168.101.1    private range    3.709ms
2.   62.115.42.118    Arelion Sweden AB (Germany)    172.783ms
3.   62.115.124.56    TELIANET (France)    193.725ms
4.   62.115.112.242    TELIANET (United States)    285.389ms
5.   62.115.123.125    TELIANET (United States)    288.032ms
6.   213.248.83.119    Arelion Sweden AB (United States)    261.788ms
7.   152.195.65.153    Edgecast Inc. (United States)    581.99ms
8.   93.184.215.14    Edgecast Inc. (United Kingdom)    286.464ms
```

### 2. Package
  Install Package:
  ```bash
  go get github.com/nirdosh17/tracer
  ```

  ```go
  package main

  import (
    "fmt"
    "sync"

    "github.com/nirdosh17/tracer"
  )

  func main() {
    host := "example.com"

    var wg sync.WaitGroup
    wg.Add(1)
    c := make(chan tracer.Hop)
    go liveReader(&wg, c)

    t := tracer.NewTracer(tracer.NewConfig())
    _, err := t.Run(host, c)
    if err != nil {
      fmt.Println("trace err: ", err)
    }
    wg.Wait()
  }

  // read live hops from channel
  func liveReader(wg *sync.WaitGroup, c chan tracer.Hop) {
    for {
      hop, ok := <-c
      // channel closed
      if !ok {
        wg.Done()
        return
      }
      fmt.Printf("%v.   %v    %v    %v\n", hop.TTL, hop.Addr, hop.Location, hop.ElapsedTime)
    }
  }

  ```

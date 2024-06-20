# Tracer
Network diagnostic tool in Go inspired by Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop(router address) until it finally reaches the destination or max TTL is reached.

## [Usage](https://pkg.go.dev/github.com/nirdosh17/tracer)

### Options
|                     |                 Details                    |  Default Value   |
|---------------------|--------------------------------------------|------------------|
|   `MaxHops`         | max network hops to trace the packet route |        64        |
|   `TimeoutSeconds`  | UDP call & ICMP response wait time         |        5         |
|   `MaxRetries`      | retrying UDP/ICMP with same TTL(hop)       |        2         |


### 1. CLI
  Download `gotrace` binary from [HERE](https://github.com/nirdosh17/tracer/releases).

  **Examples:**
  ```bash
  # tracing requires privileged access
  sudo gotrace example.com

  # with options (max hops, timeout, retries)
  sudo gotrace -hops 5 -t 5 -r 5 example.com

  # get your public ip
  gotrace myip

  # view command details
  gotrace -help
  ```

  **Find public IP:**
  ```
  $ gotrace myip
  Your Public IP:
  101.12.38.5 | ISP Name Pvt. Ltd (France)
  ```

  **Trace route:**
  ```bash
  $ sudo gotrace route example.com
  tracing example.com (93.184.215.14), 64 hops max, max retries: 2
  1.   192.168.101.1    private range    12.023833ms
  2.   62.115.42.118    Arelion Sweden AB (France)    178.632ms
  3.   62.115.122.159    TELIANET (United States)    281.884917ms
  4.   62.115.123.125    TELIANET (United States)    305.271958ms
  5.   62.115.175.71    Arelion Sweden AB (United States)    277.827958ms
  6.   152.195.64.153    Edgecast Inc. (United States)    276.270542ms
  7.   93.184.215.14    Edgecast Inc. (United Kingdom)    305.162792ms
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

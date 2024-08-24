# Tracer
Network diagnostic tool in Go inspired by Traceroute.

Makes UDP call to target host increasing the TTL(hops) of IP packet and recording the ICMP response for each hop(router address) until it finally reaches the destination or max TTL is reached.

## [Usage](https://pkg.go.dev/github.com/nirdosh17/tracer)

### 1. CLI
  ```bash
  go install github.com/nirdosh17/tracer/cmd/gotrace@latest
  ```
  _You can find the binaries [HERE](https://github.com/nirdosh17/tracer/releases/latest)._

  **Examples:**
  ```bash
  # trace route to a host
  gotrace example.com

  # with options (max hops, timeout, retries)
  gotrace route -hops 5 -t 5 -r 5 example.com

  # get your public ip
  gotrace myip

  # view command details
  gotrace -help
  ```

  **Get public IP:**

  ![Get IP](https://github.com/nirdosh17/tracer/assets/5920689/facec0bd-e9ee-4cc4-b182-e3444f2f95b5)

  **Trace route:**

  ![Traceroute](https://github.com/nirdosh17/tracer/assets/5920689/7429ca7c-9b2b-4691-8aec-e22ddc5dc858)


### 2. Use as a lib
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

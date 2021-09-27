# load balancing algorithm library

Goroutine-safe, High-performance general load balancing algorithm library.

Smooth weighted load balancing algorithm: [NGINX](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35) and [LVS](http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling)

Doublejump provides a revamped Google's jump consistent hash.

## Features

- WeightedRoundRobin
- SmoothWeightedRoundRobin
- ConsistentHash
- RoundRobin
- Random

## Installation

```go
go get -u github.com/fufuok/balancer
```

## Quickstart

```go
package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	wNodes := map[string]int{
		"A": 5,
		"B": 3,
		"C": 1,
		"D": 0,
	}
	balancer.Update(wNodes)

	// result of smooth selection is similar to: A A A B A B C A B
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
}
```

## Examples

please visit: [example](example)

### initialize the balancer

Sample data:

```go
// for WeightedRoundRobin/SmoothWeightedRoundRobin
// To be selected : Weighted
wNodes := map[string]int{
    "A": 5,
    "B": 3,
    "C": 1,
    "D": 0,
}
// for RoundRobin/Random/ConsistentHash
nodes := []string{"A", "B", "C"}
```

1. use default balancer (WRR)

   WeightedRoundRobin is the default balancer algorithm.

   ```go
   balancer.Update(wNodes)
   ```

2. use WeightedRoundRobin (WRR)

   ```go
   var lb Balancer
   lb = balancer.New(balancer.WeightedRoundRobin, wNodes)
   
   // or
   lb = balancer.New(balancer.WeightedRoundRobin, nil)
   lb.Update(wNodes)
   
   // or
   lb = balancer.NewWeightedRoundRobin(wNodes)
   
   // or
   lb = balancer.NewWeightedRoundRobin()
   lb.Update(wNodes)
   ```

3. use SmoothWeightedRoundRobin (SWRR)

   ```go
   var lb Balancer
   lb = balancer.New(balancer.SmoothWeightedRoundRobin, wNodes)
   
   // or
   lb = balancer.New(balancer.SmoothWeightedRoundRobin, nil)
   lb.Update(wNodes)
   
   // or
   lb = balancer.NewSmoothWeightedRoundRobin(wNodes)
   
   // or
   lb = balancer.NewSmoothWeightedRoundRobin()
   lb.Update(wNodes)
   ```

4. use ConsistentHash

   ```go
   var lb Balancer
   lb = balancer.New(balancer.ConsistentHash, nodes)
   
   // or
   lb = balancer.New(balancer.ConsistentHash, nil)
   lb.Update(nodes)
   
   // or
   lb = balancer.NewConsistentHash(nodes)
   
   // or
   lb = balancer.NewConsistentHash()
   lb.Update(nodes)
   ```

5. use RoundRobin (RR)

   ```go
   var lb Balancer
   lb = balancer.New(balancer.RoundRobin, nodes)
   
   // or
   lb = balancer.New(balancer.RoundRobin, nil)
   lb.Update(nodes)
   
   // or
   lb = balancer.NewRoundRobin(nodes)
   
   // or
   lb = balancer.NewRoundRobin()
   lb.Update(nodes)
   ```

6. use Random

   ```go
   var lb Balancer
   lb = balancer.New(balancer.Random, nodes)
   
   // or
   lb = balancer.New(balancer.Random, nil)
   lb.Update(nodes)
   
   // or
   lb = balancer.NewRandom(nodes)
   
   // or
   lb = balancer.NewRandom()
   lb.Update(nodes)
   ```

### gets next selected item

```go
node := lb.Select()
```

ip consistent hash:

```go
node := lb.Select("192.168.1.100")
node := lb.Select("192.168.1.100", "Test", "...")
```

### interface

```go
type Balancer interface {
	// Add add an item to be selected.
	Add(item string, weight int)

	// All get all items.
	// RoundRobin/Random/ConsistentHash: []string
	// WeightedRoundRobin/SmoothWeightedRoundRobin: map[string]int
	All() interface{}

	// Select gets next selected item.
	// key is only used for ConsistentHash
	Select(key ...string) string

	// Name load balancer name.
	Name() string

	// Remove remove an item.
	// asClean: clean up or remove only one
	Remove(item string, asClean ...bool) bool

	// RemoveAll remove all items.
	RemoveAll()

	// Reset reset the balancer.
	Reset()

	// Update reinitialize the balancer items.
	// RoundRobin/Random/ConsistentHash: []string
	// WeightedRoundRobin/SmoothWeightedRoundRobin: map[string]int
	Update(items interface{}) bool
}
```

## Benchmark

```shell
go test -bench=. -benchtime=1s -count=2
goos: linux
goarch: amd64
pkg: github.com/fufuok/balancer
cpu: Intel(R) Xeon(R) CPU E5-2667 v2 @ 3.30GHz
BenchmarkWeightedRoundRobin-4           30918976                38.35 ns/op            0 B/op          0 allocs/op
BenchmarkWeightedRoundRobin-4           30627727                38.33 ns/op            0 B/op          0 allocs/op
BenchmarkSmoothWeightedRoundRobin-4     20629270                59.60 ns/op            0 B/op          0 allocs/op
BenchmarkSmoothWeightedRoundRobin-4     21024522                75.05 ns/op            0 B/op          0 allocs/op
BenchmarkConsistentHash-4               26980348                43.72 ns/op            0 B/op          0 allocs/op
BenchmarkConsistentHash-4               26978271                43.77 ns/op            0 B/op          0 allocs/op
BenchmarkRoundRobin-4                   51317532                23.18 ns/op            0 B/op          0 allocs/op
BenchmarkRoundRobin-4                   51717782                29.64 ns/op            0 B/op          0 allocs/op
BenchmarkRandomR-4                      47053610                32.63 ns/op            0 B/op          0 allocs/op
BenchmarkRandomR-4                      46210752                25.15 ns/op            0 B/op          0 allocs/op
```

## License

Third-party library licenses:

- [doublejump]([doublejump/LICENSE at master · edwingeng/doublejump (github.com)](https://github.com/edwingeng/doublejump/blob/master/LICENSE))
- [go-jump]([go-jump/LICENSE at master · dgryski/go-jump (github.com)](https://github.com/dgryski/go-jump/blob/master/LICENSE))





*ff*
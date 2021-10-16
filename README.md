# ⚖️ load balancing algorithm library

Goroutine-safe, High-performance general load balancing algorithm library.

Smooth weighted load balancing algorithm: [NGINX](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35) and [LVS](http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling), Doublejump provides a revamped Google's jump consistent hash.

------

If you want a **faster** load balancer that supports **interface()**, please refer to another library: [fufuok/load-balancer](https://github.com/fufuok/load-balancer)

------

## 🎯 Features

- WeightedRoundRobin
- SmoothWeightedRoundRobin
- WeightedRand
- ConsistentHash
- RoundRobin
- Random

## ⚙️ Installation

```go
go get -u github.com/fufuok/balancer
```

## ⚡️ Quickstart

```go
package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	wNodes := map[string]int{
		"🍒": 5,
		"🍋": 3,
		"🍉": 1,
		"🥑": 0,
	}
	balancer.Update(wNodes)

	// result of smooth selection is similar to: 🍒 🍒 🍒 🍋 🍒 🍋 🍒 🍋 🍉
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
	fmt.Println()
}
```

## 📚 Examples

please see: [examples](examples)

### Initialize the balancer

Sample data:

```go
// for WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand
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
   var lb balancer.Balancer
   lb = balancer.New(balancer.WeightedRoundRobin, wNodes, nil)
   
   // or
   lb = balancer.New(balancer.WeightedRoundRobin, nil, nil)
   lb.Update(wNodes)
   
   // or
   lb = balancer.NewWeightedRoundRobin(wNodes)
   
   // or
   lb = balancer.NewWeightedRoundRobin()
   lb.Update(wNodes)
   ```

3. use SmoothWeightedRoundRobin (SWRR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.SmoothWeightedRoundRobin, wNodes, nil)
   
   // or
   lb = balancer.New(balancer.SmoothWeightedRoundRobin, nil, nil)
   lb.Update(wNodes)
   
   // or
   lb = balancer.NewSmoothWeightedRoundRobin(wNodes)
   
   // or
   lb = balancer.NewSmoothWeightedRoundRobin()
   lb.Update(wNodes)
   ```

4. use WeightedRand (WR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.WeightedRand, wNodes, nil)
   
   // or
   lb = balancer.New(balancer.WeightedRand, nil, nil)
   lb.Update(wNodes)
   
   // or
   lb = balancer.NewWeightedRand(wNodes)
   
   // or
   lb = balancer.NewWeightedRand()
   lb.Update(wNodes)
   ```

5. use ConsistentHash

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.ConsistentHash, nil, nodes)
   
   // or
   lb = balancer.New(balancer.ConsistentHash, nil, nil)
   lb.Update(nodes)
   
   // or
   lb = balancer.NewConsistentHash(nodes)
   
   // or
   lb = balancer.NewConsistentHash()
   lb.Update(nodes)
   ```

6. use RoundRobin (RR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.RoundRobin, nil, nodes)
   
   // or
   lb = balancer.New(balancer.RoundRobin, nil, nil)
   lb.Update(nodes)
   
   // or
   lb = balancer.NewRoundRobin(nodes)
   
   // or
   lb = balancer.NewRoundRobin()
   lb.Update(nodes)
   ```

7. use Random

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.Random, nil, nodes)
   
   // or
   lb = balancer.New(balancer.Random, nil, nil)
   lb.Update(nodes)
   
   // or
   lb = balancer.NewRandom(nodes)
   
   // or
   lb = balancer.NewRandom()
   lb.Update(nodes)
   ```

### Gets next selected item

```go
node := lb.Select()
```

ip consistent hash:

```go
node := lb.Select("192.168.1.100")
node := lb.Select("192.168.1.100", "Test", "...")
```

### Interface

```go
type Balancer interface {
	// Add add an item to be selected.
	// weight is only used for WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand, default: 1
	Add(item string, weight ...int)

	// All get all items.
	// RoundRobin/Random/ConsistentHash: []string
	// WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand: map[string]int
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
	// WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand: map[string]int
	Update(items interface{}) bool
}
```

## 🤖 Benchmarks

```shell
go test -run=^$ -benchmem -benchtime=1s -count=1 -bench=.
goos: linux
goarch: amd64
pkg: github.com/fufuok/balancer
cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
BenchmarkBalancer/WRR/10-4                              37112553                30.34 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR/10-4                             38851680                30.39 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR/10-4                               36406916                33.15 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash/10-4                             31506262                37.60 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin/10-4                       53076963                23.86 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random/10-4                           64582524                18.02 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR#01/100-4                          32221255                37.31 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR#01/100-4                          7337542                165.6 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR#01/100-4                           21253034                53.29 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash#01/100-4                         25851721                46.58 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin#01/100-4                   51670482                22.59 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random#01/100-4                       66175606                18.14 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR#02/1000-4                         28502208                42.09 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR#02/1000-4                          872499                 1391 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR#02/1000-4                          16595787                71.57 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash#02/1000-4                        19103568                63.47 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin#02/1000-4                  52725135                23.05 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random#02/1000-4                      66541184                18.24 ns/op            0 B/op          0 allocs/op
```

## ⚠️ License

Third-party library licenses:

- [doublejump]([doublejump/LICENSE at master · edwingeng/doublejump (github.com)](https://github.com/edwingeng/doublejump/blob/master/LICENSE))
- [go-jump]([go-jump/LICENSE at master · dgryski/go-jump (github.com)](https://github.com/dgryski/go-jump/blob/master/LICENSE))





*ff*
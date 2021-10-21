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
BenchmarkBalancer/WRR-10-4                              37112553                30.34 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-10-4                             38851680                30.39 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-10-4                               36406916                33.15 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-10-4                             31506262                37.60 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-10-4                       53076963                23.86 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-10-4                           64582524                18.02 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-100-4                             32221255                37.31 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-100-4                             7337542                165.6 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-100-4                              21253034                53.29 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-100-4                            25851721                46.58 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-100-4                      51670482                22.59 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-100-4                          66175606                18.14 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-1000-4                            28502208                42.09 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-1000-4                             872499                 1391 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-1000-4                             16595787                71.57 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-1000-4                           19103568                63.47 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-1000-4                     52725135                23.05 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-1000-4                         66541184                18.24 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-10000-4                           27912939                42.67 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-10000-4                             86983                14019 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-10000-4                            12691062                92.73 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-10000-4                          16084016                73.96 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-10000-4                    52327888                24.05 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-10000-4                        66457050                18.17 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-100000-4                          24896972                43.20 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-100000-4                             7046               173884 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-100000-4                            9491140                127.3 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-100000-4                         16090567                76.46 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-100000-4                   49422337                24.06 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-100000-4                       62700792                18.52 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-1000000-4                         24038544                47.01 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-1000000-4                             381              3108476 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-1000000-4                           4863207                259.8 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-1000000-4                        16194163                74.15 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-1000000-4                  52971156                22.70 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-1000000-4                      40018782                26.25 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-10-4                      12663168                92.42 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-10-4                     12709807                94.45 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-10-4                       28055660                68.66 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-10-4                     16188872                72.30 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-10-4               13612167                91.14 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-10-4                   17698065                67.72 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-100-4                     13133222                88.99 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-100-4                     6563328                181.4 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-100-4                      26578320                79.94 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-100-4                    16083211                74.67 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-100-4              12263182                91.64 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-100-4                  17741816                67.81 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-1000-4                    11212780                104.0 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-1000-4                     805879                 1474 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-1000-4                     15821539                72.37 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-1000-4                   14478384                81.82 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-1000-4             12103447                88.15 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-1000-4                 17729145                67.81 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-10000-4                   10567130                105.9 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-10000-4                     81170                14685 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-10000-4                    14379578                78.16 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-10000-4                  14215629                84.50 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-10000-4            13372892                90.45 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-10000-4                17676268                67.92 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-100000-4                  11561236                110.4 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-100000-4                     6835               175792 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-100000-4                   22145756                54.60 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-100000-4                 14285690                84.04 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-100000-4           12744205                90.57 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-100000-4               17859376                67.19 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-1000000-4                 11473530                109.0 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-1000000-4                     379              3127391 ns/op            7 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-1000000-4                  10868928                100.4 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-1000000-4                14018941                85.23 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-1000000-4          11614006                87.05 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-1000000-4              17565664                68.35 ns/op            0 B/op          0 allocs/op
```

## ⚠️ License

Third-party library licenses:

- [doublejump]([doublejump/LICENSE at master · edwingeng/doublejump (github.com)](https://github.com/edwingeng/doublejump/blob/master/LICENSE))
- [go-jump]([go-jump/LICENSE at master · dgryski/go-jump (github.com)](https://github.com/dgryski/go-jump/blob/master/LICENSE))





*ff*
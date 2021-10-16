package balancer

import (
	"math/rand"
	"strconv"
	"testing"
)

const (
	numMin = 10
	numMax = 1_000_000
)

func BenchmarkBalancer(b *testing.B) {
	for n := numMin; n <= numMax; n *= 10 {
		b.Run("WRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsMap(n)
				lb := NewWeightedRoundRobin(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("SWRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsMap(n)
				lb := NewSmoothWeightedRoundRobin(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("WR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsMap(n)
				lb := NewWeightedRand(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("Hash", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsSlice(n)
				lb := NewConsistentHash(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select("192.168.1.1")
				}
			})
		})

		b.Run("RoundRobin", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsSlice(n)
				lb := NewRoundRobin(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("Random", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsSlice(n)
				lb := NewRandom(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})
	}
}

func BenchmarkBalancerParallel(b *testing.B) {
	for n := numMin; n <= numMax; n *= 10 {
		b.Run("WRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsMap(n)
				lb := NewWeightedRoundRobin(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("SWRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsMap(n)
				lb := NewSmoothWeightedRoundRobin(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("WR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsMap(n)
				lb := NewWeightedRand(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("Hash", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsSlice(n)
				lb := NewConsistentHash(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select("192.168.1.1")
					}
				})
			})
		})

		b.Run("RoundRobin", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsSlice(n)
				lb := NewRoundRobin(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("Random", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				items := genItemsSlice(n)
				lb := NewRandom(items)
				for i := 0; i < 1000; i++ {
					lb.Select()
				}
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})
	}
}

func genItemsMap(n int) map[string]int {
	items := make(map[string]int)
	for i := 0; i < n; i++ {
		items[strconv.Itoa(i)] = rand.Intn(20)
	}
	return items
}

func genItemsSlice(n int) (items []string) {
	for i := 0; i < n; i++ {
		items = append(items, strconv.Itoa(i))
	}
	return
}

// go test -run=^$ -benchmem -benchtime=1s -count=1 -bench=.
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/balancer
// cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
// BenchmarkBalancer/WRR/10-4                              37112553                30.34 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR/10-4                             38851680                30.39 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR/10-4                               36406916                33.15 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash/10-4                             31506262                37.60 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin/10-4                       53076963                23.86 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random/10-4                           64582524                18.02 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#01/100-4                          32221255                37.31 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#01/100-4                          7337542                165.6 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#01/100-4                           21253034                53.29 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#01/100-4                         25851721                46.58 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#01/100-4                   51670482                22.59 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#01/100-4                       66175606                18.14 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#02/1000-4                         28502208                42.09 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#02/1000-4                          872499                 1391 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#02/1000-4                          16595787                71.57 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#02/1000-4                        19103568                63.47 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#02/1000-4                  52725135                23.05 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#02/1000-4                      66541184                18.24 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#03/10000-4                        27912939                42.67 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#03/10000-4                          86983                14019 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#03/10000-4                         12691062                92.73 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#03/10000-4                       16084016                73.96 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#03/10000-4                 52327888                24.05 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#03/10000-4                     66457050                18.17 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#04/100000-4                       24896972                43.20 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#04/100000-4                          7046               173884 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#04/100000-4                         9491140                127.3 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#04/100000-4                      16090567                76.46 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#04/100000-4                49422337                24.06 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#04/100000-4                    62700792                18.52 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#05/1000000-4                      24038544                47.01 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#05/1000000-4                          381              3108476 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#05/1000000-4                        4863207                259.8 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#05/1000000-4                     16194163                74.15 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#05/1000000-4               52971156                22.70 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#05/1000000-4                   40018782                26.25 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR/10-4                      12663168                92.42 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR/10-4                     12709807                94.45 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR/10-4                       28055660                68.66 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash/10-4                     16188872                72.30 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin/10-4               13612167                91.14 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random/10-4                   17698065                67.72 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#01/100-4                  13133222                88.99 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#01/100-4                  6563328                181.4 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#01/100-4                   26578320                79.94 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#01/100-4                 16083211                74.67 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#01/100-4           12263182                91.64 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#01/100-4               17741816                67.81 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#02/1000-4                 11212780                104.0 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#02/1000-4                  805879                 1474 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#02/1000-4                  15821539                72.37 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#02/1000-4                14478384                81.82 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#02/1000-4          12103447                88.15 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#02/1000-4              17729145                67.81 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#03/10000-4                10567130                105.9 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#03/10000-4                  81170                14685 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#03/10000-4                 14379578                78.16 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#03/10000-4               14215629                84.50 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#03/10000-4         13372892                90.45 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#03/10000-4             17676268                67.92 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#04/100000-4               11561236                110.4 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#04/100000-4                  6835               175792 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#04/100000-4                22145756                54.60 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#04/100000-4              14285690                84.04 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#04/100000-4        12744205                90.57 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#04/100000-4            17859376                67.19 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#05/1000000-4              11473530                109.0 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#05/1000000-4                  379              3127391 ns/op            7 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#05/1000000-4               10868928                100.4 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#05/1000000-4             14018941                85.23 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#05/1000000-4       11614006                87.05 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#05/1000000-4           17565664                68.35 ns/op            0 B/op          0 allocs/op

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
// cpu: Intel(R) Xeon(R) CPU E5-2667 v2 @ 3.30GHz
// BenchmarkBalancer/WRR/10-4                              38190100                36.54 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR/10-4                             31051455                39.56 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash/10-4                             32343030                47.23 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin/10-4                       58601172                20.14 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random/10-4                           71845088                16.46 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#01/100-4                          33369616                34.03 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#01/100-4                          6438565               237.5 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#01/100-4                         26024400                57.22 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#01/100-4                   59103342                30.60 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#01/100-4                       70480639                16.35 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#02/1000-4                         31300266                49.85 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#02/1000-4                          753825              2008 ns/op               0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#02/1000-4                        20951144                56.93 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#02/1000-4                  58551811                19.75 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#02/1000-4                      71825866                21.06 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#03/10000-4                        31281838                49.10 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#03/10000-4                          60955             24994 ns/op               0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#03/10000-4                       17948910                65.05 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#03/10000-4                 59174136                19.91 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#03/10000-4                     71194105                21.20 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#04/100000-4                       27980596                51.19 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#04/100000-4                          5926            228089 ns/op               0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#04/100000-4                      18386246                83.29 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#04/100000-4                59072595                31.17 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#04/100000-4                    72389031                16.41 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#05/1000000-4                      12541518                90.22 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#05/1000000-4                          332           3611514 ns/op               0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#05/1000000-4                     18413552                70.58 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#05/1000000-4               17556219                60.58 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#05/1000000-4                   36476977                29.91 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR/10-4                      12844932               305.4 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR/10-4                      3143989               411.2 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash/10-4                     16331000                72.65 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin/10-4               11692734               205.8 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random/10-4                   21195294                55.94 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#01/100-4                  13749801               261.5 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#01/100-4                  1890832               630.9 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#01/100-4                 16473523                72.75 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#01/100-4           11893130               238.4 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#01/100-4               20649256                58.01 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#02/1000-4                 11632530               208.0 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#02/1000-4                  232335              5108 ns/op               0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#02/1000-4                17453338                68.80 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#02/1000-4          12115618               192.4 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#02/1000-4              21093782                56.78 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#03/10000-4                11394490               283.3 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#03/10000-4                  17031             69365 ns/op               0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#03/10000-4               16612495                72.17 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#03/10000-4         12598099               225.0 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#03/10000-4             20785418                57.34 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#04/100000-4               11428856               339.9 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#04/100000-4                  5629            213868 ns/op               1 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#04/100000-4              17755728                86.29 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#04/100000-4        12290454               217.1 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#04/100000-4            22997978                51.96 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#05/1000000-4               3110058               419.4 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#05/1000000-4                  337           3573590 ns/op               4 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#05/1000000-4             17034688                65.07 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#05/1000000-4       13031574               252.0 ns/op             0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#05/1000000-4           23596234                50.20 ns/op            0 B/op          0 allocs/op

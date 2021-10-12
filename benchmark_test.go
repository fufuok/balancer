package balancer

import (
	"strconv"
	"testing"
)

func BenchmarkWeightedRoundRobin(b *testing.B) {
	lb := NewWeightedRoundRobin()
	for i := 0; i < 15; i++ {
		lb.Add("item-"+strconv.Itoa(i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lb.Select()
	}
}

func BenchmarkSmoothWeightedRoundRobin(b *testing.B) {
	lb := NewSmoothWeightedRoundRobin()
	for i := 0; i < 15; i++ {
		lb.Add("item-"+strconv.Itoa(i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lb.Select()
	}
}

func BenchmarkConsistentHash(b *testing.B) {
	lb := NewConsistentHash()
	for i := 0; i < 15; i++ {
		lb.Add("item-"+strconv.Itoa(i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lb.Select("192.168.1.100")
	}
}

func BenchmarkRoundRobin(b *testing.B) {
	lb := NewRoundRobin()
	for i := 0; i < 15; i++ {
		lb.Add("item-"+strconv.Itoa(i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lb.Select()
	}
}

func BenchmarkRandomR(b *testing.B) {
	lb := NewRandom()
	for i := 0; i < 15; i++ {
		lb.Add("item-"+strconv.Itoa(i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lb.Select()
	}
}

// go test -bench=. -benchtime=1s -count=2
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/balancer
// cpu: Intel(R) Xeon(R) CPU E5-2667 v2 @ 3.30GHz
// BenchmarkWeightedRoundRobin-4           30947391                37.95 ns/op            0 B/op          0 allocs/op
// BenchmarkWeightedRoundRobin-4           30508263                48.99 ns/op            0 B/op          0 allocs/op
// BenchmarkSmoothWeightedRoundRobin-4     20379349                60.12 ns/op            0 B/op          0 allocs/op
// BenchmarkSmoothWeightedRoundRobin-4     20407590                75.66 ns/op            0 B/op          0 allocs/op
// BenchmarkConsistentHash-4               33911062                44.45 ns/op            0 B/op          0 allocs/op
// BenchmarkConsistentHash-4               34818607                44.34 ns/op            0 B/op          0 allocs/op
// BenchmarkRoundRobin-4                   51605906                23.35 ns/op            0 B/op          0 allocs/op
// BenchmarkRoundRobin-4                   50344368                29.87 ns/op            0 B/op          0 allocs/op
// BenchmarkRandomR-4                      45422324                26.06 ns/op            0 B/op          0 allocs/op
// BenchmarkRandomR-4                      46060688                32.93 ns/op            0 B/op          0 allocs/op

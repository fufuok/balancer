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
// BenchmarkWeightedRoundRobin-4           30918976                38.35 ns/op            0 B/op          0 allocs/op
// BenchmarkWeightedRoundRobin-4           30627727                38.33 ns/op            0 B/op          0 allocs/op
// BenchmarkSmoothWeightedRoundRobin-4     20629270                59.60 ns/op            0 B/op          0 allocs/op
// BenchmarkSmoothWeightedRoundRobin-4     21024522                75.05 ns/op            0 B/op          0 allocs/op
// BenchmarkConsistentHash-4               26980348                43.72 ns/op            0 B/op          0 allocs/op
// BenchmarkConsistentHash-4               26978271                43.77 ns/op            0 B/op          0 allocs/op
// BenchmarkRoundRobin-4                   51317532                23.18 ns/op            0 B/op          0 allocs/op
// BenchmarkRoundRobin-4                   51717782                29.64 ns/op            0 B/op          0 allocs/op
// BenchmarkRandomR-4                      47053610                32.63 ns/op            0 B/op          0 allocs/op
// BenchmarkRandomR-4                      46210752                25.15 ns/op            0 B/op          0 allocs/op

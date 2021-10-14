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

// go test -run=^$ -bench=. -benchtime=1s -count=2
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/balancer
// cpu: Intel(R) Xeon(R) CPU E5-2667 v2 @ 3.30GHz
// BenchmarkWeightedRoundRobin-4           31541454                37.92 ns/op            0 B/op          0 allocs/op
// BenchmarkWeightedRoundRobin-4           30563991                49.22 ns/op            0 B/op          0 allocs/op
// BenchmarkSmoothWeightedRoundRobin-4     20565193                58.03 ns/op            0 B/op          0 allocs/op
// BenchmarkSmoothWeightedRoundRobin-4     20774737                57.12 ns/op            0 B/op          0 allocs/op
// BenchmarkConsistentHash-4               33287440                35.11 ns/op            0 B/op          0 allocs/op
// BenchmarkConsistentHash-4               33085503                35.03 ns/op            0 B/op          0 allocs/op
// BenchmarkRoundRobin-4                   51950608                35.87 ns/op            0 B/op          0 allocs/op
// BenchmarkRoundRobin-4                   51404785                22.90 ns/op            0 B/op          0 allocs/op
// BenchmarkRandomR-4                      55436377                21.16 ns/op            0 B/op          0 allocs/op
// BenchmarkRandomR-4                      55588168                21.13 ns/op            0 B/op          0 allocs/op

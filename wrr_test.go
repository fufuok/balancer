package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestWeightedRoundRobin(t *testing.T) {
	lb := NewWeightedRoundRobin()
	item := lb.Select()
	if item != "" {
		t.Fatalf("wrr expected empty, actual %s", item)
	}

	lb.Add("A", 0)
	item = lb.Select()
	if item != "" {
		t.Fatalf("wrr expected empty, actual %s", item)
	}
	lb.Add("B", 1)
	item = lb.Select()
	if item != "B" {
		t.Fatalf("wrr expected B, actual %s", item)
	}
	item = lb.Select("test")
	if item != "B" {
		t.Fatalf("wrr expected B, actual %s", item)
	}

	nodes := map[string]int{
		"A": 0,
		"B": 1,
		"C": 7,
		"D": 2,
	}
	lb = NewWeightedRoundRobin(nodes)
	count := make(map[string]int)
	for i := 0; i < 1000; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] != 100 || count["C"] != 700 || count["D"] != 200 {
		t.Fatal("wrr wrong")
	}

	lb.Reset()

	lb.Add("E", 10)
	for i := 0; i < 2000; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] != 200 || count["C"] != 1400 || count["D"] != 400 || count["E"] != 1000 {
		t.Fatal("wrr reset() wrong")
	}

	ok := lb.Remove("E")
	if ok != true {
		t.Fatal("wrr remove() wrong")
	}

	lb.Reset()

	for i := 0; i < 1000; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] != 300 || count["C"] != 2100 || count["D"] != 600 {
		t.Fatal("wrr wrong")
	}

	lb.RemoveAll()
	lb.Add("F", 2)
	lb.Add("F", 1)
	all, ok := lb.All().(map[string]int)
	if !ok || all["F"] != 1 {
		t.Fatal("wrr all() wrong")
	}

	nodes = map[string]int{
		"X": 0,
		"Y": 1,
	}
	ok = lb.Update(nodes)
	if ok != true {
		t.Fatal("wrr update wrong")
	}
	all, ok = lb.All().(map[string]int)
	if !ok || all["Y"] != 1 {
		t.Fatal("swrr all() wrong")
	}
	item = lb.Select()
	if item != "Y" {
		t.Fatal("wrr update wrong")
	}
}

func TestWeightedRoundRobin_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := map[string]int{
		"A": 5,
		"B": 1,
		"C": 4,
		"D": 0,
	}
	lb := NewWeightedRoundRobin(nodes)

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				switch lb.Select() {
				case "A":
					atomic.AddInt64(&a, 1)
				case "B":
					atomic.AddInt64(&b, 1)
				case "C":
					atomic.AddInt64(&c, 1)
				case "D":
					atomic.AddInt64(&d, 1)
				}
			}
		}()
	}
	wg.Wait()

	if atomic.LoadInt64(&a) != 500000 {
		t.Fatal("wrr wrong: a")
	}
	if atomic.LoadInt64(&b) != 100000 {
		t.Fatal("wrr wrong: b")
	}
	if atomic.LoadInt64(&c) != 400000 {
		t.Fatal("wrr wrong: c")
	}
	if atomic.LoadInt64(&d) != 0 {
		t.Fatal("wrr wrong: d")
	}
}

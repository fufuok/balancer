package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestWeightedRand(t *testing.T) {
	lb := NewWeightedRand()
	item := lb.Select()
	if item != "" {
		t.Fatalf("wr expected empty, actual %s", item)
	}

	lb.Add("A", 0)
	item = lb.Select()
	if item != "" {
		t.Fatalf("wr expected empty, actual %s", item)
	}
	lb.Add("B", 1)
	item = lb.Select()
	if item != "B" {
		t.Fatalf("wr expected B, actual %s", item)
	}
	item = lb.Select("test")
	if item != "B" {
		t.Fatalf("wr expected B, actual %s", item)
	}

	nodes := map[string]int{
		"A": 0,
		"B": 1,
		"C": 7,
		"D": 2,
	}
	lb = NewWeightedRand(nodes)
	count := make(map[string]int)
	for i := 0; i < 2000; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] <= 150 || count["C"] <= 750 || count["D"] <= 250 {
		t.Fatal("wr wrong")
	}
	if count["A"]+count["B"]+count["C"]+count["D"] != 2000 {
		t.Fatal("wr wrong")
	}

	lb.RemoveAll()
	lb.Add("F", 2)
	lb.Add("F", 1)
	all, ok := lb.All().(map[string]int)
	if !ok || all["F"] != 1 {
		t.Fatal("wr all() wrong")
	}

	lb.Remove("F")
	item = lb.Select()
	if item != "" {
		t.Fatalf("wr expected empty, actual %s", item)
	}

	nodes = map[string]int{
		"X": 0,
		"Y": 1,
	}
	ok = lb.Update(nodes)
	if ok != true {
		t.Fatal("wr update wrong")
	}
	all, ok = lb.All().(map[string]int)
	if !ok || all["Y"] != 1 {
		t.Fatal("swrr all() wrong")
	}
	item = lb.Select()
	if item != "Y" {
		t.Fatal("wr update wrong")
	}
}

func TestWeightedRand_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := map[string]int{
		"A": 5,
		"B": 1,
		"C": 4,
		"D": 0,
	}
	lb := NewWeightedRand(nodes)

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

	if atomic.LoadInt64(&a) <= 410000 {
		t.Fatal("wr wrong: a")
	}
	if atomic.LoadInt64(&b) <= 90000 {
		t.Fatal("wr wrong: b")
	}
	if atomic.LoadInt64(&c) <= 390000 {
		t.Fatal("wr wrong: c")
	}
	if atomic.LoadInt64(&d) != 0 {
		t.Fatal("wr wrong: d")
	}
	if atomic.LoadInt64(&a)+atomic.LoadInt64(&b)+atomic.LoadInt64(&c) != 1000000 {
		t.Fatal("wr wrong: sum")
	}
}

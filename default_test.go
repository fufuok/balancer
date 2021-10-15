package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestDefaultBalancer(t *testing.T) {
	item := Select()
	if item != "" {
		t.Fatalf("default balancer expected empty, actual %s", item)
	}

	Add("A", 0)
	item = Select()
	if item != "" {
		t.Fatalf("default balancer expected empty, actual %s", item)
	}
	Add("B", 1)
	item = Select()
	if item != "B" {
		t.Fatalf("default balancer expected B, actual %s", item)
	}

	nodes := map[string]int{
		"A": 0,
		"B": 1,
		"C": 7,
		"D": 2,
	}
	Update(nodes)
	count := make(map[string]int)
	for i := 0; i < 1000; i++ {
		item := Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] != 100 || count["C"] != 700 || count["D"] != 200 {
		t.Fatal("default balancer wrong")
	}

	Reset()

	Add("E", 10)
	for i := 0; i < 2000; i++ {
		item := Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] != 200 || count["C"] != 1400 || count["D"] != 400 || count["E"] != 1000 {
		t.Fatal("default balancer reset() wrong")
	}

	ok := Remove("E")
	if ok != true {
		t.Fatal("default balancer remove() wrong")
	}

	Reset()

	for i := 0; i < 1000; i++ {
		item := Select()
		count[item]++
	}
	if count["A"] != 0 || count["B"] != 300 || count["C"] != 2100 || count["D"] != 600 {
		t.Fatal("default balancer wrong")
	}

	RemoveAll()
	Add("F", 2)
	Add("F", 1)
	all, ok := All().(map[string]int)
	if !ok || all["F"] != 1 {
		t.Fatal("default balancer all() wrong")
	}

	nodes = map[string]int{
		"X": 0,
		"Y": 1,
	}
	ok = Update(nodes)
	if ok != true {
		t.Fatal("default balancer update wrong")
	}
	item = Select()
	if item != "Y" {
		t.Fatal("default balancer update wrong")
	}

	if Name() != "WeightedRoundRobin" {
		t.Fatal("default balancer name wrong")
	}

	RemoveAll()
}

func TestDefaultBalancer_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := map[string]int{
		"A": 5,
		"B": 1,
		"C": 4,
		"D": 0,
	}
	Update(nodes)

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				switch Select() {
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
		t.Fatal("default balancer wrong: a")
	}
	if atomic.LoadInt64(&b) != 100000 {
		t.Fatal("default balancer wrong: b")
	}
	if atomic.LoadInt64(&c) != 400000 {
		t.Fatal("default balancer wrong: c")
	}
	if atomic.LoadInt64(&d) != 0 {
		t.Fatal("default balancer wrong: d")
	}

	RemoveAll()
}

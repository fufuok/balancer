package balancer

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRandom(t *testing.T) {
	lb := NewRandom()
	item := Select()
	if item != "" {
		t.Fatalf("r expected empty, actual %s", item)
	}

	lb.Add("A", 0)
	item = lb.Select()
	if item != "A" {
		t.Fatalf("r expected A, actual %s", item)
	}

	nodes := []string{"A", "B", "C", "D"}
	lb = NewRandom(nodes)
	count := make(map[string]int)
	for i := 0; i < 2000; i++ {
		ok := lb.Select()
		count[ok]++
	}
	if count["A"] <= 300 || count["B"] <= 300 || count["C"] <= 300 || count["D"] <= 300 {
		t.Fatal("r wrong")
	}
	if count["A"]+count["B"]+count["C"]+count["D"] != 2000 {
		t.Fatal("r wrong")
	}

	lb.Add("A", 0)
	lb.Add("A", 0)
	lb.Add("C", 0)
	lb.Add("A", 0)
	all := lb.All().([]string)
	if strings.Join(all, "") != "ABCDAACA" {
		t.Fatalf("hash all() wrong")
	}

	lb.Remove("C")
	all = lb.All().([]string)
	if strings.Join(all, "") != "ABDAACA" {
		t.Fatalf("hash all() wrong")
	}

	lb.Remove("A", true)
	all = lb.All().([]string)
	if strings.Join(all, "") != "BDC" {
		t.Log(all)
		t.Fatalf("hash all() wrong")
	}

	lb.RemoveAll()
	lb.Add("F", 1)
	all, ok := lb.All().([]string)
	if !ok || len(all) != 1 {
		t.Fatal("r all() wrong")
	}

	nodes = []string{"X", "Y"}
	ok = lb.Update(nodes)
	if ok != true {
		t.Fatal("r update wrong")
	}
	item = lb.Select()
	if item != "X" && item != "Y" {
		t.Fatal("r update wrong")
	}
}

func TestRandom_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := []string{"A", "B", "C", "D"}
	lb := NewRandom(nodes)

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

	if atomic.LoadInt64(&a) <= 200000 {
		t.Fatal("swrr wrong: a")
	}
	if atomic.LoadInt64(&b) <= 200000 {
		t.Fatal("swrr wrong: b")
	}
	if atomic.LoadInt64(&c) <= 200000 {
		t.Fatal("swrr wrong: c")
	}
	if atomic.LoadInt64(&d) <= 200000 {
		t.Fatal("swrr wrong: d")
	}
}

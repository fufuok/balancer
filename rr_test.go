package balancer

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	lb := NewRoundRobin()
	item := lb.Select()
	if item != "" {
		t.Fatalf("rr expected empty, actual %s", item)
	}

	lb.Add("A")
	item = lb.Select()
	if item != "A" {
		t.Fatalf("rr expected A, actual %s", item)
	}

	nodes := []string{"A", "B", "C", "D"}
	lb = NewRoundRobin(nodes)
	item = lb.Select()
	if item != "A" {
		t.Fatalf("rr expected A, actual %s", item)
	}
	item = lb.Select()
	if item != "B" {
		t.Fatalf("rr expected B, actual %s", item)
	}
	item = lb.Select()
	if item != "C" {
		t.Fatalf("rr expected C, actual %s", item)
	}
	item = lb.Select()
	if item != "D" {
		t.Fatalf("rr expected D, actual %s", item)
	}
	item = lb.Select()
	if item != "A" {
		t.Fatalf("rr expected A, actual %s", item)
	}

	count := make(map[string]int)
	for i := 0; i < 2000; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 500 || count["B"] != 500 || count["C"] != 500 || count["D"] != 500 {
		t.Fatal("rr wrong")
	}

	lb.Add("E", 10)
	for i := 0; i < 1000; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 700 || count["B"] != 700 || count["C"] != 700 || count["D"] != 700 || count["E"] != 200 {
		t.Fatal("rr add() wrong")
	}

	ok := lb.Remove("E")
	if ok != true {
		t.Fatal("rr remove() wrong")
	}
	for i := 0; i < 800; i++ {
		item := lb.Select()
		count[item]++
	}
	if count["A"] != 900 || count["B"] != 900 || count["C"] != 900 || count["D"] != 900 {
		t.Fatal("rr wrong")
	}

	lb.Add("A")
	lb.Add("A")
	lb.Add("C")
	lb.Add("A")
	all := lb.All().([]string)
	if strings.Join(all, "") != "ABCDAACA" {
		t.Fatal("rr all() wrong")
	}

	lb.Remove("C")
	all = lb.All().([]string)
	if strings.Join(all, "") != "ABDAACA" {
		t.Fatal("rr all() wrong")
	}

	lb.Remove("A", true)
	all = lb.All().([]string)
	if strings.Join(all, "") != "BDC" {
		t.Fatal("rr all() wrong")
	}

	lb.RemoveAll()
	lb.Add("F", 1)
	all, ok = lb.All().([]string)
	if !ok || len(all) != 1 {
		t.Fatal("rr all() wrong")
	}

	ok = lb.Update([]string{
		"X",
		"Y",
	})
	if ok != true {
		t.Fatal("rr update wrong")
	}
	item = lb.Select()
	if item != "X" {
		t.Fatal("rr update wrong")
	}
}

func TestRoundRobin_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := []string{"A", "B", "C", "D"}
	lb := NewRoundRobin(nodes)

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

	if atomic.LoadInt64(&a) != 250000 {
		t.Fatal("rr wrong: a")
	}
	if atomic.LoadInt64(&b) != 250000 {
		t.Fatal("rr wrong: b")
	}
	if atomic.LoadInt64(&c) != 250000 {
		t.Fatal("rr wrong: c")
	}
	if atomic.LoadInt64(&d) != 250000 {
		t.Fatal("rr wrong: d")
	}
}

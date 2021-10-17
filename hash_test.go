package balancer

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	lb := NewConsistentHash()
	item := lb.Select()
	if item != "" {
		t.Fatalf("hash expected empty, actual %s", item)
	}

	lb.Add("A")
	item = lb.Select()
	if item != "A" {
		t.Fatalf("hash expected A, actual %s", item)
	}

	nodes := []string{"A", "B", "C", "D"}
	lb = NewConsistentHash(nodes)
	item = lb.Select()
	if item != "B" {
		t.Fatalf("hash expected B, actual %s", item)
	}
	item = lb.Select()
	if item != "B" {
		t.Fatalf("hash expected B, actual %s", item)
	}
	item = lb.Select("192.168.1.100")
	if item != "A" {
		t.Fatalf("hash expected A, actual %s", item)
	}
	item = lb.Select("192.168.1.101")
	if item != "C" {
		t.Fatalf("hash expected C, actual %s", item)
	}
	item = lb.Select("192.168.1.102")
	if item != "D" {
		t.Fatalf("hash expected D, actual %s", item)
	}
	item = lb.Select("192.168.1.100")
	if item != "A" {
		t.Fatalf("hash expected A, actual %s", item)
	}
	item = lb.Select("2400:da00::6666")
	if item != "C" {
		t.Fatalf("hash expected C, actual %s", item)
	}

	for i := 0; i < 2000; i++ {
		item := lb.Select("192.168.1.100")
		if item != "A" {
			t.Fatalf("hash expected A, actual %s", item)
		}
	}

	lb.Add("E")
	lb.Add("C")
	lb.Add("E")
	lb.Remove("B")
	item = lb.Select("192.168.1.100")
	if item != "A" {
		t.Fatalf("hash expected A, actual %s", item)
	}

	all := lb.All().([]string)
	if strings.Join(all, "") != "ACDECE" {
		t.Fatalf("hash all() wrong")
	}

	lb.Remove("E")
	all = lb.All().([]string)
	if strings.Join(all, "") != "ACDCE" {
		t.Fatalf("hash all() wrong")
	}

	lb.Remove("C", true)
	all = lb.All().([]string)
	if strings.Join(all, "") != "ADE" {
		t.Fatalf("hash all() wrong")
	}

	lb.RemoveAll()
	lb.Add("F", 1)
	all, ok := lb.All().([]string)
	if !ok || len(all) != 1 {
		t.Fatal("hash all() wrong")
	}

	ok = lb.Update([]string{
		"X",
		"Y",
	})
	if ok != true {
		t.Fatal("hash update wrong")
	}
	item = lb.Select()
	if item != "Y" {
		t.Fatal("hash update wrong")
	}
	item = lb.Select()
	if item != "Y" {
		t.Fatal("hash update wrong")
	}
}

func TestConsistentHash_C(t *testing.T) {
	var c int64
	nodes := []string{"A", "B", "C", "D"}
	lb := NewConsistentHash(nodes)

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				switch lb.Select("192.168.1.7") {
				case "C":
					atomic.AddInt64(&c, 1)
				default:
				}
			}
		}()
	}
	wg.Wait()

	if atomic.LoadInt64(&c) != 1000000 {
		t.Fatalf("hash expected C == 1000000, actual C == %d, item: %s", atomic.LoadInt64(&c), lb.Select("192.168.1.7"))
	}
}

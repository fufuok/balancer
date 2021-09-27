package balancer

import (
	"testing"
)

func TestBalancer(t *testing.T) {
	lb := New(WeightedRoundRobin, nil)
	if lb.Name() != "WeightedRoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(SmoothWeightedRoundRobin, nil)
	if lb.Name() != "SmoothWeightedRoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(ConsistentHash, nil)
	if lb.Name() != "ConsistentHash" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(RoundRobin, nil)
	if lb.Name() != "RoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(Random, nil)
	if lb.Name() != "Random" {
		t.Fatal("balancer.New wrong")
	}

	lb.Add("A", 0)
	best := lb.Select()
	if best != "A" {
		t.Fatal("balancer select wrong")
	}
}

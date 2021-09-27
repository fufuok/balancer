package balancer

import (
	"testing"
)

func TestBalancer(t *testing.T) {
	lb := New(WeightedRoundRobin, nil, nil)
	if lb.Name() != "WeightedRoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(SmoothWeightedRoundRobin, nil, nil)
	if lb.Name() != "SmoothWeightedRoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(ConsistentHash, nil, nil)
	if lb.Name() != "ConsistentHash" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(RoundRobin, nil, nil)
	if lb.Name() != "RoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(Random, nil, nil)
	if lb.Name() != "Random" {
		t.Fatal("balancer.New wrong")
	}

	lb.Add("A", 0)
	best := lb.Select()
	if best != "A" {
		t.Fatal("balancer select wrong")
	}
}

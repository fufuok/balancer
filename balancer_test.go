package balancer

import (
	"testing"
)

func TestBalancer(t *testing.T) {
	lb := New(WeightedRoundRobin, nil, nil)
	if lb.Name() != "WeightedRoundRobin" || lb.Name() != WeightedRoundRobin.String() {
		t.Fatal("balancer.New wrong")
	}

	lb = New(SmoothWeightedRoundRobin, nil, nil)
	if lb.Name() != "SmoothWeightedRoundRobin" || lb.Name() != SmoothWeightedRoundRobin.String() {
		t.Fatal("balancer.New wrong")
	}

	lb = New(WeightedRand, nil, nil)
	if lb.Name() != "WeightedRand" || lb.Name() != WeightedRand.String() {
		t.Fatal("balancer.New wrong")
	}

	lb = New(ConsistentHash, nil, nil)
	if lb.Name() != "ConsistentHash" || lb.Name() != ConsistentHash.String() {
		t.Fatal("balancer.New wrong")
	}

	lb = New(RoundRobin, nil, nil)
	if lb.Name() != "RoundRobin" || lb.Name() != RoundRobin.String() {
		t.Fatal("balancer.New wrong")
	}

	lb = New(Random, nil, nil)
	if lb.Name() != "Random" || lb.Name() != Random.String() {
		t.Fatal("balancer.New wrong")
	}

	if Mode(777).String() != "" {
		t.Fatal("balancer name wrong")
	}

	lb.Add("A")
	best := lb.Select()
	if best != "A" {
		t.Fatal("balancer select wrong")
	}
}

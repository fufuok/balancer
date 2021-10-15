package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	var lb balancer.Balancer

	// SmoothWeightedRoundRobin and WeightedRoundRobin use map[string]int
	wNodes := map[string]int{
		"A": 5,
		"B": 1,
		"C": 1,
		"D": 0,
	}
	lb = balancer.New(balancer.WeightedRoundRobin, wNodes, nil)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.SmoothWeightedRoundRobin, wNodes, nil)
	fmt.Println("balancer name:", lb.Name())

	// RoundRobin / Random / ConsistentHash use []string
	nodes := []string{"A", "B", "C"}
	lb = balancer.New(balancer.ConsistentHash, nil, nodes)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.Random, nil, nodes)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.RoundRobin, nil, nodes)
	fmt.Println("balancer name:", lb.Name())

	// result of RoundRobin: A B C A B C A
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// same effect
	lb = balancer.New(balancer.SmoothWeightedRoundRobin, nil, nil)
	lb.Add("A", 1)
	lb.Select()

	// or like this
	lb = balancer.New(balancer.RoundRobin, nil, nil)
	lb.Update(nodes)
	lb.Select()
}

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
	lb = balancer.New(balancer.WeightedRoundRobin, wNodes)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.SmoothWeightedRoundRobin, wNodes)
	fmt.Println("balancer name:", lb.Name())

	// RoundRobin / Random / ConsistentHash use []string
	nodes := []string{"A", "B", "C"}
	lb = balancer.New(balancer.ConsistentHash, nodes)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.Random, nodes)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.RoundRobin, nodes)
	fmt.Println("balancer name:", lb.Name())

	// result of RoundRobin: A B C A B C A
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}

	// same effect
	lb = balancer.New(balancer.SmoothWeightedRoundRobin, nil)
	lb.Add("A", 1)
	lb.Select()

	// or like this
	lb = balancer.New(balancer.RoundRobin, nil)
	lb.Update(nodes)
	lb.Select()
}

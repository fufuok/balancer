package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	nodes := []string{"A", "B", "C"}
	// lb := balancer.New(balancer.RoundRobin, nodes)
	lb := balancer.NewRoundRobin(nodes)
	fmt.Println("balancer name:", lb.Name())

	// A B C A B
	for i := 0; i < 5; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// add an item to be selected
	lb.Add("E", 0)
	// output: C
	fmt.Println(lb.Select())
	// output: E
	fmt.Println(lb.Select())

	// get all items
	nodes = lb.All().([]string)
	// [A B C E]
	fmt.Printf("%+v\n", nodes)

	nodes = append(nodes, "F")
	// reinitialize the balancer items
	lb.Update(nodes)

	// A B C E F A B
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// reset the balancer
	lb.Reset()

	// remove an item
	lb.Remove("E")

	// remove all items
	lb.RemoveAll()
}

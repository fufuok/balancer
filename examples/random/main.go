package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	nodes := []string{"A", "B", "C"}
	// lb := balancer.New(balancer.Random, nodes)
	lb := balancer.NewRandom(nodes)
	fmt.Println("balancer name:", lb.Name())

	// random
	for i := 0; i < 9; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// add an item to be selected
	lb.Add("E")
	fmt.Println(lb.Select())

	// get all items
	nodes = lb.All().([]string)
	// [A B C E]
	fmt.Printf("%+v\n", nodes)

	nodes = append(nodes, "F")
	// reinitialize the balancer items
	lb.Update(nodes)

	// random
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// remove an item
	lb.Remove("E")

	// remove all items
	lb.RemoveAll()
}

package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	nodes := []string{"A", "B", "C"}
	// lb := balancer.New(balancer.ConsistentHash, nodes)
	lb := balancer.NewConsistentHash(nodes)
	fmt.Println("balancer name:", lb.Name())

	// B B B B B
	for i := 0; i < 5; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// add an item to be selected
	lb.Add("E")
	// output: B
	fmt.Println(lb.Select())
	// output: B
	fmt.Println(lb.Select())
	// output: C
	fmt.Println(lb.Select("1.2.3.4"))

	// get all items
	nodes = lb.All().([]string)
	// [A B C E]
	fmt.Printf("%+v\n", nodes)

	nodes = append(nodes, "F")
	// reinitialize the balancer items
	lb.Update(nodes)

	// B B B B B B B
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// output: F
	fmt.Println(lb.Select("1.2.3.4"))

	lb.Remove("F")

	// output: A
	fmt.Println(lb.Select("1.2.3.4"))

	// remove all items
	lb.RemoveAll()
}

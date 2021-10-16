package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	// Weighted node
	wNodes := map[string]int{
		"A": 5,
		"B": 3,
		"C": 1,
		"D": 0,
	}
	// lb := balancer.New(balancer.WeightedRand, wNodes)
	lb := balancer.NewWeightedRand(wNodes)
	fmt.Println("balancer name:", lb.Name())

	for i := 0; i < 9; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// add an item to be selected
	lb.Add("E", 20)
	fmt.Println(lb.Select())

	// get all items
	wNodes = lb.All().(map[string]int)
	// map[A:5 B:3 C:1 D:0 E:20]
	fmt.Printf("%+v\n", wNodes)

	wNodes["E"] = 5
	wNodes["D"] = 1
	wNodes["A"] = 0
	delete(wNodes, "B")
	// reinitialize the balancer items
	lb.Update(wNodes)

	for i := 0; i < 14; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// reset the balancer items weight
	lb.Reset()

	// remove an item
	lb.Remove("E")

	// remove all items
	lb.RemoveAll()
}

package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	// WeightedRoundRobin is the default balancer algorithm.
	fmt.Println("default balancer name:", balancer.Name())

	// reinitialize the balancer items.
	// D: will not select items with a weight of 0
	wNodes := map[string]int{
		"A": 5,
		"B": 3,
		"C": 1,
		"D": 0,
	}
	balancer.Update(wNodes)

	// result of smooth selection is similar to: A A A B A B A B C
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
	fmt.Println()

	// add an item to be selected
	balancer.Add("E", 20)
	// output: E
	fmt.Println(balancer.Select())

	// get all items
	wNodes = balancer.All().(map[string]int)
	// map[A:5 B:3 C:1 D:0 E:20]
	fmt.Printf("%+v\n", wNodes)

	wNodes["E"] = 5
	wNodes["D"] = 1
	wNodes["A"] = 0
	delete(wNodes, "B")
	// reinitialize the balancer items
	balancer.Update(wNodes)

	// when the weight difference is large, it is not smooth: E E E E D E C E E E E D E C
	for i := 0; i < 14; i++ {
		fmt.Print(balancer.Select(), " ")
	}
	fmt.Println()

	// reset the balancer items weight
	balancer.Reset()

	// remove an item
	balancer.Remove("E")

	// remove all items
	balancer.RemoveAll()
}

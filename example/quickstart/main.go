package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	// To be selected : Weighted
	wNodes := map[string]int{
		"A": 5,
		"B": 3,
		"C": 1,
		"D": 0,
	}
	balancer.Update(wNodes)

	// result of smooth selection is similar to: A A A B A B C A B
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
}

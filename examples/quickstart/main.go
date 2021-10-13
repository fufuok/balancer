package main

import (
	"fmt"

	"github.com/fufuok/balancer"
)

func main() {
	wNodes := map[string]int{
		"🍒": 5,
		"🍋": 3,
		"🍉": 1,
		"🥑": 0,
	}
	balancer.Update(wNodes)

	// result of smooth selection is similar to: 🍒 🍒 🍒 🍋 🍒 🍋 🍒 🍋 🍉
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
}

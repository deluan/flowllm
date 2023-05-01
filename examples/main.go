package main

import (
	"fmt"
	"os"
	"strconv"
)

type example struct {
	name        string
	description string
	fn          func()
}

var examples []example

func registerExample(name, description string, fn func()) {
	examples = append(examples, example{name: name, description: description, fn: fn})
}

func main() {
	if len(os.Args) == 1 {
		println("usage: examples <example_number>")
		println("available examples:")
		for i, e := range examples {
			fmt.Printf(" %d) %-20s - %s\n", i, e.name, e.description)
		}
		os.Exit(0)
	}
	num, err := strconv.Atoi(os.Args[1])
	if err != nil || num < 0 || num >= len(examples) {
		println("invalid example number")
		os.Exit(1)
	}

	e := examples[num]
	e.fn()
	os.Exit(0)
}

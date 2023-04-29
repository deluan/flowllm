package main

import (
	"fmt"
	"os"
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
		println("usage: examples <example_name>")
		println("available examples:")
		for _, e := range examples {
			fmt.Printf("  %s: %s\n", e.name, e.description)
		}
		os.Exit(0)
	}
	name := os.Args[1]
	for _, e := range examples {
		if e.name == name {
			e.fn()
			os.Exit(0)
		}
	}
	println("unknown example:", name)
}

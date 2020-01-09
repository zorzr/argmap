package main

import (
	"fmt"

	"github.com/zorzr/argmap"
)

func main() {
	parser := argmap.NewArgsParser("Greeter", "Says hello")
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi"})
	parser.NewBoolFlag(argmap.BoolFlag{Name: "spanish"})

	aMap, _ := parser.Parse()
	name, _ := argmap.GetListValue(aMap, "hello", 0)
	if argmap.GetBool(aMap, "spanish") {
		fmt.Printf("Hola %s\n", name)
	} else {
		fmt.Printf("Hello %s\n", name)
	}
}

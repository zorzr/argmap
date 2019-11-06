package main

import (
	"fmt"

	"github.com/zorzr/argmap"
)

func main() {
	parser := argmap.NewArgsParser("Introducing myself", "Tells you how I'm called")
	parser.NewPositionalArg(argmap.PositionalArg{Name: "name", Required: true})
	parser.NewPositionalArg(argmap.PositionalArg{Name: "surname"})

	aMap, err := parser.Parse()
	if err != nil {
		parser.ReportError(err)
	}

	// The "name" positional is required, we are sure it's in the map
	name, _ := argmap.GetPositional(aMap, "name")
	if surname, err := argmap.GetPositional(aMap, "surname"); err == nil {
		fmt.Printf("My name is %s, %s %s\n", surname, name, surname)
	} else {
		fmt.Printf("My name is %s\n", name)
	}
}

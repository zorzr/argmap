package main

import (
	"fmt"

	"github.com/zorzr/argmap"
)

func main() {
	parser := argmap.NewArgsParser("Greeter", "Says hello")
	parser.NewStringFlag(argmap.StringFlag{
		Name:  "hello",
		Short: "hi",
		Vars:  []string{"urname"},
		Help:  "greets you",
	})

	aMap, err := parser.Parse()
	if err != nil {
		parser.ReportError(err)
	}

	name, _ := argmap.GetListValue(aMap, "hello", 0)
	fmt.Printf("Hello %s\n", name)
}

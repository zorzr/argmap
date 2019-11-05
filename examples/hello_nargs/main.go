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
		NArgs: 2,
	})

	aMap, _ := parser.Parse()
	names, err := argmap.GetSFArray(aMap, "hello")
	if err != nil {
		parser.ReportError(err)
	}

	fmt.Printf("Hello %s and %s\n", names[0], names[1])
}

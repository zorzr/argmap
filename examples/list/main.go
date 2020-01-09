package main

import (
	"fmt"

	"github.com/zorzr/argmap"
)

func main() {
	parser := argmap.NewArgsParser("Shopping list", "Records what you need to buy")
	parser.NewListFlag(argmap.ListFlag{Name: "list", Short: "l", Var: "item", Help: "items to be bought"})
	parser.NewBoolFlag(argmap.BoolFlag{Short: "c", Help: "look for the cheapest goods"})

	aMap, err := parser.Parse()
	if err != nil {
		parser.ReportError(err)
	} else if !argmap.IsPresent(aMap, "list") {
		parser.ReportError(fmt.Errorf("Please insert a list of items"))
	}

	list, _ := argmap.GetList(aMap, "list")

	s := "s"
	if len(list) == 1 {
		s = ""
	}
	fmt.Printf("[SHOPPING LIST - %d item%s]\n", len(list), s)
	if argmap.GetBool(aMap, "c") {
		fmt.Println("[Keep an eye on the budget!]")
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

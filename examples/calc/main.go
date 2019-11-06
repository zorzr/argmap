package main

import (
	"fmt"
	"strconv"

	"github.com/zorzr/argmap"
)

func conv(values []string) (float64, float64, error) {
	a, err1 := strconv.ParseFloat(values[0], 64)
	b, err2 := strconv.ParseFloat(values[1], 64)

	if err1 != nil {
		return 0, 0, err1
	} else if err2 != nil {
		return 0, 0, err2
	}
	return a, b, nil
}

func main() {
	parser := argmap.NewArgsParser("Calculator", "Solves all your problems")
	parser.NewPositionalArg(argmap.PositionalArg{Name: "action", Required: true})
	parser.NewStringFlag(argmap.StringFlag{Short: "o", NArgs: 2})

	aMap, err := parser.Parse()
	if err != nil {
		parser.ReportError(err)
	}

	action, _ := argmap.GetPositional(aMap, "action")
	if operands, err := argmap.GetSFArray(aMap, "o"); err == nil {
		if a, b, err := conv(operands); err == nil {
			switch action {
			case "add":
				fmt.Println(a + b)
			case "sub":
				fmt.Println(a - b)
			case "prod":
				fmt.Println(a * b)
			case "div":
				fmt.Println(a / b)
			default:
				fmt.Println("Error: unknown operation")
			}
		} else {
			fmt.Println("Error: operands are not numbers")
		}
	} else {
		fmt.Println("Error: not enough operands")
	}
}

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/zorzr/argmap"
)

func initParser() *argmap.ArgsParser {
	parser := argmap.NewArgsParser("Printer", "Shows you something from command line")
	parser.NewCommand(argmap.CommandParams{Name: "hello", Help: "greets the user"})
	printer, _ := parser.NewCommand(argmap.CommandParams{Name: "print", Help: "prints a string or the content of a file"})

	str, _ := printer.NewSubcommand(argmap.CommandParams{Name: "string", Help: "prints a string"})
	str.NewPositionalArg(argmap.PositionalArg{Name: "input", Help: "the input string", Required: true})

	file, _ := printer.NewSubcommand(argmap.CommandParams{Name: "file", Help: "prints the content of a file"})
	file.NewPositionalArg(argmap.PositionalArg{Name: "path", Help: "location of the file to be read", Required: true})

	return &parser
}

func main() {
	parser := initParser()
	aMap, err := parser.Parse()
	if err != nil {
		parser.ReportError(err)
	}

	cmd, cmdMap, err := argmap.GetCommandMap(aMap)
	if err != nil {
		parser.ReportError(fmt.Errorf("Please type a command to be executed"))
	}

	switch cmd {
	case "hello":
		fmt.Println("Nice to meet you!")
	case "print":
		sub, subMap, err := argmap.GetCommandMap(cmdMap)
		if err != nil {
			parser.ReportError(fmt.Errorf("Missing subcommand for command 'print'"))
		}

		switch sub {
		case "string":
			input, _ := argmap.GetPositional(subMap, "input")
			fmt.Println(input)
		case "file":
			path, _ := argmap.GetPositional(subMap, "path")
			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%s\n", data)
		}
	}
}

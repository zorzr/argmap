// Package argmap is a simple command line argument parser providing high level features
// and wide freedom of usage. The user inputs are stored in a map allowing to easily
// gather, analyze and retrieve them.
//
// Along with positional arguments, different types of flags can be defined according to
// your needs. Nothing is arbitrarily made without giving you a choice, even in case
// something goes wrong: there are built-in functions to easily report possible
// errors, but they won't be automatically called if you do not want so!
package argmap

import (
	"fmt"
	"os"
	"sort"
)

// HelpMessageGenerator type used to allow customizable help messages
type HelpMessageGenerator func(*ArgsParser, []*Command) string

// ArgsParser stores the list of possible arguments
type ArgsParser struct {
	Name        string
	Description string
	argsList    []Argument
	helpGen     HelpMessageGenerator
}

// NewArgsParser function to return an initialized struct
func NewArgsParser(name, descr string) ArgsParser {
	var helpArg = []Argument{HelpFlag{"shows help message and exits"}}

	return ArgsParser{
		Name:        name,
		Description: descr,
		argsList:    helpArg,
		helpGen:     DefaultHelp,
	}
}

// DefaultHelp produces the standard complete help message for the program
func DefaultHelp(p *ArgsParser, cmdTrace []*Command) string {
	help := fmt.Sprintf("%s\n%s\n", p.Name, p.Description)

	if cmdTrace == nil || len(cmdTrace) == 0 {
		// PROGRAM HELP
		p.SortArgsList()
		length := len(p.argsList)
		argsHelp := make([][]string, length)

		maxLeftLen := 0
		commandsIndex := length
		for i := 0; i < length; i++ {
			argsHelp[i] = p.argsList[i].GetHelpStrings()
			if len(argsHelp[i][0]) > maxLeftLen {
				maxLeftLen = len(argsHelp[i][0])
			}

			if commandsIndex == length && p.argsList[i].getOrder() == orderCommand {
				commandsIndex = i
			}
		}

		if maxLeftLen > 40 {
			maxLeftLen = 40
		}

		help += "\nArguments:\n"
		for i := 0; i < length; i++ {
			if i == commandsIndex {
				help += "\nCommands:\n"
			}

			argStr := argsHelp[i][0]
			for len(argStr) <= maxLeftLen {
				argStr += " "
			}
			help += fmt.Sprintf("  %s %s\n", argStr, argsHelp[i][1])

			if i == length-1 && commandsIndex < length {
				help += "Type -h or --help after a command for more details\n"
			}
		}
	} else {
		// COMMAND HELP
		traceString := ""
		for i := len(cmdTrace) - 1; i >= 0; i-- {
			traceString += fmt.Sprintf(" %s", cmdTrace[i].GetID())
		}

		help += fmt.Sprintf("\nReference: %s\n", traceString)
		help += cmdTrace[0].GenerateHelp()
	}

	return help
}

func parseArgs(args []string, argsList []Argument) (map[string]interface{}, error) {
	var argsMap = make(map[string]interface{})

	var posIndex = 0
	var posArgs = []int{}
	var reqPos = []string{}

	var reprMap = make(map[string]*Argument)
	for i, a := range argsList {
		if a.getOrder() <= orderPositionalOpt {
			posArgs = append(posArgs, i)
			if a.getOrder() == orderPositionalReq {
				reqPos = append(reqPos, a.GetID())
			}
			continue
		}

		for _, r := range a.Represent() {
			reprMap[r] = &argsList[i]
		}
	}

	n := len(args)
	for i := 0; i < n; i++ {
		if arg, ok := reprMap[args[i]]; ok {
			switch (*arg).getOrder() {
			// STRINGFLAG
			case orderStringFlag:
				flag := (*arg).(StringFlag)

				if i+flag.NArgs >= n {
					return nil, fmt.Errorf("Error: incorrect arguments usage")
				}

				var j int
				var values = make([]string, flag.NArgs)
				for j = 0; j < flag.NArgs; j++ {
					values[j] = args[i+j+1] // TODO: check if arg is in reprMap?
				}
				i += j

				argsMap[flag.GetID()] = values

			// BOOLFLAG
			case orderBoolFlag:
				flag := (*arg).(BoolFlag)
				argsMap[flag.GetID()] = true

			// HELPFLAG
			case orderHelpFlag:
				argsMap = map[string]interface{}{"help": true}
				return argsMap, nil

			// COMMAND
			case orderCommand:
				cmd := (*arg).(*Command)
				cmdMap, err := cmd.parseArgs(args[i+1:])
				if err != nil {
					return nil, err
				}

				if GetBool(cmdMap, "help") {
					trace := []*Command{}
					if IsPresent(cmdMap, "trace") {
						trace = cmdMap["trace"].([]*Command)
					}
					trace = append(trace, cmd)
					cmdMap["trace"] = trace
					return cmdMap, nil
				}

				argsMap[cmd.GetID()] = cmdMap
				i = n
			}
		} else {
			// POSITIONAL ARGUMENTS
			if len(posArgs) == posIndex {
				return nil, fmt.Errorf("Error: unrecognized argument '%s'", args[i])
			}

			pArg := argsList[posArgs[posIndex]].(PositionalArg)
			argsMap[pArg.GetID()] = args[i]
			posIndex++
		}
	}

	// We check if any required positional argument is missing
	// TODO: possible implementation for required flags
	for _, pos := range reqPos {
		if !IsPresent(argsMap, pos) {
			return nil, fmt.Errorf("Error: missing required positional argument '%s'", pos)
		}
	}

	return argsMap, nil
}

// GenerateHelp produces the help string to be shown when the "-h" or "--help" flags are inserted by the user.
func (p *ArgsParser) GenerateHelp() string {
	return p.helpGen(p, nil)
}

// GenerateCommandHelp produces the help string for a Command to be shown when the "-h" or "--help" flags are inserted by the user.
func (p *ArgsParser) GenerateCommandHelp(cmdTrace []*Command) string {
	return p.helpGen(p, cmdTrace)
}

// SetHelpGenerator accepts a function to be used to generate a custom help message
// to be shown when the "-h" or "--help" flags are inserted by the user.
func (p *ArgsParser) SetHelpGenerator(h HelpMessageGenerator) {
	p.helpGen = h
}

// SetHelpFlagMessage accepts a string to be used in the program help with that HelpFlag
func (p *ArgsParser) SetHelpFlagMessage(m string) {
	for i, a := range p.argsList {
		if a.getOrder() == orderHelpFlag {
			p.argsList[i] = HelpFlag{Help: m}
			return
		}
	}
}

// PrintHelp shows the complete help message for the program
func (p *ArgsParser) PrintHelp() {
	help := p.helpGen(p, nil)
	fmt.Println(help)
}

// PrintCommandHelp shows the complete help message for a program command
func (p *ArgsParser) PrintCommandHelp(cmdTrace []*Command) {
	help := p.helpGen(p, cmdTrace)
	fmt.Println(help)
}

// ReportError prints the passed error's message, shows the correct usage and quits
func (p *ArgsParser) ReportError(err error) {
	fmt.Println(err.Error())
	p.PrintHelp()
	os.Exit(0)
}

// Parse function returns a map with argument values
func (p *ArgsParser) Parse() (map[string]interface{}, error) {
	p.SortArgsList()
	argsMap, err := parseArgs(os.Args[1:], p.argsList)
	if err != nil {
		return nil, err
	}

	if GetBool(argsMap, "help") {
		if !IsPresent(argsMap, "trace") {
			p.PrintHelp()
		} else {
			cmdTrace := argsMap["trace"].([]*Command)
			p.PrintCommandHelp(cmdTrace)
		}
		os.Exit(0)
	}

	return argsMap, nil
}

// NewStringFlag checks the fields for consistency and inserts the new flag
func (p *ArgsParser) NewStringFlag(f StringFlag) error {
	if f.Name == "" && f.Short == "" {
		return fmt.Errorf("Error: at least one identifier must be specified")
	}

	if f.NArgs < 1 {
		f.NArgs = 1
	}

	if len(f.Vars) < f.NArgs {
		for len(f.Vars) < f.NArgs {
			f.Vars = append(f.Vars, "value")
		}
	} else if len(f.Vars) > f.NArgs {
		return fmt.Errorf("Error: too many value names specified (expected %d, got %d)", f.NArgs, len(f.Vars))
	}

	err := checkIdentifiers(p.argsList, f)
	if err != nil {
		return err
	}

	p.argsList = append(p.argsList, f)
	return nil
}

// NewBoolFlag checks the flag representations and inserts the new flag
func (p *ArgsParser) NewBoolFlag(f BoolFlag) error {
	if f.Name == "" && f.Short == "" {
		return fmt.Errorf("Error: at least one identifier must be specified")
	}

	err := checkIdentifiers(p.argsList, f)
	if err != nil {
		return err
	}

	p.argsList = append(p.argsList, f)
	return nil
}

// NewPositionalArg checks the argument identifier and inserts it
func (p *ArgsParser) NewPositionalArg(a PositionalArg) error {
	if a.Name == "" {
		return fmt.Errorf("Error: unspecified argument name")
	}

	err := checkIdentifiers(p.argsList, a)
	if err != nil {
		return err
	}

	p.argsList = append(p.argsList, a)
	return nil
}

// NewCommand checks the argument identifier and inserts it
func (p *ArgsParser) NewCommand(param CommandParams) (*Command, error) {
	if param.Name == "" {
		return nil, fmt.Errorf("Error: unspecified command name")
	}

	c := &Command{
		name:     param.Name,
		Help:     param.Help,
		argsList: []Argument{HelpFlag{"shows command help and exits"}},
		helpGen:  DefaultCommandHelp,
	}

	err := checkIdentifiers(p.argsList, c)
	if err != nil {
		return nil, err
	}

	p.argsList = append(p.argsList, c)
	return c, nil
}

// SortArgsList sorts the list of arguments according to their type. This allows to
// keep the list clearly ordered and avoid possible mistakes: for example, if an
// optional argument is inserted between two required ones and the user inserts only
// two positionals, we can't distinguish whether the user wanted to insert the
// optional and forgot the required or instead willingly ignored the optional one.
//
// Sorting the array of inserted arguments solves the ambiguity.
// The best design choice, however, would be to avoid too many positionals and
// handle the presence/absence of a StringFlag in the map after the parsing.
//  Order of relevance:
//      1. PositionalArg (required)
//      2. PositionalArg (optional)
//      3. StringFlag
//      4. BoolFlag
//      5. HelpFlag
func (p *ArgsParser) SortArgsList() {
	sort.Slice(p.argsList, func(i, j int) bool {
		return p.argsList[i].getOrder() < p.argsList[j].getOrder()
	})
}

// GetArgsList returns a copy of the argument list to allow the generation of custom help messages
func (p *ArgsParser) GetArgsList() []Argument {
	arr := make([]Argument, len(p.argsList))
	copy(arr, p.argsList)
	return arr
}

/************************************************************/
func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// TODO: goroutines
func checkIdentifiers(argsList []Argument, b Argument) error {
	for _, a := range argsList {
		if a.GetID() == b.GetID() {
			return fmt.Errorf("Error: identifier '%s' already exists", b.GetID())
		}
		for _, r := range b.Represent() {
			if contains(a.Represent(), r) {
				return fmt.Errorf("Error: representation '%s' already exists", r)
			}
		}
	}
	return nil
}

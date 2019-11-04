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
type HelpMessageGenerator func(*ArgsParser) string

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
func DefaultHelp(p *ArgsParser) string {
	help := fmt.Sprintf("%s\n%s\n\nArguments:\n", p.Name, p.Description)
	argsHelp := make([][]string, len(p.argsList))
	p.SortArgsList()

	maxLeftLen := 0
	for i := 0; i < len(p.argsList); i++ {
		argsHelp[i] = p.argsList[i].GetHelpStrings()
		if len(argsHelp[i][0]) > maxLeftLen {
			maxLeftLen = len(argsHelp[i][0])
		}
	}

	if maxLeftLen > 40 {
		maxLeftLen = 40
	}

	for i := 0; i < len(p.argsList); i++ {
		argStr := argsHelp[i][0]
		for len(argStr) <= maxLeftLen {
			argStr += " "
		}

		help += fmt.Sprintf("  %s %s\n", argStr, argsHelp[i][1])
	}
	return help
}

// GenerateHelp produces the help string to be shown when the "-h" or "--help" flags are inserted by the user.
func (p *ArgsParser) GenerateHelp() string {
	return p.helpGen(p)
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
	help := p.helpGen(p)
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
	var argsMap = make(map[string]interface{})
	p.SortArgsList()

	var nPosArg = 0
	var posArgs = []int{}

	var reprMap = make(map[string]*Argument)
	for i, a := range p.argsList {
		if a.getOrder() <= orderPositionalOpt {
			nPosArg++
			posArgs = append(posArgs, i)
			continue
		}

		for _, r := range a.Represent() {
			reprMap[r] = &p.argsList[i]
		}
	}

	n := len(os.Args)
	for i := 1; i < n; i++ { // we ignore the program executable name
		if arg, ok := reprMap[os.Args[i]]; ok {
			if (*arg).getOrder() == orderStringFlag {
				flag := (*arg).(StringFlag)

				if i+flag.NArgs >= n {
					return nil, fmt.Errorf("Error: incorrect arguments usage")
				}

				// TODO: can replace these blocks with argument-specific functions?
				values := make([]string, flag.NArgs)

				var j int
				for j = 0; j < flag.NArgs; j++ {
					values[j] = os.Args[i+j+1]
				}
				i += j

				argsMap[flag.GetID()] = values
			} else if (*arg).getOrder() == orderBoolFlag {
				flag := (*arg).(BoolFlag)
				argsMap[flag.GetID()] = true
			} else if (*arg).getOrder() == orderHelpFlag {
				p.PrintHelp()
				os.Exit(0)
			}
		} else {
			if nPosArg == 0 {
				return nil, fmt.Errorf("Error: unrecognized argument '%s'", os.Args[i])
			}

			pArg := p.argsList[posArgs[len(posArgs)-nPosArg]].(PositionalArg)
			argsMap[pArg.GetID()] = os.Args[i]
			nPosArg--
		}
	}

	// We check if any required positional argument is missing
	// TODO: can use IDs instead? (=> possible implementation for required flags)
	if nPosArg > 0 {
		for i := nPosArg; i > 0; i-- {
			pArg := p.argsList[posArgs[len(posArgs)-i]].(PositionalArg)
			if pArg.Required {
				return nil, fmt.Errorf("Error: missing required positional argument '%s'", pArg.Name)
			}
		}
	}

	return argsMap, nil
}

// newArgument appends the argument passed as input to the array
func (p *ArgsParser) newArgument(a Argument) {
	p.argsList = append(p.argsList, a)
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

	p.newArgument(f)
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

	p.newArgument(f)
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

	p.newArgument(a)
	return nil
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

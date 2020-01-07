package argmap

import (
	"fmt"
	"sort"
	"strings"
)

// CommandHelpGenerator type used to allow customizable help for commands
type CommandHelpGenerator func(*Command) string

// Command is both a type of argument and a parser of what comes after it
type Command struct {
	name     string
	Help     string
	argsList []Argument
	helpGen  CommandHelpGenerator
}

// CommandParams used for commands initialization
type CommandParams struct {
	Name string
	Help string
}

// GetID returns the identifier of the command
func (c Command) GetID() string {
	return c.name
}

// Represent returns the name of the command
func (c Command) Represent() []string {
	return []string{c.name}
}

// GetHelpStrings returns the two hand sides of the help message
func (c Command) GetHelpStrings() []string {
	return []string{c.name, c.Help}
}

// Defines the priority of the argument for sorting (also used to determine the argument type)
func (c Command) getOrder() int {
	return orderCommand
}

/********************************************************************/

// GenerateHelp produces the help string to be shown when the "-h" or "--help" flags are inserted by the user.
func (c *Command) GenerateHelp() string {
	return c.helpGen(c)
}

// SetHelpGenerator accepts a function to be used to generate a custom help message
// to be shown when the "-h" or "--help" flags are inserted by the user.
func (c *Command) SetHelpGenerator(h CommandHelpGenerator) {
	c.helpGen = h
}

// SetHelpFlagMessage accepts a string to be used in the program help with that HelpFlag
func (c *Command) SetHelpFlagMessage(m string) {
	for i, a := range c.argsList {
		if a.getOrder() == orderHelpFlag {
			c.argsList[i] = HelpFlag{Help: m}
			return
		}
	}
}

// SortArgsList sorts the list of arguments according to their type.
func (c *Command) SortArgsList() {
	sort.Slice(c.argsList, func(i, j int) bool {
		return c.argsList[i].getOrder() < c.argsList[j].getOrder()
	})
}

// DefaultCommandHelp produces a part of the help message for the command to be printed by the ArgsParser
func DefaultCommandHelp(c *Command) string {
	c.SortArgsList()
	length := len(c.argsList)
	argsHelp := make([][]string, length)

	maxLeftLen := 0
	subcommandsIndex := length
	for i := 0; i < length; i++ {
		argsHelp[i] = c.argsList[i].GetHelpStrings()
		if len(argsHelp[i][0]) > maxLeftLen {
			maxLeftLen = len(argsHelp[i][0])
		}

		if subcommandsIndex == length && c.argsList[i].getOrder() == orderCommand {
			subcommandsIndex = i
		}
	}

	if maxLeftLen > 40 {
		maxLeftLen = 40
	}

	help := fmt.Sprintf("    %s   %s\n\nArguments:\n", c.name, c.Help)
	for i := 0; i < length; i++ {
		if i == subcommandsIndex {
			help += "\nSubcommands:\n"
		}

		argStr := argsHelp[i][0]
		for len(argStr) <= maxLeftLen {
			argStr += " "
		}
		help += fmt.Sprintf("    %s %s\n", argStr, argsHelp[i][1])

		if i == length-1 && subcommandsIndex < length {
			help += "Type -h or --help after a command for more details\n"
		}
	}

	return help
}

// GetArgsList returns a copy of the argument list to be used for the production of custom helps
func (c *Command) GetArgsList() []Argument {
	arr := make([]Argument, len(c.argsList))
	copy(arr, c.argsList)
	return arr
}

/***************************************************************/

// NewStringFlag checks the fields for consistency and inserts the new flag
func (c *Command) NewStringFlag(f StringFlag) error {
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

	err := checkIdentifiers(c.argsList, f)
	if err != nil {
		return err
	}

	c.argsList = append(c.argsList, f)
	return nil
}

// NewBoolFlag checks the flag representations and inserts the new flag
func (c *Command) NewBoolFlag(f BoolFlag) error {
	if f.Name == "" && f.Short == "" {
		return fmt.Errorf("Error: at least one identifier must be specified")
	}

	err := checkIdentifiers(c.argsList, f)
	if err != nil {
		return err
	}

	c.argsList = append(c.argsList, f)
	return nil
}

// NewPositionalArg checks the argument identifier and inserts it
func (c *Command) NewPositionalArg(a PositionalArg) error {
	if a.Name == "" {
		return fmt.Errorf("Error: unspecified argument name")
	}

	err := checkIdentifiers(c.argsList, a)
	if err != nil {
		return err
	}

	c.argsList = append(c.argsList, a)
	return nil
}

// NewSubcommand checks the argument identifier and inserts it
func (c *Command) NewSubcommand(param CommandParams) (*Command, error) {
	if param.Name == "" {
		return nil, fmt.Errorf("Error: unspecified subcommand name")
	}

	sc := &Command{
		name:     param.Name,
		Help:     param.Help,
		argsList: []Argument{HelpFlag{"shows command help and exits"}},
		helpGen:  DefaultCommandHelp,
	}

	err := checkIdentifiers(c.argsList, sc)
	if err != nil {
		return nil, err
	}

	c.argsList = append(c.argsList, sc)
	return sc, nil
}

/******************************************************************/

func (c *Command) parseArgs(args []string) (map[string]interface{}, error) {
	c.SortArgsList()
	argsMap, err := parseArgs(args, c.argsList)
	if err != nil {
		placeholder := "[*]"
		errorString := err.Error()
		if strings.Contains(errorString, placeholder) {
			errorString = strings.Replace(errorString, placeholder, fmt.Sprintf("%s%s ", placeholder, c.name), 1)
			return nil, fmt.Errorf(errorString)
		}
		return nil, fmt.Errorf("%s for command '%s%s'", errorString, placeholder, c.name)
	}
	return argsMap, nil
}

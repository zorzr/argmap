package argmap

import (
	"fmt"
)

// Argument interface defines the basic methods an argument struct must have
//  GetID()             returns the identifier of the argument to be used in the map
//  GetHelpStrings()    returns the two sides of the help message (see declarations for details)
//  Represent()         eventual representations of the flag in the user inputs (e.g.: "-h", "--help")
type Argument interface {
	GetID() string
	GetHelpStrings() []string
	Represent() []string
	getOrder() int
}

const orderPositionalReq = 1
const orderPositionalOpt = 2
const orderStringFlag = 3
const orderListFlag = 4
const orderBoolFlag = 5
const orderHelpFlag = 9
const orderCommand = 10

/************************************************************/

// StringFlag argument
type StringFlag struct {
	Name  string
	Short string
	NArgs int
	Vars  []string
	Help  string
}

// GetID returns the identifier of the argument
func (f StringFlag) GetID() string {
	if f.Name != "" {
		return f.Name
	}
	return f.Short
}

// ShortArg returns short flag
func (f StringFlag) ShortArg() string {
	return "-" + f.Short
}

// LongArg returns full name flag
func (f StringFlag) LongArg() string {
	return "--" + f.Name
}

// Represent returns possible argument representations
func (f StringFlag) Represent() []string {
	if f.Name != "" && f.Short != "" {
		return []string{f.ShortArg(), f.LongArg()}
	} else if f.Name != "" {
		return []string{f.LongArg()}
	} else {
		return []string{f.ShortArg()}
	}
}

// GetHelpStrings returns the two hand sides of the help message
//  Example:  ["-a, --arg metavar1 metavar2", "this is an example of help message"]
func (f StringFlag) GetHelpStrings() []string {
	metaVars := ""
	for _, s := range f.Vars {
		metaVars += fmt.Sprintf("%s ", s)
	}

	var repr string
	if f.Name != "" && f.Short != "" {
		repr = fmt.Sprintf("%s, %s", f.ShortArg(), f.LongArg())
	} else if f.Name == "" {
		repr = f.ShortArg()
	} else {
		repr = f.LongArg()
	}

	leftHand := fmt.Sprintf("%s %s", repr, metaVars)
	return []string{leftHand, f.Help}
}

// Defines the priority of the argument for sorting (also used to determine the argument type)
func (f StringFlag) getOrder() int {
	return orderStringFlag
}

/*******************************************************/

// ListFlag argument
type ListFlag struct {
	Name  string
	Short string
	Var   string
	Help  string
}

// GetID returns the identifier of the argument
func (f ListFlag) GetID() string {
	if f.Name != "" {
		return f.Name
	}
	return f.Short
}

// ShortArg returns short flag
func (f ListFlag) ShortArg() string {
	return "-" + f.Short
}

// LongArg returns full name flag
func (f ListFlag) LongArg() string {
	return "--" + f.Name
}

// Represent returns possible argument representations
func (f ListFlag) Represent() []string {
	if f.Name != "" && f.Short != "" {
		return []string{f.ShortArg(), f.LongArg()}
	} else if f.Name != "" {
		return []string{f.LongArg()}
	} else {
		return []string{f.ShortArg()}
	}
}

// GetHelpStrings returns the two hand sides of the help message
//  Example:  ["-a, --arg metavar1 metavar2", "this is an example of help message"]
func (f ListFlag) GetHelpStrings() []string {
	var repr string
	if f.Name != "" && f.Short != "" {
		repr = fmt.Sprintf("%s, %s", f.ShortArg(), f.LongArg())
	} else if f.Name == "" {
		repr = f.ShortArg()
	} else {
		repr = f.LongArg()
	}

	leftHand := fmt.Sprintf("%s %s %s... ", repr, f.Var, f.Var)
	return []string{leftHand, f.Help}
}

// Defines the priority of the argument for sorting (also used to determine the argument type)
func (f ListFlag) getOrder() int {
	return orderListFlag
}

/************************************************************/

// BoolFlag argument
type BoolFlag struct {
	Name  string
	Short string
	Help  string
}

// GetID returns the identifier of the argument
func (f BoolFlag) GetID() string {
	if f.Name != "" {
		return f.Name
	}
	return f.Short
}

// ShortArg returns short flag
func (f BoolFlag) ShortArg() string {
	return "-" + f.Short
}

// LongArg returns full name flag
func (f BoolFlag) LongArg() string {
	return "--" + f.Name
}

// Represent returns possible argument representations
func (f BoolFlag) Represent() []string {
	if f.Name != "" && f.Short != "" {
		return []string{f.ShortArg(), f.LongArg()}
	} else if f.Name != "" {
		return []string{f.LongArg()}
	} else {
		return []string{f.ShortArg()}
	}
}

// GetHelpStrings returns the two hand sides of the help message
//  Example:  ["-b, --bool", "this is an example of help message"]
func (f BoolFlag) GetHelpStrings() []string {
	var leftHand string
	if f.Name != "" && f.Short != "" {
		leftHand = fmt.Sprintf("%s, %s", f.ShortArg(), f.LongArg())
	} else if f.Name == "" {
		leftHand = f.ShortArg()
	} else {
		leftHand = f.LongArg()
	}

	return []string{leftHand, f.Help}
}

// Defines the priority of the argument for sorting (also used to determine the argument type)
func (f BoolFlag) getOrder() int {
	return orderBoolFlag
}

/************************************************************/

// PositionalArg argument
type PositionalArg struct {
	Name     string
	Help     string
	Required bool
}

// GetID returns the identifier of the argument
func (a PositionalArg) GetID() string {
	return a.Name
}

// MetaArg returns a representation of the argument
//  Example:  required [optional]
func (a PositionalArg) MetaArg() string {
	if a.Required {
		return a.Name
	}
	return fmt.Sprintf("[%s]", a.Name)
}

// Represent returns no representations
// We do not look for a predefined string (like "--flag")
func (a PositionalArg) Represent() []string {
	return []string{}
}

// GetHelpStrings returns the two hand sides of the help message
//  Example:	 required	example of help message (f.Help)
//  Example:	 [optional]	example of help message (f.Help)
func (a PositionalArg) GetHelpStrings() []string {
	return []string{a.MetaArg(), a.Help}
}

// Defines the priority of the argument for sorting (also used to determine the argument type)
func (a PositionalArg) getOrder() int {
	if a.Required {
		return orderPositionalReq
	}
	return orderPositionalOpt
}

/************************************************************/

// HelpFlag argument
type HelpFlag struct {
	Help string
}

// GetID returns the identifier of the argument
func (f HelpFlag) GetID() string {
	return "help"
}

// ShortArg returns short flag
func (f HelpFlag) ShortArg() string {
	return "-h"
}

// LongArg returns full name flag
func (f HelpFlag) LongArg() string {
	return "--help"
}

// Represent returns possible argument representations
func (f HelpFlag) Represent() []string {
	return []string{f.ShortArg(), f.LongArg()}
}

// GetHelpStrings returns the two hand sides of the help message
//  Example: ["-h, --help",  "this is an example of help message"]
func (f HelpFlag) GetHelpStrings() []string {
	leftHand := fmt.Sprintf("%s, %s", f.ShortArg(), f.LongArg())
	return []string{leftHand, f.Help}
}

// Defines the priority of the argument for sorting (also used to determine the argument type)
func (f HelpFlag) getOrder() int {
	return orderHelpFlag
}

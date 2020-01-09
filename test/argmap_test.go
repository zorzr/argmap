package test

import (
	"os"
	"reflect"
	"testing"

	"github.com/zorzr/argmap"
)

const ProjectName = "argmap"
const ERRORUsage = "Error: incorrect arguments number for flag"
const ERRORUnrecognized = "Error: unrecognized argument"
const ERRORTooManyNames = "Error: too many value names specified"
const ERRORMissingPositional = "Error: missing required positional argument"

/**********************************************************************/
/*** CORRECT STRINGFLAG PARSING ***************************************/
/**********************************************************************/
func TestCorrectStringFlagFull_Short(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi", NArgs: 1, Vars: []string{"name"}, Help: "greets you"})

	// Everything fine using short
	os.Args = []string{ProjectName, "-hi", "jack"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"hello": []string{"jack"}}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestCorrectStringFlagFull_Long(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi", NArgs: 1, Vars: []string{"name"}, Help: "greets you"})

	// Everything fine using full name
	os.Args = []string{ProjectName, "--hello", "jack"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"hello": []string{"jack"}}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestCorrectStringFlagFull_NoValue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi", NArgs: 1, Vars: []string{"name"}, Help: "greets you"})

	// No value
	os.Args = []string{ProjectName, "--hello"}
	aMap, err := parser.Parse()
	if err != nil {
		if err.Error() != ERRORUsage+" '--hello'" {
			t.Error(err)
		}
	} else {
		t.Errorf("Wrong returned map: expected nil, got %s", aMap)
	}
}

func TestCorrectStringFlagFull_ExtraValue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi", NArgs: 1, Vars: []string{"name"}, Help: "greets you"})

	// One unrecognized extra value
	os.Args = []string{ProjectName, "--hello", "jack", "jill"}
	aMap, err := parser.Parse()
	if err != nil {
		if err.Error() != ERRORUnrecognized+" 'jill'" {
			t.Error(err)
		}
	} else {
		t.Errorf("Wrong returned map: expected nil, got %s", aMap)
	}
}

/**********************************************************************/
/*** STRINGFLAG INSERTION WITH LESS PARAMETERS ************************/
/**********************************************************************/
func TestCorrectStringFlagPartial_JustName(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewStringFlag(argmap.StringFlag{Name: "hello"})
	if err != nil {
		t.Error(err)
	}
}

func TestCorrectStringFlagPartial_JustShort(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewStringFlag(argmap.StringFlag{Short: "hi"})
	if err != nil {
		t.Error(err)
	}
}

func TestCorrectStringFlagPartial_Vars(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewStringFlag(argmap.StringFlag{Short: "hi", Vars: []string{"name"}})
	if err != nil {
		t.Error(err)
	}
}

func TestCorrectStringFlagPartial_NArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewStringFlag(argmap.StringFlag{Short: "hi", NArgs: 2})
	if err != nil {
		t.Error(err)
	}
}

func TestWrongStringFlag_UnspecifiedNArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewStringFlag(argmap.StringFlag{Short: "hi", Vars: []string{"name1", "name2"}})
	if err == nil || err.Error()[:len(ERRORTooManyNames)] != ERRORTooManyNames {
		t.Errorf("Expecting error, got nil or wrong one")
	}
}

/**********************************************************************/
/*** LISTFLAG INSERTION AND PARSING ***********************************/
/**********************************************************************/
func TestCorrectListFlagFull(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Help: "greets you"})
	parser.NewBoolFlag(argmap.BoolFlag{Name: "test", Short: "t", Help: "just trying"})
	parser.NewListFlag(argmap.ListFlag{Name: "list", Short: "l", Var: "item", Help: "give me stuff"})

	expMap := map[string]interface{}{"list": []string{"a", "b", "c"}}
	os.Args = []string{ProjectName, "--list", "a", "b", "c"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"list": []string{"a", "b"}, "hello": []string{"Novak"}}
	os.Args = []string{ProjectName, "-l", "a", "b", "--hello", "Novak"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"hello": []string{"Roger"}, "list": []string{"a", "b"}, "test": true}
	os.Args = []string{ProjectName, "--hello", "Roger", "-l", "a", "b", "-t"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"list": []string{"a"}, "test": true}
	os.Args = []string{ProjectName, "-t", "-l", "--list", "a"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestCorrectListFlagPartial(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewBoolFlag(argmap.BoolFlag{Name: "test", Short: "t", Help: "just trying"})
	parser.NewListFlag(argmap.ListFlag{Short: "l"})

	expMap := map[string]interface{}{"l": []string{"a"}}
	os.Args = []string{ProjectName, "-l", "a"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"l": []string{"a", "b"}, "test": true}
	os.Args = []string{ProjectName, "-l", "a", "b", "-t"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"l": []string{}, "test": true}
	os.Args = []string{ProjectName, "-l", "a", "b", "-t", "-l"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestWrongListFlag(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewBoolFlag(argmap.BoolFlag{Name: "test", Short: "t", Help: "just trying"})

	err := parser.NewListFlag(argmap.ListFlag{Short: "t"})
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}

	err = parser.NewListFlag(argmap.ListFlag{Short: "test"})
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}

	err = parser.NewListFlag(argmap.ListFlag{Name: "test"})
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

/**********************************************************************/
/*** BOOLFLAG INSERTION AND PARSING ***********************************/
/**********************************************************************/
func TestCorrectBoolFlag_JustName(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewBoolFlag(argmap.BoolFlag{Name: "hello"})
	if err != nil {
		t.Error(err)
	}

	os.Args = []string{ProjectName, "--hello"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"hello": true}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestCorrectBoolFlag_JustShort(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewBoolFlag(argmap.BoolFlag{Short: "hi"})
	if err != nil {
		t.Error(err)
	}

	os.Args = []string{ProjectName, "-hi"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"hi": true}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestCorrectBoolFlag_Full(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewBoolFlag(argmap.BoolFlag{Name: "hello", Short: "hi", Help: "greets you"})
	if err != nil {
		t.Error(err)
	}

	os.Args = []string{ProjectName, "--hello"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"hello": true}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

/**********************************************************************/
/*** POSITIONAL ARGUMENTS *********************************************/
/**********************************************************************/
func TestCorrectPositional_Required(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewPositionalArg(argmap.PositionalArg{Name: "your_name", Required: true})
	if err != nil {
		t.Error(err)
		return
	}

	os.Args = []string{ProjectName, "mario"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"your_name": "mario"}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestWrongPositional_Required(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewPositionalArg(argmap.PositionalArg{Name: "your_name", Required: true})
	if err != nil {
		t.Error(err)
		return
	}

	os.Args = []string{ProjectName}
	_, err = parser.Parse()
	if err == nil || err.Error()[:len(ERRORMissingPositional)] != ERRORMissingPositional {
		t.Errorf("Expecting error, got nil or wrong one")
	}
}

func TestCorrectPositional_Optional(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewPositionalArg(argmap.PositionalArg{Name: "your_name", Required: true})
	err := parser.NewPositionalArg(argmap.PositionalArg{Name: "greet_lang"})
	if err != nil {
		t.Error(err)
		return
	}

	os.Args = []string{ProjectName, "mario"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"your_name": "mario"}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	os.Args = []string{ProjectName, "mario", "spanish"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"your_name": "mario", "greet_lang": "spanish"}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestCorrectPositional_TwoRequiredOneOptional(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewPositionalArg(argmap.PositionalArg{Name: "greet_lang", Required: true})
	parser.NewPositionalArg(argmap.PositionalArg{Name: "your_surname"})
	parser.NewPositionalArg(argmap.PositionalArg{Name: "your_name", Required: true})

	os.Args = []string{ProjectName, "en", "mario"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"greet_lang": "en", "your_name": "mario"}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	os.Args = []string{ProjectName, "en", "mario", "kart"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap := map[string]interface{}{"greet_lang": "en", "your_name": "mario", "your_surname": "kart"}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

/**********************************************************************/
/*** COMMANDS AND SUBCOMMANDS *****************************************/
/**********************************************************************/
func TestCommandStringFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	cmd, _ := parser.NewCommand(argmap.CommandParams{Name: "run"})
	cmd.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi", NArgs: 1, Vars: []string{"name"}, Help: "greets you"})
	expMap := map[string]interface{}{"run": nil}

	os.Args = []string{ProjectName, "run", "-hi", "Luke"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap["run"] = map[string]interface{}{"hello": []string{"Luke"}}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	os.Args = []string{ProjectName, "-hi", "Luke"}
	aMap, err = parser.Parse()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}

	os.Args = []string{ProjectName, "run", "-hi"}
	aMap, err = parser.Parse()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}

	os.Args = []string{ProjectName, "run", "Luke"}
	aMap, err = parser.Parse()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestCommandMultipleFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Name: "hello", Help: "greets you"})
	parser.NewBoolFlag(argmap.BoolFlag{Name: "english"})

	cmd, _ := parser.NewCommand(argmap.CommandParams{Name: "add"})
	cmd.NewPositionalArg(argmap.PositionalArg{Name: "a", Required: true})
	cmd.NewPositionalArg(argmap.PositionalArg{Name: "b"})
	cmd.NewStringFlag(argmap.StringFlag{Name: "hello"})
	cmd.NewBoolFlag(argmap.BoolFlag{Short: "v"})

	cmd, _ = parser.NewCommand(argmap.CommandParams{Name: "run"})
	cmd.NewStringFlag(argmap.StringFlag{Name: "hello"})

	expMap := map[string]interface{}{"hello": []string{"Roger"}, "run": nil}
	os.Args = []string{ProjectName, "--hello", "Roger", "run"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap["run"] = map[string]interface{}{}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"hello": []string{"Roger"}, "add": nil}
	os.Args = []string{ProjectName, "--hello", "Roger", "add", "1", "-v", "2"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if expMap["add"] = map[string]interface{}{"a": "1", "v": true, "b": "2"}; !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}

	expMap = map[string]interface{}{"add": map[string]interface{}{"a": "1", "b": "2", "hello": []string{"Roger"}, "v": true}}
	os.Args = []string{ProjectName, "add", "1", "2", "--hello", "Roger", "-v"}
	aMap, err = parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

func TestSubcommandArguments(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	parser := argmap.NewArgsParser(ProjectName, t.Name())

	cmd, _ := parser.NewCommand(argmap.CommandParams{Name: "run"})
	cmd.NewStringFlag(argmap.StringFlag{Name: "out", Short: "o"})
	cmd.NewBoolFlag(argmap.BoolFlag{Short: "hi"})

	sub, _ := cmd.NewSubcommand(argmap.CommandParams{Name: "fast"})
	sub.NewStringFlag(argmap.StringFlag{Name: "hello", Short: "hi"})
	sub.NewStringFlag(argmap.StringFlag{Name: "out", Short: "o"})

	expMap := map[string]interface{}{"run": map[string]interface{}{"hi": true, "fast": map[string]interface{}{"hello": []string{"Roger"}, "out": []string{"file.txt"}}}}
	os.Args = []string{ProjectName, "run", "-hi", "fast", "-hi", "Roger", "-o", "file.txt"}
	aMap, err := parser.Parse()
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(aMap, expMap) {
		t.Errorf("Wrong returned map: expected %s, got %s", expMap, aMap)
	}
}

/**********************************************************************/
/*** GENERIC INSERTION ERRORS *****************************************/
/**********************************************************************/
func TestWrongArgument_ExistingIdentifier(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Short: "hi"})
	err := parser.NewStringFlag(argmap.StringFlag{Name: "hi"})
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestWrongArgument_HelpIdentifier(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	err := parser.NewBoolFlag(argmap.BoolFlag{Name: "help"})
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestWrongArgument_ExistingRepresentation(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.NewStringFlag(argmap.StringFlag{Short: "n"})
	err := parser.NewStringFlag(argmap.StringFlag{Name: "name", Short: "n"})
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

/**********************************************************************/
/*** GENERIC FUNCTIONS TESTS ******************************************/
/**********************************************************************/
func TestCustomHelp(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.SetHelpGenerator(func(p *argmap.ArgsParser, cmdTr []*argmap.Command) string { return p.Name + " custom help" })

	if parser.GenerateHelp() != ProjectName+" custom help" {
		t.Errorf("Wrong help message: got %s", parser.GenerateHelp())
	}
}

func TestCustomHelpFlagText(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.SetHelpFlagMessage("hello curious user!")

	aList := parser.GetArgsList()
	if text := aList[0].GetHelpStrings()[1]; text != "hello curious user!" {
		t.Errorf("Wrong HelpFlag text: got %s", text)
	}
}

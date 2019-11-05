package test

import (
	"os"
	"reflect"
	"testing"

	"github.com/zorzr/argmap"
)

const ProjectName = "argmap"
const ERRORUsage = "Error: incorrect arguments usage"
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
		if err.Error() != ERRORUsage {
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
		// if err.Error()[:len(ERRORUnrecognized)] != ERRORUnrecognized {
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

/**********************************************************************/
/*** GENERIC FUNCTIONS TESTS ******************************************/
/**********************************************************************/
func TestCustomHelp(t *testing.T) {
	parser := argmap.NewArgsParser(ProjectName, t.Name())
	parser.SetHelpGenerator(func(p *argmap.ArgsParser) string { return p.Name + " custom help" })

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

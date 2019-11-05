# argmap
**argmap** is a simple command line argument parser for [Go](https://golang.org) providing high level features and wide freedom of usage. The user inputs are stored in a map allowing to easily gather, analyze and retrieve them.

Along with positional arguments, different types of flags can be defined according to your needs. Nothing is arbitrarily made without giving you a choice, even in case something goes wrong: there are built-in functions to easily report possible errors, but they won't be automatically called if you do not want so!



## Features

- Positional arguments
  - ```Usage:    argmap required [optional]```
  - Can be both required or optional
  - The inserted values are stored in the map as strings
    - ```E.g.:    map["arg_name": "inserted_value"]```

- `StringFlag`  arguments
  - ```Usage:    argmap [-f|--flag] [value1] [value2]```
  - Customizable number of input values (separated by a space `' '`)
  - The two values are inserted in a map as a slice of strings
    - Easy to retrieve and manage (e.g., integer conversion)
    - ```E.g.:    map["flag": ["v1", "v2"]]```

- `BoolFlag`  arguments
  - ```Usage:    argmap [-f|--flag]```
  - If the flag is present, `true` is stored in the map
    - ```E.g.:    map["flag": true]```

- Customizable help message
- Flexible parameters management and error handling
  - Errors aren't automatically reported, but a simple function can be called if you want to
- Automatic help flag management
  - ```Usage:    argmap [-h|--help]```
  - A help message (default or custom) is generated and printed, then the execution ends



## Usage

In order to adopt *argmap*, make sure to get it using the `go get` command as follows;

```
go get -v "github.com/zorzr/argmap"
```

Afterwards, in your .go files you need to import the package. To do so, just add *argmap* to the list of imports at the beginning of the files:

```go
import "github.com/zorzr/argmap"
```



### Creating a parser

```go
parser := argmap.NewArgsParser("Program name", "Program description")
```

The `ArgsParser` is the most important struct of the package, which cares about:

- managing all the possible flags and positionals
- parsing the command line arguments and storing them in a map
- produce and show the help messages
- handle errors when asked to

When declared, as it can be observed above, you have to tell how your program is called and a brief description of what it does: these strings will be printed in the help message when invoked. You can then insert the arguments you need according to their type.



### Inserting a StringFlag

```go
parser.NewStringFlag(argmap.StringFlag{
	Name:  "name",
	Short: "n",
	NArgs: 1,
	Vars:  []string{"your_name"},
	Help:  "greets you",
})
```

In the code reported above, you can see how a `StringFlag` is defined. There are 5 fields which can be filled:

- *Name*: the long name of the argument, will be called by adding two minus signs before it (e.g., `--name` )
- *Short*: the short name of the argument, called with only one minus sign (e.g., `-n`)
  - **Note**. At least one of these two is needed to add the argument. If absent, an error is returned.
  - **Note**. If one of the two representations already exists in the parser (e.g, `--help`), an error is returned.
- *NArgs*: number of fields required after the flag call, default is 1 (e.g. `--name Jack` or `-n Jill`)
- *Vars*: optional name to be used in the help message to refer to the argument values (e.g. `your_name`)
- *Help*: help message to be displayed regarding this flag

A `StringFlag` can be created just by typing a Name or a Short name for the argument: this will be used to identify the input values in the map (see below for deeper details). For instance:

- ```
  StringFlag{Name: "name"}		map["name": ["Jack"]]
  ```

- ```
  StringFlag{Short: "n"}			map["n": ["Jill"]]
  ```



### Inserting a BoolFlag

```go
parser.NewBoolFlag(argmap.BoolFlag{Name: "bool", Short: "b", Help: "stores true if present"})
```

A `BoolFlag` is much simpler than a `StringFlag`. It has just these three simple fields:

- *Name*: the long name of the argument, will be called by adding two minus signs before it (e.g., `--bool` )
- *Short*: the short name of the argument, called with only one minus sign (e.g., `-b`)
  - **Note**. At least one of these two is needed to add the argument. If absent, an error is returned.
  - **Note**. If one of the two representations already exists in the parser (e.g, `--help`), an error is returned.
- *Help*: help message to be displayed regarding this flag

As a `StringFlag` , not all of them are necessary if you don't want to. The important is to properly choose a valid identifier (see below for further explanations on the matter).



### Inserting a PositionalArg

```go
parser.NewPositionalArg(argmap.PositionalArg{Name: "req", Required: true, Help: "required positional"})
```

A `PositionalArg` is a type of argument which is not indicated by a flag, but just with its presence in the arguments list.

- *Name*: the long name of the positional, which will be used as identifier in the map.
- *Required*: boolean, `true` if an error has to be raised if it isn't found in the user inputs (default is `false`).
- *Help*: help message to be displayed regarding this flag

In the package implementations, a `PositionalArg` can be located everywhere in the parsed command line string. These two possible usages are exactly the same (assuming that the `--flag` StringFlag has `NArgs = 1`):

```
./app.exe my_positional --flag flag_value
./app.exe --flag flag_value my_positional
```

**Note**. In order to avoid inconsistencies, required positionals must be placed *BEFORE* any other optional positional. The parser automatically sorts the list of inserted arguments in order to keep it organized and functioning in the correct way. Please check that your expected usage is correct by printing the program help message:

```
parser.PrintHelp()
```



### Possible mistakes with argument names

Two main concepts to be kept into considerations:

- **Identifier**: what is used as the key in the map. Must be unique for every argument.
  - If a Name is present, it is used as identifier
  - If not, then the Short is adopted for the same purpose
- **Representations**: what can be used to refer to a flag. Must be uniquely associated to every argument.
  - Example: `--name`, `-n` are the two representations for the StringFlag above.
  - There can't be another flag named `"name"` or with a `"n"` short name.

Brief summary:

```
StringFlag{Name: "name"}                 BoolFlag{Name: "n"}      OK
StringFlag{Name: "name", Short: "n"}     BoolFlag{Name: "n"}      OK
StringFlag{Name: "name", Short: "n"}     BoolFlag{Short: "n"}     WRONG (-n representation ambiguity)
StringFlag{Name: "n"}                    BoolFlag{Short: "n"}     WRONG ("n" identifier ambiguity)
```

You can look for possible errors by observing the return value of any functions which adds a new argument, e.g.

```go
err := parser.NewBoolFlag(argmap.BoolFlag{Name: "bool", Short: "b")
```



### I'm confused. How can I put all those pieces together?

Worry not, it's easier to use than to learn.

Look at the [hello example](https://github.com/zorzr/argmap/blob/master/examples/hello/main.go). It shows a simple usage of a StringFlag and highlights a useful function for retrieving values from the map. Here's two possible user inputs and the respective outputs:

```
$ hello -hi Jack
Hello Jack
$ hello --hello Jill
Hello Jill
```



Looking for positionals? Check out the [myname_pos example](https://github.com/zorzr/argmap/blob/master/examples/myname_pos/main.go) instead. Some randomly chosen inputs:

```
$ myname_pos James
My name is James
$ myname_pos James Bond
My name is Bond, James Bond
```



In the `examples` you can find other common usages and several tricks to make a better use of *argmap*.



## Upcoming functionalities

- Variable number of arguments for StringFlags
- Commands definition and argument parsing
- Required StringFlags (for design simplicity)
- Parallelism through goroutines



## Support

Please report if you notice any bug by opening a new issue here on Github!

Thank you!

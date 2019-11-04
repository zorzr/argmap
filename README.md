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

#### StringFlag

```go
package main

import (
	"fmt"
	"github.com/zorzr/argmap"
)

func main() {
	parser := argmap.NewArgsParser("Program name", "Program description")
	parser.NewStringFlag(argmap.StringFlag{
		Name:  "hello",
		Short: "hi",
		NArgs: 1,
		Vars:  []string{"urname"},
		Help:  "greets you",
	})

	aMap, err := parser.Parse()
	if err != nil {
		parser.ReportError(err)
	}

	if helloList, ok := aMap["hello"]; ok {
		name := helloList.([]string)[0]
		fmt.Println("Hello " + name)
	}
}

```

In the code reported above, you can see how a StringFlag is defined and how it can be parsed. Example of user input:

```
$ argmap -hi Jack
Hello Jack
$ argmap --hello Jill
Hello Jill
```



## Upcoming functionalities

- Variable number of arguments for StringFlags
- Commands definition and argument parsing
- Required StringFlags (for design simplicity)
- Parallelism through goroutines



## Support

Please report if you notice any bug by opening a new issue here on Github!

Thank you!

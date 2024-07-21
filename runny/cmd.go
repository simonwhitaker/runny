package runny

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Run() {
	runny, err := readConfig()
	if err != nil {
		color.Red("Problem reading config: %v", err)
	}

	args := os.Args[1:]
	for len(args) > 0 && args[0][0] == '-' {
		option := args[0]
		switch option {
		case "-h", "--help":
			runny.ShowHelp()
			return
		case "-v", "--verbose":
			// TODO: handle verbose case
		default:
			panic(fmt.Sprintf("Unknown option: %s", option))
		}
		args = args[1:]
	}

	// read command line args
	if len(args) > 0 {
		name := CommandName(args[0])
		runny.Execute(name, args[1:]...)
	} else {
		runny.ShowHelp()
	}
}

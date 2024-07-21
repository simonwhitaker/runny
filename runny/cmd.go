package runny

import (
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
			runny.PrintHelp()
			return
		case "-v", "--verbose":
			runny.verbose = true
		default:
			color.Red("Unknown option: %s", option)
			os.Exit(1)
		}
		args = args[1:]
	}

	// read command line args
	if len(args) > 0 {
		name := CommandName(args[0])
		runny.Execute(name, args[1:]...)
	} else {
		runny.PrintCommands()
	}
}

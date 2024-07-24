package runny

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Run() {
	exitWithError := func(err error) {
		errStr := fmt.Sprintf("%v\n", err)
		red := color.New(color.FgRed)
		red.Fprint(os.Stderr, errStr)
		os.Exit(1)
	}

	runny, err := readConfig(".runny.yaml")
	if err != nil {
		exitWithError(err)
	}

	// Parse command-line options
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
			exitWithError(fmt.Errorf("unknown option: %s", option))
		}
		args = args[1:]
	}

	// Process runny command
	if len(args) > 0 {
		name := CommandName(args[0])
		err := runny.Execute(name, args[1:]...)
		if err != nil {
			exitWithError(err)
		}
	} else {
		runny.PrintCommands()
	}
}

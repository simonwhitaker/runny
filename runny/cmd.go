package runny

import (
	"os"
)

func Run() {
	runny, err := readConfig(".runny.yaml")
	if err != nil {
		errorColor.Printf("Problem reading config: %v\n", err)
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
			errorColor.Printf("Unknown option: %s\n", option)
			os.Exit(1)
		}
		args = args[1:]
	}

	// read command line args
	if len(args) > 0 {
		name := CommandName(args[0])
		err := runny.Execute(name, args[1:]...)
		if err != nil {
			errorColor.Println(err)
			os.Exit(1)
		}
	} else {
		runny.PrintCommands()
	}
}

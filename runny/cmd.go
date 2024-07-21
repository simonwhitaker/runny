package runny

import (
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func showHelp(conf Config) {
	commands := conf.Commands
	names := make([]CommandName, len(commands))
	i := 0
	for key := range commands {
		names[i] = key
		i += 1
	}

	slices.Sort(names)
	var separator = " "
	var maxCommandLength = 60
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		separator = "\t"
		maxCommandLength = 40
	}

	for _, name := range names {
		var rawCommand = commandStringToSingleLine(commands[name].Run, maxCommandLength)

		fmt.Print(primaryColor.Sprint(name))
		fmt.Print(separator)
		fmt.Println(secondaryColor.Sprint(rawCommand))
	}

}

func Run() {
	conf, err := readConfig()
	if err != nil {
		color.Red("Problem reading config: %v", err)
	}

	args := os.Args[1:]
	for len(args) > 0 && args[0][0] == '-' {
		option := args[0]
		switch option {
		case "-h", "--help":
			showHelp(conf)
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
		if command, ok := conf.Commands[name]; ok {
			command.Execute(conf, args[1:]...)
		} else {
			color.Red("Command not found")
		}
	} else {
		showHelp(conf)
	}
}

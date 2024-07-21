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

	// read command line args
	if len(os.Args) > 1 {
		name := CommandName(os.Args[1])
		if command, ok := conf.Commands[name]; ok {
			command.Execute(conf, os.Args[2:]...)
		} else {
			color.Red("Command not found")
		}
	} else {
		showHelp(conf)
	}
}

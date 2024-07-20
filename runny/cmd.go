package runny

import (
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
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

	for _, name := range names {
		nameColor := color.New(color.Bold)
		fmt.Printf("%s %s\n", nameColor.Sprint(name), secondaryColor.Sprint(commandStringToSingleLine(commands[name].Command)))
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
			command.Execute(conf)
		} else {
			color.Red("Command not found")
		}
	} else {
		showHelp(conf)
	}
}

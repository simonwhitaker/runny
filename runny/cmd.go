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

	for _, name := range names {
		nameColor := color.New(color.Bold)
		var separator = " "
		if !term.IsTerminal(int(os.Stdout.Fd())) {
			separator = "\t"
		}
		fmt.Printf("%s%s%s\n",
			nameColor.Sprint(name),
			separator,
			secondaryColor.Sprint(commandStringToSingleLine(commands[name].Command, 40)))
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

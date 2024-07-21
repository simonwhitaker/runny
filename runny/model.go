package runny

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

type CommandName string

type CommandDef struct {
	Run   string
	Needs []CommandName
	If    string
}

type Config struct {
	Commands map[CommandName]CommandDef
	Shell    string
}

func (c *Config) GetShell() Shell {
	if len(c.Shell) > 0 {
		return NewShell(c.Shell)
	}
	return NewShell(defaultShell)
}

func (c *Config) ShowHelp() {
	commands := c.Commands
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

func (c *Config) Execute(name CommandName, args ...string) error {
	command, ok := c.Commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s", name)
	}

	// Handle pre-commands
	for _, name := range command.Needs {
		// TODO: handle invalid names
		err := c.Execute(name)
		if err != nil {
			return err
		}
	}

	// Check the If
	cond := strings.TrimSpace(command.If)
	if len(cond) > 0 {
		shell := c.GetShell()
		err := shell.Run(cond)
		if err != nil {
			// Run returns an error if the exit status is not zero. So in this case, this means the test failed.
			// secondaryColor.Printf("'%v' not true\n", cond)
			return nil
		}
	}

	// Handle the Run
	run := strings.TrimSpace(command.Run)
	if len(run) > 0 {
		shell := c.GetShell()
		err := shell.Run(run, args...)
		if err != nil {
			fmt.Printf("%s %s\n", color.RedString(string(run)), secondaryColor.Sprint(err))
			return err
		}
	}

	return nil
}

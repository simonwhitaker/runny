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
	verbose  bool
}

func (c *Config) GetShell() Shell {
	shellString := c.Shell
	if len(shellString) == 0 {
		shellString = defaultShell
	}
	shell, err := NewShell(shellString)
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
	return shell
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
	shell := c.GetShell()

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
		err := shell.Run(cond, c.verbose)
		if err != nil {
			// Run returns an error if the exit status is not zero. So in this case, this means the test failed.
			if c.verbose {
				secondaryColor.Printf("%v: '%v' not true, skipping\n", name, cond)
			}
			return nil
		}
	}

	// Handle the Run
	run := strings.TrimSpace(command.Run)
	if len(run) > 0 {
		err := shell.Run(run, c.verbose, args...)
		if err != nil {
			fmt.Printf("%s %s\n", color.RedString(string(run)), secondaryColor.Sprint(err))
			return err
		}
	}

	return nil
}

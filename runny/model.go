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
	Run      string        `json:"run,omitempty"`
	Needs    []CommandName `json:"needs,omitempty"`
	If       string        `json:"if,omitempty"`
	Env      []string      `json:"env,omitempty"`
	ArgNames []string      `json:"argnames,omitempty"`
}

type Config struct {
	Commands map[CommandName]CommandDef `json:"commands"`
	Shell    string                     `json:"shell,omitempty"`
	Env      []string                   `json:"env,omitempty"`
	verbose  bool
}

func (c *Config) GetShell() Shell {
	shellString := c.Shell
	if len(shellString) == 0 {
		shellString = defaultShell
	}
	shell, err := NewShell(shellString)
	if err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	return shell
}

func (c *Config) PrintHelp() {
	titleString := color.New(color.FgYellow, color.Bold).Sprintf("runny")
	usageString := color.New(color.Bold).Sprintf("runny [options] [command]")
	fmt.Printf(`%s -- for running things.

Usage:
  %s

Options:
  -h, --help     Show this help
  -v, --verbose  Enable verbose mode

Run without arguments to list commands.`, titleString, usageString)
	fmt.Print("\n")
}

func (c *Config) PrintCommands() {
	commands := c.Commands
	names := make([]CommandName, len(commands))
	i := 0
	for key := range commands {
		names[i] = key
		i += 1
	}

	slices.Sort(names)
	var separator = " "
	maxLineLength := 80
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		separator = "\t"
		maxLineLength = 40
	}

	for _, name := range names {
		maxCommandLength := maxLineLength - len(name) - 1
		var rawCommand = commandStringToSingleLine(commands[name].Run, maxCommandLength)

		fmt.Print(primaryColor.Sprint(name))
		fmt.Print(separator)
		fmt.Println(secondaryColor.Sprint(rawCommand))
	}
}

func (c *Config) Execute(name CommandName, args ...string) error {
	shell := c.GetShell()
	command, ok := c.Commands[name]
	if !ok {
		errorMsg := fmt.Sprintf("unknown command: %s", name)
		errorColor.Println(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	env := append(c.Env, command.Env...)

	// Check the If condition
	cond := strings.TrimSpace(command.If)
	if len(cond) > 0 {
		err := shell.Run(cond, []string{}, false, c.verbose, env)
		if err != nil {
			// Run returns an error if the exit status is not zero. So in this case, this means the test failed.
			if c.verbose {
				secondaryColor.Printf("%v: '%v' not true, skipping\n", name, cond)
			}
			return nil
		}
	}

	// Handle Needs
	for _, name := range command.Needs {
		err := c.Execute(name)
		if err != nil {
			errorColor.Print(err)
			return err
		}
	}

	// Collect args
	if len(command.ArgNames) > 0 {
		if len(args) < len(command.ArgNames) {
			return fmt.Errorf("%d named args defined but only %d supplied", len(command.ArgNames), len(args))
		}
		for _, argName := range command.ArgNames {
			env = append(env, fmt.Sprintf("%s=%s", argName, args[0]))
			args = args[1:]
		}
	}

	// Handle the Run
	run := strings.TrimSpace(command.Run)
	if len(run) > 0 {
		err := shell.Run(run, args, true, c.verbose, env)
		if err != nil {
			return err
		}
	}

	return nil
}

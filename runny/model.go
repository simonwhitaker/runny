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
	// The command to be run
	Run string `json:"run,omitempty"`
	// A list of the commands that must be run before this one
	Needs []CommandName `json:"needs,omitempty"`
	// A list of the commands that must be run after this one
	Then []CommandName `json:"then,omitempty"`
	// A conditional expression; this command only runs if the If expression evaluates to true
	If string `json:"if,omitempty"`
	// A list of environment variables to be set when running the command
	Env []string `json:"env,omitempty"`
	// A list of argument names to be passed on the command line when invoking Runny. These arguments can be accessed in
	// the Run string by prefixing them with $.
	ArgNames []string `json:"argnames,omitempty"`
	// If true, the command is suppressed from the output generated when you run Runny with no arguments, and from shell
	// completion.
	Internal bool `json:"internal,omitempty"`
}

type Config struct {
	Commands map[CommandName]CommandDef `json:"commands"`
	// The shell to be used when running commands. Defaults to /bin/bash
	Shell string `json:"shell,omitempty"`
	// A list of environment variables to be set when running all commands
	Env []string `json:"env,omitempty"`
	// Set if the -v/--verbose flag is set when invoking Runny
	verbose bool
}

func (c *Config) GetShell() (Shell, error) {
	shellString := c.Shell
	if len(shellString) == 0 {
		shellString = defaultShell
	}
	shell, err := NewShell(shellString)
	if err != nil {
		return nil, err
	}
	return shell, nil
}

func (c *Config) PrintHelp() {
	titleString := color.New(color.FgYellow, color.Bold).Sprintf("🍯 runny")
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
	names := []CommandName{}
	for key, cmd := range commands {
		if c.verbose || !cmd.Internal {
			names = append(names, key)
		}
	}

	slices.Sort(names)
	var separator = " "
	var maxLineLength int

	if term.IsTerminal(int(os.Stdout.Fd())) {
		if c.verbose {
			maxLineLength = 0
		} else {
			width, _, err := term.GetSize(int(os.Stdin.Fd()))
			if err == nil {
				maxLineLength = width
			}
		}
	} else {
		separator = "\t"
		maxLineLength = 0
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
	shell, err := c.GetShell()
	if err != nil {
		return err
	}

	command, ok := c.Commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s", name)
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

	// Handle Then
	for _, name := range command.Then {
		err := c.Execute(name)
		if err != nil {
			return err
		}
	}

	return nil
}

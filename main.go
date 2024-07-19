package main

import (
	"cmp"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var secondaryColor = color.New(color.FgHiBlack)
var defaultShell = "/bin/bash"

type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Config struct {
	Commands []Command `yaml:"commands"`
	Shell    string    `yaml:"shell,omitempty"`
}

func executeCommand(shell, rawCommand string) error {
	command := strings.TrimSpace(rawCommand)
	// FIXME: -c is bash-specific, won't work with every shell
	args := []string{"-c", command}

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func main() {
	var conf Config

	// Read .runny.yaml from the current directory
	yamlFile, err := os.ReadFile(".runny.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	shell := conf.Shell
	if len(shell) == 0 {
		shell = defaultShell
	}

	// read command line args
	if len(os.Args) > 1 {
		cmdName := os.Args[1]
		found := false
		for _, c := range conf.Commands {
			if c.Name == cmdName {
				err := executeCommand(shell, c.Command)
				if err != nil {
					fmt.Printf("%s %s\n", color.RedString(c.Command), secondaryColor.Sprint(err))
				}
				found = true
			}
		}
		if !found {
			color.Red("Command not found")
		}
	} else {
		commands := conf.Commands
		slices.SortFunc(commands, func(a, b Command) int {
			return cmp.Compare(a.Name, b.Name)
		})
		for _, cmd := range commands {
			nameColor := color.New(color.Bold)

			fmt.Printf("%s %s\n", nameColor.Sprint(cmd.Name), secondaryColor.Sprint(cmd.Command))
		}
	}
}

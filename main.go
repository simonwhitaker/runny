package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type CommandName string

var secondaryColor = color.New(color.FgHiBlack)
var defaultShell = "/bin/bash"

func singleLine(command string) string {
	command = strings.TrimSpace(command)
	lines := strings.Split(command, "\n")
	trimmedLines := []string{}
	for _, line := range lines {
		trimmedLines = append(trimmedLines, strings.TrimSpace(line))
	}
	return strings.Join(trimmedLines, "; ")
}

type CommandDef struct {
	Command string        `yaml:"command"`
	Pre     []CommandName `yaml:"pre"`
	Post    []CommandName `yaml:"post"`
}

type Config struct {
	Commands map[CommandName]CommandDef `yaml:"commands"`
	shell    string                     `yaml:"shell"`
}

func (c *Config) GetShell() string {
	if len(c.shell) > 0 {
		return c.shell
	}
	return defaultShell
}

func (c *CommandDef) Execute(conf Config) error {
	// Handle pre-commands
	for _, name := range c.Pre {
		// TODO: handle invalid names
		command := conf.Commands[name]
		err := command.Execute(conf)
		if err != nil {
			return err
		}
	}

	// Handle the command
	command := strings.TrimSpace(c.Command)
	if len(command) > 0 {
		// FIXME: -c is bash-specific, won't work with every shell
		args := []string{"-c", command}

		cmd := exec.Command(conf.GetShell(), args...)
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s %s\n", color.RedString(string(command)), secondaryColor.Sprint(err))
			return err
		}
	}

	// Handle post-commands
	for _, name := range c.Post {
		// TODO: handle invalid names
		command := conf.Commands[name]
		err := command.Execute(conf)
		if err != nil {
			return err
		}
	}

	return nil
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

	// read command line args
	if len(os.Args) > 1 {
		name := CommandName(os.Args[1])
		if command, ok := conf.Commands[name]; ok {
			command.Execute(conf)
		} else {
			color.Red("Command not found")
		}
	} else {
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
			fmt.Printf("%s %s\n", nameColor.Sprint(name), secondaryColor.Sprint(singleLine(commands[name].Command)))
		}
	}
}

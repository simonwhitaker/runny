package runny

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type CommandName string
type CommandDef struct {
	Run   string        `yaml:"run"`
	Needs []CommandName `yaml:"needs"`
}

type Config struct {
	Commands map[CommandName]CommandDef `yaml:"commands"`
	Shell    string                     `yaml:"shell"`
}

func (c *Config) GetShell() Shell {
	if len(c.Shell) > 0 {
		return NewShell(c.Shell)
	}
	return NewShell(defaultShell)
}

func (c *CommandDef) Execute(conf Config) error {
	// Handle pre-commands
	for _, name := range c.Needs {
		// TODO: handle invalid names
		command := conf.Commands[name]
		err := command.Execute(conf)
		if err != nil {
			return err
		}
	}

	// Handle the command
	command := strings.TrimSpace(c.Run)
	if len(command) > 0 {
		shell := conf.GetShell()
		err := shell.Run(command)
		if err != nil {
			fmt.Printf("%s %s\n", color.RedString(string(command)), secondaryColor.Sprint(err))
			return err
		}
	}

	return nil
}

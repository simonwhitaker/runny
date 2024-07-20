package runny

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

type CommandName string
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

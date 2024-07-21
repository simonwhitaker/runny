package runny

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

func (c *CommandDef) Execute(conf Config, args ...string) error {
	// Handle pre-commands
	for _, name := range c.Needs {
		// TODO: handle invalid names
		command := conf.Commands[name]
		err := command.Execute(conf)
		if err != nil {
			return err
		}
	}

	// Check the If
	cond := strings.TrimSpace(c.If)
	if len(cond) > 0 {
		shell := conf.GetShell()
		err := shell.Run(cond)
		if err != nil {
			// Run returns an error if the exit status is not zero. So in this case, this means the test failed.
			// secondaryColor.Printf("'%v' not true\n", cond)
			return nil
		}
	}

	// Handle the Run
	run := strings.TrimSpace(c.Run)
	if len(run) > 0 {
		shell := conf.GetShell()
		err := shell.Run(run, args...)
		if err != nil {
			fmt.Printf("%s %s\n", color.RedString(string(run)), secondaryColor.Sprint(err))
			return err
		}
	}

	return nil
}

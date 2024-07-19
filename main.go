package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var secondaryColor = color.New(color.FgHiBlack)

type Cmd struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Cmds struct {
	Cmds []Cmd `yaml:"commands"`
}

func main() {
	var cmds Cmds

	// read a YAML file from disk
	yamlFile, err := os.ReadFile(".runny.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &cmds)
	if err != nil {
		panic(err)
	}

	// read command line args
	if len(os.Args) > 1 {
		cmdName := os.Args[1]
		found := false
		for _, c := range cmds.Cmds {
			if c.Name == cmdName {
				cmdTokens := strings.Split(c.Command, " ")
				args := cmdTokens[1:]
				args = append(args, os.Args[2:]...)

				cmd := exec.Command(cmdTokens[0], args...)
				cmd.Stdout = os.Stdout
				err := cmd.Run()
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
		println("Found ", len(cmds.Cmds), " commands")
		// iterate over cmds
		for _, cmd := range cmds.Cmds {
			nameColor := color.New(color.Bold)

			fmt.Printf("%s %s\n", nameColor.Sprint(cmd.Name), secondaryColor.Sprint(cmd.Command))
		}
	}
}

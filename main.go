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

type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Config struct {
	Commands []Command `yaml:"commands"`
}

func main() {
	var conf Config

	// read a YAML file from disk
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
		cmdName := os.Args[1]
		found := false
		for _, c := range conf.Commands {
			if c.Name == cmdName {
				cmdTokens := strings.Split(c.Command, " ")

				cmd := exec.Command(cmdTokens[0], cmdTokens[:1]...)
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

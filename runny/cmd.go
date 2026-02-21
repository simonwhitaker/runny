package runny

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

func printHelp(flags *pflag.FlagSet) {
	titleString := color.New(color.FgYellow, color.Bold).Sprintf("🍯 runny")
	usageString := color.New(color.Bold).Sprintf("runny [options] [command]")
	fmt.Printf(`%s -- for running things.

Usage:
  %s

Options:
%sRun without arguments to list commands.
`, titleString, usageString, flags.FlagUsages())
}

const sampleConfig = `commands:
  hello:
    run: echo "Hello from Runny!"
  greet:
    run: echo "Hello, $name!"
    argnames:
      - name
`

func initConfig(path string) error {
	// Check if config file already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	// Write the sample config
	err := os.WriteFile(path, []byte(sampleConfig), 0644)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", path, err)
	}

	fmt.Printf("Created %s\n", path)
	return nil
}

func run() error {
	flags := pflag.NewFlagSet("runny", pflag.ContinueOnError)

	file := flags.StringP("file", "f", ".runny.yaml", "Specify a runny file")
	verbose := flags.BoolP("verbose", "v", false, "Enable verbose mode")
	help := flags.BoolP("help", "h", false, "Show this help")
	init := flags.Bool("init", false, "Create a sample .runny.yaml file")
	schema := flags.Bool("schema", false, "Print the JSON schema for .runny.yaml")

	flags.Usage = func() { printHelp(flags) }

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if *help {
		printHelp(flags)
		return nil
	}

	if *schema {
		s, err := GenerateSchema()
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	}

	if *init {
		return initConfig(*file)
	}

	runny, err := readConfig(*file)
	if err != nil {
		return err
	}
	runny.verbose = *verbose

	// Process runny command
	args := flags.Args()
	if len(args) > 0 {
		name := CommandName(args[0])
		err := runny.Execute(name, args[1:]...)
		if err != nil {
			return err
		}
	} else {
		runny.PrintCommands()
	}
	return nil
}

func Run() {
	err := run()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		red := color.New(color.FgRed)
		red.Fprint(os.Stderr, errStr)
		os.Exit(1)
	}
}

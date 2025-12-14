package runny

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printHelp() {
	titleString := color.New(color.FgYellow, color.Bold).Sprintf("🍯 runny")
	usageString := color.New(color.Bold).Sprintf("runny [options] [command]")
	fmt.Printf(`%s -- for running things.

Usage:
  %s

Options:
  -h, --help     Show this help
  -v, --verbose  Enable verbose mode
  --init         Create a sample .runny.yaml file

Run without arguments to list commands.`, titleString, usageString)
	fmt.Print("\n")
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

func run(path string) error {
	// Parse command-line options
	args := os.Args[1:]
	verbose := false
	for len(args) > 0 && args[0][0] == '-' {
		option := args[0]
		switch option {
		case "-h", "--help":
			printHelp()
			return nil
		case "-v", "--verbose":
			verbose = true
		case "--init":
			return initConfig(path)
		default:
			return fmt.Errorf("unknown option: %s", option)
		}
		args = args[1:]
	}

	runny, err := readConfig(path)
	if err != nil {
		return err
	}
	runny.verbose = verbose

	// Process runny command
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
	err := run(".runny.yaml")
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		red := color.New(color.FgRed)
		red.Fprint(os.Stderr, errStr)
		os.Exit(1)
	}
}

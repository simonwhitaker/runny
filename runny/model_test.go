package runny

import "testing"

func TestGetShell(t *testing.T) {
	c := Config{}
	shell, err := c.GetShell()
	if shell == nil || err != nil {
		t.Errorf("Expected a valid shell, got an error: %v", err)
	}

	c = Config{Shell: "/bin/bash"}
	shell, err = c.GetShell()
	if shell == nil || err != nil {
		t.Errorf("Expected a valid shell, got an error: %v", err)
	}

	c = Config{Shell: "/bin/pwsh"}
	_, err = c.GetShell()
	if err == nil {
		t.Errorf("Expected an error, but the call succeeded")
	}
}

func ExampleConfig_PrintHelp() {
	c := Config{}
	c.PrintHelp()
	// Output: 🍯 runny -- for running things.
	//
	// Usage:
	//   runny [options] [command]
	//
	// Options:
	//   -h, --help     Show this help
	//   -v, --verbose  Enable verbose mode
	//
	// Run without arguments to list commands.
}

func ExampleConfig_PrintCommands() {
	c := Config{Commands: map[CommandName]CommandDef{"foo": {Run: "bar"}}}
	c.PrintCommands()
	// Output: foo	bar
}

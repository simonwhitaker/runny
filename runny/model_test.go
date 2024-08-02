package runny

import "testing"

func TestConfigGetShell(t *testing.T) {
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

func TestCommandDefGetShell(t *testing.T) {
	cmdWithShell := CommandDef{Shell: "/bin/zsh"}
	cmdWithoutShell := CommandDef{}
	c := Config{Shell: "/bin/bash", Commands: map[CommandName]CommandDef{
		"cmdWithShell":    cmdWithShell,
		"cmdWithoutShell": cmdWithoutShell,
	}}

	shell, err := cmdWithShell.GetShell(&c)
	if shell, ok := shell.(*PosixShell); ok && shell.command != "/bin/zsh" {
		t.Errorf("Expected a valid shell, got an error: %v", err)
	}

	shell, err = cmdWithoutShell.GetShell(&c)
	if shell, ok := shell.(*PosixShell); ok && shell.command != "/bin/bash" {
		t.Errorf("Expected a valid shell, got an error: %v", err)
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

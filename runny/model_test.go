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

func TestCommandDisplayUsesDescriptionWhenPresent(t *testing.T) {
	displayText, displayColor := commandDisplay(CommandDef{Description: "foo description", Run: "echo foo"})
	if displayText != "foo description" {
		t.Fatalf("expected description text, got %q", displayText)
	}
	if displayColor != descriptionColor {
		t.Fatalf("expected descriptionColor for description")
	}
}

func TestCommandDisplayFallsBackToRun(t *testing.T) {
	displayText, displayColor := commandDisplay(CommandDef{Description: "   ", Run: "echo foo"})
	if displayText != "echo foo" {
		t.Fatalf("expected run text, got %q", displayText)
	}
	if displayColor != runValueColor {
		t.Fatalf("expected runValueColor for run")
	}
}

func ExampleConfig_PrintCommands() {
	c := Config{Commands: map[CommandName]CommandDef{"foo": {Description: "Show foo output", Run: "bar"}}}
	c.PrintCommands()
	// Output: foo	Show foo output
}

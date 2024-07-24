package runny

import (
	"testing"
)

type Input struct {
	command   string
	maxLength int
}

type Param struct {
	input    Input
	expected string
}

func TestCommandStringToSingleLine(t *testing.T) {
	params := []Param{
		{input: Input{command: "", maxLength: 80}, expected: ""},
		{input: Input{command: "foo bar wibble", maxLength: 10}, expected: "foo bar w…"},
		{input: Input{command: "foo\nbar", maxLength: 80}, expected: "foo; bar"},
		{input: Input{command: " foo \n bar ", maxLength: 80}, expected: "foo; bar"},
	}
	for _, p := range params {
		output := commandStringToSingleLine(p.input.command, p.input.maxLength)
		if output != p.expected {
			t.Fatalf("Expected %s, got %s\n", p.expected, output)
		}
	}
}

func TestReadConfig(t *testing.T) {
	conf, err := readConfig("fixtures/simple-good.yaml")
	if err != nil {
		t.Fatalf("Got error when reading in config file: %v\n", err)
	}

	expectedLenCommands := 3
	if len(conf.Commands) != expectedLenCommands {
		t.Fatalf("Expected %d commands, got %d\n", expectedLenCommands, len(conf.Commands))
	}

	expectedCommandFooRun := "ls foo"
	if conf.Commands["foo"].Run != expectedCommandFooRun {
		t.Fatalf("Expected foo command's run value to be %s, got %s", expectedCommandFooRun, conf.Commands["foo"].Run)
	}
}

package runny

import (
	"strings"
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
		{input: Input{command: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aliquam convallis. Nunc lacus. Curabitur nunc mauris, commodo vel, eleifend in, ornare sit amet, felis. Nullam mi neque, feugiat et, porttitor vitae, pharetra non, lacus. Fusce imperdiet sem quis dui.", maxLength: 0}, expected: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aliquam convallis. Nunc lacus. Curabitur nunc mauris, commodo vel, eleifend in, ornare sit amet, felis. Nullam mi neque, feugiat et, porttitor vitae, pharetra non, lacus. Fusce imperdiet sem quis dui."},
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

func TestReadInvalidConfigs(t *testing.T) {
	_, err := readConfig("fixtures/invalid-runny.yaml")
	if err == nil {
		t.Fatalf("Expected an error when reading invalid runny config, but reading was successful")
	}

	if !strings.Contains(err.Error(), "invalid runny config file") {
		t.Fatalf("unexpected error message: %v", err)
	}

	_, err = readConfig("fixtures/invalid-yaml.yaml")
	if err == nil {
		t.Fatalf("Expected an error when reading invalid YAML file, but reading was successful")
	}
}

func TestReadMissingFIle(t *testing.T) {
	_, err := readConfig("fixtures/does-not-exist.yaml")
	if err == nil {
		t.Fatalf("Expected an error when reading missing config, but reading was successful")
	}

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("unexpected error message: %v", err)
	}

}

func TestConfigWithCircularDependency(t *testing.T) {
	_, err := readConfig("fixtures/invalid-circular-dependency.yaml")
	if err == nil {
		t.Fatalf("Expected an error when reading invalid config, but reading was succeeded")
	}

	if !strings.Contains(err.Error(), "edge would create a cycle") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

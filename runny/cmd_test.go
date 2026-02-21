package runny

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func TestPrivateRun(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{os.Args[0], "-f", "fixtures/minimal.yaml"}
	err := run()
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-v", "-f", "fixtures/minimal.yaml"}
	err = run()
	if err != nil {
		t.Fatalf("Expected success when run with -v, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-h"}
	err = run()
	if err != nil {
		t.Fatalf("Expected success when run with -h, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-Z"}
	err = run()
	if err == nil {
		t.Fatalf("Expected failure when run with unknown command-line arg, got success")
	}

	os.Args = []string{os.Args[0], "-f", "fixtures/minimal.yaml", "ok"}
	err = run()
	if err != nil {
		t.Fatalf("Expected success when run with known command, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-f", "fixtures/minimal.yaml", "not-ok"}
	err = run()
	if err == nil {
		t.Fatalf("Expected failure when run with unknown command, got success")
	}

	os.Args = []string{os.Args[0], "-f", "fixtures/invalid-yaml.yaml"}
	err = run()
	if err == nil {
		t.Fatalf("Expected error, got success")
	}
}

func Example_printHelp() {
	flags := pflag.NewFlagSet("runny", pflag.ContinueOnError)
	flags.StringP("file", "f", ".runny.yaml", "Specify a runny file")
	flags.BoolP("verbose", "v", false, "Enable verbose mode")
	flags.BoolP("help", "h", false, "Show this help")
	flags.Bool("init", false, "Create a sample .runny.yaml file")
	printHelp(flags)
	// Output: 🍯 runny -- for running things.
	//
	// Usage:
	//   runny [options] [command]
	//
	// Options:
	//   -f, --file string   Specify a runny file (default ".runny.yaml")
	//   -h, --help          Show this help
	//       --init          Create a sample .runny.yaml file
	//   -v, --verbose       Enable verbose mode
	// Run without arguments to list commands.
}

package runny

import (
	"os"
	"testing"
)

func TestPrivateRun(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{os.Args[0]}
	err := run("fixtures/minimal.yaml")
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-v"}
	err = run("fixtures/minimal.yaml")
	if err != nil {
		t.Fatalf("Expected success when run with -v, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-h"}
	err = run("fixtures/minimal.yaml")
	if err != nil {
		t.Fatalf("Expected success when run with -h, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "-Z"}
	err = run("fixtures/minimal.yaml")
	if err == nil {
		t.Fatalf("Expected failure when run with unknown command-line arg, got success")
	}

	os.Args = []string{os.Args[0], "ok"}
	err = run("fixtures/minimal.yaml")
	if err != nil {
		t.Fatalf("Expected success when run with known command, got error: %v", err)
	}

	os.Args = []string{os.Args[0], "not-ok"}
	err = run("fixtures/minimal.yaml")
	if err == nil {
		t.Fatalf("Expected failure when run with unknown command, got success")
	}

	err = run("fixtures/invalid-yaml.yaml")
	if err == nil {
		t.Fatalf("Expected error, got success")
	}
}

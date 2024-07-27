package runny

import "testing"

func TestNewShell(t *testing.T) {
	shell, err := NewShell("/bin/bash")
	if shell == nil || err != nil {
		t.Errorf("Expected a valid shell, got an error: %v", err)
	}

	_, err = NewShell("/bin/pwsh")
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestPosixShellRun(t *testing.T) {
	ps := PosixShell{command: "/bin/bash"}
	err := ps.Run("ls", []string{}, false, false, []string{})
	if err != nil {
		t.Errorf("Expected success, got an error: %v", err)
	}

	err = ps.Run("ls", []string{"/"}, false, false, []string{})
	if err != nil {
		t.Errorf("Expected success when passing a valid arg, got an error: %v", err)
	}

	err = ps.Run("ls", []string{}, true, false, []string{})
	if err != nil {
		t.Errorf("Expected success when echoing stdout, got an error: %v", err)
	}

	err = ps.Run("ls", []string{}, false, true, []string{})
	if err != nil {
		t.Errorf("Expected success when running in verbose mode, got an error: %v", err)
	}

	err = ps.Run("ls", []string{}, false, false, []string{"FOO=foo"})
	if err != nil {
		t.Errorf("Expected success when passing env variables, got an error: %v", err)
	}

	err = ps.Run("ls", []string{"does-not-exist"}, false, false, []string{})
	if err == nil {
		t.Errorf("Expected error if the command fails, but no error was returned")
	}
}

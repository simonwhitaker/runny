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

package runny

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Shell interface {
	Run(command string, extraArgs []string, echoStdout, verbose bool, env []string) error
}

func NewShell(command string) (Shell, error) {
	switch path.Base(command) {
	case "pwsh", "powershell":
		return nil, fmt.Errorf("unsupported shell: %s", command)
	default:
		return PosixShell{command: command}, nil
	}
}

type PosixShell struct {
	// PosixShell might be the wrong term here. Fish for example isn't strictly POSIX-compliant. But really, any shell
	// that allows you to run a command with `shell -c command` works.
	command string
}

func (shell PosixShell) Run(command string, extraArgs []string, echoStdout, verbose bool, env []string) error {
	if len(extraArgs) > 0 {
		command = command + " " + strings.Join(extraArgs, " ")
	}

	if verbose {
		secondaryColor.Printf("Executing %s\n", command)
	}
	args := []string{"-c", command}

	cmd := exec.Command(shell.command, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stderr = os.Stderr
	if echoStdout {
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}

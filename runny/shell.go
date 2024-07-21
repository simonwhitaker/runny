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

type BashShell struct {
	command string
}

func NewShell(command string) (Shell, error) {
	switch path.Base(command) {
	case "bash":
		return BashShell{command: command}, nil
	}
	return nil, fmt.Errorf("unsupported shell: %s", command)
}

func (b BashShell) Run(command string, extraArgs []string, echoStdout, verbose bool, env []string) error {
	if len(extraArgs) > 0 {
		command = command + " " + strings.Join(extraArgs, " ")
	}

	if verbose {
		secondaryColor.Printf("Executing %s\n", command)
	}
	args := []string{"-c", command}

	cmd := exec.Command(b.command, args...)
	cmd.Env = env
	if echoStdout {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

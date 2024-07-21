package runny

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Shell interface {
	Run(string) error
}

type BashShell struct {
	command string
}

func NewShell(command string) Shell {
	switch path.Base(command) {
	case "bash":
		return BashShell{command: command}
	}
	panic(fmt.Sprintf("Shell not supported: %v", command))
}

func (b BashShell) Run(command string) error {
	args := []string{"-c", command}

	cmd := exec.Command(b.command, args...)
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

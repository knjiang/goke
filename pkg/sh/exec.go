package sh

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/craiggwilson/goke/task"
)

// ExitCode retrieves the exit code from an error.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	if eerr, ok := err.(*exec.ExitError); ok {
		return eerr.ExitCode()
	}

	return 1
}

// IsNotRan indicates if command that generated the error actually ran.
func IsNotRan(err error) bool {
	if err == nil {
		return false
	}
	if eerr, ok := err.(*exec.ExitError); ok {
		return !eerr.Exited()
	}
	return true
}

// Run the specified command piping its output to goke's output.
func Run(ctx *task.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = ctx
	cmd.Stdin = os.Stdin
	return RunCmd(ctx, cmd)
}

// RunOutput runs the specified command and get the command output.
func RunOutput(ctx *task.Context, name string, args ...string) (string, error) {
	var output bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &output
	cmd.Stderr = ctx
	cmd.Stdin = os.Stdin
	err := RunCmd(ctx, cmd)
	return strings.TrimRight(output.String(), "\r\n"), err
}

// RunCmd runs the provided command.
func RunCmd(ctx *task.Context, cmd *exec.Cmd) error {
	LogCmd(ctx, cmd)
	return cmd.Run()
}

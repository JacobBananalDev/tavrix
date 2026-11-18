//go:build linux || darwin

package runner

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
)

func shellCommand(ctx context.Context, script string) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, "/bin/sh", "-c", script), nil // #nosec G204 -- Explicit task execution intentionally runs the trusted task's shell script.
}

func configureCancellation(cmd *exec.Cmd, interactive bool) {
	if interactive {
		// Stay in the foreground group so tasks can read from the terminal.
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return err
	}
}

func exitCode(exit *exec.ExitError) int {
	if status, ok := exit.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		return 128 + int(status.Signal())
	}
	return exit.ExitCode()
}

package runner

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

func shellCommand(ctx context.Context, script string) (*exec.Cmd, error) {
	dir, err := windows.GetSystemDirectory()
	if err != nil {
		return nil, fmt.Errorf("locate Windows command shell: %w", err)
	}
	shell := filepath.Join(dir, "cmd.exe")
	cmd := exec.CommandContext(ctx, shell)
	// cmd.exe uses its own quoting rules, not Go's default Windows argument quoting.
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `"` + shell + `" /d /s /c "` + script + `"`}
	return cmd, nil
}

func configureCancellation(_ *exec.Cmd, _ bool) {}

func exitCode(exit *exec.ExitError) int { return exit.ExitCode() }

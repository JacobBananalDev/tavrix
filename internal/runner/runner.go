package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"
)

type Runner struct {
	Stdin  *os.File
	Stdout io.Writer
	Stderr io.Writer
}

func (r Runner) Run(ctx context.Context, dir, script string) (int, error) {
	cmd, err := shellCommand(ctx, script)
	if err != nil {
		return 1, err
	}
	cmd.Dir = dir
	if r.Stdin != nil {
		cmd.Stdin = r.Stdin
	}
	cmd.Stdout, cmd.Stderr = r.Stdout, r.Stderr
	cmd.WaitDelay = 2 * time.Second
	interactive := r.Stdin != nil && term.IsTerminal(int(r.Stdin.Fd()))
	configureCancellation(cmd, interactive)

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return 130, fmt.Errorf("task interrupted: %w", ctx.Err())
		}
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return exitCode(exit), err
		}
		return 1, fmt.Errorf("could not execute task: %w", err)
	}
	return 0, nil
}

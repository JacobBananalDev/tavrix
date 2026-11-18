package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/JacobBananalDev/tavrix/internal/config"
	"github.com/JacobBananalDev/tavrix/internal/runner"
)

const help = `Tavrix - A fast, lightweight task runner for developers.

Usage:
  tavrix <command> [arguments]
  tavrix <task>

Commands:
  init         Create a starter tavrix.yaml
  list         List tasks in tavrix.yaml
  run <task>   Run a task by name

Global flags:
  -h, --help   Show this help
  --version    Print the Tavrix version

Examples:
  tavrix init
  tavrix list
  tavrix run hello
  tavrix hello

Use 'tavrix run <task>' if a task name matches a built-in command.
`

func Run(ctx context.Context, args []string, version string, stdin *os.File, stdout, stderr io.Writer) int {
	return run(ctx, args, version, ".", runner.Runner{Stdin: stdin, Stdout: stdout, Stderr: stderr})
}

func run(ctx context.Context, args []string, version, dir string, r runner.Runner) int {
	out, errOut := r.Stdout, r.Stderr
	fail := func(err error) int {
		color, restore := colorSupport(errOut)
		defer restore()
		fmt.Fprintf(errOut, "%s %v\n", color.label("31", "Error:"), err)
		return 1
	}
	usage := func(message string) int {
		fmt.Fprintf(errOut, "Error: %s\nSee 'tavrix --help' for usage.\n", message)
		return 2
	}

	if len(args) == 0 {
		fmt.Fprint(out, help)
		return 0
	}

	path := filepath.Join(dir, config.Filename)
	var taskName string
	switch args[0] {
	case "-h", "--help", "--version", "init", "list":
		if len(args) != 1 {
			return usage(fmt.Sprintf("%s does not accept arguments", args[0]))
		}
		switch args[0] {
		case "-h", "--help":
			fmt.Fprint(out, help)
			return 0
		case "--version":
			fmt.Fprintf(out, "Tavrix v%s\n", version)
			return 0
		case "init":
			if err := initConfig(path); err != nil {
				return fail(err)
			}
			fmt.Fprintln(out, "Created tavrix.yaml. Try 'tavrix list' or 'tavrix hello'.")
			return 0
		}
	case "run":
		if len(args) != 2 || args[1] == "" || strings.HasPrefix(args[1], "-") {
			return usage("expected 'tavrix run <task>'")
		}
		taskName = args[1]
	default:
		if strings.HasPrefix(args[0], "-") {
			return usage(fmt.Sprintf("unknown flag %q", args[0]))
		}
		if len(args) != 1 || args[0] == "" {
			return usage("expected a single task name")
		}
		taskName = args[0]
	}

	cfg, err := config.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fail(fmt.Errorf("%s not found in the current directory; run 'tavrix init' to create it", config.Filename))
		}
		return fail(err)
	}
	if args[0] == "list" {
		color, restore := colorSupport(out)
		defer restore()
		fmt.Fprintln(out, color.label("1", "Available Tasks"))
		if len(cfg.Tasks) == 0 {
			fmt.Fprintln(out, "\n  (no tasks defined)")
			return 0
		}
		fmt.Fprintln(out)
		table := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		for _, name := range cfg.Names() {
			fmt.Fprintf(table, "  %s\t%s\n", name, cfg.Tasks[name].Description)
		}
		if err := table.Flush(); err != nil {
			return fail(fmt.Errorf("write task list: %w", err))
		}
		return 0
	}

	task, ok := cfg.Tasks[taskName]
	if !ok {
		return fail(fmt.Errorf("unknown task %q; run 'tavrix list' to see available tasks", taskName))
	}
	outColor, restoreOut := colorSupport(out)
	defer restoreOut()
	errColor, restoreErr := colorSupport(errOut)
	defer restoreErr()
	fmt.Fprintf(out, "%s %s\n", outColor.label("36", "==> Running task:"), taskName)
	start := time.Now()
	code, err := r.Run(ctx, dir, string(task.Command))
	duration := time.Since(start).Round(time.Millisecond)
	if err != nil {
		fmt.Fprintf(errOut, "%s %s (%s): %v\n", errColor.label("31", "Failed:"), taskName, duration, err)
		return code
	}
	fmt.Fprintf(out, "%s %s (%s)\n", outColor.label("32", "Completed:"), taskName, duration)
	return 0
}

func initConfig(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644) // #nosec G302 G304 -- Exclusive creation of the project's shareable starter file.
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("%s already exists; edit or move it before initializing", config.Filename)
	}
	if err != nil {
		return fmt.Errorf("create configuration: %w", err)
	}
	_, writeErr := io.WriteString(f, config.Starter)
	if err := errors.Join(writeErr, f.Close()); err != nil {
		return fmt.Errorf("write configuration (file may be incomplete): %w", err)
	}
	return nil
}

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JacobBananalDev/tavrix/internal/config"
	"github.com/JacobBananalDev/tavrix/internal/runner"
)

func invoke(t *testing.T, dir string, args ...string) (int, string, string) {
	t.Helper()
	var out, stderr bytes.Buffer
	code := run(t.Context(), args, "test", dir, runner.Runner{Stdout: &out, Stderr: &stderr})
	return code, out.String(), stderr.String()
}

func TestArguments(t *testing.T) {
	dir := t.TempDir()
	for _, tc := range []struct {
		name       string
		args       []string
		code       int
		wantOut    string
		wantStderr string
	}{
		{name: "no arguments", wantOut: "Usage:"},
		{name: "help", args: []string{"--help"}, wantOut: "Examples:"},
		{name: "version", args: []string{"--version"}, wantOut: "Tavrix vtest"},
		{name: "unknown flag", args: []string{"--unknown"}, code: 2, wantStderr: "unknown flag"},
		{name: "missing task", args: []string{"run"}, code: 2, wantStderr: "expected 'tavrix run <task>'"},
		{name: "extra argument", args: []string{"list", "extra"}, code: 2, wantStderr: "list does not accept arguments"},
		{name: "missing config", args: []string{"list"}, code: 1, wantStderr: "run 'tavrix init'"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, out, stderr := invoke(t, dir, tc.args...)
			if code != tc.code || !strings.Contains(out, tc.wantOut) || !strings.Contains(stderr, tc.wantStderr) {
				t.Fatalf("code = %d, stdout = %q, stderr = %q", code, out, stderr)
			}
		})
	}
}

func TestInitAndRun(t *testing.T) {
	dir := t.TempDir()
	if code, out, stderr := invoke(t, dir, "init"); code != 0 || !strings.Contains(out, "Created tavrix.yaml") {
		t.Fatalf("init: code = %d, stdout = %q, stderr = %q", code, out, stderr)
	}
	if code, out, stderr := invoke(t, dir, "hello"); code != 0 || !strings.Contains(out, "Hello from Tavrix!") ||
		!strings.Contains(out, "Completed: hello") || strings.Contains(out, "\x1b") {
		t.Fatalf("hello: code = %d, stdout = %q, stderr = %q", code, out, stderr)
	}
}

func TestInitRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, config.Filename)
	const original = "# do not overwrite"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := invoke(t, dir, "init"); code != 1 || !strings.Contains(stderr, "already exists") {
		t.Fatalf("init: code = %d, stderr = %q", code, stderr)
	}
	content, err := os.ReadFile(path)
	if err != nil || string(content) != original {
		t.Fatalf("configuration was changed: %q, error = %v", content, err)
	}
}

func TestTaskDispatch(t *testing.T) {
	dir := t.TempDir()
	const yaml = `version: "1"
tasks:
  z-last:
    description: Always exits with code 42
    command: exit 42
  list:
    description: A task that shadows a built-in
    command: echo explicit task
`
	if err := os.WriteFile(filepath.Join(dir, config.Filename), []byte(yaml), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, out, stderr := invoke(t, dir, "list"); code != 0 ||
		strings.Contains(out, "explicit task") || !strings.Contains(out, "A task that shadows a built-in") ||
		!strings.Contains(out, "z-last") || strings.Index(out, "  list") > strings.Index(out, "  z-last") {
		t.Fatalf("list: code = %d, stdout = %q, stderr = %q", code, out, stderr)
	}
	if code, out, stderr := invoke(t, dir, "run", "list"); code != 0 || !strings.Contains(out, "explicit task") {
		t.Fatalf("run list: code = %d, stdout = %q, stderr = %q", code, out, stderr)
	}
	if code, _, stderr := invoke(t, dir, "run", "z-last"); code != 42 || !strings.Contains(stderr, "42") {
		t.Fatalf("run z-last: code = %d, stderr = %q", code, stderr)
	}
	if code, _, stderr := invoke(t, dir, "run", "missing"); code != 1 || !strings.Contains(stderr, "unknown task") {
		t.Fatalf("run missing: code = %d, stderr = %q", code, stderr)
	}
}

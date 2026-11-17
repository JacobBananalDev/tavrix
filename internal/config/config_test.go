package config

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cfg, err := parse(strings.NewReader(`version: "1"
tasks:
  test:
    description: Run the tests
    command: go test ./...
  build:
    command: go build ./...
`))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(cfg.Names(), []string{"build", "test"}) {
		t.Fatalf("unexpected task order: %v", cfg.Names())
	}
	if task := cfg.Tasks["test"]; task.Description != "Run the tests" || task.Command != "go test ./..." {
		t.Fatalf("unexpected task: %+v", task)
	}
}

func TestInvalidConfig(t *testing.T) {
	const header = "version: \"1\"\n"
	for _, tc := range []struct{ name, input, want string }{
		{"empty", "", "configuration is empty"},
		{"malformed", header + "tasks: [", "yaml:"},
		{"missing version", "tasks: {}", "missing configuration version"},
		{"unsupported version", "version: \"2\"\ntasks: {}", "unsupported configuration version"},
		{"numeric version", "version: 1\ntasks: {}", "expected a string"},
		{"missing tasks", header, "tasks must be a mapping"},
		{"unknown top-level field", header + "tasks: {}\nextra: true", "field extra not found"},
		{"unknown task field", header + "tasks:\n  build: {command: echo hi, extra: true}", "field extra not found"},
		{"missing command", header + "tasks:\n  build: {}", `task "build": command is required`},
		{"nonstring command", header + "tasks:\n  build: {command: true}", "expected a string"},
		{"bad name", header + "tasks:\n  -build: {command: echo hi}", `invalid task name "-build"`},
		{"duplicate task", header + "tasks:\n  build: {}\n  build: {}", `mapping key "build" already defined`},
		{"extra document", header + "tasks: {}\n---\n", "expected one YAML document"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parse(strings.NewReader(tc.input))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), Filename))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error = %v, want os.ErrNotExist", err)
	}
}

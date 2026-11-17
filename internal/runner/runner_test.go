package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sample file"), []byte("sample contents\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	script := `cat "sample file"; echo diagnostic >&2; exit 23`
	if runtime.GOOS == "windows" {
		script = `type "sample file" & echo diagnostic>&2 & exit /b 23`
	}
	var out, stderr bytes.Buffer
	r := Runner{Stdout: &out, Stderr: &stderr}
	code, err := r.Run(t.Context(), dir, script)
	if code != 23 || err == nil {
		t.Fatalf("code = %d, error = %v; want 23 and an error", code, err)
	}
	if !strings.Contains(out.String(), "sample contents") || !strings.Contains(stderr.String(), "diagnostic") {
		t.Fatalf("stdout = %q, stderr = %q", &out, &stderr)
	}
}

type onWrite func()

func (f onWrite) Write(p []byte) (int, error) {
	f()
	return len(p), nil
}

func TestCancelRunningTask(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx, timeout := context.WithTimeout(ctx, 10*time.Second)
	defer timeout()
	script := `echo ready; sleep 30`
	if runtime.GOOS == "windows" {
		script = `for /l %i in (1,1,2147483647) do @echo ready`
	}
	r := Runner{Stdout: onWrite(cancel), Stderr: io.Discard}
	code, err := r.Run(ctx, t.TempDir(), script)
	if code != 130 || !errors.Is(err, context.Canceled) {
		t.Fatalf("code = %d, error = %v; want 130 and cancellation on first output", code, err)
	}
}

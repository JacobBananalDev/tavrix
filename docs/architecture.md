# Architecture

- [`cmd/tavrix`](../cmd/tavrix) is the entry point. It creates a signal-aware context, supplies the version and standard streams, and exits with the CLI's status. To override the version at build time, use `go build -ldflags="-X main.version=0.1.0" ./cmd/tavrix`.
- [`internal/cli`](../internal/cli) handles arguments, starter file creation, and terminal presentation. File creation uses an exclusive open rather than a separate existence check, so concurrent calls can't overwrite a file.
- [`internal/config`](../internal/config) decodes and validates YAML and supplies sorted task names. It does not start processes.
- [`internal/runner`](../internal/runner) runs a validated command string using the platform shell. It owns streaming, cancellation, and exit status handling. Shell selection and process control are kept in OS-specific files.

The CLI depends on config and runner; neither package depends on the other. YAML decoding uses `go.yaml.in/yaml/v3`. Terminal checks and platform APIs use `golang.org/x/term` and `golang.org/x/sys`.

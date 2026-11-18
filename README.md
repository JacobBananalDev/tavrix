# Tavrix

A fast, lightweight, cross-platform task runner for developers.

Put your project's everyday commands in a `tavrix.yaml` and run them by name. Tavrix validates the file, lists tasks alphabetically, and streams task output as it happens. It reports each task's duration and preserves the shell's exit code.

## Build from source

Tavrix supports **Linux**, **macOS**, and **Windows**. Neither contributors nor users need a system-wide Go installation. The built executable needs no Go, Bash, or Guix runtime; tasks still need the tools they invoke.

On Linux with Guix:

```sh
git clone https://github.com/JacobBananalDev/tavrix.git
cd tavrix
./scripts/shell.sh
# Inside the Guix shell:
./run.ysh build
```

Without Guix, download the [Go 1.27.4 archive](https://go.dev/dl/) for your OS and unpack it to a directory you control. From the Tavrix checkout, invoke that toolchain directly:

```sh
# Linux or macOS; substitute the path to your unpacked toolchain:
"$HOME/tools/go/bin/go" build ./cmd/tavrix
```

```powershell
# Windows PowerShell; substitute the path to your unpacked toolchain:
& "$HOME\tools\go\bin\go.exe" build ./cmd/tavrix
```

The binary is written to the checkout as `tavrix` (or `tavrix.exe` on Windows). Run it by its path, or put just the binary somewhere on your user `PATH`. The examples below assume it is on `PATH`.

## Quick start

In your project's root directory:

```sh
tavrix init       # Create tavrix.yaml; refuses if it already exists
tavrix list       # Show available tasks and descriptions
tavrix run hello  # Run the starter task
tavrix hello      # Same task, shorter syntax
tavrix --help
tavrix --version
```

Edit your file to add tasks:

```yaml
version: "1"
tasks:
  hello:
    description: Print a greeting
    command: echo Hello from Tavrix!
  test:
    description: Run the Go tests
    command: go test ./...
```

Tavrix reads `tavrix.yaml` from the current directory, not from parent directories. Built-ins (`init`, `list`, `run`) take precedence; use `tavrix run list` to execute a task named `list`. Commands run in `/bin/sh` on Linux and macOS, and `cmd.exe` on Windows. Use syntax supported by the target shell. See the [configuration guide](docs/configuration.md) for the full schema and exit behavior.

**Only run tasks from configurations you trust.** A task runs with your permissions and inherits your environment. Listing tasks validates the file without running commands. Tavrix doesn't print command text, but tasks can print anything they have access to.

## Develop and contribute

Use [`run.ysh`](run.ysh) for contributor tasks. On Linux, start the Guix shell below. Outside Guix, put Oils (`ysh`) and a user-local Go toolchain on your session's `PATH`.

```sh
./run.ysh -h                        # Show available tasks (same as --help)
./run.ysh list-tasks
./run.ysh lint --help               # Show help for a specific task
./run.ysh penv code .                # Open the repo in VS Code from this environment
./run.ysh fmt                        # Format Go packages
./run.ysh lint                       # Check YSH syntax, Go vet, Staticcheck, and Revive
./run.ysh check                      # Test, vet, and build
./run.ysh test -v ./internal/runner  # Run a focused test package
./run.ysh vuln                       # Run gosec and govulncheck
./run.ysh pre-commit                 # Check formatting, lint, scan, and run unit tests
```

`penv` runs the command from the repository root with your current environment. To open VS Code, `code` must already be on `PATH`; use the regular Guix shell rather than the container for desktop apps. The lint and security tools are pinned and run through Go without global installs. First use may need network access, and `govulncheck` fetches current advisories.

To check staged files automatically before each commit, opt in with `./run.ysh penv git config --local core.hooksPath .githooks`. Commit from an environment with Go and YSH available. The hook checks a temporary copy of the staged tree without formatting or changing it. If formatting fails, run `./run.ysh fmt` and stage the result.

`./run.ysh build` embeds `git describe --tags --always --dirty` in the binary's `--version` output. CI checks Go formatting, tests, vet, and build on all three supported platforms without requiring YSH. For larger changes, open an issue first, then submit a focused pull request with tests and updated documentation. See the [architecture guide](docs/architecture.md) for the package layout.

### Guix development shell (Linux)

Run [`scripts/shell.sh`](scripts/shell.sh) to enter a Guix development environment, or `scripts/shell.sh --container` (`-c`) for a container with networking, FHS emulation, and a writable root. The container uses the repository's `.bashrc` if one exists. The helper does not start another shell if `GUIX_ENVIRONMENT` is already set.

The helper uses the authenticated Guix revision pinned in [`channels.scm`](channels.scm). Guix supplies the bootstrap `go`, and the helper defaults `GOTOOLCHAIN` to `auto` so `go build` and `go test` can download 1.27.4 into Go's cache if the channel supplies an older version. The first run may need network access for Guix packages, the toolchain, and Go modules. Guix is contributor tooling only, not a Tavrix runtime feature.

AWS files are not mounted by the helper by default. To expose them read-only to a trusted container, set `TAVRIX_GUIX_EXPOSE_AWS=1` when launching it. Any process inside that container can read the exposed files.

## Roadmap

Potential additions after v0.1.0:

- Task dependencies, parallel execution, environment settings, `.env` files, and task-specific working directories.
- Configurable shells, YSH/Oils and Guix integration, interactive task selection, and shell completion.
- Caching, incremental execution, file watching, and execution history.
- Plugins, remote execution, and a web interface.

## License

[MIT](LICENSE).

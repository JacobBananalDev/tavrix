# Configuration reference

Create `tavrix.yaml` in the directory where you run Tavrix. There is no parent-directory search, alternate filename, or user-wide configuration in v0.1.0.

| Field | Type | Rules |
| --- | --- | --- |
| `version` | String | Required. Use `"1"` (quoted). |
| `tasks` | Mapping | Required. Use `tasks: {}` for no tasks. |
| `tasks.<name>.description` | String | Optional. One line, with no control characters. |
| `tasks.<name>.command` | String | Required. Must contain a non-whitespace character and cannot contain NUL. |

Task names start with an ASCII letter, digit, or underscore. After that they can also contain `.` and `-`. Names are case-sensitive. Tavrix rejects unknown fields, duplicate mapping keys, and multiple YAML documents. Quote strings when they might otherwise be parsed as another YAML type.

See the [example](../examples/tavrix.yaml). `tavrix list` sorts names alphabetically. If a task is named `init`, `list`, or `run`, invoke it explicitly with `tavrix run <name>`. Tavrix doesn't accept additional arguments to pass to a task.

## Shell and environment

Commands run in the current directory with inherited environment variables and stdin, stdout, and stderr. Linux and macOS use `/bin/sh -c`; Windows uses the system `cmd.exe` with AutoRun disabled. For portable tasks, favor a single command that is available on every target platform. For example, `go test ./...` works on all supported platforms if Go is installed. POSIX shell scripts and Windows command syntax are not interchangeable.

Tavrix never runs commands while parsing or listing a file. Running a task explicitly gives that command the same access to files, credentials, and network as the invoking user. Only run tasks from files you've reviewed and trust.

## Exit status and cancellation

- `0`: success (including help and version output).
- `1`: configuration, file, task lookup, or shell startup error.
- `2`: invalid CLI usage.
- A task's shell exit code is propagated. On Unix, an independently signaled shell maps to `128 + signal number`.
- `130`: Tavrix canceled a running task after an interrupt or termination request.

On Unix, canceling a task with non-terminal stdin kills its shell process group. Tasks with terminal stdin share Tavrix's foreground group so they can read input and receive terminal Ctrl+C directly. A signal sent only to Tavrix in this mode cancels the immediate shell; descendants may remain. On Windows, Tavrix cancels the immediate shell; it does not guarantee that all descendants exit. Detached processes are not tracked.

Tavrix disables its own ANSI styling on redirected streams, on `TERM=dumb`, and when `NO_COLOR` is set to a nonempty value. This does not affect programs run by a task.

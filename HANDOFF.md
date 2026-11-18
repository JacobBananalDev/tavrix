# Tavrix handoff

Paste this file into the coding agent on the other computer after opening the `JacobBananalDev/tavrix` checkout.

## Current state

- Tavrix is a Go CLI for running tasks defined in `tavrix.yaml`. The MVP includes `init`, `list`, `run <task>`, direct task invocation, help, and version output.
- At handoff, the worktree was clean before this note and `Master` and the local `origin/Master` ref both pointed to `6192c02` (`Add YSH developer tasks and pre-commit checks`). The MVP commit is `ea7908b`. There were no local release tags. Check the current state before making changes.
- Keep contributor commands behind [run.ysh](./run.ysh). Do not require a system-wide Go installation. The module requires Go 1.27.4; on Linux, [scripts/shell.sh](./scripts/shell.sh) opens the pinned Guix environment. The full pinned shell launch was not verified end to end.
- `./run.ysh -h` lists tasks; `./run.ysh penv code .` opens the repo in VS Code when `code` is available. `fmt`, `lint`, `vuln`, `unit-test`, and `pre-commit` are available. The [opt-in Git hook](./.githooks/pre-commit) checks staged files without modifying them. The full local `./run.ysh pre-commit` passed; `govulncheck` reported no affected calls and one advisory in a required module that Tavrix does not appear to call.
- Keep the quality checks Go/YSH-only for now: no Biome, `biome.sh`, SQL linting, or asset cache-busting. Guix and YSH are contributor tooling, not runtime requirements.

## Next task: fix Windows CI before tagging

1. Check the latest GitHub Actions run. In the [first run](https://github.com/JacobBananalDev/tavrix/actions/runs/44201382639), Linux and macOS passed. Windows failed at formatting, flagging all Go files before tests ran. Git stores LF and there was no EOL rule, so checkout conversion to CRLF is the likely cause.
2. Make a small follow-up change: add `*.go text eol=lf` to `.gitattributes` and verify the formatting check with a fresh checkout configured for CRLF conversion. Keep the Windows formatting check enabled. Then confirm that native Windows tests, vet, and build actually run and pass, alongside Linux and macOS.
3. Before releasing, validate the clean-clone Guix workflow without global Go. Smoke-test any binaries that will be distributed; if no binaries are ready, describe the release as source-only.
4. Once CI is green, tag the exact tested commit with an annotated `v0.1.0` and verify it with `git describe --tags --exact-match HEAD`. `./run.ysh build` already embeds the Git description as the binary version. Do not commit, push, tag, or publish on my behalf unless I explicitly ask.

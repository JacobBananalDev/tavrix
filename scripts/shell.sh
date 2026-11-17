#!/bin/bash

set -euo pipefail

if (( $# > 1 )); then
	echo "Usage: scripts/shell.sh [-c|--container]" >&2
	exit 1
fi

if (( $# == 1 )) && [[ $1 != -c && $1 != --container ]]; then
	printf 'Unrecognized argument: %s\n' "$1" >&2
	exit 1
fi

if [[ -n ${GUIX_ENVIRONMENT:-} ]]; then
	echo "The Guix shell is already active."
	exit 0
fi

if ! command -v guix >/dev/null; then
	echo "Guix is required to run this script." >&2
	exit 1
fi

cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.."
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
args=(time-machine --channels=channels.scm -- shell --manifest=manifest.scm)

if (( $# == 0 )); then
	shell_path=${SHELL:-/bin/bash}
	if entry=$(getent passwd "$(id -u)"); then
		shell_path=${entry##*:}
	fi
	login_shell=${shell_path##*/}
	case $login_shell in
		bash|fish|zsh) ;;
		*) login_shell=bash ;;
	esac
	exec guix "${args[@]}" "$login_shell" -- "$login_shell"
fi

args+=(--container --network --emulate-fhs --writable-root --preserve='^GOTOOLCHAIN$')
if [[ ${TAVRIX_GUIX_EXPOSE_AWS:-} == 1 ]]; then
	args+=("--expose=$HOME/.aws")
fi
if [[ -f .bashrc ]]; then
	exec guix "${args[@]}" -- bash --rcfile .bashrc
fi
exec guix "${args[@]}" -- bash

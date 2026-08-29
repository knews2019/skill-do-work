#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
module_dir="$script_dir/do-work-cli"
binary_path="$module_dir/do-work-cli"
minimum_go_version="1.26.1"

version_at_least() {
  awk -v minimum="$minimum_go_version" -v actual="$1" 'BEGIN {
    split(minimum, required_parts, ".")
    split(actual, actual_parts, ".")
    for (component = 1; component <= 3; component++) {
      required = required_parts[component] + 0
      found = actual_parts[component] + 0
      if (found > required) exit 0
      if (found < required) exit 1
    }
    exit 0
  }'
}

needs_build=0
if [ ! -x "$binary_path" ]; then
  needs_build=1
elif find "$module_dir" -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) -newer "$binary_path" -print -quit | grep -q .; then
  needs_build=1
fi

if [ "$needs_build" -eq 1 ]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "do-work-cli: Go $minimum_go_version or newer is required to build the command" >&2
    exit 2
  fi
  go_version_output="$(go version 2>/dev/null)" || {
    echo "do-work-cli: could not read the installed Go version; Go $minimum_go_version or newer is required" >&2
    exit 2
  }
  go_version="$(printf '%s\n' "$go_version_output" | awk '{print $3}' | sed 's/^go//')"
  if [ -z "$go_version" ] || ! version_at_least "$go_version"; then
    echo "do-work-cli: Go $minimum_go_version or newer is required (found ${go_version:-unknown})" >&2
    exit 2
  fi

  temporary_binary=""
  cleanup() {
    if [ -n "$temporary_binary" ]; then
      rm -f "$temporary_binary"
    fi
  }
  trap cleanup EXIT
  trap 'exit 129' HUP
  trap 'exit 130' INT
  trap 'exit 143' TERM
  temporary_binary="$(mktemp "$module_dir/do-work-cli.build.XXXXXX")"
  if ! (cd "$module_dir" && go build -o "$temporary_binary" ./cmd/do-work-cli); then
    echo "do-work-cli: build failed; stale output was not run" >&2
    exit 2
  fi
  chmod +x "$temporary_binary"
  mv -f "$temporary_binary" "$binary_path"
  trap - EXIT HUP INT TERM
fi

exec "$binary_path" "$@"

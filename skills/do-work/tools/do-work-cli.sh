#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
set -euo pipefail

script_dir="${BASH_SOURCE[0]%/*}"
[[ "$script_dir" != "${BASH_SOURCE[0]}" ]] || script_dir=.
script_dir="$(cd "$script_dir" && pwd -P)"
module_dir="$script_dir/do-work-cli"
minimum_go_version="1.25.0"

# Pure bash on purpose: this runs on every invocation, and restricted PATHs that
# carry only bash and go must still reach the command. A component keeps its
# leading digits ("26rc1" reads as 26); an absent component reads as 0.
version_at_least() {
  local -a required_parts found_parts
  local component required found
  IFS=. read -r -a required_parts <<<"$minimum_go_version"
  IFS=. read -r -a found_parts <<<"$1"
  for component in 0 1 2; do
    required="${required_parts[component]:-0}"; required="${required%%[^0-9]*}"
    found="${found_parts[component]:-0}"; found="${found%%[^0-9]*}"
    if [ "${found:-0}" -gt "${required:-0}" ]; then return 0; fi
    if [ "${found:-0}" -lt "${required:-0}" ]; then return 1; fi
  done
  return 0
}

if ! command -v go >/dev/null 2>&1; then
  echo "do-work-cli: Go $minimum_go_version or newer is required to run the command" >&2
  exit 2
fi
go_version_output="$(go version 2>/dev/null)" || {
  echo "do-work-cli: could not read the installed Go version; Go $minimum_go_version or newer is required" >&2
  exit 2
}
read -r _ _ go_version _ <<<"$go_version_output"
go_version="${go_version#go}"
if [ -z "$go_version" ] || ! version_at_least "$go_version"; then
  echo "do-work-cli: Go $minimum_go_version or newer is required (found ${go_version:-unknown})" >&2
  exit 2
fi

# The Go toolchain owns the build. `go tool` compiles the `tool` directive in
# go.mod, caches the linked executable in GOCACHE under a hash of every input
# (sources, go.mod, flags, toolchain), and reuses it while nothing changed. No
# binary lives in this tree and no timestamp decides staleness. `-n` prints the
# cached executable's path instead of running it, so the command is exec'd from
# the caller's own directory: `go tool -C` would otherwise run it inside the
# module, where relative arguments and the default repository root resolve to
# the wrong tree. `go run` would cache the same way but collapses every exit
# status to 1, and callers read the command's exact status.
tool_binary="$(go tool -C "$module_dir" -n do-work-cli)" || {
  echo "do-work-cli: the Go toolchain could not build the command; nothing was run" >&2
  exit 2
}
exec "$tool_binary" "$@"

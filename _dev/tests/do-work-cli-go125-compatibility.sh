#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
exact_go_toolchain="go1.25.0"

observed_go_version="$(GOTOOLCHAIN="$exact_go_toolchain" go env GOVERSION)"
if [ "$observed_go_version" != "$exact_go_toolchain" ]; then
  printf 'FAIL: requested %s but go selected %s\n' \
    "$exact_go_toolchain" "$observed_go_version" >&2
  exit 1
fi

(
  cd "$repo_root/skills/do-work/tools/do-work-cli"
  GOTOOLCHAIN="$exact_go_toolchain" go test -count=1 ./...
)

printf 'do-work-cli Go 1.25.0 compatibility tests passed\n'

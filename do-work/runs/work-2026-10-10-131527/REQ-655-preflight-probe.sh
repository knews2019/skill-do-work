#!/usr/bin/env bash
# Pre-flight baseline for REQ-655 (run before the build): the toolboxcommands package, which gains the ai-report-index verb and its tests, builds, vets and passes today.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
go vet -C "$module_dir" ./internal/toolboxcommands/
go test -C "$module_dir" -count=1 ./internal/toolboxcommands/

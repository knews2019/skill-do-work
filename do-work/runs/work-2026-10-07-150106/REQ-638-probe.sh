#!/usr/bin/env bash
# Focused tests for REQ-638: the release guard's shipped-change check.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -run '^TestRelease(RefusedWhenTheImplementationShipsNothing|Guard)' ./internal/finalization/

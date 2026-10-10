#!/usr/bin/env bash
# Pre-flight baseline for REQ-660 (run before the build): the two do-work-cli packages this REQ edits, internal/cleanup (where the worktree command goes) and internal/resultmodel (its typed result), pass today.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 ./internal/cleanup/ ./internal/resultmodel/

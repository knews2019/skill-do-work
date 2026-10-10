#!/usr/bin/env bash
# Pre-flight baseline for REQ-685 (run before the build): the gittransaction package, which holds the rollback code and tests this REQ edits, passes today.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 ./internal/gittransaction/

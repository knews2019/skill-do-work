#!/usr/bin/env bash
# Focused tests for REQ-634: the association walk, the associate handler, and the protected-inventory shim.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -run 'Associat|ProtectedInventory|CompatibilityShim|TerminalSuccess' ./internal/corehelpers/

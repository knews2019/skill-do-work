#!/usr/bin/env bash
# Focused repository gate for REQ-706 (user decision 2026-11-02: focused tests only,
# no maintainer-verify.sh, no heavy lanes). Runs the changed Go package and the
# contract scripts that cover the changed files, from whichever tree it runs in.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
go vet ./internal/toolboxcommands/
go test ./internal/toolboxcommands/ -count=1
cd "$root"
for script in shipped-package-reference-contract contract-regressions action-shell-blocks prescribed-shell-canonicalization suite-manifest-contract; do
  printf 'focused gate: %s\n' "$script"
  bash "_dev/tests/$script.sh"
done
printf 'focused gate: all green\n'

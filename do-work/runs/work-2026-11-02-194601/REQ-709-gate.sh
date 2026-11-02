#!/usr/bin/env bash
# Focused repository gate for REQ-709 (user decision 2026-11-02: focused tests only,
# no maintainer-verify.sh, no heavy lanes). Runs the changed Go packages and the
# contract scripts that cover shipped files.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
unformatted="$(gofmt -l ./internal/finalization ./internal/lifecycleadvance)"
if [ -n "$unformatted" ]; then printf 'gofmt: %s\n' "$unformatted"; exit 1; fi
go vet ./internal/finalization/ ./internal/lifecycleadvance/
go test ./internal/finalization/ ./internal/lifecycleadvance/ -count=1
cd "$root"
for script in shipped-package-reference-contract contract-regressions action-shell-blocks prescribed-shell-canonicalization suite-manifest-contract; do
  printf 'focused gate: %s\n' "$script"
  bash "_dev/tests/$script.sh"
done
printf 'focused gate: all green\n'

#!/usr/bin/env bash
# Focused repository gate for REQ-708 (user decision 2026-11-02: focused tests only,
# no maintainer-verify.sh, no heavy lanes). Runs the changed Go packages and the
# contract scripts (install-suite-behavior.sh is heavy-only, so the Go suiteinstall tests stand in).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
go vet ./internal/toolboxcommands/ ./internal/suiteinstall/
go test ./internal/toolboxcommands/ ./internal/suiteinstall/ -count=1
cd "$root"
for script in shipped-package-reference-contract contract-regressions action-shell-blocks prescribed-shell-canonicalization suite-manifest-contract; do
  printf 'focused gate: %s\n' "$script"
  bash "_dev/tests/$script.sh"
done
printf 'focused gate: all green\n'

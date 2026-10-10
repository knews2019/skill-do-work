#!/usr/bin/env bash
# GREEN checks for REQ-684 (integrator test gate): qualify carries the two distinct evidence strings, the old single message is gone, and the two tests (names contain ImplementationSummary) pass. Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
checks_source="$(<"$module_dir/internal/corehelpers/checks.go")"
grep -q 'Implementation Summary section not found' <<<"$checks_source" || { echo "checks.go lacks the 'section not found' evidence"; exit 1; }
grep -q 'Implementation Summary lists no backticked file paths' <<<"$checks_source" || { echo "checks.go lacks the 'lists no backticked file paths' evidence"; exit 1; }
if grep -q 'missing or empty' <<<"$checks_source"; then echo "checks.go still has the old 'missing or empty' message"; exit 1; fi
output="$(go test -C "$module_dir" -count=1 -v -run 'ImplementationSummary' ./internal/corehelpers/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
passed_tests="$(grep -c -- '^--- PASS: .*ImplementationSummary' <<<"$output" || true)"
[ "$passed_tests" -ge 2 ] || { printf '%s\n' "$output"; echo "expected at least 2 passing ImplementationSummary tests, saw $passed_tests"; exit 1; }

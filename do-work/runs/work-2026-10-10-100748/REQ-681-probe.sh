#!/usr/bin/env bash
# GREEN checks for REQ-681 (integrator test gate): sectionLineBounds is gone, the new test file exists, its cases (names contain AppendSectionEntry; CRLF, trailing spaces, fenced heading) pass, and the package vets. Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
state_dir="$module_dir/internal/requeststate"
test -f "$state_dir/append_section_entry_test.go" || { echo "append_section_entry_test.go is missing"; exit 1; }
if grep -q 'sectionLineBounds' "$state_dir/state_apply.go"; then echo "sectionLineBounds still exists in state_apply.go"; exit 1; fi
output="$(go test -C "$module_dir" -count=1 -v -run 'AppendSectionEntry' ./internal/requeststate/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
passed_cases="$(grep -c -- 'PASS: TestAppendSectionEntry' <<<"$output" || true)"
[ "$passed_cases" -ge 3 ] || { printf '%s\n' "$output"; echo "expected at least 3 passing AppendSectionEntry cases, saw $passed_cases"; exit 1; }
go vet -C "$module_dir" ./internal/requeststate/

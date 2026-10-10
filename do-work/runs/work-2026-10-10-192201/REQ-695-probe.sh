#!/usr/bin/env bash
# GREEN checks for REQ-695 (integrator test gate): a missing run-status run folder is named in text and JSON, a C3 row whose integrator is working does not recommend `do-work run`, the untouched C3 and no-run-folder pins stay green, the runstatus package is gofmt- and vet-clean, and both shipped C3 table cells name the integrator exception. Fails before the build (the two new tests do not exist yet).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
unformatted="$(gofmt -l "$module_dir/internal/runstatus")"
[ -z "$unformatted" ] || { printf 'gofmt -l lists: %s\n' "$unformatted"; exit 1; }
go vet -C "$module_dir" ./internal/runstatus/
test_names="TestRunStatusMissingRunFolderIsNamedInTextAndJSON TestRunStatusIntegratingHandbackDoesNotRecommendRun TestRunStatusLandedHandbackIsClassC3 TestRunStatusReportsQueueRowsWithoutARunDirectory TestRunStatusWatchFitsTwentyLinesForEightOpenRequests"
run_pattern="^($(tr ' ' '|' <<<"$test_names"))\$"
output="$(go test -C "$module_dir" -count=1 -v -run "$run_pattern" ./internal/runstatus/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for test_name in $test_names; do
  grep -qF -- "--- PASS: $test_name (" <<<"$output" || { printf 'no PASS line for %s:\n%s\n' "$test_name" "$output"; exit 1; }
done
for doc_path in skills/do-work/actions/status.md skills/do-work/docs/status-guide.md; do
  c3_line="$(grep -F 'C3' "$root/$doc_path" || true)"
  grep -qi 'integrator' <<<"$c3_line" || { printf '%s: no C3 line names the integrator exception\n' "$doc_path"; exit 1; }
done
printf 'REQ-695 GREEN: missing run folder named, integrating C3 recommends no run, docs updated\n'

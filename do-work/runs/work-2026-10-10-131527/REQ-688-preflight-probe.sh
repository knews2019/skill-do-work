#!/usr/bin/env bash
# Pre-flight baseline for REQ-688 (run before the build): the existing capture-example, containment and handler
# registration tests the build extends pass at base and are not skipped. The new example test does not exist yet.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
probe_tests=(TestBuildCapturePlanAcceptsPublishedCaptureExamples TestOutsideTextContainmentUsesLongerFence TestHandlersRegisterEveryPublicationCommand TestBuildCapturePlanBootstrapsAbsentDoWorkAndPinsRawInput)
run_pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -v -run "$run_pattern" ./internal/publication/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
echo "REQ-688 pre-flight baseline: ok"

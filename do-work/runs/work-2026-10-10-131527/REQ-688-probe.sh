#!/usr/bin/env bash
# GREEN checks for REQ-688 (integrator test gate): capture-files --example, filled in, passes a dry run with a
# triple-backtick raw input; the published capture-reference UR example passes the raw-input containment check;
# capture Step 5 names the example command. A skip is not a pass. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
probe_tests=(TestBuildCapturePlanAcceptsPublishedCaptureExamples TestCaptureFilesExampleFilledInPassesDryRun TestOutsideTextContainmentUsesLongerFence TestHandlersRegisterEveryPublicationCommand)
run_pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -v -run "$run_pattern" ./internal/publication/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
grep -qF -- "--- PASS: TestBuildCapturePlanAcceptsPublishedCaptureExamples/UR_input " <<<"$output" || { printf '%s\n' "$output"; echo "UR input example case did not pass"; exit 1; }
if grep -qF -- "--- SKIP" <<<"$output"; then printf '%s\n' "$output"; echo "a probe test was skipped"; exit 1; fi
reference="$(cat "$root/skills/do-work/actions/capture-reference.md")"
if grep -qF -- '> ````text' <<<"$reference"; then echo "capture-reference.md still shows the four-backtick text fence"; exit 1; fi
grep -qxF -- '> ```' <<<"$reference" || { echo "capture-reference.md UR example has no bare three-backtick fence line"; exit 1; }
capture_action="$(cat "$root/skills/do-work/actions/capture.md")"
grep -qF -- 'capture-files --example' <<<"$capture_action" || { echo "capture.md Step 5 does not name capture-files --example"; exit 1; }
echo "REQ-688 GREEN probe: ok"

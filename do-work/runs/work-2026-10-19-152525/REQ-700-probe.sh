#!/usr/bin/env bash
# GREEN checks for REQ-700 (integrator test gate): the second-opinion session form has a prepare
# phase that never execs the agent's own CLI and a finalize phase that publishes the agent-written
# page with every launcher guarantee. Names only this REQ's two tests; a skip is not a pass.
# Expected to FAIL at the pre-build baseline (the tests do not exist yet). Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
probe_tests=(TestSecondOpinionSessionPrepareNeverExecsOwnBackend TestSecondOpinionSessionFinalizePublishesAgentWrittenPage)
pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$module" -count=1 -v -run "$pattern" ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
echo "REQ-700 probe: both session-form tests passed"

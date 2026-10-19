#!/usr/bin/env bash
# GREEN checks for REQ-701 (integrator test gate): a subagent model in the Claude reply envelope never makes the backend model ambiguous, and a model check that fails after the page is written publishes the page under the unverified name instead of deleting it. Names only this REQ's Go tests; a skip or a missing test is not a pass. Under 30 s on a quiet machine. Expected to FAIL at base: two of the three tests do not exist yet and the third still carries its old name.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
probe_tests=(TestSecondOpinionSubagentModelIsNotAmbiguous TestSecondOpinionUnverifiedModelKeepsFinishedPage TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified)
output="$(go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -v -run "^($(IFS='|'; echo "${probe_tests[*]}"))\$" ./internal/toolboxcommands/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done

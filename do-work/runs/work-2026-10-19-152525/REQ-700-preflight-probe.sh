#!/usr/bin/env bash
# Pre-flight baseline for REQ-700 (run before the build): the existing second-opinion tests that the
# split (prepare/finalize, digest-aware preservation, extracted publication tail) must keep green, including
# the four REQ-701 tests (released in 0.306.3 before this build), plus
# the handler-count pin (flags, not a new command) and the typed-summary round trip. Named tests only;
# a skip is not a pass. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
toolbox_tests=(TestSecondOpinionNativeModelCollisionAndSkippedJudge TestSecondOpinionTargetNeverSubstitutes TestSecondOpinionOutsideDirtAndOriginalMutationRefuse TestSecondOpinionBatchIncludesMissingAndExcluded TestSecondOpinionJudgeVerdictsGateOriginal TestSecondOpinionAllNativeAdaptersAndSerialFailureContinuation TestSecondOpinionPreservesAssetsAndRejectsUnprefixedSymlinkAndMissingBacklink TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified TestSecondOpinionSubagentModelIsNotAmbiguous TestSecondOpinionUnverifiedModelKeepsFinishedPage TestSecondOpinionUnverifiedModelWithoutPageFails TestSecondOpinionBackendCommandSafetyContract TestHandlersRegisterCanonicalToolboxCommands)
toolbox_pattern="^($(IFS='|'; echo "${toolbox_tests[*]}"))\$"
output="$(go test -C "$module" -count=1 -v -run "$toolbox_pattern" ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${toolbox_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
model_output="$(go test -C "$module" -count=1 -v -run '^TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence$' ./internal/resultmodel 2>&1)" || { printf '%s\n' "$model_output"; exit 1; }
grep -q -- '^--- PASS: TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence ' <<<"$model_output" || { printf '%s\n' "$model_output"; echo "TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence did not pass (skipped or missing)"; exit 1; }
echo "REQ-700 pre-flight: existing second-opinion tests, handler pin and typed-summary round trip passed"

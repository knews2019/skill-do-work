#!/usr/bin/env bash
# Pre-flight baseline for REQ-703 (run before the build): every existing second-opinion test (plain naming, unverified
# naming, session prepare/finalize, banner, judge gate, capture), every report-catalog test the presets-folder skip must
# keep green, the handler-count pin (flags, not a new command) and the typed-summary round trip. Named tests only;
# a skip is not a pass. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
toolbox_tests=(TestSecondOpinionNativeModelCollisionAndSkippedJudge TestSecondOpinionTargetNeverSubstitutes TestSecondOpinionOutsideDirtAndOriginalMutationRefuse TestSecondOpinionBatchIncludesMissingAndExcluded TestSecondOpinionBannerOnlyRealBodyAndPerFileIdempotency TestSecondOpinionJudgeVerdictsGateOriginal TestSecondOpinionAllNativeAdaptersAndSerialFailureContinuation TestSecondOpinionRejectsUnverifiedIdentityAndInvalidAGYCompletion TestSecondOpinionPreservesAssetsAndRejectsUnprefixedSymlinkAndMissingBacklink TestSecondOpinionLatestAndLinkedBundleRefusal TestSecondOpinionOutsideRenameIncludesBothSides TestSecondOpinionCapturedRunnerCancellation TestSecondOpinionBackendCommandSafetyContract TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified TestSecondOpinionSubagentModelIsNotAmbiguous TestSecondOpinionUnverifiedModelWithoutPageFails TestSecondOpinionUnverifiedModelKeepsFinishedPage TestSecondOpinionSessionPrepareNeverExecsOwnBackend TestSecondOpinionSessionFinalizePublishesAgentWrittenPage TestSecondOpinionBacklinkAllowsSemanticAttributeAndInlineVariants TestSecondOpinionCaptureDrainsExitedLeaderChildren TestSecondOpinionCaptureStreamsProgressWhileBackendRuns TestSecondOpinionCaptureCapsBackendOutput TestSecondOpinionCapturePreservesBackendStatusAndChannels TestAIReportIndexCatalogsEveryBundleNamingStyleOnce TestAIReportIndexLinksSupersededProposalToSuccessor TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded TestAIReportIndexRefusesHandMadeCatalog TestHandlersRegisterCanonicalToolboxCommands)
toolbox_pattern="^($(IFS='|'; echo "${toolbox_tests[*]}"))\$"
output="$(go test -C "$module" -count=1 -v -run "$toolbox_pattern" ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${toolbox_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
model_output="$(go test -C "$module" -count=1 -v -run '^TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence$' ./internal/resultmodel 2>&1)" || { printf '%s\n' "$model_output"; exit 1; }
grep -q -- '^--- PASS: TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence ' <<<"$model_output" || { printf '%s\n' "$model_output"; echo "TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence did not pass (skipped or missing)"; exit 1; }
echo "REQ-703 pre-flight: existing second-opinion and report-catalog tests, handler pin and typed-summary round trip passed"

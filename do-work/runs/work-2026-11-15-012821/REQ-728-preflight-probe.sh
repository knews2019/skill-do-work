#!/usr/bin/env bash
# Pre-build baseline for REQ-728 (ai-report-judge renders one named page of a bundle): the
# existing judge tests (browser cases included, so DO_WORK_HEAVY_TESTS=1) and the second-opinion
# tests that call the judge through its stub (the finalize verdict path), run at the base revision. Every named judge test
# must PASS; a skip is not a pass. Needs a Chromium-family browser (Google Chrome on this machine).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"

fail() { printf 'REQ-728 pre-flight FAIL: %s\n' "$*" >&2; exit 1; }

judge_output="$(DO_WORK_HEAVY_TESTS=1 go test -count=1 -v -run '^TestAIReportJudge' ./internal/toolboxcommands/ 2>&1)" \
  || fail "judge tests failed: $judge_output"
for test_name in TestAIReportJudgeFailsOnPhoneWidthOverflow TestAIReportJudgeFailsOnBrokenRelativeImage \
  TestAIReportJudgePassesCleanBundleAndStopsServer TestAIReportJudgeReportsSkippedWithoutBrowser; do
  grep -qF -- "--- PASS: $test_name " <<<"$judge_output" || fail "$test_name did not PASS: $judge_output"
done

opinion_tests='TestSecondOpinionJudgeVerdictsGateOriginal TestSecondOpinionNativeModelCollisionAndSkippedJudge TestSecondOpinionSessionFinalizePublishesAgentWrittenPage TestSecondOpinionChangeListTargetsAreChecked'
opinion_pattern="^($(tr ' ' '|' <<<"$opinion_tests"))\$"
opinion_output="$(go test -count=1 -v -run "$opinion_pattern" ./internal/toolboxcommands/ 2>&1)" \
  || fail "second-opinion judge tests failed: $opinion_output"
for test_name in $opinion_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$opinion_output" || fail "$test_name did not PASS: $opinion_output"
done

printf 'REQ-728 pre-flight PASS\n'

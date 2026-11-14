#!/usr/bin/env bash
# GREEN probe for REQ-728 (ai-report-judge renders one named page of a bundle with --page).
# Checks whichever tree it runs in. Expected to FAIL at base 1389e2be: the two new tests do not
# exist yet, so their PASS lines are missing. Needs a Chromium-family browser for the heavy case.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-728 probe FAIL: %s\n' "$*" >&2; exit 1; }

# 1. The judge form in the action and the guide name the option.
grep -qF -- '--page' skills/do-work-toolbox/actions/ai-report.md || fail 'ai-report.md does not name --page'
grep -qF -- '--page' skills/do-work-toolbox/docs/ai-report-guide.md || fail 'ai-report-guide.md does not name --page'

cd skills/do-work/tools/do-work-cli
unformatted="$(gofmt -l internal/toolboxcommands/report_judge.go internal/toolboxcommands/report_judge_test.go)"
[ -z "$unformatted" ] || fail "gofmt reports: $unformatted"

# 2. The new tests and the existing judge tests (no-option behavior unchanged), browser cases included.
judge_tests='TestAIReportJudgeRendersNamedPageInPlace TestAIReportJudgeRefusesMissingOrOutsidePage TestAIReportJudgeFailsOnPhoneWidthOverflow TestAIReportJudgeFailsOnBrokenRelativeImage TestAIReportJudgePassesCleanBundleAndStopsServer TestAIReportJudgeReportsSkippedWithoutBrowser'
judge_output="$(DO_WORK_HEAVY_TESTS=1 go test -count=1 -v -run '^TestAIReportJudge' ./internal/toolboxcommands/ 2>&1)" \
  || fail "judge tests failed: $judge_output"
for test_name in $judge_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$judge_output" || fail "$test_name did not PASS (a skip is not a pass): $judge_output"
done

# 3. The second-opinion finalize path keeps its judge verdicts.
opinion_tests='TestSecondOpinionJudgeVerdictsGateOriginal TestSecondOpinionNativeModelCollisionAndSkippedJudge TestSecondOpinionSessionFinalizePublishesAgentWrittenPage TestSecondOpinionChangeListTargetsAreChecked'
opinion_pattern="^($(tr ' ' '|' <<<"$opinion_tests"))\$"
opinion_output="$(go test -count=1 -v -run "$opinion_pattern" ./internal/toolboxcommands/ 2>&1)" \
  || fail "second-opinion judge tests failed: $opinion_output"
for test_name in $opinion_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$opinion_output" || fail "$test_name did not PASS: $opinion_output"
done

printf 'REQ-728 probe PASS\n'

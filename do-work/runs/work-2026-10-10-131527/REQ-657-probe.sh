#!/usr/bin/env bash
# REQ-657 GREEN probe: the ai-report-judge verb's focused tests pass with a real
# browser (none skipped), the verb is registered, and every render-check action
# plus the routing, help and guide files point at the command.
# Checks whichever tree it runs in. Expected to FAIL at the base commit.
set -euo pipefail

repository_root="$(git rev-parse --show-toplevel)"
module_directory="$repository_root/skills/do-work/tools/do-work-cli"

test_output="$(DO_WORK_HEAVY_TESTS=1 go test -C "$module_directory" -count=1 -v \
  -run '^(TestAIReportJudge.*|TestHandlersRegisterCanonicalToolboxCommands)$' \
  ./internal/toolboxcommands/ 2>&1)" || {
  printf '%s\n' "$test_output"
  printf 'REQ-657 probe: focused tests failed\n' >&2
  exit 1
}
printf '%s\n' "$test_output"

# A skipped browser case is not a pass.
if grep -q -- '--- SKIP:' <<<"$test_output"; then
  printf 'REQ-657 probe: a focused test skipped (no browser, or the heavy gate is off)\n' >&2
  exit 1
fi
for required_test_name in \
  TestAIReportJudgeFailsOnPhoneWidthOverflow \
  TestAIReportJudgeFailsOnBrokenRelativeImage \
  TestAIReportJudgePassesCleanBundleAndStopsServer \
  TestAIReportJudgeReportsSkippedWithoutBrowser \
  TestHandlersRegisterCanonicalToolboxCommands; do
  if ! grep -q -- "--- PASS: ${required_test_name} " <<<"$test_output"; then
    printf 'REQ-657 probe: missing PASS line for %s\n' "$required_test_name" >&2
    exit 1
  fi
done

# Prose wiring: each file must name the command.
for wired_file in \
  skills/do-work-toolbox/actions/ai-report.md \
  skills/do-work-toolbox/actions/architecture-report.md \
  skills/do-work-toolbox/actions/stakeholder-report.md \
  skills/do-work/docs/command-line-guide.md; do
  if ! grep -q -F 'ai-report-judge' "$repository_root/$wired_file"; then
    printf 'REQ-657 probe: %s does not call ai-report-judge\n' "$wired_file" >&2
    exit 1
  fi
done
for routed_file in \
  skills/do-work-toolbox/SKILL.md \
  skills/do-work-toolbox/actions/help.md \
  skills/do-work-toolbox/docs/ai-report-guide.md; do
  if ! grep -q -F 'ai-report judge' "$repository_root/$routed_file"; then
    printf 'REQ-657 probe: %s does not name the ai-report judge form\n' "$routed_file" >&2
    exit 1
  fi
done
printf 'REQ-657 probe: GREEN\n'

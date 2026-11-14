#!/usr/bin/env bash
# REQ-729 GREEN probe: the judge's usage text and the ai-report judge step each state
# that a relative --out (and a relative bundle path) starts at --repo-root.
# Expected to FAIL at base (neither sentence exists yet).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
judge_go="$root/skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go"
action_md="$root/skills/do-work-toolbox/actions/ai-report.md"
rule_pattern='relative .*--out.* starts at .?--repo-root.?, not the current folder'

usage_line="$(grep -F 'Usage: ai-report-judge' "$judge_go" || true)"
if ! grep -Eq "$rule_pattern" <<<"$usage_line"; then
  printf 'REQ-729 probe FAIL: judgeUsage in report_judge.go does not state the relative --out rule\n' >&2
  exit 1
fi
action_text="$(cat "$action_md")"
if ! grep -Eq "$rule_pattern" <<<"$action_text"; then
  printf 'REQ-729 probe FAIL: ai-report.md does not state the relative --out rule\n' >&2
  exit 1
fi
# The package still compiles and the browser-free judge test passes with the new text.
test_output="$(cd "$root/skills/do-work/tools/do-work-cli" && go test -count=1 -v -run '^TestAIReportJudgeReportsSkippedWithoutBrowser$' ./internal/toolboxcommands/ 2>&1)" || {
  printf '%s\nREQ-729 probe FAIL: go test failed\n' "$test_output" >&2
  exit 1
}
if ! grep -q -- '--- PASS: TestAIReportJudgeReportsSkippedWithoutBrowser' <<<"$test_output"; then
  printf '%s\nREQ-729 probe FAIL: the judge test did not report PASS\n' "$test_output" >&2
  exit 1
fi
printf 'REQ-729 probe PASS\n'

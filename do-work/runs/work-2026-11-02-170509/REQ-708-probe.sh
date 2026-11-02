#!/usr/bin/env bash
# GREEN probe for REQ-708: install, update, the catalog and the report writers recommend
# ai-report organize only when finished months are still loose.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestReportOrganizeHintPrintsOnlyWhenWorkRemains|TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose|TestInstallNarrationNeverReachesStandardOutput)$' ./internal/toolboxcommands/ ./internal/suiteinstall/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- '--- PASS: TestReportOrganizeHintPrintsOnlyWhenWorkRemains' <<<"$output"
grep -q -- '--- PASS: TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose' <<<"$output"
grep -q -- '--- PASS: TestInstallNarrationNeverReachesStandardOutput' <<<"$output"
cd "$root"
for action in skills/do-work-toolbox/actions/ai-report.md skills/do-work-toolbox/actions/architecture-report.md skills/do-work-toolbox/actions/stakeholder-report.md; do
  text="$(cat "$action")"
  grep -q 'ai-report-organize --check' <<<"$text" || { printf '%s does not run the organize check\n' "$action"; exit 1; }
done
version="$(cat skills/do-work/actions/version.md)"
grep -q 'ai-report organize' <<<"$version"
printf 'REQ-708 GREEN: organize recommended by install, update, index and report writers\n'

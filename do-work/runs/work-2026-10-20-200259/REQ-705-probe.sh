#!/usr/bin/env bash
# GREEN checks for REQ-705 (integrator test gate): a user preset resolves from <bundle>/second-opinion-presets/ and
# the shared ai-reports/second-opinion-presets/ folder is no longer read, the catalog skip and its test are gone, and
# no shipped file outside the changelogs and lessons still names the shared folder. Names only this REQ's tests plus
# the two REQ-703 preset tests it edits or must keep green; a skip is not a pass.
# Expected to FAIL at the pre-build baseline (the bundle-only test does not exist yet). Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
commands="$module/internal/toolboxcommands"
index_source="$(cat "$commands/report_index.go")"
if grep -q -- 'secondOpinionPresetsFolder' <<<"$index_source"; then echo "report_index.go still names secondOpinionPresetsFolder (catalog skip not removed)"; exit 1; fi
index_tests="$(cat "$commands/report_index_test.go")"
if grep -q -- 'TestAIReportIndexSkipsSecondOpinionPresetsFolder' <<<"$index_tests"; then echo "the catalog skip test still exists"; exit 1; fi
stale_mentions="$(grep -rl --exclude=CHANGELOG.md --exclude='lessons-*.md' --exclude='*_test.go' -- 'ai-reports/second-opinion-presets' "$root/skills" || true)"
[ -z "$stale_mentions" ] || { printf 'shipped files still name the shared preset folder:\n%s\n' "$stale_mentions"; exit 1; }
guide_text="$(cat "$root/skills/do-work-toolbox/docs/ai-report-guide.md")"
grep -q -- 'second-opinion-presets/<name>.md' <<<"$guide_text" || { echo "ai-report-guide.md does not name the bundle preset location"; exit 1; }
probe_tests=(TestSecondOpinionPresetResolvesFromBundleOnly TestSecondOpinionPresetResolutionAndArgumentRules TestSecondOpinionDirectedRunNamesAndMarksPage TestAIReportIndexCatalogsEveryBundleNamingStyleOnce)
pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$module" -count=1 -v -run "$pattern" ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
echo "REQ-705 probe: bundle-only preset lookup, catalog skip removal and shared-folder sweep checked"

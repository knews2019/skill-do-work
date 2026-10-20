#!/usr/bin/env bash
# GREEN checks for REQ-704 (integrator test gate): a second opinion can target a tree page, tree pages get v1cNN names,
# a parent link and tree metas, the launcher checks the #rev-* change list, index.html lists the tree in pre-order, a
# plain run keeps today's name, and REQ-703's directed names gained their tree id. Also checks the shipped prompt
# placeholders and the direction exception list, and that the action and guide describe tree pages.
# Names only this REQ's four tests; a skip is not a pass.
# Expected to FAIL at the pre-build baseline (the tests, placeholders and prose do not exist yet). Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
prompt_text="$(cat "$root/skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md")"
grep -q -- '{{BASE_PAGE}}' <<<"$prompt_text" || { echo "shipped prompt has no {{BASE_PAGE}} placeholder"; exit 1; }
grep -q -- '{{CHANGE_LIST_REQUIREMENT}}' <<<"$prompt_text" || { echo "shipped prompt has no {{CHANGE_LIST_REQUIREMENT}} placeholder"; exit 1; }
grep -q -- '^DIRECTION\..*CHANGE LIST' <<<"$prompt_text" || { echo "shipped prompt DIRECTION paragraph does not keep CHANGE LIST out of the direction's reach"; exit 1; }
action_text="$(cat "$root/skills/do-work-toolbox/actions/ai-report.md")"
for phrase in 'v1c01' 'base_page' 'tree_id' 'ask-user prompt'; do
  grep -q -- "$phrase" <<<"$action_text" || { echo "ai-report.md does not mention $phrase"; exit 1; }
done
guide_text="$(cat "$root/skills/do-work-toolbox/docs/ai-report-guide.md")"
grep -q -- 'v1c' <<<"$guide_text" || { echo "ai-report-guide.md does not describe tree page names"; exit 1; }
probe_tests=(TestSecondOpinionTreeTargetsNamesAndParentHeader TestSecondOpinionChangeListTargetsAreChecked TestSecondOpinionTreeListOnIndexFollowsTreeOrder TestSecondOpinionDirectedRunNamesAndMarksPage)
pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$module" -count=1 -v -run "$pattern" ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
echo "REQ-704 probe: shipped prompt and prose checked; tree target, change-list, tree-list and directed-name tests passed"

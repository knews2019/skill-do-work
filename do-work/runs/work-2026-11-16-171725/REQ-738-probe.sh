#!/usr/bin/env bash
# GREEN probe for REQ-738 (the ai-report second-opinion launcher shows a short header and folds
# the full direction, so the page content starts on the first phone screen).
# Checks whichever tree it runs in. Expected to FAIL at base 7f2c4b94: the new test does not exist
# yet, so the first check stops the probe.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-738 probe FAIL: %s\n' "$*" >&2; exit 1; }

package_folder=skills/do-work/tools/do-work-cli
launcher_folder=internal/toolboxcommands

# 1. The RED-first test exists by name (REQ-738 D-04).
test_text="$(cat "$package_folder/$launcher_folder/report_second_opinion_test.go")"
grep -qF -- 'func TestSecondOpinionHeaderFoldsTheFullDirection(' <<<"$test_text" \
  || fail 'report_second_opinion_test.go has no TestSecondOpinionHeaderFoldsTheFullDirection'

# 2. The two contract sentences no longer describe an always-open direction block, and the action
#    still names the header markers.
action_text="$(cat skills/do-work-toolbox/actions/ai-report.md)"
guide_text="$(cat skills/do-work-toolbox/docs/ai-report-guide.md)"
if grep -qF -- 'showing the preset and the full direction;' <<<"$action_text"; then fail 'ai-report.md still says the header shows the full direction'; fi
if grep -qF -- 'the page shows the direction at its top' <<<"$guide_text"; then fail 'ai-report-guide.md still says the page shows the direction at its top'; fi
grep -qF -- '<!-- second-opinion-header -->' <<<"$action_text" || fail 'ai-report.md no longer names the header markers'

# 3. The launcher still writes the same markers, label and long-token guard.
launcher_text="$(cat "$package_folder/$launcher_folder/report_second_opinion.go")"
grep -qF -- '"<!-- second-opinion-header -->", "<!-- /second-opinion-header -->"' <<<"$launcher_text" || fail 'the launcher no longer uses the header markers'
grep -qF -- 'overflow-wrap:anywhere' <<<"$launcher_text" || fail 'the launcher header lost its long-token guard'

# 4. gofmt is clean on the three Go files of the write set.
cd "$package_folder"
unformatted="$(gofmt -l "$launcher_folder/report_second_opinion.go" "$launcher_folder/report_second_opinion_test.go" "$launcher_folder/report_second_opinion_tree_test.go")"
[ -z "$unformatted" ] || fail "gofmt reports: $unformatted"

# 5. The new test and the two tests that pin the header on real runs pass. A skip is not a pass.
header_tests='TestSecondOpinionHeaderFoldsTheFullDirection TestSecondOpinionDirectedRunNamesAndMarksPage TestSecondOpinionTreeTargetsNamesAndParentHeader'
header_pattern="^($(tr ' ' '|' <<<"$header_tests"))\$"
header_output="$(go test -count=1 -v -run "$header_pattern" ./"$launcher_folder"/ 2>&1)" \
  || fail "header tests failed: $header_output"
for test_name in $header_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$header_output" || fail "$test_name did not PASS (a skip is not a pass): $header_output"
done

printf 'REQ-738 probe PASS (3 header tests, 2 contract sentences)\n'

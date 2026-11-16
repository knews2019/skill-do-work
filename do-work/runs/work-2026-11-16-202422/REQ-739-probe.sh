#!/usr/bin/env bash
# GREEN probe for REQ-739 (the second-opinion header keeps a side margin on a phone, without
# overriding a page that styles the header itself). Checks whichever tree it runs in.
# Expected to FAIL at base 49e54593: the new test does not exist yet, so the first check stops it.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-739 probe FAIL: %s\n' "$*" >&2; exit 1; }

package_folder=skills/do-work/tools/do-work-cli
launcher_folder=internal/toolboxcommands

# 1. The RED-first test exists by name (REQ-739 D-04).
test_text="$(cat "$package_folder/$launcher_folder/report_second_opinion_test.go")"
grep -qF -- 'func TestSecondOpinionHeaderKeepsAPhoneSideMargin(' <<<"$test_text" \
  || fail 'report_second_opinion_test.go has no TestSecondOpinionHeaderKeepsAPhoneSideMargin'

# 2. The launcher keeps its markers, and the aside open tag carries only the long-token guard, so a
#    page's own rule for the header still wins (no inline padding or margin).
launcher_text="$(cat "$package_folder/$launcher_folder/report_second_opinion.go")"
grep -qF -- '"<!-- second-opinion-header -->", "<!-- /second-opinion-header -->"' <<<"$launcher_text" || fail 'the launcher no longer uses the header markers'
grep -qF -- '<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">' <<<"$launcher_text" || fail 'the aside open tag changed: it must keep only the long-token guard inline'

# 3. The action still names the header markers.
action_text="$(cat skills/do-work-toolbox/actions/ai-report.md)"
grep -qF -- '<!-- second-opinion-header -->' <<<"$action_text" || fail 'ai-report.md no longer names the header markers'

# 4. gofmt is clean on the two Go files of the write set.
cd "$package_folder"
unformatted="$(gofmt -l "$launcher_folder/report_second_opinion.go" "$launcher_folder/report_second_opinion_test.go")"
[ -z "$unformatted" ] || fail "gofmt reports: $unformatted"

# 5. The new test and the three tests that pin the header pass. A skip is not a pass.
header_tests='TestSecondOpinionHeaderKeepsAPhoneSideMargin TestSecondOpinionHeaderFoldsTheFullDirection TestSecondOpinionDirectedRunNamesAndMarksPage TestSecondOpinionTreeTargetsNamesAndParentHeader'
header_pattern="^($(tr ' ' '|' <<<"$header_tests"))\$"
header_output="$(go test -count=1 -v -run "$header_pattern" ./"$launcher_folder"/ 2>&1)" \
  || fail "header tests failed: $header_output"
for test_name in $header_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$header_output" || fail "$test_name did not PASS (a skip is not a pass): $header_output"
done

printf 'REQ-739 probe PASS (4 header tests, aside tag unchanged)\n'

#!/usr/bin/env bash
# GREEN probe for REQ-727 (ship a relation-map preset for ai-report second-opinion).
# Checks whichever tree it runs in. Expected to FAIL at base 445950c4: the preset file does not
# exist yet, so the first check stops the probe.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-727 probe FAIL: %s\n' "$*" >&2; exit 1; }

# 1. The preset ships, loads (non-empty, no {{), and says how it differs from the mind map.
preset=skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md
[ -s "$preset" ] || fail "$preset is missing or empty"
preset_text="$(cat "$preset")"
if grep -qF -- '{{' <<<"$preset_text"; then fail "$preset contains {{, which the loader refuses"; fi
grep -qiE -- 'mind ?map' <<<"$preset_text" || fail "$preset does not say how it differs from the mind map preset"

# 2. The guide's preset list names it.
guide_text="$(cat skills/do-work-toolbox/docs/ai-report-guide.md)"
# shellcheck disable=SC2016 # the backticks are literal Markdown, not an expansion
grep -qF -- '`relation-map`' <<<"$guide_text" || fail 'ai-report-guide.md does not name `relation-map`'

# 3. The shipped-presets test names the preset, so a deleted file fails instead of shrinking the loop.
cd skills/do-work/tools/do-work-cli
test_text="$(cat internal/toolboxcommands/report_second_opinion_test.go)"
grep -qF -- '"relation-map"' <<<"$test_text" || fail 'report_second_opinion_test.go does not name "relation-map"'
unformatted="$(gofmt -l internal/toolboxcommands/report_second_opinion_test.go)"
[ -z "$unformatted" ] || fail "gofmt reports: $unformatted"

# 4. The preset tests pass: prepare with --preset relation-map returns prepared with the shipped preset_file.
preset_tests='TestSecondOpinionShippedPresetsLoadFromTheToolbox TestSecondOpinionPresetResolutionAndArgumentRules TestSecondOpinionPresetResolvesFromBundleOnly'
preset_pattern="^($(tr ' ' '|' <<<"$preset_tests"))\$"
preset_output="$(go test -count=1 -v -run "$preset_pattern" ./internal/toolboxcommands/ 2>&1)" \
  || fail "preset tests failed: $preset_output"
for test_name in $preset_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$preset_output" || fail "$test_name did not PASS (a skip is not a pass): $preset_output"
done

# 5. One real page was built with the preset, in a bundle other than the example's.
cd "$root"
pages="$(grep -rlF --include='index-*.html' -- 'name="ai-report-second-opinion-preset" content="relation-map"' ai-reports || true)"
other_pages="$(grep -vF -- '2026-11-14_1547_proposal-tinka-review-surface' <<<"$pages" || true)"
[ -n "$other_pages" ] || fail 'no second-opinion page built with --preset relation-map outside the example bundle'

printf 'REQ-727 probe PASS (%s)\n' "$(head -n 1 <<<"$other_pages")"

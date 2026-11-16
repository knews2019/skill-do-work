#!/usr/bin/env bash
# Pre-build baseline for REQ-737 (the relation-map preset keeps box titles readable after a click
# on a phone), run at the base revision. It pins what must still hold after the change: the preset
# ships with checks 1 to 20 numbered in order, it stays inside the byte bound, the shipped-presets
# Go test passes, and the 0.324.0 before page is in the tstack bundle. A skip is not a pass.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-737 pre-flight FAIL: %s\n' "$*" >&2; exit 1; }

preset=skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md
[ -s "$preset" ] || fail "$preset is missing or empty"
preset_text="$(cat "$preset")"
preset_bytes="$(wc -c <"$preset" | tr -d ' ')"
[ "$preset_bytes" -le 3600 ] || fail "$preset is $preset_bytes bytes, above the 3600-byte bound"
for check_number in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do
  grep -qE -- "^${check_number}\\. " <<<"$preset_text" || fail "check $check_number is not a numbered line in $preset"
done

bundle=ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1
before_page="$bundle/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html"
[ -s "$before_page" ] || fail "the before page $before_page is missing"

cd skills/do-work/tools/do-work-cli
test_name=TestSecondOpinionShippedPresetsLoadFromTheToolbox
test_output="$(go test -count=1 -v -run "^${test_name}\$" ./internal/toolboxcommands/ 2>&1)" \
  || fail "$test_name failed: $test_output"
grep -qF -- "--- PASS: $test_name " <<<"$test_output" || fail "$test_name did not PASS (a skip is not a pass): $test_output"

printf 'REQ-737 pre-flight PASS (%s bytes, 20 checks)\n' "$preset_bytes"

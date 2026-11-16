#!/usr/bin/env bash
# GREEN probe for REQ-737 (the relation-map preset keeps box titles readable after a click on a
# phone). Checks whichever tree it runs in. Expected to FAIL at base 7f2c4b94: no numbered check
# of the 0.324.0 preset names both a click and the phone, so check 2 stops the probe.
# The title sizes themselves are measured in a browser by the builder and the reviewer; a shell
# probe pins only what the files can show.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-737 probe FAIL: %s\n' "$*" >&2; exit 1; }

# 1. The preset still ships whole: loadable (non-empty, no {{), inside the byte bound (target
#    3,500, refused above 3,600), and checks 1 to 20 numbered in order with no gap.
preset=skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md
[ -s "$preset" ] || fail "$preset is missing or empty"
preset_text="$(cat "$preset")"
if grep -qF -- '{{' <<<"$preset_text"; then fail "$preset contains {{, which the loader refuses"; fi
preset_bytes="$(wc -c <"$preset" | tr -d ' ')"
[ "$preset_bytes" -le 3600 ] || fail "$preset is $preset_bytes bytes, above the 3600-byte bound"
numbered_lines="$(grep -E -- '^[0-9]+\. ' <<<"$preset_text" || true)"
numbered_count="$(grep -cE -- '^[0-9]+\. ' <<<"$preset_text" || true)"
[ "$numbered_count" -ge 20 ] || fail "$preset has $numbered_count numbered checks, fewer than 20"
expected_number=1
while IFS= read -r numbered_line; do
  [ "${numbered_line%%.*}" = "$expected_number" ] || fail "check numbers are not in order: expected $expected_number at: $numbered_line"
  expected_number=$((expected_number + 1))
done <<<"$numbered_lines"

# 2. One numbered check states the phone rule for a click: the same line names a click or tap,
#    and the phone or the 600 px width.
click_lines="$(grep -iE -- 'click|tap' <<<"$numbered_lines" || true)"
grep -qiE -- 'phone|600 px' <<<"$click_lines" || fail "no numbered check in $preset names both a click and the phone (or 600 px)"

# 3. The shipped-presets Go test still passes, by name.
cd skills/do-work/tools/do-work-cli
test_name=TestSecondOpinionShippedPresetsLoadFromTheToolbox
test_output="$(go test -count=1 -v -run "^${test_name}\$" ./internal/toolboxcommands/ 2>&1)" \
  || fail "$test_name failed: $test_output"
grep -qF -- "--- PASS: $test_name " <<<"$test_output" || fail "$test_name did not PASS (a skip is not a pass): $test_output"

# 4. A proof page built with the relation-map preset sits in the tstack bundle beside the
#    0.324.0 before page (v1c03), which is still there.
cd "$root"
bundle=ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1
before_name=index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html
[ -s "$bundle/$before_name" ] || fail "the before page $bundle/$before_name is missing"
pages="$(grep -lF -- 'name="ai-report-second-opinion-preset" content="relation-map"' "$bundle"/index-v1c*.html || true)"
proof_pages="$(grep -vF -- "$before_name" <<<"$pages" || true)"
[ -n "$proof_pages" ] || fail "no relation-map page other than v1c03 in $bundle"

printf 'REQ-737 probe PASS (%s bytes, %s checks, %s)\n' "$preset_bytes" "$numbered_count" "$(head -n 1 <<<"$proof_pages")"

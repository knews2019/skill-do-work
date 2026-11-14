#!/usr/bin/env bash
# REQ-733 pre-build baseline: the Tinka proposal bundle is intact at base and has no v1c03c01 page yet.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
bundle="$repo_root/ai-reports/2026-11-14_1547_proposal-tinka-review-surface"
failures=0
fail() { printf 'FAIL: %s\n' "$1"; failures=$((failures + 1)); }

while read -r expected_blob relative_name; do
  [[ -f "$bundle/$relative_name" ]] || { fail "missing $relative_name"; continue; }
  actual_blob="$(git hash-object "$bundle/$relative_name")"
  [[ "$actual_blob" == "$expected_blob" ]] || fail "$relative_name changed bytes"
done <<'LIST'
d215b8eba00c6ca399af9318901632aa4369f212 index-v1c01-claude-claude-fable-5-2-custom-worked-example-tally.html
46b1010aec016f680cf5fb8553db03f36d6915b4 index-v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines.html
7d35c89fcafa362014877e69b485a874f5e6f823 index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html
ce7c5c56e64e668e360654f1027cc1f8508271fe index.html
b12a97516edcfa38d0e5305f95f280e09805807c presentation.css
f109728ebacfabca5f73baded221822b977138fc presentation.js
ec3063e2dfaadd8e87e0f408410ba668a9056732 v1c01-claude-claude-fable-5-2-custom-worked-example-tally-reply.md
b228374da55b86ecdc472e852fa28c4a527a133e v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines-reply.md
eb860181ecf48fa0bd0d3f7d0cc7afaf16aefaea v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map-reply.md
LIST

shopt -s nullglob
existing_children=("$bundle"/index-v1c03c01-*.html)
[[ ${#existing_children[@]} -eq 0 ]] || fail "a v1c03c01 page already exists at base"

if [[ $failures -eq 0 ]]; then
  printf 'REQ-733 preflight PASS\n'
else
  exit 1
fi

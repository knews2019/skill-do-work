#!/usr/bin/env bash
# REQ-734 pre-build baseline: the Tinka proposal bundle is intact at base (REQ-733's v1c03c01 page
# included) and has no v1c04 table page yet.
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
9fafa0d2dbf9bf2354ee24c4d81778af71a647ab index.html
d215b8eba00c6ca399af9318901632aa4369f212 index-v1c01-claude-claude-fable-5-2-custom-worked-example-tally.html
46b1010aec016f680cf5fb8553db03f36d6915b4 index-v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines.html
7d35c89fcafa362014877e69b485a874f5e6f823 index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html
20bec32cc321f0e15052d09ada7cfe5e3e96befb index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html
b12a97516edcfa38d0e5305f95f280e09805807c presentation.css
f109728ebacfabca5f73baded221822b977138fc presentation.js
2b0ef49a6b4e6bbc043c22714feb82dc599ed1f9 screenshots/canvas-lineage-overview.png
2ec8a613b60c2e681f01aba3277e581385044e08 screenshots/s1-probe-html-tile-full-view.png
b462f6b437a6cdb4f4b0f3efb1afbecdd20feaa9 screenshots/tile-capture-full-view.png
345146f39a4432e243faf8a333d8c5bab1cf09c6 screenshots/tile-reply-markdown-full-view.png
ec3063e2dfaadd8e87e0f408410ba668a9056732 v1c01-claude-claude-fable-5-2-custom-worked-example-tally-reply.md
b228374da55b86ecdc472e852fa28c4a527a133e v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines-reply.md
eb860181ecf48fa0bd0d3f7d0cc7afaf16aefaea v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map-reply.md
b69f6f43657e8d75f63cf54b1c85d0d8071a6d91 v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map-reply.md
LIST

shopt -s nullglob
existing_pages=("$bundle"/index-v1c04-*.html)
[[ ${#existing_pages[@]} -eq 0 ]] || fail "a v1c04 page already exists at base"

if [[ $failures -eq 0 ]]; then
  printf 'REQ-734 preflight PASS\n'
else
  exit 1
fi

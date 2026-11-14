#!/usr/bin/env bash
# REQ-733 GREEN probe: one finalized v1c03c01 map page beside the first map page, the existing
# bundle files keep their bytes, and index.html changed only inside its iteration list.
# Fails at base by design (no v1c03c01 page yet).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
bundle="$repo_root/ai-reports/2026-11-14_1547_proposal-tinka-review-surface"
base_page="index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html"
failures=0
fail() { printf 'FAIL: %s\n' "$1"; failures=$((failures + 1)); }

# 1. Every pre-existing file except index.html keeps its exact bytes.
while read -r expected_blob relative_name; do
  [[ -f "$bundle/$relative_name" ]] || { fail "missing $relative_name"; continue; }
  actual_blob="$(git hash-object "$bundle/$relative_name")"
  [[ "$actual_blob" == "$expected_blob" ]] || fail "$relative_name changed bytes"
done <<'LIST'
d215b8eba00c6ca399af9318901632aa4369f212 index-v1c01-claude-claude-fable-5-2-custom-worked-example-tally.html
46b1010aec016f680cf5fb8553db03f36d6915b4 index-v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines.html
7d35c89fcafa362014877e69b485a874f5e6f823 index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html
b12a97516edcfa38d0e5305f95f280e09805807c presentation.css
f109728ebacfabca5f73baded221822b977138fc presentation.js
2b0ef49a6b4e6bbc043c22714feb82dc599ed1f9 screenshots/canvas-lineage-overview.png
2ec8a613b60c2e681f01aba3277e581385044e08 screenshots/s1-probe-html-tile-full-view.png
b462f6b437a6cdb4f4b0f3efb1afbecdd20feaa9 screenshots/tile-capture-full-view.png
345146f39a4432e243faf8a333d8c5bab1cf09c6 screenshots/tile-reply-markdown-full-view.png
ec3063e2dfaadd8e87e0f408410ba668a9056732 v1c01-claude-claude-fable-5-2-custom-worked-example-tally-reply.md
b228374da55b86ecdc472e852fa28c4a527a133e v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines-reply.md
eb860181ecf48fa0bd0d3f7d0cc7afaf16aefaea v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map-reply.md
LIST

# 2. index.html outside the second-opinion-tree block is byte-identical to base.
outside_tree="$(sed '/<!-- second-opinion-tree -->/,/<!-- \/second-opinion-tree -->/d' "$bundle/index.html" | shasum -a 256)"
[[ "${outside_tree%% *}" == "dbb595ce5acd18dc8b90b9ca9855d1f3d4ed7b9c70df0f7b8e9fc629a630a757" ]] || fail "index.html changed outside the iteration list"

# 3. Exactly one v1c03c01 page and reply, and no staging leftovers.
shopt -s nullglob
new_pages=("$bundle"/index-v1c03c01-*.html)
new_replies=("$bundle"/v1c03c01-*-reply.md)
staging_leftovers=("$bundle"/.*staging*)
[[ ${#new_pages[@]} -eq 1 ]] || fail "expected one index-v1c03c01-*.html, found ${#new_pages[@]}"
[[ ${#new_replies[@]} -eq 1 ]] || fail "expected one v1c03c01-*-reply.md, found ${#new_replies[@]}"
[[ ${#staging_leftovers[@]} -eq 0 ]] || fail "staging files left in the bundle: ${staging_leftovers[*]}"

if [[ ${#new_pages[@]} -eq 1 ]]; then
  page_text="$(cat "${new_pages[0]}")"
  page_name="$(basename "${new_pages[0]}")"
  # 4. The launcher finalized it as a child of the first map page.
  grep -q '<meta name="ai-report-second-opinion-id" content="v1c03c01">' <<<"$page_text" || fail "page lacks the v1c03c01 id meta"
  grep -qF "<meta name=\"ai-report-second-opinion-base\" content=\"$base_page\">" <<<"$page_text" || fail "page base meta is not the first map page"
  # 5. Inline SVG, no library: no remote script or stylesheet.
  grep -q '<svg' <<<"$page_text" || fail "page has no inline SVG"
  if grep -qiE '<script[^>]+src="(https?:)?//' <<<"$page_text"; then fail "page loads a remote script"; fi
  if grep -qiE '<link[^>]+href="(https?:)?//' <<<"$page_text"; then fail "page loads a remote stylesheet"; fi
  # 6. index.html links the new page inside its iteration list.
  index_text="$(cat "$bundle/index.html")"
  grep -qF "href=\"$page_name\"" <<<"$index_text" || fail "index.html iteration list does not link $page_name"
fi

if [[ $failures -eq 0 ]]; then
  printf 'REQ-733 probe PASS\n'
else
  exit 1
fi

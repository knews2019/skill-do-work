#!/usr/bin/env bash
# REQ-734 GREEN probe: one finalized v1c04 supports-and-limits table page built on index.html,
# carrying all 37 relations (r01-r37, D-04) with no script of its own (D-05); every existing bundle
# file keeps its bytes, index.html changed only inside its iteration list, and no comparison
# picture landed in the bundle (D-07). Fails at base by design (no v1c04 page yet).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
bundle="$repo_root/ai-reports/2026-11-14_1547_proposal-tinka-review-surface"
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

# 2. index.html outside the second-opinion-tree block is byte-identical to base.
outside_tree="$(sed '/<!-- second-opinion-tree -->/,/<!-- \/second-opinion-tree -->/d' "$bundle/index.html" | shasum -a 256)"
[[ "${outside_tree%% *}" == "dbb595ce5acd18dc8b90b9ca9855d1f3d4ed7b9c70df0f7b8e9fc629a630a757" ]] || fail "index.html changed outside the iteration list"

# 3. Exactly one v1c04 page and reply, no staging leftovers, and no other new file (no picture).
shopt -s nullglob
new_pages=("$bundle"/index-v1c04-*.html)
new_replies=("$bundle"/v1c04-*-reply.md)
staging_leftovers=("$bundle"/.*staging*)
[[ ${#new_pages[@]} -eq 1 ]] || fail "expected one index-v1c04-*.html, found ${#new_pages[@]}"
[[ ${#new_replies[@]} -eq 1 ]] || fail "expected one v1c04-*-reply.md, found ${#new_replies[@]}"
[[ ${#staging_leftovers[@]} -eq 0 ]] || fail "staging files left in the bundle: ${staging_leftovers[*]}"
bundle_file_count="$(find "$bundle" -type f | wc -l | tr -d ' ')"
[[ "$bundle_file_count" == "17" ]] || fail "bundle holds $bundle_file_count files, expected 17 (15 at base plus page and reply)"

if [[ ${#new_pages[@]} -eq 1 ]]; then
  page_text="$(cat "${new_pages[0]}")"
  page_name="$(basename "${new_pages[0]}")"
  # 4. The launcher finalized it as a child of the original report.
  grep -q '<meta name="ai-report-second-opinion-id" content="v1c04">' <<<"$page_text" || fail "page lacks the v1c04 id meta"
  grep -qF '<meta name="ai-report-second-opinion-base" content="index.html">' <<<"$page_text" || fail "page base meta is not index.html"
  # 5. The two columns and every one of the 37 relations.
  grep -q '>Supports' <<<"$page_text" || fail "page has no Supports heading"
  grep -q '>Limits' <<<"$page_text" || fail "page has no Limits heading"
  for relation_number in $(seq -w 1 37); do
    grep -qF "data-relation=\"r$relation_number\"" <<<"$page_text" || fail "relation r$relation_number is missing"
  done
  if grep -qE 'data-relation="r(3[89]|[4-9][0-9])"' <<<"$page_text"; then fail "page numbers a relation past r37"; fi
  # 6. No script of its own and no remote file: the only script allowed is the bundle's presentation.js.
  script_tags="$(grep -oE '<script[^>]*>' <<<"$page_text" || true)"
  if grep -vqxF '<script src="presentation.js">' <<<"$script_tags"; then
    [[ -z "$script_tags" ]] || fail "page has a script other than presentation.js"
  fi
  if grep -qiE '<link[^>]+href="(https?:)?//' <<<"$page_text"; then fail "page loads a remote stylesheet"; fi
  # 7. index.html links the new page inside its iteration list.
  index_text="$(cat "$bundle/index.html")"
  grep -qF "href=\"$page_name\"" <<<"$index_text" || fail "index.html iteration list does not link $page_name"
fi

if [[ $failures -eq 0 ]]; then
  printf 'REQ-734 probe PASS\n'
else
  exit 1
fi

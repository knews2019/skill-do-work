#!/usr/bin/env bash
# REQ-725 GREEN probe: the memory-note destination report exists and gives every note in the
# maintainer's local memory folder exactly one row with a valid destination.
# Expected to FAIL at base: no ai-reports/*_proposal-REQ-725-* bundle exists yet.
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
memory_folder="/Users/t2/.claude/projects/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/memory"

fail() { printf 'REQ-725 probe FAIL: %s\n' "$1" >&2; exit 1; }

[ -d "$memory_folder" ] || fail "memory folder absent on this machine: $memory_folder"

shopt -s nullglob
report_pages=("$repo_root"/ai-reports/*_proposal-REQ-725-*/index.html "$repo_root"/ai-reports/*/*_proposal-REQ-725-*/index.html)
shopt -u nullglob
[ "${#report_pages[@]}" -eq 1 ] || fail "expected exactly one ai-reports/*_proposal-REQ-725-*/index.html, found ${#report_pages[@]}"
report_page="${report_pages[0]}"
page_text="$(cat "$report_page")"
# Each count below wraps grep in { grep ... || true; } because grep exits 1 on zero matches and
# pipefail would end the script with no message instead of the named failure.

grep -q '<meta name="ai-report-kind" content="proposal">' <<<"$page_text" || fail "missing ai-report-kind proposal meta"

# Every note file (MEMORY.md is the index, not a note) has exactly one row.
note_count=0
for note_path in "$memory_folder"/*.md; do
  note_name="$(basename "$note_path")"
  [ "$note_name" = "MEMORY.md" ] && continue
  note_count=$((note_count + 1))
  row_count="$({ grep -o -F "data-memory-note=\"$note_name\"" <<<"$page_text" || true; } | wc -l | tr -d ' ')"
  [ "$row_count" -eq 1 ] || fail "note $note_name has $row_count rows, expected 1"
done
[ "$note_count" -gt 0 ] || fail "memory folder holds no notes"

# No extra rows: the row total equals the note count.
total_rows="$({ grep -o 'data-memory-note="[^"]*"' <<<"$page_text" || true; } | wc -l | tr -d ' ')"
[ "$total_rows" -eq "$note_count" ] || fail "report has $total_rows note rows, folder has $note_count notes"

# Every row carries one of the three destinations, and nothing else.
destination_count="$({ grep -o -E 'data-destination="(move-to-repository|move-to-skill|delete)"' <<<"$page_text" || true; } | wc -l | tr -d ' ')"
any_destination_count="$({ grep -o 'data-destination="[^"]*"' <<<"$page_text" || true; } | wc -l | tr -d ' ')"
[ "$destination_count" -eq "$note_count" ] || fail "valid destinations: $destination_count, notes: $note_count"
[ "$any_destination_count" -eq "$note_count" ] || fail "destination attributes: $any_destination_count, notes: $note_count"

printf 'REQ-725 probe PASS: %s notes, one row each, in %s\n' "$note_count" "${report_page#"$repo_root"/}"

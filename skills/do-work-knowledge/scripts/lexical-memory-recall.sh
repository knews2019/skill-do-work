#!/usr/bin/env bash
# Rank bounded lexical memory matches from a sanitized query.
set -u

if [ "$#" -ne 2 ]; then
  printf 'Usage: %s <memory-directory> <query-text>\n' "$0" >&2
  exit 2
fi

memory_directory="$1"
query_text="$2"
if [ ! -d "$memory_directory" ]; then
  printf 'Memory directory does not exist: %s\n' "$memory_directory" >&2
  exit 2
fi

temporary_root="$(mktemp -d "${TMPDIR:-/tmp}/do-work-memory-recall.XXXXXX")" || exit 2
trap 'rm -rf "$temporary_root"' EXIT HUP INT TERM
token_file="$temporary_root/tokens"
source_file="$temporary_root/sources"
match_file="$temporary_root/matches"

printf '%s' "$query_text" \
  | LC_ALL=C tr '[:upper:]' '[:lower:]' \
  | LC_ALL=C tr -cs 'a-z0-9' '\n' \
  | awk 'length($0) >= 3 && !seen[$0]++' > "$token_file"
if [ ! -s "$token_file" ]; then
  exit 0
fi

: > "$source_file"
if [ -f "$memory_directory/working-memory.md" ]; then
  printf '%s\n' "$memory_directory/working-memory.md" >> "$source_file"
fi
if [ -d "$memory_directory/logs" ]; then
  find "$memory_directory/logs" -type f -name '*.md' -print | LC_ALL=C sort >> "$source_file"
fi
if [ ! -s "$source_file" ]; then
  exit 0
fi

if cutoff_week="$(date -u -v-7d +%F 2>/dev/null)"; then :
else cutoff_week="$(date -u -d '7 days ago' +%F 2>/dev/null)" || cutoff_week='9999-99-99'
fi
if cutoff_month="$(date -u -v-30d +%F 2>/dev/null)"; then :
else cutoff_month="$(date -u -d '30 days ago' +%F 2>/dev/null)" || cutoff_month='9999-99-99'
fi

: > "$match_file"
while IFS= read -r memory_source; do
  awk -v cutoff_week="$cutoff_week" -v cutoff_month="$cutoff_month" '
    FILENAME == ARGV[1] { tokens[++token_count] = $0; next }
    /^##[[:space:]]/ { current_heading = $0 }
    {
      lowered = tolower($0)
      distinct_hits = 0
      for (token_index = 1; token_index <= token_count; token_index++) {
        if (index(lowered, tokens[token_index]) > 0) distinct_hits++
      }
      if (distinct_hits == 0) next
      source_name = FILENAME
      if (source_name ~ /\/working-memory\.md$/) {
        date_label = "working memory"
        weight = 4
      } else {
        date_label = source_name
        sub(/^.*\//, "", date_label)
        sub(/\.md$/, "", date_label)
        if (date_label >= cutoff_week) weight = 3
        else if (date_label >= cutoff_month) weight = 2
        else weight = 1
      }
      if (current_heading == "") current_heading = "(no heading)"
      gsub(/\t/, " ", current_heading)
      content = $0
      gsub(/\t/, " ", content)
      printf "%d\t%s:%d\t%s\t%s\t%s\n", distinct_hits * weight, FILENAME, FNR, date_label, current_heading, content
    }
  ' "$token_file" "$memory_source" >> "$match_file"
done < "$source_file"

LC_ALL=C sort -t "$(printf '\t')" -k1,1nr -k2,2 "$match_file" | head -n 8

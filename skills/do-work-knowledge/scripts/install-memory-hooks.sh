#!/usr/bin/env bash
# Compose the two optional memory hooks into existing Claude settings.
set -u

if [ "$#" -ne 2 ]; then
  printf 'Usage: %s <project-root> <memory-hooks-fragment>\n' "$0" >&2
  exit 2
fi

project_root="$1"
hooks_fragment="$2"
settings_file="$project_root/.claude/settings.json"
backup_file="$settings_file.pre-memory-module"
merge_file="$settings_file.merge-tmp"

if [ ! -f "$hooks_fragment" ]; then
  printf 'Memory hook fragment is missing: %s\n' "$hooks_fragment" >&2
  exit 2
fi
if ! command -v jq >/dev/null 2>&1; then
  cat "$hooks_fragment"
  printf 'hooks: MANUAL STEP\n'
  exit 0
fi

mkdir -p "$(dirname "$settings_file")" || exit 2
if [ -f "$settings_file" ]; then
  jq . "$settings_file" >/dev/null 2>&1 || {
    printf 'Claude settings are invalid JSON: %s\n' "$settings_file" >&2
    exit 2
  }
else
  printf '{}\n' > "$settings_file" || exit 2
fi

append_session_start=1
append_stop=1
grep -q 'memory-session-start.sh' "$settings_file" && append_session_start=0
grep -q 'memory-stop-capture.sh' "$settings_file" && append_stop=0
if [ "$append_session_start" -eq 0 ] && [ "$append_stop" -eq 0 ]; then
  printf 'memory hooks already present — skipping merge\n'
  exit 0
fi

cp "$settings_file" "$backup_file" || exit 2
if jq --slurpfile fragment "$hooks_fragment" \
  --argjson add_start "$append_session_start" --argjson add_stop "$append_stop" \
  '(if $add_start == 1 then .hooks.SessionStart = ((.hooks.SessionStart // []) + $fragment[0].hooks.SessionStart) else . end)
   | (if $add_stop == 1 then .hooks.Stop = ((.hooks.Stop // []) + $fragment[0].hooks.Stop) else . end)' \
  "$settings_file" > "$merge_file" \
  && jq . "$merge_file" >/dev/null 2>&1 \
  && grep -q 'memory-session-start.sh' "$merge_file" \
  && grep -q 'memory-stop-capture.sh' "$merge_file" \
  && mv "$merge_file" "$settings_file"; then
  rm -f "$backup_file"
  printf 'hooks: installed\n'
  exit 0
fi

rm -f "$merge_file"
mv "$backup_file" "$settings_file" 2>/dev/null || true
printf 'Memory hook merge failed; original settings restored.\n' >&2
exit 2

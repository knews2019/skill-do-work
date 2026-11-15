#!/usr/bin/env bash
# REQ-736 GREEN probe: the ai-report judge line in the toolbox help listing shows --page,
# in the same form ai-report.md documents. Expected to FAIL at base (line has no --page).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
help_file="$repo_root/skills/do-work-toolbox/actions/help.md"
action_file="$repo_root/skills/do-work-toolbox/actions/ai-report.md"

help_text="$(cat "$help_file")"
action_text="$(cat "$action_file")"

# The source form the help line must match still exists.
grep -qF 'judge <bundle-dir> [--page <file.html>]' <<<"$action_text" || {
  echo "FAIL: ai-report.md no longer documents judge [--page <file.html>]"; exit 1; }

# Exactly one judge line in help.md, and it carries the --page option.
judge_lines="$(grep -c 'ai-report judge ' <<<"$help_text" || true)"
[ "$judge_lines" = "1" ] || { echo "FAIL: expected 1 judge line in help.md, found $judge_lines"; exit 1; }
grep -qF 'ai-report judge <dir> [--page <file.html>]' <<<"$help_text" || {
  echo "FAIL: help.md judge line does not show [--page <file.html>]"; exit 1; }

# The neighbouring listing lines are untouched.
for neighbour in \
  'ai-report organize [--dry-run] File finished months into month folders; fix links' \
  'ai-report revise <dir|latest>  New revision; the old report is kept'; do
  grep -qF "$neighbour" <<<"$help_text" || { echo "FAIL: neighbour line changed: $neighbour"; exit 1; }
done

echo "PASS: REQ-736 help.md judge line shows --page"

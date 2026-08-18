#!/usr/bin/env bash
# session-start.sh — Injects do-work skill status at session start.
#
# Install as a SessionStart hook in .claude/settings.json (merge with existing hooks):
#
#   {
#     "hooks": {
#       "SessionStart": [
#         {
#           "hooks": [
#             {
#               "type": "command",
#               "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/session-start.sh\""
#             }
#           ]
#         }
#       ]
#     }
#   }
#
# The command is anchored to $CLAUDE_PROJECT_DIR — Claude Code runs hooks from the project
# root, not from this skill directory, so a bare "hooks/session-start.sh" would resolve to
# <project-root>/hooks/... and fail with "No such file or directory". The path also assumes
# the canonical install location .claude/skills/do-work/; if you installed do-work elsewhere,
# adjust it to match. Do NOT "simplify" this back to a relative path — it has regressed before.

set -u

SKILL_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="$(sed -n 's/^\*\*Current version\*\*:[[:space:]]*//p' \
  "$SKILL_ROOT/actions/version.md" 2>/dev/null)"
VERSION="${VERSION:-unknown}"
PENDING="$(find "${CLAUDE_PROJECT_DIR:-.}/do-work/queue" \
  -maxdepth 1 -name 'REQ-*.md' 2>/dev/null | wc -l | tr -d ' ')"
PENDING="${PENDING:-0}"

echo "do-work v${VERSION} loaded. ${PENDING} pending REQ(s). Say 'do-work help' for commands."

# Mechanical housekeeping: reap stale REQ reservation markers (redundant once the
# REQ file exists; abandoned after two days). Guarded so a partial install or a
# cleanup failure can never break the status banner above.
CLEANUP_SCRIPT="$SKILL_ROOT/scripts/cleanup-req-reservations.sh"
if [ -f "$CLEANUP_SCRIPT" ]; then
  CLEANUP_SUMMARY="$(bash "$CLEANUP_SCRIPT" "${CLAUDE_PROJECT_DIR:-.}" 2>/dev/null)" || CLEANUP_SUMMARY=''
  if [ -n "$CLEANUP_SUMMARY" ]; then
    echo "$CLEANUP_SUMMARY"
  fi
fi

# Mechanical housekeeping: repair detectably wrong queue/working timestamps
# (future beyond the 2-minute skew allowance, impossible orderings) before any
# agent or board render reads them — see scripts/repair-req-timestamps.sh for
# the detection and derivation rules. Guarded like the cleanup above so a
# partial install can never break the status banner; `|| true` instead of
# discarding, because on a tripped guard the script's failure lines ARE the
# audit trail and must reach the banner.
REPAIR_SCRIPT="$SKILL_ROOT/scripts/repair-req-timestamps.sh"
if [ -f "$REPAIR_SCRIPT" ]; then
  REPAIR_SUMMARY="$(bash "$REPAIR_SCRIPT" "${CLAUDE_PROJECT_DIR:-.}" 2>/dev/null)" || true
  if [ -n "$REPAIR_SUMMARY" ]; then
    echo "$REPAIR_SUMMARY"
  fi
fi

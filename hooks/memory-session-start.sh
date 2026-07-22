#!/usr/bin/env bash
# memory-session-start.sh — Injects the memory engine's frozen snapshot at session start.
#
# Installed by `do-work install memory-module` (or merge manually from hooks/memory-hooks.json):
#
#   {
#     "hooks": {
#       "SessionStart": [
#         {
#           "hooks": [
#             {
#               "type": "command",
#               "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/memory-session-start.sh\""
#             }
#           ]
#         }
#       ]
#     }
#   }
#
# The command is anchored to $CLAUDE_PROJECT_DIR — Claude Code runs hooks from the project
# root, not from this skill directory, so a bare "hooks/memory-session-start.sh" would resolve
# to <project-root>/hooks/... and fail. The path also assumes the canonical install location
# .claude/skills/do-work/; if you installed do-work elsewhere, adjust it to match. Do NOT
# "simplify" this back to a relative path — the sibling hooks regressed on this before.
#
# Exits 0 silently when no memory/ store exists — the hook must never break a session
# in a repo that hasn't run `do-work install memory-module`.

set -euo pipefail

MEMORY_DIR="${CLAUDE_PROJECT_DIR:-.}/memory"
WORKING_MEMORY_FILE="$MEMORY_DIR/working-memory.md"

[ -f "$WORKING_MEMORY_FILE" ] || exit 0

TODAY_LOG="$MEMORY_DIR/logs/$(date -u +%F).md"

echo "<background-memory>"
echo "Frozen memory snapshot (see .claude/skills/do-work/actions/memory.md). Treat as silent background context: do not greet, recap, or mention it unless it becomes relevant. Writes made this session surface at the NEXT session start."
echo
cat "$WORKING_MEMORY_FILE"
if [ -f "$TODAY_LOG" ]; then
  echo
  echo "## Today's log ($(date -u +%F))"
  cat "$TODAY_LOG"
fi
echo "</background-memory>"

# Best-effort inject ledger line — never fail the hook over instrumentation.
printf '{"ts":"%s","engine":"memory","event":"inject","query":"","hits":0,"source":"hooks/memory-session-start.sh","note":""}\n' \
  "$(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$MEMORY_DIR/usage-ledger.jsonl" 2>/dev/null || true

exit 0

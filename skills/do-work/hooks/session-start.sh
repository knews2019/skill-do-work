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

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

set -euo pipefail

SKILL_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION_FILE="$SKILL_ROOT/actions/version.md"
QUEUE_DIR="${CLAUDE_PROJECT_DIR:-.}/do-work/queue"

# Extract version (format pinned by CLAUDE.md — line starting with `**Current version**:` in actions/version.md)
VERSION=$(grep -m1 '^\*\*Current version\*\*:' "$VERSION_FILE" 2>/dev/null | sed 's/^\*\*Current version\*\*:[[:space:]]*//')
if [ -z "$VERSION" ]; then
  echo "do-work: could not parse version from $VERSION_FILE (expected line starting with '**Current version**:')" >&2
  VERSION="unknown"
fi

# Count pending REQs
PENDING=0
if [ -d "$QUEUE_DIR" ]; then
  PENDING=$(find "$QUEUE_DIR" -maxdepth 1 -name 'REQ-*.md' 2>/dev/null | wc -l | tr -d ' ')
fi

# Check for active pipeline
PIPELINE_MSG=""
PIPELINE_FILE="${CLAUDE_PROJECT_DIR:-.}/do-work/pipeline.json"
if [ -f "$PIPELINE_FILE" ]; then
  if command -v jq &>/dev/null; then
    ACTIVE=$(jq -r '.active // false' "$PIPELINE_FILE" 2>/dev/null || echo "false")
    if [ "$ACTIVE" = "true" ]; then
      PIPELINE_MSG=" | Pipeline active — run 'do-work pipeline' to resume"
    fi
  elif grep -q '"active"[[:space:]]*:[[:space:]]*true' "$PIPELINE_FILE" 2>/dev/null; then
    PIPELINE_MSG=" | Pipeline active — run 'do-work pipeline' to resume"
  fi
fi

echo "do-work v${VERSION} loaded. ${PENDING} pending REQ(s)${PIPELINE_MSG}. Say 'do-work help' for commands."

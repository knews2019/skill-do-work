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
#               "command": "bash hooks/session-start.sh"
#             }
#           ]
#         }
#       ]
#     }
#   }

set -euo pipefail

SKILL_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION_FILE="$SKILL_ROOT/actions/version.md"
QUEUE_DIR="${CLAUDE_PROJECT_DIR:-.}/do-work/queue"

# Extract version
VERSION=$(grep -m1 '^\*\*Current version\*\*:' "$VERSION_FILE" 2>/dev/null | sed 's/.*: //' || echo "unknown")

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
    ACTIVE=$(jq -r '.active // false' "$PIPELINE_FILE" 2>/dev/null)
    if [ "$ACTIVE" = "true" ]; then
      PIPELINE_MSG=" | Pipeline active — run 'do work pipeline' to resume"
    fi
  elif grep -q '"active"[[:space:]]*:[[:space:]]*true' "$PIPELINE_FILE" 2>/dev/null; then
    PIPELINE_MSG=" | Pipeline active — run 'do work pipeline' to resume"
  fi
fi

echo "do-work v${VERSION} loaded. ${PENDING} pending REQ(s)${PIPELINE_MSG}. Say 'do work help' for commands."

#!/usr/bin/env bash
# memory-stop-capture.sh — Appends a deduplicated capture of the session's final exchange
# to the memory engine's daily log when the session stops.
#
# Installed by `do-work install memory-module` (or merge manually from hooks/memory-hooks.json):
#
#   {
#     "hooks": {
#       "Stop": [
#         {
#           "hooks": [
#             {
#               "type": "command",
#               "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/memory-stop-capture.sh\""
#             }
#           ]
#         }
#       ]
#     }
#   }
#
# The command is anchored to $CLAUDE_PROJECT_DIR — Claude Code runs hooks from the project
# root, not from this skill directory. The path assumes the canonical install location
# .claude/skills/do-work/; adjust if you installed elsewhere. Do NOT "simplify" this back
# to a relative path — the sibling hooks regressed on this before.
#
# Contract (actions/memory-reference.md, "Stop-Capture Hash Dedup Spec"):
#   - verbatim tail of the final user+assistant exchange, ~1,500 chars, third-person framed
#   - sha256-prefix (8 hex chars) dedup key in the heading — idempotent across re-fires
#   - ALWAYS exits 0. Capture is never worth blocking a session end; every failure path
#     below falls through to exit 0. This hook must NEVER emit a blocking decision the
#     way pipeline-guard.sh does — _dev/tests/contract-regressions.sh enforces this.

# Deliberately no `set -e`: a parse failure must not abort before the final exit 0.
set -u

INPUT="$(cat 2>/dev/null || true)"

# Never loop on hook-driven continuations
if printf '%s' "$INPUT" | grep -q '"stop_hook_active"[[:space:]]*:[[:space:]]*true' 2>/dev/null; then
  exit 0
fi

MEMORY_DIR="${CLAUDE_PROJECT_DIR:-.}/memory"
[ -d "$MEMORY_DIR/logs" ] || exit 0

# Locate the transcript — prefer jq, fall back to sed
TRANSCRIPT_PATH=""
if command -v jq &>/dev/null; then
  TRANSCRIPT_PATH="$(printf '%s' "$INPUT" | jq -r '.transcript_path // empty' 2>/dev/null || true)"
else
  TRANSCRIPT_PATH="$(printf '%s' "$INPUT" | sed -n 's/.*"transcript_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)"
fi
[ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ] || exit 0

# Extract the final user and assistant message texts from the JSONL transcript.
CAPTURE_TEXT=""
if command -v jq &>/dev/null; then
  extract_last_message_text() {
    # $1 = entry type (user|assistant). Content may be a string or an array of blocks.
    jq -rs --arg entry_type "$1" '
      [.[] | select(.type == $entry_type and (.message.content != null))] | last
      | .message.content
      | if type == "string" then . else ([.[]? | .text? // empty] | join(" ")) end
    ' "$TRANSCRIPT_PATH" 2>/dev/null || true
  }
  LAST_USER_TEXT="$(extract_last_message_text user)"
  LAST_ASSISTANT_TEXT="$(extract_last_message_text assistant)"
  [ "$LAST_USER_TEXT" = "null" ] && LAST_USER_TEXT=""
  [ "$LAST_ASSISTANT_TEXT" = "null" ] && LAST_ASSISTANT_TEXT=""
  CAPTURE_TEXT="$(printf 'User: %s\n\nAgent: %s' "$LAST_USER_TEXT" "$LAST_ASSISTANT_TEXT")"
else
  # Best-effort fallback when jq is absent: grab the raw text fields from the
  # transcript tail. Cruder than the jq path — install jq for clean captures.
  CAPTURE_TEXT="$(tail -c 8000 "$TRANSCRIPT_PATH" | sed -n 's/.*"text"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | tail -6 || true)"
fi

# Nothing meaningful extracted → skip silently
CAPTURE_TEXT="$(printf '%s' "$CAPTURE_TEXT" | head -c 1500)"
[ -n "$(printf '%s' "$CAPTURE_TEXT" | tr -d '[:space:]')" ] || exit 0
case "$CAPTURE_TEXT" in
  "User: "*"Agent: ") exit 0 ;;   # both messages came back empty
esac

# Hash for idempotency (sha256sum, shasum fallback)
if command -v sha256sum &>/dev/null; then
  CAPTURE_HASH="$(printf '%s' "$CAPTURE_TEXT" | sha256sum | cut -c1-8)"
elif command -v shasum &>/dev/null; then
  CAPTURE_HASH="$(printf '%s' "$CAPTURE_TEXT" | shasum -a 256 | cut -c1-8)"
else
  exit 0   # no hash tool → cannot dedup safely → skip capture entirely
fi

TODAY_LOG="$MEMORY_DIR/logs/$(date -u +%F).md"
if [ -f "$TODAY_LOG" ] && grep -q "session capture $CAPTURE_HASH" "$TODAY_LOG" 2>/dev/null; then
  exit 0   # already captured
fi

{
  printf '\n## %s UTC session capture %s\n\n' "$(date -u +%H:%M)" "$CAPTURE_HASH"
  printf 'Session capture — final exchange between the user and the agent:\n\n'
  printf '%s\n' "$CAPTURE_TEXT"
} >> "$TODAY_LOG" 2>/dev/null || exit 0

# Best-effort capture ledger line
printf '{"ts":"%s","engine":"memory","event":"capture","query":"","hits":0,"source":"hooks/memory-stop-capture.sh","note":"%s"}\n' \
  "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$CAPTURE_HASH" >> "$MEMORY_DIR/usage-ledger.jsonl" 2>/dev/null || true

exit 0

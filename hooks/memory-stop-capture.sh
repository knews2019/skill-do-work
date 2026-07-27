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

# Capture size budget. CAPTURE_TEXT_BUDGET is what's left for the two message texts
# once the "User: " / "\n\nAgent: " framing (15 bytes) is accounted for, so a fully
# budgeted capture composes to exactly CAPTURE_BUDGET_BYTES. CAPTURE_SIDE_FLOOR is the
# amount each side is guaranteed before the other can claim any slack.
CAPTURE_BUDGET_BYTES=1500
CAPTURE_TEXT_BUDGET=$(( CAPTURE_BUDGET_BYTES - 15 ))
CAPTURE_SIDE_FLOOR=$(( CAPTURE_TEXT_BUDGET / 2 ))

# Extract the final user and assistant message texts from the JSONL transcript.
CAPTURE_TEXT=""
if command -v jq &>/dev/null; then
  extract_last_message_text() {
    # $1 = entry type (user|assistant). Content may be a string or an array of blocks.
    #
    # Claude Code records TOOL RESULTS as `type: "user"` entries whose content is an
    # array of `tool_result` blocks — no top-level `.text`. Taking `last` over all
    # `user` entries therefore lands on a tool result for any session whose final turn
    # used a tool (measured: 7 of 8 recent transcripts on the author's machine), and
    # the capture stores an empty `User:` side. The assistant side has the mirror
    # problem: a turn ending in a `tool_use` block yields no text either.
    #
    # So: pull text ONLY from blocks explicitly typed `text`, skip `isMeta` entries
    # (caveats and slash-command wrappers are not the human's prompt), drop every
    # entry whose extracted text is blank, and take the last one that survives.
    jq -rs --arg entry_type "$1" '
      [ .[]
        | select(.type == $entry_type and (.message.content != null) and ((.isMeta // false) | not))
        | .message.content
        | if type == "string" then .
          else ([.[]? | select(.type? == "text") | .text? // empty] | join(" "))
          end
        | select(type == "string" and test("\\S"))
      ] | last // ""
    ' "$TRANSCRIPT_PATH" 2>/dev/null || true
  }
  LAST_USER_TEXT="$(extract_last_message_text user)"
  LAST_ASSISTANT_TEXT="$(extract_last_message_text assistant)"
  [ "$LAST_USER_TEXT" = "null" ] && LAST_USER_TEXT=""
  [ "$LAST_ASSISTANT_TEXT" = "null" ] && LAST_ASSISTANT_TEXT=""

  # Trim the two sides INDEPENDENTLY, never the composed string. Truncating
  # "User: … Agent: …" as one blob means a long prompt eats the entire assistant
  # reply — and the reply is the half holding the decisions and outcome this
  # capture exists to preserve. Each side is guaranteed at least half the budget;
  # whichever side comes in under its half yields the slack to the other, so a
  # short prompt still lets a long answer through.
  byte_length_of() { printf '%s' "$1" | wc -c | tr -d '[:space:]'; }
  truncate_to_bytes() {
    # $1 = text, $2 = byte budget. Marks the cut so a reader knows text is missing.
    if [ "$(byte_length_of "$1")" -le "$2" ]; then printf '%s' "$1"; return; fi
    printf '%s [truncated]' "$(printf '%s' "$1" | head -c "$(( $2 - 12 ))")"
  }
  user_text_bytes="$(byte_length_of "$LAST_USER_TEXT")"
  assistant_text_bytes="$(byte_length_of "$LAST_ASSISTANT_TEXT")"
  if [ "$(( user_text_bytes + assistant_text_bytes ))" -gt "$CAPTURE_TEXT_BUDGET" ]; then
    if [ "$user_text_bytes" -le "$CAPTURE_SIDE_FLOOR" ]; then
      LAST_ASSISTANT_TEXT="$(truncate_to_bytes "$LAST_ASSISTANT_TEXT" "$(( CAPTURE_TEXT_BUDGET - user_text_bytes ))")"
    elif [ "$assistant_text_bytes" -le "$CAPTURE_SIDE_FLOOR" ]; then
      LAST_USER_TEXT="$(truncate_to_bytes "$LAST_USER_TEXT" "$(( CAPTURE_TEXT_BUDGET - assistant_text_bytes ))")"
    else
      LAST_USER_TEXT="$(truncate_to_bytes "$LAST_USER_TEXT" "$CAPTURE_SIDE_FLOOR")"
      LAST_ASSISTANT_TEXT="$(truncate_to_bytes "$LAST_ASSISTANT_TEXT" "$CAPTURE_SIDE_FLOOR")"
    fi
  fi
  CAPTURE_TEXT="$(printf 'User: %s\n\nAgent: %s' "$LAST_USER_TEXT" "$LAST_ASSISTANT_TEXT")"
else
  # Best-effort fallback when jq is absent: grab the raw text fields from the
  # transcript tail. Cruder than the jq path — install jq for clean captures.
  CAPTURE_TEXT="$(tail -c 8000 "$TRANSCRIPT_PATH" | sed -n 's/.*"text"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | tail -6 || true)"
fi

# Nothing meaningful extracted → skip silently. This blanket cap is a backstop for the
# no-jq fallback above, which has no separate user/assistant sides to budget; the jq
# path already composes to <= CAPTURE_BUDGET_BYTES, so it is a no-op there.
CAPTURE_TEXT="$(printf '%s' "$CAPTURE_TEXT" | head -c "$CAPTURE_BUDGET_BYTES")"
[ -n "$(printf '%s' "$CAPTURE_TEXT" | tr -d '[:space:]')" ] || exit 0
case "$CAPTURE_TEXT" in
  "User: "*"Agent: ") exit 0 ;;   # both messages came back empty
esac

# Redact credential-shaped content BEFORE anything is persisted or hashed — memory
# files are committed plaintext (actions/memory-reference.md, "Capture Redaction
# Spec"). The pattern list is illustrative, not exhaustive: it catches common token
# shapes, and `memory remember` curation remains the real gate. A private-key block
# spans lines, so line-based redaction can't be trusted — drop the capture entirely.
case "$CAPTURE_TEXT" in
  *"PRIVATE KEY-----"*) exit 0 ;;
esac
REDACTED_TEXT="$(printf '%s' "$CAPTURE_TEXT" | sed -E \
  -e 's/(gh[pousr]|github_pat)_[A-Za-z0-9_]{16,}/[REDACTED]/g' \
  -e 's/sk-[A-Za-z0-9_-]{16,}/[REDACTED]/g' \
  -e 's/AKIA[0-9A-Z]{16}/[REDACTED]/g' \
  -e 's/xox[baprs]-[A-Za-z0-9-]{10,}/[REDACTED]/g' \
  -e 's/eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{5,}\.[A-Za-z0-9_-]{5,}/[REDACTED]/g' \
  -e 's/[Bb]earer[[:space:]]+[A-Za-z0-9._~+\/=-]{16,}/Bearer [REDACTED]/g' \
  -e 's/(([Pp]assword|[Pp]asswd|[Ss]ecret|[Tt]oken|[Aa][Pp][Ii][_-]?[Kk]ey)["'"'"']?[[:space:]]*[:=][[:space:]]*)["'"'"']?[^[:space:]"'"'"']{6,}/\1[REDACTED]/g' \
  2>/dev/null)" || exit 0   # redaction failed → never persist the unredacted text
CAPTURE_TEXT="$REDACTED_TEXT"
[ -n "$(printf '%s' "$CAPTURE_TEXT" | tr -d '[:space:]')" ] || exit 0

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

# Blockquote every body line before it lands in the log. The capture is verbatim
# transcript text, and assistant responses routinely contain Markdown headings — an
# unquoted `## Findings` inside a capture reads as a new log section, which let
# hooks/memory-session-start.sh treat the rest of the capture as curated content and
# inject it at the next session start, bypassing the prompt-injection guard the
# exclusion exists to enforce (actions/memory-reference.md → Daily-Log Entry Conventions).
# `> ` cannot begin a heading in any Markdown parser and is honest about what the
# body is: quoted transcript, not authored log prose. The reader-side awk is hardened
# independently, because logs already on disk were written without this prefix.
QUOTED_CAPTURE_TEXT="$(printf '%s\n' "$CAPTURE_TEXT" | sed 's/^/> /')" || exit 0

{
  printf '\n## %s UTC session capture %s\n\n' "$(date -u +%H:%M)" "$CAPTURE_HASH"
  printf 'Session capture — final exchange between the user and the agent:\n\n'
  printf '%s\n' "$QUOTED_CAPTURE_TEXT"
} >> "$TODAY_LOG" 2>/dev/null || exit 0

# Best-effort capture ledger line
printf '{"ts":"%s","engine":"memory","event":"capture","query":"","hits":0,"source":"hooks/memory-stop-capture.sh","note":"%s"}\n' \
  "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$CAPTURE_HASH" >> "$MEMORY_DIR/usage-ledger.jsonl" 2>/dev/null || true

exit 0

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
  # Inject only the CURATED lines of today's log. Raw `## … session capture …`
  # sections are verbatim third-party/transcript text captured by the Stop hook —
  # injecting them here would put unvetted content into context before any
  # prompt-injection guard can load. They stay reachable only via `memory recall`,
  # which loads crew-members/prompt-injection.md first (actions/memory.md).
  #
  # Heading grammar alone cannot end a capture section, because raw capture text can
  # contain a line that looks exactly like one — `## 12:34 UTC note` is trivially
  # spoofable by anything that reaches the transcript. The boundary must be something
  # body text provably cannot produce, so the format decides:
  #
  #   quoted (current) — the section's first non-blank line is CAPTURE_BODY_SENTINEL and
  #     every body line is `> `-prefixed. A boundary is then a heading-grammar line that
  #     is NOT quoted, which no body line can be. Curated notes after a capture still
  #     inject normally.
  #   legacy (pre-0.139.4, already on disk) — bodies are unquoted, so NO boundary can be
  #     trusted. Suppress to end-of-file: once capture_format is `legacy`, no later heading
  #     can be a boundary — not even a new quoted capture's — so suppression runs to the end
  #     of that day's log. It deliberately also hides curated entries written after a legacy
  #     capture that day; a bounded, one-day loss of convenience is the right trade against
  #     injecting raw transcript text, and it self-clears at the next UTC day's fresh log.
  #
  # Keep this string byte-identical to CAPTURE_BODY_SENTINEL in hooks/memory-stop-capture.sh.
  CAPTURE_BODY_SENTINEL='<!-- do-work:capture-body quoted -->'
  CURATED_LOG_LINES="$(awk -v sentinel="$CAPTURE_BODY_SENTINEL" '
    {
      is_heading = ($0 ~ /^## [0-9][0-9]:[0-9][0-9] UTC /)
      is_quoted  = ($0 ~ /^>/)
      is_boundary = 0
      if (is_heading) {
        if (!in_capture_section) is_boundary = 1
        else if (capture_format == "quoted" && !is_quoted) is_boundary = 1
      }
      if (is_boundary) {
        in_capture_section = ($0 ~ /session capture/)
        capture_format = ""
      } else if (in_capture_section && capture_format == "" && $0 ~ /[^[:space:]]/) {
        if ($0 == sentinel) capture_format = "quoted"; else capture_format = "legacy"
      }
      if (!in_capture_section) print
    }
  ' "$TODAY_LOG" 2>/dev/null || true)"
  if [ -n "$(printf '%s' "$CURATED_LOG_LINES" | tr -d '[:space:]')" ]; then
    echo
    echo "## Today's log ($(date -u +%F)) — curated entries only; raw session captures load via \`do-work memory recall\`"
    printf '%s\n' "$CURATED_LOG_LINES"
  fi
fi
echo "</background-memory>"

# Best-effort inject ledger line — never fail the hook over instrumentation.
printf '{"ts":"%s","engine":"memory","event":"inject","query":"","hits":0,"source":"hooks/memory-session-start.sh","note":""}\n' \
  "$(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$MEMORY_DIR/usage-ledger.jsonl" 2>/dev/null || true

exit 0

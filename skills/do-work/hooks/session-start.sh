#!/usr/bin/env bash
# session-start.sh — retained SessionStart path; Go owns the hook protocol.
set -u

SKILL_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CLI_LAUNCHER="$SKILL_ROOT/tools/do-work-cli.sh"

if [ ! -f "$CLI_LAUNCHER" ]; then
  printf 'do-work SessionStart failed: canonical launcher is missing at %s; reinstall the suite with Go 1.26.1 or newer.\n' "$CLI_LAUNCHER" >&2
  exit 2
fi

exec bash "$CLI_LAUNCHER" \
  --repo-root "${CLAUDE_PROJECT_DIR:-.}" \
  session-start --skill-root "$SKILL_ROOT"

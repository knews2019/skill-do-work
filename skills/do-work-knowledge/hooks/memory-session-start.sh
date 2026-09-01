#!/usr/bin/env bash
# memory-session-start.sh — retained SessionStart path; the sibling core CLI owns behavior.
set -u

HOOK_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_LAUNCHER="$HOOK_DIR/../../do-work/tools/do-work-cli.sh"

if [ ! -f "$CLI_LAUNCHER" ]; then
  printf 'memory SessionStart failed: canonical do-work launcher is missing at %s; reinstall the suite with Go 1.26.1 or newer.\n' "$CLI_LAUNCHER" >&2
  exit 2
fi

exec bash "$CLI_LAUNCHER" \
  --repo-root "${CLAUDE_PROJECT_DIR:-.}" \
  memory-session-start

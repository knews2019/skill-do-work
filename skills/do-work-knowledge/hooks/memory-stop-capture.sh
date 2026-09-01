#!/usr/bin/env bash
# memory-stop-capture.sh — retained nonblocking Stop path; the sibling core CLI owns behavior.
set -u

HOOK_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_LAUNCHER="$HOOK_DIR/../../do-work/tools/do-work-cli.sh"

if [ ! -f "$CLI_LAUNCHER" ]; then
  printf 'memory Stop capture skipped: canonical do-work launcher is missing at %s; reinstall the suite with Go 1.26.1 or newer.\n' "$CLI_LAUNCHER" >&2
  exit 0
fi

bash "$CLI_LAUNCHER" \
  --repo-root "${CLAUDE_PROJECT_DIR:-.}" \
  memory-stop-capture
launcher_status=$?
if [ "$launcher_status" -ne 0 ]; then
  printf 'memory Stop capture skipped: canonical do-work launcher failed with status %s; verify Go 1.26.1 or newer and reinstall the suite.\n' "$launcher_status" >&2
fi
exit 0

#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
# memory-stop-capture.sh — retained nonblocking Stop path; the sibling core CLI owns behavior.
set -u

HOOK_DIR="$(cd "$(dirname "$0")" && pwd)"
CLI_LAUNCHER="$HOOK_DIR/../../do-work/tools/do-work-cli.sh"

if [ ! -f "$CLI_LAUNCHER" ]; then
  printf 'memory Stop capture skipped: canonical do-work launcher is missing at %s; reinstall the suite with Go 1.25.0 or newer.\n' "$CLI_LAUNCHER" >&2
  exit 0
fi

launcher_error="$(bash "$CLI_LAUNCHER" \
  --repo-root "${CLAUDE_PROJECT_DIR:-.}" \
  memory-stop-capture 2>&1 >/dev/null)"
launcher_status=$?
if [ "$launcher_status" -ne 0 ]; then
  printf 'memory Stop capture skipped: canonical do-work launcher failed with status %s (%s); verify Go 1.25.0 or newer and reinstall the suite.\n' \
    "$launcher_status" "${launcher_error:-no diagnostic}" >&2
elif [ -n "$launcher_error" ]; then
  printf '%s\n' "$launcher_error" >&2
fi
exit 0

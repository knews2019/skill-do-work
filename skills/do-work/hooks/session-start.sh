#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
# session-start.sh — retained SessionStart path; Go owns the hook protocol.
set -u

SKILL_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CLI_LAUNCHER="$SKILL_ROOT/tools/do-work-cli.sh"

if [ ! -f "$CLI_LAUNCHER" ]; then
  printf 'do-work SessionStart failed: canonical launcher is missing at %s; reinstall the suite with Go 1.25.0 or newer.\n' "$CLI_LAUNCHER" >&2
  exit 2
fi

REPO_ROOT="$(git -C "${CLAUDE_PROJECT_DIR:-.}" rev-parse --show-toplevel 2>/dev/null || echo "${CLAUDE_PROJECT_DIR:-.}")"
if [ -f "$REPO_ROOT/_dev/tests/gate-runner.sh" ]; then
  TMP_DIR="${TMPDIR:-/tmp}"
  LOG_ROOT="${TMP_DIR%/}/do-work-gate-runs"
  mkdir -p "$LOG_ROOT"
  REPO_ID="$(printf '%s' "$REPO_ROOT" | shasum -a 256 2>/dev/null | cut -c1-12 || echo 'default')"
  PID_FILE="$LOG_ROOT/runner-$REPO_ID.pid"
  RUNNER_RUNNING=no
  if [ -f "$PID_FILE" ]; then
    EXISTING_PID="$(cat "$PID_FILE" 2>/dev/null || true)"
    if [ -n "$EXISTING_PID" ] && kill -0 "$EXISTING_PID" 2>/dev/null; then
      RUNNER_RUNNING=yes
    fi
  fi
  if [ "$RUNNER_RUNNING" = 'no' ]; then
    printf 'gate runner logging to %s\n' "$LOG_ROOT"
    nohup bash "$REPO_ROOT/_dev/tests/gate-runner.sh" >/dev/null 2>&1 &
    printf '%s\n' "$!" > "$PID_FILE"
  fi
fi

exec bash "$CLI_LAUNCHER" \
  --repo-root "${CLAUDE_PROJECT_DIR:-.}" \
  session-start --skill-root "$SKILL_ROOT"


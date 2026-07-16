#!/usr/bin/env bash
# preflight.sh — mechanical form of actions/work.md Step 5.75 (Routes B and C).
# Environment sanity check before the builder starts coding. Every finding is a
# WARNING, never a blocker — the exit code is always 0 so the work loop continues.
#
# Usage: tools/checks/preflight.sh [test-command ...]
#   e.g.  tools/checks/preflight.sh npm test
#   With no test command, the baseline check is skipped (generic detection of the
#   right command is the orchestrator's judgment, guided by the REQ's prime files).
#
# Output: "WARN: ..." / "OK: ..." lines on stdout.
# Side effect: when a test command is given, writes do-work/working/baseline.json
# recording the command and any failing output tail, so Step 6.5 can separate
# pre-existing failures from new regressions.
set -uo pipefail

# --- 1. Git clean (uncommitted changes outside do-work/ get swept into commits) ---
if git rev-parse --git-dir >/dev/null 2>&1; then
  dirty_files="$(git status --porcelain --untracked-files=all | awk '{print $2}' | grep -v '^do-work/' || true)"
  if [ -n "$dirty_files" ]; then
    echo "WARN: uncommitted changes detected — the commit step may stage unrelated files:"
    printf '  %s\n' $dirty_files
  else
    echo "OK: working tree clean outside do-work/"
  fi
else
  echo "WARN: not a git repository — no clean-tree or diff-based checks available"
fi

# --- 2. Test baseline (pre-existing failures must not be blamed on the builder) ---
if [ "$#" -gt 0 ]; then
  baseline_dir="do-work/working"
  mkdir -p "$baseline_dir"
  baseline_output="$("$@" 2>&1)" && baseline_status=0 || baseline_status=$?
  if [ "$baseline_status" -eq 0 ]; then
    echo "OK: test baseline passing ($*)"
    # A stale failures file from an earlier failing preflight would make Step 6.5
    # misclassify a new regression as pre-existing — clear it on a passing baseline.
    rm -f "$baseline_dir/baseline-failures.txt"
  else
    echo "WARN: baseline tests failing BEFORE any changes — builder is not to blame for these:"
    printf '%s\n' "$baseline_output" | tail -n 20 | sed 's/^/  /'
  fi
  python3 - "$*" "$baseline_status" <<'PYEOF' 2>/dev/null || \
    printf '{"test_command": "%s", "exit_status": %s}\n' "$*" "$baseline_status" > "$baseline_dir/baseline.json"
import json, sys
baseline_record = {"test_command": sys.argv[1], "exit_status": int(sys.argv[2])}
with open("do-work/working/baseline.json", "w") as handle:
    json.dump(baseline_record, handle, indent=2)
PYEOF
  if [ "$baseline_status" -ne 0 ]; then
    printf '%s\n' "$baseline_output" > "$baseline_dir/baseline-failures.txt"
    echo "OK: baseline recorded in $baseline_dir/baseline.json + baseline-failures.txt (Step 6.5 compares against this)"
  fi
else
  echo "SKIP: no test command supplied — baseline check skipped"
fi

# --- 3. Dependencies present ---
if [ -f package.json ] && [ ! -d node_modules ]; then
  echo "WARN: package.json exists but node_modules/ does not — dependencies may not be installed"
fi
if [ -f requirements.txt ] && [ -z "${VIRTUAL_ENV:-}" ] && ! python3 -c "import sys; sys.exit(0 if sys.prefix != sys.base_prefix else 1)" 2>/dev/null; then
  echo "WARN: requirements.txt exists but no active virtualenv detected"
fi

exit 0

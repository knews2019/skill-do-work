#!/usr/bin/env bash
# preflight.sh — mechanical form of actions/work.md Step 5.75 (Routes B and C).
# Environment sanity check before the builder starts coding. Every finding is a
# WARNING, never a blocker — the exit code is always 0 so the work loop continues.
#
# Usage: tools/checks/preflight.sh [test-command ...]
#   Separate words:  tools/checks/preflight.sh npm test
#   One quoted string (run via `sh -c`, so pipes/&&/cd work):
#                    tools/checks/preflight.sh "npx jest --silent"
#                    tools/checks/preflight.sh "cd app && npm test"
#   With no test command, the baseline check is skipped (generic detection of the
#   right command is the orchestrator's judgment, guided by the REQ's prime files).
#
#   Run from the project root (the directory containing do-work/) — the baseline
#   artifacts are written to do-work/working/ relative to the current directory,
#   and Step 6.5 reads them from the root. Use the quoted form above when the test
#   command itself has to run somewhere else.
#
# Output: "WARN: ..." / "OK: ..." lines on stdout.
# Side effect: when a test command is given, writes do-work/working/baseline.json
# recording the command and any failing output tail, so Step 6.5 can separate
# pre-existing failures from new regressions. A command that could not be launched
# at all (exit 126/127) records "launched": false and writes NO failures file —
# a fictional red baseline would let Step 6.5 excuse a real regression as
# pre-existing, which is the exact mis-attribution this check exists to prevent.
set -uo pipefail

# --- 1. Repository state (pre-existing changes can contaminate diff-based evidence) ---
if git rev-parse --git-dir >/dev/null 2>&1; then
  dirty_files=()
  dirty_file_count=0
  status_check_failed=0
  append_dirty_path() {
    local candidate_path="$1" existing_path
    case "$candidate_path" in
      do-work|do-work/*) return ;;
    esac
    if [ "$dirty_file_count" -gt 0 ]; then
      for existing_path in "${dirty_files[@]}"; do
        [ "$existing_path" != "$candidate_path" ] || return
      done
    fi
    dirty_files[dirty_file_count]="$candidate_path"
    dirty_file_count=$((dirty_file_count + 1))
  }

  status_output_file="$(mktemp "${TMPDIR:-/tmp}/do-work-preflight-status.XXXXXX" 2>/dev/null || true)"
  if [ -z "$status_output_file" ]; then
    echo "WARN: could not allocate temporary storage for the clean-tree check"
    status_check_failed=1
  elif ! git -c status.renames=copies status --porcelain=v1 --untracked-files=all -z \
    > "$status_output_file"; then
    echo "WARN: git status could not read the working tree"
    status_check_failed=1
    rm -f "$status_output_file"
  else
    while IFS= read -r -d '' status_record; do
      index_status="${status_record:0:1}"
      worktree_status="${status_record:1:1}"
      changed_path="${status_record:3}"
      append_dirty_path "$changed_path"
      case "$index_status$worktree_status" in
        R*|*R|C*|*C)
          rename_origin_path=''
          IFS= read -r -d '' rename_origin_path || true
          [ -z "$rename_origin_path" ] || append_dirty_path "$rename_origin_path"
          ;;
      esac
    done < "$status_output_file"
    rm -f "$status_output_file"
  fi

  if [ "$dirty_file_count" -gt 0 ]; then
    echo "WARN: pre-existing uncommitted changes detected outside do-work/ — preserve them and, unless they prevent the active REQ, exclude them from its staging; account for them in repository-wide qualification/review evidence:"
    printf '  %s\n' "${dirty_files[@]}"
  elif [ "$status_check_failed" -eq 0 ]; then
    echo "OK: working tree clean outside do-work/"
  fi
else
  echo "WARN: not a git repository — no clean-tree or diff-based checks available"
fi

# --- 2. Test baseline (pre-existing failures must not be blamed on the builder) ---
if [ "$#" -gt 0 ]; then
  baseline_dir="do-work/working"
  mkdir -p "$baseline_dir"
  # One argument containing whitespace is a whole command line, not an executable
  # name — `"$@"` would look for a file literally called "npx jest --silent" and
  # exit 127. Hand that form to `sh -c` so the documented quoted usage works.
  if [ "$#" -eq 1 ] && [ "$1" != "${1#*[[:space:]]}" ]; then
    baseline_output="$(sh -c "$1" 2>&1)" && baseline_status=0 || baseline_status=$?
  else
    baseline_output="$("$@" 2>&1)" && baseline_status=0 || baseline_status=$?
  fi
  # 127 = command not found, 126 = found but not executable. Neither is a test
  # result: the suite never ran, so there is no baseline to record. Recording one
  # anyway is worse than recording nothing — Step 6.5 would read the fiction as
  # "already failing before we started" and wave a real regression through.
  if [ "$baseline_status" -eq 126 ] || [ "$baseline_status" -eq 127 ]; then
    baseline_launched=false
  else
    baseline_launched=true
  fi
  if [ "$baseline_launched" = false ]; then
    echo "WARN: could not run the test command — no baseline recorded ($*):"
    printf '%s\n' "$baseline_output" | tail -n 20 | sed 's/^/  /'
    echo "  (pass the command as separate words, or as one quoted string — see the usage header)"
    # Same reasoning as the passing branch: a stale failures file from an earlier
    # run would be compared against as if it described this session.
    rm -f "$baseline_dir/baseline-failures.txt"
  elif [ "$baseline_status" -eq 0 ]; then
    echo "OK: test baseline passing ($*)"
    # A stale failures file from an earlier failing preflight would make Step 6.5
    # misclassify a new regression as pre-existing — clear it on a passing baseline.
    rm -f "$baseline_dir/baseline-failures.txt"
  else
    echo "WARN: baseline tests failing BEFORE any changes — builder is not to blame for these:"
    printf '%s\n' "$baseline_output" | tail -n 20 | sed 's/^/  /'
  fi
  python3 - "$*" "$baseline_status" "$baseline_launched" "$baseline_dir/baseline.json" <<'PYEOF' 2>/dev/null || \
    printf '{"test_command": "%s", "exit_status": %s, "launched": %s}\n' "$*" "$baseline_status" "$baseline_launched" > "$baseline_dir/baseline.json"
import json, sys
baseline_record = {
    "test_command": sys.argv[1],
    "exit_status": int(sys.argv[2]),
    "launched": sys.argv[3] == "true",
}
with open(sys.argv[4], "w") as handle:
    json.dump(baseline_record, handle, indent=2)
PYEOF
  if [ "$baseline_launched" = true ] && [ "$baseline_status" -ne 0 ]; then
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

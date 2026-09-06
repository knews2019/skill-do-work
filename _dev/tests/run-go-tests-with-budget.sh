#!/usr/bin/env bash
# Run one Go test selection and reject any source test file at or above the budget.
set -euo pipefail

if [ "$#" -lt 1 ]; then
  printf 'usage: %s <module-directory> [go-test-argument ...]\n' "$0" >&2
  exit 2
fi

module_directory="$1"
shift
test_budget_seconds="${DO_WORK_TEST_FILE_BUDGET_SECONDS:-30}"
excluded_test_prefixes="${DO_WORK_GO_TEST_EXCLUDE_PREFIXES:-}"
result_file="$(mktemp "${TMPDIR:-/tmp}/do-work-go-test-budget.XXXXXX")"
cleanup_result() {
  rm -f -- "$result_file"
}
trap cleanup_result EXIT

escape_go_test_regex() {
  local input="$1"
  local output=""
  local i char
  for ((i = 0; i < ${#input}; i++)); do
    char="${input:i:1}"
    case "$char" in
      "\\"|"."|"*"|"+"|"?"|"^"|"$"|"("|")"|"["|"]"|"{"|"}"|"|")
        output+="\\$char"
        ;;
      *)
        output+="$char"
        ;;
    esac
  done
  printf '%s' "$output"
}

started_at="$(date +%s)"
test_status=0
if [ -n "$excluded_test_prefixes" ]; then
  skip_pattern=""
  old_ifs="$IFS"
  IFS=','
  for prefix in $excluded_test_prefixes; do
    prefix="${prefix#"${prefix%%[![:space:]]*}"}"
    prefix="${prefix%"${prefix##*[![:space:]]}"}"
    [ -z "$prefix" ] && continue
    escaped="$(escape_go_test_regex "$prefix")"
    if [ -z "$skip_pattern" ]; then
      skip_pattern="$escaped"
    else
      skip_pattern="$skip_pattern|$escaped"
    fi
  done
  IFS="$old_ifs"
  if [ -n "$skip_pattern" ]; then
    set -- -skip "^($skip_pattern)" "$@"
  fi
fi
(
  cd "$module_directory"
  go test -json -count=1 "$@"
) > "$result_file" || test_status=$?
elapsed_seconds=$(( $(date +%s) - started_at ))

if [ "$test_status" -ne 0 ]; then
  python3 - "$result_file" <<'PY'
import json
import pathlib
import re
import sys

for line in pathlib.Path(sys.argv[1]).read_text(encoding="utf-8").splitlines():
    try:
        event = json.loads(line)
    except json.JSONDecodeError:
        print(line)
        continue
    output = event.get("Output")
    if output:
        print(output, end="")
PY
  printf 'go-test budget: FAIL module=%s wall=%ss exit=%s\n' \
    "$module_directory" "$elapsed_seconds" "$test_status" >&2
fi

python3 - "$result_file" "$module_directory" "$test_budget_seconds" "$elapsed_seconds" "$test_status" "$excluded_test_prefixes" <<'PY'
import json
import os
import pathlib
import re
import sys

result_path = pathlib.Path(sys.argv[1])
module_directory = sys.argv[2]
budget_seconds = float(sys.argv[3])
wall_seconds = int(sys.argv[4])
test_status = int(sys.argv[5])
excluded_test_prefixes = sys.argv[6] if len(sys.argv) > 6 else ""
enforce_budget = os.environ.get("DO_WORK_TEST_ENFORCE_BUDGET", "yes") == "yes" and test_status == 0
durations = []
test_file_by_name = {}
for test_file in pathlib.Path(module_directory).rglob("*_test.go"):
    source_text = test_file.read_text(encoding="utf-8")
    for test_name in re.findall(r"^func\s+(Test\w+)\s*\(", source_text, re.MULTILINE):
        test_file_by_name[test_name] = test_file.relative_to(module_directory).as_posix()

for line in result_path.read_text(encoding="utf-8").splitlines():
    try:
        event = json.loads(line)
    except json.JSONDecodeError:
        continue
    test_name = event.get("Test")
    if event.get("Action") in {"pass", "fail"} and test_name and "/" not in test_name:
        durations.append((float(event.get("Elapsed", 0)), event.get("Package", ""), test_name))

if excluded_test_prefixes and not durations and test_status == 0:
    print("no fast Go tests remain after applying the heavy prefixes", file=sys.stderr)
    sys.exit(1)

file_durations = {}
for elapsed, _, test_name in durations:
    test_file = test_file_by_name.get(test_name, "<unknown test file>")
    file_durations[test_file] = file_durations.get(test_file, 0.0) + elapsed
duration_log_path = os.environ.get("DO_WORK_TEST_DURATION_LOG")
duration_run_id = os.environ.get("DO_WORK_TEST_RUN_ID")
repo_root_env = os.environ.get("DO_WORK_TEST_REPO_ROOT")
if duration_log_path and duration_run_id and repo_root_env:
    repository_root = pathlib.Path(repo_root_env)
    module_path = pathlib.Path(module_directory)
    module_relative = module_path.relative_to(repository_root)
    concurrent_count = os.environ.get("DO_WORK_TEST_OTHER_GATE_PROCESSES", "0")
    with pathlib.Path(duration_log_path).open("a", encoding="utf-8") as duration_log:
        for test_file, elapsed in sorted(file_durations.items()):
            logged_file = (module_relative / test_file).as_posix()
            duration_log.write(
                f"{duration_run_id}\t{logged_file}\t{elapsed:.2f}\t{concurrent_count}\n"
            )
over_budget = [entry for entry in file_durations.items() if entry[1] >= budget_seconds] if enforce_budget else []
slowest_file, slowest_duration = max(file_durations.items(), key=lambda entry: entry[1], default=("none", 0.0))
print(
    f"go-test budget: module={module_directory} wall={wall_seconds}s "
    f"tests={len(durations)} slowest-file={slowest_file}:{slowest_duration:.2f}s "
    f"limit={'<' + format(budget_seconds, 'g') + 's' if enforce_budget else 'none (heavy)'}"
)
for test_file, elapsed in sorted(over_budget, key=lambda entry: entry[1], reverse=True):
    print(
        f"FAIL: {test_file} accumulated {elapsed:.2f}s; "
        f"each test file must finish under {budget_seconds:g}s",
        file=sys.stderr,
    )
if over_budget:
    raise SystemExit(1)
PY

if [ "$test_status" -ne 0 ]; then
  exit "$test_status"
fi

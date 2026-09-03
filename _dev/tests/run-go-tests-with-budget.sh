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

started_at="$(date +%s)"
test_status=0
if [ -n "$excluded_test_prefixes" ]; then
  test_run_pattern="$({
    cd "$module_directory"
    go test -list '^Test'
  } | python3 -c '
import re
import sys

prefixes = tuple(filter(None, sys.argv[1].split(",")))
names = [line.strip() for line in sys.stdin if line.startswith("Test")]
selected = [name for name in names if not name.startswith(prefixes)]
if not selected:
    raise SystemExit("no fast Go tests remain after applying the heavy prefixes")
print("^(" + "|".join(re.escape(name) for name in selected) + ")$")
' "$excluded_test_prefixes")"
  set -- -run "$test_run_pattern" "$@"
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
  exit "$test_status"
fi

python3 - "$result_file" "$module_directory" "$test_budget_seconds" "$elapsed_seconds" <<'PY'
import json
import pathlib
import re
import sys

result_path = pathlib.Path(sys.argv[1])
module_directory = sys.argv[2]
budget_seconds = float(sys.argv[3])
wall_seconds = int(sys.argv[4])
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
    if event.get("Action") == "pass" and test_name and "/" not in test_name:
        durations.append((float(event.get("Elapsed", 0)), event.get("Package", ""), test_name))

file_durations = {}
for elapsed, _, test_name in durations:
    test_file = test_file_by_name.get(test_name, "<unknown test file>")
    file_durations[test_file] = file_durations.get(test_file, 0.0) + elapsed
over_budget = [entry for entry in file_durations.items() if entry[1] >= budget_seconds]
slowest_file, slowest_duration = max(file_durations.items(), key=lambda entry: entry[1], default=("none", 0.0))
print(
    f"go-test budget: module={module_directory} wall={wall_seconds}s "
    f"tests={len(durations)} slowest-file={slowest_file}:{slowest_duration:.2f}s "
    f"limit=<{budget_seconds:g}s"
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

#!/usr/bin/env bash
# A failed Go test must still produce file-attributed timing before its status escapes.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
budget_runner="${DO_WORK_TEST_BUDGET_RUNNER:-$repo_root/_dev/tests/run-go-tests-with-budget.sh}"
fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/go-test-budget-behavior.XXXXXX")"
trap 'rm -rf -- "$fixture_root"' EXIT

printf 'module example.invalid/budget-fixture\n\ngo 1.26.1\n' > "$fixture_root/go.mod"
cat > "$fixture_root/budget_test.go" <<'GO'
package budgetfixture

import "testing"

func TestPassingFile(t *testing.T) {}

func TestFailingFile(t *testing.T) {
	t.Fatal("intentional fixture failure")
}
GO
duration_log="$fixture_root/test-durations.tsv"
printf 'run_id\tfile\tseconds\tother_gate_processes\n' > "$duration_log"

run_fixture() {
  local run_id="$1"
  local test_pattern="$2"
  DO_WORK_TEST_DURATION_LOG="$duration_log" \
    DO_WORK_TEST_RUN_ID="$run_id" \
    DO_WORK_TEST_OTHER_GATE_PROCESSES=0 \
    DO_WORK_TEST_REPO_ROOT="$fixture_root" \
    DO_WORK_TEST_ENFORCE_BUDGET=yes \
    DO_WORK_TEST_FILE_BUDGET_SECONDS=30 \
    bash "$budget_runner" "$fixture_root" -run "$test_pattern" ./...
}

run_fixture passing-run '^TestPassingFile$' >/dev/null
failing_status=0
run_fixture failing-run '^TestFailingFile$' >/dev/null 2>&1 || failing_status=$?
if [ "$failing_status" -ne 1 ]; then
  printf 'FAIL: failing Go fixture returned %s; expected the original go test status 1.\n' \
    "$failing_status" >&2
  exit 1
fi

for run_id in passing-run failing-run; do
  if ! awk -F '\t' -v run_id="$run_id" '
    $1 == run_id && $2 == "budget_test.go" && $3 ~ /^[0-9]+([.][0-9]+)?$/ && $4 == "0" { matches++ }
    END { exit(matches == 1 ? 0 : 1) }
  ' "$duration_log"; then
    printf 'FAIL: %s did not log exactly one attributed budget_test.go duration row.\n' \
      "$run_id" >&2
    exit 1
  fi
done

printf 'Go test budget behavior probes passed.\n'

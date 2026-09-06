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

func TestHeavyAlpha(t *testing.T) {
	t.Fatal("TestHeavyAlpha must be excluded by heavy prefixes")
}

func TestHeavyBeta(t *testing.T) {
	t.Fatal("TestHeavyBeta must be excluded by heavy prefixes")
}

func TestSpecialLiteral(t *testing.T) {}
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

# Probe 1: Prefix exclusion skips heavy tests without running discovery
exclude_status=0
DO_WORK_TEST_DURATION_LOG="$duration_log" \
  DO_WORK_TEST_RUN_ID="exclude-heavy-run" \
  DO_WORK_TEST_OTHER_GATE_PROCESSES=0 \
  DO_WORK_TEST_REPO_ROOT="$fixture_root" \
  DO_WORK_TEST_ENFORCE_BUDGET=yes \
  DO_WORK_TEST_FILE_BUDGET_SECONDS=30 \
  DO_WORK_GO_TEST_EXCLUDE_PREFIXES="TestHeavy,TestFailing" \
  bash "$budget_runner" "$fixture_root" -run '^(TestPassingFile|TestHeavyAlpha|TestHeavyBeta)$' ./... >/dev/null 2>&1 || exclude_status=$?

if [ "$exclude_status" -ne 0 ]; then
  printf 'FAIL: exclude-heavy-run failed (%s); heavy tests were not excluded.\n' "$exclude_status" >&2
  exit 1
fi

# Probe 2: Exhaustive exclusion refuses when no fast tests remain
refusal_status=0
refusal_output="$(
  DO_WORK_TEST_DURATION_LOG="$duration_log" \
    DO_WORK_TEST_RUN_ID="refusal-run" \
    DO_WORK_TEST_OTHER_GATE_PROCESSES=0 \
    DO_WORK_TEST_REPO_ROOT="$fixture_root" \
    DO_WORK_TEST_ENFORCE_BUDGET=yes \
    DO_WORK_TEST_FILE_BUDGET_SECONDS=30 \
    DO_WORK_GO_TEST_EXCLUDE_PREFIXES="Test" \
    bash "$budget_runner" "$fixture_root" ./... 2>&1
)" || refusal_status=$?

if [ "$refusal_status" -ne 1 ]; then
  printf 'FAIL: exhaustive exclusion returned %s; expected status 1.\n' "$refusal_status" >&2
  exit 1
fi
if ! grep -Fq 'no fast Go tests remain after applying the heavy prefixes' <<< "$refusal_output"; then
  printf 'FAIL: exhaustive exclusion output missing refusal diagnostic: %s\n' "$refusal_output" >&2
  exit 1
fi

# Probe 3: Regex metacharacters in prefixes are escaped and do not crash or over-match
meta_status=0
DO_WORK_TEST_DURATION_LOG="$duration_log" \
  DO_WORK_TEST_RUN_ID="meta-escape-run" \
  DO_WORK_TEST_OTHER_GATE_PROCESSES=0 \
  DO_WORK_TEST_REPO_ROOT="$fixture_root" \
  DO_WORK_TEST_ENFORCE_BUDGET=yes \
  DO_WORK_TEST_FILE_BUDGET_SECONDS=30 \
  DO_WORK_GO_TEST_EXCLUDE_PREFIXES='TestSpecial.,TestBracket[' \
  bash "$budget_runner" "$fixture_root" -run '^TestSpecialLiteral$' ./... >/dev/null 2>&1 || meta_status=$?

if [ "$meta_status" -ne 0 ]; then
  printf 'FAIL: metacharacter escape run returned %s; regex escaping failed.\n' "$meta_status" >&2
  exit 1
fi

printf 'Go test budget behavior probes passed.\n'

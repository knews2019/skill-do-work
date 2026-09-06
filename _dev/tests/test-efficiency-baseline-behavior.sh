#!/usr/bin/env bash
# Behavioral verification for test-efficiency baseline and duration logging extensions.
# Pins exit status preservation, opt-in log isolation, subprocess counting,
# and Go duration attribution.
set -euo pipefail

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"

fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/baseline-behavior.XXXXXX")"
cleanup_fixture() {
  rm -rf -- "$fixture_root"
}
trap cleanup_fixture EXIT

# shellcheck source=_dev/tests/test-duration-log.sh disable=SC1091
source "$repo_root/_dev/tests/test-duration-log.sh"

# Probes 1 and 2 drive measure_command_efficiency, the writer the baseline tool actually
# calls. They used to drive record_test_efficiency, a bash twin of the writer's own
# append that nothing called, so a regression in the real writer's opt-in or header
# would have passed them.

# 1. Opt-in isolation: without DO_WORK_TEST_EFFICIENCY_LOG, efficiency logging remains dormant.
default_log="$repo_root/do-work/test-efficiency.tsv"
default_log_before="$( [ -e "$default_log" ] && cksum "$default_log" || printf 'absent')"
(
  unset DO_WORK_TEST_EFFICIENCY_LOG
  measure_command_efficiency "dormant-case" "warm" true >/dev/null
)
default_log_after="$( [ -e "$default_log" ] && cksum "$default_log" || printf 'absent')"
if [ "$default_log_before" != "$default_log_after" ]; then
  printf 'FAIL: measure_command_efficiency wrote to the default log without the opt-in variable set.\n' >&2
  exit 1
fi

# 2. Opt-in recording: when DO_WORK_TEST_EFFICIENCY_LOG is set, writes valid TSV.
eff_tsv="$fixture_root/efficiency.tsv"
export DO_WORK_TEST_EFFICIENCY_LOG="$eff_tsv"
(
  measure_command_efficiency "optin-case" "warm" true >/dev/null
)

if [ ! -s "$eff_tsv" ]; then
  printf 'FAIL: efficiency TSV was not written when DO_WORK_TEST_EFFICIENCY_LOG was set.\n' >&2
  exit 1
fi

header_line="$(sed -n '1p' "$eff_tsv")"
expected_header="$(printf 'run_id\tcase_name\tcache_condition\trevision\ttoolchain\tconcurrency\twall_seconds\tuser_cpu_seconds\tsys_cpu_seconds\ttotal_cpu_seconds\tgo_accum_seconds\tsubprocess_counts\texit_status')"
if [ "$header_line" != "$expected_header" ]; then
  printf 'FAIL: invalid efficiency TSV header: %s\n' "$header_line" >&2
  exit 1
fi

if ! grep -q "optin-case" "$eff_tsv"; then
  printf 'FAIL: optin-case record missing from %s\n' "$eff_tsv" >&2
  exit 1
fi

# 3. Exit status preservation in measure_command_efficiency
(
  failing_status=0
  measure_command_efficiency "failing-test" "cold" bash -c 'exit 42' >/dev/null 2>&1 || failing_status=$?
  if [ "$failing_status" -ne 42 ]; then
    printf 'FAIL: measure_command_efficiency returned %s; expected child exit status 42.\n' "$failing_status" >&2
    exit 1
  fi
)

# 4. Subprocess counting and measurement output
(
  output_json="$(measure_command_efficiency "subproc-test" "warm" git --version)"
  if ! python3 -c '
import sys, json
data = json.loads(sys.argv[1])
assert data["case"] == "subproc-test"
assert data["exit_status"] == 0
assert data["wall_seconds"] > 0
assert data["total_cpu_seconds"] >= 0
assert "git" in data["subprocess_counts"]
' "$output_json"; then
    printf 'FAIL: measure_command_efficiency output did not match expected structure: %s\n' "$output_json" >&2
    exit 1
  fi
)

# 5. Go JSON event duration distinguishing
(
  # Produce mock Go JSON stream through measure_command_efficiency
  mock_stream_script="$fixture_root/mock_go_json.sh"
  cat > "$mock_stream_script" <<'EOF'
#!/bin/sh
printf '{"Time":"2026-09-06T00:00:00Z","Action":"run","Package":"pkg/test","Test":"TestA"}\n'
printf '{"Time":"2026-09-06T00:00:01Z","Action":"pass","Package":"pkg/test","Test":"TestA","Elapsed":1.5}\n'
printf '{"Time":"2026-09-06T00:00:01Z","Action":"run","Package":"pkg/test","Test":"TestB"}\n'
printf '{"Time":"2026-09-06T00:00:02Z","Action":"pass","Package":"pkg/test","Test":"TestB","Elapsed":2.0}\n'
printf '{"Time":"2026-09-06T00:00:02Z","Action":"pass","Package":"pkg/test","Elapsed":3.0}\n'
EOF
  chmod +x "$mock_stream_script"

  go_output_json="$(measure_command_efficiency "go-json-test" "cold" "$mock_stream_script" -json)"
  if ! python3 -c '
import sys, json
data = json.loads(sys.argv[1])
assert data["case"] == "go-json-test"
val_ga = data.get("go_accum_seconds")
val_pw = data.get("package_wall_seconds")
assert val_ga == "3.500", "expected 3.500 got " + str(val_ga)
assert val_pw == "3.000", "expected 3.000 got " + str(val_pw)
' "$go_output_json"; then
    printf 'FAIL: Go JSON duration parser did not correctly distinguish accumulated from package wall: %s\n' "$go_output_json" >&2
    exit 1
  fi
)

printf 'Test efficiency baseline behavior probes passed.\n'

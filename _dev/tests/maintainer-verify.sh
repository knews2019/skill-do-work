#!/usr/bin/env bash
set -euo pipefail

# Version FLOORS, not exact pins: this gate must be runnable on a machine that has a
# newer toolchain than the one it was written against. The exact pins these replaced were
# buying one thing — gofmt has no version flag and its output can change between Go
# releases, so a newer toolchain can reformat a file this repo considers formatted and the
# gofmt lane will name it. That is a legible failure with a legible fix (run the newer
# gofmt, commit the result, or raise this floor), which is the trade this floor accepts.
minimum_go_version="go1.26.1"
minimum_shellcheck_version="0.11.0"
test_file_budget_seconds=30
script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"
test_duration_log_path="$repo_root/do-work/test-durations.tsv"
test_run_id="$(/bin/date -u +%Y%m%dT%H%M%SZ)-$$"
other_gate_processes="$(
  /bin/ps -Ao pid=,comm=,args= | /usr/bin/awk -v own_pid="$$" '
    $1 != own_pid && $2 ~ /(^|\/)(ba)?sh$/ && $0 ~ /\/maintainer-verify\.sh( |$)/ { count++ }
    END { print count + 0 }
  '
)"
export DO_WORK_TEST_DURATION_LOG="$test_duration_log_path"
export DO_WORK_TEST_RUN_ID="$test_run_id"
export DO_WORK_TEST_OTHER_GATE_PROCESSES="$other_gate_processes"
export DO_WORK_TEST_REPO_ROOT="$repo_root"
self_test_exit_root=''
# Self-test fixture context. Script-scoped rather than local to run_self_test because the
# per-run assertion helpers below read it; only the self-test path ever assigns it.
self_test_root=''
fixture_root=''
fixture_script=''
fixture_go_root=''
with_node_bin=''
# The exact line each engine-gated lane prints when it did not run. The whole heavy
# tier and the isolated --heavy-lane path print the same bytes, so a caller that reads
# one of them reads the other; the lane runner keys `skipped` on the SKIP: prefix.
browser_lane_skip_line='SKIP: no browser is available; strict browser behavior lane was not run. Set QUEUE_KANBAN_BROWSER to name one.'
node_lane_skip_line='SKIP: Node is unavailable; strict JavaScript behavior lane was not run.'

# Reports whether a Node engine is available for the strict JavaScript behavior lane.
node_engine_available() {
  command -v node >/dev/null 2>&1
}

# Reports whether a browser engine is available for the strict browser behavior lane.
# QUEUE_KANBAN_BROWSER names an engine that is not on PATH under a well-known name.
browser_engine_available() {
  local browser_probe_candidate
  if [ -n "${QUEUE_KANBAN_BROWSER:-}" ]; then
    return 0
  fi
  for browser_probe_candidate in google-chrome google-chrome-stable chromium chromium-browser chrome; do
    if command -v "$browser_probe_candidate" >/dev/null 2>&1; then
      return 0
    fi
  done
  return 1
}

print_usage_and_exit() {
  printf 'usage: %s [--heavy|--heavy-lane <lane-id>|--self-test]\n' "$0" >&2
  exit 2
}

run_budgeted_go_tests() {
  local module_directory="$1"
  local enforce_test_budget=yes
  shift
  if [ -n "${MAINTAINER_VERIFY_SELFTEST_LOG:-}" ]; then
    (
      cd "$module_directory"
      go test -count=1 "$@"
    )
    return
  fi
  if [ "$verification_tier" = heavy ]; then
    enforce_test_budget=no
  fi
  DO_WORK_TEST_ENFORCE_BUDGET="$enforce_test_budget" \
  DO_WORK_TEST_FILE_BUDGET_SECONDS="$test_file_budget_seconds" \
    bash "$repo_root/_dev/tests/run-go-tests-with-budget.sh" "$module_directory" "$@"
}

# Runs one expensive Go test stage, or reports it from stored evidence when the
# complete relevant inputs of that stage are provably unchanged. The decision,
# the revocation and the recording are three separate CLI calls on purpose: the
# revocation has to happen before the stage runs and has to outlive a stage that
# never finishes, so a failed or interrupted run can never leave an older green
# reusable. Anything this function cannot determine selects execution.
#
# A reused stage writes no per-file duration rows and enforces no per-file budget
# for that run: it inherits the budget verdict of the run whose evidence it is
# reusing, and the REUSED line says so. The budget catches tests getting slower
# as they are added, which is an input-determined property the fingerprint
# covers; it does not catch a breach caused purely by machine contention, which
# is the failure this reuse deliberately stops reporting on an unchanged tree.
run_stage_with_evidence() (
  # The manifest's runtime probes exclude global and system Git configuration.
  # Execute under that same configuration so a cached pass and a fresh test see
  # the same Git behavior. The subshell keeps this isolation local to the stage.
  export GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null
  local stage_id="$1"
  local module_directory="$2"
  shift 2
  local module_relative_path="${module_directory#"$repo_root/"}"
  local cli_launcher="$repo_root/skills/do-work/tools/do-work-cli.sh"
  # A module the prefix strip cannot make repository-relative would be sealed
  # under a path the manifest cannot declare, which the manifest comparison below
  # would refuse on every run — reuse silently off, forever, with a green gate.
  # Stop instead: that is a wiring error, not a test result.
  if [ "$module_relative_path" = "$module_directory" ]; then
    printf 'maintainer-verify: stage %s: module %s is not inside %s\n' \
      "$stage_id" "$module_directory" "$repo_root" >&2
    return 2
  fi
  local decision_line='' decision_status=0 stage_status=0
  local disposition='' reason='' fingerprint='' recorded_at='' extra_fields=''
  local stage_argv
  stage_argv=(run-go-tests-with-budget.sh "$module_relative_path" "$@")

  # The self-test drives this gate through a command shim whose `bash` case
  # accepts only the aggregate suite, so a decision child would exit 64 there and
  # fail the fixture. --self-test proves the gate's stage LIST and must keep
  # counting nine stages exactly once; reuse is proved by its own probe,
  # _dev/tests/fast-stage-reuse-behavior.sh, where it can be tested far more
  # thoroughly.
  if [ -n "${MAINTAINER_VERIFY_SELFTEST_LOG:-}" ]; then
    run_budgeted_go_tests "$module_directory" "$@"
    return
  fi

  # The decision child runs from the repository root, so its own working
  # directory is fixed, with five names removed from its environment.
  # DO_WORK_TEST_RUN_ID is a timestamp-plus-PID label,
  # DO_WORK_TEST_OTHER_GATE_PROCESSES a `ps` sample and DO_WORK_TEST_DURATION_LOG
  # an output path; OLDPWD names the caller's previous directory and SHLVL counts
  # how deeply the caller's shells are nested. All five change between runs, none
  # decides any assertion, and the fingerprint seals the whole environment, so
  # with them present no stage could ever reuse — SHLVL alone was measured
  # defeating it between a terminal and a wrapper script. DO_WORK_TEST_ENFORCE_BUDGET
  # and DO_WORK_TEST_FILE_BUDGET_SECONDS are deliberately NOT removed: they change
  # the verdict. Status and output are captured separately, because a decision
  # command that could not run must read as "execute", never as "nothing
  # changed".
  if [ "${DO_WORK_FAST_STAGE_REUSE:-on}" = 'off' ]; then
    # The measurement escape hatch skips reuse and recording, but must still
    # revoke any old pass before execution can fail or be interrupted.
    decision_line='executed reuse_disabled - -'
  else
    decision_line="$(
      cd "$repo_root" &&
        env -u DO_WORK_TEST_RUN_ID -u DO_WORK_TEST_OTHER_GATE_PROCESSES \
          -u DO_WORK_TEST_DURATION_LOG -u OLDPWD -u SHLVL \
          bash "$cli_launcher" --repo-root "$repo_root" \
          decide-fast-stage --stage "$stage_id" -- "${stage_argv[@]}"
    )" || decision_status=$?
    if [ "$decision_status" -ne 0 ]; then
      decision_line='executed decision_unavailable - -'
    fi
  fi
  if ! read -r disposition reason fingerprint recorded_at extra_fields <<< "$decision_line"; then
    disposition=''
  fi
  if [ -z "$disposition" ] || [ -z "$reason" ] || [ -z "$fingerprint" ] || \
    [ -z "$recorded_at" ] || [ -n "$extra_fields" ]; then
    disposition='executed'
    reason='decision_unparseable'
    fingerprint='-'
  fi

  if [ "$disposition" = 'reused' ]; then
    printf 'maintainer-verify: stage %s: REUSED (%s, recorded %s; per-file budget verdict inherited from that run)\n' \
      "$stage_id" "$reason" "$recorded_at"
    return 0
  fi
  printf 'maintainer-verify: stage %s: EXECUTING (%s)\n' "$stage_id" "$reason"
  # A prior green this gate cannot revoke is the false-green shape the reuse
  # rule exists to prevent, so an unrevocable record fails the gate rather than
  # letting the stage run unprotected.
  bash "$cli_launcher" --repo-root "$repo_root" \
    invalidate-fast-stage --stage "$stage_id" > /dev/null || return $?
  run_budgeted_go_tests "$module_directory" "$@" || stage_status=$?
  if [ "$stage_status" -ne 0 ]; then
    return "$stage_status"
  fi
  if [ "$fingerprint" != '-' ]; then
    # A recording that does not land costs the next run one execution, which is
    # the safe direction, so it reports rather than failing a stage that passed.
    if ! (
      cd "$repo_root" &&
        env -u DO_WORK_TEST_RUN_ID -u DO_WORK_TEST_OTHER_GATE_PROCESSES \
          -u DO_WORK_TEST_DURATION_LOG -u OLDPWD -u SHLVL \
          bash "$cli_launcher" --repo-root "$repo_root" \
          record-fast-stage --stage "$stage_id" --fingerprint "$fingerprint" \
          --stage-exit-status 0 -- "${stage_argv[@]}" > /dev/null
    ); then
      printf 'maintainer-verify: stage %s: evidence not recorded; the next run executes it again\n' \
        "$stage_id" >&2
    fi
  fi
  return 0
)

# Returns 0 when $1 is at or above the floor $2. Compares dot-separated components as
# integers, so 0.11.0 clears a 0.9.9 floor where a lexical compare would not. A missing
# component reads as 0 (`0.11` clears `0.11.0`), and a trailing non-digit run is dropped
# so a prerelease like `1.27rc1` compares as 27.
version_at_least() {
  local candidate_rest="$1"
  local floor_rest="$2"
  local candidate_part
  local floor_part

  while [ -n "$candidate_rest" ] || [ -n "$floor_rest" ]; do
    candidate_part="${candidate_rest%%.*}"
    floor_part="${floor_rest%%.*}"
    candidate_part="${candidate_part%%[!0-9]*}"
    floor_part="${floor_part%%[!0-9]*}"
    [ -z "$candidate_part" ] && candidate_part=0
    [ -z "$floor_part" ] && floor_part=0
    if [ "$candidate_part" -gt "$floor_part" ]; then
      return 0
    fi
    if [ "$candidate_part" -lt "$floor_part" ]; then
      return 1
    fi
    case "$candidate_rest" in *.*) candidate_rest="${candidate_rest#*.}" ;; *) candidate_rest='' ;; esac
    case "$floor_rest" in *.*) floor_rest="${floor_rest#*.}" ;; *) floor_rest='' ;; esac
  done
  return 0
}

fail_self_test() {
  printf 'FAIL: maintainer-verify self-test: %s\n' "$1" >&2
  return 1
}

cleanup_self_test_exit() {
  if [ -n "$self_test_exit_root" ]; then
    rm -rf -- "$self_test_exit_root"
    self_test_exit_root=''
  fi
}

write_command_shim() {
  local shim_path="$1"

  {
    printf '%s\n' '#!/bin/bash'
    printf '%s\n' 'set -eu'
    cat <<'SHIM'
command_name="${0##*/}"
failure_stage="${MAINTAINER_VERIFY_FAIL_STAGE:-}"
stage_log="${MAINTAINER_VERIFY_SELFTEST_LOG:?}"

record_stage() {
  printf '%s\n' "$1" >> "$stage_log"
  if [ "$failure_stage" = "$1" ]; then
    exit 41
  fi
}

case "$command_name" in
  go)
    if [ "${1:-}" = 'version' ]; then
      if [ "$#" -ne 1 ]; then
        exit 64
      fi
      printf '%s\n' 'go-version' >> "$stage_log"
      if [ "$failure_stage" = 'go-version' ]; then
        printf '%s\n' 'go version go1.20.0 fixture/arch'
      elif [ -n "${MAINTAINER_VERIFY_SELFTEST_NEWER_TOOLS:-}" ]; then
        printf '%s\n' 'go version go1.99.0 fixture/arch'
      else
        printf '%s\n' 'go version go1.26.1 fixture/arch'
      fi
      exit 0
    fi
    if [ "$1" = 'env' ]; then
      if [ "$#" -ne 2 ] || [ "$2" != 'GOROOT' ]; then
        exit 64
      fi
      printf '%s\n' "${MAINTAINER_VERIFY_SELFTEST_GOROOT:?}"
      exit 0
    fi
    case "$PWD" in
      */skills/do-work-board/tools/queue-kanban)
        if [ "$#" -eq 2 ] && [ "$1" = 'vet' ] && [ "$2" = './...' ]; then
          record_stage 'board-vet'
        elif [ "$#" -eq 3 ] && [ "$1" = 'test' ] && [ "$2" = '-count=1' ] && [ "$3" = './...' ]; then
          if [ "${QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR:-}" = '1' ]; then
            record_stage 'board-test-strict'
          else
            record_stage 'board-test'
          fi
        else
          printf 'unexpected board go argv (%s): %s\n' "$#" "$*" >&2
          exit 64
        fi
        ;;
      */skills/do-work/tools/do-work-cli)
        if [ "$#" -eq 2 ] && [ "$1" = 'vet' ] && [ "$2" = './...' ]; then
          record_stage 'cli-vet'
        elif [ "$#" -eq 3 ] && [ "$1" = 'test' ] && [ "$2" = '-count=1' ] && [ "$3" = './...' ]; then
          record_stage 'cli-test'
        elif [ "$#" -eq 4 ] && [ "$1" = 'test' ] && [ "$2" = '-count=1' ] && [ "$3" = '-short' ] && [ "$4" = './...' ]; then
          record_stage 'cli-test'
        else
          printf 'unexpected cli go argv (%s): %s\n' "$#" "$*" >&2
          exit 64
        fi
        ;;
      *) exit 64 ;;
    esac
    ;;
  gofmt)
    if [ "$#" -ne 3 ] || [ "$1" != '-l' ] || [ "$2" != '--' ] || \
      [ "$3" != '_dev/tests/fixture-go.go' ]; then
      exit 64
    fi
    if [ "$failure_stage" = 'gofmt-unformatted' ]; then
      printf '%s\n' 'gofmt-lint' >> "$stage_log"
      printf '%s\n' '_dev/tests/fixture-go.go'
      exit 0
    fi
    record_stage 'gofmt-lint'
    ;;
  shellcheck)
    if [ "${1:-}" = '--version' ]; then
      if [ "$#" -ne 1 ]; then
        exit 64
      fi
      printf '%s\n' 'shellcheck-version' >> "$stage_log"
      printf '%s\n' 'ShellCheck - shell script analysis tool'
      if [ "$failure_stage" = 'shellcheck-version' ]; then
        printf '%s\n' 'version: 0.9.9'
      elif [ -n "${MAINTAINER_VERIFY_SELFTEST_NEWER_TOOLS:-}" ]; then
        printf '%s\n' 'version: 9.9.9'
      else
        printf '%s\n' 'version: 0.11.0'
      fi
      exit 0
    fi
    if [ "$#" -ne 3 ] || [ "$1" != '--severity=warning' ] || \
      [ "$2" != '--' ] || [ "$3" != '_dev/tests/fixture-shell.sh' ]; then
      exit 64
    fi
    record_stage 'shellcheck-lint'
    ;;
  git)
    if [ "$#" -ne 6 ] || [ "$1" != '-C' ] || \
      [ "$2" != "${MAINTAINER_VERIFY_EXPECTED_REPO_ROOT:?}" ] || \
      [ "$3" != 'ls-files' ] || [ "$4" != '-z' ] || [ "$5" != '--' ]; then
      exit 64
    fi
    case "$6" in
      '*.sh') printf '%s\0' '_dev/tests/fixture-shell.sh' ;;
      '*.go') printf '%s\0' '_dev/tests/fixture-go.go' ;;
      *) exit 64 ;;
    esac
    ;;
  bash)
    if [ "${1##*/}" != 'contract-regressions.sh' ] || [ "$#" -ne 1 ]; then
      exit 64
    fi
    record_stage 'aggregate'
    ;;
  node)
    exit 0
    ;;
  *) exit 64 ;;
esac
SHIM
  } > "$shim_path"
  chmod +x "$shim_path"
}

count_stage_occurrences() {
  local stage_name="$1"
  local log_path="$2"
  local stage_line
  local stage_count=0

  while IFS= read -r stage_line; do
    if [ "$stage_line" = "$stage_name" ]; then
      stage_count=$((stage_count + 1))
    fi
  done < "$log_path"
  printf '%s\n' "$stage_count"
}

assert_success_stages() {
  local log_path="$1"
  local expect_strict="$2"
  local expected_stage
  local actual_count
  local total_count=0
  local expected_count=9
  local expected_board_stage='board-test'
  local stage_line

  # Only --heavy with Node present carries the strict JavaScript marker, so the shim records
  # that one board test run under a different stage name; there is no separate strict lane.
  # Both tiers run the same nine stages — the tier changes the board stage, not the count.
  if [ "$expect_strict" = 'yes' ]; then
    expected_board_stage='board-test-strict'
  fi
  for expected_stage in \
    go-version shellcheck-version shellcheck-lint gofmt-lint aggregate \
    board-vet "$expected_board_stage" cli-vet cli-test; do
    actual_count="$(count_stage_occurrences "$expected_stage" "$log_path")"
    if [ "$actual_count" -ne 1 ]; then
      fail_self_test "$expected_stage ran $actual_count times; want exactly once"
      return 1
    fi
  done
  while IFS= read -r stage_line; do
    total_count=$((total_count + 1))
  done < "$log_path"
  if [ "$total_count" -ne "$expected_count" ]; then
    fail_self_test "success run recorded $total_count stages; want $expected_count"
    return 1
  fi
}

# Runs the shimmed fixture and asserts its stage list. The gate arguments after the fourth
# positional select the tier, so one call covers a whole tier's expected stages.
run_success_fixture() {
  local fixture_label="$1"
  local fixture_path_directory="$2"
  local expect_strict="$3"
  local newer_tools_marker="$4"
  shift 4
  local fixture_log="$self_test_root/$fixture_label.log"
  local fixture_output="$self_test_root/$fixture_label.out"

  : > "$fixture_log"
  if ! PATH="$fixture_path_directory" \
    QUEUE_KANBAN_BROWSER='' \
    MAINTAINER_VERIFY_SELFTEST_LOG="$fixture_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
    MAINTAINER_VERIFY_SELFTEST_GOROOT="$fixture_go_root" \
    MAINTAINER_VERIFY_SELFTEST_NEWER_TOOLS="$newer_tools_marker" \
    /bin/bash "$fixture_script" "$@" > "$fixture_output" 2>&1; then
    sed 's/^/  /' "$fixture_output" >&2
    sed 's/^/  stage: /' "$fixture_log" >&2
    fail_self_test "the $fixture_label fixture exited nonzero"
    return 1
  fi
  assert_success_stages "$fixture_log" "$expect_strict"
}

# Injects a failure at one stage and requires the gate to carry it out. A stage whose
# nonzero status the gate swallowed would exit zero here.
assert_failure_stage() {
  local stage_name="$1"
  shift
  local failure_log="$self_test_root/failure-$stage_name.log"
  local failure_output="$self_test_root/failure-$stage_name.out"

  : > "$failure_log"
  if PATH="$with_node_bin" \
    QUEUE_KANBAN_BROWSER='' \
    MAINTAINER_VERIFY_SELFTEST_LOG="$failure_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
    MAINTAINER_VERIFY_SELFTEST_GOROOT="$fixture_go_root" \
    MAINTAINER_VERIFY_FAIL_STAGE="$stage_name" \
    /bin/bash "$fixture_script" "$@" > "$failure_output" 2>&1; then
    fail_self_test "$stage_name failure exited zero"
    return 1
  fi
  if [ "$stage_name" = 'gofmt-unformatted' ] && \
    ! grep -q '_dev/tests/fixture-go.go' "$failure_output"; then
    fail_self_test 'the unformatted-Go failure did not name the offending file'
    return 1
  fi
}

# Asserts that an isolated --heavy-lane run with no engine available prints the lane's
# own SKIP line and exits 0. A lane that exits red instead gives the caller no "did not
# run" evidence to record, and a lane that exits 0 silently is indistinguishable from a
# lane that passed.
assert_isolated_lane_skips() {
  local lane_id="$1"
  local expected_skip_line="$2"
  local engineless_path_directory="$3"
  local lane_skip_output="$self_test_root/lane-skip-$lane_id.out"

  if ! PATH="$engineless_path_directory" \
    QUEUE_KANBAN_BROWSER='' \
    /bin/bash "$fixture_script" --heavy-lane "$lane_id" > "$lane_skip_output" 2>&1; then
    sed 's/^/  /' "$lane_skip_output" >&2
    fail_self_test "the isolated $lane_id lane exited nonzero with no engine available"
    return 1
  fi
  if ! grep -qF -- "$expected_skip_line" "$lane_skip_output"; then
    sed 's/^/  /' "$lane_skip_output" >&2
    fail_self_test "the isolated $lane_id lane did not print its explicit skip"
    return 1
  fi
}

run_self_test() {
  local without_node_bin
  local generic_shim
  local stage_name
  local mutated_fixture_script
  local mutation_output
  local lane_probe_script
  local lane_cli_binary
  local lane_duration_log

  self_test_root="$(mktemp -d)"
  self_test_exit_root="$self_test_root"
  trap cleanup_self_test_exit EXIT
  fixture_root="$self_test_root/repository"
  fixture_script="$fixture_root/_dev/tests/maintainer-verify.sh"
  fixture_go_root="$self_test_root/go-root"
  with_node_bin="$self_test_root/with-node-bin"
  without_node_bin="$self_test_root/without-node-bin"
  mkdir -p \
    "$fixture_root/_dev/tests" \
    "$fixture_root/skills/do-work-board/tools/queue-kanban" \
    "$fixture_root/skills/do-work/tools/do-work-cli" \
    "$fixture_go_root/bin" "$with_node_bin" "$without_node_bin"
  cp "$script_path" "$fixture_script"
  chmod +x "$fixture_script"
  printf '%s\n' '#!/usr/bin/env bash' > "$fixture_root/_dev/tests/fixture-shell.sh"
  printf '%s\n' '#!/usr/bin/env bash' > "$fixture_root/_dev/tests/contract-regressions.sh"
  printf '%s\n' 'package fixture' > "$fixture_root/_dev/tests/fixture-go.go"

  generic_shim="$self_test_root/command-shim"
  write_command_shim "$generic_shim"
  for stage_name in go shellcheck git bash node; do
    ln -s "$generic_shim" "$with_node_bin/$stage_name"
  done
  for stage_name in go shellcheck git bash; do
    ln -s "$generic_shim" "$without_node_bin/$stage_name"
  done
  ln -s "$generic_shim" "$fixture_go_root/bin/gofmt"

  # The default tier is the gate the loop runs: no strict marker even with Node present, so
  # the board stage records unstrict and Node presence changes no lane selection.
  run_success_fixture default "$with_node_bin" no ''
  # Floors, not exact pins: a toolchain newer than the floor must still pass the default tier.
  run_success_fixture newer-tools "$with_node_bin" no yes
  # --heavy is the only tier that carries the strict JavaScript marker.
  run_success_fixture heavy "$with_node_bin" yes '' --heavy
  # Without Node the heavy tier drops to the ordinary board run and must say which lane it
  # gave up, since a silent drop is indistinguishable from a lane that passed.
  run_success_fixture heavy-no-node "$without_node_bin" no '' --heavy
  if ! grep -q 'SKIP: Node is unavailable; strict JavaScript behavior lane was not run.' \
    "$self_test_root/heavy-no-node.out"; then
    fail_self_test 'the no-Node heavy path did not print its explicit skip'
    return 1
  fi
  # The same contract one lane at a time: the planner selects lanes individually, so an
  # isolated lane with no engine must announce the skip the whole tier announces.
  assert_isolated_lane_skips queue-kanban-browser "$browser_lane_skip_line" "$without_node_bin" || return 1
  assert_isolated_lane_skips queue-kanban-javascript "$node_lane_skip_line" "$without_node_bin" || return 1

  for stage_name in \
    go-version shellcheck-version shellcheck-lint gofmt-lint gofmt-unformatted \
    aggregate board-vet board-test \
    cli-vet cli-test; do
    assert_failure_stage "$stage_name"
  done
  assert_failure_stage board-test-strict --heavy

  mutated_fixture_script="$self_test_root/maintainer-verify-mutated.sh"
  sed '/^run_verification()/,$ s/QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1/QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=0/' \
    "$fixture_script" > "$mutated_fixture_script"
  if cmp -s "$fixture_script" "$mutated_fixture_script"; then
    fail_self_test 'strict-marker mutation did not alter the production command'
    return 1
  fi
  mv "$mutated_fixture_script" "$fixture_script"
  chmod +x "$fixture_script"
  mutation_output="$self_test_root/wrong-strict-regex.out"
  if /bin/bash "$fixture_script" --self-test > "$mutation_output" 2>&1; then
    fail_self_test 'a disabled strict JavaScript marker left --self-test green'
    return 1
  fi
  cp "$script_path" "$fixture_script"
  chmod +x "$fixture_script"
  if ! cmp -s "$script_path" "$fixture_script"; then
    fail_self_test 'the strict-marker mutation fixture was not restored'
    return 1
  fi

  lane_probe_script="$self_test_root/lane-probe.sh"
  lane_cli_binary="$self_test_root/lane-cli"
  lane_duration_log="$self_test_root/lane-durations.tsv"
  printf '%s\n' '#!/usr/bin/env bash' 'set -eu' \
    'test "${DO_WORK_MAINTAINER_TIER:-}" = heavy' \
    'test -x "${DO_WORK_TEST_DO_WORK_CLI_BINARY:?}"' > "$lane_probe_script"
  printf '%s\n' '#!/usr/bin/env bash' 'exit 0' > "$lane_cli_binary"
  chmod +x "$lane_probe_script" "$lane_cli_binary"
  if ! DO_WORK_TEST_DURATION_LOG="$lane_duration_log" \
    DO_WORK_TEST_RUN_ID=lane-self-test \
    DO_WORK_TEST_DO_WORK_CLI_BINARY="$lane_cli_binary" \
    run_heavy_probe_lane "$lane_probe_script"; then
    fail_self_test 'the selected shell-lane fixture failed'
    return 1
  fi
  if [ "$(grep -c '^lane-self-test' "$lane_duration_log")" -ne 1 ]; then
    fail_self_test 'one selected shell lane must record exactly one duration row'
    return 1
  fi

  trap - EXIT
  cleanup_self_test_exit
  printf 'Maintainer verification self-test passed.\n'
}

run_verification() {
  local verification_tier="$1"
  local gate_started_seconds="$SECONDS"
  local go_version_output
  local go_version_text
  local shellcheck_version_output
  local shellcheck_version_line
  local shellcheck_version_text=''
  local tracked_shell_file
  local tracked_shell_files=()
  local gofmt_command
  local tracked_go_file
  local tracked_go_files=()
  local unformatted_go_files

  for required_command in git go shellcheck bash; do
    if ! command -v "$required_command" >/dev/null 2>&1; then
      printf 'maintainer-verify: required command is unavailable: %s\n' "$required_command" >&2
      return 1
    fi
  done

  printf 'maintainer-verify: checking Go %s or newer\n' "$minimum_go_version"
  go_version_output="$(go version)"
  read -r _ _ go_version_text _ <<< "$go_version_output"
  if ! version_at_least "${go_version_text#go}" "${minimum_go_version#go}"; then
    printf 'maintainer-verify: Go version is %s; require %s or newer\n' \
      "$go_version_text" "$minimum_go_version" >&2
    return 1
  fi
  # gofmt carries no version flag, so take it from the toolchain this run just
  # version-checked, via its GOROOT, rather than from PATH: a stray older gofmt formats
  # differently and would make this gate's verdict depend on which Go happens to come
  # first on PATH. That matters more under a floor than it did under an exact pin.
  gofmt_command="$(go env GOROOT)/bin/gofmt"
  if [ ! -x "$gofmt_command" ]; then
    printf 'maintainer-verify: the selected Go toolchain has no usable formatter at %s\n' \
      "$gofmt_command" >&2
    return 1
  fi

  printf 'maintainer-verify: checking ShellCheck %s or newer\n' "$minimum_shellcheck_version"
  shellcheck_version_output="$(shellcheck --version)"
  while IFS= read -r shellcheck_version_line; do
    case "$shellcheck_version_line" in
      'version: '*) shellcheck_version_text="${shellcheck_version_line#version: }" ;;
    esac
  done <<< "$shellcheck_version_output"
  if [ -z "$shellcheck_version_text" ] || \
    ! version_at_least "$shellcheck_version_text" "$minimum_shellcheck_version"; then
    printf 'maintainer-verify: ShellCheck version is %s; require %s or newer\n' \
      "${shellcheck_version_text:-unknown}" "$minimum_shellcheck_version" >&2
    return 1
  fi

  while IFS= read -r -d '' tracked_shell_file; do
    tracked_shell_files+=("$tracked_shell_file")
  done < <(git -C "$repo_root" ls-files -z -- '*.sh')
  if [ "${#tracked_shell_files[@]}" -eq 0 ]; then
    printf 'maintainer-verify: git reported no tracked shell files\n' >&2
    return 1
  fi
  printf 'maintainer-verify: ShellCheck warning-level lint (%s tracked files)\n' \
    "${#tracked_shell_files[@]}"
  (
    cd "$repo_root"
    shellcheck --severity=warning -- "${tracked_shell_files[@]}"
  )

  while IFS= read -r -d '' tracked_go_file; do
    if [ -f "$repo_root/$tracked_go_file" ]; then
      tracked_go_files+=("$tracked_go_file")
    fi
  done < <(git -C "$repo_root" ls-files -z -- '*.go')
  if [ "${#tracked_go_files[@]}" -eq 0 ]; then
    printf 'maintainer-verify: git reported no tracked Go files\n' >&2
    return 1
  fi
  printf 'maintainer-verify: gofmt formatting check (%s tracked files)\n' \
    "${#tracked_go_files[@]}"
  # gofmt exits zero even when it lists unformatted files, so the verdict is the
  # emptiness of its output, never its status.
  unformatted_go_files="$(cd "$repo_root" && "$gofmt_command" -l -- "${tracked_go_files[@]}")"
  if [ -n "$unformatted_go_files" ]; then
    printf 'maintainer-verify: tracked Go files are not gofmt-formatted:\n' >&2
    printf '%s\n' "$unformatted_go_files" >&2
    printf 'maintainer-verify: run "gofmt -w" on the files listed above\n' >&2
    return 1
  fi

  printf 'maintainer-verify: aggregate contract suite\n'
  DO_WORK_MAINTAINER_TIER="$verification_tier" \
    bash "$repo_root/_dev/tests/contract-regressions.sh"

  printf 'maintainer-verify: queue-kanban go vet\n'
  (
    cd "$repo_root/skills/do-work-board/tools/queue-kanban"
    go vet ./...
  )
  # One board test run either way. The default tier runs the package plainly: no strict
  # marker, so a run whose JavaScript probes all skipped is a legitimate green. --heavy
  # carries the marker, which makes the package's TestMain refuse a green result whose
  # probes never executed; without Node the probes skip and the run says so.
  if [ "$verification_tier" != 'heavy' ]; then
    printf 'maintainer-verify: queue-kanban uncached tests\n'
    # The selectors stay on this call so the decision child inherits them too:
    # they decide which tests run at all, and the fingerprint seals the whole
    # environment, so a decision computed without them would seal a different
    # command than the one that runs.
    (
      QUEUE_KANBAN_JAVASCRIPT_PROBES=off \
        QUEUE_KANBAN_BROWSER_PROBES=off \
        DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior \
        run_stage_with_evidence queue-kanban-fast-tests \
          "$repo_root/skills/do-work-board/tools/queue-kanban" ./...
    )
  elif node_engine_available; then
    printf 'maintainer-verify: queue-kanban uncached tests with strict JavaScript behavior probes\n'
    (
      QUEUE_KANBAN_BROWSER_PROBES=off \
        QUEUE_KANBAN_JAVASCRIPT_PROBES=on \
        QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 \
        run_budgeted_go_tests "$repo_root/skills/do-work-board/tools/queue-kanban" ./...
    )
  else
    printf '%s\n' "$node_lane_skip_line"
    printf 'maintainer-verify: queue-kanban uncached ordinary tests\n'
    (
      QUEUE_KANBAN_JAVASCRIPT_PROBES=off \
        QUEUE_KANBAN_BROWSER_PROBES=off \
        DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior \
        run_budgeted_go_tests "$repo_root/skills/do-work-board/tools/queue-kanban" ./...
    )
  fi
  # The browser behavior lane is heavy-only, and inside that tier it is guarded exactly as
  # the Node lane above is: run it when an engine is present, print an explicit SKIP naming
  # what did not run when it is not. Either way this script exits 0 — the lane's own
  # zero-probe guard is what stops a skipped run from being mistaken for a green one when it
  # IS selected. QUEUE_KANBAN_BROWSER names an engine that is not on PATH under a well-known
  # name.
  if [ "$verification_tier" = 'heavy' ]; then
    if browser_engine_available; then
      printf 'maintainer-verify: queue-kanban strict browser behavior lane\n'
      (
        QUEUE_KANBAN_JAVASCRIPT_PROBES=off \
          QUEUE_KANBAN_BROWSER_PROBES=on \
          QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 \
          run_budgeted_go_tests "$repo_root/skills/do-work-board/tools/queue-kanban" \
            -run '^TestBrowserBehavior' -v .
      )
    else
      printf '%s\n' "$browser_lane_skip_line"
    fi
  fi

  printf 'maintainer-verify: do-work-cli go vet\n'
  (
    cd "$repo_root/skills/do-work/tools/do-work-cli"
    go vet ./...
  )
  printf 'maintainer-verify: do-work-cli uncached tests\n'
  if [ "$verification_tier" = 'heavy' ]; then
    DO_WORK_HEAVY_TESTS=1 run_budgeted_go_tests "$repo_root/skills/do-work/tools/do-work-cli" ./...
  else
    run_stage_with_evidence do-work-cli-fast-tests \
      "$repo_root/skills/do-work/tools/do-work-cli" -short ./...
  fi

  # Whole-gate wall time was recorded nowhere before this line, and requirement 5
  # of REQ-591 asks for it. SECONDS counts from this shell's start, which is this
  # script's start, so the subtraction needs no extra process.
  printf 'maintainer-verify: gate wall %ss\n' "$((SECONDS - gate_started_seconds))"
  printf 'Maintainer verification passed.\n'
}

run_heavy_probe_lane() (
  local probe_script="$1"
  local lane_scratch
  local lane_cli_binary="${DO_WORK_TEST_DO_WORK_CLI_BINARY:-}"
  local started_at
  local elapsed_seconds
  local lane_status=0
  # Re-source inside the lane subshell so an isolated caller's log/run overrides
  # are the ones this row records.
  # shellcheck source=_dev/tests/test-duration-log.sh
  source "$repo_root/_dev/tests/test-duration-log.sh"
  if [ -z "$lane_cli_binary" ]; then
    lane_scratch="$(mktemp -d)"
    trap 'rm -rf "$lane_scratch"' EXIT
    lane_cli_binary="$lane_scratch/do-work-cli"
    if ! (
      cd "$repo_root/skills/do-work/tools/do-work-cli"
      go build -ldflags='-s -w' -o "$lane_cli_binary" ./cmd/do-work-cli
    ); then
      printf 'maintainer-verify: could not build the heavy-lane do-work-cli\n' >&2
      return 1
    fi
  elif [ ! -x "$lane_cli_binary" ]; then
    printf 'maintainer-verify: supplied heavy-lane do-work-cli is not executable: %s\n' "$lane_cli_binary" >&2
    return 1
  fi
  started_at="$(date +%s)"
  DO_WORK_MAINTAINER_TIER=heavy \
    DO_WORK_TEST_DO_WORK_CLI_BINARY="$lane_cli_binary" \
    bash "$probe_script" || lane_status=$?
  elapsed_seconds=$(( $(date +%s) - started_at ))
  printf 'test-file duration: %s %ss (limit none (heavy))\n' \
    "${probe_script#"$repo_root/"}" "$elapsed_seconds"
  record_test_duration "${probe_script#"$repo_root/"}" "$elapsed_seconds"
  return "$lane_status"
)

# Each stable lane is independently executable. The typed planner emits only these
# argv values; --heavy remains the explicit force-all gate above.
run_heavy_lane() {
  local verification_tier=heavy
  local lane_id="$1"

  case "$lane_id" in
    queue-kanban-javascript)
      # An isolated lane announces "did not run" exactly as the whole heavy tier does.
      # Without this the strict Go lane ran, skipped inside, and tripped its own
      # zero-probe guard: a red exit and no SKIP line for a caller to record.
      if ! node_engine_available; then
        printf '%s\n' "$node_lane_skip_line"
        return 0
      fi
      QUEUE_KANBAN_BROWSER_PROBES=off \
        QUEUE_KANBAN_JAVASCRIPT_PROBES=on \
        QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 \
        run_budgeted_go_tests "$repo_root/skills/do-work-board/tools/queue-kanban" \
          -run '^TestJavaScriptBehavior' -v .
      ;;
    queue-kanban-browser)
      if ! browser_engine_available; then
        printf '%s\n' "$browser_lane_skip_line"
        return 0
      fi
      QUEUE_KANBAN_JAVASCRIPT_PROBES=off \
        QUEUE_KANBAN_BROWSER_PROBES=on \
        QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 \
        run_budgeted_go_tests "$repo_root/skills/do-work-board/tools/queue-kanban" \
          -run '^TestBrowserBehavior' -v .
      ;;
    do-work-cli-integrations)
      DO_WORK_HEAVY_TESTS=1 \
        run_budgeted_go_tests "$repo_root/skills/do-work/tools/do-work-cli" ./...
      ;;
    staged-skills)
      run_heavy_probe_lane "$repo_root/_dev/tests/staged-skills-contract.sh"
      ;;
    updater)
      run_heavy_probe_lane "$repo_root/_dev/tests/update-script-behavior.sh"
      ;;
    installer)
      run_heavy_probe_lane "$repo_root/_dev/tests/install-suite-behavior.sh"
      ;;
    *)
      printf 'maintainer-verify: unknown heavy lane: %s\n' "$lane_id" >&2
      return 2
      ;;
  esac
}

if [ "$#" -gt 2 ]; then
  print_usage_and_exit
fi
case "${1:-}" in
  --self-test)
    [ "$#" -eq 1 ] || print_usage_and_exit
    run_self_test
    ;;
  --heavy)
    [ "$#" -eq 1 ] || print_usage_and_exit
    run_verification heavy
    ;;
  --heavy-lane)
    [ "$#" -eq 2 ] || print_usage_and_exit
    run_heavy_lane "$2"
    ;;
  '')
    [ "$#" -eq 0 ] || print_usage_and_exit
    run_verification fast
    ;;
  *) print_usage_and_exit ;;
esac

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
script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"

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
        elif [ "$#" -eq 3 ] && [ "$1" = 'test' ] && \
          [ "$2" = '-count=1' ] && [ "$3" = './...' ]; then
          record_stage 'board-test'
        elif [ "$#" -eq 6 ] && [ "$1" = 'test' ] && \
          [ "$2" = '-count=1' ] && [ "$3" = '-run' ] && \
          [ "$4" = '^TestMaintainerStrictJavaScriptBehaviorLane$' ] && \
          [ "$5" = '-v' ] && [ "$6" = '.' ]; then
          record_stage 'board-strict'
        else
          exit 64
        fi
        ;;
      */skills/do-work-toolbox/tools/audit-metrics)
        if [ "$#" -eq 2 ] && [ "$1" = 'vet' ] && [ "$2" = './...' ]; then
          record_stage 'audit-vet'
        elif [ "$#" -eq 3 ] && [ "$1" = 'test' ] && \
          [ "$2" = '-count=1' ] && [ "$3" = './...' ]; then
          record_stage 'audit-test'
        else
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
  local stage_line

  if [ "$expect_strict" = 'yes' ]; then
    expected_count=10
  fi
  for expected_stage in \
    go-version shellcheck-version shellcheck-lint gofmt-lint aggregate \
    board-vet board-test audit-vet audit-test; do
    actual_count="$(count_stage_occurrences "$expected_stage" "$log_path")"
    if [ "$actual_count" -ne 1 ]; then
      fail_self_test "$expected_stage ran $actual_count times; want exactly once"
      return 1
    fi
  done
  actual_count="$(count_stage_occurrences 'board-strict' "$log_path")"
  if { [ "$expect_strict" = 'yes' ] && [ "$actual_count" -ne 1 ]; } || \
    { [ "$expect_strict" = 'no' ] && [ "$actual_count" -ne 0 ]; }; then
    fail_self_test "board-strict ran $actual_count times with expect_strict=$expect_strict"
    return 1
  fi
  while IFS= read -r stage_line; do
    total_count=$((total_count + 1))
  done < "$log_path"
  if [ "$total_count" -ne "$expected_count" ]; then
    fail_self_test "success run recorded $total_count stages; want $expected_count"
    return 1
  fi
}

run_self_test() {
  local self_test_root
  local fixture_root
  local fixture_script
  local fixture_go_root
  local with_node_bin
  local without_node_bin
  local generic_shim
  local stage_name
  local success_log
  local success_output
  local no_node_log
  local no_node_output
  local newer_tools_log
  local newer_tools_output
  local failure_log
  local failure_output
  local mutated_fixture_script
  local mutation_output

  self_test_root="$(mktemp -d)"
  trap 'rm -rf -- "$self_test_root"' EXIT
  fixture_root="$self_test_root/repository"
  fixture_script="$fixture_root/_dev/tests/maintainer-verify.sh"
  fixture_go_root="$self_test_root/go-root"
  with_node_bin="$self_test_root/with-node-bin"
  without_node_bin="$self_test_root/without-node-bin"
  mkdir -p \
    "$fixture_root/_dev/tests" \
    "$fixture_root/skills/do-work-board/tools/queue-kanban" \
    "$fixture_root/skills/do-work-toolbox/tools/audit-metrics" \
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

  success_log="$self_test_root/success.log"
  success_output="$self_test_root/success.out"
  : > "$success_log"
  if ! PATH="$with_node_bin" \
    QUEUE_KANBAN_BROWSER='' \
    MAINTAINER_VERIFY_SELFTEST_LOG="$success_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
    MAINTAINER_VERIFY_SELFTEST_GOROOT="$fixture_go_root" \
    /bin/bash "$fixture_script" > "$success_output" 2>&1; then
    sed 's/^/  /' "$success_output" >&2
    fail_self_test 'the all-success fixture exited nonzero'
    return 1
  fi
  assert_success_stages "$success_log" yes

  no_node_log="$self_test_root/no-node.log"
  no_node_output="$self_test_root/no-node.out"
  : > "$no_node_log"
  if ! PATH="$without_node_bin" \
    QUEUE_KANBAN_BROWSER='' \
    MAINTAINER_VERIFY_SELFTEST_LOG="$no_node_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
    MAINTAINER_VERIFY_SELFTEST_GOROOT="$fixture_go_root" \
    /bin/bash "$fixture_script" > "$no_node_output" 2>&1; then
    sed 's/^/  /' "$no_node_output" >&2
    fail_self_test 'the no-Node fixture exited nonzero'
    return 1
  fi
  assert_success_stages "$no_node_log" no
  if ! grep -q 'SKIP: Node is unavailable; strict JavaScript behavior lane was not run.' "$no_node_output"; then
    fail_self_test 'the no-Node success path did not print its explicit skip'
    return 1
  fi

  newer_tools_log="$self_test_root/newer-tools.log"
  newer_tools_output="$self_test_root/newer-tools.out"
  : > "$newer_tools_log"
  if ! PATH="$with_node_bin" \
    QUEUE_KANBAN_BROWSER='' \
    MAINTAINER_VERIFY_SELFTEST_LOG="$newer_tools_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
    MAINTAINER_VERIFY_SELFTEST_GOROOT="$fixture_go_root" \
    MAINTAINER_VERIFY_SELFTEST_NEWER_TOOLS=yes \
    /bin/bash "$fixture_script" > "$newer_tools_output" 2>&1; then
    sed 's/^/  /' "$newer_tools_output" >&2
    fail_self_test 'a toolchain newer than the floor was rejected — these gates are floors, not exact pins'
    return 1
  fi
  assert_success_stages "$newer_tools_log" yes

  for stage_name in \
    go-version shellcheck-version shellcheck-lint gofmt-lint gofmt-unformatted \
    aggregate board-vet board-test board-strict audit-vet audit-test; do
    failure_log="$self_test_root/failure-$stage_name.log"
    failure_output="$self_test_root/failure-$stage_name.out"
    : > "$failure_log"
    if PATH="$with_node_bin" \
      MAINTAINER_VERIFY_SELFTEST_LOG="$failure_log" \
      MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
      MAINTAINER_VERIFY_SELFTEST_GOROOT="$fixture_go_root" \
      MAINTAINER_VERIFY_FAIL_STAGE="$stage_name" \
      /bin/bash "$fixture_script" > "$failure_output" 2>&1; then
      fail_self_test "$stage_name failure exited zero"
      return 1
    fi
    if [ "$stage_name" = 'gofmt-unformatted' ] && \
      ! grep -q '_dev/tests/fixture-go.go' "$failure_output"; then
      fail_self_test 'the unformatted-Go failure did not name the offending file'
      return 1
    fi
  done

  mutated_fixture_script="$self_test_root/maintainer-verify-mutated.sh"
  sed '/^run_verification()/,$ s/TestMaintainerStrictJavaScriptBehaviorLane/TestDefinitelyWrongStrictLane/' \
    "$fixture_script" > "$mutated_fixture_script"
  if cmp -s "$fixture_script" "$mutated_fixture_script"; then
    fail_self_test 'strict-regex mutation did not alter the production command'
    return 1
  fi
  mv "$mutated_fixture_script" "$fixture_script"
  chmod +x "$fixture_script"
  mutation_output="$self_test_root/wrong-strict-regex.out"
  if /bin/bash "$fixture_script" --self-test > "$mutation_output" 2>&1; then
    fail_self_test 'a wrong strict JavaScript test regex left --self-test green'
    return 1
  fi
  cp "$script_path" "$fixture_script"
  chmod +x "$fixture_script"
  if ! cmp -s "$script_path" "$fixture_script"; then
    fail_self_test 'the strict-regex mutation fixture was not restored'
    return 1
  fi

  rm -rf -- "$self_test_root"
  trap - EXIT
  printf 'Maintainer verification self-test passed.\n'
}

run_verification() {
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
    tracked_go_files+=("$tracked_go_file")
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
  bash "$repo_root/_dev/tests/contract-regressions.sh"

  printf 'maintainer-verify: queue-kanban go vet\n'
  (
    cd "$repo_root/skills/do-work-board/tools/queue-kanban"
    go vet ./...
  )
  printf 'maintainer-verify: queue-kanban uncached ordinary tests\n'
  (
    cd "$repo_root/skills/do-work-board/tools/queue-kanban"
    go test -count=1 ./...
  )
  if command -v node >/dev/null 2>&1; then
    printf 'maintainer-verify: queue-kanban strict JavaScript behavior lane\n'
    (
      cd "$repo_root/skills/do-work-board/tools/queue-kanban"
      go test -count=1 -run '^TestMaintainerStrictJavaScriptBehaviorLane$' -v .
    )
  else
    printf 'SKIP: Node is unavailable; strict JavaScript behavior lane was not run.\n'
  fi
  # The browser behavior lane, guarded exactly as the Node lane above is: run it when
  # an engine is present, print an explicit SKIP naming what did not run when it is
  # not. Either way this script exits 0 — the lane's own zero-probe guard is what
  # stops a skipped run from being mistaken for a green one when it IS selected.
  # QUEUE_KANBAN_BROWSER names an engine that is not on PATH under a well-known name.
  browser_probe_binary=""
  if [ -n "${QUEUE_KANBAN_BROWSER:-}" ]; then
    browser_probe_binary="$QUEUE_KANBAN_BROWSER"
  else
    for browser_probe_candidate in google-chrome google-chrome-stable chromium chromium-browser chrome; do
      if command -v "$browser_probe_candidate" >/dev/null 2>&1; then
        browser_probe_binary="$browser_probe_candidate"
        break
      fi
    done
  fi
  if [ -n "$browser_probe_binary" ]; then
    printf 'maintainer-verify: queue-kanban strict browser behavior lane\n'
    (
      cd "$repo_root/skills/do-work-board/tools/queue-kanban"
      go test -count=1 -run '^TestMaintainerStrictBrowserBehaviorLane$' -v .
    )
  else
    printf 'SKIP: no browser is available; strict browser behavior lane was not run. Set QUEUE_KANBAN_BROWSER to name one.\n'
  fi

  printf 'maintainer-verify: audit-metrics go vet\n'
  (
    cd "$repo_root/skills/do-work-toolbox/tools/audit-metrics"
    go vet ./...
  )
  printf 'maintainer-verify: audit-metrics uncached tests\n'
  (
    cd "$repo_root/skills/do-work-toolbox/tools/audit-metrics"
    go test -count=1 ./...
  )

  printf 'Maintainer verification passed.\n'
}

case "${1:-}" in
  --self-test)
    if [ "$#" -ne 1 ]; then
      printf 'usage: %s [--self-test]\n' "$0" >&2
      exit 2
    fi
    run_self_test
    ;;
  '')
    run_verification
    ;;
  *)
    printf 'usage: %s [--self-test]\n' "$0" >&2
    exit 2
    ;;
esac

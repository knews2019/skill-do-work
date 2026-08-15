#!/usr/bin/env bash
set -euo pipefail

required_go_version="go1.26.1"
required_shellcheck_version="0.11.0"
script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"

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
        printf '%s\n' 'go version go1.99.0 fixture/arch'
      else
        printf '%s\n' 'go version go1.26.1 fixture/arch'
      fi
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
  shellcheck)
    if [ "${1:-}" = '--version' ]; then
      if [ "$#" -ne 1 ]; then
        exit 64
      fi
      printf '%s\n' 'shellcheck-version' >> "$stage_log"
      printf '%s\n' 'ShellCheck - shell script analysis tool'
      if [ "$failure_stage" = 'shellcheck-version' ]; then
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
      [ "$3" != 'ls-files' ] || [ "$4" != '-z' ] || \
      [ "$5" != '--' ] || [ "$6" != '*.sh' ]; then
      exit 64
    fi
    printf '%s\0' '_dev/tests/fixture-shell.sh'
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
  local expected_count=8
  local stage_line

  if [ "$expect_strict" = 'yes' ]; then
    expected_count=9
  fi
  for expected_stage in \
    go-version shellcheck-version shellcheck-lint aggregate \
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
  local with_node_bin
  local without_node_bin
  local generic_shim
  local stage_name
  local success_log
  local success_output
  local no_node_log
  local no_node_output
  local failure_log
  local failure_output
  local mutated_fixture_script
  local mutation_output

  self_test_root="$(mktemp -d)"
  trap 'rm -rf -- "$self_test_root"' EXIT
  fixture_root="$self_test_root/repository"
  fixture_script="$fixture_root/_dev/tests/maintainer-verify.sh"
  with_node_bin="$self_test_root/with-node-bin"
  without_node_bin="$self_test_root/without-node-bin"
  mkdir -p \
    "$fixture_root/_dev/tests" \
    "$fixture_root/skills/do-work-board/tools/queue-kanban" \
    "$fixture_root/skills/do-work-toolbox/tools/audit-metrics" \
    "$with_node_bin" "$without_node_bin"
  cp "$script_path" "$fixture_script"
  chmod +x "$fixture_script"
  printf '%s\n' '#!/usr/bin/env bash' > "$fixture_root/_dev/tests/fixture-shell.sh"
  printf '%s\n' '#!/usr/bin/env bash' > "$fixture_root/_dev/tests/contract-regressions.sh"

  generic_shim="$self_test_root/command-shim"
  write_command_shim "$generic_shim"
  for stage_name in go shellcheck git bash node; do
    ln -s "$generic_shim" "$with_node_bin/$stage_name"
  done
  for stage_name in go shellcheck git bash; do
    ln -s "$generic_shim" "$without_node_bin/$stage_name"
  done

  success_log="$self_test_root/success.log"
  success_output="$self_test_root/success.out"
  : > "$success_log"
  if ! PATH="$with_node_bin" \
    MAINTAINER_VERIFY_SELFTEST_LOG="$success_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
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
    MAINTAINER_VERIFY_SELFTEST_LOG="$no_node_log" \
    MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
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

  for stage_name in \
    go-version shellcheck-version shellcheck-lint aggregate \
    board-vet board-test board-strict audit-vet audit-test; do
    failure_log="$self_test_root/failure-$stage_name.log"
    failure_output="$self_test_root/failure-$stage_name.out"
    : > "$failure_log"
    if PATH="$with_node_bin" \
      MAINTAINER_VERIFY_SELFTEST_LOG="$failure_log" \
      MAINTAINER_VERIFY_EXPECTED_REPO_ROOT="$fixture_root" \
      MAINTAINER_VERIFY_FAIL_STAGE="$stage_name" \
      /bin/bash "$fixture_script" > "$failure_output" 2>&1; then
      fail_self_test "$stage_name failure exited zero"
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

  for required_command in git go shellcheck bash; do
    if ! command -v "$required_command" >/dev/null 2>&1; then
      printf 'maintainer-verify: required command is unavailable: %s\n' "$required_command" >&2
      return 1
    fi
  done

  printf 'maintainer-verify: checking Go %s\n' "$required_go_version"
  go_version_output="$(go version)"
  read -r _ _ go_version_text _ <<< "$go_version_output"
  if [ "$go_version_text" != "$required_go_version" ]; then
    printf 'maintainer-verify: Go version is %s; require exactly %s\n' \
      "$go_version_text" "$required_go_version" >&2
    return 1
  fi

  printf 'maintainer-verify: checking ShellCheck %s\n' "$required_shellcheck_version"
  shellcheck_version_output="$(shellcheck --version)"
  while IFS= read -r shellcheck_version_line; do
    case "$shellcheck_version_line" in
      'version: '*) shellcheck_version_text="${shellcheck_version_line#version: }" ;;
    esac
  done <<< "$shellcheck_version_output"
  if [ "$shellcheck_version_text" != "$required_shellcheck_version" ]; then
    printf 'maintainer-verify: ShellCheck version is %s; require exactly %s\n' \
      "${shellcheck_version_text:-unknown}" "$required_shellcheck_version" >&2
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

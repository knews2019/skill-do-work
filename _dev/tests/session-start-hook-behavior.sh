#!/usr/bin/env bash
# Behavioral regression probes for the core SessionStart status hook.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
hook_source="$repo_root/skills/do-work/hooks/session-start.sh"
fixture_root="$(mktemp -d)" || {
  printf 'FAIL: could not allocate SessionStart hook fixture root.\n' >&2
  exit 1
}
trap 'rm -rf "$fixture_root"' EXIT

failure_count=0

run_hook_case() {
  local case_name="$1"
  local version_mode="$2"
  local queue_count="$3"
  local expected_output="$4"
  local case_root="$fixture_root/$case_name"
  local skill_root="$case_root/skill"
  local project_root="$case_root/project"
  local stderr_path="$case_root/stderr.txt"
  local actual_output=''
  local hook_status=0
  local request_index=0

  mkdir -p "$skill_root/hooks" "$skill_root/actions" "$project_root"
  cp "$hook_source" "$skill_root/hooks/session-start.sh"
  case "$version_mode" in
    valid) printf '**Current version**: 9.8.7\n' > "$skill_root/actions/version.md" ;;
    reformatted) printf '**Version**: 9.8.7\n' > "$skill_root/actions/version.md" ;;
    missing) ;;
    *)
      printf 'FAIL: %s fixture has unknown version mode %s.\n' "$case_name" "$version_mode" >&2
      failure_count=$((failure_count + 1))
      return
      ;;
  esac

  if [[ "$queue_count" -gt 0 ]]; then
    mkdir -p "$project_root/do-work/queue"
    while [[ "$request_index" -lt "$queue_count" ]]; do
      request_index=$((request_index + 1))
      printf -- '---\nid: REQ-%03d\nstatus: pending\n---\n' "$request_index" \
        > "$project_root/do-work/queue/REQ-$(printf '%03d' "$request_index")-fixture.md"
    done
  fi

  actual_output="$(CLAUDE_PROJECT_DIR="$project_root" \
    bash "$skill_root/hooks/session-start.sh" 2> "$stderr_path")" || hook_status=$?
  if [[ "$hook_status" -ne 0 ]]; then
    printf 'FAIL: %s: hook exited %s; stderr: %s\n' \
      "$case_name" "$hook_status" "$(tr '\n' ' ' < "$stderr_path")" >&2
    failure_count=$((failure_count + 1))
    return
  fi
  if [[ "$actual_output" != "$expected_output" ]]; then
    printf 'FAIL: %s: expected <%s>, got <%s>.\n' \
      "$case_name" "$expected_output" "$actual_output" >&2
    failure_count=$((failure_count + 1))
  fi
}

run_hook_case \
  happy-path valid 2 \
  "do-work v9.8.7 loaded. 2 pending REQ(s). Say 'do-work help' for commands."
run_hook_case \
  missing-version-file missing 0 \
  "do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands."
run_hook_case \
  reformatted-version-line reformatted 0 \
  "do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands."
run_hook_case \
  missing-queue-directory valid 0 \
  "do-work v9.8.7 loaded. 0 pending REQ(s). Say 'do-work help' for commands."

# With the cleanup script installed, a marker made redundant by its landed REQ file
# is reaped at session start and the hook appends the cleanup summary to the banner.
# The four cases above run without scripts/ present, pinning that a partial install
# still emits the plain banner.
cleanup_case_root="$fixture_root/reservation-cleanup"
cleanup_skill_root="$cleanup_case_root/skill"
cleanup_project_root="$cleanup_case_root/project"
mkdir -p "$cleanup_skill_root/hooks" "$cleanup_skill_root/actions" "$cleanup_skill_root/scripts" \
  "$cleanup_project_root/do-work/queue" "$cleanup_project_root/do-work/.req-reservations"
cp "$hook_source" "$cleanup_skill_root/hooks/session-start.sh"
cp "$repo_root/skills/do-work/scripts/cleanup-req-reservations.sh" "$cleanup_skill_root/scripts/"
printf '**Current version**: 9.8.7\n' > "$cleanup_skill_root/actions/version.md"
printf -- '---\nid: REQ-001\nstatus: pending\n---\n' \
  > "$cleanup_project_root/do-work/queue/REQ-001-fixture.md"
: > "$cleanup_project_root/do-work/.req-reservations/REQ-000001"
cleanup_expected_output="do-work v9.8.7 loaded. 1 pending REQ(s). Say 'do-work help' for commands.
do-work: removed 1 stale REQ reservation marker(s) from do-work/.req-reservations/ — stage and commit the deletion(s)."
cleanup_actual_output="$(CLAUDE_PROJECT_DIR="$cleanup_project_root" \
  bash "$cleanup_skill_root/hooks/session-start.sh" 2> "$cleanup_case_root/stderr.txt")" || {
  printf 'FAIL: reservation-cleanup: hook exited nonzero; stderr: %s\n' \
    "$(tr '\n' ' ' < "$cleanup_case_root/stderr.txt")" >&2
  failure_count=$((failure_count + 1))
}
if [[ "$cleanup_actual_output" != "$cleanup_expected_output" ]]; then
  printf 'FAIL: reservation-cleanup: expected <%s>, got <%s>.\n' \
    "$cleanup_expected_output" "$cleanup_actual_output" >&2
  failure_count=$((failure_count + 1))
fi
if [[ -e "$cleanup_project_root/do-work/.req-reservations/REQ-000001" ]]; then
  printf 'FAIL: reservation-cleanup: the redundant marker survived the session start.\n' >&2
  failure_count=$((failure_count + 1))
fi

if [[ "$failure_count" -gt 0 ]]; then
  exit 1
fi

printf 'SessionStart hook behavior probes passed.\n'

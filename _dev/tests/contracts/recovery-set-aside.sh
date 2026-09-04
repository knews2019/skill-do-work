#!/usr/bin/env bash
# Recovery reads finalization records per REQ (REQ-515). The pins here are the
# per-record wording that replaced the whole-run gate; REQ-456's stuck
# finalization tail parked 31 pending REQs while that gate stood.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
core_actions="$repo_root/skills/do-work/actions"
fail_count=0

# The token is the command's, not the prose's. Read it from the Go constant so a
# rename in the CLI reddens here instead of leaving the actions naming a code
# nothing emits.
set_aside_code="$(sed -n 's/^const SetAsideReasonCode = "\(.*\)"$/\1/p' \
  "$repo_root/skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go")"
if [ -z "$set_aside_code" ]; then
  printf 'FAIL: internal/finalization/finalization_commands.go no longer declares SetAsideReasonCode; the actions below cite a code nothing emits.\n' >&2
  fail_count=$((fail_count + 1))
fi

# The retired gate. Its exact shape is what made one REQ's refused record the
# whole run's stop, so its return must redden even if the surrounding paragraph
# is rewritten.
if grep -qF 'Continue only on typed success whose terminal phase is' "$core_actions/work-reference.md"; then
  printf 'FAIL: actions/work-reference.md restored the whole-run finalization gate; recovery reads finalizations one record at a time (REQ-515).\n' >&2
  fail_count=$((fail_count + 1))
fi

require_action_phrase() {
  local action_file="$1"
  local required_phrase="$2"
  local why="$3"
  if [ ! -f "$core_actions/$action_file" ]; then
    printf 'FAIL: actions/%s is missing; %s\n' "$action_file" "$why" >&2
    fail_count=$((fail_count + 1))
    return
  fi
  if ! grep -qF -- "$required_phrase" "$core_actions/$action_file"; then
    printf 'FAIL: actions/%s must state %s — missing: %s\n' "$action_file" "$why" "$required_phrase" >&2
    fail_count=$((fail_count + 1))
  fi
}

# Every entry point into the loop reads the same per-record contract, so serial
# `run` and the sole-authority verb cannot drift into different stop rules.
if [ -n "$set_aside_code" ]; then
  for recovery_action in work.md work-reference.md run-with-recovery.md; do
    require_action_phrase "$recovery_action" "$set_aside_code" \
      'which reason code marks a REQ the run set aside'
  done
fi
require_action_phrase work.md 'never as one pass/fail gate for the run' \
  'that a refused record excludes its own REQ rather than stopping the loop'
require_action_phrase work-reference.md 'one record at a time' \
  'that finalization records are read per REQ'
require_action_phrase work-reference.md 'Set-aside-by-recovery section' \
  'that the composed exit summary renders the set-aside REQs'
require_action_phrase work-reference.md 'reason codes, comma-separated' \
  "each set-aside REQ's own reason codes in the exit summary"
require_action_phrase work-reference.md 'recover: <resolving verb>' \
  'a resolving verb beside every set-aside REQ in the exit summary'
require_action_phrase work-reference.md 'its finding names no REQ' \
  'that the surviving whole-run stop is the finding no REQ owns'
require_action_phrase run-with-recovery.md 'one record at a time' \
  'the same per-record reading the ordinary run uses'

[ "$fail_count" -eq 0 ] || exit 1
printf 'recovery set-aside contract probes passed.\n'

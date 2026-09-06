#!/usr/bin/env bash
# Fixture execution proofs for cleanup-req-reservations.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# cleanup-req-reservations: a marker whose REQ file is committed is removed at
# any zero-padding width and archive depth; non-marker entries stay untouched.
reservation_project="$fixture_root/reservation-project"
fixture_repo_init "$reservation_project"
mkdir -p "$reservation_project/do-work/.req-reservations" \
  "$reservation_project/do-work/queue" \
  "$reservation_project/do-work/archive/UR-001"
printf 'req\n' > "$reservation_project/do-work/queue/REQ-203-fixture.md"
printf 'req\n' > "$reservation_project/do-work/archive/UR-001/REQ-7-archived-fixture.md"
: > "$reservation_project/do-work/.req-reservations/REQ-000203"
: > "$reservation_project/do-work/.req-reservations/REQ-000007"
: > "$reservation_project/do-work/.req-reservations/README"
fixture_repo_commit_all "$reservation_project" 'captured fixtures'
ln -s REQ-000203 "$reservation_project/do-work/.req-reservations/REQ-000777"
reservation_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations redundant-marker case returned nonzero'
grep -q 'removed 2 stale REQ reservation marker' <<<"$reservation_output" \
  || fail_case 'cleanup-req-reservations redundant-marker case did not report exactly two removals'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000203" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case kept the queue-claimed marker'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000007" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case kept the archive-claimed marker'
[ -e "$reservation_project/do-work/.req-reservations/README" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case deleted a non-marker file'
[ -L "$reservation_project/do-work/.req-reservations/REQ-000777" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case deleted a symlinked marker'

# cleanup-req-reservations: in a git work tree, a REQ file merely present on
# disk is a capture still staging — its marker must survive until the capture
# commits, then be reaped. Deleting early breaks capture's prescribed
# `git add do-work/.req-reservations/REQ-NNNNNN`.
printf 'req\n' > "$reservation_project/do-work/queue/REQ-500-inflight-fixture.md"
: > "$reservation_project/do-work/.req-reservations/REQ-000500"
reservation_uncommitted_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations uncommitted-capture case returned nonzero'
[ -z "$reservation_uncommitted_output" ] \
  || fail_case 'cleanup-req-reservations uncommitted-capture case printed output while the capture was mid-flight'
[ -f "$reservation_project/do-work/.req-reservations/REQ-000500" ] \
  || fail_case 'cleanup-req-reservations uncommitted-capture case deleted a mid-capture marker'
fixture_repo_commit_all "$reservation_project" 'captured REQ-500'
reservation_committed_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations committed-capture case returned nonzero'
grep -q 'removed 1 stale REQ reservation marker' <<<"$reservation_committed_output" \
  || fail_case 'cleanup-req-reservations committed-capture case did not reap the landed marker'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000500" ] \
  || fail_case 'cleanup-req-reservations committed-capture case kept the landed marker'

# cleanup-req-reservations: a young marker with no REQ file is a capture still in
# flight — it must survive, and a run that removes nothing must print nothing.
: > "$reservation_project/do-work/.req-reservations/REQ-000999"
reservation_inflight_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations in-flight-marker case returned nonzero'
[ -z "$reservation_inflight_output" ] \
  || fail_case 'cleanup-req-reservations in-flight-marker case printed output for a no-op run'
[ -f "$reservation_project/do-work/.req-reservations/REQ-000999" ] \
  || fail_case 'cleanup-req-reservations in-flight-marker case deleted a young unmatched marker'

# cleanup-req-reservations: the same unmatched marker aged past two days is an
# abandoned capture and is removed.
touch -m -t 202001010000 "$reservation_project/do-work/.req-reservations/REQ-000999"
reservation_timeout_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations timeout-marker case returned nonzero'
grep -q 'removed 1 stale REQ reservation marker' <<<"$reservation_timeout_output" \
  || fail_case 'cleanup-req-reservations timeout-marker case did not report the timeout removal'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000999" ] \
  || fail_case 'cleanup-req-reservations timeout-marker case kept the abandoned marker'

# cleanup-req-reservations: a repo without a reservation directory is a silent no-op.
reservation_absent_project="$fixture_root/reservation-absent-project"
mkdir -p "$reservation_absent_project/do-work/queue"
reservation_absent_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_absent_project")" \
  || fail_case 'cleanup-req-reservations missing-directory case returned nonzero'
[ -z "$reservation_absent_output" ] \
  || fail_case 'cleanup-req-reservations missing-directory case printed output with nothing to clean'

# cleanup-req-reservations: a symlinked reservation directory is refused, so the
# automatic hook can never delete files outside the project through the link —
# a regular child reached through a symlinked parent passes a per-file -L check.
reservation_symlink_project="$fixture_root/reservation-symlink-project"
reservation_external_store="$fixture_root/reservation-external-store"
mkdir -p "$reservation_symlink_project/do-work/queue" "$reservation_external_store"
printf 'req\n' > "$reservation_symlink_project/do-work/queue/REQ-42-fixture.md"
: > "$reservation_external_store/REQ-000042"
touch -m -t 202001010000 "$reservation_external_store/REQ-000042"
ln -s "$reservation_external_store" "$reservation_symlink_project/do-work/.req-reservations"
reservation_symlink_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_symlink_project")" \
  || fail_case 'cleanup-req-reservations symlinked-directory case returned nonzero'
[ -z "$reservation_symlink_output" ] \
  || fail_case 'cleanup-req-reservations symlinked-directory case printed output for a refused store'
[ -f "$reservation_external_store/REQ-000042" ] \
  || fail_case 'cleanup-req-reservations symlinked-directory case deleted through the symlinked store'

prescribed_shell_finish

#!/usr/bin/env bash
# Fixture execution proofs for capture-screenshot.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# capture-screenshot: a destination occupied by a DIRECTORY must fail closed and keep
# the staged source. `ln` refuses an occupied FILE — which is where the no-clobber
# guarantee comes from — but nests on a directory and exits zero, and under --staged
# that zero is read as permission to delete the staged source. The dispatch holds the
# only copy, so a false success there destroys it.
capture_occupied_root="$fixture_root/capture-occupied"
mkdir -p "$capture_occupied_root/stage" "$capture_occupied_root/assets/result.png"
printf 'the only copy' > "$capture_occupied_root/stage/source.png"
printf 'occupant\n' > "$capture_occupied_root/assets/result.png/pre-existing.txt"
"$core_scripts/capture-screenshot.sh" --staged \
  "$capture_occupied_root/stage/source.png" "$capture_occupied_root/assets/result.png" >/dev/null 2>&1 \
  && fail_case 'capture-screenshot occupied-destination case reported success for a publication that nested'
[ -f "$capture_occupied_root/stage/source.png" ] \
  || fail_case 'capture-screenshot occupied-destination case destroyed the staged source it never published'
[ -d "$capture_occupied_root/assets/result.png" ] \
  || fail_case 'capture-screenshot occupied-destination case did not leave the occupying directory in place'
[ "$(cat "$capture_occupied_root/assets/result.png/pre-existing.txt")" = occupant ] \
  || fail_case 'capture-screenshot occupied-destination case disturbed the occupying directory contents'
leaked_private_paths="$(find "$capture_occupied_root/assets/result.png" -name '*.copying.*' -print -quit)" \
  || fail_case 'capture-screenshot occupied-destination case could not search the occupying directory'
[ -n "$leaked_private_paths" ] \
  && fail_case 'capture-screenshot occupied-destination case abandoned its private copy inside the occupant'

# capture-screenshot: coordinate two writers so the loser cannot publish the winner's private copy.
# Both children perform one finite local-file publication and return; neither waits
# for a later kill. The timeout below diagnoses a regression in that operation.
capture_root="$fixture_root/capture"
mkdir -p "$capture_root/a" "$capture_root/b" "$capture_root/assets"
printf 'dispatch-a' > "$capture_root/a/source.png"
printf 'dispatch-b' > "$capture_root/b/source.png"
capture_destination="$capture_root/assets/result.png"
(
  "$core_scripts/capture-screenshot.sh" --staged "$capture_root/a/source.png" "$capture_destination"
) >"$capture_root/a.out" 2>"$capture_root/a.err" & race_a_pid=$!
(
  "$core_scripts/capture-screenshot.sh" --staged "$capture_root/b/source.png" "$capture_destination"
) >"$capture_root/b.out" 2>"$capture_root/b.err" & race_b_pid=$!
race_wait_ticks=0
while { kill -0 "$race_a_pid" 2>/dev/null || kill -0 "$race_b_pid" 2>/dev/null; } \
  && [ "$race_wait_ticks" -lt 1000 ]; do
  sleep 0.01
  race_wait_ticks=$((race_wait_ticks + 1))
done
if kill -0 "$race_a_pid" 2>/dev/null || kill -0 "$race_b_pid" 2>/dev/null; then
  kill "$race_a_pid" "$race_b_pid" 2>/dev/null || true
  wait "$race_a_pid" 2>/dev/null || true
  wait "$race_b_pid" 2>/dev/null || true
  fail_case 'capture-screenshot coordinated-race case timed out'
  race_a_status=1
  race_b_status=1
else
  wait "$race_a_pid"; race_a_status=$?
  wait "$race_b_pid"; race_b_status=$?
fi
if [ "$race_a_status" -eq 0 ] && [ "$race_b_status" -ne 0 ]; then
  capture_winner=a
  capture_loser=b
elif [ "$race_b_status" -eq 0 ] && [ "$race_a_status" -ne 0 ]; then
  capture_winner=b
  capture_loser=a
else
  capture_winner=''
  capture_loser=''
  fail_case 'capture-screenshot coordinated-race case did not install exactly one writer'
fi
[ -n "$capture_winner" ] && [ "$(cat "$capture_destination" 2>/dev/null)" = "dispatch-$capture_winner" ] \
  || fail_case 'capture-screenshot coordinated-race case published cross-dispatch bytes'
[ -n "$capture_winner" ] && [ ! -e "$capture_root/$capture_winner/source.png" ] \
  && [ "$(cat "$capture_root/$capture_loser/source.png" 2>/dev/null)" = "dispatch-$capture_loser" ] \
  || fail_case 'capture-screenshot coordinated-race case did not preserve only the loser source'
leaked_private_paths="$(find "$capture_root/assets" -name 'result.png.copying.*' -print -quit)" \
  || fail_case 'capture-screenshot coordinated-race case could not search the assets directory'
[ -n "$leaked_private_paths" ] \
  && fail_case 'capture-screenshot coordinated-race case leaked private scratch'

prescribed_shell_finish

#!/usr/bin/env bash
# gate.sh <revision> <log>: run the repository gate at a revision in the detached drain-head checkout.
set -uo pipefail
main=/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
head="$main/.git/work-run-work-2026-11-15-012821/drain-head"
git -C "$head" checkout -q --detach "$1" || exit 97
cd "$head" || exit 98
export DO_WORK_FAST_STAGE_REUSE=off
start=$(date +%s)
bash _dev/tests/maintainer-verify.sh > "$2" 2>&1
status=$?
echo "GATE-EXIT=$status WALL=$(( $(date +%s) - start ))s REV=$(git rev-parse HEAD)" >> "$2"
exit $status

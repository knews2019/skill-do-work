#!/usr/bin/env bash
# Runtime pins for the board launcher scripts behind the run-kanban and kanban-static recipes.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
fail_count=0
contract_fail() {
  printf 'FAIL: %s\n' "$1" >&2
  fail_count=$((fail_count + 1))
}

# The run-kanban and kanban-static recipes are one-line calls to these two shipped scripts, so the
# scripts are what this contract drives. Each case runs a copy of the real script from a fixture
# board folder whose queue-kanban "binary" is a stub that records its argv, so no case can start
# the real server. go, lsof, ps, kill and sleep are exported functions (kill is a builtin, so a
# PATH stub could not replace it), and "$BASH" runs the child so the exported functions decode.
for board_script in run-kanban.sh kanban-static.sh; do
  if [ ! -f "$repo_root/skills/do-work-board/tools/$board_script" ]; then
    printf 'FAIL: skills/do-work-board/tools/%s is missing; the run-kanban and kanban-static recipes call it.\n' "$board_script" >&2
    exit 1
  fi
done
fixture_root="$(mktemp -d)"
trap 'rm -rf -- "$fixture_root"' EXIT
fixture_board="$fixture_root/do-work-board"
mkdir -p "$fixture_board/tools/queue-kanban"
cp "$repo_root/skills/do-work-board/tools/run-kanban.sh" "$repo_root/skills/do-work-board/tools/kanban-static.sh" "$fixture_board/tools/"
cat > "$fixture_board/tools/queue-kanban/queue-kanban" <<'STUB'
#!/bin/sh
printf '%s\n' "$*" >> "$FIXTURE_ARGV_LOG"
if [ "$1" = generate ] && [ "$2" = --out ]; then
  mkdir -p "$3" && printf '<html></html>\n' > "$3/index.html"
fi
STUB
chmod +x "$fixture_board/tools/queue-kanban/queue-kanban"
export FIXTURE_ARGV_LOG="$fixture_root/queue-kanban-argv.log"
export FIXTURE_KILL_LOG="$fixture_root/kill.log"
fixture_project="$fixture_root/project"
mkdir -p "$fixture_project"

go() { return 0; }
sleep() { return 0; }
export -f go sleep

reset_fixture_logs() {
  rm -f "$FIXTURE_ARGV_LOG" "$FIXTURE_KILL_LOG"
}

# 1. A stale queue-kanban listener is stopped, and serve then receives the project root and port.
# The lsof stub reports pid 4242 until kill records it, then reports nothing.
reset_fixture_logs
if stale_listener_output="$(
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/tmp/queue-kanban\n'
    elif [ -s "$FIXTURE_KILL_LOG" ]; then
      return 1
    else
      printf '4242\n'
    fi
  }
  ps() { printf '/tmp/queue-kanban serve --port 8090\n'; }
  kill() { printf '%s\n' "$*" >> "$FIXTURE_KILL_LOG"; }
  export -f lsof ps kill
  "$BASH" "$fixture_board/tools/run-kanban.sh" "$fixture_project" 8090 2>&1
)"; then
  stale_listener_status=0
else
  stale_listener_status=$?
fi
if [ "$stale_listener_status" -ne 0 ]; then
  contract_fail "run-kanban must serve after stopping a stale queue-kanban listener; got exit $stale_listener_status: $stale_listener_output"
fi
if ! grep -qxF 'queue-kanban: stopping previous session on :8090 (pid 4242): /tmp/queue-kanban serve --port 8090' <<<"$stale_listener_output"; then
  contract_fail "run-kanban must announce the stale listener it stops: $stale_listener_output"
fi
if [ "$(cat "$FIXTURE_KILL_LOG" 2>/dev/null)" != 4242 ]; then
  contract_fail 'run-kanban must call kill exactly once, with the stale pid 4242.'
fi
if [ "$(cat "$FIXTURE_ARGV_LOG" 2>/dev/null)" != "serve --open --repo-root $fixture_project --port 8090" ]; then
  contract_fail "run-kanban must start serve --open with the project root and port; got: $(cat "$FIXTURE_ARGV_LOG" 2>/dev/null)"
fi

# 2. A queue-kanban PID that remains a listener throughout the bounded wait must make run-kanban
# fail before build+serve, and the diagnosis must name both the PID and full command. sleep is a
# no-op in the fixture, so all 320 iterations run without slowing the suite.
reset_fixture_logs
if stuck_listener_output="$(
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/tmp/queue-kanban\n'
    else
      printf '4242\n'
    fi
  }
  ps() { printf '/tmp/queue-kanban serve --port 8090\n'; }
  kill() { return 0; }
  export -f lsof ps kill
  "$BASH" "$fixture_board/tools/run-kanban.sh" "$fixture_project" 8090 2>&1
)"; then
  stuck_listener_status=0
else
  stuck_listener_status=$?
fi
if [ "$stuck_listener_status" -ne 1 ]; then
  contract_fail "run-kanban must refuse startup when a listener remains after the bounded wait; got exit $stuck_listener_status."
fi
if ! grep -qF 'pid 4242' <<<"$stuck_listener_output"; then
  contract_fail 'run-kanban stuck-listener refusal must name pid 4242.'
fi
if ! grep -qF '/tmp/queue-kanban serve --port 8090' <<<"$stuck_listener_output"; then
  contract_fail 'run-kanban stuck-listener refusal must name the listener command.'
fi
if [ -e "$FIXTURE_ARGV_LOG" ]; then
  contract_fail 'run-kanban must not start serve while a listener remains.'
fi

# 3. Preserve the older safety boundary: a foreign executable is refused immediately and never
# passed to kill.
reset_fixture_logs
if foreign_listener_output="$(
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/usr/bin/python3\n'
    else
      printf '3131\n'
    fi
  }
  ps() { printf '/usr/bin/python3 -m http.server 8090\n'; }
  kill() { printf '%s\n' "$*" >> "$FIXTURE_KILL_LOG"; }
  export -f lsof ps kill
  "$BASH" "$fixture_board/tools/run-kanban.sh" "$fixture_project" 8090 2>&1
)"; then
  foreign_listener_status=0
else
  foreign_listener_status=$?
fi
# This fixture asserts the refusal status and kill boundary; its captured text is intentionally ignored.
: "$foreign_listener_output"
if [ "$foreign_listener_status" -ne 1 ] || [ -s "$FIXTURE_KILL_LOG" ]; then
  contract_fail 'run-kanban must refuse a foreign listener without calling kill.'
fi

# 4. kanban-static writes the snapshot and exactly one local exclude line across two runs, and
# announces the write once. -x matches the whole line, because just echoes recipe text.
static_project="$fixture_root/static-project"
git init -q "$static_project"
static_exclude_file="$static_project/.git/info/exclude"
exclude_message='kanban-static: added /build/queue-kanban-board/ to .git/info/exclude (local-only ignore)'
first_static_status=0
first_static_output="$("$BASH" "$fixture_board/tools/kanban-static.sh" "$static_project" 2>&1)" || first_static_status=$?
second_static_status=0
second_static_output="$("$BASH" "$fixture_board/tools/kanban-static.sh" "$static_project" 2>&1)" || second_static_status=$?
if [ "$first_static_status" -ne 0 ] || [ "$second_static_status" -ne 0 ]; then
  contract_fail "kanban-static must succeed twice; got exits $first_static_status and $second_static_status: $first_static_output $second_static_output"
fi
if [ ! -f "$static_project/build/queue-kanban-board/index.html" ]; then
  contract_fail 'kanban-static must generate build/queue-kanban-board/index.html in the project root.'
fi
static_exclude_lines="$(grep -cxF '/build/queue-kanban-board/' "$static_exclude_file" || true)"
if [ "$static_exclude_lines" != 1 ]; then
  contract_fail "kanban-static must leave exactly one /build/queue-kanban-board/ exclude line after two runs; found $static_exclude_lines."
fi
if ! grep -qxF "$exclude_message" <<<"$first_static_output"; then
  contract_fail "the first kanban-static run must announce the exclude write: $first_static_output"
fi
if grep -qxF "$exclude_message" <<<"$second_static_output"; then
  contract_fail 'the second kanban-static run must not announce a second exclude write.'
fi

# 5. A failed build stops kanban-static before the exclude write.
failed_static_project="$fixture_root/failed-static-project"
git init -q "$failed_static_project"
if failed_static_output="$(
  go() { return 1; }
  export -f go
  "$BASH" "$fixture_board/tools/kanban-static.sh" "$failed_static_project" 2>&1
)"; then
  failed_static_status=0
else
  failed_static_status=$?
fi
: "$failed_static_output"
if [ "$failed_static_status" -eq 0 ] \
  || grep -qxF '/build/queue-kanban-board/' "$failed_static_project/.git/info/exclude" 2>/dev/null; then
  contract_fail "kanban-static must exit non-zero on a failed build and write no exclude line; got exit $failed_static_status."
fi
# REQ-633: Panel B's idle-gap rule (board) and the calibration log's max_stamp_gap_minutes
# (core CLI) must read one stamp set: the board's phase milestones minus completed_at and release_at.
board_gap_stamps="$(grep -o 'fieldName: "[a-z_]*"' "$repo_root/skills/do-work-board/tools/queue-kanban/durations.go" \
  | sed 's/^fieldName: "\(.*\)"$/\1/' | grep -v -x -e completed_at -e release_at | tr '\n' ' ')"
cli_gap_stamps="$(sed -n 's/^var CalibrationGapStampFields = \[\]string{\(.*\)}$/\1/p' \
  "$repo_root/skills/do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go" | tr -d '",' | tr -s ' ')"
if [ -z "$cli_gap_stamps" ] || [ "${board_gap_stamps% }" != "$cli_gap_stamps" ]; then
  printf 'FAIL: calibration gap stamps (%s) differ from the board phase milestones (%s).\n' "$cli_gap_stamps" "$board_gap_stamps" >&2
  fail_count=$((fail_count + 1))
fi
[ "$fail_count" -eq 0 ] || exit 1
printf 'queue-kanban contract probes passed.\n'

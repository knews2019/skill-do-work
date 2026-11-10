#!/usr/bin/env bash
# Live proof for REQ-717 (suite-managed justfile recipes call a script in one line).
# Runs the shipped board template in a throwaway consumer repository, so nothing is written to the
# tree under test: the board package is copied (without its compiled binary) to
# <consumer>/.claude/skills/do-work-board and the template becomes the consumer's justfile.
#   1. just run-kanban <port> twice: the second run stops the first queue-kanban server on that port
#      and serves in its place.
#   2. just kanban-static twice: build/queue-kanban-board/index.html exists, .git/info/exclude holds
#      exactly one /build/queue-kanban-board/ line, and only the first run announces the write.
#   3. The tree's own root justfile: run-kanban rejects a bad port (proves the recipe reaches the
#      port validation) and kanban-static names a script that exists (or is still inline at base).
# A stub `open`/`xdg-open` first on PATH keeps the browser closed. Every server this script starts is
# stopped on exit; the script fails if one is still listening.
# It passes at the base revision too: the change is behavior-preserving, so this is a
# characterization, not a red-green test.
set -euo pipefail

tree_root="$(git rev-parse --show-toplevel)"
proof_root="$(mktemp -d "${TMPDIR:-/tmp}/req717-live-proof.XXXXXX")"
consumer_root="$proof_root/consumer"
stub_directory="$proof_root/stub-bin"
started_pids=()
proof_port=''

stop_started_servers() {
  local started_pid
  for started_pid in ${started_pids[@]+"${started_pids[@]}"}; do
    kill "$started_pid" 2>/dev/null || true
  done
  if [ -n "$proof_port" ]; then
    local leftover_pid
    for leftover_pid in $(lsof -ti tcp:"$proof_port" -sTCP:LISTEN 2>/dev/null || true); do
      case "$(ps -p "$leftover_pid" -o comm= 2>/dev/null || true)" in
        *queue-kanban*) kill "$leftover_pid" 2>/dev/null || true ;;
      esac
    done
  fi
  rm -rf -- "$proof_root"
}
trap stop_started_servers EXIT

fail() {
  printf 'LIVE-PROOF FAIL: %s\n' "$1" >&2
  exit 1
}

listener_pid_on() {
  lsof -ti tcp:"$1" -sTCP:LISTEN 2>/dev/null | head -n1 || true
}

mkdir -p "$stub_directory" "$consumer_root/.claude/skills" "$consumer_root/do-work"
printf '#!/bin/sh\nexit 0\n' > "$stub_directory/open"
cp "$stub_directory/open" "$stub_directory/xdg-open"
chmod +x "$stub_directory/open" "$stub_directory/xdg-open"
export PATH="$stub_directory:$PATH"

cp -R "$tree_root/skills/do-work-board" "$consumer_root/.claude/skills/do-work-board"
rm -f "$consumer_root/.claude/skills/do-work-board/tools/queue-kanban/queue-kanban"
cp "$tree_root/skills/do-work-board/justfile.template" "$consumer_root/justfile"
git -C "$consumer_root" init -q
git -C "$consumer_root" -c user.name=proof -c user.email=proof@example.invalid \
  commit -q --allow-empty -m 'proof consumer'

consumer_just() {
  just --justfile "$consumer_root/justfile" --working-directory "$consumer_root" "$@"
}

# --- 1. run-kanban replaces a stale instance on its port -------------------------------------
for candidate_port in 18717 18727 18737 18747 18757; do
  if [ -z "$(listener_pid_on "$candidate_port")" ]; then
    proof_port="$candidate_port"
    break
  fi
done
[ -n "$proof_port" ] || fail 'no free proof port among 18717..18757'

wait_for_listener() {
  local expected_other_than="$1" wait_steps=0 found_pid
  while [ "$wait_steps" -lt 1200 ]; do
    found_pid="$(listener_pid_on "$proof_port")"
    if [ -n "$found_pid" ] && [ "$found_pid" != "$expected_other_than" ]; then
      printf '%s\n' "$found_pid"
      return 0
    fi
    sleep 0.1
    wait_steps=$((wait_steps + 1))
  done
  return 1
}

consumer_just run-kanban "$proof_port" > "$proof_root/first-run.log" 2>&1 &
started_pids+=("$!")
first_listener_pid="$(wait_for_listener none)" \
  || fail "first run-kanban never listened on :$proof_port ($(cat "$proof_root/first-run.log"))"
started_pids+=("$first_listener_pid")

consumer_just run-kanban "$proof_port" > "$proof_root/second-run.log" 2>&1 &
started_pids+=("$!")
second_listener_pid="$(wait_for_listener "$first_listener_pid")" \
  || fail "second run-kanban never replaced pid $first_listener_pid ($(cat "$proof_root/second-run.log"))"
started_pids+=("$second_listener_pid")

second_run_log="$(cat "$proof_root/second-run.log")"
grep -qF "queue-kanban: stopping previous session on :$proof_port (pid $first_listener_pid)" <<<"$second_run_log" \
  || fail "second run did not announce the stop of pid $first_listener_pid: $second_run_log"
if kill -0 "$first_listener_pid" 2>/dev/null; then
  fail "stale pid $first_listener_pid is still alive after the replacement"
fi
printf 'PASS run-kanban: pid %s replaced pid %s on :%s\n' "$second_listener_pid" "$first_listener_pid" "$proof_port"

# --- 2. kanban-static writes the snapshot and exactly one exclude line -----------------------
first_static_output="$(consumer_just kanban-static 2>&1)" || fail "kanban-static failed: $first_static_output"
second_static_output="$(consumer_just kanban-static 2>&1)" || fail "second kanban-static failed: $second_static_output"
[ -s "$consumer_root/build/queue-kanban-board/index.html" ] || fail 'kanban-static wrote no build/queue-kanban-board/index.html'
exclude_file="$(git -C "$consumer_root" rev-parse --git-path info/exclude)"
case "$exclude_file" in /*) ;; *) exclude_file="$consumer_root/$exclude_file" ;; esac
exclude_line_count="$(grep -cxF '/build/queue-kanban-board/' "$exclude_file" || true)"
[ "$exclude_line_count" = 1 ] || fail "exclude holds $exclude_line_count /build/queue-kanban-board/ lines, want 1"
exclude_message='kanban-static: added /build/queue-kanban-board/ to .git/info/exclude (local-only ignore)'
grep -qxF "$exclude_message" <<<"$first_static_output" || fail "first kanban-static did not announce the exclude write"
if grep -qxF "$exclude_message" <<<"$second_static_output"; then
  fail 'second kanban-static announced a second exclude write'
fi
if [ -n "$(git -C "$consumer_root" status --porcelain --untracked-files=all -- build)" ]; then
  fail 'build/ snapshot still shows in git status'
fi
printf 'PASS kanban-static: snapshot written, one exclude line, idempotent\n'

# --- 3. the tree's own root justfile reaches the same logic ----------------------------------
root_bad_port_status=0
root_bad_port_output="$(just --justfile "$tree_root/justfile" --working-directory "$tree_root" run-kanban 'not-a-port' 2>&1)" \
  || root_bad_port_status=$?
[ "$root_bad_port_status" -ne 0 ] || fail 'root run-kanban accepted an invalid port'
grep -qF "queue-kanban: invalid port 'not-a-port' - must be digits only" <<<"$root_bad_port_output" \
  || fail "root run-kanban did not print the invalid-port message: $root_bad_port_output"
root_static_command="$(just --justfile "$tree_root/justfile" --working-directory "$tree_root" --dry-run kanban-static 2>&1)"
root_static_script="$(grep -oE 'skills/do-work-board/tools/[A-Za-z0-9._-]+\.sh' <<<"$root_static_command" | head -n1 || true)"
if [ -n "$root_static_script" ]; then
  [ -f "$tree_root/$root_static_script" ] || fail "root kanban-static names a missing script: $root_static_script"
  printf 'PASS root justfile: run-kanban validates the port; kanban-static calls %s\n' "$root_static_script"
else
  printf 'PASS root justfile: run-kanban validates the port; kanban-static is still inline (base revision)\n'
fi

printf 'REQ-717 live proof passed.\n'

#!/usr/bin/env bash
# Fixture execution proofs for generate-report-image.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# Bounded stand-in for `wait` on a backgrounded wrapper. `kill -TERM` followed by a bare
# `wait` blocks forever against a wrapper that does not finish, which is how a stuck backend
# wedges the whole gate — and the gate's exit status is the only thing anyone reads. Wait with
# a deadline instead: on expiry, name the processes still alive, fail the case, and kill them.
# Sets wrapper_wait_status. Arguments: wrapper PID, case label, deadline seconds, then any
# further PIDs to name in the diagnostic.
wait_for_wrapper_or_fail() {
  local wrapper_pid="$1"
  local case_label="$2"
  local deadline_seconds="$3"
  shift 3
  local deadline_report="$fixture_root/.wrapper-deadline-$wrapper_pid"
  local watchdog_pid
  local survivor_pid
  local surviving_pids
  rm -f "$deadline_report"
  (
    sleep "$deadline_seconds"
    surviving_pids=""
    for survivor_pid in "$wrapper_pid" "$@"; do
      kill -0 "$survivor_pid" 2>/dev/null && surviving_pids="$surviving_pids $survivor_pid"
    done
    printf '%s\n' "${surviving_pids# }" > "$deadline_report"
    for survivor_pid in "$wrapper_pid" "$@"; do
      kill -KILL "$survivor_pid" 2>/dev/null || true
    done
  ) &
  watchdog_pid=$!
  wait "$wrapper_pid" 2>/dev/null
  wrapper_wait_status=$?
  kill "$watchdog_pid" 2>/dev/null || true
  wait "$watchdog_pid" 2>/dev/null || true
  # The report file exists only when the deadline actually expired, so its presence — not the
  # wait status, which a killed wrapper also makes nonzero — is what distinguishes a stuck
  # wrapper from an interrupted one.
  [ -e "$deadline_report" ] || return 0
  fail_case "$case_label did not finish within ${deadline_seconds}s; still alive: $(cat "$deadline_report")"
  return 1
}

# generate-report-image: a direct backend receives inert prompt text and publishes
# from a private adjacent path only after success.
image_bin="$fixture_root/image-bin"
mkdir -p "$image_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) printf "%s" "$2" > "$IMAGE_PROMPT_LOG"; shift 2 ;; *) shift ;; esac; done' 'printf "%s" "$output_path" > "$IMAGE_STAGED_PATH_LOG"' 'printf new-png > "$output_path"' > "$image_bin/imagegen"
chmod +x "$image_bin/imagegen"
image_output="$fixture_root/report.png"
printf old-png > "$image_output"
PATH="$image_bin:$PATH" \
  IMAGE_PROMPT_LOG="$fixture_root/image-prompt" \
  IMAGE_STAGED_PATH_LOG="$fixture_root/image-staged-path" \
  "$toolbox_scripts/generate-report-image.sh" "$image_output" 'blue style' 'diagram $(touch injected-image)' \
  || fail_case 'generate-report-image direct-backend case returned nonzero'
[ "$(cat "$image_output")" = new-png ] && [ ! -e "$fixture_root/injected-image" ] \
  || fail_case 'generate-report-image direct-backend case did not atomically replace the target with inert-prompt output'
[ "$(cat "$fixture_root/image-staged-path")" != "$image_output" ] \
  || fail_case 'generate-report-image direct-backend case wrote directly to the final target'
find "$fixture_root" -name '.report.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image direct-backend case leaked private staging after success'

# generate-report-image: a failed backend may leave an old target recoverable, but
# the invocation must fail and clean its partial private output.
image_failure_bin="$fixture_root/image-failure-bin"
mkdir -p "$image_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' 'printf partial > "$output_path"' 'exit 7' > "$image_failure_bin/imagegen"
chmod +x "$image_failure_bin/imagegen"
printf stable-old-png > "$image_output"
PATH="$image_failure_bin:$PATH" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=0 \
  "$toolbox_scripts/generate-report-image.sh" "$image_output" style failure >/dev/null 2>&1 \
  && fail_case 'generate-report-image stale-target case accepted a failed backend'
[ "$(cat "$image_output")" = stable-old-png ] \
  || fail_case 'generate-report-image stale-target case changed the recoverable old target'
find "$fixture_root" -name '.report.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image stale-target case leaked private staging after failure'

# generate-report-image caller contract: mixed results retain both PIDs/statuses and
# wait for every job before evaluating current outputs.
image_mixed_bin="$fixture_root/image-mixed-bin"
mkdir -p "$image_mixed_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' 'case "$image_prompt" in *mixed-success*) printf success > "$output_path"; : > "$MIXED_SUCCESS_DONE"; exit 0 ;; *) printf partial > "$output_path"; : > "$MIXED_FAILURE_DONE"; exit 9 ;; esac' > "$image_mixed_bin/imagegen"
chmod +x "$image_mixed_bin/imagegen"
mixed_success_target="$fixture_root/mixed-success.png"
mixed_failure_target="$fixture_root/mixed-failure.png"
printf stale > "$mixed_failure_target"
PATH="$image_mixed_bin:$PATH" MIXED_SUCCESS_DONE="$fixture_root/mixed-success.done" MIXED_FAILURE_DONE="$fixture_root/mixed-failure.done" \
  "$toolbox_scripts/generate-report-image.sh" "$mixed_success_target" style mixed-success &
image_generation_pids[0]=$!
PATH="$image_mixed_bin:$PATH" MIXED_SUCCESS_DONE="$fixture_root/mixed-success.done" MIXED_FAILURE_DONE="$fixture_root/mixed-failure.done" \
  "$toolbox_scripts/generate-report-image.sh" "$mixed_failure_target" style mixed-failure &
image_generation_pids[1]=$!
image_generation_statuses=()
image_index=0
while [ "$image_index" -lt "${#image_generation_pids[@]}" ]; do
  image_status=0
  wait "${image_generation_pids[$image_index]}" || image_status=$?
  image_generation_statuses[$image_index]="$image_status"
  image_index=$((image_index + 1))
done
[ "${image_generation_statuses[0]}" -eq 0 ] && [ "${image_generation_statuses[1]}" -ne 0 ] \
  || fail_case 'generate-report-image mixed-status case did not retain success and failure independently'
[ -e "$fixture_root/mixed-success.done" ] && [ -e "$fixture_root/mixed-failure.done" ] \
  || fail_case 'generate-report-image mixed-status case did not wait for every launched job'
[ "$(cat "$mixed_success_target")" = success ] && [ "$(cat "$mixed_failure_target")" = stale ] \
  || fail_case 'generate-report-image mixed-status case did not separate current output from a stale failed target'

# generate-report-image: interruption cleans the invocation-private file and leaves
# the old target untouched.
image_interrupt_bin="$fixture_root/image-interrupt-bin"
mkdir -p "$image_interrupt_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' 'printf "%s" "$output_path" > "$IMAGE_INTERRUPT_PATH_LOG"' 'printf partial > "$output_path"' ': > "$IMAGE_INTERRUPT_READY"' 'trap "exit 143" TERM INT HUP' 'while :; do sleep 0.1; done' > "$image_interrupt_bin/imagegen"
chmod +x "$image_interrupt_bin/imagegen"
interrupt_target="$fixture_root/interrupted.png"
printf stable-interrupt > "$interrupt_target"
PATH="$image_interrupt_bin:$PATH" IMAGE_INTERRUPT_PATH_LOG="$fixture_root/interrupted-stage" IMAGE_INTERRUPT_READY="$fixture_root/interrupted-ready" \
  "$toolbox_scripts/generate-report-image.sh" "$interrupt_target" style interruption >/dev/null 2>&1 &
interrupt_helper_pid=$!
background_process_ids="$background_process_ids $interrupt_helper_pid"
interrupt_wait_ticks=0
while [ ! -e "$fixture_root/interrupted-ready" ] && [ "$interrupt_wait_ticks" -lt 200 ]; do
  sleep 0.01
  interrupt_wait_ticks=$((interrupt_wait_ticks + 1))
done
if [ ! -e "$fixture_root/interrupted-ready" ]; then
  fail_case 'generate-report-image interruption case never reached the backend wait'
else
  kill -TERM "$interrupt_helper_pid" 2>/dev/null || true
  wait_for_wrapper_or_fail "$interrupt_helper_pid" 'generate-report-image interruption case' 10
  interrupt_status="$wrapper_wait_status"
  [ "$interrupt_status" -ne 0 ] || fail_case 'generate-report-image interruption case returned success'
  interrupted_stage_path="$(cat "$fixture_root/interrupted-stage")"
  [ ! -e "$interrupted_stage_path" ] || fail_case 'generate-report-image interruption case leaked private staging'
  [ "$(cat "$interrupt_target")" = stable-interrupt ] || fail_case 'generate-report-image interruption case changed the old target'
fi

# generate-report-image: the sandbox-bypassed executable is unreachable without the
# exact opt-in value, while exact 1 still exercises the explicitly authorized branch.
agentic_bin="$fixture_root/agentic-bin"
mkdir -p "$agentic_bin"
for agentic_tool in bash chmod cp mktemp mv rm; do
  ln -s "$(command -v "$agentic_tool")" "$agentic_bin/$agentic_tool"
done
printf '%s\n' '#!/usr/bin/env bash' ': > "$AGENTIC_INVOKED_MARKER"' 'printf agentic-png > generated.png' > "$agentic_bin/codex"
chmod +x "$agentic_bin/codex"
for agentic_opt_in in unset 0 true 01; do
  agentic_marker="$fixture_root/agentic-$agentic_opt_in.invoked"
  if [ "$agentic_opt_in" = unset ]; then
    (unset DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND; PATH="$agentic_bin" AGENTIC_INVOKED_MARKER="$agentic_marker" TMPDIR="$fixture_root" \
      "$toolbox_scripts/generate-report-image.sh" "$fixture_root/agentic-$agentic_opt_in.png" style description >/dev/null 2>&1)
  else
    PATH="$agentic_bin" AGENTIC_INVOKED_MARKER="$agentic_marker" TMPDIR="$fixture_root" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND="$agentic_opt_in" \
      "$toolbox_scripts/generate-report-image.sh" "$fixture_root/agentic-$agentic_opt_in.png" style description >/dev/null 2>&1
  fi
  agentic_status=$?
  [ "$agentic_status" -ne 0 ] || fail_case "generate-report-image agentic opt-in case accepted non-exact value $agentic_opt_in"
  [ ! -e "$agentic_marker" ] || fail_case "generate-report-image agentic opt-in case invoked codex for value $agentic_opt_in"
done
agentic_marker="$fixture_root/agentic-1.invoked"
PATH="$agentic_bin" AGENTIC_INVOKED_MARKER="$agentic_marker" TMPDIR="$fixture_root" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1 \
  "$toolbox_scripts/generate-report-image.sh" "$fixture_root/agentic-1.png" style description \
  || fail_case 'generate-report-image agentic opt-in case rejected exact value 1'
[ -e "$agentic_marker" ] && [ "$(cat "$fixture_root/agentic-1.png")" = agentic-png ] \
  || fail_case 'generate-report-image agentic opt-in case did not publish the explicitly authorized output'
find "$fixture_root" \( -name '.*.generating.*' -o -name 'do-work-ai-report-image.*' \) -print -quit | grep -q . \
  && fail_case 'generate-report-image agentic opt-in case leaked private paths'

# generate-report-image, interrupted directly: the helper owns the process tree it
# started, so signalling it must leave neither the backend nor the backend's own
# descendant alive, and must reap them before the private stage is removed. A
# bare-PID kill passes the file assertions below while the descendant keeps running.
image_tree_bin="$fixture_root/image-tree-bin"
mkdir -p "$image_tree_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf "%s" "$output_path" > "$IMAGE_TREE_STAGE_LOG"' \
  'printf partial > "$output_path"' \
  '( while :; do sleep 0.2; done ) &' \
  'printf "%s\n" "$!" > "$IMAGE_TREE_DESCENDANT_PID"' \
  'printf "%s\n" "$$" > "$IMAGE_TREE_BACKEND_PID"' \
  ': > "$IMAGE_TREE_READY"' \
  'while :; do sleep 0.2; done' \
  > "$image_tree_bin/imagegen"
chmod +x "$image_tree_bin/imagegen"
image_tree_target="$fixture_root/process-tree.png"
printf stable-tree > "$image_tree_target"
PATH="$image_tree_bin:$PATH" \
  IMAGE_TREE_STAGE_LOG="$fixture_root/process-tree-stage" \
  IMAGE_TREE_BACKEND_PID="$fixture_root/process-tree-backend.pid" \
  IMAGE_TREE_DESCENDANT_PID="$fixture_root/process-tree-descendant.pid" \
  IMAGE_TREE_READY="$fixture_root/process-tree-ready" \
  "$toolbox_scripts/generate-report-image.sh" "$image_tree_target" style process-tree >/dev/null 2>&1 &
image_tree_helper_pid=$!
background_process_ids="$background_process_ids $image_tree_helper_pid"
image_tree_ready_ticks=0
while [ "$image_tree_ready_ticks" -lt 200 ]; do
  [ -e "$fixture_root/process-tree-ready" ] \
    && [ -s "$fixture_root/process-tree-backend.pid" ] \
    && [ -s "$fixture_root/process-tree-descendant.pid" ] && break
  sleep 0.05
  image_tree_ready_ticks=$((image_tree_ready_ticks + 1))
done
if [ ! -s "$fixture_root/process-tree-backend.pid" ] || [ ! -s "$fixture_root/process-tree-descendant.pid" ]; then
  fail_case 'generate-report-image process-tree case never reached a running backend and descendant'
else
  image_tree_backend_pid="$(cat "$fixture_root/process-tree-backend.pid")"
  image_tree_descendant_pid="$(cat "$fixture_root/process-tree-descendant.pid")"
  background_process_ids="$background_process_ids $image_tree_backend_pid $image_tree_descendant_pid"
  kill -TERM "$image_tree_helper_pid" 2>/dev/null || true
  wait_for_wrapper_or_fail "$image_tree_helper_pid" 'generate-report-image process-tree case' 10 \
    "$image_tree_backend_pid" "$image_tree_descendant_pid"
  image_tree_status="$wrapper_wait_status"
  [ "$image_tree_status" -eq 143 ] \
    || fail_case "generate-report-image process-tree case returned $image_tree_status instead of the TERM status 143"
  image_tree_survivor_ticks=0
  while [ "$image_tree_survivor_ticks" -lt 40 ]; do
    image_tree_survivors=0
    for image_tree_recorded_pid in "$image_tree_backend_pid" "$image_tree_descendant_pid"; do
      kill -0 "$image_tree_recorded_pid" 2>/dev/null && image_tree_survivors=$((image_tree_survivors + 1))
    done
    [ "$image_tree_survivors" -eq 0 ] && break
    sleep 0.05
    image_tree_survivor_ticks=$((image_tree_survivor_ticks + 1))
  done
  [ "$image_tree_survivors" -eq 0 ] \
    || fail_case "generate-report-image process-tree case left $image_tree_survivors backend process(es) or descendant(s) alive"
  image_tree_stage_path="$(cat "$fixture_root/process-tree-stage")"
  [ ! -e "$image_tree_stage_path" ] || fail_case 'generate-report-image process-tree case leaked private staging'
  [ "$(cat "$image_tree_target")" = stable-tree ] || fail_case 'generate-report-image process-tree case changed the old target'
fi

# generate-report-image: `mv` treats an existing destination directory as a container,
# so an output path occupied by a directory nests the staged image inside it and still
# exits zero. Publication must fail closed and leave that directory untouched.
image_directory_bin="$fixture_root/image-directory-bin"
mkdir -p "$image_directory_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf new-png > "$output_path"' \
  > "$image_directory_bin/imagegen"
chmod +x "$image_directory_bin/imagegen"
image_directory_parent="$fixture_root/directory-target"
image_directory_target="$image_directory_parent/report.png"
mkdir -p "$image_directory_target"
printf owned-by-someone-else > "$image_directory_target/keep.txt"
PATH="$image_directory_bin:$PATH" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=0 \
  "$toolbox_scripts/generate-report-image.sh" "$image_directory_target" style directory-collision >/dev/null 2>&1 \
  && fail_case 'generate-report-image output-is-a-directory case reported success'
[ "$(cat "$image_directory_target/keep.txt" 2>/dev/null)" = owned-by-someone-else ] \
  || fail_case 'generate-report-image output-is-a-directory case did not preserve the occupying directory byte-for-byte'
[ "$(ls -A "$image_directory_target" 2>/dev/null)" = keep.txt ] \
  || fail_case 'generate-report-image output-is-a-directory case left its staged image nested inside the occupying directory'
find "$image_directory_parent" -name '.report.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image output-is-a-directory case leaked private staging'

# generate-report-image: an interruption arriving as early as the invocation is observable —
# at the moment the private staging file appears, which is before any backend is launched —
# is still reported as an interruption and still cleans up. This pins the deferral the wrapper
# uses to close its publish-the-PID window: a deferred status that is dropped instead of
# re-raised turns the interruption into an ordinary backend failure (exit 1), and a backend
# started after the deferral began must still be reaped.
image_early_bin="$fixture_root/image-early-bin"
mkdir -p "$image_early_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf "%s\n" "$$" > "$IMAGE_EARLY_BACKEND_PID"' \
  'printf partial > "$output_path"' \
  'trap "exit 143" TERM INT HUP' \
  'while :; do sleep 0.1; done' \
  > "$image_early_bin/imagegen"
chmod +x "$image_early_bin/imagegen"
image_early_target="$fixture_root/early-interrupt.png"
printf stable-early > "$image_early_target"
PATH="$image_early_bin:$PATH" \
  IMAGE_EARLY_BACKEND_PID="$fixture_root/early-backend.pid" \
  "$toolbox_scripts/generate-report-image.sh" "$image_early_target" style early-interruption >/dev/null 2>&1 &
image_early_helper_pid=$!
background_process_ids="$background_process_ids $image_early_helper_pid"
image_early_stage_ticks=0
while [ "$image_early_stage_ticks" -lt 500 ]; do
  find "$fixture_root" -name '.early-interrupt.png.generating.*' -print -quit | grep -q . && break
  sleep 0.002
  image_early_stage_ticks=$((image_early_stage_ticks + 1))
done
kill -TERM "$image_early_helper_pid" 2>/dev/null || true
wait_for_wrapper_or_fail "$image_early_helper_pid" 'generate-report-image early-interruption case' 10
[ "$wrapper_wait_status" -eq 143 ] \
  || fail_case "generate-report-image early-interruption case returned $wrapper_wait_status instead of the TERM status 143"
image_early_backend_pid="$(cat "$fixture_root/early-backend.pid" 2>/dev/null || printf '')"
if [ -n "$image_early_backend_pid" ]; then
  background_process_ids="$background_process_ids $image_early_backend_pid"
  image_early_survivor_ticks=0
  while [ "$image_early_survivor_ticks" -lt 40 ]; do
    kill -0 "$image_early_backend_pid" 2>/dev/null || break
    sleep 0.05
    image_early_survivor_ticks=$((image_early_survivor_ticks + 1))
  done
  kill -0 "$image_early_backend_pid" 2>/dev/null \
    && fail_case "generate-report-image early-interruption case left backend $image_early_backend_pid alive"
fi
find "$fixture_root" -name '.early-interrupt.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image early-interruption case leaked private staging'
[ "$(cat "$image_early_target")" = stable-early ] \
  || fail_case 'generate-report-image early-interruption case changed the old target'

# generate-report-image: every interruption case above waits with a deadline, so a wrapper
# that will not finish fails the probe with a diagnostic naming what is still alive instead of
# wedging the gate forever. Exercised against a deliberately TERM-deaf stand-in, inside a
# command substitution so the deadline's intentional failure stays out of this file's tally.
deadline_probe_report="$( {
  deadline_probe_ready="$fixture_root/deadline-probe-ready"
  ( trap '' TERM; : > "$deadline_probe_ready"; while :; do sleep 0.1; done ) >/dev/null 2>&1 &
  deadline_probe_pid=$!
  deadline_probe_ticks=0
  while [ ! -e "$deadline_probe_ready" ] && [ "$deadline_probe_ticks" -lt 200 ]; do
    sleep 0.01
    deadline_probe_ticks=$((deadline_probe_ticks + 1))
  done
  kill -TERM "$deadline_probe_pid" 2>/dev/null || true
  wait_for_wrapper_or_fail "$deadline_probe_pid" 'generate-report-image deadline probe' 1
} 2>&1 )"
case "$deadline_probe_report" in
  *'deadline probe did not finish within 1s; still alive: '[0-9]*) ;;
  *) fail_case "generate-report-image deadline case did not name the surviving process (got: $deadline_probe_report)" ;;
esac

prescribed_shell_finish

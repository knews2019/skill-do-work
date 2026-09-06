#!/usr/bin/env bash
# Long-lived fixture shells stop themselves after 30s, beyond the cleanup assertions.
# Parent traps and the existing wait deadlines only provide earlier cleanup.
# Fixture execution proofs for generate-report-image-batch.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"
toolbox_scripts="$fixture_root/toolbox-scripts"
mkdir -p "$toolbox_scripts"
cat > "$toolbox_scripts/generate-report-image.sh" <<EOF
#!/usr/bin/env bash
set -euo pipefail
launcher_arguments=(--format text)
if [[ -n "\${DO_WORK_COMPATIBILITY_REPO_ROOT:-}" ]]; then
  launcher_arguments+=(--repo-root "\$DO_WORK_COMPATIBILITY_REPO_ROOT")
fi
exec bash "$repo_root/skills/do-work/tools/do-work-cli.sh" "\${launcher_arguments[@]}" generate-report-image "\$@"
EOF
cat > "$toolbox_scripts/generate-report-image-batch.sh" <<EOF
#!/usr/bin/env bash
set -euo pipefail
launcher_arguments=(--format text)
if [[ -n "\${DO_WORK_COMPATIBILITY_REPO_ROOT:-}" ]]; then
  launcher_arguments+=(--repo-root "\$DO_WORK_COMPATIBILITY_REPO_ROOT")
fi
compatibility_arguments=("\$@")
report_index=0
if [[ "\${compatibility_arguments[0]:-}" == --dry-run || "\${compatibility_arguments[0]:-}" == --commit ]]; then
  report_index=1
fi
if [[ -n "\${compatibility_arguments[\$report_index]:-}" && "\${compatibility_arguments[\$report_index]}" != /* ]]; then
  compatibility_arguments[\$report_index]="\$PWD/\${compatibility_arguments[\$report_index]}"
fi
exec bash "$repo_root/skills/do-work/tools/do-work-cli.sh" "\${launcher_arguments[@]}" generate-report-image-batch "\${compatibility_arguments[@]}"
EOF
chmod +x "$toolbox_scripts/generate-report-image.sh" "$toolbox_scripts/generate-report-image-batch.sh"

# Canonical publication is repository-rooted.
fixture_repo_init "$fixture_root"
export DO_WORK_COMPATIBILITY_REPO_ROOT="$fixture_root"

process_id_is_live() {
  process_state="$(/bin/ps -o stat= -p "$1" 2>/dev/null)" || return 1
  process_state="${process_state//[[:space:]]/}"
  case "$process_state" in ''|Z*) return 1 ;; *) return 0 ;; esac
}

# generate-report-image-batch: the batch joins each target name to its own private
# staging directory, so a target carrying a path separator would write outside that
# boundary, and a pair with no colon would silently make the prompt the filename. Both
# are usage errors that must be refused before any staging directory exists.
batch_arguments_root="$fixture_root/batch-arguments"
mkdir -p "$batch_arguments_root/ai-reports/report"
"$toolbox_scripts/generate-report-image-batch.sh" "$batch_arguments_root/ai-reports/report" style \
  '../escaped.png:<prompt 1>' >/dev/null 2>&1
batch_escape_status=$?
[ "$batch_escape_status" -eq 2 ] \
  || fail_case "generate-report-image-batch path-separator target case returned $batch_escape_status instead of the usage status 2"
"$toolbox_scripts/generate-report-image-batch.sh" "$batch_arguments_root/ai-reports/report" style \
  '01-architecture.png' >/dev/null 2>&1
batch_unpaired_status=$?
[ "$batch_unpaired_status" -eq 2 ] \
  || fail_case "generate-report-image-batch unpaired-argument case returned $batch_unpaired_status instead of the usage status 2"
leaked_private_paths="$(find "$batch_arguments_root/ai-reports/report" -name '.generated.staging.*' -print -quit)" \
  || fail_case 'generate-report-image-batch usage-error cases could not search the report directory for invocation-private staging'
[ -n "$leaked_private_paths" ] \
  && fail_case 'generate-report-image-batch usage-error cases allocated invocation-private staging'

# generate-report-image-batch: the shipped batch owns staging, launch, wait-all,
# freshness, and publication. An all-failed batch must leave no public or private
# directory and still return zero so the caller falls back to SVG/Mermaid; a mixed
# batch must publish only the status-zero, non-empty current image after waiting
# every job, and report the published directory on stdout.
run_ai_report_batch_replay() {
  replay_name="$1"
  replay_bin="$2"
  replay_root="$fixture_root/$replay_name"
  mkdir -p "$replay_root/ai-reports/<report-slug>"
  (
    cd "$replay_root" || exit 2
    PATH="$replay_bin:$PATH" \
      REPLAY_FIRST_DONE="$replay_root/first.done" \
      REPLAY_SECOND_DONE="$replay_root/second.done" \
      REPLAY_COLLISION_DESTINATION="${REPLAY_COLLISION_DESTINATION:-}" \
      "$toolbox_scripts/generate-report-image-batch.sh" \
        'ai-reports/<report-slug>' \
        'replay style' \
        '01-architecture.png:<prompt 1>' \
        '02-dataflow.png:<prompt 2>'
  ) > "$fixture_root/$replay_name.stdout" 2> "$fixture_root/$replay_name.stderr"
}

image_all_failed_bin="$fixture_root/image-all-failed-bin"
mkdir -p "$image_all_failed_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'case "$image_prompt" in *"<prompt 1>"*) : > "$REPLAY_FIRST_DONE" ;; *) : > "$REPLAY_SECOND_DONE" ;; esac' \
  'printf partial > "$output_path"' \
  'exit 9' \
  > "$image_all_failed_bin/imagegen"
chmod +x "$image_all_failed_bin/imagegen"
run_ai_report_batch_replay image-all-failed "$image_all_failed_bin" \
  || fail_case 'ai-report all-failed batch replay returned nonzero instead of falling back'
[ ! -e "$fixture_root/image-all-failed/ai-reports/<report-slug>/generated" ] \
  || fail_case 'ai-report all-failed batch replay published an empty generated/ directory'
leaked_private_paths="$(find "$fixture_root/image-all-failed/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit)" \
  || fail_case 'ai-report all-failed batch replay could not search the report directory for invocation-private staging'
[ -n "$leaked_private_paths" ] \
  && fail_case 'ai-report all-failed batch replay leaked invocation-private staging'
[ -e "$fixture_root/image-all-failed/first.done" ] && [ -e "$fixture_root/image-all-failed/second.done" ] \
  || fail_case 'ai-report all-failed batch replay did not wait for every launched job'
[ ! -s "$fixture_root/image-all-failed.stdout" ] \
  || fail_case 'ai-report all-failed batch replay reported a published directory on stdout'

image_batch_mixed_bin="$fixture_root/image-batch-mixed-bin"
mkdir -p "$image_batch_mixed_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'case "$image_prompt" in *"<prompt 1>"*) printf current-success > "$output_path"; : > "$REPLAY_FIRST_DONE"; exit 0 ;; *) printf partial > "$output_path"; : > "$REPLAY_SECOND_DONE"; exit 9 ;; esac' \
  > "$image_batch_mixed_bin/imagegen"
chmod +x "$image_batch_mixed_bin/imagegen"
run_ai_report_batch_replay image-batch-mixed "$image_batch_mixed_bin" \
  || fail_case 'ai-report mixed batch replay returned nonzero'
[ "$(cat "$fixture_root/image-batch-mixed/ai-reports/<report-slug>/generated/01-architecture.png" 2>/dev/null)" = current-success ] \
  && [ ! -e "$fixture_root/image-batch-mixed/ai-reports/<report-slug>/generated/02-dataflow.png" ] \
  || fail_case 'ai-report mixed batch replay did not publish only the status-backed successful image'
leaked_private_paths="$(find "$fixture_root/image-batch-mixed/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit)" \
  || fail_case 'ai-report mixed batch replay could not search the report directory for invocation-private staging'
[ -n "$leaked_private_paths" ] \
  && fail_case 'ai-report mixed batch replay leaked invocation-private staging'
[ -e "$fixture_root/image-batch-mixed/first.done" ] && [ -e "$fixture_root/image-batch-mixed/second.done" ] \
  || fail_case 'ai-report mixed batch replay did not wait for every launched job'
# The caller learns that generated/ was published from stdout — a script cannot set the
# caller's own GEN variable, so the published path is the batch's success signal.
mixed_published_directory="$(cat "$fixture_root/image-batch-mixed.stdout" 2>/dev/null)"
[ -n "$mixed_published_directory" ] \
  && [ "$mixed_published_directory" -ef "$fixture_root/image-batch-mixed/ai-reports/<report-slug>/generated" ] \
  || fail_case 'ai-report mixed batch replay did not report the published generated/ directory on stdout'

# generate-report-image-batch, interrupted: the batch owns the process tree it started, so
# signalling the caller must leave no helper — and no helper's descendant — alive, and
# must reap them before staging is removed. A file-only cleanup passes the directory
# assertions below while backends keep running against the deleted stage.
image_interrupt_batch_bin="$fixture_root/image-interrupt-batch-bin"
mkdir -p "$image_interrupt_batch_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf partial > "$output_path"' \
  '( fixture_deadline=$((SECONDS + 30)); while (( SECONDS < fixture_deadline )); do sleep 0.2; done ) &' \
  'helper_descendant_pid=$!' \
  'case "$image_prompt" in' \
  '  *"<prompt 1>"*) printf "%s\n" "$$" > "$REPLAY_FIRST_PID"; printf "%s\n" "$helper_descendant_pid" > "$REPLAY_FIRST_CHILD_PID"; : > "$REPLAY_FIRST_DONE" ;;' \
  '  *) printf "%s\n" "$$" > "$REPLAY_SECOND_PID"; printf "%s\n" "$helper_descendant_pid" > "$REPLAY_SECOND_CHILD_PID"; : > "$REPLAY_SECOND_DONE" ;;' \
  'esac' \
  'fixture_deadline=$((SECONDS + 30)); while (( SECONDS < fixture_deadline )); do sleep 0.2; done' \
  > "$image_interrupt_batch_bin/imagegen"
chmod +x "$image_interrupt_batch_bin/imagegen"

interrupt_batch_root="$fixture_root/image-interrupt-batch"
mkdir -p "$interrupt_batch_root/ai-reports/<report-slug>"
(
  cd "$interrupt_batch_root" \
    && exec env PATH="$image_interrupt_batch_bin:$PATH" \
      REPLAY_FIRST_DONE="$interrupt_batch_root/first.done" \
      REPLAY_SECOND_DONE="$interrupt_batch_root/second.done" \
      REPLAY_FIRST_PID="$interrupt_batch_root/first.pid" \
      REPLAY_SECOND_PID="$interrupt_batch_root/second.pid" \
      REPLAY_FIRST_CHILD_PID="$interrupt_batch_root/first.child.pid" \
      REPLAY_SECOND_CHILD_PID="$interrupt_batch_root/second.child.pid" \
      "$toolbox_scripts/generate-report-image-batch.sh" \
        'ai-reports/<report-slug>' \
        'replay style' \
        '01-architecture.png:<prompt 1>' \
        '02-dataflow.png:<prompt 2>'
) &
interrupt_batch_pid=$!
interrupt_ready_ticks=0
while [ "$interrupt_ready_ticks" -lt 200 ]; do
  [ -s "$interrupt_batch_root/first.child.pid" ] && [ -s "$interrupt_batch_root/second.child.pid" ] && break
  sleep 0.05
  interrupt_ready_ticks=$((interrupt_ready_ticks + 1))
done
[ -s "$interrupt_batch_root/first.child.pid" ] && [ -s "$interrupt_batch_root/second.child.pid" ] \
  || fail_case 'ai-report interrupted batch replay never reached both running helpers'
for interrupt_pid_file in first.pid second.pid first.child.pid second.child.pid; do
  background_process_ids="$background_process_ids $(cat "$interrupt_batch_root/$interrupt_pid_file" 2>/dev/null || true)"
done
kill -TERM "$interrupt_batch_pid" 2>/dev/null || true
# A batch that will not finish must fail this case, never hang the gate — the gate's exit
# status is the only thing anyone reads, and two stalls were traced to an interruption probe
# that waited forever (REQ-325). The batch swaps its interruption traps around each helper
# launch, so a regression there is exactly what would leave this wait blocked; the watchdog
# turns that into a KILL, which the status assertion below reads as a failure.
( sleep 10; kill -KILL "$interrupt_batch_pid" 2>/dev/null || true ) &
interrupt_batch_watchdog_pid=$!
interrupt_batch_status=0
wait "$interrupt_batch_pid" || interrupt_batch_status=$?
kill "$interrupt_batch_watchdog_pid" 2>/dev/null || true
wait "$interrupt_batch_watchdog_pid" 2>/dev/null || true
[ "$interrupt_batch_status" -eq 143 ] \
  || fail_case "ai-report interrupted batch replay exited $interrupt_batch_status instead of the TERM status 143"
interrupt_survivor_ticks=0
while [ "$interrupt_survivor_ticks" -lt 40 ]; do
  interrupt_survivors=0
  for interrupt_pid_file in first.pid second.pid first.child.pid second.child.pid; do
    interrupt_recorded_pid="$(cat "$interrupt_batch_root/$interrupt_pid_file" 2>/dev/null || true)"
    [ -n "$interrupt_recorded_pid" ] || continue
    kill -0 "$interrupt_recorded_pid" 2>/dev/null && interrupt_survivors=$((interrupt_survivors + 1))
  done
  [ "$interrupt_survivors" -eq 0 ] && break
  sleep 0.05
  interrupt_survivor_ticks=$((interrupt_survivor_ticks + 1))
done
[ "$interrupt_survivors" -eq 0 ] \
  || fail_case "ai-report interrupted batch replay left $interrupt_survivors helper process(es) or descendant(s) alive"
leaked_private_paths="$(find "$interrupt_batch_root/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit)" \
  || fail_case 'ai-report interrupted batch replay could not search the report directory for invocation-private staging'
[ -n "$leaked_private_paths" ] \
  && fail_case 'ai-report interrupted batch replay leaked invocation-private staging'
[ ! -e "$interrupt_batch_root/ai-reports/<report-slug>/generated" ] \
  || fail_case 'ai-report interrupted batch replay published generated/'

# The retired shell implementation's ps/sleep liveness seams are covered by
# report_image_process_test.go against the canonical owned-process implementation.
# They are retained below as historical characterization, but are not applicable
# to the launcher because it no longer resolves either command.
if false; then
# generate-report-image-batch: a controlled ps seam keeps the batch group zombie-only
# while preserving a real descendant-bearing backend. The helper's own backend group can
# independently report live for the TERM-deaf KILL check below.
image_batch_liveness_bin="$fixture_root/image-batch-liveness-bin"
mkdir -p "$image_batch_liveness_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf partial > "$output_path"' \
  'case "$image_prompt" in *term-deaf*) printf "%s\n" "$$" > "$IMAGE_BATCH_LIVENESS_BACKEND_PID"; ( trap "" TERM; fixture_deadline=$((SECONDS + 30)); while (( SECONDS < fixture_deadline )); do sleep 1; done ) & trap "" TERM ;; *) ( fixture_deadline=$((SECONDS + 30)); while (( SECONDS < fixture_deadline )); do sleep 1; done ) & trap "sleep 0.3; exit 143" TERM ;; esac' \
  'printf "%s\n" "$!" > "$IMAGE_BATCH_LIVENESS_DESCENDANT_PID"' \
  ': > "$IMAGE_BATCH_LIVENESS_READY"' \
  'fixture_deadline=$((SECONDS + 30)); while (( SECONDS < fixture_deadline )); do sleep 1; done' \
  > "$image_batch_liveness_bin/imagegen"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'printf "%s\n" "$1" >> "$IMAGE_BATCH_LIVENESS_SLEEP_LOG"' \
  'exec /bin/sleep "$@"' \
  > "$image_batch_liveness_bin/sleep"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'case "$1" in' \
  '  -o) printf "%s\n" "$4"; printf "%s\n" "$4" >> "$IMAGE_BATCH_LIVENESS_GROUP_IDS" ;;' \
  '  -eo)' \
  '    backend_group_id="$(cat "$IMAGE_BATCH_LIVENESS_BACKEND_PID" 2>/dev/null || printf "")"' \
  '    while IFS= read -r group_id; do' \
  '      process_state="$IMAGE_BATCH_LIVENESS_FAKE_STATE"' \
  '      [ "$group_id" != "$backend_group_id" ] || process_state="$IMAGE_BATCH_LIVENESS_BACKEND_STATE"' \
  '      printf "%s %s\n" "$group_id" "$process_state"' \
  '    done < "$IMAGE_BATCH_LIVENESS_GROUP_IDS"' \
  '    ;;' \
  '  *) exit 2 ;;' \
  'esac' \
  > "$image_batch_liveness_bin/ps"
chmod +x "$image_batch_liveness_bin/imagegen" "$image_batch_liveness_bin/sleep" "$image_batch_liveness_bin/ps"
image_batch_liveness_runtime_tools="$fixture_root/image-batch-liveness-tools"
mkdir -p "$image_batch_liveness_runtime_tools"
for runtime_script in generate-report-image.sh generate-report-image-batch.sh; do
  sed 's|process_status_command=/bin/ps|process_status_command="${DO_WORK_PROCESS_STATUS_COMMAND}"|' \
    "$toolbox_scripts/$runtime_script" > "$image_batch_liveness_runtime_tools/$runtime_script"
  chmod +x "$image_batch_liveness_runtime_tools/$runtime_script"
done

run_batch_liveness_replay() {
  batch_liveness_name="$1"
  batch_liveness_prompt="$2"
  batch_liveness_fake_state="$3"
  batch_liveness_backend_state="$4"
  batch_liveness_root="$fixture_root/$batch_liveness_name"
  mkdir -p "$batch_liveness_root/ai-reports/<report-slug>"
  (
    cd "$batch_liveness_root" \
      && exec env PATH="$image_batch_liveness_bin:$PATH" \
        IMAGE_BATCH_LIVENESS_READY="$batch_liveness_root/ready" \
        IMAGE_BATCH_LIVENESS_DESCENDANT_PID="$batch_liveness_root/descendant.pid" \
        IMAGE_BATCH_LIVENESS_SLEEP_LOG="$batch_liveness_root/sleeps" \
        IMAGE_BATCH_LIVENESS_GROUP_IDS="$batch_liveness_root/groups" \
        IMAGE_BATCH_LIVENESS_BACKEND_PID="$batch_liveness_root/backend.pid" \
        IMAGE_BATCH_LIVENESS_FAKE_STATE="$batch_liveness_fake_state" \
        IMAGE_BATCH_LIVENESS_BACKEND_STATE="$batch_liveness_backend_state" \
        DO_WORK_PROCESS_STATUS_COMMAND="$image_batch_liveness_bin/ps" \
        "$image_batch_liveness_runtime_tools/generate-report-image-batch.sh" \
          'ai-reports/<report-slug>' style "01-image.png:$batch_liveness_prompt"
  ) >/dev/null 2>&1 &
  batch_liveness_pid=$!
  background_process_ids="$background_process_ids $batch_liveness_pid"
  batch_liveness_ready_ticks=0
  while [ ! -e "$batch_liveness_root/ready" ] && [ "$batch_liveness_ready_ticks" -lt 200 ]; do
    sleep 0.01
    batch_liveness_ready_ticks=$((batch_liveness_ready_ticks + 1))
  done
  [ -e "$batch_liveness_root/ready" ] \
    || { fail_case "generate-report-image-batch $batch_liveness_name case never reached the descendant-bearing backend"; return; }
  batch_liveness_descendant_pid="$(cat "$batch_liveness_root/descendant.pid")"
  background_process_ids="$background_process_ids $batch_liveness_descendant_pid"
  kill -TERM "$batch_liveness_pid" 2>/dev/null || true
  batch_liveness_status=0
  wait "$batch_liveness_pid" || batch_liveness_status=$?
  [ "$batch_liveness_status" -eq 143 ] \
    || fail_case "generate-report-image-batch $batch_liveness_name case returned $batch_liveness_status instead of the TERM status 143"
}

run_batch_liveness_replay liveness-zombies liveness-zombies Z Z
batch_liveness_grace_ticks="$(grep -cx '0.1' "$fixture_root/liveness-zombies/sleeps" 2>/dev/null || true)"
[ "$batch_liveness_grace_ticks" -le 1 ] \
  || fail_case "generate-report-image-batch exited-group case waited $batch_liveness_grace_ticks grace ticks instead of stopping promptly after the TERM-obedient helper and descendant exited"

run_batch_liveness_replay liveness-term-deaf term-deaf Z S
batch_liveness_grace_ticks="$(grep -cx '0.1' "$fixture_root/liveness-term-deaf/sleeps" 2>/dev/null || true)"
[ "$batch_liveness_grace_ticks" -eq 10 ] \
  || fail_case "generate-report-image-batch TERM-deaf case used $batch_liveness_grace_ticks grace ticks instead of the helper's full 10-tick budget before KILL"
batch_liveness_survivor_ticks=0
while process_id_is_live "$(cat "$fixture_root/liveness-term-deaf/descendant.pid")" \
  && [ "$batch_liveness_survivor_ticks" -lt 40 ]; do
  sleep 0.05
  batch_liveness_survivor_ticks=$((batch_liveness_survivor_ticks + 1))
done
process_id_is_live "$(cat "$fixture_root/liveness-term-deaf/descendant.pid")" \
  && fail_case 'generate-report-image-batch TERM-deaf case left the TERM-deaf descendant alive after KILL'
fi

# generate-report-image-batch, destination appears at the final boundary: `mv` treats an
# existing directory as a container, so the check-then-rename window can nest the
# private stage inside a colliding generated/ and still exit zero. The mv shim below
# creates the destination in exactly that window.
image_publish_collision_bin="$fixture_root/image-publish-collision-bin"
mkdir -p "$image_publish_collision_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'case "$image_prompt" in *"<prompt 1>"*) : > "$REPLAY_FIRST_DONE" ;; *) : > "$REPLAY_SECOND_DONE" ;; esac' \
  'printf current-success > "$output_path"' \
  'if [ -n "${REPLAY_COLLISION_DESTINATION:-}" ] && mkdir "$REPLAY_COLLISION_DESTINATION" 2>/dev/null; then printf owned-by-someone-else > "$REPLAY_COLLISION_DESTINATION/keep.txt"; fi' \
  > "$image_publish_collision_bin/imagegen"
chmod +x "$image_publish_collision_bin/imagegen"
publish_collision_generated="$fixture_root/image-publish-collision/ai-reports/<report-slug>/generated"
REPLAY_COLLISION_DESTINATION="$publish_collision_generated" \
  run_ai_report_batch_replay image-publish-collision "$image_publish_collision_bin" \
  && fail_case 'ai-report publish-collision replay reported success after the destination appeared'
[ "$(cat "$publish_collision_generated/keep.txt" 2>/dev/null)" = owned-by-someone-else ] \
  || fail_case 'ai-report publish-collision replay did not preserve the colliding destination byte-for-byte'
[ "$(ls -A "$publish_collision_generated" 2>/dev/null)" = keep.txt ] \
  || fail_case 'ai-report publish-collision replay left its staged batch nested inside the colliding destination'
leaked_private_paths="$(find "$fixture_root/image-publish-collision/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit)" \
  || fail_case 'ai-report publish-collision replay could not search the report directory for invocation-private staging'
[ -n "$leaked_private_paths" ] \
  && fail_case 'ai-report publish-collision replay leaked invocation-private staging'

# The Go command installs its signal context before allocating staging; that ordering
# and pre-cancelled launch behavior are exercised directly in report-image Go tests.
if false; then
# generate-report-image-batch: a staging directory created before the interruption traps
# exist is a directory no trap owns — the signal takes its default action, the EXIT trap
# never runs, and .generated.staging.* is left in the user's report directory. The mktemp
# shim below holds the batch inside exactly that window, so the interruption is an event
# rather than a race; the stub backend exits at once so a dropped interruption ends the
# batch instead of hanging this case.
image_early_batch_bin="$fixture_root/image-early-batch-bin"
mkdir -p "$image_early_batch_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf "%s\n" "$$" >> "$EARLY_BATCH_HELPER_RAN"' \
  'printf partial > "$output_path"' \
  'exit 9' \
  > "$image_early_batch_bin/imagegen"
chmod +x "$image_early_batch_bin/imagegen"
# The batch resolves mktemp through PATH, so this shim can create the staging directory the
# batch asked for, announce that it exists, and return only once the case has fired its
# interruption. Every other mktemp call — the per-image helper's own staging — is handed to
# the real one untouched, so the shim cannot change what the rest of the run does.
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'case "$*" in' \
  '  *.generated.staging.*) ;;' \
  '  *) exec "$EARLY_BATCH_REAL_MKTEMP" "$@" ;;' \
  'esac' \
  'printf "%s\n" "$$" > "$EARLY_BATCH_SHIM_PID"' \
  'early_batch_staging="$("$EARLY_BATCH_REAL_MKTEMP" "$@")" || exit 1' \
  'printf "%s\n" "$early_batch_staging" > "$EARLY_BATCH_STAGE_LOG"' \
  'early_batch_shim_ticks=0' \
  'while [ ! -e "$EARLY_BATCH_GO" ] && [ "$early_batch_shim_ticks" -lt 400 ]; do' \
  '  sleep 0.05' \
  '  early_batch_shim_ticks=$((early_batch_shim_ticks + 1))' \
  'done' \
  'printf "%s\n" "$early_batch_staging"' \
  > "$image_early_batch_bin/mktemp"
chmod +x "$image_early_batch_bin/mktemp"

early_batch_root="$fixture_root/image-early-batch"
mkdir -p "$early_batch_root/ai-reports/<report-slug>"
(
  cd "$early_batch_root" \
    && exec env PATH="$image_early_batch_bin:$PATH" \
      EARLY_BATCH_REAL_MKTEMP="$(command -v mktemp)" \
      EARLY_BATCH_STAGE_LOG="$early_batch_root/staging.path" \
      EARLY_BATCH_SHIM_PID="$early_batch_root/shim.pid" \
      EARLY_BATCH_GO="$early_batch_root/go" \
      EARLY_BATCH_HELPER_RAN="$early_batch_root/helper.ran" \
      "$toolbox_scripts/generate-report-image-batch.sh" \
        'ai-reports/<report-slug>' \
        'replay style' \
        '01-architecture.png:<prompt 1>'
) >/dev/null 2>&1 &
early_batch_pid=$!
background_process_ids="$background_process_ids $early_batch_pid"
early_batch_ready_ticks=0
while [ ! -s "$early_batch_root/staging.path" ] && [ "$early_batch_ready_ticks" -lt 500 ]; do
  sleep 0.01
  early_batch_ready_ticks=$((early_batch_ready_ticks + 1))
done
if [ ! -s "$early_batch_root/staging.path" ]; then
  fail_case 'ai-report early-interrupted batch replay never reached its staging window'
  : > "$early_batch_root/go"
  wait "$early_batch_pid" 2>/dev/null || true
else
  background_process_ids="$background_process_ids $(cat "$early_batch_root/shim.pid" 2>/dev/null || true)"
  kill -TERM "$early_batch_pid" 2>/dev/null || true
  # Released only after the signal is on its way: a batch that stages under its traps defers
  # the interruption until this mktemp returns, so the shim has to be able to return.
  : > "$early_batch_root/go"
  early_batch_status=0
  wait "$early_batch_pid" || early_batch_status=$?
  [ "$early_batch_status" -eq 143 ] \
    || fail_case "ai-report early-interrupted batch replay exited $early_batch_status instead of the TERM status 143"
  leaked_private_paths="$(find "$early_batch_root/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit)" \
    || fail_case 'ai-report early-interrupted batch replay could not search the report directory for invocation-private staging'
  [ -n "$leaked_private_paths" ] \
    && fail_case 'ai-report early-interrupted batch replay leaked invocation-private staging'
  [ ! -e "$early_batch_root/ai-reports/<report-slug>/generated" ] \
    || fail_case 'ai-report early-interrupted batch replay published generated/'
  [ ! -e "$early_batch_root/helper.ran" ] \
    || fail_case 'ai-report early-interrupted batch replay launched a helper after the interruption'
fi
fi

prescribed_shell_finish

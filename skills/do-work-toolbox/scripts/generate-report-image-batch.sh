#!/usr/bin/env bash
# Generate one report's image batch in parallel and publish generated/ only after verified success.
set -u

if [ "$#" -lt 3 ]; then
  printf 'Usage: %s <report-directory> <style-brief> <target-name>:<prompt> [<target-name>:<prompt> ...]\n' "$0" >&2
  exit 2
fi

report_directory_argument="$1"
style_brief="$2"
shift 2

# Each remaining argument pairs a bare target filename with its prompt, split on the
# first colon only — prompts contain colons, target names may not contain a slash
# because the batch joins them to its own private staging directory.
for image_specification in "$@"; do
  case "$image_specification" in
    *:*) ;;
    *) printf 'Image argument must be <target-name>:<prompt>: %s\n' "$image_specification" >&2; exit 2 ;;
  esac
  case "${image_specification%%:*}" in
    ''|*/*) printf 'Image target must be a bare filename: %s\n' "${image_specification%%:*}" >&2; exit 2 ;;
  esac
done

report_directory="$(cd "$report_directory_argument" && pwd)" || exit 1
generated_directory="$report_directory/generated"
[ ! -e "$generated_directory" ] || { printf 'REFUSING: generated/ already exists\n' >&2; exit 1; }

# The per-image helper is a sibling of this script, not a PATH lookup: callers hand the
# batch a report directory, not an installation root.
script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" || exit 1
report_image_helper="$script_directory/generate-report-image.sh"

cleanup_report_image_stage() {
  if [ -n "${image_generation_stage:-}" ]; then
    rm -rf -- "$image_generation_stage"
  fi
}
trap cleanup_report_image_stage EXIT

image_generation_pids=()
image_generation_groups=()
image_generation_targets=()
# Reach ps by absolute path and strip whitespace with parameter expansion: a caller may
# hand this batch a minimal PATH, and a bare `ps` would silently fail there and degrade
# every launch to bare-PID signalling. Without ps no group can be proved, which the
# verification below already treats as bare-PID-only signalling.
if [ -x /bin/ps ]; then
  process_status_command=/bin/ps
elif [ -x /usr/bin/ps ]; then
  process_status_command=/usr/bin/ps
else
  process_status_command=""
fi
caller_process_group_id=""
if [ -n "$process_status_command" ]; then
  caller_process_group_id="$($process_status_command -o pgid= -p "$$" 2>/dev/null || true)"
  caller_process_group_id="${caller_process_group_id//[[:space:]]/}"
fi

deferred_batch_interrupt_status=""
# The traps below act on an interruption; this is what they do, named once so the deferral
# can re-raise exactly the same thing.
interrupt_report_image_batch() {
  terminate_report_image_batch
  exit "$1"
}

# A trapped signal runs between commands, so an interruption arriving between a helper launch
# and the appends that publish its PID and process group would reach
# terminate_report_image_batch with nothing in the arrays: the helper and its backend outlive
# the batch, and the EXIT trap then deletes the staging directory they are writing into.
# Defer the interruption across that window instead — record the status, let the handles be
# recorded, then act on it with the whole batch reachable.
defer_interrupts_across_helper_launch() {
  deferred_batch_interrupt_status=""
  trap 'deferred_batch_interrupt_status=129' HUP
  trap 'deferred_batch_interrupt_status=130' INT
  trap 'deferred_batch_interrupt_status=143' TERM
}

resume_interrupts_after_helper_launch() {
  trap 'interrupt_report_image_batch 129' HUP
  trap 'interrupt_report_image_batch 130' INT
  trap 'interrupt_report_image_batch 143' TERM
  [ -z "$deferred_batch_interrupt_status" ] \
    || interrupt_report_image_batch "$deferred_batch_interrupt_status"
}

launch_report_image() {
  image_target="$1"
  image_description="$2"
  defer_interrupts_across_helper_launch
  set -m   # job control gives each helper its own process group, so its descendants are reachable
  "$report_image_helper" "$image_target" "$style_brief" "$image_description" &
  image_helper_pid=$!
  set +m
  image_helper_group=""
  if [ -n "$process_status_command" ]; then
    image_helper_group="$($process_status_command -o pgid= -p "$image_helper_pid" 2>/dev/null || true)"
    image_helper_group="${image_helper_group//[[:space:]]/}"
  fi
  case "$image_helper_group" in ''|*[!0-9]*) image_helper_group="" ;; esac
  # Signal a group only when this helper leads one of its own; otherwise fall back to
  # the bare PID. Never signal the caller's group — that would kill the report run.
  if [ "$image_helper_group" != "$image_helper_pid" ] || [ "$image_helper_group" = "$caller_process_group_id" ]; then
    image_helper_group=""
  fi
  image_generation_pids[${#image_generation_pids[@]}]="$image_helper_pid"
  image_generation_groups[${#image_generation_groups[@]}]="$image_helper_group"
  image_generation_targets[${#image_generation_targets[@]}]="$image_target"
  resume_interrupts_after_helper_launch
}

report_image_batch_is_alive() {
  batch_index=0
  while [ "$batch_index" -lt "${#image_generation_pids[@]}" ]; do
    if [ -n "${image_generation_groups[$batch_index]}" ]; then
      kill -0 -- "-${image_generation_groups[$batch_index]}" 2>/dev/null && return 0
    else
      kill -0 "${image_generation_pids[$batch_index]}" 2>/dev/null && return 0
    fi
    batch_index=$((batch_index + 1))
  done
  return 1
}

signal_report_image_batch() {
  batch_signal="$1"
  batch_index=0
  while [ "$batch_index" -lt "${#image_generation_pids[@]}" ]; do
    if [ -n "${image_generation_groups[$batch_index]}" ]; then
      kill -"$batch_signal" -- "-${image_generation_groups[$batch_index]}" 2>/dev/null || true
    else
      kill -"$batch_signal" "${image_generation_pids[$batch_index]}" 2>/dev/null || true
    fi
    batch_index=$((batch_index + 1))
  done
}

terminate_report_image_batch() {
  signal_report_image_batch TERM
  batch_grace_ticks=0
  while report_image_batch_is_alive && [ "$batch_grace_ticks" -lt 10 ]; do
    sleep 0.1
    batch_grace_ticks=$((batch_grace_ticks + 1))
  done
  if report_image_batch_is_alive; then
    signal_report_image_batch KILL
  fi
  batch_index=0
  while [ "$batch_index" -lt "${#image_generation_pids[@]}" ]; do
    wait "${image_generation_pids[$batch_index]}" 2>/dev/null || true
    batch_index=$((batch_index + 1))
  done
}
# Reap the batch first, then let the EXIT trap remove staging: an interrupted batch
# must not leave helpers writing into a directory it is about to delete.
trap 'interrupt_report_image_batch 129' HUP
trap 'interrupt_report_image_batch 130' INT
trap 'interrupt_report_image_batch 143' TERM

# Created only once the traps above can clean it up. Staged earlier — as it was until an
# early-interruption probe caught it — an interruption arriving before those traps exist
# takes its default action, so the EXIT trap never runs and .generated.staging.* is left
# behind in the caller's report directory.
image_generation_stage="$(umask 077; mktemp -d "$report_directory/.generated.staging.XXXXXX")" || exit 1

for image_specification in "$@"; do
  launch_report_image \
    "$image_generation_stage/${image_specification%%:*}" \
    "${image_specification#*:}"
done

image_generation_statuses=()
image_index=0
while [ "$image_index" -lt "${#image_generation_pids[@]}" ]; do
  image_status=0
  wait "${image_generation_pids[$image_index]}" || image_status=$?
  image_generation_statuses[$image_index]="$image_status"
  image_index=$((image_index + 1))
done

# An image is current only when its own helper status is zero and its staged target is
# non-empty; presence alone never proves freshness.
image_generation_success_count=0
image_index=0
while [ "$image_index" -lt "${#image_generation_targets[@]}" ]; do
  image_target="${image_generation_targets[$image_index]}"
  if [ "${image_generation_statuses[$image_index]}" -ne 0 ] || [ ! -s "$image_target" ]; then
    rm -f -- "$image_target"
    printf 'MISSING: %s → fall back to SVG/Mermaid for that section\n' "${image_target##*/}" >&2
  else
    image_generation_success_count=$((image_generation_success_count + 1))
  fi
  image_index=$((image_index + 1))
done

if [ "$image_generation_success_count" -gt 0 ]; then
  [ ! -e "$generated_directory" ] || { printf 'REFUSING: generated/ appeared before publication\n' >&2; exit 1; }
  mv "$image_generation_stage" "$generated_directory" || exit 1
  # `mv` treats an existing destination directory as a container, so a generated/ that
  # appeared after the check above would leave the stage nested inside it and still
  # report success. Verify the rename actually published, and fail closed if it nested.
  nested_image_generation_stage="$generated_directory/${image_generation_stage##*/}"
  if [ -e "$nested_image_generation_stage" ]; then
    rm -rf -- "$nested_image_generation_stage"
    image_generation_stage=""
    printf 'REFUSING: generated/ appeared during publication — staged batch discarded, existing generated/ left unchanged\n' >&2
    exit 1
  fi
  image_generation_stage=""
  # The published directory is the batch's success signal; absolute helper outputs were
  # published here, and the report HTML still embeds relative generated/… paths.
  printf '%s\n' "$generated_directory"
else
  # An all-failed batch is not an error: it removes its exact private directory, prints
  # nothing, and returns zero so the caller falls back to SVG/Mermaid.
  cleanup_report_image_stage
  image_generation_stage=""
fi
trap - EXIT HUP INT TERM
exit 0

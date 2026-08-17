#!/usr/bin/env bash
# Generate one verified report image through the allowed backend chain.
set -u

if [ "$#" -ne 3 ]; then
  printf 'Usage: %s <absolute-output-png> <style-brief> <sanitized-description>\n' "$0" >&2
  exit 2
fi

output_path="$1"
style_brief="$2"
visual_description="$3"
case "$output_path" in
  /*) ;;
  *) printf 'Output path must be absolute.\n' >&2; exit 2 ;;
esac

output_directory="${output_path%/*}"
output_filename="${output_path##*/}"
if [ ! -d "$output_directory" ] || [ -z "$output_filename" ]; then
  printf 'Output directory must already exist.\n' >&2
  exit 2
fi

staged_output_path="$(mktemp "$output_directory/.${output_filename}.generating.XXXXXX")" || exit 2
agent_directory=""
backend_process_id=""
backend_process_group_id=""
# Callers hand this script a minimal PATH, so reach ps by absolute path and strip
# whitespace with parameter expansion rather than `tr`. Without ps no group can be
# proved, which the verification below already treats as bare-PID-only signalling.
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

# A backend launched under job control leads its own process group, which is the only
# handle that reaches its descendants. Record that group only when it is provably the
# backend's own and not this script's — an unverified group means bare-PID signalling.
record_backend_process_group() {
  backend_process_group_id=""
  if [ -n "$process_status_command" ]; then
    backend_process_group_id="$($process_status_command -o pgid= -p "$backend_process_id" 2>/dev/null || true)"
    backend_process_group_id="${backend_process_group_id//[[:space:]]/}"
  fi
  case "$backend_process_group_id" in ''|*[!0-9]*) backend_process_group_id="" ;; esac
  if [ "$backend_process_group_id" != "$backend_process_id" ] \
    || [ "$backend_process_group_id" = "$caller_process_group_id" ]; then
    backend_process_group_id=""
  fi
}

backend_process_is_alive() {
  if [ -n "$backend_process_group_id" ]; then
    kill -0 -- "-$backend_process_group_id" 2>/dev/null
  else
    kill -0 "$backend_process_id" 2>/dev/null
  fi
}

signal_backend_process() {
  backend_signal="$1"
  if [ -n "$backend_process_group_id" ]; then
    kill -"$backend_signal" -- "-$backend_process_group_id" 2>/dev/null || true
  else
    kill -"$backend_signal" "$backend_process_id" 2>/dev/null || true
  fi
}

terminate_backend_process() {
  [ -n "$backend_process_id" ] || return 0
  signal_backend_process TERM
  backend_grace_ticks=0
  while backend_process_is_alive && [ "$backend_grace_ticks" -lt 10 ]; do
    sleep 0.1
    backend_grace_ticks=$((backend_grace_ticks + 1))
  done
  if backend_process_is_alive; then
    signal_backend_process KILL
  fi
  wait "$backend_process_id" 2>/dev/null || true
}

cleanup_report_image_paths() {
  # Reap the process tree first: nothing this invocation started may still be writing
  # into the staged file or the agent directory removed below.
  terminate_backend_process
  if [ -n "$staged_output_path" ]; then
    rm -f "$staged_output_path"
  fi
  if [ -n "$agent_directory" ]; then
    rm -rf "$agent_directory"
  fi
}
trap cleanup_report_image_paths EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

publish_report_image() {
  [ -s "$staged_output_path" ] || return 1
  mv "$staged_output_path" "$output_path" || return 1
  # `mv` treats an existing destination directory as a container, so an output path
  # occupied by a directory nests the staged image inside it and still exits zero. No
  # backend can repair that, so discard only the nested stage and fail the invocation
  # outright rather than returning to a fallback that would re-stage over nothing.
  nested_staged_output_path="$output_path/${staged_output_path##*/}"
  if [ -e "$nested_staged_output_path" ]; then
    rm -f -- "$nested_staged_output_path"
    staged_output_path=""
    printf 'REFUSING: %s is a directory — staged image discarded, existing directory left unchanged\n' "$output_path" >&2
    exit 1
  fi
  staged_output_path=""
}

image_prompt="$(printf '%s Content: %s' "$style_brief" "$visual_description")"
if command -v imagegen >/dev/null 2>&1; then
  set -m   # job control gives the backend its own process group, so its descendants are reachable
  imagegen --output "$staged_output_path" --prompt "$image_prompt" >/dev/null 2>&1 &
  backend_process_id=$!
  set +m
  record_backend_process_group
  wait "$backend_process_id"
  direct_backend_status=$?
  backend_process_id=""
  backend_process_group_id=""
  if [ "$direct_backend_status" -eq 0 ] && publish_report_image; then
    exit 0
  fi
fi

[ "${DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND:-0}" = "1" ] || exit 1
command -v codex >/dev/null 2>&1 || exit 1
: > "$staged_output_path"
agent_directory="$(mktemp -d "${TMPDIR:-/tmp}/do-work-ai-report-image.XXXXXX")" || exit 2
chmod 700 "$agent_directory" || exit 2
set -m   # the recorded PID is the wrapping subshell, so only its group reaches codex itself
(
  cd "$agent_directory" || exit 2
  codex exec --dangerously-bypass-approvals-and-sandbox \
    "Generate a 16:9 image and save the PNG EXACTLY to ./generated.png. $image_prompt" >/dev/null 2>&1
) &
backend_process_id=$!
set +m
record_backend_process_group
wait "$backend_process_id"
agent_status=$?
backend_process_id=""
backend_process_group_id=""
if [ "$agent_status" -eq 0 ] && [ -s "$agent_directory/generated.png" ]; then
  cp "$agent_directory/generated.png" "$staged_output_path" || exit 2
  publish_report_image
  exit $?
fi
exit 1

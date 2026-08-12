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

cleanup_report_image_paths() {
  if [ -n "$backend_process_id" ]; then
    kill "$backend_process_id" 2>/dev/null || true
    wait "$backend_process_id" 2>/dev/null || true
  fi
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
  staged_output_path=""
}

image_prompt="$(printf '%s Content: %s' "$style_brief" "$visual_description")"
if command -v imagegen >/dev/null 2>&1; then
  imagegen --output "$staged_output_path" --prompt "$image_prompt" >/dev/null 2>&1 &
  backend_process_id=$!
  wait "$backend_process_id"
  direct_backend_status=$?
  backend_process_id=""
  if [ "$direct_backend_status" -eq 0 ] && publish_report_image; then
    exit 0
  fi
fi

[ "${DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND:-0}" = "1" ] || exit 1
command -v codex >/dev/null 2>&1 || exit 1
: > "$staged_output_path"
agent_directory="$(mktemp -d "${TMPDIR:-/tmp}/do-work-ai-report-image.XXXXXX")" || exit 2
chmod 700 "$agent_directory" || exit 2
(
  cd "$agent_directory" || exit 2
  codex exec --dangerously-bypass-approvals-and-sandbox \
    "Generate a 16:9 image and save the PNG EXACTLY to ./generated.png. $image_prompt" >/dev/null 2>&1
) &
backend_process_id=$!
wait "$backend_process_id"
agent_status=$?
backend_process_id=""
if [ "$agent_status" -eq 0 ] && [ -s "$agent_directory/generated.png" ]; then
  cp "$agent_directory/generated.png" "$staged_output_path" || exit 2
  publish_report_image
  exit $?
fi
exit 1

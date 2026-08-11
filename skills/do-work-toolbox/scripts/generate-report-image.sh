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

image_prompt="$(printf '%s Content: %s' "$style_brief" "$visual_description")"
if command -v imagegen >/dev/null 2>&1; then
  imagegen --output "$output_path" --prompt "$image_prompt" >/dev/null 2>&1 || true
  [ -s "$output_path" ] && exit 0
fi

[ "${DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND:-0}" = "1" ] || exit 1
command -v codex >/dev/null 2>&1 || exit 1
agent_directory="$(mktemp -d "${TMPDIR:-/tmp}/do-work-ai-report-image.XXXXXX")" || exit 2
trap 'rm -rf "$agent_directory"' EXIT HUP INT TERM
chmod 700 "$agent_directory" || exit 2
(
  cd "$agent_directory" || exit 2
  codex exec --dangerously-bypass-approvals-and-sandbox \
    "Generate a 16:9 image and save the PNG EXACTLY to ./generated.png. $image_prompt" >/dev/null 2>&1
)
agent_status=$?
if [ "$agent_status" -eq 0 ] && [ -s "$agent_directory/generated.png" ]; then
  cp "$agent_directory/generated.png" "$output_path" || exit 2
fi
[ -s "$output_path" ]

#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
# Compatibility entry point for the canonical typed queue selector.
set -euo pipefail

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi

repository_root='.'
selector_skip='false'
while [ "$#" -gt 0 ]; do
  case "$1" in
    --repo-root)
      if [ "$#" -lt 2 ]; then
        printf 'select-simple-reqs: --repo-root needs a directory\n' >&2
        exit 2
      fi
      repository_root="$2"
      shift 2
      ;;
    --skip-impact-negligible)
	  selector_skip='true'
      shift
      ;;
    *)
      printf 'select-simple-reqs: unrecognized argument %s\n' "$1" >&2
      printf 'usage: select-simple-reqs.sh [--repo-root DIR] [--skip-impact-negligible]\n' >&2
      exit 2
      ;;
  esac
done

selector_status=0
if [ "$selector_skip" = 'true' ]; then
  selector_output="$("$script_directory/do-work-cli.sh" --repo-root "$repository_root" next --simple --skip-impact-negligible)" || selector_status=$?
else
  selector_output="$("$script_directory/do-work-cli.sh" --repo-root "$repository_root" next --simple)" || selector_status=$?
fi
if [ "$selector_status" -ne 0 ]; then
  printf '%s\n' "$selector_output"
  exit "$selector_status"
fi

if printf '%s\n' "$selector_output" | grep -qx 'run_set: '; then
  printf 'No pending REQ currently qualifies for a cheaper model.\n'
fi
printf '%s\n' "$selector_output"

# Keep the legacy diagnostic channel while the compatibility path remains.
printf '%s\n' "$selector_output" | while IFS= read -r output_line; do
  case "$output_line" in
    'skipped SCHEMA-FALLBACK:'*) printf '%s\n' "$output_line" >&2 ;;
  esac
done

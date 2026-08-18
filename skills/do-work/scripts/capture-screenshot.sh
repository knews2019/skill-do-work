#!/usr/bin/env bash
# Verify and install one screenshot without overwriting another dispatch.
set -u

capture_mode="${1:-}"
if [ "$#" -ne 3 ] || { [ "$capture_mode" != "--staged" ] && [ "$capture_mode" != "--keep-source" ]; }; then
  printf 'Usage: %s (--staged|--keep-source) <source> <destination>\n' "$0" >&2
  exit 2
fi

source_path="$2"
destination_path="$3"
destination_directory="$(dirname "$destination_path")"
copy_path=""

cleanup_private_copy() {
  if [ -n "$copy_path" ]; then
    rm -f "$copy_path"
  fi
}
trap cleanup_private_copy EXIT HUP INT TERM

mkdir -p "$destination_directory" || exit 2
copy_path="$(mktemp "${destination_path}.copying.XXXXXX")" || {
  printf 'Screenshot temporary-copy allocation failed; staged source preserved: %s\n' "$source_path" >&2
  exit 2
}

if ! cp "$source_path" "$copy_path" || ! cmp -s "$source_path" "$copy_path"; then
  printf 'Screenshot copy or verification failed; staged source preserved: %s\n' "$source_path" >&2
  exit 1
fi
if ! ln "$copy_path" "$destination_path"; then
  printf 'Screenshot destination already exists or no-clobber install failed; staged source preserved: %s\n' "$source_path" >&2
  exit 1
fi
# `ln` refuses an occupied FILE, and that refusal is where the no-clobber guarantee
# comes from — but it treats a directory in the destination's place as a container,
# nesting the private copy inside it and still exiting zero. Verify the path actually
# written BEFORE anything reads this status as permission to remove the staged source:
# under --staged that source is the dispatch's only copy, so a false success here
# destroys the screenshot it was asked to install. Discard only this invocation's own
# nested copy; the occupying directory is someone else's and is left untouched.
nested_destination_path="$destination_path/${copy_path##*/}"
if [ -e "$nested_destination_path" ]; then
  rm -f -- "$nested_destination_path"
  printf 'REFUSING: %s is a directory — screenshot not installed, staged source preserved: %s\n' \
    "$destination_path" "$source_path" >&2
  exit 1
fi

if ! rm -f "$copy_path"; then
  printf 'Permanent screenshot copy verified, but temporary copy could not be removed: %s\n' "$copy_path" >&2
fi
copy_path=""

if [ "$capture_mode" = "--staged" ]; then
  staging_directory="$(dirname "$source_path")"
  if rm "$source_path"; then
    if [ -d "$staging_directory" ] \
      && [ -z "$(find "$staging_directory" -mindepth 1 -maxdepth 1 -print -quit)" ]; then
      rmdir "$staging_directory" || printf \
        'Permanent screenshot copy verified, but empty staging directory could not be removed: %s\n' \
        "$staging_directory" >&2
    fi
  else
    printf 'Permanent screenshot copy verified, but staged source could not be removed: %s\n' \
      "$source_path" >&2
  fi
fi

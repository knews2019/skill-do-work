#!/usr/bin/env bash
# Stage one deletion while inspecting cached metadata only, never former contents.
set -u

if [ "$#" -ne 1 ]; then
  printf 'Usage: %s <deleted-path>\n' "$0" >&2
  exit 2
fi

deleted_path="$1"
cached_deletion_file="$(mktemp)" || exit 2
trap 'rm -f "$cached_deletion_file"' EXIT HUP INT TERM

read_cached_deletion() {
  cached_status=''
  cached_path=''
  cached_extra=''
  {
    IFS= read -r -d '' cached_status || true
    IFS= read -r -d '' cached_path || true
    IFS= read -r -d '' cached_extra || true
  } < "$cached_deletion_file"
}

git --literal-pathspecs diff --cached --name-status --no-renames -z -- "$deleted_path" > "$cached_deletion_file" || exit 2
read_cached_deletion
if [ "$cached_status" = 'D' ] && [ "$cached_path" = "$deleted_path" ] && [ -z "$cached_extra" ]; then
  exit 0
fi

git --literal-pathspecs add -u -- "$deleted_path" || exit 2
git --literal-pathspecs diff --cached --name-status --no-renames -z -- "$deleted_path" > "$cached_deletion_file" || exit 2
read_cached_deletion
if [ "$cached_status" != 'D' ] || [ "$cached_path" != "$deleted_path" ] || [ -n "$cached_extra" ]; then
  printf 'Cached state is not one exact deletion: %s\n' "$deleted_path" >&2
  exit 2
fi

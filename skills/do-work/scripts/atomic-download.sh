#!/usr/bin/env bash
# Download privately beside the target and publish only a complete file.
set -u

if [ "$#" -ne 2 ]; then
  printf 'Usage: %s <source-url> <target-path>\n' "$0" >&2
  exit 2
fi

source_url="$1"
target_path="$2"
download_path=""

cleanup_download() {
  if [ -n "$download_path" ]; then
    rm -f "$download_path"
  fi
}
trap cleanup_download EXIT HUP INT TERM

download_path="$(mktemp "${target_path}.download.XXXXXX")" || {
  printf 'Could not allocate download beside target: %s\n' "$target_path" >&2
  exit 2
}

curl -fsSL -o "$download_path" "$source_url"
download_status=$?
if [ "$download_status" -ne 0 ]; then
  printf 'Download failed; target was not published: %s\n' "$target_path" >&2
  exit "$download_status"
fi
if ! mv "$download_path" "$target_path"; then
  printf 'Downloaded file could not be published: %s\n' "$target_path" >&2
  exit 2
fi
download_path=""

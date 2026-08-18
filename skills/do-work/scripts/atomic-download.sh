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

# Opt-in credentials: absent or empty leaves the request exactly as it was.
# `set --` rebuilds the optional argument list; the two real arguments were
# already captured above, and "$@" is safe to expand empty under `set -u`.
download_token="${GH_TOKEN:-${GITHUB_TOKEN:-}}"
if [ -n "$download_token" ]; then
  set -- -H "Authorization: Bearer $download_token"
else
  set --
fi

# `--retry` alone treats HTTP 429 as transient from curl 7.51.0; `--retry-all-errors`
# would raise the floor to 7.71 for no gain here.
curl -fsSL --retry 3 --retry-delay 2 --retry-max-time 60 "$@" -o "$download_path" "$source_url"
download_status=$?
if [ "$download_status" -ne 0 ]; then
  printf 'Download failed; target was not published: %s\n' "$target_path" >&2
  exit "$download_status"
fi
if ! mv "$download_path" "$target_path"; then
  printf 'Downloaded file could not be published: %s\n' "$target_path" >&2
  exit 2
fi
# `mv` treats a directory in the destination's place as a container rather than a
# collision, so a target occupied by a directory nests the download inside it and
# still exits zero — and the caller reads that status as proof the file landed.
# Verify the path actually written, discard only this invocation's own nested
# file, and leave the occupying directory exactly as it was.
nested_download_path="$target_path/${download_path##*/}"
if [ -e "$nested_download_path" ]; then
  rm -f -- "$nested_download_path"
  download_path=""
  printf 'REFUSING: %s is a directory — download discarded, existing directory left unchanged\n' "$target_path" >&2
  exit 1
fi
download_path=""

#!/usr/bin/env bash
# Publish one verified portfolio draft to canonical and optional immutable snapshot paths.
set -u

publication_mode="${1:-}"
if { [ "$publication_mode" = "--canonical-only" ] && [ "$#" -ne 3 ]; } \
  || { [ "$publication_mode" = "--with-snapshot" ] && [ "$#" -ne 4 ]; } \
  || { [ "$publication_mode" != "--canonical-only" ] && [ "$publication_mode" != "--with-snapshot" ]; }; then
  printf 'Usage: %s --canonical-only <source> <canonical>\n' "$0" >&2
  printf '       %s --with-snapshot <source> <canonical> <snapshot-candidate>\n' "$0" >&2
  exit 2
fi

source_path="$2"
canonical_path="$3"
canonical_directory="$(dirname "$canonical_path")"
canonical_filename="$(basename "$canonical_path")"
private_path=""

cleanup_private_path() {
  if [ -n "$private_path" ]; then
    rm -f "$private_path"
  fi
}
trap cleanup_private_path EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

if [ ! -f "$source_path" ]; then
  printf 'Portfolio source is not a regular file: %s\n' "$source_path" >&2
  exit 2
fi
if ! mkdir -p "$canonical_directory"; then
  printf 'Portfolio canonical directory could not be created: %s\n' "$canonical_directory" >&2
  exit 2
fi

private_path="$(mktemp "$canonical_directory/.${canonical_filename}.publishing.XXXXXX")" || {
  printf 'Portfolio private-copy allocation failed beside: %s\n' "$canonical_path" >&2
  exit 2
}
if ! cp "$source_path" "$private_path" || ! cmp -s "$source_path" "$private_path"; then
  printf 'Portfolio source copy could not be verified; no output was published.\n' >&2
  exit 1
fi

snapshot_path=""
if [ "$publication_mode" = "--with-snapshot" ]; then
  snapshot_candidate="$4"
  snapshot_directory="$(dirname "$snapshot_candidate")"
  snapshot_filename="$(basename "$snapshot_candidate")"
  case "$snapshot_filename" in
    *.*)
      snapshot_stem="${snapshot_candidate%.*}"
      snapshot_extension=".${snapshot_candidate##*.}"
      ;;
    *)
      snapshot_stem="$snapshot_candidate"
      snapshot_extension=""
      ;;
  esac

  if ! mkdir -p "$snapshot_directory"; then
    printf 'Portfolio snapshot directory could not be created: %s\n' "$snapshot_directory" >&2
    exit 2
  fi

  suffix_number=1
  while :; do
    if [ "$suffix_number" -eq 1 ]; then
      snapshot_path="$snapshot_candidate"
    else
      snapshot_path="${snapshot_stem}-${suffix_number}${snapshot_extension}"
    fi

    if ln "$private_path" "$snapshot_path"; then
      break
    fi
    if [ -e "$snapshot_path" ] || [ -L "$snapshot_path" ]; then
      suffix_number=$((suffix_number + 1))
      continue
    fi
    printf 'Portfolio snapshot could not be published exclusively; canonical was not changed: %s\n' \
      "$snapshot_path" >&2
    exit 1
  done
fi

if ! mv "$private_path" "$canonical_path"; then
  if [ -n "$snapshot_path" ]; then
    printf 'Portfolio snapshot was published, but canonical refresh failed; snapshot retained: %s\n' \
      "$snapshot_path" >&2
  else
    printf 'Portfolio canonical refresh failed: %s\n' "$canonical_path" >&2
  fi
  exit 1
fi
private_path=""

printf '%s\n' "$canonical_path"
if [ -n "$snapshot_path" ]; then
  printf '%s\n' "$snapshot_path"
fi

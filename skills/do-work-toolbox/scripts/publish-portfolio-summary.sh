#!/usr/bin/env bash
# Publish one verified portfolio draft to canonical and optional immutable snapshot paths.
#
# Snapshot and canonical are published from the same verified bytes but never from the
# same inode: a hard link would make the immutable snapshot follow every later in-place
# edit of canonical, which is the opposite of what a snapshot is for. Each output gets
# its own private copy, verified against the source before either is published.
#
# `ln` and `mv` both treat an existing directory operand as a container rather than a
# collision, so each publication is verified after the fact: a directory in the snapshot
# candidate's place advances to the next numeric suffix, a directory in canonical's place
# fails closed, and neither may leave a private file nested inside it.
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
private_canonical_path=""
private_snapshot_path=""

cleanup_private_copies() {
  if [ -n "$private_canonical_path" ]; then
    rm -f "$private_canonical_path"
  fi
  if [ -n "$private_snapshot_path" ]; then
    rm -f "$private_snapshot_path"
  fi
}
trap cleanup_private_copies EXIT
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

allocate_private_copy() {
  allocated_path="$(mktemp "$canonical_directory/.${canonical_filename}.publishing.XXXXXX")" || return 1
  if ! cp "$source_path" "$allocated_path" || ! cmp -s "$source_path" "$allocated_path"; then
    rm -f "$allocated_path"
    return 1
  fi
  printf '%s' "$allocated_path"
}

private_canonical_path="$(allocate_private_copy)" || {
  printf 'Portfolio source copy could not be verified; no output was published.\n' >&2
  exit 1
}

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

  private_snapshot_path="$(allocate_private_copy)" || {
    printf 'Portfolio source copy could not be verified; no output was published.\n' >&2
    exit 1
  }

  suffix_number=1
  while :; do
    if [ "$suffix_number" -eq 1 ]; then
      snapshot_path="$snapshot_candidate"
    else
      snapshot_path="${snapshot_stem}-${suffix_number}${snapshot_extension}"
    fi

    if ln "$private_snapshot_path" "$snapshot_path"; then
      nested_snapshot_path="$snapshot_path/${private_snapshot_path##*/}"
      if [ -e "$nested_snapshot_path" ]; then
        # The candidate is a directory, so `ln` linked into it instead of colliding.
        rm -f "$nested_snapshot_path"
        suffix_number=$((suffix_number + 1))
        continue
      fi
      if [ ! -f "$snapshot_path" ]; then
        printf 'Portfolio snapshot is not a regular file after publication: %s\n' "$snapshot_path" >&2
        exit 1
      fi
      # Drop the private name so the snapshot is the sole link to its own inode.
      rm -f "$private_snapshot_path"
      private_snapshot_path=""
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

report_retained_snapshot() {
  if [ -n "$snapshot_path" ]; then
    printf 'Portfolio snapshot was published and is retained: %s\n' "$snapshot_path" >&2
  fi
}

if [ -d "$canonical_path" ] && [ ! -L "$canonical_path" ]; then
  printf 'Portfolio canonical path is a directory; refusing to publish into it: %s\n' "$canonical_path" >&2
  report_retained_snapshot
  exit 1
fi

if ! mv "$private_canonical_path" "$canonical_path"; then
  if [ -n "$snapshot_path" ]; then
    printf 'Portfolio snapshot was published, but canonical refresh failed; snapshot retained: %s\n' \
      "$snapshot_path" >&2
  else
    printf 'Portfolio canonical refresh failed: %s\n' "$canonical_path" >&2
  fi
  exit 1
fi

nested_canonical_path="$canonical_path/${private_canonical_path##*/}"
if [ -e "$nested_canonical_path" ]; then
  # A directory appeared in canonical's place after the check above, so `mv` moved the
  # private copy inside it and still reported success. Remove only our own file.
  rm -f "$nested_canonical_path"
  private_canonical_path=""
  printf 'Portfolio canonical path became a directory during publication; it was left unchanged: %s\n' \
    "$canonical_path" >&2
  report_retained_snapshot
  exit 1
fi
private_canonical_path=""

if [ ! -f "$canonical_path" ]; then
  printf 'Portfolio canonical publication is not a regular file: %s\n' "$canonical_path" >&2
  report_retained_snapshot
  exit 1
fi

printf '%s\n' "$canonical_path"
if [ -n "$snapshot_path" ]; then
  printf '%s\n' "$snapshot_path"
fi

#!/usr/bin/env bash
# uncommitted-inventory.sh — mechanical form of the Step 1 preflight shared by
# actions/commit.md and actions/inspect.md: enumerate every uncommitted path,
# categorize it, and tag the secret-shaped ones.
#
# Usage: tools/checks/uncommitted-inventory.sh [repo-root]
# Output: one TAB-separated row per path — "<tag>\t<path>"
#           M  modified (tracked, content changed)
#           A  added (staged-new or untracked)
#           D  deleted
#           X  excluded: secret-shaped name, reported but not to be read
# Exit 0: rows emitted. Exit 1: clean working tree (no rows).
# Exit 2: not a git repository, or usage error.
#
# A renamed path is tagged M, not A: the caller's treatment of M ("read the
# git diff") is the correct one for a move, whereas A tells it to read the
# whole file as if the content were new.
#
# Known limit: paths are emitted verbatim, so a filename containing a literal
# newline produces a row that spans lines. The -z read below means such a path
# is never corrupted or truncated — it just cannot be rendered on one line, and
# a line-oriented format is what the prose callers consume.
#
# X is a tag, never a silent drop. Both callers must still report excluded
# paths — the user needs to know a secret-shaped file is sitting uncommitted —
# they just must not read or diff the contents. A script that omitted the row
# would reintroduce exactly the silence both prose copies warn against.
set -uo pipefail

repository_root="${1:-.}"
if [ ! -d "$repository_root" ]; then
  echo "usage: $0 [repo-root]" >&2
  exit 2
fi
cd "$repository_root" || exit 2

if ! git rev-parse --git-dir >/dev/null 2>&1; then
  echo "NOT-A-GIT-REPO: $repository_root" >&2
  exit 2
fi

# A secret-shaped BASENAME. Matched on the basename rather than the whole path
# so a directory that merely contains "secret" does not tag every ordinary file
# beneath it. The broad *credentials* / *secret* forms subsume the narrower ones
# the prose also spells out.
#
# `.env*` is a prefix glob, deliberately: an earlier `.env|.env.*` spelling let
# `.envrc` through — a direnv file, routinely full of exported secrets — because
# neither branch matches a suffix with no dot. Under-matching here is the one
# failure this script exists to prevent, so the pattern tracks the advertised
# `.env*` exactly. `*.env` is the deliberate extra: `production.env` is an env
# file by any reading, and both callers advertise it alongside the prefix form.
is_secret_shaped() {
  local candidate_basename="${1##*/}"
  case "$candidate_basename" in
    .env*|*.env)               return 0 ;;
    *credentials*)             return 0 ;;
    *.pem|*.key|*.p12|*.pfx)   return 0 ;;
    *secret*)                  return 0 ;;
  esac
  return 1
}

emitted_any_row=0

# --untracked-files=all is load-bearing, not cosmetic. Plain
# `git status --porcelain` collapses a wholly-untracked directory into a single
# "?? dir/" row and never lists the files inside it, so every file in a
# brand-new directory would escape the secret-shaped exclusion above. That is a
# secret-leak path, and actions/stray-check.md's Red Flags record that it has
# been hit.
#
# -z terminates each record with NUL so paths containing spaces, quotes, or
# newlines survive verbatim; without it git quotes such paths and the consumer
# reads the quoting as part of the name. The output is read through process
# substitution rather than "$(...)": bash strips NUL bytes from command
# substitution, which would silently collapse every record into one. Process
# substitution also keeps the loop in this shell, so emitted_any_row survives —
# a pipe would run it in a subshell and the count would always read 0.
while IFS= read -r -d '' status_record; do
  index_status="${status_record:0:1}"
  worktree_status="${status_record:1:1}"
  changed_path="${status_record:3}"

  # A rename or copy record is followed by a SECOND NUL-terminated field holding
  # the origin path. Read and discard it, or the origin gets parsed as the next
  # record's status bytes and every row after it shifts by one.
  case "$index_status$worktree_status" in
    R*|C*|*R|*C) IFS= read -r -d '' _rename_origin_path || true ;;
  esac

  if is_secret_shaped "$changed_path"; then
    printf 'X\t%s\n' "$changed_path"
    emitted_any_row=1
    continue
  fi

  # Deleted wins over modified: a path deleted in either the index or the
  # worktree cannot be read, so a caller that saw "M" would try to diff a file
  # that is gone. Added covers both staged-new (A) and untracked (??).
  case "$index_status$worktree_status" in
    *D*)  path_tag='D' ;;
    '??') path_tag='A' ;;
    A*)   path_tag='A' ;;
    *)    path_tag='M' ;;
  esac

  printf '%s\t%s\n' "$path_tag" "$changed_path"
  emitted_any_row=1
done < <(git status --porcelain --untracked-files=all -z)

[ "$emitted_any_row" -eq 1 ] || exit 1
exit 0

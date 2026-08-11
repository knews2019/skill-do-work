#!/usr/bin/env bash
# Display an ordinary or merge commit without losing the merge diff.
set -u

if [ "$#" -ne 1 ]; then
  printf 'Usage: %s <commit>\n' "$0" >&2
  exit 2
fi

commit_revision="$1"
if ! git rev-parse --verify -q "${commit_revision}^{commit}" >/dev/null; then
  printf 'Commit does not resolve: %s\n' "$commit_revision" >&2
  exit 2
fi

if git rev-parse --verify -q "${commit_revision}^2" >/dev/null; then
  git show --first-parent -m "$commit_revision" --
else
  git show "$commit_revision" --
fi

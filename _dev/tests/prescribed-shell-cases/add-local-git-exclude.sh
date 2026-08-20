#!/usr/bin/env bash
# Fixture execution proofs for add-local-git-exclude.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# add-local-git-exclude: a subdirectory caller must append once to Git's actual exclude.
exclude_repo="$fixture_root/exclude-repo"
fixture_repo_init "$exclude_repo"
mkdir -p "$exclude_repo/nested" "$exclude_repo/cache/data"
(cd "$exclude_repo/nested" && "$core_scripts/add-local-git-exclude.sh" ../cache/data/file.bin '**/cache/data/' >/dev/null) || fail_case 'add-local-git-exclude subdirectory case returned nonzero'
exclude_file="$exclude_repo/$(git -C "$exclude_repo" rev-parse --git-path info/exclude)"
[ "$(grep -Fc '**/cache/data/' "$exclude_file")" -eq 1 ] || fail_case 'add-local-git-exclude subdirectory case did not append exactly once'

prescribed_shell_finish

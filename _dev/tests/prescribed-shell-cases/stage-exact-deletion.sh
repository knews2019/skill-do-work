#!/usr/bin/env bash
# Fixture execution proofs for stage-exact-deletion.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# stage-exact-deletion: stage only the named pathological deletion.
deletion_repo="$fixture_root/deletion-repo"
fixture_repo_init "$deletion_repo"
printf 'delete\n' > "$deletion_repo/secret name.key"
printf 'keep\n' > "$deletion_repo/other.txt"
fixture_repo_commit_all "$deletion_repo" base
rm "$deletion_repo/secret name.key"
printf 'changed\n' >> "$deletion_repo/other.txt"
(cd "$deletion_repo" && "$core_scripts/stage-exact-deletion.sh" 'secret name.key') || fail_case 'stage-exact-deletion pathological-name case returned nonzero'
[ "$(git -C "$deletion_repo" diff --cached --name-only)" = 'secret name.key' ] || fail_case 'stage-exact-deletion pathological-name case staged another path'

# stage-exact-deletion: pathspec-looking filenames must remain literal and isolated.
magic_deletion_repo="$fixture_root/magic-deletion-repo"
fixture_repo_init "$magic_deletion_repo"
magic_deleted_path=':(glob)*'
printf 'magic\n' > "$magic_deletion_repo/$magic_deleted_path"
printf 'other\n' > "$magic_deletion_repo/other.txt"
fixture_repo_commit_all "$magic_deletion_repo" base
rm "$magic_deletion_repo/$magic_deleted_path" "$magic_deletion_repo/other.txt"
(cd "$magic_deletion_repo" && "$core_scripts/stage-exact-deletion.sh" "$magic_deleted_path")
magic_deletion_status=$?
[ "$magic_deletion_status" -eq 0 ] || fail_case "stage-exact-deletion literal-pathspec case returned $magic_deletion_status instead of 0"
magic_cached_deletions="$(git -C "$magic_deletion_repo" diff --cached --name-status --no-renames)"
expected_magic_deletion="$(printf 'D\t%s' "$magic_deleted_path")"
[ "$magic_cached_deletions" = "$expected_magic_deletion" ] || fail_case 'stage-exact-deletion literal-pathspec case staged another path'

prescribed_shell_finish

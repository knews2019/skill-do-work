#!/usr/bin/env bash
# Update a project-local four-module do-work suite.
set -euo pipefail

upstream_url='https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz'

fail() {
  printf 'do-work update: %s\n' "$*" >&2
  exit 1
}

version_order() {
  awk -v local_version="$1" -v remote_version="$2" '
    BEGIN {
      split(local_version, local_parts, ".")
      split(remote_version, remote_parts, ".")
      for (part_number = 1; part_number <= 3; part_number++) {
        local_part = local_parts[part_number] + 0
        remote_part = remote_parts[part_number] + 0
        if (remote_part > local_part) { print 1; exit }
        if (remote_part < local_part) { print -1; exit }
      }
      print 0
    }
  '
}

read_action_version() {
  sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$1" | head -n 1
}

project_root=''
if [ "$#" -eq 2 ] && [ "$1" = '--project-root' ]; then
  project_root="$2"
else
  fail 'usage: do-work-update.sh --project-root <project-root>'
fi

[ -d "$project_root" ] || fail "project root does not exist: $project_root"
project_root="$(cd "$project_root" && pwd -P)"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
skill_root="$(cd "$script_dir/.." && pwd -P)"
installed_manifest_validator="$script_dir/validate-suite-manifest.sh"
installed_suite_installer="$script_dir/install-do-work-suite.sh"

case "$skill_root" in
  "$project_root"/*) ;;
  *) fail "skill is outside this project ($skill_root is not within $project_root); refusing to update a shared install" ;;
esac

[ -s "$skill_root/SKILL.md" ] || fail "SKILL.md is missing at $skill_root"
[ -f "$skill_root/actions/version.md" ] || fail "actions/version.md is missing at $skill_root"
local_version="$(read_action_version "$skill_root/actions/version.md")"
printf '%s\n' "$local_version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' \
  || fail 'could not read a semantic local version'

git_root="$(git -C "$project_root" rev-parse --show-toplevel 2>/dev/null || true)"
[ -n "$git_root" ] || fail 'the project must be a Git repository so a failed suite update can be recovered'
git_root="$(cd "$git_root" && pwd -P)"
[ "$git_root" = "$project_root" ] \
  || fail "--project-root must name the Git worktree root ($git_root)"

update_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-update.XXXXXX")"
upstream_tarball="$update_tmp/upstream.tar.gz"
fresh_upstream="$update_tmp/fresh"
mkdir -p "$fresh_upstream"

remote_version='unknown'
trap 'rm -rf "$update_tmp"' EXIT

printf 'Checking do-work updates…\n'
curl -fsSL -o "$upstream_tarball.download" "$upstream_url" \
  || fail 'upstream tarball download failed; no files were changed'
mv "$upstream_tarball.download" "$upstream_tarball"
tar xzf "$upstream_tarball" -C "$fresh_upstream" --strip-components=1 \
  || fail 'upstream archive could not be extracted; no files were changed'

[ -x "$installed_manifest_validator" ] \
  || fail 'the installed suite manifest validator is missing or not executable'
bash "$installed_manifest_validator" --root "$fresh_upstream" \
  || fail 'suite manifest validation failed; no files were changed'
remote_version="$(sed -n '1p' "$fresh_upstream/VERSION")"
[ -x "$installed_suite_installer" ] \
  || fail 'the installed full-suite installer is missing or not executable'

printf '%s\n' "$remote_version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' \
  || fail 'could not read a semantic upstream version'

case "$(version_order "$local_version" "$remote_version")" in
  0)
    printf "You're up to date (v%s)\n" "$local_version"
    exit 0
    ;;
  -1) fail "upstream version v$remote_version is older than installed v$local_version" ;;
esac

printf 'Update available: v%s (you have v%s), archive layout: four-module suite.\n' \
  "$remote_version" "$local_version"
suite_install_status=0
DO_WORK_INSTALL_CANCEL_EXIT_STATUS=3 \
  bash "$installed_suite_installer" --project-root "$project_root" --archive "$upstream_tarball" \
  || suite_install_status=$?
if [ "$suite_install_status" -eq 3 ]; then
  printf 'Update cancelled; no files were changed.\n'
  exit 0
fi
[ "$suite_install_status" -eq 0 ] \
  || fail 'full-suite installation failed; managed paths were recovered'
installed_version="$(read_action_version "$project_root/.claude/skills/do-work/actions/version.md" 2>/dev/null || true)"
[ "$installed_version" = "$remote_version" ] \
  || fail "post-update version verification failed (expected v$remote_version, found v${installed_version:-unknown})"

printf 'Updated to v%s at %s using the four-module suite.\n' "$remote_version" "$project_root"

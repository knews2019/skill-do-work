#!/usr/bin/env bash
# Updates a project-local do-work install. Invoked by `just run-do-work-update`.
set -euo pipefail

upstream_url='https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz'
shipped_paths=(SKILL.md actions crew-members prompts interviews specs docs hooks tools CHANGELOG.md README.md next-steps.md justfile)

fail() {
  printf 'do-work update: %s\n' "$*" >&2
  exit 1
}

version_is_newer() {
  awk -v local_version="$1" -v remote_version="$2" '
    BEGIN {
      split(local_version, local_parts, ".")
      split(remote_version, remote_parts, ".")
      for (part_number = 1; part_number <= 3; part_number++) {
        local_part = (part_number in local_parts) ? local_parts[part_number] + 0 : 0
        remote_part = (part_number in remote_parts) ? remote_parts[part_number] + 0 : 0
        if (remote_part > local_part) exit 0
        if (remote_part < local_part) exit 1
      }
      exit 1
    }
  '
}

append_install_diff() {
  local fresh_root="$1"
  local installed_root="$2"
  local destination="$3"
  local shipped_path

  for shipped_path in "${shipped_paths[@]}"; do
    if [ -e "$fresh_root/$shipped_path" ] && [ ! -e "$installed_root/$shipped_path" ]; then
      # Upstream ships this but it's wholly absent from the install: `diff` on a
      # missing argument only errors to stderr (swallowed below), so flag it here.
      printf 'Missing from install (shipped upstream): %s\n' "$shipped_path" >> "$destination"
    elif [ -e "$fresh_root/$shipped_path" ] || [ -e "$installed_root/$shipped_path" ]; then
      # No --new-file, so local-only files (compiled binary, Finder .DS_Store,
      # editor swapfile) show as "Only in <install>" extras and are dropped by the
      # fixed-string filter — grep -vF, not -v, so a regex metachar in the install
      # path can't defeat it. Real failures ("Only in <fresh>", content diffs) stay.
      diff -ru "$fresh_root/$shipped_path" "$installed_root/$shipped_path" \
        | grep -vF "Only in $installed_root" >> "$destination" || true
    fi
  done
}

project_root=''
if [ "$#" = 2 ] && [ "$1" = '--project-root' ]; then
  project_root="$2"
else
  fail 'usage: do-work-update.sh --project-root <project-root>'
fi

[ -d "$project_root" ] || fail "project root does not exist: $project_root"
project_root="$(cd "$project_root" && pwd -P)"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
skill_root="$(cd "$script_dir/.." && pwd -P)"

case "$skill_root" in
  "$project_root"|"$project_root"/*) ;;
  *) fail "skill is outside this project ($skill_root is not within $project_root); refusing to update a shared install" ;;
esac

[ -s "$skill_root/SKILL.md" ] || fail "SKILL.md is missing at $skill_root"
[ -f "$skill_root/actions/version.md" ] || fail "actions/version.md is missing at $skill_root"

local_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$skill_root/actions/version.md" | head -n 1)"
[ -n "$local_version" ] || fail 'could not read the local version'

update_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-update.XXXXXX")"
backup_path=''
backup_ready=''
files_in_place=''
cleanup() {
  rm -rf "$update_tmp"
  # If the run aborts inside the destructive region — after the rollback copy is made
  # but before the new files are all in place — `set -e` exits before reaching the
  # explicit rollback prints below. Surface the backup here so a mid-extract failure
  # (e.g. ENOSPC: the cp -R backup just doubled the install's on-disk size) never leaves
  # a half-updated install with no pointer to the rollback copy.
  if [ -n "$backup_ready" ] && [ -z "$files_in_place" ]; then
    printf 'Update did not complete; the install may be partially modified.\nRollback copy: %s\n' "$backup_path" >&2
  fi
}
trap cleanup EXIT
upstream_tarball="$update_tmp/upstream.tar.gz"
fresh_upstream="$update_tmp/fresh"
mkdir -p "$fresh_upstream"

printf 'Checking do-work updates…\n'
curl -fsSL -o "$upstream_tarball.download" "$upstream_url" \
  || fail 'upstream tarball download failed; no files were changed'
mv "$upstream_tarball.download" "$upstream_tarball"
tar xzf "$upstream_tarball" -C "$fresh_upstream" --strip-components=1 \
  --exclude='_dev' --exclude='do-work' --exclude='ai-reports' --exclude='.vscode' --exclude='decisions'

remote_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$fresh_upstream/actions/version.md" | head -n 1)"
[ -n "$remote_version" ] || fail 'could not read the upstream version'

if ! version_is_newer "$local_version" "$remote_version"; then
  printf "You're up to date (v%s)\n" "$local_version"
  exit 0
fi

printf 'Update available: v%s (you have v%s).\n' "$remote_version" "$local_version"

dirty_files=''
if git -C "$skill_root" rev-parse --git-dir >/dev/null 2>&1; then
  dirty_files="$(git -C "$skill_root" status --porcelain -- "${shipped_paths[@]}")"
fi

diff_file="$update_tmp/install.diff"
append_install_diff "$fresh_upstream" "$skill_root" "$diff_file"

if [ -n "$dirty_files" ]; then
  printf 'Shipped skill files have uncommitted changes:\n%s\n' "$dirty_files" >&2
fi
if [ -s "$diff_file" ]; then
  printf 'Reviewing installed-versus-upstream skill changes before overwrite:\n'
  cat "$diff_file"
fi

printf 'Continue with the update? This creates a rollback copy first. [y/N] '
# EOF / non-interactive stdin (piped, CI, </dev/null) makes `read` return non-zero; under
# `set -e` a bare `read` would abort here before the case below can default to No. Treat
# EOF as an explicit No so the cancel path runs and the script exits 0.
read -r confirmation || confirmation=''
case "$confirmation" in
  y|Y|yes|YES) ;;
  *) printf 'Update cancelled; no files were changed.\n'; exit 0 ;;
esac

backup_path="$skill_root.preupdate-$(date -u +%Y%m%dT%H%M%SZ).bak"
cp -R "$skill_root" "$backup_path"
backup_ready=1  # rollback copy exists; the destructive writes below are now recoverable

find "$skill_root/prompts" -maxdepth 1 -name '*.md' ! -name 'README.md' -delete 2>/dev/null || true
find "$skill_root/interviews" -maxdepth 1 -name '*.md' -delete 2>/dev/null || true
tar xzf "$upstream_tarball" -C "$skill_root" --strip-components=1 \
  --exclude='_dev' --exclude='do-work' --exclude='ai-reports' --exclude='.vscode' --exclude='decisions'
# Remove stale vendored maintainer docs (older installs shipped CLAUDE.md/AGENTS.md into
# the skill dir before they were export-ignored from the tarball). ONLY in a nested
# install: when skill_root IS the project root (the justfile recipe's root fallback),
# these are the project's own instruction files, not the skill's — deleting them would
# destroy project instructions, so leave them in place.
if [ "$skill_root" != "$project_root" ]; then
  rm -f "$skill_root/CLAUDE.md" "$skill_root/AGENTS.md"
fi
files_in_place=1  # new files are in place; failures past here have their own rollback prints

installed_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$skill_root/actions/version.md" | head -n 1)"
[ "$installed_version" = "$remote_version" ] \
  || fail "post-update verification failed (expected v$remote_version, found v${installed_version:-unknown}); rollback copy: $backup_path"

post_diff="$update_tmp/post-update.diff"
append_install_diff "$fresh_upstream" "$skill_root" "$post_diff"
if [ -s "$post_diff" ]; then
  printf 'Update completed, but the installed files differ from the reviewed upstream tree:\n' >&2
  cat "$post_diff" >&2
  printf 'Rollback copy: %s\n' "$backup_path" >&2
  exit 1
fi

printf 'Updated to v%s at %s\n' "$remote_version" "$skill_root"
printf 'Rollback copy: %s\n' "$backup_path"

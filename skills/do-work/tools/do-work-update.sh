#!/usr/bin/env bash
# Update a project-local do-work install from either archive layout.
set -euo pipefail

upstream_url='https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz'
bridge_capability='suite-layout-v2'
legacy_shipped_paths=(SKILL.md actions crew-members prompts interviews specs docs hooks tools CHANGELOG.md README.md next-steps.md)

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

if [ "$#" -eq 1 ] && [ "$1" = '--capabilities' ]; then
  printf '%s\n' "$bridge_capability"
  exit 0
fi

project_root=''
if [ "$#" -eq 2 ] && [ "$1" = '--project-root' ]; then
  project_root="$2"
else
  fail 'usage: do-work-update.sh --project-root <project-root> | --capabilities'
fi

[ -d "$project_root" ] || fail "project root does not exist: $project_root"
project_root="$(cd "$project_root" && pwd -P)"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
skill_root="$(cd "$script_dir/.." && pwd -P)"
installed_manifest_validator="$script_dir/validate-suite-manifest.sh"
installed_suite_installer="$script_dir/install-do-work-suite.sh"

case "$skill_root" in
  "$project_root"|"$project_root"/*) ;;
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

plan_sources=()
plan_destinations=()
managed_relative_paths=()
preexisting_paths=()
deletion_source="$update_tmp/deleted-upstream-path"
write_started=''
update_verified=''
remote_version='unknown'

path_preexisted() {
  local candidate_path="$1" recorded_path
  for recorded_path in "${preexisting_paths[@]}"; do
    [ "$recorded_path" = "$candidate_path" ] && return 0
  done
  return 1
}

recover_managed_paths() {
  local relative_path destination_path candidate_path tracked_files recovery_failed=''

  set +e
  for relative_path in "${managed_relative_paths[@]}"; do
    tracked_files="$(git -C "$project_root" ls-files -- "$relative_path")"
    if [ -n "$tracked_files" ]; then
      git -C "$project_root" restore --source=HEAD --staged --worktree -- "$relative_path" \
        || recovery_failed=1
    fi
  done

  for destination_path in "${plan_destinations[@]}"; do
    if [ -e "$destination_path" ] || [ -L "$destination_path" ]; then
      while IFS= read -r -d '' candidate_path; do
        if ! path_preexisted "$candidate_path"; then
          rm -rf -- "$candidate_path" || recovery_failed=1
        fi
      done < <(find "$destination_path" -depth -print0 2>/dev/null)
    fi
  done
  set -e

  if [ -n "$recovery_failed" ]; then
    printf 'Automatic recovery was incomplete; inspect only these managed paths before retrying:\n' >&2
    printf '  %s\n' "${managed_relative_paths[@]}" >&2
    return 1
  fi
  printf 'Restored the previous managed installation after the failed update.\n' >&2
}

cleanup() {
  local exit_status=$?
  trap - EXIT
  if [ -n "$write_started" ] && [ -z "$update_verified" ]; then
    printf 'Update did not complete; recovering the managed skill paths.\n' >&2
    recover_managed_paths || exit_status=1
  fi
  rm -rf "$update_tmp"
  exit "$exit_status"
}
trap cleanup EXIT

printf 'Checking do-work updates…\n'
curl -fsSL -o "$upstream_tarball.download" "$upstream_url" \
  || fail 'upstream tarball download failed; no files were changed'
mv "$upstream_tarball.download" "$upstream_tarball"
tar xzf "$upstream_tarball" -C "$fresh_upstream" --strip-components=1 \
  || fail 'upstream archive could not be extracted; no files were changed'

archive_layout='legacy all-in-one skill'
suite_installer=''
if [ -e "$fresh_upstream/VERSION" ] \
  || [ -e "$fresh_upstream/suite" ] \
  || [ -e "$fresh_upstream/skills" ]; then
  archive_layout='four-module suite'
  [ -x "$installed_manifest_validator" ] \
    || fail 'the installed bridge manifest validator is missing or not executable'
  bash "$installed_manifest_validator" --root "$fresh_upstream" \
    || fail 'suite manifest validation failed; no files were changed'
  remote_version="$(sed -n '1p' "$fresh_upstream/VERSION")"
  [ -x "$installed_suite_installer" ] \
    || fail 'the installed bridge full-suite installer is missing or not executable'
  suite_installer="$installed_suite_installer"
else
  [ -s "$fresh_upstream/SKILL.md" ] \
    || fail 'legacy archive is missing SKILL.md; no files were changed'
  [ -f "$fresh_upstream/actions/version.md" ] \
    || fail 'legacy archive is missing actions/version.md; no files were changed'
  remote_version="$(read_action_version "$fresh_upstream/actions/version.md")"

  if [ "$skill_root" = "$project_root" ]; then
    skill_relative_root=''
  else
    skill_relative_root="${skill_root#"$project_root"/}"
  fi
  for legacy_path in "${legacy_shipped_paths[@]}"; do
    [ -e "$fresh_upstream/$legacy_path" ] || continue
    [ ! -L "$fresh_upstream/$legacy_path" ] \
      || fail "legacy archive path is a symlink: $legacy_path"
    plan_sources+=("$fresh_upstream/$legacy_path")
    plan_destinations+=("$skill_root/$legacy_path")
    if [ -n "$skill_relative_root" ]; then
      managed_relative_paths+=("$skill_relative_root/$legacy_path")
    else
      managed_relative_paths+=("$legacy_path")
    fi
  done
  if [ "$skill_root" != "$project_root" ] && [ -e "$fresh_upstream/justfile" ]; then
    plan_sources+=("$fresh_upstream/justfile")
    plan_destinations+=("$skill_root/justfile")
    managed_relative_paths+=("$skill_relative_root/justfile")
  fi
  # Preserve the legacy updater's one intentional stale-file cleanup. These maintainer docs
  # shipped in old nested installs but are not skill content; root fallback paths are the
  # consuming project's own instructions and are therefore never managed here.
  if [ "$skill_root" != "$project_root" ]; then
    for stale_maintainer_doc in CLAUDE.md AGENTS.md; do
      if [ -e "$skill_root/$stale_maintainer_doc" ] || [ -L "$skill_root/$stale_maintainer_doc" ]; then
        plan_sources+=("$deletion_source")
        plan_destinations+=("$skill_root/$stale_maintainer_doc")
        managed_relative_paths+=("$skill_relative_root/$stale_maintainer_doc")
      fi
    done
  fi
fi

printf '%s\n' "$remote_version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' \
  || fail 'could not read a semantic upstream version'

case "$(version_order "$local_version" "$remote_version")" in
  0)
    printf "You're up to date (v%s)\n" "$local_version"
    exit 0
    ;;
  -1) fail "upstream version v$remote_version is older than installed v$local_version" ;;
esac

if [ -n "$suite_installer" ]; then
  printf 'Update available: v%s (you have v%s), archive layout: %s.\n' \
    "$remote_version" "$local_version" "$archive_layout"
  suite_install_status=0
  DO_WORK_INSTALL_CANCEL_EXIT_STATUS=3 \
    bash "$suite_installer" --project-root "$project_root" --archive "$upstream_tarball" \
    || suite_install_status=$?
  if [ "$suite_install_status" -eq 3 ]; then
    update_verified=1
    printf 'Update cancelled; no files were changed.\n'
    exit 0
  fi
  [ "$suite_install_status" -eq 0 ] \
    || fail 'full-suite installation failed; managed paths were recovered'
  installed_version="$(read_action_version "$project_root/.claude/skills/do-work/actions/version.md" 2>/dev/null || true)"
  [ "$installed_version" = "$remote_version" ] \
    || fail "post-update version verification failed (expected v$remote_version, found v${installed_version:-unknown})"
  update_verified=1
  printf 'Updated to v%s at %s using the %s.\n' "$remote_version" "$project_root" "$archive_layout"
  exit 0
fi

[ "${#plan_sources[@]}" -gt 0 ] || fail 'archive declares no managed update paths'

# Suite validation covers its four module roots. Legacy archives have no manifest, so reject
# symlinks anywhere in their managed source paths before constructing a destination write.
for source_path in "${plan_sources[@]}"; do
  [ "$source_path" = "$deletion_source" ] && continue
  if [ -L "$source_path" ]; then
    fail "managed source path is a symlink: $source_path"
  fi
  if [ -d "$source_path" ] && find "$source_path" -type l -print -quit | grep -q .; then
    fail "managed source path contains a symlink: $source_path"
  fi
done

# A textual manifest destination can still escape through an existing parent symlink. Resolve
# each nearest existing parent physically before showing the diff or accepting confirmation.
for destination_path in "${plan_destinations[@]}"; do
  destination_parent="$(dirname "$destination_path")"
  while [ ! -d "$destination_parent" ]; do
    next_parent="$(dirname "$destination_parent")"
    [ "$next_parent" != "$destination_parent" ] \
      || fail "could not resolve managed destination parent: $destination_path"
    destination_parent="$next_parent"
  done
  destination_parent="$(cd "$destination_parent" && pwd -P)"
  case "$destination_parent" in
    "$project_root"|"$project_root"/*) ;;
    *) fail "managed destination resolves outside the project: $destination_path" ;;
  esac
done

printf 'Update available: v%s (you have v%s), archive layout: %s.\n' \
  "$remote_version" "$local_version" "$archive_layout"

dirty_files="$(git -C "$project_root" status --porcelain -- "${managed_relative_paths[@]}")"
if [ -n "$dirty_files" ]; then
  printf 'Managed skill paths have uncommitted changes:\n%s\n' "$dirty_files" >&2
  printf 'Continuing discards those changes and restores from committed Git content if recovery is needed.\n' >&2
fi

diff_file="$update_tmp/install.diff"
: > "$diff_file"
for plan_index in "${!plan_sources[@]}"; do
  source_path="${plan_sources[$plan_index]}"
  destination_path="${plan_destinations[$plan_index]}"
  printf '\n--- managed destination: %s ---\n' \
    "${managed_relative_paths[$plan_index]}" >> "$diff_file"
  if [ "$source_path" = "$deletion_source" ]; then
    printf 'Upstream deletion: %s\n' "$destination_path" >> "$diff_file"
  else
    diff_status=0
    diff -ruN "$destination_path" "$source_path" >> "$diff_file" 2>&1 || diff_status=$?
    [ "$diff_status" -le 1 ] \
      || fail "could not compare managed destination ${managed_relative_paths[$plan_index]}"
  fi
done

printf 'Reviewing the complete managed update before overwrite:\n'
cat "$diff_file"
printf 'Continue with this one %s update? [y/N] ' "$archive_layout"
read -r confirmation || confirmation=''
case "$confirmation" in
  y|Y|yes|YES) ;;
  *) printf 'Update cancelled; no files were changed.\n'; exit 0 ;;
esac

# Record every path that exists inside the validated destinations. Recovery later removes only
# paths absent from this inventory, then restores tracked content from HEAD. It never invokes a
# broad checkout/clean and therefore cannot cross into project-owned runtime or configuration.
for destination_path in "${plan_destinations[@]}"; do
  if [ -e "$destination_path" ] || [ -L "$destination_path" ]; then
    while IFS= read -r -d '' existing_path; do
      preexisting_paths+=("$existing_path")
    done < <(find "$destination_path" -print0)
  fi
done

write_started=1
# Confirmation explicitly authorizes discarding dirty managed content. Normalize tracked
# managed paths in both index and worktree first so a previously staged customization cannot
# survive underneath the new files or reappear during a later Git restore.
for relative_path in "${managed_relative_paths[@]}"; do
  if [ -n "$(git -C "$project_root" ls-files -- "$relative_path")" ]; then
    git -C "$project_root" restore --source=HEAD --staged --worktree -- "$relative_path"
  fi
done

for plan_index in "${!plan_sources[@]}"; do
  source_path="${plan_sources[$plan_index]}"
  destination_path="${plan_destinations[$plan_index]}"
  rm -rf -- "$destination_path"
  [ "$source_path" = "$deletion_source" ] && continue
  mkdir -p "$(dirname "$destination_path")"
  if [ -d "$source_path" ]; then
    mkdir -p "$destination_path"
    cp -R "$source_path/." "$destination_path/"
  else
    cp -p "$source_path" "$destination_path"
  fi
done

for plan_index in "${!plan_sources[@]}"; do
  source_path="${plan_sources[$plan_index]}"
  destination_path="${plan_destinations[$plan_index]}"
  if [ "$source_path" = "$deletion_source" ]; then
    if [ -e "$destination_path" ] || [ -L "$destination_path" ]; then
      fail "post-update deletion verification failed for ${managed_relative_paths[$plan_index]}"
    fi
    continue
  fi
  if [ -d "$source_path" ]; then
    diff -qr "$source_path" "$destination_path" >/dev/null \
      || fail "post-update byte verification failed for ${managed_relative_paths[$plan_index]}"
  else
    cmp -s "$source_path" "$destination_path" \
      || fail "post-update byte verification failed for ${managed_relative_paths[$plan_index]}"
  fi
done

installed_version="$(read_action_version "$project_root/.claude/skills/do-work/actions/version.md" 2>/dev/null || true)"
if [ "$skill_root" = "$project_root" ] && [ "$archive_layout" = 'legacy all-in-one skill' ]; then
  installed_version="$(read_action_version "$project_root/actions/version.md")"
fi
[ "$installed_version" = "$remote_version" ] \
  || fail "post-update version verification failed (expected v$remote_version, found v${installed_version:-unknown})"

update_verified=1
printf 'Updated to v%s at %s using the %s.\n' "$remote_version" "$project_root" "$archive_layout"

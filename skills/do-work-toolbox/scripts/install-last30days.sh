#!/usr/bin/env bash
# Install or repair the project-local last30days skill and verify every guarantee.
set -u

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
  printf 'Usage: %s (check|install) <project-root> [source-repository]\n' "$0" >&2
  exit 2
fi

install_mode="$1"
project_root="$2"
source_repository="${3:-https://github.com/mvanhorn/last30days-skill}"
target_directory="$project_root/.claude/skills/last30days"
clone_directory=""
staging_directory=""
backup_parent=""
publication_started=0
publication_complete=0

last30days_tree_is_complete() {
  candidate_directory="$1"
  [ -s "$candidate_directory/SKILL.md" ] \
    && [ -s "$candidate_directory/scripts/last30days.py" ]
}

restore_previous_tree() {
  if [ -e "$target_directory" ] || [ -L "$target_directory" ]; then
    rm -rf "$target_directory"
  fi
  if [ -n "$backup_parent" ] && [ -e "$backup_parent/previous" ]; then
    mv "$backup_parent/previous" "$target_directory" || return 1
  fi
  publication_started=0
}

cleanup_install_paths() {
  if [ "$publication_complete" -eq 0 ] \
    && { [ "$publication_started" -eq 1 ] \
      || { [ -n "$backup_parent" ] && [ -e "$backup_parent/previous" ]; }; }; then
    if ! restore_previous_tree; then
      printf 'last30days: rollback FAILED; prior tree remains at %s\n' "$backup_parent/previous" >&2
      backup_parent=""
    fi
  fi
  if [ -n "$clone_directory" ]; then
    rm -rf "$clone_directory"
  fi
  if [ -n "$staging_directory" ]; then
    rm -rf "$staging_directory"
  fi
  if [ -n "$backup_parent" ]; then
    rm -rf "$backup_parent"
  fi
}
trap cleanup_install_paths EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

if [ "$install_mode" != "check" ] && [ "$install_mode" != "install" ]; then
  printf 'Unknown last30days mode: %s\n' "$install_mode" >&2
  exit 2
fi

if [ "$install_mode" = "install" ] && ! last30days_tree_is_complete "$target_directory"; then
  target_parent="$project_root/.claude/skills"
  mkdir -p "$target_parent" || exit 2
  clone_directory="$(mktemp -d "${TMPDIR:-/tmp}/do-work-last30days.XXXXXX")" || exit 2
  source_directory="$clone_directory/skills/last30days"
  if ! git clone --depth 1 "$source_repository" "$clone_directory" \
    || ! last30days_tree_is_complete "$source_directory"; then
    printf 'last30days: clone/source validation FAILED\n' >&2
    exit 1
  fi

  staging_directory="$(mktemp -d "$target_parent/.last30days.staging.XXXXXX")" || exit 2
  if ! cp -R "$source_directory/." "$staging_directory/" \
    || ! last30days_tree_is_complete "$staging_directory"; then
    printf 'last30days: clone/copy FAILED\n' >&2
    exit 1
  fi
  rm -rf "$clone_directory"
  clone_directory=""

  if [ -e "$target_directory" ] || [ -L "$target_directory" ]; then
    backup_parent="$(mktemp -d "$target_parent/.last30days.backup.XXXXXX")" || exit 2
    if ! mv "$target_directory" "$backup_parent/previous"; then
      printf 'last30days: existing-tree backup FAILED\n' >&2
      exit 1
    fi
  fi

  publication_started=1
  if mv "$staging_directory" "$target_directory"; then
    publication_complete=1
    staging_directory=""
  else
    if ! restore_previous_tree; then
      printf 'last30days: publication FAILED; prior tree remains at %s\n' "$backup_parent/previous" >&2
      backup_parent=""
      exit 1
    fi
    printf 'last30days: publication FAILED\n' >&2
    exit 1
  fi
fi

if [ "$install_mode" = "install" ] && git -C "$project_root" rev-parse --git-dir >/dev/null 2>&1; then
  if git -C "$project_root" ls-files --error-unmatch -- .claude/skills/last30days >/dev/null 2>&1; then
    printf 'last30days: tracked vendored path must be untracked before local ignore can protect it.\n' >&2
    exit 1
  fi
  script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  local_exclude_script="$script_directory/../../do-work/scripts/add-local-git-exclude.sh"
  (
    cd "$project_root" || exit 2
    "$local_exclude_script" .claude/skills/last30days/SKILL.md '**/.claude/skills/last30days/' >/dev/null
  ) || exit $?
fi

verification_failed=0
if last30days_tree_is_complete "$target_directory"; then
  printf 'runtime payload: OK\n'
else
  printf 'runtime payload: FAILED\n'
  verification_failed=1
fi
if git -C "$project_root" rev-parse --git-dir >/dev/null 2>&1; then
  if git -C "$project_root" check-ignore -q .claude/skills/last30days/SKILL.md; then
    printf 'ignore rule: OK\n'
  else
    printf 'ignore rule: FAILED\n'
    verification_failed=1
  fi
else
  printf 'ignore rule: n/a (not a git repo)\n'
fi

found_python=""
for python_candidate in python3.13 python3.12 python3 python; do
  if command -v "$python_candidate" >/dev/null 2>&1 \
    && "$python_candidate" -c 'import sys; raise SystemExit(0 if sys.version_info >= (3, 12) else 1)' 2>/dev/null; then
    found_python="$python_candidate"
    break
  fi
done
if [ -n "$found_python" ]; then
  printf 'python 3.12+: OK (%s)\n' "$found_python"
else
  printf 'python 3.12+: FAILED\n'
  verification_failed=1
fi
exit "$verification_failed"

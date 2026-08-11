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
skill_file="$target_directory/SKILL.md"
clone_directory=""

cleanup_clone() {
  if [ -n "$clone_directory" ]; then
    rm -rf "$clone_directory"
  fi
}
trap cleanup_clone EXIT HUP INT TERM

if [ "$install_mode" = "install" ] && [ ! -s "$skill_file" ]; then
  clone_directory="$(mktemp -d "${TMPDIR:-/tmp}/do-work-last30days.XXXXXX")" || exit 2
  if ! git clone --depth 1 "$source_repository" "$clone_directory" \
    || [ ! -s "$clone_directory/skills/last30days/SKILL.md" ] \
    || ! mkdir -p "$target_directory" \
    || ! cp -R "$clone_directory/skills/last30days/." "$target_directory/"; then
    printf 'last30days: clone/copy FAILED\n' >&2
    exit 1
  fi
  rm -rf "$clone_directory"
  clone_directory=""
elif [ "$install_mode" != "check" ] && [ "$install_mode" != "install" ]; then
  printf 'Unknown last30days mode: %s\n' "$install_mode" >&2
  exit 2
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
if [ -s "$skill_file" ]; then printf 'skill file: OK\n'; else printf 'skill file: FAILED\n'; verification_failed=1; fi
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

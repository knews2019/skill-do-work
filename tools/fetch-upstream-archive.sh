#!/usr/bin/env bash
# Fetch the upstream suite archive by whichever route works, and publish it only whole.
#
# Route 1 is the anonymous tarball over HTTP, fetched through the shipped
# atomic-download primitive so it inherits that helper's retry and optional
# credentials. Route 2 is a shallow clone repacked with `git archive`.
#
# The git transport sits behind a different rate limiter than codeload, which is
# what makes route 2 load-bearing rather than decorative: a sustained 429 defeats
# retry alone. `git archive` — never a worktree copy — is mandatory, because only
# it honors `export-ignore`; `cp -R`, `rsync`, and tarring a clone all ship the
# maintainer-only tree into consumer installs, and nothing downstream catches that.
set -u

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
  printf 'Usage: %s <archive-target-path> <upstream-tarball-url> [upstream-repo-url]\n' "$0" >&2
  exit 2
fi

archive_target_path="$1"
upstream_tarball_url="$2"
upstream_repo_url="${3:-}"
archive_stage_path=""
clone_directory=""

fetch_cleanup() {
  if [ -n "$archive_stage_path" ]; then
    rm -f "$archive_stage_path"
  fi
  if [ -n "$clone_directory" ]; then
    rm -rf "$clone_directory"
  fi
}
trap fetch_cleanup EXIT HUP INT TERM

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
# This file is mirrored into two layouts with different relative depths, so the
# primitive is located by probing both rather than by a single fixed path.
atomic_download_script=""
for atomic_download_candidate in \
  "$script_dir/../scripts/atomic-download.sh" \
  "$script_dir/../skills/do-work/scripts/atomic-download.sh"; do
  if [ -f "$atomic_download_candidate" ]; then
    atomic_download_script="$atomic_download_candidate"
    break
  fi
done

# Derive the git route from the tarball URL when the caller did not name it.
# GitHub's branch tarballs unpack to "<repo>-<branch>/", which both callers'
# `tar --strip-components=1` already assumes, so the archive prefix matches it.
archive_prefix='upstream-main/'
case "$upstream_tarball_url" in
  */archive/refs/heads/*.tar.gz)
    upstream_branch="${upstream_tarball_url##*/archive/refs/heads/}"
    upstream_branch="${upstream_branch%.tar.gz}"
    upstream_repository_base="${upstream_tarball_url%%/archive/refs/heads/*}"
    archive_prefix="${upstream_repository_base##*/}-${upstream_branch}/"
    if [ -z "$upstream_repo_url" ]; then
      upstream_repo_url="${upstream_repository_base}.git"
    fi
    ;;
esac

http_route_outcome='not attempted'
git_route_outcome='not attempted'

if [ -n "$atomic_download_script" ]; then
  if bash "$atomic_download_script" "$upstream_tarball_url" "$archive_target_path" 2>/dev/null \
    && tar tzf "$archive_target_path" >/dev/null 2>&1; then
    printf 'upstream archive fetched over HTTP\n'
    exit 0
  fi
  http_route_outcome='failed (host unreachable, rate limited, or archive unreadable)'
else
  http_route_outcome='unavailable (atomic-download.sh not found beside this script)'
fi

if [ -z "$upstream_repo_url" ]; then
  git_route_outcome='unavailable (no repository URL supplied and none derivable from the tarball URL)'
elif ! command -v git >/dev/null 2>&1; then
  git_route_outcome='unavailable (git is not installed)'
else
  clone_directory="$(mktemp -d "${TMPDIR:-/tmp}/do-work-upstream-clone.XXXXXX")" || clone_directory=""
  archive_stage_path="$(mktemp "${archive_target_path}.fetching.XXXXXX")" || archive_stage_path=""
  if [ -z "$clone_directory" ] || [ -z "$archive_stage_path" ]; then
    git_route_outcome='failed (could not allocate private working paths)'
  else
    rm -rf "$clone_directory"
    if GIT_TERMINAL_PROMPT=0 git clone --depth 1 --quiet "$upstream_repo_url" "$clone_directory" 2>/dev/null \
      && git -C "$clone_directory" archive --format=tar.gz --prefix="$archive_prefix" HEAD > "$archive_stage_path" 2>/dev/null \
      && [ -s "$archive_stage_path" ] \
      && tar tzf "$archive_stage_path" >/dev/null 2>&1 \
      && mv "$archive_stage_path" "$archive_target_path"; then
      archive_stage_path=""
      printf 'upstream archive fetched with git (HTTP route %s)\n' "$http_route_outcome"
      exit 0
    fi
    git_route_outcome='failed (clone, repack, or publication did not complete)'
  fi
fi

printf 'upstream archive could not be fetched. HTTP route: %s. Git route: %s.\n' \
  "$http_route_outcome" "$git_route_outcome" >&2
printf 'Set DO_WORK_UPSTREAM_URL to a reachable archive URL to route around a blocked host.\n' >&2
exit 1

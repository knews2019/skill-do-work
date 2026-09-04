#!/usr/bin/env bash
# Ratchet the single-home contract for prescribed shell primitive rationale.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
canonical_guide="$repo_root/skills/do-work/docs/prescribed-shell-primitives.md"
failure_count=0

for prescribed_script in \
  skills/do-work/scripts/show-commit-diff.sh \
  skills/do-work/scripts/add-local-git-exclude.sh \
  skills/do-work/scripts/atomic-download.sh \
  skills/do-work/scripts/capture-screenshot.sh \
  skills/do-work/scripts/run-blocked-check.sh \
  skills/do-work/scripts/protected-inventory.sh \
  skills/do-work/scripts/stage-exact-deletion.sh \
  skills/do-work-knowledge/scripts/lexical-memory-recall.sh \
  skills/do-work-knowledge/scripts/install-memory-hooks.sh \
  skills/do-work/tools/fetch-upstream-archive.sh
do
  if [ ! -x "$repo_root/$prescribed_script" ]; then
    printf 'FAIL: prescribed shell script is missing or not executable: %s\n' "$prescribed_script" >&2
    failure_count=$((failure_count + 1))
  fi
done

# A tool script must delegate its downloads to the shared primitives rather than
# hand-rolling one — the canonicalization campaign swept actions/ and scripts/ and
# never named tools/, which is how the same curl ended up written four times.
#
# One exemption, and it is keyed on a condition rather than a filename: text inside a
# *quoted heredoc* is emitted for someone else to run, not executed by this script.
# That is what makes the installer's BOOTSTRAP block legitimate — nothing is installed
# when it runs, so it cannot call a helper that does not exist yet.
#
# fetch-upstream-archive.sh is deliberately NOT exempt. It is the primitive's tool-side
# home, but it delegates to scripts/atomic-download.sh rather than calling curl itself,
# so it needs no exemption today — and pre-granting one would authorize in advance the
# exact duplication this check exists to prevent.
while IFS= read -r tool_script_path; do
  direct_download_lines="$(awk '
    match($0, /<<-?\047[A-Za-z_][A-Za-z0-9_]*\047/) {
      heredoc_delimiter = substr($0, RSTART, RLENGTH)
      gsub(/^<<-?\047|\047$/, "", heredoc_delimiter)
      inside_quoted_heredoc = 1
      next
    }
    inside_quoted_heredoc && $0 == heredoc_delimiter { inside_quoted_heredoc = 0; next }
    inside_quoted_heredoc { next }
    /^[[:space:]]*#/ { next }
    /(^|[;&|(]|\$\()[[:space:]]*curl[[:space:]]/ { printf "%d:%s\n", NR, $0 }
  ' "$tool_script_path")"
  if [ -n "$direct_download_lines" ]; then
    printf 'FAIL: tool script downloads directly instead of delegating to the shared fetcher: %s\n' \
      "${tool_script_path#"$repo_root/"}" >&2
    printf '%s\n' "$direct_download_lines" >&2
    failure_count=$((failure_count + 1))
  fi
done < <(find "$repo_root"/skills/*/tools -type f -name '*.sh' | sort)

if [[ ! -f "$canonical_guide" ]]; then
  printf 'FAIL: core prescribed-shell guide is missing.\n' >&2
  exit 1
fi

for required_heading in \
  '## Per-file untracked inventory' \
  '## Merge-aware commit diff' \
  '## Commit file listing' \
  '## Local Git ignore' \
  '## Verified exact publication' \
  '## Atomic download publication' \
  '## Portfolio summary publication' \
  '## Report image batch publication' \
  '## Raw text before shell quoting' \
  '## Diff output filtering' \
  '## State across command blocks'
do
  if ! grep -Fqx "$required_heading" "$canonical_guide"; then
    printf 'FAIL: prescribed-shell guide is missing heading: %s\n' "$required_heading" >&2
    failure_count=$((failure_count + 1))
  fi
done

core_pointer='../docs/prescribed-shell-primitives.md'
sibling_pointer='../../do-work/docs/prescribed-shell-primitives.md'
for core_site in \
  skills/do-work/actions/commit.md \
  skills/do-work/actions/capture.md \
  skills/do-work/actions/review-work.md \
  skills/do-work/actions/work.md \
  skills/do-work/actions/work-reference.md \
  skills/do-work/crew-members/background-agents.md
do
  if ! grep -Fq "$core_pointer" "$repo_root/$core_site"; then
    printf 'FAIL: %s does not point at the core prescribed-shell guide.\n' "$core_site" >&2
    failure_count=$((failure_count + 1))
  fi
done

for sibling_site in \
  skills/do-work-board/actions/board.md \
  skills/do-work-knowledge/actions/memory-reference.md \
  skills/do-work-knowledge/actions/setup-memory.md \
  skills/do-work-knowledge/crew-members/background-agents.md \
  skills/do-work-toolbox/actions/ai-report.md \
  skills/do-work-toolbox/actions/inspect.md \
  skills/do-work-toolbox/actions/install.md \
  skills/do-work-toolbox/actions/present-work.md \
  skills/do-work-toolbox/actions/stray-check.md \
  skills/do-work-toolbox/crew-members/background-agents.md
do
  if ! grep -Fq "$sibling_pointer" "$repo_root/$sibling_site"; then
    printf 'FAIL: %s does not point at the core prescribed-shell guide.\n' "$sibling_site" >&2
    failure_count=$((failure_count + 1))
  fi
done

stale_patterns_file="$(mktemp)" || exit 1
trap 'rm -f "$stale_patterns_file"' EXIT
printf '%s\n' \
  'plain `git status --porcelain` collapses' \
  'plain `git show` prints a combined diff that is usually empty' \
  'pattern with an interior slash is root-anchored' \
  'curl -o` writes the final path incrementally' \
  'container rather than a collision' \
  'follow every later in-place edit' \
  'Shell state does not survive' \
  'never interpolate raw user text inside shell quoting' \
  '`diff -x PATTERN`' \
  > "$stale_patterns_file"

old_implementations_file="$(mktemp)" || exit 1
trap 'rm -f "$stale_patterns_file" "$old_implementations_file"' EXIT
printf '%s\n' \
  'screenshot_copy_path=' \
  'cached_deletion_file=' \
  'gen_image()' \
  'CLONE_DIR=' \
  'append_session_start=1' \
  'command -v gtimeout' \
  'exclude_file="$(git rev-parse --git-path info/exclude' \
  > "$old_implementations_file"

while IFS= read -r shipped_markdown; do
  [[ "$shipped_markdown" == "$canonical_guide" ]] && continue
  [[ "$(basename "$shipped_markdown")" == CHANGELOG.md ]] && continue
  while IFS= read -r stale_pattern; do
    if grep -Fq "$stale_pattern" "$shipped_markdown"; then
      printf 'FAIL: %s restates canonical prescribed-shell rationale <%s>; keep local intent and point at the guide.\n' \
        "${shipped_markdown#"$repo_root/"}" "$stale_pattern" >&2
      failure_count=$((failure_count + 1))
    fi
  done < "$stale_patterns_file"
  while IFS= read -r old_implementation; do
    if grep -Fq "$old_implementation" "$shipped_markdown"; then
      printf 'FAIL: %s retains promoted shell implementation <%s>; keep intent plus the shipped script invocation.\n' \
        "${shipped_markdown#"$repo_root/"}" "$old_implementation" >&2
      failure_count=$((failure_count + 1))
    fi
  done < "$old_implementations_file"
done < <(find "$repo_root/skills" -type f -name '*.md' -print | LC_ALL=C sort)

if [[ "$failure_count" -gt 0 ]]; then
  exit 1
fi

printf 'Prescribed shell primitive canonicalization checks passed.\n'

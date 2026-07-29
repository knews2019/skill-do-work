#!/usr/bin/env bash
# qualify.sh — mechanical parts of actions/work.md Step 6.3 (Qualify
# Implementation). Covers checklist items 1 (files exist / show in diff),
# 4 (P-A-U box audit vs debug artifacts), and the grep half of 5 (wiring).
# Items 2 (substantive), 3 (requirements traced), and 6 (data actually flows)
# are judgment — this script feeds them evidence, it does not decide them.
#
# Usage: tools/checks/qualify.sh <req-file>
# Env: DO_WORK_DIFF_RANGE — when set (worktree dispatch mode, e.g. "<pre>..<merge_hash>"),
#   the file-change and debug-artifact checks read `git diff <range>` instead of
#   the working+staged diff. Unset/empty (the serial default) reads working+staged,
#   byte-for-byte as before. Setting it only changes WHICH diff is read. A set-but-
#   unresolvable range is a hard FAIL (exit 1) naming the range — never a silent OK.
# Exit 0: all mechanical checks pass. Exit 1: at least one FAIL line.
# Exit 2: usage error (missing/unreadable <req-file>).
# WARN lines (wiring not found, etc.) do not fail the run — they are handed to
# the orchestrator's judgment, which owns the exception list (entry points,
# framework-convention routes, barrel re-exports, dynamic imports, ...).
set -uo pipefail

request_file="${1:-}"
if [ ! -f "$request_file" ]; then
  echo "usage: $0 <req-file>" >&2
  exit 2
fi

# Worktree dispatch mode passes this REQ's merge range (<pre>..<merge_hash>). Empty/unset
# is the serial default: read the working+staged diff, unchanged.
diff_range="${DO_WORK_DIFF_RANGE:-}"

failure_count=0

summary_bullets="$(awk '
  /^## Implementation Summary$/ {inside=1; next}
  inside && /^## / {inside=0}
  inside && /^- `/ {print}
' "$request_file")"

if [ -z "$summary_bullets" ]; then
  echo "FAIL: no '## Implementation Summary' file list — run after Step 6.25"
  exit 1
fi

git_available=0
git rev-parse --git-dir >/dev/null 2>&1 && git_available=1
[ "$git_available" -eq 0 ] && echo "WARN: not a git repository — diff-based checks degraded to existence checks"

# --- Range validation: an unresolvable range FAILs, it never reads as clean ---
# `git diff <unresolvable-range>` exits 128 with EMPTY stdout, and this script runs
# without `set -e`. So before this gate, a bogus range emptied both $changed_file_list
# and the debug-artifact grep and every check passed vacuously — the script printed
# `OK: mechanical qualification passed` with git's fatal scrolled off above it, and the
# debug-artifact check was silently disabled. The classic bogus value is an endpoint the
# orchestrator lost between command blocks (shell state does not survive them), which
# yields ".." or "..<hash>"; hence the per-endpoint diagnostic below. Validate HERE,
# before any check consumes the range. Serial mode never enters this block.
if [ -n "$diff_range" ]; then
  if [ "$git_available" -eq 0 ]; then
    echo "FAIL: DO_WORK_DIFF_RANGE is set ('$diff_range') but this is not a git repository — the merge range cannot be read"
    exit 1
  fi
  range_lower_bound="${diff_range%%..*}"
  range_upper_bound="${diff_range##*..}"
  range_failure_detail=""
  case "$diff_range" in
    *..*) ;;
    *) range_failure_detail="not a two-dot commit range (expected <pre>..<merge_hash>)" ;;
  esac
  if [ -z "$range_failure_detail" ] && [ -z "$range_lower_bound" ]; then
    range_failure_detail="the lower bound (<pre>) is empty — that hash was lost between command blocks; re-type it as a literal, never carry it in a shell variable"
  fi
  if [ -z "$range_failure_detail" ] && [ -z "$range_upper_bound" ]; then
    range_failure_detail="the upper bound (<merge_hash>) is empty — that hash was lost between command blocks; re-type it as a literal, never carry it in a shell variable"
  fi
  if [ -z "$range_failure_detail" ] && ! git rev-parse --verify --quiet "${range_lower_bound}^{commit}" >/dev/null 2>&1; then
    range_failure_detail="the lower bound '$range_lower_bound' is not a commit in this repository"
  fi
  if [ -z "$range_failure_detail" ] && ! git rev-parse --verify --quiet "${range_upper_bound}^{commit}" >/dev/null 2>&1; then
    range_failure_detail="the upper bound '$range_upper_bound' is not a commit in this repository"
  fi
  if [ -z "$range_failure_detail" ] && ! git diff --name-only "$diff_range" >/dev/null 2>&1; then
    range_failure_detail="git refused to read it as a diff range"
  fi
  if [ -n "$range_failure_detail" ]; then
    echo "FAIL: DO_WORK_DIFF_RANGE does not resolve: '$diff_range' — $range_failure_detail."
    echo "  Every diff-based check below would read an EMPTY diff and pass vacuously, so this is a hard failure, not a warning."
    exit 1
  fi
fi

non_dowork_count=0

# Computed once, consumed by the per-file checks below. Piping `git diff`
# straight into `grep -q` is a pipefail trap: -q exits on the first match, the
# upstream git dies with SIGPIPE, and the pipeline's non-zero status made a
# file that IS in the diff read as absent (false WARN on every modified file).
# In worktree dispatch mode $diff_range holds this REQ's merge range
# (<pre>..<merge_hash>, <pre> being the integration tip captured just before this REQ's
# FIRST --no-ff merge and <merge_hash> the latest such merge commit — which carries any
# integration seam, folded in before the merge was committed); reading it keeps the same
# "this-REQ-only" guarantee the working+staged default gives serially, because <pre> sits
# immediately before this REQ's first merge. Range definition and re-merge semantics:
# actions/work-reference.md -> Worktree Dispatch Mode (Step 1).
changed_file_list=""
if [ "$git_available" -eq 1 ]; then
  if [ -n "$diff_range" ]; then
    changed_file_list="$(git diff --name-only "$diff_range" | sort -u)"
  else
    changed_file_list="$({ git diff --name-only; git diff --staged --name-only; } | sort -u)"
  fi
fi

# --- Check 1: every listed file matches its claimed state on disk / in diff ---
while IFS= read -r summary_line; do
  # Portable extraction (no GNU-only grep -P): first backtick-quoted token, then the verb.
  file_path="$(printf '%s' "$summary_line" | sed -n 's/^[^`]*`\([^`]*\)`.*/\1/p')"
  change_verb="$(printf '%s' "$summary_line" | grep -oE '\((new|modified|modify|deleted)\)' | head -1 | tr -d '()')"
  [ -z "$file_path" ] && continue
  case "$file_path" in do-work/*) continue;; esac
  non_dowork_count=$((non_dowork_count + 1))
  case "$change_verb" in
    new)
      if [ ! -f "$file_path" ]; then
        echo "FAIL: listed (new) but not on disk: $file_path"; failure_count=$((failure_count + 1))
      fi ;;
    modified|modify)
      # Deliberately only working+staged diffs: Step 6.3 runs BEFORE this REQ's
      # commit, and including the previous commit (HEAD~1) would let a no-op
      # builder pass on the back of the last REQ's work. In worktree dispatch
      # mode the same guarantee comes from $diff_range's lower bound <pre> (the
      # tip just before this REQ's FIRST merge), not from working+staged — HEAD~1
      # would still be wrong there.
      if [ ! -f "$file_path" ]; then
        echo "FAIL: listed (modified) but not on disk: $file_path"; failure_count=$((failure_count + 1))
      elif [ "$git_available" -eq 1 ] && ! printf '%s\n' "$changed_file_list" | grep -xF "$file_path" >/dev/null; then
        echo "WARN: listed (modified) but not in working/staged diff (or the merge range, in worktree mode): $file_path"
      fi ;;
    deleted)
      if [ -f "$file_path" ]; then
        echo "FAIL: listed (deleted) but still on disk: $file_path"; failure_count=$((failure_count + 1))
      elif [ "$git_available" -eq 1 ] && ! printf '%s\n' "$changed_file_list" | grep -xF "$file_path" >/dev/null; then
        echo "WARN: listed (deleted) and absent from disk, but no deletion in working/staged diff (or the merge range, in worktree mode): $file_path — verify the path is not a typo and the file was deleted by THIS REQ"
      fi ;;
    *) echo "WARN: no (new|modified|deleted) verb on summary line: $summary_line" ;;
  esac

  # --- Check 5 (grep half): a (new) source file nothing references is dead until judged ---
  if [ "$change_verb" = "new" ] && [ -f "$file_path" ]; then
    file_base="$(basename "$file_path")"
    file_stem="${file_base%.*}"
    if [ "$git_available" -eq 1 ]; then
      reference_hits="$(git grep -l -F "$file_stem" -- . 2>/dev/null | grep -vxF "$file_path" | grep -v '^do-work/' || true)"
      if [ -z "$reference_hits" ]; then
        echo "WARN: (new) file has no static reference anywhere: $file_path — judge against the Step 6.3 exception list (entry points, config, tests, framework routes, barrels, dynamic imports)"
      fi
    fi
  fi
done <<< "$summary_bullets"

# --- "Only do-work/ paths" rule from Step 6.25: no project files means no implementation ---
if [ "$non_dowork_count" -eq 0 ]; then
  echo "FAIL: Implementation Summary lists only do-work/ paths — the REQ was not implemented (design-artifact REQs excepted; see Step 6.25)"
  failure_count=$((failure_count + 1))
fi

# --- Check 4: P-A-U box audit + debug artifacts in the diff ---
unchecked_boxes="$(grep -cE '^[[:space:]]*-[[:space:]]\[ \][[:space:]]\*\*\[(PLAN|APPLY|UNIFY)\]' "$request_file" || true)"
if [ "${unchecked_boxes:-0}" -gt 0 ]; then
  echo "FAIL: $unchecked_boxes P-A-U checkbox(es) still unchecked — the builder did not complete those phases"
  failure_count=$((failure_count + 1))
fi
if [ "$git_available" -eq 1 ]; then
  # do-work/ is excluded at the pathspec level, NOT with a `grep -v 'do-work/'`
  # on the piped lines: added-content lines carry no file path, so a content
  # grep cannot scope by file — it silently matched REQ prose that merely
  # *mentions* console.log/TODO (the REQ file is part of this diff) and FAILed
  # clean implementations. `+++` headers are dropped so a filename containing
  # TODO cannot trip the artifact grep either.
  if [ -n "$diff_range" ]; then
    debug_artifact_lines="$(git diff "$diff_range" -- . ':(exclude)do-work/' | grep -E '^\+' | grep -vE '^\+\+\+ ' | grep -nE 'console\.log|debugger|(^|[^[:alnum:]_])print\(|TODO|FIXME' || true)"
  else
    debug_artifact_lines="$({ git diff -- . ':(exclude)do-work/'; git diff --staged -- . ':(exclude)do-work/'; } | grep -E '^\+' | grep -vE '^\+\+\+ ' | grep -nE 'console\.log|debugger|(^|[^[:alnum:]_])print\(|TODO|FIXME' || true)"
  fi
  if [ -n "$debug_artifact_lines" ] && grep -qE '^[[:space:]]*-[[:space:]]\[x\][[:space:]]\*\*\[UNIFY\]' "$request_file"; then
    echo "FAIL: [UNIFY] is checked but the diff adds debug artifacts — un-check it and flag:"
    printf '%s\n' "$debug_artifact_lines" | head -10 | sed 's/^/  /'
    failure_count=$((failure_count + 1))
  fi
fi

if [ "$failure_count" -eq 0 ]; then
  echo "OK: mechanical qualification passed — judgment checks 2 (substantive), 3 (requirements traced), 6 (data flows) remain with the orchestrator"
  exit 0
fi
exit 1

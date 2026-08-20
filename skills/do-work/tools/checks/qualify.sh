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
#
# What the artifact scans read is NOT only a diff (REQ-263). In serial mode they also
# read untracked, non-ignored files whole, because Step 6.3 runs before this REQ's commit
# and a new source file the builder never staged appears in no diff at all. Ownership of
# a file's process exit — the condition that sorts an added output line into WARN or FAIL
# — is judged at the BASE revision for a path that exists there, and on the file's own
# content for one that does not. Both behaviors are pinned in
# _dev/tests/prescribed-shell-cases/qualify.sh.
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

# --- Ownership base: the revision that decides whether a file owns its process exit ---
# The probe below used to read the POST-change working copy, so an exit idiom added in
# the SAME diff as a debug print retroactively turned library code into a reporter and
# downgraded the FAIL to a WARN (REQ-254's review reproduced it; REQ-263 closes it).
# Ownership is judged at the revision the change started from instead: a file that was
# library code before this REQ stays library code for this REQ's verdict. A file that
# does not exist at the base is new and is judged on its own content — which is what
# keeps a legitimately new checker, whose prints and whose exit arrive together in one
# new file, a WARN rather than a FAIL. That case is why the categorical rule
# ("exit added in this diff ⇒ FAIL") is wrong and is not what this implements.
ownership_base=""
if [ "$git_available" -eq 1 ]; then
  if [ -n "$diff_range" ]; then
    ownership_base="${diff_range%%..*}"
  elif git rev-parse --verify --quiet HEAD >/dev/null 2>&1; then
    ownership_base="HEAD"
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
# A REQ with NO P-A-U section at all used to sail through this whole check (REQ-264).
# Every UNIFY-gated FAIL below keys on a CHECKED [UNIFY] box, and the unchecked-box FAIL
# keys on an UNCHECKED one, so a file carrying neither satisfies both by absence: the
# audit is disarmed rather than passed, and nothing said so. That is the shape REQ-254's
# own qualification "Passed" with — its review re-ran the range armed and got FAILs.
# work.md Step 6 states P-A-U phasing is mandatory, so the absence is a defect in the
# REQ, not a supported mode; it is a WARN and not a FAIL because the missing section is
# the orchestrator's paperwork rather than evidence about the code, and because REQs
# written before the section existed are still legitimately qualifiable.
# The arming condition is the [UNIFY] BOX, in any state — not the section, and not any box.
# Counting all three was the first attempt and left a hole a review caught: a legacy or
# hand-edited REQ that keeps PLAN and APPLY but drops UNIFY has a nonzero box count, so the
# warning stayed quiet, while every artifact FAIL below still keys on a CHECKED [UNIFY] and
# stayed unreachable. That REQ exited 0 with a clean OK — the exact defect this warning
# exists to expose, one step narrower.
# A [UNIFY] box present but UNCHECKED needs no warning: the unchecked-box FAIL right below
# fires on it, so that state is already loud. Only absence is silent.
unify_box_total="$(grep -cE '^[[:space:]]*-[[:space:]]\[( |x|~)\][[:space:]]\*\*\[UNIFY\]' "$request_file" || true)"
if [ "${unify_box_total:-0}" -eq 0 ]; then
  pau_box_total="$(grep -cE '^[[:space:]]*-[[:space:]]\[( |x|~)\][[:space:]]\*\*\[(PLAN|APPLY|UNIFY)\]' "$request_file" || true)"
  if [ "${pau_box_total:-0}" -eq 0 ]; then
    echo "WARN: no 'AI Execution State (P-A-U Loop)' section in this REQ — Check 4's box audit is DISARMED, not passed: every [UNIFY]-gated FAIL below is unreachable, so a debug artifact in the diff cannot fail this run. Add the section (work.md Step 6 makes P-A-U phasing mandatory) and re-run before trusting an OK."
  else
    echo "WARN: this REQ has a P-A-U section but no [UNIFY] box — Check 4's box audit is DISARMED, not passed: every [UNIFY]-gated FAIL below keys on a checked [UNIFY] line and is unreachable without one, so a debug artifact in the diff cannot fail this run. Add the missing box (work.md Step 6 makes P-A-U phasing mandatory) and re-run before trusting an OK."
  fi
fi
unchecked_boxes="$(grep -cE '^[[:space:]]*-[[:space:]]\[ \][[:space:]]\*\*\[(PLAN|APPLY|UNIFY)\]' "$request_file" || true)"
if [ "${unchecked_boxes:-0}" -gt 0 ]; then
  echo "FAIL: $unchecked_boxes P-A-U checkbox(es) still unchecked — the builder did not complete those phases"
  failure_count=$((failure_count + 1))
fi
if [ "$git_available" -eq 1 ]; then
  # Artifact tokens split by a property, never a file list (REQ-254). Markers of
  # unfinished or debug-only work (debugger, TODO, FIXME — vocabulary illustrative,
  # meant to grow) are artifacts wherever they land: nothing shipped keeps them on
  # purpose, so they FAIL on sight. Output primitives (print(, console.log — same
  # caveat) are different in kind: the same token writes a checker's success line
  # and a forgotten debug dump, and a scan cannot see intent. The distinguishing
  # condition: printed output belongs to whoever owns the process exit. A file that
  # ends its own process (the exit idioms in the regex below approximate that
  # condition) is a program with a terminal audience — a check, a CLI — so an added
  # output line there is presumptively the file's own reporting, surfaced as a WARN
  # for judgment. A file that never ends its process is library code: nothing reads
  # its stdout by contract, so the same line FAILs as leftover instrumentation.
  # This scan once FAILed a checker's only success line (REQ-244's remediation) and
  # the override went on the record — a gate that cries wolf trains people to wave
  # through FAILs.
  unfinished_marker_regex='debugger|TODO|FIXME'
  output_primitive_regex='console\.log|(^|[^[:alnum:]_])print\('
  # Code-shaped exit occurrences only (REQ-263). Two narrowings over the whole-file
  # text grep this replaced, each earned by a reproduced false WARN:
  #   1. The bare `exit N` form must be STATEMENT-shaped — beginning a statement (line
  #      start, or after `;`/`&`/`|`, or after `then`/`else`/`do`) AND terminated right
  #      after its status by end-of-line, `;`, `)`, `}`, or a comment. Prose reads
  #      "...the caller should exit 1 on failure": a word before it, words after it, so
  #      it no longer counts as ownership. A bare status-less `exit` did not match the
  #      previous regex either, and still does not — this narrows, it never widens.
  #   2. Full-line comments are dropped before the match (see file_owns_process_exit),
  #      so `# exits 1 on failure` is not ownership either.
  # The parenthesised forms are already code-shaped by their own punctuation.
  # DOCUMENTED RESIDUAL, pinned by a fixture in _dev/tests/prescribed-shell-cases/qualify.sh:
  # a prose line that both begins a statement and ends immediately after the status — a
  # docstring line reading exactly `exit 1` — still reads as ownership. Separating that
  # from real code needs a per-language parser, which is more machinery than this gate
  # is worth; the fixture exists so the boundary is a stated limit rather than a surprise.
  process_exit_regex='(^[[:space:]]*|[;&|][[:space:]]*|[[:space:]](then|else|do)[[:space:]]+)exit[[:space:]]+[0-9$][^[:space:];)}#]*[[:space:]]*([;)}#]|$)|sys\.exit[[:space:]]*\(|raise[[:space:]]+SystemExit|os\._exit[[:space:]]*\(|process\.exit[[:space:]]*\('

  # Does this path own its process exit? Judged at $ownership_base when the path exists
  # there, on the working copy when it does not (a new or untracked file).
  # `grep -c` rather than `grep -q` at the end of a pipe on purpose: -q quits on the
  # first match, the upstream grep dies of SIGPIPE, and pipefail then reports the
  # pipeline as failed — the exact trap already documented above $changed_file_list.
  # A renamed file does not exist at the base under its NEW path, so a bare cat-file probe
  # falls through to the post-change working copy and reads the file as brand new — which
  # hands a rename that also adds an exit idiom and a debug print exactly the reporter
  # exemption the base-revision rule exists to deny. A review caught this. Ownership is a
  # property of the file's IDENTITY, and git's rename detection is precisely a statement
  # about identity, so `--find-renames` is the right tool here — unlike the relocation
  # question below, which is about content and deliberately does not use it.
  # Read in shell rather than through `awk -v`, which processes escape sequences in the
  # value and would mangle a path containing a backslash.
  resolve_base_path_before_rename() {
    local current_path="$1"
    local rename_status_stream
    local change_status
    local old_path
    local new_path
    if [ -n "$diff_range" ]; then
      rename_status_stream="$(git diff --find-renames --name-status "$diff_range" 2>/dev/null || true)"
    else
      rename_status_stream="$({ git diff --find-renames --name-status; \
        git diff --staged --find-renames --name-status; } 2>/dev/null || true)"
    fi
    [ -z "$rename_status_stream" ] && return 1
    while IFS="$(printf '\t')" read -r change_status old_path new_path; do
      case "$change_status" in R*) ;; *) continue;; esac
      if [ -n "$new_path" ] && [ "$new_path" = "$current_path" ]; then
        printf '%s' "$old_path"
        return 0
      fi
    done <<< "$rename_status_stream"
    return 1
  }
  file_owns_process_exit() {
    local probe_path="$1"
    local probe_content
    local probe_hit_count
    local base_probe_path="$probe_path"
    local renamed_from_path
    if [ -n "$ownership_base" ] && ! git cat-file -e "${ownership_base}:${probe_path}" 2>/dev/null; then
      renamed_from_path="$(resolve_base_path_before_rename "$probe_path" || true)"
      if [ -n "$renamed_from_path" ]; then
        base_probe_path="$renamed_from_path"
      fi
    fi
    if [ -n "$ownership_base" ] && git cat-file -e "${ownership_base}:${base_probe_path}" 2>/dev/null; then
      probe_content="$(git show "${ownership_base}:${base_probe_path}" 2>/dev/null)"
    elif [ -f "$probe_path" ]; then
      probe_content="$(cat -- "$probe_path")"
    else
      return 1
    fi
    probe_hit_count="$(printf '%s\n' "$probe_content" \
      | grep -vE '^[[:space:]]*(#|//|\*)' \
      | grep -cE "$process_exit_regex" || true)"
    [ "${probe_hit_count:-0}" -gt 0 ]
  }

  # --- Relocated lines are not added lines (REQ-301) ---
  # Every scan below reads `^+` out of a diff (or a whole untracked file), so text that was
  # MOVED reads exactly like text that was WRITTEN. Every REQ that relocates code therefore
  # FAILed the artifact audit on markers that already existed: REQ-258 hit it on four
  # deliberate fixture TODO strings, byte-identical in the pre-change tree, and had to
  # override the FAIL with evidence. That is the dangerous direction — a gate that cries
  # wolf on a whole category of change teaches builders to wave its FAILs through, and this
  # is the gate that catches real leftover instrumentation.
  #
  # Of the two candidate fixes, this is the second: subtract from the flagged set any line
  # whose content already exists in the pre-change tree. It is chosen over `git diff -C
  # --find-copies-harder` because git's rename/copy detection is file-level, so it sees a
  # file split (REQ-258's shape) but not a hunk moved within a file or between two files
  # that both already existed — and those are relocations too. The condition is the content,
  # not the file topology.
  #
  # The REQ warns this approach "can mask a genuinely re-added marker elsewhere in the same
  # file." Mere presence at base is indeed too coarse — it cannot tell a MOVE from a
  # DUPLICATE, and REQ-263's own same-diff-exit fixture proved it by appending a second
  # copy of a line the tree already had and being excused for it. So the test is not
  # presence but COUNT: a line is relocated only when the change did not increase how many
  # times that exact text occurs in the tree. A move removes one occurrence and adds one
  # (count unchanged); a genuine addition raises the count, wherever it lands.
  #
  # Belt and braces: even a relocated line is DOWNGRADED to a named WARN, never dropped.
  # Nothing the scans found is ever silently discarded — only text whose occurrence count
  # actually grew can FAIL.
  #
  # `do-work/` is excluded from both counts for the same reason the scans exclude it: a REQ
  # file that merely discusses a TODO must not license one in shipped code. The post-change
  # count passes `--untracked` so a relocation into a not-yet-staged file counts on both
  # sides of the comparison — that is REQ-258's shape and the whole case for this fix.
  fresh_matched_lines=""
  relocated_matched_lines=""
  count_matching_lines_in_tree() {
    # $1 = revision to search, or empty for the working tree (+ untracked).
    # `git grep -c` prints `path:count` per matching file, so sum the last field. It exits
    # 1 on no match, which is a legitimate zero here, hence the `|| true` inside the
    # substitution rather than a bare pipeline.
    local search_revision="$1"
    local search_pattern="$2"
    if [ -n "$search_revision" ]; then
      printf '%s' "$(git grep -c -F -e "$search_pattern" "$search_revision" \
        -- . ':(exclude)do-work/' 2>/dev/null | awk -F: '{total += $NF} END {print total + 0}' || true)"
    else
      printf '%s' "$(git grep -c -F --untracked -e "$search_pattern" \
        -- . ':(exclude)do-work/' 2>/dev/null | awk -F: '{total += $NF} END {print total + 0}' || true)"
    fi
  }
  partition_matched_lines_by_relocation() {
    local matched_block="$1"
    local matched_line
    local line_content
    local base_occurrences
    local current_occurrences
    fresh_matched_lines=""
    relocated_matched_lines=""
    while IFS= read -r matched_line; do
      [ -z "$matched_line" ] && continue
      # Drop `grep -n`'s line number, then the diff's leading `+` when there is one
      # (the untracked scan reads whole files, so its lines carry no `+`).
      line_content="${matched_line#*:}"
      line_content="${line_content#+}"
      base_occurrences=0
      current_occurrences=0
      if [ -n "$ownership_base" ] && [ -n "$line_content" ]; then
        # -e marks the pattern explicitly, so content starting with `-` is data, not options.
        base_occurrences="$(count_matching_lines_in_tree "$ownership_base" "$line_content")"
        current_occurrences="$(count_matching_lines_in_tree "" "$line_content")"
      fi
      if [ "${base_occurrences:-0}" -gt 0 ] \
        && [ "${current_occurrences:-0}" -le "${base_occurrences:-0}" ]; then
        relocated_matched_lines="${relocated_matched_lines}${matched_line}"$'\n'
      else
        fresh_matched_lines="${fresh_matched_lines}${matched_line}"$'\n'
      fi
    done <<< "$matched_block"
  }
  # do-work/ is excluded at the pathspec level, NOT with a `grep -v 'do-work/'`
  # on the piped lines: added-content lines carry no file path, so a content
  # grep cannot scope by file — it silently matched REQ prose that merely
  # *mentions* console.log/TODO (the REQ file is part of this diff) and FAILed
  # clean implementations. `+++` headers are dropped so a filename containing
  # TODO cannot trip the artifact grep either.
  if [ -n "$diff_range" ]; then
    debug_artifact_lines="$(git diff "$diff_range" -- . ':(exclude)do-work/' | grep -E '^\+' | grep -vE '^\+\+\+ ' | grep -nE "$unfinished_marker_regex" || true)"
  else
    debug_artifact_lines="$({ git diff -- . ':(exclude)do-work/'; git diff --staged -- . ':(exclude)do-work/'; } | grep -E '^\+' | grep -vE '^\+\+\+ ' | grep -nE "$unfinished_marker_regex" || true)"
  fi
  if [ -n "$debug_artifact_lines" ] && grep -qE '^[[:space:]]*-[[:space:]]\[x\][[:space:]]\*\*\[UNIFY\]' "$request_file"; then
    partition_matched_lines_by_relocation "$debug_artifact_lines"
    if [ -n "$fresh_matched_lines" ]; then
      echo "FAIL: [UNIFY] is checked but the diff adds debug artifacts — un-check it and flag:"
      printf '%s' "$fresh_matched_lines" | head -10 | sed 's/^/  /'
      failure_count=$((failure_count + 1))
    fi
    if [ -n "$relocated_matched_lines" ]; then
      echo "WARN: debug-artifact marker(s) in the diff whose exact text already exists in the pre-change tree — read as relocated, not added, so they do not fail this run; confirm the move was the intent:"
      printf '%s' "$relocated_matched_lines" | head -10 | sed 's/^/  /'
    fi
  fi
  # The output-primitive half needs file attribution (the ownership condition reads
  # the file that gained the line), so it walks the changed files instead of the
  # raw diff stream. do-work/ is skipped per path — the same boundary as above.
  while IFS= read -r changed_path; do
    [ -z "$changed_path" ] && continue
    case "$changed_path" in do-work/*) continue;; esac
    if [ -n "$diff_range" ]; then
      added_output_lines="$(git diff "$diff_range" -- "$changed_path" | grep -E '^\+' | grep -vE '^\+\+\+ ' | grep -nE "$output_primitive_regex" || true)"
    else
      added_output_lines="$({ git diff -- "$changed_path"; git diff --staged -- "$changed_path"; } | grep -E '^\+' | grep -vE '^\+\+\+ ' | grep -nE "$output_primitive_regex" || true)"
    fi
    [ -z "$added_output_lines" ] && continue
    if file_owns_process_exit "$changed_path"; then
      # The WARN prints its matched lines exactly as the FAIL branch does (REQ-263).
      # Without them "confirm from the diff" cost a manual dig for the very lines the
      # check had already found, which is how a WARN gets waved through.
      echo "WARN: added print(/console.log line(s) in $changed_path read as the file's own reporting — it owns its process exit, so printed output is presumptively contract, not a debug artifact; confirm from these lines:"
      printf '%s\n' "$added_output_lines" | head -10 | sed 's/^/  /'
    elif grep -qE '^[[:space:]]*-[[:space:]]\[x\][[:space:]]\*\*\[UNIFY\]' "$request_file"; then
      partition_matched_lines_by_relocation "$added_output_lines"
      if [ -n "$fresh_matched_lines" ]; then
        echo "FAIL: [UNIFY] is checked but the diff adds output line(s) to $changed_path, which never ends its own process — no terminal audience, so they read as leftover instrumentation:"
        printf '%s' "$fresh_matched_lines" | head -10 | sed 's/^/  /'
        failure_count=$((failure_count + 1))
      fi
      if [ -n "$relocated_matched_lines" ]; then
        echo "WARN: output line(s) in $changed_path whose exact text already exists in the pre-change tree — read as relocated, not added, so they do not fail this run; confirm the move was the intent:"
        printf '%s' "$relocated_matched_lines" | head -10 | sed 's/^/  /'
      fi
    fi
  done <<< "$changed_file_list"

  # --- Both artifact scans, over untracked files, whole-file (REQ-263 addendum) ---
  # Serial mode only. Step 6.3 runs BEFORE this REQ's commit, so a new source file the
  # builder never staged is in neither `git diff` nor `git diff --staged`: both scans
  # above never saw it, Check 1's (new) branch only tests that it is on disk, and a
  # checked [UNIFY] could therefore ship with leftover instrumentation inside it.
  # `git ls-files --others --exclude-standard` is the source: it lists untracked files
  # individually — unlike `git status --porcelain`, which collapses a wholly-untracked
  # directory into one row — and it drops correctly-ignored paths, so it doubles as the
  # ignore filter (_dev/primes/prime-shell-commands.md).
  # An untracked file has never been committed, so every line in it is added; scanning it
  # whole does not weaken the added-lines-only contract the diffed paths get.
  # This is a SEPARATE walk, deliberately: folding untracked paths into
  # $changed_file_list would change Check 1's (modified)/(deleted) WARN behavior, which
  # must stay exactly as it was. do-work/ is skipped per path, the same boundary as above.
  # In worktree dispatch mode $diff_range reads committed work, so this block is skipped.
  if [ -z "$diff_range" ]; then
    unify_box_checked=0
    grep -qE '^[[:space:]]*-[[:space:]]\[x\][[:space:]]\*\*\[UNIFY\]' "$request_file" && unify_box_checked=1
    while IFS= read -r untracked_path; do
      [ -z "$untracked_path" ] && continue
      case "$untracked_path" in do-work/*) continue;; esac
      [ -f "$untracked_path" ] || continue
      untracked_marker_lines="$(grep -nE "$unfinished_marker_regex" -- "$untracked_path" || true)"
      if [ -n "$untracked_marker_lines" ] && [ "$unify_box_checked" -eq 1 ]; then
        partition_matched_lines_by_relocation "$untracked_marker_lines"
        if [ -n "$fresh_matched_lines" ]; then
          echo "FAIL: [UNIFY] is checked but untracked file $untracked_path carries debug artifacts — un-check it and flag:"
          printf '%s' "$fresh_matched_lines" | head -10 | sed 's/^/  /'
          failure_count=$((failure_count + 1))
        fi
        if [ -n "$relocated_matched_lines" ]; then
          echo "WARN: debug-artifact marker(s) in untracked file $untracked_path whose exact text already exists in the pre-change tree — read as relocated, not added, so they do not fail this run; confirm the move was the intent:"
          printf '%s' "$relocated_matched_lines" | head -10 | sed 's/^/  /'
        fi
      fi
      untracked_output_lines="$(grep -nE "$output_primitive_regex" -- "$untracked_path" || true)"
      [ -z "$untracked_output_lines" ] && continue
      if file_owns_process_exit "$untracked_path"; then
        echo "WARN: print(/console.log line(s) in untracked file $untracked_path read as the file's own reporting — it owns its process exit, so printed output is presumptively contract, not a debug artifact; confirm from these lines:"
        printf '%s\n' "$untracked_output_lines" | head -10 | sed 's/^/  /'
      elif [ "$unify_box_checked" -eq 1 ]; then
        partition_matched_lines_by_relocation "$untracked_output_lines"
        if [ -n "$fresh_matched_lines" ]; then
          echo "FAIL: [UNIFY] is checked but untracked file $untracked_path carries output line(s), and it never ends its own process — no terminal audience, so they read as leftover instrumentation:"
          printf '%s' "$fresh_matched_lines" | head -10 | sed 's/^/  /'
          failure_count=$((failure_count + 1))
        fi
        if [ -n "$relocated_matched_lines" ]; then
          echo "WARN: output line(s) in untracked file $untracked_path whose exact text already exists in the pre-change tree — read as relocated, not added, so they do not fail this run; confirm the move was the intent:"
          printf '%s' "$relocated_matched_lines" | head -10 | sed 's/^/  /'
        fi
      fi
    done <<< "$(git ls-files --others --exclude-standard)"
  fi
fi

if [ "$failure_count" -eq 0 ]; then
  echo "OK: mechanical qualification passed — judgment checks 2 (substantive), 3 (requirements traced), 6 (data flows) remain with the orchestrator"
  exit 0
fi
exit 1

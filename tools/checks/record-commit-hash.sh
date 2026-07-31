#!/usr/bin/env bash
# record-commit-hash.sh — the write-back half of actions/work.md's Commit Phase (Step 9):
# put the implementation commit's hash into an archived REQ's frontmatter `commit:` field.
#
# WHY THIS IS A SCRIPT AND NOT PROSE. The step used to be described in
# actions/work-reference.md and executed free-form ("replace the commit: line, or add it if
# missing"). In one consumer repo those free-form edits truncated SIX archived REQ files to
# 0 bytes — 9 KB to 26 KB of irreplaceable decision trail each. Every one of the metadata
# commits that landed the damage read "1 file changed, N deletions(-)", and nothing looked at
# that number before committing; the commit message claimed success, so the loss stayed
# invisible for weeks. Recovery was possible only because the content still existed at the
# preceding implementation commit. Every guard below exists to make that one edit impossible
# to get wrong, which is why they all run BEFORE the file is replaced and why a tripped guard
# leaves the REQ byte-identical to how it was found.
#
# Usage:
#   tools/checks/record-commit-hash.sh <req-file> <hash>
#   tools/checks/record-commit-hash.sh --verify <req-file> <hash>
#
#   <req-file>  the ARCHIVED REQ file (do-work/archive/.../REQ-NNN-slug.md).
#   <hash>      the IMPLEMENTATION commit's short hash. Serially that is
#               `git rev-parse --short HEAD` — the Step 9 commit just made. In worktree
#               dispatch mode it is the --no-ff <merge_hash> held since Step 6, never HEAD
#               (HEAD there names the changelog commit, not the implementation).
#   Run from the project root — the directory containing do-work/.
#
#   --verify inspects what the metadata commit ACTUALLY committed. Run it after `git commit`.
#   It is the only check that can catch a content-mutating pre-commit hook (a formatter, a
#   lint --fix, a whitespace stripper) rewriting the file after every pre-commit guard has
#   already passed — which is itself a sufficient mechanism for the incident above. What it
#   proves is the committed PATCH: that HEAD introduced exactly one new line for this file, that
#   the line is the expected `commit:` field, and that the only line it removed is the one HEAD^'s
#   own frontmatter held there (none at all, on an insert). Matching a removed `commit:` line by
#   SHAPE instead would let a hook delete a BODY `commit:` line unnoticed — archived REQs really
#   do have those, where they quote the schema. It does NOT merely compare the committed
#   blob against the worktree — that comparison alone proves nothing against the commonest
#   hook shape, which rewrites the worktree file and re-stages it so both sides move together.
#   Where the patch cannot be isolated (root commit, merge HEAD, or the file added by this
#   very commit) the mode says so explicitly rather than reporting a guarantee it did not make.
#
# This script does NOT stage or commit. `git add` and `git commit` stay with the caller, so
# that a tripped guard is a stop signal on an unmodified file rather than a rollback.
#
# Exit 0 — the field records <hash>:
#     "OK: ..."   this run wrote it (or the edit was already on disk but not yet committed).
#     "NOOP: ..." it was already recorded AND already committed. Make no metadata commit.
# Exit 1 — a guard tripped. STOP: do not retry, do not hand-edit around it, do not commit.
#     The REQ file is left as it was found unless a printed line says otherwise.
# Exit 2 — usage error: bad argument count or hash shape, or <req-file> is not a usable REQ
#     file (missing, 0 bytes, a symlink, CRLF, no frontmatter, or no valid id:).
#
# This script never runs `git show <hash>`: on a merge commit — worktree dispatch mode's
# <merge_hash> — that prints an EMPTY combined diff, so it would validate nothing. It only
# rev-parses. Consumers reading `commit:` as a diff source must use
# `git show --first-parent -m <hash>`.
set -uo pipefail

# The byte arithmetic below depends on ${#var} counting BYTES rather than locale characters,
# and byte-oriented matching keeps awk from erroring on an invalid multibyte sequence inside
# REQ prose. Only ASCII patterns are ever matched.
export LC_ALL=C

verify_mode=0
if [ "${1:-}" = "--verify" ]; then
  verify_mode=1
  shift
fi

# Exactly two: an unquoted path containing a space arrives as three arguments, and silently
# editing $1 would rewrite the wrong file.
if [ "$#" -ne 2 ]; then
  echo "usage: $0 [--verify] <req-file> <hash>   — exactly 2 arguments; quote a path containing spaces" >&2
  exit 2
fi
request_file="$1"
commit_hash="$2"

# 7-40 lowercase hex. This is also what rejects a copy-paste of the literal placeholder
# <hash> out of the procedure's own code block.
if ! printf '%s' "$commit_hash" | grep -Eq '^[0-9a-f]{7,40}$'; then
  echo "usage: $0 [--verify] <req-file> <hash> — '$commit_hash' is not a 7-40 character lowercase hex hash." >&2
  echo "  Pass the resolved hash, never the literal placeholder <hash>. Lowercase it if git printed capitals." >&2
  exit 2
fi
if [ -L "$request_file" ]; then
  echo "usage: $0 [--verify] <req-file> <hash> — '$request_file' is a symlink; this script rewrites by atomic rename, which would replace the link with a regular file." >&2
  exit 2
fi
if [ ! -f "$request_file" ]; then
  echo "usage: $0 [--verify] <req-file> <hash> — '$request_file' is not a file" >&2
  exit 2
fi

git_available=0
git rev-parse --git-dir >/dev/null 2>&1 && git_available=1

# ---------------------------------------------------------------------------
# Frontmatter readers. Scope is line 1 `---` up to the next `---`, so a `commit:` in body
# prose or inside a fenced YAML sample is unreachable — several archived REQs quote the
# schema, which is exactly how a file-wide sed corrupts one. Keys are anchored at column 0:
# an indented `  commit:` is a nested key under some other mapping, not the schema field.
#
# Defined ahead of --verify (not just ahead of the write path) so BOTH halves decide what
# counts as the frontmatter `commit:` line with the same parser. --verify reads the PARENT
# blob through these; a second, hand-rolled reader there would be free to drift from the
# writer's idea of frontmatter, which is precisely the confusion a body `commit:` exploits.
# ---------------------------------------------------------------------------
frontmatter_line_for() {
  awk -v field_name="$2" '
    NR == 1 && $0 == "---" { inside_frontmatter = 1; next }
    inside_frontmatter && $0 == "---" { exit }
    inside_frontmatter && index($0, field_name ":") == 1 { print; exit }
  ' "$1"
}
frontmatter_count_for() {
  awk -v field_name="$2" '
    NR == 1 && $0 == "---" { inside_frontmatter = 1; next }
    inside_frontmatter && $0 == "---" { exit }
    inside_frontmatter && index($0, field_name ":") == 1 { hit_count++ }
    END { print hit_count + 0 }
  ' "$1"
}
frontmatter_value_for() {
  local field_line field_value
  field_line="$(frontmatter_line_for "$1" "$2")"
  [ -z "$field_line" ] && return 0
  field_value="${field_line#*:}"
  field_value="${field_value%%#*}"                                # trailing YAML comment
  field_value="${field_value#"${field_value%%[![:space:]]*}"}"    # ltrim
  field_value="${field_value%"${field_value##*[![:space:]]}"}"    # rtrim
  field_value="${field_value%\"}"; field_value="${field_value#\"}"
  field_value="${field_value%\'}"; field_value="${field_value#\'}"
  printf '%s' "$field_value"
}

# ---------------------------------------------------------------------------
# --verify: read back what actually got committed.
# ---------------------------------------------------------------------------
if [ "$verify_mode" -eq 1 ]; then
  if [ "$git_available" -eq 0 ]; then
    echo "FAIL: --verify needs a git repository; '$request_file' is not in one, so there is no committed blob to read back."
    exit 1
  fi
  committed_full_name="$(git ls-files --full-name -- "$request_file" 2>/dev/null)"
  if [ -z "$committed_full_name" ]; then
    echo "FAIL: $request_file is not tracked by git — the metadata commit did not land, so there is nothing to verify."
    exit 1
  fi
  committed_bytes="$(git cat-file -s "HEAD:$committed_full_name" 2>/dev/null || echo 0)"
  worktree_bytes="$(wc -c < "$request_file" | tr -d '[:space:]')"
  # `grep -c` reads to EOF, so it cannot SIGPIPE the upstream `git show` the way `grep -q`
  # would under pipefail; `|| true` absorbs grep's no-match exit status. Matching the whole
  # line (-x) is what keeps a `commit:` inside body prose from counting.
  committed_hash_lines="$(git show "HEAD:$committed_full_name" 2>/dev/null | grep -c -x "commit: $commit_hash" || true)"
  if [ "$committed_bytes" -ne "$worktree_bytes" ] || [ "$committed_hash_lines" -ne 1 ]; then
    echo "FAIL: the committed content of $request_file is not what was verified."
    echo "  committed: $committed_bytes bytes · on disk: $worktree_bytes bytes · 'commit: $commit_hash' lines in the commit: $committed_hash_lines (expected 1)"
    echo "  Something rewrote the file between the guards and the commit — a content-mutating pre-commit hook is the usual cause."
    echo "  Recover:  git checkout HEAD~1 -- $request_file   (the pre-metadata-commit content), then fix that hook."
    echo "  Do NOT re-run the write-back until the hook is fixed, and never reach for --no-verify."
    exit 1
  fi

  # The read-back above is necessary but NOT sufficient, and treating it as sufficient is how
  # this mode came to over-promise. It compares the committed blob against the WORKTREE, and
  # the commonest content-mutating hook (lint-staged and friends) rewrites the worktree file
  # and re-stages it — both sides move together, so ANY amount of body corruption compares
  # equal for as long as one `commit: <hash>` line survives. What actually proves the claim is
  # the committed patch: a metadata commit may introduce exactly one new line for this path,
  # and that line must be the commit: field.
  #
  # Each skip below is a shape where no single-line change can be isolated, so the check is
  # declared skipped rather than silently passed. `^` is quoted throughout: it is a glob
  # operator in zsh and an escape character in cmd.exe, and this idiom gets copy-pasted.
  patch_check_skip_reason=""
  if ! git rev-parse --verify -q 'HEAD^{commit}' >/dev/null 2>&1; then
    patch_check_skip_reason="HEAD does not resolve to a commit"
  elif ! git rev-parse --verify -q 'HEAD^' >/dev/null 2>&1; then
    patch_check_skip_reason="HEAD is a root commit, so there is no parent to diff it against"
  elif git rev-parse --verify -q 'HEAD^2' >/dev/null 2>&1; then
    patch_check_skip_reason="HEAD is a merge commit, which a metadata commit never is — this is not the commit the procedure describes"
  elif ! git cat-file -e "HEAD^:$committed_full_name" 2>/dev/null; then
    patch_check_skip_reason="$committed_full_name does not exist in HEAD^, so this commit ADDED it (an archive move committed together with the hash) and there is no one-line change to isolate"
  fi

  # --unified=0 drops context lines; --no-ext-diff / --no-textconv keep a configured diff
  # driver from deciding what this guard gets to see; --no-renames keeps a single-path diff
  # deterministic.
  committed_patch_text=""
  if [ -z "$patch_check_skip_reason" ]; then
    committed_patch_text="$(git --no-pager diff --unified=0 --no-color --no-ext-diff --no-textconv --no-renames 'HEAD^' HEAD -- "$committed_full_name" 2>/dev/null)"
    # Nothing to inspect means HEAD is not the metadata commit — --verify was run later than
    # the procedure calls for. Say that plainly instead of failing (there is no evidence of
    # damage) and instead of passing quietly (there is no evidence of correctness either).
    if [ -z "$committed_patch_text" ]; then
      patch_check_skip_reason="HEAD does not touch $committed_full_name, so HEAD is not this file's metadata commit — run --verify immediately after that commit"
    fi
  fi

  if [ -n "$patch_check_skip_reason" ]; then
    echo "INFO: the committed-patch check was skipped — $patch_check_skip_reason."
    echo "  Only the read-back ran, and it cannot see a same-size body rewrite (a hook that re-stages its own edit moves the worktree too)."
    echo "  If a content-mutating hook is suspected here, read the diff yourself: git show --stat HEAD -- $request_file"
    echo "OK (read-back only): $request_file reads back at $committed_bytes bytes with commit: $commit_hash recorded once."
    exit 0
  fi

  # Net added / net removed: an added line cancelled by an identical removed line is a
  # relocation, not a content change. That cancellation is what lets the legitimate
  # missing-trailing-newline rewrite (the last body line reappears verbatim) through without
  # loosening the count to "one or two". Hunk bodies only — `@@` gates out the ---/+++ headers
  # without a prefix match that a body line beginning `++` could impersonate.
  # shellcheck disable=SC2016  # an awk program: $0 must reach awk unexpanded
  net_line_filter='
    /^@@/ { in_hunk = 1; next }
    !in_hunk { next }
    substr($0, 1, 1) == "+" { text = substr($0, 2); added_order[++added_total] = text; added_count[text]++; next }
    substr($0, 1, 1) == "-" { text = substr($0, 2); removed_order[++removed_total] = text; removed_count[text]++; next }
    END {
      if (wanted_side == "added") {
        for (i = 1; i <= added_total; i++) {
          text = added_order[i]
          if (removed_count[text] > 0) { removed_count[text]--; continue }
          print text
        }
      } else {
        for (i = 1; i <= removed_total; i++) {
          text = removed_order[i]
          if (added_count[text] > 0) { added_count[text]--; continue }
          print text
        }
      }
    }'
  net_added_lines="$(printf '%s\n' "$committed_patch_text" | awk -v wanted_side=added "$net_line_filter")"
  net_removed_lines="$(printf '%s\n' "$committed_patch_text" | awk -v wanted_side=removed "$net_line_filter")"
  net_added_count="$(printf '%s' "$net_added_lines" | awk 'END { print NR + 0 }')"
  net_removed_count="$(printf '%s' "$net_removed_lines" | awk 'END { print NR + 0 }')"

  # What the removal side is allowed to be, read off the PARENT's frontmatter rather than
  # matched by shape. `commit:*` was too loose: an archived REQ may quote the schema in its
  # BODY (several do — see the reader comment above), so a hook that drops a body `commit:`
  # line while the write-back inserts the frontmatter one nets +1/-1 and read as a legitimate
  # replace, passing with a message that claimed the patch was one line. Process substitution
  # feeds the parent blob to the same reader the write path uses; `git cat-file -e` on
  # `HEAD^:` already succeeded above (an absent parent file is a declared skip), so an empty
  # result here means the parent genuinely had no frontmatter `commit:` field.
  parent_frontmatter_commit_line="$(frontmatter_line_for \
    <(git cat-file blob "HEAD^:$committed_full_name" 2>/dev/null) commit)"
  if [ -n "$parent_frontmatter_commit_line" ]; then
    expected_removed_count=1                                        # a replace
    expected_removed_lines="$parent_frontmatter_commit_line"
  else
    expected_removed_count=0                                        # an insert
    expected_removed_lines=""
  fi

  # Tracked separately from the added side so the diagnosis names the half that actually failed:
  # a body rewrite fails on additions, and pointing at body `commit:` lines there would mislead.
  removal_side_ok=1
  [ "$net_removed_count" -eq "$expected_removed_count" ] || removal_side_ok=0
  [ "$net_removed_lines" = "$expected_removed_lines" ] || removal_side_ok=0

  committed_patch_ok=1
  [ "$net_added_count" -eq 1 ] || committed_patch_ok=0
  [ "$net_added_lines" = "commit: $commit_hash" ] || committed_patch_ok=0
  [ "$removal_side_ok" -eq 1 ] || committed_patch_ok=0

  if [ "$committed_patch_ok" -ne 1 ]; then
    expected_removal_description="0 lines (HEAD^ has no frontmatter commit: field, so this is an insert)"
    if [ "$expected_removed_count" -eq 1 ]; then
      expected_removal_description="exactly 1: HEAD^'s own frontmatter line '$parent_frontmatter_commit_line'"
    fi
    echo "FAIL: the metadata commit changed more than the 'commit:' line of $request_file."
    echo "  net added: $net_added_count line(s) (expected exactly 1: 'commit: $commit_hash')"
    echo "  net removed: $net_removed_count line(s) (expected $expected_removal_description)"
    # Only when the removed text really is a `commit:` line — on a plain body rewrite it is not,
    # and naming the schema-quoting hazard there would send the reader after the wrong cause.
    if [ "$removal_side_ok" -ne 1 ]; then
      case "$net_removed_lines" in
        commit:*) echo "  That removed 'commit:' line is NOT the one HEAD^ had in its frontmatter, so it came from the BODY — several archived REQs quote the schema there, and losing one is silent content loss." ;;
      esac
    fi
    printf '%s\n' "$committed_patch_text" | head -n 40 | sed 's/^/    /'
    echo "  A read-back of the blob alone would have PASSED this: a pre-commit hook that rewrites"
    echo "  the file and re-stages it moves the worktree too, so the sizes agree while the body is gone."
    echo "  Recover:  git checkout 'HEAD^' -- $request_file   (the pre-metadata-commit content), then fix that hook."
    echo "  Do NOT re-run the write-back until the hook is fixed, and never reach for --no-verify."
    exit 1
  fi

  verified_patch_note=""
  if [ "$expected_removed_count" -eq 1 ]; then
    verified_patch_note=" replacing HEAD^'s '$parent_frontmatter_commit_line'"
  fi
  echo "OK: $request_file reads back at $committed_bytes bytes, and HEAD's patch for it is exactly the one line 'commit: $commit_hash'$verified_patch_note."
  exit 0
fi

# ---------------------------------------------------------------------------
# Input guards. An ALREADY-empty REQ is the aftermath of the truncation this script exists
# to prevent, not an input — recording a hash into it would commit the loss.
# ---------------------------------------------------------------------------
if [ ! -s "$request_file" ]; then
  echo "usage: $0 [--verify] <req-file> <hash> — '$request_file' is 0 bytes. That is the truncation this script exists to prevent, not an input." >&2
  echo "  Recover the content first:  git show HEAD:$request_file   (or from the implementation commit: git show $commit_hash:$request_file)" >&2
  exit 2
fi
# `awk '$0 == "---"'` cannot recognise a CRLF delimiter — $0 keeps the \r — so on a CRLF file
# the frontmatter would go undetected, the edit would land nowhere, and the run would look
# like a silent no-op. Refuse rather than guess.
if grep -q $'\r' "$request_file"; then
  echo "usage: $0 [--verify] <req-file> <hash> — '$request_file' has CRLF line endings; normalise to LF first (frontmatter delimiters are matched exactly)." >&2
  exit 2
fi

request_directory="$(dirname "$request_file")"

if [ "$(frontmatter_count_for "$request_file" id)" -eq 0 ]; then
  echo "usage: $0 [--verify] <req-file> <hash> — '$request_file' has no frontmatter block with an 'id:' field; this is not a REQ file." >&2
  exit 2
fi
request_id="$(frontmatter_value_for "$request_file" id)"
if ! printf '%s' "$request_id" | grep -Eq '^(REQ|UR)-[0-9]+$'; then
  echo "usage: $0 [--verify] <req-file> <hash> — frontmatter id: '$request_id' is not REQ-NNN." >&2
  exit 2
fi

existing_commit_count="$(frontmatter_count_for "$request_file" commit)"
# Two `commit:` keys is ambiguous under last-wins YAML parsing: replacing one leaves the
# other for readers to find, and replacing both is a 2/2 diff every size guard would reject.
if [ "$existing_commit_count" -gt 1 ]; then
  echo "FAIL: $request_file frontmatter has $existing_commit_count 'commit:' lines — ambiguous (a YAML reader takes the last one)."
  echo "  Nothing was written. Delete the duplicates by hand, then re-run this script."
  exit 1
fi
existing_commit_line="$(frontmatter_line_for "$request_file" commit)"
existing_commit_value="$(frontmatter_value_for "$request_file" commit)"
pre_edit_status="$(frontmatter_value_for "$request_file" status)"

case "$request_file" in
  */archive/*) ;;
  *) echo "WARN: $request_file is not under an archive/ directory — the hash write-back belongs on the ARCHIVED copy (Step 8 moves the file before Step 9)." ;;
esac
# WARN, never FAIL: the stamping rule applies to EVERY terminal flip, and failed requests get
# committed too — so `failed` and `cancelled` REQs legitimately carry a commit hash. Only an
# unfinished REQ is worth remarking on, and even then the value is left untouched.
case "$pre_edit_status" in
  completed|completed-with-issues|failed|cancelled) ;;
  *) echo "WARN: status: '$pre_edit_status' is not a terminal status (completed / completed-with-issues / failed / cancelled) — recording a commit hash on an unfinished REQ. Proceeding; the value is left untouched." ;;
esac

# ---------------------------------------------------------------------------
# Git probes. Every git-dependent guard is optional-by-detection, so a non-git repo or an
# out-of-git do-work/ tree still gets the full content guard set below.
# ---------------------------------------------------------------------------
head_exists=0
path_tracked=0
tracked_full_name=""
if [ "$git_available" -eq 1 ]; then
  git rev-parse --verify --quiet HEAD >/dev/null 2>&1 && head_exists=1
  # ls-files FIRST: a tracked file can also match an ignore pattern, and a tracked file must
  # still be diffed against HEAD. check-ignore is only the secondary signal, and it is
  # cwd-relative while interior-slash patterns are root-anchored — another reason not to
  # lead with it.
  if git ls-files --error-unmatch -- "$request_file" >/dev/null 2>&1; then
    path_tracked=1
    tracked_full_name="$(git ls-files --full-name -- "$request_file")"
  fi

  # A hash the repo cannot resolve makes the REQ a completion anomaly on `do-work board` and
  # breaks review-work / present-work traceability.
  if ! git rev-parse --verify --quiet "${commit_hash}^{commit}" >/dev/null 2>&1; then
    echo "FAIL: '$commit_hash' does not resolve to a commit in this repository. Nothing was written."
    echo "  Serial mode: the hash is \`git rev-parse --short HEAD\` run right after the Step 9 commit."
    echo "  Worktree dispatch mode: it is the <merge_hash> literal held since Step 6, NOT HEAD (HEAD is the changelog commit)."
    echo "  If git called the hash ambiguous, pass a longer prefix."
    exit 1
  fi
  # WARN, not FAIL: broken provenance is not data loss, and both legitimate shapes pass —
  # serially the hash IS HEAD, and in worktree mode <merge_hash> is HEAD's ancestor because
  # the changelog commit sits on top of it.
  if [ "$head_exists" -eq 1 ] && ! git merge-base --is-ancestor "$commit_hash" HEAD 2>/dev/null; then
    echo "WARN: $commit_hash is not an ancestor of HEAD — the recorded implementation commit is not in this branch's history (an unmerged builder branch?). Recording it anyway."
  fi
  # `^2` is quoted: ^ is a glob operator in zsh and an escape character in cmd.exe, and this
  # idiom gets copied into interactive shells.
  if git rev-parse --verify --quiet "${commit_hash}^2" >/dev/null 2>&1; then
    echo "INFO: $commit_hash is a merge commit (worktree dispatch mode's <merge_hash>). Consumers reading commit: as a DIFF source must use \`git show --first-parent -m $commit_hash\` — a plain \`git show\` on a merge prints an empty combined diff."
  fi
else
  echo "INFO: not a git repository — the hash-resolution, size-floor and diff guards are skipped; the content guards below still run in full."
fi

# ---------------------------------------------------------------------------
# Pre-edit content guards. These are the only guards that can see damage done BEFORE this
# script ran: the pre-image every later guard compares against IS the damaged file, so it
# cannot detect its own corruption.
# ---------------------------------------------------------------------------
pre_edit_bytes="$(wc -c < "$request_file" | tr -d '[:space:]')"
pre_edit_lines="$(wc -l < "$request_file" | tr -d '[:space:]')"
# A file with no trailing newline gains one from awk's print. Left unaccounted for, that
# byte shows the last body line as changed and trips the single-line-change guard on a
# perfectly legitimate file.
trailing_newline_added=0
if [ -n "$(tail -c 1 "$request_file")" ]; then
  trailing_newline_added=1
  echo "INFO: $request_file has no trailing newline; one will be added (1 byte), accounted for in the size guards below."
fi

pending_insertions=0
pending_deletions=0
if [ "$path_tracked" -eq 1 ] && [ "$head_exists" -eq 1 ]; then
  head_blob_bytes="$(git cat-file -s "HEAD:$tracked_full_name" 2>/dev/null || true)"
  # THE incident guard, expressed as a size floor: a metadata commit changes one line, so an
  # archived REQ can never approach half its committed size. This catches 26 KB -> 0 bytes
  # and every partial truncation, whatever caused it.
  if [ -n "$head_blob_bytes" ] && [ "$head_blob_bytes" -gt 0 ] && [ "$((pre_edit_bytes * 2))" -lt "$head_blob_bytes" ]; then
    echo "FAIL: $request_file is $pre_edit_bytes bytes on disk but $head_blob_bytes bytes in HEAD — content was lost BEFORE this script ran."
    echo "  Nothing was written. Recording a hash now is exactly how archived REQs got committed at 0 bytes."
    echo "  Recover first:  git checkout HEAD -- $request_file    (inspect with: git diff HEAD -- $request_file)"
    exit 1
  fi
  # --no-renames keeps the numbers deterministic for a single-path diff.
  pending_numstat="$(git diff --numstat --no-renames HEAD -- "$request_file")"
  if [ -n "$pending_numstat" ]; then
    pending_insertions="$(printf '%s\n' "$pending_numstat" | awk 'NR == 1 { print $1 }')"
    pending_deletions="$(printf '%s\n' "$pending_numstat" | awk 'NR == 1 { print $2 }')"
    case "${pending_insertions}${pending_deletions}" in
      *[!0-9]*)
        echo "FAIL: git reports $request_file as binary (numstat '$pending_insertions'/'$pending_deletions') — an archived REQ must be UTF-8 text. Nothing was written."
        exit 1 ;;
    esac
  fi
elif [ "$git_available" -eq 1 ] && [ "$path_tracked" -eq 0 ]; then
  echo "INFO: $request_file is not tracked by git (this project may keep do-work/ out of git) — the size-floor and diff guards are skipped; the content guards below still run in full."
fi

# ---------------------------------------------------------------------------
# The edit. Temps live in the REQ's own directory so `mv` is a same-filesystem atomic
# rename, and `cp -p` seeds each temp with the ORIGINAL's mode so mktemp's 0600 cannot ride
# the rename onto the REQ file.
# ---------------------------------------------------------------------------
temp_file=""
backup_file=""
# shellcheck disable=SC2329  # invoked indirectly by the EXIT trap below
remove_temp_files() {
  [ -n "$temp_file" ] && rm -f "$temp_file"
  [ -n "$backup_file" ] && rm -f "$backup_file"
  return 0
}
trap remove_temp_files EXIT

post_edit_bytes="$pre_edit_bytes"
changed_line_count=0
edit_kind="none"

if [ "$existing_commit_value" != "$commit_hash" ]; then
  if [ "$existing_commit_count" -eq 1 ]; then edit_kind="replaced"; else edit_kind="inserted"; fi

  temp_file="$(mktemp "$request_directory/.record-commit-hash.XXXXXX")" || {
    echo "FAIL: cannot create a temp file in $request_directory — is it writable? Nothing was written."; exit 1; }
  backup_file="$(mktemp "$request_directory/.record-commit-hash.orig.XXXXXX")" || {
    echo "FAIL: cannot create a backup file in $request_directory. Nothing was written."; exit 1; }
  cp -p "$request_file" "$backup_file" || { echo "FAIL: could not back up $request_file. Nothing was written."; exit 1; }
  cp -p "$request_file" "$temp_file"   || { echo "FAIL: could not seed the temp file's mode from $request_file. Nothing was written."; exit 1; }

  # awk -v processes backslash escapes in its value; a hex-validated hash contains none, and
  # the path is passed as an operand rather than through -v.
  awk -v new_hash="$commit_hash" '
    # The frontmatter is buffered and re-emitted; the body streams through verbatim and is
    # never inspected, which is what makes a body `commit:` structurally unreachable.
    NR == 1 && $0 == "---" { inside_frontmatter = 1; print; next }
    inside_frontmatter && $0 == "---" {
      if (!seen_commit) {
        # Insert point: right after completed_at:, the field the schema pairs with commit as
        # the terminal-flip stamp; otherwise the last line of the block.
        if (anchor_index == 0) anchor_index = buffered_count
        for (line_index = 1; line_index <= buffered_count; line_index++) {
          print buffered_line[line_index]
          if (line_index == anchor_index) print "commit: " new_hash
        }
        if (buffered_count == 0) print "commit: " new_hash
      } else {
        for (line_index = 1; line_index <= buffered_count; line_index++) print buffered_line[line_index]
      }
      inside_frontmatter = 0; frontmatter_closed = 1; print "---"; next
    }
    inside_frontmatter {
      buffered_count++
      if (index($0, "commit:") == 1) {
        if (seen_commit) duplicate_commit = 1
        buffered_line[buffered_count] = "commit: " new_hash
        seen_commit = 1
      } else {
        buffered_line[buffered_count] = $0
        if (index($0, "completed_at:") == 1) anchor_index = buffered_count
      }
      next
    }
    { print }
    END {
      # Without this flush an unterminated frontmatter block would make THIS SCRIPT the
      # truncation it exists to prevent. Flush the buffer, then fail.
      if (inside_frontmatter) {
        for (line_index = 1; line_index <= buffered_count; line_index++) print buffered_line[line_index]
        exit 4
      }
      if (!frontmatter_closed) exit 5
      if (duplicate_commit) exit 6
    }
  ' "$request_file" > "$temp_file"
  awk_status=$?
  if [ "$awk_status" -ne 0 ]; then
    case "$awk_status" in
      4) echo "FAIL: $request_file has an opening '---' with no closing '---' — the frontmatter block is unterminated." ;;
      5) echo "FAIL: $request_file has no frontmatter block (line 1 must be exactly '---')." ;;
      6) echo "FAIL: $request_file frontmatter has more than one 'commit:' line — ambiguous." ;;
      *) echo "FAIL: the frontmatter rewrite failed (awk exit $awk_status)." ;;
    esac
    echo "  $request_file is UNCHANGED on disk. Nothing was staged or committed."
    exit 1
  fi

  # --- Exact post-edit arithmetic. Any other changed byte fails this. ---
  new_commit_line="commit: $commit_hash"
  # Adding the missing trailing newline also rewrites the final line, so it costs one extra
  # insertion AND one extra deletion in every diff-shaped view — the byte/line arithmetic and
  # the numstat allowance below both have to carry that term or a legitimate file trips them.
  if [ "$edit_kind" = "replaced" ]; then
    expected_bytes=$(( pre_edit_bytes - ${#existing_commit_line} + ${#new_commit_line} + trailing_newline_added ))
    expected_lines=$(( pre_edit_lines + trailing_newline_added ))
    expected_changed_lines=$(( 2 + 2 * trailing_newline_added ))
    allowed_insertions=$(( 1 + trailing_newline_added ))
    allowed_deletions=$(( 1 + trailing_newline_added ))
  else
    expected_bytes=$(( pre_edit_bytes + ${#new_commit_line} + 1 + trailing_newline_added ))
    expected_lines=$(( pre_edit_lines + 1 + trailing_newline_added ))
    expected_changed_lines=$(( 1 + 2 * trailing_newline_added ))
    allowed_insertions=$(( 1 + trailing_newline_added ))
    allowed_deletions=$(( 0 + trailing_newline_added ))
  fi
  post_edit_bytes="$(wc -c < "$temp_file" | tr -d '[:space:]')"
  post_edit_lines="$(wc -l < "$temp_file" | tr -d '[:space:]')"
  # `grep -c` reads to EOF, so it cannot SIGPIPE the upstream diff the way `grep -q` would
  # under pipefail; `|| true` absorbs grep's no-match exit status.
  changed_line_count="$(diff "$request_file" "$temp_file" | grep -c '^[<>]' || true)"

  if [ ! -s "$temp_file" ] \
     || [ "$post_edit_bytes" -ne "$expected_bytes" ] \
     || [ "$post_edit_lines" -ne "$expected_lines" ] \
     || [ "$changed_line_count" -ne "$expected_changed_lines" ]; then
    echo "FAIL: refusing to write $request_file — the rewrite changed more than the single 'commit:' line."
    echo "  pre-edit:  $pre_edit_lines lines, $pre_edit_bytes bytes"
    echo "  post-edit: $post_edit_lines lines, $post_edit_bytes bytes  (expected $expected_lines lines, $expected_bytes bytes) — REJECTED and discarded"
    echo "  changed lines: $changed_line_count (expected $expected_changed_lines)"
    echo "  $request_file is UNCHANGED on disk. Nothing was staged or committed."
    echo "  Do NOT retry this command and do NOT hand-edit around it: a discrepancy this size means the file or this"
    echo "  script's assumptions are wrong, and committing anyway is how six archived REQs were truncated to 0 bytes."
    echo "  Inspect the file, recover it if needed (git show HEAD:$request_file), then run this once."
    exit 1
  fi

  # --- Re-parse the rewritten file: structure, and the fields that must not move ---
  post_edit_parse_ok=1
  [ "$(head -n 1 "$temp_file")" = "---" ] || post_edit_parse_ok=0
  [ "$(frontmatter_count_for "$temp_file" commit)" -eq 1 ] || post_edit_parse_ok=0
  [ "$(frontmatter_value_for "$temp_file" commit)" = "$commit_hash" ] || post_edit_parse_ok=0
  [ "$(frontmatter_value_for "$temp_file" id)" = "$request_id" ] || post_edit_parse_ok=0
  [ -n "$(frontmatter_value_for "$temp_file" status)" ] || post_edit_parse_ok=0
  [ "$(frontmatter_value_for "$temp_file" status)" = "$pre_edit_status" ] || post_edit_parse_ok=0
  if [ "$post_edit_parse_ok" -ne 1 ]; then
    echo "FAIL: the rewritten frontmatter does not re-parse as expected — id: and status: must be unchanged and commit: must be exactly '$commit_hash', once."
    echo "  $request_file is UNCHANGED on disk. Nothing was staged or committed."
    exit 1
  fi

  mv "$temp_file" "$request_file" || {
    echo "FAIL: the atomic rename into $request_file failed. $request_file is UNCHANGED. Nothing was staged or committed."; exit 1; }
  temp_file=""   # ownership transferred to $request_file

  # Post-rename numstat: the pending delta may only have grown by this one line. This is the
  # "1 insertion / 1 deletion, never 0/126" guard, stated against the recorded baseline so it
  # stays correct when the baseline was not already clean.
  if [ "$path_tracked" -eq 1 ] && [ "$head_exists" -eq 1 ]; then
    post_numstat="$(git diff --numstat --no-renames HEAD -- "$request_file")"
    post_insertions="$(printf '%s\n' "$post_numstat" | awk 'NR == 1 { print $1 + 0 }')"
    post_deletions="$(printf '%s\n' "$post_numstat" | awk 'NR == 1 { print $2 + 0 }')"
    if [ "${post_deletions:-0}" -gt "$((pending_deletions + allowed_deletions))" ] || [ "${post_insertions:-0}" -gt "$((pending_insertions + allowed_insertions))" ]; then
      echo "FAIL: git diff --numstat HEAD -- $request_file reports +${post_insertions}/-${post_deletions}; this edit may only add $allowed_insertions line(s) and delete $allowed_deletions (baseline was +$pending_insertions/-$pending_deletions)."
      if cp -p "$backup_file" "$request_file"; then
        echo "  $request_file has been RESTORED to its pre-edit content. Nothing was staged or committed."
      else
        echo "  !! The restore FAILED. $request_file holds the rejected edit. Recover it with: git checkout HEAD -- $request_file"
      fi
      exit 1
    fi
  fi
fi

# ---------------------------------------------------------------------------
# Idempotency. "Already recorded" is NOT "already committed": a pre-commit hook that rejected
# an earlier run leaves exactly that state, and reporting NOOP there would strand the edit
# uncommitted forever.
# ---------------------------------------------------------------------------
if [ "$edit_kind" = "none" ]; then
  if [ "$path_tracked" -eq 1 ] && [ "$head_exists" -eq 1 ] && ! git diff --quiet HEAD -- "$request_file" 2>/dev/null; then
    echo "OK: $request_file already records commit: $commit_hash on disk, but that edit is not in HEAD — an earlier metadata commit did not land. Stage and commit it now."
    exit 0
  fi
  echo "NOOP: $request_file already records commit: $commit_hash. Nothing to do — make no metadata commit."
  exit 0
fi

edit_description="$edit_kind"
[ "$edit_kind" = "inserted" ] && edit_description="inserted after completed_at:"
[ "$edit_kind" = "replaced" ] && edit_description="replaced (was '${existing_commit_value:-<empty>}')"
echo "OK: commit: $commit_hash $edit_description in $request_file — $pre_edit_bytes → $post_edit_bytes bytes, $changed_line_count changed line(s); all content guards passed."
echo "  Now stage and commit exactly this file:"
# The path is single-quoted so these lines stay copy-pasteable when it contains a space.
echo "    git add -- '$request_file'"
echo "    git commit -m \"[$request_id] record commit hash $commit_hash\" -- '$request_file'"
exit 0

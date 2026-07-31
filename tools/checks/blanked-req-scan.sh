#!/usr/bin/env bash
# blanked-req-scan.sh — find REQ/UR files whose CONTENT was destroyed, and resolve where each
# one can be recovered from.
#
# WHY THIS EXISTS. In one consumer repo the Step 9 "record commit hash" metadata commit
# replaced six archived REQ files with nothing — 9 KB to 26 KB of decision trail each. The
# only symptom anyone saw, weeks later, was the Kanban board parking them under Needs input /
# Blocked as *untitled* with an `unrecognized status ""` warning, because an empty file has no
# `status:` frontmatter. That framing is worse than useless here: it invites someone to "fix"
# the status field, which cements the loss. This scan reports the real defect — the body is
# gone — and, because the window closes when git gc collects the unreferenced objects, it
# resolves the recovery commit at the same time.
#
# `tools/checks/record-commit-hash.sh` is the guard that stops this happening; this is the
# detector for damage that already happened.
#
# Usage: tools/checks/blanked-req-scan.sh [--porcelain | --restore [--dry-run]]
#   Run from the project root (the directory containing do-work/).
#   --porcelain  emit only the machine-readable BLANKED records, no human report.
#   --restore    repair each damaged file: write back the content from its recovery commit,
#                then re-apply the recorded implementation hash by calling
#                tools/checks/record-commit-hash.sh (never by editing frontmatter here —
#                one implementation of that edit, carrying its guards, is the whole point).
#                actions/cleanup.md Pass 6 is the consent gate; this script does not prompt.
#   --dry-run    with --restore: print exactly what would be restored and write nothing.
#
# Output. For each damaged file, a human-readable block, plus one tab-separated record:
#
#   BLANKED<TAB><path><TAB><recovery-source-sha><TAB><recorded-hash>
#
# `<recovery-source-sha>` is `-` when git holds no non-empty version of the file, and
# `<recorded-hash>` is `-` when no blanking commit message carried one. Callers consuming
# these records must handle both.
#
# Exit 0: nothing damaged (or no do-work/ directory to scan). Under --restore, every damaged
#         file was FULLY repaired — content back AND the recorded commit: hash re-applied.
# Exit 1: at least one damaged file was found. This is a finding, not an error — callers that
#         treat non-zero as a crash will misreport a successful scan. Under --restore it also
#         covers a PARTIAL repair: content restored but its recorded hash could not be
#         re-applied. That file is committable with wrong provenance, so it must not be
#         reported as repaired (actions/cleanup.md Pass 6 keys its report on this exit code).
# Exit 2: usage error.
#
# Read-only. This script never writes, moves, or commits anything, which is what lets
# actions/forensics.md call it while keeping its own read-only contract.
set -uo pipefail
export LC_ALL=C

porcelain_mode=0
restore_mode=0
dry_run_mode=0
while [ "$#" -gt 0 ]; do
  case "$1" in
    --porcelain) porcelain_mode=1 ;;
    --restore) restore_mode=1 ;;
    --dry-run) dry_run_mode=1 ;;
    *) echo "usage: $0 [--porcelain | --restore [--dry-run]]" >&2; exit 2 ;;
  esac
  shift
done
if [ "$restore_mode" -eq 1 ] && [ "$porcelain_mode" -eq 1 ]; then
  echo "usage: $0 — --porcelain and --restore are mutually exclusive" >&2; exit 2
fi
if [ "$dry_run_mode" -eq 1 ] && [ "$restore_mode" -eq 0 ]; then
  echo "usage: $0 — --dry-run only applies with --restore (the scan alone never writes)" >&2; exit 2
fi

# The restore re-applies the hash through the guarded write-back rather than editing
# frontmatter itself, so resolve it up front and refuse early if it is missing — a partial
# restore that leaves the commit: field wrong is worse than not starting.
record_commit_hash_script="$(dirname "$0")/record-commit-hash.sh"
if [ "$restore_mode" -eq 1 ] && [ ! -x "$record_commit_hash_script" ]; then
  echo "FAIL: $record_commit_hash_script is missing or not executable — the restore re-applies the commit: field through it and will not hand-edit frontmatter instead." >&2
  exit 2
fi

if [ ! -d "do-work" ]; then
  [ "$porcelain_mode" -eq 0 ] && echo "No do-work/ directory here — nothing to scan. Run from the project root."
  exit 0
fi

git_available=0
git rev-parse --git-dir >/dev/null 2>&1 && git_available=1

# A file is damaged when it is 0 bytes, or when it has no parseable frontmatter — a body that
# survived but lost its header is still a file the pipeline cannot read, and the remedy is the
# same (recover content, not edit a field).
has_parseable_frontmatter() {
  awk '
    NR == 1 && $0 != "---" { exit 1 }
    NR > 1 && $0 == "---" { found_close = 1; exit }
    NR > 200 { exit 1 }
    END { if (!found_close) exit 1 }
  ' "$1" 2>/dev/null
}

# Walk this file's history newest-first and return the newest commit whose blob is non-empty.
# --full-history because default history simplification can hide a commit that touched the
# file on a side branch — exactly where a bad merge-time edit would live.
resolve_recovery_source() {
  local file_path="$1" tracked_name candidate_sha blob_size
  tracked_name="$(git ls-files --full-name -- "$file_path" 2>/dev/null)"
  [ -z "$tracked_name" ] && return 1
  while IFS= read -r candidate_sha; do
    [ -z "$candidate_sha" ] && continue
    blob_size="$(git cat-file -s "$candidate_sha:$tracked_name" 2>/dev/null || echo 0)"
    if [ "${blob_size:-0}" -gt 0 ]; then
      printf '%s' "$candidate_sha"
      return 0
    fi
  done < <(git log --full-history --format=%H -- "$file_path" 2>/dev/null)
  return 1
}

# The commit that destroyed the file is the newest one whose blob is empty; its subject is
# where Step 9 recorded the implementation hash ("[REQ-NNN] record commit hash <hash>").
# Never `git show --name-only` here: it prints the commit header and message before the file
# list, so a message line can pass a filename grep and become a phantom path.
resolve_recorded_hash() {
  local file_path="$1" tracked_name candidate_sha blob_size subject_text extracted_hash
  tracked_name="$(git ls-files --full-name -- "$file_path" 2>/dev/null)"
  [ -z "$tracked_name" ] && return 1
  while IFS= read -r candidate_sha; do
    [ -z "$candidate_sha" ] && continue
    blob_size="$(git cat-file -s "$candidate_sha:$tracked_name" 2>/dev/null || echo 0)"
    if [ "${blob_size:-0}" -gt 0 ]; then
      return 1   # walked past the damage without finding a recorded hash
    fi
    subject_text="$(git log -1 --format=%s "$candidate_sha" 2>/dev/null)"
    extracted_hash="$(printf '%s' "$subject_text" | sed -n 's/.*record commit hash \([0-9a-f]\{7,40\}\).*/\1/p')"
    if [ -n "$extracted_hash" ]; then
      printf '%s' "$extracted_hash"
      return 0
    fi
  done < <(git log --full-history --format=%H -- "$file_path" 2>/dev/null)
  return 1
}

# Restore one damaged file from its recovery commit, then re-apply the recorded hash through
# the guarded write-back. Prints ONE outcome token on stdout — `restored`, `partial`, or
# `skipped` — so the caller can tally; everything human-readable goes to stderr, keeping the
# token clean. It has to be a token rather than a 1/0 count because "content is back but its
# commit: hash is not" is neither a success nor a no-op: the file is committable, and
# committing it records wrong provenance under a report that claims full repair. This runs in
# a command substitution, so a variable set here cannot reach the caller — the token is the
# only channel. Content is written to a temp file in the target's own directory and moved into
# place only after it is confirmed non-empty — the recovery must never itself blank the file.
restore_one_file() {
  local target_path="$1" recovery_sha="$2" recorded_hash="$3"
  local tracked_name recovered_bytes restore_temp write_back_output

  if [ "$recovery_sha" = "-" ]; then
    echo "SKIP: $target_path — no recoverable content in git history; nothing to restore." >&2
    printf 'skipped'
    return 0
  fi
  tracked_name="$(git ls-files --full-name -- "$target_path" 2>/dev/null)"
  if [ -z "$tracked_name" ]; then
    echo "SKIP: $target_path — not tracked by git, so its history cannot be read." >&2
    printf 'skipped'
    return 0
  fi
  recovered_bytes="$(git cat-file -s "$recovery_sha:$tracked_name" 2>/dev/null || echo 0)"

  if [ "$dry_run_mode" -eq 1 ]; then
    echo "WOULD RESTORE: $target_path — $recovered_bytes bytes from commit $recovery_sha, then set commit: ${recorded_hash}." >&2
    printf 'skipped'
    return 0
  fi

  restore_temp="$(mktemp "$(dirname "$target_path")/.blanked-restore.XXXXXX")" || {
    echo "FAIL: cannot create a temp file next to $target_path; nothing was written." >&2
    printf 'skipped'
    return 0
  }
  if ! git cat-file -p "$recovery_sha:$tracked_name" > "$restore_temp" 2>/dev/null; then
    rm -f "$restore_temp"
    echo "FAIL: could not read $recovery_sha:$tracked_name out of git; $target_path is unchanged." >&2
    printf 'skipped'
    return 0
  fi
  if [ ! -s "$restore_temp" ]; then
    rm -f "$restore_temp"
    echo "FAIL: the recovered content for $target_path is empty; refusing to write it. $target_path is unchanged." >&2
    printf 'skipped'
    return 0
  fi
  if ! mv "$restore_temp" "$target_path"; then
    rm -f "$restore_temp"
    echo "FAIL: could not move the recovered content into $target_path; it is unchanged." >&2
    printf 'skipped'
    return 0
  fi
  echo "RESTORED: $target_path — $recovered_bytes bytes from commit $recovery_sha." >&2

  if [ "$recorded_hash" = "-" ]; then
    # Not a failure: no hash was ever recoverable, so there is nothing this run could have
    # applied. The content — the irreplaceable part — is fully back, so it counts as restored;
    # the missing provenance is reported for the operator instead of failing the whole repair.
    echo "  NOTE: the commit that emptied it carried no 'record commit hash' message, so commit: was left as recovered. Identify the implementation hash by hand and apply it with tools/checks/record-commit-hash.sh." >&2
    printf 'restored'
    return 0
  fi
  # The write-back's own output is the diagnosis (which guard tripped, and what to do). Passing
  # it through indented, rather than swallowing it, is what makes "run it yourself" unnecessary.
  if write_back_output="$("$record_commit_hash_script" "$target_path" "$recorded_hash" 2>&1)"; then
    echo "  commit: set to $recorded_hash (via the guarded write-back)." >&2
    printf 'restored'
    return 0
  fi
  echo "FAIL: $target_path — content restored, but re-applying commit: $recorded_hash did NOT succeed." >&2
  printf '%s\n' "$write_back_output" | sed 's/^/    /' >&2
  echo "  This file now holds the right content with the WRONG commit: provenance. Fix it with tools/checks/record-commit-hash.sh before committing — do not hand-edit the frontmatter." >&2
  printf 'partial'
}

damaged_count=0
restored_count=0
partial_count=0

# `find` recurses, so this surfaces loose archive REQs and UR-nested ones in one pass; a
# top-level glob would silently miss every UR folder. UR-*.md files are included because the
# same write-back touches them and an unreadable UR is the same class of loss.
while IFS= read -r candidate_file; do
  [ -f "$candidate_file" ] || continue
  file_bytes="$(wc -c < "$candidate_file" | tr -d '[:space:]')"
  damage_kind=""
  if [ "${file_bytes:-0}" -eq 0 ]; then
    damage_kind="0 bytes — the body is gone"
  elif ! has_parseable_frontmatter "$candidate_file"; then
    damage_kind="$file_bytes bytes but no parseable frontmatter block"
  fi
  [ -z "$damage_kind" ] && continue

  damaged_count=$((damaged_count + 1))
  recovery_sha="-"
  recorded_hash="-"
  if [ "$git_available" -eq 1 ]; then
    recovery_sha="$(resolve_recovery_source "$candidate_file" || true)"
    [ -z "$recovery_sha" ] && recovery_sha="-"
    recorded_hash="$(resolve_recorded_hash "$candidate_file" || true)"
    [ -z "$recorded_hash" ] && recorded_hash="-"
  fi

  if [ "$porcelain_mode" -eq 0 ]; then
    echo "DATA LOSS: $candidate_file — $damage_kind."
    if [ "$recovery_sha" != "-" ]; then
      recovered_bytes="$(git cat-file -s "$recovery_sha:$(git ls-files --full-name -- "$candidate_file")" 2>/dev/null || echo '?')"
      echo "  Recoverable: $recovered_bytes bytes at commit $recovery_sha."
      if [ "$recorded_hash" != "-" ]; then
        echo "  The commit that emptied it recorded implementation hash $recorded_hash."
      else
        echo "  No 'record commit hash' message found on the commit that emptied it — the implementation hash must be identified by hand."
      fi
      echo "  Restore with: do-work cleanup   (Pass 6 restores it and re-applies the hash, after asking)"
    elif [ "$git_available" -eq 1 ]; then
      echo "  No recoverable content in git history — every recorded version of this file is empty."
      echo "  If it was never committed with content, the body cannot be recovered here; check backups or re-capture the work."
    else
      echo "  Not a git repository — no history to recover from. Check backups."
    fi
  fi
  if [ "$restore_mode" -eq 1 ]; then
    case "$(restore_one_file "$candidate_file" "$recovery_sha" "$recorded_hash")" in
      restored) restored_count=$((restored_count + 1)) ;;
      # Content is back, so it is no longer damaged — but it is not repaired either, and the
      # summary below turns any partial into a non-zero exit.
      partial)  partial_count=$((partial_count + 1)) ;;
      *) ;;
    esac
    continue
  fi

  printf 'BLANKED\t%s\t%s\t%s\n' "$candidate_file" "$recovery_sha" "$recorded_hash"
done < <(
  {
    find do-work/archive -type f \( -name 'REQ-*.md' -o -name 'UR-*.md' \) 2>/dev/null
    find do-work/queue do-work/working -type f \( -name 'REQ-*.md' -o -name 'UR-*.md' \) 2>/dev/null
  } | sort -u
)

if [ "$damaged_count" -eq 0 ]; then
  [ "$porcelain_mode" -eq 0 ] && echo "No blanked or unparseable REQ/UR files found."
  exit 0
fi

# A completed repair is a success, not a finding — exit 0 so cleanup's Pass 6 can tell "I
# fixed it" from "there is damage here". A dry run reports damage without fixing it, so it
# keeps the finding exit code. A PARTIAL repair is a finding too: its content is back but its
# commit: field is not what the blanking commit recorded, and exiting 0 there would hand
# cleanup a "fully repaired" report over a file that is about to be committed with wrong
# provenance — the one outcome that puts a fresh lie into the trail this script exists to save.
if [ "$restore_mode" -eq 1 ] && [ "$dry_run_mode" -eq 0 ]; then
  echo "Fully repaired $restored_count of $damaged_count damaged file(s). Re-run the scan to confirm, then commit the repaired paths."
  if [ "$partial_count" -gt 0 ]; then
    echo "$partial_count more had their content restored but NOT their recorded commit: hash — see the FAIL lines above. Apply the hash before committing those paths."
  fi
  unresolved_count=$((damaged_count - restored_count - partial_count))
  if [ "$unresolved_count" -gt 0 ]; then
    echo "$unresolved_count could not be restored at all — see the SKIP/FAIL lines above."
  fi
  [ "$restored_count" -eq "$damaged_count" ] || exit 1
  exit 0
fi

if [ "$porcelain_mode" -eq 0 ]; then
  echo
  echo "$damaged_count file(s) have lost their content. This is the signature of an unguarded"
  echo "commit-hash write-back (see tools/checks/record-commit-hash.sh). Recover them before the"
  echo "unreferenced git objects are collected."
fi
exit 1

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
# Usage: tools/checks/blanked-req-scan.sh [--porcelain]
#   Run from the project root (the directory containing do-work/).
#   --porcelain  emit only the machine-readable BLANKED records, no human report.
#
# Output. For each damaged file, a human-readable block, plus one tab-separated record:
#
#   BLANKED<TAB><path><TAB><recovery-source-sha><TAB><recorded-hash>
#
# `<recovery-source-sha>` is `-` when git holds no non-empty version of the file, and
# `<recorded-hash>` is `-` when no blanking commit message carried one. Callers consuming
# these records must handle both.
#
# Exit 0: nothing damaged (or no do-work/ directory to scan).
# Exit 1: at least one damaged file was found. This is a finding, not an error — callers that
#         treat non-zero as a crash will misreport a successful scan.
# Exit 2: usage error.
#
# Read-only. This script never writes, moves, or commits anything, which is what lets
# actions/forensics.md call it while keeping its own read-only contract.
set -uo pipefail
export LC_ALL=C

porcelain_mode=0
case "${1:-}" in
  --porcelain) porcelain_mode=1 ;;
  "") ;;
  *) echo "usage: $0 [--porcelain]" >&2; exit 2 ;;
esac

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

damaged_count=0

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

if [ "$porcelain_mode" -eq 0 ]; then
  echo
  echo "$damaged_count file(s) have lost their content. This is the signature of an unguarded"
  echo "commit-hash write-back (see tools/checks/record-commit-hash.sh). Recover them before the"
  echo "unreferenced git objects are collected."
fi
exit 1

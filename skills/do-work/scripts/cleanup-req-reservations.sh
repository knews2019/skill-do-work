#!/usr/bin/env bash
# Mechanically reap stale REQ reservation markers from do-work/.req-reservations/.
#
# `queue-kanban next-req` creates one durable marker per allocated number so
# concurrent captures never share an id. A marker only has to outlive the window
# between allocation and the capture landing; afterwards it is dead weight:
#
#   - Once a commit holds a REQ file claiming the number (under do-work/queue/,
#     do-work/working/, or anywhere in do-work/archive/), the file itself keeps
#     the number out of the allocator's reach (next-req counts REQ files and
#     markers alike), so the marker is redundant and is removed. Committed, not
#     merely present: a working-tree REQ file may be a capture still staging its
#     own marker, and deleting that marker mid-capture would break the capture's
#     prescribed `git add` of it. Outside a git work tree no capture flow can be
#     racing, so present-on-disk is enough there.
#   - A marker two days old or more with no such REQ file is an abandoned
#     capture and is removed. The number may then be reissued — that reuse is
#     the accepted trade for not accumulating markers forever.
#
# Younger unmatched markers are captures still in flight and are kept. Anything
# in the directory that is not a plain REQ-<digits> regular file is left alone,
# and a symlinked do-work/ or reservation directory is refused outright (the
# allocator refuses the same), so the hook can never delete through a link to
# somewhere outside the project. Deletions are plain `rm`: git shows them as
# ordinary deletions to stage and commit with the next housekeeping commit. The
# core SessionStart hook runs this every session start; it also runs standalone:
#
#   bash scripts/cleanup-req-reservations.sh [project-root]
#
# Fail-soft: no do-work/, no marker directory, or nothing stale exits 0 and
# prints nothing.
set -u

project_root="${1:-${CLAUDE_PROJECT_DIR:-.}}"
reservation_directory="$project_root/do-work/.req-reservations"
if [ -L "$project_root/do-work" ] || [ -L "$reservation_directory" ]; then
  exit 0
fi
[ -d "$reservation_directory" ] || exit 0

# Every number a REQ file's name claims, one per line. Leading zeros are
# stripped so REQ-000203 (marker width) matches REQ-203-slug.md.
claimed_number_pattern='s|.*/REQ-0*\([0-9][0-9]*\)[^0-9].*|\1|p'
if git -C "$project_root" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  claimed_numbers="$(git -C "$project_root" ls-tree -r --name-only HEAD -- \
    do-work/queue do-work/working do-work/archive 2>/dev/null \
    | grep '\.md$' | sed -n "$claimed_number_pattern")"
else
  claimed_numbers="$(find "$project_root/do-work/queue" "$project_root/do-work/working" \
    "$project_root/do-work/archive" -type f -name 'REQ-*.md' 2>/dev/null \
    | sed -n "$claimed_number_pattern")"
fi

removed_count=0
for marker_path in "$reservation_directory"/REQ-*; do
  { [ -f "$marker_path" ] && [ ! -L "$marker_path" ]; } || continue
  marker_digits="${marker_path##*/REQ-}"
  case "$marker_digits" in
    *[!0-9]* | '') continue ;;
  esac
  marker_number=$((10#$marker_digits))
  if printf '%s\n' "$claimed_numbers" | grep -qx -- "$marker_number"; then
    rm -f -- "$marker_path" && removed_count=$((removed_count + 1))
    continue
  fi
  # No landed REQ file: apply the two-day timeout to this marker alone.
  # -mtime +1 matches an age of 48 hours or more; a fresh clone resets mtimes,
  # which only delays a timeout deletion — it can never trigger one early.
  if [ -n "$(find "$marker_path" -mtime +1 -print 2>/dev/null)" ]; then
    rm -f -- "$marker_path" && removed_count=$((removed_count + 1))
  fi
done

if [ "$removed_count" -gt 0 ]; then
  printf 'do-work: removed %s stale REQ reservation marker(s) from do-work/.req-reservations/ — stage and commit the deletion(s).\n' "$removed_count"
fi
exit 0

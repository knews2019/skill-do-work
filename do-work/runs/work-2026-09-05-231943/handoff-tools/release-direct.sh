#!/usr/bin/env bash
# Release a direct fix that has no REQ behind it: bump VERSION and its mirrors and insert one
# changelog entry, through the canonical `release` publication command. The fix itself must
# already be in the working tree; this script only adds the release bookkeeping, and the
# caller commits both together. Mirrors finalize-req.sh's release block, minus the lifecycle.
#
# usage: release-direct.sh <new-version> <changelog-entry-file> [--dry-run]
set -uo pipefail

scratch="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="${DO_WORK_REPO:-$(cd "$scratch/../../../.." && pwd)}"
cli() { bash "$repo/skills/do-work/tools/do-work-cli.sh" --repo-root "$repo" "$@"; }

new_version="$1"; changelog_entry_file="$2"; dry_run="${3:-}"
work_dir="${TMPDIR:-/tmp}/release-direct-$new_version"
rm -rf "$work_dir"; mkdir -p "$work_dir"
cd "$repo" || exit 1

old_version="$(cat "$repo/VERSION")"; old_version="${old_version%$'\n'}"
entry_title="$(sed -n '1s/^## [0-9.]* — \(.*\) ([0-9-]*)$/\1/p' "$changelog_entry_file")"
[ -n "$entry_title" ] || { echo "FATAL: could not read the entry title from $changelog_entry_file" >&2; exit 1; }
[ "$(grep -c "^## $new_version " "$repo/CHANGELOG.md")" -eq 0 ] || { echo "FATAL: $new_version already has an entry" >&2; exit 1; }

cp "$repo/VERSION" "$work_dir/version-old"
printf '%s\n' "$new_version" > "$work_dir/version-new"
cp "$repo/skills/do-work/actions/version.md" "$work_dir/version-md-old"
sed "s/^\*\*Current version\*\*: ${old_version//./\\.}\$/**Current version**: $new_version/" \
  "$repo/skills/do-work/actions/version.md" > "$work_dir/version-md-new"
grep -q "^\*\*Current version\*\*: $new_version\$" "$work_dir/version-md-new" || \
  { echo "FATAL: version.md rewrite did not take" >&2; exit 1; }

cp "$repo/CHANGELOG.md" "$work_dir/changelog-old"
anchor="## $old_version"
awk -v anchor="$anchor" -v entry_file="$changelog_entry_file" '
  index($0, anchor) == 1 && !done { while ((getline line < entry_file) > 0) print line; print ""; done = 1 }
  { print }
' "$repo/CHANGELOG.md" > "$work_dir/changelog-new"
[ "$(grep -c "^## $new_version " "$work_dir/changelog-new")" -eq 1 ] || \
  { echo "FATAL: the new entry heading does not appear exactly once" >&2; exit 1; }

python3 - "$work_dir" "$old_version" "$new_version" "$anchor" "$entry_title" <<'PY'
import json, sys
work, old, new, anchor, title = sys.argv[1:6]
def target(path, oldf, newf):
    return {"path": path, "expected_payload": {"source_path": f"{work}/{oldf}"},
            "new_payload": {"source_path": f"{work}/{newf}"},
            "old_version": old, "new_version": new}
def changelog(path):
    return {"path": path, "expected_payload": {"source_path": f"{work}/changelog-old"},
            "new_payload": {"source_path": f"{work}/changelog-new"},
            "insertion_anchor": anchor, "entry_key": new, "entry_title": title}
json.dump({"operation": "release", "release": {
    "maintainer_release": True, "old_version": old, "new_version": new,
    "project_owned_targets": ["VERSION", "CHANGELOG.md"],
    "required_mirrors": ["skills/do-work/VERSION", "skills/do-work/CHANGELOG.md",
                          "skills/do-work/actions/version.md"],
    "targets": [target("VERSION", "version-old", "version-new"),
                target("skills/do-work/VERSION", "version-old", "version-new"),
                target("skills/do-work/actions/version.md", "version-md-old", "version-md-new")],
    "changelogs": [changelog("CHANGELOG.md"), changelog("skills/do-work/CHANGELOG.md")]}},
    open(f"{work}/release-manifest.json", "w"), indent=2)
PY

if [ "$dry_run" = "--dry-run" ]; then
  cli --format text release --manifest "$work_dir/release-manifest.json" --dry-run
else
  cli --format text release --manifest "$work_dir/release-manifest.json"
fi

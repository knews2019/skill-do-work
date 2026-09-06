#!/usr/bin/env bash
# Finalize one do-work request whose implementation is already merged.
#
# Mechanical by design: every judgment (the version number, the changelog prose, the commit
# message) is an argument, and everything the finalizer refuses on is computed here rather
# than typed. See FINALIZATION-RECIPE.md beside this file for why each step is in this order.
#
# usage: finalize-req.sh <REQ-id> <request-path> <merge-hash> <commit-message> [<new-version> <changelog-entry-file>]
set -uo pipefail

scratch="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# The tools directory lives at do-work/runs/<run>/handoff-tools/ inside the repository.
repo="${DO_WORK_REPO:-$(cd "$scratch/../../../.." && pwd)}"
cli() { bash "$repo/skills/do-work/tools/do-work-cli.sh" --repo-root "$repo" "$@"; }

request_id="$1"; request_path="$2"; merge_hash="$3"; commit_message="$4"
new_version="${5:-}"; changelog_entry_file="${6:-}"
writer_label="vm:$repo"
work_dir="${TMPDIR:-/tmp}/finalize-$request_id"
rm -rf "$work_dir"; mkdir -p "$work_dir"
cd "$repo" || exit 1

now="$(cli --format text now)"
[ -n "$now" ] || { echo "FATAL: could not read the canonical instant" >&2; exit 1; }

# The index must be empty and no stale journal may exist, or the finalizer refuses mid-transaction.
git diff --cached --quiet --exit-code || { echo "FATAL: the git index is not empty" >&2; exit 1; }
[ ! -e "$repo/.git/do-work-finalization/$request_id.json" ] || \
  { echo "FATAL: a finalization journal already exists for $request_id" >&2; exit 1; }

# 1. Enumerate the lifecycle targets without mutating anything.
cli --format json complete "$request_id" --request-path "$request_path" --dry-run \
  --writer "$writer_label" --at "$now" --terminal-status completed \
  --implementation-hash "$merge_hash" > "$work_dir/plan.json" 2>"$work_dir/plan.err"
plan_status=$?
[ "$plan_status" -eq 0 ] || { echo "FATAL: dry run exited $plan_status" >&2; cat "$work_dir/plan.err" >&2; exit 1; }

release_manifest=""
release_paths=()
if [ -n "$new_version" ]; then
  old_version="$(cat "$repo/VERSION")"; old_version="${old_version%$'\n'}"
  entry_title="$(sed -n '1s/^## [0-9.]* — \(.*\) ([0-9-]*)$/\1/p' "$changelog_entry_file")"
  [ -n "$entry_title" ] || { echo "FATAL: could not read the entry title from $changelog_entry_file" >&2; exit 1; }

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

  release_manifest="$work_dir/release-manifest.json"
  python3 - "$work_dir" "$old_version" "$new_version" "$anchor" "$entry_title" "$release_manifest" <<'PY'
import json, sys
work, old, new, anchor, title, out = sys.argv[1:7]
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
    open(out, "w"), indent=2)
PY
  release_paths=(VERSION CHANGELOG.md skills/do-work/VERSION skills/do-work/CHANGELOG.md skills/do-work/actions/version.md)
fi
# Lesson satellites and other already-edited files the transaction must carry (work.md Step 8
# substep 4: "include their exact paths in the manifest"). Space-separated, repo-relative.
read -r -a extra_paths <<<"${EXTRA_COMMIT_PATHS:-}"
release_paths+=("${extra_paths[@]}")

# 2. Digests LAST, after every edit to the request file.
request_digest="$(sha256sum "$request_path" | cut -d' ' -f1)"
checkpoint_digest="$(sha256sum do-work/CHECKPOINT.md | cut -d' ' -f1)"

# 3. Author the one manifest. commit_paths is both a minimum (the plan) and a ceiling.
python3 - "$work_dir" "$request_id" "$request_path" "$writer_label" "$now" \
         "$request_digest" "$checkpoint_digest" "$commit_message" "$merge_hash" \
         "$release_manifest" "${release_paths[@]}" <<'PY'
import json, sys
(work, rid, rpath, writer, now, rsha, csha, msg, merge, relman, *relpaths) = sys.argv[1:]
plan = json.load(open(f"{work}/plan.json"))
paths = sorted({c["path"] for c in plan.get("changes", [])} | set(relpaths))
manifest = {"request_id": rid, "request_path": rpath, "writer_label": writer,
            "transition": "complete", "terminal_status": "completed", "completed_at": now,
            "expected_request_sha256": rsha, "expected_checkpoint_sha256": csha,
            "commit_paths": paths, "commit_message": msg,
            "provenance_mode": "supplied_commit", "implementation_hash": merge}
if relman:
    manifest["release_manifest_path"] = relman
    manifest["release_at"] = now
json.dump(manifest, open(f"{work}/finalization-manifest.json", "w"), indent=2)
print("commit_paths:", *paths, sep="\n  ")
PY

# 4. The one call that performs every mutation.
cli --format text advance "$request_id" --request-path "$request_path" \
  --finalization-manifest "$work_dir/finalization-manifest.json"

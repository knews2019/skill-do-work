#!/usr/bin/env python3
"""Build release + finalization manifests.
usage: REQ_ID=REQ-NNN REQ_PATH=... UR_ID=UR-NNN UR_CLOSES=0|1 make-manifests.py <repo-root> <out-dir> <old> <new> <merge-hash-full> <completed-at> <entry-title> <archive-req-path> <extra-commit-path>..."""
import sys, os, json, hashlib
root, out, old, new, merge, completed_at, title, archive_req = sys.argv[1:9]
extra = sys.argv[9:]
def sha(p): return hashlib.sha256(open(p,'rb').read()).hexdigest()
def payload(name):
    p = os.path.join(out, name); return {"source_path": p, "sha256": sha(p)}
req_id = os.environ['REQ_ID']; req_path = os.environ['REQ_PATH']; ur = os.environ['UR_ID']
release = {
  "operation": "release",
  "commit_message": f"[{req_id}] release {new}: {title}",
  "release": {
    "maintainer_release": True,
    "old_version": old, "new_version": new,
    "project_owned_targets": ["VERSION", "CHANGELOG.md"],
    "required_mirrors": ["skills/do-work/VERSION", "skills/do-work/actions/version.md", "skills/do-work/CHANGELOG.md"],
    "targets": [
      {"path": "VERSION", "expected_payload": payload('VERSION.old'), "new_payload": payload('VERSION.new'), "old_version": old, "new_version": new},
      {"path": "skills/do-work/VERSION", "expected_payload": payload('VERSION.old'), "new_payload": payload('VERSION.new'), "old_version": old, "new_version": new},
      {"path": "skills/do-work/actions/version.md", "expected_payload": payload('version-md.old'), "new_payload": payload('version-md.new'), "old_version": old, "new_version": new},
    ],
    "changelogs": [
      {"path": p, "expected_payload": payload('CHANGELOG.old'), "new_payload": payload('CHANGELOG.new'),
       "insertion_anchor": "CHANGELOG-2026-04-07-up-to-v0.49.0.md)\n\n",
       "entry_key": new, "entry_title": title}
      for p in ("CHANGELOG.md", "skills/do-work/CHANGELOG.md")
    ],
  },
}
no_release = os.environ.get('NO_RELEASE') == '1'
rp = os.path.join(out, 'release-manifest.json')
if not no_release: json.dump(release, open(rp,'w'), indent=2)
commit_paths = sorted(set([
  req_path, archive_req,
  
  "do-work/CHECKPOINT.md", "do-work/calibration-log.tsv",
] + ([] if no_release else ["VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md"]) + extra + ([f"do-work/user-requests/{ur}/input.md", f"do-work/archive/{ur}/input.md"] if os.environ.get('UR_CLOSES')=='1' else [])))
writer = open(os.path.join(root,'do-work/CHECKPOINT.md')).read().split('writer: ')[1].splitlines()[0].strip()
fin = {
  "request_id": req_id, "request_path": req_path, "writer_label": writer,
  "transition": "complete", "terminal_status": "completed", "completed_at": completed_at,
  "expected_request_sha256": sha(os.path.join(root, req_path)),
  "expected_checkpoint_sha256": sha(os.path.join(root, 'do-work/CHECKPOINT.md')),
  "commit_paths": commit_paths,
  "commit_message": (f"[{req_id}] complete: {title}" if no_release else f"[{req_id}] complete: {title} ({new})"),
  "provenance_mode": "supplied_commit", "implementation_hash": merge,
}
if not no_release:
  fin["release_manifest_path"] = rp; fin["release_at"] = completed_at
fp = os.path.join(out, 'finalization-manifest.json'); json.dump(fin, open(fp,'w'), indent=2)
print(rp); print(fp); print('writer', writer); print('commit_paths', len(commit_paths))

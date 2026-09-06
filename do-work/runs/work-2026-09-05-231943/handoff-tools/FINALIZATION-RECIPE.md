# Finalization recipe (verified against the Go source)

writer_label for THIS checkout: `vm:/home/user/skill-do-work`
CLI = `bash skills/do-work/tools/do-work-cli.sh --repo-root /home/user/skill-do-work --format json`

## Order (supplied_commit, implementation already merged)

1. `$CLI advance <REQ>` must report phase `finalize`. Anything else and a `complete` manifest is rejected.
2. `fold-timing-summary` if this run recorded timing (edits the REQ file — must precede the digest).
3. Enumerate lifecycle targets WITHOUT mutating:
   `$CLI complete <REQ> --request-path <REQPATH> --dry-run --writer 'vm:/home/user/skill-do-work' --at <now> --terminal-status completed --implementation-hash <merge>`
   Every `changes[].path` is the required part of commit_paths.
4. `git diff --cached --quiet --exit-code` must exit 0. `.git/do-work-finalization/<REQ>.json` must not exist from a different attempt.
5. Author the release payload files + release manifest (below) if this is a release.
6. Compute digests LAST: `sha256sum <REQPATH> do-work/CHECKPOINT.md`
7. Write the finalization manifest.
8. `$CLI advance <REQ> --request-path <REQPATH> --finalization-manifest <path>`
9. Consume: outcome success AND exactly one `finalizations` record matching, phase `cleanup_complete`, empty blocked_paths and reason_codes.
10. Cleanup by operative name from the integration branch:
    `git worktree remove <path>` (never --force), `git branch -d <name>` (never -D), `git worktree prune`.

## Finalization manifest keys (DisallowUnknownFields — any extra key is a hard error)

request_id, request_path, writer_label, transition ("complete"|"fail"),
terminal_status ("completed"|"completed-with-issues", only when complete),
completed_at (RFC3339), expected_request_sha256, expected_checkpoint_sha256 (64 hex each),
commit_paths (non-empty), commit_message, provenance_mode ("supplied_commit"|"primary_commit"),
implementation_hash (required for supplied_commit, FORBIDDEN for primary_commit; must be an ancestor of HEAD),
release_manifest_path + release_at (both or neither),
failure_error + failure_type ("intent"|"spec"|"code"|"environment", only when transition=fail)

commit_paths is BOTH a minimum (must be a superset of the dry-run plan plus release postimages)
AND a ceiling (only declared paths get staged). Any declared path that is dirty WILL be committed.
The run directory is NOT part of the lifecycle plan — commit run artifacts separately, beforehand.

## Release manifest (separate JSON file, also DisallowUnknownFields)

This repo's owned mirror set is fixed at five paths; declare all five or the release refuses
with RELEASE-MIRROR-UNDECLARED:
  VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md,
  skills/do-work/actions/version.md
`skills/do-work-board/tools/queue-kanban/VERSION` is NOT owned — never touch it.

{"operation":"release","release":{
  "maintainer_release": true,
  "old_version":"X","new_version":"Y",
  "project_owned_targets":["VERSION","CHANGELOG.md"],
  "required_mirrors":["skills/do-work/VERSION","skills/do-work/CHANGELOG.md","skills/do-work/actions/version.md"],
  "targets":[{"path":..., "expected_payload":{"source_path":...}, "new_payload":{"source_path":...},
              "old_version":"X","new_version":"Y"}, ...],
  "changelogs":[{"path":..., "expected_payload":{...},"new_payload":{...},
                 "insertion_anchor":"## X","entry_key":"Y","entry_title":"<Title>"}, ...]}}

Refusals to avoid:
- expected_payload bytes must EQUAL the file on disk now — do NOT hand-edit VERSION or CHANGELOG first.
- both changelog targets must have byte-identical new bytes and byte-identical expected bytes.
- insertion_anchor must occur EXACTLY ONCE in the expected preimage.
- entry_key and entry_title must be absent from the old bytes and occur exactly once in the new.
- new_version must be strictly greater than old_version (3-part semver).

## Precedent commits
- `37b61ca [REQ-582] release: ...` — minimal supplied_commit release shape.
- `ae974c2 [REQ-590] complete: ...` — a UR-closing one, carrying the UR archive destination and input.md.
There is NO metadata commit on the supplied_commit path.

---
source_type: req_lesson
req_id: REQ-064
req_path: do-work/archive/UR-010/REQ-064-restore-blanked-reqs-in-cleanup.md
date: 2026-07-31
domain: general
module: tools/checks
tags: [checks, restore, blanked, archived, reqs]
---

# Lessons from REQ-064: Restore blanked archived REQs from git history in cleanup

## What the REQ was about

Give `tools/checks/blanked-req-scan.sh` a `--restore` mode that recovers each blanked REQ's content
from git history and re-applies its recorded `commit:` hash, and expose it as a consent-gated
`### Pass 6` in `actions/cleanup.md`. This automates the recovery that had to be done by hand for six
files in the consumer repo.

## Solution summary

`blanked-req-scan.sh` gained a repair mode that reads each damaged file's content out of the recovery commit it already resolves, writes it back through a temp file in the target's own directory plus an atomic rename, refuses to write recovered content that is itself empty, and then re-applies the recorded `commit:` hash by *calling* `record-commit-hash.sh` rather than editing frontmatter — so the one guarded implementation of that edit stays the only one. `--dry-run` reports the same plan and writes nothing. `actions/cleanup.md` exposes it as Pass 6, gated on the user's consent, report-only when non-interactive, and running last so it repairs files at the paths the earlier passes just moved them to.

## What worked

- Building `--restore` on top of what the REQ-063 detector already resolves — the recovery sha and the recorded hash were computed and printed; the repair mode just consumes them. No second history walk exists to drift from the first. Delegating the `commit:` edit to `record-commit-hash.sh` rather than re-implementing a one-line frontmatter write was the constraint that made the whole thing small.

## What didn't work

- The first instinct on "refuse to restore over a non-empty file unless the size floor says it is truncated" was to add a size check inside the restore. That would have been a second, independently-drifting definition of "damaged" sitting under the detector's — the refusal already falls out of the detector gate (D-02). Also: the pass-count sweep is bigger than it looks. "Six passes" appears once, but the *count* is restated nine times in `actions/cleanup.md` under different phrasings ("all 6 passes", "Passes 0–5", "two narrow exceptions") — grepping the literal string finds one of them.

## Worth knowing

- The restore fixture must use a **real** commit hash. `record-commit-hash.sh` refuses a hash git can't resolve, so a made-up one fails the write-back for an unrelated reason and the probe proves nothing about the restore.
- `restore_one_file` prints its count to stdout and everything human-readable to stderr, because the caller consumes it as `$(...)`. Adding a plain `echo` for the operator inside that function corrupts the count.
- Pass 6 deliberately runs last (D-03): the earlier passes move files, and restoring before the move would split one content edit into a delete-plus-add in the cleanup commit.

## Back-reference

See `do-work/archive/UR-010/REQ-064-restore-blanked-reqs-in-cleanup.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `069c943`.

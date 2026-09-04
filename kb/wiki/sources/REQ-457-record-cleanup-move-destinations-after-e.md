---
title: "Lessons from REQ-457: Record cleanup move destinations after exclusive creation"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-04/REQ-457-record-cleanup-move-destinations-after-e.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-457: Record cleanup move destinations after exclusive creation

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

Make transaction-created-path ownership identify the filesystem object created by this invocation, rather than trusting a pathname that can later resolve to another writer's object. Register cleanup move destinations only after exclusive creation, and keep create/replace/move rollback confined after parent swaps at every later mutation point.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modified)

## What worked

- Taking the ownership event to be the *successful exclusive create* rather than the plan collapsed three separately-reported instances (cleanup's race, publication's parent swap, the BKB scaffold's escaped writes) into one invariant with one fix each. The sweep framing in the REQ was right.
- Neutering a guard and watching the matching test go red is the only evidence that a test is a lock-in. It is what caught F2: two of the shipped mechanisms had no test at all, and every package stayed green with them removed.

## What didn't work

- `os.SameFile` alone is not object identity. Deleting `kb/raw/_inbox_queue.md` and recreating it in the same directory reuses the inode, so the captured RED passed the identity check and rollback deleted the replacement anyway. Content digest had to become part of the binding (D-01).
- Widening a guard without asking what *absence* means introduced F1: `ENOENT` collapsed into "not owned", so a created path that was simply gone reported a preserved replacement and raised `committed_state_risk` on a clean rollback. A boolean was the wrong return type; the tri-state was (D-07).
- A record-time snapshot alone would have disowned this transaction's own second write, because `atomicfile.ReplaceExisting` publishes by rename and changes the inode. Re-capturing on our own recorder calls is what separates our republication from a foreign swap — a foreign writer never routes through the recorder.

## Worth knowing

- Created *directories* already had the correct identity shape (`publishedDirectories` + `os.SameFile`), and so did `bkb_init.go`'s scaffold rollback. Only the created-*file* branch of the shared recorder trusted a pathname. When one branch of a guard looks right, check whether its sibling was ever written.
- `publishedTracked` is populated solely for paths already in `dirtyTrackedPaths`, which is why the created-path guard that reads it was dead for every fresh creation. A guard that only ever runs on a set you never join is indistinguishable from no guard.
- Directories still cannot carry a content digest, so scaffold directory rollback is inode-only and a rmdir+mkdir inode reuse remains theoretically reachable for a `recursive: true` entry (M8). D-01's reasoning stops at files.

## Back-reference

See `do-work/archive/REQ-457-record-cleanup-move-destinations-after-exclusive-creation.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `b877eb6`.

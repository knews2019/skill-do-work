---
id: REQ-457
title: '[impact-critical] Record cleanup move destinations after exclusive creation'
status: claimed
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-450, REQ-451, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456]
batch: accepted-validate-feedback-root-causes
sweep: true
sweep_key: transaction-created-path-rollback-identity
claimed_at: 2026-09-02T23:27:17Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-02T23:49:40Z
  basis:
    - Route B
    - 6-file write set
    - 3 subsystems involved
    - 7 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
---

# Make Rollback Ownership Follow the Created Filesystem Object

## What

Make transaction-created-path ownership identify the filesystem object created by this invocation, rather than trusting a pathname that can later resolve to another writer's object. Register cleanup move destinations only after exclusive creation, and keep create/replace/move rollback confined after parent swaps at every later mutation point.

This sweep now owns both premature cleanup-destination recording and REQ-413's post-record parent-swap rollback failure. Both share one invariant: rollback may remove only the same filesystem object this invocation created.

## Instances

- Cleanup records a move destination before exclusive creation, so a losing writer can delete the winner's file during rollback.
- Publication creates and records a repository path, then a later parent swap can redirect pathname-only rollback to an outside same-named file.
- `internal/knowledgecommands/bkb_init.go`: BKB scaffold rollback checks identity and then separately removes by pathname; Git subprocesses ignore the opened root, and incomplete writes can escape ownership recording. Final-boundary replacements can therefore be deleted or mutated despite the recorded identity. (found by REQ-416 / UR-081)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #16 — P1 — source:** `internal/cleanup/cleanup_apply.go:296-298`

> ````text
> [P1] Record move destinations only after creating them — [prj].claude/skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:296-298
> When two cleanup processes, or another writer, race to the same absent destination, this records the path before CreateExclusiveAt succeeds. The losing create returns EEXIST, but transaction rollback then treats the winner's file as its own creation and removes it.
> Register the destination only after this process publishes it, before later fallible source deletion. See prime-do-work-cli.md:19-20 (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L19-L20).
> ````

- **Evidence:** `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:293-299` calls `RecordCreated` before `moveWithoutOverwrite`; the exclusive destination create occurs at lines `359-362`. Rollback deletes every recorded created path at `internal/gittransaction/git_transaction.go:463-475`. In the race, writer A publishes, writer B receives `EEXIST`, and B's rollback deletes A's file.
- **Surface-cost result:** N/A — this is a direct ownership and ordering correction at the existing exclusive-create seam.

## Detailed Requirements

- Call the exclusive destination create before registering the destination as created by this transaction.
- Register the destination immediately after successful create and before deleting the source.
- If recorder registration fails, remove only the destination just created by this process and leave the source intact.
- If exclusive creation fails with `EEXIST`, never record or remove that destination during rollback.
- Preserve no-overwrite semantics and exact transaction rollback behavior for paths the process genuinely owns.
- Hold or revalidate rooted object identity for every created path through rollback after each later create, replace, or move mutation.
- Never follow a swapped parent outside the repository or remove a different writer's replacement object.

## Constraints

- Never weaken exclusive creation to a check-then-create sequence.
- Preserve the exact-path transaction contract in `prime-do-work-cli.md`.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm. The ownership event is successful exclusive creation; record exactly after that event.

## Red-Green Proof

**RED prompt/case:** Coordinate two cleanup writers around exclusive creation, then separately swap a created path's parent after it is recorded and fail each later mutation index while protecting same-named outside objects.
**Why RED now:** Cleanup can record ownership before creation, and shared rollback later resolves recorded paths by pathname; both let one transaction delete an object it did not create.
**GREEN when:** Losing-writer rollback preserves the winner, post-record parent swaps never delete outside or replacement objects, and recorder-failure cleanup removes only the object created by the same invocation while preserving its source.
**Validation:** User confirmed after validate-feedback accepted Finding #16.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #16, captured by UR-085.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and constraints are fully specified by the finding, but the ownership mechanism the fix must hook into — where created-path identity is captured and where rollback consumes it — had to be discovered across three packages before any edit could be scoped.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Ownership recording.** `internal/gittransaction/git_transaction.go` holds `MutationRecorder`. `RecordCreated` (line 229) validates the path against `creatablePaths`, adds it to `createdPaths`, then captures identities for *directories* only (`captureCreatedDirectoryIdentities`, line 244). No identity is captured for the created **file** itself.

**Rollback.** `rollbackFailure` (line 882) walks `createdPaths` deepest-first and calls `root.Remove(path)` (line 1002). The only identity guard on that branch is `recorder.publishedTracked[path]`, which `RecordTouched` fills **solely for paths already in `dirtyTrackedPaths`** — never for a freshly created path. So a created path is removed by pathname with no proof the object at that pathname is the one this invocation created. Created *directories* already have the correct shape (`publishedDirectories` + `os.SameFile`, line 1013), and `bkb_init.go`'s `rootedScaffoldWriter.rollback` has it too — the file branch of the shared recorder is the one that does not.

**Record/create ordering across all thirteen `RecordCreated` call sites.** Eleven create first and record second (publication, requeststate, toolbox report-image/note/architecture/portfolio). Exactly two record *before* creating:
- `internal/cleanup/cleanup_apply.go:291` — records `operation.DestinationPath`, then `moveWithoutOverwrite` performs the exclusive create at line 359. This is the finding's race: writer B's `CreateExclusiveAt` returns `EEXIST` after writer A published, and B's rollback then removes A's file.
- `internal/knowledgecommands/interview_commands.go:1182` — records, then `createRootedFile`. Same class; there is even a `rootedCreateTestHook` between the two, which is the seam an adversarial test uses.

**Escaped ownership in the BKB scaffold.** `rootedScaffoldWriter.createFile` (line 423) opens the file `O_CREATE|O_EXCL`, then writes, closes, and `Lstat`s — appending to `writer.created` only after all of that succeeds. A failed write, close, or `Lstat` leaves a file on disk that rollback never sees. `createDirectory` (line 400) has the same shape around `Mkdir`.

**Why `os.SameFile` alone is the right guard, and where it is not enough.** `atomicfile.ReplaceExisting` publishes by rename, so any legitimate later replacement of a created path by this same transaction changes the inode. A pure record-time snapshot would therefore disown our own second write. The recorder is told about every one of its own later mutations (`RecordTouched`/`RecordCreated` are called for each), so re-capturing on those calls — and only on those calls — distinguishes our own re-publication from a foreign writer's swap, which never routes through the recorder.

*Generated in-session (no separate explore agent — single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — capture and revalidate created-path identity; gate created-path rollback removal on it
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modify) — record the move destination only after exclusive creation succeeds, before source deletion
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modify) — record after `createRootedFile`, not before
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modify) — record scaffold ownership immediately after exclusive creation so incomplete writes cannot escape it
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — RED/GREEN for post-record parent swap and pathname-only removal
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modify) — RED/GREEN for the losing-writer `EEXIST` race and recorder-failure cleanup

**Files I will NOT touch:** `internal/atomicfile/atomic_file.go` (its exclusive-create contract is what the fix depends on and must not be weakened), `internal/publication/publication_commands.go` (already creates before recording; it is fixed by the shared recorder change alone), `prime-do-work-cli.md` (exact-path transaction contract is preserved, not changed).

**Acceptance criteria (restated from REQ):**
- [ ] The exclusive destination create runs before the destination is registered as created by this transaction
- [ ] Registration happens immediately after a successful create and before the source is deleted
- [ ] A failed registration removes only the destination this process just created and leaves the source intact
- [ ] An `EEXIST` exclusive create never records the destination and never removes it during rollback
- [ ] No-overwrite semantics and exact-path rollback are preserved for paths the process genuinely owns
- [ ] Rooted object identity is held or revalidated for every created path through rollback, after each later create, replace, or move mutation
- [ ] Rollback never follows a swapped parent outside the repository and never removes another writer's replacement object

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ⚠ RED — `bash _dev/tests/maintainer-verify.sh` exits 1 at `008f3d3` with two pre-existing failures, recorded in `do-work/working/baseline-failures.txt`:
- `internal/knowledgecommands` → `TestBKBInitRollbackPreservesReplacedObjectsAndParents` — **this REQ's own captured RED** (the `bkb_init.go` instance in `## Instances`). It must go GREEN here.
- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback` ("media commit hook survived cancellation") — unrelated to this REQ's ownership invariant; excluded from this REQ's gate attribution and routed to its own REQ.

**Dependencies:** ⚠ this checkout shipped none of the gate's required toolchain — Go 1.26.1, ShellCheck 0.11.0, and `just` (board template needs ≥ 1.4x) were all installed before the baseline could run. Not a repository change.

*Checked by work action*

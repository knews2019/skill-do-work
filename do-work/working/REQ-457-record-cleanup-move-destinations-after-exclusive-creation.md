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
- [x] **[PLAN]:** Read `prime-do-work-cli.md` and its `lessons-do-work-cli.md` satellite (the `[family: final-boundary-identity]` trap covers this change), plus `general.md`, `coding-guardrails.md`, `communication-style.md`, `testing.md`. Approach: bind each created path to the exact filesystem object published there (inode + sha256), re-capture on this transaction's own later mutations of that path, revalidate the others at each mutation point, and gate rollback removal on that binding. Separately, move every remaining record-before-create call site behind its successful exclusive create.
- [x] **[APPLY]:** Six files, all inside the declared `## Scope` list; no file outside it was touched.
- [x] **[UNIFY]:** `git diff --stat` → 6 files, +309/-28. Orchestrator re-ran `go build ./...`, `go vet ./...` (clean) and `gofmt -l .` (no output) independently of the builder's report, and read the full diff of `git_transaction.go`, `cleanup_apply.go` and `interview_commands.go` for debug artifacts — none present; every added comment states an invariant rather than narrating a step.

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
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modify) — record after the rooted create succeeds, not before
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modify) — record scaffold ownership immediately after exclusive creation so incomplete writes cannot escape it
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — RED/GREEN for post-record parent swap and pathname-only removal
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modify) — RED/GREEN for the losing-writer exclusive-create race and recorder-failure cleanup

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

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modified)

**What was done:** `MutationRecorder` now binds every created path to the exact object this invocation published there — `os.FileInfo` plus a sha256 digest, in a new `createdObjects` map — captures that binding at `RecordCreated`, re-captures it when this same transaction mutates a path it already created, revalidates every other created path at each mutation point, and lets `rollbackFailure` remove a created path only while the binding still holds. Cleanup's `moveWithoutOverwrite` now takes a `registerDestination` callback it runs after the exclusive create and before source deletion, so a losing writer that receives `EEXIST` records nothing and cannot roll back over the winner's file; a failed registration removes only the object this process just created and leaves the source intact. The two remaining record-before-create call sites (`interview_commands.go`, and the BKB scaffold writer's `createFile`/`createDirectory`) now record ownership immediately after their exclusive creation succeeds, which also closes the gap where an incomplete write escaped ownership recording entirely.

## Decisions

- **D-01** — Content digest is part of created-object identity, not `os.SameFile` alone. **DECIDE & STATE.** The captured RED proved the need: deleting `kb/raw/_inbox_queue.md` and immediately recreating it in the same directory reuses the inode, so `os.SameFile` returned true and rollback deleted the replacement. Both the shared recorder and the BKB scaffold now bind inode plus sha256.
- **D-02** — A created path with no recorded identity is preserved only when an object stands there; an absent path is skipped silently. **DECIDE & STATE.** Reporting an error for a path that does not exist would turn every legitimate "the create never happened" rollback into `rollback_incomplete` for no safety gain. Nothing is removed by pathname either way.
- **D-03** — Two existing `gittransaction` fixtures were reordered rather than exempted. **DECIDE & STATE.** `TestPreCommitFailureRestoresTrackedAndRemovesOnlyCreatedTargets` and `TestRollbackRemovesOnlyRecordedCreatedDirectoriesDeepestFirst` called `RecordCreated` before writing the file, which the new contract treats as unowned. Both now create then record, matching all thirteen production call sites; their assertions are unchanged.
- **D-04** — Revalidating another created path's identity fails the mutation instead of silently disowning it. **DECIDE & STATE.** The error routes into `rollbackFailure`, which preserves the swapped object and reports it while rolling the rest back — detection at the next mutation point, as the REQ requires, reaching the same end state a rollback-time-only check would.
- **D-05** — A registered destination is not removed by `moveWithoutOverwrite` when source deletion fails. **ESCALATE.** That path previously removed the destination directly; the transaction now owns it and rollback removes it only after proving identity. **Value:** the removal goes through the same identity proof as every other created path, so a swapped destination is preserved rather than deleted. **Risk:** a future caller that registers and then suppresses rollback would leak the destination; the pre-registration direct removal is still guarded by `destinationRegistered` for the nil-callback case. Reversible in one line.
- **D-06** — `RecordTouched` now opens the rooted repository handle on every call, not only for dirty tracked paths. **DECIDE & STATE.** One extra `openat` per recorder call, which is what makes re-capture and revalidation rooted.

## Discovered Tasks

- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/install-suite/HUP` is flaky under full-suite parallel load: it failed once during a `go test ./...` run, then passed 3/3 in isolation and in a repeat full run, and did not fail in the pre-change baseline. A signal-timing race in the install confirmation path, unrelated to this REQ.
- `moveWithoutOverwrite`'s two pre-registration failure paths (destination directory changed, source changed before deletion) still remove the destination by name through the destination handle with no identity check. Nothing is registered yet at those points so no transaction ownership is at stake, but it is the same class of pathname-trusting removal this sweep targets.
- `skills/do-work/tools/checks/scope-drift.sh` reads every backticked token inside a `## Scope` "Files I will touch" bullet as a declared path, so an identifier named in a bullet's trailing description (`EEXIST`, `createRootedFile`) is reported as declared-but-never-touched. Worked around here by keeping identifiers out of those descriptions; the script should take only the leading path of each bullet.

## Qualification

**Passed** — 6 files verified, 7 requirements traced, P-A-U confirmed.

Mechanical: `tools/checks/qualify.sh` → `OK: mechanical qualification passed`. `tools/checks/scope-drift.sh` → `OK: Implementation Summary matches the Scope declaration` (exit 0); the six touched files are exactly the six declared, with nothing outside them.

Independent (orchestrator-run, not the builder's report):
- `go build ./... && go vet ./...` clean; `gofmt -l .` printed nothing.
- Read the full diff of `git_transaction.go`, `cleanup_apply.go` and `interview_commands.go`. No debug artifacts, no stubs, no placeholder returns. The new `createdObjects` map is consumed in three places (capture, revalidate, rollback gate), so the data path is live rather than recorded-and-ignored.
- Requirement trace: exclusive-create-before-registration and registration-before-source-deletion are both visible in `moveWithoutOverwrite`'s new `registerDestination` block; the registration-failure branch removes only `destinationName` through the already-open rooted destination handle and returns before `os.Remove(sourcePath)`; an `EEXIST` return happens above that block, so nothing is recorded; `CreateExclusiveAt` is unchanged, so no-overwrite semantics are intact; identity is captured in `RecordCreated`, re-captured in `RecordTouched` for an already-created path, and revalidated for the others on every recorder call; every removal goes through the rooted `*os.Root`, so a swapped parent cannot redirect it outside the repository.

## Testing

**Tests run:** `go build ./... && go vet ./... && gofmt -l .`; `go test ./...` (whole `do-work-cli` module); `go test -race ./internal/gittransaction ./internal/knowledgecommands ./internal/cleanup`; canonical repository gate `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Module green except two failures, neither attributable to this diff (attribution below). Race detector clean on all three touched packages.

**Red-green validation:**
- `TestBKBInitRollbackPreservesReplacedObjectsAndParents` (`internal/knowledgecommands`, subtests `git/object` + `standalone/object`) — the REQ's captured RED for the `bkb_init.go` instance: ✗ `bkb_init_test.go:304`, rollback action `removed owned kb/raw/_inbox_queue.md` → ✓ all four subtests pass. The test file was **not** modified.
- `TestLosingMoveWriterRollbackPreservesTheWinnersDestination` (`internal/cleanup`, new) — traces the Red-Green Proof's "coordinate two cleanup writers around exclusive creation": ✗ `losing rollback destroyed the winner's file: "" open …/do-work/archive/REQ-601-done.md: no such file or directory` → ✓.
- `TestFailedDestinationRegistrationRemovesOnlyTheCreatedDestination` (`internal/cleanup`, new) — traces "recorder-failure cleanup removes only the object created by the same invocation while preserving its source": ✗ `destination was registered before it was created` → ✓.
- `TestCreatedTargetRollbackPreservesAnotherWritersReplacement` (`internal/gittransaction`, new) — traces "separately swap a created path's parent after it is recorded": ✗ rollback reported `removed created target created.txt` → ✓ rollback reports `created target changed after publication; preserved replacement: created.txt` and the file still reads `second writer`.

**New tests added:**
- `internal/cleanup/cleanup_apply_test.go`: `TestLosingMoveWriterRollbackPreservesTheWinnersDestination`, `TestFailedDestinationRegistrationRemovesOnlyTheCreatedDestination`
- `internal/gittransaction/git_transaction_test.go`: `TestCreatedTargetRollbackPreservesAnotherWritersReplacement`

**Existing tests updated (cross-REQ impact):**
- `internal/gittransaction/git_transaction_test.go`: `TestPreCommitFailureRestoresTrackedAndRemovesOnlyCreatedTargets` and `TestRollbackRemovesOnlyRecordedCreatedDirectoriesDeepestFirst` recorded the created path before writing the file. Under the new contract that ordering is unowned, so both fixtures now create then record — matching all thirteen production call sites. Their assertions are unchanged; only the setup order moved (D-03).
- `internal/cleanup/cleanup_apply_test.go`: two direct `moveWithoutOverwrite` calls updated for the new `registerDestination` parameter.

**Canonical repository gate — attribution.** `bash _dev/tests/maintainer-verify.sh` exits 1 both before and after this change. Neither remaining failure belongs to this REQ, and the one the REQ *did* own is now green:
- `internal/knowledgecommands` → **fixed by this REQ** (was in the recorded red baseline).
- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`: in the recorded red baseline at `008f3d3`, identical failure text, and it exercises `RecordTouched` on an existing tracked file so it never reaches created-path rollback. Tracked as **REQ-524**.
- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/install-suite/INT`: not in the recorded baseline, but in a package this diff does not touch. Re-run once per the flake rule: `-count=5` in isolation passes 5/5, and the captured stderr shows the installer still rendering its managed-install diff when the signal arrived — a test-side synchronization race that surfaces under full-module parallel load. Tracked as **REQ-525**.

*Verified by work action*

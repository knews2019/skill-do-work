---
id: REQ-598
status: claimed
domain: backend
created_at: 2026-09-06T06:25:19Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
route: C
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
estimate:
  p50_active_minutes: 60
  confidence: low
  calculated_at: 2026-09-06T07:40:57Z
  basis:
    - Route C
    - 3-file write set
    - 1 subsystem, a transaction boundary
    - 4 acceptance criteria
    - no no-handle test exists in the package
maintenance: false
depends_on: [REQ-558, REQ-602]
related: [REQ-558]
write_set: [skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go, skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go, _dev/tests/audit-lockins.sh]
title: 'Close the live nil-handle panic in transaction rollback, and decide the handle once instead of eleven times'
claimed_at: 2026-09-06T07:40:57Z
---

# Close the Live Nil-Handle Panic in Transaction Rollback

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

`rollbackFailure` opens a rooted filesystem handle with `os.OpenRoot`. On failure it records the
error and **deliberately keeps going**, so a failed transaction still unstages its paths and still
returns a typed incomplete rollback. That nil handle then fans out to eleven consumers across four
loops. Eight of them test it. **One does not**: `rollbackFailure` calls
`quarantineAndRollbackPrivate(root, state, published)` with no check, and that function dereferences
the handle with `root.Mkdir`. Any transaction with an identity-recorded private untracked target
panics mid-rollback when that open fails.

## Why

A panic during rollback is the worst available failure: the transaction has already mutated the tree and
the process dies before restoring it. The tool's own report of an incomplete rollback — the thing the
keep-going decision at the open exists to produce — is lost.

Beyond the one missing guard, the eight that exist are a per-consumer answer to a question that should
be decided once. `privateStateStillOriginal` and `rootedCreateRegular` take the same handle and never
test it; they are safe only because guards at three call sites stand in front of them, so a future edit
to any of those sites breaks them silently.

## Context

Found by the three-agent reachability trace run for REQ-558, which asked for eight of the nine guards to
be deleted. The trace established the opposite — the guards are load-bearing — and two of the three
agents independently reproduced this panic. REQ-558 deleted the one genuinely unreachable guard and
pinned the rest at eight. The full trace is in
`do-work/runs/work-2026-09-05-231943/REQ-558-handback.md` and the workflow's per-guard reports.

## Detailed Requirements

- Close the panic. The minimum is a guard at the `quarantineAndRollbackPrivate` call site matching its
  eight siblings; the better answer is the one below.
- **Decide the handle once.** At the open site, either abandon the rooted half of rollback with a typed
  finding, or run it under a proven handle. Either makes all eight surviving guards genuinely dead, and
  subsumes the missing one. The keep-going behaviour that lets a failed transaction still unstage its
  paths must be preserved — that is why the open does not return today.
- **Write the package's first no-handle rollback test.** No test in the package reaches any no-handle
  branch today, which is why nine guards could sit there with one of them missing and nothing noticed.
  The test is the point of this request as much as the fix is.
- If the guards become dead, delete them and say so; if they stay, say why. Either way update the count
  the REQ-558 lock-in pins, in the same change.

## Constraints

- One source file, its test file, and the lock-in count in `_dev/tests/audit-lockins.sh`, which
  REQ-558 pinned at 8 as a floor and a ceiling. **The pin blocks this request by design**: the minimum
  fix adds a ninth guard and the preferred fix removes all eight, and either fails the fast tier until
  the pin moves in the same change. REQ-558's review reproduced that — its stated minimum fix, applied
  verbatim, fails `audit-lockins.sh` with `9 nil-root guards`. Move the pin to the new count in this
  request, and say in the block's comment why it moved.
- This is a transaction boundary. Every change is behaviour-preserving on the paths that work today, and
  the package's existing transaction and rollback tests must stay green unchanged.
- Do not remove the keep-going behaviour at the open site without replacing what it delivers.

## Red-Green Proof

**RED case:** a transaction with an identity-recorded private untracked target, rolled back with the
rooted open forced to fail. Today: panic at `root.Mkdir`.
**GREEN when:** the same case returns a typed incomplete rollback naming the unusable handle, the
package's transaction and rollback tests are green unchanged, and a new test drives the no-handle branch.

## Open Questions

None.

## Triage

**Route: C** — Explore, plan, then build.

**Reasoning:** The defect is settled — two of REQ-558's three trace agents reproduced it end to end
through the public `ExecuteTransaction` API, and REQ-558's review reproduced it again. What is not
settled is the shape of the fix, and the two shapes the request names are not equivalent. A ninth guard
at the one unguarded call site is a four-line change that keeps the per-consumer pattern and its eight
siblings. Deciding the handle once at the open — abandoning the rooted half of rollback with a typed
finding, or running it under a proven handle — makes all eight guards dead and is a change to
`rollbackFailure`'s error handling inside a transaction boundary, where the keep-going behaviour that
lets a failed transaction still unstage its paths must survive. That choice deserves competing plans and
a judge, not a builder's preference.

**Planning:** Required. Two independent plans, one per shape, judged on what each preserves and what
each costs, with the package's first no-handle test as a fixed requirement of both.

**REQ-558's traces are this request's exploration.** They enumerated every consumer of the handle,
every path into every guard, and the exact condition that reaches the dirty-tracked branch. They live in
the run directory beside REQ-558's hand-back and are cited rather than re-derived.

## Plan

**Two competing plans, judged on what each preserves; Plan B chosen with two amendments.** Full record
in the run directory as the judged-plan workflow output (`scratchpad/req598-verdict.json` holds the
verdict; `scratchpad/req598-final.patch`, 557 lines, is the verified tree and applies cleanly on the
current head; `scratchpad/final_test_func.go` is the exact test).

**Plan A, the minimum:** one nil check at the entry of `quarantineAndRollbackPrivate`, in the shape of
its sibling `rollbackDirtyTracked`; the eight guards stay; the pin moves 8 to 9; a test that renames the
worktree root away inside the mutation callback so the rollback open fails with `ENOENT`. Prototyped and
green. What it costs: nine per-consumer decisions remain with nothing in the code to remind a tenth; the
no-handle path keeps messages that describe observations nobody made ("preserved replacement" when the
handle was simply missing); and its test cannot show the keep-going half at all, because renaming the
root kills every `git -C` along with the open.

**Plan B, decide once at the open:** `rollbackFailure` opens through a package-level
`openRollbackRoot` (`= os.OpenRoot`, the house pattern of `privateTransactionTestHook`). On failure it
records the open error and runs `rollbackWithoutRoot`: every Git-side unstage and restore-from-HEAD plus
the pathname restore of existing-untracked opt-ins, and one `<kind> left in place; rollback root is
unavailable: <path>` error per target the rooted half would have touched. On success it defers `Close`
and runs `rollbackWithRoot`: today's three loops moved verbatim with the nil tests removed. Loop 4
(rolled-back paths, index check, result assembly) stays in `rollbackFailure`. All eight guards are
deleted; the missing ninth is subsumed because `quarantineAndRollbackPrivate` is reachable only from
`rollbackWithRoot`. Four helpers move out so both halves share them. The lock-in's Finding 3 pins the
guard count at **zero**, rewritten rather than renumbered.

**The test:** `TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest`, driven through
`ExecuteTransaction`, with the seam returning `&os.PathError{Op: "open", Err: syscall.EACCES}` (the shape
a real execute-only root produces for a non-root process, verified as uid 65534; the amendment fixes the
prototype's `openat`). It declares one target of every kind that exists today (tracked, dirty-tracked
opt-in, existing-untracked opt-in, identity-recorded private, created file, created directory) and
asserts the exact error list, the exact action list, an empty staged set, and the bytes of every file, so
it is the whole contract of the no-handle walk. RED at head plus seam: panic at `root.Mkdir` via
`rollbackFailure` and `ExecuteTransaction`. GREEN on the final tree; passes under `-race` and as uid 65534.

**Why B over A:** both close the panic and both RED/GREEN reproduce. B was chosen because it could be
shown behaviour-preserving: the existing suite is green with zero edits, a panic canary at the top of
`rollbackWithoutRoot` never fires across the existing suite and does fire from the new test, the `-U0`
diff removes exactly the eight guard lines and adds no guard, and a six-kind differential gives identical
outcomes at head-plus-seam and on B for every kind that does not panic. A's test proves only the
reporting half; the request says the test is the point as much as the fix.

**Ordered steps for the builder** (from the verdict, to be followed with its own RED evidence): record
the clean-head gate pair first; add the seam and the test, run it red on the seam-only tree (no stash);
then the restructure; then the lock-in rewrite; then guards, canary, `-U0` diff, differential, suite,
`-race`, `GOOS=windows go vet`, gate.

*Generated by two Plan agents and one judge; the judge verified both prototypes in fresh clones*

## Exploration

**Covered by the judged plan.** Both planners and the judge read `rollbackFailure`'s four loops line by
line against the head and agreed which work needs the handle: branch A (dirty-tracked) needs it for
`privateStateStillOriginal` and `rollbackDirtyTracked` but not for `git restore --staged`; branch B
(private) needs it for everything; branch C (existing-untracked) needs it nowhere; branch D (tracked)
needs it only for `trackedPublicationStillOwned`; loop 2 needs it for the ownership check, `Lstat`,
`inspectCreatedObject` and `Remove` but not for `git rm --cached`; loop 3 needs it for everything; loop 4
needs it nowhere. The opens at lines 138, 221, 324, 338, 370, 467 and 917 all return on failure, so
`rollbackFailure`'s open at 994 is the only keep-going open in the file. Three forcing mechanisms were
measured on this machine: `chmod 0111` (real, git alive, but uid 0 opens it and this suite runs as uid 0,
so it skips exactly here); rename the root away (real, privilege-free, but every `git -C` dies with it);
the seam (deterministic, git alive, mirrors the real `PathError`).

*Generated by the judged-plan workflow*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — the open decided once; two halves; four shared helpers; eight guards and one stale comment deleted
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — the seam import and the six-kind no-handle test
- `_dev/tests/audit-lockins.sh` (modify) — Finding 3 pins zero nil-root guards

**Files I will NOT touch:** anything outside the package and the lock-in. The changelog and version
mirrors belong to finalization (REQ-558 shipped two versions unnamed when that was missed; the record
says so here so it is not missed again).

**The challenge, recorded rather than asked (the maintainer is not watching):** this is +155/-81 where
a seven-line guard would close the panic, and the maintainer's rule is that machinery is not an
achievement. The simplest shape of all, passing `ExecuteTransaction`'s already-open handle into
`rollbackFailure` and deleting the reopen, would remove the open, the seam and the whole no-handle half;
it was excluded because this request's fixed requirement is a no-handle rollback test, which that shape
makes unwritable, and because it changes which directory object rollback acts on if the root is swapped
mid-transaction. Resolution: B, because it deletes more than it adds in decisions (eight guards and a
stale enumeration go; one decision point and one test arrive) and because it was the only plan whose
behaviour preservation could be demonstrated rather than argued. The simplest shape is named for the
maintainer in Decisions.

**Acceptance criteria:**
- [ ] A rollback whose root open fails no longer panics: the no-handle test is red at head plus seam and
  green after, driven through `ExecuteTransaction`
- [ ] The no-handle walk still unstages every staged path and restores every tracked target from HEAD,
  and reports each rooted target it could not touch by kind and path
- [ ] Every path that worked before is byte-identical in behaviour: existing suite green with no edits,
  canary never fires from it, `-U0` diff removes only the eight guard lines
- [ ] `_dev/tests/audit-lockins.sh` Finding 3 passes at zero and fails when a guard is re-added
- [ ] `gofmt`, `go vet`, `GOOS=windows go vet`, `-race` clean; gate shows no new red beside the two
  environmental `heavyverification` failures recorded at clean head

## Pre-Flight

**Gate at the current head is green in the maintainer wrapper** (`Maintainer verification passed.`,
fast tier, on the tree that merged REQ-597); the judge's clone sees two `heavyverification` failures at
clean head that are environmental to the reviewers' sandbox and are to be recorded as the pair beside
the post-change run, so that "no new red" is checkable.

**Lanes that read these files:** the `gittransaction` package suite (green, `ok`), `audit-lockins.sh`
Finding 3 (pinned at 8 by REQ-558; moves to zero here), and the CLI module's `go vet`/`gofmt` lanes. The
final patch applies cleanly on the current head (`git apply --check`), and the three files it touches
are unchanged since its base.

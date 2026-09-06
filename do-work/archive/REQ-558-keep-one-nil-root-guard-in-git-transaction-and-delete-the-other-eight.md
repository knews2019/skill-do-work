---
id: REQ-558
title: '[impact-negligible] Keep one nil-root guard in git_transaction.go and delete the other eight'
status: completed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
estimate:
  p50_active_minutes: 30
  confidence: low
  calculated_at: 2026-09-06T05:51:12Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - transaction boundary with no nil-branch coverage
suggested_spec:
depends_on: [REQ-557]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
route: B
write_set: [skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T05:50:44Z
completed_at: 2026-09-06T07:40:28Z
commit: 86c9d4cf091fef712c1230ad036f3f721e84ad60
release_at: 2026-09-06T07:40:28Z
---

# Keep one nil-root guard in git_transaction.go and delete the other eight

## What
`internal/gittransaction/git_transaction.go` tests one `*os.Root` value for nil nine times; no test exercises any nil branch, no introducing commit or REQ names a nil-root failure, three caller⇒callee pairs test the same value twice on one path, and two functions taking the same value never test it. Keep one nil test at the single point where nil is producible (`rollbackTransaction`'s `os.OpenRoot`), pass a non-nil root or an explicit no-root branch downward, and delete the other eight guards.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Three independent reachability traces before any edit — two by call graph, one
  adversarial. Approach: delete only what the trace proves unreachable, keep everything else with a
  named reason, keep anything the trace could not settle, and pin the real surviving count rather than
  the requested one.
- [x] **[APPLY]:** Two files, both declared: one guard deleted from `git_transaction.go` with its
  precondition stated as a doc comment, and one assertion block added to `audit-lockins.sh`. No test
  file beyond the lock-in.
- [x] **[UNIFY]:** `git diff --numstat` — `git_transaction.go` +10/-3, `audit-lockins.sh` +45 (corrected after review: the
  first version quoted `--stat`'s 13, which is the changed-line total, as insertions). Both read in
  full. The deletion's thirteen call sites were each traced to an open that returns on failure or to a
  surviving guard. `go build`, `go vet` and `gofmt -l` all clean; `go test ./internal/gittransaction/`
  green and unchanged at 2.353s against 2.205s before. No debug artifacts.

## Why
Defensive checks with no incident behind them are the agent-creep class; here they accreted per REQ on the top do-work-cli hotspot (17 commits, file CCN 392, `ExecuteTransaction` CCN 87) and none is covered by a test.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 3, sweep_key `nil-root-guards-git-transaction`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -25. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- Guards at `inspectCreatedObject` (b877eb69, REQ-457: the REQ never mentions a nil root), `rollbackDirtyTracked` (0a5d4e44, REQ-491: same), `trackedPublicationStillOwned` (a43b2587, no REQ id), `rootedOpenSnapshot` (01d920dd, no REQ id), and two `root != nil && privateStateStillOriginal(root, state)` sites (01d920dd).
- Redundant pairs verified by reading: the two `privateStateStillOriginal` callers, `trackedPublicationStillOwned`, and `rollbackDirtyTracked` all reach `rootedOpenSnapshot`, which tests the same value again.
- Inconsistency evidence: `privateStateStillOriginal` and `rootedCreateRegular` take the same `root *os.Root` and do not test it.
- Behaviour preserved on every rollback path: the package's transaction and rollback tests are the safety net and stay green unchanged.
- Reproduce at dc8a64e3 (prints 9 guards and `NO TEST covers any nil-root branch`): `rg -n 'root [=!]= nil' skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go; rg -l 'OpenRoot|rollback root is unavailable|rooted filesystem handle is unavailable' skills/do-work/tools/do-work-cli/internal/gittransaction/*_test.go || echo 'NO TEST covers any nil-root branch'`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- This file only; the rollback path is a transaction boundary, so every change is behaviour-preserving and proved by the existing tests, never by a new mock.
- Lock-in limit: nil-root guards in git_transaction.go: 1 after this REQ (today 9).

## Dependencies
Depends on REQ-557 so no other audit REQ writes under `internal/` while the transaction file is refactored. Last of the batch.

## Builder Guidance
Firm on one producible-nil point; latitude on whether downstream functions take a non-nil root or an explicit no-root branch.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine `root == nil` / `root != nil` sites and no test covers any of them.
**GREEN when:** It prints one site; `go test ./internal/gittransaction/` green unchanged; the lock-in pins nil-root guards in that file at 1.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for nil-root-guards-git-transaction.*

## Triage

**Route: B** — Explore then build.

**Reasoning:** The request's baseline holds exactly at HEAD, which is unusual in this batch: the
reproduce command prints nine `root [=!]= nil` sites and `NO TEST covers any nil-root branch`. But the
change is inside a transaction boundary, and the claim that matters is not "there are nine guards" — it
is "nil is producible at exactly one of them". Proving that means tracing every path that can reach each
of the nine with a nil value, and no test exercises any of those branches, so the compiler and the
existing suite cannot settle it. That is discovery.

**Planning:** Skipped. One file plus one lock-in assertion; the work is whatever the trace establishes.

**Deleting a guard is not the same shape as deleting a duplicate.** REQ-557 removed copies that were
provably interchangeable. Here each guard is the only thing standing between a nil dereference and a
rollback path. A guard that is genuinely unreachable costs nothing to delete and a guard that is not
costs a panic during rollback, which is the worst moment available. The exploration's job is to tell
those two apart per site, and to say plainly where it cannot.

## Plan

**Planning not required** — Route B: one source file plus one lock-in assertion, and the edit set is
whatever the reachability trace establishes.

*Skipped by work action*

## Exploration

Three read-only agents traced the nine guards independently — two by call graph, one adversarially, by
trying to reach each guard with a nil handle rather than by reasoning about it. Full reports in the run
directory. **They agree, and they contradict the request.**

**The request's inference is "one producer of nil, therefore one guard". The producer is real; the
inference is not.** `rollbackFailure` opens its handle with `os.OpenRoot` and, on failure, records the
error and **deliberately keeps going** — that is what lets a failed transaction still unstage its paths
and still return a typed incomplete rollback. So a nil handle flows onward from that one site into
**eleven consumers across four independent loops**, and each consumer decides for itself what an
unusable handle means for the target in hand: `createdObjectReplaced`, `false`, an error, or skipping
the check. There is no chokepoint that dominates them.

**Eight of the nine guards are reached with nil in unmodified end-to-end runs**, and deleting any one of
seven turns a reported incomplete rollback into a nil-pointer panic. For six of the eight the mechanism
is just the open failing while `rollbackFailure` keeps going. The two in the dirty-tracked branch need
one more condition, stated after review: that branch runs `git -C <root> restore --staged` first and
`continue`s when git fails, so the open must fail while git still works — an execute-only repository
root, traversable but unreadable, does exactly that. `(*os.Root)` panics on a nil
receiver — confirmed empirically for `Lstat`, `Open`, `Remove`, `Mkdir`, `Rename`, `MkdirAll`,
`OpenFile` and `Close`. Two examples:

- `rollbackDirtyTracked` (line 1176) has one caller, which passes the handle unchecked; the `root != nil`
  test four lines above it guards only a sibling branch. Both of its own branches dereference — `Lstat`
  directly, or `Mkdir` through `quarantineAndRollbackPrivate`.
- `trackedPublicationStillOwned` (line 1193) is worse than either answer: delete it and one branch is
  unchanged, because a downstream guard returns the same `false`, while the other panics. Half-silent and
  half-fatal is the worst shape a deletion can have.

**Exactly one guard is genuinely unreachable.** `rootedOpenSnapshot` (line 1276) is entered through
three wrappers from thirteen call sites. Nine carry a handle from an open that returns immediately on
failure; three sit behind guards that return first; **and the fourth — inside
`quarantineAndRollbackPrivate` — sits behind no guard at all and is unreachable with nil only because
`root.Mkdir` nine lines earlier panics first, which is the defect REQ-598 files** (corrected after
review; the first version said all four were guarded). It re-tested a question its callers had already
settled — **and it becomes reachable the moment any one of four guards this request wanted deleted is
gone**: `inspectCreatedObject`'s, either `privateStateStillOriginal` short circuit, or
`trackedPublicationStillOwned`'s (corrected after review from "two").

**Two of the request's supporting claims run the other way.** It cites
`privateStateStillOriginal` and `rootedCreateRegular` taking the same handle without testing it as
evidence that fewer guards are needed. They are safe only because guards at three call sites stand in
front of them, so that pair is an argument for keeping those guards, not for deleting them.

**And the trace found a live defect the request did not know about.** `rollbackFailure` hands the
possibly-nil handle to `quarantineAndRollbackPrivate` with **no check at all**, and that function
dereferences it with `root.Mkdir`. Any transaction with an identity-recorded private untracked target
panics mid-rollback when the open fails. Two of the three agents reproduced it independently. That is
captured as **REQ-598**; it is the opposite class from this request, it needs the package's first
no-handle test, and this request's constraints forbid both.

*Generated by three Explore agents, judged against the code by the builder*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — delete the one guard the trace proves unreachable, and state its precondition as a doc comment
- `_dev/tests/audit-lockins.sh` (modify) — one assertion block pinning the surviving count as a floor and a ceiling

**The delivered change is one deletion, not eight, and the pinned number is 8, not 1.** That is a
deviation from the request's stated outcome, decided on evidence and stated here before the edit rather
than discovered in review. The request asked for a count; the trace answered a question the count stood
in for, and the two answers differ. Seven of the eight deletions it named would each replace a reported
incomplete rollback with a nil-pointer panic.

**Files I will NOT touch:** any test file beyond the lock-in, which the request forbids and which is
also why the missing guard found on the way cannot be fixed here — the fix needs the package's first
no-handle rollback test. `quarantineAndRollbackPrivate`'s unguarded call site, captured as REQ-598.
`privateStateStillOriginal` and `rootedCreateRegular`, which the request cites as inconsistent and which
the trace shows are safe only because of guards it wanted deleted.

**Acceptance criteria:**
- [ ] Every deleted guard is proved unreachable by tracing every caller, not by counting
- [ ] Every kept guard has a named reason, and no guard is kept on a maybe
- [ ] Behaviour preserved on every rollback path, proved by the package's existing transaction and
  rollback tests staying green and unchanged
- [ ] The lock-in pins the real surviving count in both directions, and fails loudly when its scan
  cannot run
- [ ] The deviation from the requested count is stated as a decision with its evidence
- [ ] The live defect found on the way is captured, not fixed here

## Pre-Flight

**Green gate at `46da9507`**, the revision the builder branched from.
`bash _dev/tests/maintainer-verify.sh` printed `Maintainer verification passed.` and exited 0.

**The safety net is the package's existing suite, and its limits are stated rather than assumed.**
`go test ./internal/gittransaction/` is green before and after — but it exercises **no** no-handle
branch. Verified: `rg -l 'OpenRoot|rollback root is unavailable|rooted filesystem handle is unavailable'`
over the package's `*_test.go` files returns nothing. So a green suite is not evidence that a deleted
guard was dead, which is precisely why this request needed a trace and not a count.

**The baseline holds exactly at HEAD**, which was not true of several siblings in this batch: the
reproduce command prints nine `root [=!]= nil` sites and `NO TEST covers any nil-root branch`.

**`(*os.Root)` panics on a nil receiver**, confirmed empirically for `Lstat`, `Open`, `Remove`, `Mkdir`,
`Rename`, `MkdirAll`, `OpenFile` and `Close` rather than assumed from the type. That is the fact every
keep verdict rests on.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`
- `_dev/tests/audit-lockins.sh`

**What was done: one guard deleted, eight kept, and the count pinned at 8.** The request asked for eight
deletions and a pin of 1. Seven of those eight deletions would each replace a reported incomplete
rollback with a nil-pointer panic.

**The one deletion.** `rootedOpenSnapshot`'s guard is entered through three wrappers from thirteen call
sites, all in this file and none in a test. Nine carry a handle from an `os.OpenRoot` that returns
immediately when the open fails; the other four sit behind guards that are kept and that return before
the call. It re-tested a question its callers had already settled. Go cannot express a non-nil pointer in
the type, so the deleted test is replaced by a doc comment stating the precondition and naming today's
no-handle branches as the current set rather than a closed list — and naming the one caller the deleted
test never covered anyway.

**The eight kept, each with its own reason**, none on a maybe. Two reasons are restated after review,
because the first version described the file before the change rather than the one that shipped:
`inspectCreatedObject` — with `rootedOpenSnapshot`'s check gone, this is now the only thing that turns a
nil handle into `createdObjectReplaced` instead of a panic inside `root.Lstat` (before the change,
deleting it alone was byte-identical in behaviour, which is exactly why it looked redundant); the
`Close` at line 998, which would panic while returning and destroy the rollback report; the two
`root != nil &&` short circuits, which stand in front of a direct `root.Lstat`; the created-paths
`Lstat` at 1114, which has no downstream guard; `rollbackFailure`'s `!recorded || root == nil`, the only
thing between a nil handle and both a `Lstat` and a `Remove`; `rollbackDirtyTracked`, whose one caller
passes the handle unchecked and whose both branches dereference; and `trackedPublicationStillOwned`,
whose deletion is **fully fatal in the delivered file** — both branches now dereference, since the
downstream check that used to make one branch silent is the one this request deleted. (The
"half-silent and half-fatal" description belonged to the pre-change file and now lives in the
Exploration, where it is a fact about that file.)

**The lock-in pins 8 as a floor and a ceiling**, scanning only `git_transaction.go` with the audit's own
pattern and reading `rg`'s exit status rather than judging a piped total: a status above 1 means the scan
could not run and fails loudly, and a status of 1 means every guard is gone and fails too.

## Decisions — implementation

- **D-01 — the delivered count is 8, not the 1 the request names. DECIDE & STATE.** The request reasoned
  from one producer of nil to one guard. The producer records its error and keeps going by design, so
  the nil handle reaches eleven consumers in four loops and each answers for its own target. There is no
  chokepoint. Pinning 1 would have required seven deletions that each introduce a panic in rollback.
- **D-02 — a guard the trace cannot settle is kept.** No trace returned `cannot-establish`, so this rule
  bound nothing here, but it was the rule going in: a guard on a rollback path is not deleted on a maybe.
- **D-03 — the deleted guard's precondition is a doc comment naming an open set, not a closed list.**
  `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale applies to prose that enumerates
  callers as much as to shell. The comment says what must hold and names today's branches as the current
  set.
- **D-04 — the live defect found on the way is captured, not fixed.** `quarantineAndRollbackPrivate` is
  called with the possibly-nil handle and dereferences it. Fixing it means adding a guard, which is the
  opposite of this request's class, and proving the fix means writing the package's first no-handle test,
  which this request's constraints forbid. It is **REQ-598**.
- **This request is a release, and two releases already shipped its bytes.** `git_transaction.go`
  ships under `skills/`, so the change carries a version bump and a changelog entry, which the
  finalization manifest owns. Stated after review: because this request was held for review while
  siblings finalized, 0.305.5 and 0.305.6 both went out on top of its merge with no entry naming it. The
  behaviour is preserved on every path, so nothing a user sees changed; the entry this request writes
  says it was already present in those two versions.

## Qualification

**Passed.** Read from the merge range `46da9507..86c9d4cf`, two files, 55 insertions and 3 deletions.
Canonical `qualify` and `scope-drift` both satisfied.

- **The request asked for eight deletions and got one, and the review could not break that.** The
  synthesizer restored the deleted guard as a `panic("CANARY…")` in a scratch clone and ran the full
  package suite plus nine purpose-built probes that inject nil at the exact producer and cover every
  branch of `rollbackFailure` that touches the handle. The canary never fired. Then it ablated each of
  the eight kept guards individually with the same nil injection: **every ablation panics in a real
  `ExecuteTransaction` run**, and guard-hit instrumentation confirmed each of the eight is actually
  entered with a nil handle. Seven of those eight are guards the request named for deletion.
- **The thirteen call sites of the deleted check were enumerated independently of the record**: nine
  hold a handle from one of the seven `os.OpenRoot` sites that return on failure, three sit behind kept
  guards, and the fourth — inside `quarantineAndRollbackPrivate` — sits behind no guard and is
  unreachable with nil only because `root.Mkdir` panics nine lines earlier. That fourth is REQ-598.
- **Behaviour is preserved on every path that worked before.** The package suite is green and
  unchanged; the one probe that fails, a private untracked target with a recorded publication under a
  failed open, fails identically at the branch point — it is the pre-existing panic, reproduced end to
  end through the public API.
- **The lock-in did not hold the property its header claimed, and it does now.** It scanned for the
  bare text `root [=!]= nil`, so deleting a load-bearing guard and adding one comment line containing
  that text kept the count at eight and every gate green — the review did exactly that with
  `trackedPublicationStillOwned`'s guard and watched a rollback-crashing deletion pass the repository.
  The scan is now anchored to a guard shape on a non-comment line in either operand order. Proven:
  deleted guard plus a comment mentioning `root != nil` → **7, red**; `if nil == root` rewrite → 8,
  green, as a no-op should be; a ninth guard added → **9, red**.
- **The pin's own message told the next maintainer the guard set was complete**, which is the belief
  REQ-598 exists to correct. It now says one consumer is still unguarded and names the request that owns
  it, and the header says the pin is expected to move when that lands. REQ-598's write set gains the
  lock-in file, so its author has a declared path to re-pin rather than a red fast tier and nothing it
  may edit.
- **Six record claims are corrected in place**: the UNIFY diff count, "the other four sit behind
  guards", "two of the guards" (any one of four is enough), the missing condition under which the two
  dirty-tracked guards are reachable, and two keep-reasons that described the file before the change
  rather than the one that shipped. The verdict survived every one of them; the reasons did not.

## Testing

**The evidence for this request is a reachability trace, not a test run, and the record says so.** No
test in the package reaches any no-handle branch, so `go test ./internal/gittransaction/` — `ok 2.353s`
before, `ok 2.353s` after, unchanged — proves only that nothing that worked stopped working.

**The trace was then attacked by three reviewers and held.** Canary in the deleted guard's position:
never fired across 44 package tests and nine nil-injection probes. Each kept guard ablated alone: eight
panics, each named to its dereference — `root.Lstat` in `rootedOpenSnapshot` for four of them, the
deferred `Close` at the return statement, the direct `Lstat` at 1114, the created-directories `Lstat`,
and `root.Mkdir` inside `quarantineAndRollbackPrivate`.

**The lock-in, red in both directions and for the right reasons:** a deleted guard paid for by a
comment → 7, red; a ninth guard → 9, red; a `nil == root` no-op → 8, green; restored → green.
`bash -n` and `shellcheck --severity=warning` exit 0.

**Gate at the builder's head:** `Maintainer verification passed.`, exit 0, wall 89s. `go build`,
`go vet` and `gofmt -l` clean.

## Review

**Overall: 88%** | 2026-09-06T09:20:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 88% |
| Code Quality | 85% |
| Test Adequacy | 70% |
| Scope | 98% |
| Risk | Low |
| Acceptance | Pass with corrections |

**Verdict: Pass with corrections.** "The trace's conclusion is correct and I could not break it. The
request was wrong and the builder was right to refuse it." What did not hold up: the record's prose in
the sections whose only job is to justify the deviation, and the lock-in's scan, which counted matching
text rather than guards so that the exact move it was written to catch still passed.

Where the reviewers disagreed, and what was picked:

- Whether the deleted guard is reachable. Two reviewers enumerated call sites statically; one placed a
  canary and drove nine probes. All three agree it is unreachable; the canary is the evidence that
  settles it.
- Whether any kept guard is dead. One reviewer argued each of the eight could be dead from reading. The
  synthesizer ablated each and drove the probes: every ablation panics. Picked the ablations.

**Important findings:**

- The lock-in's floor counts matching text, not guards: delete a load-bearing guard, add one comment
  containing `root != nil`, and every gate in the repository stays green. — impact-rule-change → fixed
  in remediation
- The pin at exactly 8 blocks REQ-598 — its minimum fix adds a ninth guard, its preferred fix removes
  all eight — and REQ-598's write set had no path to move it. — impact-rule-change → REQ-598's write set
  extended, the lock-in's header says the pin is expected to move

**Minor findings:**

- Three of the four remaining call sites sit behind kept guards; the record said four. The fourth is
  unguarded and safe only because of the REQ-598 panic. — corrected
- `trackedPublicationStillOwned`'s keep-reason described the pre-change file; in the delivered file its
  deletion is fully fatal. — corrected
- `inspectCreatedObject`'s keep-reason described a control flow that does not occur in either revision;
  the guard is load-bearing for a different reason. — corrected
- "Two of the guards" understated a threshold that is any one of four. — corrected
- The record never named the condition under which the two dirty-tracked guards are reachable: an
  execute-only repository root, where the open fails but `git -C` still works. — corrected
- UNIFY quoted `--stat`'s changed-line total as insertions. — corrected
- The pin's FAIL message said "one per consumer, and no more" when one consumer has no guard. — fixed
- D-04 framed the release as ahead of the change; two releases already shipped it unnamed. — corrected
- The new doc comment's mood: "callers satisfy" reads as fact and then names an exception four lines
  later. — report only; the content is right and the branch list is open

**Requirements checklist:**

- [x] Every deleted guard is proved unreachable by tracing every caller — delivered, and attacked with a
  canary
- [x] Every kept guard has a named reason and none is kept on a maybe — delivered, each proven by
  ablation; two reasons restated after review
- [x] Behaviour preserved on every rollback path — delivered
- [ ] → [x] The lock-in pins the real count in both directions and fails loudly when its scan cannot run
  — **not delivered at review**: a comment could pay for a deleted guard; **delivered in remediation**
- [x] The deviation from the requested count is stated as a decision with its evidence — delivered
- [x] The live defect found on the way is captured, not fixed — delivered as REQ-598, now with a path to
  move the pin

**Acceptance testing**

**Result: Pass with corrections at review, Pass after remediation.** Three reviewers built nil-injection
probes through the public `ExecuteTransaction` API rather than calling private functions, and left the
working tree clean.

**Follow-ups created:** none new; REQ-598 gains a file in its write set.

*Reviewed by review-work action*

## Lessons Learned

- **A request that reasons from a count to a refactor has to be re-derived, not executed.** "One
  producer of nil, therefore one guard" was a plausible sentence about a real producer. The producer
  records its error and keeps going *on purpose*, so the nil reaches eleven consumers. Nine guards was
  not redundancy; it was a per-consumer answer to a per-consumer question. The audit that generated this
  request counted the guards and did not trace one.
- **A guard that looks redundant because a sibling catches the same case is only redundant while the
  sibling exists.** `inspectCreatedObject`'s guard could be deleted alone with byte-identical behaviour
  before this change, because `rootedOpenSnapshot`'s guard returned an error the switch absorbed. This
  request deleted that sibling. The first guard is now the only thing between a nil handle and a panic.
  Redundancy is a property of a set, and deleting from the set changes it for every member left.
- **The safety net must be named honestly or it is not one.** The package suite was green before and
  after, and it reaches no no-handle branch at all. A record that had said "the existing tests prove
  this is safe" would have been believed. The record said the suite proves only that nothing that
  worked stopped working, and that is what let the review go looking for the canary.
- **A ratchet on text is a ratchet on text.** Counting `root [=!]= nil` counted a comment. The review
  deleted a guard whose absence crashes a rollback, wrote one comment line, and every gate in the
  repository passed. Anchor a count to the shape of the thing being counted — an `if` on a non-comment
  line — and say in the block what shape it still cannot see.
- **A pin that will block the next request should say so.** The count was pinned at eight while the
  request that adds a ninth guard, or removes all eight, was already queued — by this same record. The
  pin was right; leaving its successor with no declared path to move it was not.
- **Describe the file that shipped, not the one you traced.** Two keep-reasons were written from the
  pre-change control flow and were wrong about the delivered one, in the direction of making the
  deletion look safer than it was. When a change removes a guard, re-read every reason that mentioned
  that guard's neighbours.

## Orientation

`internal/gittransaction/git_transaction.go` opens eight rooted filesystem handles. Seven return the
moment `os.OpenRoot` fails. The eighth, in `rollbackFailure`, **records the failure and keeps going on
purpose**, so a failed transaction still unstages its paths and still returns a typed incomplete-rollback
result when the worktree root cannot be opened. That one possibly-nil handle fans out to eleven
consumers across four loops, and each consumer decides for itself what an unusable handle means for the
target in hand — `createdObjectReplaced`, `false`, a typed error, or skipping the check.

**Eight nil checks guard those consumers, and every one is load-bearing.** Each was ablated alone under a
nil injection and each ablation panics in a real `ExecuteTransaction` run. `rootedOpenSnapshot` no longer
tests the handle: every path into it either holds a handle from an open that returned on failure or has
already taken one of those eight branches, and its doc comment states that precondition and names
today's branches as the current set, not a closed list.

**One consumer has no guard: `quarantineAndRollbackPrivate`, called from `rollbackFailure` with the
possibly-nil handle, dereferences it at `root.Mkdir`.** Any transaction with an identity-recorded private
untracked target panics mid-rollback when the open fails. That is REQ-598, which also owns the better
answer — deciding the handle once at the open, which would make all eight guards dead — and the
package's first no-handle rollback test.

The pin is Finding 2 in `_dev/tests/audit-lockins.sh`: exactly eight `if`-shaped nil tests on
non-comment lines, floor and ceiling. It cannot see a guard rewritten as a helper call. Its header says
the pin is expected to move when REQ-598 lands, and REQ-598's write set includes the file.

Reaching the two guards in the dirty-tracked branch needs one condition beyond the open failing: that
branch runs `git -C <root> restore --staged` first and continues on git failure, so the root must be
execute-only — traversable but unreadable — for the open to fail while git still works.

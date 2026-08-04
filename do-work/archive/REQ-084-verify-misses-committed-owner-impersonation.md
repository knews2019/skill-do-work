---
id: REQ-084
title: verify's queue-state probe misses a builder that committed its do-work edits
status: completed
claimed_at: 2026-08-03T23:54:52Z
completed_at: 2026-08-04T00:00:52Z
created_at: 2026-08-03T17:09:21Z
kb_status: pending
user_request: UR-016
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-072, REQ-082, REQ-083]
batch: audit-remediation-external
addendum_to: REQ-072
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go, actions/forensics.md]
---

# `verify`'s Queue-State Probe Misses a Builder That Committed Its `do-work` Edits

## What

The owner-impersonation probe runs only `git status --porcelain … -- do-work/` inside the builder's
worktree, so it sees uncommitted changes and nothing else. A builder that edits `do-work/` and
**commits** on its own branch leaves a clean worktree and passes verification — which is the more
likely shape, since a builder commits its work by design.

REQ-072 requirement 3 asked for something wider, verbatim: "any `worktree-agent-*` worktree whose
`do-work/` **differs from the main tree's** (an owner impersonation — a builder that wrote queue
state)."

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `tools/queue-kanban/prime-do-kanban.md` (including the REQ-083 entry added
  hours earlier), `crew-members/general.md`, `coding-guardrails.md`, `testing.md`, and REQ-082's
  archived hand-back rule. Approach: resolve the integration ref once per run from the repo-root
  checkout, then per leftover run `git diff --name-only <ref>...<branch> -- do-work/` **in addition to**
  the untouched porcelain check, emitting a sibling category. Place the committed check *before* the
  `hasWorktree` guard, since committed edits live in the branch and outlive the worktree.
- [x] **[APPLY]:** Code written as planned. Three files, all declared in `## Scope`.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file.
  - `tools/queue-kanban/verify.go` — checked `worktreeDirtyQueueState` is byte-identical (body and doc
    comment both untouched, as Builder Guidance requires); the two new functions sit above it; the one
    new constant joined the aligned block; the loop gained one guarded stanza and the ref resolution
    above it. No other probe touched. No debug artifacts, no `fmt.Println`, no `TODO`.
  - `tools/queue-kanban/verify_test.go` — checked the three new tests plus one helper; no existing test
    modified; the REQ-083 fixture helpers are reused rather than duplicated.
  - `actions/forensics.md` — checked the probe-table edit: one row became two, no other row disturbed.
  - Linters: `gofmt -w` (applied), `go vet ./...` (clean), `go test ./...` (pass).

## Why

This probe exists to catch the one thing worktree isolation cannot prevent structurally: a builder
writing the state the orchestrator alone owns. Scoped to uncommitted changes, it catches the case where
the builder was interrupted mid-write and misses the case where it finished — so it detects the
accident and not the behavior.

The narrowing was never recorded as a decision. REQ-072's `## Review` paraphrases the delivered probe
as covering "any such worktree carrying uncommitted `do-work/` changes," which reads as a restatement
of the requirement rather than a departure from it. No D-XX, no Open Question, no note under suggested
testing. That is what makes this worth a REQ rather than a one-line patch: the gap is invisible from
the archive.

## Context

`tools/queue-kanban/verify.go:427-447`, with its stated reasoning:

```go
// worktreeDirtyQueueState returns the do-work/ paths a builder's worktree has
// modified. Uncommitted changes there are the detectable signature of the "state
// stays home" rule being broken; a stale committed snapshot (which a worktree
// legitimately carries where the consumer commits do-work/) is not.
func worktreeDirtyQueueState(worktreePath string) []string {
	command := exec.Command("git", "-C", worktreePath, "status", "--porcelain", "--untracked-files=all", "--", "do-work/")
```

**The stated concern is real; the narrowing is not the only way to address it.** Where a consumer
commits `do-work/`, the worktree carries a snapshot from its branch point, and the main tree moves on
constantly as the orchestrator claims and archives — so a plain two-tree comparison would fire on
almost every run. A **merge-base** comparison does not: it shows only what the *branch* changed, and is
blind to how far main has moved.

**Reproduced.** Same synthetic repo and worktree as REQ-083. A forged
`do-work/queue/REQ-999-forged.md` written inside the builder's worktree:

```
=== A: uncommitted ===
! worktree-wrote-queue-state: <path> has uncommitted changes under do-work/ (do-work/queue/REQ-999-forged.md) — a builder wrote queue state the orchestrator alone owns

=== B: same file, committed in the worktree ===
$ git status --porcelain -- do-work/        # (empty — worktree is clean)
(no worktree-wrote-queue-state finding at all)

--- what a merge-base diff sees in state B: ---
$ git diff --name-only main...HEAD -- do-work/
do-work/queue/REQ-999-forged.md
```

**Interaction with REQ-082.** That REQ grants builders exactly one main-tree write — their own
`REQ-NNN-handback.md` — written by absolute main-tree path and never committed. Under those rules the
hand-back file cannot appear in a merge-base diff of the builder's branch, so it must not need an
exemption here. If it does, one of the two REQs is wrong; see requirement 6.

## Detailed Requirements

1. **Add a merge-base comparison alongside the porcelain check.** `git diff --name-only
   <integration-branch>...<builder-branch> -- do-work/` (three dots — symmetric difference against the
   merge base) is the shape proven above. Keep the porcelain check: it catches uncommitted edits, which
   a merge-base diff cannot see.
2. **Report the two states distinguishably.** "Committed queue edits on the builder's branch" and
   "uncommitted queue edits in the builder's worktree" call for different remedies — the first is in
   the branch about to be merged, the second is loose in a working tree. Reuse
   `verifyCategoryWorktreeWroteQueueState` or add a sibling constant; either way the detail line must
   say which it found.
3. **Neither state is `Fixable`.** Both are already correctly unflagged today. Discarding a builder's
   commits is never mechanical — keep that, and do not let REQ-083's classification work bleed into
   this probe.
4. **Do not fire on the legitimate stale snapshot.** The fixture must include a worktree whose
   `do-work/` is merely *behind* the main tree (orchestrator committed queue changes on the integration
   branch after the branch point) and assert **no finding**. This is the case the original narrowing
   was protecting, and the regression test for it is what makes widening the probe safe.
5. **Name the integration branch explicitly, not `HEAD`.** The probe runs against the repo root; derive
   the comparison point from the repo-root checkout and say so in a comment. The `-d`-tests-against-HEAD
   trap (`actions/work-reference.md:294`) is the same class of error.
6. **Confirm the REQ-082 hand-back file does not trip it.** Whichever of the two REQs lands second
   verifies this in its Qualification rather than assuming it. If the hand-back file *does* appear in
   the diff, that is a REQ-082 violation to report — not an exemption to add here.
7. **Degrade like every other probe.** A git command that cannot answer becomes a reported
   `SkippedProbes` entry, never silence. `VerifyReport`'s doc comment is explicit that a silently
   skipped invariant would read as a clean one.
8. **Keep `verify` read-only.** Additional read commands only; nothing written.

## Constraints

- **`do-work/` is not always tracked.** On the common install it is untracked, so a merge-base diff
  finds nothing there and the porcelain check (with `--untracked-files=all`, already present) is the
  only probe that can see anything. Both paths must work, and neither may error on the other's install
  shape.
- **Do not compare working trees directly.** `diff -r` between two checkouts is the tempting shape and
  it is wrong here: it fires on the legitimate snapshot divergence, and `CLAUDE.md`'s prescribed-command
  traps include `diff -x` matching directories as well as files.
- **Do not renumber or restructure the other probes**, same as REQ-083 — this is a change inside one
  helper.
- Wired consumers (`actions/forensics.md` Check 14, the Step 9 accelerator) read the report; if a
  category string changes, grep both.

## Dependencies

`addendum_to: REQ-072`, whose requirement 3 this completes. `related: REQ-083` (same file, same
function neighbourhood — no `depends_on`; whichever lands second re-reads before editing) and
`related: REQ-082` for requirement 6.

## Builder Guidance

**Certainty: Firm.** The requirement text is verbatim from REQ-072, the failure is reproduced, and the
fix shape is proven against a live repo. The one genuinely open choice is whether the two states share
a category constant or get two — log it as a D-XX either way.

Do not remove or soften the existing doc comment's reasoning about the legitimate stale snapshot. It is
correct, it is the reason requirement 4 exists, and a later reader widening this probe further needs it.

## Red-Green Proof

**RED prompt/case:** A Go test that builds a temp repo with a `worktree-agent-*` worktree, writes
`do-work/queue/REQ-999-forged.md` inside it, **commits it on the builder's branch**, then calls
`runVerifyProbes` and asserts a `worktree-wrote-queue-state` finding is present.

**Why RED now:** The worktree is clean after the commit, so `git status --porcelain -- do-work/`
returns nothing and `worktreeDirtyQueueState` returns nil — no finding is emitted. Reproduced by hand
already (state B above).

**GREEN when:** That assertion passes; the uncommitted case from state A still reports (so the
porcelain check was kept, not replaced); the stale-snapshot fixture from requirement 4 reports
**nothing**; `go test ./...` in `tools/queue-kanban/` stays green; and the manual probe reproduces:
`git diff --name-only main...HEAD -- do-work/` naming the forged file is now surfaced as a finding.

**Validation:** Inferred during capture, then reproduced — both the miss and the working fix shape were
executed against a synthetic repo before this REQ was written.

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the
external audit, and the batch constraints.

## Triage

**Route: B** - Medium

**Reasoning:** Requirements are firm and verbatim from REQ-072, the failure is reproduced, and the fix
shape (`git diff --name-only <integration>...<branch>`) is proven. What needs discovery is how to
resolve the integration ref safely from the repo root and how the new check composes with
`appendWorktreeFindings` as REQ-083 just left it.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Spec loaded:** `specs/bug-fix.md` — reproduce, isolate, minimal fix, regression test, root-cause note,
grep for the same pattern elsewhere.

**Re-read `appendWorktreeFindings` after REQ-083** (per this REQ's Dependencies note — landing second
means re-reading, not assuming). REQ-083 changed only the *leftover* finding's category and `Fixable`
routing; the dirty-`do-work/` sub-probe at the tail of the loop is byte-identical to what this REQ
described, and `worktreeDirtyQueueState` was not touched. The two REQs do not interact beyond sharing
the function.

**The gap, precisely.** `worktreeDirtyQueueState` shells out to `git -C <worktree> status --porcelain
--untracked-files=all -- do-work/`. A commit clears the working tree, so the probe's input goes empty
and it returns nil — indistinguishable from "clean". The existing doc comment's reasoning is correct
and must survive: a worktree legitimately carries a *stale snapshot* of `do-work/` from its branch
point, and the main tree moves constantly as the orchestrator claims and archives.

**Why three dots is the whole fix.** `git diff A...B` diffs from `merge-base(A,B)` to `B` — it reports
what B's branch *added*, and is blind to everything A gained after the branch point. That is exactly
the distinction the doc comment protects: a builder's own queue edits show; the orchestrator's newer
queue commits on the integration branch do not. A two-tree comparison (`A..B` or `diff -r`) would fire
on every run, which is presumably why the original narrowed to porcelain instead.

**Resolving the integration ref (requirement 5).** `git -C <repoRoot> rev-parse --abbrev-ref HEAD`
names the repo-root checkout's branch. Run from the repo root, `HEAD` and the integration branch are
the same thing — but naming it explicitly is what stops the comparison silently re-anchoring if the
command is ever run with a different `-C`, which is the whole `git branch -d` trap. A detached
repo-root checkout returns the literal string `HEAD`; the commit id from `git rev-parse HEAD` names the
same point just as explicitly and is the fallback.

**Requirement 6 — confirmed structurally, not assumed.** REQ-082's archived text (line 173, restated at
415) settles it: the builder writes its hand-back "by absolute main-tree path — and never commits it."
Two independent reasons it cannot trip either state — it lands in the **main tree**, so it is not in the
builder's branch for a merge-base diff to find, and it is **never committed**, so it could not appear
there even if it were. It is equally invisible to the porcelain check, which runs inside the worktree.
No exemption is needed, which is what REQ-082 requirement 7 predicted. Pinned with a fixture assertion
rather than left as an argument.

**Both install shapes (Constraints).** Where `do-work/` is untracked (the common install), the
merge-base diff simply returns nothing and the porcelain check with `--untracked-files=all` remains the
only probe that can see anything — the new check adds a no-op read, not an error. Where the consumer
commits `do-work/`, both checks are live and cover the two halves.

**Degradation (requirement 7).** The new check has a genuinely reachable failure: a `worktree-agent-*`
worktree in detached HEAD has no branch of that name, so the diff cannot resolve — the same case
REQ-083 classifies as `worktree-merge-state-undetermined`. That must become a `SkippedProbes` entry,
not silence.

## Scope

**Files I will touch:**
- `tools/queue-kanban/verify.go` (modify) — add the merge-base comparison and the integration-ref
  resolver; emit a distinguishable finding; keep the porcelain check unchanged
- `tools/queue-kanban/verify_test.go` (modify) — the committed-edit RED case, the retained uncommitted
  case, the stale-snapshot no-finding regression test, and the REQ-082 hand-back assertion
- `actions/forensics.md` (modify) — Check 14's probe row says "uncommitted `do-work/` changes"; the
  probe now covers committed ones too, so the row goes stale the moment this lands (declared up front,
  same as REQ-083's D-02)

**Files I will NOT touch:**
- `worktreeDirtyQueueState`'s existing body or its doc-comment reasoning — explicitly protected by
  Builder Guidance; the new check is added alongside it
- REQ-083's merge-state classification — requirement 3 forbids letting it bleed into this probe
- The other seven probes, and any file split of `verify.go`

**Acceptance criteria (restated from REQ):**
- [ ] A merge-base comparison runs alongside — not instead of — the porcelain check (req 1)
- [ ] Committed and uncommitted states are reported distinguishably, and the detail says which (req 2)
- [ ] Neither state is `Fixable` (req 3)
- [ ] A worktree merely *behind* the main tree produces no finding, pinned by a fixture (req 4)
- [ ] The integration ref is named explicitly from the repo-root checkout, with a comment (req 5)
- [ ] REQ-082's hand-back file is confirmed not to trip it (req 6)
- [ ] A git command that cannot answer becomes a reported `SkippedProbes` entry (req 7)
- [ ] `verify` still writes nothing (req 8)

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/verify.go` (modified)
- `tools/queue-kanban/verify_test.go` (modified)
- `actions/forensics.md` (modified)

**What was done:** `appendWorktreeFindings` now resolves the integration ref once per run
(`resolveIntegrationBranchRef`, from the repo-root checkout, falling back to the commit id on a
detached checkout) and runs `worktreeCommittedQueueState` — `git diff --name-only <ref>...<branch> --
do-work/` — for every `worktree-agent-*` name, emitting the new
`worktree-committed-queue-state` finding when the builder's branch carries queue edits. The existing
porcelain check and `worktreeDirtyQueueState` are untouched, so both states are now covered and
reported distinguishably. Neither is `Fixable`. A ref that cannot be resolved becomes a
`SkippedProbes` entry rather than silence.

**Root cause (per `specs/bug-fix.md`):** the probe was scoped to its *detection method* rather than to
the rule it enforces. "A builder must not write queue state" was implemented as "the worktree has
uncommitted `do-work/` changes" — which catches the builder interrupted mid-write and misses the one
that finished, i.e. it detects the accident and not the behavior. The narrowing was a real defensive
response to a real problem (the legitimate stale snapshot), but it traded away the main case instead of
picking a comparison that excludes only the snapshot. Three-dot diff semantics are that comparison.

## Decisions

- **D-01 (DECIDE & STATE)** — *Requirement 2: two category constants, not one shared constant with a
  varying detail line.* `verifyCategoryWorktreeCommittedQueueState` joins the existing
  `verifyCategoryWorktreeWroteQueueState`. Reasoning: the remedies are genuinely different actions in
  different places — discard working-tree edits in the worktree, versus drop commits from a branch
  before it is merged — and a category is what a consumer greps for and what a report groups by. Two
  states that need two remedies are two findings. The detail line *also* says which, per the
  requirement, but the constant is what makes it machine-distinguishable.
- **D-02 (DECIDE & STATE)** — *The committed check runs before the `hasWorktree` guard, so it covers
  branch-only leftovers too.* The requirement did not specify placement. A builder's committed queue
  edits live in the branch, so they remain detectable — and just as wrong — after the worktree is
  removed; gating them behind a surviving worktree would reintroduce a narrower version of the very
  miss this REQ fixes.
- **D-03 (DECIDE & STATE)** — *Extend scope to `actions/forensics.md`, and split its one probe row into
  two.* Check 14's row said the probe covers "any such worktree with uncommitted `do-work/` changes",
  which this change makes false. After REQ-083's merge-state clause the single row was carrying two
  distinct probes and had grown unreadable, so the fix is one row per probe — leftover classification,
  and queue-state writes — rather than a third clause. Declared in `## Scope` before coding.

## Discovered Tasks

- **[low]** `worktreeDirtyQueueState` (`tools/queue-kanban/verify.go`) returns `nil` both for "no dirty
  paths" and for "the git command failed", so a failure there is silence rather than a `SkippedProbes`
  entry — the degradation contract requirement 7 states. Left untouched deliberately: Builder Guidance
  for this REQ explicitly protects that function, and the failure is close to unreachable (the worktree
  was just enumerated by git, so `git status` inside it does not realistically fail). Worth folding in
  if that helper is ever revisited.

## Qualification

Passed — 3 files verified, 8 requirements traced, P-A-U confirmed.

- `tools/checks/qualify.sh` → OK. `tools/checks/scope-drift.sh` → OK, both set-differences empty.
- **Substantive (check 2):** two new functions with real git invocations plus a guarded stanza in the
  loop; three new tests with real commits. Not whitespace.
- **Requirements traced (check 3):** 1 → `worktreeCommittedQueueState` added; `worktreeDirtyQueueState`
  confirmed byte-identical by diff, and `TestVerifyStillFlagsUncommittedQueueStateInAWorktree` proves
  the porcelain path still fires. 2 → separate constant + the word "committed" in the detail, asserted.
  3 → neither finding sets `Fixable`; asserted in both tests, and `FixableCount` is unaffected. 4 →
  `TestVerifyIgnoresAStaleQueueSnapshotAndTheHandBackFile` moves the main tree on after the branch
  point and asserts zero findings. 5 → `resolveIntegrationBranchRef` with the comment explaining why
  `HEAD` is not passed through. 6 → below. 7 → two `SkippedProbes` paths, one for the unresolvable ref
  and one per-branch. 8 → `TestVerifyWritesNothing` still passes; only read commands added.
- **Requirement 6 — the REQ-082 hand-back file, confirmed rather than assumed.** REQ-082's archived
  text (line 173, restated at 415) grants the builder exactly one main-tree write: its own
  `REQ-NNN-handback.md`, "by absolute main-tree path — and never commits it." Two independent reasons
  it cannot trip either state, and the probe needs no exemption: (a) it is written to the **main tree**,
  so it is not in the builder's branch and a merge-base diff of that branch cannot contain it; (b) it is
  **never committed**, so it could not appear in that diff even if it were written inside the worktree.
  It is equally invisible to the porcelain check, which runs with `-C <worktreePath>` and so never looks
  at the main tree. This is not left as an argument — the stale-snapshot test plants a real hand-back
  file at `do-work/runs/work-.../REQ-012-handback.md` in the main tree and asserts zero findings of
  either category. REQ-082 requirement 7's prediction holds; nothing to report against it.
- **Flowing (check 6):** not hollow — the new path was observed end-to-end on a live synthetic repo
  (below), not only through the fixture.
- **Contamination check (Step 10):** REQ-083 touched the same three files. The overlap is expected and
  declared — both REQs carry `related: REQ-072` and sit in the same function neighbourhood, and this
  REQ's Dependencies section instructs whichever lands second to re-read `appendWorktreeFindings`
  first, which was done. No REQ-083 content was modified: its merge-state classification, its three
  categories, and its three tests are untouched, and requirement 3's "do not let REQ-083's
  classification bleed into this probe" holds — the new finding sets no `Fixable` and consults no merge
  state.

## Testing

**Tests run:** `cd tools/queue-kanban && go test ./...` (+ `go vet ./...`, `gofmt`)
**Result:** ✓ All passing — 18/18 in the verify suite, whole module green

**Red-green validation:** traced to `## Red-Green Proof`, whose RED case was implemented verbatim.

- `TestVerifyFlagsQueueStateCommittedOnABuilderBranch`: ✗ before implementation
  (`got 0 committed-queue-state findings, want 1` — the report showed only the leftover finding, no
  impersonation finding at all) → ✓ after
- `TestVerifyStillFlagsUncommittedQueueStateInAWorktree`: passed before and after — a deliberate
  regression guard, not a red-green pair. It is what proves requirement 1's "alongside, not instead of."
- `TestVerifyIgnoresAStaleQueueSnapshotAndTheHandBackFile`: passed before and after — likewise a guard.
  It has teeth only because the committed test proves the probe *does* fire, so its silence here is a
  discrimination result rather than a probe that never fires.

The RED was observed, not assumed: the constant landed first with no logic, so the test failed on its
assertion rather than on a build error.

**New tests added:**
- `TestVerifyFlagsQueueStateCommittedOnABuilderBranch` — the RED case (req 1, 2, 3)
- `TestVerifyStillFlagsUncommittedQueueStateInAWorktree` — the porcelain check survives (req 1)
- `TestVerifyIgnoresAStaleQueueSnapshotAndTheHandBackFile` — the legitimate stale snapshot stays silent
  (req 4) and REQ-082's hand-back file does not trip it (req 6)
- Helper `writeFixtureFile`; REQ-083's `newWorktreeFixtureRepo` / `addFixtureWorktree` reused

**Manual acceptance (the REQ's GREEN condition):** rebuilt the binary and reran the hand reproduction
on a synthetic repo. `git -C <worktree> status --porcelain -- do-work/` prints nothing (the old probe's
entire input), while verify now reports
`worktree-committed-queue-state: worktree-agent-REQ-999-forger has committed changes under do-work/ on
its branch (do-work/queue/REQ-999-forged.md)`. Both halves of the REQ's state A / state B reproduction
confirmed.

*Verified by work action*

## Review

**Overall: 97%** | 2026-08-04T00:15:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Acceptance:** Pass — the committed forgery is surfaced on a live synthetic repo naming the exact
path, while the porcelain-only view of the same state stays empty; this repo's own tree reports clean.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Restatement Sweep (MUST — run):** The diff redefines what the queue-state probe *covers* (uncommitted
edits → uncommitted **and** committed) and adds a category string.

- *Category strings:* `grep -rn 'worktree-wrote-queue-state\|worktree-committed-queue-state'` across
  `.md`/`.js`/`.sh` outside `do-work/` returns nothing — no prose, template, or tool matches on either.
- *Coverage wording:* `actions/forensics.md:203-204` was the one place stating the narrow scope
  ("any such worktree with uncommitted `do-work/` changes") and is updated in this REQ. Swept every
  other `uncommitted` mention outside `do-work/`: `actions/version.md:73`, `actions/install.md:326`,
  `README.md:32` (the updater), `actions/commit.md:76`, `actions/review-work.md:30,66` (diff reading),
  and `actions/work-reference.md:308` (`git worktree remove` refusing on a dirty worktree, which is
  cleanup's behavior, not verify's) — none restates this probe. `CLAUDE.md:163`'s "queue and worktree
  invariants" is generic and still true.
- *The rule itself* — `actions/work-reference.md` → Worktree Dispatch Mode, "State stays home" and
  "Sole integrator" — is unchanged by this diff; verify now enforces more of it, so the canonical
  statement and its enforcement moved *closer* together, not apart.
- *One stale restatement, carried over:* `docs/forensics-guide.md:23` (Minor, below).

### Findings

**Minor — `docs/forensics-guide.md:23` summarizes verify's probes without either worktree refinement.**
The user-facing guide lists the checks ending "…and `worktree-agent-*` leftovers", which now omits both
REQ-083's merge-state split and this REQ's committed-state half. It makes no false claim — it is a
summary, and it never promised fixability semantics — so no reader acts on a wrong contract. Same
finding REQ-083 raised, unchanged in severity; the two are one docs edit. Report-only per the sweep's
severity rule.

**Minor — `worktreeDirtyQueueState` still conflates "clean" with "git failed."** Requirement 7's
degradation contract is met for everything this REQ adds, but the pre-existing porcelain helper returns
`nil` in both cases. Deliberately untouched: Builder Guidance explicitly protects that function, and the
failure is close to unreachable (git enumerated the worktree moments earlier). Filed under
`## Discovered Tasks` as `[low]` rather than fixed, which is the right call for a protected function on
a near-unreachable path.

### Notes on the dimensions

- **Requirements 100%** — all eight traced in Qualification. Requirement 6 is the one worth calling out:
  it demanded confirmation rather than assumption, and it got both a structural argument from REQ-082's
  archived text *and* a fixture that plants a real hand-back file and asserts silence. An argument alone
  would have satisfied the letter of it and been the weaker artifact.
- **Code Quality 96%** — the three-dot choice is the whole fix and it is explained where a later reader
  will meet it, including *why* the tempting two-tree comparison is wrong. `resolveIntegrationBranchRef`
  handles the detached-checkout case rather than assuming a branch exists. Deducted a little for
  `appendWorktreeFindings` continuing to grow: it now runs three sub-probes and is the longest function
  in the file. REQ-072's review already flagged `verify.go` as a split candidate, and both REQs were
  explicitly forbidden from restructuring it — correct for these REQs, but the pressure is real and
  compounding, and the next change here should probably extract rather than add.
- **Test Adequacy 96%** — the strongest part of this REQ is that the two *passing* guards
  (porcelain-survives, stale-snapshot-silent) are only meaningful because the RED test proves the probe
  fires at all; without it, both would pass against a probe that never fires. Deducted because the
  `SkippedProbes` paths added for requirement 7 have no direct assertion — they are reachable (a
  detached `worktree-agent-*` worktree hits the per-branch one) but nothing pins their text.
- **Scope 100%** — `scope-drift.sh` clean. The one file beyond the captured `write_set`
  (`actions/forensics.md`) was declared before coding and logged as D-03. `worktreeDirtyQueueState` was
  verified byte-identical, honoring the explicit protection in Builder Guidance.
- **Risk None** — read-only additions. One extra `git diff` per leftover, and leftovers are typically
  0–3. The probe can only *add* findings, never remove or reclassify existing ones, so no consumer's
  current behavior changes. Both install shapes were reasoned through and neither errors on the other's
  (untracked `do-work/` → empty diff, not a failure). Most importantly the healthy-tree check was run,
  per the prime's own REQ-072 lesson: this repo tracks `do-work/`, so the new probe is *live* here, and
  `verify --repo-root .` reports no findings — the exact failure mode that lesson exists to catch
  (fixtures and code sharing one misreading) does not reproduce.

### Suggested additional testing

- **Assert the skip lines.** The `SkippedProbes` entries added for requirement 7 are the one new path
  without direct coverage; a detached `worktree-agent-*` worktree reaches the per-branch one and could
  pin its wording cheaply.
- **A live fan-out wave.** REQ-085 (queued) runs a real two-builder scenario — the natural place to
  confirm that a legitimate wave, where the orchestrator commits queue state on the integration branch
  while builders work, produces no impersonation findings. The stale-snapshot fixture models that, but
  a real wave is the honest test of it.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- **Scoping the probe to the rule, not to a detection method.** The root cause here was not a missed
  case — it was a probe named after `git status` instead of after "a builder must not write queue
  state." Once the question was re-asked as "what does the *branch* change", three-dot diff semantics
  answered it in one command. Worth asking of any probe: does its name describe the rule or the tool?
- **The two passing tests carry the real weight, but only because a failing one exists.** The
  stale-snapshot and porcelain-survives guards would both pass against a probe that never fires; they
  become meaningful the moment the RED test proves it does. Writing all three together is what makes
  the widening safe rather than merely asserted.
- **Requirement 6 asked for confirmation and the fixture beat the argument.** The structural reasoning
  from REQ-082 was sound and would have been accepted — but planting a real hand-back file in the main
  tree and asserting silence costs four lines and cannot rot the way an argument can.

**What didn't:**
- **The first live reproduction ran against a stale binary and printed nothing.** `queue-kanban` had
  been built during REQ-083, so the manual GREEN check silently exercised the old code and looked like
  the fix had failed. Any manual reproduction in this module must rebuild first — the binary is
  gitignored and there is nothing to remind you it is old.

**Worth knowing:**
- **`git diff A...B` (three dots) is the only comparison that is safe here.** It diffs merge-base(A,B)
  to B, so it reports what the builder's branch changed and is blind to how far the integration branch
  has moved. Two dots or `diff -r` would fire on every run in a repo that commits `do-work/`, because
  the orchestrator claims and archives constantly — which is exactly why the original probe narrowed to
  porcelain instead. The narrowing was a reasonable response to a real problem; it just traded away the
  main case rather than picking a sharper comparison.
- **The integration ref must be named, not passed as `HEAD`.** From the repo root they are the same
  thing, which is what makes the shortcut tempting — but `HEAD` means "whatever checkout this command
  runs in", so inside a worktree the comparison silently becomes branch-against-itself. Same class as
  `git branch -d`'s trap.
- **`appendWorktreeFindings` now runs three sub-probes and is the file's longest function.** Two REQs in
  a row were forbidden from restructuring `verify.go`, correctly — but the next change in this
  neighbourhood should extract rather than append.

## Orientation

**Now `do-work forensics` catches a builder that wrote queue state and committed it** — previously only
the interrupted-mid-write case was detectable, so the probe caught the accident and missed the
behavior. Lives in the queue-kanban board tool's `verify` subsystem
(`tools/queue-kanban/prime-do-kanban.md`), completing REQ-072 requirement 3.

No change to the system's shape: no new module, data flow, or renamed concept. `verify` stays read-only
with the same write surfaces, its report struct and exit-code semantics are untouched, and the rule
being enforced (`actions/work-reference.md` → Worktree Dispatch Mode, "state stays home" / "sole
integrator") is unchanged — only how much of it is mechanically checked. Prime staleness spot-check:
`tools/queue-kanban/prime-do-kanban.md`'s referenced paths all still exist; not made stale.

---
*Source: external audit finding F4 (P2), accepted by `do-work validate-feedback` triage after
empirical reproduction. The audit's characterization of REQ-072 requirement 3 was checked verbatim
against the archived REQ and is accurate.*

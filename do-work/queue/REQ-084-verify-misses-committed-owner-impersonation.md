---
id: REQ-084
title: verify's queue-state probe misses a builder that committed its do-work edits
status: pending
created_at: 2026-08-03T17:09:21Z
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
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---
*Source: external audit finding F4 (P2), accepted by `do-work validate-feedback` triage after
empirical reproduction. The audit's characterization of REQ-072 requirement 3 was checked verbatim
against the archived REQ and is accurate.*

---
id: REQ-083
title: verify reports every builder worktree as a fixable orphan, including active and unmerged ones
status: completed
claimed_at: 2026-08-03T23:41:00Z
completed_at: 2026-08-03T23:52:10Z
commit: f6c1514
created_at: 2026-08-03T17:09:21Z
kb_status: pending
user_request: UR-016
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-072, REQ-084]
batch: audit-remediation-external
addendum_to: REQ-072
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go, actions/forensics.md]
---

# `verify` Reports Every Builder Worktree as a Fixable Orphan, Including Active and Unmerged Ones

## What

`appendWorktreeFindings` emits an `orphan-worktree` finding with `Fixable: true` for every
`worktree-agent-*` worktree and branch it finds, with no merged / unmerged / still-building
distinction. So an in-flight builder is reported as a leftover, and unmerged builder work — which
`actions/cleanup.md` Pass 5 will not delete without asking — is advertised in the report footer as
mechanically fixable.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `tools/queue-kanban/prime-do-kanban.md`, `crew-members/general.md`,
  `coding-guardrails.md`, `testing.md`, and `actions/cleanup.md` → Pass 5. Approach: classify each
  leftover with one extra read command (`git merge-base --is-ancestor <name> HEAD`, run with
  `-C repoRoot`), map the three outcomes onto three categories through a pure routing function, and
  set `Fixable` from that map instead of the hardcoded `true`. The enumeration, the dedupe, the
  `Detail` string, and the dirty-`do-work/` sub-probe are untouched.
- [x] **[APPLY]:** Code written as planned. Three files, all declared in `## Scope`.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file.
  - `tools/queue-kanban/verify.go` — checked the const block renames land in the same aligned block;
    the new `classifyWorktreeMergeState` / `routeWorktreeLeftover` sit next to their only caller; the
    single call site is the only edit inside `appendWorktreeFindings`; no other probe touched. No
    debug artifacts, no `fmt.Println`, no `TODO`.
  - `tools/queue-kanban/verify_test.go` — checked the four new helpers are fixture-only, that both new
    tests skip rather than fail without `git`, and that no existing test was modified. `os/exec` is the
    one added import and it is used.
  - `actions/forensics.md` — checked the single table row, one line, no other row disturbed.
  - Linters: `gofmt -w` (applied), `go vet ./...` (clean), `go test ./...` (pass).

## Why

The `Fixable` flag has an unusually explicit contract, written in the same file:

> `Fixable` means specifically: `do-work cleanup` can mechanically resolve it. Everything else is a
> human decision, and must not be advertised otherwise — an inflated fixable count sends the user to a
> command that will not help. — `tools/queue-kanban/verify.go:53-55`

Unmerged builder work is precisely a human decision — Pass 5 is consent-gated and may hold the only
copy of that work. So the probe violates the contract stated twenty lines above it, and the report's
closing `N fixable: run do-work cleanup` line points the user at a command that will ask instead of
act. The neighbouring `version-changelog-mismatch` probe gets this right for its own expected-transient
state: it carries the caveat in the remedy and is **not** marked fixable.

The in-flight false positive is the milder half but arrives at the worst moment. Fan-out integration is
serial (`actions/work-reference.md` → Fan-Out Dispatch), so during a wave of N builders the Step 9
accelerator (`actions/work.md:606`) runs `verify` while N−1 siblings are still building, and every one
of them is reported as a fixable orphan.

## Context

`tools/queue-kanban/verify.go:354-369` — the append is unconditional:

```go
for _, leftoverName := range orderedNames {
	...
	report.Findings = append(report.Findings, VerifyFinding{
		Category: verifyCategoryOrphanWorktree,
		Detail:   fmt.Sprintf("%s%s exists — %s", ...),
		Fixable:  true,
		Remedy:   "cleanup Pass 5 removes a merged one and asks before discarding an unmerged one; a builder still in flight is not a leftover",
	})
```

The remedy string already knows about all three states. It is prose; nothing acts on it.

**Reproduced** (synthetic repo, one **active** builder worktree created with `git worktree add -b
worktree-agent-REQ-001-active`, nothing merged, nothing dirty):

```
! orphan-worktree [fixable]: worktree-agent-REQ-001-active exists — <path>
    → cleanup Pass 5 removes a merged one and asks before discarding an unmerged one; a builder still in flight is not a leftover
1 fixable: run do-work cleanup
verify_exit=1
```

**Merged-ness is cheaply available.** `git branch --merged <integration-branch>` (or
`git merge-base --is-ancestor <branch> HEAD`) answers it with one read command, in the same
`exec.Command("git", "-C", …)` shape the file already uses throughout.

**What verify genuinely cannot know** is whether a builder is *still running* — there is no lock, no
heartbeat, and REQ-073 forbids adding one. That is a real boundary, and the fix must report honestly
within it rather than inventing a liveness signal.

**Why no check caught it.** REQ-072's own review recorded the gap: "the worktree probes skip when git
cannot answer, so `listWorktreeAgentWorktrees` / `worktreeDirtyQueueState` are exercised on this live
repo (correctly finding nothing) but have no fixture that plants a `worktree-agent-*` leftover — a
temp-repo test with a real `git worktree add` would close it, and is noted under suggested testing
rather than built." That test is now required, not suggested.

## Detailed Requirements

1. **Set `Fixable: true` only for a branch merged into the integration branch.** Everything else is a
   human decision and must report as not-fixable, per the field's own doc comment.
2. **Split the finding by state**, so the report says which case it found rather than leaving the
   reader to decode one category from a remedy sentence. At minimum: merged residue (fixable) versus
   unmerged (not fixable). Keep the category strings greppable and name them in the constant block at
   `verify.go:16-26` alongside the existing ones.
3. **Be honest about liveness.** `verify` cannot distinguish "builder still running" from "builder died
   and left this behind" — no lock exists and none may be added. State that boundary in the finding's
   remedy and in the doc comment. **Do not** add a heartbeat, a lock file, a PID check, an mtime
   heuristic, or a claim registry to guess it.
4. **Decide and state whether an unmerged leftover should be a finding at all during a live run**, and
   record the reasoning where a later reader will find it. Options that fit the existing design: keep
   it as a non-fixable finding (matching how `version-changelog-mismatch` handles its expected
   mid-release state), or demote it to a reported-but-non-failing note. Do not silently suppress it —
   `VerifyReport`'s doc comment is explicit that silence reads as "checked and clean."
5. **Resolve `--repo-root` to the right integration branch.** Merged-ness is HEAD-relative: `git branch
   -d`'s own trap, already documented at `actions/work-reference.md:294` ("`-d` tests merged-ness
   against the current HEAD … from anywhere else a perfectly-merged branch refuses and an unmerged one
   can pass"). Compare against the repo-root checkout's HEAD, not the worktree's, and say so in a
   comment.
6. **Add the temp-repo fixture REQ-072 deferred.** A real `git worktree add` in a `t.TempDir()` repo,
   covering: merged branch with worktree, merged branch without worktree, unmerged branch with
   worktree, unmerged branch without worktree, and a non-`worktree-agent-*` worktree that must be
   ignored entirely. Skip the test cleanly when `git` is unavailable, matching how the probe itself
   degrades.
7. **Assert the `Fixable` contract, not just the finding.** A test that pins `FixableCount()` for an
   unmerged leftover at zero is what stops this regressing — the finding's presence is the easy half.
8. **Keep `verify` read-only.** It reports and routes; it repairs nothing. Adding a merged-ness query
   adds read commands only.

## Constraints

- **Never `-D`, never `--force`, and never delete anything.** `verify` writes nothing at all
  (`tools/queue-kanban/verify.go:92-96`, and `CLAUDE.md` → Shipped Tooling's two-write-surface
  sentence). This REQ changes classification only.
- **The 3-hour staleness constant is not the model to copy here.** `staleClaimThreshold` bounds how
  long a dead *claim* goes unnoticed and explicitly authorizes nothing; do not reach for a time
  threshold as a liveness proxy for worktrees.
- **A developer's own worktrees stay untouched and unreported.** Only names carrying
  `worktreeAgentNamePrefix` are in scope (`verify.go:36-40`) — that boundary is unchanged.
- **Do not renumber or restructure the other seven probes.** `verify.go` is already the largest file in
  the module and REQ-072's review flagged it as a split candidate at twelve probes; this REQ adds
  classification inside one probe, not a reorganization.
- Wired consumers must keep working: `actions/forensics.md` Check 14 and the Step 9 accelerator both
  read this report. If a category string changes, grep both.

## Dependencies

`addendum_to: REQ-072`, which introduced the probe. `related: REQ-084` — same file, same function
neighbourhood, deliberately no `depends_on` (neither needs the other's output). Whichever lands second
re-reads `appendWorktreeFindings` before editing; the merge is the non-interference proof, not the
overlaps badge.

## Builder Guidance

**Certainty: Firm on requirements 1, 3, 6 and 7; open on requirements 2 and 4.** The `Fixable`
correction is not a judgment call — the field's doc comment settles it. How many categories the split
produces, and whether an unmerged leftover should fail the exit code during a live run, are genuine
design choices; log whichever you pick as a D-XX with reasoning, and pick the shape that keeps the
report readable rather than the one with the most states.

Read `actions/cleanup.md` → Pass 5 before deciding what "fixable" means here. The routing has to match
what that pass will actually do, which is the whole point of the flag.

## Red-Green Proof

**RED prompt/case:** A Go test that builds a temp repo, runs `git worktree add -b
worktree-agent-REQ-001-unmerged`, commits something on that branch so it is genuinely unmerged, then
calls `runVerifyProbes` and asserts `report.FixableCount() == 0`.

**Why RED now:** `verify.go:367` sets `Fixable: true` unconditionally, so the count is 1 and the
assertion fails. Reproduced by hand already — a synthetic repo with one active builder worktree prints
`orphan-worktree [fixable]` and `1 fixable: run do-work cleanup`.

**GREEN when:** That assertion passes; a *merged* leftover in the same fixture still reports as
fixable (so the fix is a classification, not a blanket demotion); the full fixture table from
requirement 6 passes; `go test ./...` in `tools/queue-kanban/` stays green; and re-running the manual
probe on a repo with one active builder no longer prints `N fixable: run do-work cleanup`.

**Validation:** Inferred during capture, then reproduced — the RED output above was executed against a
synthetic repo before this REQ was written.

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the
external audit, and the batch constraints.

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" and the "where" are both pinned by the REQ (`appendWorktreeFindings`,
`verify.go:354-369`), but the classification has to match what `actions/cleanup.md` Pass 5 actually
does and fit the file's existing probe/finding shape — that is discovery, not planning. No
architectural decision is open; requirements 2 and 4 are local design choices inside one function.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Spec loaded:** `specs/bug-fix.md` (`suggested_spec: bug-fix`) — reproduce, isolate, minimal fix,
regression test, root-cause note, then grep for the same pattern elsewhere.

**The probe as it stands** (`tools/queue-kanban/verify.go:328-383`). `appendWorktreeFindings`
enumerates `worktree-agent-*` worktrees (`listWorktreeAgentWorktrees`, keyed by directory basename)
and branches (`listWorktreeAgentBranches`), unions them into `orderedNames`, and emits one
`verifyCategoryOrphanWorktree` finding per name with `Fixable: true` hardcoded at line 367. The
three-state truth lives only in the `Remedy` prose at line 368, which nothing reads.

**What `Fixable` is contractually allowed to mean** — `verify.go:53-55`: "`do-work cleanup` can
mechanically resolve it." Read against `actions/cleanup.md` → Pass 5, that is exactly one of the
three states:

| Leftover state | Pass 5 behavior (`actions/cleanup.md:130-131`) | Fixable? |
| --- | --- | --- |
| Branch merged into the integration branch | `git worktree remove` + `git branch -d` both succeed — "pure residue" | yes |
| Unmerged / dirty | Refusal; Pass 5 "stops being mechanical", loads `clear-questions.md` and asks | no |
| Merge state unanswerable | Pass 5's "If you cannot establish which branch that is, delete nothing" | no |

**Merged-ness is one read command** in the shape the file already uses eight times:
`exec.Command("git", "-C", repoRoot, "merge-base", "--is-ancestor", <name>, "HEAD")` — exit 0 merged,
exit 1 unmerged, anything else (missing ref, detached worktree whose branch is gone) unanswerable.
Pinning it with `-C repoRoot` is requirement 5: `actions/work-reference.md` → Worktree Dispatch Mode,
*Cleanup — happy path* documents the same HEAD-relative trap for `git branch -d`.

**No consumer greps the category string.** `grep -rn 'orphan-worktree'` outside `do-work/` returns
exactly one hit — the constant's own definition. `actions/forensics.md:203` and `actions/work.md:607`
run the binary and read the rendered report; Check 14's table row describes the probe in prose
(`worktree-agent-* leftovers, and any such worktree with uncommitted do-work/ changes`) and cites
Pass 5 as the definition it reuses. That row is the one restatement this change makes imprecise.

**Test shape.** `verify_test.go` seeds plain temp directories, so today every worktree probe is
skipped there (`git worktree list` fails outside a repo) — which is precisely the fixture hole
REQ-072's review recorded. A real `git init` + `git worktree add` fixture is new ground for this
file; `findingsMentioning` and `snapshotTreeContents` are the reusable helpers.

**No liveness signal exists and none may be added** (requirement 3, and REQ-073). `staleClaimThreshold`
is explicitly not a liveness test, and its own doc comment says so.

## Scope

**Files I will touch:**
- `tools/queue-kanban/verify.go` (modify) — split `appendWorktreeFindings`'s single finding by merge
  state; add the merge-state classifier and its category constants
- `tools/queue-kanban/verify_test.go` (modify) — the temp-repo `git worktree` fixture table and the
  `FixableCount() == 0` contract assertion
- `actions/forensics.md` (modify) — Check 14's probe-table row, one line, so the documented probe
  matches the split it now reports (declared up front rather than taken as drift; see D-02)

**Files I will NOT touch:**
- `actions/cleanup.md` — Pass 5 is the source this change conforms *to*; its behavior is unchanged
- The other seven probes in `verify.go`, and any file-splitting of it (explicit constraint)
- `crew-members/`, `actions/work.md`, `actions/work-reference.md` — the worktree contract they state
  is unchanged; only verify's classification of it changes

**Acceptance criteria (restated from REQ):**
- [ ] `Fixable: true` is set only for a branch merged into the integration branch (req 1)
- [ ] The finding is split by state with greppable category constants in the `verify.go:16-26` block (req 2)
- [ ] The liveness boundary is stated in the remedy and in a doc comment; no heartbeat, lock, PID
      check, mtime heuristic, or claim registry is added (req 3)
- [ ] The unmerged-during-a-live-run question is decided and the reasoning recorded (req 4)
- [ ] Merged-ness is resolved against the repo-root checkout's HEAD, with a comment saying why (req 5)
- [ ] A temp-repo fixture covers all five cases from req 6 and skips cleanly without `git`
- [ ] A test pins `FixableCount()` at 0 for an unmerged leftover (req 7)
- [ ] `verify` still writes nothing (req 8)

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/verify.go` (modified)
- `tools/queue-kanban/verify_test.go` (modified)
- `actions/forensics.md` (modified)

**What was done:** `appendWorktreeFindings` now asks `git merge-base --is-ancestor <name> HEAD`
(pinned to the repo root with `-C`) for each `worktree-agent-*` leftover and routes the answer through
`routeWorktreeLeftover`, which sets `Fixable: true` for merged residue only. `verifyCategoryOrphanWorktree`
is replaced by three categories — `merged-worktree-leftover`, `unmerged-worktree-leftover`,
`worktree-merge-state-undetermined` — so the report states the case it found instead of leaving it in
remedy prose. Tests gained a real `git init` + `git worktree add` fixture covering all five states from
requirement 6 plus a `FixableCount() == 0` contract assertion. Check 14's probe row in
`actions/forensics.md` was updated to describe the split.

**Root cause (per `specs/bug-fix.md`):** the probe was written to *enumerate* leftovers, and `Fixable`
was set once for the whole enumeration. The three-state distinction existed from the start — it was
written into the `Remedy` sentence — but as prose, where no code could act on it and no test could
assert it. The defect is a fact encoded in the wrong medium, not a missed case.

## Decisions

- **D-01 (DECIDE & STATE)** — *Requirement 4: an unmerged leftover stays a reported, non-fixable
  finding during a live run; it is not suppressed and not demoted to a non-failing note.* Reasoning:
  `VerifyReport`'s own doc comment says silence reads as "checked and clean," and verify has no way to
  know a run is active (that is requirement 3's boundary), so any suppression rule would have to guess
  — and every wrong guess hides genuinely stranded builder work, which is the expensive direction of
  the error. The precedent in the same file is `version-changelog-mismatch`: an expected transient
  state, reported anyway, with the transient case named in the remedy and `Fixable` left false. The
  reasoning is recorded in `routeWorktreeLeftover`'s doc comment, where a later reader meets it.
- **D-02 (DECIDE & STATE)** — *Extend scope to `actions/forensics.md`.* Check 14's probe table
  describes this probe in prose and cites Pass 5 as the definition it reuses; splitting the finding
  makes that row imprecise the moment it lands. Declared in `## Scope` and mirrored to `write_set`
  before coding rather than taken as drift afterwards. One row, no other check touched.
- **D-03 (DECIDE & STATE)** — *Three categories, and retire the `orphan-worktree` string rather than
  reusing it for the merged half.* Keeping the old name for only one of the three states would
  silently change what an existing name means, which is worse for a grep-oriented consumer than a
  clean rename. `grep -rn 'orphan-worktree'` outside `do-work/` returned only the constant's own
  definition, so nothing downstream matches on it. The third state (git declining to answer) is kept
  separate from "unmerged" because reporting an unanswered question as a definite answer is the same
  class of defect this REQ fixes; both route to not-fixable, so the extra category costs a line and
  buys an honest report.

## Discovered Tasks

- **[normal]** `_dev/tests/update-script-behavior.sh` has 7 failing probes on the base branch, all
  covering `tools/do-work-update.sh`: a mid-update failure exits 0 instead of non-zero and prints none
  of the four documented recovery lines (`Update did not complete`, `may be partially updated`, the
  `git checkout --` restore command, the `git clean -nd --` cleanup command), and the dirty-install
  path neither exits non-zero nor prints its `restores the COMMITTED content` warning. Verified
  pre-existing: the same 8 FAIL lines appear with this REQ's changes stashed. Not touched here — it is
  a different subsystem and unrelated to this REQ. Worth deciding whether the updater regressed or the
  probes were written ahead of it.

## Qualification

Passed — 3 files verified, 8 requirements traced, P-A-U confirmed.

- `tools/checks/qualify.sh` → OK (mechanical checks 1, 4, and the grep half of 5).
- `tools/checks/scope-drift.sh` → OK: Implementation Summary matches the Scope declaration exactly.
  Both set-differences empty — `actions/forensics.md` was declared before coding, not after.
- **Substantive (check 2):** `verify.go` +102/-~14 is logic, not whitespace — a new type, two new
  functions, and a changed call site. `verify_test.go` +161 is a real git fixture, not a stub.
- **Requirements traced (check 3):** 1 → `routeWorktreeLeftover`'s `Fixable` is true only in the
  `worktreeMergeStateMerged` arm. 2 → three constants in the `verify.go:16-28` block. 3 →
  `classifyWorktreeMergeState`'s doc comment plus the unmerged remedy; the negative half verified by
  grepping the added lines for `ModTime|Stat(|Pid|time.Now|time.Since|lock|heartbeat` — the only hit
  is the doc comment saying no such thing exists. 4 → D-01 and `routeWorktreeLeftover`'s doc comment.
  5 → `exec.Command("git", "-C", repoRoot, …)` with the comment citing `git branch -d`'s HEAD-relative
  trap. 6 → `TestVerifyClassifiesWorktreeLeftoversByMergeState`, five leftovers, `t.Skip` without git.
  7 → `TestVerifyDoesNotAdvertiseAnUnmergedWorktreeAsFixable` pins `FixableCount()` at 0. 8 →
  `TestVerifyWritesNothing` still passes; only read commands were added, and the grep above confirms
  no write call.
- **Flowing (check 6):** not a hollow implementation — `classifyWorktreeMergeState` actually shells out
  and its three branches were each observed on a live synthetic repo, not only in the fixture.

## Testing

**Tests run:** `cd tools/queue-kanban && go test ./...` (+ `go vet ./...`, `gofmt`)
**Result:** ✓ All passing — 14/14 in the verify suite, whole module green

**Red-green validation:** traced to `## Red-Green Proof`, whose RED case was implemented verbatim
(temp repo → `git worktree add -b worktree-agent-REQ-001-unmerged` → commit on the branch →
`runVerifyProbes` → assert `FixableCount() == 0`).

- `TestVerifyDoesNotAdvertiseAnUnmergedWorktreeAsFixable`: ✗ before implementation
  (`FixableCount = 1, want 0`, report showing `orphan [fixable]` and `1 fixable: run do-work cleanup`)
  → ✓ after
- `TestVerifyClassifiesWorktreeLeftoversByMergeState`: ✗ before implementation (all four leftovers
  classified unmerged; `FixableCount = 4, want 2`) → ✓ after

The RED was **observed, not assumed**: the category constants were added first with the old
unconditional `Fixable: true` left in place, so the tests compiled and failed on their assertions
rather than on a build error.

**New tests added:**
- `TestVerifyDoesNotAdvertiseAnUnmergedWorktreeAsFixable` — the `Fixable` contract assertion (req 7)
- `TestVerifyClassifiesWorktreeLeftoversByMergeState` — the five-state fixture table (req 6)
- `TestVerifyReportsAnUndeterminedMergeStateSeparately` — the third state's own coverage, added during
  Step 7 after the review flagged it as an untested new code path (see `## Review`)
- Fixture helpers `runGitInFixture`, `newWorktreeFixtureRepo`, `addFixtureWorktree`,
  `commitFixtureWork` — the first real git repo in `verify_test.go`, closing the hole REQ-072's review
  recorded

**Manual acceptance (the REQ's GREEN condition):** re-ran the hand reproduction on a synthetic repo.
An unmerged builder branch now prints `unmerged-worktree-leftover` with no `[fixable]` marker and no
`N fixable: run do-work cleanup` line; adding a merged leftover to the same repo prints
`merged-worktree-leftover [fixable]` and `1 fixable: run do-work cleanup`. Both directions observed.

**Suite-level check:** `_dev/tests/contract-regressions.sh` reports 7 failing update-script probes.
Confirmed **pre-existing, not a regression** — the identical 8 FAIL lines appear with this REQ's three
files stashed. Filed under `## Discovered Tasks`. The Step 5.75 baseline (`do-work/working/baseline.json`,
`launched: true`, `exit_status: 0`) covered the Go suite, which is green before and after.

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-04T00:05:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Acceptance:** Pass — both directions observed on a live synthetic repo: an unmerged leftover prints
no `[fixable]` marker and no `N fixable` line; a merged one in the same repo prints both.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Restatement Sweep (MUST — run):** The diff redefines two things stated in more than one place: the
probe's category vocabulary (`orphan-worktree` retired for three names) and what `[fixable]` covers for
a worktree leftover.

- *Category strings:* `grep -rn 'orphan-worktree'` outside `do-work/` returns only the constant's own
  definition. No prose, test, template, or tool matches on it — nothing stale.
- *`[fixable]` semantics:* `actions/forensics.md:192` states the marker generically ("the ones
  `do-work cleanup` can mechanically resolve") — still exactly true, and truer than before.
  `actions/cleanup.md` Pass 5 is the source this change conforms to, unchanged. `docs/cleanup-guide.md:25`
  already described the merged/unmerged split correctly and is now consistent with verify rather than
  ahead of it. `CLAUDE.md:163`'s "reports and routes; it never repairs" holds.
  `actions/work.md:607`'s accelerator text claims only the version/title checks — unaffected.
- *One stale restatement found* — `docs/forensics-guide.md:23` (Minor, below).

### Findings

**Minor — `docs/forensics-guide.md:23` still lists the probe undifferentiated.** The user-facing guide
enumerates verify's checks ending "…and `worktree-agent-*` leftovers", which is now the only place
describing the probe without its merge-state split — its action-file twin (`actions/forensics.md:203`)
was updated in this REQ. Still literally true and it makes no fixability claim, so no reader acts on a
wrong contract; it is imprecision, not misdirection. Report-only per the sweep's severity rule (Minor
findings do not earn a follow-up REQ). Worth folding into the next docs pass over that guide.

**Minor — the third state shipped untested, and this review caught it.** `routeWorktreeLeftover`'s
`worktreeMergeStateUndetermined` arm had no coverage: requirement 6's five-case table lists only merged
and unmerged shapes, so the fixture satisfied the requirement while leaving the new branch unexercised.
Fixed in place rather than deferred — `TestVerifyReportsAnUndeterminedMergeStateSeparately` plants a
`git worktree add --detach` checkout whose name matches the convention but has no branch behind it, and
pins both the category and `FixableCount() == 0`. An 8-line test inside an already-declared file is not
worth a follow-up REQ. Recorded here because a self-caught gap that leaves no trace reads, later, like
a gap nobody looked for.

### Notes on the dimensions

- **Requirements 100%** — all eight traced to specific code in Qualification above, including the two
  negative requirements (no liveness heuristic; verify stays read-only), which were checked by grepping
  the added lines rather than by assertion.
- **Code Quality 95%** — matches the file's conventions: two-word names, doc comments that explain
  *why* (the HEAD-relative trap, the liveness boundary, the D-01 reasoning) rather than what, and the
  same `exec.Command("git", "-C", …)` shape used eight times already. The exit-code discrimination is
  correct and load-bearing: exit 1 is git's answer, anything else is git declining, and conflating them
  would recreate this REQ's defect one level down. Deducted for `runError.(*exec.ExitError)` being a
  bare type assertion where `errors.As` is the modern idiom — the file wraps no errors, so it is safe
  today, and matching local style over introducing a new import is the right call; noting it as taste,
  not a defect.
- **Test Adequacy 95%** — red-green honored and *observed*, not asserted: the constants landed first so
  the RED was an assertion failure, not a build error. All three states now covered, plus the
  ignore-a-developer's-worktree case. Deducted because the fixture cannot cover the case that matters
  most operationally — a builder genuinely mid-flight — which is a real boundary, not an omission
  (requirement 3).
- **Scope 100%** — `scope-drift.sh` reports both set-differences empty. The one file beyond the
  captured `write_set` (`actions/forensics.md`) was declared in `## Scope` and mirrored into
  `write_set` before any code was written, and logged as D-02.
- **Risk None** — read-only additions to a read-only tool. One extra `git merge-base` invocation per
  leftover; leftovers are typically 0–3, so the cost is unmeasurable. No caller's interface changed:
  `VerifyFinding`, `VerifyReport`, `FixableCount`, and `ExitCode` are untouched, and both wired
  consumers read rendered text they still parse identically. The change can only *lower* the fixable
  count, so no consumer is newly pointed at a command that will not help — the failure direction this
  REQ existed to close.

### Suggested additional testing

- **A real fan-out wave.** The in-flight false positive arrives when Step 9's accelerator runs verify
  while sibling builders are still building. REQ-085 (queued) runs exactly that scenario and is the
  natural place to confirm the report now reads correctly mid-wave.
- **A consumer repo on an older tarball.** Check 14 and the Step 9 accelerator report verify's output
  verbatim; worth one look at a real consumer install to confirm the new category strings read clearly
  to someone who has not read this REQ.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- **Landing the constants before the logic bought a real RED.** Adding the three category names first,
  with the old unconditional `Fixable: true` still in place, turned a compile error into an assertion
  failure — `FixableCount = 1, want 0` against a rendered report showing `[fixable]` and
  `1 fixable: run do-work cleanup`. A build error proves the test is new; only an assertion failure
  proves the bug is real. This is the session's fourth consecutive case of observing the RED paying
  for itself in under a minute.
- **Reading `actions/cleanup.md` Pass 5 first turned a design question into a lookup.** "What should
  fixable mean here" reads open-ended until you notice Pass 5 has exactly three branches — mechanical,
  consent-gated, and refuse-to-act-without-a-merge-target. The three categories fell out of that table
  rather than being invented; the REQ's Builder Guidance pointing at Pass 5 was the highest-value line
  in the file.

**What didn't:**
- **The requirement list was a floor, exactly as this session's checkpoint warned.** Requirement 6
  enumerates five fixture cases, none of which exercises the third state the split introduces — so a
  fixture that satisfied the requirement in full still left a new code path untested. The review caught
  it; the requirement never would have. A REQ that says "split into N states" and separately lists
  fixtures is two lists that can disagree, and here they did.
- **The first live reproduction misled for a minute.** Creating the worktree at `../wt-active` while
  naming its branch `worktree-agent-REQ-001-active` printed `branch only (no worktree)` and looked like
  a detection bug. It was the fixture: worktree dispatch mode keeps the directory basename and the
  branch name *the same string*, and `listWorktreeAgentWorktrees` keys on the basename. Any manual
  worktree reproduction has to honor that or it tests a shape the skill never produces.

**Worth knowing:**
- **`git merge-base --is-ancestor` has three outcomes, not two, and the third one matters.** Exit 0 is
  merged, exit 1 is git's answer "no", and anything else (usually 128) is git *declining to answer* —
  most often a worktree whose branch is gone. Folding 128 into "unmerged" would report an unanswered
  question as a definite answer, which is the same defect class this REQ fixed, one level down.
- **Merged-ness is HEAD-relative and the question must be pinned to the repo root.** `git -C repoRoot`
  is doing real work here, not tidiness: asked from a builder's worktree, a merged branch reads
  unmerged and vice versa. Same trap `git branch -d` carries in `actions/work-reference.md` →
  Worktree Dispatch Mode, *Cleanup — happy path*.
- **`verify_test.go` now has a real git-repo fixture** (`newWorktreeFixtureRepo` and friends). Before
  this REQ every test used a plain temp directory, which silently *skipped* all worktree probes — the
  coverage hole REQ-072's review predicted. Any future probe that shells out to git should build on
  these helpers rather than adding a second fixture style.
- **`_dev/tests/contract-regressions.sh` is red on the base branch** — 7 update-script probes. Confirm
  pre-existing by stashing before treating any of it as your own regression. Filed as a Discovered Task.

## Orientation

**Now `do-work forensics` and the Step 9 release accelerator tell you which builder leftovers you can
actually clean up** — merged residue is marked `[fixable]` and counted toward
`N fixable: run do-work cleanup`, while unmerged work (possibly a builder still running) and
unanswerable cases are reported and left to you. Lives in the queue-kanban board tool's `verify`
subsystem (`tools/queue-kanban/prime-do-kanban.md`).

The system's shape is unchanged — no new module, data flow, or contract. `verify` remains read-only
with the same two write surfaces the skill documents, and its report struct and exit-code semantics are
untouched; only one probe's classification changed. Prime staleness spot-check:
`tools/queue-kanban/prime-do-kanban.md`'s referenced paths all still exist, and its testing guidance
still matches (`go test ./...` in the module). Not made stale by this change.

---
*Source: external audit finding F3 (P1). Accepted by `do-work validate-feedback` triage with its
severity corrected: the audit claimed Step 9's verifier would fail the run, but `actions/work.md:606`
never gates on verify's exit code. The defect is the `Fixable` contract violation and the missing
classification.*

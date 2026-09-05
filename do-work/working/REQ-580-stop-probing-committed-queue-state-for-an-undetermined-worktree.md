---
id: REQ-580
title: 'Stop probing committed queue state for a worktree whose merge state is already undetermined'
status: claimed
created_at: 2026-09-05T00:19:58Z
user_request: UR-118
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-579]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
claimed_at: 2026-09-05T00:38:08Z
route: A
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-09-05T00:48:26Z
  basis:
    - Route A
    - 2-file write set
    - 4 acceptance criteria
---

# Stop Probing Committed Queue State for a Worktree Whose Merge State Is Already Undetermined

## What

In `appendWorktreeFindings` (verify.go), a leftover worktree whose branch is gone produces two messages: a `worktree-merge-state-undetermined` finding, and a skipped-probe line saying the committed-queue-state probe failed with "no such branch". Both come from one fact. When the merge-state classification is already undetermined for a leftover, skip the committed-queue-state probe for that leftover and say inside the undetermined finding that the queue-state check could not run either. One row per worktree, and the "not checked" fact is still reported.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, the board prime and satellite, coding-guardrails and testing. Traced `appendWorktreeFindings`, `routeWorktreeLeftover`, `classifyWorktreeMergeState` and the undetermined arm of `classifyWorktreeLeftover`. Plan was one extra term on the existing `if` plus an extended remedy string.
- [x] **[APPLY]:** Exactly the two declared files changed; `git status --short` in the builder worktree printed only those two as modified.
- [x] **[UNIFY]:** `git diff --stat` = 2 files, 63 insertions, 2 deletions. `gofmt -l .` and `go vet ./...` both silent; `go test -count=1 ./...` returned ok in 87.660s. Grepped the repo for other holders of the old remedy sentence: no shipped test, script or contract asserts on the changed text.

## Why

The user asked whether the card and the line under it were the same information. They are two checks with one root cause, so the reader sees the same worktree name twice and has to work out that nothing new is being said. A good layout (REQ-579, render verify findings as compact rows) should not have to hide a redundant producer.

## Detailed Requirements

- The committed-queue-state probe (`worktreeCommittedQueueState`) runs only when `classifyWorktreeLeftover` did not return `worktreeLeftoverMergeStateUndetermined` for that leftover. The other dispositions keep the probe exactly as today.
- The undetermined finding's remedy gains a clause stating that, because git could not resolve the branch, whether the builder committed queue state under `do-work/` could not be checked either. "Unknown never reads as clean" holds: the fact is moved into the finding, not dropped.
- No new `SkippedProbes` entry is written for that leftover's committed-queue-state probe. Skipped entries for other reasons (git missing, integration ref unresolvable) are unchanged.
- The text renderer (`renderVerifyReport`) and the board payload need no change; they print what the report carries.

## Red-Green Proof
**RED prompt/case:** Extend `TestVerifyReportsAnUndeterminedMergeStateSeparately` (or add a sibling) on the same detached-worktree fixture: assert `report.SkippedProbes` contains no entry mentioning `committed-queue-state probe for worktree-agent-REQ-005-detached`, and that the single undetermined finding's remedy mentions that committed queue state could not be checked.
**Why RED now:** `appendWorktreeFindings` runs `worktreeCommittedQueueState` for every leftover regardless of disposition; on a detached or branchless worktree `git diff <ref>...<name>` fails and the failure is appended to `SkippedProbes`, so the fixture today yields one finding plus one skipped line for the same name.
**GREEN when:** The fixture yields exactly one undetermined finding for the name, zero skipped-probe lines for it, the finding is still not `Fixable`, and every other test in `verify_test.go` (in particular the committed-queue-state and uncommitted-queue-state cases around line 1087 and 1118) still passes.
**Validation:** Inferred during capture from the user's question ("are these finding duplicated?") and their "ok, do it" to the proposal that named this as D4.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban model" and its family `unknown-reads-as-clean` is the exact rule this REQ must keep. Over the 2000-token budget on its own; the rule is restated in Detailed Requirements instead.

*Source: "are these finding duplicated? the big box and the small line is the same info?" / "ok, do it"*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact function (`appendWorktreeFindings`), the exact condition to add, the exact remedy string to extend, and its own two-file write set. Its Detailed Requirements already state what must not change (the renderer and the board payload). Nothing needed discovering, so exploration and planning would only restate the request.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)

**What was done:** appendWorktreeFindings no longer runs the committed-queue-state probe for a leftover whose merge state is already undetermined, and the undetermined finding's remedy now says the same unresolved branch stopped that check. One fact produces one row, and the unchecked part is reported rather than dropped.

**Per-file detail.** In verify.go the probe condition gained one term, so it now requires both a resolvable integration ref and a disposition other than the merge-state-undetermined one; routeWorktreeLeftover's undetermined remedy gained the "check could not run" clause. Eleven changed lines, seven of them an explanatory comment. In verify_test.go one new test function was added, 54 lines, and no existing test was edited.

**Implementation range:** `7bbbc325..d89efc0b`. Builder commit `bd29594dca16ff9ce869f88e16bf7343616a6e4f`.

**Exact new remedy text** (only the sentence after "merge target for" is new):

> git could not say whether this is merged (typically a worktree whose branch is gone) — inspect it by hand; cleanup Pass 5 deletes nothing it cannot establish a merge target for. The same unresolved branch stopped the committed-queue-state check, so whether a builder committed queue state under do-work/ on it is unknown, not clean

The finding's Category, Detail and non-fixable flag are unchanged.

## Decisions

- **D-01 — gate on the disposition, not a re-read.** `appendWorktreeFindings` already holds `disposition` from `classifyWorktreeLeftover`, so the guard is one extra term on the existing condition. No new call to `classifyWorktreeMergeState`, no new helper.
- **D-02 — the "not checked" fact lives in the remedy string, not a new finding field.** The REQ asks for one row, and says the renderer and board payload need no change. A new struct field would have forced both. The remedy is the surface a reader already reads for what to do next.
- **D-03 — "unresolved branch", not "missing branch".** `classifyWorktreeMergeState` returns undetermined for any git exit code other than 1, which is most often a gone branch but is not guaranteed to be one. The new clause matches the existing sentence's hedge instead of over-claiming a cause.
- **D-04 — the clause ends "unknown, not clean".** The report's contract is that silence reads as checked-and-clean, so the moved fact has to name the unknown explicitly. The test asserts on the word `unknown`, pinning the property rather than the whole sentence.
- **D-05 — a new sibling test rather than extending `TestVerifyReportsAnUndeterminedMergeStateSeparately`.** The REQ allowed either. The existing test pins a different claim (the classification is separate and not fixable); merging the de-duplication claim into it would make one failure mean two things.
- **D-06 — the `SkippedProbes` assertion matches the leftover name, not the exact probe sentence.** Asserting a full error string is brittle. Matching the name proves the stronger property anyway: no skipped-probe line for this worktree from any cause.

## Qualification

Passed the request-bound advance qualify gate for `7bbbc325..d89efc0b` against the merged range. Exactly the two declared files changed, 63 insertions and 2 deletions, no undeclared touch and no path under `do-work/` on the builder branch. The change is one added term on an existing condition plus one extended string; no debug artifact, no assertion deleted, no suppression, no widened scope. The P-A-U boxes are supported by the evidence quoted in the hand-back rather than merely checked.

## Testing

**Red-green validation:** The test was written first and failed on three real assertions before any change to `verify.go`. The first is the duplicate line the REQ targets: `committed-queue-state probe for worktree-agent-REQ-005-detached: \`git diff master...worktree-agent-REQ-005-detached -- do-work/\` failed (no such branch, or unrelated histories)`. The other two are the remedy not naming the check that could not run, and the unrun check not reading as unknown. No compile or import error — three genuine assertion failures. After the fix all three pass.

The "one row per worktree" count assertion already held before the fix, which is precisely the defect: the duplicate lived in the `SkippedProbes` list, not in the findings list.

**Controls preserved:** `TestVerifyFlagsQueueStateCommittedOnABuilderBranch` is unchanged and still passes — its fixture `worktree-agent-REQ-010-impersonator` has a real branch, so its merge state is determinable, the probe still runs, and it still reports the forged `REQ-999-forged.md`. The new test's doc comment names it so the pair cannot drift. `TestVerifyStillFlagsUncommittedQueueStateInAWorktree` and `TestVerifyIgnoresAStaleQueueSnapshotAndTheHandBackFile` also still pass.

**Module verification** from `skills/do-work-board/tools/queue-kanban`: `go test -count=1 ./...` returned `ok ... 87.660s`; `gofmt -l .` and `go vet ./...` both printed nothing.

## Discovered Tasks

- **T-01 — this change is release-shaped and needs the release payload.** `_dev/primes/prime-releases.md` treats a commit under `skills/` as a release, and `CHANGELOG.md`, its mirror at `skills/do-work/CHANGELOG.md` and the version files are owner-only writes outside this REQ's write set. The user-visible wording change is the remedy text quoted above. → owned by this REQ's finalization, not a follow-up.
- **T-02 — REQ-579's screenshot description goes stale.** `do-work/queue/REQ-579-render-verify-findings-and-skipped-probes-as-compact-rows.md:66` describes a captured board screenshot whose "2 probe(s) could not run" disclosure is exactly the two committed-queue-state lines this REQ removes. The screenshot still documents the board as of 2026-09-05, but REQ-579 must not expect to reproduce those bullets from an undetermined worktree. → report only.
- **T-03 — no sibling defect in the same function.** The other two skipped-probe producers in `appendWorktreeFindings` (git missing, integration ref unresolvable) are per-run rather than per-leftover, so they cannot duplicate a finding the same way. → report only.

## Review

**Overall: 95%**
**Acceptance: Pass.** Independent review approved the change and reproduced every claim itself rather than reading the hand-back. Six findings, all report-only, none Important.

### Independent Review Record

**Approve** — the redundant row is gone, the "not checked" fact still reaches the reader in both the CLI and the board, and RED, GREEN and the control were all reproduced by the reviewer.

Route A | Merge range `7bbbc325..d89efc0b` | Builder `bd29594d`

#### What's built

`appendWorktreeFindings` (verify.go:1141) runs the committed-queue-state probe only when the integration ref resolved and the leftover's disposition is not the merge-state-undetermined one. A detached or branchless leftover now produces one finding row and zero skipped-probe lines. The undetermined remedy (verify.go:1059) gained the clause naming the check that could not run, ending "unknown, not clean". Two files, 63 insertions, 2 deletions, exactly the declared write set.

#### Acceptance testing the reviewer performed

- **R1 — RED reproduced without touching the repository.** Exported the module at `d89efc0b` with `git archive` into a scratch directory, overwrote only `verify.go` with the pre-change bytes from `7bbbc325`, confirmed by diff, and got three real assertion failures in 0.48s. Content byte-identical to the hand-back's claim. No compile or import error, so the RED is genuine.
- **R2 — GREEN reproduced.** The working tree is byte-identical to `d89efc0b` for this module. `gofmt -l .` silent, `go vet ./...` silent, `go test -count=1 ./...` ok in 91.772s.
- **R3 — the control proved by mutation, not by passing.** Widening the gate to `if false && ...` in a scratch copy made `TestVerifyFlagsQueueStateCommittedOnABuilderBranch` fail with `got 0 committed-queue-state findings, want 1`, while the new test still passed. The pair is real, and the new test alone would not have caught a widened gate.
- **R4 — the fact reaches a reader.** Printed the real CLI render for the detached fixture: one finding row carrying the new clause in its remedy, and `SkippedProbes` holding only the unrelated calibration-log entry.
- **R5 — consumers grepped, not assumed.** `SkippedProbes` has exactly two production consumers, `renderVerifyReport` and `attachVerifyFindings` via `web/board-cards.js`. No JSON field, contract script or `_dev/tests` check reads it. Both consumers carry `Remedy`, so the moved fact renders in both surfaces.

#### Gating-condition analysis

The gate is correct and no other disposition should join it. Undetermined is the only disposition where the branch ref may fail to resolve: the unmerged-branch case requires `merge-base` to have answered, and the state-unknown, in-flight, uncommitted-work and finished-residue cases are all reached only after `merge-base` returned 0. The one case that would break the reasoning, an unborn HEAD, is already caught upstream by the integration-ref check.

#### Findings

**Important:** None.

**Minor (all report only):**

- `VerifyReport`'s doc comment (verify.go:82-90) still says unverified coverage lives in exactly two lists. There is now a third resting place — inside a finding's remedy — and the comment was not amended. A future reader could treat `SkippedProbes` as exhaustive. *impact-rule-change*
- D-03 argues for hedging the cause, but the shipped text does not hedge: both the new remedy clause and the new code comment assert the unresolved branch as fact, while `classifyWorktreeMergeState` returns undetermined for any git exit code other than 1. *impact-negligible*
- Discovered task T-03 misdescribes the function. `appendWorktreeFindings` has four other `SkippedProbes` producers, not two, and one is per-leftover rather than per-run: the worktree removability probe at verify.go:1119-1121. T-03's conclusion still holds, but for a different reason — that pair is complementary by design, because its remedy points the reader at the skipped line. *impact-negligible*
- No lesson or prime note records the new exception. The satellite's `unknown-reads-as-clean` entry still prescribes routing the unknowable arm to a `SkippedProbes` line. The rule's substance holds; the mechanism it names now has one documented exception. Charged to the Lessons-Capture Phase, not to the builder. *impact-rule-change*

**Nit (report only):** the hand-back's RED quote cites test lines 799/806/809 against the shipped 803/810/813, a consistent 4-line offset from a later doc-comment expansion; the three assertions are identical. The new test computes the rendered report but uses it only in failure messages, so it pins the struct field rather than the reader-facing surface.

#### Restatement sweep

One thing is redefined: when the committed-queue-state probe runs, and where its non-coverage is stated. Swept and found not stale: `work-reference.md:488`, `prime-do-kanban.md` § Traps first bullet, the state-unknown remedy at verify.go:1051-1053, and REQ-579's screenshot description. One consequence worth knowing: `forensics.md`'s mapping table sends only skipped probes to its Skipped or Unverified Coverage section, so in a forensics report this particular unverified coverage now appears under Warnings instead. No information lost; it changed sections.

#### Scores

| Dimension | Score | Notes |
|---|---|---|
| Requirements | 100% | All four detailed requirements plus the UR-118 batch constraint, independently verified |
| Code Quality | 88% | Minimal and idiomatic; two wording and doc-comment accuracy gaps |
| Test Adequacy | 92% | RED reproduced, control proved by mutation; no assertion on the rendered surface |
| Scope | 100% | Exactly the two declared files, no drift, no debug artifacts |
| Risk | None | Read-only probe; other dispositions provably unaffected |
| Acceptance | Pass | Full module suite green in 91.8s; end-to-end render verified |

#### Follow-ups created

None — all six findings are report only.

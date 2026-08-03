---
id: REQ-083
title: verify reports every builder worktree as a fixable orphan, including active and unmerged ones
status: pending
created_at: 2026-08-03T17:09:21Z
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
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go]
---

# `verify` Reports Every Builder Worktree as a Fixable Orphan, Including Active and Unmerged Ones

## What

`appendWorktreeFindings` emits an `orphan-worktree` finding with `Fixable: true` for every
`worktree-agent-*` worktree and branch it finds, with no merged / unmerged / still-building
distinction. So an in-flight builder is reported as a leftover, and unmerged builder work — which
`actions/cleanup.md` Pass 5 will not delete without asking — is advertised in the report footer as
mechanically fixable.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
*Source: external audit finding F3 (P1). Accepted by `do-work validate-feedback` triage with its
severity corrected: the audit claimed Step 9's verifier would fail the run, but `actions/work.md:606`
never gates on verify's exit code. The defect is the `Fixable` contract violation and the missing
classification.*

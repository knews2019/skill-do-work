---
id: REQ-450
title: 'Exclude already-claimed requests before selection'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-451, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Exclude Already-Claimed Requests Before Selection

## What

Exclude a queued request before eligibility when its request record carries `claimed_at` or a foreign writer entry in `do-work/CHECKPOINT.md` still claims it. Return typed already-claimed evidence without inferring lease expiry or writer liveness.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this missing claim-aware selector eligibility root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #1 — P1 — source:** `internal/nextselection/next_selection.go:159-160`

> ````text
> [P1] Exclude claimed requests before selecting — [prj].claude/skills/do-work/tools/do-
> work-cli/internal/nextselection/next_selection.go:159-160
> When a pending queue record has a live claimed_at, or remains named under a writer: entry in do-work/CHECKPOINT.md, this path
> proceeds directly to other filters and can select already-owned work for automatic dispatch. The selector contract requires every
> candidate to be unclaimed (.claude/skills/do-work/actions/work-reference.md:395); load checkpoint claim evidence and reject
> claimed records before selection.
> ````

- **Finding #9 — P1 — source:** `internal/nextselection/next_selection.go:154-160`

> ````text
> [P1] Exclude pending records that still carry a claim — [prj].claude/skills/do-work/tools/do-work-
> cli/internal/nextselection/next_selection.go:154-160
> A queue file with status: pending but a non-empty claimed_at passes this gate and can be selected for another build. Auto-wave's
> documented predicate explicitly requires no live claimed_at (actions/work-reference.md line 395), and the replaced simple selector
> also rejected this state, so inspect record.ClaimedAt before admitting the candidate.
> ````

- **Finding #17 — P2 — source:** `internal/nextselection/next_selection.go:154-157`

> ````text
> [P2] Exclude queued records that already carry a claim stamp — [prj].claude/skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go:154-157
> For a queued record whose status is still pending but whose claimed_at is non-empty, such as after an interrupted or manual claim, this status gate accepts it because ClaimedAt is never inspected. Default, fan-out, and simple selection can consequently hand already-
> claimed work to another builder, contrary to the unclaimed selection contract in actions/work.md:35 (.claude/skills/do-work/actions/work.md#L35); emit a typed already-claimed exclusion before eligibility.
> ````

- **Evidence:** `RequestRecord.ClaimedAt` is parsed at `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go:49,270`, while `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go:154-166` checks status, assignment, and impact without checking `ClaimedAt`. `RepositorySnapshot` lacks checkpoint claim evidence at `internal/repositorymodel/repository_model.go:84-97`, and discovery ignores `CHECKPOINT.md` at `repository_model.go:187-222`. The selector contract at `skills/do-work/actions/work-reference.md:391-395` requires both no `claimed_at` and no claim by another writer in the synced checkpoint; ADR-018 lines 35-37 and 49-55 records the duplicate-writer incident and governing rule.
- **Surface-cost result:** Earned — the repository documents a synced-checkpoint double-claim incident. A small typed checkpoint projection and claim exclusion cost less than duplicate dispatch; cover pending plus `claimed_at`, a queued REQ in a foreign writer entry, and unrelated checkpoint entries.

## Detailed Requirements

- Reject a pending or queued candidate carrying a non-empty `claimed_at` before any eligibility or dispatch decision.
- Load the typed checkpoint claim evidence required to reject a request named by another writer.
- Emit a typed already-claimed exclusion with enough provenance for callers to explain the result.
- Apply the exclusion consistently to default, fan-out, simple, and targeted selection wherever the unclaimed contract applies.
- Do not invent lease expiry, heartbeat, liveness, or stale-claim heuristics.
- Ignore unrelated checkpoint entries.

## Constraints

- Preserve explicit targeting's documented overrides without turning it into an override for ownership.
- Preserve the selector's no-rescan result-model contract.

## Dependencies

No request prerequisite. Shared selector files with other UR-085 requests do not establish dependency ordering.

## Builder Guidance

Certainty level: Firm. Model stored ownership evidence directly and keep policy out of the projection.

## Red-Green Proof

**RED prompt/case:** Select (1) a pending queue record with non-empty `claimed_at`, (2) a pending queue record named under another writer in `do-work/CHECKPOINT.md`, and (3) a control record with only unrelated checkpoint entries.
**Why RED now:** Selection parses the claim stamp but never applies it, and repository discovery supplies no checkpoint claim evidence.
**GREEN when:** The first two records return typed already-claimed exclusions before eligibility, while the unrelated control remains eligible.
**Validation:** User confirmed after validate-feedback accepted Findings #1/#9/#17.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Findings #1, #9, and #17, captured by UR-085.*

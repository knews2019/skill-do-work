---
id: REQ-623
title: 'Reduce orchestration instruction duplication'
status: pending
created_at: 2026-09-08T21:49:05Z
user_request: UR-130
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md"]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: ["REQ-622"]
related: ["REQ-622", "REQ-624"]
batch: instruction-history-cleanup
write_set: ["skills/do-work/SKILL.md", "skills/do-work/actions/work.md", "skills/do-work/actions/work-reference.md", "skills/do-work/actions/review-work.md"]
---

# Reduce Orchestration Instruction Duplication

## What
Remove proven duplication from the core SKILL.md, work.md, work-reference.md, and review-work.md. Preserve command contracts, judgment, authored evidence, recovery guidance, and behavior. Verify representative normal and recovery paths and report the text reduction; make no unmeasured runtime-savings claim.

## Detailed Requirements
- Reduce proven restatements in the four declared core files around their canonical owners. Preserve when to invoke commands, required inputs, typed-result interpretation, authored evidence, judgment, and recovery guidance.
- Record before/after text counts and show which surviving owner covers each removed responsibility.

## Constraints
Preserve behavior. Read instrumentation is optional for this cleanup; whole-file counts and a CLI finding code alone do not prove either runtime savings or that prose is dispensable.

## Dependencies
Depends on REQ-622 (Correct verified prose drift), so the cleanup starts from reconciled wording. REQ-624 (Addendum: Trim the shipped changelog to 50 entries) follows this request by the approved serial order.

## Builder Guidance
Use the smallest justified subtraction. No pending or pending-answers candidate in any UR shares this request's root cause; the seven existing queued requests address different release, commit-integrity, and diagnostics defects. The prior guard dedupe in REQ-023 (Dedupe intra-file guard restatements in note, scan-ideas, commit, quick-wins) touched different actions; this request targets the four named core files.

## Completion Proof
Use existing focused contract checks and representative walkthroughs of a simple request, a complex request, and an interrupted/recovery path. Each removed restatement must have a surviving owner with the same responsibility; restore guidance if its absence changes a decision. Report measured text reduction without a runtime-savings claim unless that outcome is actually measured.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens; condition-preserving prose extraction in the core action files. Partial slug coverage prevents targeted selection.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens; prescribed command blocks and moved prose anchors. Partial slug coverage prevents targeted selection.

## Full Context
See `do-work/user-requests/UR-130/input.md` for the user instruction, adopted report, and batch decisions. The approved prompt excerpts are in `do-work/user-requests/UR-130/assets/rev-v2-approved-scope.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the listed primes and agent rules; record a brief approach.
- [ ] **[APPLY]:** Make the scoped changes.
- [ ] **[UNIFY]:** Review every changed file and record the focused checks and outcomes.

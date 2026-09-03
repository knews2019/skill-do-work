---
id: REQ-471
title: 'Flow and reader consistency plus documentation for gate-blocked set-aside'
status: cancelled
created_at: 2026-09-01T04:29:16Z
user_request: UR-087
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-469]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-468, REQ-469, REQ-470, REQ-472]
batch: non-blocking-orchestration
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/roadmap.md, skills/do-work/actions/clarify.md, skills/do-work/actions/cleanup.md, skills/do-work/docs/work-guide.md, _dev/tests/contract-regressions.sh]
completed_at: 2026-09-03T20:40:52Z
---

# Flow and Reader Consistency Plus Documentation for Gate-Blocked Set-Aside

## What

Apply the non-blocking set-aside behavior consistently across every flow — default, targeted, wave, fan-out, crash-recovery, checkpoint, cleanup, roadmap, clarify, and composed-summary — and sweep stale language that still says an unrelated canonical-gate failure must preserve a claim and stop the session. Queue summaries and roadmap output must clearly distinguish blocked work, pending user decisions, dependency-gated work, and runnable work.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- "Apply the behavior consistently to default, targeted, wave, fan-out, crash-recovery, checkpoint, cleanup, roadmap, clarify, and composed-summary flows."
- Acceptance test: "Queue summaries and roadmap output clearly distinguish blocked work, pending user decisions, dependency-gated work, and runnable work."
- "Update the authoritative orchestration instructions, schema/restatements, board/readers, documentation, and regression tests together. Search for stale language that says an unrelated canonical-gate failure must preserve a claim and stop the entire session."

## Constraints

- Most flows inherit the behavior once the hold becomes a `blocked` flip: targeted runs already probe blocked REQs, cleanup already leaves `blocked` untouched, clarify Step 5.5 already confirms blocked conditions, crash recovery already preserves `blocked`, and the Composed Exit Summary already has a blocked-external section. The work here is verifying each flow against the new set-aside and adding only what is missing — e.g. a summary/roadmap line shape for a gate-blocked REQ (its `blocked_by` names the gate) and fan-out/crash-recovery/checkpoint notes.
- Readers already distinguish the four kinds (roadmap's Ready / Needs clarification / Blocked buckets with per-shape remedy lines; the composed summary's per-cause sections; the board's PendingReady/PendingWaiting split and status pills). Assert and extend where the gate-blocked shape needs its own line; do not rebuild the taxonomy. The board's merged "Needs input · Blocked" column stays as is — triage judged it out of scope.
- Stale-language sweep — the condition is the rule, but known sites from triage: `actions/work.md` Step 6.5 item 4 residue, the Step 6.5 orchestrator-checklist line, the Error Handling table row, and the `_dev/tests/contract-regressions.sh` canonical-gate lane (REQ-469 owns the primary edits; this REQ sweeps the remainder, including `docs/work-guide.md` and any other docs restating the hold). Archived REQs are immutable provenance — never edited; `CHANGELOG.md` is owner-only.
- Contract-regressions assertions pinning any edited prose change in the same commit.

## Dependencies

- REQ-469 (blocked set-aside) — this REQ propagates its behavior to the surrounding flows, readers, and docs.

## Builder Guidance

Certainty: Firm on the outcome; latitude on how much each flow needs (verification may show several need nothing — record that, don't invent edits). Delete-before-you-add applies to the stale-language sweep.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** `grep -rn "preserve the claimed REQ" skills/ _dev/` still returns hold-language sites outside REQ-469's core edits, and the Composed Exit Summary / roadmap have no line shape for a gate-blocked REQ.
**Why RED now:** The set-aside behavior would exist in Step 6.5 but be contradicted or unreported by surrounding flows, summaries, and docs.
**GREEN when:** No shipped instruction or doc says an unrelated canonical-gate failure preserves a claim and stops the session; each of the ten flows is verified consistent (with edits or an explicit no-change finding); summaries/roadmap render the gate-blocked shape; contract assertions updated in the same commit and `bash _dev/tests/contract-regressions.sh` exits zero.
**Validation:** Inferred during capture (from the spec's acceptance tests)

## Folded From REQ-472

Hand triage 2026-09-03: REQ-472's scenario list was folded into REQ-469 § Folded From REQ-472, with each scenario assigned to the REQ that owns it. This REQ proves the scenarios marked with its id there in its own Testing section; no separate test REQ exists.

## Full Context
See `do-work/user-requests/UR-087/input.md` for complete verbatim input.

---
*Source: UR-087 — "Apply the behavior consistently to default, targeted, wave, fan-out, crash-recovery, checkpoint, cleanup, roadmap, clarify, and composed-summary flows."*

## Cancelled

- **When:** 2026-09-03T20:40:52Z
- **Why:** folded into REQ-510's sweep as one line; the gate-blocked set-aside it documents was superseded by the repository-gate deferral lifecycle that shipped as REQ-491 to REQ-494 on 2026-09-02. Maintainer's 2026-09-03 triage.
- **Decided by:** user, via `do-work abandon`

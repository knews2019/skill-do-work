---
id: REQ-564
title: 'Reuse matching per-lane verification evidence for four hours'
status: claimed
created_at: 2026-09-03T22:58:23Z
user_request: UR-109
domain: testing
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-563]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-539, REQ-563]
batch: smart-heavy-verification
claimed_at: 2026-09-04T20:02:18Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-04T20:03:07Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - persistence changes
    - cross-route regression gates
dispatch_at: 2026-09-04T20:03:07Z
---

# Reuse Matching Per-Lane Verification Evidence for Four Hours

## What

Cache successful heavy-lane results for at most four hours and reuse a lane's evidence only when a deterministic fingerprint still matches. Allow unaffected lanes to reuse valid evidence while affected lanes rerun, and record the disposition of every lane.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Recent successful heavy-lane work should not be repeated when its complete inputs are unchanged, but age alone is too weak to prove that evidence still applies.

## Context

REQ-563 (Select affected heavy-test lanes from request changes) supplies stable lane identities and the selected plan. This request adds bounded, fingerprint-validated reuse to that plan without turning a recent timestamp into authorization.

## Detailed Requirements

- Cache each successful lane result for at most four hours.
- Reuse evidence only when a deterministic fingerprint of the lane's command, test inputs, fixtures, toolchain, and required environment still matches.
- Time alone must never authorize reuse.
- Rerun a lane when its fingerprint differs, its evidence is older than four hours, its prior result was not successful, or fingerprint coverage is uncertain.
- Allow unaffected lanes to reuse matching evidence while affected lanes rerun in the same verification plan.
- Record whether each lane was executed or reused.

## Constraints

- Build this after REQ-539's aggregate split and after REQ-563 establishes lane selection and identity.
- The four-hour window is a maximum age, not a guarantee of reuse.
- Evidence reuse must fail closed when any required fingerprint input cannot be determined.

## Dependencies

- Requires REQ-563 (Select affected heavy-test lanes from request changes), which is itself downstream of REQ-539.

## Builder Guidance

Firm on the fingerprint inputs, maximum age, per-lane partial reuse, fail-closed behavior, and executed-versus-reused record. The cache representation and persistence location are builder decisions; keep invalidation deterministic and auditable.

## Red-Green Proof

**RED prompt/case:** Record a successful heavy-lane result, verify again within four hours with an identical fingerprint, then change one fixture or required toolchain/environment input and verify again.
**Why RED now:** The current gate has no per-lane evidence cache or deterministic reuse decision, so the identical lane reruns and there is no safe invalidation behavior to prove.
**GREEN when:** The identical lane is reused and recorded as reused, the lane with a changed fingerprint executes and is recorded as executed, unaffected lanes can still reuse their own matching evidence, and evidence older than four hours executes even when its fingerprint otherwise matches.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens; relevant to verification evidence and workflow behavior, but the bare partial-slugged satellite exceeds the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5924 tokens; relevant to deterministic structured evidence and cache validation, but the bare partial-slugged satellite exceeds the 2000-token budget.

## Full Context

See `do-work/user-requests/UR-109/input.md` for complete verbatim input.

---
*Source: Replace the all-or-nothing heavy-test trigger with change-aware lane selection and reusable per-lane evidence. Select only heavy lanes affected by the request's changed paths, explain why each lane was selected, and fall back to the complete heavy suite whenever coverage is uncertain. Cache each successful lane result for at most four hours, but reuse it only when a deterministic fingerprint of its command, test inputs, fixtures, toolchain, and required environment still matches; time alone must never authorize reuse. Allow unaffected lanes to reuse evidence while affected lanes rerun, preserve --heavy as a force-all override, and record whether each lane was executed or reused. Build this after REQ-539's aggregate split.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The caching rule is specified, but the lane fingerprint and the existing heavy-verification plan and record sites have to be located before the reuse window can be added. Outcome clear, location needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Handoff State (session interrupted 2026-09-04T20:2xZ)

**Builder branch:** `worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours` at `27b74c8`, pushed to origin. **Unmerged, and this commit is an orchestrator-authored WIP snapshot, not a builder hand-back.**

The builder was interrupted mid-write with everything uncommitted. The orchestrator committed the working tree verbatim so an ephemeral container could not destroy it. There is no hand-back file, no file manifest, no decisions record, and **no red-green evidence for anything in it**.

Snapshot contents: a new `heavy_evidence.go` and its test, a reworked `heavy_run.go`, changes across `heavy_commands`, `heavy_verification` and `resultmodel`, an extended `_dev/tests/heavy-lanes.json` fixture, and one-line edits to `work.md` and `clarify.md`.

**Treat every line of it as unverified.** Re-run the package tests before trusting any of it, and require revert-and-show-red proof on each new test before merging. A test that cannot fail is the failure mode this session already hit once.

**Sizing note, so the work is not oversold.** This REQ is queued as a queue-speed improvement. Measured this session: the repository gate ran nine times at roughly ten minutes each, about half the run's wall clock — but those were **fast-tier** runs, not heavy lanes. This REQ targets heavy lanes, so its benefit here is real but narrower than the title suggests. Do not widen scope to chase the fast tier; record it as a discovered task if a clean shape appears.

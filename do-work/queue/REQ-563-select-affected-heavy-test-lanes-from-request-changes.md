---
id: REQ-563
title: 'Select affected heavy-test lanes from request changes'
status: pending
created_at: 2026-09-03T22:58:23Z
user_request: UR-109
domain: testing
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-539]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-539, REQ-564]
batch: smart-heavy-verification
---

# Select Affected Heavy-Test Lanes From Request Changes

## What

Replace the all-or-nothing heavy-test trigger with change-aware lane selection. Select only heavy lanes affected by a request's changed paths, explain every selection, and fall back to the complete heavy suite whenever coverage is uncertain.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The broad heavy gate can rerun unrelated expensive lanes and block an otherwise valid request on timing failures outside the paths it changed.

## Context

REQ-539 (Cut the contract file to the incident core and split the aggregate into fast and heavy) is the prerequisite that gives the heavy aggregate explicit lane boundaries. This request uses those boundaries to make selection change-aware; REQ-564 (Reuse matching per-lane verification evidence for four hours) adds reuse after lane identity exists.

## Detailed Requirements

- Select only heavy lanes affected by the request's changed paths.
- Explain why each lane was selected.
- Fall back to the complete heavy suite whenever coverage is uncertain.
- Preserve `--heavy` as a force-all override that is not narrowed by changed-path selection.
- Expose stable lane identities that per-lane evidence can reference.

## Constraints

- Build this after REQ-539's aggregate split.
- A change-aware selection must never silently trade away coverage; uncertainty selects the complete heavy suite.
- Keep the behavior compatible with the existing explicit heavy-test entry point.

## Dependencies

- Requires REQ-539 (Cut the contract file to the incident core and split the aggregate into fast and heavy).
- REQ-564 (Reuse matching per-lane verification evidence for four hours) requires this request's stable lane-selection result.

## Builder Guidance

Firm on the observable selection, explanation, fallback, and force-all behavior. The representation of the changed-path-to-lane mapping is a builder decision; prefer a condition-based mapping that fails closed when it cannot prove coverage.

## Red-Green Proof

**RED prompt/case:** Run request verification for a change limited to one known heavy lane, then for a changed path that no lane mapping recognizes.
**Why RED now:** The current heavy trigger is all-or-nothing, so the known change launches unrelated heavy lanes and has no change-aware uncertainty fallback to inspect.
**GREEN when:** The known change selects only its affected heavy lane or lanes and reports the reason for each, while the unrecognized path selects the complete heavy suite with an explicit uncertainty reason; explicit `--heavy` still selects every heavy lane.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens; relevant to routing and evidence behavior in action files, but the bare partial-slugged satellite exceeds the 2000-token budget.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; relevant to safe heavy-lane command selection, but the bare partial-slugged satellite exceeds the 2000-token budget.

## Full Context

See `do-work/user-requests/UR-109/input.md` for complete verbatim input.

---
*Source: Replace the all-or-nothing heavy-test trigger with change-aware lane selection and reusable per-lane evidence. Select only heavy lanes affected by the request's changed paths, explain why each lane was selected, and fall back to the complete heavy suite whenever coverage is uncertain. Cache each successful lane result for at most four hours, but reuse it only when a deterministic fingerprint of its command, test inputs, fixtures, toolchain, and required environment still matches; time alone must never authorize reuse. Allow unaffected lanes to reuse evidence while affected lanes rerun, preserve --heavy as a force-all override, and record whether each lane was executed or reused. Build this after REQ-539's aggregate split.*

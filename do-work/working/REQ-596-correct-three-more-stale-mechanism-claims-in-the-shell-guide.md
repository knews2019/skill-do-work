---
id: REQ-596
status: claimed
domain: general
created_at: 2026-09-06T05:18:15Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: B
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T05:21:33Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
maintenance: true
depends_on: [REQ-595]
related: [REQ-555, REQ-595]
write_set: [skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Correct three more stale mechanism claims in the prescribed-shell guide, in sections REQ-595 never opened'
claimed_at: 2026-09-06T05:20:24Z
---

# Correct Three More Stale Mechanism Claims in the Prescribed-Shell Guide

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Three sentences in `skills/do-work/docs/prescribed-shell-primitives.md` describe behaviour the Go code
does not have. All three are in sections REQ-595 did not open, so that request left them alone rather
than widening a second time. Each was verified against the code by an independent reviewer.

## Why

Same finding class as REQ-555 and REQ-595: the guide is the pointer target from sixteen shipped files
and currently describes implementations that do not exist. One of the three loosens a secret-quarantine
instruction, which is the only one with a safety edge.

## Context

Found during the independent three-lens review of REQ-595. That review scored 86% and confirmed all
three against the code; they are outside REQ-595's declared class, so they are captured here.

## Detailed Requirements

- **Line 124, atomic download.** The guide says the download verifies what it wrote and removes its own
  nested artifact. The code pre-checks instead, and a rename cannot nest. Describe what
  `atomic-download` really does.
- **Line 72, inventory ordering.** "An archived REQ outranks an in-flight one" is not what the code
  does: the active flag is stored at `internal/corehelpers/inventory.go:307` and never read, and the
  only comparison is on the completion instant.
- **Line 57, secret quarantine.** The guide's pattern is `*credentials*`, narrower than the code's
  `strings.Contains(base, "credential")`. A reader following the guide's by-hand fallback would miss a
  file the tool quarantines. Correct the pattern to match the code.
- Check the sections these three sit in for the same class in the same pass, the way REQ-595 checked all
  fourteen Mechanics cells. The sentence that was checked says nothing about its neighbours.

## Constraints

- Guide prose only: no behaviour change, no code change, no route-column change, no heading change.
- `_dev/tests/audit-lockins.sh` Finding 7 pins the route column, the orchestration claim and the
  Mechanics column's shell vocabulary. Do not weaken any of the three.
- If a claim cannot be verified against code, say so in the request rather than guessing.

## Dependencies

Depends on REQ-595, which rewrote three Mechanics cells and two prose claims in the same file.

## Open Questions

None.

## Triage

**Route: B** — Explore then build.

**Reasoning:** The three named claims are already settled — each was verified against the code by the
REQ-595 reviewer and re-verified here before the route was chosen. What makes exploration real is the
request's fourth requirement: check the sections these three sit in for the same class in the same pass.
That is the requirement that turned REQ-595 from a one-cell edit into a three-cell one, and the same
argument applies here — the sentence that was reported says nothing about its neighbours.

**Planning:** Skipped. One file, and the work is whatever the sweep finds.

**One of the three has a safety edge and is the reason this is not "later".** The guide's by-hand
secret-quarantine fallback lists `*credentials*`; the code matches `credential`. A reader following the
guide when the tool is unavailable would not quarantine `credential.json`. The fallback is what someone
executes by hand, so a pattern narrower than the code's is a miss, not a nit.

## Plan

**Planning not required** — Route B: one file, and the edits are whatever the section sweep confirms.

*Skipped by work action*

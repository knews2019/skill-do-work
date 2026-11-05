---
id: REQ-711
title: 'Repair runtime fingerprint baseline failures'
status: claimed
route: B
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
write_set: [skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go]
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-05T16:37:19Z
  basis: [Route B, 2-file write set, 3 acceptance criteria, cross-route regression gates]
created_at: 2026-11-05T16:34:59Z
user_request: UR-163
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: runtime-fingerprint-baseline
depends_on: []
related: [REQ-710]
claimed_at: 2026-11-05T16:35:25Z
---

# Repair runtime fingerprint baseline failures

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97 affecting REQ-710 (found by REQ-710 / UR-163)

## Repository Gate Repair Intake

- **Parent:** REQ-710
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97
- **Repair dependency:** REQ-711
- **Diagnostic evidence:** "baseline revision 40e9e4f6ee0f9ed19404d1515e7750f3a453a398"
- **Diagnostic evidence:** "canonical gate exited 1 twice before implementation"
- **Diagnostic evidence:** "TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes: default runtime must have a determinable fingerprint (GOOS, GOARCH, CGO_ENABLED, DO_WORK_TEST_DO_WORK_CLI_BINARY)"
- **Diagnostic evidence:** "TestShippedGitIsolationPreservesGenericLaneInheritance: shipped runtime probe did not isolate host Git configuration"

## Triage

**Route: B** — Specific pre-existing runtime-fingerprint failures; exploration will locate the environment-sensitive input and minimal repair.


## Plan

Planning not required (Route B). Reproduce the runtime probe failure, isolate test Go settings without weakening reuse safety, and rerun focused and canonical checks.

## Exploration

Explore agent traced the failures to host effective GOFLAGS=-trimpath. The production helper correctly refuses nonempty GOFLAGS; fixture newRuntimeEvidenceRepository inherits the host GOENV. The GOENV subtest passes because it supplies its own file. Git-isolation failure has the same cause, not a production Git-isolation defect. Isolate the fixture baseline from unrelated Go environment settings; retain tests that deliberately alter settings after fixture setup.

Required-lessons index consulted: CLI satellite (24067 tokens, partial coverage) exceeds projection budget. No required_lessons stamp; builder reads touch-relevant prime/satellite.

## Scope

**Files I will touch:**
- skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go

**Acceptance criteria:** Existing named tests pass with hostile host Go settings; tests continue to exercise deliberate per-case setting changes; production refuses opaque flags/workspaces unchanged. Canonical gate must pass before REQ-710 resumes.

## Red-Green Proof

RED: Existing TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes and TestShippedGitIsolationPreservesGenericLaneInheritance fail under inherited GOFLAGS=-trimpath, observed in two baseline gate runs. GREEN: deterministic fixture baseline allows these existing tests to pass while deliberate runtime changes still invalidate reuse. Add a focused regression for the host setting if needed to pin fixture isolation.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Isolate only fixture environment, preserving production fail-closed behavior.
- [ ] **[APPLY]:** Implement the declared fixture correction.
- [ ] **[UNIFY]:** Verify focused tests, diff and integrated canonical gate.

## Pre-Build Evidence

Canonical gate failed twice (exit 1, then 1) before implementation with intake fingerprint 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97. Focused advance preflight also failed twice (1, 1) with the same named tests. This repository_gate_repair has the exact expected red baseline and may implement under the repair exception; no recursive deferral.

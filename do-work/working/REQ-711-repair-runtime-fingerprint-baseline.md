---
id: REQ-711
title: 'Repair runtime fingerprint baseline failures'
status: claimed
route: B
builder_handback_at: 2026-11-05T16:42:44Z
integration_at: 2026-11-05T16:42:44Z
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
review_at: 2026-11-05T16:47:07Z
commit: 2c3825997ffd303be982f91f304365b9fb808d20
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
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go` (modify)

**Acceptance criteria:** Existing named tests pass with hostile host Go settings; tests continue to exercise deliberate per-case setting changes; production refuses opaque flags/workspaces unchanged. Canonical gate must pass before REQ-710 resumes.

## Red-Green Proof

RED: Existing TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes and TestShippedGitIsolationPreservesGenericLaneInheritance fail under inherited GOFLAGS=-trimpath, observed in two baseline gate runs. GREEN: deterministic fixture baseline allows these existing tests to pass while deliberate runtime changes still invalidate reuse. Add a focused regression for the host setting if needed to pin fixture isolation.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Isolate only fixture environment, preserving production fail-closed behavior.
- [x] **[APPLY]:** Implement the declared fixture correction.
- [x] **[UNIFY]:** Builder checked gofmt, focused RED/GREEN, hostile GOENV replay and diff; orchestrator checked the four-line integrated diff. Canonical gate remains the next testing step.

## Pre-Build Evidence

Canonical gate failed twice (exit 1, then 1) before implementation with intake fingerprint 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97. Focused advance preflight also failed twice (1, 1) with the same named tests. This repository_gate_repair has the exact expected red baseline and may implement under the repair exception; no recursive deferral.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go` (modified) — fixture sets GOENV=off, empty GOFLAGS and GOWORK=off before intentional per-case overrides.

**What was done:** Isolated the runtime-evidence test fixture from host Go configuration. Production fail-closed fingerprinting is unchanged. Builder commit 3c6efd59ec172b8481cbca0ccdacc044573bd078 merged as 2c3825997ffd303be982f91f304365b9fb808d20; qualification range d7c13a8cab8eb115e735b7041243107dfbc94bbf..2c3825997ffd303be982f91f304365b9fb808d20.

Builder handback: do-work/runs/2026-11-05-req711/REQ-711-handback.md. Existing tests reproduced RED before edit; all six relevant tests passed after edit with inherited host settings and temporary GOENV containing GOFLAGS=-trimpath. No persistent Go settings changed. No missing lesson paths reported.

## Qualification

Advance qualify and scope-drift records are satisfied for exact REQ/path and range d7c13a8cab8eb115e735b7041243107dfbc94bbf..2c3825997ffd303be982f91f304365b9fb808d20. Initial scope parser warning was missing Markdown backticks in the authored Scope, corrected without changing its declared path. Four-line fixture diff is substantive, matches requirements, and preserves production refusal and intentional per-case overrides.

## Testing

**Red-green validation:** Named runtime-identity and Git-isolation tests failed before edits (canonical baseline 1 then 1; focused preflight 1 then 1), passed after fixture isolation. Builder replay with hostile GOENV passed six tests in 4.362s. Integrated focused advance probe passed in 4.370s (exit 0, no timeout); reviewer independent subset passed in 5.305s. Production opaque-override tests remain green.
**Repository gate:** bash _dev/tests/maintainer-verify.sh passed directly, exit 0, 119s on merge 2c3825997ffd303be982f91f304365b9fb808d20. Advance recorded matching request/path green gate evidence.
**Heavy lanes selected:** do-work-cli-integrations, staged-skills, updater, installer; pending drain. No test expectation weakened; fixture host isolation is the intentional cross-test correction.

## Review

**Overall: 100%** | 2026-11-05T16:47:07Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** None
**Acceptance:** Pass — implementation and targeted integration independently reviewed; parent also verified integrated canonical fast gate. Heavy lanes pending. No deployment stage applies to this export-ignored fixture.
**Restatement sweep:** nothing redefined outside fixture setup.
**Suggested testing:** selected heavy lanes.
**Follow-ups created:** None

## Lessons Learned

**What worked:** Deterministic fixture-owned Go settings removed host contamination while retaining deliberate per-case overrides.
**What didn't:** Assuming host defaults were fingerprintable conflated opaque Go configuration with Git isolation.
**Worth knowing:** Keep production refusal of arbitrary flags; isolate test baselines instead. Existing CLI satellite already records host GOFLAGS behavior, so no duplicate lesson entry.

## Orientation

The runtime-evidence tests now pass independently of host Go configuration. This prerequisite clears the baseline failure that prevented building optional AI-Evaluation-Loop modes.

## Heavy Verification Plan

- Base revision: d7c13a8cab8eb115e735b7041243107dfbc94bbf
- Target revision: 2c3825997ffd303be982f91f304365b9fb808d20
- do-work-cli-integrations: ["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "do-work-cli-integrations"] — skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: ["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "staged-skills"] — skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go matched subtree skills
- updater: ["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "updater"] — skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go matched subtree skills/do-work/tools/do-work-cli
- installer: ["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "installer"] — skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go matched subtree skills/do-work/tools/do-work-cli

---
id: REQ-564
title: 'Reuse matching per-lane verification evidence for four hours'
status: completed
route: B
review_at: 2026-09-04T20:57:34Z
re_review_at: 2026-09-04T21:07:26Z
kb_status: pending
builder_handback_at: 2026-09-04T21:05:09Z
integration_at: 2026-09-04T21:05:09Z
remediation_at: 2026-09-04T21:05:09Z
write_set:
  - _dev/tests/heavy-lanes.json
  - _dev/tests/heavy-runtime-fingerprint.py
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/work.md
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_fingerprint_probe.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_reuse_regression_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_runtime_contract_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
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
status_changed_at: 2026-09-04T20:43:09Z
claimed_at: 2026-09-04T20:43:09Z
completed_at: 2026-09-04T21:09:33Z
commit: e9b12501713d2478a9d6faf3aaf9d9040940832b
release_at: 2026-09-04T21:09:33Z
---

# Reuse Matching Per-Lane Verification Evidence for Four Hours

## What

Cache successful heavy-lane results for at most four hours and reuse a lane's evidence only when a deterministic fingerprint still matches. Allow unaffected lanes to reuse valid evidence while affected lanes rerun, and record the disposition of every lane.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Handoff State (session stopped 2026-09-04T20:3xZ)

**Builder branch:** `worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours` at `19103f5`, pushed to origin. **Unmerged, but complete.**

The builder finished after all and wrote a full hand-back at `do-work/runs/work-2026-09-04-200249/REQ-564-handback.md`. Its branch carries three commits: `27b74c8`, an orchestrator-authored WIP snapshot taken when the builder appeared interrupted, then the builder's own `4526ab9` and `19103f5` on top of it. The snapshot is superseded — read the branch tip, not `27b74c8`.

**Resume at the Step 6 hand-back merge.** Do not re-dispatch a builder.

**Verification the builder reports** (verify, do not assume): `go vet ./...` clean, `gofmt -l ./internal` empty, `go test -count=1 ./...` all packages ok, Windows cross-compile and vet clean, `contract-regressions.sh` exit 0 with no FAIL lines, and **twelve isolating reverts each proved red**. It did not run a real heavy lane — each is roughly ten minutes — so the merged-tree gate is still owed.

**Sizing note, so the work is not oversold.** This REQ is queued as a queue-speed improvement. Measured this session: the repository gate ran nine times at roughly ten minutes each, about half the run's wall clock — but those were **fast-tier** runs, not heavy lanes. This REQ targets heavy lanes, so its benefit here is real but narrower than the title suggests.

## Triage

**Route: B** — resume completed builder branch `19103f5`; independent review before integration found unsafe reuse cases requiring remediation.

## Plan

Planning not required. Keep the per-lane cache and four-hour contract; correct invalidation and complete the inputs that authorize reuse.

## Pre-integration Review

Acceptance failed on four important findings: a failed/skipped forced rerun retained an older green; unknown-file fallback could still reuse stale evidence; browser/toolchain/environment/external executable inputs were incompletely sealed; a run-wide timestamp let later lanes reuse expired records. Toolchain probes also needed bounded execution. These are corrected within this request, before release.

## Remediation Plan

An isolated builder at `codex/req564-remediation` fixes the unsafe reuse cases with targeted RED/GREEN regressions. Scope extends the original heavyverification/resultmodel code, tests and action/prime prose, plus `_dev/tests/heavy-lanes.json` and the runtime-fingerprint helper required to determine shipped-lane inputs. No checkpoint/finalization implementation is part of this request. Browser evidence remains non-reusable while its complete runtime cannot be identified.

## Exploration

The original handback located manifest lane identity/coverage, the owned lane runner, typed result projections and Git-private evidence storage. Independent review then traced the actual helper, untracked-source, toolchain, environment and forced-rerun inputs. Remediation `0b470b91` seals those inputs and preserves custom-manifest environment semantics.

## Scope

**Files I will touch:**
- `_dev/tests/heavy-lanes.json`
- `_dev/tests/heavy-runtime-fingerprint.py`
- `skills/do-work/actions/clarify.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence_test.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_fingerprint_probe.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_reuse_regression_test.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_runtime_contract_test.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

**Acceptance criteria:** reuse only successful matching evidence no older than four hours; invalidate before a failed/skipped new attempt; independently decide every lane; fail closed on uncertain inputs; preserve force-all execution and report executed/reused disposition.

## Implementation Summary

**Files changed:**
- `_dev/tests/heavy-lanes.json` (modified)
- `_dev/tests/heavy-runtime-fingerprint.py` (new)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_fingerprint_probe.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_reuse_regression_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_runtime_contract_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** Heavy lanes now persist successful evidence and reuse it only with a matching deterministic fingerprint within four hours. Reruns revoke old greens before launch, and each lane checks the clock when considered. Fingerprints include declared/shared/unclassified committed inputs, relevant regular untracked files, inherited environment, effective tools/module settings, and supplied native CLI bytes. Probes are bounded. The shipped fixture commands explicitly isolate Git configuration; custom lanes retain their inherited environment. Browser evidence remains non-cacheable because its runtime is uncertain.

**Integration range:** `d4e4a985..e9b12501`; includes original builder `19103f5` and remediation `0b470b91`. Read every file in this manifest; Go tests discover new `_test.go` files automatically. Final focused tests/vet and diff checks pass, with no debug artifacts.

## Decisions

Accepted the original builder's Git-private per-lane store and opt-out flag. Superseded its incomplete browser/runtime-fingerprint judgment: incomplete evidence executes rather than reuses. Full environment hashing avoids an incomplete selector-variable list. Shipped Git isolation is visible in the manifest argv and applies equally to each lane and probe. Regular untracked inputs are content-sealed so normal ignored launcher binaries do not disable all useful caching. Symlinks, gitlinks and opaque runtime extensions remain uncertain.

The required runtime helper and regressions extend the initial code-only scope because complete toolchain/fixture coverage is an explicit acceptance requirement. No unrelated queue task was implemented.

## Qualification

Passed — canonical `advance` returned satisfied qualify and scope-drift records for `d4e4a985..e9b12501`. Static new-file warnings are Go package functions called by symbols and automatically discovered `_test.go` files, verified by package execution. The Python reporter intentionally writes one digest to stdout and a failure explanation to stderr; those are its protocol, not debug output. All requirements map to the real lane-run/cache/probe paths.

## Testing

**Fast verification:** the full unpiped `bash _dev/tests/maintainer-verify.sh` passed on merged `e9b12501`: ShellCheck, gofmt, aggregate contracts, both module vet checks and uncached tests. Canonical `advance` then returned satisfied scope-drift, focused-test and green-gate records; focused argv was `go test -count=1 ./internal/heavyverification ./internal/resultmodel`. Independent reviewer also ran both focused packages (17.582s / 0.532s).

**Red-green validation:** original handback records twelve isolated reversions covering reuse, input changes, expiry/non-extension, unsuccessful evidence, declaration validation and result/command projection. Remediation added observed RED cases for green → forced failure/skip → later reuse, later-lane expiry, unknown helpers, inherited exclusion variables, untracked/ignored test sources and symlink inputs. Restored implementation passes all these tests. The actual shipped fingerprint diagnostic observed identical digests twice for all five cacheable lanes with the normal ignored CLI binary present; browser correctly remained uncertain.

**Cross-REQ assertions:** manifest tests now recognize the explicit shipped environment argv and deliberate browser exception; result rendering tests require executed/reused disposition. These are intentional contract changes.

## Review

**Acceptance: Pass.** Independent final review of `d4e4a985..e9b12501` found no remaining substantive issues after remediation. Requirements, code quality, test adequacy and scope each scored 100%; risk low. Relevant merged files match the independently tested remediation commit. Partial reuse and forced execution remain intact. No automatic follow-ups.

## Lessons Learned

**What worked:** revoke prior success before a new attempt; evaluate time per lane; seal actual input bytes and effective configuration; use isolated regression cases for each invalidation path.
**What did not:** hashing a few environment names and executable path strings missed real inputs, and a run-wide clock exceeded the cache ceiling for later lanes. The original browser exception was unsafe and was replaced by execution.
**Worth knowing:** effective `go env` includes a random scratch path in GOGCCFLAGS; hash its determining settings/toolchain instead. Browser runtime inputs are not fully determinable, so browser lanes intentionally do not cache. The deferred lesson/index payload is published only with its archive target.

## Orientation

Heavy verification can reuse successful matching lane evidence for up to four hours, report each decision and rerun only lanes whose evidence changed or became uncertain. Runtime probes are bounded; custom manifest commands retain their inherited environment. The prime and work/clarify instructions describe the new runner contract.

**Real reuse proof:** after the real JavaScript heavy lane passed, a second identical invocation returned exit 0, `disposition: reused`, `disposition_reason: fingerprint_match`, unskipped, in about 0.09 seconds. The cached execution evidence remains tied to its original measured revision/time.

**Final heavy verification:** all six selected lanes executed and passed unskipped at `e9b12501713d2478a9d6faf3aaf9d9040940832b`: queue-kanban-javascript (7s), queue-kanban-browser (85s), do-work-cli-integrations (49s), staged-skills (22s), updater (53s), installer (29s). The runner returned success and no findings.

---
id: REQ-498
title: 'Make orchestrator finalization resumable'
status: completed-with-issues
created_at: 2026-09-02T13:07:19Z
user_request: UR-096
domain: backend
route: C
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_journal.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/commit.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
kb_status: pending
claimed_at: 2026-09-02T13:46:20Z
planning_at: 2026-09-02T14:00:49Z
dispatch_at: 2026-09-02T14:10:30Z
builder_handback_at: 2026-09-02T14:50:42Z
integration_at: 2026-09-02T14:51:46Z
review_at: 2026-09-02T15:11:25Z
remediation_at: 2026-09-02T15:14:26Z
re_review_at: 2026-09-02T16:22:06Z
estimate:
  p50_active_minutes: 70
  confidence: low
  calculated_at: 2026-09-02T13:47:59Z
  basis:
    - Route C
    - 12-file write set
    - 3 new files
    - 3 subsystems involved
    - 11 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
completed_at: 2026-09-02T16:27:19Z
commit: 1249e856
release_at: 2026-09-02T16:27:20Z
---

# Make Orchestrator Finalization Resumable

## What
Replace the crash-prone archive/release/commit tail with one CLI-owned, Git-private journaled finalization flow, and recover safe unfinished tails before selecting another REQ.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the durable brief, both prime lesson satellites, crew rules, and bug-fix spec; chose protected-inventory-first discovery that converts exact legacy groups into the existing journal phase engine.
- [x] **[APPLY]:** Added strict provenance modes, ordered typed finalization records, protected no-journal discovery, exact prepared-diff commit recognition, release-at images, exact pre-primary rollback, and action delegation.
- [x] **[UNIFY]:** Reviewed all 14 files, ran focused/race/full Go suites, vet, Go 1.25 compatibility, contract regressions, diff checks, removed generated residue, committed `e191b266`, and left the branch clean.

## Why
An interruption after archive removes the working claim that current recovery scans, while leaving shared checkpoint or release state dirty. Every later claim then fails on the same shared target and the orchestrator cannot make progress automatically.

## Detailed Requirements
- Add strict typed `finalize --manifest` and `recover-finalization --discover` CLI commands.
- Persist an exact Git-private journal before lifecycle mutation and advance it through prepared, lifecycle, release, primary-commit, metadata-commit, verification, and cleanup phases.
- Reuse canonical lifecycle, release, protected-inventory, and Git mutation authorities.
- Preserve sufficient exact pre/post evidence for idempotent recovery without duplicating archive moves, calibration rows, release entries, version bumps, or commits.
- Support serial provenance from the primary commit and worktree provenance from a supplied merge hash.
- Resume journals before ordinary working-REQ recovery, selection, or claim, then continue the same run automatically when shared state is safe.
- Discover legacy unjournaled tails only when project and lifecycle ownership are unambiguous; shared metadata requires semantic REQ evidence and never generic latest-owner association.
- Preserve unrelated unstaged changes; block on ambiguous shared state, foreign staged entries, dirty checkpoint state, or protected paths.
- Return typed phase, resume/discovery, commits, blockers, reason codes, and exact verification commands.
- Update work and commit action contracts to delegate finalization and startup recovery to the CLI.
- Keep existing lifecycle and release commands backward compatible and retain the single-releaser model.

## Constraints
- Journals are local Git-private state.
- Never guess shared metadata ownership or commit secret-classified content.
- Exact-path commits must not absorb unrelated staged or unstaged work.
- The current session should stop after capturing this intent and committing one safe, coherent implementation slice.

## Red-Green Proof
**RED prompt/case:** Interrupt a successful REQ after canonical archive/checkpoint mutation but before its primary or metadata commit, then invoke `do-work run` again.
**Why RED now:** Recovery scans only `do-work/working/`; the archived REQ is invisible and its dirty shared checkpoint makes every next claim refuse.
**GREEN when:** The next run resumes or safely associates the unfinished finalization exactly once, leaves no duplicate lifecycle/release effects, records provenance, and proceeds to the next selectable REQ without manual `do-work commit` intervention.
**Validation:** User confirmed through the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-096/input.md` for the complete verbatim plan and stopping instruction.

## Open Questions
None.

## Triage

**Route: C** - Complex

**Reasoning:** The request adds two typed CLI commands, journaled multi-phase recovery, protected Git mutation, legacy-state association, and changes to multiple action contracts with end-to-end interruption tests.

## Plan

Complete the safe vertical slice after commit `761d8e6a`, which already registered `finalize` and `recover-finalization`, added the Git-private journal, composed request-state/release/exact-commit/provenance authorities, and proved one lifecycle interruption.

1. Add the captured RED fixtures first: interruption across every durable journal phase, the REQ-494-shaped no-journal tail, foreign/protected staged paths, ambiguous shared state, and two independent safe legacy groups.
2. Tighten manifest provenance to explicit `primary_commit` versus `supplied_commit` modes and project ordered per-finalization typed records from one helper.
3. Bind `release_at` into recorded release images, add exact pre-primary rollback, and retain roll-forward-only recovery after the primary commit.
4. Implement bounded `recover-finalization --discover`: replay journals oldest-first, classify legacy candidates from one protected repository/Git snapshot, require whole-path ownership plus REQ-specific semantic evidence for shared metadata, journal each accepted group, and reuse the phase engine.
5. Invoke discovery at the start of `do-work run` before checkpoint/working recovery, selection, or claim; continue only from typed success with no blocked paths or reason codes.
6. Replace the separate complete/fail, release, staging, primary-commit, and provenance tail with one exact finalization manifest and `finalize` call while preserving action-owned semantic judgment.
7. Make `do-work commit` recover finalization tails before ordinary association, then lock the ordering and delegation contracts in `_dev/tests/contract-regressions.sh`.
8. Verify focused finalization behavior, race tests across finalization/Git/request-state/publication, Go 1.25 compatibility, action contracts, and the unpiped maintainer gate.

**Planned files:** the existing finalization types/commands/prepare/apply files, a new discovery file and recovery test file, result-model projection, `work.md`, `work-reference.md`, `commit.md`, action contract tests, and the CLI prime. `CHANGELOG.md`, `VERSION`, shipped version mirrors, and `actions/version.md` remain integrator-only.

**Consumer contract:** each finalization record carries request/archive/journal identity; durable phase; resumed/discovered flags; exact commit allowlist; created and settled commit hashes; exact blockers/reason codes; and exact next/verification argv. Actions branch only on typed command success plus empty blocker/reason collections.

**Plan validation:** All eleven detailed requirements map to tasks above and the typed consumer fields are explicit. Warning: the plan has eight implementation tasks, beyond the three-task quality guide, because CLI recovery, legacy association, and two action delegations form one transactional acceptance flow. The builder should keep them as one vertical slice but avoid widening discovery beyond the REQ-494 evidence classes.

*Generated by Plan agent; full planning artifact: `do-work/runs/work-2026-09-02-134759/REQ-498-plan.md`*

## Exploration

The existing foundation already owns strict journal decoding and private payloads, canonical request-state and release composition, exact-path commits, primary/metadata provenance, oldest-first journal replay, and one interruption test. The remaining gap is one coherent vertical slice: provenance/result tightening, exact rollback and terminal projection, bounded legacy discovery, startup orchestration, and delegation from `work` and `commit`.

Reuse `internal/requeststate` for lifecycle bytes, `internal/publication` for release plans, `internal/gittransaction.CommitExactPaths` for commit guards and hooks, `corehelpers.ReadProtectedInventory` for M/A/D/X/XD classification, and `requestmodel` for identity/status parsing. Do not use `AssociateProjectPaths` as the ownership decision because its latest-completed tie-break is forbidden for shared finalization recovery.

The plan omitted `finalization_journal.go`; it is required because strict manifest and recovered-journal provenance validation live there. No edits are needed in command registration, request-state, publication, Git transaction, or protected-inventory packages. Contract tests must assert active directive order rather than token presence.

Primary integration hazard: `REQ-489` concurrently fixes checkpoint section/entry handling used by finalization lifecycle images. Merge `REQ-489` first, then verify `REQ-498` against the composed request-state and finalization suites.

*Generated by Explore agent; full findings: `do-work/runs/work-2026-09-02-134759/REQ-498-exploration.md`*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go` (modify) — provenance and aggregate-record contracts
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modify) — `--discover` orchestration and replay aggregation
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go` (modify) — strict provenance fixture and exact rollback expectation required by the new contract
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_journal.go` (modify) — strict manifest/journal validation and recovery identity
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modify) — prepared evidence and discovery journal construction
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modify) — image-bound release stamp, exact rollback, terminal projections
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (new) — bounded legacy-tail classifier
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (new) — interruption, discovery, ambiguity, and idempotence matrix
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — ordered normalized finalization records
- `skills/do-work/actions/work.md` (modify) — startup recovery and one finalization call
- `skills/do-work/actions/work-reference.md` (modify) — working recovery versus finalization recovery contract
- `skills/do-work/actions/commit.md` (modify) — recover tails before ordinary association
- `_dev/tests/contract-regressions.sh` (modify) — active delegation and ordering assertions
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — finalized ownership map

**Files I will NOT touch:** command registration, request-state/publication/Git-transaction authorities, `CHANGELOG.md`, version mirrors, queue/archive/checkpoint state from the builder worktree, or a general dirty-tree association framework.

**Scope expansion accepted at handback:** `finalization_commands_test.go` was omitted from the brief, but the strict incoming provenance requirement made its foundation manifests invalid and exact hook rollback changed its expected terminal result. The minimum fixture/expectation update is requirements-implied; no other undeclared path changed.

**Acceptance criteria (restated from REQ):**
- [ ] Strict planned finalization and `recover-finalization --discover` share one Git-private journal/phase engine.
- [ ] Retries neither duplicate lifecycle/release/calibration/version/commit effects nor guess shared ownership.
- [ ] Serial provenance uses the primary commit; worktree provenance uses a validated supplied merge hash.
- [ ] Startup replays/resolves safe tails before checkpoint recovery, selection, or claim, then continues automatically only on typed clean success.
- [ ] Legacy discovery admits whole-path unambiguous project ownership and REQ-specific shared metadata only; protected, staged, foreign, or ambiguous state refuses byte-identically.
- [ ] Unrelated unstaged changes remain untouched.
- [ ] Typed records expose ordered identity, phase, resumed/discovered state, exact commit paths/hashes, blockers/reasons, and next/verification commands with non-null slices.
- [ ] `work` and `commit` delegate to the canonical finalization engine while lifecycle/release commands remain backward compatible and single-releaser semantics remain unchanged.
- [ ] The captured REQ-494 flow turns RED before implementation and GREEN afterwards, including successful selection/claim of the next REQ without manual commit intervention.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2409 tokens; matches finalization, request-state, publication, Git-transaction, and typed-evidence changes, but the partially slugged satellite cannot be narrowed below the 2000-token budget. Read anyway under the prime's touch-conditional Lessons rule.
- `_dev/primes/lessons-action-files.md` — 3436 tokens; matches work/commit pipeline delegation and downstream contract readers, but the partially slugged satellite cannot be narrowed below the 2000-token budget. Read anyway under the action-file prime's touch-conditional Lessons rule.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (created)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_journal.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (created)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** The finalization engine now discovers safe legacy no-journal tails, replays them through the existing phase journal, binds provenance and exact commit identity, rolls back pre-primary owned state exactly, exposes ordered typed records, and runs before selection or ordinary commit association. Work and commit actions delegate their lifecycle/release/commit/provenance tails to this resumable authority.

## Root Cause

Finalization authority existed in separate lifecycle, release, staging, commit, and provenance steps. Once an interruption moved a REQ out of `working/`, startup no longer had a durable owner record from which to resume the remaining shared-state tail, so later claims encountered dirty shared metadata with no safe automatic recovery path.

## Qualification

Passed — the exact merge range `e8e5a79d..75648a49` contains the 14 declared implementation files, including the justified finalization command-test expansion. Mechanical qualification confirmed substantive project changes, complete P-A-U evidence, traced requirements, scope alignment, and no queue-state leakage from the builder branch.

## Testing

**Merged-state checks:** `go test -count=1 ./internal/finalization`; `go test -race ./internal/finalization ./internal/gittransaction ./internal/requeststate ./internal/publication`; `go vet ./...`; `go test -count=1 ./...`; `bash _dev/tests/do-work-cli-go125-compatibility.sh`; `bash _dev/tests/contract-regressions.sh`; and `bash _dev/tests/maintainer-verify.sh`.

**Result:** Every focused, race, static, full-module, Go 1.25, action-contract, and canonical repository gate check passed on the merged tree. The maintainer gate's optional strict browser lane reported the repository's normal no-browser skip; queue-kanban and strict JavaScript lanes passed.

**RED→GREEN:** `TestRecoverFinalizationDiscoversLegacyNoJournalTail` failed before implementation because `recover-finalization --discover` returned `FINALIZATION-USAGE` and created no recovery commit or provenance. It passes after implementation, proving one discovered record, distinct primary/metadata commits, canonical `commit:` provenance, untouched unrelated work, journal cleanup, idempotent replay, and successful selection/claim of the next REQ.

## Review — Attempt 1

**Overall: 50%** | **Acceptance: Fail** | **Risk: Critical**

Independent review accepted the journal foundation, explicit provenance modes, startup delegation, and exact declared scope, but found the motivating semantic legacy tail unsafe and under-tested. Required remediation is consolidated under incomplete semantic legacy-finalization ownership:

- Prove lifecycle-aware ownership for calibration, UR/follow-up, changelog/version/release, `release_at`, implementation bytes, and whole-file/hunk boundaries; do not omit dirty release state or absorb foreign changes.
- Search the full prepared-descendant chain for an exact matching primary commit so an earlier independent group does not strand a later group.
- Preserve unrelated unstaged protected files while refusing staged protected/foreign state with distinct typed reasons.
- Return complete ordered terminal/refusal records, including created versus settled commits, blockers/reasons, verification/cleanup completion, and non-null slices.
- Add the full semantic, multi-group, protected-state, phase interruption, provenance-mode, idempotence, and duplicate-effect matrix.
- Remove active direct-`complete`/release/stage/provenance restatements from work action prose.

*Full independent artifact: `do-work/runs/work-2026-09-02-134759/REQ-498-review.md`.*

## Remediation

The single remediation pass replaced generic discovery with semantic grouping, added image-set integrity and durable verified/cleanup phases, searched the full descendant chain for exact primary commits, preserved unstaged protected state, distinguished staged refusal reasons, completed ordered result evidence, expanded phase/provenance/idempotence tests, and removed the originally identified direct-tail sequences. The authoritative remediation merge is `1249e856`.

## Re-review

**Overall: 50%** | **Acceptance: Fail** | **Risk: Critical**

The remediation closed multi-group recovery, protected-state classification, and typed terminal/refusal evidence. It substantially expanded the behavior matrix, but fresh review found the central semantic boundary still incomplete:

- Release discovery validates only dirty hard-coded mirrors and does not prove configured required-mirror completeness or consumer version/lock sources.
- A tracked follow-up whose current `addendum_to` points here can still carry foreign edits without creation/fold preimage proof.
- Hook-failure resume, already-green/no-release, planned-release, partial-mirror, and dirty tracked-follow-up negative cases remain absent.
- Two operative restatements remain in the work/commit action overview prose.

The one-remediation limit is exhausted. REQ-499 is widened to close the semantic ownership and acceptance-matrix gaps before adding its sole-releaser override; REQ-504 records the remaining prose instances. Terminal disposition: `completed-with-issues`.

*Full independent artifact: `do-work/runs/work-2026-09-02-134759/REQ-498-rereview.md`.*

## Lessons Learned

**What worked:** Reusing the journal phase engine, exact image identity, and typed result authority produced reliable phase replay, multi-group commit recognition, protected-state preservation, and actionable terminal evidence.

**What didn't:** Inferring legacy ownership from only the dirty paths and current document metadata was too weak. A positive end-to-end fixture did not prove completeness: missing configured mirrors and pre-existing follow-up edits require negative preimage/required-set cases.

**Worth knowing:** Recovery association is safe only when it proves both sides of the boundary—every required member is present and every admitted byte belongs to the REQ. “All observed paths look coherent” is not equivalent to “the semantic set is complete.”

## Orientation

[MAP CHANGED] Resumable finalization now lives in the do-work CLI finalization subsystem and is invoked before selection and ordinary commit association. Planned journal recovery, durable terminal phases, typed evidence, and most legacy-tail recovery are available; strict legacy release/follow-up ownership remains explicitly tracked in widened REQ-499 before sole-releaser recovery ships.

---
*Source: implement and capture the resumable orchestrator finalization plan*

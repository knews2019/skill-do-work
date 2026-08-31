---
id: REQ-413
title: 'Implement capture-file, answer, release, version, and changelog transactions'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-412]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-08-31T22:07:13Z
route: C
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/capture_files.go
  - skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/release.go
  - skills/do-work/tools/do-work-cli/internal/publication/release_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go
  - skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/stakeholder-answers.md
  - skills/do-work/actions/verify-requests.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - _dev/tests/contract-regressions.sh
estimate:
  p50_active_minutes: 105
  confidence: low
  calculated_at: 2026-08-31T22:08:35Z
  basis:
    - Route C
    - 24-file write set
    - 12 new files
    - 7 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - persistence changes
    - cross-route regression gates
    - full-suite verification
---

# Implement Capture-File, Answer, Release, Version, and Changelog Transactions

## What
Move deterministic publication and resolution phases for capture, answers, and releases into Go.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Accepted a three-task plan for typed publication manifests/planners, one atomic apply/command boundary, and delegation of every capture/answer/release writer.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement structured `do-work-capture-files` publication with atomic UR/REQ/assets writes and reservation handling.
- Implement `do-work-answer` for applying answered questions with outside-text containment and status resolution.
- Implement parameterized `do-work-release` transactions for version files and changelog updates.
- Preserve exact input containment, linkage, timestamps, collision refusal, dry-run, optional commit, and rollback.

## Constraints
- LLM capture/answer/release actions may decide content but must delegate deterministic file and state mutations to the CLI.

## Dependencies
Depends on REQ-412 (shared lifecycle transaction behavior).

## Builder Guidance
Certainty level: Firm. Separate typed publication inputs from repository mutations so callers can inspect the full intended transaction before apply.

## Red-Green Proof
**RED prompt/case:** Publish a structured capture containing outside text and assets, apply an answer, and perform a parameterized release across collision and rollback fixtures.
**Why RED now:** These deterministic writes are currently assembled and mutated by action prose and shell helpers.
**GREEN when:** Each operation is atomic, byte-safe, collision-aware, text/JSON actionable, and leaves no partial files after a pre-commit failure.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Triage

**Route: C** — Complex

**Reasoning:** This request introduces three public mutation domains plus shared publication/release primitives, spans capture and answer actions, and must preserve atomic filesystem, Git, containment, version, changelog, and rollback contracts. The cross-action persistence surface requires explicit planning and exploration.

**Planning:** Required

## Plan

1. Add a standard-library `internal/publication` family with strict typed manifests and pure planners for `capture-files`, `answer`, and `release`. RED fixtures cover malformed/escaping/symlink inputs, stale IDs/statuses, target collisions, unsafe outside bytes, invalid whole-record answer dispositions, and inconsistent/non-monotonic release replacements.
2. Apply each complete sorted plan through one Git transaction, register all three commands, and add only the narrow lossless request-body editing seam required for unique question updates. Prove byte-identical dry-run, dirty/index refusal, exact commit paths, pre-commit rollback, committed-state risk, containment, timestamps, modes, and text/JSON parity.
3. Delegate capture, clarify, stakeholder-answer, verify-repair, and work release writers to the canonical commands while retaining content, semantic-answer, bump-level, changelog-voice, and commit-choreography judgment in actions. Ratchet registration, active delegation, sole-writer ownership, and no fallback in contract regressions and the CLI prime.

**Architecture decisions:** One publication family exposes three explicit manifests rather than a generic copy API; every operation plans the full target set before one apply; raw outside bytes are the containment authority; answer disposition derives from the whole record; release validates caller-chosen replacements instead of choosing policy; shared foundations expand only for a focused RED.

**Testing approach:** Command-level UNKNOWN-COMMAND RED, pure planner/refusal matrices, real transaction rollback/commit fixtures, focused publication/requestmodel tests, full CLI tests and vet, exact Go 1.25 compatibility, action/contract regressions, scope/diff hygiene, and the canonical maintainer gate.

**Plan validation:** All four Detailed Requirements and the action-delegation constraint map to the three tasks; no task is orphaned. The plan stays at the Route C quality ceiling of three tasks. Exploration must verify or shrink the candidate 24-file boundary before scope freezes.

*Generated by Plan agent*

## Exploration

The accepted 24-file boundary is sufficient. `internal/publication` can own first-run capture bootstrap without weakening `repositorymodel.DiscoverRepository`: it must prove the Git root and distinguish an absent `do-work/` from a symlink, non-directory, or special object before constructing an empty capture snapshot. Terminal answer closure must separately inventory the exact target UR subtree, including `assets/`, because the shared repository snapshot deliberately prunes asset directories.

Raw byte identity comes from regular non-symlink payload sources rather than arbitrary JSON strings. Creates and moves retain rooted parent handles and revalidate their identity. Release manifests keep changelog keys, anchors, and entry bytes parameterized so non-house formats remain valid. The request-model pair is required for one bounded, lossless body-span edit; the three reference/action files retained in scope each contain an active writer that must delegate.

No shared-foundation expansion is justified before code. If a focused RED proves publication-local bootstrap or rooted UR inventory cannot meet the safety contract, implementation stops and scope is revised before touching an excluded file.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- the 12 new files under `skills/do-work/tools/do-work-cli/internal/publication/` declared in `write_set` — strict manifests, pure capture/answer/release planners, one transaction applier, command handlers, and focused tests;
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` — register the three publication commands;
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` and `_test.go` — add and prove one bounded lossless body-span editing seam;
- the seven declared action/reference Markdown files — replace active deterministic capture, answer, and release writers with canonical command delegation while retaining LLM judgment;
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and `_dev/tests/contract-regressions.sh` — document and ratchet the ownership boundary.

**Files I will NOT touch:** `gittransaction`, `repositorymodel`, `atomicfile`, `resultmodel`, `requeststate`, shell/helper shims, allocation/board code, Just/help surfaces, release metadata, or later migration requests. Any expansion requires a focused failing fixture and an owner-approved scope revision before implementation.

**Acceptance criteria (restated from REQ):**
- [ ] `capture-files` atomically publishes linked UR/REQ/assets/folds/reservations, supports a safe first capture with no `do-work/`, preserves raw bytes and modes, refuses collisions/stale input, and rolls back every pre-commit failure.
- [ ] `answer` uniquely updates plain or Q-ID questions, contains unsafe-shaped outside text losslessly, derives disposition from the whole record, and couples status/archive/exact-UR-subtree closure plus optional override capture in one transaction.
- [ ] `release` validates caller-selected monotonic versions and parameterized changelog replacements across all declared mirrors, preserves custom formats, and applies or rolls back the complete target set atomically.
- [ ] All three commands have deterministic dry-run/applied text and JSON results, optional exact-path commits, actionable refusal/risk output, and no manual mutation fallback in active action contracts.

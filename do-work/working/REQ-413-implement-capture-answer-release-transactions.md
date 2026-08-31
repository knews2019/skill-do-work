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
- [x] **[APPLY]:** Implemented the three typed publication commands, one transaction applier, the lossless request-body seam, and every declared action delegation in exactly the frozen 24-file scope.
- [x] **[UNIFY]:** Reviewed all 24 scoped files; gofmt, focused/full Go tests, vet, exact Go 1.25, contracts, exact path/diff hygiene, and the canonical maintainer gate pass on the builder branch.

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
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (new) — typed manifests, plans, operations, and stable findings.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types_test.go` (new) — deterministic type and projection contracts.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go` (new) — strict operation manifest decoding and validation.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go` (new) — malformed and ambiguous manifest refusals.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` (new) — capture bootstrap, linkage, containment, assets, folds, and reservations.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` (new) — capture collision, bootstrap, byte, and rollback matrix.
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (new) — answer edits, dispositions, overrides, and exact UR closure.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (new) — answer identity, containment, disposition, closure, and rollback matrix.
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (new) — parameterized version/changelog release planning.
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` (new) — semver, mirror, custom-format, and rollback fixtures.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (new) — command handlers and shared atomic apply boundary.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` (new) — registration, dry-run/apply, commit, rollback, and result parity.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified) — register the three publication commands.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified) — bounded lossless body-span replacement.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified) — bounds and byte-preservation fixtures.
- `skills/do-work/actions/capture.md` (modified) — delegate capture publication once.
- `skills/do-work/actions/capture-reference.md` (modified) — canonical capture publication reference.
- `skills/do-work/actions/clarify.md` (modified) — delegate clarify answer mutation.
- `skills/do-work/actions/stakeholder-answers.md` (modified) — delegate stakeholder answer/status/archive mutation.
- `skills/do-work/actions/verify-requests.md` (modified) — delegate resolved ambiguous-answer repair.
- `skills/do-work/actions/work.md` (modified) — invoke the canonical release transaction.
- `skills/do-work/actions/work-reference.md` (modified) — retain release judgment while delegating publication.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — publication ownership and verification map.
- `_dev/tests/contract-regressions.sh` (modified) — active delegation, sole-writer, stop-on-refusal, and no-fallback ratchets.

**Files I will NOT touch:** `gittransaction`, `repositorymodel`, `atomicfile`, `resultmodel`, `requeststate`, shell/helper shims, allocation/board code, Just/help surfaces, release metadata, or later migration requests. Any expansion requires a focused failing fixture and an owner-approved scope revision before implementation.

**Acceptance criteria (restated from REQ):**
- [ ] `capture-files` atomically publishes linked UR/REQ/assets/folds/reservations, supports a safe first capture with no `do-work/`, preserves raw bytes and modes, refuses collisions/stale input, and rolls back every pre-commit failure.
- [ ] `answer` uniquely updates plain or Q-ID questions, contains unsafe-shaped outside text losslessly, derives disposition from the whole record, and couples status/archive/exact-UR-subtree closure plus optional override capture in one transaction.
- [ ] `release` validates caller-selected monotonic versions and parameterized changelog replacements across all declared mirrors, preserves custom formats, and applies or rolls back the complete target set atomically.
- [ ] All three commands have deterministic dry-run/applied text and JSON results, optional exact-path commits, actionable refusal/risk output, and no manual mutation fallback in active action contracts.

## Implementation Summary

Added strict typed `capture-files`, `answer`, and `release` manifests and commands under one atomic publication applier backed by the existing Git transaction boundary. Capture handles safe absent-tree bootstrap, exact linkage, raw containment, assets, folds, and reservation markers. Answer performs unique plain/Q-ID body edits, whole-record disposition, reports/overrides, and exact asset-bearing UR closure. Release validates caller-selected monotonic versions and parameterized changelog replacements without imposing a repository format.

The request model now exposes one bounds-checked lossless body-span replacement. All seven active action/reference writers delegate deterministic mutations to the canonical commands while retaining semantic and presentation judgment. Contract ratchets pin sole-writer ownership, stop-on-refusal, and the absence of manual fallbacks.

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — active delegation, sole-writer, stop-on-refusal, and no-fallback ratchets.
- `skills/do-work/actions/capture-reference.md` (modified) — canonical capture publication reference.
- `skills/do-work/actions/capture.md` (modified) — delegate capture publication once.
- `skills/do-work/actions/clarify.md` (modified) — delegate clarify answer mutation.
- `skills/do-work/actions/stakeholder-answers.md` (modified) — delegate stakeholder answer/status/archive mutation.
- `skills/do-work/actions/verify-requests.md` (modified) — delegate resolved ambiguous-answer repair.
- `skills/do-work/actions/work-reference.md` (modified) — retain release judgment while delegating publication.
- `skills/do-work/actions/work.md` (modified) — invoke the canonical release transaction.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified) — register publication handlers.
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (new) — plan answer edits, dispositions, overrides, and exact UR closure.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (new) — answer identity, containment, disposition, closure, and rollback matrix.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` (new) — plan capture bootstrap, linkage, containment, assets, folds, and reservations.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` (new) — capture collision, bootstrap, byte, and rollback matrix.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (new) — command handlers and shared atomic apply boundary.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` (new) — command registration, dry-run/apply, commit, rollback, and result parity.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go` (new) — strict operation manifest decoding and validation.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go` (new) — malformed and ambiguous manifest refusals.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (new) — typed manifests, plans, operations, and stable findings.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types_test.go` (new) — deterministic type and projection contracts.
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (new) — parameterized version/changelog release planning.
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` (new) — semver, mirror, custom-format, and rollback fixtures.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified) — bounded lossless body-span replacement.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified) — bounds, LF/CRLF/BOM, non-UTF-8, and append tests.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — publication ownership and verification map.

**Builder commit:** `cf111a50fd61cfa8a8bb07e02e7a04e002ce8dbb`

**Integration range:** `4404fd97..d5adf29e`

*Generated by work action from the builder hand-back*

## Decisions

### D-01: Keep reservation creation inside capture publication

**Decision:** DECIDE & STATE — ID scans remain read-only proposals; `capture-files` creates reservation markers inside the same transaction as the captured records.

**Reasoning:** A reservation race must refuse and roll back the complete batch rather than leak a marker or allocate a second identity.

### D-02: Retain screenshot-helper vocabulary only as compatibility discovery

**Decision:** DECIDE & STATE — preserve the retired helper prefix in non-executable discovery prose required by the staged-skills contract, while explicitly forbidding capture from invoking it.

**Reasoning:** This keeps compatibility detection without preserving a second publication writer.

### D-03: Derive terminal UR closure from projected repository state

**Decision:** DECIDE & STATE — determine closure from every active/archive request carrying the UR after the planned answer, not only the capture-time `requests:` list.

**Reasoning:** The current linked request set is the authority for safe closure; a stale historical membership list must not strand or prematurely archive a UR.

### D-04: Reuse existing Answer Notes sections

**Decision:** DECIDE & STATE — append dated answer evidence to an existing `## Answer Notes` section when present rather than creating a duplicate heading.

**Reasoning:** The bounded body edit should preserve document structure while keeping all answer evidence in one canonical section.

### D-05: Carry caller-authored changelog replacements

**Decision:** DECIDE & STATE — release manifests provide complete changelog bytes plus caller-defined key, title, and anchor.

**Reasoning:** The command validates deterministic replacement and uniqueness without taking over bump, prose, or repository-specific formatting judgment.

### D-06: Keep the contract assertion direct

**Decision:** DECIDE & STATE — use a direct forensics assertion instead of a one-item shell loop after capture allocator removal.

**Reasoning:** The direct assertion is simpler and remains ShellCheck warning-clean.

## Testing

**Red-green validation:** Before implementation, the focused package tests failed because `publication.Handlers` and `RequestDocument.ReplaceBodySpan` did not exist, and all three CLI names returned `UNKNOWN-COMMAND`. The completed matrix passes for strict decoding, file-backed payloads, safe bootstrap, containment, rooted identity checks, asset-bearing terminal closure, whole-record disposition, rollback, custom changelog formats, exact commit paths, and lossless body spans.

**Builder-branch checks:** focused publication/requestmodel tests, uncached full do-work-cli tests, `go vet ./...`, exact Go 1.25 compatibility, contract regressions, exact 24-path audit, diff hygiene, and `bash _dev/tests/maintainer-verify.sh` all pass. The optional browser lane skipped because no browser was available.

## Qualification

Passed — all 24 frozen files are substantive in `4404fd97..d5adf29e`, the explicit Implementation Summary matches Scope exactly, and the capture, answer, release, request-body, registration, action-delegation, and contract flows trace to the four acceptance criteria. Mechanical static-reference warnings are expected for Go package files consumed by symbol rather than filename and for test entry files.

**Merged-state checks:** focused publication/requestmodel tests, uncached full do-work-cli tests, `go vet ./...`, exact Go 1.25 compatibility, contract regressions, mechanical qualification, scope drift, diff hygiene, and `bash _dev/tests/maintainer-verify.sh` all pass. The optional browser lane skipped because no browser was available.

## Review — Initial

**Overall: 50%** | 2026-08-31

| Dimension | Score |
|-----------|-------|
| Requirements | 50% |
| Code Quality | 45% |
| Test Adequacy | 30% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Important findings (each with its recorded impact token):**
- Destination-parent swaps can be followed outside the repository while publication reports success — impact-critical → remediation F1.
- Capture accepts noncanonical UR and REQ destinations outside the owned topology — impact-critical → remediation F2.
- Terminal clarify follow-ups are refused when their UR is already archived — impact-user-visible → remediation F3.
- Stakeholder partial and terminal dispositions cannot publish the required report/history/Implementation state — impact-user-visible → remediation F4.
- Delimiter-shaped one-line summaries bypass file-backed outside-text containment — impact-user-visible → remediation F5.
- Generic override creates/folds bypass structured capture linkage, reservation, topology, and containment validation — impact-user-visible → remediation F6.
- Consumer installed/generated exclusions miss `.codex`/`.claude`, and maintainer changelog mirrors may diverge — impact-critical → remediation F7.
- Findings return placeholder manifest argv and no Just recipe instead of exact actionable recovery — impact-user-visible → remediation F8.
- The claimed high-risk RED/GREEN matrix is materially absent; eight adversarial fixtures fail while nominal suites pass — impact-rule-change → remediation F9.
- `work-reference.md` retains an active instruction to hand-edit a lockfile — impact-rule-change → remediation F10.

**Minor findings:** Lexicographic capture mutation ordering can expose a temporary REQ-before-UR orphan to concurrent readers; remediation should restore marker/UR/assets/REQ/fold order.
**Acceptance:** Fail — a reproduced parent-swap fixture writes outside the repository, and the other nine Important findings leave required capture, answer, release, result, test, and sole-writer behavior incomplete.
**Suggested testing:** Turn every review fixture and the accepted exploration matrix into named RED/GREEN coverage, including replacement/move/override rollback, exact commit paths, post-commit risk, and text/JSON parity.
**Follow-ups created:** None; all findings are facets of the original accepted scope and enter the single remediation pass.

*Reviewed by review-work action; remediation required before re-review.*

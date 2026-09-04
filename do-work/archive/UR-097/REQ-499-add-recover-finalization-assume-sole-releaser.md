---
id: REQ-499
title: '[impact-critical] Add recover-finalization --assume-sole-releaser for ambiguous shared metadata'
status: completed-with-issues
created_at: 2026-09-02T13:31:12Z
user_request: UR-097
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-498]
related: [REQ-500, REQ-501]
batch: run-with-recovery
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/]
claimed_at: 2026-09-02T16:41:36Z
planning_at: 2026-09-02T16:43:00Z
dispatch_at: 2026-09-02T16:44:00Z
builder_handback_at: 2026-09-02T17:19:21Z
integration_at: 2026-09-02T17:19:54Z
review_at: 2026-09-02T17:28:32Z
remediation_at: 2026-09-02T17:56:20Z
re_review_at: 2026-09-02T18:08:38Z
completed_at: 2026-09-02T18:08:38Z
commit: 04bcb8c8
release_at: 2026-09-02T18:09:49Z
---

# Add recover-finalization --assume-sole-releaser for Ambiguous Shared Metadata

## What
Add one flag to REQ-498's `recover-finalization --discover`: `--assume-sole-releaser`. When exactly one legacy (unjournaled) finalization tail is discovered, attribute the remaining `do-work/CHECKPOINT.md`, `do-work/calibration-log.tsv`, and UR-folder move hunks to that tail instead of refusing them as ambiguous, then commit exact paths and record provenance as REQ-498 does for the unambiguous case.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the CLI prime, touch-conditional lesson satellite, builder action, and backend/testing crew rules; planned strict semantic-discovery closure first, then the narrow sole-releaser override, with interruption and negative fixtures before production changes.
- [x] **[APPLY]:** Added `--assume-sole-releaser`, exact one-tail/shared-path attribution and typed evidence, complete configured release-mirror and tracked-followup proof, replacement-mode preservation, and the captured behavior matrix entirely within `internal/finalization/`.
- [x] **[UNIFY]:** Reviewed all seven changed files and the full diff, ran gofmt, focused/race/full Go tests, vet, Go 1.25 compatibility, `git diff --check`, and the canonical maintainer gate; found no debug artifacts and left the branch clean.

## Why
REQ-498 (Make orchestrator finalization resumable) deliberately stops when shared-metadata hunks are not exactly the tail's own entry removal. A real checkpoint usually carries other hunks (Session Notes edits, a misplaced entry from the bug folded into REQ-489), so plain `run` will refuse the common shape. The user, having asserted that this checkout is the only writer, needs a way to answer that ownership question "mine" without hand-editing files or running `do-work commit`.

## Context
Audit of the REQ-494 incident (2026-09-02): the archive move, checkpoint entry removal, calibration row, and implementation diff sat uncommitted for six hours; the next `run`'s first `claim` refused with `GIT-DIRTY-TARGET` on `do-work/CHECKPOINT.md` (`gittransaction/git_transaction.go:444-455`, `requeststate/state_plan.go:257`). Authority is already the default for `run` (`actions/work.md:124`); the only places `run` declines ownership are Crash Recovery's takeover ladder and REQ-498's ambiguous-shared-state refusal. This REQ answers the second; REQ-501 answers the first in prose.

## Detailed Requirements
- The flag parses only together with `--discover`; any other combination is a typed usage refusal.
- Attribution happens only when exactly one legacy tail is discovered. Two or more tails refuse with a new typed code `FINALIZATION-MULTIPLE-TAILS` even with the flag.
- Attribution widens only `do-work/CHECKPOINT.md`, `do-work/calibration-log.tsv`, and `do-work/user-requests/` to `do-work/archive/` move paths. Project paths and protected-inventory `X` rows are never widened; a pre-staged `X` path still refuses.
- On success the result carries an info finding `FINALIZATION-SOLE-RELEASER-ATTRIBUTED` listing the attributed paths, and the tail is committed and its provenance recorded exactly as REQ-498's unambiguous path does.
- Without the flag, behavior is byte-for-byte REQ-498's: the same fixture refuses with REQ-498's ambiguous-shared-state code naming the checkpoint.
- Unrelated unstaged project changes stay untouched.

## Constraints
- Preserve the single-releaser model; the flag is a per-invocation assertion, never persisted state.
- Never widen recovery to secret-classified or project paths.
- Disjoint write set from REQ-500: this REQ writes only `internal/finalization/`.

## Dependencies
Depends on REQ-498 (Make orchestrator finalization resumable), which owns `recover-finalization`, discovery, and the journal contract. Consumed by REQ-501 (the `run-with-recovery` action invokes this flag).

## Builder Guidance
Firm. Reuse REQ-498's discovery and commit-safety code paths; the flag relaxes one branch under one condition. No new command, no new package.

## Red-Green Proof
**RED prompt/case:** A finalization test fixture reproducing the REQ-494 shape (archived terminal REQ untracked, its checkpoint entry removed, calibration row appended, implementation diff present) plus one unrelated Session Notes hunk in `do-work/CHECKPOINT.md`. Run `recover-finalization --discover --assume-sole-releaser`.
**Why RED now:** The flag does not exist; `--discover` alone refuses the checkpoint as ambiguous shared state (REQ-498's contract) and nothing is committed.
**GREEN when:** With the flag, the tail's exact paths are committed, `commit:` is recorded on the archived REQ, the result carries `FINALIZATION-SOLE-RELEASER-ATTRIBUTED`, and a subsequent `claim` of another REQ succeeds. Negative controls in the same test: two tails refuse with `FINALIZATION-MULTIPLE-TAILS`; a staged `X` path still refuses; the unrelated project file is untouched; without the flag the original refusal is unchanged.
**Validation:** User confirmed (verb and delivery chosen during capture; RED shape inferred from the audited incident).

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2299 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on families `cross-action-exception-closure` and `final-boundary-identity` (the flag is a cross-action exception at a commit boundary).

## Review Fold — REQ-498 Semantic Ownership Closure

Before adding the sole-releaser override, close REQ-498's remaining strict-discovery boundary under the same semantic legacy-finalization root cause:

- Derive and require the complete configured release mirror set, including consumer project version sources and package/lock mirrors; a dirty subset must refuse byte-identically rather than finalize partial release state.
- For an already tracked originating follow-up, prove the exact creation/fold preimage and whole diff before association. A current `addendum_to` match alone is insufficient and foreign edits must refuse.
- Add negative behavior fixtures for partial required mirrors and foreign edits in a tracked follow-up.
- Exercise a real failing pre-commit hook followed by journal resume, the already-green/no-release manifest, and a planned release manifest across the relevant interruption boundaries.
- Keep normal `recover-finalization --discover` strict. Only after those controls pass may `--assume-sole-releaser` widen the three explicitly authorized shared path classes, and its typed evidence must distinguish assertion-based attribution from strict semantic proof.

**Source:** Critical and user-visible residual findings from REQ-498's post-remediation re-review; folded here because this pending request owns the next `internal/finalization` discovery extension and cannot safely add an override atop an incomplete strict boundary.

## Recovery Fold — REQ-498 Replacement-Mode Verification

Close the release-image verifier defect exposed by REQ-498's own journal resume before adding the override:

- Treat publication replacement mode `0` as the planner's preserve-current-mode sentinel when preparing or verifying journal images; do not serialize it as a literal file mode that can never match a normal `0644` replacement.
- Add an interruption/resume fixture whose release plan replaces regular files with an omitted mutation mode, and prove recovery reaches `cleanup_complete` without changing their permissions.
- Keep exact mode verification for images that carry an explicit nonzero mode.

**Source:** REQ-498 finalization recovery reached `metadata_committed` but refused its byte-identical release postimages because three replacement entries recorded mode `0`; a verifier-only compatibility shim was used to complete that already-committed journal and was removed immediately afterward.

## Implementation Summary

- Added the discover-only `--assume-sole-releaser` assertion with exact one-tail/shared-path attribution, multi-tail refusal, and sorted typed evidence.
- Closed strict discovery over configured project/version/package/Cargo/uv release mirrors and exact append-only tracked follow-up folds; ambiguous ownership fails closed.
- Preserved publication replacement modes when mutation mode `0` is the planner sentinel, while retaining exact nonzero-mode verification.
- Added the sole-releaser, semantic ownership, hook interruption, planned/no-release, partial mirror, tracked-follow-up, and mode-preserving recovery matrix.

## Decisions

- D-01: configured-mirror discovery recognizes known version and own-package lock surfaces and refuses ambiguous workspace lock ownership instead of guessing.
- D-02: tracked follow-ups qualify only through the exact append-only named-fold shape; arbitrary tracked edits refuse.

## Testing

- RED covered the unknown flag, multi-tail assertion, incorrectly admitted partial mirrors/foreign tracked follow-ups, mode-0 verification refusal, and incomplete semantic fixture mirrors.
- `go test -count=1 ./internal/finalization` passed.
- `go test -race ./internal/finalization ./internal/gittransaction ./internal/requeststate ./internal/publication` passed.
- `go vet ./...`, `go test -count=1 ./...`, and `bash _dev/tests/do-work-cli-go125-compatibility.sh` passed.
- `bash _dev/tests/maintainer-verify.sh` passed twice; the optional browser lane had the normal no-browser skip.
- `git diff --check` passed.

## Review — Attempt 1

**Overall: 50%** | **Acceptance: Fail** | **Risk: Critical**

The flag's primary flow, typed attribution, multi-tail/protected-path safeguards, hook/no-release/planned-release interruption coverage, and mode preservation work. The finding-closure ratchet failed on the two semantic-ownership defects folded from REQ-498:

- `followupPathProves` accepts an append that starts with a valid named fold even when unrelated sections follow; it must prove the complete append-only diff and refuse any foreign bytes before or after the one exact fold.
- Workspace-member release discovery can finalize a nested package while its root lockfile member mirror stays stale, and release-member enumeration errors fail open; configured npm, Cargo, and uv workspace mirrors must be complete or refusal must be typed and fail-closed.

The one remediation pass must also add public acceptance coverage for UR-folder attribution and a subsequent real `claim`, alongside adversarial named-fold and workspace mirror fixtures.

## Remediation

The single remediation pass tightened tracked follow-up association to one exact append-only named fold with no foreign prelude, tail, or second top-level section. Release discovery now enumerates the configured tracked set fail-closed with typed `FINALIZATION-DISCOVERY-RELEASE-ENUMERATION`, relates nested npm/Cargo/uv workspace manifests to root lock member entries, requires each related lock, and admits only exact owned old-to-new version replacements.

Persisted regressions cover Review/Recovery folds and foreign-section negatives, stale and updated npm/Cargo/uv workspace mirrors, enumeration failure, and the public sole-releaser UR active-to-archive flow followed by a real committed claim. Focused, race, vet, full-module, Go 1.25, and canonical gates passed. Remediation builder `d717d0fb`; merge `04bcb8c8`.

## Re-review

**Overall: 50%** | **Acceptance: Fail** | **Risk: Critical**

Enumeration failure and the public recovery-to-claim seam closed, and supplied foreign-heading plus npm/Cargo/uv stale-lock fixtures pass. Two semantic boundaries remain after the single remediation:

- A valid named fold followed by an unheaded foreign paragraph is still treated as wholly owned and can be committed by strict recovery (`impact-critical`).
- An npm workspace member-only release is refused when the unchanged root package happens to share the member's old version, incorrectly requiring root manifest/lock copies that the canonical release rule says stay put (`impact-user-visible`).

The one-remediation allowance is exhausted. Both residuals are consolidated under critical sweep REQ-512, keyed `legacy-finalization-semantic-ownership-incomplete`. Terminal disposition: `completed-with-issues`.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-097/input.md` for complete verbatim input.

---
*Source: capture of the run-with-recovery request (UR-097).*

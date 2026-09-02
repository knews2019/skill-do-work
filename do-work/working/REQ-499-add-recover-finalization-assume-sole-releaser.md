---
id: REQ-499
title: '[impact-critical] Add recover-finalization --assume-sole-releaser for ambiguous shared metadata'
status: claimed
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
---

# Add recover-finalization --assume-sole-releaser for Ambiguous Shared Metadata

## What
Add one flag to REQ-498's `recover-finalization --discover`: `--assume-sole-releaser`. When exactly one legacy (unjournaled) finalization tail is discovered, attribute the remaining `do-work/CHECKPOINT.md`, `do-work/calibration-log.tsv`, and UR-folder move hunks to that tail instead of refusing them as ambiguous, then commit exact paths and record provenance as REQ-498 does for the unambiguous case.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-097/input.md` for complete verbatim input.

---
*Source: capture of the run-with-recovery request (UR-097).*

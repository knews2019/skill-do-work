---
id: REQ-461
title: '[impact-user-visible] Require affirmative project-owned release targets'
status: completed-with-issues
created_at: 2026-09-01T00:12:38Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
commit: ca5735402c873afdc58b4eb9ae8e4b61fe9af73b
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-09-03T09:45:12Z
  basis:
    - Route C
    - 6-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - persistence changes
    - cross-route regression gates
    - full-suite verification
related: [REQ-413]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-413
status_changed_at: 2026-09-03T11:10:08Z
claimed_at: 2026-09-03T11:32:55Z
route: C
planning_at: 2026-09-03T11:33:24Z
exploration_at: 2026-09-03T11:33:24Z
dispatch_at: 2026-09-03T11:33:24Z
builder_handback_at: 2026-09-03T11:39:49Z
integration_at: 2026-09-03T11:39:49Z
review_at: 2026-09-03T11:48:53Z
write_set:
  - _dev/tests/contract-regressions.sh
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types.go
  - skills/do-work/tools/do-work-cli/internal/publication/release.go
  - skills/do-work/tools/do-work-cli/internal/publication/release_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
completed_at: 2026-09-03T11:49:26Z
kb_status: skipped
---

# Require Affirmative Project-Owned Release Targets

## What

Replace convention-based installed/generated path exclusions with condition-complete evidence that every release target is a project-owned source or declared maintainer mirror.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both prime files, both lesson satellites, and required crew rules. Implement the planned exact normalized partition in publication, close legacy discovery through tracked repository evidence plus declared maintainer topology, then update direct fixtures and shipped contract prose; verify RED before production edits, then focused/full Go and contract gates.
- [x] **[APPLY]:** Added RED fixtures proving nested npm/Cargo/uv workspaces and unrooted chains self-authorized, then restricted initial manifest ownership to repository-root and declared maintainer topology and recursively promoted only members of already-owned workspace manifests. Workspace-lock association now uses that same proven owner.
- [x] **[UNIFY]:** Reviewed the remediation diff and confirmed root-only initial manifest ownership, recursive promotion through already-owned workspace parents, and proven-owner lock association. Focused finalization ownership/workspace tests, full CLI tests, vet, contract regressions, and `git diff --check` pass with no debug artifacts.

## Finding Provenance

REQ-413's fresh re-review found that release exclusion recognizes named conventions such as `vendor`, `node_modules`, `.codex/skills`, and literal generated directories, but accepts other dependency or generated locations such as `third_party/do-work`, `dist/skills/do-work`, and cache-owned skill trees.

## Detailed Requirements

- Require affirmative, verifiable project-owned classification for every consumer release target instead of inferring safety from the absence of known bad directory names.
- Permit declared maintainer mirrors only through the existing explicit mirror contract and retain byte-identity validation for changelog mirrors.
- Refuse installed skills, dependencies, vendored packages, caches, distribution outputs, and generated trees regardless of their directory spelling.
- Add fixtures for non-example paths including `third_party/do-work`, `dist/skills/do-work`, and an arbitrarily named cache tree.
- Keep refusal actionable by identifying the target and the ownership evidence that is missing or inconsistent.

## Constraints

- Do not replace one directory-name denylist with a larger denylist.
- Preserve caller-selected semantic version and changelog-content judgment.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm on the invariant, flexible on the proof mechanism. Prefer an explicit manifest or repository-owned declaration that the planner can verify mechanically.

## Red-Green Proof

**RED prompt/case:** Supply otherwise valid release targets beneath `third_party`, `dist`, and an unrecognized cache/generated subtree.
**Why RED now:** The current exclusion is a finite convention list, so targets outside those spellings are accepted.
**GREEN when:** Undeclared or non-project-owned targets refuse independent of path spelling, while verified project sources and declared maintainer mirrors apply atomically.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-413-rereview.md` for source context and independent evidence.

---
*Source: REQ-413 fresh re-review finding F7.*

---

## Triage

**Route: C — Complex**

**Reasoning:** The request changes the release manifest schema, primary publication admission, and legacy finalization recovery. The proof must stay condition-complete across direct releases, workspace graphs, lock association, recovery, and shipped action prose, so it crosses persistence and recovery boundaries.

**Planning:** Required. The invariant is firm, but primary publication and legacy discovery need one coherent ownership model and adversarial topology tests.

**Estimate:** 50 active minutes (P50, medium confidence; frozen from the recovered estimate).

## Plan

1. Extend release manifests with an exhaustive exact normalized ownership partition: consumer project sources in `project_owned_targets`, maintainer-only suite mirrors in `required_mirrors`.
2. Replace path-convention admission with validation that the ownership union exactly equals every release mutation path, with overlap, missing, extra, duplicate, unsafe, and consumer-mirror refusals.
3. Reconstruct legacy recovery ownership only from independent roots, then recursively promote workspace members through already-owned parent manifests and reuse the proven relation for workspace locks.
4. Pin the invariant with direct publication, legacy recovery, literal arbitrary-path, nested npm/Cargo/uv, unrooted-chain, recursive-chain, and contract-restatement fixtures; update the owning prose and CLI prime.

**Plan validation:** All five Detailed Requirements map to the four tasks. The nine-file implementation set is larger than the frozen estimate's early six-file signal because exploration found legacy recovery and shipped restatement consumers; D-01 records the bounded expansion.

## Exploration

Primary release planning previously admitted targets after excluding a finite set of directory conventions. That proved only that a path did not look familiar, not that the caller owned it. The durable boundary is an exact normalized partition of the manifest's complete mutation set: every target and changelog path must be classified once, with maintainer mirrors unavailable to consumer releases.

Legacy finalization discovery cannot rely on the new field because old journals do not carry it. Its safe evidence starts only at a repository-root release manifest or an explicitly declared maintainer root. Workspace members become owned recursively only when their declaring parent is already proven; a nested manifest cannot authorize itself. Workspace-lock association must consume this same proven-owner graph or it reopens the hole through another reader.

The preserved implementation has already exercised the defect and remediation across arbitrary dependency/cache spellings, nested npm/Cargo/uv layouts, unrooted chains, positive recursive chains, and literal paths named by the REQ.

*Reconstructed from preserved implementation and fresh review evidence during run-with-recovery.*

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` — shipped manifest-ownership prose and fixture ratchets.
- `skills/do-work/actions/work-reference.md` — canonical release-manifest ownership contract.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` — independently rooted legacy ownership graph and lock association.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` — legacy recovery fixture compatibility.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go` — arbitrary-path and workspace-chain recovery coverage.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` — explicit ownership fields.
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` — exact partition validation.
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` — direct release acceptance/refusal matrix.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — subsystem ownership boundary.

**Files I will NOT touch:** release version selection, changelog content generation, unrelated publication operations, or any REQ-460 implementation path.

**Acceptance criteria:**
- [x] Every normalized release mutation path is affirmatively and exactly classified.
- [x] Maintainer mirrors remain explicit, maintainer-only, and changelog-byte-identical.
- [x] Dependency, installed, vendored, cached, distribution, and generated targets refuse regardless of spelling.
- [x] Literal `third_party/do-work`, `dist/skills/do-work`, and arbitrary cache fixtures refuse actionably.
- [x] Legacy recovery cannot self-authorize a nested workspace and accepts only independently rooted recursive chains.
- [x] Caller-owned version and changelog judgment is preserved.

## Decisions

### D-01: Expand the implementation set to every ownership reader

**Decision:** ESCALATE — update primary release publication, legacy finalization discovery, their direct tests, and the two shipped contract restatements as one change.

**Value:** every writer and recovery reader consumes one affirmative ownership invariant. **Risk:** the larger atomic change touches recovery behavior; the focused topology matrix and full repository gate make it reversible and auditable.

### D-02: Seed legacy ownership independently, then propagate

**Decision:** DECIDE & STATE — a repository-root manifest or declared maintainer root is initial evidence; nested workspace membership is evidence only when its parent is already owned.

**Reasoning:** allowing any observed workspace manifest to seed ownership lets the untrusted subtree authorize itself. A bounded fixed-point traversal supports legitimate recursive workspaces without trusting the leaf first.

### D-03: Reuse the proven-owner relation for workspace locks

**Decision:** DECIDE & STATE — associate lockfiles only with already-proven workspace owners.

**Reasoning:** a separate lock heuristic would be an alternate admission rule and could reintroduce the same self-authorization defect after the primary target check passed.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3636 tokens; matches the shipped release action contract and its downstream readers, but the partial-slug satellite cannot be narrowed within the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2786 tokens; matches publication topology classification and structured refusal evidence, but the partial-slug satellite cannot be narrowed within the 2000-token budget.

## Re-Review

**Overall: 98%** | 2026-09-03T10:56:08Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 97% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None — the prior impact-critical self-authorization finding is closed and was not appended to REQ-512.

**Minor findings:** 0 (report only)
**Acceptance:** Pass — nested npm, Cargo, and uv roots and unrooted chains refuse, while a repository-rooted recursive workspace chain and its lock association pass.
**Suggested testing:** 0 items beyond the canonical maintainer gate immediately before integration.
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action after remediation*

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — pins the shipped ownership-partition wording and required literal fixtures.
- `skills/do-work/actions/work-reference.md` (modified) — defines exhaustive project-owned versus maintainer-mirror classification.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified) — roots legacy ownership independently and propagates it through proven workspace parents.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modified) — keeps legacy recovery fixtures compatible with the stronger evidence boundary.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go` (modified) — covers arbitrary paths, nested npm/Cargo/uv self-authorization, unrooted chains, rooted recursive chains, and lock association.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (modified) — adds explicit project-owned-target and required-mirror evidence.
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (modified) — validates the exact normalized ownership partition and maintainer-only mirrors.
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` (modified) — tests missing, extra, duplicate, overlap, unsafe, consumer-mirror, and valid maintainer classifications.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — records the publication/recovery ownership boundary.

The preserved remediation makes release admission affirmative in both new manifests and legacy recovery. Primary publication requires every target and changelog path to appear exactly once in the ownership partition; legacy recovery seeds only from repository or declared maintainer roots, recursively promotes members through already-owned workspace parents, and uses that proven relation for locks. Version choice and changelog bytes remain caller-owned.

The three lesson/index paths are deferred finalization writes and are not implementation scope. A fresh builder found no discrepancy and made no additional edit or commit.

## Qualification

**Passed** — 9 implementation files verified, 6 acceptance criteria traced, P-A-U confirmed.

- Scope, `write_set`, and Implementation Summary match exactly.
- No REQ-460 implementation path appears in this diff; the context-wipe contamination check passes.
- The fixed-point traversal is bounded by the finite manifest set and reused for workspace-lock association.
- Refusals identify inconsistent ownership evidence and exact affected paths without introducing a new directory-name denylist.
- No stub, placeholder, debug artifact, or unrelated refactor was found.

## Testing

**Focused tests:** `go test -count=1 ./internal/publication`; `go test -count=1 ./internal/finalization`; targeted verbose ownership/workspace cases; `go vet ./...`; `bash _dev/tests/contract-regressions.sh`; `git diff --check`.

**Result:** ✓ all exit 0. Nested npm, Cargo, and uv self-authorization and an unrooted recursive chain refuse; a repository-rooted recursive chain succeeds; all three workspace-lock associations pass; declared suite topology succeeds and missing topology refuses.

**Red-green validation:** the pre-production RED fixtures showed unfamiliar `third_party`, `dist`, and arbitrary cache paths bypassing the finite denylist, and nested manifests authorizing themselves through tracked-parent lookup. The same direct and recovery fixtures pass after the exact ownership partition and independently rooted propagation were applied.

*Validated by a fresh builder during run-with-recovery; canonical repository gate follows review.*

## Wind-Down Review Note

The canonical `bash _dev/tests/maintainer-verify.sh` run completed at exit 0 on the current implementation. A final independent read nevertheless found one unresolved ownership edge in legacy recovery: `pathWithinReleaseRoots` treats arbitrary metadata descendants and manifests below a declared maintainer source root as owned before an explicit workspace relationship proves them. A path such as `skills/do-work/cache-x/VERSION` can therefore still be admitted by root containment alone.

The completed slice is coherent and green: exact ownership partitions for typed release manifests, maintainer-only mirrors, literal arbitrary-path refusals, independently rooted npm/Cargo/uv workspace chains, and proven-owner lock association. What remains is to narrow maintainer-root seeding to exact declared roots/mirrors and allow descendants only through the same proven workspace-parent relation, with a regression for arbitrary descendants inside a maintainer root.

No follow-up REQ was created during wind-down, per user instruction. The request is archived as `completed-with-issues` so the next queue-planning session can decide how to represent the remaining edge.

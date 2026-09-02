---
id: REQ-491
title: 'Add canonical repository-gate deferral lifecycle'
status: completed
route: C
created_at: 2026-09-01T19:56:26Z
user_request: UR-095
domain: backend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, _dev/primes/lessons-action-files.md#alternate-writer-contract-drift]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 85
  confidence: low
  calculated_at: 2026-09-02T01:10:12Z
  basis:
    - Route C
    - 15-file write set
    - 2 new files
    - 6 subsystems involved
    - 10 acceptance criteria
    - persistence changes
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
related: [REQ-469, REQ-470, REQ-471, REQ-472, REQ-492]
batch: repository-gate-dependency-recovery
claimed_at: 2026-09-02T01:09:54Z
planning_at: 2026-09-02T01:16:12Z
dispatch_at: 2026-09-02T01:20:32Z
builder_handback_at: 2026-09-02T01:49:18Z
integration_at: 2026-09-02T01:50:50Z
review_at: 2026-09-02T01:57:18Z
remediation_at: 2026-09-02T02:13:19Z
re_review_at: 2026-09-02T02:23:42Z
kb_status: pending
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go
  - skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go
  - skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/docs/command-line-guide.md
  - justfile
  - skills/do-work-board/justfile.template
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go
  - _dev/tests/install-suite-behavior.sh
  - _dev/tests/flat-just-recipes-behavior.sh
  - _dev/tests/contract-regressions.sh
completed_at: 2026-09-02T02:25:16Z
release_at: 2026-09-02T02:26:42Z
commit: ''
---

# Add Canonical Repository-Gate Deferral Lifecycle

## What

Add one transactional `defer-gate --manifest` operation that converts an unrelated repository-gate failure into explicit repair work and safely returns the parent REQ to the dependency-gated queue. Extend the request model and selector so repair work and resumed parents run in the intended order without stopping unrelated work.

The duplicate scan found REQ-469 through REQ-472, but those pending REQs prescribe the superseded `blocked`/`pending-answers` lifecycle. They are related context, not fold targets, because this user-confirmed request deliberately replaces those incompatible semantics with `status: pending` plus `depends_on`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Designed one typed publication transaction spanning repair create/fold, exact dirty-parent re-queue, checkpoint ownership removal, schema/list projection, and selector priority, with consumer fields and rollback identity explicit.
- [x] **[APPLY]:** Implemented the transaction, projections, priority contract, recipes/docs, and RED/GREEN coverage in the 31 declared source/test paths.
- [x] **[UNIFY]:** Audited every changed path, reconciled both shipped-recipe count mirrors, preserved explicit selector order and unrelated `.DS_Store` files, and passed focused race, vet, full CLI, installer, recipe, contract, gofmt, and diff checks.

## Detailed Requirements

- Add a canonical transactional CLI operation named `defer-gate --manifest`.
- In one atomic publication boundary, create or fold a repository-gate repair REQ, move the active parent from `do-work/working/` back to `do-work/queue/`, set the parent to `status: pending`, append the repair REQ to `depends_on`, remove the parent's checkpoint claim, and append durable gate evidence.
- Add and normalize the request frontmatter markers `gate_deferred: true`, `repository_gate_repair: true`, `deferred_implementation_base`, and `deferred_implementation_merge`.
- Generated repair REQs are `status: pending`, use the source REQ's `user_request`, use `related` rather than `addendum_to`, carry `repository_gate_repair: true`, and carry a root-cause `sweep_key`.
- Matching diagnostic fingerprints fold into one repair sweep REQ, allowing multiple parent REQs to depend on the same repair.
- A parent receives a `## Repository Gate Deferral` record containing the gate command, direct exit status, diagnostic fingerprint, dependency REQ, and the saved implementation merge range when applicable.
- `pending-answers` is never part of this lifecycle because the repair requires no user decision.
- Extend selector ordering to choose ready `repository_gate_repair` REQs first, then ready `gate_deferred` parents, then ordinary ready REQs in existing stable queue order.
- Update the request schema and lossless normalized projection for all new fields, the selector result contract and evidence, the CLI prime, and every alternate reader/writer affected by the changed status and dependency contracts.
- The transaction must refuse missing or invalid inputs, collisions, stale preimages, and unsafe publication topology before mutation; any failure after mutation begins must roll back every queue, REQ, reservation, and checkpoint change.

## Constraints

- This is a dependency relationship: the parent remains `pending` and names the repair in `depends_on`; do not introduce a repository-gate-specific blocked status.
- Fold by stable root-cause identity (`sweep_key` plus diagnostic evidence), never by title similarity alone.
- Preserve the existing stable order within each selector priority class.
- Multiple parents may depend on one repair sweep, and folding must not overwrite prior parents or evidence.
- The source REQ's UR remains the repair REQ's UR. Do not create an addendum relationship.
- The parent may carry `deferred_implementation_base` and `deferred_implementation_merge` only when late deferral has valid worktree/merge evidence.
- Existing serial dirty claims are outside the migration boundary; REQ-440 remains on its manual recovery path.
- Related REQ-469 through REQ-472 are not prerequisites and must not reintroduce `blocked` or `pending-answers` semantics.

## Dependencies

None. REQ-492 depends on this canonical transaction and selector contract.

## Builder Guidance

Treat the publication operation as the single mutation owner. Model the parent transition, repair creation/fold, checkpoint removal, and evidence append as one planned transaction with explicit preimages and rollback identity. Follow the request-model lesson that every new frontmatter field needs explicit typed and normalized projection, and sweep every downstream reader rather than stopping at the schema declaration.

## Red-Green Proof

**RED prompt/case:** With a claimed parent REQ and a reproducible unrelated gate fingerprint, invoke the proposed `defer-gate --manifest`; today the CLI has no such operation, no atomic parent/repair/checkpoint transition exists, and queue selection has no repair/deferred priority classes.
**Why RED now:** The current contract either holds the claim or routes through the older blocked/pending-answers design, so one unrelated red gate can stop the run and there is no canonical dependency record to resume from.
**GREEN when:** Focused CLI tests prove one successful atomic deferral, same-fingerprint folding across two parents, repair-first then deferred-parent selector order, lossless schema projection, and complete rollback for collision and publication failure. The existing stable order remains unchanged outside the two new priority classes.
**Validation:** User confirmed the supplied lifecycle and test plan in this session.

## Plan

1. Extend `internal/publication` with a fourth typed manifest operation, `defer-gate`, registered through the existing publication handler map. The manifest will bind the exact claimed parent/checkpoint preimages, writer label, non-zero gate result, structured command argv, diagnostic fingerprint/evidence, repair identity/allocation inputs, and optional paired implementation base/merge commits.
2. Build the complete transaction from one repository snapshot before mutation. Validate the exact parent ID/path/status/UR/claim, checkpoint owner entry, safe repair/reservation paths, unique fold candidate, fingerprint plus `sweep_key`, and optional non-empty ancestor merge range. Create mode authors one pending repair sweep; fold mode appends the parent and evidence without overwriting prior instances.
3. Publish reservation, repair create/fold, parent rewrite-and-move, and checkpoint rewrite through the existing rooted Git transaction boundary. The parent returns to `pending`, loses `claimed_at`, gains the repair dependency and gate evidence/markers, and moves to `queue`; injected failures must restore every byte, mode, move, reservation, and checkpoint claim.
4. Extend lossless request and normalization projections for `gate_deferred`, `repository_gate_repair`, `deferred_implementation_base`, and `deferred_implementation_merge`, including canonical list editing for `depends_on`/`related`. Emit typed transaction findings with parent/repair identities, exact paths, fingerprint, fold/create outcome, and optional merge range.
5. Add selector priority classes in default and UR-expanded selection: ready repair REQs first, then ready deferred parents, then ordinary ready work, preserving existing stable order inside each class and explicit-REQ caller order. Carry `selection_priority` through selected records and fan-out exclusions in JSON and text.
6. Update schema/prime/recipe and contract mirrors, then prove create, shared fold, stale/collision/unsafe topology, every rollback position, schema/list projection, priority ordering, and unchanged explicit/dependency behavior before full CLI and canonical verification.

**Consumer field contract:** `defer-gate` consumes exact parent ID/path, expected parent and checkpoint payloads, expected claimed state, writer label, direct non-zero gate status, structured argv, fingerprint, repair ID/path/reservation, and optional paired merge commits. Its typed result exposes parent then repair IDs, every affected path/change, create-versus-fold outcome, fingerprint, repair dependency, and merge range. Selector consumers receive request ID/path, provenance, original state, dependency evidence, and `selection_priority`; action-owned mutation never infers identity from display text.

**Plan validation:** All Detailed Requirements map to the six tasks and every task maps back to transaction, schema, selector, or verification requirements. Warning: six tasks and a projected 20+ file surface exceed the five-task heuristic, but splitting would leave the public transaction without its required lossless schema, selector priority, or contract mirrors; REQ-492 already owns the separate orchestration integration.

*Generated by Plan agent*

## Exploration

The existing publication handler map, strict manifest decoder, contained payload loader, `PublicationPlan`, and rooted apply path are the correct command seam; `cmd/do-work-cli/main.go` already registers the whole handler map dynamically. `requestmodel` supplies lossless scalar/body edits but needs a canonical list edit. Repository discovery already carries request records and checkpoint claim evidence, so no `repositorymodel` change is required.

The critical new boundary is a claimed parent whose REQ bytes are intentionally dirty with planning and implementation history. Generic Git transactions currently refuse tracked dirty targets and restore from `HEAD`; `defer-gate` therefore needs a narrow exact-path opt-in that rejects staged targets, snapshots the dirty bytes/mode/identity, and restores that preimage on failure. Root and shipped Just inventories plus the command-line guide enumerate public commands; the generic managed-section recipe scanner needs no change.

Reusable test seams include publication failure injection, Git transaction rollback hooks, request-model lossless fixtures, selector stable ordering, and JSON/text result parity. Fold tests must prove multiple parents share one fingerprint-keyed repair without overwriting prior `related` or evidence records.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go`
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go`
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go`
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go`
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go`
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go`
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/docs/command-line-guide.md`
- `justfile`
- `skills/do-work-board/justfile.template`
- `skills/do-work-board/tools/queue-kanban/model_test.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go`
- `_dev/tests/install-suite-behavior.sh`
- `_dev/tests/flat-just-recipes-behavior.sh`
- `_dev/tests/contract-regressions.sh`

**Acceptance criteria:** One strict manifest transaction creates or uniquely folds a fingerprint-keyed pending repair, returns the exact claimed parent to pending with a repair dependency and durable gate evidence, removes only its writer claim, and rolls every path back to its exact preimage on any failure. The four new frontmatter fields project losslessly; ready repairs sort before ready deferred parents before ordinary work without changing explicit-token or within-class order. Stale, staged, colliding, ambiguous, unsafe, malformed, and invalid merge-evidence inputs refuse before mutation. Public command recipes/docs and all affected schema/selector/result contracts remain synchronized; no blocked or pending-answers lifecycle is introduced.

## Pre-Flight

**Git:** Clean outside `do-work/`.

**Tests:** The uncached CLI preflight recorded a transient `internal/corehelpers` shared-fixture failure while another session was active; the package passed immediately when rerun alone, and the canonical maintainer gate passed immediately before this claim. Preserve the recorded baseline evidence for attribution, but final integration still requires a fully green uncached module and canonical gate.

**Dependencies:** Available.

## Root Cause

The work pipeline had no deterministic owner for turning an unrelated repository-gate failure into dependency work. Publication could create or replace individual records, request state could move one REQ, and selection understood dependencies, but no transaction bound repair allocation/folding, parent re-queue, checkpoint ownership removal, and durable evidence together. The selector also had no stable priority for the repair/resume lifecycle.

## Decisions

- **D-01 — KEEP:** Implement `defer-gate` in the typed publication family so one rooted Git transaction owns every create, fold, move, and checkpoint replacement; orchestration judgment remains caller-owned.
- **D-02 — KEEP:** Add a narrow exact-path opt-in for intentional unstaged tracked preimages. Staged paths and commit mode remain refused, while rollback restores the caller-bound bytes, complete mode, and identity rather than `HEAD`.
- **D-03 — KEEP:** Make priority an explicit result field and sort only selection modes whose order is computed by the selector; directly named REQ order remains caller-owned.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/flat-just-recipes-behavior.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `justfile` (modified)
- `skills/do-work-board/justfile.template` (modified)
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (created)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (created)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** Added the strict `defer-gate --manifest` transaction, exact fingerprint-and-sweep-key repair creation/folding, parent dependency/evidence re-queue, exact owned checkpoint removal, four typed frontmatter projections, canonical list writes, topology-aware rollback, selector priority evidence/order, and synchronized recipes, documentation, installer counts, and contracts. The remediation classifies tracked-dirty, tracked-clean, and untracked targets independently; emits canonical `sweep: true` repairs with `## Instances`; preserves explicit anchors while prioritizing mixed UR expansions; rejects prefix identity collisions and merge commits outside current `HEAD`; and covers create/fold/rollback across the supported repository states.

## Review

**Verdict:** Request changes (50%, acceptance fail; critical risk).

The independent review found that the initial implementation passed its existing suite but did not satisfy valid repository topologies and exact-identity guarantees. The single remediation pass must:

- classify tracked-dirty, tracked-clean, and untracked parent/checkpoint inputs independently, including create, fold, and rollback;
- compare diagnostic fingerprints and checkpoint writer labels exactly rather than by substring;
- emit canonical repair sweeps with `sweep: true` and the `## Instances` projection;
- priority-order UR expansions without disturbing explicitly named REQ order in mixed targets;
- prove the deferred merge belongs to current `HEAD` ancestry;
- synchronize text-result semantics and preserve durable RED/GREEN and rollback evidence.

Focused review verification passed, but the missing topology and collision cases are blocking. Completion and archive remain on hold pending remediation and independent re-review.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; relevant because a prescribed CLI transaction and gate commands change, but the satellite exceeds the budget and is only partially family-slugged.

## Full Context

See `do-work/user-requests/UR-095/input.md` for complete verbatim input.

---
*Source: UR-095 — "Add canonical repository-gate deferral lifecycle"*

## Triage

**Route: C** - Complex

**Reasoning:** This request adds a new multi-file transactional command, lossless schema fields, repair folding, dependency mutations, selector priority, rollback guarantees, and alternate-reader contract updates. Planning and exploration are required before implementation.

**Planning:** Required

## Remediation

The first review rejected the implementation because its passing fixtures covered only tracked-dirty parent/checkpoint inputs and used incomplete identity, sweep, ordering, and merge-evidence contracts. The one permitted remediation pass added independent tracked-dirty/tracked-clean/untracked parent and checkpoint classification, exact parsed fingerprint and writer identity, canonical sweep instances, mixed-target priority/deduplication, current-`HEAD` merge ancestry, and broader result parity and rollback tests.

The post-remediation review confirmed those original blockers closed, but found two remaining facets of one preflight-topology root cause: an existing tracked-dirty repair is not classified for fold, and an occupied parent queue destination is discovered during apply rather than refused during planning. These Important findings are preserved in one review-generated sweep rather than hidden or repaired through a second remediation cycle.

## Re-Review

**Verdict:** Approve with follow-ups (79%, acceptance partial; low integrity risk, material lifecycle-completeness gap).

- **Important — impact-user-visible:** same-fingerprint deferral cannot fold into a tracked-dirty repair because repair-path dirty classification is absent.
- **Important — impact-rule-change:** a pre-existing parent queue destination is rejected only when the move applies, rather than during the no-mutation planning boundary.
- **Minor:** the earlier P-A-U line still named 29 paths; corrected to the final 31-file scope.

Exact fingerprint/writer matching, canonical sweep projection, board consumption, mixed explicit-REQ/UR selection, merge ancestry, strict decoding, text/JSON parity, rollback integrity for the covered topologies, and the no-`blocked`/no-`pending-answers` lifecycle all passed re-review. The two Important findings share the root cause `defer-gate-preflight-topology-incomplete` and are consolidated into REQ-493 (Complete repository-gate deferral preflight topology), one mandatory follow-up sweep.

## Qualification

**Result:** Pass. The final diff matches the declared 31-file Scope and write set; `git diff --check`, `gofmt`, focused package tests, and independent scope review found no unrelated implementation paths. The pre-existing `.DS_Store` files and knowledge-base work remain unstaged and untouched.

## Testing

**Red-green validation:** Before REQ-491, `defer-gate --manifest` and repair/deferred selector priority did not exist. The implemented request-specific tests now pass for strict manifests, create/fold transactions, exact identity collisions, canonical sweep instances, supported parent/checkpoint repository states, rollback positions, mixed REQ/UR ordering, schema projection, result parity, and merge ancestry. The re-review’s two uncovered topology facets are named RED cases in the follow-up sweep.

**Focused:** `go test -race -count=1` passed for publication, Git transaction, request model, schema normalization, next selection, and result model.

**Modules:** uncached `go test -count=1 ./...` passed in the CLI module; the queue-kanban module passed in 148 seconds, including the alternate sweep consumer.

**Canonical:** `bash _dev/tests/maintainer-verify.sh` passed directly, including contracts, shell behavior, installer behavior, queue-kanban vet/tests/strict JavaScript, CLI vet, and uncached CLI tests. The strict browser lane was skipped because no browser was configured, as allowed by the gate.

## Lessons Learned

**What worked:** Exact preimage-bound transactions plus race, alternate-consumer, and repository-topology fixtures exposed failures that happy-path publication tests could not. A fresh re-review caught two missing topology states after the initial remediation suite was green.

**What didn't:** Treating “dirty target support” as a parent/checkpoint concern was too narrow; every existing mutation target, including a folded repair and a move destination, needs classification before planning is considered complete.

**Worth knowing:** `defer-gate` is a cross-reader contract: publication, Git rollback, request/schema projections, selector ordering, recipes, and board sweep parsing must move together. REQ-493 carries the remaining topology-class closure.

## Orientation

[MAP CHANGED] Repository-gate failure deferral now lives in the do-work CLI publication subsystem as a typed atomic transaction, with repair/deferred priority in next-selection and sweep visibility in the board model. The CLI prime and command guide point to the new flow; the remaining dirty-repair/destination-preflight edge class is explicitly queued rather than implied complete.

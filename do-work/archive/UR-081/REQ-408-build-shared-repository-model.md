---
id: REQ-408
title: 'Build shared request, schema, dependency, atomic-file, and repository packages'
status: completed
claimed_at: 2026-08-30T18:12:30Z
completed_at: 2026-08-30T19:22:59Z
commit: ac2e3acd
created_at: 2026-08-29T20:28:26Z
route: C
write_set: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go, skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go, skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_unix.go, skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_windows.go, skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_unsupported.go, skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go, skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go, skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go, skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go, skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go, skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go, skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go, skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go]
estimate:
  p50_active_minutes: 100
  confidence: low
  calculated_at: 2026-08-30T18:12:57Z
  basis:
    - Route C
    - 18-file write set
    - 12 new files
    - 5 subsystems involved
    - 4 acceptance criteria
    - dependency depth 2
    - persistence changes
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
kb_status: promoted
kb_entry: REQ-408-build-shared-request-schema-dependency-a.md
tdd: true
suggested_spec:
depends_on: [REQ-407]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Build Shared Request, Schema, Dependency, Atomic-File, and Repository Packages

## What
Create the reusable repository model required by every request and queue command.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read REQ-408, UR-081, the required crew rules, and existing queue-kanban compatibility sources; ordered work as prime → atomic file → schema → request model → repository model → dependency graph.
- [x] **[APPLY]:** Implemented each package after its behavioral RED inside the declared 14-file scope; no command registration, module dependency, board, later-command, or pipeline-metadata edits.
- [x] **[UNIFY]:** Reviewed all 14 files, confirmed package direction and standard-library-only imports, scanned for debug/TODO/trailing-whitespace artifacts, ran focused/full/race tests, vet, gofmt, and Windows compile-only verification.

## Detailed Requirements
- Implement shared REQ/UR frontmatter parsing and writing, schema aliases and normalization, canonical timestamps, ID allocation, dependency graphs, atomic files, and repository discovery/modeling.
- Preserve unknown fields and bytes where commands are not authorized to rewrite them.
- Support queue, working, root archive, nested archived UR, reservation, and user-request layouts.
- Give downstream commands typed evidence and exact paths without rescanning through ad hoc shell pipelines.

## Constraints
- Use the Go standard library unless an existing dependency is demonstrably necessary.
- Treat this package layer as the single source of truth for later command families.

## Dependencies
Depends on REQ-407 (the Go installation/runtime path must exist before expanding the model).

## Builder Guidance
Certainty level: Firm. Reuse existing queue-kanban schema behavior where compatible, but do not merge the board binary into this module.

## Red-Green Proof
**RED prompt/case:** Parse representative current and legacy REQ/UR fixtures, malformed frontmatter, timestamp variants, reservation races, and dependency cycles with the absent shared packages.
**Why RED now:** Each shell/action path currently reconstructs parts of the repository model independently.
**GREEN when:** Unit and fixture tests expose one normalized typed model, preserve required bytes/fields, allocate collision-free IDs, and return deterministic dependency results.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

---

## Triage

**Route: C** - Complex

**Reasoning:** This REQ establishes the typed repository model used by every later request and queue command, spanning parsing, preservation, allocation, dependency graphs, atomic writes, and repository discovery.

**Planning:** Required

## Plan

Create five standard-library internal packages in the existing `do-work-cli` module, with no CLI command registration and no `go.mod` dependency changes.

1. **Atomic-file contract:** Add RED tests and implement same-directory replacement of existing regular files plus exclusive reservation-marker creation. Preserve file modes, clean temporary files, refuse symlinks/special files and changed targets, and split replacement primitives by platform.
2. **Schema normalization:** Add table-driven RED tests and implement the Schema Read Contract for canonical values, aliases, defaults, exact warning evidence, terminal predicates, dependency satisfaction, and canonical-key precedence.
3. **Lossless REQ/UR model:** Add fixtures for current, legacy, malformed, BOM, CRLF, nested estimate, list, quoting, duplicate-key, and timestamp shapes. Implement byte-offset parsing and field-local edits that retain comments, ordering, line endings, unknown fields, body bytes, and unrelated frontmatter bytes.
4. **Repository discovery, paths, and allocation:** Add layouts for queue, working, loose/nested archive, active URs, reservations, collisions, stray files, and excluded trees. Load one typed snapshot with exact paths/evidence and reserve collision-free REQ ids using exclusive creation and retry.
5. **Deterministic dependency graph:** Add RED tests for roots, chains, branches, aliases, missing targets, terminal statuses, duplicate edges, cycles, reverse edges, readiness, and depth. Derive stable evidence and ordering exclusively from the repository snapshot.
6. **Unify:** Verify acyclic package imports, standard-library-only implementation, no rescans or board invocation, no command registration, and no generated artifacts.

**Planned files:**
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_nonwindows.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_windows.go` (new)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (new)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (new)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (new)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` (new)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go` (new)

**Architectural decisions:** `queue-kanban` stays a separate module and binary; compatibility behavior is characterized rather than imported from its `package main`. Frontmatter support is a lossless parser for documented live REQ/UR YAML shapes, not a general YAML implementation. Reads never rewrite aliases; only authorized writes emit canonical fields. Repository discovery scans once and preserves collisions as evidence rather than choosing a winner. Atomic replacement refuses symlinks and target-identity changes.

**Requirement coverage:** frontmatter parsing/writing and byte preservation map to tasks 1 and 3; normalization to task 2; timestamps to task 3; allocation and reservation races to tasks 1 and 4; dependency results to task 5; repository layouts and exact-path evidence to task 4; standard-library and single-source constraints apply throughout.

**Testing:** RED-before-GREEN per package, then focused package tests, `go vet ./...`, uncached full module tests, a Windows compile-only atomicfile check in a temporary directory, and the unpiped canonical `bash _dev/tests/maintainer-verify.sh` gate.

**Plan validation:** Every Detailed Requirement maps to at least one task; every task traces to the REQ. Warning: the plan has 6 tasks — quality degrades past 3 and the action specifically flags 5+. Keep each unit independently testable and avoid broadening into REQ-409+ commands.

*Generated by Plan agent*

## Exploration

The existing behavior to characterize lives in the separate `queue-kanban` module: `frontmatter.go` and `testing.go` for BOM/CRLF-aware frontmatter parsing and local field edits; `model.go` for schema normalization, dependency aliases, terminal predicates, commit aliases, timestamp parsing, and coercion; `walk.go` for repository discovery and excluded trees; `allocate.go` for filename-plus-frontmatter allocation and cross-process reservations; `atomic_write.go` plus platform replacement files for safe publication; and `dependency_test.go` for current readiness semantics.

The new CLI module has no `go.sum` and must stay standard-library-only. Package layering is `requestmodel` over `schemanormalization`, `repositorymodel` over `requestmodel` and `atomicfile`, and `dependencygraph` over typed repository records. The implementation must retain original bytes and offsets rather than split/join whole files, keep collisions as evidence, observe `.req-reservations` despite hidden-directory pruning, and refuse unsupported atomic replacement platforms. The planned atomic split was corrected from `nonwindows` to the existing proven `unix`, `windows`, and fail-closed `unsupported` shape.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (new) — compact module/package map and verification index
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go` (new) — safe replacement and exclusive reservation API
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go` (new) — atomicity, preservation, refusal, race, and reservation proofs
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_unix.go` (new) — Unix replacement primitive
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_windows.go` (new) — Windows replacement primitive
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_unsupported.go` (new) — fail-closed unsupported-platform primitive
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (new) — aliases, defaults, predicates, and key precedence
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (new) — Schema Read Contract characterization
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (new) — lossless REQ/UR documents, typed fields, edits, and timestamps
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (new) — live/legacy/malformed byte fixtures and mutation proofs
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (new) — repository discovery, typed snapshot, layouts, collisions, and allocation
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (new) — tree-layout, evidence, and reservation-race proofs
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` (new) — deterministic dependency readiness, cycles, depth, and reverse edges
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go` (new) — roots, chains, aliases, missing targets, terminal states, and cycles

**Files I will NOT touch:** `queue-kanban`, `go.mod`, `cmd/do-work-cli/main.go`, shell shims, later command packages, and REQ-409+ lifecycle behavior.

**Acceptance criteria (restated from REQ):**
- [ ] Current and legacy REQ/UR frontmatter is parsed into typed evidence and authorized fields can be written without changing unknown fields, unrelated bytes, comments, ordering, line endings, or body bytes.
- [ ] Schema aliases, normalization, canonical timestamps, terminal predicates, and canonical-key precedence match the documented contracts.
- [ ] Repository discovery covers queue, working, root and nested archive, active user requests, and durable reservation layouts in one typed snapshot with exact paths and collision evidence.
- [ ] REQ allocation is collision-free across filenames, frontmatter ids, archived/live layouts, reservations, and concurrent races.
- [ ] Dependency graphs return deterministic readiness, unmet/dangling evidence, reverse edges, cycles, and depth without rescanning.
- [ ] Atomic writes refuse unsafe targets, preserve required metadata, and use platform-specific replacement without adding dependencies.

## Decisions

### D-01: Create the missing do-work-cli prime index

**Decision:** DECIDE & STATE — add `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` within this REQ and keep it under the prime-file size and noise limits.

**Reasoning:** REQ-408 establishes the reusable package layer for the CLI, and the implementation rules require a primary utility without a prime to gain one before future command families build on it. The prime will index package roles, source-of-truth files, and test commands without duplicating implementation details.

### D-02: Bound frontmatter support to documented live shapes

**Decision:** DECIDE & STATE — implement a byte-preserving parser for the REQ/UR shapes observed and documented in this repository, not a general YAML engine.

**Reasoning:** General YAML support would require a new dependency or a speculative parser surface. Known fields fail closed when unsupported syntax prevents a typed read; unrelated bytes remain preserved.

### D-03: Preserve malformed and duplicate-key evidence

**Decision:** DECIDE & STATE — retain original bytes and expose the last duplicate as the typed read, matching existing recovery behavior.

**Reasoning:** Rewriting malformed input during a read would destroy forensic evidence, while dropping the record would hide live queue state.

### D-04: Count filename and frontmatter IDs

**Decision:** DECIDE & STATE — collision detection and allocation treat both identities as claims on a number.

**Reasoning:** A renamed file with stale frontmatter still owns both numbers; reissuing either would manufacture an avoidable duplicate.

### D-05: Fail dependency readiness closed on ambiguous records

**Decision:** DECIDE & STATE — keep every colliding record in the repository snapshot and refuse to treat an ambiguous dependency target as satisfied.

**Reasoning:** Choosing a winner would hide the collision evidence later doctor/forensics commands need and could run work against the wrong record.

### D-06: Use durable fixed-width reservation marker names

**Decision:** DECIDE & STATE — reservation files use six-digit names, while returned request IDs retain the repository's minimum three-digit display width.

**Reasoning:** Fixed-width markers sort naturally and remain width-agnostic during numeric parsing; returned IDs preserve existing user-facing conventions.

### D-07: Limit mode preservation to ordinary permission bits

**Decision:** DECIDE & STATE — REQ-408 preserves ordinary permission bits; the special-mode-bit correction remains isolated in REQ-426.

**Reasoning:** REQ-426 is the already-captured addendum for setuid, setgid, and sticky-bit behavior. Folding it here would duplicate queued work and blur review provenance.

### D-08: Use content digests as authoritative change evidence

**Decision:** DECIDE & STATE — compare SHA-256 content evidence before atomic publication; inode, size, and mtime remain supplementary checks.

**Reasoning:** An in-place writer can retain inode identity, file size, and timestamp. Digest comparison is the smallest reliable way to refuse overwriting those edits.

### D-09: Anchor reservation publication at both containment levels

**Decision:** DECIDE & STATE — create markers through rooted `do-work` and `.req-reservations` handles, then verify both identities after publication and remove the marker through the original handle on mismatch.

**Reasoning:** Path-based validation followed by path-based creation leaves a symlink-swap window. Rooted handles keep the write inside the directory that was actually validated.

### D-10: Make collision evidence authoritative for dependency ambiguity

**Decision:** DECIDE & STATE — any ID named by `CollisionEntries`, including a cross filename/frontmatter claim, is ambiguous for node lookup, readiness, unmet-dependency evidence, and depth.

**Reasoning:** Indexing only the selected typed record can hide a second path that claims the same identity. Readiness must fail closed wherever the repository snapshot has collision evidence.

### D-11: Read discovered documents through rooted regular-file handles

**Decision:** DECIDE & STATE — REQ and UR bytes must come from an `os.Root`-contained regular-file handle whose identity is checked before and after the read; symlinks and non-regular files are warning evidence, never parsed queue state.

**Reasoning:** A path that is lexically inside `do-work/` can still resolve outside it through a symlink. Rooted reads keep discovery evidence contained in the validated repository tree.

### D-12: Project every downstream schema field with normalization evidence

**Decision:** DECIDE & STATE — the typed request record exposes normalized values and `FieldResult` evidence for domain, route, impact, effort estimate, maintenance, TDD, error type, KB status, testing status, and builder-decided, while retaining generic parser evidence.

**Reasoning:** Later commands need one shared model, not another round of field-by-field parsing and normalization.

### D-13: Preserve unknown values when no default exists

**Decision:** DECIDE & STATE — an unrecognized field with no schema default retains its normalized raw value in `ResolvedValue` while remaining explicitly unrecognized.

**Reasoning:** Returning an empty value contradicted the warning that the input was reported unchanged and discarded evidence downstream commands may need to diagnose.

### D-14: State the portable atomic-publication boundary precisely

**Decision:** DECIDE & STATE — replacement guarantees complete atomic publication and best-effort refusal for target changes observed during pre-publication validation; it does not claim portable compare-and-swap against arbitrary non-cooperating writers.

**Reasoning:** Standard-library Unix and Windows replacement primitives accept no expected inode, version, or digest. Every separate validation has a final check/use interval, and cooperative locks cannot constrain arbitrary writers.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_unix.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_windows.go` (new)
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_replace_unsupported.go` (new)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (new)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (new)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (new)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` (new)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go` (new)

**What was done:** Added the standard-library shared repository-model layer for safe atomic publication and reservations, schema normalization, byte-preserving REQ/UR parsing and field edits, one-pass repository discovery/allocation, and deterministic dependency evidence. Added a compact prime index for later command families.

## Qualification

**Attempt 1: Failed.** Mechanical qualification and scope drift passed, all 14 files are substantive, and the six REQ acceptance areas trace to implementation. Two race-safety requirements remain unproven:

1. `atomicfile.ReplaceExisting` rechecks only file identity. An in-place edit between validation and replacement keeps the same inode, so the current code overwrites another writer despite the plan's changed-target refusal.
2. `repositorymodel.ReserveNextRequestID` validates `.req-reservations` with `Lstat` and later creates by absolute path. A directory symlink swap between those operations can redirect the exclusive marker outside the repository; the explored board implementation uses a rooted directory handle specifically to close this window.

Return to implementation with assertion-level regression tests for both races, then rerun qualification.

**Attempt 2: Passed.** All 14 files exist, are substantive, and match Scope exactly. Six acceptance areas trace to concrete packages and behavioral tests. P-A-U is complete; no debug artifacts or undeclared dependencies are present. The static-reference warnings are expected for new library packages and their tests before REQ-409+ register consumers. Independent focused, race, vet, full-module, and Windows compile checks all exited 0 after the two race regressions turned GREEN.

**Attempt 3 after review remediation: Passed.** Collision evidence now flows into graph ambiguity/readiness/depth, discovered REQ/UR reads are rooted and reject symlinks, every downstream schema field has normalized typed evidence, and unknown no-default values remain available. The atomic contract and prime now state the defensible portable boundary. Mechanical qualification and scope drift exited 0; all 14 files remain substantive and exactly scoped.

## Testing

**Tests run:** `go test -count=1 ./internal/atomicfile ./internal/schemanormalization ./internal/requestmodel ./internal/repositorymodel ./internal/dependencygraph`; `go test -race -count=1 ./internal/atomicfile ./internal/repositorymodel`; `go vet ./...`; `go test -count=1 ./...`; Windows atomicfile compile-only; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing. The canonical repository gate exited 0; its optional browser lane remained in the documented default-skipped state.

**Red-green validation:**
- `atomic_file_test.go`: whole replacement first returned old contents; in-place same-inode mutation was initially overwritten → ✓ complete publication and SHA-256 change refusal pass
- `schema_normalization_test.go`: aliases/defaults and no-default warning assertions failed → ✓ full Schema Read Contract cases pass
- `request_model_test.go`: parsed body was empty and literal-block/comment preservation failed → ✓ byte-preserving parse/edit cases pass
- `repository_model_test.go`: discovery returned no records and filename/frontmatter collision evidence was absent → ✓ layouts, collision evidence, allocation races, and rooted-directory swap refusal pass
- `dependency_graph_test.go`: dependency records had empty status/edges and were not ready → ✓ readiness, terminal gating, reverse edges, cycles, and depth pass
- Review remediation: cross filename/frontmatter collision initially left a dependent ready; symlinked request content entered the snapshot; downstream normalized values were empty; and unknown no-default resolution was empty → ✓ ambiguity fails closed, contained reads reject symlinks, complete typed evidence is populated, and unknown values are retained

**New tests added:**
- `internal/atomicfile/atomic_file_test.go`
- `internal/schemanormalization/schema_normalization_test.go`
- `internal/requestmodel/request_model_test.go`
- `internal/repositorymodel/repository_model_test.go`
- `internal/dependencygraph/dependency_graph_test.go`

*Verified by work action*

## Review

**Overall: 82%** | 2026-08-30T19:22:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 88% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings:**
- Filename-only collision claims with divergent frontmatter IDs are safely blocked but labeled missing rather than ambiguous because nil-node handling precedes collision evidence — impact-user-visible → REQ-428
- The normalized typed record omits `caveman`, leaving one Schema Read Contract field outside the promised shared projection — impact-rule-change → REQ-429

**Prior findings reassessed:** The atomic contract is now precise; cross filename/frontmatter collision readiness, rooted symlink containment, the previously listed typed fields, and unknown no-default preservation are fixed.

**Minor findings:** 1 (folded YAML block scalars are retained but not decoded as folded content; no current repository fixture uses the shape)
**Acceptance:** Partial — all core operations and safety gates pass; the two remaining evidence-completeness gaps are isolated in mandatory follow-ups.
**Suggested testing:** filename-only collision classification, `caveman` projection completeness, and eventual Windows runtime replacement behavior.
**Follow-ups created:** REQ-428, REQ-429; **sweeps appended to:** None

*Reviewed by review-work action after remediation*

## Lessons Learned

**What worked:**
- Byte-offset frontmatter edits and rooted file/directory handles kept preservation and containment contracts testable without external dependencies.
- Assertion-level adversarial fixtures exposed collision handoff and directory-swap defects that normal package tests missed.

**What didn't:**
- Inode, size, and timestamp checks alone did not detect restored-metadata in-place edits; content evidence was required.
- The first review phrased portable atomic replacement as compare-and-swap against arbitrary writers, a guarantee the standard-library replacement primitives cannot provide; narrowing the contract made the real safety boundary reviewable.

**Worth knowing:** Collision evidence is a first-class repository fact and should be consulted before any dependency winner is selected. Unknown schema values remain evidence even when they are unrecognized.

## Orientation

The do-work CLI now has a shared, standard-library repository-model layer for lossless REQ/UR documents, schema evidence, contained discovery and allocation, deterministic dependency graphs, and atomic publication. It lives under the `do-work-cli` subsystem described by `prime-do-work-cli.md`. **[MAP CHANGED]** — later command families can consume one typed snapshot instead of rebuilding queue state independently. The prime's file, directory, command, and relative-reference pointers all resolve in the current tree.

---
id: REQ-485
title: 'Canonicalize REQ reservation marker filenames across allocation flows'
status: claimed
priority: now
created_at: 2026-09-01T12:11:03Z
user_request: UR-092
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
batch: go-no-llm-command-platform
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-09-03T21:44:31Z
  basis:
    - Route C
    - 8-file write set
    - 4 subsystems involved
    - 5 acceptance criteria
    - persistence changes
    - cross-route regression gates
    - full-suite verification
claimed_at: 2026-09-03T21:43:41Z
route: C
planning_at: 2026-09-03T21:53:38Z
exploration_at: 2026-09-03T21:53:38Z
write_set:
  - skills/do-work-board/tools/queue-kanban/allocate.go
  - skills/do-work-board/tools/queue-kanban/allocate_test.go
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/capture_files.go
  - skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/reservations.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go
  - skills/do-work/tools/do-work-cli/internal/hookcommands/session_start_test.go
  - skills/do-work/actions/capture.md
---

# Canonicalize REQ Reservation Marker Filenames Across Allocation Flows

## What

Make every REQ-number reservation flow use one canonical marker filename so
exclusive-create collision detection actually collides. Today queue-kanban
`next-req` writes zero-padded `REQ-000482` while capture-files manifests carry
unpadded `REQ-482`; the two flows compute the same candidate from the same max-scan
and both succeed, defeating the guard entirely.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in
any UR sharing this root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Pick one canonical marker filename and enforce it in every writer: queue-kanban
  `allocate.go` (`next-req`), the capture-files marker-path guidance in
  `skills/do-work/actions/capture.md`, and any other flow that creates markers.
- Read-side acceptance of both legacy spellings everywhere markers are consumed:
  the allocation max-scan, `skills/do-work/scripts/cleanup-req-reservations.sh`'s
  committed-REQ and timeout reaping, and any capture-side duplicate scan — an
  existing marker in either spelling must block its number and must reap.
- Writers refuse or normalize a non-canonical marker path in a capture manifest
  rather than silently creating a second spelling.
- Lock-in test reproducing the 2026-09-01 collision shape: two flows reserving the
  same number through different spellings must collide, not both succeed.
- The board parser lock-step rule (`_dev/primes/prime-kanban-board.md`) governs if
  the board reads marker names anywhere; verify and follow it.

## Red-Green Proof

**RED prompt/case:** Reserve a number via `queue-kanban next-req`, then submit a
capture-files manifest whose reservation path is the unpadded spelling of the same
number.
**Why RED now:** Both succeed — reproduced 2026-09-01: `next-req` created
`REQ-000482` at 11:50Z and a concurrent capture created `REQ-482` and committed a
second REQ-482 file (78847fe4) at 11:56Z; only manual inspection prevented a
duplicate REQ id from being committed twice.
**GREEN when:** The second reservation is refused with an actionable finding, the
lock-in test pins the cross-spelling collision, and legacy-spelling markers still
count in the max-scan and still reap.
**Validation:** Observed collision, this session; evidence in the UR input.

## Full Context

See `do-work/user-requests/UR-092/input.md` for complete verbatim input.

---
*Source: UR-092 (Canonicalize REQ reservation marker filenames across allocation flows)*

## Triage

**Route: C** - Complex

**Reasoning:** The canonical spelling must be selected once and applied atomically across board allocation, capture validation, cleanup compatibility, docs, and collision tests while preserving both legacy read spellings. The cross-writer contract needs a plan and source inventory.

**Planning:** Required

## Plan

1. Add literal cross-writer RED fixtures proving the board allocator's six-digit marker and capture's stored-ID marker can reserve the same number independently.
2. Make every current writer derive the canonical stored-ID basename (`REQ-%03d`, minimum width), and require canonical capture/defer manifest paths while retaining exclusive creation as the ownership event.
3. Add exact, width-agnostic numeric reservation readers so allocation maxima, publication collision/fold handling, and cleanup recognize both canonical and legacy aliases without accepting suffix junk.
4. Prove canonical collision, legacy blocking/folding, coexistence cleanup, malformed preservation, and rooted/race safeguards independently in both Go modules.
5. Update capture and shipped board guidance, then run both module suites plus repository contract/release gates.

**Plan validation:** The plan covers every writer, reader, cleanup, manifest, migration, and documentation acceptance criterion. The separate Go modules cannot share an internal helper, so literal cross-contract tests are required to prevent helper-coupled false greens.

*Generated from delegated exploration; full evidence: `do-work/runs/work-2026-09-03-214500/REQ-485-exploration.md`.*

## Exploration

The two allocators independently emit fixed-six markers, while capture and defer-gate derive paths from the stored ID; cleanup then recognizes only six digits. Canonical writers must converge on the stored-ID spelling, while exact numeric readers remain width-compatible during migration. Publication must check numeric aliases before planning a create, defer fold must accept a matching legacy marker, and cleanup must revalidate and reap each concrete alias independently.

## Scope

**Files I will touch:** the allocator, repository-model, publication, cleanup, and delegated answer/hook tests listed in `write_set`, plus `skills/do-work/actions/capture.md` and the shipped board prime.

**Files I will not touch:** board card/parser/UI code, generated board output, request-file numbering, marker contents, migration commands, or the cleanup shell launcher.

**Acceptance criteria:** New writers use one minimum-three-digit basename; either legacy or canonical evidence blocks the number and participates in max scans/folds; both aliases reap under the existing authority/age policy; malformed adjacent names remain excluded; current writers collide at the canonical exclusive-create boundary.

## Pre-Flight

**Git:** The wave baseline was clean at `c27d349a` after the three claims, estimates, run manifest, and briefs were committed.

**Tests:** Direct canonical fast gate passed and was recorded at the shared wave baseline before any builder dispatch.

**Dependencies:** None. The board and CLI are separate modules and must each be verified independently.

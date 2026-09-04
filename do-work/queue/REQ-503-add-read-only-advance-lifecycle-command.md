---
id: REQ-503
title: 'Add the read-only advance lifecycle command'
status: pending
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-489, REQ-498, REQ-499, REQ-500, REQ-501, REQ-502, REQ-567]
batch: orchestrator-simplification
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
route: C
planning_at: 2026-09-04T15:02:48Z
exploration_at: 2026-09-04T15:04:18Z
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-04T14:57:40Z
  basis:
    - Route C
    - 4-file write set
    - 2 new files
    - 2 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
gate_deferred: 'true'
---

# Add the Read-Only advance Lifecycle Command

## What
Add `do-work-cli advance REQ-NNN`: for the REQ's route, report the current lifecycle phase, the evidence still missing, and the exact next command, and refuse an impossible transition with a typed finding. Read-only in this REQ: it composes the existing commands (`next`, `claim`, `estimate-p50`, `preflight`, `scope-drift`, `qualify`, `run-blocked-check`, `complete`, `finalize`, `recover-finalization`) into one state machine but mutates nothing itself.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Every later move in this chain deletes prose by pointing at `advance`. Landing it read-only first means it changes no behavior, so it needs no prose deletion and every other REQ can depend on it.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- Phases per route (A, B, C) derived from the current `work.md` step order; the report's step table is the source, with mechanical phases CLI-driven and judgment phases reported as "agent judgment: <what>".
- Output is typed: phase, missing evidence with the file or field that would satisfy it, `next_argv`, `verification_argv`; text and JSON formats like every other command.
- Refusals are typed findings (for example `ADVANCE-EVIDENCE-MISSING`, `ADVANCE-PHASE-UNKNOWN`); no paragraphs.
- No mutation in this REQ; `advance` never writes a REQ, the checkpoint, or Git.
- Prime file `prime-do-work-cli.md` documents the command and its phase table.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Root of the chain: waits on REQ-489, REQ-498, REQ-499, REQ-500, REQ-501 and REQ-502 so recovery and finalization exist before the state machine composes them.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Run `do-work-cli --format json advance REQ-NNN` against a fixture REQ in each of: queue, claimed without triage, planned Route C, implemented without qualification, archived without provenance.
**Why RED now:** The command does not exist.
**GREEN when:** Each fixture returns the expected phase, the missing evidence, and a `next_argv` that names the existing command for that phase; an unknown state returns a typed refusal; `go test ./internal/lifecycleadvance` passes and no file outside the write set changes.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6265 tokens, over the claim-time budget; `slugged: partial` so no targeted form. Matched on typed evidence projection and cross-action lifecycle composition.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** This adds a typed state-machine command that composes lifecycle evidence across queue, working, and archive states, with route-specific phases and exact next-command contracts. The new package, command registration, cross-route fixtures, and later action consumers make planning and source exploration necessary before implementation.

**Planning:** Required.

## Plan

1. Add `internal/lifecycleadvance/advance_commands_test.go` first with queue, claimed Route A/B/C, implemented, qualified, reviewed, archived-provenance-gap, complete, malformed, ambiguous, and impossible-order fixtures. Assert exact phases, evidence coordinates, argv, typed refusals, and byte-for-byte read-only behavior.
2. Add `internal/lifecycleadvance/advance_commands.go` as a thin read-only composition layer: discover one repository snapshot, resolve one identity-consistent REQ, classify ordered phase evidence from normalized frontmatter and exact Markdown sections, and emit existing commands as argv without invoking mutating handlers.
3. Extend `internal/resultmodel/result_model.go` and its tests with a first-class advance projection carrying request identity/path, tree/status/route provenance, stable phase and phase kind, structured missing-evidence coordinates, tokenized next argv, and verification argv. Render the same ordered fields in text and JSON; do not hide state in prose findings.
4. Register the handler once in `cmd/do-work-cli/main.go`, then update `prime-do-work-cli.md` with the package owner, read-only boundary, route phase table, evidence sources, and typed refusal codes.
5. Run focused lifecycle/result-model tests, command registration tests, `go vet ./...`, the uncached module suite, and the repository’s unpiped canonical gate.

**Architectural decisions:** Use one repository snapshot and no sibling-command subprocesses. Treat intermediate `##` sections as evidence because the current action defines them that way. Judgment phases carry empty `next_argv` and name the evidence an agent must record; mechanical phases carry complete tokenized argv. A clean pre-flight has no durable section, so Scope-without-Implementation-Summary is a combined preflight/implementation boundary rather than a fabricated `preflight_at`. Preparing a finalization manifest remains agent judgment because its exact path is not discoverable from a REQ.

**Consumer field contract:** identity (`request_id`, exact `request_path`), provenance (`tree_section`), observed state (`status`, `route`), derived state (`phase`, `phase_kind`), structured requirements (`missing_evidence[]` with kind/path/field-or-section/expected), action (`next_argv`), replay (`verification_argv`), and ordinary top-level outcome/findings.

**Plan validation:** All five detailed requirements map to tasks above and every task traces to one of them. The captured scope is contradictory: `CommandResult` has no typed advance field, so `internal/resultmodel/result_model.go` and `internal/resultmodel/result_model_test.go` are required additions. The final scope must therefore contain six files, not the captured directory-plus-two-file shorthand. The request mentions `run-blocked-check`, but the current action makes targeted `next` the sole blocked-probe owner; follow that current source of truth and pin it in tests.

*Generated by Plan agent.*

## Exploration

`repositorymodel.DiscoverRepository` already provides one complete snapshot, and each `RequestFile` carries the exact relative path, tree section, filename ID, parsed document, normalized record, source bytes, and parse failure. `requeststate.ResolveTarget(snapshot, id, "")` can reuse the existing unique-identity authority without a package cycle. Estimate evidence is available through `FieldEvidenceByName["estimate"].NestedValues`; intermediate state must be detected with a narrow exact-heading scanner over `RequestDocument.BodyBytes()` because the existing section helpers are private.

The command should follow the existing package pattern: constant, `Handlers()`, strict parser, direct handler, and typed finding constructors. Tests should use real temporary `do-work/{queue,working,archive}` fixtures, invoke the handler for the phase matrix, invoke `commandruntime.NewRuntime` for text/JSON parity, compare argv slices exactly, and snapshot files plus Git state to prove both formats are read-only. Route absence is triage judgment; unknown nonblank routes, duplicate authoritative fields, malformed IDs, contradictory section order, and impossible tree/status pairs must refuse rather than fall through.

The six-file scope is sufficient. No action prose, request parser, repository discovery, existing command package, or contract-regression file should change in this foundation REQ.

*Generated by Explore agent.*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (new) — discover and classify the read-only lifecycle state.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (new) — phase matrix, refusal, argv, registration, and immutability tests.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — first-class typed advance projection, normalization, and text rendering.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — text/JSON parity and non-null collection tests.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify) — register the command family once.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — document ownership, route phases, evidence, and refusal boundaries.

**Files I will NOT touch:** `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, `_dev/tests/contract-regressions.sh`, existing lifecycle command packages, or request/repository parsing packages. Later chain members own prose and predicate deletion.

**Acceptance criteria (restated from REQ):**
- [ ] Route A, B, and C phases follow the current action order and distinguish mechanical commands from named agent judgments.
- [ ] JSON and text expose typed identity, phase, structured missing evidence, tokenized `next_argv`, and `verification_argv` from one result.
- [ ] Missing, malformed, ambiguous, contradictory, or impossible evidence produces stable typed findings.
- [ ] `advance` never mutates a REQ, the checkpoint, Git, or any other repository byte.
- [ ] The CLI prime documents the command and phase table; this foundation deletes no action prose.

## Repository Gate Deferral

- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** sha256:3af85b84722557f94ddfd466fc32136086fb5fed306e478bd344f689902472ff
- **Repair dependency:** REQ-567
- **Diagnostic evidence:** "shipped-package-reference-contract: obsolete archive link for REQ-491"
- **Diagnostic evidence:** "shipped-package-reference-contract: obsolete archive link for REQ-492"
- **Diagnostic evidence:** "shipped-package-reference-contract: obsolete archive link for REQ-493"

---
id: REQ-503
title: 'Add the read-only advance lifecycle command'
status: completed
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
preflight_at: 2026-09-04T15:22:08Z
dispatch_at: 2026-09-04T15:22:56Z
builder_handback_at: 2026-09-04T15:39:51Z
integration_at: 2026-09-04T15:40:31Z
status_changed_at: 2026-09-04T20:59:59Z
commit: f38a78b0ad80b34b3f5cd332b31e21ae63a7602d
heavy_verified_at: 2026-09-04T20:59:59Z
heavy_verified_revision: f38a78b0ad80b34b3f5cd332b31e21ae63a7602d
review_at: 2026-09-04T22:04:28Z
kb_status: pending
claimed_at: 2026-09-04T21:59:20Z
completed_at: 2026-09-04T22:04:50Z
release_at: 2026-09-04T22:04:50Z
---

# Add the Read-Only advance Lifecycle Command

## What
Add `do-work-cli advance REQ-NNN`: for the REQ's route, report the current lifecycle phase, the evidence still missing, and the exact next command, and refuse an impossible transition with a typed finding. Read-only in this REQ: it composes the existing commands (`next`, `claim`, `estimate-p50`, `preflight`, `scope-drift`, `qualify`, `run-blocked-check`, `complete`, `finalize`, `recover-finalization`) into one state machine but mutates nothing itself.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The Plan and Explore agents defined the six-file command/result-model/registration/prime boundary, route matrix, typed refusals, and test-first implementation order.
- [x] **[APPLY]:** Added the snapshot-driven read-only command, typed projection, registration, phase/refusal tests, and prime documentation in exactly the six declared files.
- [x] **[UNIFY]:** Reviewed the complete six-file diff; focused tests, module tests, vet, the real CLI seam, diff checks, and the canonical repository gate passed on the builder branch.

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

## Pre-Flight

**Git:** Clean outside `do-work/`.

**Tests baseline:** `go test -count=1 ./...` in the CLI module launched successfully and retained the pre-existing `internal/corehelpers` failure in `TestProtectedInventoryPersistsLaterXAndRequiresStartedState`; every other package passed. This is the same unrelated baseline failure observed before the repository-gate repair.

**Repository gate:** `check-green-gate` returned typed success with `matches: true` for `bash _dev/tests/maintainer-verify.sh` at baseline revision `62df4d447c6525c2e34631fbe1fde6503d9a2ef8`, so no redundant pre-build gate run was required.

**Dependencies:** Installed; Go 1.26.1 and the repository gate dependencies were verified by the recorded green run.

*Checked by work action.*

## Implementation Summary

- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (new) — resolves one canonical request snapshot and classifies its next route-aware lifecycle phase without mutation.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (new) — covers Route A/B/C phase progression, typed refusals, exact argv, JSON/text output, and byte-for-byte read-only behavior.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — adds the first-class lifecycle advance projection and deterministic rendering.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — locks text/JSON parity and non-null collections for the new projection.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified) — registers the lifecycle advance command family once.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — documents command ownership, the phase table, evidence boundaries, read-only guarantee, and refusal codes.

## Decisions

- **D-01 — canonical identity:** Reuse one repository snapshot and `requeststate.ResolveTarget` so lifecycle classification cannot drift from existing collision and identity authority.
- **D-02 — local heading evidence:** Keep a narrow exact-heading scanner inside the new package because existing section helpers are private and parser expansion is outside this foundation's scope.
- **D-03 — durable evidence only:** Project preflight at the Scope-to-Implementation-Summary boundary because successful preflight has no durable timestamp or section.
- **D-04 — typed judgment boundary:** Leave `next_argv` empty for judgment phases; only existing deterministic commands receive complete tokenized argv.
- **D-05 — current blocked owner:** Point blocked requests at targeted `next`, the current owner of the bounded probe, rather than the obsolete standalone blocked-check composition.
- **D-06 — Route A exception:** Project Route A testing as agent judgment immediately after qualification because Route A does not run scope drift.

## Discovered Tasks

None.

## Qualification

Passed — 6 files verified, 5 acceptance criteria traced, and P-A-U confirmed. The exact merge range `d1e87a31d358d27479816f9181555887e71af3f6..f38a78b0ad80b34b3f5cd332b31e21ae63a7602d` matches the declared scope. The generic new-file wiring heuristic warned on the two `advance_commands` basenames, but both files are wired by Go package compilation and the `lifecycleadvance.Handlers()` registration imported from `cmd/do-work-cli/main.go`; focused and black-box command tests exercise that seam.

## Testing

- RED: before production implementation, `go test -count=1 ./internal/lifecycleadvance` exited 1 in 2.6 seconds because every fixture reached the real CLI and received typed `UNKNOWN-COMMAND`; a later test-first Route A boundary case also failed with mechanical `scope-drift` instead of testing judgment.
- GREEN: `go test -count=1 ./internal/lifecycleadvance ./internal/resultmodel ./cmd/do-work-cli` passed on merged `main`; package times were 1.183, 0.340, and 0.458 seconds.
- Live seam: `do-work-cli --format json advance REQ-503` returned typed Route C `scope-drift` evidence and exact tokenized next/replay argv without changing repository state.
- Canonical repository gate: `bash _dev/tests/maintainer-verify.sh` passed on merged `main`; 647 CLI tests completed with the slowest test file at 23.46 seconds, below the 30-second budget.
- The pre-existing core-helper baseline failure did not recur in either the builder's full module run or the merged canonical gate; no new regression remains.
- Green-gate evidence was recorded at revision `8201a1dc01295ae0fb6f2154f4871d0cc7b44d94` for the exact canonical argv.

## Heavy Verification Plan

- Base revision: `d1e87a31d358d27479816f9181555887e71af3f6`
- Target revision: `f38a78b0ad80b34b3f5cd332b31e21ae63a7602d`
- `do-work-cli-integrations`: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — all six declared files matched subtree `skills/do-work/tools/do-work-cli`.
- `staged-skills`: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all six declared files matched subtree `skills`.
- `updater`: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — all six declared files matched subtree `skills/do-work/tools/do-work-cli`.
- `installer`: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — all six declared files matched subtree `skills/do-work/tools/do-work-cli`.

## Open Questions

- [x] Run the selected heavy lane commands at `f38a78b0ad80b34b3f5cd332b31e21ae63a7602d`; did every command exit 0? → Confirmed: All 4 selected heavy lanes passed without skips at f38a78b0ad80b34b3f5cd332b31e21ae63a7602d.

Recommended: Yes

Also: No — <failing lane>


## Answer Notes

- 2026-09-04 - [ ] Run the selected heavy lane commands at `f38a78b0ad80b34b3f5cd332b31e21ae63a7602d`; did every command exit 0?: Confirmed: All 4 selected heavy lanes passed without skips at f38a78b0ad80b34b3f5cd332b31e21ae63a7602d.
> ```
> Exact-revision heavy verification via do-work clarify. Stored base, target, selected lanes, argv and coverage reasons matched the recomputed plan. All lane results came from the detached checkout at f38a78b0ad80b34b3f5cd332b31e21ae63a7602d.
> All 4 selected heavy lanes passed without skips at f38a78b0ad80b34b3f5cd332b31e21ae63a7602d.
> Initial attempt: staged-skills, updater and installer each exited 1 after 0 seconds before their tests started, reporting an invalid timing-log header. Preserved the original log and initialized a fresh log using the repository test-duration-log.sh helper. Reran only those three lanes at the same revision; all passed. The earlier passing CLI integration result remains applicable. No tracked source was changed.
> Scope: verification results only; implementation changes, fresh review and archiving remain for do-work run. Date and timestamp follow skills/do-work/actions/work-reference.md, Timestamp rule and its date-only paragraph.
> ```

## Heavy Verification Result

Target revision: `f38a78b0ad80b34b3f5cd332b31e21ae63a7602d`
Execution revision: `f38a78b0ad80b34b3f5cd332b31e21ae63a7602d`

- do-work-cli-integrations: exit 0, 53s — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0, 23s — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0, 52s — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0, 24s — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Review

**Overall: 78.75%** | 2026-09-04T22:01:46Z

**Verdict: Approve with follow-ups.** Typed identity, phase, evidence coordinates, next/replay argv, deterministic rendering, command registration, and the read-only boundary are delivered. Acceptance is partial because legitimate Route A completion and fenced Markdown examples are misclassified.

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

Overall calculation: `(90 + 85 + 80 + 100) / 4 - 10 = 78.75`; the ten-point deduction is the documented Partial modifier.

**Important findings:**

- At the reviewed revision, `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:243` requires `Lessons Learned` for every route and refuses an existing `Orientation` without it. `skills/do-work/actions/work.md:627` expressly permits straightforward Route A requests to skip lessons. A Route A fixture with Triage, Plan, Implementation Summary, Qualification, Testing, Review, and Orientation returned exit 1 / `ADVANCE-EVIDENCE-MISSING` instead of finalization-manifest judgment. Preserve that optional Route A path and cover it with a focused regression. This is an incorrect advisory refusal, with the existing prose still providing the legitimate completion path — impact-user-visible → report only
- At the reviewed revision, `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:334` scans every matching `##` line without excluding fenced code. A Route A request awaiting implementation returned `qualify` after adding only a fenced Markdown example containing `## Implementation Summary` beneath What; the control without that example correctly returned implementation judgment. Treat fenced examples as content rather than completed phase evidence and cover both backtick and tilde fences. This misdirects the advisory next step; downstream qualification remains an independent authority — impact-user-visible → report only

**Minor findings:** None.

**Requirements checklist:**

- Route A/B/C normal phase progression and mechanical/judgment separation: delivered, except the optional Route A lessons branch above.
- Typed phase, exact request identity/path, missing-evidence coordinates, next argv and verification argv in text/JSON: delivered and exercised through the real binary.
- Typed malformed/ambiguous/impossible-state refusals: delivered for the covered cases; fenced content can be mistaken for real lifecycle evidence as recorded above.
- Read-only repository/REQ/checkpoint/Git behavior: delivered; the implementation discovers and projects without calling mutation handlers, and the byte-digest test passed in both output formats.
- Prime ownership and route table: delivered. The six production/test/doc files match the expanded Scope; the additional range entry is the REQ's own lifecycle evidence. No action prose or sentence predicates were deleted, as this explicitly scoped foundation requires; subsequent REQs own those moves.

**Acceptance: Partial.** In an isolated detached worktree at the exact target, `go test -count=1 ./internal/lifecycleadvance ./internal/resultmodel ./cmd/do-work-cli` passed (1.486s, 0.700s, 0.502s). These tests invoke the real CLI across queue, Route A/B/C, archive provenance, malformed identity, contradictory evidence, text/JSON, and byte-for-byte immutability cases. An independently built exact-target binary reproduced both findings above; the ordinary Route A implementation control passed. Four saved exact-target heavy lanes already passed without skips and were not rerun. The detached worktree was clean and removed without force; temporary binary and fixture directories were removed. No background task remains.

**Restatement sweep:** Compared the new phase table and classifier with the exact-target work action and reference templates, and searched shipped consumers for `advance`, its typed result, and refusal tokens. Registration, renderer and prime agree with the new public projection. The optional Route A lessons mismatch is recorded above. The phase classifier's omission of fenced-content boundaries was confirmed with the real binary rather than inferred from tests. Current later revisions still contain both underlying classifier paths; intentional mutation extensions are outside this review's attribution.

**Self-validation:** All three P-A-U boxes are checked. The scope expansion and canonical blocked-probe choice are documented decisions. The stored RED failure exercised the missing command, and GREEN exercises registration rather than merely compiling helpers. Extra edge cases were tested because the existing matrix did not exercise optional lessons or Markdown examples. No source edits, queue capture, or commits were made by this reviewer.

**Suggested testing:** 2 items — Route A optional-lessons completion and lifecycle-looking headings inside fenced examples, as specified in the findings.

**Follow-ups created:** None (2 findings report only).

*Reviewed by review-work action; orchestrated artifact for the owning work action to persist.*

## Lessons Learned

**What worked:** One canonical repository snapshot kept identity and read-only projection aligned; exact-revision review separated this foundation from its later mutation extensions.

**What didn't:** A line-based heading scanner mistook fenced examples for completed evidence, and the phase matrix omitted the optional Route A lessons branch. Both remain report-only review findings under the action's Partial-acceptance policy.

**Worth knowing:** Phase classifiers need the route's exceptions and Markdown content boundaries as well as its normal section order.

## Orientation

[MAP CHANGED] The lifecycle CLI can report a request's next phase, missing evidence, and replayable command from one canonical snapshot. This read-only foundation supports the later queue, evidence-gate, and finalization extensions. Two advisory phase-classification edge cases remain documented in Review. The CLI prime's referenced package paths still resolve.

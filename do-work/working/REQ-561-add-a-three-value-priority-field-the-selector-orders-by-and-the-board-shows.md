---
id: REQ-561
title: 'Add a three-value priority field the selector orders by and the board shows'
status: claimed
created_at: 2026-09-03T20:38:55Z
user_request: UR-107
domain: backend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
route: C
estimate:
  p50_active_minutes: 65
  confidence: low
  calculated_at: 2026-09-03T21:09:41Z
  basis:
    - Route C
    - 11-file write set
    - 4 subsystems involved
    - 4 acceptance criteria
    - browser evidence
    - persistence changes
    - cross-route regression gates
    - full-suite verification
related: [REQ-530]
required_lessons: [_dev/primes/lessons-action-files.md#alternate-writer-contract-drift, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#alternate-writer-contract-drift]
write_set:
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/actions/capture.md
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go
  - skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go
  - skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
  - skills/do-work-board/actions/board.md
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/serve_test.go
  - skills/do-work-board/tools/queue-kanban/priority_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - CHANGELOG.md
  - VERSION
  - skills/do-work/CHANGELOG.md
  - skills/do-work/actions/version.md
  - do-work/queue/
claimed_at: 2026-09-03T20:58:09Z
planning_at: 2026-09-03T21:10:13Z
exploration_at: 2026-09-03T21:17:07Z
dispatch_at: 2026-09-03T21:18:56Z
builder_handback_at: 2026-09-03T21:35:57Z
integration_at: 2026-09-03T21:37:20Z
---

# Add a Three-Value Priority Field the Selector Orders By and the Board Shows

## What

Add one optional frontmatter field to the REQ schema, `priority`, with the closed set `now`, `next`, `later`; absent reads as `next`, anything else normalizes to `next` with a warning under the Schema Read Contract. The selector orders ready work inside its ordinary class by priority (`now` before `next` before `later`) and keeps the existing queue order inside each value; the gate-repair and deferred-parent classes stay above it and `depends_on` stays the hard gate. The board sorts the pending column by priority inside its ready and waiting groups and shows a small `now` or `later` tag on the card (no tag for `next`, the default), the way impact tags are shown today. Capture may set the field from the user's words at Step 1 and addenda may change it. The landing commit stamps the current queue per the 23:20 triage table in `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html`: the build-now set `now`, the deferred set `later`, the rest untouched.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Added and validated the schema → typed model → selector → board → capture → release plan, then corrected the UR-target, browser-test, and addendum seams during exploration.
- [x] **[APPLY]:** Implemented the closed priority enum, distinct typed projection, ordinary-only stable ordering, Ready/Waiting board ordering and tags, capture contracts, release mirrors, and exact pending queue stamp.
- [x] **[UNIFY]:** Reviewed the 23-file builder manifest plus release/queue integration paths; focused/full Go tests, vet, contract checks, strict light/dark Chromium evidence, mirror checks, and diff hygiene passed. Generated output, screenshots, foreign claims, and REQ-530's archive were untouched.

## Why

"can we make a priority list that would reflect in the kanban board as well". A targeted run keeps caller order but lives in a shell command; the board cannot show it and the next session cannot see it.

## Context

- Selection order today (`work-reference.md`, Selection Order): `repository_gate_repair: true` first, ready `gate_deferred: true` parents second, ordinary work third, existing queue order inside each class. Priority sorts inside the third class only.
- The board's pending column already splits ready from waiting on dependencies and shows impact tags on cards; priority is one more sort key and one more tag, not a new column.
- REQ-530 (order ready work by the newest REQ it unblocks) is superseded: with a priority field the maintainer says what is urgent instead of the selector inferring it from REQ numbers. Its rule is not adopted as the tie-break; existing queue order stays the tie-break so the change is one field and nothing else. Cancel REQ-530 with this REQ's landing hash.
- The triage table in the velocity report (commit a524cf16 and later) is the first priority list: 28 REQs to build now (the REQ-502 chain, REQ-559, REQ-560, REQ-542, REQ-539, REQ-475, REQ-483, REQ-496, REQ-512, REQ-544, REQ-514, REQ-515, REQ-485, REQ-527, REQ-534, REQ-547, REQ-490, REQ-535, REQ-536, REQ-545), 13 deferred (REQ-549 to REQ-558, REQ-482, REQ-486, and REQ-530 until cancelled). Read the report's current table at build time; it may have moved.

## Detailed Requirements

- Schema: `priority` documented in `work-reference.md` Request File Schema and the Schema Read Contract table (closed set, default `next`, normalization warning), and in `capture-reference.md`'s Simple REQ template as a commented optional line; `capture.md` Step 1 gets a one-line priority assessment: emit only when the user's words rank the work, never invent.
- do-work-cli: `schemanormalization` reads and normalizes the field; `nextselection` orders the ordinary class by it and projects the value on selected and excluded records so a caller never infers it from display text. Table-driven tests: default absent, each value, an invalid value, and a `later` REQ that is the only ready dependency of a `now` REQ (the dependency still runs first).
- Board: `model.go` parses the field (lock-step with the parser rule in prime-kanban-board.md), `generate.go` sorts the pending column by priority inside ready and waiting, `board-cards.js` renders the tag, `board.css` styles it in both themes. A model test pins the sort; the static snapshot and the live board agree.
- Stamp: the landing commit edits `priority:` on the pending REQ files named by the triage table; those are capture-editable pending files, staged by explicit path, never a claimed file.
- Version bump and changelog entry; the changelog title says what shipped.

## Constraints

- One integrating commit; the schema, selector and board land together so a stamped queue is never malformed to any reader.
- `depends_on` is never overridden by priority; a `now` REQ behind an unfinished dependency waits.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** Put `priority: now` on a pending REQ, run `do-work-cli next --format json` and regenerate the board.
**Why RED now:** the field is unknown: the selector ignores it and orders by queue number, and the board shows no tag and no ordering.
**GREEN when:** `next` lists the `now` REQ before other ordinary ready work and projects `priority` on the record; a `now` REQ behind an unmet dependency is still excluded; the board's pending column shows the `now` REQ at the top of its ready group with a tag; `later` REQs sink to the bottom of their group; and after the stamp the board's pending column reads in the triage table's order.
**Validation:** Inferred during capture; the maintainer approved the capture ("capture UR-107").

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5562 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes the board model and cards.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over budget and `slugged: partial`. Matched because this REQ changes schema normalization and selection.
- `_dev/primes/lessons-action-files.md` — 4050 tokens, over budget and `slugged: partial`. Matched because this REQ adds a schema field read by several actions.

## Full Context
See `do-work/user-requests/UR-107/input.md` for complete verbatim input.

## Triage

**Route: C** - Complex

**Reasoning:** The change adds one schema contract across capture, canonical selection, typed result projection, board parsing, ordering, client rendering, queue state, and release metadata. It also requires TDD and browser evidence.

**Planning:** Required

## Plan

1. Add the `priority` Schema Read Contract row and typed request-model projection, preserving raw and recognition evidence.
2. Add a distinct typed `priority` field to selected and excluded results; keep `selection_priority` as the existing class signal.
3. Sort only ordinary eligible records `now` → `next` → `later` after dependency gating and before fan-out bounding; preserve explicit target order and current stable ties.
4. Parse and warn in the board's independent schema table, then stable-sort Pending Ready and Pending Waiting separately and rebuild their compatibility union.
5. Project the effective/raw priority into static and live board payloads; render compact `now`/`later` badges and no `next` badge.
6. Prove schema, projection, class precedence, dependency gating, stable ties, fan-out exclusions, static/live agreement, and real DOM order with RED/GREEN tests and exact-URL browser evidence.
7. Update capture and schema documentation only after executable behavior is green; `work.md` remains unchanged because it already consumes canonical selector order.
8. At integration, stamp only still-pending queue files from the latest velocity report and reconcile the already-cancelled REQ-530 provenance without re-cancelling it.

Full requirement-to-file mapping and verification commands: `do-work/runs/work-2026-09-03-210100/REQ-561-plan.md`.

**Plan validation:** Every captured requirement maps to a task and every task maps back to schema, selection, board, capture, stamp, or release acceptance. Consumer fields are explicit: `priority` remains separate from `selection_priority`, and board effective/raw/unrecognized evidence is retained. Warning: 8 task groups exceed the usual three-task quality range, but splitting them would create revisions where stamped queue state has old readers, so this REQ keeps the atomic contract.

*Generated by Plan agent*

## Exploration

The implementation seams match the plan, with three corrections. UR-expanded ordering belongs in `next_targets_test.go` as well as final selector tests; priority UI gets a dedicated `priority_browser_probe_test.go` rather than coupling to the clipboard contract; and queued addenda must explicitly support setting, changing, or removing authored priority. The schema normalizer already owns trimming, case-folding, defaults, and exact warnings, so the new field is a data-row plus typed projections rather than a second parser.

Selection already has the correct pipeline: eligibility and dependency checks, stable class sorting, then fan-out bounding. The implementation adds authored priority only to the ordinary-class comparator and leaves explicit REQ target order untouched; each UR expansion adds priority before its existing depth/id tie-break. The board already establishes numeric tie order and shares one static/live payload path, so it can stable-sort Pending Ready and Pending Waiting independently in `bucketColumns`, rebuild their compatibility union, and keep every other column unchanged.

The builder worktree predates REQ-502's checkpoint-cleanup merge, but their source manifests do not overlap. Release metadata and the still-pending queue stamp remain orchestrator-only integration seams. REQ-530 is already archived as cancelled, so this lane will record any canonical post-terminal amendment refusal instead of cancelling it twice.

*Generated by Explore agent; full findings: `do-work/runs/work-2026-09-03-210100/REQ-561-exploration.md`*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work-reference.md` (modify) — schema and selection contract
- `skills/do-work/actions/capture-reference.md` (modify) — base and addendum template guidance
- `skills/do-work/actions/capture.md` (modify) — user-authored priority assessment and addendum updates
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modify) — executable enum contract
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modify) — enum/default/warning matrix
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modify) — typed value and evidence
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modify) — typed-projection inventory
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — selected/excluded JSON projection
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — result defaults and distinct axes
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go` (modify) — authored-priority constants and rank
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (modify) — UR-member ordering only
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modify) — ordinary ordering and evidence propagation
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (modify) — UR and explicit-target ordering
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modify) — class, dependency, fan-out, warning, and stable-order regressions
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modify) — real CLI JSON projection
- `skills/do-work-board/actions/board.md` (modify) — parser/order/rendering contract
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — independent parser and pending-group sort
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modify) — parser and group-order regressions
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — generated priority evidence
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — static fixture projection
- `skills/do-work-board/tools/queue-kanban/serve_test.go` (modify) — live/static agreement
- `skills/do-work-board/tools/queue-kanban/priority_browser_probe_test.go` (new) — exact-URL Chromium order, badge, theme, and geometry evidence
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — now/later badges
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — badge styling
- `CHANGELOG.md` (modify) — release entry
- `VERSION` (modify) — suite version
- `skills/do-work/CHANGELOG.md` (modify) — installed changelog mirror
- `skills/do-work/actions/version.md` (modify) — installed version mirror
- `do-work/queue/` (modify selected pending files only) — stamp the report's explicit now/later sets by exact path during integration

**Files I will NOT touch:** `skills/do-work/actions/work.md`, primes/lessons, clipboard/browser harness tests, generated board output, screenshots, foreign claimed files, or REQ-530's terminal archive unless a canonical amendment path accepts the exact landing hash.

**Builder boundary:** the isolated builder owns only the source/action/test files above through `board.css`; the orchestrator owns version/changelog mirrors, exact queue records, terminal-provenance reconciliation, and final integration.

**Acceptance criteria:** the schema defaults and warns canonically; JSON exposes distinct class and authored-priority fields; dependencies and special classes remain authoritative; fan-out runs after ordinary priority ordering; static and live Pending Ready/Waiting groups agree; only now/later badges render; and the current pending queue reflects the handoff triage without editing a claimed or held REQ.

## Pre-Flight

- `skills/do-work/tools/checks/preflight.sh`: working tree clean outside `do-work/`; no source drift detected.
- Canonical gate evidence: exact revision `056d54d68fa7ead5e370cb184c1c5505038b4b5a` matches a persisted green `bash _dev/tests/maintainer-verify.sh` run.
- Integration hazards: the builder starts from claim commit `4e765265`; integrate serially onto current main and verify the exact branch manifest. REQ-502's merged files do not overlap this scope.

## Implementation Summary

Added `priority: now | next | later` to the executable schema and independent board parser with absent/invalid fallback to `next`. Canonical selection now retains repair/deferred-parent class precedence and dependency gating, then orders only ordinary work by authored priority before fan-out. Selected and excluded JSON expose authored `priority` separately from `selection_priority`; explicit REQ order remains caller-authored and UR expansions sort internally.

The board stable-sorts Pending Ready and Pending Waiting independently, keeps their union Ready-then-Waiting, projects invalid provenance, and renders only `now`/`later` badges. Capture emits priority only from explicit ranking language and lets queued addenda set, change, or remove it. The release advances to 0.273.0 and stamps the report's 27 still-pending build-now REQs as `now` and 12 deferred REQs as `later`.

## Decisions

- **D-01 — Keep both priority axes.** `selection_priority` remains the internal scheduling class; authored `priority` is a separate typed field. Overloading the old field would erase repair/deferred-parent evidence.
- **D-02 — Preserve explicit target order.** Priority affects default scans and each UR expansion, never an explicitly ordered REQ token list.
- **D-03 — Reuse the shared static/live projection test.** `serve_test.go` decodes generated static data and live server data from real fixture files; a second `generate_test.go` fixture would duplicate the same seam.
- **D-04 — Do not amend REQ-530 by hand.** REQ-530 was already archived as cancelled before this landing and the canonical CLI exposes no post-terminal amendment command. Its supersession is recorded here with this integration rather than fabricating a second cancellation or mutating terminal bytes outside an authority.

## Discovered Tasks

- No critical follow-up. A second live-HTTP DOM probe would duplicate the tested shared client/payload route and remains report-only.

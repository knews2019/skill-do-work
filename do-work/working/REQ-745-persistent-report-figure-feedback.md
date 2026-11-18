---
id: REQ-745
title: 'Save likes, dislikes and notes on reports and figures'
status: claimed
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: frontend
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-feedback
related: [REQ-746, REQ-747, REQ-748]
depends_on: [REQ-744]
claimed_at: 2026-11-18T19:09:13Z
route: C
estimate:
  p50_active_minutes: 60
  confidence: low
  calculated_at: 2026-11-18T19:13:11Z
  basis:
    - Route C
    - 7-file write set
    - 2 new files
    - 3 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - browser evidence
    - persistence changes
    - cross-route regression gates
planning_at: 2026-11-18T19:13:12Z
write_set:
  - justfile
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go
  - skills/do-work-toolbox/docs/ai-report-guide.md
---
# Save likes, dislikes and notes on reports and figures

## What
Open the gallery with `just run-report-gallery`. From that page mark a whole report or an individual figure as liked or disliked. A short note is optional.

## Detailed Requirements
Keep report and figure feedback distinguishable, including when a user likes a report and dislikes one figure in it. Persist marks so both the reopened page and the round-generation agent can read the same feedback.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** Open with the Just recipe, like a report without a note, dislike one figure with a note, close and reopen. Both distinct targets, their marks and the note are restored and machine-readable.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8851 tokens) covers alternate report writers and gallery routing. The index marks it partially slugged.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (25210 tokens) covers CLI result identity, rooted paths and alternate writers. The index marks it partially slugged.
- `_dev/primes/lessons-shell-commands.md` (9341 tokens) covers the Just command, argv and launch/persistence probes. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.

## Triage

**Route: C** - Complex

**Reasoning:** Durable project feedback requires a writable gallery route, controls for report and figure targets, and a machine-readable record shared with round generation.

**Planning:** Required

## Plan

1. Add the maintainer Just entry point and a loopback gallery run mode on the existing CLI. Reuse the gallery renderer for curated and installed checkouts. Preserve the static regeneration and refusal contract. Keep the live page and feedback on a single local origin.
2. Save report- and figure-level marks with an optional note in a project-root JSON record at `do-work/report-gallery-feedback.json`. Use the report page path plus a stable per-report figure key for figure targets. Upsert one mark per target, preserve sibling targets, and make writes atomic. Document the location and record consumed by REQ-747.
3. Add independent gallery controls and verify the Just launch, both target types, reload restoration, on-disk readback, clear error states, and existing gallery behavior. Exercise the real existing phase-by-phase report as trial evidence without writing real user marks during development.

**Validation.** The three tasks cover launch, target separation, persistence and the next-round consumer. Feedback entries need target kind, report page identity, figure identity when applicable, mark and optional note. The browser should confirm a save only after the endpoint acknowledges the matching target and stored state.

*Prepared inline by the pre-build owner*

## Exploration

- `justfile` has no `run-report-gallery` recipe. The existing `ai-report-gallery` CLI publishes `docs/report-gallery/index.html` under a generator-header guard and provides read-only `--list` and `--resolve`. Its main-checkout route and the installed consumer resolve assets under different roots.
- The Go page template owns the curated worked-card and generic sample views. `report_gallery_worked.go` reads a ten-card preview manifest and validates each source. Reuse report paths and preview relationships. A PNG preview is an image of a page and provides no intrinsic figure identity. Inspect published HTML as data and derive stable figure targets there.
- The catalog walker already finds actual reports in direct and monthly bundles. Use `ai-reports/2026-11-11_1411_commit-22f82fed-work-procedures-load-by-phase/` and source commit `22f82fed` for the common real-evidence candidate. Round generation and genuine marks belong to REQ-747/748.
- An installed gallery without a preview manifest still needs the same live feedback behavior when there are generated reports. Keep the static file usable and preserve generation, list, resolve, and the hand-made-page ownership guard. Add feedback script/event hooks separately from the starter prompt copy code.

*Prepared inline by the pre-build owner*

## Scope

**Files I will touch**
- `justfile` (modify) — maintainer `run-report-gallery` entry point
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modify) — live mode dispatch through the existing gallery command
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modify) — live report and figure controls and feedback client hookup
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` (new) — local serving, target validation and persistent feedback
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback_test.go` (new) — independent targets, write integrity and route coverage
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modify) — static/curated renderer regression assertions where needed
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — launch, persistence and machine-readable handoff

**Acceptance criteria**
- [ ] `just run-report-gallery` opens the working gallery.
- [ ] The gallery can save a report like with no note and a dislike on one figure with a short note.
- [ ] Closing and reopening restores each mark on its original target.
- [ ] A round-generation agent can read the same distinct targets and note from the persisted data.

## Pre-Flight

**Git.** Main at f3840b087913d8501267dda1c1395580304f2ba4. Changes are limited to current lifecycle records and run evidence.
**Tests baseline.** Eight focused gallery/catalog tests passed. The canonical focused baseline launched and exited 0.
**Canonical repository gate.** Direct exit 0 at f3840b087913d8501267dda1c1395580304f2ba4 in 158 seconds under the run lock.
**Dependencies.** The Go and shell gate ran successfully. The fresh installed consumer and Just probes launched.
**Feature RED.** The Just recipe and live CLI mode are absent. The curated gallery has a real worked card and figure but no feedback controls or persisted record.
**Evidence.** [Typed preflight](../runs/work-2026-11-18-191126/REQ-745-advance-preflight.json) · [Focused baseline](../runs/work-2026-11-18-191126/REQ-745-focused-baseline.json) · [Feature RED](../runs/work-2026-11-18-191126/REQ-745-red.log) · [Shared gate](../runs/work-2026-11-18-191126/canonical-gate-status.log)

*Checked by work action*

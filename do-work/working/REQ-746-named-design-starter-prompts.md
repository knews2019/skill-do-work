---
id: REQ-746
title: 'Copy a starter prompt that selects a specific report design'
status: claimed
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-feedback
related: [REQ-745, REQ-747, REQ-748]
depends_on: [REQ-744]
claimed_at: 2026-11-18T19:09:13Z
route: C
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-11-18T19:13:11Z
  basis:
    - Route C
    - 6-file write set
    - 1 new files
    - 3 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - browser evidence
    - cross-route regression gates
planning_at: 2026-11-18T19:14:10Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generated.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go
  - skills/do-work-toolbox/actions/ai-report.md
  - docs/report-gallery/index.html
---
# Copy a starter prompt that selects a specific report design

## What
Give every design a copyable do-work-toolbox ai-report starter prompt. Pasting it makes the writer use the named design directly.

## Detailed Requirements
Cover each shipped sample, each user-added sample and each generated report design. Show the full usable prompt and a copy control on the gallery page. An explicit selection wins even when automatic ranking would choose another sample. Preserve automatic ranking for ordinary requests.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** For each of the three design sources, copy the starter prompt and use it for a target whose automatic match would differ. The writer uses the named design with the new target’s own evidence.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8851 tokens) covers alternate report writers and gallery routing. The index marks it partially slugged.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (25210 tokens) covers CLI result identity, rooted paths and alternate writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.

## Triage

**Route: C** - Complex

**Reasoning:** Shipped and project samples already have a named-design path. Generated reports need a design identity that the gallery, resolver, and installed writer can all use.

**Planning:** Required

## Plan

1. Keep the existing installed `ai-report --design` route, sample cards, and title/slug round trip. Add an exact generated identity distinct from sample names. Discover generated bundles with the existing report walk, including month folders and curated worked cards. Use `ai-reports/<bundle>/gallery-design.md` for a generated design's design-only record and resolve existing cataloged report pages as legacy sources when that record is absent.
2. Render a visible starter and copy control for each generated design and surface shipped sample prompts alongside project sample prompts on the curated gallery. Preserve worked-card name copy, fragments, preview anchors, and distinct DOM targets. Keep generated cards and selection data shared between the ordinary and curated views.
3. Update the installed writer at its named-design resolver and the evidence-to-visuals seam. Treat an explicit generated source as authoritative for design moves, keep source facts out of new output, and leave automatic ranking driven by shipped/project samples. Verify clipboard text, resolution, marker/choice and different target evidence for all three source types plus the unchanged automatic case.

**Validation.** Gallery and writer consume the same canonical design identity. For the existing resolver row keep the ordering and meaning of identity, fit text and worktree-relative record/source path. Failure and ambiguous results must stop explicit resolution. The copy control must emit the argument the resolver accepts.

*Prepared inline by the pre-build owner*

## Exploration

- `report_gallery.go` already resolves a shipped or user sample by slug or uniquely matching title. `report_gallery_page.go` emits `do-work-toolbox ai-report --design "…"` on those sample cards, including shell quoting, and the installed `ai-report.md` resolves the name before Step 3 ranks alternatives.
- With `docs/report-gallery/preview-manifest.json` present the renderer filters shipped cards and the curated worked section only copies display names. An ordinary installed gallery exposes samples and latest-by-kind links but has no generated design starter. `reportIndexWalk` already gives primary pages and second-opinion paths for direct and month-folder bundles.
- Keep generated design records inside their owning report bundles at `ai-reports/<bundle>/gallery-design.md` (including the optional month parent). Use a namespaced identity shared by render and resolve, with ambiguity handled explicitly. Curated legacy cards should resolve to their own report's design, not a generic kind-wide preset.
- Use the real phase-by-phase report from commit `22f82fed` as the curated generated-source case. Exercise the installed route in a disposable project, including a project sample and a generated bundle with different target evidence. The sibling feedback builder owns the live route and feedback bindings.

*Prepared inline by the pre-build owner*

## Scope

**Files I will touch**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modify) — integrate explicit generated lookup with existing sample resolution
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modify) — starter display and copy bindings for every design source
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go` (modify) — connect curated cards to canonical generated identities
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generated.go` (new) — discovery and source resolution for generated designs
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go` (modify) — identity, route and copy-text proof
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — installed writer's design record authoring and explicit selection
- `docs/report-gallery/index.html` (regenerate) — apply the canonical renderer to the checked-in curated route

**Acceptance criteria**
- [ ] Shipped, user-added and generated designs each show a full starter prompt and copy control on their gallery surface.
- [ ] Pasted prompts select the named design even when automatic matching would choose another.
- [ ] The report uses the new target's own facts, figures and units.
- [ ] Ordinary requests still use existing automatic sample ranking.

## Pre-Flight

**Git.** Main at f3840b087913d8501267dda1c1395580304f2ba4. Changes are limited to current lifecycle records and run evidence.
**Tests baseline.** Eight focused gallery/catalog tests passed. The canonical focused baseline launched and exited 0.
**Canonical repository gate.** Reused the exact-revision green record at f3840b087913d8501267dda1c1395580304f2ba4 from REQ-745.
**Dependencies.** The Go and shell gate ran successfully. The fresh installed consumer and Just probes launched.
**Feature RED.** The fresh installed gallery shows 10 shipped starters and one project starter, and all resolve. Curated rendering shows only the project starter plus worked-card name copy. The generated design has no starter or named resolution in either route.
**Evidence.** [Typed preflight](../runs/work-2026-11-18-191126/REQ-746-advance-preflight.json) · [Focused baseline](../runs/work-2026-11-18-191126/REQ-746-focused-baseline.json) · [Feature RED](../runs/work-2026-11-18-191126/REQ-746-red.log) · [Shared gate](../runs/work-2026-11-18-191126/canonical-gate-status.log)

*Checked by work action*

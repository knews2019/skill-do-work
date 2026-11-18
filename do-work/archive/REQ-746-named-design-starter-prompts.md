---
id: REQ-746
title: 'Copy a starter prompt that selects a specific report design'
status: completed
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
commit: 2003d9d06cbecbde4fccf320e1f4cfb886b3f76f
heavy_verified_at: 2026-11-18T23:18:27Z
heavy_verified_revision: 2003d9d06cbecbde4fccf320e1f4cfb886b3f76f
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
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generated.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go
  - skills/do-work-toolbox/actions/ai-report.md
  - docs/report-gallery/index.html
completed_at: 2026-11-18T23:19:19Z
release_at: 2026-11-18T23:19:19Z
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
- [x] **[APPLY]:** Implement or execute the scoped work.
- [x] **[UNIFY]:** Reconcile the merged renderer, verify the file manifest and rerun the focused acceptance cases.

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

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modify) — integrate explicit generated lookup with existing sample resolution
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` (modify) — join live reviews to the shared cards by canonical report path
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modify) — starter display and copy bindings for every design source
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go` (modify) — connect curated cards to canonical generated identities
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generated.go` (new) — discovery and source resolution for generated designs
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go` (modify) — identity, route and copy-text proof
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modify) — verify the closed archive, latest-kind summary and unduplicated curated route
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — installed writer's design record authoring and explicit selection
- `docs/report-gallery/index.html` (regenerate) — apply the canonical renderer to the checked-in curated route

**Acceptance criteria**
- [x] Shipped, user-added and generated designs each show a full starter prompt and copy control on their gallery surface.
- [x] Pasted prompts select the named design even when automatic matching would choose another.
- [x] The report uses the new target's own facts, figures and units.
- [x] Ordinary requests still use existing automatic sample ranking.

## Pre-Flight

**Git.** Main at f3840b087913d8501267dda1c1395580304f2ba4. Changes are limited to current lifecycle records and run evidence.
**Tests baseline.** Eight focused gallery/catalog tests passed. The canonical focused baseline launched and exited 0.
**Canonical repository gate.** Reused the exact-revision green record at f3840b087913d8501267dda1c1395580304f2ba4 from REQ-745.
**Dependencies.** The Go and shell gate ran successfully. The fresh installed consumer and Just probes launched.
**Feature RED.** The fresh installed gallery shows 10 shipped starters and one project starter, and all resolve. Curated rendering shows only the project starter plus worked-card name copy. The generated design has no starter or named resolution in either route.
**Evidence.** [Typed preflight](../runs/work-2026-11-18-191126/REQ-746-advance-preflight.json) · [Focused baseline](../runs/work-2026-11-18-191126/REQ-746-focused-baseline.json) · [Feature RED](../runs/work-2026-11-18-191126/REQ-746-red.log) · [Shared gate](../runs/work-2026-11-18-191126/canonical-gate-status.log)

*Checked by work action*

## Implementation Summary

**Files changed:**
- `docs/report-gallery/index.html` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generated.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go` (modified)

**What was done:** Merged the named-design writer and copyable starters from `8e43a728` into the live gallery from 0.329.0. The rendered card and live feedback share the report's canonical primary path. A generated report without a sidecar reuses the layout of its own primary page. Older reports sit behind a closed disclosure and the ten curated worked reports keep their existing cards and previews.

## Qualification

Qualification and scope both returned `satisfied` over `e1f4e357..2003d9d0`. One filename-based wiring warning was reviewed. The dispatcher calls `resolveReportGalleryGenerated`, the renderer calls `reportGalleryGeneratedSources`, and the focused tests exercise both entry points. The nine-path diff is limited to the declared gallery and writer paths. The preview manifest, captures, original reports, trial evidence and feedback storage schema are unchanged.

## Testing

- `bash _dev/tests/maintainer-verify.sh` passed in 149 seconds, with 428 board tests, 1,002 CLI tests, vet, formatting and the shipped contracts. Request-bound `advance` marked scope, the original eight-test focused probe and the exact-revision green gate satisfied.
- `go test ./internal/toolboxcommands -run 'Test(AI)?Report(Gallery|Index)' -count=1 -v` passed all 18 top-level tests. Toolbox vet and the implementation-range whitespace check passed.
- **Red-green validation:** The original test lacked curated, shipped and generated starters on the base. The integration extension found two authored cards, each missing one control, and no legacy starter. Both cases passed after joining cards by primary path. An existing latest-by-kind test now checks that summary and separately asserts the closed older-report disclosure.
- The actual source and fresh installed Just recipes served the gallery. In Chrome, the monthly authored and legacy cards each copied their own starters. The authored card saved a report like and a figure dislike with a note. The same card copied its starter and restored both marks from disk after restarting the process and browser. The second figure stayed unmarked. Its link and the parent, child and opinion links reached their original targets.
- Curated Chrome copied all ten worked starters and verified their preview and fragment hooks. The source archive held 40 other reports behind a closed disclosure. The checked-in static gallery regenerated byte-for-byte and passed 28/28 checks. The installed static gallery had 23 starters and no live controls.
- Used four Chrome clipboard prompts with the fresh installed action, all targeting REQ-750's authentic verification receipts. Named shipped, added, authored and legacy design choices differed from the automatic `verdict-with-checks` choice. Each of the five bundles used 7 checks per view, 28 passing outcomes, 0 findings, and widths of 1440 px and 390 px. All five passed the installed four-view render judge.

*Verified by work action*

## Review

**Acceptance:** Pass. The in-session review followed starter text from card to installed resolver, selected design, and the HTML's choice, visuals and evidence. It also followed each review target from card to request, saved record and reopened browser state. The static and live views share card construction, and the automatic listing retains its ordinary sample-only contract.

**Scope and risk:** The two added work paths remove the redundant review-card inventory and update the assertions to cover the combined path. The gallery still discovers reports at startup. A legacy design uses the report's own primary page, and an invalid authored record leaves the report reviewable without a starter. `--design` selects a layout for completed-work reports.

**Follow-ups created:** None.

*Reviewed in-session*

## Lessons Learned

**What worked:** Assert starter and feedback controls within a single card selected by primary report path before invoking either. Independent checks for a copy button and for saved feedback let duplicate cards through.

**Worth knowing:** The queue board's strict browser lane checks known names on PATH. On this Mac it needed the installed Chrome executable specified through `QUEUE_KANBAN_BROWSER`.

## Orientation

[MAP CHANGED] Gallery cards share a canonical primary report path for design selection and live feedback. Curated cards take precedence over the older-report disclosure. Generated starters resolve to an authored sidecar or the primary page, while the ordinary sample list drives automatic ranking.

## Decisions

- **D-01 — Design identity.** Use the bundle basename when unique so moving it into a month preserves the starter. On a collision, qualify the identity with its full bundle path.
- **D-02 — Combined cards.** Carry the validated primary path separately from the clickable URL and use that path to join controls. Keep records with valid sidecars visible, and list other non-curated project reports once behind a closed disclosure.
- **D-03 — Writer proof.** Use the archived REQ-750 checks and captures as the new target evidence. Keep the combined-route latency fixture and the September 6 trial evidence separate.
- **D-04 — Release.** The additive starter coverage warrants 0.330.0. The suite owns the root version and changelog and the three installed mirrors.

## Discovered Tasks

None.

## Heavy Verification Plan

Base `e1f4e35740389c1ec028397991ae1fc540053201`, target `2003d9d06cbecbde4fccf320e1f4cfb886b3f76f`. The typed planner selected every lane because coverage is uncertain for `docs/report-gallery/index.html`.

- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Heavy Verification Result

Target and execution revision `2003d9d06cbecbde4fccf320e1f4cfb886b3f76f`. All six lanes executed and exited zero under the serial gate lock.

- `queue-kanban-javascript` — 8 seconds.
- `queue-kanban-browser` — 81 seconds, 39 tests. The initial run skipped this lane because Chrome was outside PATH. The rerun supplied the installed executable through `QUEUE_KANBAN_BROWSER`.
- `do-work-cli-integrations` — 72 seconds.
- `staged-skills` — 35 seconds.
- `updater` — 73 seconds.
- `installer` — 32 seconds.

[Integration receipt](../runs/work-2026-11-18-191126/REQ-746-integration-receipt.md)

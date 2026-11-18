---
id: REQ-745
title: 'Save likes, dislikes and notes on reports and figures'
status: completed
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
  - _dev/tests/install-suite-behavior.sh
  - skills/do-work-board/justfile.template
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go
  - skills/do-work-toolbox/docs/ai-report-guide.md
builder_handback_at: 2026-11-18T19:36:57Z
integration_at: 2026-11-18T19:37:04Z
review_at: 2026-11-18T22:43:41Z
kb_status: pending
remediation_at: 2026-11-18T22:54:31Z
re_review_at: 2026-11-18T22:58:15Z
commit: 8567dede1b2c30ed5b6867088700159cbedaf17a
heavy_verified_at: 2026-11-18T23:02:14Z
heavy_verified_revision: 8567dede1b2c30ed5b6867088700159cbedaf17a
completed_at: 2026-11-18T23:02:31Z
release_at: 2026-11-18T23:02:31Z
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
- [x] **[APPLY]:** Implement or execute the scoped work.
- [x] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

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

**Files I will touch:**
- `justfile` (modify) — maintainer run-report-gallery entry point
- `skills/do-work-board/justfile.template` (modify) — installer-managed consumer entry point
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modify) — live mode dispatch through the existing gallery command
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modify) — live report and figure controls and feedback client hookup
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` (new) — local serving, target validation and persistent feedback
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback_test.go` (new) — independent targets, write integrity and route coverage
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modify) — static/curated renderer regression assertions where needed
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — launch, persistence and machine-readable handoff
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modify) — include the new managed recipe in collision coverage
- `_dev/tests/install-suite-behavior.sh` (modify) — verify the named gallery recipe in a fresh bootstrap

**Acceptance criteria (restated from REQ):**
- [x] `just run-report-gallery` opens the working gallery.
- [x] The gallery can save a report like with no note and a dislike on one figure with a short note.
- [x] Closing and reopening restores each mark on its original target.
- [x] A round-generation agent can read the same distinct targets and note from the persisted data.

## Pre-Flight

**Git.** Main at f3840b087913d8501267dda1c1395580304f2ba4. Changes are limited to current lifecycle records and run evidence.
**Tests baseline.** Eight focused gallery/catalog tests passed. The canonical focused baseline launched and exited 0.
**Canonical repository gate.** Direct exit 0 at f3840b087913d8501267dda1c1395580304f2ba4 in 158 seconds under the run lock.
**Dependencies.** The Go and shell gate ran successfully. The fresh installed consumer and Just probes launched.
**Feature RED.** The Just recipe and live CLI mode are absent. The curated gallery has a real worked card and figure but no feedback controls or persisted record.
**Evidence.** [Typed preflight](../runs/work-2026-11-18-191126/REQ-745-advance-preflight.json) · [Focused baseline](../runs/work-2026-11-18-191126/REQ-745-focused-baseline.json) · [Feature RED](../runs/work-2026-11-18-191126/REQ-745-red.log) · [Shared gate](../runs/work-2026-11-18-191126/canonical-gate-status.log)

*Checked by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/install-suite-behavior.sh` (modified)
- `justfile` (modified)
- `skills/do-work-board/justfile.template` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)

**What was done:** Added a loopback gallery launch for source and installed checkouts. The live renderer discovers curated and generated reports, exposes independent report and figure controls, and reads and atomically updates a versioned project JSON record. Static publication keeps its existing renderer and ownership guard.

## Qualification

The cumulative merge range `e6c2229e..8567dede` contains ten project files and passes typed qualification and scope checks. The two new-file warnings reflect filename lookup. The gallery dispatcher calls `serveReportGallery` in its own Go package, and `go test` discovers the test file. The prepared Scope heading lacked its required colon, and inline code formatting made the recipe name parse as a path. Corrected the declaration and retested it. The preparation and trial-evidence commits precede this range.

Both Just entry points reach the same gallery dispatcher. The live forms use structured report and figure identities to round-trip separate entries. Static output keeps its ownership guard. The corrective merges exercise the installed recipe by name in both collision and fresh-bootstrap coverage.

## Testing

The request-bound scope, eight-test focused baseline, and canonical green gate all returned `satisfied` for `8567dede`. The exact-revision `bash _dev/tests/maintainer-verify.sh` gate passed in 158 seconds, including 428 board tests, 1,000 CLI tests, vet, formatting and the shipped contracts. Ten focused gallery tests passed on the earlier corrective merge. The final merge changes only installer tests.

**Repository gate retry:** The first full gate at `e90015a2` exited 1 because the installer collision test expected 41 managed definitions and found 42. An isolated run reproduced it. The corrective merge made the collision test pass all 42 cases and the full gate exited 0 at `aa2203c9`. At that revision the heavy installer probe found a second stale count. The final merge keeps the fresh-install template comparison and pins the recipe by name, then the full gate again exited 0.

**Red-green validation:**

- `TestAIReportGalleryServeRoundTrip` first failed with `TOOLBOX-USAGE`. It passed after merging and covers process restart, independent targets, a same-ID figure in a second report, updates, JSON readback and refused writes.
- Preparation RED found no `just run-report-gallery` recipe. The source recipe and the recipe from a fresh install both launched the gallery. Saved a report like with the note field empty and a dislike with a note on the real phase-by-phase `how-diagram-title` figure. Closed each tab, stopped Just, relaunched and opened a fresh tab. Both restored the original marks and note. The installed project discovered the report in a month folder without a preview manifest.
- Direct GET after restart matched the independently read `version: 1` JSON record in each fixture. A locked installed save preserved the exact previous bytes and kept the unsaved note available for retry.
- Static regeneration was byte-for-byte identical. The existing gallery gate passed all seven checks across four views. The live source preview still opened the actual report with its procedure diagram and proof section.

**New tests added:** `report_gallery_feedback_test.go` covers durable writes, figure identity and restart. **Existing tests updated:** `report_gallery_test.go` checks the curated live and static renderer, `install_transaction_test.go` includes the new recipe in REQ-419’s collision coverage, and `install-suite-behavior.sh` checks the fresh bootstrap recipe by name.

*Verified by work action*

## Review

**Overall: 96%** | 2026-11-18T22:43:41Z

| Dimension | Score |
| --- | --- |
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** Figure fallback IDs can change when their underlying visual content or label changes. The gallery refreshes the report inventory at startup, so a running gallery needs a restart to see newly generated reports — impact-user-visible → report only.

**Acceptance:** Pass — the in-session review verified both Just launches, restored browser state after a process and tab restart, and matching API and disk records. The integrator recorded the OS handler’s requested URLs and used fresh T3 browser tabs on macOS. The builder exercised the default browser launcher. All marks were fixtures. The review of the final corrective merge confirmed that the gallery runtime bytes and reviewed behavior were unchanged.

**Restatement sweep:** Checked the Just recipe, CLI usage, guide, schema and managed recipe assertions. The catalog and static output guard still agree with their callers.

**Suggested testing:** 0 items.
**Follow-ups created:** None (1 finding report only).

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Run both source and installed Just recipes on the same real report. Restart the server and parse the JSON independently to distinguish file persistence from browser memory.

**What didn't:** The focused gallery suite could not catch the two old managed recipe counts. The static capture gate must run against assets in its own checkout because it resolves trusted capture URLs relative to the runner.

**Worth knowing:** The Scope heading needs its canonical colon. Scope bullets should use backticks only around file paths because the checker reads every backticked token.

## Orientation

[MAP CHANGED] The report gallery runs as a loopback service with separate report and figure feedback forms. The gallery and future report rounds share one versioned JSON record in the project.

## Decisions

- **D-01 — Installed boundary.** The managed Just template belongs in the effective scope because it owns the fresh consumer’s launch recipe.
- **D-02 — Figure identity.** Prefer an authored visual ID, then an SVG title ID, then stable figure content. Scope that identity to the exact project-relative report path.
- **D-03 — Live inventory.** Discover project reports at launch. Serve active report content in an opaque-origin sandbox and retain the existing static output path.
- **D-04 — Save conflicts.** Serialize writes with an exclusive claim. Retain saved bytes and the unsaved draft on a refused update so the user can retry.
- **D-05 — Corrective scope.** Extend the installed collision case count and explicitly test the named recipe. In the fresh-bootstrap probe, use the existing full template comparison plus a named recipe assertion in place of a second hardcoded count.
- **D-06 — Release.** The suite owns the version and changelog mirrors. The new recipe and feedback store are additive and warrant release 0.329.0.

## Remediation

The first review accepted the gallery after the ordinary gate passed at `aa2203c9`. The subsequent installer heavy lane failed on its stale bootstrap recipe count. The corrective merge `8567dede` updates both installer assertions. The final range has a fresh qualification, test gate, and in-session re-review.

## Heavy Verification Plan

Target range `e6c2229eba1605dd4cf35212d51edb10d651c288..8567dede1b2c30ed5b6867088700159cbedaf17a`. The root Justfile lacks a lane mapping, so the planner selected all six lanes. This integration ran the request-scoped drain under the serial full-gate lock.

- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. coverage is uncertain for: justfile.
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. coverage is uncertain for: justfile.
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. coverage is uncertain for: justfile.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. coverage is uncertain for: justfile.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. coverage is uncertain for: justfile.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. coverage is uncertain for: justfile.

## Heavy Verification Result

Target and execution revision `8567dede1b2c30ed5b6867088700159cbedaf17a`. All six lanes ran under the serial lock with exit 0.

- `queue-kanban-javascript` — executed in 7 seconds.
- `queue-kanban-browser` — executed in 86 seconds.
- `do-work-cli-integrations` — executed in 74 seconds.
- `staged-skills` — executed in 36 seconds.
- `updater` — executed in 75 seconds.
- `installer` — executed in 38 seconds.

[Integration receipt](../runs/work-2026-11-18-191126/REQ-745-integration-receipt.md)

## Timing

Observed 2026-11-18T19:18:42Z to 2026-11-18T22:56:35Z: 3h 37m 53s total, 26m 02s attributed across 4 events, 3h 11m 51s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 18m 15s | 1 |
| verification-gate | 7m 47s | 3 |

Slowest stage: builder-work / build, 18m 15s, outcome success.
Slowest command: verification-gate / maintainer-verify, 2m 42s, exit 1, bash (2 argv tokens).

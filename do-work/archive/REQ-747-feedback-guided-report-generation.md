---
id: REQ-747
title: 'Generate traceable report variations with GPT-6.2-astra in T3 Code'
status: completed
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: general
route: C
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-feedback
related: [REQ-745, REQ-746, REQ-748]
depends_on: [REQ-745, REQ-746]
claimed_at: 2026-11-18T23:21:30Z
planning_at: 2026-11-18T23:22:33Z
dispatch_at: 2026-11-18T23:25:49Z
builder_handback_at: 2026-11-18T23:40:25Z
integration_at: 2026-11-18T23:40:36Z
review_at: 2026-11-18T23:47:52Z
commit: 4f9410eff23b3a42dde67cb18ee9ab8e75df274d
heavy_verified_at: 2026-11-18T23:51:44Z
heavy_verified_revision: 4f9410eff23b3a42dde67cb18ee9ab8e75df274d
estimate:
  p50_active_minutes: 70
  confidence: low
  calculated_at: 2026-11-18T23:22:33Z
  basis:
    - Route C
    - 8-file write set
    - 3 new files
    - 3 subsystems involved
    - 6 acceptance criteria
    - dependency depth 2
    - browser evidence
    - persistence changes
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-toolbox/SKILL.md
  - skills/do-work-toolbox/actions/ai-report.md
  - skills/do-work-toolbox/actions/ai-report-gallery-generate-reference.md
  - skills/do-work-toolbox/actions/help.md
  - skills/do-work-toolbox/docs/ai-report-guide.md
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generate.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generate_test.go
completed_at: 2026-11-18T23:51:56Z
release_at: 2026-11-18T23:51:56Z
---
# Generate traceable report variations with GPT-6.2-astra in T3 Code

## What
Provide an agent command that reads saved feedback and generates the next round using GPT-6.2-astra through T3 Code.

## Detailed Requirements
Each generated report records its parent report, the specific feedback marks used, and the changes it made. Keep report and figure marks attributable to their original targets. For the initial generation record the real source and any absence of prior marks honestly. Every generated design gets the same copyable named-design starter prompt. Factual or unit errors disqualify a report regardless of design approval.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8,851 tokens) covers the report writer and action routing. The partially slugged index row allows only a whole-file match, which exceeds the budget. The action-file prime still directs the builder to its relevant lessons.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** With known marks, generate a variation through the named model and provider. Inspect its parent, marks and change record and try its starter prompt. A liked example with an incorrect unit is disqualified.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [x] **[APPLY]:** Implemented and merged the scoped files.
- [x] **[UNIFY]:** Inspected the eight shipped paths and merge diff. The lifecycle Testing and Review sections hold the end-to-end evidence.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.

---

## Triage

**Route: C** - Complex

**Reasoning:** This action connects persistent gallery feedback, T3 generation, reports with reusable design identities and factual review. The named-model execution and candidate selection both need direct evidence.

**Planning:** Required

## Plan

1. Route a named gallery generation action through the toolbox. Pin the source bytes, exact parent path and feedback rows before generation. Use T3 delegation for the requested Astra model and obtain the actual model, terminal status and execution IDs from T3.
2. Write report and figure changes with anchors and a `gallery-design.md` through the existing writer contract. Review the rendered claims against the pin, keep the verdict bound to the reviewed bytes, and gate selection before considering design marks.
3. On an isolated project, generate an initial report and child using the action. Save synthetic report and figure marks through the gallery. Verify copied starters and a rendered change. Inject a wrong unit for a liked candidate and prove an explicit selection exclusion.

## Exploration

The ordinary `ai-report` action owns the writer and render judge. `report_gallery_feedback.go` owns the versioned read path and canonical feedback identities. `report_gallery_generated.go` resolves the published `report:` names from cataloged pages and `gallery-design.md`, using qualified names for ambiguous basenames. The September 6 source receipts distinguish elapsed wall seconds, accounted CPU seconds and tool invocation counts. The coordinator's contract and integration receipts identify the existing checkpoints for all three surfaces.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/SKILL.md` (modify) to expose the agent command.
- `skills/do-work-toolbox/actions/ai-report.md` (modify) to route the generation form.
- `skills/do-work-toolbox/actions/ai-report-gallery-generate-reference.md` (new) for the workflow.
- `skills/do-work-toolbox/actions/help.md` (modify) for discovery.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) for usage.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modify) for any deterministic command dispatch.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generate.go` (new) for deterministic gating and records.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generate_test.go` (new) for focused regressions.

**Acceptance criteria (restated from REQ):**
- [ ] The named action generates via real GPT-6.2-astra T3 delegation and records completed execution identity.
- [ ] The initial report pins evidence with an explicit absent parent and no invented marks.
- [ ] The child preserves the evidence pin, exact parent and report and figure mark snapshots and records anchored changes.
- [ ] Every generated design is reusable using its copyable named starter.
- [ ] Factual and unit review bars an incorrect candidate despite a like, including on selection.
- [ ] Synthetic acceptance leaves trial feedback and trial generation to the coordinator.

## Pre-Flight

**Git:** Clean at the preflight start on main `48c6a5532b1d93e1003fb322714fe37735380775`. Current edits are REQ-747 lifecycle and brief only.
**Tests baseline:** All eight gallery and catalog checks passed. `bash _dev/tests/maintainer-verify.sh` exited 0 in 150 seconds with 428 board tests and 1,002 CLI tests. The request-bound preflight and green-gate records are satisfied.
**Dependencies:** Go, ShellCheck, Chrome and Playwright are available. T3 Code offers Codex child tasks in the current environment.

*Checked by work action*

## Decisions

- **D-01** One explicit parent per child keeps the feedback attributable and makes rejection reasons visible. Each invocation generates one candidate.
- **D-02** The caller takes model identity from the completed T3 task and independently reviews the rendered content. A generation sidecar records provenance and a separate digest-bound review sidecar controls eligibility.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/ai-report-gallery-generate-reference.md` (new)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generate.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_generate_test.go` (new)

**What was done:** Added an installed gallery generation action and the supporting CLI to pin evidence, carry exact report and figure feedback, validate anchored changes, record T3 identity and exclude parents without a current factual pass. The action uses the existing gallery starter resolver.

## Qualification

Passed typed qualification and scope comparison for `d60c991e4001e6f31703ac9a2e3c3f6e33b5dd7d..4f9410eff23b3a42dde67cb18ee9ab8e75df274d`. The two filename-based wiring warnings resolve through `handleAIReportGallery` calling `handleReportGalleryGeneration` and Go discovering `TestAIReportGalleryGenerationLineage`. The installed route and the native test exercised both.

## Testing

**Tests run:** Direct `bash _dev/tests/maintainer-verify.sh` on the merged implementation. The request-bound baseline probe ran the eight existing gallery/catalog checks.

**Result:** Maintainer gate exited 0 in 141 seconds with 428 board and 1,003 CLI tests. The canonical green record is bound to the implementation commit. Scope and focused comparison were satisfied.

**Red-green validation:** The builder ran `TestAIReportGalleryGenerationLineage` before implementation and reproduced `TOOLBOX-USAGE` for `--prepare-generation`. It passed after implementation and in the merged test suite. The fresh installation control also returned the original usage error.

**New tests added:** `TestAIReportGalleryGenerationLineage` checks absent, unverified and stale reviews, execution identity, exact figure and report lineage, changed evidence, rejection of a liked wrong-unit parent and reusable designs.

**Heavy verification plan:** The typed planner selected CLI integrations, staged skills, updater and installer because they cover the shipped CLI and action files. See the exact lane plan below.

## Review

**Overall: 97%** | 2026-11-18T23:47:52Z

| Dimension | Score |
| --- | --- |
| Requirements | 100% |
| Code Quality | 94% |
| Test Adequacy | 94% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** The bundle-only render judge cannot fetch the sibling parent URL. Chrome verified the exact parent and child links in the live gallery origin. The child review records `parent-links-only` and four rendered views. impact-rule-change → report only

**Acceptance:** Pass for merged implementation, fresh consumer install and disposable live workflow. Native T3 metadata verified two completed real GPT-6.2-astra runs. Chrome saved a synthetic report like and named-figure dislike with a note, and copied the starters. The child applied the note. The wrong-unit test was first stale, then explicitly excluded by a current factual failure while retaining its starter.

**Restatement sweep:** The new generation form agrees across the toolbox router, help, caller reference and guide. The existing feedback JSON and qualified design resolver contracts retain their prior forms.

**Suggested testing:** Genuine trial feedback is the coordinator's live user step.

**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Save the rendered point of comparison alongside the T3 receipt and exact feedback entries. Use the live gallery origin to verify links the bundle-only judge cannot serve.

**What did not:** The layout check accepted the deliberately wrong unit. The independent factual review of the rendered claim caught it and the selection gate rejected it.

**Worth knowing:** Eligibility is tied to the reviewed bundle bytes. Re-review a changed bundle before using it as a parent. The reader-mark procedure depends on actual reader input for the trial.

## Orientation

In T3 Code, the report gallery can generate an evidence-pinned report or one feedback-guided child per invocation. Each candidate has a reusable starter. A current factual pass is required to use a candidate as a parent.

## Heavy Verification Plan

Base revision: `d60c991e4001e6f31703ac9a2e3c3f6e33b5dd7d`
Target revision: `4f9410eff23b3a42dde67cb18ee9ab8e75df274d`

- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — gallery CLI changes.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped action and CLI changes.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — shipped CLI changes.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — shipped CLI changes.

## Heavy Verification Result

Target and execution revision `4f9410eff23b3a42dde67cb18ee9ab8e75df274d`. All four selected lanes executed with evidence reuse disabled.

- `do-work-cli-integrations` — 70 seconds, exit 0.
- `staged-skills` — 35 seconds, exit 0.
- `updater` — 71 seconds, exit 0.
- `installer` — 32 seconds, exit 0.

[Integration receipt](../runs/work-2026-11-18-191126/REQ-747-integration-receipt.md)

## Timing

Observed 2026-11-18T23:25:49Z to 2026-11-18T23:41:19Z: 15m 30s total, 15m 30s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 15m 30s | 1 |

Slowest stage: builder-work / T3 GPT-6.2-sol implementation, 15m 30s, outcome success.

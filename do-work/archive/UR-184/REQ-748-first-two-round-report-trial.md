---
id: REQ-748
title: 'Run two rounds of five reports on one real evidence bundle'
status: completed
heavy_verified_at: 2026-11-19T14:24:29Z
heavy_verified_revision: c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb
review_at: 2026-11-19T14:17:57Z
kb_status: pending
commit: c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb
integration_at: 2026-11-19T14:14:35Z
builder_handback_at: 2026-11-19T14:07:12Z
dispatch_at: 2026-11-19T14:03:05Z
route: B
write_set:
  - ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html
  - ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/presentation.css
  - ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/presentation.js
  - ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-design.md
  - ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-generation.json
  - ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-factual-review.json
  - docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html.png
  - ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html
  - ai-reports/2026-11-19_1559_req748-r2-latency-ranking/presentation.css
  - ai-reports/2026-11-19_1559_req748-r2-latency-ranking/presentation.js
  - ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-design.md
  - ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-generation.json
  - ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-factual-review.json
  - docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html.png
  - ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html
  - ai-reports/2026-11-19_1559_req748-r2-replicate-variability/presentation.css
  - ai-reports/2026-11-19_1559_req748-r2-replicate-variability/presentation.js
  - ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-design.md
  - ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-generation.json
  - ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-factual-review.json
  - docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html.png
  - ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html
  - ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/presentation.css
  - ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/presentation.js
  - ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-design.md
  - ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-generation.json
  - ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-factual-review.json
  - docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html.png
  - ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html
  - ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/presentation.css
  - ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/presentation.js
  - ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-design.md
  - ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-generation.json
  - ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-factual-review.json
  - docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html.png
  - ai-reports/catalog.json
  - ai-reports/index.html
  - docs/report-gallery/index.html
estimate:
  calculated_at: 2026-11-18T23:55:11Z
  p50_active_minutes: 85
  confidence: medium
  basis:
    - Route B
    - 25-file write set
    - 22 new files
    - 2 subsystems involved
    - 5 acceptance criteria
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-gallery-feedback
related: [REQ-745, REQ-746, REQ-747]
depends_on: [REQ-747]
claimed_at: 2026-11-18T23:53:11Z
status_changed_at: 2026-11-19T13:57:59Z
completed_at: 2026-11-19T14:24:39Z
---
# Run two rounds of five reports on one real evidence bundle

## What
Carry out the first feedback trial with ten reports, five per round. The builder chooses a suitable existing real evidence bundle and records which one was used.

## Detailed Requirements
Keep the same bundle across both rounds. Deliver the first five for the user to mark, then pause. Generate the second five only after the user has marked the first five. Use the saved feedback and record lineage and changes. Check facts and units against the bundle, and disqualify an erroneous report even if the user liked it.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** Observe five first-round reports and no second-round generation before the user’s marks. After genuine user feedback, observe five second-round reports tied to the same evidence, with traceable changes and correctness decisions.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8851 tokens) covers gallery routing and report writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [x] **[APPLY]:** Implement or execute the scoped work.
- [x] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See the [verbatim request and clarification](input.md).

## Triage

**Route B (exploration-guided implementation).** The released gallery workflow supplies the operations. This trial needs evidence verification, two waves of five isolated authoring runs, a serialized gallery rebuild, and a pause for actual feedback.

**Planning** is covered by the existing generation procedure.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The 0.331.0 release provides evidence preparation, starter resolution, execution receipts, bundle digests, rendered factual reviews, and live gallery marks. The September 6 verifier re-extracted 27 rows and checked the golden outputs. The original run is the authority for the values shown. Earlier acceptance bundles were disposable.

Use bottleneck-hunt for elapsed ranking, verdict-with-checks for replicate spread, before-and-after for recorded warm/cold comparisons, parts-over-time for tool counts, and two-measures-side-by-side for CPU and elapsed time. Each Astra author works in a separate T3 worktree. Main collected the reviewed first-round bundles. The coordinator preserved five whole-report likes from the user’s “Like all five”. The second round resolves each exact reviewed parent and uses that saved mark with the same six evidence pins. Each author prepares a distinct variation in a T3 worktree. Main independently reviews the rendered facts and integrates the gallery.

## Scope

**Files I will touch:**
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html` — second-round report and anchored changes.
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/presentation.css` — local report styles.
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/presentation.js` — local report controls.
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-design.md` — named reusable design.
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-generation.json` — evidence, feedback and T3 execution receipt.
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-factual-review.json` — independent review pinned to rendered bundle bytes.
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html.png` — gallery capture of the final report.
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html` — second-round report and anchored changes.
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/presentation.css` — local report styles.
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/presentation.js` — local report controls.
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-design.md` — named reusable design.
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-generation.json` — evidence, feedback and T3 execution receipt.
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-factual-review.json` — independent review pinned to rendered bundle bytes.
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html.png` — gallery capture of the final report.
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html` — second-round report and anchored changes.
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/presentation.css` — local report styles.
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/presentation.js` — local report controls.
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-design.md` — named reusable design.
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-generation.json` — evidence, feedback and T3 execution receipt.
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-factual-review.json` — independent review pinned to rendered bundle bytes.
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html.png` — gallery capture of the final report.
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html` — second-round report and anchored changes.
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/presentation.css` — local report styles.
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/presentation.js` — local report controls.
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-design.md` — named reusable design.
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-generation.json` — evidence, feedback and T3 execution receipt.
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-factual-review.json` — independent review pinned to rendered bundle bytes.
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html.png` — gallery capture of the final report.
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html` — second-round report and anchored changes.
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/presentation.css` — local report styles.
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/presentation.js` — local report controls.
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-design.md` — named reusable design.
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-generation.json` — evidence, feedback and T3 execution receipt.
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-factual-review.json` — independent review pinned to rendered bundle bytes.
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html.png` — gallery capture of the final report.
- `ai-reports/catalog.json` — index the five variations.
- `ai-reports/index.html` — render the updated catalog.
- `docs/report-gallery/index.html` — publish the five new cards and controls.

**Acceptance criteria:**
- [ ] Preserve the five frozen parents and six source pins across both rounds.
- [ ] Generate five distinct second-round children after the genuine whole-report likes, with one verified Astra execution and one exact parent per child.
- [ ] Record each mark and anchored change, and independently check the rendered facts and units.
- [ ] Publish gallery previews, working report and figure feedback controls, deep links, and copyable starters that resolve the named designs.
- [ ] Complete canonical finalization of REQ-748.

## Interrupted round-one publication

The coordinator stopped the trial worker on the explicit handoff request. Continue from the five existing page bundles on main. The evidence manifest is `do-work/runs/work-2026-11-18-191126/trial-source/gallery-input.json`. Each bundle contains a design record, the T3 execution receipt and a factual/render review. The original candidate for Warm/cold labels had a correction. The review and diagnostic artifacts are in the local Git directory at `do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-one.lkQZp8/`, including the preserved failed version.

| First-round page | Recorded factual/render verdict | Current digest equals review |
| --- | --- | --- |
| `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-latency-ranking/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-replicate-variability/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html` | pass/pass | True |

Next integrator: verify the five main-checkout digests and any changed bytes, complete the live gallery check without writing marks, then publish the five links to the user. Preserve the recorded T3 authorship. Do not generate the first five again. Wait until the user has marked the first five before round two. Then generate five children with the same evidence manifest, verify provenance, changed figures, factual and unit verdicts, and finish canonical finalization. Use ordinary direct repository test commands.

## Round-one publication and feedback handoff

Verified all five reports from `160a96b5` in the live gallery. The [round-one handoff](../../runs/work-2026-11-18-191126/REQ-748-round-one/README.md) links the five pages and explains how to save marks and resume. The [manifest](../../runs/work-2026-11-18-191126/REQ-748-round-one/manifest.json) records the verified source pins, Astra terminal runs, reviewed bundle digests, live links, exact starters, and the original warm-label rejection.

The next wave waits for the user to save at least one current report or figure mark on each of the five cards. Verify the saved entries with the manifest and read-only preparation command in the handoff before clearing the hold. Round two remains pending.


## Reader feedback — gallery comparison

The user gave this feedback while reviewing the gallery on 2026-11-19.

> request-lifecycle, 
> cli-command-map is kind of nice, but samples are nicer.
>
> Also why don't we have preview image for all of them?

Use the samples as the stronger visual reference and retain the favorable comparison to request-lifecycle and cli-command-map. This comment refers to gallery examples. The five first-round report preferences remain to be collected. The coordinator is fixing missing actual page screenshots on generated and older report cards before the next round.

## Blocked

- [2026-11-19] blocked on "User saves at least one Like or Dislike on each of the five round-one reports in the live report gallery" — cleared by user via clarify

## Pre-Flight

**Git:** The claimed base is `caf2a8b4bbdc9b43215586d8c9e6d2957d97422c`. The request preparation and owner-written integration records are local metadata.
**Tests baseline:** Canonical preflight passed the source verification. The direct maintainer gate passed at the claimed base with 427 board and 1,004 CLI tests.
**Dependencies:** The shipped CLI, T3 Codex provider with GPT-6.2-astra, and Chromium are available.

*Checked by work action*

## Decisions

- **D-01** Each whole-report like guided one child to keep its parent's core presentation and add a focused comparison. Astra authors chose the new presentation within those directions. The feedback contained no note or figure-specific mark.
- **D-02** The review exercised feedback writes against a disposable copy of the live gallery. The tracked feedback file preserved the five genuine likes exactly.

## Implementation Summary

**Files changed:**
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-design.md` (new)
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-factual-review.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-generation.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html` (new)
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/presentation.css` (new)
- `ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/presentation.js` (new)
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-design.md` (new)
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-factual-review.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-generation.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html` (new)
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/presentation.css` (new)
- `ai-reports/2026-11-19_1559_req748-r2-latency-ranking/presentation.js` (new)
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-design.md` (new)
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-factual-review.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-generation.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html` (new)
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/presentation.css` (new)
- `ai-reports/2026-11-19_1559_req748-r2-replicate-variability/presentation.js` (new)
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-design.md` (new)
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-factual-review.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-generation.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html` (new)
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/presentation.css` (new)
- `ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/presentation.js` (new)
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-design.md` (new)
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-factual-review.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-generation.json` (new)
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html` (new)
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/presentation.css` (new)
- `ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/presentation.js` (new)
- `ai-reports/catalog.json` (modified)
- `ai-reports/index.html` (modified)
- `docs/report-gallery/index.html` (modified)
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html.png` (new)
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html.png` (new)
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html.png` (new)
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html.png` (new)
- `docs/report-gallery/previews/ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html.png` (new)

**What was done:** Published five Astra-authored second-round reports, linked each to the exact reviewed first-round report and its whole-report like, and captured the new gallery previews. Independent factual reviews pin the five approved bundles to their current digests.

## Qualification

Passed. Canonical `advance` returned satisfied qualification and scope-drift records for the exact integrated range `4ffc9b99061e239f0f0c2ccb7bd6279b4d482fa7..c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb`. The 38 declared project paths match the integrated diff. The six evidence pins, five parent bundles and saved likes were verified byte-for-byte.

## Testing

**Tests run:** The source verifier re-extracted the 27 pinned TSV rows and all nine groups. The focused integration probe verified five unchanged parents, exact mark identity, five verified T3 executions, 38 scoped project paths, and reviewed bundle digests. The independent factual checks and Chromium gallery controls passed. Node syntax checks passed.

**Result:** The direct maintainer gate at `c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb` passed in 133 seconds with 427 board and 1,004 CLI tests. Canonical `advance` recorded the focused probe and green gate as satisfied.

**Trial sequence:** At `e4c28bf5`, the five likes were recorded and no round-two page existed. The actual second-round task runs followed that commit. An isolated replay verified that an unmarked parent and the preserved factual failure each stop selection before generation.

**Heavy verification plan:**
- Range `4ffc9b99061e239f0f0c2ccb7bd6279b4d482fa7..c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb`.
- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. The planner selected this lane because the report files fall outside its declared coverage.
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. The planner selected this lane because the report files fall outside its declared coverage.
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. The planner selected this lane because the report files fall outside its declared coverage.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. The planner selected this lane because the report files fall outside its declared coverage.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. The planner selected this lane because the report files fall outside its declared coverage.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. The planner selected this lane because the report files fall outside its declared coverage.

*Verified by work action*

## Review

**Overall: 100%** | 2026-11-19T14:17:57Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None on the integrated bytes. Review clarified the latency sample order, variability's recorded-label heading, and the warm comparison title and figure label. The author commits preserve their original versions.
**Acceptance:** Pass — all five original bundles and source pins were preserved. Five second-round bundles passed independent source and unit review. Local Chromium rendered 20 wide/phone light/dark captures. The loopback gallery verified ten previews, five exact clipboard starters, eleven child figure links and all twelve before/after targets. Ten isolated report and figure saves persisted across reload. The acceptance run also verified that absent feedback and the preserved statistical-independence error both exclude a parent.
**Restatement sweep:** nothing redefined. This was the sole wave member. The diff adds data and reports under the existing writer, feedback, and design contracts.
**Suggested testing:** 0 items.
**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Pair each current parent digest with its exact feedback target before authoring. Recompute the candidate's quantities from the underlying rows, and use the gallery origin to verify sibling links that a bundle-only renderer cannot reach.

**What didn't:** The bundle-only render judge reports a 404 for a valid sibling link. Its four captures and the same link's visible anchor under the gallery origin establish the more precise `parent-links-only` result.

**Worth knowing:** Warm and cold describe labels assigned by the runner. The archived failed claim about statistical independence remains disqualified even when the feedback store contains a like on the report path.

## Orientation

The report gallery's first feedback trial contains two rounds of five reports on the same September 6 test-efficiency measurements. The second round follows the reader's five saved likes. Each variation has a reusable design starter, traceable changes and a factual review tied to its current bundle digest.

## Heavy Verification Plan

Base revision: `4ffc9b99061e239f0f0c2ccb7bd6279b4d482fa7`
Target revision: `c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb`

- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. The planner selected this lane because the report files fall outside its declared coverage.
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. The planner selected this lane because the report files fall outside its declared coverage.
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. The planner selected this lane because the report files fall outside its declared coverage.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. The planner selected this lane because the report files fall outside its declared coverage.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. The planner selected this lane because the report files fall outside its declared coverage.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. The planner selected this lane because the report files fall outside its declared coverage.

## Heavy Verification Result

Target and execution revision `c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb`. All six selected lanes executed successfully.

- `queue-kanban-javascript` — 7 seconds, executed, exit 0.
- `queue-kanban-browser` — 78 seconds, executed, exit 0.
- `do-work-cli-integrations` — 68 seconds, executed, exit 0.
- `staged-skills` — 35 seconds, executed, exit 0.
- `updater` — 69 seconds, executed, exit 0.
- `installer` — 31 seconds, executed, exit 0.

The initial pass skipped the board browser lane. The rerun used the installed Chrome app executable and ran 39 tests. The first and retry receipts retain both observations.

## Timing

Observed 2026-11-19T14:15:14Z to 2026-11-19T14:17:27Z: 2m 13s total, 2m 13s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 2m 13s | 1 |

Slowest command: verification-gate / integrated-maintainer-gate, 2m 13s, exit 0, bash (2 argv tokens).

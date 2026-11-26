---
id: REQ-757
title: 'Addendum: Show repository reports before gallery samples'
status: claimed
dispatch_at: 2026-11-25T23:12:29Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-25T23:12:22Z
  basis:
  - trivial short-circuit
route: A
created_at: 2026-11-25T22:17:40Z
user_request: UR-189
domain: frontend
prime_files: []
tdd: true
impact: impact-user-visible
effort_estimate: effort-mechanical
addendum_to: REQ-749
maintenance: false
claimed_at: 2026-11-25T23:11:32Z
builder_handback_at: 2026-11-26T00:08:46Z
integration_at: 2026-11-26T00:09:52Z
review_at: 2026-11-26T00:15:57Z
kb_status: pending
commit: 9348c4b15220d527f6325476b911abd2f00d0e07
---
# Show repository reports before gallery samples
## What
Open the gallery on reports from the current repository, newest first. Label the sample designs clearly and place them after all repository report listings. When the repository has no reports, show the samples immediately.

## Prior Implementation

REQ-749 (build the gallery from ten real report pages) added a curated worked-report section when a manifest is present. The shared template in `report_gallery_page.go` currently emits the sample section before generated and older reports, including in repositories without a worked-report manifest.

## Red-Green Proof

**RED prompt/case:** Generate and open a gallery for a repository with two dated reports, shipped samples, and no worked-report manifest.
**Why RED now:** The first grid shows shipped sample cards. The repository reports appear further down under Generated report designs.
**GREEN when:** The newer repository report appears first, the older one follows, and the sample designs follow the repository listings. With no repository reports, samples are the first cards.
**Validation:** Inferred during capture from the screenshot and discussion.

## Assets

[Gallery on opening](../user-requests/UR-189/assets/REQ-757-screenshot-1-gallery-opening.png) shows the gallery in sa2-sentence-aligner2. The first grid is headed “10 report designs” and starts with Before and After, Bottleneck Hunt, Cause Chain, and Changing Mix. Generated designs is a navigation link at the top.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Use catalog order for one repository report grid and keep samples in their own section.
- [x] **[APPLY]:** Update the shared static and live renderer, keep worked card controls, and regenerate the checked-in gallery.
- [x] **[UNIFY]:** Reviewed all seven files listed below against the merge range. Native lint and formatting checks passed in the fast suite. The range passed `git diff --check`.
*Source: [Original request](../user-requests/UR-189/input.md)*

## Triage

Route A. The gallery template and before-and-after ordering are identified, including the empty repository case.

## Plan

**Planning not required** — Route A, direct implementation.

## Implementation Summary

**Files changed**

- `_dev/tests/report-gallery-gate.cjs` (modified) checks fragment uniqueness and focuses the gallery before each full-page capture.
- `docs/report-gallery/index.html` (modified) regenerates the tracked gallery and includes the newly merged instruction-maintenance report and its captured preview.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) documents report order, preview behavior, and sample placement.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modified) renders one repository report grid before samples and reserves worked card fragments.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modified) covers mixed reports, the empty repository, unique fragments, and live feedback.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go` (modified) keeps rendered starter prompts with the presentation card.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified) shares the existing date extraction and stable catalog sort.

The gallery lists annotated, legacy, and curated repository reports together by bundle date, followed by sample designs. When no reports are present, samples supply the first cards.

Implementation range `5eb8a9fab53ce8695c1f008fd9bf4e2eec95a2a1..9348c4b15220d527f6325476b911abd2f00d0e07`. The merge's second parent is the corrected hand-back `c5505b7b35a810287441099d8152a806930cbd50`.

## Decisions

- **D-01.** Use one repository grid with inline worked report kinds to preserve chronological order.
- **D-02.** Apply the catalog date rule to manifest-only reports. Break date ties by path and place undated reports last.
- **D-03.** Reserve bare fragments for worked cards. Prefix conflicting UI IDs with `gallery_` and update their references.
- **D-04.** Focus the gallery before each capture and allow 30 seconds for the screenshot. Preserve the interaction timeout.
- **D-05.** Regenerate the checked-in gallery inside the merge from the current catalog, including the REQ-760 report and preview.

## Qualification

Passed. Request-bound `advance` accepted all seven files in the merged range without findings. Both merge parents match the captured base and corrected builder tip. The maintenance report and its preview were already tracked at the base. No queue bookkeeping is in the range.

## Testing

**Request-bound focused tests.** Ran `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run 'Test(AIReportGallery|ReportGallery|AIReportIndex|ReportIndex)' -count=1`. The probe passed, exit 0, with 0.573 seconds reported by Go. The mixed report fixture verifies date order, unique worked fragments, and one live report feedback form per report. The empty repository fixture verifies that samples supply the first cards.

**Canonical fast gate.** Ran `bash _dev/tests/maintainer-verify.sh` directly under the run's full-gate lock. One run passed, exit 0 in 145 seconds. Both Go test stages executed. The request-bound green record names `9348c4b15220d527f6325476b911abd2f00d0e07`.

**Integrated browser checks.** Ran `GALLERY_PLAYWRIGHT=/opt/homebrew/lib/node_modules/playwright node _dev/tests/report-gallery-gate.cjs --gallery docs/report-gallery/index.html --manifest docs/report-gallery/preview-manifest.json --out /tmp/req-757-integrated.WUPGWA/views`. All 28 checks passed in 123 seconds at 1440px and 390px, each in light and dark. All four full-page captures were written. A separate Chrome check confirmed 63 unique repository report links in catalog order, including the 10 worked cards, followed by 10 sample cards. It checked unique HTML IDs, opened the newest maintenance report from its preview, decoded its own tracked 1040 × 740 image, and exercised the sample enlargement dialog. [Four-view evidence](/tmp/req-757-integrated.WUPGWA/views/gates.json)

**Red-green validation**
- `TestAIReportGalleryOpensOnRepositoryReportsBeforeSamples` failed as expected against the pre-merge renderer overlaid in a scratch directory. The newer and older reports did not appear in the required order. It passed on the merged renderer, including the empty-repository case.
- The builder reproduced manifest-only ordering and duplicate fragment failures on the initial implementation. The corrected mixed fixture passed on the merged tree.

**New and updated tests**
- Added the report-first and empty-repository cases. Updated the REQ-743 gallery test and the generated and legacy preview checks to look for older reports in the main grid.
- Extended the REQ-749 worked-gallery coverage to interleave mixed reports and check live feedback. The REQ-750 browser validator still checks interactions and now also requires one element per worked fragment.

**Heavy verification plan.** The typed planner selected `queue-kanban-javascript`, `queue-kanban-browser`, `do-work-cli-integrations`, `staged-skills`, `updater`, and `installer` for `5eb8a9fab53ce8695c1f008fd9bf4e2eec95a2a1..9348c4b15220d527f6325476b911abd2f00d0e07`. Coverage is uncertain for `docs/report-gallery/index.html`. The commands and reasons are recorded below for the shared drain.

*Verified by work action*

## Review

**Overall 98%** | 2026-11-26T00:15:57Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Acceptance:** Pass for implementation and local integration. Compared the first-parent range against UR-189 and the captured RED/GREEN case. The old renderer failed the ordering assertion. The merged renderer passed the mixed and empty-repository fixtures, and Chrome confirmed the checked-in page's order and the maintenance report's local preview and destination. Reviewed the live feedback and fragment references against their updated DOM targets. The separate heavy lanes remain scheduled for the shared drain.

**Restatement sweep:** Redefined the gallery's opening order, report grouping, sample placement, and worked fragment ownership. Wave one ends here because the working records for REQ-758 and REQ-759 are archived under UR-190. Used both archived Review sweep lines to check the inherited capture constraint and fold checks, optional verification, gap reporting, repair disposition, and opt-in feedback gate across actions, guides, help, prompts, and contract checks. The old Nit definition based on passing mentions is gone from the verifier and its guide.

**Minor findings**

- Toolbox help summarizes the gallery as samples, kinds, and presets. The README's first-upgrade example also describes a gallery generated from installed samples and omits repository reports. The report-first behavior is documented in the AI report guide — impact-negligible → report only
- The capture companion still says the verifier “grades” folds in its prose-backlog and fold-record explanations. The verifier now reports gaps and the fold linkage instructions still match — impact-negligible → report only
- The README and core help full-cycle prompts still stop on “verification fails” without naming which gaps should stop that explicitly requested workflow. Ordinary capture and the tutorial both make this verification optional — impact-rule-change → report only

**Suggested testing:** 0 items.
**Follow-ups created:** None (3 findings report only).

*Reviewed by review-work action*

## Lessons Learned

**What worked**
- Regenerate a checked-in gallery on the integration tree. The builder's page predates the separately merged report, and checking only its old fixture would miss that report.

**What didn't**
- The earlier browser check failed to capture the long mobile page while the link-probe tab was active. Bring the gallery forward before taking the full-page capture.

**Worth knowing**
- The date comes from the bundle folder, including for manifest-only reports. Ties use the catalog's path order.

## Orientation

The report gallery opens on repository reports in date order. Sample designs follow and supply the first cards in an empty repository.

## Heavy Verification Plan

- Base revision `5eb8a9fab53ce8695c1f008fd9bf4e2eec95a2a1`
- Target revision `9348c4b15220d527f6325476b911abd2f00d0e07`
- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. coverage is uncertain for: docs/report-gallery/index.html.
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. coverage is uncertain for: docs/report-gallery/index.html.
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. coverage is uncertain for: docs/report-gallery/index.html.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. coverage is uncertain for: docs/report-gallery/index.html.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. coverage is uncertain for: docs/report-gallery/index.html.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. coverage is uncertain for: docs/report-gallery/index.html.

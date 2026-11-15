---
id: REQ-737
title: 'Addendum: relation-map preset keeps box titles readable after a click on a phone'
status: pending
created_at: 2026-11-15T12:45:58Z
user_request: UR-180
addendum_to: REQ-727
domain: frontend
prime_files: ["_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-738]
batch: relation-map-phone-layout
write_set: ["skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c*.html", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c*-reply.md", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html", "ai-reports/catalog.json", "ai-reports/index.html", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION"]
---
# Addendum: Relation-Map Preset Keeps Box Titles Readable After a Click on a Phone
## What
Change the shipped `relation-map` preset so that, on a phone, a click on a box leaves the clicked box's title and the related boxes' titles readable. Today one click zooms the map out to about 0.2 to 0.34, and titles become about 4 to 7.5 px.
## Prior Implementation
REQ-727 (ship a relation-map preset for ai-report second-opinion) added `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md`, a 3,057-byte preset with 20 numbered checks, and named it in the "A Second Opinion" preset list of `skills/do-work-toolbox/docs/ai-report-guide.md`. It was released as 0.324.0 on 2026-11-15 (merge `eec3bad0`, release commit `d50775cd`). The proof page was built through a Session second opinion on the tstack bundle: `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`. The builder measured it with one Playwright script that clicked every shown box from the first view and from "Show everything" at 1440 x 1000 and 390 x 844.

The phone limit is recorded as D-16 in `do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md` (line 197) and as review finding F1 (line 297). The cause is check 15 at `relation-map.md:17`: "A click never leaves the clicked box or a related shown box outside the frame, and it zooms in only as far as all of them still fit." On a 390 px phone, that rule forces a click on a box with many relations (O2 in the test page has 11) to zoom out to 0.195 to 0.326. At that scale 0 to 1 relation words fit and titles are 4 to 7.5 px. The list under the map still names every relation. REQ-727's D-16 named a phone-specific layout, such as stacking the related boxes, as the likely fix, and said it needs the maintainer's taste.
## Detailed Requirements
- Below 600 px wide, a click on a box leaves its title, and the titles of the related boxes the page shows, readable without a manual zoom.
- This extends check 15 and narrows it on phones only: on a phone, readable titles win over "every related box inside the frame". Related boxes that do not fit at a readable size are reachable another way the preset names, for example a stacked phone layout of the related boxes or one tap from the detail panel. The builder picks the form and shows it to the maintainer before the REQ is called done.
- At desktop width, check 15 keeps its current behavior.
- The other 19 checks keep holding at both widths, in particular check 17 (phone first view, no sideways page scroll), check 18 (every box as text and every relation as one table row below the map) and check 19 (no-JavaScript view).
- The change stays in the preset text. Keep it short: the preset is 3,057 bytes today; stay near that size (about 3,500 bytes at most, the bound REQ-727 planned with). Removing a check only to save bytes is not allowed, because each check prevents a defect REQ-733 (rebuild the relation map page) found.
- Build one real page with the changed preset, as REQ-727 did, and measure it.
## Constraints
- Only the preset text changes. The preset loader, the second-opinion prompt, the other presets, the launcher Go code and the report kit (`skeleton.html`, `presentation.css`, `presentation.js`) stay unchanged.
- The launcher's direction block above the page is REQ-738 (the launcher's direction block keeps the map on the first phone screen), not this REQ.
- Existing pages in other bundles are not rebuilt.
- Release rules from `_dev/primes/prime-releases.md` apply: a CHANGELOG entry and a VERSION bump.
## Dependencies
None. REQ-738 is related but independent: it changes the launcher, not the preset. This REQ measures the title size after a click, which does not depend on the direction block above the map.
## Red-Green Proof
**RED prompt/case:** Open the 0.324.0 test page `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html` at 390 x 844 px. Click each box of the first view, then each box after "Show everything". Measure each shown title's on-screen size (its font size times the map's current scale).
**Why RED now:** Check 15 makes a phone click zoom out to 0.195 to 0.342, so the clicked and related titles render at about 4 to 7.5 px (REQ-727 D-16 and F1).
**GREEN when:** On a page built with the changed preset, at 390 x 844 px, after every one of those clicks: (1) the clicked box's title is fully in view and renders at about 12 px or more; (2) every related box title the page shows after the click renders at about 12 px or more, and every related box not shown is reachable in one step the page names; (3) the page never scrolls sideways, and the list under the map still names every relation. At 1440 x 1000 px, 0 clicked or related boxes are out of frame, as REQ-727 measured. The existing shipped-presets test (`TestSecondOpinionShippedPresetsLoadFromTheToolbox` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`) still passes, and the page passes `ai-report-judge`.
**Validation:** Inferred during capture. The 12 px floor matches the desktop first-view title size REQ-727 measured (12.4 px); the user asked only for "readable titles after a click".
## Builder Guidance
- Certainty: the goal (readable titles after a click on a phone) is the user's. The layout form is the builder's choice and needs the maintainer's eye: show phone captures of the new click view before the REQ is called done.
- Measure with one script that clicks every shown box from each starting view at both widths, not a sample of clicks (REQ-727's lesson in `do-work/runs/work-2026-11-15-012821/REQ-727-handback.md` line 125). Re-measure after finalize, because the launcher's direction block moves the map.
- Decided at capture: build the proof page on the same tstack bundle as REQ-727, so before and after compare directly. A new tree page or a rebuild of v1c03 are both fine.
- Keep language clear and do not complicate the preset (the maintainer's scope cue for REQ-727 on 2026-11-14).
## Full Context
See `do-work/user-requests/UR-180/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (bare; the row is `slugged: partial`, so the targeted `#session-is-the-backend` form is not legal) — cost 24,714 tokens against the 2,000-token budget — matches "a session that is its own second-opinion backend": the builder runs a Session second opinion for the proof page.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: relation-map preset needs a phone layout: readable titles after a click, and the launcher's preset text must not push the map below the first screen*

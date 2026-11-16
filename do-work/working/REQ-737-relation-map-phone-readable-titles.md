---
id: REQ-737
title: 'Addendum: relation-map preset keeps box titles readable after a click on a phone'
status: claimed
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
write_set: ["skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-*.html", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c04-*-reply.md", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html", "ai-reports/catalog.json", "ai-reports/index.html", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-16T17:17:19Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-11-16T17:20:55Z
  basis:
    - Route B
    - 9-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
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

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and the one shipped file are named, and no code changes. The work is the wording of one phone rule in the preset, then one real page built with it and measured after every click at two widths. The form of the phone layout is open, so the builder needs the earlier page, its measurements and the launcher's Session sequence before it writes.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Required-lessons consult:** read `_dev/primes/lessons-releases.md` in full. Its family canonical-link-outlives-its-target applies only at finalize: UR-180 closes when the last of REQ-737 and REQ-738 is archived, so the integrator runs the link sweep then. The row for `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` in `do-work/lessons-index.md` still matches ("a session that is its own second-opinion backend") and is still dropped for budget (section above, unchanged). The brief names its two session-is-the-backend bullets at `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:17` and `:18` as a fixed-string grep. No other row matches.

**Where the change goes:**
- The preset: `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md`, 22 lines, 3,057 bytes: one opening paragraph, one blank line, then checks 1 to 20, one sentence per line. Check 15 is line 17. Check 17 (line 19) is the only check that names 600 px, and check 18 (line 20) names the phone for the table rows. No numbered check names both a click and a phone, and that gap is this REQ.
- Room: about 440 bytes up to the 3,500 target (D-02).
- No test reads the preset's words. `TestSecondOpinionShippedPresetsLoadFromTheToolbox` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go:1007`) names relation-map in its list at line 1026 and only checks that the file loads: not empty, no double opening brace. The guide (`skills/do-work-toolbox/docs/ai-report-guide.md:124`) names the preset in one clause that says nothing about clicks, so it needs no change.
- The launcher prints the whole preset on the page: a meta tag and the header block from `secondOpinionDecorate` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:1602`, the Direction paragraph at line 1634). REQ-738 changes that block. This REQ does not.

**The before page** (`ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`, 617 lines, 110 KB, built by 0.324.0):
- `fitScale` at line 390, `select` at line 492 with its phone test at line 505 (a media query for 599 px and less), `fitShown` at line 518. The click fit follows check 15 at every width, and that produces the small phone scale.
- The box title is 22 px in map units, so a title of about 12 px needs a map scale of about 0.55 or more. REQ-727 measured the phone first view at scale 0.588 (12.9 px), and after a click 0.326 (D, O1) and 0.195 (O2, O3). O2 has 11 relations. The map is 1800 x 2480 units with 23 boxes and 28 relations.
- The page's own CSS at line 16 styles the launcher's block by its aria-label (REQ-727 D-13). REQ-738 changes that block's markup, so the new page must not style it (D-07).
- Head metas at lines 103 to 124 record the preset name and the full direction. The header block is lines 125 to 146.

**How the page is made:** the Session sequence in `skills/do-work-toolbox/actions/ai-report.md:81`: prepare, write the staging page and the staging reply, finalize. The launcher never overwrites a page, and the bundle already holds v1c01 to v1c03, so a new run on base index.html becomes v1c04 (D-01). A re-judge uses the judge's page form (`skills/do-work-toolbox/actions/ai-report.md:190`). Prepare refuses to start when Git shows a change outside the bundle, so the preset commit comes first.

**Measurement to repeat:** REQ-727's builder used one Playwright script against a local Python HTTP server on the bundle folder (`do-work/runs/work-2026-11-15-012821/REQ-727-handback.md`, Verification). The script itself was scratch and was not kept, so this builder writes a new one and keeps its path and raw output for the hand-back (D-06). Reference numbers at 1440 x 1000: 0 of 33 clicked or related boxes out of frame, first-view title 12.4 px, 12 of 12 relation words at load.

*Generated by the pre-dispatch agent (exploration done in-session, no Explore subagent)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (modify) — the phone rule for a click: readable titles win below 600 px, and the page names how the other related boxes are reached
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-*.html` (new, launcher-named) — the proof page built with the changed preset (D-01)
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c04-*-reply.md` (new, launcher-named) — the editor's reply the launcher publishes beside it
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` (modify, launcher only) — the iteration list the launcher rebuilds
- `ai-reports/catalog.json` (modify, integrator only) — catalog rebuild after the merge (D-08)
- `ai-reports/index.html` (modify, integrator only) — catalog rebuild after the merge (D-08)
- `CHANGELOG.md` (modify, integrator only) — release entry
- `skills/do-work/CHANGELOG.md` (modify, integrator only) — byte-identical mirror
- `VERSION` (modify, integrator only) — version bump (D-09)

**Files I will NOT touch:** the preset loader and every line of `report_second_opinion.go` and its test file, `ai-report-second-opinion-prompt.md`, `ai-report-guide.md`, the other twelve presets, the report kit (`skeleton.html`, `presentation.css`, `presentation.js`), the v1c03 page and every other existing file in any `ai-reports/` bundle except the launcher's iteration list, any `do-work/` path.

**Acceptance criteria (restated from REQ):**
- [ ] The preset says, as a check a reader can measure, that below 600 px wide a click leaves the clicked box's title and the shown related titles readable, and it names how a related box that does not fit is reached in one step.
- [ ] Check 15 keeps its desktop behavior, all 20 checks are still there and numbered, and the preset stays at about 3,500 bytes or less (D-02).
- [ ] One proof page is built with the changed preset on the tstack bundle through a Session second opinion, and the finalize judge and a standalone judge both report pass.
- [ ] At 390 x 844, after a click on every shown box from the first view and from "Show everything": the clicked title is fully in view at about 12 px or more, every shown related title is at about 12 px or more, every related box not shown is reachable in one step the page names, the page never scrolls sideways, and the list under the map names every relation (D-05).
- [ ] At 1440 x 1000, 0 clicked or related boxes are out of frame after every click, as REQ-727 measured.
- [ ] The other 19 checks still hold at both widths on the proof page, in particular checks 17, 18 and 19.
- [ ] The shipped-presets Go test passes, and only the preset changes among shipped files.
- [ ] The maintainer sees phone captures of the click view before the REQ is called done (D-04).

## Decisions

(from pre-dispatch)

- **D-01** (DECIDE & STATE): The proof page is a new tree page on base index.html of `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1`, made with the Session form, backend claude, the builder's own declared model slug, and the relation-map preset alone. The launcher gives it the next free id, expected v1c04. The 0.324.0 page v1c03 stays untouched as the before page, so both can be opened side by side. The launcher never overwrites a page, so a rebuild in place was not possible without deleting a published file.
- **D-02** (DECIDE & STATE): Byte bound. The target is 3,500 bytes, as REQ-727 planned. The GREEN probe refuses a preset above 3,600 bytes, because the target is approximate and one word must not fail a gate. It also requires checks 1 to 20 numbered in order with no gap. A check 21 is allowed. Removing or merging a check to save bytes is not.
- **D-03** (DECIDE & STATE): The layout form stays the builder's choice, as the REQ says. Pre-dispatch names three forms so the choice is explicit, and the builder reports the one chosen and the ones rejected. Form 1, focus and panel jumps: on a phone a click keeps the scale at about 0.55 or more, frames the clicked box with the related boxes that fit, and the detail panel lists every related box as a button that moves the map to it. Form 2, stacked related boxes: on a phone a click redraws the related boxes as one column beside or under the clicked box at a readable size. Form 3, no map move: on a phone a click only opens the panel. Pre-dispatch recommends form 1. It needs one or two sentences in check 15, it changes nothing at desktop width, and the panel already lists every relation (check 8). Form 2 needs a second layout engine in every generated page. Form 3 leaves the map unchanged, so the related boxes are not shown at all.
- **D-04** (DECIDE & STATE): The maintainer's look comes before integration, not before the build. The builder saves phone and desktop captures of the click view under the run folder (`do-work/runs/work-2026-11-16-171725/REQ-737-captures/`), and the coordinator shows them to the maintainer before integration starts. If the maintainer wants another form, that is a new build pass on the same branch.
- **D-05** (DECIDE & STATE): How "about 12 px" is measured. The title size is the title element's rendered font size on screen: its font size times the map's current scale, or the height the browser reports for it. The target is 12 px, and a value of 11.5 px or more counts, because the user asked for "readable" and the number was inferred at capture. The builder reports the smallest value over all clicks. "Fully in view" means the title's rectangle lies inside the map frame and inside the 390 x 844 screen after the page's own scroll.
- **D-06** (DECIDE & STATE): The measurement script is kept. The builder writes one Playwright script, saves it and its raw output beside the captures in the run folder, and the hand-back gives its path and the command that runs it. REQ-727's script was scratch and is gone, and this REQ's integrator must run the same script again after REQ-738 lands.
- **D-07** (DECIDE & STATE): The proof page does not style the launcher's direction block. REQ-727's page did (its D-13, a CSS rule keyed on the block's aria-label). REQ-738 changes that block's markup, and it integrates first, so a page rule written against the 0.324.0 markup could fight the new block when the integrator decorates the page again.
- **D-08** (DECIDE & STATE): The catalog pair (`ai-reports/catalog.json`, `ai-reports/index.html`) is rebuilt by the integrator after the merge, and the integrator replaces the two launcher-named globs in Scope with the resolved names. Same handling as REQ-727 D-10 and D-11.
- **D-09** (DECIDE & STATE): Release size is the integrator's call at finalize. Pre-dispatch proposes a patch bump on top of whatever REQ-738 released, because this REQ removes a known limit of a preset that 0.324.0 shipped and adds no command or file.
- **D-10** (DECIDE & STATE): tdd stays false. The shipped change is prose for a model, and no Go test can read a layout from it. The red-green proof is the browser measurement: RED on the v1c03 page (4 to 7.5 px after a phone click), GREEN on the new page. The GREEN probe pins what a shell can pin: the phone click rule exists as a numbered check, the byte bound, the 20 checks, the named Go test, and a second relation-map page in the bundle.

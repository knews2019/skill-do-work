---
id: REQ-737
title: 'Addendum: relation-map preset keeps box titles readable after a click on a phone'
status: completed
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
dispatch_at: 2026-11-16T17:31:30Z
builder_handback_at: 2026-11-16T20:00:45Z
integration_at: 2026-11-16T20:01:13Z
kb_status: pending
heavy_verified_at: 2026-11-16T20:12:30Z
heavy_verified_revision: 3f6ebf383b7290efefbaf7ea00c31bb127933d17
review_at: 2026-11-16T20:12:13Z
commit: 8e6d94e9981664297ba37dfd9239b76ccd725da3
completed_at: 2026-11-16T20:13:24Z
release_at: 2026-11-16T20:13:24Z
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
- [x] **[PLAN]:** Pass 1: one phone sentence in check 15; build the proof page through the Session form; measure RED on v1c03 and GREEN on the new page with one kept script. Passes 2 and 3 add the maintainer's page-top and label rules to checks 12 and 17 and rebuild the page on the merged 0.324.1 launcher. (from the builder hand-back)
- [x] **[APPLY]:** Preset commits e9d30a3f (check 15 phone sentence) and b8542cb5 (checks 12 and 17); pass 3 merged main at ab5e07c1 (437232b4), dropped the first-pass pages (9f4a0200) and built the final page v1c04 (9375501d). (from the builder hand-back)
- [x] **[UNIFY]:** git diff --stat main...HEAD shows 4 files: relation-map.md (3 lines changed), the new page (692+), its reply (20+), the bundle's index.html (1+); no staging file, debug output or capture in the diff; shipped-presets Go test, shipped-package-reference-contract.sh and contract-regressions.sh pass. (from the builder hand-back)
*Source: relation-map preset needs a phone layout: readable titles after a click, and the launcher's preset text must not push the map below the first screen*

---

## Addendum (mid-run)

2026-11-16T19:34Z. This extends the Detailed Requirements after the maintainer saw the phone captures of the first build (D-11, form 1) and of REQ-738's short header. The maintainer picked two options in the question dialog; their exact text follows.

On REQ-738 (the launcher's short header), the maintainer chose:

> ```
> Approve + REQ-737 note (Recommended): Integrate REQ-738 as built. Send the REQ-737 builder a mid-run note: its new proof page must show a real part of the map on the first phone screen, about half the screen.
> ```

On REQ-737's layout form (D-11), the maintainer chose:

> ```
> Keep form 1, fix page top (Recommended): Keep 'readable titles, Show taps for the rest'. Send the builder back once: the map must start in the top half of the first phone screen (shorter intro on phones, less empty space above the columns), and no relation label may run off the frame edge.
> ```

What this adds to the GREEN criteria, as the coordinator reads it: at 390 x 844, with the REQ-738 header (about 87 px) above the page, the map frame's top edge sits in the top half of the first screen (y about 422 or less) at load, and no relation label is cut by the frame edge after any click. Form 1 stays; the Show taps for related boxes outside the frame stay.

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
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html` (new, launcher-named) — the proof page built with the changed preset (D-01)
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c04-claude-claude-opus-5-5-relation-map-phone-first-map-reply.md` (new, launcher-named) — the editor's reply the launcher publishes beside it
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

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-16-171725/REQ-737-preflight-probe.sh) and the repository gate green at 6f00ab48 (134 s, detached checkout, stage reuse off, load 1.9). Three earlier base runs (at 7f2c4b94 and 6f00ab48) were red only on the 30 s per-file budget, each on a different file (finalization_recovery_test.go 31.15 s at load 28; report_second_opinion_test.go 30.89 s at load 11; finalization_req499_test.go 34.79 s at load 6), from other projects and the builders loading the machine.
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-16-171725)*

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

(from the builder hand-back, verbatim; D-11 to D-19 are pass 1, D-20 to D-23 pass 2, D-24 pass 3. The first-pass and second-pass pages they name were deleted on the branch in pass 3.)

- **D-11** (ESCALATE): Layout form 1 with a slide toward the left-out boxes, as described under Captures. Value: titles stay at 12.3 px after every phone click, relation words toward off-frame boxes stay readable, and the preset needed one sentence (232 bytes). Risk: with this map's column spacing, a 356 px frame at a readable scale holds one column, so most phone clicks show the clicked box plus same-column neighbours, and 70 of 87 related boxes need a Show tap. If the maintainer prefers form 2, it is a new build pass on this branch, and the preset sentence changes.
- **D-12** (DECIDE & STATE): Check 15 keeps its desktop sentence byte for byte and adds the phone sentence in the same check, so no check 21. Checks 8, 17 and 18 do not contradict it: related boxes are still shown (only some sit outside the frame), and the first view and table rows are untouched.
- **D-13** (DECIDE & STATE): The phone sentence says "every related title in the frame" and "a button that brings that box into view on the map". "Into view" covers the page scroll, because on a phone the panel sits below the map.
- **D-14** (DECIDE & STATE): The page reuses v1c03's content and map technique, rebuilt by a scratch script that strips the launcher's decoration and the page's own CSS rule for the direction block (D-07), then edits the click code. The facts in its corrections were re-checked against the current repository (line numbers unchanged).
- **D-15** (DECIDE & STATE): The measurement script is a Node Playwright script (global `playwright`, run with `NODE_PATH="$(npm root -g)"`), not the browser MCP tool, so the integrator can rerun it with one command. It therefore created nothing in the main tree's `.playwright-mcp/`; that folder does not exist, and nothing was deleted there.
- **D-16** (DECIDE & STATE): The script clicks a box with the real mouse when its title shows on screen and with dispatched pointer events otherwise (7 boxes in the phone first view lie outside the frame). The row records which.
- **D-17** (DECIDE & STATE): The panel button changed from "Open X" (select X) to "Show X" (move the map, keep the selection), because the preset now asks for a button that brings the box into view. Tapping X on the map still opens it.
- **D-18** (DECIDE & STATE): The model slug is `claude-opus-5-5`, recorded as `session-declared`. Commit trailers use Opus 5.5 per the brief's current text (an earlier copy said Fable 5.2).
- **D-19** (DECIDE & STATE): My first finalize call repeated `--backend/--model/--preset` and was refused as `TOOLBOX-USAGE`; `--finalize <staging-page>` alone is the accepted form. Nothing changed on disk from the refusal.
- **D-20** (DECIDE & STATE): The rebuilt page is a new tree page, v1c05, not a rebuild of v1c04 in place. The launcher never overwrites a page, and an in-place rebuild means deleting a published file and editing the launcher's iteration list by hand. v1c04 stays on the branch as the first-pass page. If the integrator wants only one proof page, deleting v1c04 and its `index.html` line is a two-path change; the probe passes either way.
- **D-21** (DECIDE & STATE): The two new rules go into checks 12 and 17, no check 21. The frame-top rule is about the phone first view, which check 17 owns. The arrowhead rule is a case of "never covers a box", which check 12 owns.
- **D-22** (DECIDE & STATE): The phone intro moves below the map with CSS order instead of being cut or folded, so no words leave the page and the desktop layout keeps its order.
- **D-23** (DECIDE & STATE): "Cut at the frame edge" in the addendum turned out to be the label under the arrowhead stack. Both are now measured. The edge count was already 0, and the arrowhead count went from 10 to 0 on the phone.
- **D-24** (DECIDE & STATE): The final page took id **v1c04** again, because the launcher hands out the next id not present in the bundle once v1c04 and v1c05 were deleted. Its file name differs from the first-pass v1c04 by its slug (`phone-first-map`, not `phone-readable-map`), so no link from an earlier commit's page names this file. Prepare (staging `.claude-claude-opus-5-5-staging-2753489393.html`) and `--finalize <staging-page>` alone were run, from the second-pass staging content rebuilt by the same scratch builder. Finalize: `judge_verdict pass`, `tree_id v1c04`, `model_evidence session-declared`, `banner_state inserted`.

(integrator)

- **D-25** (DECIDE & STATE): D-11's escalation is settled. The maintainer kept form 1 and sent the builder back once for the page top and the labels (`## Addendum (mid-run)`), then saw `final-390-load.png` and `final-390-click-O2.png` from pass 3 and approved integration (question dialog, 2026-11-16, per the integrate brief).
- **D-26** (DECIDE & STATE): No builder-work timing event. The build ran in three passes over about two and a half hours with maintainer waits between them, so a dispatch-to-hand-back interval would book waiting time as builder work.
- **D-27** (DECIDE & STATE): Re-passing the recorded pre-flight argv (as the brief says) rewrote `do-work/working/baseline.json`; I restored it at once with `git checkout -- do-work/working/baseline.json`. The green-gate record it wrote lives under `.git/`.
- **D-28** (DECIDE & STATE): The catalog rebuild is the merge follow-up: `ai-report-index` on the merged tree listed 50 bundles and added one row (the final v1c04 page, preset `relation-map`) under the tstack bundle, committed as `3f6ebf38`. As in REQ-727 (its D-22), qualification, the gate and the heavy lanes use the range `38b278f2..3f6ebf38` (builder merge plus catalog), and the recorded implementation hash is the builder merge `8e6d94e9`.

- **D-29** (DECIDE & STATE): Release 0.324.2, a patch: the preset changes three checks and adds no command or file. The entry claims only the phone-click limit of the 0.324.0 "Known limits" line, because REQ-738's 0.324.1 already removed the other (review F2).

- **D-30** (DECIDE & STATE): The brief asks for the integration report in this finalization's commit, but the report names the finalization commit, which does not exist until it is made. As REQ-738 did (its D-17), the report is written after finalization and committed by exact path as one plain `[REQ-737] integration report` commit. The review file and the progress log ride in the finalization commit.

<!-- D-XX counter: last used D-30. Next decision: D-31. -->

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (modified) — check 15 gains a phone sentence (readable titles win below 600 px, Show buttons reach the other related boxes); check 12 adds "wholly inside the frame" and arrowheads as taken space; check 17 adds "the map frame starts in the top half of the first screen" and "at the frame's top"; 3,057 to 3,426 bytes, still 20 checks
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html` (created) — the proof page, built through the Session form on the merged 0.324.1 launcher
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c04-claude-claude-opus-5-5-relation-map-phone-first-map-reply.md` (created) — the editor's reply beside it
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` (modified) — the launcher's iteration list gained the v1c04 line
- `ai-reports/catalog.json` and `ai-reports/index.html` (modified, integrator follow-up 3f6ebf38) — the catalog pair lists the new page under the tstack bundle

**What was done:** On a phone, a page made with the relation-map preset now keeps the clicked and in-frame related titles at about 12 px after a click (measured 12.32 px in 33 of 33 clicks, was 4.52 px on v1c03), and each related box outside the frame is one Show tap away (70 of 70 steps work). The map frame starts at y = 398 on a 390 x 844 screen under the 0.324.1 short header, and no relation word sits on an arrowhead or is cut by the frame edge. Desktop behavior of check 15 is unchanged (0 of 33 clicks leave a box out of frame). Pages created and deleted on the branch in passes 1 and 2 (the first v1c04 and v1c05 pages and their replies) are not in the cumulative diff, so Scope names only the final page.

## Discovered Tasks

(from the builder hand-back)

- The brief's finalize line ("the same launcher line with `--finalize` in place of the prepare flags") is refused by the launcher: finalize takes only the staging path — impact-negligible → report only
- With the 0.324.0 header block unstyled (D-07), the direction block was 1,892 px tall at 390, so the phone map started on the fourth screen until REQ-738 landed — impact-user-visible → report only (resolved by REQ-738's short header in 0.324.1: frame top now y = 398)
- REQ-738's short header has no side gutter at 390 px: "Iteration v1c04" and the summary triangle start at x = 0 (`final-390-load.png`). The page may not style the block (D-07), so this is a launcher change — impact-user-visible → report only (the coordinator captures it as its own new request after this integration; it is REQ-738 review F3)

## Qualification

**Gate records:** `advance --diff-range 38b278f2..3f6ebf38` (builder merge 8e6d94e9 plus the catalog follow-up, D-28) → qualify gate satisfied on the first call. Scope-drift reports only SCOPE-DECLARED-NOT-TOUCHED for `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` and `VERSION`: Scope marks all three "integrator only", and the finalization release writes them. No undeclared touch. Before the call, the two launcher-named globs in Scope were replaced with the final file names.

**Diff** (`git diff 38b278f2..3f6ebf38 --stat`): 6 files: the preset (3 lines changed), the final v1c04 page (692+), its reply (20+), the bundle's `index.html` (1 line), `ai-reports/catalog.json` (+7, one v1c04 row with preset `relation-map`) and `ai-reports/index.html` (1 line). Nothing from REQ-738 is in the range, because the branch merged main at ab5e07c1 and that commit is already in 38b278f2.

**Requirement trace (against the diff and the files at 3f6ebf38):**
- Readable titles after a phone click: check 15's new sentence says below 600 px the clicked title and every related title in the frame stay at about 12 px or more. Measured on the published page (`measure-output-final.json`): 12.32 px smallest clicked and in-frame related title in 33 of 33 clicks at 390 x 844 (v1c03: 4.52 px).
- Related boxes not shown reachable in one named step: check 15 names the panel's button that brings the box into view; 70 of 70 Show steps work, each title at 12.32 px after the tap.
- Desktop behavior of check 15 unchanged: the first sentence of check 15 is byte-identical; 0 of 33 desktop clicks leave a clicked or related box out of frame (as REQ-727 measured).
- Other checks hold: 20 numbered checks in order, no `{{`; check 17 (no sideways scroll 0 at both widths, first view titles 12.93 px), check 18 (28 = 28 rows, 0 too wide), check 19 (no-JavaScript view: tools hidden, 28 of 28 words on their lines). Checks 12 and 17 gained the maintainer's two addendum rules (D-21).
- Mid-run addendum: map frame top y = 398 at load at 390 x 844 under the real 0.324.1 short header (top half: about 422 or less); 0 words on arrowheads, 0 cut at the frame edge, 0 on boxes, 0 broken, at both widths; form 1 kept.
- Byte bound: 3,426 bytes (target about 3,500, probe bound 3,600); no check removed or merged.
- Only the preset changes among shipped files: `git diff --name-only 38b278f2..3f6ebf38 -- skills/ _dev/ suite/ do-work/` lists only the preset. Loader, prompt, other presets, launcher Go code and report kit untouched; v1c03 and every other bundle file untouched apart from the iteration list line.
- Proof page and judges: launcher finalize `judge_verdict pass`, standalone judge `success` with findings `[]` (hand-back pass 3).
- Maintainer saw the phone captures and approved integration (D-25).

**Scope comparison:** declared 9 paths, touched 6, all declared; the 3 untouched are finalization-owned. No debug output, staging file, capture or queue file in the range.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at 3f6ebf383b7290efefbaf7ea00c31bb127933d17 (builder merge 8e6d94e9 plus the catalog follow-up, D-28) in the run's detached checkout, stage reuse off (`helpers/gate.sh`), started 2026-11-16T20:02:41Z at 1-minute load 2.82 with no other gate running; then `advance REQ-737 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-11-16-171725/REQ-737-probe.sh`.
**Result:** ✓ Gate passed on the first run (wall 164 s; do-work-cli 987 Go tests, slowest file `finalization_recovery_test.go` 21.87 s; queue-kanban 427 tests, slowest file 23.37 s; both under the 30 s budget). GREEN probe passed under advance (`run-blocked-check` and `green-gate` satisfied): `REQ-737 probe PASS (3426 bytes, 20 checks, …index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html)`. Builder: `TestSecondOpinionShippedPresetsLoadFromTheToolbox` PASS, `shipped-package-reference-contract.sh` and `contract-regressions.sh` exit 0 in every pass.

**Red-green validation (tdd: false; D-10, the browser measurement is the proof):**
- `REQ-737-captures/measure-relation-map.js` at 390 x 844, 33 clicks: ✗ on v1c03 (0.324.0 preset): smallest clicked and related title 4.52 px, scale 0.205 to 0.342 → ✓ on the final v1c04: 12.32 px in 33 of 33 clicks, 70 of 70 Show steps work, frame top y = 398 at load, 0 words on arrowheads or cut at the frame edge, `summary.green` true at both widths (`measure-output-v1c03.json`, `measure-output-final.json`).
- GREEN probe `REQ-737-probe.sh`: ✗ at base (step 2: no numbered check names a click and a phone) → ✗ after the preset commit (step 4: no relation-map page other than v1c03, as designed) → ✓ at 9375501d and again under advance at 3f6ebf38.

**New tests added:** none (prose preset; D-10). The measurement script and the GREEN probe are run artifacts, not shipped tests.

**Existing tests updated (cross-REQ impact):** none. `TestSecondOpinionShippedPresetsLoadFromTheToolbox` passes unchanged.

**Heavy verification plan:**
- Range: 38b278f2fd28c454f146a2de8cfec4977b4c4b45..3f6ebf383b7290efefbaf7ea00c31bb127933d17
- Six lanes, each selected because coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles): queue-kanban-javascript, queue-kanban-browser, do-work-cli-integrations, staged-skills, updater, installer; each `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane>`

*Verified by work action (integrator, run work-2026-11-16-171725)*

## Review

**Overall: 95%** | 2026-11-16T20:12:13Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: After a phone click or Show tap, the framing centers vertically and leaves empty space above the column title in the 620 px frame: 88 px (O2), 155 px (O1), 223 px (O2 then Show S1). Only the first view pins to the top (`index-v1c04-…-phone-first-map.html:583`; `frameSelection` :519, `bringIntoView` :564 use `centerOn`). The first screen meets the addendum (12 px). — impact-user-visible → report only
- F2: The hand-back's proposed CHANGELOG entry is titled 0.324.1 (already REQ-738's) and claims both 0.324.0 limits. The release should claim only the phone-click limit. — impact-negligible → report only (handled in the release entry text, D-29)
- F3 (REQ-738's, not scored): the short header has no side gutter, x = 0 at 390 and 1440 (REQ-738 review F3). The coordinator captures it as a new request. — impact-user-visible → report only
- F4 (Nit): check 15's "readable titles win" leaves the overridden rule implicit. It is the sentence just before it. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages. My rerun of `measure-relation-map.js` at `3f6ebf38` matches the hand-back exactly: frame top y = 398 (phone) and 385 (desktop); 12.32 px titles in 33 of 33 phone clicks; 70 of 70 Show steps; 0 words on arrowheads or cut at the frame (also after drags and zooms); header fold closed at load and opens; green at both widths; shipped-presets Go test ok. Anti-bloat: 0 extras.
**Restatement sweep:** redefined relation-map checks 12, 15 and 17 (phone click rule, words off arrowheads and inside the frame, phone frame in the top half). No stale shipped restatement in `ai-report-guide.md:124`, `ai-report.md`, CLI help or the lessons files. The released 0.324.0 "Known limits" phone-click line and the 0.324.1 "y = 830" figure are noted only. Inherited from REQ-738: the header form; the v1c04 page agrees.
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: 38b278f2fd28c454f146a2de8cfec4977b4c4b45
- Target: 3f6ebf383b7290efefbaf7ea00c31bb127933d17 (builder merge 8e6d94e9 plus the catalog follow-up, D-28)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles)
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for the four `ai-reports/` paths (no lane covers report bundles)

## Heavy Verification Result

- Target: 3f6ebf383b7290efefbaf7ea00c31bb127933d17; execution revision 3f6ebf383b7290efefbaf7ea00c31bb127933d17 (the run's detached checkout, browser variable set, no HEAVY-RUN-LANE-SKIPPED, started 2026-11-16T20:06:06Z at 1-minute load 5.82 while the reviewer ran; every lane green on the first run)
- queue-kanban-javascript: executed, exit 0, 10 s
- queue-kanban-browser: executed, exit 0, 93 s
- do-work-cli-integrations: executed, exit 0, 79 s
- staged-skills: executed, exit 0, 42 s
- updater: executed, exit 0, 76 s
- installer: executed, exit 0, 31 s

## Lessons Learned

**What worked:** Building the final proof page after merging main into the builder branch (pass 3), so the frame top was measured under the real 0.324.1 header (y = 398) instead of derived from a page with no header (400). The kept measurement script let the reviewer rerun the same 33-click matrix and get identical numbers.
**What didn't:** The first phone framing at a readable scale showed only the clicked box, because a 356 px frame at 0.56 holds one column of this map; the builder found it from a capture, not from the size numbers. A label read as "cut by the frame edge" in pass 1 was really under the arrowhead stack (D-23).
**Worth knowing:** After a phone click or Show tap, the page centers the selection vertically, so 88 to 223 px of empty space can sit above the column title in the frame (review F1). The preset's top-of-frame rule covers only the first view.

## Orientation

The relation-map preset (`skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md`) is prose the `ai-report second-opinion` launcher hands to the editor model; its numbered checks are the only contract a generated map page follows, and checks 12, 15 and 17 now carry the phone rules.

## Timing

Observed 2026-11-16T20:00:30Z to 2026-11-16T20:12:14Z: 11m 44s total, 15m 53s attributed across 4 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 9m 06s | 2 |
| review | 6m 04s | 1 |
| handback-merge | 43s | 1 |

Slowest stage: verification-gate / heavy drain, 6m 08s, outcome success.

The heavy drain ran beside the review (both started at about 20:06), so the attributed time is larger than the observed span.

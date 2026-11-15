---
id: REQ-727
title: 'Ship a relation-map preset for ai-report second-opinion'
status: claimed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
depends_on: [REQ-733, REQ-734]
related: [REQ-733, REQ-726, REQ-728, REQ-729, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-*.html", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c03-*-reply.md", "ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html", "ai-reports/catalog.json", "ai-reports/index.html", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION"]
status_changed_at: 2026-11-15T10:59:37Z
claimed_at: 2026-11-15T11:01:29Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-11-15T11:05:01Z
  basis:
    - Route B
    - 11-file write set
    - 3 subsystems involved
    - 9 acceptance criteria
dispatch_at: 2026-11-15T11:10:24Z
builder_handback_at: 2026-11-15T11:35:32Z
---
# Ship a Relation-Map Preset for ai-report second-opinion
## What
Add one shipped second-opinion preset, `relation-map`, so `ai-report second-opinion --preset relation-map <bundle>` writes a page that shows a report as one zoomable map of its items and the relations between them.
## Why
The maintainer wants report content that is easy to understand by zooming in and out and by seeing relations. One example page was built by hand and passed the render judge.
## Context
- The example is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`, built with `--editor-direction`.
- A preset is one Markdown file with a direction to the editor. The loader reads `<bundle>/second-opinion-presets/<name>.md`, then the shipped folder `skills/do-work-toolbox/docs/second-opinion-presets/`. No loader change is needed.
- `mindmap.md` in that folder is the nearest existing preset. A mind map is a tree from one center; a relation map shows relations between any two items.
## Detailed Requirements
- Read `mindmap.md` first, and say in the new preset how the two differ.
- Boxes are the report's own items (for example options, risks, findings, limits, files, requests). Lines are relations the report states. The editor says on the page that the lines are its reading of the report.
- At most four kinds of line, each with one color and a plain name, and a key that explains them.
- Semantic zoom: zoomed out, a box shows its code and title; zoomed in, a summary; a click opens the detail with its sources and its relations, and fades unrelated boxes.
- Drag to move, pinch or Control-scroll to zoom, buttons and keys for the same, and a button that shows everything. Plain scrolling must still scroll the page.
- The page works without JavaScript: the map shows at full size and scrolls. Every box is also given as text, and every relation as one table row.
- No library. Inline SVG and plain JavaScript only, as the second-opinion prompt already requires.
- About 30 boxes is the upper limit for one map. With more, group items and say so.
- On a narrow screen the map starts on the main item, not on the whole picture.
- Every fact still comes from the report or the repository, as in `bold.md`.
- Files: `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (new); the preset list in "A Second Opinion" in `skills/do-work-toolbox/docs/ai-report-guide.md`; `TestSecondOpinionShippedPresetsLoadFromTheToolbox` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`, if it names the shipped presets; `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
## Constraints
- The preset loader, the second-opinion prompt, the other presets and every existing `ai-reports/` bundle stay unchanged.
- The report kit (`skeleton.html`, `presentation.css`, `presentation.js`) stays unchanged.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** `ai-report second-opinion --prepare --preset relation-map` on a fixture bundle, through the existing shipped-presets test.
**Why RED now:** The call is refused because no preset named `relation-map` ships.
**GREEN when:** The same call returns state `prepared` with `preset_file` in the shipped folder; the guide names the preset; one real page is built with the preset on a bundle other than the example, passes `ai-report-judge`, and is shown to the maintainer before the request is called done.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Open Questions
- [x] Has the maintainer looked at the example map page, and should it become a standard option? → Not looked yet, hold it.
  Note (2026-11-14): The request is captured as blocked and waits for the maintainer. The maintainer was told the risk of a yes: the map was readable for 28 boxes, and larger reports are untested. `do-work clarify` releases it when the maintainer confirms.
## Builder Guidance
- Certainty: the direction is the maintainer's; the page itself has not been approved. Do not build before the block is released.
- Scope cue from the maintainer on 2026-11-14: keep language clear, do not complicate it.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Addendum (2026-11-14)

Captured with UR-177 (rebuild the relation map page). The maintainer saw the first map page and asked whether it is the best way to present the content. It is not, so REQ-733 builds an improved page first.

- This request now depends on REQ-733.
- The example the maintainer looks at before releasing the block is the page REQ-733 builds, not the first map page.
- The preset asks the editor for the improved design: a left-to-right layout, a first view of about nine main boxes with the rest opening on click, line labels that stay readable at every zoom level, and only light dimming of unrelated boxes. Where this differs from the Detailed Requirements above, this addendum wins.

## Addendum 2 (2026-11-14)

Captured with UR-178 (build the proposal as a table and compare it with the map). The maintainer wants to compare the map with a supports-and-limits table before choosing a form.

- This request now also depends on REQ-734.
- The maintainer picks map or table after seeing both. If the table wins, this request is not built as written: the maintainer cancels it or captures a new one for the table form.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (bare; the row is `slugged: partial`, so the targeted `#session-is-the-backend` form is not legal) — cost 24,500 tokens against the 2,000-token budget — matches "a session that is its own second-opinion backend": the builder runs a Session second opinion for the real page. The brief names `lessons-do-work-cli.md:17` as optional reading.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-relation-map-preset.md*

## Blocked

- [2026-11-15] blocked on "Not looked yet, hold it: the maintainer looks at the example map page and confirms that it should become a standard option" — cleared by user via clarify

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and the files are named, and the loader needs no change. The preset text is the real work: it must turn the defects that took REQ-733 (rebuild the relation map page) three review rounds into instructions any report can follow, so the example page and its reviews need exploration first. One real page must then be built with the preset and pass the render judge.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Required-lessons consult:** read `_dev/primes/lessons-releases.md` in full (family `canonical-link-outlives-its-target` applies: this REQ closes UR-176, see D-12). `do-work/lessons-index.md` row `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` matches "a session that is its own second-opinion backend", but that satellite is `slugged: partial` at 24,500 tokens, so it is dropped for budget (section below); the brief names its one bullet `session-is-the-backend` (`lessons-do-work-cli.md:17`) as optional reading. No other row matches: the kanban satellite covers the board, not report pages.

**Where the change goes:**
- Shipped presets: `skills/do-work-toolbox/docs/second-opinion-presets/*.md`, twelve files of one paragraph each (273 to 467 bytes). `mindmap.md` is a tree from one center with click-to-reveal nodes; `bold.md` carries the "every fact still comes from the report or the repository" line the new preset repeats.
- Loader: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:195-229`. The bundle folder wins, then the folder beside the prompt; the text is trimmed and refused when empty or when it contains `{{` (`:225`). No change needed.
- The whole direction is printed at the top of every directed page: a meta tag and a `white-space:pre-line` paragraph in the launcher header block (`report_second_opinion.go:1626` and `:1634`). Line breaks in the preset therefore show on the page, and a long preset pushes the map down (D-02).
- The prompt: `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md:19-20`. A direction overrides ROLE, GOAL, every RULE and PROCESS, so Rule 4's "inline SVG or Mermaid" (`:26`) no longer holds once a preset is given. The preset must say "no library" itself.
- Test: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go:1007-1037`, `TestSecondOpinionShippedPresetsLoadFromTheToolbox`. It reads the real toolbox folder, requires a named list at `:1026` (the ten learning-method presets), then prepares every shipped file and checks `state`, `preset` and `preset_file`. A name added to the `:1026` list is the RED: "does not ship".
- Guide: `skills/do-work-toolbox/docs/ai-report-guide.md:124`, the "A Second Opinion" paragraph lists the shipped presets by group.
- Not touched by this REQ but naming presets: the report kit's closing band (`skeleton.html`, REQ-712) suggests a teaching method; the kit stays unchanged by constraint (D-07).

**The example and what its reviews found** (`ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html`, 97 KB, reviews `do-work/runs/work-2026-11-15-012821/REQ-733-review.md`, `-review2.md`, `-review3.md`):
- Labels are an unscaled HTML overlay, 110 to 190 px wide (`:52`, `:311-313`); shown boxes are taken space and dropped labels are counted, "N of M relation words in view" (`:428`).
- A click zooms in to `LABEL_SCALE` 0.575 only when the box and its related boxes all fit (`:461-468`); the fit test is written twice (H2) and a click after a manual zoom-out zooms in unevenly (H1).
- The view never zooms itself below `READABLE` 0.42 (`:441-449`); below 600 px the first view frames the decision and options (`:457`).
- `focusin` resets the frame's scroll to 0 before measuring the focused box (`:559-564`); the scroll-handler pin alone did not work (REQ-733 Lessons Learned).
- No-JS: `.map-tools[hidden]{display:none}` (`:20`) hides dead buttons; the line words need script, so the no-JS map has none (F5, open).
- Dimming: unrelated boxes at 0.72 opacity, unrelated lines at 0.22 (`:49-50`).
- REQ-734's comparison: with no click the map shows 0 of 37 relations on the first screen at 1440 and at 390; at 390 the map is a 1910 px drawing in a 331 px frame that needs panning, and the relation table needs a sideways swipe.

**Bundle for the real page:** `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1` (current revision, about 15 coded items: options O1-O3, seeds S1-S6, questions Q1-Q3, limits; two tree pages v1c01 and v1c02 already), so the new page is v1c03 there (D-01). The Session sequence is `skills/do-work-toolbox/actions/ai-report.md:81` onward; a re-judge after a fix uses `ai-report-judge --page` (`ai-report.md:190`, shipped by REQ-728).

*Generated by the pre-dispatch agent (exploration done in-session, no Explore subagent)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (new) — the preset: an opening paragraph that says how it differs from the mind map, then the numbered checks of D-03 and D-04
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — one clause in the "A Second Opinion" preset list names relation-map
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — relation-map joins the named list in TestSecondOpinionShippedPresetsLoadFromTheToolbox (D-06)
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-*.html` (new, launcher-named) — the real page built with the preset (D-01)
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c03-*-reply.md` (new, launcher-named) — the editor's reply the launcher publishes beside it
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` (modify, launcher only) — the iteration list the launcher rebuilds
- `ai-reports/catalog.json` (modify, integrator only) — catalog rebuild after the merge (D-10)
- `ai-reports/index.html` (modify, integrator only) — catalog rebuild after the merge (D-10)
- `CHANGELOG.md` (modify, integrator only) — release entry
- `skills/do-work/CHANGELOG.md` (modify, integrator only) — byte-identical mirror
- `VERSION` (modify, integrator only) — minor bump (D-09)

**Files I will NOT touch:** the preset loader and every other line of `report_second_opinion.go`, `ai-report-second-opinion-prompt.md`, the other twelve presets, the report kit (`skeleton.html`, `presentation.css`, `presentation.js`), `skills/do-work-toolbox/actions/ai-report.md`, `skills/do-work-toolbox/actions/help.md`, every existing file in any `ai-reports/` bundle except the launcher's iteration list, any `do-work/` path.

**Acceptance criteria (restated from REQ):**
- [ ] `ai-report second-opinion --prepare --preset relation-map` on a fixture bundle returns state `prepared` with `preset_file` in the shipped folder, through the shipped-presets test, which was RED first.
- [ ] The guide's "A Second Opinion" preset list names `relation-map`.
- [ ] The preset says how a relation map differs from the mind map: a tree from one center versus relations between any two items.
- [ ] The preset asks for every item of the Detailed Requirements as changed by Addendum 1. Boxes are the report's own items, and lines are relations the report states, marked on the page as the editor's reading. At most four line kinds, each with one color, a plain name and a key. Semantic zoom. Drag, pinch or Control-scroll, buttons, keys and "Show everything", while plain scrolling still scrolls the page. A no-JS page with every box as text and every relation as a table row. No library. About 30 boxes at most, else group and say so. A narrow screen starts on the main item. A left-to-right layout, a first view of about nine boxes, labels readable at every zoom, and light dimming. Every fact comes from the report or the repository.
- [ ] The preset turns each of the seven open defects of REQ-733 and REQ-734 into a checkable instruction (D-04).
- [ ] A decision records which of them a preset cannot guarantee (D-05).
- [ ] One real page is built with the preset on a bundle other than the example, and the finalize judge reports `pass`. The hand-back opens with its path so the maintainer can look at it (D-08), and it records a browser measurement of the seven checks at 1440x1000 and 390x844.
- [ ] The loader, the prompt, the other presets, the report kit and every existing bundle file stay unchanged, except the launcher's iteration list.
- [ ] The preset tests pass and gofmt is clean on the test file.

## Decisions

- **D-01** (DECIDE & STATE): The real page for the GREEN proof is built on `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1` with its `index.html` as the base, so it becomes tree page v1c03 of that bundle. The builder runs the Session form (`--prepare`, then `--finalize`) in its own worktree with `--backend claude --model <the builder's own declared model slug>` and `--preset relation-map` alone, without `--editor-direction`, so the page tests the preset text and nothing else. If that report holds fewer than eight items with stated relations, the builder picks another current bundle (never the Tinka example) and records why.
- **D-02** (DECIDE & STATE): The preset is one short opening paragraph, then numbered checks of one sentence each, every check stated as something a reader can see or measure. Target: about 3,500 characters or less. The launcher prints the whole direction at the top of the page (`report_second_opinion.go:1634`), so every extra line pushes the map down. The builder reports the preset's byte count and the height of the direction block on the real page at 1440 and 390.
- **D-03** (DECIDE & STATE): The layout rule for any report, which REQ-733 hard-coded (its reply Q4): columns run left to right in the order of the report's argument. First what is known (evidence, findings, sources), then what is proposed or decided (the subject, the decision, the options), then what tests or limits it (risks, objections, open questions, limits). The first view shows about nine main boxes: the subject or decision, the items the report ranks or recommends, and what tests them most directly. A click reveals the boxes related to the clicked box, and "Show everything" shows all of them.
- **D-04** (DECIDE & STATE): The preset turns each open defect into a check. (1) Relation words visible only after a click: with no click, the first view shows the words of the relations between its own boxes. (2) Labels covering box titles: a relation word never covers a box, because shown boxes count as taken space, and a word with no free place is left out and counted ("N of M relation words in view") while the detail panel still lists every relation. (3) Words broken inside narrow labels: a label is never narrower than its longest word, and no word in a label is broken or hyphenated. (4) A forced zoom that pushes related boxes out of view: a click never leaves the clicked box or a related shown box outside the frame, and it zooms in only as far as all of them still fit. (5) Tab focus on off-screen boxes: on focus, reset the map frame's own scroll to 0, then move the map so the focused box is fully in view. (6) A phone view that needs panning: below 600 px the first view frames the main item, its titles can be read without zoom or sideways panning, the page never scrolls sideways, and each relation-table row stacks so it needs no sideways swipe. (7) No relation words without JavaScript: the no-JS map draws each relation's words as SVG text on its line and hides the tool buttons (a `[hidden]` rule that beats the buttons' display rule). The checks from the REQ text also go in: light dimming (unrelated titles stay readable), a deselect that restores the panel, and every control and fallback the Detailed Requirements and Addendum name.
- **D-05** (ESCALATE): What a preset cannot guarantee. The judge (`ai-report-judge`) checks page width, themes, overflow and links. It does not check label overlap, broken words, zoom framing, focus, or the no-JS line words, so checks (2) to (5) and (7) are instructions that an editor can still miss. The preset cannot shorten the launcher's direction block, which sits above the map. It cannot make an editor without a browser test its page. It cannot make a map of more than about five boxes explorable on a phone without dragging: (6) covers only the first view and the text forms. The lines are the editor's reading of the report, and nothing checks them against it. Value: the preset states its limits honestly, and the builder's browser measurement of the real page shows the checks work at least once. Risk: a generated map can still carry a defect the judge cannot see. This is reversible: editing the preset is one text file, and the maintainer looks at the real page (D-08).
- **D-06** (DECIDE & STATE): The test change adds `"relation-map"` to the named list in `TestSecondOpinionShippedPresetsLoadFromTheToolbox` (`report_second_opinion_test.go:1026`). The comment and the error text change from "learning-method preset" to "named preset", because a relation map is a form, not a teaching method. The RED is the test failing with "named preset \"relation-map\" does not ship". The GREEN is the same test passing, with the prepare row in state `prepared` and `preset_file` in the shipped folder. No new test function is added.
- **D-07** (DECIDE & STATE): Docs change only in the "A Second Opinion" paragraph of `ai-report-guide.md:124`: one clause names `relation-map`. The report kit's closing band (`skeleton.html`), `ai-report.md`, `help.md` and the prompt stay unchanged by the REQ's constraints. So the closing band does not suggest the new preset, and the builder records that as a report-only discovered task.
- **D-08** (ESCALATE): The REQ says the real page "is shown to the maintainer before the request is called done". This run is unattended. The hand-back opens with the page path, and the coordinator puts that path first in its run summary. The REQ completes on the finalize judge `pass` and the review. Value: the run does not stall on a person. Risk: the maintainer may dislike the page after completion. The fix is then an edit to the preset or a follow-up request, both reversible.
- **D-09** (DECIDE & STATE): A new shipped preset is a new user-visible option, so the integrator makes a minor release (0.323.0 to 0.324.0 at today's HEAD). The builder only proposes the changelog entry.
- **D-10** (DECIDE & STATE): REQ-730 (the catalog lists second-opinion pages) makes `ai-reports/catalog.json` and `ai-reports/index.html` list tree pages. After the merge, the integrator rebuilds the catalog pair with `ai-report-index`, as REQ-734 did (its D-16). The builder does not touch the catalog.
- **D-11** (DECIDE & STATE): Scope names the launcher-named page and reply as globs, because their final names are known only after finalize. The integrator replaces the globs with the resolved names, as in REQ-734 (its D-17).
- **D-12** (DECIDE & STATE): This REQ is the last open member of UR-176. Its six other members are archived flat: REQ-726 and REQ-728 to REQ-732. Finalizing this REQ closes UR-176 and moves those records into `do-work/archive/UR-176/`. A sweep on 2026-11-15 found no shipped file that links those flat paths. The only matches were fixture strings in two tests (`finalization_recovery_test.go:509` and `recovery_set_aside_test.go:113`), and those are not links. The integrator runs the sweep again at finalize (`_dev/primes/lessons-releases.md`, family `canonical-link-outlives-its-target`).

<!-- D-XX counter: last used D-12. Next decision: D-13. -->

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-727-preflight-probe.sh) and the repository gate green at 445950c4 (130 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*

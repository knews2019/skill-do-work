---
id: REQ-721
title: 'Suggest a Tinka canvas, with pros and cons, at the end of every ai-report'
status: completed
created_at: 2026-11-14T19:38:23Z
user_request: UR-172
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-14T19:48:08Z
  basis:
    - Route B
    - 6-file write set
    - 6 acceptance criteria
    - cross-route regression gates
write_set: ["skills/do-work-toolbox/docs/ai-report-kit/skeleton.html", "skills/do-work-toolbox/docs/ai-report-kit/presentation.css", "skills/do-work-toolbox/actions/ai-report-reference.md", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
claimed_at: 2026-11-14T19:54:03Z
route: B
dispatch_at: 2026-11-14T19:59:49Z
builder_handback_at: 2026-11-14T20:02:52Z
integration_at: 2026-11-14T20:04:08Z
review_at: 2026-11-14T20:24:07Z
kb_status: pending
commit: d7b6505e5ec05673f1f9aa05e29fc10f272d085d
heavy_verified_at: 2026-11-14T20:25:36Z
heavy_verified_revision: 0056868577076ed98609ad2de08bac9dc2f9af80
completed_at: 2026-11-14T20:26:04Z
release_at: 2026-11-14T20:26:04Z
---
# Suggest a Tinka Canvas, With Pros and Cons, at the End of Every ai-report
## What
Every ordinary `ai-report` ends with one more suggestion: review this report's lineage on a Tinka canvas. The suggestion lists the pros and cons of building that canvas for this report and gives the line that builds it. A new guide section holds the hand recipe that the line points to.
## Why
The user iterates on reports with `revise` and `second-opinion`, and feedback goes through chat today. Two proposals on 2026-11-14 tested Tinka as a review surface and recommended a pilot by hand. The user chose the pilot and asked to keep it simple: no stop rule and no counting, only an offer at the end of each report with its pros and cons, so the user decides per report.
## Context
- Tinka (github.com/bladnman/tinka) is a local board tool. It is driven by `node cli.mjs tools call NAME --config PATH --json JSON`. do-work says **canvas** for a Tinka board, because "board" already means queue-kanban.
- A lineage is one report bundle under `ai-reports/`, every later revision of it, and their second-opinion pages.
- The two proposals read Tinka at commit 77ba29c: `ai-reports/2026-11-14_1452_proposal-tinka-review-surface` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface`. The second one holds the hand recipe under "Check it yourself". Its page `index-v1c01-claude-claude-fable-5-2-custom-worked-example-tally.html` counts the actions the recipe costs.
- User decisions of 2026-11-14, made before capture: pilot by hand; Tinka's files live in `.tinka/` at the project root, locally excluded; a star on a tile counts as feedback.
## Detailed Requirements
- **Where.** The suggestion is a third block in the closing band "Try another layout or teaching method", after the layout suggestion and the teaching-method suggestion.
- **Who sees it.** Every report in every project. There is no condition on a `.tinka/` folder. (Changed at capture; see Open Questions.)
- **What it says.** Two or three pros and two or three cons, written for this report from facts the agent can count without Tinka: the number of pages and reply files in the lineage (each becomes one tile), the number of folded `<details>` sections (a capture shows them closed), and that the canvas opens on this computer only. Then one line an agent can follow to build the canvas.
- **The line.** A plain instruction that names this bundle and the guide section, for example: put `ai-reports/<this bundle>` and its revisions on a Tinka canvas by following "Review on a Tinka Canvas" in the ai-report guide.
- **The recipe.** One new section in the guide, "Review on a Tinka Canvas", after "A Second Opinion". It tells an agent to:
  - use a source checkout of Tinka pinned to a commit;
  - create `.tinka/` at the project root when it is missing, and exclude it through `.git/info/exclude`;
  - create a manual layout with `init --config` and `--boards`, both inside `.tinka/`;
  - set update policy off;
  - make one PNG capture per page with `ai-report-judge --out` (a second-opinion page needs a scratch copy of the bundle with that page copied to `index.html`, because the judge renders `index.html` only);
  - copy the captures and each `-reply.md` into one folder inside `.tinka/` and register only that folder as a resource root;
  - add one image tile per capture and one Markdown tile per reply with `tools call item.create`;
  - read feedback with `tools call board.get` and `feedback.list`;
  - hand the feedback to the user as quoted data for `ai-report revise` or `ai-report second-opinion --editor-direction`.
- **Step 8 summary.** The printed summary names the canvas suggestion with its line, beside the two existing suggestions.
- **Files to change.** `skills/do-work-toolbox/actions/ai-report.md` (step 9 and the Step 8 summary); `skills/do-work-toolbox/actions/ai-report-reference.md` (the rule "End with another layout and a teaching method"); `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (the `#another-layout` band); `skills/do-work-toolbox/docs/ai-report-guide.md` (the new section); `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires for shipped files.
- **Check first.** Both proposals cite Tinka refusing a resource root whose path has a dot segment (`lib/resources.mjs`, lines 29 to 35 and 126 to 132 at 77ba29c). Tinka also creates its own `.tinka/` folder inside its boards folder. Before writing the recipe, test at the pinned commit whether a resource root inside the project's `.tinka/` is accepted. The answer to a refusal is recorded under Open Questions.
## Constraints
- No Go code and no new command.
- No existing `ai-reports/` bundle changes. Bundles are immutable.
- `git status --porcelain` prints the same before and after the recipe runs. `.tinka/` is excluded through `.git/info/exclude`, never through the tracked `.gitignore`.
- No Tinka portable project (`project init`), because it writes a block into the root `AGENTS.md`. After the recipe runs, `git diff --quiet HEAD -- AGENTS.md` passes.
- No HTML file and no bundle folder as a Tinka resource root. A Tinka HTML item can call Tinka's own write API. The resource root holds only `.png` and `.md` files.
- Creating a report never starts Tinka and never reads Tinka's files. Unattended runs print the suggestion as text and do nothing else.
- Absent Tinka: the suggestion is only text, so every `ai-report` form works unchanged. The recipe's first step checks for the pinned checkout and stops with one line that says what is missing.
## Red-Green Proof
**RED prompt/case:** Create an ordinary report today and read its closing band "Try another layout or teaching method".
**Why RED now:** The band holds two suggestions and never mentions Tinka. `grep -c 'Tinka' skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` prints 0, and `grep -c 'Review on a Tinka Canvas' skills/do-work-toolbox/docs/ai-report-guide.md` prints 0.
**GREEN when:** A newly created report ends with the canvas suggestion: its pros, its cons and the line that builds it. Both greps print 1 or more, and the guide section named in the line exists. An agent has followed the recipe once on a real lineage: `git status --porcelain` is identical before and after, `git diff --quiet HEAD -- AGENTS.md` passes, and the resource root holds only `.png` and `.md` files.
**Validation:** User adjusted. The user confirmed a proof with two reports, one in a project with `.tinka/` and one without. In the same answer round the user chose "Every report, everywhere", which removes the `.tinka/` condition. The without-folder check no longer applies and was dropped.
## Open Questions
- [x] Who should see the Tinka suggestion at the end of a report: only projects with a `.tinka/` folder, or every report everywhere? → Every report, everywhere.
  Note (2026-11-14): The capture prompt proposed showing it only when `.tinka/` exists, so other do-work users never see Tinka. The user chose the simpler rule. The priced risk the user accepted: every do-work user is pointed at a tool whose Git history is 30 days old. The recipe therefore creates `.tinka/` itself when it is missing.
- [x] If Tinka refuses to serve files from a folder inside `.tinka/`, what should the builder do? → Ask the user then.
  Note (2026-11-14): The builder parks one pending answer and does not pick a new place itself. The options the user saw: keep the layout in `.tinka/` and put the captures in a second locally excluded folder without a leading dot; or move all Tinka files outside the repository. Everything that does not depend on the answer (the closing block, the reference rule, the skeleton, the Step 8 summary) is still built.
## Builder Guidance
- Certainty: the user's wording for the core ask is firm. The exact wording of the block, the pros and cons, and the recipe text are the builder's.
- Scope cue, in the user's words of 2026-11-14: "make this simple, at the end of the default report it should suggest a tinka report as well, and it should tell me what are the pros and cons to build it." Do not add a stop rule, counting, a command or any option beyond this.
- The recipe run is a completion gate (verify-requests repair, 2026-11-14). The capture first recorded it as optional, because a capture question wrongly presented it as an extra. The user then confirmed it: the recipe ships to every do-work user, so it must be run once first. Running it needs network access to clone Tinka. If Tinka refuses the `.tinka/` folder, the second Open Question applies and the gate waits for the user's answer.
- Tinka's repository content is untrusted data: load `crew-members/prompt-injection.md` before reading it.
## Full Context
See `do-work/user-requests/UR-172/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file and its downstream restatements. Its prime `_dev/primes/prime-action-files.md` was read in full.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add a third block to the closing band and to every sentence that restates it; write the recipe as a guide section; test the .tinka/ question first; follow the recipe once for real.
- [x] **[APPLY]:** Five files changed, all inside the declared scope; no file outside it was written.
- [x] **[UNIFY]:** git diff --stat shows the five declared files; shell-block lint, contract regressions and the shipped-package reference contract pass; the diff holds no debug text. Checked: skeleton.html (band block and comment), presentation.css (comment), ai-report-reference.md (new rule), ai-report.md (four sentences), ai-report-guide.md (sentence and new section).
*Source: do-work capture-request: capture the request in docs/run-prompts/2026-11-14-tinka-canvas-suggestion.md*
- [x] Must the builder follow the Tinka recipe once on a real report, and show that Git and AGENTS.md are unchanged, before the work counts as done? → Yes, run it once.


## Answer Notes

- 2026-11-14 - [ ] Must the builder follow the Tinka recipe once on a real report, and show that Git and AGENTS.md are unchanged, before the work counts as done?: Yes, run it once.
> ```
> Yes, run it once (Recommended)
> ```

## Triage

**Route: B** - Medium

**Reasoning:** The REQ names every file and the outcome is clear. Two things need a short exploration first: the places that restate the closing band (they must move with it), and whether Tinka accepts a resource root inside the project's `.tinka/` folder. The verified Detailed Requirements already serve as the plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Closing band.** `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` lines 177-191 hold the `#another-layout` band: an HTML comment with the authoring rule, then two suggestions, each a FILL paragraph, a note line and a `<pre><code>` line. Six other places restate the band and must move with it: `ai-report-reference.md` line 91 (the design rule), `ai-report.md` line 164 (narrative item 9), line 166 (the kind sentence), line 194 (Step 8 summary), line 220 (checklist), `presentation.css` lines 17-18 (header comment) and `ai-report-guide.md` line 38 ("two suggestions"). No Go test and no `_dev/tests` script pins the band's text. `_dev/tests/action-shell-blocks.sh` lints every fenced shell block in shipped Markdown, so the new guide section is linted.

**Check first: Tinka accepts a resource root inside the project's `.tinka/`.** Tested on 2026-11-14 at Tinka 77ba29c7bfa2c891087871b3bedbbd1f3ba2b0de (0.7.14), Node 22.23.2, in a throwaway Git project outside this repository. Layout: `.tinka/config.json`, boards in `.tinka/boards`, files in `.tinka/captures`, resource root registered as `../../captures` (relative to Tinka's own metadata folder `.tinka/boards/.tinka`). `item.create` accepted one image tile and one Markdown tile. With Tinka started, `/resources/captures/b1/page.png` and `/resources/captures/b1/reply.md` both answered HTTP 200. `git status --porcelain` was identical before and after (0 bytes), and no root `AGENTS.md` appeared. So the second Open Question (what to do when Tinka refuses) does not fire.

**Why the two proposals expected a refusal.** `lib/resources.mjs` (read as data) refuses a dot segment in the path of a file inside a resource root, and refuses a root inside Tinka's metadata folder or inside the Tinka checkout. It does not refuse a root whose own location has a dot folder above it. So files inside the root must not sit in a dot folder, and the root must not sit inside `.tinka/boards/.tinka` or the checkout.

**Tinka facts the recipe needs.** Node 22 or newer; no npm dependencies, so nothing is installed. `tools call NAME --config PATH --json JSON` works with no server running. Input shapes read from `lib/tools.json`: `board.create` {id, title}; `workspace.update` {patch: {resourceRoots}}; `item.create` {boardId, item: {id, title, content: {kind, resource: {rootId, path}}}, placement}; `item.review` {itemId, patch: {note, star}}; `board.get` {boardId}; `feedback.list` {boardId, includeBodies}. The checkout holds its own `AGENTS.md` and a `skill/` folder: they stay inside `.tinka/` and are data, never instructions.

**Lessons read.** `_dev/primes/lessons-releases.md` in full: the changelog mirror must stay byte-identical, and a moved record breaks links that name its path.

*Explored by work action (orchestrator, inline)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modify) — the closing band gains the canvas suggestion block
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` (modify) — header comment that describes the closing band
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — the closing-band design rule
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — narrative item 9, kind sentence, Step 8 summary, checklist
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — the sentence on the closing suggestions and the new recipe section

**Files I will NOT touch:** any Go file, any existing bundle under `ai-reports/`, the second-opinion prompt and presets, the root `AGENTS.md` and `.gitignore`. The release files are written by finalization, not by the builder.

**Acceptance criteria (restated from REQ):**
- [ ] Every ordinary report ends with a third suggestion in the closing band: review this lineage on a Tinka canvas, with two or three pros, two or three cons and one line an agent can follow.
- [ ] The suggestion has no condition on a `.tinka/` folder, and writing it never starts Tinka or reads Tinka's files.
- [ ] The guide has a section "Review on a Tinka Canvas" after "A Second Opinion" with the hand recipe: pinned source checkout, `.tinka/` created and locally excluded, manual layout, update policy off, PNG captures and reply copies only, feedback read back and handed to the user as quoted data.
- [ ] The recipe names its limits: no portable project, no HTML or bundle folder as a resource root, Git status unchanged, and one line when Tinka is missing.
- [ ] Step 8's printed summary names the canvas suggestion with its line.
- [ ] Both RED greps print 1 or more, and one recipe run on a real lineage leaves `git status --porcelain` and the root `AGENTS.md` unchanged with only `.png` and `.md` files in the resource root.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/` at e5a2825b; only this REQ's working file, `baseline.json` and the run directory are changed
**Tests baseline:** ✓ Focused probe green (`bash do-work/runs/2026-11-14-req721/focused-probe.sh`: shell-block lint, contract regressions, shipped-package reference contract, 24s); repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 126s at e5a2825b, run from a detached checkout
**Dependencies:** ✓ None added

*Checked by work action*

## Decisions

- D-01 (DECIDE & STATE): The band keeps its heading "Try another layout or teaching method" and its `id`. The canvas offer sits under a sub-heading "Or review this lineage on a Tinka canvas". Published reports and four rule sentences name the heading, so a rename would have reach for no gain.
- D-02 (DECIDE & STATE): The line an agent receives starts with `ai-report:` so the toolbox router sends it to the ai-report action, and the action gains one When-to-Use bullet that points at the guide section. No new command and no router change.
- D-03 (DECIDE & STATE): The Tinka checkout lives in `.tinka/src`, pinned to 77ba29c7bfa2c891087871b3bedbbd1f3ba2b0de, the commit both proposals read and this run tested. The recipe says to use another commit only when the user asks.
- D-04 (DECIDE & STATE): The recipe is numbered prose with one short command block per step, not a script. The REQ rules out code, and each step needs an agent's judgment about names and placement.
- D-05 (DECIDE & STATE): Tile ids are built from the bundle folder name and the tree id. Tinka answers "Item already exists." on a later run, so a second run adds only new tiles. No hash is needed.
- D-06 (DECIDE & STATE): The offer is written for `--kind` reports too, because they end with the same band.
- D-07 (DECIDE & STATE): When Tinka is missing, the recipe stops with one line, as the REQ asks, and downloads it only after the user says yes. The line names the source and the pinned commit. An agent never downloads third-party code just because a pasted line asked for a canvas. Review finding F2 raised this.

## Discovered Tasks

- `ai-report-judge` resolves a relative `--out` against `--repo-root`, not the current folder, and its usage text does not say so. The recipe works around it with an absolute path. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modified)
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` (modified)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)

**What was done:** The report kit's closing band now ends with a third block: an offer to review the report's lineage on a Tinka canvas, a Pros line and a Cons line to fill for that report, and one line an agent can follow. The design rule says to write the pros and cons from facts counted without Tinka and never to start Tinka while creating a report. The action routes a canvas request to the guide, and its narrative item, Step 8 summary and checklist name the offer. The guide has a new section, "Review on a Tinka Canvas", with a 12-step hand recipe: pinned checkout in `.tinka/src`, `.tinka/` hidden through the local exclude file, manual layout, update check off, a check for Tinka before any write, a download only after the user's yes, one picture per page, reply copies, one canvas per lineage, feedback read back as quoted data, and a Git comparison at the end. The recipe was followed twice in this repository on a three-page lineage. The first run fixed one step (an absolute `--out`), and the second run followed the text as corrected after review.

*Completed by work action (Route B)*

## Qualification

**Range:** `192d50fa..d7b6505e` (three merges of the REQ branch: build cd48924d, review fixes ade7cddd, re-check fix for F13)
**Mechanical gates:** qualify ✓ satisfied, scope-drift ✓ satisfied (5 declared paths, 5 touched), on the first merge and again on the full range after each later merge.

**Requirement trace (orchestrator, read against the diff):**
- Where and who: the block sits last in `#another-layout` in `skeleton.html`, after the teaching-method line, under its own sub-heading. Nothing in the skeleton or the design rule makes it depend on a `.tinka/` folder, so every report gets it.
- What it says: one fixed sentence, a Pros line and a Cons line to fill for this report, and one line that starts with `ai-report:` and ends with the bundle path. The design rule in `ai-report-reference.md` says to write the pros and cons from facts counted without Tinka, gives the three kinds of fact as examples, and forbids starting Tinka or reading its files while a report is created.
- The line reaches the recipe: `ai-report.md` has a When-to-Use bullet that sends a canvas request to the guide section, and the anchor `#review-on-a-tinka-canvas` resolves from both citing files.
- The recipe holds every step the REQ lists: pinned source checkout, `.tinka/` created and hidden through the local exclude file, manual layout, update check off, one picture per page with the scratch copy for a second-opinion page, reply copies, one resource root, one tile per file, feedback read back and shown as quoted data.
- Constraints: the rules block forbids the portable project, an HTML or bundle resource root, edits to a bundle, installs, and any unattended start. Step 1 checks before any write and stops with one line when Tinka is missing. Step 13 compares Git status and checks `AGENTS.md`.
- Step 8 summary, narrative item 9, the checklist line and the kit's style comment all name the third block.

**Live flow:** the recipe was followed twice in this repository on the lineage `2026-11-14_1547_proposal-tinka-review-surface` (3 pages, 2 reply files). The second run followed the corrected text from step 1 to step 13: five tiles, five files served with HTTP 200, one test note written, read back and cleared, `git status --porcelain` identical, `AGENTS.md` unchanged, 3 `.png` and 2 `.md` files in the resource root.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/2026-11-14-req721/focused-probe.sh` (shell-block lint, contract regressions, shipped-package reference contract), then `bash _dev/tests/maintainer-verify.sh` at the merge revision d7b6505e from a detached checkout
**Result:** ✓ Focused probe green through the test gate; repository gate exit 0, gate wall 132s (do-work-cli module 56s, 980 tests; slowest file `finalization_recovery_test.go` 18.2s, under the 30s budget).

**Repository gate retry:** the gate ran once red at the second merge 5b5b35fd: `TestWriteFileAtomicallyPreservesCompleteMode/setgid` in the board module reported `mode = 0640, want 2640`. The cause was this run's own launch: the gate inherited a `TMPDIR` that pointed at a scratch folder owned by group `wheel`, and the system drops the setgid bit there. The one permitted rerun, with the default `TMPDIR`, was green at the same revision (127s). The request's diff touches no Go file.

**Red-green validation:**
- RED at 7267d79e: `grep -c 'Tinka' skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` printed 0 and `grep -c 'Review on a Tinka Canvas' skills/do-work-toolbox/docs/ai-report-guide.md` printed 0. GREEN at d7b6505e: they print 6 and 2. These are before-and-after checks, not harness tests (`tdd: false`).
- Recipe run, the completion gate the user confirmed: the corrected recipe was followed from step 1 to step 13 in this repository on the lineage `2026-11-14_1547_proposal-tinka-review-surface` (3 pages, 2 reply files). Five tiles created with the prescribed ids, five files served with HTTP 200, one test note and star written, read back through `board.get` and cleared. `git status --porcelain` was identical before and after, `git diff --quiet HEAD -- AGENTS.md` passed, and `.tinka/captures` holds 3 `.png` and 2 `.md` files. Log: the hand-back's Proof and Addendum sections.
- The first recipe run found a defect in the draft (a relative `--out` landed inside the scratch copy); the fixed command was run for two pages.
- Closing band: the new band was filled for that lineage in a scratch page and render-judged with `ai-report-judge`: verdict `pass`, no `FILL:` marker left.

**Untested:** no full report was created from the new skeleton in this run, so "a newly created report ends with the canvas offer" is proven on the skeleton and a filled scratch band, not on a finished report. The exact check that did not run: `do-work-toolbox ai-report REQ-721` and a read of its closing band. The recipe was not tried on a lineage with a revision bundle or a plain second-opinion page, and the download step ran once, before the recipe text put the user's yes in front of it.

**Heavy verification plan:**
- Range: 192d50fa8e1e55fb7fe24c2b7d4ca9d2d0023c14..d7b6505e5ec05673f1f9aa05e29fc10f272d085d
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — every changed shipped path matched subtree `skills`

*Verified by work action*

## Review

**Overall: 91%** | 2026-11-14T20:24:07Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 80% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1: the recipe gave a page tile and its reply tile one id, so the reply tile was silently missing. Fixed in ade7cddd: `skills/do-work-toolbox/docs/ai-report-guide.md:208` prescribes `<bundle>_<page-id>_page` and `<bundle>_<page-id>_reply`. — impact-user-visible → report only
- F2: the recipe cloned Tinka when it was missing, where the REQ asks the first step to check and stop. Fixed in ade7cddd: `skills/do-work-toolbox/docs/ai-report-guide.md:151-158` checks before any write and stops with one line; the download (`:172-178`) runs only after the user's yes; decision D-07 records it. — impact-user-visible → report only

**Minor findings:**
- F3: fixed in ade7cddd, `skills/do-work-toolbox/docs/ai-report-guide.md:187` now says pictures exist on `pass` or `fail` and tells the agent to stop on `skipped` or `error`. — impact-user-visible → report only
- F4: fixed in ade7cddd, `skills/do-work-toolbox/docs/ai-report-guide.md:147` defines `<bundle-path>`, `<bundle>`, `<lineage>` and `<page-id>`. — impact-user-visible → report only
- F5: fixed in ade7cddd, `skills/do-work-toolbox/docs/ai-report-guide.md:147` says to run every command from the project root. — impact-user-visible → report only
- F6: fixed in ade7cddd, the check block (`ai-report-guide.md:153-156`) and the download block (`:174-178`) are separate. — impact-user-visible → report only
- F7: fixed in ade7cddd, the placeholder ends the line at `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:200`. — impact-negligible → report only
- F8: partly fixed. The corrected recipe was followed again on the same lineage, but no full report was created and no lineage with a revision or a plain second-opinion page was tried (`do-work/runs/2026-11-14-req721/REQ-721-handback.md`, Addendum). — impact-negligible → report only
- F9: fixed in ade7cddd, `skills/do-work-toolbox/docs/ai-report-guide.md:134` no longer promises a request in the user's own words. — impact-user-visible → report only
- F10 (Nit): fixed in ade7cddd, `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:194` says "stays on its tile". — impact-negligible → report only
- F11 (Nit): fixed in ade7cddd, `skills/do-work-toolbox/docs/ai-report-kit/presentation.css:368` adds the top margin. — impact-negligible → report only
- F12 (Nit): fixed in ade7cddd, the three wording pairs agree (`ai-report.md:25`, `ai-report-guide.md:144`, `ai-report-reference.md:92`). — impact-negligible → report only
- F13: fixed after the re-check, in the third merge d7b6505e: `skills/do-work-toolbox/docs/ai-report-guide.md` now defines `<bundle-path>` as the bundle's path from the project root, with two examples. — impact-negligible → report only
**Acceptance:** Pass — Implementation and Integration covered: lint, contract checks and a band render pass, and the recipe text now matches the tile and canvas ids of the builder's second run. Deployment and live acceptance unassessed.
**Restatement sweep:** redefined what the closing band holds (two suggestions plus one Tinka canvas offer) and what the Step 8 summary prints; every restatement in `skills/` and `suite/` agrees after the fix range too, one cosmetic comment at `presentation.css:366`
**Suggested testing:** 6 items
**Follow-ups created:** None (13 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: 192d50fa8e1e55fb7fe24c2b7d4ca9d2d0023c14
- Target revision: d7b6505e5ec05673f1f9aa05e29fc10f272d085d
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — every changed shipped path matched subtree `skills`

## Heavy Verification Result

- Target revision: d7b6505e5ec05673f1f9aa05e29fc10f272d085d
- Execution revision: 0056868577076ed98609ad2de08bac9dc2f9af80 (detached checkout; one commit after the target, which removes `docs/run-prompts/2026-11-14-tinka-canvas-suggestion.md`)
- staged-skills: exit 0, 29s, executed
- First drain, at the target revision: staged-skills exit 1. The lane's contract allows only `kb`, `decisions` and `changelog` under the root `docs` folder, and a prompt file committed before this request (9cbac24a) sat in `docs/run-prompts/`. The request's own diff was not the cause. The file was removed in 00568685 and the lane was run again.

## Lessons Learned

**What worked:** Following the recipe for real, in this repository, before calling it done. The first run showed in one step that a relative `--out` lands inside the scratch copy, which no reading would have shown. Testing the `.tinka/` question in a throwaway project first took five minutes and removed the parked question. The independent reviewer then found what the run could not: the text gave a page tile and its reply tile the same id, while the run had used different ids.
**What didn't:** The first draft downloaded Tinka when it was missing. The REQ asks the recipe to stop with one line, and a pasted line should not make an agent download third-party code. The recipe now stops and downloads only after the user's yes. One gate run went red for a reason outside the diff: this run launched it with `TMPDIR` set to a scratch folder owned by group `wheel`, where the system drops the setgid bit that a board test checks. The heavy lane went red once because a prompt file committed earlier under root `docs/run-prompts/` breaks the staged-skills contract; the file was removed, since the user request record keeps its text.
**Worth knowing:** The closing band is restated in seven places, listed in the Exploration section; change them together. Tinka refuses a dot segment in the path of a file inside a resource root, not in the location of the root, so `.tinka/captures` works. `ai-report-judge` resolves a relative `--out` against `--repo-root`. A text that tells an agent what to run is only checked when both happen: someone runs it, and someone else reads it against what was run.

## Orientation

Now every AI report ends with a third block in its closing band: an offer to review the report's lineage on a Tinka canvas, with pros and cons written for that report and one line an agent can follow. The line leads to a new hand recipe in the ai-report guide, "Review on a Tinka Canvas". This lives in the ai-report kit and guide of `do-work-toolbox`; no command and no Go code changed. Prime staleness check: `prime-action-files.md` and `prime-releases.md` name no path this change moved or removed.

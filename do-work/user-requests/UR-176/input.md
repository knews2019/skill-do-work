---
id: UR-176
title: 'Follow-ups after the Tinka work: decision record, map preset, judge, catalog, run-prompts folder'
created_at: 2026-11-14T22:47:58Z
requests: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-730, REQ-731, REQ-732]
word_count: 2305
---
## Summary
The maintainer abandoned Tinka (a local board tool that was tried as a review canvas for report pages) and asked to capture five follow-up prompts written on 2026-11-15. The prompts are files under `_dev/run-prompts/` at commit ef4f4dc6; their full text is in the verbatim input below, after the maintainer's five lines. Two prompts hold two requests each, so the capture made seven requests.

## Extracted Requests

| REQ | Prompt file | Request |
| --- | --- | --- |
| REQ-726 | `record-no-review-canvas-decision.md` | Write one decision record: do-work adopts no external review canvas (Tinka) for ai-report bundles. |
| REQ-727 | `relation-map-preset.md` | Add one shipped second-opinion preset, `relation-map`, so `ai-report second-opinion --preset relation-map <bundle>` writes a page that shows a report as one zoomable map of its items and the relations between them. |
| REQ-728 | `judge-renders-one-named-page.md` | Add an option to `ai-report-judge`, for example `--page <file.html>`, that renders one named page of a bundle in place. |
| REQ-729 | `judge-renders-one-named-page.md` | Make the rule for a relative `--out` of `ai-report-judge` visible: either resolve it against the current folder, or state in the usage text and in the action that it is resolved against `--repo-root`. |
| REQ-730 | `catalog-lists-second-opinion-pages.md` | Under each bundle's row, the report catalog (`ai-reports/index.html` and `ai-reports/catalog.json`) lists that bundle's second-opinion pages with tree id, preset, model and date. |
| REQ-731 | `catalog-lists-second-opinion-pages.md` | After `ai-report` creates a report, its closing summary prints the one command that rebuilds the catalog, when the new bundle is not in `ai-reports/catalog.json`. |
| REQ-732 | `run-prompts-folder-rule.md` | Let the staged-skills contract accept `docs/run-prompts/` under the root `docs` folder, keep that folder out of the shipped package, and move the prompt files from `_dev/run-prompts/` there. |

## Batch Constraints
- No order was requested between the requests. None depends on another.
- No existing `ai-reports/` bundle changes in any request.
- Requests that change files under `skills/` are releases; the decision record and the run-prompts folder change are not.
- REQ-727 is captured as blocked: it waits until the maintainer has looked at the example map page.

## Capture Clarifications
Asked during capture and answered by the maintainer on 2026-11-14. The answers are recorded on the requests they belong to.

1. Has the maintainer looked at the example map page, and should it become a standard option (REQ-727)? The maintainer's answer:

> ```
> Not looked yet, hold it
> ```

2. How should a new report get into the catalog: a printed reminder, or a rebuild on every new report (REQ-731)? The maintainer's answer:

> ```
> Print a reminder (Recommended)
> ```

   This settles the choice that the catalog prompt left open under its second gap.

3. Where do run prompts live in this repository (REQ-732)? The maintainer's answer:

> ```
> Allow docs/run-prompts/ (Recommended)
> ```

   This settles the choice that the run-prompts prompt parked as a pending answer.

## Full Verbatim Input
> ```
> do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-record-no-review-canvas-decision.md
> do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-relation-map-preset.md
> do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-judge-renders-one-named-page.md
> do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md
> do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-run-prompts-folder-rule.md
> 
> Contents of _dev/run-prompts/2026-11-15-record-no-review-canvas-decision.md at commit ef4f4dc6:
> 
> # Capture prompt: record the decision to adopt no external review canvas
> 
> ## How to run this prompt
> 
> Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.
> 
> ## The request
> 
> Write one decision record: do-work adopts no external review canvas (Tinka) for ai-report bundles. The record exists so a later session does not propose it again without new facts.
> 
> ### Context
> 
> - Tinka (github.com/bladnman/tinka) is a local board tool for trying designs with an AI agent. Two proposals tested it on 2026-11-14 as a place to review report pages: `ai-reports/2026-11-14_1452_proposal-tinka-review-surface` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface`.
> - REQ-721 (archived under `do-work/archive/UR-172/`) shipped a canvas offer in every report plus a hand recipe as version 0.316.0. Version 0.316.1 fixed the canvas id. Version 0.317.0 removed both the same day. The three changelog entries are in `CHANGELOG.md`.
> - What was observed:
>   - A report page is about 5,000 pixels tall. As a picture tile it is a thin strip that cannot be read, and folded sections stay closed (38 to 46% of the text, measured in the second proposal).
>   - Tinka answers only on 127.0.0.1, so a phone cannot open it.
>   - A live HTML tile served from Tinka's own address can write feedback by itself. A test page did it.
>   - A live tile served from a second local port cannot write (tested 2026-11-14: the browser blocked the read and one write, Tinka refused two writes with HTTP 403). But it is live only after two clicks and one page at a time, which a browser tab already does.
>   - The marks Tinka adds (star, vote, note) move feedback from chat to a second place and add steps.
> - What would reopen the question: Tinka shows several live pages side by side with an isolation that stops them from writing feedback, or do-work gets a project whose work is many rounds of screen designs.
> - The page `index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` in the second proposal's bundle is the direction taken instead: a zoomable relation map inside the report.
> 
> ### Files to change
> 
> - `docs/decisions/records/adr-023-<slug>.md` (new), in the format of `adr-022-one-broken-pipe-does-not-stop-the-factory.md`.
> - `docs/decisions/log.md`, `docs/decisions/_master_index.md` and the fitting topic index under `docs/decisions/topics/` (probably `_index_pipeline-deliverables.md`), as the existing records are linked.
> 
> ### What must not change
> 
> - Nothing under `skills/` or `suite/`. This is a maintainer record, so it is no release.
> - No existing `ai-reports/` bundle.
> - No new command and no new file outside `docs/decisions/`.
> 
> ### Verification
> 
> RED first: `ls docs/decisions/records | grep -c 'review-canvas'` prints 0.
> 
> GREEN: the record exists with status accepted and names the two proposals, REQ-721 and the three versions as sources; the log and the topic index link it; `git diff --stat HEAD -- skills/ suite/` is empty; `bash _dev/tests/contract-regressions.sh` passes.
> 
> Contents of _dev/run-prompts/2026-11-15-relation-map-preset.md at commit ef4f4dc6:
> 
> # Capture prompt: ship a relation-map preset for ai-report second-opinion
> 
> ## How to run this prompt
> 
> Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.
> 
> ## Gate
> 
> Capture this only after the maintainer has looked at the example page and wants it as a standard option. The example is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`. If the capture gives no sign of that, ask once and park the request.
> 
> ## The request
> 
> Add one shipped second-opinion preset, `relation-map`, so `ai-report second-opinion --preset relation-map <bundle>` writes a page that shows a report as one zoomable map of its items and the relations between them.
> 
> ### Context
> 
> - The maintainer wants report content that is easy to understand by zooming in and out and by seeing relations. The example page was built by hand with `--editor-direction` and passed the render judge.
> - A preset is one Markdown file with a direction to the editor. The loader reads `<bundle>/second-opinion-presets/<name>.md`, then the shipped folder `skills/do-work-toolbox/docs/second-opinion-presets/`. No loader change is needed. `mindmap.md` there is the nearest existing preset: read it first and say in the new preset how the two differ (a mind map is a tree from one center; a relation map shows relations between any two items).
> 
> ### What the preset must ask the editor for
> 
> - Boxes are the report's own items (for example options, risks, findings, limits, files, requests). Lines are relations the report states. The editor says on the page that the lines are its reading of the report.
> - At most four kinds of line, each with one color and a plain name, and a key that explains them.
> - Semantic zoom: zoomed out, a box shows its code and title; zoomed in, a summary; a click opens the detail with its sources and its relations, and fades unrelated boxes.
> - Drag to move, pinch or Control-scroll to zoom, buttons and keys for the same, and a button that shows everything. Plain scrolling must still scroll the page.
> - The page works without JavaScript: the map shows at full size and scrolls. Every box is also given as text, and every relation as one table row.
> - No library. Inline SVG and plain JavaScript only, as the second-opinion prompt already requires.
> - About 30 boxes is the upper limit for one map. With more, group items and say so.
> - On a narrow screen the map starts on the main item, not on the whole picture.
> - Every fact still comes from the report or the repository, as in `bold.md`.
> 
> ### Files to change
> 
> - `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (new).
> - `skills/do-work-toolbox/docs/ai-report-guide.md`: the preset list in "A Second Opinion".
> - `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`: `TestSecondOpinionShippedPresetsLoadFromTheToolbox`, if it names the shipped presets.
> - `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
> 
> ### What must not change
> 
> - The preset loader, the second-opinion prompt, the other presets, and every existing `ai-reports/` bundle.
> - The report kit (`skeleton.html`, `presentation.css`, `presentation.js`).
> 
> ### Verification
> 
> RED first: `--prepare --preset relation-map` on a fixture bundle is refused because the preset does not ship (the existing shipped-presets test is the place for it).
> 
> GREEN: the same call returns state `prepared` with `preset_file` in the shipped folder; the guide names the preset; one real page is built with the preset on a bundle other than the example, passes `ai-report-judge`, and is shown to the maintainer before the request is called done.
> 
> Contents of _dev/run-prompts/2026-11-15-judge-renders-one-named-page.md at commit ef4f4dc6:
> 
> # Capture prompt: ai-report-judge can render one named page of a bundle
> 
> ## How to run this prompt
> 
> Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.
> 
> ## The request
> 
> Two changes to `ai-report-judge`, both found while building report pages on 2026-11-14.
> 
> 1. **Render one named page.** Today the judge renders only `index.html` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go`, lines 25, 111 and 168). To check or capture a second-opinion page such as `index-v1c01-<name>.html`, an agent must copy the whole bundle and write that page over the copy's `index.html`. Add an option, for example `--page <file.html>`, that renders that page of the bundle in place. `secondOpinionRender` in `report_second_opinion.go` (from line 1706) already does the copy-and-swap internally; reuse or replace it as the builder judges.
> 2. **A relative `--out`.** A relative `--out` is resolved against `--repo-root`, not against the current folder, and the usage text does not say so. With `--repo-root` pointing at a scratch copy, the captures landed inside the copy. Either resolve a relative `--out` against the current folder, or state the rule in the usage text and in the action. Park the choice as a pending answer if existing callers depend on today's behavior.
> 
> ### Files to change
> 
> - `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` and its test file.
> - `skills/do-work-toolbox/actions/ai-report.md` (the judge form) and `skills/do-work-toolbox/docs/ai-report-guide.md`, where the judge is described.
> - `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
> 
> ### What must not change
> 
> - Without the new option the judge behaves exactly as today: same captures, same `judge.json`, same exit codes.
> - The judge never writes inside the bundle. Bundles stay byte-identical.
> - The second-opinion finalize path keeps its verdicts for the same pages.
> 
> ### Verification
> 
> RED first: a Go test runs the judge with `--page index-v1c01-x.html` on a fixture bundle with two pages and fails on the unknown option.
> 
> GREEN: the test passes; the result's served address ends with the named page and the captures show that page; a named page that is missing or outside the bundle is refused with a typed finding; the fixture bundle's digest is unchanged; `go test ./...` passes in the do-work-cli module.
> 
> Contents of _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md at commit ef4f4dc6:
> 
> # Capture prompt: the report catalog shows second-opinion pages and new reports
> 
> ## How to run this prompt
> 
> Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.
> 
> ## The request
> 
> Two gaps in the report catalog (`ai-reports/index.html` and `ai-reports/catalog.json`), both seen on 2026-11-14.
> 
> 1. **Second-opinion pages are not in the catalog.** It has one entry per bundle (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`, around lines 166 to 172, `pickReportEntry`). The bundle `ai-reports/2026-11-14_1547_proposal-tinka-review-surface` holds `index.html` and three second-opinion pages (`v1c01`, `v1c02`, `v1c03`), and the catalog cannot show them. Under each bundle's row, list its second-opinion pages with tree id, preset, model and date. Each page's head already carries these as `ai-report-second-opinion-*` meta tags.
> 2. **A new report is missing from the catalog until someone rebuilds it.** Ordinary creation writes only the bundle. On 2026-11-14 two new proposal bundles were absent from `catalog.json` (46 bundles, no match for either). First read why creation does not rebuild the catalog (`skills/do-work-toolbox/actions/completed-work-presentation-reference.md`, Collision-Safe Publication). If that reason still holds, do the smaller change: Step 8 prints the `ai-report index` line when the new bundle is not in the catalog. If it does not hold, creation rebuilds the catalog as the revise form does. Park the choice as a pending answer when the reason is unclear.
> 
> ### Files to change
> 
> - `report_index.go` and its test file in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`.
> - `skills/do-work-toolbox/actions/ai-report.md` (catalog forms, Step 8) and `skills/do-work-toolbox/docs/ai-report-guide.md` ("Report Catalog").
> - `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
> 
> ### What must not change
> 
> - `catalog.json` changes are additive: every existing field keeps its name and meaning, so `revise`, `find`, `organize` and `second-opinion latest` keep working.
> - No bundle is edited. The catalog reads bundles and never writes into them.
> - A page without the tree meta tags (a plain second-opinion page) is still listed, by its file name.
> 
> ### Verification
> 
> RED first: a Go test builds a fixture bundle with `index.html` and one `index-v1c01-*.html` page and asserts that the catalog entry lists that page; it fails because the entry has no such list.
> 
> GREEN: the test passes; `ai-reports/index.html` shows the pages under their bundle; a second test covers the chosen fix for gap 2; `go test ./...` passes in the do-work-cli module; existing catalog tests pass unchanged.
> 
> Contents of _dev/run-prompts/2026-11-15-run-prompts-folder-rule.md at commit ef4f4dc6:
> 
> # Capture prompt: decide where run prompts live in this repository
> 
> ## How to run this prompt
> 
> Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.
> 
> ## The request
> 
> Settle one conflict between a maintainer rule and a repository check, with the smallest change.
> 
> ### Context
> 
> - The maintainer's global rule says: write a prompt for another session to `docs/run-prompts/<date>-<slug>.md` and commit it.
> - In this repository that file breaks a check. `_dev/tests/staged-skills-contract.sh` (the loops near lines 75 and 99 to 114, message "legacy root runtime must be retired at modular cutover") allows only `docs/kb/`, `docs/decisions/` and `docs/changelog/` under the root `docs` folder. On 2026-11-14 the file `docs/run-prompts/2026-11-14-tinka-canvas-suggestion.md` turned the heavy lane `staged-skills` red five commits after it was added, and was removed in commit 00568685. The fast gate did not catch it.
> - Since then run prompts for this repository are kept in `_dev/run-prompts/`, which is maintainer-only.
> 
> ### The choice (park it as a pending answer; recommended first)
> 
> 1. **Allow `docs/run-prompts/` in the contract** (recommended). It is a maintainer folder like the three allowed ones, and the global rule then works here as in every other project. Add it to both lists in the contract, check that it is export-ignored like the other maintainer docs, and move the files from `_dev/run-prompts/` there.
> 2. **Keep `_dev/run-prompts/`.** Change nothing in the contract. Write one line in the maintainer `CLAUDE.md` that names this folder as the home of run prompts in this repository.
> 
> ### Files to change
> 
> - Choice 1: `_dev/tests/staged-skills-contract.sh`, the export-ignore rules if needed, and the moved prompt files.
> - Choice 2: `CLAUDE.md` only.
> 
> ### What must not change
> 
> - Nothing under `skills/` or `suite/`, so this is no release.
> - The contract keeps refusing every other new file under the root `docs` folder.
> 
> ### Verification
> 
> RED first (choice 1): with a file `docs/run-prompts/probe.md` present, `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` fails with the legacy-root message.
> 
> GREEN (choice 1): the same lane passes with that file present and still fails with a file `docs/other/probe.md`; the probe files are removed afterwards. GREEN (choice 2): `CLAUDE.md` names the folder and the lane passes.
> ```

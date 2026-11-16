# Builder Brief — REQ-738 (the ai-report second-opinion launcher shows a short header and folds the full direction)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-738-fold-second-opinion-direction`.
- Worktree (the ONLY tree you write, with the three exceptions below): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-738-fold-second-opinion-direction`. The coordinator creates it from main HEAD (`7f2c4b94` or later, VERSION 0.324.0) right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-08, Exploration and Scope: read `do-work/working/REQ-738-second-opinion-direction-block-first-screen.md` from the MAIN tree (your worktree may hold an older copy).
- Main-tree paths you may write, by absolute path, never staged and never committed by you:
  1. Hand-back file: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-738-handback.md`
  2. Captures folder: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-738-captures/`
  3. Deleting, by name, the files your own browser session created in the main tree's `.playwright-mcp/`.
- A second builder runs at the same time in its own worktree: REQ-737 (the relation-map preset keeps box titles readable after a click on a phone). It changes preset text only. You share no source file with it. Integration is serial: your REQ first, REQ-737 last.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-releases.md`, `_dev/primes/prime-action-files.md` and the required lesson `_dev/primes/lessons-releases.md`.
3. Dropped only for budget, read them with a fixed-string grep, not as whole files: the three bullets `[family: writer-anchored-marker-position]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (lines 103-105: no whitespace outside the markers, strip on the masked surface, decorating twice gives the same bytes) and the second `[family: session-is-the-backend]` bullet (line 18: the block moves the map, so measure after decorating).
4. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
5. The code: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` lines 1563-1648 (`secondOpinionStripBlocks`, `secondOpinionDecorate`), and the two tests that pin the header: `report_second_opinion_test.go` lines 942-1002 and `report_second_opinion_tree_test.go` lines 43-90.
6. The kit rules the new element meets: `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` lines 196-199 and `presentation.js` lines 44-50.

## Decided (do not reopen)
- **D-01 (recommended form):** the block stays where it is, between the same two markers, as `<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">`, a direct child of body. A published page styles exactly that selector (`body>aside[aria-label="Second-opinion direction"]` and its `p`), so keep the element, the label and `p` elements for the short text. Inside: a short always-visible part (the iteration line with the parent link on a tree page, and the preset name or `custom`), then one closed `details` whose `summary` is a short label such as "Full direction" and whose body is the `white-space:pre-line` direction text. The header's `details` and `summary` carry inline styles that cancel the kit's card look. No script. If your measurement shows this form cannot meet the numbers, pick another form, record it as a decision (D-09 onward) and name the forms you rejected.
- **D-02:** the browser proof uses scratch copies decorated by a throwaway Go test file. Published pages under `ai-reports/` are not rewritten. Details under "Browser proof".
- **D-03:** "about 80 px" is approximate. Up to about 100 px counts as met when the map frame top is below y = 844 on the phone. Report exact numbers. Measure the opened fold too; it has no limit.
- **D-04:** the RED test is one new function named exactly `TestSecondOpinionHeaderFoldsTheFullDirection` in `report_second_opinion_test.go`. The GREEN probe looks for that name.
- **D-05:** the integrator picks the version bump (a patch is suggested). You only propose the changelog entry.
- **D-06:** the maintainer sees your phone captures before integration starts. Save them as described under "Captures".
- **D-07:** the browser matrix below is the reduced one for this REQ: you do not change a map, so the click sweep only shows that nothing else moved.
- **D-08:** edit the two existing tests only where your form breaks an assert. A declared file you did not need goes under "Scope changes" as "declared, not needed".

## What to build
Three commits.

**Commit 1 (RED, test only).** Add `TestSecondOpinionHeaderFoldsTheFullDirection` to `report_second_opinion_test.go`. Call `secondOpinionDecorate` directly, so the test costs milliseconds (the file already runs close to its 30-second budget). Read the shipped relation-map preset text from `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (see how `TestSecondOpinionShippedPresetsLoadFromTheToolbox` finds the toolbox folder). The comment above the test names the real failure it pins: a 3,057-byte direction printed as 1,302 px of open text above the page on a phone (REQ-727 D-17). Pin these, and nothing decorative:
1. On a tree page decorated with the relation-map text: take the bytes between the two markers, remove the folded element, and the rest does not contain the direction text. The folded element contains the whole escaped direction text.
2. The visible rest contains the preset name, the iteration id and the parent link.
3. The six head metas equal one literal expected string (preset, direction, id, base, model, date), so they stay byte-identical to 0.324.0.
4. Decorating twice gives the same bytes.
5. The old-block strip case: a page that holds a literal 0.324.0 header block (copy its exact shape from line 1634 and 1641 of today's code into the test as a string) decorates to the same bytes as the same page without it.
6. A `bold` page and an editor-direction-only page show `bold` or `custom` in the visible part and the full direction in the folded part.
7. The block carries the long-token guard (`overflow-wrap:anywhere`).
Run the test and record the failure text. It must fail on check 1, not on a compile error or a missing file.

**Commit 2 (GREEN, the launcher).** Change only what `secondOpinionDecorate` writes inside the block (today's lines 1631, 1634 and 1641) and its function comment (lines 1597-1601). Do not change the strip calls, the metas, the insert position, the markers, or any other function. No whitespace outside the markers. Escape every value as today. Update an assert in the two existing header tests only if your form breaks it. Run the new test (passes) and the whole package (see Verify).

**Commit 3 (contract text).** Two sentences, nothing else in those long paragraphs:
1. `skills/do-work-toolbox/actions/ai-report.md` line 89: replace "showing the preset and the full direction" with words that say what the block shows now (the preset at the top, the full direction folded below it, in your chosen form). The probe fails while the old words `showing the preset and the full direction;` are still there. Keep the marker names in the sentence.
2. `skills/do-work-toolbox/docs/ai-report-guide.md` line 124: replace "the page shows the direction at its top" with plain words for the same fact (for example: the page names the preset at its top and keeps the full direction behind a fold).

**Not yours:** `VERSION`, both `CHANGELOG.md` files, anything under `do-work/` (even though the REQ's `write_set` lists the release files: the integrator writes them); the preset loader, the prompt builder, the iteration list; `ai-report-second-opinion-prompt.md`; every preset (REQ-737 owns `relation-map.md`); the report kit; the judge; every file under `ai-reports/`.

## Browser proof (scratch copies, never committed)
Work in a scratch folder outside every Git tree. Decorate with a throwaway Go test file in your worktree (for example `internal/toolboxcommands/zz_req738_scratch_test.go`) that reads a page, calls `secondOpinionDecorate`, and writes the result to the scratch folder. Delete that file before your last commit and check `git status` shows it gone.

Scratch pages to make:
- **S1 (the proof page):** copy the bundle folder `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/` to scratch. Decorate its `index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html` with the direction rebuilt from its own metas: preset `relation-map`, the shipped preset file as the text, id `v1c03`, base `index.html`, model `claude claude-opus-5-5`, date `2026-11-15`. This page holds a 0.324.0 block, so it is also the old-block strip case on real bytes. Check that the six metas of the result equal the published page's metas byte for byte.
- **S2 (custom):** in the same scratch bundle, decorate `index-v1c01-codex-gpt-6.1-terra-custom-matrix-to-implementation.html` the same way from its metas (preset slot `custom`, the direction meta is the editor direction). This page has six `details` of its own.
- **S3 (kit page, bold):** a copy of `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` with `presentation.css` and `presentation.js` beside it, decorated with preset `bold` and the shipped `bold.md` text as a tree page.
- **S4 (long token):** any small page decorated with an editor direction that holds one 400-character token with no space.

Serve the scratch folder with `python3 -I -m http.server <free port> --bind 127.0.0.1` (the browser tool blocks `file:` URLs). Measure with ONE script, at 390 x 844 and at 1440 x 1000, light theme. Put the script in scratch, and put its path and its raw output in the hand-back. Per width it reports:
1. S1 at load, fold closed: header height (the `aside` box plus its margins), the map frame's top y, the page's sideways scroll (`scrollWidth` against `clientWidth`). Targets: header about 80 px or less (D-03), map frame top below 844 at 390 wide, no sideways scroll. For comparison, 0.324.0 gave 1,302 px and y = 2,069 on the phone, 514 px and y = 853 on the desktop.
2. S1 with the fold opened: header height, sideways scroll, and that the direction text reads with its line breaks.
3. S1 click sweep, fold closed: click every shown box from the first view and again after "Show everything". Per click: the clicked title's px size, each shown related title's px size, boxes out of frame, sideways page scroll. You do not change the map, so these numbers are for comparison with REQ-727's and for the REQ-737 integrator; a title size is not a pass or fail for you. Any new sideways scroll or a script error is.
4. S1 Tab order: the first Tab stops (the parent link, then the `summary`) are in view and have a visible focus ring; Enter or Space on the `summary` opens the fold.
5. S1 with JavaScript off: the short header shows, the fold opens by a click, the list under the map is still there.
6. S2: the header shows `custom`, the fold holds the editor direction, the page's own six `details` look and work as before.
7. S3: the header's computed style is not the kit card (no 2 px border, no 14 px radius, no bold 1.1 rem summary), its height is about 80 px or less, and a click on "Open all details" opens the kit's three `details` and leaves the header's fold closed.
8. S4: no sideways scroll with the fold closed and with it open.

Run the render judge on S1: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge <scratch bundle> --page index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html --out <scratch>/judge`. It must report `pass`. If the judge refuses a bundle outside the repository, copy the scratch bundle into an untracked folder of your worktree's `ai-reports/`, judge it there, and delete the folder before your last commit; record that as a decision.

**Captures (D-06).** Save final PNG captures into the captures folder named above, by absolute path: S1 at 390 x 844 at load (fold closed, the map in the first screen), S1 at 390 x 844 with the fold opened, S1 at 1440 x 1000 at load, S3 at 390 x 844 at load. Add S2 if it looks different. List each file with one line on what it shows in the hand-back section `## Captures for the maintainer`, and name there the form you chose and the forms you rejected. Delete every other capture before hand-back. Captures the Playwright tools write land in the main tree's gitignored `.playwright-mcp/`: delete the files your session created there by name, nothing else, and record that as a decision. Stop the server.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-738]`. End each commit message with the line `Co-Authored-By: Claude Fable 5.2 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not ask the user a question: an open choice becomes a decision in the hand-back.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. Two builders share the machine, so a 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.
- Nothing under `/Users/t2/.claude/` is ever edited.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l internal/toolboxcommands/report_second_opinion.go internal/toolboxcommands/report_second_opinion_test.go internal/toolboxcommands/report_second_opinion_tree_test.go` prints nothing; `go vet ./internal/toolboxcommands/` is clean; `go test -count=1 ./internal/toolboxcommands/` passes (the whole package, as the REQ's Builder Guidance asks, because the change touches every directed page; record its time).
2. `bash _dev/tests/shipped-package-reference-contract.sh`, `bash _dev/tests/contract-regressions.sh`, `bash _dev/tests/quiet-grep-pipeline-audit.sh`, and with `DO_WORK_MAINTAINER_TIER=heavy`: `bash _dev/tests/staged-skills-contract.sh` (about 35 s).
3. `git diff --stat main...HEAD` lists only files of the Scope: the launcher file, one or two test files, `ai-report.md`, `ai-report-guide.md`. `git status --short` shows no scratch test file and no untracked folder under `ai-reports/`.
4. `git diff main...HEAD -- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` touches only `secondOpinionDecorate` and its comment.
5. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-738-probe.sh`, run with your worktree as the current directory, prints `REQ-738 probe PASS (...)`. It fails at base by design: the new test does not exist, so its first check stops it. After commit 2 it still fails at check 2 (the contract sentences) until commit 3.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb.
- `## Scope changes` — including a declared file you did not need (D-08).
- `## Red-green evidence` — the new test failing after commit 1 (the failure text) and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; the measurement script's path and raw output; the eight measurements at both widths with numbers; the judge verdict; for a rerun under load, both runs.
- `## Captures for the maintainer` — each PNG with one line on what it shows; the form chosen and the forms rejected.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-09 onward (D-01 to D-08 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — see below.
- `## Proposed CHANGELOG entry` — a descriptive title that says what shipped (for example "Second-Opinion Pages Fold the Full Direction Under a Short Header"; check that no earlier entry uses it) and the line "What changes for an installed copy". There: a page that `ai-report second-opinion` builds with a preset or an editor direction now opens with a short header, the full direction is behind a fold on the same page, the head metas are unchanged, and pages already published keep their old header until a new pass rebuilds them.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
No source file is shared with REQ-737. The seam is in behaviour, and your hand-back section must give the REQ-737 integrator what it needs:
| Seam | What you report |
| --- | --- |
| REQ-737's proof page is decorated by the 0.324.0 launcher in its builder's worktree. After your REQ lands, its integrator decorates that page again with your code and measures again. | The exact new block for a tree page with preset `relation-map` (the bytes between the markers, with the direction shortened), the header height you measured at both widths, the map frame top y on S1 at both widths, and the selector to measure the header by. |
| A published page may style `body>aside[aria-label="Second-opinion direction"]` and its `p` (the tstack v1c03 page does). | Whether your form keeps that element, label and `p` children, and how S1 looks with the page's own rule applied. |
| Published pages are not rewritten. | Confirm `git diff --stat main...HEAD -- ai-reports/` prints nothing. |
| The 30-second per-file budget of `report_second_opinion_test.go`. | The time of `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/` before and after your change. |

## Notes for the integrator (the builder does nothing about these)
- **Release:** the integrator writes `CHANGELOG.md`, its byte-identical mirror `skills/do-work/CHANGELOG.md`, and `VERSION` (D-05: a patch from 0.324.0 is suggested).
- **Maintainer look:** the coordinator shows the PNG files in `REQ-738-captures/` to the maintainer before integration starts (D-06).
- **Catalog:** no rebuild; no page under `ai-reports/` changes.
- **GREEN probe:** `do-work/runs/work-2026-11-16-171725/REQ-738-probe.sh` (`-- --probe-file`). Pre-flight probe: `do-work/runs/work-2026-11-16-171725/REQ-738-preflight-probe.sh`; it passed at `7f2c4b94` on 2026-11-16.
- **After this REQ lands:** REQ-737's integrator re-decorates and re-measures its proof page (REQ-737 Builder Guidance).

Your final message: under 2500 characters. State the branch tip hash, the form you chose, the header heights and map top y at both widths, the verification results, the capture file names, and anything the integrator must decide.

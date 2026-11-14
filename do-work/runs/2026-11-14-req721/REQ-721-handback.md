# REQ-721 hand-back

Builder: the orchestrator session, on the worktree rung. Branch `worktree-agent-REQ-721-ai-report-closing-suggests-tinka-canvas`, commit cd48924d.

## Files changed

- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modified) — the closing band gains a third block: the Tinka canvas offer, a Pros line, a Cons line and the line that builds it. The authoring comment says the offer is text only.
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` (modified) — header comment names the third block.
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — new design rule "Then offer a Tinka canvas".
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — new When-to-Use bullet that routes a canvas request to the guide, narrative item 9, Step 8 summary, checklist.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — the sentence on the closing suggestions, and the new section "Review on a Tinka Canvas" with a 12-step hand recipe.

## Proof

- RED before: both greps printed 0 at 7267d79e. GREEN after: `grep -c 'Tinka' skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` prints 6, `grep -c 'Review on a Tinka Canvas' skills/do-work-toolbox/docs/ai-report-guide.md` prints 2.
- Check first: Tinka 77ba29c accepts and serves a resource root inside the project's `.tinka/` (throwaway project, HTTP 200 for a PNG and a Markdown file).
- Recipe run on a real lineage, in this repository, following the guide's steps: bundle `ai-reports/2026-11-14_1547_proposal-tinka-review-surface` (3 pages, 2 reply files). Five tiles created, all five files served with HTTP 200, one test note and star written and read back through `board.get`, then cleared. `git status --porcelain` identical before and after, `git diff --quiet HEAD -- AGENTS.md` passed, and `.tinka/captures` holds 3 `.png` and 2 `.md` files.
- The run found one defect in the first draft of the recipe: a relative `--out` is resolved against `--repo-root`, so the capture of a second-opinion page landed inside the scratch copy. The recipe now passes an absolute `--out`, and the fixed command was run for two pages.
- Focused checks in the worktree: shell-block lint, contract regressions and the shipped-package reference contract all pass.

## Decisions

- D-01 (DECIDE & STATE): The band keeps its heading "Try another layout or teaching method" and its `id`. The canvas offer sits under a sub-heading "Or review this lineage on a Tinka canvas". Published reports and four rule sentences name the heading, so a rename would have reach for no gain.
- D-02 (DECIDE & STATE): The line an agent receives starts with `ai-report:` so the toolbox router sends it to the ai-report action, and the action gains one When-to-Use bullet that points at the guide section. No new command and no SKILL.md change.
- D-03 (DECIDE & STATE): The Tinka checkout lives in `.tinka/src`, pinned to 77ba29c7bfa2c891087871b3bedbbd1f3ba2b0de, the commit both proposals read and this run tested. The recipe says to use another commit only when the user asks.
- D-04 (DECIDE & STATE): The recipe is numbered prose with one short command block per step, not a script. The REQ rules out code, and each step needs an agent's judgment about names and placement.
- D-05 (DECIDE & STATE): Tile ids are built from the bundle folder name and the tree id. Tinka then answers "Item already exists." on a later run, so a second run adds only new tiles. No hash is needed.
- D-06 (DECIDE & STATE): The offer is written for `--kind` reports too, because they end with the same band.

## Discovered Tasks

- `ai-report-judge` resolves a relative `--out` against `--repo-root`, not the current folder. The command's usage text does not say so. impact-negligible → report only
- The recipe run left `.tinka/` in this repository (checkout, layout, canvas `tinka-review-surface-1547` with five tiles) and one line in `.git/info/exclude`. This is the pilot setup the user chose, not a defect. impact-negligible → report only

## Lesson evidence

- `_dev/primes/lessons-releases.md` read in full (whole satellite). `_dev/primes/prime-action-files.md` and `_dev/primes/prime-releases.md` read in full. `_dev/primes/lessons-action-files.md` was over budget and not read.

## P-A-U

- PLAN: add a third block to the closing band and its restatements; write the recipe as a guide section; test the `.tinka/` question first; run the recipe once for real.
- APPLY: five files, all inside the declared scope.
- UNIFY: `git diff --stat` shows the five declared files, 120 insertions and 6 deletions before the recipe fix; shell-block lint, contract regressions and reference contract pass; no debug text in the diff.

## Addendum (after review)

Branch commit ade7cddd fixes review findings F1 to F7 and F9 to F12.

- F1: the page tile and the reply tile now get different ids, `<bundle>_<page-id>_page` and `<bundle>_<page-id>_reply`. These are the ids the real run used.
- F2 and F6: step 1 now checks for Node and the pinned checkout before any write, and stops with one line when Tinka is missing. The download is its own step and runs only after the user's yes. D-07 records this.
- F3: the judge sentence now matches its verdicts: pictures on `pass` or `fail`, none on `skipped` or `error`.
- F4 and F5: one paragraph defines `<bundle-path>`, `<bundle>`, `<lineage>` and `<page-id>`, and says to run every command from the project root.
- F7, F10, F11: the band's line ends with its placeholder, the offer says a note stays on its tile, and the sub-heading has a top margin.
- F9: the guide no longer promises a request in the user's own words.
- F12: three wording pairs now agree.
- F8 (no new report was created) is answered with two checks, not with a full report. The corrected recipe was followed again in this repository from step 1 to step 13 on the same lineage, with the canvas id and tile ids the recipe now prescribes: five tiles, five files served with HTTP 200, one test note written, read back and cleared, Git status identical, `AGENTS.md` unchanged. The new band was filled for that lineage in a scratch page and render-judged: verdict `pass`. No full report was created in this run, and the lineage has no revision and no plain second-opinion page.

- D-07 (DECIDE & STATE): When Tinka is missing, the recipe stops with one line, as the REQ asks, and downloads only after the user says yes. The line names the source and the pinned commit. An agent never downloads third-party code just because a pasted line asked for a canvas.

## Addendum 2 (after the re-check)

Branch commit after ade7cddd fixes re-check finding F13: `<bundle-path>` is now defined as the path from the project root, with two examples. The reviewer's re-check marked F1 to F7 and F9 to F12 fixed and F8 partly fixed.

# Builder Brief — REQ-727 (ship a relation-map preset for ai-report second-opinion)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-727-relation-map-preset`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-727-relation-map-preset`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-12, Exploration and Scope: read `do-work/working/REQ-727-relation-map-second-opinion-preset.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-727-handback.md`
- No other builder runs at the same time. REQ-733 (rebuild the relation map page) and REQ-734 (the table page and the map-versus-table comparison) are merged and archived; the maintainer confirmed the map as a standard option through clarify.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-releases.md` and the required lesson `_dev/primes/lessons-releases.md`. Dropped only for budget: the bullet `[family: session-is-the-backend]` at `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:17` (you are the second-opinion backend; read it with a fixed-string grep).
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin). The maintainer's scope cue for this REQ: keep language clear, do not complicate it.
4. The existing presets in `skills/do-work-toolbox/docs/second-opinion-presets/`, `mindmap.md` and `bold.md` first.
5. The sources the preset must learn from: `do-work/archive/UR-177/REQ-733-rebuild-relation-map-page.md` (Decisions D-08 to D-20, Discovered Tasks, Lessons Learned), its three reviews `do-work/runs/work-2026-11-15-012821/REQ-733-review.md`, `-review2.md`, `-review3.md`, and `do-work/archive/UR-178/REQ-734-proposal-table-and-map-comparison.md` (Review, "Map vs table"). The example page is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html`; the REQ's Exploration names the lines that hold each technique.
6. `skills/do-work-toolbox/actions/ai-report.md`, Second-opinion form, Session sequence (from line 81), and the judge `--page` form (line 190).

## Decided (do not reopen)
- **D-01:** the real page goes on `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1`, base `index.html`, Session form, `--backend claude --model <your own declared model slug>` (recorded `session-declared`), `--preset relation-map` alone. It becomes tree page v1c03 there. Fallback rule is in the REQ.
- **D-02:** the preset is one short opening paragraph, then numbered checks of one sentence each, about 3,500 characters or less. The launcher prints the whole direction at the top of the page, so length costs the reader.
- **D-03:** the layout rule for any report (columns in argument order; a first view of about nine main boxes; a click reveals related boxes; "Show everything").
- **D-04:** the seven defect checks, worded in the REQ. Each must be in the preset as a check a reader can see or measure.
- **D-05:** what a preset cannot guarantee stays a decision in the REQ, not text in the preset. The preset may say once that the editor checks its page in a browser when it can.
- **D-06:** the test change (named list plus "named preset" wording; no new test function).
- **D-07:** docs change only in `ai-report-guide.md:124`. Not `ai-report.md`, not `help.md`, not the report kit.
- **D-08:** your hand-back opens with the real page's path; the maintainer looks at it after the run.

## What to build
Three commits.

**Commit 1 (RED, test only).** In `report_second_opinion_test.go`, `TestSecondOpinionShippedPresetsLoadFromTheToolbox`: add `"relation-map"` to the named list at line 1026. Change the comment above it (line 1025), the error text (line 1028) and the function comment (line 1004, "a learning-method preset the user curated") so they say "named preset" and name UR-176 beside UR-164. Run the test and record the failure (`named preset "relation-map" does not ship`).

**Commit 2 (GREEN, the preset).**
1. `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (new). The opening paragraph says what the page is, and it says how it differs from `mindmap`: a mind map is a tree from one center, while a relation map shows relations between any two of the report's items. Then the numbered checks. They cover every item of the REQ's Detailed Requirements as changed by Addendum 1, the layout rule of D-03, and the seven checks of D-04. Also include: light dimming (unrelated box titles stay readable); a deselect that returns the detail panel to its first content; "No library: inline SVG and plain JavaScript only, no Mermaid, no remote file", because a direction overrides the prompt's Rule 4; and "Every fact still comes from the report or the repository." The text must not contain `{{`. Plain words, one idea per sentence.
2. `skills/do-work-toolbox/docs/ai-report-guide.md:124`: one clause in the preset list names `relation-map` (for example "and `relation-map` draws the report's items and the relations it states between them as one zoomable map"). Change nothing else in that paragraph.
3. Run the test: it passes. Run the GREEN probe: it still fails at check 5 (no real page yet), which is expected after this commit.

**Commit 3 (the real page, no shipped change).** Commit 2 must be committed first, because the launcher refuses to start when Git shows a change outside the bundle.
1. From the worktree root, run the Session sequence of `ai-report.md` on the D-01 bundle with `--preset relation-map`. You are the editor: write the staging page by following the preset text exactly as the prompt shows it.
2. Before finalize, serve the bundle folder with `python3 -I -m http.server <free port> --bind 127.0.0.1` (the browser tool blocks `file:` URLs), and measure the staging page with a browser script at 1440x1000 and 390x844, light theme. Measure each D-04 check. (1) Relation words in view at load, with no click. (2) Labels that cover a box, in the first view, after "Show everything" plus a click, and after one click on a phone. (3) Words broken inside a label. (4) A click that leaves the clicked box or a related shown box outside the frame. (5) Tab to every box at a high zoom: the focused box is fully in the frame. (6) On the phone: the first view's titles are readable with no sideways pan, there is no page overflow, and relation-table rows need no sideways swipe. (7) With JavaScript off: words on the lines, the tool buttons hidden, every box as text, every relation as a table row. Fix the page until each check holds, then finalize. The finalize judge must report `pass`.
3. Measure the published page again with the same script after finalize, because the launcher adds its direction block at the top. Record the direction block's height at both widths.
4. Stop the server. Commit the launcher's outputs by exact path: the new `index-v1c03-*.html`, its `v1c03-*-reply.md`, and the bundle's `index.html`. Do not touch `ai-reports/catalog.json` or `ai-reports/index.html` (D-10, the integrator's job).
5. If the preset text turns out wrong or unclear while you edit, do not finalize a page the preset did not produce. Record the problem as a decision. Remove only the staging files your prepare call created (the prepare result names them). Fix the preset in a new commit on your branch, then prepare again.

**Not yours:** `VERSION`, both `CHANGELOG.md` files (the integrator writes the release, a minor bump per D-09); anything under `do-work/`; the catalog pair; the loader and every other line of `report_second_opinion.go`; `ai-report-second-opinion-prompt.md`; the other twelve presets; `skeleton.html`, `presentation.css`, `presentation.js`; `ai-report.md`, `help.md`; every existing file in any `ai-reports/` bundle except the iteration list the launcher rebuilds.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them. One is already known: the report kit's closing band does not suggest the new preset (D-07).
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-727]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Browser captures go to your own scratch folder and are deleted before hand-back. The browser tool writes `.playwright-mcp/` files into the main tree on navigation: delete the files your session created there by name, and nothing else.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded. `report_second_opinion_test.go` already runs close to 30 s (REQ-705's lesson).

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l internal/toolboxcommands/report_second_opinion_test.go` prints nothing; `go vet ./internal/toolboxcommands/`; `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/` passes (the whole second-opinion family; record its time).
2. `bash _dev/tests/shipped-package-reference-contract.sh`, `bash _dev/tests/contract-regressions.sh`, and with `DO_WORK_MAINTAINER_TIER=heavy`: `bash _dev/tests/staged-skills-contract.sh` (about 35 s).
3. `git diff --stat main...HEAD -- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md skills/do-work-toolbox/docs/ai-report-kit/` prints nothing. Running `git diff --stat main...HEAD -- skills/do-work-toolbox/docs/second-opinion-presets/` shows only `relation-map.md`.
4. A standalone re-judge of the published page: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1 --page <the new page's file name> --out <your scratch folder>/judge` reports `pass`. Delete the output folder afterwards.
5. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-727-probe.sh` run with your worktree as the current directory prints `REQ-727 probe PASS (...)`. It fails at base by design: the preset file is missing, so its first check stops it.

## Hand-back file content
The first line under the title is the real page's path from the repository root, so the maintainer can open it. Then use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the test failing after commit 1 and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result. List the seven D-04 measurements at both widths, for the staging page and for the published page, with numbers (for example "0 labels on boxes in 12 views"). Give the preset's byte count, the direction block height, and the finalize and standalone judge verdicts. For a rerun under load, give both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-13 onward (D-01 to D-12 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — what you changed in each shared file below, with line ranges after your edit.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Relation-Map Preset for ai-report Second Opinion"; check that no earlier entry uses it) and the line "What changes for an installed copy". There, `ai-report second-opinion --preset relation-map` writes a zoomable map of the report's items and their relations, the guide lists it, and nothing else changes.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
No sibling builds in this pass, so these are files other work last changed, not a parallel merge:
| Shared file | REQ-727 edits | Last changed by |
| --- | --- | --- |
| `skills/do-work-toolbox/docs/ai-report-guide.md` | one clause in the line-124 paragraph | REQ-728, REQ-729, REQ-730, REQ-731 (judge `--page`, relative `--out`, catalog lines); keep their lines as they are |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` | the function comment at line 1004 and lines 1025-1028 only (comment, list, error text) | REQ-712 (the named list) |
| `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` | the launcher's iteration list only | earlier second-opinion runs v1c01 and v1c02 |

## Notes for the integrator (the builder does nothing about these)
- **UR-176 closes with this REQ.** Its other six members are archived flat (REQ-726, REQ-728 to REQ-732). Finalization moves them into `do-work/archive/UR-176/`. Run the canonical-link sweep from `_dev/primes/lessons-releases.md` at finalize. On 2026-11-15 no shipped file linked those flat paths (D-12).
- **Release:** minor bump, 0.323.0 to 0.324.0 at today's HEAD (D-09). Root `CHANGELOG.md` is copied byte for byte to `skills/do-work/CHANGELOG.md`.
- **Catalog:** rebuild `ai-reports/catalog.json` and `ai-reports/index.html` with `ai-report-index` after the merge, so the new tree page is listed (D-10, as REQ-734 D-16).
- **Scope globs:** replace the two launcher-named globs in `## Scope` with the resolved names after the merge (D-11, as REQ-734 D-17).
- **Showing the page:** put the real page's path first in the run summary for the maintainer (D-08).
- **GREEN probe:** `do-work/runs/work-2026-11-15-012821/REQ-727-probe.sh` (`-- --probe-file`). The pre-flight argv is in the coordinator's message.

Your final message: under 2500 characters. State the branch tip hash, the real page's path, the verification results, and anything the integrator must decide.

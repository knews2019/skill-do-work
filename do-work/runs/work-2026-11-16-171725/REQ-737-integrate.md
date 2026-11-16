# REQ-737 integrator brief — run work-2026-11-16-171725

- REQ: REQ-737 ('Addendum: relation-map preset keeps box titles readable after a click on a phone'; `addendum_to: REQ-727`)
- REQ path: do-work/working/REQ-737-relation-map-phone-readable-titles.md
- UR: UR-180 (do-work/user-requests/UR-180/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/
- Shared guide (read it fully first): do-work/runs/work-2026-11-16-171725/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-16-171725/REQ-737-handback.md. It has three parts: the first pass, `## Addendum` (second pass, after the maintainer's answer) and `## Addendum (pass 3)` (the rebuild on the merged REQ-738 launcher). The pass-3 numbers are the final ones.
- Operative name / branch: worktree-agent-REQ-737-relation-map-phone-titles (base 6f00ab48; the branch merged main at ab5e07c1 in pass 3)
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-737-relation-map-phone-titles
- Dispatch instant: 2026-11-16T17:31:30Z. The build had three passes over several hours with maintainer waits between them, so SKIP the builder-work timing event and say so.
- Route: B
- tdd: false
- Closes UR-180: YES (UR_CLOSES=1); archive path do-work/archive/UR-180/REQ-737-relation-map-phone-readable-titles.md
- Release: yes, patch, 0.324.1 to 0.324.2 (preset text only; no new command or file). Re-read VERSION right before prep-payloads.
- Members already finalized in this run: REQ-738 (the launcher's short header), released 0.324.1, flat archive `do-work/archive/REQ-738-second-opinion-direction-block-first-screen.md`.
- Set aside: none.
- Last successful integration of the run (run the Restatement Sweep): YES. Pass the reviewer the sibling REQ-738 and the parent REQ-727 (archived at `do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md`).
- Mid-run addendum text: already written in the REQ as `## Addendum (mid-run)` (the maintainer's two choices, in their exact words). The review judges the build against it.
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-16-171725/REQ-737-probe.sh
- Pre-flight was recorded by the coordinator (`## Pre-Flight` is written). A plain `advance REQ-737` may still answer `preflight` / `ADVANCE-GATE-INPUT-REQUIRED`; re-pass `advance REQ-737 --request-path do-work/working/REQ-737-relation-map-phone-readable-titles.md --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-11-16-171725/REQ-737-preflight-probe.sh` to see success. If that rewrites `do-work/working/baseline.json`, restore it (`git checkout -- do-work/working/baseline.json`), as the REQ-738 integrator did (its D-15).
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/integrate-REQ-737/
- Commit the builder's untracked run files in your step-0 run-artifacts commit: `REQ-737-handback.md`, `REQ-737-progress.log`, and `REQ-737-captures/` (PNG files, the measurement script and its JSON outputs).
- Write your final report to do-work/runs/work-2026-11-16-171725/REQ-737-integration-report.md and add it, the review file and the progress log to commit_paths.
- The branch's merge of main: the cumulative range `<pre>..<merge>` then shows only REQ-737's own changes, because main is already in `<pre>`. Check that `git diff --stat <pre>..<merge>` lists the preset, ONE proof page, its reply, and the bundle's `index.html`, and nothing from REQ-738.
- Scope: replace the two launcher-named globs in `## Scope` (`index-v1c*.html`, `v1c*-reply.md`) with the real file names of the final proof page before qualification. Pages created and deleted on the branch (v1c04, v1c05) are not in the cumulative diff and need no Scope line; say so in one sentence of the Implementation Summary.
- Catalog: after the merge, run the ai-report index command on the merged tree and commit the regenerated `ai-reports/catalog.json` and `ai-reports/index.html` pair on main as the merge follow-up. The tstack bundle's row must list the new proof page.
- UR-180 closes with this REQ. In the same finalization transaction, move the flat-archived sibling `do-work/archive/REQ-738-second-opinion-direction-block-first-screen.md` into `do-work/archive/UR-180/`, and put the flat path, its destination, `do-work/user-requests/UR-180/input.md` and its archive destination in commit_paths. If the finalizer moves closed-UR siblings itself, believe the command and follow its refusals.
- Link sweep before the manifests: `git grep -n "archive/REQ-738" -- ':!do-work/'`. One known hit: the REQ-738 lesson bullet in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` line 19. Re-point it to `archive/UR-180/` and add that file and `do-work/lessons-index.md` (refresh its row) to commit_paths.
- Reviewer (Restatement Sweep candidates): every shipped place that restates the relation-map preset's checks or the phone limit REQ-727 recorded (D-16, F1): `skills/do-work-toolbox/docs/ai-report-guide.md`, `skills/do-work-toolbox/actions/ai-report.md`, help text, the 0.324.0 CHANGELOG entry (do not rewrite released entries; note only), and the lessons files. Also tell the reviewer to open the final proof page at 390 x 844 and 1440 x 1000 and re-measure: the map frame top on the first phone screen at load (the maintainer asked for the top half, y about 422 or less), titles after a click at about 12 px or more, no relation word on an arrowhead or cut by the frame edge, the header fold works. It runs the builder's `REQ-737-captures/measure-relation-map.js` (the command is in the hand-back's Integration seams), keeps captures in its own scratch folder and deletes them after.
- Report only for the reviewer: after a phone click about 95 px of empty space stays above the column label inside the frame (coordinator's look at `v1c05-390-click-O2.png`).
- REQ-738 review findings F1, F3 and F4 are report-only discovered tasks of REQ-738; do not fold them into REQ-737.
- CHANGELOG: use the builder's proposed entry, updated for passes 2 and 3 (checks 12, 15 and 17 changed; the map starts in the top half of the first phone screen). Name files in plain words or with paths that resolve from the installed core package; run `_dev/tests/shipped-package-reference-contract.sh` before finalization.
- Never edit anything under /Users/t2/.claude/.

# REQ-738 integrator brief — run work-2026-11-16-171725

- REQ: REQ-738 ('Second-opinion direction block keeps the page content on the first phone screen')
- REQ path: do-work/working/REQ-738-second-opinion-direction-block-first-screen.md
- UR: UR-180 (do-work/user-requests/UR-180/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/
- Shared guide (read it fully first): do-work/runs/work-2026-11-16-171725/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-16-171725/REQ-738-handback.md
- Operative name / branch: worktree-agent-REQ-738-fold-second-opinion-direction (tip 6432bb22, base 6f00ab48)
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-738-fold-second-opinion-direction
- Dispatch instant (for the builder-work timing event): 2026-11-16T17:31:30Z. The hand-back landed at 17:43:57Z, long before you start, so SKIP the builder-work timing event and say so (the recorder times start-to-now and would charge the builder with the wait).
- Route: B
- tdd: true
- Closes UR-180: NO (UR_CLOSES=0); flat archive path do-work/archive/REQ-738-second-opinion-direction-block-first-screen.md
- Release: yes, patch, 0.324.0 to 0.324.1 (the launcher's header changes form; no new command, flag or file). Re-read VERSION right before prep-payloads.
- Members already finalized in this run: none. REQ-737 (the relation-map preset's phone click view) is in a second build pass on its own branch; its working REQ file, captures folder, hand-back and progress log are not yours. Leave them alone, dirty or not.
- Set aside: none.
- Last successful integration of the run (run the Restatement Sweep): NO. REQ-737 integrates after you and carries it.
- Mid-run addendum text for REQ-738: none. The maintainer saw the phone captures and approved REQ-738 as built (question dialog, 2026-11-16).
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-16-171725/REQ-738-probe.sh
- Pre-flight is recorded by the coordinator before you start (see "Base gate" below). Never re-run pre-flight.
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/integrate-REQ-738/
- Write your final report to do-work/runs/work-2026-11-16-171725/REQ-738-integration-report.md and add it, the hand-back, the review file, your progress log and the `REQ-738-captures/` PNG files to commit_paths so they are kept with the run.
- Base gate: three base runs at 6f00ab48 went red only on the 30 s per-file budget, each time on a different file, at loads 28, 11 and 6 (other projects and the builders). The coordinator's notes for the fourth run are in the pre-flight record. If your gate goes red only on that budget, rerun once on a quiet machine and record both runs.
- Budget watch: `report_second_opinion_test.go` reached 30.89 s once under load at base; the builder measured `-run '^TestSecondOpinion'` at 19.2 s before and 23.6 s after the change (the new test itself 0.00 s). If the gate goes red on that file, rerun the file alone and record its time.
- Reviewer: tell it to open a scratch copy of the tstack v1c03 page decorated by the merged launcher at 390 x 844 and 1440 x 1000, check the header height (about 87 px), the map top (y = 830 at 390), the closed fold works with keyboard and without JavaScript, the report kit's "Open all details" button leaves it alone (S3 page), and no sideways scroll. Captures go in its own scratch folder and are deleted after.
- Builder decisions D-09 and D-10 (the visible label is now "Preset: <name> (full direction)", and the folded text no longer starts with "Direction:") change visible words; check that no shipped text or test still quotes "Second-opinion preset:" or "Direction:" as the header's words (`git grep -n "Second-opinion preset:" -- skills/`).
- Discovered task from the hand-back (no side gutter for the header on pages without their own header style; same gap in 0.324.0): report only. Name it in your final message.
- CHANGELOG: use the builder's proposed entry (title `Second-Opinion Pages Fold the Full Direction Under a Short Header`), edited for truth. Name files in plain words or with paths that resolve from the installed core package; run `_dev/tests/shipped-package-reference-contract.sh` before finalization.
- Never edit anything under /Users/t2/.claude/.
- Pre-flight was recorded at 2026-11-16T19:4xZ with `advance REQ-738 --request-path do-work/working/REQ-738-second-opinion-direction-block-first-screen.md --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-11-16-171725/REQ-738-preflight-probe.sh` (success; `## Pre-Flight` written). A plain `advance REQ-738` still answers phase `preflight` / `ADVANCE-GATE-INPUT-REQUIRED`; that is known. Re-pass that exact argv to see success, then go to the hand-back merge.
- `do-work/working/baseline.json` was rewritten by the REQ-737 pre-flight after yours; it is committed in the run-artifacts commit before you start. Leave it alone.

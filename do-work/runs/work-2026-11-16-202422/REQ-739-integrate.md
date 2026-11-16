# REQ-739 integrator brief — run work-2026-11-16-202422

- REQ: REQ-739 ('the second-opinion header keeps a side margin on a phone')
- REQ path: do-work/working/REQ-739-second-opinion-header-side-margin.md
- UR: UR-181 (do-work/user-requests/UR-181/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-202422/
- Shared guide (read it fully first): do-work/runs/work-2026-11-16-202422/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-16-202422/REQ-739-handback.md
- Operative name / branch: worktree-agent-REQ-739-header-phone-side-margin (tip 6584e473, base e71b4aa2)
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-739-header-phone-side-margin
- Dispatch instant (for the builder-work timing event): 2026-11-16T20:30:25Z. The hand-back already landed before you start, so SKIP the builder-work timing event and say so.
- Route: B
- tdd: true
- Closes UR-181: YES (UR_CLOSES=1); archive path do-work/archive/UR-181/REQ-739-second-opinion-header-side-margin.md (UR-181 has only this REQ)
- Release: yes, patch, 0.324.2 to 0.324.3. Re-read VERSION right before prep-payloads.
- Members already finalized in this run: none (single-member run). Set aside: none.
- Last successful integration of the run (run the Restatement Sweep): YES. Pass the reviewer REQ-738 (archived at do-work/archive/UR-180/REQ-738-second-opinion-direction-block-first-screen.md), which shipped the header this REQ changes.
- Mid-run addendum text: none. The maintainer saw the captures (P5-v1c04-req737-final-390.png, P1-kit-page-390.png) and approved integration as built (question dialog, 2026-11-16).
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-16-202422/REQ-739-probe.sh
- Pre-flight was recorded by the coordinator (`## Pre-Flight` written). A plain `advance REQ-739` may still answer `preflight` / `ADVANCE-GATE-INPUT-REQUIRED`; re-pass `advance REQ-739 --request-path do-work/working/REQ-739-second-opinion-header-side-margin.md --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-11-16-202422/REQ-739-preflight-probe.sh` to see success. If that rewrites `do-work/working/baseline.json`, restore it with `git checkout -- do-work/working/baseline.json`.
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/integrate-REQ-739/
- Commit the builder's untracked run files in your step-0 run-artifacts commit: REQ-739-handback.md, REQ-739-progress.log, REQ-739-captures/ (PNG files, the measurement script and its raw output), and this brief.
- Write your final report to do-work/runs/work-2026-11-16-202422/REQ-739-integration-report.md and add it, the review file and the progress log to commit_paths.
- Reviewer: tell it to decorate scratch copies (never published pages) at 390 x 844 and 1440 x 1000 and check: header left text edge 16 px on a kit page and on the REQ-737 v1c04 page, the tstack v1c03 page's own header rule still wins (31 px), no sideways scroll, fold works, kit "Open all details" leaves the fold closed. Captures in its own scratch folder, deleted after.
- Known report-only items from earlier runs, not this REQ's scope: after a phone click about 88-223 px of empty space sits above the column title (REQ-737 review F1); fold-arrow hidden on pages that hide `summary` markers (REQ-738 review F1).
- CHANGELOG: use the builder's proposed entry, edited for truth. Name files in plain words or with paths that resolve from the installed core package; run `_dev/tests/shipped-package-reference-contract.sh` before finalization.
- Never edit anything under /Users/t2/.claude/.

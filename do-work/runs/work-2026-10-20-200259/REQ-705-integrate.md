# REQ-705 integrator brief — run work-2026-10-20-200259

- REQ: REQ-705 (user presets for second opinions live only inside the report bundle (the user answer to the old confirm question))
- REQ path: do-work/working/REQ-705-confirm-uncommitted-preset-refusal.md
- UR: UR-160 (do-work/user-requests/UR-160/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-20-200259/
- Shared guide (read it fully first): do-work/runs/work-2026-10-20-200259/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-10-20-200259/REQ-705-handback.md
- Operative name / branch: worktree-agent-REQ-705-bundle-local-second-opinion-presets
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-705-bundle-local-second-opinion-presets
- Dispatch instant (for the builder-work timing event): 2026-10-20T21:31:30Z
- Route: B (standard review)
- Wave membership: REQ-703, REQ-704 (serial). Set aside: none.
- Closes UR-160: no (UR_CLOSES=0; flat archive path do-work/archive/REQ-705-confirm-uncommitted-preset-refusal.md)
- Last successful integration of the wave (run the Restatement Sweep): no (REQ-704 is the last member and carries the wave-end sweep)
- Lesson satellite for a Go lesson: skills/do-work/tools/do-work-cli/lessons-do-work-cli.md (action-file lessons: _dev/primes/lessons-action-files.md)
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- Wave membership line above: the members are REQ-703 (done, 0.306.6), REQ-705 (you), REQ-704 (after you).
- Pre-flight was recorded by the coordinator; `## Pre-Flight` is in the REQ. Never re-run pre-flight.
- tdd: true. The Testing section carries RED/GREEN evidence for TestSecondOpinionPresetResolvesFromBundleOnly from the hand-back.
- RUN/REQ-703-probe.sh names a test this REQ deletes; never run it. Your probe is RUN/REQ-705-probe.sh.
- Release: VERSION is 0.306.6; this REQ ships 0.306.7 (patch). Re-read VERSION right before prep-payloads.
- The REQ-703 lesson link stays flat for now; REQ-704's integrator re-points it when UR-160 closes. If you add a lesson bullet, link your flat archive path the same way.

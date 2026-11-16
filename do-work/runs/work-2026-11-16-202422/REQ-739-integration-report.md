# REQ-739 integration report (run work-2026-11-16-202422)

REQ-739 (the second-opinion header keeps a 16 px side margin on a phone) is done and released as 0.324.3. It closes UR-181.

- Archive: `do-work/archive/UR-181/REQ-739-second-opinion-header-side-margin.md`. The UR-181 input is archived beside it.
- Commits: run artifacts `540925d4` (pre), merge `ef5dc795`, finalization `9bd69ce6` (completed, cleanup_complete, no blocked paths, no reason codes), and this report `[REQ-739] integration report`.
- Gate at the merge: passed in 128 s at a load of 3.11. The GREEN probe passed.
- Review: 98%, Pass (Route B, Restatement Sweep over REQ-738). It found no stale restatement and 0 anti-bloat items. F1 (Nit, impact-negligible → report only): the launcher's `<style>` sits in the body, which strict HTML does not allow. D-01 chose this on purpose. The reviewer's Chromium numbers match the builder's: 16 px on the kit and v1c04 pages, 31 px on v1c03 (that page's own rule wins), and no sideways scroll.
- Heavy lanes: all 4 executed and passed. do-work-cli-integrations took 65 s, staged-skills 40 s, updater 71 s and installer 34 s.
- Timing: the handback-merge, 2 verification-gate and review events were folded into `## Timing` (14m 21s). The builder-work event was skipped because the hand-back landed before I started (D-14).
- Lesson: one `writer-anchored-marker-position` bullet in `lessons-do-work-cli.md`, with its index row refreshed. `contract-regressions.sh` and `shipped-package-reference-contract.sh` pass after finalization.
- Cleanup: the builder worktree and branch were removed without force. The drain-head checkout was left for you.
- Decisions: D-12 committed the hand-back in the step-0 commit, as the brief asked. D-13 copied the measurement script and output into the captures folder. D-15 committed this report after finalization, because it names that commit. D-16 picked the patch release, with the CHANGELOG edited for truth (about 24 px on pages that keep the default body margin).
- Left dirty: nothing.
- Report-only discovered task: the published tstack v1c03 page still shows the 0.324.0 unfolded header. A new pass rebuilds it.

Next: the coordinator runs `advance --checkpoint` and the run cleanup.

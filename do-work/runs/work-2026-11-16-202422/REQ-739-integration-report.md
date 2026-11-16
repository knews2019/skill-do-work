# REQ-739 integration report (run work-2026-11-16-202422)

REQ-739 (the second-opinion header keeps a 16 px side margin on a phone) is done and released as 0.324.3. It closes UR-181.

- Archive: `do-work/archive/UR-181/REQ-739-second-opinion-header-side-margin.md`. UR-181 input is archived at `do-work/archive/UR-181/input.md`.
- Commits: run artifacts `540925d4` (pre). Merge `ef5dc795` (ef5dc79566d7bf7f074f2d1ce3ed13488720cacc). Finalization `9bd69ce6` (9bd69ce6819a1d93f7a891fa8b57cad1d5252640): status completed, phase cleanup_complete, no blocked paths, no reason codes.
- Gate: `maintainer-verify.sh` at the merge passed in 128 s, at a 1-minute load of 3.11. The slowest file was `report_second_opinion_test.go` at 16.33 s. The GREEN probe passed.
- Review: 98%, Pass (Route B, with the Restatement Sweep over REQ-738). It found no stale restatement and 0 anti-bloat items. F1 (Nit, impact-negligible → report only): the launcher's `<style>` sits in the body, which strict HTML does not allow. Browsers apply it, and REQ-739's D-01 chose it on purpose. The reviewer measured the same numbers in Chromium: 16 px on the kit and v1c04 pages, 31 px and 43.5 px on v1c03 (that page's own rule wins), no sideways scroll, and the header fold stays closed after "Open all details".
- Heavy lanes at ef5dc795: all 4 ran and passed (do-work-cli-integrations 65 s, staged-skills 40 s, updater 71 s, installer 34 s). None was skipped or reused.
- Timing events: handback-merge, verification-gate ×2 (the gate with the probe, and the heavy drain), and review, folded into `## Timing` (14m 21s). The builder-work event was skipped because the hand-back landed before the integrator started (D-14).
- Lesson: one bullet added in the `writer-anchored-marker-position` family of `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`. The `do-work/lessons-index.md` token count was refreshed to 25210. `contract-regressions.sh` and `shipped-package-reference-contract.sh` pass after finalization.
- Cleanup: the builder worktree was removed and its branch deleted with `-d`, and the worktree list was pruned. The shared drain-head checkout was left for the coordinator.
- Decisions: D-12 committed the hand-back in the step-0 commit, as the brief asked. D-13 copied `measure.js` and `measure-output.txt` into `REQ-739-captures/`. D-15 commits this report alone after finalization, because it names the finalization commit. D-16 picked the patch release and edited the CHANGELOG entry for truth (about 24 px on pages that keep the browser's default body margin).
- Left dirty under ROOT: nothing. This report is committed by exact path in `[REQ-739] integration report`.
- Report-only discovered task: the published tstack v1c03 page still shows the 0.324.0 unfolded header (1326.7 px tall at 390 px). By design, published pages are rebuilt only by a new pass.

Next step for the coordinator: run `advance --checkpoint` and the run cleanup. Nothing else is pending for REQ-739.

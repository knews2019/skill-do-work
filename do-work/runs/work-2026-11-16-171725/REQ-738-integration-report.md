REQ-738 (the second-opinion launcher shows a short header and folds the full direction) is finalized and released as **0.324.1**.

- **Archive:** `do-work/archive/REQ-738-second-opinion-direction-block-first-screen.md` (flat; UR-180 stays open for REQ-737).
- **Commits:** run artifacts f995337a (pre), merge 943c43858982436dcdae54ec9ccc80caed0e5468, finalization ad8842ea61489adfffd32320b2e65b6f3a0a1e7b (no metadata commit), then this report as one plain commit.
- **Gate:** green on the first run at 943c4385 (129 s, load 1.95, slowest Go file 17.65 s). GREEN probe passed under advance.
- **Review:** Approve, 96%, acceptance Pass (implementation and integration stages, browser-measured: header 87.4 px, map top y = 830 at 390 x 844). Findings, all report only: F1 a page rule that hides summary markers would also hide the header's ▸; F2 a lesson line stated the old header size (fixed in the release commit, D-16); F3 no side gutter on kit and plain pages (pre-existing, same as the discovered task); F4 the test does not pin the kit-cancelling inline styles. No fix re-merge.
- **Heavy lanes at 943c4385, all executed, exit 0:** do-work-cli-integrations 61 s, staged-skills 36 s, updater 66 s, installer 26 s.
- **Timing events:** handback-merge, verification-gate (gate), review, verification-gate (heavy drain); folded into `## Timing`. Builder-work event skipped per the brief.
- **Lesson:** one `[family: session-is-the-backend]` bullet beside its siblings in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`; lessons-index row refreshed (24908 tokens). `contract-regressions.sh` and `shipped-package-reference-contract.sh` pass after finalization.
- **Worktree:** builder worktree removed, branch `worktree-agent-REQ-738-fold-second-opinion-direction` deleted, worktrees pruned.
- **Decisions:** D-15 restored `do-work/working/baseline.json` after the pre-flight re-pass rewrote it; D-16 fixed review F2 in the release commit; D-17 commits this report after finalization because it names the finalization commit.
- **Left dirty under ROOT:** only REQ-737's files (`REQ-737-handback.md`, `REQ-737-captures/`, `REQ-737-progress.log`), REQ-737's builder's, untouched.
- **Report-only discovered task:** on a page without its own header rule (report kit, custom page), the header text has no side gutter and starts at x = 0 on a phone; 0.324.0 has the same gap. A fix needs its own decision: an inline padding would override the tstack page's own rule.

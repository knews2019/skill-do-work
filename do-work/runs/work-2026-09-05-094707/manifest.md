# Queue run

Status: in-progress

User scope: `do-work run REQ-507` (hand the archive and commit tails to finalize), serial.

- REQ-507: claimed at 12d264c2; canonical gate red twice on ShellCheck SC2148 in another session's committed probe script; deferred through `defer-gate` to repair REQ-584 (commit cdd19977). Saved implementation range 8e3dbf01..ad8bceb7 retained for resume.
- REQ-584: claimed at c5ee4a8c; Route A; already-green repair no-op (the probe was given a shebang by another session in 4d47c821; a second unrelated red — stale archive links in lessons-do-work-cli.md — was repaired by another session in 9f914188). Direct gate green at f9659f0f; green-gate record recorded; validator tdd_allowed/review_allowed true. Independent review → REQ-584-review.md pending.
- Foreign dirt left alone and never staged: do-work/working/REQ-574-*.md (another session's live claim), ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html, .playwright-cli/, output/.

- REQ-584 completed and archived (commit aea9d619, metadata 14436f51; no release, already-green no-op). Builder worktree for it was unused and removed.
- REQ-507 re-claimed at 1012e5e2 (explicit target). Saved-range resume proof: drift on 8 of 12 paths → reuse rejected, pair deleted, prior evidence demoted to "Prior Evidence". Direct canonical gate green at 1012e5e2 (recorded green-gate); preflight baseline green (lifecycleadvance, finalization, resultmodel). Foreign do-work/working/baseline.json preserved and restored.
- REQ-507 builder build_507 dispatched at 2026-09-05T10:10:08Z on worktree-agent-REQ-507-hand-archive-and-commit-tails-to-finalize (.git/work-run-20260905/); expected artifact REQ-507-handback.md.

- REQ-507 remediation builder returned no commits (all four acceptance criteria hold on current main; hand-back promoted into the record and removed). Qualification over 8e3dbf01..ad8bceb7 satisfied; focused probe green; direct canonical gate green (run at 5cb094dd, recorded at e6005c0e, do-work-only drift between); independent review 96% Pass, five report-only findings; six heavy lanes executed green at c6f3c77d.
- REQ-507 completed and delivered: release 0.294.0 "Closing a Request Is One Command", commit d3be4122 (supplied provenance ad8bceb7), archive do-work/archive/UR-098/REQ-507-hand-archive-and-commit-tails-to-finalize.md; UR-098 closed and archived with all members.
- Observed collision: the sibling session's commit 985fa736 swept this run's uncommitted REQ-507 working edits into its own "[work run] REQ-574 owner evidence" commit. Nothing was lost; noted for the maintainer.

Status: complete (targeted run REQ-507; repair REQ-584 folded in).

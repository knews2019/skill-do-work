# Queue run

Status: in-progress

User scope: `do-work run REQ-507` (hand the archive and commit tails to finalize), serial.

- REQ-507: claimed at 12d264c2; canonical gate red twice on ShellCheck SC2148 in another session's committed probe script; deferred through `defer-gate` to repair REQ-584 (commit cdd19977). Saved implementation range 8e3dbf01..ad8bceb7 retained for resume.
- REQ-584: claimed at c5ee4a8c; Route A; already-green repair no-op (the probe was given a shebang by another session in 4d47c821; a second unrelated red — stale archive links in lessons-do-work-cli.md — was repaired by another session in 9f914188). Direct gate green at f9659f0f; green-gate record recorded; validator tdd_allowed/review_allowed true. Independent review → REQ-584-review.md pending.
- Foreign dirt left alone and never staged: do-work/working/REQ-574-*.md (another session's live claim), ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html, .playwright-cli/, output/.

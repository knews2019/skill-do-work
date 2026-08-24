# REQ-354 builder brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-354-open-the-detail-drawer-from-a-durations-mark`

Exploration is accepted. Implement test-first only in:

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`

Reuse `openDetail("req", id)` without editing `board-detail.js`. Factor one jitter-aware nearest-mark
selection over `markIndex` for hover and overlay click. Make every circle a semantic SVG button with
one roving Tab stop across the complete projected set; Left/Right reach every mark and Enter/Space
open the same drawer. Preserve hover text and all REQ-351/352 behavior.

Add RED/GREEN browser coverage for sole-tab-stop semantics, exhaustive arrow reach, Enter/Space,
and a trusted CDP click through the overlay at a busy-day jittered mark. The fixture must prove raw
and jitter-aware targeting would differ, then assert the drawer/readout use the jittered mark.

Run focused tests with the retained Chromium, the full module/canonical gates, and inspect responsive
generated-board behavior. Commit on the worktree branch and write handback only to
`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-174255/REQ-354-handback.md`.

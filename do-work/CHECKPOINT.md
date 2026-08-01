---
session_ended: 2026-08-01T14:00:00Z
last_completed: REQ-070
queue_state: 0 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 4
session_depth: moderate
---

# Session Checkpoint

## Completed This Session

- REQ-067: Target ID Resolution contract + UR ids in `do-work run` (Route B, v0.159.0) — commit `1e653bc`
- REQ-068: UR ids in `do-work abandon` / `reserve` / `release` (Route B, v0.160.0) — commit `180e523`
- REQ-069: Adopt exclusive-session model, remove concurrency machinery (Route C, v0.161.0) — commit `76cdf39`
- REQ-070: REQ ids in `do-work roadmap` (Route B, v0.162.0) — commit `03349a8`

Plus a PR-review fix commit (v0.161.1, `24950ff`): reserve mode-table UR-token
routing + `do-work run` dual-provenance precedence (Codex findings on REQ-067/068).

## Still Queued

- **Queue empty.** No pending, pending-answers, blocked, reserved, or in-progress REQs.
- **UR-011 closed** — REQ-067/068/070 consolidated into `do-work/archive/UR-011/`.
- **UR-012 closed** — REQ-069 consolidated into `do-work/archive/UR-012/`.

## Session Notes

- **All four REQs are skill-instruction maintenance passes** (`maintenance: true`), landed on
  branch `claude/do-work-run-u77g5u` (PR #126). Each is `tdd: true`, proven RED→GREEN by a
  `_dev/tests/contract-regressions.sh` assertion; the suite is the mechanical gate for every one.
- **Pre-existing test baseline:** `_dev/tests/update-script-behavior.sh` fails on the untouched
  tree in this sandbox (mid-update / dirty-install recovery probes) — an environment quirk,
  reproduced before any change, excluded from every REQ's pass/fail gate. `tools/do-work-update.sh`
  was never touched.
- **REQ-069 was the large one:** removed ~8,200 words across `work.md` / `work-reference.md` /
  `cleanup.md` (three-file total 38,837 → 30,637). The dangling-reference sweep reached six files
  beyond the declared `write_set` (board, capture-reference, work-guide, board-guide, kanban prime,
  clear-questions) — recorded as decision D-02 on the REQ. `tools/queue-kanban/model.go` was left
  untouched per constraint; `write_set` survives as a display-only field for the board's badge.
- **User steer mid-session:** "read-only actions can run in parallel" — folded into REQ-069's
  Execution Model rule (the exclusive-session boundary governs writers only).

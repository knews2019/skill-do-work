---
id: REQ-650
title: 'The disk probe measures the repo root only'
status: pending
created_at: 2026-10-09T16:06:43Z
user_request: UR-143
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-651, REQ-652, REQ-653]
batch: october-review-triage
write_set: ["skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/disk_space_unix.go", "skills/do-work-board/tools/queue-kanban/disk_space_windows.go", "skills/do-work-board/tools/queue-kanban/disk_space_test.go", "skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
---
# The Disk Probe Measures the Repo Root Only
## What
The `low-disk-space` verify probe measures the filesystem holding the repo root and then every `worktree-agent-*` worktree, deduplicating by device id so two directories on one disk report once. Reduce it to one measurement of the repo root. Delete the worktree loop, the device-id dedupe and their tests; keep the thresholds, the single finding, and the Testing-page line.
## Why
Consumer review (validate-feedback 2026-10-09, F8, accepted): "the disk space check is over-engineered; it treats the computer like a complex network of different volumes when a simple check of the project folder would have sufficed." The earning incident in REQ-625 (the low-disk-space probe) §Why was growth inside the repo: browser-QA screenshots grew the repo to about 20 GB and the disk nearly filled. A repo-root-only check catches that. No archived document names an incident where a worktree on another volume filled up; the multi-device branch rests only on "worktrees can sit on another disk than the repo root" (`verify.go:216-217`). The maintainer chose "Repo root only".
## Verified Facts (from triage)
- `skills/do-work-board/tools/queue-kanban/verify.go:279-290` measures the repo root first, then each worktree path in name order; the worktree paths come from the shared `git worktree list --porcelain` read (`verify.go:217-223`, `1491-1510`).
- `verify.go:313-316` skips a device already reported through `measuredDevices[measurement.deviceIdentity]`; the comments at `:238-241` and `:265-269` say "so two directories on one disk are reported once" and "One finding per device".
- Each path costs one `Statfs` plus one `os.Stat` for the device id (`disk_space_unix.go:19-44`).
- The Testing page shows only the repo root's reading (`verify.go:293-302`, `generate.go:709-720`), so the Testing line does not change.
- REQ-625's review noted there was no test for "a worktree on the repo root's device, the common layout" (`REQ-625.md:241`).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:23` records the incident under family `disk-space-blind-spot`; it stays.
## Detailed Requirements
1. `verify.go`: the probe measures exactly one path, the repo root. Remove the worktree-path loop, the `measuredDevices` map, the device-identity field on the measurement type if nothing else reads it, and the comments that describe per-device reporting. If the worktree read at `:216-223` feeds only this probe, remove that call too; if the VERIFY band's other probes share it, leave the shared read alone.
2. `disk_space_unix.go` and `disk_space_windows.go`: drop the `os.Stat` device-identity read when nothing else uses it. Keep `syscall.Statfs` and `GetDiskFreeSpaceExW`, the build tags, and `disk_space_unsupported.go`.
3. `disk_space_test.go`: remove the dedupe and multi-worktree cases. Keep the threshold cases, the unsupported-platform case and the root-reading case. Add the RED test below.
4. Thresholds (warn below 10 GiB, critical below 3 GiB), `diskSpaceLevelFor`, the finding's Detail and Remedy text, the `.` path reduction and the Testing-page payload are unchanged.
5. Docs: `skills/do-work-board/docs/board-guide.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` no longer say "one finding per device" or "each builder worktree"; they say the probe measures the repo root.
6. Changelog entry (written by the integrator at release) says why: the earning incident was growth inside the repo, and the multi-device branch had no incident behind it.
## Constraints
- Deletion only. No new flag, no new finding, no change to thresholds or to what the run loop reads (it reads nothing from this probe).
- Read `_dev/primes/prime-kanban-board.md` and `prime-do-kanban.md` before editing; the board's versioning and build rules apply.
- Browser lane: set `QUEUE_KANBAN_BROWSER` when running the heavy lane; a skipped lane is not a pass.
## Dependencies
None. May run beside REQ-651, REQ-652 and REQ-653.
## Builder Guidance
Certainty is high: the triage located every line. Latitude: whether the measurement type keeps a device field for the Windows helper's shape, and how the root-only test is phrased.
## Red-Green Proof
**RED prompt/case:** A queue-kanban test builds a repo with one `worktree-agent-REQ-1-*` worktree on a second device and runs the disk probe with a fake measurer that records the paths it was asked to measure.
**Why RED now:** The measurer is called twice (repo root and the worktree) and the probe emits one finding per device; `grep -n measuredDevices verify.go` finds the dedupe map.
**GREEN when:** The measurer is called exactly once, with the repo root; `grep -n measuredDevices skills/do-work-board/tools/queue-kanban/verify.go` prints nothing; `go test ./...` in `skills/do-work-board/tools/queue-kanban/` is green; `queue-kanban verify` on this repo prints at most one disk finding.
**Validation:** User confirmed. The maintainer chose "Repo root only (Recommended)" in the triage questions and approved the plan.
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8110 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: family `disk-space-blind-spot` is this probe's own incident, and the satellite governs the queue-kanban model this REQ edits.
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over budget; `slugged: partial`). Matching reason: its owning prime governs the board tool this REQ edits.
## Full Context
See `do-work/user-requests/UR-143/input.md` for complete verbatim input (both pastes, the maintainer's answers, and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer review "the disk space check is over-engineered; it treats the computer like a complex network of different volumes when a simple check of the project folder would have sufficed", accepted by `do-work-toolbox validate-feedback` on 2026-10-09 as F8; maintainer answer "Repo root only (Recommended)".*

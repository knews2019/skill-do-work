# Builder brief — REQ-650 (The disk probe measures the repo root only)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-650-the-disk-probe-measures-the-repo-root-only
- Branch: worktree-agent-REQ-650-the-disk-probe-measures-the-repo-root-only, created from main HEAD 2eb14357 with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-650-the-disk-probe-measures-the-repo-root-only.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-6, Constraints, Red-Green Proof, Required Lessons — Dropped for Budget. Requirement 6 (changelog) belongs to the integrator; you PROPOSE the changelog entry and a lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-143/input.md (the review, the maintainer's answers, the triage; read the F8 verdict and the REQ-625 incident).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-09-160815/REQ-650-handback.md
- Route A, tdd: true, impact-user-visible, effort-mechanical, domain backend. Go deletion plus docs.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-kanban-board.md and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (read both before touching the board: versioning, parser lock-step, build outputs, traps) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release is the integrator's). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md read every bullet carrying `[family: disk-space-blind-spot]` and the Traps section of prime-do-kanban.md.

## The change (decided at capture and in the approved plan, do not reopen)
In `skills/do-work-board/tools/queue-kanban/` of your worktree:
1. `verify.go`: the `low-disk-space` probe measures exactly one path, the repo root. Delete the worktree-path loop (around lines 279-290), the `measuredDevices` map and its skip (around 313-316), and the comments describing per-device reporting (around 238-241 and 265-269). The worktree read around 216-223 stays if any other probe shares it (check its callers); delete it only if the disk probe was its only reader. Keep the thresholds (warn below 10 GiB, critical below 3 GiB), `diskSpaceLevelFor`, the finding's Detail and Remedy text, the `.` path reduction, and the Testing-page payload in `generate.go` (lines ~126-135 and ~709-720; do not edit generate.go unless a signature you removed forces it).
2. `disk_space_unix.go` and `disk_space_windows.go`: drop the `os.Stat` device-identity read and the device field if nothing else reads it. Keep `syscall.Statfs`, `GetDiskFreeSpaceExW`, the build tags and `disk_space_unsupported.go`.
3. `disk_space_test.go`: remove the dedupe and multi-worktree cases; keep threshold, unsupported-platform and root-reading cases. Write the RED test first: a repo with one `worktree-agent-REQ-1-*` worktree (or a stubbed worktree list) and a fake measurer that records the paths it was asked to measure; it fails today (two calls, one finding per device) and passes when the measurer is called exactly once with the repo root. Record the failing run and the passing run in the hand-back.
4. Docs: `skills/do-work-board/docs/board-guide.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` no longer say "one finding per device" or "each builder worktree"; they say the probe measures the repo root. `lessons-do-kanban.md` line ~23 (the incident) stays.
- Deletion only: no new flag, finding, threshold, or run-loop read. If prime-do-kanban.md or the board's versioning rule asks for a version/comment update in the Go tool, follow it.

## Write boundary
Exactly: `verify.go`, `disk_space_unix.go`, `disk_space_windows.go`, `disk_space_test.go`, `skills/do-work-board/docs/board-guide.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`; plus `generate.go` or `generate_test.go` only if a removed symbol forces it (say so in the hand-back under Decisions). Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `cd skills/do-work-board/tools/queue-kanban && go vet ./ && go test ./` (about 45 s; browser tests skip without a browser, say so).
- `grep -n measuredDevices verify.go` prints nothing; `grep -rn "per device\|each builder worktree" skills/do-work-board/docs/board-guide.md skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` prints nothing about the disk probe.
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: 2eb14357.
- File manifest: each file with (modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, go vet/go test exit codes and wall time, and each file checked.
- Red-green evidence: the new test's name, its failure output before the change, its pass after.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` in that file's bullet shape (`- [family: <slug>] [REQ-650: <one-line lesson>](<relative archive link, the integrator fixes the path>)`), and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words, one or two sentences on why: the earning incident was growth inside the repo and the multi-device branch had no incident, then specific bullets). The integrator writes both.
- Integration seams (none expected). Test wall times.

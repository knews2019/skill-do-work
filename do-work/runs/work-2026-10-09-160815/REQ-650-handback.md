# Hand-back — REQ-650 (The disk probe measures the repo root only)

**Branch:** worktree-agent-REQ-650-the-disk-probe-measures-the-repo-root-only
**Worktree:** /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-650-the-disk-probe-measures-the-repo-root-only
**Base:** 2eb14357  **Head:** b504b4a3 (single commit)

## File manifest

All under `skills/do-work-board/tools/queue-kanban/`.

- `verify.go` (modified): `appendDiskSpaceFindings(report, repoRoot, measure)` measures the repo root once. Deleted the worktree-path loop, the `measuredDevices` map and skip, the `deviceIdentity` field, the "disk-space probe for worktrees" skip line, and the per-device comments. The shared `gitState` worktree read stays: `appendWorktreeFindings` reads it.
- `disk_space_unix.go` (modified): dropped the `os.Stat` device-id read and the `os` import. `syscall.Statfs`, Bavail handling and build tags unchanged.
- `disk_space_windows.go` (modified): dropped the FNV volume-name hash, the `hash/fnv` and `strings` imports and the comment about it. `GetDiskFreeSpaceExW` unchanged.
- `disk_space_test.go` (modified): removed `TestDiskSpaceProbeReportsOneFindingPerDevice` and `TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory`. `freeOnDevice` became `freeOnDisk`. `TestDiskSpaceProbeKeepsTheRepoRootSkipReason` and the unsupported-platform test lost their worktree halves. Added `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot`.
- `generate_test.go` (modified, forced): `TestMain` and `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` set the removed field. The payload test seeded a worktree disk finding through the old four-argument call. It now seeds a second root finding and asserts `.` appears twice and no other disk Subject appears.
- `prime-do-kanban.md` (modified): the disk-space trap now says the probe measures the filesystem holding the repo root, the only path it measures. The worktree-Subject sentence is gone.
- `skills/do-work-board/docs/board-guide.md`: not changed. See D-02.
- `disk_space_unsupported.go`, `generate.go`, `lessons-do-kanban.md`: not changed.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the brief, the REQ, UR-143 F8, the crew rules, both board primes, prime-releases and the `disk-space-blind-spot` lesson. `gitState.worktreePathsByName` is also read by `appendWorktreeFindings`, so the shared read stays and only the disk probe's local use goes. The `deviceIdentity` field has no reader once the dedupe goes, so it is removed on both platforms. The signature drops the worktree map. The RED test goes through `runVerifyProbes` with a real git repo and one `worktree-agent-REQ-1-fixture` worktree, so it compiles against both the old and new code and fails on the assertion. No lesson contradicts the plan.
- [x] **[APPLY]:** RED test written and run first (failing output below). Then verify.go, the two platform files, the test files and the prime, exactly as planned. No new flag, finding, threshold or run-loop read.
- [x] **[UNIFY]:**
  - `git diff --stat 2eb14357..b504b4a3`: 6 files changed, 108 insertions(+), 200 deletions(-). `git diff --check` clean.
  - `go vet ./` exit 0. Cross-OS vet with GOOS=windows, linux and netbsd exit 0. The netbsd run covers the unsupported build.
  - `go test -count=1 ./` exit 0, wall 78 s (`ok ... 76.132s`). Browser tests skipped because `QUEUE_KANBAN_BROWSER` was not set. A skipped browser lane is not a pass, but this change touches no `web/` asset or rendered page.
  - `gofmt -l .` prints nothing.
  - `grep -n measuredDevices verify.go` prints nothing. `grep -rn "per device\|each builder worktree"` over board-guide.md and prime-do-kanban.md prints nothing.
  - Live CLI: a binary built into the session scratchpad ran `verify --repo-root <main tree>`. It printed zero `low-disk-space` findings with 48 GiB free, which meets "at most one disk finding". Exit 1 came from the in-flight builder worktree findings.
  - Files checked: verify.go (deletion and comments), disk_space_unix.go, disk_space_windows.go (imports, cross-vet), disk_space_test.go, generate_test.go, prime-do-kanban.md. The diff holds no debug artifacts and no build outputs.

## Red-green evidence

Test: `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot` (`disk_space_test.go`).

RED at base 2eb14357 with only the test added:

```
--- FAIL: TestDiskSpaceProbeMeasuresOnlyTheRepoRoot (0.19s)
    disk_space_test.go:321: measured ["/var/folders/.../TestDiskSpaceProbeMeasuresOnlyTheRepoRoot1811141353/001" "/private/var/folders/.../TestDiskSpaceProbeMeasuresOnlyTheRepoRoot1811141353/002/worktree-agent-REQ-1-fixture"], want exactly the repo root ["/var/folders/.../TestDiskSpaceProbeMeasuresOnlyTheRepoRoot1811141353/001"]
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.538s
```

GREEN at b504b4a3:

```
--- PASS: TestDiskSpaceProbeMeasuresOnlyTheRepoRoot (0.20s)
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.929s   (disk and payload tests, -run filter)
ok  	github.com/knews2019/skill-do-work/queue-kanban	76.132s  (full suite)
```

## Decisions

- **D-01 DECIDE & STATE:** `generate_test.go` was edited, outside the write_set. The removed `deviceIdentity` field and the old four-argument signature forced it. The payload test still requires two low-disk-space findings reduced to `.`, so it still proves the reduction.
- **D-02 DECIDE & STATE:** `board-guide.md` was left unchanged. It never said "one finding per device" or "each builder worktree", and its only disk line already says the toolbar shows "the repo root's free disk space". Adding text would go against the deletion-only constraint.
- **D-03 DECIDE & STATE:** the measurement type keeps only free and total bytes. The REQ allowed keeping a device field for the Windows helper's shape, but nothing reads it. Removing it also removes the FNV hash and its known limit for folder-mounted volumes.
- **D-04 DECIDE & STATE:** the "disk-space probe for worktrees: <error>" skip line is gone along with the loop. A worktree-list failure is still reported by the worktree probes' own "worktree probes: ..." skip.
- **D-05 DECIDE & STATE:** no version bump or comment change in the Go tool. Under prime-kanban-board.md the tool rides the suite version and the root changelog, so the integrator owns that.

## Discovered Tasks

- `skills/do-work/CHANGELOG.md` entry for 0.305.60 still says the probe measures "the repo root and each builder worktree, one finding per device". It is a historical entry and should stay as written. impact-cosmetic → report only

## Lessons read

- Satellite `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, the `[family: disk-space-blind-spot]` bullet (line 23, the 0.305.60 incident). Left unchanged.
- `prime-do-kanban.md` Traps section in full. `_dev/primes/prime-kanban-board.md` Conventions and Traps.
- `_dev/primes/lessons-kanban-board.md` has no `disk-space-blind-spot` bullet. The REQ dropped it whole for budget, and it was not read.

## Proposed lesson bullet (integrator writes it)

```
- [family: disk-space-blind-spot] [REQ-650: a probe watches where its earning incident happened; the repo-root disk probe had grown a per-worktree, per-device branch that no incident ever needed, so it now measures the repo root only](../../../../do-work/archive/UR-143/REQ-650-the-disk-probe-measures-the-repo-root-only.md)
```

This is the second bullet in the family. Under general.md it promotes one generalized `[family: disk-space-blind-spot]` line into the prime's Traps. That is the integrator's call, together with the `do-work/lessons-index.md` row refresh.

## Proposed CHANGELOG entry (integrator writes it)

```
## The disk-space check measures the project folder only

The board's low-disk-space check measured the repo root and every builder worktree and merged readings by disk. The problem that earned the check was the repo itself growing to about 20 GB, and no problem ever came from a worktree on another disk, so the extra machinery is gone.

- `queue-kanban verify`'s `low-disk-space` probe measures one path, the filesystem holding the repo root. Thresholds are unchanged: a warning below 10 GiB free and critical below 3 GiB.
- The per-worktree measurements, the device-id dedupe and the "disk-space probe for worktrees" skip line are removed.
- The Testing page's disk line and the finding's text are unchanged.
```

## Integration seams

None. Neither generate.go nor the run loop changed, and no other REQ's files were touched.

## Test wall times

| Run | Wall |
|---|---|
| RED, single test | 2.3 s |
| GREEN, disk and payload tests | 2.2 s |
| `go vet ./` | under 1 s |
| `go test -count=1 ./` | 78 s |

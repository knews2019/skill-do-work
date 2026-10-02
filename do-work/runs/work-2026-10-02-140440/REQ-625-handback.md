# Hand-back — REQ-625
**Branch:** worktree-agent-REQ-625-low-disk-space-verify-probe  **Base:** 8812a445  **Head:** 2a698a96
## File manifest
- skills/do-work-board/tools/queue-kanban/verify.go (modified) — `verifyCategoryLowDiskSpace`, the two threshold constants, `diskSpaceMeasurement`, `errDiskSpaceUnsupported`, `diskSpaceMeasurer`, `appendDiskSpaceFindings`, `formatGibibytes`, and the call in `collectVerifyFindings` after `appendWorktreeFindings`
- skills/do-work-board/tools/queue-kanban/disk_space_unix.go (new) — `measureDiskSpace` via `syscall.Statfs` + `Stat_t.Dev`, tags aix/darwin/dragonfly/freebsd/ios/linux
- skills/do-work-board/tools/queue-kanban/disk_space_windows.go (new) — `measureDiskSpace` via `GetDiskFreeSpaceExW`, device id = FNV-64a of upper-cased volume name
- skills/do-work-board/tools/queue-kanban/disk_space_unsupported.go (new) — returns `errDiskSpaceUnsupported` wrapped with `%w`
- skills/do-work-board/tools/queue-kanban/disk_space_test.go (new) — six probe tests with fake measurements
- skills/do-work-board/tools/queue-kanban/generate_test.go (modified) — `TestMain` plenty-free fake; `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` seeds low-space root + worktree findings and scans Subject
- skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (modified) — one Traps bullet for the probe
- skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md (modified) — `[family: disk-space-blind-spot] 0.305.60` entry at the top of the list
## Integration seams
- `do-work/lessons-index.md`, the `lessons-do-kanban.md` row: add family `disk-space-blind-spot` to the slug set and refresh the token estimate ((bytes+3)//4 of the new file). Orchestrator's file; not edited.
- Release (CHANGELOG.md + mirror, VERSION x2, actions/version.md): the lessons entry is keyed `0.305.60`; if the release lands on another number, change that one token.
## P-A-U
- [PLAN] Follow the fixed design. A per-platform `measureDiskSpace` fills a small struct, and `appendDiskSpaceFindings` takes the measurer as a parameter so unit tests pass fakes directly. `collectVerifyFindings` passes the package variable `diskSpaceMeasurer`, which `TestMain` sets to a plenty-free fake so existing real-disk tests stay deterministic.
- [APPLY] All edits stayed inside the eight Scope files. Nothing under `do-work/` was written except this hand-back. CHANGELOG, VERSION, version.md, go.mod, go.sum, verify_test.go, atomic_replace_*.go, web/, work.md and work-reference.md are untouched.
- [UNIFY] `git diff --stat 8812a445..HEAD`:
  ```
  .../tools/queue-kanban/disk_space_test.go          | 223 +++++++++++++++++++++
  .../tools/queue-kanban/disk_space_unix.go          |  44 ++++
  .../tools/queue-kanban/disk_space_unsupported.go   |  14 ++
  .../tools/queue-kanban/disk_space_windows.go       |  54 +++++
  .../tools/queue-kanban/generate_test.go            |  35 +++-
  .../tools/queue-kanban/lessons-do-kanban.md        |   2 +
  .../tools/queue-kanban/prime-do-kanban.md          |   1 +
  skills/do-work-board/tools/queue-kanban/verify.go  | 117 +++++++++++
  8 files changed, 488 insertions(+), 2 deletions(-)
  ```
  Linters: `gofmt -l .` empty; `go vet .` plus cross-vet for 12 GOOS/GOARCH pairs all exit 0. Each changed file was re-read in the diff; no debug artifacts. A throwaway real-disk check (written, run, deleted before commit) measured 52.0 GiB free of 177.5 GiB on this Mac, matching `df -g`, and a missing directory returned `statfs: no such file or directory`. Prime file is 56 lines, under the ~60 budget.
## Red-green evidence
RED was run against a compiling stub (constants, types, sentinel, nil measurer, empty `appendDiskSpaceFindings`), so every failure is an assertion, not a build error. RED was not committed, to keep the branch bisectable. GREEN at 2a698a96.
- TestDiskSpaceProbeReportsLowFreeSpaceAtEachThreshold: RED `2_GiB_free_is_critical: disk_space_test.go:77: expected exactly one low-disk-space finding, got []` and `9_GiB_free_is_a_warning: ... got []` → GREEN at 2a698a96
- TestDiskSpaceProbeThresholdBoundariesAreStrict: RED `disk_space_test.go:109: exactly 3 GiB free should be one warning-level finding, got []` → GREEN at 2a698a96
- TestDiskSpaceProbeReportsOneFindingPerDevice: RED `disk_space_test.go:140: two worktrees on one device should yield one finding, got []` → GREEN at 2a698a96
- TestDiskSpaceProbeUnsupportedPlatformIsSkippedNotClean: RED `disk_space_test.go:161: SkippedProbes = [], want exactly ["disk-space probe: unsupported on darwin"]` → GREEN at 2a698a96
- TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory: RED `disk_space_test.go:179: the repo root should still be reported, got []` and `disk_space_test.go:183: SkippedProbes = [], want exactly ["disk-space probe for /fixture/worktrees/worktree-agent-REQ-901: permission denied"]` → GREEN at 2a698a96
- TestDiskSpaceProbeReachesTheVerifyReport (the REQ's RED case through `runVerifyProbes`): RED `disk_space_test.go:206: expected one critical low-disk-space finding for the repo root, got []` and `disk_space_test.go:221: expected one unsupported disk-space skip, got ["calibration-log probe: ..." "worktree probes: git worktree list failed (no worktree support, or not a git repo)"]` → GREEN at 2a698a96
- TestGeneratedVerifyPayloadCarriesNoAbsolutePaths (extended): RED `generate_test.go:3203: low-disk-space findings missing from the payload, so it proves nothing; subjects: map[]` → GREEN at 2a698a96
## Tests run
All from `skills/do-work-board/tools/queue-kanban` with `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off`.
- `go test -count=1 -run 'DiskSpace|GeneratedVerifyPayloadCarriesNoAbsolutePaths|Verify' .` → exit 0, 71 tests passed (top-level and subtests)
- `go test -count=1 ./...` → exit 0, 420 passed, 101 skipped (JavaScript and browser probe lanes off), 0 failed
- `gofmt -l .` → no output
- `go vet .` → exit 0
- `GOOS=windows GOARCH=amd64 go vet .` → exit 0
- `GOOS=linux go vet .` → exit 0
- `GOOS=freebsd go vet .` → exit 0
- `GOOS=openbsd go vet .` → exit 0
- `GOOS=netbsd go vet .` → exit 0
- `GOOS=solaris go vet .` → exit 2, `go: unsupported GOOS/GOARCH pair solaris/arm64` (host is arm64; solaris exists only on amd64). `GOOS=solaris GOARCH=amd64 go vet .` → exit 0
- `GOOS=js GOARCH=wasm go vet .` → exit 0
- The brief's `&&` chain as written → exit 2 at the solaris step, for the reason above; every step run individually passes.
- Extra cross-vet: aix/ppc64, dragonfly/amd64, illumos/amd64, linux/386, linux/arm → each exit 0
## Decisions
- D-01: `errDiskSpaceUnsupported` is declared in verify.go, not in the unsupported file. verify.go references it on every platform, so it cannot live in a tag-gated file. The unsupported file wraps it with `%w` and the GOOS. DECIDE & STATE.
- D-02: `aix` stays in the unix tag list, as the brief says, although the REQ text said to exclude it. Go 1.26's `syscall` has `Statfs` on aix with `Bavail`/`Bsize`/`Blocks`, and `GOOS=aix GOARCH=ppc64 go vet .` passes. netbsd (no `syscall.Statfs`), openbsd (fields are `F_`-prefixed), illumos and solaris fall to the unsupported file. DECIDE & STATE.
- D-03: The negative-`Bavail` clamp converts through `int64` once, which compiles for both signed (dragonfly, freebsd) and unsigned field types without a per-OS branch. DECIDE & STATE.
- D-04: The free percentage prints with one decimal. The REQ's own example, 2 GiB of 500 GiB, is 0.4%, which a whole-number format would print as "0%". Builder latitude on Detail wording. DECIDE & STATE.
- D-05: When git is not on PATH, the disk probe adds no skip line of its own. `appendWorktreeFindings` already reports "worktree probes: git is not on PATH", and the repo root is still measured. An enumeration error does add "disk-space probe for worktrees: <err>", per the brief. DECIDE & STATE.
- D-06: A device counts as measured after any successful measurement, healthy or not. A second directory on the same device has the same free space, so it can add nothing. DECIDE & STATE.
- D-07: A zero total yields 0.0% instead of NaN. It is a one-line guard against a pseudo-filesystem answer. DECIDE & STATE.
## Discovered Tasks
- The brief's cross-vet chain uses `GOOS=solaris go vet .` with no GOARCH. That fails on any arm64 host, which includes this Mac, and stops the `&&` chain before js/wasm. The fix is `GOOS=solaris GOARCH=amd64`, in the exploration/brief template or wherever this command is recorded. (Process; not a code defect.)
- `collectVerifyFindings` now runs `git worktree list` twice per verify: once in `appendWorktreeFindings` and once for the disk probe. That is one extra git subprocess per serve request. Passing one enumeration to both probes would remove it. The fixed design chose the simpler wiring, so this was not changed. (Low value; optional.)
## Lessons read
- skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md — whole satellite, skimmed in full; the `unknown-reads-as-clean` and `subject-not-restated-in-detail` families shaped the skip behavior and the Detail wording.
- _dev/primes/prime-kanban-board.md and skills/do-work-board/tools/queue-kanban/prime-do-kanban.md — read in full.
- _dev/primes/lessons-kanban-board.md — not read (the REQ dropped it for budget; no finding-detail convention beyond the Subject rule was needed).
## Blockers
- None

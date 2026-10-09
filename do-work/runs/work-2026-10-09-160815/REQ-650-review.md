## Review: REQ-650

**Approve** — the disk probe now measures only the repo root, and every kept behaviour is the same as before.
Route A | merge 674fa943 (range 54056e00..674fa943, builder commit b504b4a3)

### What's built
- `appendDiskSpaceFindings` (`skills/do-work-board/tools/queue-kanban/verify.go:256`) calls the measurer once, with the repo root. The worktree loop, the `measuredDevices` dedupe, the `deviceIdentity` field and both platform files' device reads are gone.
- Nothing is missing against the REQ.

### Decisions / risks for you
- None. D-04 (the "disk-space probe for worktrees" skip line is gone) is safe: a worktree-list failure is still reported as "worktree probes: ..." at `verify.go:1241-1243`.

### Findings

**Nit:**
- F1. `verify.go:258` builds `repoRootReading` before the error switch, but only the success path uses it. Building it after the switch would read more directly. — impact-negligible → report only
- F2. `generate_test.go:3165` now seeds a second repo-root finding through `appendDiskSpaceFindings`. It repeats the live probe's finding, so it adds little beyond the `.` count of 2. It is harmless and keeps the test's two-source shape. — impact-negligible → report only

### Requirements Checklist

- [x] R1. The probe measures exactly one path, the repo root; the loop, dedupe map, device field and per-device comments are removed. `grep -rn "measuredDevices\|deviceIdentity"` over `skills/` prints nothing. Delivered.
- [x] R2. The shared worktree read stays: `gitState.worktreePathsByName` still feeds `appendWorktreeFindings` (`verify.go:1246`). Only the disk probe's local copy and its skip line went. Delivered.
- [x] R3. `disk_space_unix.go` drops `os.Stat` and the `os` import; `disk_space_windows.go` drops the FNV hash and the `hash/fnv` and `strings` imports. `syscall.Statfs`, `GetDiskFreeSpaceExW`, build tags and `disk_space_unsupported.go` are unchanged. `go vet .` exit 0. Delivered.
- [x] R4. Thresholds, `diskSpaceLevelFor`, the Detail and Remedy strings, both skip-reason strings ("not measured on GOOS", "not measured: err"), both skipped-probe strings and the `RepoRootDiskSpace` reading match the old index-0 branch line by line. The unsupported case still returns early. The error case used to `continue` into the worktree loop, which is now empty, so returning is the same outcome. Delivered.
- [x] R5. Tests: `TestDiskSpaceProbeReportsOneFindingPerDevice` and `TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory` covered only the deleted branch. The threshold, boundary, unsupported, root-reading and skip-reason tests stay with their worktree inputs removed. The skip-reason test now also checks that an unmeasured root creates no finding. Delivered.
- [x] R6. Docs: `prime-do-kanban.md:36` now says "the filesystem holding the repo root, the only path it measures". `board-guide.md:78` and `actions/board.md:98` already said "the repo root's free disk space" and needed no change (D-02). Delivered.
- [N/A] R7. The changelog entry is written by the integrator at release.

### Acceptance Testing

**Result: Pass**
- `go test -count=1 -run 'TestDiskSpace|TestGeneratedVerifyPayloadCarriesNoAbsolutePaths' .` at HEAD: `ok` (0.5 s). `go vet .`: exit 0.
- RED/GREEN for `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot` is real. It only sets `freeBytes` and `totalBytes`, so it compiles against the base type too. It drives `runVerifyProbes` over a real git repo with one `worktree-agent-REQ-1-fixture` worktree. At base the old loop measured both paths, which matches the recorded RED output (`measured [<repo root> <worktree path>]`). At HEAD the measurer is called once. It compares against `report.RepoRoot`, so the macOS `/private` symlink does not break it.
- The `generate_test.go` edit is the minimum forced change (D-06): `TestMain`'s fake drops the field, and the payload test still requires two `.` subjects and nothing else. That still proves the absolute-path reduction.

### Suggested Additional Testing

- None beyond the planned heavy lanes (queue-kanban-javascript, queue-kanban-browser with `QUEUE_KANBAN_BROWSER` set, staged-skills).

### Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All six code/doc requirements delivered |
| Code Quality | 97% | Straight deletion; one ordering nit (F1) |
| Test Adequacy | 96% | Real RED/GREEN; removed tests matched the removed branch |
| Scope | 100% | `generate_test.go` forced and declared via D-06 |
| Risk | None | Callers: one production call site, tests only |
| Acceptance | Pass | Targeted tests and vet green at HEAD |

### Domain Review

Backend rules: no new I/O, no new error paths. One fewer `os.Stat` per verify run. No concerns.

### Follow-ups created
None (2 findings report only)

## Review

**Overall: 98%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None. Nits: F1 `verify.go:258` builds `repoRootReading` before the error switch that never uses it — impact-negligible → report only; F2 `generate_test.go:3165` seeded second repo-root finding repeats the live probe's — impact-negligible → report only
**Acceptance:** Pass — the measurer is called once with the repo root; disk tests and vet are green at 674fa943, and RED at base 2eb14357 matches the test code.
**Restatement sweep:** redefined the low-disk-space probe's scope ("repo root and each worktree, one finding per device" became "repo root only"); grepped skills/, _dev/, docs, README and both CHANGELOGs for "per device", "each builder worktree", "deviceIdentity", "measuredDevices", "disk-space probe for worktrees", and disk near worktree-agent. Only hits are the historical 0.305.60 entry in skills/do-work/CHANGELOG.md:222 and its root mirror CHANGELOG.md:222 (expected to stay). prime-do-kanban.md:36, board-guide.md:78 and actions/board.md:98 all say repo root. No stale restatements.
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

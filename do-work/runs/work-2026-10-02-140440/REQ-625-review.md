## Review: REQ-625

**Approve** — the low-disk-space probe is built as specified, every test and cross-platform vet passes, and the three pushed-back items are absent.
Route B | builder 2a698a96, merge add2c162 (range 8812a445..add2c162)

### What's built
- `queue-kanban verify` and the board's VERIFY band now report a `low-disk-space` finding when the repo root or a `worktree-agent-*` worktree has under 10 GiB free (warning) or under 3 GiB free (critical), one finding per device.
- Platforms without a free-space call report the probe as skipped, never as clean. A failing directory is a skip for that directory only.
- Still open, and owned by finalization: the release itself. VERSION is still 0.305.59, and the lessons entry is keyed 0.305.60, so the release must land on 0.305.60 or that token must change.

### Decisions / risks for you
- The builder kept `aix` in the unix tag list although the REQ said to exclude it (hand-back decision D-02). Evidence supports this: `GOOS=aix GOARCH=ppc64 go vet .` exits 0 here, and the exploration found `syscall.Statfs` with `Bavail` on aix. The REQ's premise was wrong, not the build.
- Each verify run, and so each board serve request, now runs `git worktree list` twice. The REQ constraint said "one Statfs call per measured directory and nothing heavier". The cost is one small git subprocess. It is recorded below as a Minor finding, not a blocker.

### Findings

**Important:**
- None.

**Minor:**
- `skills/do-work-board/tools/queue-kanban/verify.go:204-210` calls `listWorktreeAgentWorktrees` a second time after `appendWorktreeFindings` already did. That adds one `git worktree list` subprocess per verify and per serve request, against the REQ's "nothing heavier" constraint. Passing one enumeration to both probes removes it. The builder listed this as a discovered task. — impact-negligible → report only
- `skills/do-work-board/tools/queue-kanban/verify.go:206-208`: in a project where `git worktree list` fails (not a git repo), the report now shows two skip lines for one cause: "worktree probes: `git worktree list` failed …" and "disk-space probe for worktrees: `git worktree list` failed …". Both are true, but the second repeats the first. — impact-user-visible → report only
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go`: the dedupe test covers two worktrees on one device, but no test covers a worktree on the same device as the repo root, which is the common layout (this machine's live run has exactly that). The code path is the same map lookup, so the risk is low. — impact-negligible → report only

**Nit:**
- `skills/do-work-board/tools/queue-kanban/verify.go:296-299`: a filesystem that reports zero total bytes (some pseudo or FUSE mounts) produces a critical "0.0 GiB free of 0.0 GiB (0.0%)" finding. The guard from decision D-07 avoids NaN but still reports it as critical. Treating zero total as "unmeasurable, skip" would be more honest. — impact-negligible → report only
- `skills/do-work-board/tools/queue-kanban/disk_space_unix.go:37-41` multiplies block counts by `Bsize`. On Linux, `statfs` also has `Frsize`, and some tools use it for block counts. The two values are equal on common Linux filesystems (ext4, xfs, btrfs), so this is not a confirmed defect. I could not test it on Linux here. — impact-negligible → report only

### Requirements Checklist

- [x] 1. Build-tagged helper split like `atomic_replace_*.go`: unix `syscall.Statfs` (tags aix, darwin, dragonfly, freebsd, ios, linux), Windows `GetDiskFreeSpaceExW` via `syscall.NewLazyDLL("kernel32.dll")`, unsupported file with a typed error. `Bavail` is used, all `Statfs_t` arithmetic stays in the platform file, and `go.mod` is unchanged. Delivered. The tag list follows measured compile evidence instead of the REQ's aix wording (D-02, documented).
- [x] 2. `appendDiskSpaceFindings` is called from `collectVerifyFindings` after `appendWorktreeFindings`. It measures the root and every listed worktree, dedupes by `Stat_t.Dev` on unix and an FNV-64a hash of the volume name on Windows, emits at most one finding per device, returns early with one skip line on the unsupported error, and skips only the failing directory on other errors. Delivered. The signature takes the worktree map and the measurer as parameters, which the Builder Guidance allowed.
- [x] 3. Thresholds are two named constants, `lowDiskSpaceWarningBytes = 10 << 30` and `lowDiskSpaceCriticalBytes = 3 << 30`, with no percentage rule. Delivered. The REQ's suggested name was `lowDiskSpaceWarnBytes`; the builder had naming latitude.
- [x] 4. Finding shape: category `low-disk-space`, Subject is the repo root (reduces to `.`) or the worktree name, `Fixable: false`, Detail in the specified form with GiB values and one-decimal percent (D-04), and the Remedy is byte-identical to the REQ's text and names no project folder. No largest-directories walk. Delivered.
- [x] 5. No new path code. The extended payload test proves Subject, Detail and Remedy carry no absolute path for this finding. Delivered.
- [x] 6. Read-only: the probe only calls `Statfs`/`os.Stat` or `GetDiskFreeSpaceExW` and deletes nothing. Delivered.
- [x] 7. Docs: one Traps bullet in `prime-do-kanban.md`, one `[family: disk-space-blind-spot]` entry at the top of `lessons-do-kanban.md` in the existing shape. The `do-work/lessons-index.md` row is refreshed in the main tree (families sorted, 7095 tokens, which matches 28377 bytes / 4 rounded up), uncommitted and owned by the orchestrator. Delivered.
- [ ] 8. Release (changelog entry, version bump). Not in the merge range. VERSION is still 0.305.59. N/A for the builder: finalization owns it. Not counted against the score.
- [x] Out of scope, confirmed absent: no header chip (no `web/` file in the diff, no `generatedBoardData` field added), no largest-directories walk (no directory iteration in the probe), no `work.md` or `work-reference.md` change.
- [x] Scope acceptance criteria: all eight are met (thresholds and fixable flag, strict boundaries, per-device dedupe, unsupported and per-directory skips, Subject reduction plus payload test, generic remedy, no new dependency plus cross-vet, prime and lessons entries).

### Code review notes

- **Unix arithmetic.** `int64(Bavail)` then a `> 0` check clamps signed negative values on freebsd and dragonfly to zero, and compiles for unsigned fields elsewhere. A uint64 above 2^63 blocks would also clamp to zero, which is not reachable in practice.
- **Strict boundaries.** `freeBytes < critical` and `freeBytes < warning` mean exactly 3 GiB is a warning and exactly 10 GiB is clean, as the test pins.
- **Dedupe.** A device is marked measured only after a successful measurement, the root is measured first, and worktrees follow in sorted name order, so the output is deterministic.
- **Windows shape.** It matches `atomic_replace_windows.go`: package-level lazy procedure, `UTF16PtrFromString` with a wrapped error, `Call`, a zero result, and the `syscall.Errno(0)` "failed without an error code" branch. The mounted-folder caveat for volume identity is documented in the file comment.
- **Restatement sweep.** The diff adds a verify category and a probe. Forensics Check 14 maps findings by output class, not by a category list, so it is not stale. `board.md` defers to Check 14. `boardRenderedVerifyCategories` is a suppression list and does not need the new category. No doc, JavaScript file or test enumerates the verify category set. Skipped-probe lines with absolute paths are reduced by `attachVerifyFindings` at `generate.go:674-676`, so the per-directory skip text does not leak paths into the board. No stale restatement found.
- **Naming and comments.** Names with reach use two or more words. The `measure` parameter and `target` loop variable are short locals under the carve-out. Comments match the code.

### Acceptance Testing

**Result: Pass**

All commands ran from `skills/do-work-board/tools/queue-kanban` in the main checkout at add2c162.

| Command | Exit |
|---|---|
| `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...` | 0 (65.5 s) |
| `gofmt -l .` | 0, no output |
| `go vet .` | 0 |
| `GOOS=windows GOARCH=amd64 go vet .` | 0 |
| `GOOS=openbsd go vet .` | 0 |
| `GOOS=netbsd`, `linux`, `freebsd` `go vet .` | 0 each |
| `GOOS=solaris GOARCH=amd64`, `illumos/amd64`, `aix/ppc64`, `dragonfly/amd64`, `js/wasm` `go vet .` | 0 each |
| `go run . verify --repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2` | 1 |

- **Live CLI.** The verify run printed one finding, `worktree-present-run-in-flight` for this REQ's own builder worktree, which is the expected state while REQ-625 is in `do-work/working/`. That finding is why the exit code is 1. There was no `low-disk-space` finding and no disk-space skip line, which is correct: `df -h` shows 50 GiB free of 178 GiB. The builder worktree sits on the same device as the repo root, so the dedupe path ran.
- **Real-disk measurement.** A throwaway test injected with `go test -overlay` (no file written into the package) called `measureDiskSpace` directly. It returned 49.7 GiB free of 177.5 GiB, which matches `df`. A missing directory returned `statfs: no such file or directory`. A zero-total fake produced the critical finding described in the Nit above. The package tree was not modified.
- **Red-green.** The hand-back records RED against a compiling stub for all six new tests and the extended payload test, each failing on an assertion, and GREEN at 2a698a96. The RED stub was not committed, which keeps the branch bisectable. The GREEN half was reproduced by the test run above.
- **Test adequacy.** Each test pins the failure its comment names. The threshold test fails if the probe stays silent. The boundary test fails on an off-by-one and also locks the constant values. The device test fails on duplicates and on a path used as Subject. The unsupported test fails if the platform reads as clean or emits one skip per directory. The error test fails if one bad directory hides the root. The reach test fails if the probe is unwired from `runVerifyProbes`. The extended payload test now scans Subject and requires two `.` subjects and one worktree-name subject, so it fails on a dead probe. `TestMain` installs a plenty-free fake, so existing zero-finding tests stay deterministic on a nearly full machine.
- **Regressions.** The full package suite passed, including existing verify tests that assert zero findings through `collectVerifyFindings`.

### Suggested Additional Testing

- Run `verify` on a Linux host and on Windows once, to confirm the real measurement matches `df` or the drive properties. Cross-vet proves only that the code compiles there.
- Run the board's `serve` mode on a nearly full disk, or a small RAM disk, to see the finding render in the VERIFY band.
- Run `verify` with a builder worktree on a second volume, such as an external disk, to see one finding per device for real.

### Scores (on the record — not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All seven builder requirements delivered; release (requirement 8) is finalization's and is excluded |
| Code Quality | 92% | Clean, matches the Windows pattern; one redundant git enumeration and one zero-total edge |
| Test Adequacy | 92% | Every test pins a named failure with red-green evidence; no root-plus-worktree same-device case |
| Scope | 100% | Exactly the eight declared files; seven decisions recorded in the hand-back |
| Risk | Low | One extra `git worktree list` subprocess per serve request |
| Acceptance | Pass | Suite, gofmt, vet and 11 cross-vets green; live CLI clean on a 50 GiB-free disk |

### Follow-ups created
None (5 findings report only)

Adjacent observation, not caused by this diff: `skills/do-work-board/actions/board.md:35` says forensics Check 14 holds "the canonical probe descriptions", but Check 14 contains no probe descriptions and names the board tool as the authority.

---

## Review

**Overall: 96%** | 2026-10-02T14:24:35Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- verify.go:204-210 runs `git worktree list` a second time per verify and serve request; pass one enumeration to both probes — impact-negligible → report only
- verify.go:206-208 prints a second skip line for the same `git worktree list` failure in a non-git project — impact-user-visible → report only
- disk_space_test.go has no case for a worktree on the repo root's device, the common layout — impact-negligible → report only
- Nit: verify.go:296-299 reports a zero-total filesystem as critical instead of skipping it — impact-negligible → report only
- Nit: disk_space_unix.go:37-41 uses `Bsize`, not `Frsize`, on Linux; equal on common filesystems, untested here — impact-negligible → report only

**Acceptance:** Pass — full package suite, gofmt, go vet and 11 cross-platform vets exit 0; live `verify` on this repo shows no disk finding at 50 GiB free and exits 1 only for the expected in-flight worktree finding.
**Suggested testing:** 3 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

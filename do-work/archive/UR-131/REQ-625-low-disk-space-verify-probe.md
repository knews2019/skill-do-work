---
id: REQ-625
title: 'Add a low-disk-space probe to the board VERIFY band'
status: completed
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 13-file write set
  - 2 subsystems involved
  - 8 acceptance criteria
  calculated_at: 2026-10-02T14:06:02Z
created_at: 2026-10-02T13:58:27Z
user_request: UR-131
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: ["skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/disk_space_unix.go", "skills/do-work-board/tools/queue-kanban/disk_space_windows.go", "skills/do-work-board/tools/queue-kanban/disk_space_unsupported.go", "skills/do-work-board/tools/queue-kanban/disk_space_test.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md"]
integration_at: 2026-10-02T14:20:39Z
builder_handback_at: 2026-10-02T14:19:49Z
dispatch_at: 2026-10-02T14:12:37Z
review_at: 2026-10-02T14:26:59Z
commit: add2c162e69ca70d22b7362e121028de4dcae3e7
kb_status: pending
heavy_verified_at: 2026-10-02T14:28:15Z
heavy_verified_revision: add2c162e69ca70d22b7362e121028de4dcae3e7
claimed_at: 2026-10-02T14:03:40Z
completed_at: 2026-10-02T14:32:48Z
release_at: 2026-10-02T14:32:48Z
---

# Add a Low-Disk-Space Probe to the Board VERIFY Band

## What
Add one read-only verify probe to the board tool that measures free disk space on the filesystem holding the repo root (and each enumerated `worktree-agent-*` worktree on a different device) and emits a `low-disk-space` finding below fixed thresholds. The finding reaches the CLI `verify` report and the board's VERIFY band through the existing `attachVerifyFindings` path, with no new UI.

## Why
During a long `do-work run --fan-out 2` on a game repo (October 2026), browser-QA runs wrote hundreds of megabytes of screenshots per comparison. The repo grew to about 20 GB in a few hours and the disk nearly filled, and nothing in the suite noticed. An earlier run (REQ-504, commit 1a726b6f) also lost its heavy test lanes to a full disk. Builders, gates and git fail in confusing ways at zero bytes free, and the VERIFY band is where the run operator already looks.

## Detailed Requirements
1. **Helper, build-tagged.** A `diskSpaceFor(directory) (free, total uint64, device uint64, err)` style helper split like `atomic_replace_*.go`: `syscall.Statfs` on unix, `GetDiskFreeSpaceExW` via `syscall.NewLazyDLL("kernel32.dll")` on Windows, and an unsupported file returning a typed "unsupported on GOOS" error. **The unix tag list must be narrower than `atomic_replace_unix.go`**: aix, illumos and solaris have no `syscall.Statfs` and belong to the unsupported file. All `Statfs_t` arithmetic stays inside the platform file (field types differ per OS); use `Bavail`, not `Bfree`. No new module dependency; `go.mod` stays as it is.
2. **Probe.** `appendDiskSpaceFindings(&report, repoRoot)` in `verify.go`, called from `collectVerifyFindings` after `appendWorktreeFindings`. It measures the repo root and every path returned by `listWorktreeAgentWorktrees`, dedupes by device id (`os.Stat(dir).Sys().(*syscall.Stat_t).Dev` cast to uint64 on unix; volume path or serial on Windows, inside the platform file), and emits at most one finding per device. On the unsupported error it appends to `SkippedProbes` (the invariant is real and went unverified) and never fails the run. A `Statfs` error on one directory is a skipped probe for that directory, not a crash.
3. **Thresholds.** Two named constants: `lowDiskSpaceWarnBytes = 10 GiB` and `lowDiskSpaceCriticalBytes = 3 GiB`. No percentage rule (dropped at triage: on a 20 GB disk a 10 % warn is 2 GiB, so critical would fire before warn).
4. **Finding shape.** `Category: verifyCategoryLowDiskSpace = "low-disk-space"`, `Subject: <measured directory>` (the repo root, which reduces to `.`, or the worktree path, which reduces to its name), `Fixable: false`. Detail: `"<free> free of <total> (<pct>%) — below the <warning|critical> threshold (<threshold>)"` with human-readable GiB values. Remedy, generic and project-agnostic: `"free space: clear regenerable QA output, finished builder worktrees (do-work cleanup), browser caches; `du -sh * | sort -h` at the repo root shows the largest directories"`. The suite must not name any project's output folders. **No largest-directories walk** in the probe (pushed back at triage: sizing direct children means walking them fully, so the 200 ms cap would make every result "(partial)" on exactly the repos that need it, on every serve request).
5. **Paths.** Nothing new: `attachVerifyFindings` already reduces Subject, Detail and Remedy. Confirm by test that the board payload carries no absolute path for this finding.
6. **Read-only.** The probe measures and reports. It deletes nothing. Cleanup stays human-consented through `actions/cleanup.md`; `Fixable` stays false because cleanup cannot resolve disk space.
7. **Docs.** Add the probe to `prime-do-kanban.md` where `verify.go`'s probes are named. Add one lesson to `lessons-do-kanban.md` under `[family: disk-space-blind-spot]` with the file's existing entry shape: a long fan-out run with browser QA can fill a disk in hours, and the board is where the operator looks. Refresh the `do-work/lessons-index.md` row for that satellite (families and token count) in the same commit.
8. **Release.** Shipped files change, so this is a release: changelog entry with a descriptive title and a version bump per `_dev/primes/prime-releases.md`. Forensics Check 14 counts every verify finding as a warning and needs no change.

## Out of scope (pushed back at triage, do not build)
- An always-visible free-space chip in the board header.
- A top-3 largest-directories walk inside the finding.
- Any sentence in `work.md` or `work-reference.md` that pauses dispatch on this finding. If a dispatch pause is wanted later, it belongs in queue-mode `advance` as a typed exclusion in the core CLI, never as prose keyed on a board-package finding.

## Constraints
- Read-only probe; no deletion.
- No new Go module dependency.
- Thresholds are constants, not configuration.
- The probe runs on every serve request (findings are computed outside the mtime cache), so it must stay at one `Statfs` call per measured directory and nothing heavier.

## Dependencies
None. `listWorktreeAgentWorktrees`, `SkippedProbes`, `attachVerifyFindings` and `reduceAbsolutePaths` all exist today.

## Builder Guidance
High certainty on scope: the triage report fixed the shape and the exclusions above. Builder latitude on helper naming, on the exact Detail wording, and on where the Windows device identity comes from. Tests use an injected measurement function (a package-level `var` the test swaps), never the real filesystem, because real free space is nondeterministic.

## Red-Green Proof
**RED prompt/case:** With the measurement function faked to report 2 GiB free of 500 GiB for the repo root, `collectVerifyFindings` over a clean synthetic board returns no `low-disk-space` finding. With the function faked to return the unsupported error, `SkippedProbes` has no disk-space entry.
**Why RED now:** No probe exists; `collectVerifyFindings` ends at `appendWorktreeFindings` (`verify.go:195`).
**GREEN when:** The same inputs yield exactly one finding with `Category "low-disk-space"`, `Subject "."` after reduction, `Fixable false`, Detail naming the critical threshold; 20 GiB free yields no finding; 9 GiB yields a warning-level Detail; two worktree paths on one device yield one finding; the unsupported error yields one `SkippedProbes` entry and zero findings; the board payload test `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` still passes with the new finding present.
**Validation:** Inferred during capture

## Tests (each names the failure it pins)
- Fake measurement above and below the thresholds: pins "probe silently clean when disk is low".
- Boundary at exactly 3 GiB and exactly 10 GiB: pins an off-by-one at the constants.
- Two worktrees on the same device produce one finding: pins duplicate findings per device.
- Unsupported platform lands in `SkippedProbes`, never silently clean: pins "unknown reads as clean".
- Extend `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` (`generate_test.go:3125`): pins a machine path leaking into a shareable snapshot.

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (6908 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: changing queue-kanban model and testing. Read it anyway if time allows; the `unknown-reads-as-clean` family is the closest prior.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, over budget; `slugged: partial`). Matched: finding-detail text.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** A per-platform `measureDiskSpace` fills a small struct; `appendDiskSpaceFindings` takes the measurer as a parameter so unit tests pass fakes directly; `collectVerifyFindings` passes the package variable `diskSpaceMeasurer`, which `TestMain` sets to a plenty-free fake so existing real-disk tests stay deterministic. (Builder, from the hand-back.)
- [x] **[APPLY]:** All edits stayed inside the eight Scope files; nothing under `do-work/` was written except the hand-back; release files, go.mod, go.sum, verify_test.go, atomic_replace_*.go, web/ untouched. (Builder; orchestrator confirmed against `git diff --stat 8812a445..add2c162`.)
- [x] **[UNIFY]:** `git diff --stat 8812a445..2a698a96`: 8 files, 488 insertions, 2 deletions. `gofmt -l .` empty; `go vet .` and cross-vet for 12 GOOS/GOARCH pairs exit 0; each changed file re-read in the diff, no debug artifacts; a throwaway real-disk check matched `df -g` and was deleted before commit. (Builder, from the hand-back.)

## Full Context
See `do-work/user-requests/UR-131/input.md` for the complete verbatim input and the original proposal.

*Source: capture the ones that are accepted*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and the touched files are fixed by the REQ, but the package's conventions for injecting a fake measurement into a verify probe, the synthetic-board fixtures the probe tests use, and the Windows syscall shape need discovery before dispatch. One new helper split plus one probe, no architectural change, so no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors: `do-work/runs/work-2026-10-02-140440/REQ-625-exploration.md`. The findings that change what the builder writes:

- **Subject for a worktree is its name, not its path.** Builder worktrees live outside the repo, and `reduceAbsolutePaths` (`generate.go:689-699`) turns any path outside the repo root into `<path outside this repository>`. `appendWorktreeFindings` already uses the map key from `listWorktreeAgentWorktrees` as Subject (`verify.go:1171`). The repo root's Subject is the absolute `repoRoot`, which reduces to `.`.
- **The REQ's build-tag premise was partly wrong.** Cross-vet with go1.26.1: `syscall.Statfs` and `Statfs_t.Bavail` compile on aix, darwin, dragonfly, freebsd, ios, linux. They do not compile on netbsd (`undefined: syscall.Statfs`), openbsd (fields are `F_bavail`), illumos, solaris. Unix tag: `aix || darwin || dragonfly || freebsd || ios || linux`; the unsupported file negates that list and windows.
- **Field types differ per OS** (`Bsize` uint32 on darwin, int64 on linux; freebsd `Bavail` is int64 and can be negative, clamp at 0). Convert with `uint64(...)` inside the platform file; `uint64(stat.Dev)` compiles on every listed OS.
- **Existing zero-finding tests hit the real disk once the probe is live** (`verify_test.go:86,130,1615` assert no findings through `collectVerifyFindings`). The package has a `TestMain` (`generate_test.go:203-216`) and no `t.Parallel()`, so a package-level measurement variable set to a "plenty free" fake in `TestMain` is race-free and keeps them deterministic.
- **No package-level function-variable seam exists today**; injection is by parameter (`buildBoard(..., lookupGitCommitDate)`). Compromise: `appendDiskSpaceFindings` takes the measurement function as a parameter (unit tests call it directly with fakes), and `collectVerifyFindings` passes a package-level `var diskSpaceMeasurement = measureDiskSpace` that `TestMain` overrides.
- **No registration needed for a new category.** `boardRenderedVerifyCategories` (`generate.go:647-652`) is a suppression list; the CLI renderer and the web strip print any category verbatim.
- **The payload test checks Detail and Remedy, not Subject** (`generate_test.go:3172-3178`); extend it to Subject, and seed a low-space fake so it cannot pass on a dead probe.
- **No byte formatter exists**; add a small `formatGibibytes`. `TestVerifyFindingDetailsDoNotRepeatTheirSubject` (`verify_test.go:2446`) fails if Detail starts with the Subject, so Detail starts with the free amount.
- **Docs:** `prime-do-kanban.md:11` names `verify.go` with no per-probe list; add a Traps bullet near lines 31-34. `lessons-do-kanban.md` prepends version-keyed entries at the top of the list in the form `- [family: slug] X.Y.Z: **bold lesson.** prose`. `do-work/lessons-index.md:11` is the row to refresh (tokens = (bytes+3)//4, families sorted, coverage stays partial); that file is under `do-work/` and is the orchestrator's to edit.
- **Tests run:** `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...` from the package, plus cross-`go vet` for windows, linux, freebsd, openbsd, netbsd, solaris and js/wasm, which is the only check that the tag split compiles everywhere.
- **Release preimage:** VERSION 0.305.59 in three places; CHANGELOG.md and its mirror are byte-identical; entry shape `## X.Y.Z — Title (YYYY-MM-DD)`, one lead paragraph, 2-3 bullets.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — verifyCategoryLowDiskSpace, the two threshold constants, formatGibibytes, appendDiskSpaceFindings, the diskSpaceMeasurement package variable, and the call in collectVerifyFindings
- `skills/do-work-board/tools/queue-kanban/disk_space_unix.go` (new) — measureDiskSpace via syscall.Statfs, device id via Stat_t.Dev
- `skills/do-work-board/tools/queue-kanban/disk_space_windows.go` (new) — measureDiskSpace via GetDiskFreeSpaceExW through syscall.NewLazyDLL("kernel32.dll")
- `skills/do-work-board/tools/queue-kanban/disk_space_unsupported.go` (new) — typed errDiskSpaceUnsupported
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (new) — the probe tests with fake measurements
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — TestMain plenty-free fake; extend TestGeneratedVerifyPayloadCarriesNoAbsolutePaths
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify) — one Traps bullet naming the probe
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modify) — one [family: disk-space-blind-spot] entry at the top of the list

**Files I will NOT touch:** `do-work/lessons-index.md` (orchestrator refreshes the row at finalization), `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, VERSION, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (release paths, written by finalization), `go.mod`/`go.sum` (no new dependency), `verify_test.go` (new tests go in their own file to stay under the per-file budget), `atomic_replace_*.go`, `work.md`, `work-reference.md`, any web/ file.

**Acceptance criteria (restated from REQ):**
- [ ] A fake measurement reporting 2 GiB free of 500 GiB for the repo root yields exactly one low-disk-space finding, Fixable: false, Detail naming the critical threshold; 9 GiB yields a warning-level Detail; 20 GiB yields no finding
- [ ] Boundary tests at exactly 3 GiB and exactly 10 GiB pin the constants
- [ ] Two worktree directories on one device yield one finding (dedupe by device id)
- [ ] The unsupported-platform error yields one SkippedProbes entry and zero findings; a per-directory measurement error is a skipped probe for that directory, not a crash
- [ ] Subject is . for the repo root and the worktree name for a worktree after reduction; TestGeneratedVerifyPayloadCarriesNoAbsolutePaths covers Subject and passes with a low-space finding present
- [ ] Remedy is generic and names no project folder; the probe deletes nothing
- [ ] No new module dependency; go vet passes for windows, linux, freebsd, openbsd, netbsd, solaris, js/wasm
- [ ] `prime-do-kanban.md` names the probe; `lessons-do-kanban.md` carries the disk-space-blind-spot entry

## Pre-Flight

**Git:** ✓ Integration tip 8812a445 on `main`; the only dirt is this REQ's own `do-work/` trail (working REQ, run directory) — no third-party paths
**Tests baseline:** ✓ Board package green (`go test -C skills/do-work-board/tools/queue-kanban -count=1 ./...` with JavaScript and browser probes off, launched by `advance`, preflight satisfied)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 14:08:29Z, both Go stages EXECUTING (reuse disabled), gate wall 129s; green-gate record satisfied
**Dependencies:** ✓ Go toolchain present; no new module dependency planned

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — verifyCategoryLowDiskSpace, lowDiskSpaceWarningBytes/lowDiskSpaceCriticalBytes, diskSpaceMeasurement, errDiskSpaceUnsupported, diskSpaceMeasurer, appendDiskSpaceFindings, formatGibibytes, wired after appendWorktreeFindings in collectVerifyFindings
- `skills/do-work-board/tools/queue-kanban/disk_space_unix.go` (new) — syscall.Statfs + Stat_t.Dev, tags aix/darwin/dragonfly/freebsd/ios/linux
- `skills/do-work-board/tools/queue-kanban/disk_space_windows.go` (new) — GetDiskFreeSpaceExW via kernel32; device id is FNV-64a of the upper-cased volume name
- `skills/do-work-board/tools/queue-kanban/disk_space_unsupported.go` (new) — wraps errDiskSpaceUnsupported with the GOOS
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (new) — six probe tests with fake measurements, including the REQ's RED case through runVerifyProbes
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — TestMain plenty-free fake; TestGeneratedVerifyPayloadCarriesNoAbsolutePaths seeds low-space root and worktree findings and scans Subject
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified) — one Traps bullet for the probe
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modified) — [family: disk-space-blind-spot] 0.305.60 entry at the top of the list

**What was done:** Added a read-only low-disk-space verify probe that measures the repo root and each worktree-agent-* worktree, one finding per device below 10 GiB (warning) or 3 GiB (critical), skipped rather than silent on unsupported platforms, reaching the CLI report and the board's VERIFY band through the existing attachVerifyFindings path. Merge range 8812a445..add2c162 (builder commit 2a698a96, merge add2c162). The `do-work/lessons-index.md` row refresh is the orchestrator's integration seam, applied in the main tree.

## Qualification

**Diff range:** 8812a445..add2c162 (builder commit 2a698a96, merge add2c162)
**Gate records:** qualify satisfied; scope-drift satisfied (second call; the first call read backticked identifier names in Scope and the Implementation Summary as paths, so those sections were reworded, no code change).
**Warnings judged:** QUALIFY-NEW-FILE-UNWIRED on disk_space_windows.go — expected: the file is selected by its `//go:build windows` tag and defines the same measureDiskSpace the unix and unsupported files define, so no static reference outside it can exist. Not dead code.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. Thresholds are two typed constants with strict `<`; one finding per device with the repo root measured first and worktrees in sorted name order; the unsupported sentinel returns after one skip line; a per-directory error skips only that directory; Subject is the absolute repo root or the worktree name; Remedy names no project folder; nothing is deleted; go.mod unchanged. The signed-Bavail clamp converts through int64 once (D-03). The three pushed-back items (header chip, directory walk, run-loop prose) are absent.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY was cross-checked against `git diff --stat 8812a445..add2c162` (eight Scope files, nothing else outside do-work/).
**Live data flow:** `collectVerifyFindings` is the single producer for both the CLI report and `attachVerifyFindings`, so the finding reaches the VERIFY band with no new path.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (DO_WORK_FAST_STAGE_REUSE=off) on the merged tree at add2c162
**Result:** ✓ All passing — exit 0, gate wall 131s; stage queue-kanban-fast-tests EXECUTING (409 tests, wall 42s, slowest file 21.15s < 30s); stage do-work-cli-fast-tests EXECUTING (866 tests, wall 60s, slowest file 21.49s < 30s). Green-gate record satisfied by advance.

**Focused tests:** `go test -count=1 -run 'DiskSpace|GeneratedVerifyPayloadCarriesNoAbsolutePaths|Verify' .` in the board package with JavaScript and browser probes off → exit 0 (advance probe record satisfied). A first probe over the whole package (`./...`, 42s wall) exceeded the probe timeout (status 124); the whole package is covered by the repository gate above, so the probe was narrowed to the REQ's tests rather than rerun.

**Red-green validation:** traced to `## Red-Green Proof`; RED was taken against a compiling stub (constants, types, sentinel, empty probe), GREEN at builder commit 2a698a96:
- TestDiskSpaceProbeReportsLowFreeSpaceAtEachThreshold (2 GiB critical, 9 GiB warning, 20 GiB clean): ✗ `expected exactly one low-disk-space finding, got []` → ✓
- TestDiskSpaceProbeThresholdBoundariesAreStrict (exactly 3 GiB warning, exactly 10 GiB clean): ✗ → ✓
- TestDiskSpaceProbeReportsOneFindingPerDevice: ✗ `two worktrees on one device should yield one finding, got []` → ✓
- TestDiskSpaceProbeUnsupportedPlatformIsSkippedNotClean: ✗ `SkippedProbes = [], want exactly ["disk-space probe: unsupported on darwin"]` → ✓
- TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory: ✗ → ✓
- TestDiskSpaceProbeReachesTheVerifyReport (the REQ's RED case through runVerifyProbes): ✗ `expected one critical low-disk-space finding for the repo root, got []` → ✓
- TestGeneratedVerifyPayloadCarriesNoAbsolutePaths (extended, now scans Subject and requires a low-disk-space entry): ✗ `low-disk-space findings missing from the payload, so it proves nothing` → ✓

**New tests added:**
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (six tests above)

**Existing tests updated (cross-REQ impact):**
- `generate_test.go` TestMain (from REQ-284's payload work): installs a plenty-free fake measurer so the three real-disk zero-finding tests in verify_test.go stay deterministic on a nearly full machine — intentional
- `generate_test.go` TestGeneratedVerifyPayloadCarriesNoAbsolutePaths: extended to Subject and to a seeded low-space root plus worktree — intentional

**Cross-platform compile proof (builder, each exit 0):** go vet for windows/amd64, linux, freebsd, openbsd, netbsd, solaris/amd64, js/wasm, aix/ppc64, dragonfly/amd64, illumos/amd64, linux/386, linux/arm; gofmt -l empty. Note: bare `GOOS=solaris go vet` fails on an arm64 host (no solaris/arm64 pair), so the amd64 pair is the one that proves it.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 8812a445..add2c162
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

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

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-02-140440/REQ-625-handback.md`) so they survive run-directory cleanup. All seven are DECIDE & STATE.

- D-01: errDiskSpaceUnsupported is declared in verify.go, not in the unsupported file, because verify.go references it on every platform and a tag-gated file cannot own it; the unsupported file wraps it with `%w` and the GOOS.
- D-02: aix stays in the unix tag list although the REQ text said to exclude it. Go 1.26's syscall has Statfs with Bavail on aix and `GOOS=aix GOARCH=ppc64 go vet .` passes; netbsd, openbsd, illumos and solaris fall to the unsupported file.
- D-03: the negative-Bavail clamp converts through int64 once, which compiles for both signed and unsigned field types without a per-OS branch.
- D-04: the free percentage prints with one decimal, because the REQ's own example (2 GiB of 500 GiB) is 0.4% and a whole-number format would print 0%.
- D-05: when git is not on PATH the disk probe adds no skip line of its own; the worktree probes already report it and the repo root is still measured. An enumeration error does add a disk-space skip line.
- D-06: a device counts as measured after any successful measurement, healthy or not, since a second directory on the same device has the same free space.
- D-07: a zero total yields 0.0% instead of NaN (one-line guard).

Orchestrator decision:
- D-09: finalization refused the release with RELEASE-WITHOUT-SHIPPED-CHANGE because the guard counted only VERSION-carrying modules as shipped roots, and this REQ changes only `skills/do-work-board`. Finalizing without a release would ship the probe with no version or changelog entry, against `_dev/primes/prime-releases.md`. Cleared by a focused fix in `skills/do-work/tools/do-work-cli` (every declared module source is a shipped root; version ownership unchanged), pinned by a new guard test case, committed on its own before finalization and covered by this release's changelog. Outside the declared Scope; recorded here rather than silently absorbed. DECIDE & STATE (reversible, evidence in the test).
- D-08: the focused-test probe was narrowed from the whole board package (42s, over the probe timeout) to the REQ's own tests; the whole package is proven by the repository gate on the merged tree. DECIDE & STATE.

## Discovered Tasks

From the builder's hand-back and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- The brief's cross-vet chain used `GOOS=solaris go vet .` with no GOARCH, which fails on an arm64 host and stops the chain; the amd64 pair is the one that proves solaris. Process note for the next brief that records this command. — impact-negligible → report only
- collectVerifyFindings runs `git worktree list` twice per verify (once in appendWorktreeFindings, once for the disk probe); passing one enumeration to both probes would remove one subprocess per serve request. — impact-negligible → report only
- In a non-git project the report prints two skip lines for one `git worktree list` failure (worktree probes and disk-space probe for worktrees). — impact-user-visible → report only
- No test covers a worktree on the repo root's own device, the common layout; the code path is the same map lookup. — impact-negligible → report only
- A filesystem reporting zero total bytes is reported as a critical 0.0 GiB finding instead of being skipped as unmeasurable. — impact-negligible → report only
- `skills/do-work-board/actions/board.md:35` says forensics Check 14 holds "the canonical probe descriptions", but Check 14 names the board tool as the authority and holds no descriptions. Adjacent, not caused by this diff. — impact-negligible → report only

## Lessons Learned

**What worked:** Cross-compiling with `go vet` per GOOS before writing the tag list; it overturned the REQ's premise (aix supports Statfs, netbsd and openbsd do not) before any code was written. Taking RED against a compiling stub so every failure is an assertion, not a build error, and not committing the stub so the branch stays bisectable.
**What didn't:** A focused-test probe over the whole board package exceeded the advance probe timeout; the probe has to name the REQ's own tests and leave the whole package to the repository gate. Backticked identifier names in Scope and the Implementation Summary are read by the qualifier as file paths (nineteen false "claimed path missing" errors); only real paths may be backticked there. `git merge` refuses while run artifacts sit staged in the index, so the artifacts commit must land before the merge, not inside it.
**Worth knowing:** A worktree's absolute path lies outside the repo and reduces to `<path outside this repository>` on the board, so its Subject must be the worktree name. The probe's measurer is a package variable only for the test seam; `TestMain` installs a plenty-free fake because three older tests assert zero findings through the real disk.

## Orientation

Now the board's VERIFY band and `queue-kanban verify` report when the disk under the repo root or a builder worktree is running out (warning under 10 GiB, critical under 3 GiB), so an unattended fan-out run shows the problem where the operator already looks; lives in the board tool's verify probes (`_dev/primes/prime-kanban-board.md`, `prime-do-kanban.md`). Not a map change: one more probe behind the existing collectVerifyFindings → attachVerifyFindings path. Prime spot-check: `prime-do-kanban.md` gained the probe's Traps bullet and its referenced paths exist; `prime-kanban-board.md` and `prime-releases.md` unchanged and current.

## Heavy Verification Plan

- Base revision: 8812a445844ce8e947986ebf3dbd415346e76df2
- Target revision: add2c162e69ca70d22b7362e121028de4dcae3e7 (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

## Heavy Verification Result

- Target revision: add2c162e69ca70d22b7362e121028de4dcae3e7
- Execution revision: add2c162e69ca70d22b7362e121028de4dcae3e7 (detached drain checkout `.git/work-run-2026-10-02-140440/drain-head`, run 14:25:43Z to 14:27:51Z, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 8s
- queue-kanban-browser: exit 0, executed (not skipped), 76s
- staged-skills: exit 0, executed, 41s

Green: every selected lane present, exit 0, none skipped, none reused.

## Timing

Observed 2026-10-02T14:04:40Z to 2026-10-02T14:28:15Z: 23m 35s total, 28m 38s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| exploration-preflight | 7m 57s | 1 |
| builder-work | 7m 12s | 1 |
| verification-gate | 6m 50s | 2 |
| review | 5m 09s | 1 |
| handback-merge | 1m 30s | 1 |

Slowest stage: exploration-preflight / triage, estimate, exploration agent, scope, pre-flight gate, 7m 57s, outcome success.

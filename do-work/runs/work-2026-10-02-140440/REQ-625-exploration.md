## Exploration

All paths are relative to `skills/do-work-board/tools/queue-kanban/` unless they start with `_dev/`, `do-work/`, or a repo-root file name.

### Concerns first (each changes what the builder writes)

- **C1. Worktree Subject does NOT "reduce to its name".** Builder worktrees live OUTSIDE the repo (`skills/do-work/actions/work-reference.md:420`: sibling `../<repo>-worktrees/worktree-agent-…`; live example: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-625-…`). `reduceAbsolutePaths` (`generate.go:689-699`) turns any path outside the repo root into `<path outside this repository>`, not its basename. So: set `Subject` to the worktree NAME (the map key from `listWorktreeAgentWorktrees`, exactly as `appendWorktreeFindings` does with `leftoverName`, `verify.go:1171`), and keep the absolute path out of Detail. Only the repo root's Subject should be the absolute `repoRoot` (reduces to `.` via `generate.go:695`).
- **C2. REQ premise about aix is wrong, about netbsd/openbsd is missing.** Cross-compiled with `go vet` (go1.26.1): `syscall.Statfs` + `Statfs_t.Bavail` compile on darwin, ios (needs CGO_ENABLED=1 to vet), linux, freebsd, dragonfly, **aix**. They do NOT compile on: **netbsd** (`undefined: syscall.Statfs`), **openbsd** (fields are `F_bavail`/`F_bsize`/`F_blocks`, not `Bavail`), illumos, solaris (`undefined: syscall.Statfs_t`). Safe unix tag: `//go:build aix || darwin || dragonfly || freebsd || ios || linux` (aix optional; keeping it unsupported is harmless). Unsupported file must cover `!<that list> && !windows`, which includes netbsd/openbsd/illumos/solaris/plan9/js/wasip1.
- **C3. Statfs_t field types differ:** darwin `Bsize uint32, Blocks/Bavail uint64`; linux `Bsize int64`; freebsd `Bsize uint64, Bavail int64` (can be negative: clamp at 0); dragonfly all `int64`; aix `Bsize uint64`. Convert every field with `uint64(...)` inside the platform file. `Stat_t.Dev` is also per-OS (int32 on darwin) — `uint64(stat.Dev)` compiles everywhere in the unix list.
- **C4. Existing zero-finding tests hit the REAL disk once the probe is live.** These assert `len(report.Findings) == 0` through `runVerifyProbes` / `collectVerifyFindings`: `verify_test.go:86` (`TestVerifyPassesOnACleanTree`), `verify_test.go:130` (`TestVerifyAllowsWriteSetOverlapWithoutMutatingRequests`), `verify_test.go:1615`. On a machine below 10 GiB free they would fail with a `low-disk-space` finding. The package has a `TestMain` at `generate_test.go:203-216`; installing a "plenty free" fake there (or in each of those tests) keeps them deterministic. No test file uses `t.Parallel()` (grep count 0 across `*_test.go`), so swapping a package-level var is race-free.
- **C5. No existing package-level `var x = func` seam exists.** Injection in this package is by PARAMETER (`buildBoard(..., lookupGitCommitDate)`, `model.go:437`; tests pass `stubGitLookup := func(string, string) (time.Time, bool) {...}`, `board_synthetic_test.go:72,207`). `gitBinaryAvailable` is a `sync.Once` struct, not swappable (`model.go:1567-1580`). The lessons file prefers parameter seams over mutable package vars (`lessons-do-kanban.md`, REQ-017 entry). The REQ explicitly asks for a package-level var, and threading a parameter through `collectVerifyFindings` would touch every caller (`generate.go:659`, `serve.go:163` via attachVerifyFindings, `verify.go:166`, tests). Recommended compromise: `appendDiskSpaceFindings(report, repoRoot, measure diskSpaceMeasurement)` takes the function as a parameter (unit tests call it directly with a fake, like `appendCompletionAnomalyFindings` is called directly at `verify_test.go:1647`), and `collectVerifyFindings` passes a package var `measureDiskSpace` that `TestMain` / end-to-end tests swap. The only existing swap-and-restore idiom is `t.Cleanup(func() { markdownToHtmlRenderer.SetParser(originalParser) })` at `generate_test.go:910`.
- **C6. The payload test checks Detail and Remedy only, not Subject** (`generate_test.go:3172-3178`). Extend the loop to `finding.Subject` too, or the worktree-path leak in C1 is not pinned.
- **C7. The REQ's Remedy contains a backticked `du -sh * | sort -h`.** No slash, so `remainingAbsolutePath` (`generate.go:714`) leaves it intact. Fine as is.
- **C8. `listWorktreeAgentWorktrees` needs git.** `appendWorktreeFindings` guards with `gitBinaryAvailable()` and skips (`verify.go:1121-1130`). The disk probe should still measure the repo root when git is missing or the listing fails; only the worktree enumeration is lost (a SkippedProbes line for that half, or silently measure root only — builder's call, but silence reads as checked-and-clean per `verify.go:105-114`).
- **C9. Line drift:** REQ cites `verify.go:195` for `appendWorktreeFindings` in `collectVerifyFindings`; it is now `verify.go:196`. `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` is at `generate_test.go:3125` as cited.

### 1. verify.go

- Imports `verify.go:3-13`: fmt, os, os/exec, path/filepath, regexp, sort, strconv, strings, time. No syscall yet; the `Stat_t` cast must live in the platform file (it does not exist on windows).
- Category block `verify.go:15-45`: `const ( verifyCategoryXxx = "kebab-case" )`, gofmt-aligned. Add `verifyCategoryLowDiskSpace = "low-disk-space"` after `verifyCategoryCalibrationRowUnreconcilable` (`verify.go:44`).
- **No registration needed anywhere else.** `boardRenderedVerifyCategories` (`generate.go:647-652`) is a SUPPRESSION list (categories the board already shows another way); a new category must NOT be added there. The CLI renderer prints any category verbatim (`verify.go:1395`). The web strip renders blindly (`web/board-cards.js:633-662`); label is `categoryWords(category)` = lowercase, dashes to spaces (`web/board-cards.js:746-748`). Repo-wide grep finds no category enumeration in Go, JS, shell or docs beyond verify.go and tests.
- Thresholds/constants style: documented named consts, e.g. `staleClaimThreshold = 3 * time.Hour` (`verify.go:72-79`), `maxNamedItemsInFindingDetail` (`verify.go:57`).
- `VerifyFinding` `verify.go:96-102` (Category, Detail, Subject, Fixable, Remedy). `VerifyReport` `verify.go:116-127` (RepoRoot, Findings, SkippedProbes, NotApplicableProbes).
- `collectVerifyFindings(repoRoot string, board *Board, now time.Time) VerifyReport` `verify.go:179-199`: builds `report := VerifyReport{RepoRoot: repoRoot}`, then 14 `appendXxx(&report, ...)` calls, last is `appendWorktreeFindings(&report, repoRoot, board)` at `verify.go:196`, then `return report`.
- `appendWorktreeFindings(report *VerifyReport, repoRoot string, board *Board)` `verify.go:1120-1220`: git guard + skip `verify.go:1121-1124`; `listWorktreeAgentWorktrees` + skip `verify.go:1126-1131`; per-name findings with `Subject: leftoverName` `verify.go:1171-1177`.
- `listWorktreeAgentWorktrees(repoRoot string) (map[string]string, error)` `verify.go:1287-1306`: name to absolute path, from `git worktree list --porcelain`, filtered on `worktreeAgentNamePrefix` (`verify.go:85`).
- SkippedProbes phrasing: lowercase `"<probe name>: <reason>"`, e.g. `"worktree probes: git is not on PATH"` (`verify.go:1122`), `fmt.Sprintf("worktree probes: %v", listError)` (`verify.go:1128-1129`), `fmt.Sprintf("worktree removability probe for %s: %v", name, err)` (`verify.go:1168-1169`). Suggested: `"disk-space probe: unsupported on <GOOS>"`, `"disk-space probe for <dir>: <err>"`. Rendered as `  - skipped <text>` (`verify.go:1403-1405`).
- Detail prose uses an em dash already (`verify.go:1200,1214`).
- **No byte/size formatter exists in the package** (grep for GiB/formatBytes/humanize: none). Closest model: `formatApproximateDuration` `verify.go:1364-1374`. Builder writes e.g. `formatGibibytes(bytes uint64) string` -> `"%.1f GiB"`.
- `joinCappedItemList(items []string) string` `verify.go:59-69`: first 5 then "and N more". Not needed for this probe.
- CLI renderer `renderVerifyReport` `verify.go:1378-1416`: `  ! <category>[ fixable]: <subject> <detail>` then `      → <remedy>`.

### 2. Tests

- Synthetic repo helper: `writeVerifyFixture(t, []verifyFixtureFile{{"rel/path", "content"}}) string` `verify_test.go:24-41` (temp dir with `do-work/queue/`). Clean base files: `cleanVersionFile` `verify_test.go:57`, `cleanChangelog` `verify_test.go:59-70`; plus a clean REQ `verify_test.go:78`.
- Board build without git: `buildBoard(repoRoot, moment, defaultRecentWindow, lookupGitCommitDate)` `verify_test.go:1960`, or a stub lookup `board_synthetic_test.go:207`. A bare `&Board{}` works for probes that only read columns (`verify_test.go:1627`).
- Category assertion helper: `findingsMentioning(report, category) []VerifyFinding` `verify_test.go:44-52`. Direct-call pattern for one probe: `report := VerifyReport{}; appendCompletionAnomalyFindings(&report, root, board)` then scan `report.SkippedProbes` with `strings.Contains` (`verify_test.go:1646-1659`). End-to-end pattern: `collectVerifyFindings(repoRoot, board, moment)` (`verify_test.go:1966-1975`).
- Git worktree fixtures if the builder wants a real two-worktree case: `newWorktreeFixtureRepo` `verify_test.go:573`, `addFixtureWorktree(t, repoRoot, parent, name)` `verify_test.go:596`. For the same-device dedupe test a fake measurement that returns one device id for two dirs is enough; no git needed if the worktree list is also injectable, otherwise use these fixtures.
- Subject rule test: `TestVerifyFindingDetailsDoNotRepeatTheirSubject` `verify_test.go:2446-2525` fails if `Detail` has `Subject` as a PREFIX (`verify_test.go:2508`). Detail starting with "<free> free of …" is safe.

### 3. Platform split

- `atomic_replace_unix.go:1`: `//go:build aix || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris`, imports only `os`.
- `atomic_replace_windows.go:1`: `//go:build windows`; `var replaceFileProcedure = syscall.NewLazyDLL("kernel32.dll").NewProc("ReplaceFileW")` (`:11`); `syscall.UTF16PtrFromString` with wrapped error (`:17-24`); `proc.Call(uintptr(unsafe.Pointer(p)), …)` (`:26-33`); success when `result != 0`; `callError == syscall.Errno(0)` means "failed without an error code" (`:34-40`). For `GetDiskFreeSpaceExW` pass three `*uint64` out-pointers (freeBytesAvailableToCaller, totalBytes, totalFreeBytes). Device identity on Windows: `GetVolumePathNameW` (kernel32) or `filepath.VolumeName(abs)`; the latter needs no syscall but misses mount points.
- `atomic_replace_unsupported.go:1`: the negation of both; returns `fmt.Errorf("... unsupported on %s", runtime.GOOS)` (`:12-14`). Make the disk version a typed sentinel (e.g. `errDiskSpaceUnsupported`) so the probe can `errors.Is` it.
- `go.mod`: requires only goldmark, x/text, yaml.v3. **`golang.org/x/sys` is NOT a dependency.** `syscall` is imported only by `atomic_replace_windows.go` and `serve.go:17` (signals). No file uses `Statfs` or `Stat_t` today.

### 4. generate.go

- `attachVerifyFindings(data *generatedBoardData, board *Board, now time.Time)` `generate.go:658-677`: calls `collectVerifyFindings(board.RepoRoot, board, now)`, skips suppressed categories, reduces Detail, Subject, Remedy and every SkippedProbes string through `reduceAbsolutePaths`. Callers: `generate.go:481` (static), `serve.go:163` (every live request, outside the mtime cache).
- `reduceAbsolutePaths(text, repoRoot)` `generate.go:689-699`: `repoRoot+sep` -> "", `repoRoot` -> `.`, any other absolute path -> `<path outside this repository>` (regex `generate.go:714`).
- `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` `generate_test.go:3125-3180`: fixture with a stale claim; builds board, `buildGeneratedBoardData`, seeds a synthetic unreduced worktree finding, calls `attachVerifyFindings`, then asserts (a) `encodeBoardDataForJsAssignment` output does not contain `repoRoot` (`:3168-3170`), (b) no Detail/Remedy matches `remainingAbsolutePath` (`:3172-3178`). Extension: swap the measurement fake to report low space for the repo root (and a worktree name) before `attachVerifyFindings`, assert a `low-disk-space` entry is present (so the test cannot pass on a dead probe), and add `finding.Subject` to the checked texts.

### 5. CLI and board rendering

Covered in section 1: nothing to register. The strip comment at `web/board-cards.js:634` says "sixteen categories" (stale prose, not a list; leave it).

### 6. Docs

- `prime-do-kanban.md:11`: `- \`verify.go\` — the read-only invariant probes and their report (wired into …forensics.md Check 14)`. There is no per-probe list; add a short clause on that line or a new Traps bullet near `prime-do-kanban.md:31-34` (the verify traps). Traps bullets use either `- [family: slug] …` (`:27`) or `- **Bold claim** — explanation`.
- `lessons-do-kanban.md`: 47 lines. Header `:1-12`. Newest version-keyed entries are prepended at the top of the list (`:10` is `- 0.303.5: **bold lesson.** prose…`, `:12` `- 0.295.1: …`). Family-marked form: `- [family: unknown-reads-as-clean] REQ-458: **bold lesson.** prose` (`:44`), `- [family: paired-predicate-drift] 0.305.9: **…**` (`:43`). New entry: `- [family: disk-space-blind-spot] 0.305.60: **…** …`, placed at the top of the list.
- `do-work/lessons-index.md:3` formula: tokens = `(wc -c + 3) / 4` integer division; families = exact sorted set of `[family: …]` markers; slugged `full` only if every bullet has a marker. Row `do-work/lessons-index.md:11`: families `paired-predicate-drift, subject-not-restated-in-detail, unknown-reads-as-clean`, tokens `6908`, `slugged: partial`. After the edit: families become `disk-space-blind-spot, paired-predicate-drift, subject-not-restated-in-detail, unknown-reads-as-clean`; recompute tokens; coverage stays `partial`. No script generates this index (grep found none).

### 7. Running the tests

- Gate lane: `_dev/tests/maintainer-verify.sh:724-736` runs the board package with `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior` through `run_budgeted_go_tests` (`:70-87`) which calls `_dev/tests/run-go-tests-with-budget.sh`; per-test-file budget 30 s (`maintainer-verify.sh:12`, `run-go-tests-with-budget.sh:12`).
- Fast builder command (from the worktree's `skills/do-work-board/tools/queue-kanban`):

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run 'Verify|DiskSpace|GeneratedVerifyPayload' .
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...
GOOS=windows GOARCH=amd64 go vet . && GOOS=linux go vet . && GOOS=freebsd go vet . && GOOS=openbsd go vet . && GOOS=netbsd go vet . && GOOS=solaris go vet . && GOOS=js GOARCH=wasm go vet .
```

  The cross-`vet` line is the only check that the tag split compiles on every platform; without it a wrong tag ships silently.

### 8. Release preimage

- `VERSION` = `0.305.59`; `skills/do-work/VERSION` = `0.305.59`; `skills/do-work/actions/version.md:5` = `**Current version**: 0.305.59`. Next patch: `0.305.60` (a new probe could argue minor; follow `_dev/primes/prime-releases.md`).
- `skills/do-work-board/tools/queue-kanban/VERSION` = `0.236.20`, stale since REQ-331 (`406c64ef`); not in the write_set and `_dev/primes/prime-kanban-board.md:12` says the tool rides the skill version. Leave it.
- `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` are byte-identical (`cmp`). Header lines 1-10 (title, two intro paragraphs, five archive links), then the newest entry at `CHANGELOG.md:12`:

```
## 0.305.59 — Preserve Indented and Commented Section Boundaries (2026-09-24)

Timing replacement and claim recovery now preserve requirements under …

- Restore section boundaries without changing the original heading or body bytes.
- Cover both writers with regression tests …
```

  Shape: `## X.Y.Z — Descriptive Title (YYYY-MM-DD)`, one lead paragraph, 2-3 bullets.

*Generated by Explore agent*

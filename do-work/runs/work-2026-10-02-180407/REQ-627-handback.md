# Hand-back — REQ-627
**Branch:** worktree-agent-REQ-627-testing-page-free-disk-space  **Base:** 953c346e  **Head:** feb206d0

## File manifest
- skills/do-work-board/tools/queue-kanban/verify.go (modified) — `VerifyReport.RepoRootDiskSpace *diskSpaceReading`, filled on the repo root's first pass in `appendDiskSpaceFindings` (healthy, `not measured on <GOOS>`, or `not measured: <err>`); new `diskSpaceLevelFor` is now the one switch over the two threshold constants, used by the finding too
- skills/do-work-board/tools/queue-kanban/generate.go (modified) — `generatedDiskSpace` type + `DiskSpace` field (`json:"diskSpace,omitempty"`), set in `attachVerifyFindings` with directory and skip reason passed through `reduceAbsolutePaths`
- skills/do-work-board/tools/queue-kanban/web/board-testing.js (modified) — pure `diskSpaceLineFor(diskSpace, liveApiAvailable)` returning `{ text, level }` plus one call in `renderTestingView()` that sets the text and level class
- skills/do-work-board/tools/queue-kanban/web/template.html (modified) — `<span class="testing-disk-space" id="testing-disk-space"></span>` in the Testing toolbar
- skills/do-work-board/tools/queue-kanban/web/board.css (modified) — `.testing-disk-space` (`--ink-soft`, right-aligned with `margin-left: auto`), `-warning` (`--accent-pending`), `-critical` (`--accent-blocked`)
- skills/do-work-board/tools/queue-kanban/disk_space_test.go (modified) — 3 tests
- skills/do-work-board/tools/queue-kanban/generate_test.go (modified) — extended `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths`, added `TestGeneratedVerifyPayloadReducesTheDiskSpaceSkipReason`
- skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go (modified) — `TestJavaScriptBehaviorTestingDiskSpaceLineReadsEveryCase`
- skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (modified) — one sentence on the disk-probe Traps bullet

## Integration seams
None. `serve.go` already calls `attachVerifyFindings` on every request outside the mtime cache, so the field is recomputed per request with no change there.

## P-A-U
- [PLAN] Keep the repo root's first measurement on the verify report before the skip, dedupe and threshold logic, so one measure call feeds both the finding and the readout. Project it in `attachVerifyFindings` (the single home both static and serve call) with the level from a Go helper over the existing constants. Render it with a small sliceable JS helper keyed off the existing `liveTestingApi` flag for the "(at generation)" label.
- [APPLY] Only the 9 Scope files changed. No edits to serve.go, platform measurement files, board-controls.js, release files, go.mod/go.sum or anything under do-work/.
- [UNIFY] `git diff --stat HEAD~1`: 9 files changed, 264 insertions(+), 8 deletions(-). `gofmt -l .` empty; `go vet .` clean; cross-vet windows/amd64, linux, openbsd, js/wasm clean. I re-read each diff: finding text is unchanged (`thresholdName` still "critical"/"warning", same Detail format), and skip lines and their order are unchanged. Light-theme colours come from the tokens' existing light overrides, so no extra light rule was needed. No debug output was left in the diff.

## Red-green evidence
RED against a compiling stub (field and type present, never set; `diskSpaceLevelFor` returned ""; JS helper returned `{ text: "", level: "" }`). The stub was not committed.
- TestDiskSpaceProbeKeepsTheHealthyRepoRootReading: `RepoRootDiskSpace = <nil>, want {directory:/fixture/repo freeBytes:429496729600 totalBytes:536870912000 skipReason:}` → GREEN at feb206d0
- TestDiskSpaceProbeKeepsTheRepoRootSkipReason: `unsupported RepoRootDiskSpace = <nil>, want {... skipReason:not measured on darwin}` / `failed RepoRootDiskSpace = <nil>, want {... skipReason:not measured: permission denied}` → GREEN at feb206d0
- TestDiskSpaceLevelMatchesTheFindingThresholds: `diskSpaceLevelFor(2147483648) = "", want "critical"` (and 4 more boundary lines) → GREEN at feb206d0
- TestGeneratedVerifyPayloadCarriesNoAbsolutePaths: `payload DiskSpace = <nil>, want {FreeBytes:2147483648 TotalBytes:536870912000 FreeText:2.0 GiB TotalText:500.0 GiB Level:critical Directory:. SkipReason:}` → GREEN at feb206d0
- TestGeneratedVerifyPayloadReducesTheDiskSpaceSkipReason: `payload DiskSpace = <nil>, want {... Level:neutral Directory:. SkipReason:not measured: statfs .: permission denied}` → GREEN at feb206d0
- TestJavaScriptBehaviorTestingDiskSpaceLineReadsEveryCase: `healthyLive: got {Text: Level:}, want {Text:disk: 50.0 GiB free of 500.0 GiB Level:neutral}` (and the same for healthyStatic, warningLive, criticalLive, skippedStatic, absentLive) → GREEN at feb206d0
- The 2 GiB REQ-625 finding is still pinned by the existing tests `TestDiskSpaceProbeReportsLowFreeSpaceAtEachThreshold`, `TestDiskSpaceProbeReachesTheVerifyReport` and `TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory`. All pass unchanged.

## Tests run
All runs used QUEUE_KANBAN_BROWSER_PROBES=off, from skills/do-work-board/tools/queue-kanban.
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=off go test -count=1 -run 'DiskSpace|GeneratedVerifyPayload|Verify' .` → exit 0 (`ok 6.218s`), 7.1s wall
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -run TestJavaScriptBehavior .` → exit 0 (`ok 5.951s`), 6.2s wall
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=off go test -count=1 ./...` (-json) → exit 0, 643 passing tests/subtests, 0 failures, package 53.9s wall. The slowest test is the pre-existing TestLegacyStrictEntryPointsRejectZeroProbes at 19.2s. None of the new tests is over 2s.
- `gofmt -l .` → empty; `go vet .` → exit 0
- `GOOS=windows GOARCH=amd64 go vet . && GOOS=linux go vet . && GOOS=openbsd go vet . && GOOS=js GOARCH=wasm go vet .` → exit 0
- `go run . generate --out <scratch>/static627 --repo-root <main repo>` → wrote the static board (604 REQs). `board-data.js` carries:
  `"diskSpace":{"freeBytes":65860677632,"totalBytes":190601715712,"freeText":"61.3 GiB","totalText":"177.5 GiB","level":"neutral","directory":"."}`
  The field contains no absolute path. `board-data.js` does contain `/Users/t2/...` strings, but only inside rendered REQ/UR body text that was already there before this change, not in any verify or disk field.
- Not run: the browser/visual check. The line was not looked at in a real browser (browser probes were off per the brief).

## Decisions
- D-01: A skipped or failed repo-root reading ships with `level: "neutral"`, and free, total and their text fields are left out. The JS helper also falls back to neutral when the field is absent. Reason: "not measured" is not a low-space alarm, and the probe's own skip line in the VERIFY band already marks it as unverified. DECIDE & STATE.
- D-02: The finding's threshold switch now goes through `diskSpaceLevelFor`, so the readout level and the finding level cannot drift apart. The finding output is unchanged: the "neutral" branch maps to the old `default: continue`. DECIDE & STATE.
- D-03: `Directory` is always set, including on skip, so the payload still names what was measured (`.`). DECIDE & STATE.
- D-04: The line is right-aligned in the toolbar (`margin-left: auto`) so it does not move the tester controls. It wraps with the toolbar on narrow widths. DECIDE & STATE.

## Discovered Tasks
None

## Lessons read
- skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md — family-targeted (`disk-space-blind-spot` bullet only)

## Blockers
None

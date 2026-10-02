## Review

**Overall: 96%** | 2026-10-02T18:47:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. The Testing page shows the repo root's free and total space on one always-visible line, fed by the REQ-625 probe's single measurement, coloured by the same thresholds, labelled "(at generation)" on a static snapshot, and never carrying an absolute path.

**Requirements walk (merge range 877a7cd3..172d0e5c, 9 files, all inside the declared Scope):**
- [x] Payload field from the same measurement: `appendDiskSpaceFindings` keeps the repo root's first (`targetIndex == 0`) reading on `VerifyReport.RepoRootDiskSpace`, and `attachVerifyFindings` projects it to `generatedBoardData.DiskSpace`. One `measure` call per directory, unchanged. `diskSpaceMeasurer` and `formatGibibytes` are reused.
- [x] serve recomputes per request: `serve.go:163` already calls `attachVerifyFindings` outside the mtime cache. Verified live: three `GET /board-data.js` calls returned three different `freeBytes` (65432760320, 65485340672, 65480245248).
- [x] Static snapshot labelled: `diskSpaceLineFor` appends " (at generation)" when `liveTestingApi` is false. Seen in a browser on the generated static board.
- [x] One line, GiB with one decimal, neutral / warning (<10 GiB) / critical (<3 GiB), "not measured on <GOOS>" when skipped: the line is set on every `renderTestingView` call, the level comes from Go's `diskSpaceLevelFor`, and the skip text comes through `skipReason`.
- [x] No absolute path: `Directory` and `SkipReason` both pass through `reduceAbsolutePaths`. The healthy, unsupported and error paths are covered by `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` (extended) and `TestGeneratedVerifyPayloadReducesTheDiskSpaceSkipReason`. The live serve payload has `"directory":"."` and no `/Users` in the field.
- [x] REQ-625 finding behaviour unchanged: the finding text (`thresholdName` is still "critical"/"warning"), the skip lines, their order, the early `return` on unsupported, the `continue` on error, and the device dedupe are byte-identical in effect. "neutral" maps to the old `default: continue`. The existing REQ-625 tests pass unchanged.

**Restatement Sweep:** the diff redefines the use of the disk probe's result (healthy readings are now kept) and adds a payload field. Grepped `disk-space`, `low-disk`, `disk space`, `diskSpace`, `verifyFindings`/`verifySkipped`, `Testing view`/`Testing page` across shipped docs. `prime-do-kanban.md` was updated in the diff. `forensics.md` maps verify output by output class, not by probe list, so it is unaffected. `board.md:13` ("serve mode rebuilds from disk on every browser reload … does not push updates to an open tab") agrees with the new line's behaviour. No stale restatement was found. The one gap is an omission, recorded as M1.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1: `skills/do-work-board/docs/board-guide.md` § Testing view (line 72) and `skills/do-work-board/actions/board.md` Step 6 describe the Testing toolbar but do not mention the new disk line. This is an omission, not a contradiction, so a reader is not misled. — impact-negligible → report only
- M2: In serve mode the line shows the free space measured when the page loaded and has no time hint, so a tab left open for hours shows an old figure. The Testing view re-renders after testing actions, but `boardData` is not refetched. This matches the REQ (serve recomputes on every request, which means every `board-data.js` fetch) and the documented reload-to-refresh contract in `board.md:13`, so it is acceptable as built. — impact-negligible → report only

**Nit findings:**
- N1: `verify.go` `appendDiskSpaceFindings` now maps the level string back to threshold bytes with a second string-keyed `switch` ("critical"/"warning"). This works, but the threshold pair now lives in two switches keyed on magic strings. — impact-negligible → report only
- N2: The DOM wiring in `renderTestingView` (text plus `testing-disk-space-<level>` class on `#testing-disk-space`) is not pinned by an automated test. Only the pure `diskSpaceLineFor` helper is pinned. I verified the wiring by hand in a browser. — impact-negligible → report only

**Acceptance:** Pass. The focused tests pass, a static board and a live `serve` instance were checked in a real browser, and the line's text, position, colour (light and dark) and narrow-width wrap are all correct.

Evidence:
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run 'DiskSpace|GeneratedVerifyPayload|TestJavaScriptBehaviorTestingDiskSpace' -v .` → `ok` (12 tests PASS, including the 6 REQ-625 disk tests unchanged and the 6 new or extended ones).
- Static board (`go run . generate`, served over a local HTTP server, `#testing`): the text was `disk: 61.0 GiB free of 177.5 GiB (at generation)`, class `testing-disk-space-neutral`, right-aligned at the end of the toolbar after the read-only note. Light colours were neutral rgb(89,99,111), warning rgb(143,94,16), critical rgb(189,81,56). Dark colours were neutral rgb(153,161,173), warning rgb(216,162,74), critical rgb(217,122,89) on rgb(19,23,32). The only console error was a favicon 404.
- At 360px width the line wraps to its own row inside the toolbar, and the toolbar does not overflow (scrollWidth equals clientWidth, 311). The page's 612px document width at 360px is pre-existing: it comes from the header (`board-identity`, `board-controls`) and is the same with the line hidden.
- Live `serve` on 127.0.0.1:18627: the text was `disk: 60.9 GiB free of 177.5 GiB`, with no "(at generation)" and the read-only note hidden. The payload is recomputed per request, as shown under Requirements.

**Suggested testing:** 3 items
- Manual: run the live board with a fake low reading (or on a nearly full volume) and confirm the warning and critical colours read clearly against the toolbar in both themes. I forced the classes via DevTools rather than using real low readings.
- Platform: on Windows or an unsupported unix, confirm the line reads `disk: not measured on <GOOS>` and is not empty.
- Regression: with a builder worktree on a low device and a healthy repo root, confirm the VERIFY band still shows only the worktree's finding while the Testing line stays neutral (the two surfaces measure different directories).

**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

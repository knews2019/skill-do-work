---
id: REQ-627
title: 'Show free disk space on the Testing page'
status: claimed
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
  - Route B
  - 7-file write set
  - 6 acceptance criteria
  calculated_at: 2026-10-02T18:09:53Z
created_at: 2026-10-02T18:02:37Z
user_request: UR-132
domain: frontend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-626"]
batch: board-links-and-disk-space
required_lessons: ["skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md#disk-space-blind-spot"]
write_set: ["skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/generate.go", "skills/do-work-board/tools/queue-kanban/web/board-testing.js", "skills/do-work-board/tools/queue-kanban/web/template.html", "skills/do-work-board/tools/queue-kanban/web/board.css", "skills/do-work-board/tools/queue-kanban/disk_space_test.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
review_at: 2026-10-02T18:47:04Z
integration_at: 2026-10-02T18:44:09Z
builder_handback_at: 2026-10-02T18:43:58Z
dispatch_at: 2026-10-02T18:39:40Z
heavy_verified_at: 2026-10-02T18:50:16Z
heavy_verified_revision: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086
claimed_at: 2026-10-02T18:03:42Z
status_changed_at: 2026-10-02T18:36:46Z
commit: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086
---

# Show Free Disk Space on the Testing Page

## What
The Testing page always shows one line with the repo root's free and total disk space, for example `disk: 50.2 GiB free of 177.5 GiB`, fed by the measurement REQ-625 (the low-disk-space verify probe) already takes, and coloured by the same thresholds (warning under 10 GiB, critical under 3 GiB).

## Why
REQ-625 shipped in 0.305.60, but a verify finding appears only below the thresholds, so on a healthy machine the free space is shown nowhere. The user wants to see the number before it becomes a finding, on the Testing page.

## Detailed Requirements
- Carry the healthy measurement into the board payload: a small field on `generatedBoardData` (free bytes, total bytes, and the measured directory reduced through `reduceAbsolutePaths`, or a skip reason when the platform cannot measure). Reuse `diskSpaceMeasurer` and `formatGibibytes` from verify.go; do not measure twice if one call can feed both the probe and the payload.
- `serve` recomputes the value on every request, outside the mtime cache, the same way the verify findings are (free space changes while no file changes).
- The static snapshot carries the value measured at generation time, labelled as such (for example "at generation").
- The Testing page renders one always-visible line near the top: free of total in GiB with one decimal, coloured neutral above 10 GiB, warning below 10 GiB, critical below 3 GiB. When the measurement was skipped, the line says so ("disk: not measured on <GOOS>") rather than disappearing.
- No absolute path reaches the page. The existing payload test `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` or a sibling covers the new field.
- The probe's finding behaviour from REQ-625 is unchanged; this REQ adds a readout, not a second threshold.

## Constraints
No new dependency. One syscall per request. Nothing is deleted. Shipped files change, so this is a release. The field is display-only; nothing schedules or gates on it.

## Dependencies
None. Builds on REQ-625, which is archived (commit add2c162); independent of REQ-626 (Give every board page and lens its own URL).

## Builder Guidance
High certainty on placement (Testing page, user decision) and thresholds. Builder latitude on the payload field name and the exact markup. A Go test with a fake measurer plus a JavaScript behaviour test in the existing harness are the expected proof.

## Red-Green Proof
**RED prompt/case:** Open the Testing page on a machine with 50 GiB free: no free-space figure appears anywhere; the generated board payload has no disk-space field.
**Why RED now:** REQ-625 keeps only findings below threshold; the healthy measurement is discarded.
**GREEN when:** The Testing page shows `disk: <free> GiB free of <total> GiB` on a healthy machine (neutral colour); with a fake measurer at 9 GiB the same line turns warning and at 2 GiB critical; with the unsupported sentinel it reads "not measured"; the payload carries the field with no absolute path.
**Validation:** User adjusted (page chosen at capture)

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7095 tokens, over the 2000 budget; `slugged: partial`). Narrowed at claim time to its `disk-space-blind-spot` family, which is now in `required_lessons`.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, over budget; `slugged: partial`). Matched: static output.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: keep the repo root's first measurement on the verify report, project it in the shared verify-attach function with a Go-computed level, render it with a small sliceable helper keyed off the live-testing flag.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: only the nine Scope files changed.
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: 9 files, 264 insertions, 8 deletions; gofmt empty, go vet and four cross-vets clean; finding text and skip-line order re-read unchanged.

## Full Context
See `do-work/user-requests/UR-132/input.md` for complete verbatim input.

*Source: in one of the pages list the free space too since it's already collected*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome, placement and thresholds are fixed by the REQ, but how the REQ-625 probe's measurement can feed both the finding and a payload field in one call, where `serve` recomputes verify findings outside its mtime cache, how the static snapshot is labelled, and how the Testing page renders need discovery before dispatch. No architectural change, so no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors: `do-work/runs/work-2026-10-02-180407/REQ-627-exploration.md`. The findings that change what the builder writes:

- **One measurement can feed both outputs.** The disk probe already measures the repo root first and simply discards a healthy result. Keeping that first measurement (or its skip reason) on the verify report costs no second call.
- **The payload hook already runs in both modes.** The verify-findings attach function in `skills/do-work-board/tools/queue-kanban/generate.go` is called by static generation and, on every serve request, outside the mtime cache. A field set there reaches both, so `skills/do-work-board/tools/queue-kanban/serve.go` needs no change and leaves the declared write set.
- **The static flag already exists.** The Testing page already reads the live-testing-API flag from the payload; "at generation" keys off it, and the payload already carries the generation time.
- **Colour needs CSS.** No warning or critical class exists for the Testing toolbar; amber and red colour tokens exist in `skills/do-work-board/tools/queue-kanban/web/board.css`, which joins the Scope. The line's element goes in the Testing toolbar in `skills/do-work-board/tools/queue-kanban/web/template.html`.
- **Tests:** a fake-measurer helper exists in `skills/do-work-board/tools/queue-kanban/disk_space_test.go`; the package's test main installs a 400 GiB fake; the no-absolute-paths payload test scans only finding text, so the new field needs an explicit check. No behaviour test drives the real Testing render, so the line's text and class come from a small helper function the Node probe can slice out.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — keep the repo root's measurement or skip reason on the verify report
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — the payload field, set where verify findings are attached, with the directory reduced
- `skills/do-work-board/tools/queue-kanban/web/board-testing.js` (modify) — render the disk line with its threshold class and the generation label
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — the line's element in the Testing toolbar
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — neutral, warning and critical styles for the line
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (modify) — the report keeps the healthy, unsupported and failed repo-root measurement
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — payload carries the field with no absolute path
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (modify) — line text and class at 50, 9 and 2 GiB, unsupported, and the static label
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify) — the disk-probe Traps bullet names the payload field

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/serve.go` (its per-request attach call already carries the field), the platform measurement files, release paths, anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] The payload carries free bytes, total bytes and the reduced directory, or a skip reason, from the same measurement the probe takes
- [ ] serve recomputes the value on every request, outside the mtime cache
- [ ] The static snapshot carries the generation-time value, labelled as such
- [ ] The Testing page shows one always-visible line, free of total in GiB with one decimal, neutral at or above 10 GiB, warning below 10, critical below 3, and "not measured on <GOOS>" when skipped
- [ ] No absolute path reaches the page, covered by a payload test
- [ ] The probe's finding behaviour is unchanged

## Pre-Flight

**Git:** ✓ Integration tip 953c346e on `main`, which contains REQ-626's merge 505ec74a; the only dirt is this REQ's own `do-work/` trail (working REQ, brief) — no third-party paths
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` (DO_WORK_FAST_STAGE_REUSE=off) exit 0 at 953c346e, both Go stages EXECUTING, gate wall 135s
**Dependencies:** ✓ Go toolchain and node present; no new dependency planned

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-testing.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** The disk-space probe now keeps the repo root's first measurement (or why it could not be taken) on the verify report, so one measurement feeds both the REQ-625 finding and a new board payload field carrying free and total bytes, their GiB text, a Go-computed level (neutral, warning below 10 GiB, critical below 3 GiB) and the path-reduced directory. The Testing page shows it on one always-visible toolbar line, coloured by level, marked "(at generation)" on a static snapshot and "not measured" when skipped. The threshold switch now lives in one helper that both the finding and the level read. Merge range 877a7cd3..172d0e5c (builder commit feb206d0, merge 172d0e5c).

## Qualification

**Diff range:** 877a7cd3..172d0e5c (builder commit feb206d0, merge 172d0e5c)
**Gate records:** qualify satisfied; scope-drift satisfied (nine files, exactly the declared Scope).
**Warnings judged:** none.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. The repo root is measured once and its reading (or skip reason) is kept on the report before the device dedupe and threshold switch; the finding text, skip lines and their order are unchanged, and the threshold switch moved into one level helper that both the finding and the payload read. The payload field is set in the shared verify-attach function, which static generation calls once and serve calls on every request outside the mtime cache, so serve needed no change. Directory and skip reason pass through the existing path reduction. The page line is always set, neutral when skipped or absent, with "(at generation)" keyed off the existing live-testing flag. No new dependency; nothing deleted.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY was cross-checked against `git diff --stat 877a7cd3..172d0e5c` (nine Scope files, nothing under do-work/).
**Live data flow:** the payload field reaches the page through the same board-data script both modes load; the builder's static generation of the real repo shows `"diskSpace":{... "freeText":"61.3 GiB","totalText":"177.5 GiB","level":"neutral","directory":"."}`.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (DO_WORK_FAST_STAGE_REUSE=off) on the merged tree at 172d0e5c
**Result:** ✓ All passing — exit 0, gate wall 144s; stage queue-kanban-fast-tests EXECUTING (413 tests, wall 43s, slowest file 20.62s < 30s); stage do-work-cli-fast-tests EXECUTING (866 tests, wall 73s, slowest file 23.22s < 30s). Green-gate record satisfied by advance.

**Focused tests:** `go test -count=1 -run 'DiskSpace|GeneratedVerifyPayload|TestJavaScriptBehaviorTestingDiskSpace' .` in the board package with JavaScript probes on and browser probes off → exit 0 (advance probe record satisfied).

**Builder runs (from the hand-back, branch feb206d0):** the focused Go run (7.1s), the JavaScript behaviour lane (6.2s), the whole package with probes off (643 passing, 53.9s), gofmt, go vet and cross-vets for windows, linux, openbsd and js/wasm; a static board of this repo carried `"diskSpace":{"freeBytes":65860677632,"totalBytes":190601715712,"freeText":"61.3 GiB","totalText":"177.5 GiB","level":"neutral","directory":"."}`.

**Red-green validation:** traced to `## Red-Green Proof`; RED was taken against a compiling stub, GREEN at builder commit feb206d0:
- TestDiskSpaceProbeKeepsTheHealthyRepoRootReading: ✗ `RepoRootDiskSpace = <nil>` → ✓
- TestDiskSpaceProbeKeepsTheRepoRootSkipReason: ✗ `unsupported RepoRootDiskSpace = <nil>, want {... skipReason:not measured on darwin}` → ✓
- TestDiskSpaceLevelMatchesTheFindingThresholds: ✗ `diskSpaceLevelFor(2147483648) = "", want "critical"` → ✓
- TestGeneratedVerifyPayloadCarriesNoAbsolutePaths (extended): ✗ `payload DiskSpace = <nil>` → ✓
- TestGeneratedVerifyPayloadReducesTheDiskSpaceSkipReason: ✗ `payload DiskSpace = <nil>` → ✓
- TestJavaScriptBehaviorTestingDiskSpaceLineReadsEveryCase (50 GiB neutral, warning, critical, static label, skipped, absent): ✗ `healthyLive: got {Text: Level:}` → ✓
- REQ-625's finding behaviour stays pinned by its unchanged tests, which pass.

**New tests added:**
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (three tests above)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (one new payload test)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (one behaviour test)

**Existing tests updated (cross-REQ impact):**
- `skills/do-work-board/tools/queue-kanban/generate_test.go` TestGeneratedVerifyPayloadCarriesNoAbsolutePaths (from REQ-625): now also asserts the disk field and its reduced directory — intentional

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 877a7cd3..172d0e5c
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

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

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-02-180407/REQ-627-handback.md`). All are DECIDE & STATE.

- D-01: A skipped or failed repo-root reading ships with the neutral level and no byte or text fields; the page helper also falls back to neutral when the field is absent. "Not measured" is not a low-space alarm, and the probe's own skip line already marks it unverified in the VERIFY band.
- D-02: The finding's threshold switch now goes through the new level helper, so the readout level and the finding level cannot drift; the finding output is unchanged.
- D-03: The directory is always set, including on skip, so the payload names what was measured.
- D-04: The line is right-aligned in the toolbar so it does not move the tester controls; it wraps with the toolbar on narrow widths.

Orchestrator decision:
- D-05: After a machine restart, `recover --take-over` returned this claim to the queue and stripped its Triage, estimate, Plan, Exploration and Scope sections. The REQ was re-claimed through queue-mode advance and those sections were restored byte-for-byte from the handoff commit 1794c268; only the new claim timestamps were kept. DECIDE & STATE.

## Discovered Tasks

From the builder's hand-back (none) and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- M1: the board guide's Testing section and the board action's Step 6 do not mention the disk line. — impact-negligible → report only
- M2: in serve mode the line shows the value from page load until reload, matching the board's reload-to-refresh contract. — impact-negligible → report only
- N1: the finding maps the level string back to threshold bytes in a second string-keyed switch. — impact-negligible → report only
- N2: the Testing view's DOM wiring of the line is not pinned by a test; the pure helper is. — impact-negligible → report only

## Lessons Learned

**What worked:** Keeping the first measurement the probe already takes instead of measuring again; moving the threshold switch into one level helper made the finding and the readout share one definition. Computing the level in Go kept the thresholds out of JavaScript.
**What didn't:** An unquoted shell heredoc used to append REQ prose ran a backticked command name as a command substitution and silently dropped it from the text; append prose with a quoted heredoc or a script.
**Worth knowing:** The repo root's reading is always the first measured target, so it is set even when the platform is unsupported (the probe returns right after). The Testing line reflects the page load in serve mode; the payload itself is fresh on every request.

## Orientation

Now the board's Testing page always shows the repo root's free and total disk space, coloured by the same thresholds as the VERIFY-band finding; lives in the board tool's verify probe and Testing view (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). Not a map change: one more payload field from the existing verify-attach path. Prime spot-check: `prime-do-kanban.md`'s disk-probe Traps bullet now names the payload field and its referenced paths exist; `prime-kanban-board.md` and `prime-releases.md` unchanged and current.

## Heavy Verification Plan

- Base revision: 877a7cd34ea897e9b948e47af42312ac160fe37d
- Target revision: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086 (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

## Heavy Verification Result

- Target revision: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086
- Execution revision: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086 (detached drain checkout `.git/work-run-2026-10-02-180407/drain-head`, one run for REQ-626 and REQ-627, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 7s
- queue-kanban-browser: exit 0, executed (not skipped), 78s
- staged-skills: exit 0, executed, 40s

Green: every selected lane present, exit 0, none skipped, none reused.

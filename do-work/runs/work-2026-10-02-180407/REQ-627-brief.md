# Builder brief — REQ-627: Show free disk space on the Testing page

**Worktree (your working directory, the only tree you write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-627-testing-page-free-disk-space`
**Branch:** `worktree-agent-REQ-627-testing-page-free-disk-space`, based on the main tip at dispatch (it already contains REQ-626's merge 505ec74a). Commit on this branch only, message prefix `[REQ-627]`. Never push.
**Route:** B. **TDD: yes** (RED before GREEN, evidence required). **Impact:** user-visible. **Estimate P50:** 25 active minutes.
**Hand-back file (the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-02-180407/REQ-627-handback.md`. Never stage or commit it.

## Read first (absolute paths; the REQ and exploration live in the MAIN tree, read-only)

1. The REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-627-testing-page-free-disk-space.md` — What, Detailed Requirements, Red-Green Proof, Exploration, Scope. The `## Scope` "Files I will touch" list is your write boundary.
2. Exploration with file:line anchors: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-02-180407/REQ-627-exploration.md`.
3. In your worktree: `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`. Required lesson (read the bullet lines marked `[family: disk-space-blind-spot]` only): `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`.
4. Crew rules in your worktree under `skills/do-work/crew-members/`: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `frontend.md`, `testing.md`. Naming rule: two words minimum for anything with reach.

## Fixed design (decided by triage and exploration; do not re-open)

- **One measurement, kept.** In `verify.go`, give `VerifyReport` one new field holding the repo root's reading: free bytes, total bytes, the measured directory (the absolute repo root, reduced later), and a skip reason string. `appendDiskSpaceFindings` fills it from the FIRST target (the repo root) — healthy or not — before the device dedupe and threshold switch. On the unsupported sentinel the skip reason is `not measured on <GOOS>`; on another repo-root error it is `not measured: <err>` (the error text is reduced later). The finding behaviour of REQ-625 must be byte-identical: same findings, same skip lines, same order. No second measure call anywhere.
- **Thresholds have one home.** The level (`neutral` at or above 10 GiB, `warning` below 10 GiB, `critical` below 3 GiB, strict `<`, the existing two constants) is computed in Go, not in JavaScript. Use the existing `formatGibibytes` for the two display strings.
- **Payload field** on `generatedBoardData` in `generate.go`, set inside `attachVerifyFindings` (the single home both static generation and every serve request already call, outside the mtime cache — so `serve.go` needs no change): e.g. `DiskSpace *generatedDiskSpace json:"diskSpace,omitempty"` with `freeBytes`, `totalBytes`, `freeText`, `totalText`, `level`, `directory` (through `reduceAbsolutePaths`, so the repo root reads `.`), and `skipReason` (through `reduceAbsolutePaths`). Two-word Go names; camelCase JSON tags like the neighbouring fields.
- **Testing page line.** One always-visible element in the Testing toolbar in `web/template.html` (e.g. `<span class="testing-disk-space" id="testing-disk-space"></span>`). In `web/board-testing.js`, a small pure helper (e.g. `diskSpaceLineFor(diskSpace, liveApiAvailable)` returning `{ text, level }`, no braces inside string literals so the Node harness can slice it) and one call from `renderTestingView()` that sets the text and a level class. Text: `disk: 50.2 GiB free of 177.5 GiB`; on a static snapshot (live testing API flag false) append ` (at generation)`; skipped: `disk: <skipReason>` (e.g. `disk: not measured on plan9`); field absent: `disk: not measured`. Never hidden.
- **Styles** in `web/board.css`: neutral uses the existing soft ink token, warning the amber pending token, critical the blocked red token (see exploration §4); follow how `.testing-toolbar-note` and `.testing-toolbar-error` are written, including the light-theme overrides if those tokens have them.
- **Docs:** extend the existing disk-probe Traps bullet in `prime-do-kanban.md` with one sentence: the repo root's reading also reaches the payload field and the Testing page line, from the same measurement.

## Tests (TDD; each named for the failure it pins)

- `disk_space_test.go`: the report keeps a healthy repo-root reading (400 GiB free → field filled, zero findings); the unsupported sentinel → skip reason `not measured on <GOOS>`, zero findings, one skip line as before; a repo-root error → skip reason set, worktrees still measured; the 2 GiB case still yields exactly the REQ-625 finding (behaviour unchanged).
- `generate_test.go`: extend or add a sibling of `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` that asserts the payload's disk field is present with directory `.`, level from the fake, and no absolute path in any of its strings; also an error-path skip reason containing the absolute repo root comes out reduced.
- `javascript_behavior_d_test.go`: `TestJavaScriptBehavior…` driving the real helper sliced from the generated page: 50 GiB → neutral text exactly `disk: 50.0 GiB free of 500.0 GiB` (or whatever your fake gives); level `warning` and `critical` pass through; static snapshot adds `(at generation)`, live does not; skipped and absent cases read as above.
- RED first against a compiling stub (field exists, never set; helper returns empty), keep the failing output verbatim, do not commit the stub.

## Never touch

Anything under `do-work/` in either tree except your hand-back file; `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`; `serve.go`; the platform measurement files `disk_space_unix.go`, `disk_space_windows.go`, `disk_space_unsupported.go`; `web/board-controls.js`; generated build output; `go.mod`/`go.sum`. A file outside Scope that you believe you need: stop and report it in the hand-back instead of editing it.

## Testing (record each command's own exit line)

From `skills/do-work-board/tools/queue-kanban` in your worktree:

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run 'DiskSpace|GeneratedVerifyPayload|Verify' .
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run TestJavaScriptBehavior .
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...
gofmt -l . ; go vet .
GOOS=windows GOARCH=amd64 go vet . && GOOS=linux go vet . && GOOS=openbsd go vet . && GOOS=js GOARCH=wasm go vet .
```

Then generate a static board from the main repo into a scratch directory with the tool's own `generate`/`static` command (see the prime) and confirm `board-data.js` carries the field with directory `.` and no absolute path; quote the field in the hand-back.

## Hand-back format (write the file, then reply with ONE line under 300 characters)

```
# Hand-back — REQ-627
**Branch:** worktree-agent-REQ-627-testing-page-free-disk-space  **Base:** <hash>  **Head:** <hash>
## File manifest
- path (new|modified|deleted) — one line each
## Integration seams
None, or the exact line and the file it belongs in
## P-A-U
- [PLAN] <two or three sentences of the approach you took>
- [APPLY] <confirmation scope stayed inside the declared files>
- [UNIFY] <git diff --stat output, linters run, files checked>
## Red-green evidence
- <test name>: RED output (verbatim, trimmed) → GREEN at <commit>
## Tests run
- <command> → exit <n>, <n> tests, wall time
## Decisions
- D-01: ... (reasoning; mark DECIDE & STATE or ESCALATE with Value/Risk)
## Discovered Tasks
- ... or None
## Lessons read
- each lessons file and whether whole-satellite or family-targeted
## Blockers
- None, or what stopped you
```

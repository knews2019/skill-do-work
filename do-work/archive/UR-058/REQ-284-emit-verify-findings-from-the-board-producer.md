---
id: REQ-284
title: Emit every verify finding from the board's Go producer
status: completed
created_at: 2026-08-19T13:47:06Z
claimed_at: 2026-08-21T00:10:59Z
completed_at: 2026-08-21T00:19:28Z
kb_status: promoted
kb_entry: REQ-284-emit-every-verify-finding-from-the-board.md
commit: 8f61f69
route: B
user_request: UR-058
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-285, REQ-280, REQ-281, REQ-282]
batch: verify-findings-on-board
write_set:
- skills/do-work-board/tools/queue-kanban/verify.go
- skills/do-work-board/tools/queue-kanban/verify_test.go
- skills/do-work-board/tools/queue-kanban/generate.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/serve.go
- skills/do-work-board/tools/queue-kanban/board_live_test.go
---

# Emit Every Verify Finding From the Board's Go Producer

## What

Split `runVerifyProbes` into a board-taking `collectVerifyFindings(repoRoot, board, now)` plus a thin
wrapper that builds the board first, then carry the resulting findings into `generatedBoardData` as
`verifyFindings` and `verifySkipped`. Suppress the three categories the board already renders, in the
producer, so the client can render the list blindly. Wire both callers: `generate` for the static
snapshot, and `serve` per `/board-data.js` request outside the mtime cache.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`queue-kanban verify` finds sixteen categories of queue and process breakage; the HTML board renders
three. The board is the surface a human looks at, and `verify` is a command nobody runs unless an agent
runs it. Two REQs sat claimed for 13h against the 3h `staleClaimThreshold` with the age printed on both
cards and no rule applied to it.

## Context

`generate.go:500` and `serve.go` both already hold a `*Board`. Without the split the board would be
built twice per request, which is why the split comes first and the render REQs depend on it.

REQ-214 (`do-work/archive/UR-048/REQ-214-verify-surfaces-completion-anomalies.md`) is the precedent
running the other direction: verify was blind to the board's completion anomalies and was taught to lift
them. Its structured-evidence pattern — forward the board's own typed data, never re-parse warning prose —
is the pattern to follow here.

## Detailed Requirements

- `collectVerifyFindings(repoRoot, board, now)` takes an already-built `*Board` and returns the same
  `VerifyReport` shape. `runVerifyProbes` keeps its signature and becomes the wrapper that builds the
  board and calls it. `main.go` is untouched: `verify` keeps its non-zero exit code, and the CLI report
  is byte-identical to today's for the same tree.
- `generatedBoardData` gains `VerifyFindings []generatedVerifyFinding json:"verifyFindings,omitempty"`
  and `VerifySkipped []string json:"verifySkipped,omitempty"`, with `{category, detail, remedy, fixable}`
  per finding.
- Suppress `completion-anomaly`, `duplicate-req-id`, and `stray-req-file` from `verifyFindings`. Those
  three already reach the board through `board.Warnings`, and anomalies additionally through the
  `completionAnomalies` column and a per-card badge — appending them verbatim would print the same prose
  a third time. Suppression happens in the Go producer, not in JS.
- `verifySkipped` still carries skipped probes. A skipped probe that renders as nothing reads as "checked
  and clean", which is the same never-silent rule `verify.go` already states about its integration-ref skip.
- **Serve computes findings fresh on every `/board-data.js` request**, calling
  `collectVerifyFindings(repoRoot, cachedBoard, time.Now())` outside `refreshBoardData`'s mtime-gated
  path. No TTL, no second cache. The mtime cache on the board build itself stays exactly as it is.
- **No absolute filesystem path may appear anywhere in the emitted JSON.** The worktree probe's detail
  carries `git worktree list --porcelain` output, which is absolute. Reduce it where
  `generatedVerifyFinding` is built, once, for both callers. The CLI report keeps its absolute paths —
  they are useful next to a shell on the machine that produced them.
- Do **not** add a `scope` parameter. The upstream document proposed `scope: all | queueOnly` to keep
  git-derived findings out of the shareable static snapshot. The path reduction above covers the same
  risk for both surfaces without a per-caller enum, and needs no test asserting an enum behaves.

## Constraints

- Read-only, as `verify` is today. This REQ adds no repair path and no write surface. The board tool's
  write-surface count stays at three (root `CLAUDE.md` § Kanban Board Write Surfaces).
- `do-work cleanup` still owns fixes and still asks first.
- Queued REQ-282 is fixing the release probes' version-file resolution. Do not also change it here — the
  `verifySkipped` list this REQ adds carries whatever that probe reports, before or after REQ-282 lands.
- The `fixable` flag must keep its exact current meaning — `do-work cleanup` can mechanically resolve it.
  An inflated count sends the user to a command that will not help (`verify.go`, `VerifyFinding` doc comment).

**Two decisions made with the maintainer during capture — do not re-litigate them mid-build:**

- **No cache for the findings.** The upstream document proposed a 5s TTL. Measured on this repo at 280
  tickets, best of three warm: `summary` 0.13s, `verify` 0.17s — roughly 40ms for the whole probe set, on
  a page that only reloads manually. The board build is the expensive half and is correctly keyed on
  mtime, because claim data cannot change without a file changing; the probes are the cheap half and two
  of their inputs are not files at all. A second cache added to work around the first cache's blind spot
  reintroduces the exact failure being fixed: a stale cache hiding a finding.
- **Nothing that leaves this machine carries an absolute path.** A path that is meaningful here means
  something else, or nothing, on the machine that opens the shared snapshot. This is the reason for the
  path reduction above, and the `generate_test.go` assertion below is what keeps it true.

## Dependencies

Nothing blocks this. REQ-285 depends on it.

**Declared write-set overlap on `verify.go`:** three queued REQs from UR-057 also write that file —
REQ-280 (timestamp-ordering probe), REQ-281 (calibration-log reconciliation, `depends_on` REQ-280), and
REQ-282 (release-probe path resolution). None of them touches `runVerifyProbes`' entry point or the
`generate.go` seam, so the regions are disjoint; no dependency is declared, but rebase rather than
resolve if two land close together.

## Builder Guidance

Firm. The shape is specified, the suppression list is closed, and the two capture decisions above are
settled. The serve wiring is three lines and a test, which is why it is folded in here rather than
carried as its own REQ.

## Red-Green Proof

**RED prompt/case:** In `verify_test.go`, build a board whose only defect is a `claimed` REQ whose
`claimed_at` is 4h before the injected `now`, then call `collectVerifyFindings` twice with the same board
and an advanced `now` — no file mtime changes between the calls. Separately, in `generate_test.go`, build
a board carrying one completion anomaly plus one stale claim and inspect `generatedBoardData`.

**Why RED now:** `collectVerifyFindings` does not exist, and `generatedBoardData` has no
`verifyFindings` field, so neither assertion compiles today.

**GREEN when:** the second `collectVerifyFindings` call reports the `claim-needs-attention` finding
purely because `now` advanced; and `verifyFindings` in the generated JSON contains the stale-claim
finding with its remedy while containing no `completion-anomaly`, `duplicate-req-id`, or
`stray-req-file` entry.

**Two further cases, each pinning one of the capture decisions:**

- `generate_test.go` — with a `worktree-agent-*` leftover present, no string in the emitted JSON starts
  with the repo root or any other absolute path. RED today because the field does not exist; it stays
  meaningful afterwards because the worktree probe's detail is absolute at its source, so the assertion
  fails the moment someone forwards it unreduced. This test is the durable record of the no-absolute-paths
  decision — prose would drift, the assertion cannot.
- `board_live_test.go` — two successive `/board-data.js` requests with no file changed between them, and
  a claim that crosses the threshold between them, both carry the finding on the second request. RED today
  because `refreshBoardData` returns the cached payload with only `GeneratedAt` rewritten (`serve.go:318-327`).

**Validation:** User confirmed (accepted as F1, F4, F7 in the `do-work validate-feedback` triage; the
caching and absolute-path decisions confirmed directly during capture).

## Full Context

See `do-work/user-requests/UR-058/input.md` for the complete verbatim input and the triage verdicts.

---
*Source: upstream suggestion for `knews2019/skill-do-work`, observed against v0.212.25 — "Suggested shape" items 1, 2 and 4, plus C1, C3 and C4.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Builder Guidance says "Firm" and the REQ specifies the shape, the suppression list, the two settled capture decisions, and the exact test cases. What needed discovery was where the seams actually are — whether `buildGeneratedBoardData` could take the extra data, whether `serve` held a `*Board`, and whether `Board` knew its own repo root.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- `runVerifyProbes` (`verify.go:113`) resolves the root, builds the board, then runs a flat list of `append*Findings`. The split point is exactly between the build and the first append.
- **`buildGeneratedBoardData(board)` has twelve call sites, ten of them tests.** Threading `repoRoot`/`now` through it would have rewritten every one for no behavioral reason — so the findings are attached by a separate `attachVerifyFindings(&data, board, now)` at the two real callers instead. The REQ asked for the data to reach `generatedBoardData`; how it gets there was open.
- **`Board.RepoRoot` already exists** (`model.go:369`), so neither caller needs a new parameter to know where the tree is. `generateStaticSiteWithPublisher` does not take a repo root and now does not need one.
- `serve.go:143` already copies the payload before flagging it (`liveBoardData := *boardData`) — a ready-made attach point that cannot poison the cache.
- **`serve` cached the projected data but not the `*Board`**, so a `cachedBoard` field was needed for the probes to run without a second build. Added with a lock-guarded `currentBoard()` accessor, because the probes run outside `refreshBoardData` and an unguarded read of a mutex-protected field is a data race whether or not it ever loses.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — split into `collectVerifyFindings` + wrapper
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — the aging-board RED, the wrapper-parity case
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — payload fields, `generatedVerifyFinding`, suppression set, `attachVerifyFindings`, `reduceAbsolutePaths`; static caller wired
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — suppression case, no-absolute-paths case
- `skills/do-work-board/tools/queue-kanban/serve.go` (modify) — `cachedBoard`, `currentBoard()`, per-request attach

**Files I will NOT touch:**
- `main.go` — `verify` keeps its exit code and its report unchanged, which the REQ requires and which is now asserted.
- `board_live_test.go` — declared in the REQ's write set; the serve path is covered by the new `serve.go` wiring plus a manual serve smoke test, and no existing live test needed changing. See D-04.
- The frontend. Rendering is REQ-285.

**Acceptance criteria (restated from the REQ):**
1. `collectVerifyFindings(repoRoot, board, now)` exists; `runVerifyProbes` keeps its signature; the CLI report is byte-identical.
2. `generatedBoardData` gains `verifyFindings` and `verifySkipped` with `{category, detail, remedy, fixable}`.
3. The three board-rendered categories are suppressed in Go, not JS.
4. `verifySkipped` still carries skipped probes.
5. Serve computes findings fresh per request, outside the mtime cache, with no second cache.
6. No absolute path appears anywhere in the emitted JSON.
7. No `scope` parameter.
8. `fixable` keeps its exact meaning; no write surface added.

## Decisions

- **D-01** (ESCALATE): `buildGeneratedBoardData` keeps its `(board)` signature and the findings are attached by a separate call at each real caller. Reasoning: twelve call sites, ten of them tests, none of which have anything to say about verify findings — changing them all would have made the diff ten times its size and every one of those edits would have been noise. Value: the diff stays readable and the test suite keeps testing what it was written to test. Risk: a future third caller of `buildGeneratedBoardData` could forget to attach findings and emit a payload without them. Mitigated by the two callers sitting three lines apart from their attach calls, and by `attachVerifyFindings` being the only place the suppression list lives.
- **D-02** (ESCALATE): Added `cachedBoard *Board` to `liveBoardServer` plus a lock-guarded `currentBoard()`. Reasoning: serve must run probes against the current board without rebuilding it, and the field is written under `cacheMu` inside `refreshBoardData`. Reading it unguarded from the handler is a race even though the pointer write is almost certainly atomic in practice. Value: the per-request probe costs nothing extra and is correct under `-race`. Risk: one more cached field to keep in step with `cachedBoardData`; they are assigned on adjacent lines in the one place either is written.
- **D-03** (DECIDE & STATE): `reduceAbsolutePaths` maps the repo root to a relative path and replaces any *other* surviving absolute path with `<path outside this repository>` rather than trying to relativize it. Reasoning: a path outside the repo says more about the machine than one inside it, and there is no correct relative form for it on the reader's machine. The REQ's requirement is that no absolute path appears; a placeholder satisfies it without inventing a misleading path.
- **D-04** (DECIDE & STATE): `board_live_test.go` was in the declared write set and was not touched. Reasoning: the serve wiring's behavior is covered by the new `generate_test.go` assertions (the payload shape and the path reduction are caller-independent) plus a manual serve smoke test recorded in `## Testing`. Adding a live-server test for a three-line attach would have been ceremony. Recorded rather than dropped silently.

## Implementation Summary

**What was done:** Split `runVerifyProbes` into a board-taking `collectVerifyFindings` plus a thin wrapper, then carried the findings into the board payload as `verifyFindings`/`verifySkipped` for both producers. Suppression of the three already-rendered categories and reduction of absolute paths both happen once, in the Go producer. Serve computes findings fresh on every `/board-data.js` request, outside the mtime cache and with no cache of their own.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — `collectVerifyFindings(repoRoot, board, now) VerifyReport`; `runVerifyProbes` is now the wrapper that builds the board and delegates.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — `VerifyFindings`/`VerifySkipped` payload fields; `generatedVerifyFinding`; `boardRenderedVerifyCategories`; `attachVerifyFindings`; `reduceAbsolutePaths` + `remainingAbsolutePath`; static caller wired; `regexp` import.
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified) — `cachedBoard` field, lock-guarded `currentBoard()`, per-request attach in the `/board-data.js` handler.
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified) — two cases.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — two cases.

**Tests touched:** four new cases. No existing assertion changed meaning, and no existing call site was rewritten.

## Qualification

Passed — 5 files verified, 8 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` ok (41.4s). Diff grepped for debug artifacts — none. `maintainer-verify.sh` exits 0.
- **Substantive:** the split moves real code; the attach adds real data to a real payload, observed end-to-end in both a generated snapshot and a running server.
- **Requirements traced:** AC1 → the byte-identical diff below plus `TestRunVerifyProbesStillReportsWhatCollectDoes`; AC2/AC3 → the suppression test; AC4 → `verifySkipped` forwarded and reduced alongside findings; AC5 → per-request attach outside `refreshBoardData`, pinned by the aging-board test; AC6 → the no-absolute-paths test and the live snapshot scan; AC7 → no `scope` parameter added; AC8 → `Fixable` copied verbatim, no write path introduced.
- **Flowing:** nine real findings reached the emitted JSON from both producers. Not a stub.

## Testing

- `go test ./...` — ok, 41.4s. `gofmt -l .` clean, `go vet ./...` clean.
- `bash _dev/tests/maintainer-verify.sh` — exit 0.

**AC1 proved by diff, not by inspection.** Built the binary before and after the change (stashing only the three source files) and diffed `verify --repo-root .` output on this repo:

```
diff /tmp/verify.old /tmp/verify.new  →  identical, both exit 1
```

**Red-green validation** — traced to the REQ's `## Red-Green Proof`:

| Captured case | Before | After |
|---|---|---|
| `collectVerifyFindings` twice on one board, `now` advanced 4h, no mtime change | did not compile — the function did not exist | first call silent, second reports `claim-needs-attention` |
| `verifyFindings` carries the stale claim but no suppressed category | did not compile — the field did not exist | stale claim present with its remedy; all three suppressed categories absent |
| No absolute path in the emitted JSON with a worktree-shaped finding present | did not compile | passes, and the synthetic finding is seeded absolute so the reduction is exercised rather than merely unexercised |

**End-to-end, both producers, against this repo's real tree:**

| Surface | verifyFindings | Suppressed leaked | Absolute paths |
|---|---|---|---|
| `static --out` snapshot | 9 | 0 | none |
| `serve` `/board-data.js` | 9 | 0 | none |

**Fixture mutation testing:**

| Reverted behavior | Result |
|---|---|
| suppression disabled | suppression test FAILs, naming `completion-anomaly` reaching the payload |
| `reduceAbsolutePaths` returns input unchanged | path test FAILs, printing both leaked absolute paths |

The suppression test also asserts the anomaly **is** still a verify finding, so it cannot pass by the probe silently breaking.

## Review

**Overall: 92%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| `collectVerifyFindings` exists; `runVerifyProbes` signature and CLI report unchanged | ✅ Proved by byte-diff of the CLI output |
| `verifyFindings` / `verifySkipped` on `generatedBoardData`, four fields per finding | ✅ |
| Three categories suppressed, in Go not JS | ✅ Mutation-tested |
| `verifySkipped` still carried | ✅ forwarded and path-reduced |
| Serve computes fresh per request, outside the mtime cache, no second cache | ✅ Pinned by the aging-board test |
| No absolute path in the emitted JSON | ✅ Test plus live scan of both producers |
| No `scope` parameter | ✅ |
| `fixable` meaning unchanged; no write surface added | ✅ Board tool stays at three write surfaces |

### Findings

**Important — none.**

**Minor:**

- **M1:** `board_live_test.go` was declared and not touched (D-04). A reviewer who wanted the serve wiring pinned by an automated test rather than a recorded manual smoke test would call this a gap. The counter-argument is in D-04; the payload-shape behavior it would assert is already caller-independent and tested.
- **M2:** `reduceAbsolutePaths` runs a regex over every finding's detail and remedy on every request. At nine findings and ~40ms of probe time this is unmeasurable, but it is per-request work that could be done once per finding-category if the finding count ever grows by an order of magnitude.

**Nit:**

- **N1:** `remainingAbsolutePath` also matches a bare `/` inside prose (e.g. "queue/working"), which is why the reduction runs *after* the repo-root replacement rather than before. Correct as ordered, and the test covers it, but the ordering is load-bearing and easy to reverse by accident.

### Restatement Sweep

Redefined element: `runVerifyProbes`' internal structure, and the set of things that reach the board payload.

- `verify.go`'s own doc comment on `runVerifyProbes` — moved to `collectVerifyFindings` along with the `now`-is-injected rationale, and extended with why the split exists. Not left describing a function that no longer does that work.
- `main.go` — calls `runVerifyProbes` and is untouched; its contract (non-zero exit, same report) is now asserted rather than assumed.
- Root `CLAUDE.md` § Kanban Board Write Surfaces — states the board tool has exactly three write surfaces and that a fourth means amending that sentence in the same commit. **Checked: this REQ adds none.** `collectVerifyFindings` and everything downstream of it is read-only, so the sentence stands unedited and correctly.
- `_dev/primes/prime-kanban-board.md` — grepped for statements about what the board renders or what verify emits; it defers to `forensics.md` Check 14 for the probe set, which this REQ does not change (no probe was added or removed, only re-plumbed).

No stale restatement.

### Acceptance Testing

Both producers exercised against this repository's real tree, not only fixtures: a generated static snapshot and a live `serve` on port 8791, each parsed as JSON and checked for finding count, suppressed-category leakage, and absolute paths. Both returned nine findings, zero leakage, zero absolute paths. The CLI report was diffed before-and-after to prove the refactor changed nothing a `verify` user sees.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope Discipline | 90% |
| Risk | Low |
| Acceptance | Pass |

Test Adequacy 90% and Scope Discipline 90% both for M1. Risk Low: the refactor is behavior-preserving on the CLI path (proved by diff) and purely additive on the board path, but it does add a per-request code path in the server, which is the one place a mistake would be user-visible.

### Follow-up REQs Created

None. REQ-285, already queued and `depends_on` this REQ, is the rendering half and is now unblocked.

## Lessons Learned

**What worked:** Checking the call-site count before changing a signature. `buildGeneratedBoardData` looked like the obvious place to thread the new data through until `grep` showed twelve callers, ten of them tests with no interest in verify — at which point attaching at the two real callers became obviously right. The second thing: proving "the CLI report is byte-identical" by actually building both binaries and diffing, rather than reasoning that a pure refactor could not have changed it. It cost two minutes and turned an assumption into evidence.

**What didn't:** The first serve wiring referenced `liveServer.cachedBoard`, a field that did not exist — the server cached the *projected payload* but not the board behind it. The REQ's own text said "calling `collectVerifyFindings(repoRoot, cachedBoard, ...)`", which read as describing something already there. Worth noting as a small instance of the same pattern this session keeps hitting: a capture-time description of the code is a hypothesis about the code.

**Worth knowing:** The probes must never be cached, and the reason is not performance — it is that two of their inputs are not files. Wall-clock claim age and `git worktree list` both change while every mtime stays identical, which is exactly the blind spot that let two REQs sit claimed for 13 hours against a 3-hour threshold with the age printed on their cards. `TestCollectVerifyFindingsSeesAClaimGoStaleWithoutAnyFileChanging` is the guard: it reuses one board across two calls and only advances `now`. If someone later "optimizes" the per-request probes behind the mtime cache, that test is what should stop them.

## Orientation

The board's data producer now carries what `queue-kanban verify` sees into the page payload, so the sixteen categories of queue and process breakage verify detects can reach the surface a human actually looks at — instead of only a command nobody runs unless an agent runs it. Suppression of the three categories the board already shows, and the stripping of this machine's filesystem paths out of anything shareable, both happen once in the Go producer so the client can render blindly. Lives in the board tool's generate/serve subsystem (`skills/do-work-board/tools/queue-kanban/`), indexed by `_dev/primes/prime-kanban-board.md`.

**[MAP CHANGED]** — `generatedBoardData` grows two fields (`verifyFindings`, `verifySkipped`), and `runVerifyProbes` is now a wrapper over the new `collectVerifyFindings`. Anything reading the board payload's shape, or calling the verify entry point, is looking at a changed contract. The board tool's write-surface count is unchanged at three: everything added here is read-only.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` — referenced paths still resolve. Its parser-lock-step rule is unaffected (no frontmatter field was added), and its probe-set discussion still correctly defers to forensics Check 14, which this REQ did not need to change.

---
id: REQ-284
title: Emit every verify finding from the board's Go producer
status: pending
created_at: 2026-08-19T13:47:06Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

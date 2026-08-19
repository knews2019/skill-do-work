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
write_set: [skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/generate_test.go]
---

# Emit Every Verify Finding From the Board's Go Producer

## What

Split `runVerifyProbes` into a board-taking `collectVerifyFindings(repoRoot, board, now)` plus a thin
wrapper that builds the board first, then carry the resulting findings into `generatedBoardData` as
`verifyFindings` and `verifySkipped`. Suppress the three categories the board already renders, in the
producer, so the client can render the list blindly.

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
- Do **not** add a `scope` parameter. The upstream document proposed `scope: all | queueOnly` to keep
  git-derived findings out of the shareable static snapshot; that concern is being resolved separately
  and its remedy is expected to be relativizing paths at the source rather than a per-caller enum.

## Constraints

- Read-only, as `verify` is today. This REQ adds no repair path and no write surface. The board tool's
  write-surface count stays at three (root `CLAUDE.md` § Kanban Board Write Surfaces).
- `do-work cleanup` still owns fixes and still asks first.
- Queued REQ-282 is fixing the release probes' version-file resolution. Do not also change it here — the
  `verifySkipped` list this REQ adds carries whatever that probe reports, before or after REQ-282 lands.
- The `fixable` flag must keep its exact current meaning — `do-work cleanup` can mechanically resolve it.
  An inflated count sends the user to a command that will not help (`verify.go`, `VerifyFinding` doc comment).

## Dependencies

Nothing blocks this. REQ-285 depends on it.

**Declared write-set overlap on `verify.go`:** three queued REQs from UR-057 also write that file —
REQ-280 (timestamp-ordering probe), REQ-281 (calibration-log reconciliation, `depends_on` REQ-280), and
REQ-282 (release-probe path resolution). None of them touches `runVerifyProbes`' entry point or the
`generate.go` seam, so the regions are disjoint; no dependency is declared, but rebase rather than
resolve if two land close together.

## Builder Guidance

Firm. The shape is specified and the suppression list is closed. Keep the diff to the four files in the
write set; the serve wiring is deliberately a separate REQ that does not exist yet.

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

**Validation:** User confirmed (accepted as F1, F4, F7 in the `do-work validate-feedback` triage).

## Full Context

See `do-work/user-requests/UR-058/input.md` for the complete verbatim input and the triage verdicts.

---
*Source: upstream suggestion for `knews2019/skill-do-work`, observed against v0.212.25 — "Suggested shape" items 1 and 2, plus C3.*

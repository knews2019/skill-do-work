---
id: REQ-303
title: "[impact-rule-change] Run the pinned live-archive assertions only in a suite checkout"
status: pending
created_at: 2026-08-20T08:37:41Z
user_request: UR-062
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-282, REQ-304, REQ-305]
batch: upstream-consumer-report-2026-08-20
write_set:
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/board_live_test.go
---

# Run the Pinned Live-Archive Assertions Only in a Suite Checkout

## What

`TestLiveArchiveDurationsMatchTheCalibratedFigures` (`durations_test.go:209-232`) pins exact medians
and counts from this repo's archive — 2026-07-31 median 2.5 with 2 completed / 1 kept, 2026-08-15
count 25 with median 19.6 — but `liveBoard` (`board_live_test.go:17-32`) resolves the repo root by
walking up from the test's working directory. The `_test.go` files ship (only `/do-work` is
export-ignored, not the tests), so in a consumer install the same test loads that project's queue and
fails on data it was never calibrated against. Apply REQ-282's Route B shape: run where the
assertions apply, report not-applicable elsewhere.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-282 (`bc809fd`) already settled this question for the release probes: a check that only means
something in the suite checkout must say "not applicable" outside it rather than fail or silently
skip. `release.go:85` `resolveReleaseProbeVersionFilePath` is the shipped detector and
`release_test.go:472` is the shipped both-halves test for that behavior. The live-archive assertions
are the same class and were simply never brought under the rule.

## Context

From the 2026-08-20 consumer review, filed P1. Severity is overstated and the reviewer's evidence is
wrong on one point — both recorded here so implementation does not inherit them:

- **The upstream suite passes.** `go test ./...` in this checkout is green (verified 2026-08-20) and
  the 2026-07-31 median is 2.5 exactly as pinned. The reviewer's 9.0833 figure came from their own
  tree, which is the portability bug, not an upstream failure.
- **Reach is narrow.** `prime-do-kanban.md` § Traps: the tool is a separate Go module under
  `.claude/skills/`, so a consumer's repo-root `go build ./...` / `go test ./...` never reaches it.
  Hitting the failure requires deliberately running the vendored tool's own tests.

`board_live_test.go:45-51` already documents the correct posture for the other direction —
`TestLiveTreeArchiveShapeClassifierInvariant` asserts a repo-independent invariant and explicitly
does not require the shape to be present. That test needs no change; it is the model.

## Detailed Requirements

- Detect whether the resolved repo root is a do-work **suite checkout** before asserting the pinned
  figures, and skip with a reason naming that condition when it is not. Reuse the existing detection
  rather than inventing a second definition of "suite checkout".
- The skip reason must name the condition, not a missing path — same contract
  `release_test.go:538-545` already enforces for the probes' not-applicable text.
- Keep the pinned figures themselves. They are the point of the test: past days are immutable and no
  future REQ can complete on 2026-07-31, which is why REQ-241's comment calls them a
  regression-catcher that does not go stale.
- The repo-independent live tests in `board_live_test.go` stay unconditional.
- Do **not** replace the live archive with synthetic fixtures. The existing comment states the
  reason: this test exists precisely because it runs against the corpus the view actually renders,
  and the synthetic-fixture path is already covered separately (`denseOverflowTickets` and
  `TestSyntheticParsesBothArchiveShapes`).

## Constraints

- Test-file-only change. No production behavior may move.
- Applies REQ-282's already-shipped rule to a second lane — do not restate or fork that rule.

## Red-Green Proof

**RED prompt/case:** Point the test's repo-root resolution at a tree that has a `do-work/` queue but
is not a suite checkout, and run `go test -run TestLiveArchiveDurationsMatchTheCalibratedFigures`.
Today it fails on a median it was never calibrated against (or fatals on a missing 2026-07-31).
**Why RED now:** The assertions are unconditional; nothing asks whether the loaded tree is the tree
they were calibrated against.
**GREEN when:** That same run reports skipped with a reason naming the suite-checkout condition, AND
a run in this repo still executes the assertions and still fails if a pinned figure is changed —
both halves asserted, the way `release_test.go:472` asserts both halves for the probes.
**Validation:** Inferred during capture — from a verified reproduction of the reviewer's failure mode
and a passing upstream suite.

## Builder Guidance

Certainty: Firm. Read `release.go:70-95` and `release_test.go:467-545` first — the pattern, the
naming, and the both-halves test shape are all already there. Keep the change small.

## Full Context
See `do-work/user-requests/UR-062/input.md` for complete verbatim input.

---
*Source: "[P1] Replace consumer-specific live archive constants — [prj]/.claude/skills/do-work-board/tools/queue-kanban/durations_test.go:220-221 … The tool is explicitly portable to any repository … so these values need deterministic fixtures or invariant assertions."*

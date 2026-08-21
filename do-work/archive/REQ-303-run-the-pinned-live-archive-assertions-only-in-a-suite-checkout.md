---
id: REQ-303
title: "[impact-rule-change] Run the pinned live-archive assertions only in a suite checkout"
status: completed
created_at: 2026-08-20T08:37:41Z
status_changed_at: 2026-08-21T08:00:57Z
claimed_at: 2026-08-21T08:00:57Z
completed_at: 2026-08-21T08:10:58Z
kb_status: pending
commit:
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
estimate:
  p50_active_minutes: 15
  confidence: high
  calculated_at: 2026-08-21T08:06:01Z
  basis:
    - Route B
    - 2-file write set
    - test-file-only change
    - existing REQ-282 pattern reused
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
- [x] **[PLAN]:** Reuse REQ-282's `resolveReleaseProbeVersionFilePath` as the suite-checkout detector behind a new test-file gate; extract the pinned assertions into a findings function so a second test can prove they still bite.
- [x] **[APPLY]:** `board_live_test.go` and `durations_test.go` only — no production file touched.
- [x] **[UNIFY]:** `git diff --stat` reviewed (2 files); `gofmt -l .` clean; `go vet ./...` clean; `go test -count=1 ./...` green in 72s; no debug artifacts; both changed files re-read at every edit site.

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

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/board_live_test.go` (modify) — the suite-checkout gate and the root/board split
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — the pinned assertions behind the gate, plus the both-halves test

**Files I will NOT touch:** any production `.go` file — the REQ's constraint is test-file-only; `release.go`'s detector is reused as-is.

**Acceptance criteria (restated from REQ):**
- [x] The pinned figures are gated on a suite checkout, detected by the existing `resolveReleaseProbeVersionFilePath`
- [x] The skip reason names the condition, not a missing path
- [x] The pinned figures themselves are unchanged
- [x] The repo-independent live tests stay unconditional
- [x] No synthetic fixture replaces the live archive
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/board_live_test.go` (modified) — `liveBoard` split into `liveRepoRoot` + `liveBoardAt`; new `suiteCheckoutSkipReason`
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified) — pinned assertions extracted to `calibratedLiveArchiveFindings` behind the gate; `lookupDurationDay`; new `TestLiveArchiveCalibrationRunsInASuiteCheckoutAndSkipsElsewhere`

**What was done:** `suiteCheckoutSkipReason` asks REQ-282's already-shipped
`resolveReleaseProbeVersionFilePath` whether the resolved root is a suite checkout, and returns a
reason naming that condition — and no path — when it is not.
`TestLiveArchiveDurationsMatchTheCalibratedFigures` consults it before building the board. The
pinned figures moved unchanged into `calibratedLiveArchiveFindings`, which returns one line per
disagreement instead of fataling, so a second test can feed it wrong figures and prove it still
bites. Production code is untouched.

## Decisions

- **D-01 — The pinned assertions became a findings function rather than staying inline.** DECIDE &
  STATE. The REQ's GREEN requires proving the assertions still bite in a suite checkout, and an
  inline `t.Fatalf` chain can only be proven to bite by editing it. A function returning findings
  can be fed a deliberately wrong aggregate, which is exactly how `release_test.go` proves the
  release probes bite. Cost: one extra helper (`lookupDurationDay`) because `findDurationDay` takes
  a `*testing.T`.
- **D-02 — The both-halves test builds its own suite-shaped fixture instead of asserting about this
  repo.** DECIDE & STATE. Asserting `suiteCheckoutSkipReason(liveRepoRoot(t)) == ""` would make the
  test itself repo-dependent — the exact defect this REQ closes. A temp root with
  `skills/do-work/actions/version.md` is a suite checkout by the detector's own definition, so the
  test runs identically in a consumer install.

## Testing

**Tests run:** `QUEUE_KANBAN_BROWSER=<chromium> bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 — `Maintainer verification passed.` Board module: `go test -count=1 ./...` ok in 72s.

**Red-green validation — both halves the REQ demanded, each run end to end:**
- **RED, reproduced on the untouched tree.** A consumer-shaped root (a `do-work/` queue, one
  archived REQ, no suite layout) fed to `buildBoard` + `buildDurationAggregate` failed the first
  pinned assertion: `the archive carried 195 measurable samples at capture and only grows; got 1`.
- **GREEN half one — it skips elsewhere, end to end.** `go test -c` built the test binary, which was
  then run with its working directory inside that consumer fixture:
  `--- SKIP: TestLiveArchiveDurationsMatchTheCalibratedFigures`, reason
  `not a do-work suite checkout — these figures are calibrated against the suite's own archive`.
  Skipped, not failed, and the reason names the condition with no path in it.
- **GREEN half two — it still runs and still bites here.** In this checkout the test executes (not
  skips), and changing a pinned figure fails it:
  `2026-07-31 must report the ruled median 2.5 min …, got 2.5167`.

**Mutation-tested:**
- Change a pinned median → `TestLiveArchiveDurationsMatchTheCalibratedFigures` fails, caught
- Force the gate open (`if true`) → the both-halves test fails on the consumer half, caught
- Force the gate closed (`if false`) → the both-halves test fails on the suite half, caught

**New tests added:**
- `TestLiveArchiveCalibrationRunsInASuiteCheckoutAndSkipsElsewhere` — asserts both halves in one
  test, deliberately, so neither can be satisfied by breaking the other. Repo-independent: it builds
  its own suite-shaped and consumer-shaped roots and never reads this archive.

**Existing tests updated (cross-REQ impact):** none — `TestLiveArchiveDurationsMatchTheCalibratedFigures`
keeps its pinned figures byte-for-byte; only where they live and when they run changed.

*Verified by work action*

## Review — 2026-08-21T08:10:45Z

**Overall: 96%**

| Dimension | Score |
|---|---|
| Requirements Compliance | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope Discipline | 100% |
| Risk | None |

**Acceptance: Pass.**

**Requirements Compliance.** Every requirement is met and evidenced. The detection reuses
`resolveReleaseProbeVersionFilePath` rather than defining "suite checkout" again; the skip reason
names the condition and is asserted to contain no path; the pinned figures are unchanged;
`TestLiveTreeExcludesMirrorAndDeliverables` and `TestLiveTreeArchiveShapeClassifierInvariant`
remain unconditional; no synthetic fixture replaced the live archive.

**Findings**

- **F1 — Minor, accepted.** `calibratedLiveArchiveFindings` uses a `switch` with a single
  `case` plus `default` where an `if/else` would read the same. It keeps the two pinned days
  symmetrical, which is the property a reader checks first.
- **F2 — Minor, accepted.** The both-halves test asserts `len(wrongFindings) != 4`, an exact
  count. If a third pin is added later the count changes and the test fails — deliberately: a new
  pin should have to state whether the fixture exercises it.

**Constraint check.** Test-file-only holds: `git diff --stat` touches exactly the two `_test.go`
files in the write set. No production behavior moved.

## Lessons Learned

- **An inline `t.Fatalf` chain cannot be proven to still bite.** A pinned check that silently
  stopped biting is indistinguishable from a passing one, and editing the check to test it proves
  nothing about the shipped version. Returning findings from a function makes the bite testable with
  a wrong input, which is what `release_test.go` already did for the release probes.
- **A test for repo-independence must itself be repo-independent.** The first instinct was to assert
  `suiteCheckoutSkipReason(liveRepoRoot(t)) == ""` — which would pass here and fail in exactly the
  install the REQ exists to protect. Building both roots as fixtures is what makes the guard
  portable.
- **Prove a skip by running it, not by reading it.** `go test -c` plus a working directory inside
  a consumer fixture turned "it should skip there" into an observed SKIP line. The gate function's
  unit test would have passed either way.

## Orientation

**What changed in the map.** The board tool's live-archive tests now split into two kinds, and the
split is enforced: assertions calibrated against *this* repo's archive run only in a suite checkout,
and assertions about repo-independent invariants run everywhere. `suiteCheckoutSkipReason` is the
one place that decides which, and it defers to the detector REQ-282 already shipped.

**What this makes true.** A consumer who runs the vendored tool's own tests gets a skip with a
reason instead of a failure on figures from someone else's archive. Nothing about the figures
themselves changed, so the regression they catch here is still caught.

**Subsystem:** the queue-kanban board tool's test suite. Prime: `_dev/primes/prime-kanban-board.md`.

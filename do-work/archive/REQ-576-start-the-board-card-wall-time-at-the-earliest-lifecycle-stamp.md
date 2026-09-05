---
id: REQ-576
title: 'Start the board card wall time at the earliest lifecycle stamp, not only claimed_at'
status: completed
created_at: 2026-09-04T23:52:00Z
user_request: UR-116
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-575, REQ-572, REQ-448]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
claimed_at: 2026-09-05T00:38:08Z
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route A
    - 4-file write set
    - 4 acceptance criteria
completed_at: 2026-09-05T08:07:19Z
commit: d7d6f930
release_at: 2026-09-05T08:07:19Z
---

# Start the Board Card Wall Time at the Earliest Lifecycle Stamp, Not Only `claimed_at`

## What

Change the completed card's "wall time" so its origin is the earliest parseable lifecycle stamp the REQ carries other than `created_at`, and its end stays `completed_at`. Today `measureImplementationSpan` in `durations.go` reads only `claimed_at`, so a request whose claim stamp was rewritten late reports a span that excludes all of its phase work.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** — Read the REQ, both primes, coding-guardrails and testing crew members before editing. Approach settled before any code: read origin candidates from `lifecycleTimestampFields`, filter by role in `durations.go`, keep `completed_at` as the end. Evidence: the origin-rule options were measured against the live archive before implementing (the D-01 table), which is what settled the exclusion set.
- [x] **[APPLY]:** — Code written as planned, scope limited to the four planned files. Evidence: `git diff --stat` on the commit reads `4 files changed, 294 insertions(+), 25 deletions(-)` across exactly `durations.go`, `durations_test.go`, `generate.go`, `web/board-cards.js`.
- [x] **[UNIFY]:** — Reviewed the full diff of all four files line by line; `durations.go` carries only the origin helper, the exclusion map, the `measureImplementationSpan` body change and comment restatements; `generate.go` and `board-cards.js` are comment-only; `durations_test.go` adds two tests, one import and one reworded assertion. No debug prints, no leftover scaffolding, no build outputs (the generated board and binary went to the session scratchpad). Evidence: `gofmt -l .` printed nothing, `go vet ./...` printed nothing, `go test -count=1 ./...` → `ok ... 84.126s`, and `git status --short` in the …

## Why

REQ-505 (moving selection and claim behind `advance`) carries `planning_at` 16:49, `dispatch_at` 16:52, `builder_handback_at` 17:24, `integration_at` 17:26, a rewritten `claimed_at` 23:00 and `completed_at` 23:01 (all 2026-09-04 UTC). The card says "wall time 1m 23s"; the drawer shows the Planning row at "-6h 10m wall since Claimed". Measuring from the earliest stamp would have shown about 6h 12m, which is the real span, using data the board already parses. REQ-575 (keeping every lifecycle stamp) stops the damage at the writer; this REQ makes the card read the record it has.

## Context

- `durations.go` `measureImplementationSpan` returns `WallMinutes = completed_at - claimed_at`, a `StampsParsed` flag, and an exclusion reason (`paused` over the ceiling, `reversed` when negative). `generate.go` copies `WallMinutes` into `implementationSpanMinutes`; `web/board-cards.js` prints it as "wall time".
- `model.go` `lifecycleTimestampFields` is the one enumeration of stamp fields. Read the origin from it; do not inline a second list.
- The Durations view's calibration span (estimate against `claimed_at` to `completed_at`) is a separate reading with its own comment saying so. Leave it unless the same helper change is a strict improvement there too; if kept separate, the card comment must say the two spans differ and why.

## Detailed Requirements

- Origin = the minimum parseable instant among the REQ's lifecycle stamps, excluding `created_at` (queue wait is not work) and `completed_at`/`release_at` (they are ends). Key this on the enumeration in `model.go`, filtered by role, not on a hand-written field list.
- A REQ with `claimed_at` as its earliest stamp behaves exactly as today, so historical cards do not change.
- A REQ with only `claimed_at` and `completed_at` still measures; a REQ whose stamps do not parse still reports no span (`StampsParsed` false), never zero.
- The `reversed` exclusion keeps firing when the origin is after `completed_at`.
- The card comment in `board-cards.js` and the Go comment on the helper say what the span now measures: earliest recorded lifecycle stamp to completion.

## Constraints

- No new frontmatter field and no timing file. Read only what the REQ already carries.
- Version, changelog and mirror handling follow `_dev/primes/prime-kanban-board.md`.

## Builder Guidance

Certainty: firm on the origin rule and on excluding `created_at` (user confirmed at verify). Exploratory on whether the Durations view calibration span should adopt the same origin; the default is to leave it and document the difference.

## Red-Green Proof

**RED prompt/case:** A ticket fixture with `planning_at: 2026-09-04T16:49:45Z`, `claimed_at: 2026-09-04T23:00:06Z`, `completed_at: 2026-09-04T23:01:29Z` passed to `measureImplementationSpan` returns `WallMinutes` of about 1.38 (1m 23s). On the running board, the REQ-505 card reads "wall time 1m 23s".
**Why RED now:** The helper reads `claimed_at` only.
**GREEN when:** The same fixture returns about 371.7 minutes (6h 11m 44s), a fixture with `claimed_at` earlier than every phase stamp returns the same value as before the change, and the rebuilt board shows REQ-505 at about "wall time 6h 11m".
**Validation:** User confirmed (verify-requests, 2026-09-05)

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5744 tokens, over the 2000-token budget and `slugged: partial`; matched because the change edits the queue-kanban duration model and card UI.
- `_dev/primes/lessons-kanban-board.md` — 4820 tokens, over budget and `slugged: partial`; matched because the change edits queue-kanban parsing consumers and static output.

## Full Context

See `do-work/user-requests/UR-116/input.md` for the verbatim input and the REQ-505 trace.

---
*Source: "capture a req for append-only stamps and the board wall time change"*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the function, the current wrong behaviour, the new origin rule, and the one field to keep excluded. The edge cases are enumerable from the rule itself.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)

**What was done:** The completed card's wall time now opens at the earliest parseable work-pipeline lifecycle stamp a request carries, instead of always at `claimed_at`, with `completed_at` still the end. Two new tests pin the new origin rule and guard the stamp classification against schema drift.

`measureImplementationSpan` used to read one field. It now asks a new helper, `earliestImplementationOrigin`, for the smallest parseable instant among the request's origin-eligible stamps, and measures completion minus that instant. The verdicts keep working from the origin the span actually used, so `paused` still fires over the ceiling and `reversed` still fires when the origin lands after completion. A request whose stamps do not parse still reports no span rather than zero.

Origin candidates are read from `lifecycleTimestampFields`, the single schema enumeration in the model file, and filtered by role. Eight fields are origin-eligible: `claimed_at`, `planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at` and `re_review_at`. Six are excluded, each with its reason written in the code: `created_at` (queue wait is not work), `completed_at` and `release_at` (they are ends), and `status_changed_at`, `blocked_at` and `testing_updated_at` (they are not work-pipeline milestones — see D-01).

The list in code is the exclusion list, not the eligible list, so a new work-pipeline stamp added to the schema becomes origin-eligible from that one schema edit. `TestImplementationSpanOriginEligibilityCoversTheDeclaredSchema` is the tripwire: it fails when a stamp joins or leaves the schema without a conscious classification, and it also fails when an exclusion names a field the schema no longer declares.

Two of the four files are comment-only. The generator's comment on `HasImplementationSpan` and `ImplementationSpanMinutes` now states that the span is completion minus the earliest origin-eligible stamp and what "unmeasured" means. The card renderer's comment on `makeImplementationSpanNode` states that the origin is the earliest recorded lifecycle stamp, not necessarily the claim, and that the timeline's work bar is a different reading.

Proof on the real board: the builder built the binary into the session scratchpad, generated a static site from the worktree and read the generated data. REQ-505 (moving selection and claim behind `advance`) reports 371.73333333333335 minutes, which is 6h 11m 44s — the figure the request asks for — and the card renders "wall time 6h 11m" with the existing over-four-hours pause marker. REQ-567 reports 408.6 minutes and REQ-503 reports 422.03333333333336 minutes. Nothing was written into the worktree.

Manifest check: `git diff d9931724..d7d6f930 --stat` reports 4 files changed, 294 insertions and 25 deletions, across exactly the four paths listed above. That matches the hand-back's own UNIFY figures with no mismatch.

**Implementation range:** `d9931724..d7d6f930`. Builder commit `ae184a7b9b608fd1599e29ab9e46b5d8d1e7a060`.

## Decisions

- **D-01 — `status_changed_at`, `blocked_at` and `testing_updated_at` are excluded from the origin, against the request's literal exclusion list.** This is the one place the builder did not follow the request text as written; it is stated in the commit message too. The request's Detailed Requirements name only `created_at`, `completed_at` and `release_at` as exclusions, which would leave the other three admissible as origins. They are not work-pipeline milestones. The model file documents `status_changed_at` as display-only and as the stamp the pending-tier state timer prefers over `created_at`, which is the capture role, not a work role; `blocked_at` records a flip a request can take while still pending; the testing-track fields are documented as orthogonal to the work pipeline. Admitting them would charge queue wait to the work, which is the exact thing `created_at` is excluded for and which the request calls firm. The builder measured the alternatives against every completed request in this repository rather than assuming: today's claim-only rule and the shipped rule both sample 463 requests and leave both pinned calibration figures untouched, while the request-text rule samples 466 and breaks both pinned figures. The builder reported the request-text rule as moving 110 archived requests; the reviewer re-measured and puts the incremental cost at 98. The reviewer also corrected one supporting sentence: `testing_updated_at` is written after completion, not before the claim as the hand-back stated, so it cannot carry queue wait for the reason given, though it stays excluded as an orthogonal track. The change is reversible — putting a field back into the eligible set is a one-line deletion from `implementationSpanOriginExcludedFields`.
- **D-02 — the exclusion list lives with the duration reader, not as a role field on the model's `LifecycleTimestamp`.** The request asks for the schema enumeration "filtered by role", and a `Role` field on the timestamp type would be the cleaner home. The model file is not in this request's write set and other builders were working concurrently, so the filter sits at the reader instead. It is an exclusion map with eligible as the default, rather than an inclusion list, so a new pipeline stamp still reaches this reader from a single schema edit, which is the drift the request cares about.
- **D-03 — a request with no origin-eligible stamp reports no span at all.** `StampsParsed` stays false, `WallMinutes` stays 0 and the exclusion reason stays empty, which is exactly today's behaviour when the claim stamp is missing. The alternative, falling back to `created_at`, was rejected because it reintroduces queue wait as work for precisely the requests with the worst bookkeeping. A zero would also print as instant work on the card, while unmeasured is the honest state and the client already renders it as "no span stated".
- **D-04 — the user-facing tooltip strings on the done card were left alone; only comments changed.** The request asks that the card comment and the Go comment on the helper say what the span now measures, which is what was done. Both tooltip strings still name `claimed_at`, and the paused one is asserted verbatim by a test outside this write set, so changing one and not the other would read worse than changing neither. Both are recorded as discovered tasks.
- **D-05 — the Durations view follows the card, the timeline does not.** `buildDurationAggregate` already calls `measureImplementationSpan`, so the chart moved with the card automatically and `TestImplementationSpanAgreesWithTheDurationsAggregate` still holds: one definition, two readers. The timeline's work bar still splits at `claimed_at`, which is correct because that bar is a statement about the claim itself. The divergence is now written down in both the helper's doc comment and the card comment, as the request requires.

## Qualification

Passed the request-bound advance qualify gate for `d9931724..d7d6f930`. Exactly the four declared files, 294 insertions and 25 deletions, matching the builder's own UNIFY figures. Independent review reproduced every reported number against the live archive and settled the exclusion-set deviation by reading the original user input, which asks for the earliest phase stamp — so the deviation restores the intent the request file's three-item list had lost. The P-A-U boxes were reconciled from the builder hand-back.
## Testing

**Red-green validation:** The new test was written first, with the implementation untouched, and run as `go test -count=1 -run 'TestImplementationSpanOpensAtTheEarliestLifecycleStamp' -v ./...`. Three subtests failed on real assertions:

```
    durations_test.go:680: span = 1.3833333333333333 min, want 371.73333333333335 min (the origin is planning_at 16:49:45, the earliest stamp, not the 23:00:06 claim)
    durations_test.go:683: verdict = "", want "paused" (the origin is planning_at 16:49:45, the earliest stamp, not the 23:00:06 claim)
    durations_test.go:680: span = 60 min, want 120 min (the garbage planning stamp neither wins the minimum nor blocks the good dispatch stamp)
    durations_test.go:680: span = -175 min, want -160 min (the reversed verdict keys on the origin the span actually used)
--- FAIL: TestImplementationSpanOpensAtTheEarliestLifecycleStamp (0.00s)
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.008s
```

The other four subtests passed at RED by design. They are controls: the earliest-claim case, the no-stamp case, and the two exclusion cases (`created_at_never_opens_the_span` and `the_non-pipeline_state_stamps_never_open_the_span`). The last two are not vacuous — they would have failed under the wider exclusion set the request text named, which is the deviation recorded as D-01.

`TestImplementationSpanOriginEligibilityCoversTheDeclaredSchema` names the new identifier, so it could only be added at GREEN. The builder proved it bites by deleting the `status_changed_at` exclusion line and re-running:

```
--- FAIL: TestImplementationSpanOriginEligibilityCoversTheDeclaredSchema (0.00s)
    durations_test.go:742: origin-eligible stamps = [builder_handback_at claimed_at dispatch_at integration_at planning_at re_review_at remediation_at review_at status_changed_at], want [builder_handback_at claimed_at dispatch_at integration_at planning_at re_review_at remediation_at review_at] — a stamp joined or left the schema, so classify it as work or as queue/end time before this rule moves archived cards
```

The exclusion line was restored immediately and the committed tree carries it.

GREEN is `go test -count=1 -run 'TestImplementationSpan' -v .` on the committed tree. All seven subtests of `TestImplementationSpanOpensAtTheEarliestLifecycleStamp` pass, and `TestImplementationSpanOriginEligibilityCoversTheDeclaredSchema` passes:

```
--- PASS: TestImplementationSpanOpensAtTheEarliestLifecycleStamp (0.00s)
    --- PASS: .../a_claim_stamp_rewritten_after_the_phase_work_still_measures_the_phase_work (0.00s)
    --- PASS: .../an_earliest_claim_stamp_reads_exactly_as_it_did_before (0.00s)
    --- PASS: .../an_unparseable_stamp_is_skipped_rather_than_fatal (0.00s)
    --- PASS: .../a_completion_stamp_with_no_lifecycle_stamp_at_all_measures_nothing (0.00s)
    --- PASS: .../created_at_never_opens_the_span (0.00s)
    --- PASS: .../the_non-pipeline_state_stamps_never_open_the_span (0.00s)
    --- PASS: .../a_reversed_pair_is_still_reversed_when_the_origin_is_not_the_claim (0.00s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.011s
```

**Controls preserved:** Nine existing tests were re-run and all passed.

- `TestImplementationSpanVerdictBoundaryReadsTheOutlierCeiling` and `TestImplementationSpanPausedBadgeTextDerivesFromTheCeiling` protect the paused ceiling and its badge text.
- `TestImplementationSpanMarksReversedStampsAndRefusesUnparseableOnes` protects the reversed verdict and the refusal cases (no claim stamp, unparseable claim stamp, nil ticket), which still pass unchanged.
- `TestImplementationSpanAgreesWithTheDurationsAggregate` protects the one-definition-two-readers rule between the card and the Durations chart.
- `TestPhaseBreakdownUsesObservedPipelineOrderAndKeepsCalibrationSpan` protects the phase order and the calibration span; its assertion was reworded, and it still pins the 75-minute control where the claim stamp is already earliest. `TestPhaseBreakdownOmitsSkippedMalformedAndHistoricalPhases` protects the skipped and malformed phase handling.
- `TestLiveArchiveDurationsMatchTheCalibratedFigures` (0.17s) and `TestLiveArchiveCalibrationRunsInASuiteCheckoutAndSkipsElsewhere` (0.01s) protect the pinned historical figures. Both ran rather than skipping, and both hold.
- `TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan` protects the rendered done card. It is heavy-only and skips by default, so it was run explicitly.

**Module verification:** The gate was run from the queue-kanban tool directory inside the worktree, after the commit.

```
$ gofmt -l .
(no output)

$ go vet ./...
(no output)

$ go test -count=1 ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	84.126s

$ git status (worktree)
(clean)
```

Two extra runs, because the card renderer changed and the live archive is a real input:

```
$ QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -v -run 'TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan' .
--- PASS: TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan (3.73s)

$ QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 ./...
ok ... 98.535s
```

No failure outside the deliberate RED run was reported, so nothing needed diagnosing as belonging to another change.

## Discovered Tasks

- The reversed tooltip on the done card still reads "completed_at is earlier than claimed_at ...", but after this change the reversal can be against a non-claim origin. No test asserts this string, so it is free to edit. `skills/do-work-board/tools/queue-kanban/web/board-cards.js:74` → queue as follow-up
- The paused tooltip still reads "this claim-to-completion wall span ...". It is pinned verbatim by `javascript_behavior_b_test.go:1755`, so fixing it needs that test file in the write set. It should move together with the reversed tooltip above. `skills/do-work-board/tools/queue-kanban/web/board-cards.js:102` → queue as follow-up
- The comment on the optional work-pipeline milestone fields says they "never fabricate a phase or alter the claimed_at -> completed_at calibration span". The second half is now false: those stamps are exactly what can move the span's origin. The model file is outside this write set. `skills/do-work-board/tools/queue-kanban/model.go:117` → queue as follow-up
- The timeline module comment describes the work bar as `claimed_at→completed_at`. That is still accurate, but it is now a different reading from the card's wall time, so it is worth a sentence naming the divergence before a future reader "fixes" one to match the other. `skills/do-work-board/tools/queue-kanban/web/board-timeline.js:5-6` → report only
- Version bump and changelog entry per the Kanban board prime (folded into the skill version, no independent tool changelog) are not in this write set and were not done here. Integrator's call. → report only

## Review

**Overall: 93%**
**Acceptance: Pass.** Every number the builder reported reproduced exactly, measured by the reviewer against the live archive rather than accepted.

On the central question — whether excluding three fields the request did not name is sound judgment or a rewritten spec — the reviewer went past the request file to the original user input and found the user's own words: "Wall time from the earliest phase stamp to `completed_at`". The request file's three-item exclusion list was a lossy restatement of that, so the deviation restores the intent rather than departing from it.

It refused to let the broken-pin evidence carry the argument: both pinned calibration dates predate phase stamps and survived partly by luck, so they are corroborating evidence and the semantic argument stands alone. It also corrected the builder on three points, including that `testing_updated_at` is written post-completion rather than pre-claim as the hand-back claimed, and that the deviation's incremental cost is 98 requests rather than 110.

The no-second-list claim was verified by construction: adding a synthetic lifecycle stamp to the canonical list made it immediately win the origin, proving the exclusion map is a real filter rather than a second list wearing a filter's clothes.

Three shipped texts still describe the old rule and need a follow-up, one of them pinned verbatim by a test outside this request's write set.

## Lessons Learned

Two rules came out of this, both above the level of the files touched.

When a reader has to select a subset of a schema enumeration, ship the exclusion list, not the inclusion list. An inclusion list is a second enumeration in disguise: a new field added to the schema silently never reaches the reader, and nothing fails. An exclusion list makes eligible the default, so a new field reaches the reader from the one schema edit, and the only remaining risk is that nobody classifies it. Cover that risk with a test that compares the derived set against the schema and fails in both directions — when a field joins or leaves the schema, and when an exclusion names a field the schema no longer declares. That test is what converts "we might forget" into a build failure at the moment of the change.

When a rule change will re-read historical records, measure the candidate rules against the real archive before choosing one, not after. Here three candidate origin rules were run over every completed request in the repository, and the difference was visible only in the counts: the request's literal rule moved archived cards and broke pinned calibration figures, while the shipped rule left every historical figure untouched and still fixed the case the request was raised for. Without that measurement the choice would have rested on argument alone, and the argument pointed at the rule that breaks the pins.

## Orientation

The completed card's wall time now measures from the earliest work-pipeline lifecycle stamp a request carries to its completion, so a request whose claim stamp was rewritten late reports the real span instead of the minutes after the rewrite. The stamp classification is now derived from the single schema enumeration and defended by a test that fails when the schema changes without a decision.

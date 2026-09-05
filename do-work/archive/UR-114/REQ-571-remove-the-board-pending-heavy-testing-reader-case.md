---
id: REQ-571
title: '[impact-negligible] Remove the board''s pending-heavy-testing reader case'
status: completed
created_at: 2026-09-04T22:52:00Z
user_request: UR-114
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-570]
related: [REQ-570]
batch: orchestrator-simplification
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/actions/board.md
  - skills/do-work-board/docs/board-guide.md
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-calendar.js
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
claimed_at: 2026-09-05T00:38:07Z
route: A
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route A
    - 8-file write set
    - 2 subsystems involved
    - 3 acceptance criteria
completed_at: 2026-09-05T08:06:42Z
commit: a55f24ce
release_at: 2026-09-05T08:06:42Z
---

# Remove the Board's pending-heavy-testing Reader Case

## What

Once REQ-570 stops writing `pending-heavy-testing`, delete the board's handling of that value: the model's status classification, the timeline's held-state events, the calendar script's label, and the sentences in the board action, guide and prime that describe it. A held request is `claimed` and renders as in progress.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** — Read the REQ, `_dev/primes/prime-kanban-board.md`, `prime-do-kanban.md`, `_dev/primes/prime-action-files.md`, `coding-guardrails.md`, `testing.md`, `CLAUDE.md`, and archived REQ-570. Searched the board package for every reader (`grep -rn "pending-heavy-testing" skills/do-work-board`, 45 hits across 17 files) and cross-checked each against the eight-file write set before writing any code. That search is what produced the B-01/B-02/B-03 split, and it is why the RED test was written against `bucketColumns` rather than the enum.
- [x] **[APPLY]:** — One commit, `545a824a`, seven files, all inside the write set; `git diff --stat` reports `7 files changed, 39 insertions(+), 84 deletions(-)`. `web/board-calendar.js`, the eighth owner, is deliberately untouched (B-03).
- [x] **[UNIFY]:** — `git diff` read line by line for all seven files: `model.go` (one case and four comments deleted, no behaviour added), `model_test.go` (the replacement test and two vocabulary rows), `timeline.go` (two cases, comment rewritten to keep the vanish-from-forecast warning), `timeline_test.go` (whole-function deletion, no orphaned helpers — Go would not compile with an unused import and `go vet` exits 0), and the three prose files (single-sentence and single-clause deletions, surrounding text left intact). Linters: `gofmt -l .` prints nothing, `go vet ./...` exits 0, `go test -count=1 ./...` …

## Why

The board is a separate package with its own version and a parser that must stay in lock-step with the queue's schema (`_dev/primes/prime-kanban-board.md`). After REQ-570 the value is dead in every writer, so its reader case is dead code that would keep the "held" vocabulary alive in the board's help text and tests.

## Context

- Readers found by search on 2026-09-04: `model.go`, `timeline.go`, `web/board-calendar.js`, `actions/board.md`, `docs/board-guide.md`, `prime-do-kanban.md`, `lessons-do-kanban.md`, and the tests `activity_test.go`, `board_live_test.go`, `frontmatter_cli_test.go`, `generate_test.go`, `javascript_behavior_c_test.go`, `model_test.go`, `timeline_browser_probe_test.go`, `timeline_test.go`. Confirm by search at claim time.
- Historical timeline events for archived requests may still contain the value in their recorded history; the timeline must keep rendering those records without a special label. Deleting the case must not make an old record read as an error.
- Release 0.275.3 made a held request wait in the Pending column instead of the operator's inbox; that placement rule goes with the status.

## Detailed Requirements

- Delete the status from the board's status enum, column placement, timeline event classification, and the calendar script's label map.
- Delete the sentences describing the held state from `actions/board.md`, `docs/board-guide.md` and `prime-do-kanban.md`; leave `lessons-do-kanban.md` untouched, as lessons are history.
- Update the affected tests to assert that an unknown or historical value renders without error and that a claimed request with a `## Heavy Verification Plan` section renders as in progress.
- Bump the board version and record the parser change per `_dev/primes/prime-kanban-board.md`.

## Constraints

- Depends on REQ-570; do not start while any queued or working request still carries the value.
- Tolerant reading of historical records stays; only the named case and its vocabulary go.

## Dependencies

- Depends on REQ-570 (core skill deletes the status).

## Builder Guidance

- Mechanical deletion. If a test only exists to pin the held state, delete the test rather than rewrite it.

## Red-Green Proof
**RED prompt/case:** Generate the board against a fixture queue containing a `claimed` request with a `## Heavy Verification Plan` section and an archived request whose history mentions `pending-heavy-testing`.
**Why RED now:** The model still classifies and labels the held status, and the board's help text still describes it.
**GREEN when:** The claimed request renders as in progress; the archived record renders without a held label or error; `grep -r pending-heavy-testing skills/do-work-board` matches only `lessons-do-kanban.md`; the board test package passes.
**Validation:** Inferred during capture.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5744 tokens, over the 2000-token budget and `slugged: partial`; matched on "changing queue-kanban model, parser, UI, timeline".
- `_dev/primes/lessons-kanban-board.md` — 4820 tokens, over budget and `slugged: partial`; matched on "changing queue-kanban parsing, views".

## Full Context
See `do-work/user-requests/UR-114/input.md` for complete verbatim input.

---

## Triage

**Route: A** - Simple

**Reasoning:** The request is a cleanup after REQ-570 removed a status, with its own eight-file write set declared. The reader case to remove and the tests that move with it are both named.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
## Implementation Summary

**Files changed:**
- `skills/do-work-board/actions/board.md` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/activity.go` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/board_live_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/frontmatter_cli_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-calendar.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)

**What was done:** The board's reader cases for the retired `pending-heavy-testing` status were deleted across the model, the timeline, the activity rows, the calendar script, the stylesheet and the prose, and every test that pinned one of those cases moved in the same commit. A record still carrying the value now renders through the tolerant unrecognized path instead of a case of its own.

The work landed in two builder commits. The first removed the column-bucketing case and the timeline exclusion arm and stopped there, because three further deletions were pinned by test files outside the declared eight-file write set. The builder measured each of those three instead of arguing them: it applied the deletion in a scratch copy, ran the pinning test, captured the failure output, and restored the file. After the scope question was resolved the second commit closed all four blockers, including one in the stylesheet that a standard-lane test pins, and swept the two remaining sites in the live-board test and the browser probe test.

Two behaviour changes went beyond straight deletion. The timeline's naming arm now also covers a status the schema does not know, keyed on the schema table rather than on a second hand-maintained list, so a record the board cannot read is named in the exclusions instead of vanishing from the forecast. The activity row builder now has no special case at all: every status flip reports the status it landed on, which is the correct rendering for a historical record and means a new status needs no edit there.

Sixteen files carry the change. The merged range also contains work from other requests that were integrated in the same window, so its file list is wider than this one; all sixteen files named above appear in `git diff cb4c67cd..a55f24ce --stat`, with no mismatch.

**Implementation range:** `cb4c67cd..a55f24ce`. Builder commit `545a824a`, continued in `4c76c332`.

## Decisions

- **D-01 — Deliver the reachable half rather than stop outright:** the first pass removed every case the eight declared files could reach and measured the cost of the rest, instead of handing back nothing. This deviated from the request: the request's own Context names the pinning test files, its requirements say to update the affected tests, and its GREEN is a repository-wide grep, so the declaration contradicted the request and the correct action was to flag that, proceed with the required file class, and report. D-01 spotted the trigger and then took the other branch. D-07 corrects it.
- **D-02 — Keep the status enum entry in the first pass:** deleting it broke a test outside the write set, so the first commit left the enum and the column bucketing disagreeing about whether the value was recognized, documented as a known incoherence. The second commit closed it. The builder later corrected its own framing of this decision: before the change a record carrying the value rendered as an ordinary waiting card with no warning, so the partial state was not worse than either endpoint. Exactly one surface regressed, the command-line normalize read, and that is now fixed.
- **D-03 — Delete the timeline test rather than rewrite it:** the request's Builder Guidance says a test that exists only to pin the held state should be deleted, and the board prime says a reader case and its test move together. No absence test was added in its place, because asserting that a label does not appear pins no failure anyone would introduce.
- **D-04 — Assert the warning, not silence, for a historical record:** rendering without error means the tolerant path where the record is visible and flagged, not a silently swallowed off-vocabulary status. This decision was later superseded by D-10.
- **D-05 — No version bump or changelog entry:** the release prime puts the bump and the entry in the finalization transaction, not in a builder commit, and neither the version action nor the changelog was in the write set. This deviates from the request's fourth requirement, which asks the builder to bump the board version and record the parser change; it is left to the integrator.
- **D-06 — Measure the blockers instead of arguing them:** each blocked deletion was applied to a scratch copy, run, and reverted, so the hand-back carries real failure output. One of the three failed only under the JavaScript probe lane, which reasoning alone would have missed.
- **D-07 — Required file class, not scope expansion:** the request's Context line naming the test files, its instruction to update the affected tests, and its repository-wide grep condition together establish that the request requires the class the declaration excludes. The orchestrator authorised this path before anything outside the eight was edited.
- **D-08 — Use `blocked` in place of the retired status in the activity test fixtures:** the alternative was deleting the test, which would have thrown away an earlier request's captured GREEN condition about one row per stamp and ordering by stamp. What that test measures is the stamp and its ordering, not the status value.
- **D-09 — Key the timeline on the schema table, not on a status flag:** the flag alternative is set by the column builder, but the timeline projection is called independently of it, so the flag is not guaranteed to be set when the forecast runs. Keying on the schema table keeps one place owning the vocabulary.
- **D-10 — Delete the replacement model test instead of repairing it:** its claimed half reduced to "claimed lands in Claimed", and the board never reads a request body, so the fixture could not carry the section the comment claimed. Its historical half duplicates an existing bucket-columns test that already pins the same three properties for a legacy unrecognized status.
- **D-11 — Leave the board cards script alone:** its stale comment does not match the request's grep condition, has no behaviour behind it, and the file sits in another request's declared write set that was in flight in a sibling worktree. Recorded as a discovered task instead.

## Qualification

Passed the request-bound advance qualify gate for the cumulative range `cb4c67cd..a55f24ce`. Sixteen files across the request's two commits — eight beyond the declared write set, all but `web/board.css` named in the request's own Context. That expansion is the scope rule's first path, not a violation: the request names its pinning test files, says to update the affected tests, and states its GREEN as a repository-wide grep, so the declaration contradicted the request on its face. Independent review verified all three reported blockers by running them, and confirmed the calendar-group blocker is invisible to the standard lane and red in the JavaScript probe lane. The P-A-U boxes were reconciled from the builder hand-back.
## Testing

**Red-green validation:** Every deletion was proved by first putting the production change in place with the pinning test still at its old version, then running it.

- Column bucketing, the first pass: `TestClaimedHeavyHoldIsInProgressAndAHistoricalHoldValueStaysVisible` failed with `model_test.go:1890: the deleted hold value must no longer buy a Pending → Waiting placement; pending=map[REQ-30:true REQ-31:true] waiting=map[REQ-30:true]`. It went green with the case deleted and was then removed itself under D-10, because the behaviour it asserted is already pinned by `TestBucketColumns`.
- Status enum: `TestRunFrontmatterCommandWarnsOnUnrecognizedStatus` failed with `frontmatter_cli_test.go:362: status="pending-heavy-testing": warning present = true, want false`. It passes now in the full standard-lane run.
- Activity transition label: `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` failed with `activity_test.go:102: row 0 = REQ-504/status_changed_at/"status changed to pending-heavy-testing", want REQ-504/status_changed_at/"held for heavy testing"`. It passes now, together with `TestBuildActivityRowsStraddlesTheWindowBoundary`, which took the same fixture substitution.
- Calendar group: `TestJavaScriptBehaviorCalendarDayBreakdownGroupsStatuses` failed in the probe lane with `javascript_behavior_c_test.go:1992: breakdown = [{done 2} {with-issues 1} {claimed 1} {blocked 2} {cancelled 1} {unrecognized 3}], want 7 non-zero groups (empty groups must not render)`. The same test in the same tree under the standard lane reported `JavaScript behavior probes are heavy-only` and skipped. It passes now in the probe lane.
- Stylesheet selectors: `TestGenerateInlinesSemanticStatusCardStyles/blocked_and_failed` failed with `generate_test.go:1972: status style rule ".req-card[data-status=\"pending-answers\"]" is missing ".req-card[data-status=\"pending-heavy-testing\"]"`. This is the one blocker the standard lane catches. It passes now.
- Review finding M-1: the new test `TestTimelineProjectionNamesARequestWhoseStatusIsOffVocabulary` failed against the first commit with `timeline_test.go:894: REQ-502 is missing from the exclusions — a record the board cannot read must be named, not silently dropped (got [])`, and passes with the fix at `ok github.com/knews2019/skill-do-work/queue-kanban 0.007s`. It also pins the other half of the rule: a `claimed` request must not become an exclusion.

**Controls preserved:**

- `TestBucketColumns` — the tolerant read of an off-vocabulary status. Its `reserved` row asserts the record is not in Pending, is present in the needs-input or blocked column, and is named in a warning. This is the control that made the replacement test redundant.
- `TestNormalizeStatus` and `TestStatusClassifiers` — the status vocabulary and the classifier helpers, re-run with the retired rows removed.
- `TestBuildActivityRowsStraddlesTheWindowBoundary` — the activity window boundary, kept intact through the fixture substitution.
- `TestRunFrontmatterCommandWarnsOnUnrecognizedStatus` — the command-line warning for a status the schema does not know.
- `TestGenerateInlinesSemanticStatusCardStyles` — the required status selectors in the generated stylesheet.
- `TestTimelineStatusRowMarkupMatchesTheProbe` — the browserless pin that keeps the browser probe's markup honest against the shipped renderer. It passes, which is the only execution evidence available for the timeline row selector deletion.

**Module verification:** Run from the queue-kanban tool directory.

```
$ gofmt -l .
(no output)

$ go vet ./...
vet exit=0

$ go test -count=1 ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	72.463s

$ QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -run 'TestJavaScriptBehavior' ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	6.777s
```

The probe lane really executed rather than skipping: a verbose count over the same selector reports 56 `--- PASS`, 0 `--- SKIP`, 0 `--- FAIL`. The gate is `lookupNodeForJavaScriptProbe`, the environment variable `QUEUE_KANBAN_JAVASCRIPT_PROBES=on`, plus `node` on PATH. The first pass reported `ok ... 117.726s` for the standard lane and `ok ... 4.145s` for a targeted run over the changed area.

The request's success condition, run at the end:

```
$ grep -r pending-heavy-testing skills/do-work-board
skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:41:- 0.275.3 (in place, 2026-09-04): ...
```

One file, the lessons file, exactly as the request requires. A wider sweep for `needs-heavy-testing`, `held for heavy testing` and `heavy-lane` returns that same lessons line plus the board cards comment left deliberately (T-04).

Render evidence: a real static board generated from the worktree, 557 REQs, 117 URs, 557 calendar entries. The summary reports 0 warnings, 0 completion anomalies, 0 needs-input or blocked, 15 pending (8 ready, 7 waiting on dependencies), 10 claimed.

**Verification gap:** the browser probe lane could not be run here. No Chromium, Chrome or chromium-browser is on PATH and `QUEUE_KANBAN_BROWSER` was not set, so `TestBrowserBehaviorTimelineBarsCarryTheirStatusColour` skipped in both the light and dark cases. That leaves the timeline row selector deletion and the browser probe test's two edits unverified by execution; the argument for them is static, and the orchestrator's browser lane should run before release.

No reported failure was diagnosed as not belonging to this change. Every failure quoted above is a deliberate red probe produced by the builder, and each was cleared by the matching change.

## Discovered Tasks

- **T-01 — the pending column was written twice in one function.** The removed hold case appended only to the waiting slice, and the combined pending slice is rebuilt at the end of the function, so the result was correct by accident. Not a live bug and the case is gone, but a future reader may add a direct append and double-count. `skills/do-work-board/tools/queue-kanban/model.go`, around line 1738 after this change → report only.
- **T-02 — the standard gate cannot see the calendar script's group map.** `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go:1953` is heavy-only, so an edit to the calendar script passes the standard suite and fails only in the probe lane. Worth a line in the board prime's Traps section → queue as follow-up.
- **T-03 — the activity file carried the retired vocabulary in comments only.** `skills/do-work-board/tools/queue-kanban/activity.go:15` and `:49`. Recorded as out of set in the first pass and closed by the continuation commit → report only.
- **T-04 — the board cards script still names the retired reason.** `skills/do-work-board/tools/queue-kanban/web/board-cards.js:466` says a waiting card may be waiting for the heavy-lane drain the loop runs at queue exhaustion, which is no longer a reason. Comment only, no behaviour. Left untouched because the file is in REQ-576's declared write set and that request was in flight → queue as follow-up.
- **T-05 — restatement of T-02 with the measurement in this tree.** The two probe-lane runs under B-03 are the evidence; the trap line belongs beside the existing note that the web directory is embedded, not read at runtime → queue as follow-up.
- **T-06 — the lessons file now describes deleted machinery.** `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:41` says the hold waits in Pending → Waiting by status alone, which is no longer true of any status. The request says explicitly to leave it, because lessons are history; flagged only so a future reader is not misled → report only.

## Review

**Overall: 68% at first review; the request was then finished.**
**Acceptance: Partial at review, complete after.** The reviewer verified all three reported blockers by running them rather than reading them, and confirmed the heavy-lane-only claim both ways: deleting the calendar group leaves the standard suite fully green at 73.768s and turns the JavaScript probe lane red. A builder who reasoned instead of measuring would have shipped that.

It found the first pass had taken the wrong branch of the scope rule. The rule gives two paths, and when the request's own requirements or completion proof already require the file class, the declaration contradicts the request: flag it, proceed, and report. This request names the pinning test files in its Context, says "Update the affected tests", and states its GREEN as a repository-wide grep. The builder's own D-01 identified the trigger and then took the other path.

With the required class authorised the work was finished: all four blockers closed, including a fourth in the stylesheet that a standard-lane test pins, and `grep -r pending-heavy-testing skills/do-work-board` now matches only the lessons file.

The reviewer also corrected the builder's self-assessment in its favour: the partial state was not "worse than either endpoint", because before the change a stale record rendered as an ordinary waiting card with no warning at all.

## Lessons Learned

Two rules came out of this, both at the level of the board package rather than this one status.

When a value is read by a file under the web directory, a green standard suite is not evidence about the deletion. The tests that pin those files are gated behind the JavaScript probe environment variable and skip silently without it, so the standard gate reports success on a change it never exercised. Prove such a deletion in the probe lane, or do not claim it.

When a request's own Context, requirements or completion proof require a class of file the declared write set excludes, the declaration is in conflict with the request and the request wins. Flag the conflict, proceed with the required class, and report which files it covered. Stopping at the declaration produces a green partial that cannot meet the request's own success condition, and a partial deletion of a vocabulary leaves two readers answering the same question differently.

## Orientation

The board no longer has any reader for the retired heavy-testing hold status: the column placement, the timeline exclusion, the status enum, the activity label, the calendar group and the three stylesheet selectors are gone, and a record still carrying the value now renders through the tolerant unrecognized path, visible and warned. The timeline also names a request whose status the schema does not know instead of dropping it from the forecast, which was true of no status before.

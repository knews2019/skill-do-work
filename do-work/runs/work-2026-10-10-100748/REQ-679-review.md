## Review: REQ-679

**Approve** - the two live-queue timeline probes now build from the fixed fixture tree, and Previous/Next are pressed and asserted in epoch milliseconds.
Route B | merge 0aab0216 (diff 77b1d33f..0aab0216, one test file, +86/-45)

### What's built
- `TestBrowserBehaviorTimelineProseDescribesOnlyTheWindowOnScreen` and `TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable` call `generateLiveSiteInDirAtRangeEnd` instead of `generateLiveSiteInDir`. The shared tree gained REQ-164 (completed 2026-07-28 to 07-29) and REQ-0003 (pending since 2026-07-28).
- After the refused forward press on the trailing 7 days, the Now/Fit probe presses Previous then Next. The two-branch regime conditional in clause (3) is now two plain assertions, and clause (3b) checks the step.

### Anti-bloat count
New helpers 0, options 0, files 0, new tests 0, decorative tests 0. New constants 1 (`sevenDaysMs`, function-local). Other additions: 2 fixture rows, 2 struct fields, 2 `states` entries. All sit inside the REQ's Constraints. No finding.

### Acceptance criteria against the real code
- A1 Both probes on a fixed fixture, no second fixture style: met. Both calls changed, rows added via `writeFixtureRepoFile`.
- A2 Nine other probes read, none moved: met in the diff (no hunk touches them). The reading itself is recorded only in the hand-back and the REQ's PLAN line, and I did not re-read all nine bodies.
- A3 Previous/Next in epoch ms: met. Previous must be enabled (Fatalf), both endpoints equal trailing minus 7 days, span equal, DrawnSegments non-zero, Next enabled, and Next returns the exact start instants.
- A4 Forward refusal outright: met. Two unconditional checks (arrow disabled, readout unchanged after press).
- A5 No patch chain, no consumer ids: met. REQ-164 and REQ-0003 appear only as fixture ids.
- A6 Test-only: met. One file, no `web/` or production Go.

### Vacuity and panic check
- Nil dereference: the shared `states` loop is extended with both new states and calls `t.Fatalf` if StartMs, EndMs or SpanMs is nil, before clause (3b) runs. Safe.
- Vacuity: Previous-enabled and Next-enabled-after-Previous are Fatalf, not skipped. The `afterPrev.DrawnSegments == 0` check cannot pass on an empty chart. The Next-returns-start check compares to `trailing`, which was itself guarded, so it cannot pass by both being null.
- Fixture fit: clause (4) needs the REQ-164 filtered fit under half of the unfiltered span. REQ-164 spans about one day and the unfiltered range reaches back 45 days, so it holds, and the span only widens with time. Clause (6) needs more than one hour of room right of the filtered extent in the all-days range. REQ-164 ends 2026-07-29 and the range ends at about now, so it holds. REQ-164 is the only fixture id that matches the filter. Both probes passed.
- Remaining soft spot: endpoints are parsed from the minute-truncated readout, so "exact" means exact to the minute. A 7-day shift keeps the same truncation, so this is correct, and the wrong-step mutation (reported in Testing) fails it.

### Findings
**Important:** None.

**Minor:**
- F1 The clause (3) comment says "The fixture tree has no forecast past now". The tree has two pending REQs and I did not verify that no forecast bars are drawn. The assertion holds in the run, but the comment states a reason that is unproven. -> impact-negligible -> report only

**Nit:**
- F2 The clause (6) comment still says the filtered width "is live queue data" and tells the old mid-drain story. It is history for the narrowing design and still reads true as the reason for the zoom-floor approach, but a reader may take it for the current regime. -> impact-negligible -> report only
- F3 The Now/Fit probe has no explicit assertion that REQ-0003 keeps the trailing 7 days drawn; only the previous week's `DrawnSegments` is checked. Fine as is. -> impact-negligible -> report only

### Stale-comment sweep
Searched the file for live-queue wording. The "busiest stretch of this repo's own archive" and "depends on the live queue" comments were rewritten in the diff. Remaining hits: line 1797 (fixture helper comment: "whatever state the live queue is in", accurate for the helper's purpose) and line 2337 (F2).

### Acceptance Testing
**Result: Pass** (implementation and integration stages only)
- Ran `go test -count=1 -v -run 'TimelineNowAndFit|TimelineProse|TimelineTrailing'` with the browser lane on and Chrome set: the three browser tests PASS (1.39 s, 1.41 s, 1.60 s), none skipped. The only SKIP lines are `TestJavaScriptBehaviorTimeline*` (other lane).
- Not re-run by me: the scratch-worktree run with the week's rows deleted and the one-day mutation. Both come from the builder's Testing section.
- Deployment and live acceptance: not applicable (test-only).

### Suggested Additional Testing
- Re-run the scratch worktree with archive rows for 2026-07-27 to 08-02 removed, to confirm the prose probe stays green.
- Heavy lanes in the Testing plan (queue-kanban-javascript, queue-kanban-browser, staged-skills) before release.

### Scores
**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All six criteria met in code |
| Code Quality | 92% | Two comment nits |
| Test Adequacy | 95% | Real RED/GREEN evidence, guarded assertions |
| Scope | 100% | One declared file |
| Risk | None | Test-only |
| Acceptance | Pass | Three browser probes pass, not skipped |

### Follow-ups created
None (3 findings report only)

## Review

**Overall: 96%** | 2026-10-10T11:45:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token):**
None

**Minor findings:**
- F1 clause (3) comment claims "no forecast past now" without proof; assertion itself holds - impact-negligible -> report only
- F2 clause (6) comment still says filtered width "is live queue data" - impact-negligible -> report only
- F3 no explicit drawn-segments check on the trailing window itself - impact-negligible -> report only

**Acceptance:** Pass - implementation and integration stages: the three browser probes passed with the lane on and none skipped; deployment and live acceptance not applicable.
**Restatement sweep:** nothing redefined
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

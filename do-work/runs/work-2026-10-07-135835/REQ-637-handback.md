# REQ-637 hand-back (board Panel B compares the stamp gap in whole minutes, and the drawer row gets a new label)

- Branch: `worktree-agent-REQ-637-board-gap-readers-compare-whole-minutes-and-relabel-drawer-row` (no name collision)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-637-board-gap-readers-compare-whole-minutes-and-relabel-drawer-row`
- Base: `13932fc2`
- Commits: `0c948168` (code + tests), `a49fde43` (board guide only, separate so it can be dropped; see D-03)

## File manifest (all modified, none new)
- `skills/do-work-board/tools/queue-kanban/durations.go` — `dayMedianExclusionReason` compares `int(largestStampGap.Minutes()) > int(activityGapCeiling.Minutes())`. Header comment: the phaseless case (one gap, whole claimed_at → completed_at span, excluded over 2h even when continuous), the whole-minute agreement with the log, and the drawer row's new name.
- `skills/do-work-board/tools/queue-kanban/durations_test.go` — two new cases in `TestImplementationSpanVerdictReadsTheLargestStampGap`.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` — label "Largest idle gap" → "Largest gap between events". Measurement unchanged.
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` — expected label updated (Requirement 4: test asserting the old label). Test function name kept, because archived REQ-632 records cite it.
- `skills/do-work-board/docs/board-guide.md` — row name updated, plus one sentence: the row does not count the open stretch to now; the card's `last activity` line shows that.

## Red-green evidence
- `TestImplementationSpanVerdictReadsTheLargestStampGap` RED:
  `a 2h00m40s phase gap is kept: in whole minutes it is the ceiling, not over it: verdict "idle-gap", want "" (span 150 min)` and
  `no phase stamps, span 59 seconds over the ceiling: still the ceiling in whole minutes, kept: verdict "idle-gap", want "" (span 121 min)`. GREEN: PASS. The existing case "a phase gap one minute over the ceiling is excluded" (2h01m00s) stays green, so the exclusion side is still pinned. The 2h00m40s case is the REQ's literal Red-Green fixture (claimed 2026-08-03T01:00:00Z, dispatch 03:00:40Z, completed 03:30:00Z).
- `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` (needs `QUEUE_KANBAN_JAVASCRIPT_PROBES=on`, otherwise it skips) RED:
  `gap row = "Largest idle gap": "1h 07m (dispatch → builder handback)", want "Largest gap between events": ...`. GREEN: PASS.

## P-A-U
- [x] **[PLAN]:** F4: truncate both sides to whole minutes in `dayMedianExclusionReason`, deriving 120 from `activityGapCeiling` (no restated number). RED cases: the REQ's 2h00m40s phase-gap fixture and a phaseless ceiling+59s span (the upper edge of truncation); the existing ceiling+1m case pins the exclusion side. Fix the header comment for the phaseless case. F6: relabel the drawer row in `board-detail.js` to "Largest gap between events", update the JS probe's expected label, and the guide paragraph that names the row. No change to `activity_correlation.go` or to the rule.
- [x] **[APPLY]:** Done as planned. Tests first (both RED confirmed), then `durations.go` and `board-detail.js`. Comment in `durations.go` reflowed and the drawer-row reference renamed. Guide paragraph updated in its own commit.
- [x] **[UNIFY]:** `git diff --stat`: 5 files, +25/−11. `gofmt -l .` empty. `go vet ./...` clean. Reviewed each file: durations.go (one comparison line + comment only), durations_test.go (two table cases), board-detail.js (label string only), javascript_behavior_a_test.go (two expected-label strings), board-guide.md (one paragraph). No debug artifacts. Full module suite with JS probes on: pass, 0 failures. Generated a static board from the worktree: the built `index.html` contains "Largest gap between events" once and no "Largest idle gap". Repo grep for "Largest idle gap" outside `do-work/`: only the historical CHANGELOG.md entry (history, left as is).

## Decisions
- D-01 (DECIDE & STATE): label is "Largest gap between events", the REQ's example. Short, and "between events" says the open gap to now is not included.
- D-02 (DECIDE & STATE): kept the test function name `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap`. Archived REQ-632 records cite it by name; renaming gains little.
- D-03 (DECIDE & STATE): edited `skills/do-work-board/docs/board-guide.md`, which is outside the write_set. It names the row in bold and would be stale otherwise. It is a separate commit (`a49fde43`) so the orchestrator can drop it if the boundary must hold.
- D-04 (DECIDE & STATE): the Panel B rule text `dayMedianExclusionRule` ("largest idle gap over 2h") is unchanged. It is the Panel B rule's name, not the drawer row, and REQ-633 chose it on purpose.

## Discovered Tasks
- impact-minor: `skills/do-work-board/tools/queue-kanban/activity_correlation.go:18` comment says the drawer "states its largest idle gap"; `skills/do-work/actions/estimate-reference.md` § Calibration says "the board's drawer idle-gap row". Both still describe the row by its old name. Outside this REQ's write set → report only
- impact-minor: `strict_behavior_regression_test.go` took 34.27s in the JS-probe run, over the 30s per-file budget. Untouched by this REQ; the machine had sibling sessions running → report only

## Lessons read
- `_dev/primes/lessons-releases.md` (whole)
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — `[family: paired-predicate-drift]` bullets only
- Primes: `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`. `_dev/primes/lessons-kanban-board.md` not read (dropped for budget, no family named).

## Proposed lesson (for `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`)
- [family: paired-predicate-drift] REQ-637: **two readers of one documented rule compared different units at the boundary.** The calibration log writes `max_stamp_gap_minutes` in whole minutes rounded down, and Panel B compared the exact `time.Duration` against the 2h ceiling. A 2h00m40s gap was excluded on the board and kept by a re-fit reading the logged 120. Both readers passed their own tests because every boundary fixture used exact minutes. The fix truncates both sides with `int(…Minutes())` derived from `activityGapCeiling`; the pin is a sub-minute fixture (2h00m40s and ceiling+59s kept, ceiling+1m excluded). When a rule has a logged form and a live form, test the boundary with a value the logged form rounds.

## Proposed CHANGELOG entry text
Board: Panel B and the calibration log now agree on the 2h gap boundary, and the drawer row has a clearer name.
- The Durations page's Panel B now compares the largest stamp gap in whole minutes rounded down, the same way the calibration log records `max_stamp_gap_minutes`. Before, a 2h00m40s gap left a request out of the day medians while a re-fit reading the logged 120 kept it. A 2h01m gap is still excluded.
- The drawer row "Largest idle gap" is now "Largest gap between events". It measures only between recorded events, never the open stretch from the last event to now. We relabelled it instead of adding that open stretch so the row means the same thing on claimed and done cards, and because the card's `last activity` line already shows the open stretch.
- The Panel B code comment now says that a request with no phase stamps has one gap, its whole claim-to-completion span, so it is excluded when that span is over 2h. The rule itself is unchanged.

## Integration seams
None. No change to `activity_correlation.go`, `verify.go`, or the payload shape.

## Test commands and results
- `go test -run TestImplementationSpanVerdictReadsTheLargestStampGap .` — RED then PASS (above).
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -run TestJavaScriptBehaviorDetailStatesTheLargestIdleGap -v .` — RED then PASS.
- `gofmt -l .` — empty. `go vet ./...` — clean.
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -json ./...` in `skills/do-work-board/tools/queue-kanban` — pass, 0 fails, package 104.8s. 32 tests skipped (browser probes; `QUEUE_KANBAN_BROWSER_PROBES` not set, so the browser lane was not run — no browser probe asserts the drawer gap row). Per-file wall time (sum of top-level tests): strict_behavior_regression_test.go 34.27s (over budget, see Discovered Tasks), generate_test.go 15.97s, citations_test.go 10.41s, verify_test.go 5.81s, javascript_behavior_c_test.go 2.64s, timeline_browser_probe_test.go 2.29s, javascript_behavior_a_test.go 1.01s, javascript_behavior_b_test.go 0.95s, durations_browser_probe_test.go 0.78s, javascript_behavior_d_test.go 0.76s, user_request_progress_browser_probe_test.go 0.69s, every other file under 0.5s (durations_test.go 0.06s).
- `queue-kanban generate` from the worktree into scratchpad: built page shows the new label once, old label absent.

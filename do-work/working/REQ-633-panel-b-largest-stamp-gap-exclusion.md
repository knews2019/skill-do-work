---
id: REQ-633
title: 'Panel B and the calibration log exclude by largest stamp gap, not raw span'
status: claimed
created_at: 2026-10-05T20:11:59Z
user_request: UR-135
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: ["REQ-632"]
batch: board-activity-evidence
depends_on: [REQ-632]
write_set: ["skills/do-work-board/tools/queue-kanban/durations.go", "skills/do-work-board/tools/queue-kanban/durations_test.go", "skills/do-work-board/tools/queue-kanban/generate.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/timeline.go", "skills/do-work-board/tools/queue-kanban/timeline_test.go", "skills/do-work-board/tools/queue-kanban/web/board-durations.js", "skills/do-work-board/tools/queue-kanban/web/board-timeline.js", "skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js", "skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go", "skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go", "skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go", "skills/do-work-board/docs/board-guide.md", "skills/do-work/actions/estimate-reference.md", "skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go", "skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go", "skills/do-work/tools/do-work-cli/internal/requeststate/*_test.go", "do-work/calibration-log.tsv"]
claimed_at: 2026-10-05T20:55:46Z
---
# Panel B and the Calibration Log Exclude by Largest Stamp Gap, Not Raw Span

## What
Replace the 4-hour raw-span exclusion with a rule on the largest gap between consecutive lifecycle stamps, 2 hours, applied identically by the two readers that claim one definition (`skills/do-work/actions/estimate-reference.md` → Calibration and `skills/do-work-board/tools/queue-kanban/durations.go`), rename the verdict string from `paused` to `idle-gap` in every reader, and add a `max_stamp_gap_minutes` column to `do-work/calibration-log.tsv` so the next re-fit can apply the same rule.

## Why
The 4h rule throws away long continuous sessions. Two of 190 were excluded at fit time; in the maintainer's repo three of the last five would be. As Route C runs grow, the raw-span rule discards the most informative samples and the timeline forecast inherits the bias. The calibration log carries only `wall_minutes`, so without a gap column the next re-fit cannot apply a gap rule at all. The maintainer chose one constant, 2h, with the phase named in the label, and asked that the archive's per-phase distribution be checked before the number is fixed.

## Detailed Requirements
1. Archive check first. Before changing code, compute the per-phase p95 of stamp gaps over every REQ in `do-work/archive/` using the intervals `buildPhaseBreakdown` emits, and record the table in this REQ's plan. If any phase other than dispatch → handback shows real multi-hour holes, stop and raise an Open Question proposing a second tier (builder phase 2h, other phases 45 min) instead of coding one constant. Otherwise proceed with one constant.
2. Rule (`durations.go`): `activityGapCeiling = 2 * time.Hour` replaces `analysisOutlierCeiling`. `dayMedianExclusionReason` takes the largest gap between consecutive parsed lifecycle stamps (the whole claimed → completed span when a REQ has no phase stamps) and returns `idle-gap`, `reversed`, or empty. Stamps only, on purpose: the calibration-log appender runs inside the lifecycle transaction and must not shell out to git, and both readers must agree byte for byte.
3. Rename the verdict string everywhere it is read: `generate.go` (the `implementationSpanReason` comment and the two assignment sites), `web/board-durations.js` (summary sentence, quantile skip, axis title, mark tooltip, day tooltip, table cell), `web/board-timeline.js` (forecast exclusion and assumptions text), `web/board-user-request-summary.js` (excluded count), `timeline.go` comments, `skills/do-work-board/docs/board-guide.md`. Panel B prose becomes "largest idle gap over 2h · excluded from day medians". The "2h" reaches the client as text shipped from Go, never as a number the client re-applies (lesson REQ-219).
4. Doc: rewrite the read-time rule in `estimate-reference.md` → Calibration to "exclude a span when the largest gap between consecutive lifecycle stamps exceeds 2h, or the span is negative", add one sentence naming `durations.go` as its second reader, and note that the 188-of-190 figure predates this rule. Update the `durations.go` header comment that calls itself the second reader so the two say the same thing.
5. Calibration log: add `max_stamp_gap_minutes` as the last column. Both appenders (`requeststate/state_plan.go` and `state_apply.go`) compute it from the REQ's stamps at write time. Ship a one-off idempotent backfill that fills old rows from archive frontmatter (location is the builder's call: `skills/do-work/scripts/` if it ships, else `_dev/`). The header row is the version marker; no separate marker.
6. Tests: a boundary pair straddling 2h by one minute (lesson REQ-374: only a straddling pair catches a second ceiling); a 4h21m span with a 67-minute largest gap is kept; a Route A REQ with no phase stamps falls back to the full span; a new TSV row carries the column; the backfill is idempotent on a second run; every JavaScript behaviour test that asserted `paused` now asserts `idle-gap` and the new prose.
7. Release: changelog entry and version bump per `_dev/primes/prime-releases.md`.

## Constraints
- Panel B reads archive frontmatter, not the TSV, so its medians shift as soon as the rule lands. That is accepted.
- The card-side gap from REQ-632 (Board cards show last correlated activity) uses stamps ∪ commits; this rule uses stamps only. A commit can only shrink a gap, so Panel B may exclude a REQ whose drawer shows a smaller gap. Nothing on the card shows the exclusion, so nothing visibly disagrees; say this once in the doc.
- `reversed` keeps its name and meaning.
- Shipped files change in three packages (board, core action, core CLI), so this is a release; the backfill script, if shipped, follows `_dev/primes/prime-shell-commands.md`.

## Dependencies
Depends on REQ-632 (Board cards show last correlated activity and drop the assumed-pause badge): the badge text field must already be gone, and the drawer row this REQ's wording refers to must exist.

## Builder Guidance
Certainty is high on the shape (one constant, stamps-only, both readers, new column) and medium on the number: 2h is the maintainer's pick pending the archive check in requirement 1. Builder latitude: the backfill script's language and location, the exact Panel B sentence wording, and whether `ImplementationSpan.ExclusionReason` is renamed in Go or only its values.

## Red-Green Proof
**RED prompt/case:** A completed REQ with a 4h21m span and a 67-minute largest stamp gap, and a second completed REQ with a 2h01m gap between `dispatch_at` and `builder_handback_at` inside a 3h span.
**Why RED now:** Panel B excludes the first (`excludedReason: "paused"`) and keeps the second; `do-work/calibration-log.tsv` has no column from which a re-fit could tell them apart.
**GREEN when:** The first is kept and the second is excluded with reason `idle-gap`; the Panel B summary names "largest idle gap over 2h"; the TSV header ends with `max_stamp_gap_minutes` and every row has a value; `grep -rn '"paused"' skills/do-work-board/tools/queue-kanban/*.go skills/do-work-board/tools/queue-kanban/web/*.js` returns nothing; `estimate-reference.md` and `durations.go` state the same rule.
**Validation:** User confirmed (plan approved 2026-10-05)

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: REQ-374 (only a straddling pair catches a second ceiling) and REQ-219 (ship the verdict, never the inputs) govern this change and are restated in requirements 3 and 6.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7095 tokens, over the 2000 budget; `slugged: partial`). Matched: changing queue-kanban model and UI.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18027 tokens, over the 2000 budget; `slugged: partial`). Matched: this REQ changes the calibration-log appenders in `requeststate`.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over the 2000 budget; `slugged: partial`). Matched: a shipped backfill script is shipped shell.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-135/input.md` for complete verbatim input.

*Source: C. The Panel B calibration exclusion can keep its own rule, but it should be decided from the same gap evidence, not from the raw span … Both readers (estimate-reference.md and durations.go) must change together — the doc calls itself the single definition.*

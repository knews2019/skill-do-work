---
id: UR-062
title: Upstream consumer review comments on the board views and qualify
created_at: 2026-08-20T08:37:41Z
requests: [REQ-303, REQ-304, REQ-305]
word_count: 447
---

# Upstream Consumer Review Comments on the Board Views and Qualify

## Summary

Five review findings pasted for triage through `do-work-toolbox validate-feedback`, from the same
consumer-report stream as UR-057's `upstream-consumer-report-2026-08-19` batch. All five were
verified against the code and all five were accepted; two of the reviewer's stated premises were
refuted in the process and the corrections are recorded per REQ.

## Extracted Requests

| Finding | Verdict | Disposition |
|---|---|---|
| P1 — consumer-specific live archive constants in `durations_test.go` | Accept, severity overstated | REQ-303 |
| P2 — negative waits render as valid spans in the Timeline | Accept, scoped | REQ-304 |
| P2 — filtered timeline against a global forecast | Accept, cheaper remedy | REQ-305 |
| P2 — rounded duration remainders do not carry | Accept | fixed in place, 0.216.1 (`1311300`) |
| P2 — untracked files escape qualify's artifact scan | Accept | Addendum on queued REQ-263 |

## Batch Constraints

- The rounding finding was fixed in place rather than captured: its test seat already existed
  (`generate_test.go`'s node probes over the sliced formatters), it collided with nothing queued,
  and it was under the threshold where a REQ pays for itself. Shipped as 0.216.1.
- The untracked-files finding was appended to queued REQ-263 rather than captured as its own REQ:
  REQ-263's write set is `skills/do-work/tools/checks/qualify.sh` plus
  `_dev/tests/prescribed-shell-scripts-behavior.sh`, exactly the files that finding's remedy
  touches, and two REQs writing that pair would collide.
- REQ-304 overlaps queued REQ-280, which already owns adding the
  `created_at <= claimed_at <= completed_at` ordering probe to `queue-kanban verify`. REQ-304 is the
  RENDERING half only; the anomaly verdict itself must not be recomputed in the timeline
  (`timeline.go:22-24` — the board decides what counts as broken bookkeeping in exactly one place).
- Two reviewer premises were refuted during triage and must not be carried into implementation as
  fact: (1) `go test ./...` does NOT fail in this checkout — the whole suite passes and the
  2026-07-31 median is 2.5 as pinned; the 9.0833 figure came from the reviewer's own tree, which is
  the portability bug itself. (2) The untracked-file gap is wider than reported — the
  `debugger|TODO|FIXME` scan misses untracked files too, not only the `print(`/`console.log` half.

## Full Verbatim Input

do-work validate-feedback: Full review comments:

  - [P1] Replace consumer-specific live archive constants — [prj]/.claude/skills/do-work-board/tools/queue-
    kanban/durations_test.go:220-221
    liveBoard loads the consuming repository's actual do-work/ tree, but this test asserts exact medians and counts from one upstream
    archive. In the current checkout, go test ./... already fails here because the July 31 median is 9.0833 rather than 2.5; consumers
    with different histories can fail or lack the date entirely. The tool is explicitly portable to any repository (.claude/skills/do-
    work-board/tools/queue-kanban/prime-do-kanban.md:23-25), so these values need deterministic fixtures or invariant assertions.

  - [P2] Render negative waits as broken spans — [prj]/.claude/skills/do-work-board/tools/queue-kanban/web/
    board-timeline.js:640-643
    When claimed_at precedes created_at—or an unclaimed ticket has a future created_at—waitMinutes is negative, but this unconditional
    call passes both endpoints to drawSegment, which sorts them with min/max and paints a normal waiting bar. Only negative work spans
    receive the broken marker below, and the aggregate's anomaly flag only copies completion anomalies, so reversed waits are visually
    presented as valid; add equivalent negative-wait handling and anomaly reporting.

  - [P2] Keep filtered timelines and forecasts in sync — [prj]/.claude/skills/do-work-board/tools/queue-kanban/
    web/board-timeline.js:554-554
    When a search, domain, or status filter leaves a nonempty subset, rows and the summary/table are filtered, but projection remains
    global and renderTimelineForecast does not use its rows argument. The view can therefore say it contains three REQs while
    forecasting and listing exclusions for the entire queue, including hidden IDs; recompute/filter the projection or suppress and
    clearly label the global forecast while filters are active.

  - [P2] Carry rounded duration remainders into the next unit — [prj]/.claude/skills/do-work-board/tools/queue-
    kanban/web/board-durations.js:144-144
    For spans whose fractional remainder rounds to 60, such as 119.5 minutes, this renders 1h 60m instead of 2h 0m. The same rollover
    bug exists in timelineFormatSpanMinutes for minutes and hours and in Go's formatDurationLabelMinutes, affecting chart labels,
    hover text, tables, forecasts, and label-width planning; round first and carry overflow into the next unit in all mirrored
    formatters.

  - [P2] Scan untracked implementation files for output artifacts — [prj]/.claude/skills/do-work/tools/checks/
    qualify.sh:210-210
    In serial mode this loop only walks changed_file_list, which is built from working and staged git diff --name-only output and
    therefore excludes untracked files. A newly created source file containing print() or console.log is never inspected, so a checked
    [UNIFY] can pass with leftover instrumentation; include git ls-files --others --exclude-standard or the Implementation Summary's
    new-file paths. This script is the required serial debug-artifact gate per .claude/skills/do-work/actions/work.md:441-443.

---
*Captured: 2026-08-20T08:37:41Z*

---
session_ended: 2026-08-24T13:15:00Z
last_completed: REQ-346
queue_state: [15 pending, 1 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 3 in-progress]
reqs_processed_this_session: 9
session_depth: deep
---

# Session Checkpoint

## Completed This Session

Nine REQs archived, releases 0.236.34 through 0.236.38.

- REQ-339: Count every prescribed-shell case the runner reports (Route B, 89%)
- REQ-340: Finish the report-image interruption sweep (Route B, 75%, Partial)
- REQ-342: Neutralize user text written into a REQ body (Route C, 86%)
- REQ-343: Let verify see a structurally damaged REQ file (Route B, 93%)
- REQ-344: Quote user text written into a frontmatter value (Route C, 78%)
- REQ-345: Stop a queued REQ from failing the timeline landing probe (Route C, 95%)
- REQ-341: Give the browser probe lane trusted input (Route C, 96%)
- REQ-346: Name the UR on every Durations sample (Route B, 87%)

## In Progress (interrupted)


Three REQs are claimed for a bounded fan-out wave. Builders run in isolated worktrees. The user
stopped the run before they handed back. **Their branches may hold uncommitted
or committed work that was never merged** — check each worktree before
re-dispatching, and prefer resuming the existing branch over starting fresh.

## Still Queued

15 pending, 1 `pending-answers`. Thirteen of these did not exist when the run
started — they are review findings, each verified against real code or a real
render before capture.

**Durations batch (UR-069), dependency-ordered:**
- REQ-348 (Timeline UR row grouping), REQ-349 (panel A scale/density), REQ-350 (axis window control), REQ-351 (retire in-lane labels, depends REQ-346 ✓), REQ-352 (headline numbers, depends REQ-350), REQ-353 (hide dead knobs), REQ-354 (drawer from a mark)
- **Seven of the nine write `web/board-durations.js`.** They serialize on it — do not fan them out.

**Review follow-ups:**
- REQ-355 (bound the case-header rule at the script name) — `impact-negligible`
- REQ-357 (**`pending-answers`** — two judgment calls on the structural probe awaiting the maintainer)
- REQ-359 (suppress the duplicated status finding from the board strip)
- REQ-360 (**sweep**, `neutralization-contract-reach` — five instances across REQ-342's and REQ-344's sibling contracts)
- REQ-361 (make the neutralization lock-in catch narrowing)
- REQ-362 (stop a multi-path bullet disabling the scope-drift check)
- REQ-363 (exempt supported UR-less REQ shapes from structural damage)
- REQ-364 (pin the UR lane's geometry)
- REQ-365 (a tdd REQ must name a test file in its write set)

## Session Notes

**Gate state: green**, with the strict browser lane actually run rather than
skipped (`TestMaintainerStrictBrowserBehaviorLane` PASS, 27.69s). Note that a
gate run without a browser on `PATH` SKIPs that lane — most of this session's
earlier green reports were of the narrower kind. Supply
`QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium-1194/chrome-linux/chrome`.

**Environment:** the container needs `GOTOOLCHAIN=go1.26.1` (system Go is
1.24.7, below the repo's floor), ShellCheck 0.11.0 and `just`, all installed
this session.

**Three corrections worth carrying forward, each a case of a check measuring
less than it claimed:**

1. REQ-340's instance 3 was refuted, and I confirmed the refutation — both of us
   measured `kill -0 <pid>` when the code takes the group branch. The group
   branch burns the full second. REQ-356 reopens it.
2. REQ-342's lock-in was reported as catching narrowing. It catches deletion and
   phrase substitution; six narrowings by added qualifier ship green, two of
   which reintroduce the original defect. REQ-361 carries the fix.
3. REQ-343's real-tree false-positive check passed because neither offending
   shape exists in the tree today. It measured the corpus, not the writers.
   REQ-363 carries the carve-outs.

**One self-inflicted:** a blind string replace in my own bookkeeping matched both
a real `## Scope Extensions` heading and a citation of it inside backticks,
forging a heading in REQ-344's record — the exact damage class UR-068 exists to
prevent. Repaired and recorded there.

**Promoted to the prime rather than a REQ:** the surface behind this board's SVG
is `<body>`, not any `--surface-*` token. REQ-321 found it under the timeline
bars, REQ-346 under the durations lane; two builders paying for the same
property is what `_dev/primes/prime-kanban-board.md` § Conventions exists to
stop.

---
id: REQ-283
title: Route verify through the do-work-board skill that owns it
status: completed
claimed_at: 2026-08-20T22:43:20Z
completed_at: 2026-08-20T22:46:10Z
commit: 2308afd
created_at: 2026-08-19T13:42:45Z
user_request: UR-057
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-279, REQ-280, REQ-281, REQ-282]
batch: upstream-consumer-report-2026-08-19
effort_estimate: trivial
write_set:
- skills/do-work-board/SKILL.md
- skills/do-work-board/actions/board.md
---

# Route Verify Through the do-work-board Skill That Owns It

## What

`skills/do-work-board/SKILL.md`'s routing table has two rows, `help` and `board`, and passes `serve`, `static`, `summary`, `cli`, `--port N`, `--out DIR` through to the board action. `skills/do-work-board/actions/board.md`'s mode table maps `cli` to `open-work` and never mentions `verify` — `grep -n -i verify` on it returns nothing.

So `do-work-board verify` falls to SKILL.md's rule "an unknown command prints board help and stops," for a real subcommand of the tool this package owns. Add the routing row and the board-action mode.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add `verify` to the board router and argument hint, add the read-only mode to the board action, reuse the existing build-and-run flow, and point probe descriptions at Check 14.
- [x] **[APPLY]:** Code written exactly as planned. Scope stayed within the board router and board action.
- [x] **[UNIFY]:** Reviewed the 2-file diff and git diff --stat; ran gofmt-independent documentation checks, go test, go build, and the verify command. No debug artifacts or Go changes are present.

## Why

`verify` is not a corner of the tool. `actions/forensics.md` Check 14 documents its whole probe set, and `actions/work-reference.md:63` and `:345` both lean on it — the `duplicate-req-id` probe as the detector for colliding captures across checkouts, and the committed-queue-state probe as the guard at Step 9's merge. The package that ships it cannot route to it, so the obvious phrasing gets board help instead.

## Context

The report's severity framing — "it made a documented diagnostic unrunnable as typed" — overstates it, and the REQ should not inherit that. Check 14 documents a build-then-run of the binary directly and that command works today. What is actually wrong is discoverability through the owning skill's own router.

## Requirements

- Add a routing row to `skills/do-work-board/SKILL.md`: trigger `verify` (plus natural aliases such as "check invariants", "probes"), routed to `./actions/board.md`, and add `verify` to the token list SKILL.md passes through.
- Add a `verify` row to `actions/board.md`'s Input mode table: read-only, prints the tool's `verify` output, exit 1 means findings rather than an error.
- Add the mode to the action's step flow reusing the build-then-run contract every other mode already uses (Step 1 locate, Step 2 Go precondition, Step 3 repo root).
- Update SKILL.md's `argument-hint` frontmatter so the advertised argument set matches the routing table.
- Point at `../../do-work/actions/forensics.md` Check 14 as the canonical description of the probe set rather than restating the probes, so the two cannot drift.

## Constraints

- **Do not restate the probe table** in `board.md`. Check 14 owns it; a second copy is exactly the drift the report warned about and the reason it asked for a pointer.
- Read-only toward the pipeline, like every other board mode. `board.md`'s three-write-surface statement stays true and unedited — `verify` writes nothing, so it adds no fourth surface and CLAUDE.md's write-surface sentence needs no amendment.
- Keep the change to routing and documentation. No Go changes: the subcommand already exists and works.

## Red-Green Proof

**RED prompt/case:** Follow `skills/do-work-board/SKILL.md` as written for the input `do-work-board verify`.
**Why RED now:** The routing table has no matching row, so the documented rule "an unknown command prints board help and stops" fires and the agent prints board help — for a subcommand `main.go:79` implements and `actions/forensics.md` Check 14 documents in full.
**GREEN when:** Following SKILL.md for the same input routes to `actions/board.md`, whose mode table has a `verify` row that builds the tool and runs `queue-kanban verify --repo-root <project-root>`, reports exit 1 as findings rather than an error, and links Check 14 for the probe descriptions. `grep -n -i verify skills/do-work-board/actions/board.md` returns matches.
**Validation:** Inferred during capture; both the routing table and the mode table were read and neither mentions `verify`.

## Full Context

See `do-work/user-requests/UR-057/input.md` for the complete verbatim upstream report.

---
*Source: upstream defect report D4, severity low, from `g1w-game-find-the-difference` running v0.212.25 — verbatim claim: "`verify` is unreachable through the skill that owns it … Neither mentions `verify` … So `do-work-board verify` hits the rule 'an unknown command prints board help and stops' — for a real subcommand of the tool this package owns." Accepted by `do-work-toolbox validate-feedback` triage (2026-08-19), with its "unrunnable as typed" framing corrected: Check 14's direct binary invocation works, and the gap is routing through the owning skill. Evidence: `skills/do-work-board/SKILL.md` routing table (two rows, no `verify`); `grep -n -i verify skills/do-work-board/actions/board.md` empty; `skills/do-work-board/tools/queue-kanban/main.go:79` implements the subcommand; `skills/do-work/actions/forensics.md:181-200` documents its probes; `skills/do-work/actions/work-reference.md:63,345` depend on two of them. Surface-cost: N/A — a routing table row and a mode entry, no defensive surface.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names both documentation files and the existing `verify` binary subcommand. It requires routing and a build-then-run mode entry only.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/SKILL.md` (modified) — routes `verify` aliases and advertises the mode in the argument hint.
- `skills/do-work-board/actions/board.md` (modified) — documents the read-only verify mode, its build/run step, and Check 14 ownership.

## Qualification

Passed — 2 project files verified, all routing and execution requirements traced, and no production Go code changed.

## Testing

- `go test ./...` passed in `skills/do-work-board/tools/queue-kanban`.
- Fresh `go build -o queue-kanban .` passed.
- `./queue-kanban verify --repo-root "$(git rev-parse --show-toplevel)"` returned `OK: no findings` with exit 0.
- The router and action now both contain `verify`, the accepted aliases, and the Check 14 pointer.

## Review

**Overall: 100%**

**Requirements:** Pass — routing, argument hint, mode table, execution flow, read-only semantics, exit-1 meaning, and canonical Check 14 pointer are all present.

**Code quality:** Pass — documentation-only, consistent with existing board modes, and no duplicated probe table.

**Test adequacy:** Pass — the owned Go test suite and the actual verify command pass.

**Scope:** Pass — only the two declared files changed.

**Acceptance:** Pass — `verify` is now discoverable and runnable through its owning skill.

## Lessons Learned

**What worked:** The existing binary already owned the behavior, so the fix stayed at the routing/documentation boundary.
**What didn't:** Running the binary from its nested module without an explicit project root makes it look for `do-work/` in the module directory; the action's explicit root resolution is essential.
**Worth knowing:** Verify findings use exit 1 as a diagnostic result, not as a launcher failure.

## Orientation

The board skill now exposes its existing cross-file invariant verifier as a first-class `do-work-board verify` mode, with execution and probe ownership clearly separated.

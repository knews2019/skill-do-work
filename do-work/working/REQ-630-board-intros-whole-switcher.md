---
id: REQ-630
title: 'Board intros describe the board by what its switcher covers'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-02T21:52:56Z
created_at: 2026-10-02T21:51:57Z
user_request: UR-134
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-631"]
batch: ur-133-leftovers
write_set: ["skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
dispatch_at: 2026-10-02T21:53:17Z
builder_handback_at: 2026-10-02T21:53:17Z
integration_at: 2026-10-02T21:53:17Z
claimed_at: 2026-10-02T21:52:48Z
---
# Board Intros Describe the Board by What Its Switcher Covers

## What
Rewrite the opening line of `skills/do-work-board/docs/board-guide.md` and of `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` so neither names only some of the board's pages.

## Why
The REQ-628 review (finding M1) found both intros still read like "a Kanban board ... plus a queue activity calendar and a testing track", three of six pages. The lesson "describe the board by what its switcher covers, not by naming some of its tabs" (REQ-232, `_dev/primes/lessons-kanban-board.md`) says this drifts every time a page is added.

## Detailed Requirements
- Each intro describes the board by what its switcher covers (for example "the Board view plus the other pages in its switcher"), or names all six pages; prefer the wording that does not need editing when a page is added.
- Keep the rest of each opening paragraph (read-only toward the pipeline, the testing record) unchanged.
- Grep shipped files under `skills/do-work-board/` for other partial page lists in intros and report any in Discovered Tasks.

## Constraints
Docs only. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-631.

## Red-Green Proof
**RED prompt/case:** Line 3 of `board-guide.md` and of `prime-do-kanban.md` names a calendar and a testing track as the board's extras, leaving out Activity, Timeline and Durations.
**Why RED now:** Both intros predate the six-page switcher.
**GREEN when:** Neither intro names a partial set of pages.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: views. The REQ-232 bullet this REQ cites carries no family marker, so no narrower entry exists; the REQ's Why names it and the builder reads it directly.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: replace each partial page list with "plus the other pages in its page switcher", which needs no edit when a page is added.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: only the two declared files changed.
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: `git diff --stat` 2 files, 2 insertions, 2 deletions; no Markdown linter here; both lines re-read; grep for other partial page lists under skills/do-work-board found two, reported in Discovered Tasks.

## Full Context
See `do-work/user-requests/UR-134/input.md` for complete verbatim input.

*Source: ok, do board intros and the recover message*

---

## Triage

**Route: A** - Simple

**Reasoning:** Docs-only change to two named lines, with the wording rule given by the REQ.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** Both intros now describe the board as a Kanban board "plus the other pages in its page switcher" instead of naming a calendar and a testing track, so a new page needs no intro edit. The rest of each paragraph is unchanged. Merge range 875d49b4..bdb5c43c (builder commit 96c0465b, merge bdb5c43c).

## Qualification

**Diff range:** 875d49b4..bdb5c43c (builder commit 96c0465b, merge bdb5c43c)
**Gate records:** qualify satisfied. Route A, so no Scope comparison.
**Warnings judged:** none.
**Orchestrator read of the diff:** both lines now name the board plus "the other pages in its page switcher", which follows the REQ-232 lesson and stays true when a page is added; the rest of each paragraph is byte-identical.
**P-A-U honesty:** the orchestrator played the builder role on branch `worktree-agent-REQ-630-board-intros`; APPLY cross-checked against `git diff --stat 875d49b4..bdb5c43c` (two files, nothing under do-work/).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 00ad771a (merge bdb5c43c plus REQ trail)
**Result:** ✓ All passing — exit 0, gate wall 149s. Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-02-215500/helpers/probe-630.sh` (the Red-Green Proof as a check: neither old partial list remains, both intros carry the switcher wording) → exit 0 (advance probe record satisfied). Docs-only, so no harness test; the probe fails on the pre-change tree at 875d49b4.

**Red-green validation:** traced to `## Red-Green Proof`: RED both intros named a calendar and a testing track at 875d49b4 → GREEN neither names a partial set at bdb5c43c.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 875d49b4..bdb5c43c
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

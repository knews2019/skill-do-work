---
id: REQ-307
title: "[impact-negligible] Standing prose-reconciliation sweep"
status: completed
created_at: 2026-08-20T13:21:13Z
user_request: UR-063
domain: general
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: []
completed_at: 2026-08-20T22:37:38Z
commit:
---

# Standing Prose-Reconciliation Sweep

> **The mechanism described below was retired before this REQ ever drained.** Everything from `## What`
> through `## Context` is the original request, kept as the record of what was asked. What shipped is
> `do-work/prose-backlog.md` — see `## Implementation Summary` at the end.

## What

**This REQ never closes.** It is the destination REQ-306 creates for a class of finding that
should not mint its own REQ: a prose-only discrepancy in the skill's own text — a stale count, a
wrong cross-reference number, a comment describing a mechanism that was superseded, a restatement
that drifted from the thing it restates.

The class is defined by what the fix touches, not by how it was found: **no behavior changes, no
checker's predicate changes, no rule's stated condition changes.** If a fix would alter any of
those three, it is not this REQ's business and earns its own capture (the boundary's canonical home
is `skills/do-work/actions/capture-reference.md` → **Fold-First Rule**, shipped by REQ-306; do not
re-derive it here).

`write_set` is deliberately empty and stays empty. A standing sweep touches whatever files its
current instances name, so a declared write set would be stale the moment an instance is appended.
The board's overlaps badge reads `write_set`, so an empty field means *unknown* rather than
*nothing* — which is the honest state for a REQ whose scope is a growing list.

## Cadence

**The default work scan never selects this REQ** (`skills/do-work/actions/work.md` Step 1) — it
drains only when explicitly named (`do-work run REQ-307`) or opportunistically. Two guides for
choosing when to name it:

- **Roughly eight entries** is a batch large enough that one drain and one commit are amortized
  across it, and small enough that the instances have not aged past the point where the reader can
  still tell whether each one is still true. Advisory, not a gate; the queue status summary's
  `(standing sweep …: K open instances)` line is the signal.
- **Opportunistic folding beats waiting.** When any REQ is claimed whose declared scope already
  contains a file this list names, fix that instance in that REQ and tick it here. That is
  cheaper than a dispatch and it is the mechanism the Fold-First Rule points at first.

A drain does not close this REQ. It ticks the drained instances, commits, and leaves the REQ
`pending` with whatever remains. If the list is empty, there is nothing to drain and that is the
healthy state — an empty standing sweep is not a stale REQ.

**Every instance is re-verified at drain time, not trusted from capture.** An instance recorded
weeks ago may have been fixed in passing by an unrelated commit; three of the four REQs processed
on 2026-08-20 found at least one stated premise already stale. Tick a stale instance with that
evidence rather than editing text that is already correct.

## Instances

Both instances moved to `do-work/prose-backlog.md` unfixed and re-verified. Nothing here is
actionable — read the backlog file, not this list.

## Context

Created by explicit instruction under UR-063, alongside REQ-306, which is the rule that routes
findings here. This REQ exists so that rule has a destination on its first day rather than a
forward reference.

**REQ-273 is evidence against the rule it violates, and that is the argument for batching this
class.** `CLAUDE.md`'s Kanban-board write-surface paragraph ends: "Adding a fourth write surface
means amending this sentence in the same commit; that is the co-location rule applied to itself."
When the third surface was added, `CLAUDE.md` was amended and `frontmatter_cli.go` was amended —
and `open_work.go` and `testing.go` were not. So co-location worked for the two sites someone was
looking at and missed the two they were not. A rule that depends on remembering every restatement
at edit time will keep producing exactly this residue; batching it accepts that and makes the
residue cheap instead of pretending it will not accumulate.

That is also why this REQ carries `impact-negligible` honestly rather than defensively: nothing
here misleads a machine. Each instance misleads a *reader*, one hop from the truth, which is worth
fixing and is not worth a dispatch each.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Implementation Summary

**Delivered as `do-work/prose-backlog.md`, not as a never-closing REQ.** This REQ's purpose was to
give the Fold-First Rule's prose-only destination somewhere to land on its first day. That
destination now exists as a plain checklist file outside the pipeline, and the mechanism that made
it a REQ is removed: `standing: true`, the never-closing status, the never-auto-select rule, the
terminal-set and UR-closure carve-outs, the `## Drains` provenance line, Step 8's substep skipping,
Step 9's staging substitution, and the board's `standing` parsing. The canonical rule is
`skills/do-work/actions/capture-reference.md` -> **Fold-First Rule**, destination 3. The decision is
recorded as ADR-021.

Both instances migrated verbatim to `do-work/prose-backlog.md` and were re-verified against the
tree while migrating: the "forensics Check 11" miscitation still holds at all four sites, and the
"two write surfaces" count still holds at `open_work.go` (now line 23, was 22) and `testing.go:42`.
Neither is fixed here -- they are open backlog items, and fixing them is an ordinary REQ.

**The intent survives; only the mechanism was cancelled.** The prose-only class still has a
guaranteed destination and no recorded finding was dropped. What changed is that the destination
stopped being a pipeline REQ, which is what forced every selector, filter, terminal-set reader, and
provenance rule to learn a special case.

Files changed: `do-work/prose-backlog.md` (new), plus the mechanism removal across
`skills/do-work/actions/`, `skills/do-work/tools/select-simple-reqs.sh`,
`skills/do-work-toolbox/actions/code-review.md`, and
`skills/do-work-board/tools/queue-kanban/`.

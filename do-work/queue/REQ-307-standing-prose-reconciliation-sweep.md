---
id: REQ-307
title: "Standing prose-reconciliation sweep"
status: pending
created_at: 2026-08-20T13:21:13Z
user_request: UR-063
domain: general
sweep: true
sweep_key: prose-only-discrepancy-reconciliation
standing: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: []
---

# Standing Prose-Reconciliation Sweep

## What

**This REQ never closes.** It is the destination REQ-306 creates for a class of finding that
should not mint its own REQ: a prose-only discrepancy in the skill's own text — a stale count, a
wrong cross-reference number, a comment describing a mechanism that was superseded, a restatement
that drifted from the thing it restates.

The class is defined by what the fix touches, not by how it was found: **no behavior changes, no
checker's predicate changes, no rule's stated condition changes.** If a fix would alter any of
those three, it is not this REQ's business and earns its own capture (REQ-306 states that boundary
canonically; do not re-derive it here).

`write_set` is deliberately empty and stays empty. A standing sweep touches whatever files its
current instances name, so a declared write set would be stale the moment an instance is appended.
The board's overlaps badge reads `write_set`, so an empty field means *unknown* rather than
*nothing* — which is the honest state for a REQ whose scope is a growing list.

## Cadence

**Drain when the instance list reaches roughly eight entries, or when a claimed REQ's write set
already includes an instance's file — whichever comes first.** Both halves matter:

- **Eight** is a batch large enough that one dispatch, one review, and one commit are amortized
  across it, and small enough that the instances have not aged past the point where the reader can
  still tell whether each one is still true. The number is a default, not a gate; nine is fine and
  so is six.
- **Opportunistic folding beats waiting.** When any REQ is claimed whose declared scope already
  contains a file this list names, fix that instance in that REQ and tick it here. That is
  cheaper than a dispatch and it is the mechanism REQ-306's rule points at first.

A drain does not close this REQ. It ticks the drained instances, commits, and leaves the REQ
`pending` with whatever remains. If the list is empty, there is nothing to drain and that is the
healthy state — an empty standing sweep is not a stale REQ.

**Every instance is re-verified at drain time, not trusted from capture.** An instance recorded
weeks ago may have been fixed in passing by an unrelated commit; three of the four REQs processed
on 2026-08-20 found at least one stated premise already stale. Tick a stale instance with that
evidence rather than editing text that is already correct.

## Instances

- [ ] **`repair-req-timestamps.sh:7`, `:22`, `:136`, and `actions/work-reference.md:285`** — all four
  cite "forensics Check 11" for the future-dated-timestamp check. Check 11 is *Unrecognized Status
  Values*; the future-stamp check is `### 12. Future-Dated Timestamps` (`actions/forensics.md:156`).
  Verified against the tree on 2026-08-20. Folded from REQ-272.
- [ ] **`skills/do-work-board/tools/queue-kanban/open_work.go:22` and `testing.go:42`** — both say the
  board tool has "two write surfaces". It has three. `frontmatter_cli.go:34` already says "exactly
  three", so the count is stated correctly in one Go file and staled in two others. Verified against
  the tree on 2026-08-20. Folded from REQ-273.

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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

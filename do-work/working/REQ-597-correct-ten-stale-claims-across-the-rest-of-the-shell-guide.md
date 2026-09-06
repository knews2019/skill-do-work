---
id: REQ-597
status: claimed
domain: general
created_at: 2026-09-06T05:49:08Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-06T06:59:15Z
  basis:
    - Route B
    - 3-file write set
    - 4 acceptance criteria
maintenance: true
depends_on: [REQ-596]
related: [REQ-555, REQ-595, REQ-596]
title: 'Correct ten stale claims across the rest of the prescribed-shell guide'
claimed_at: 2026-09-06T06:59:15Z
status_changed_at: 2026-09-06T12:38:06Z
---

# Correct Ten Stale Claims Across the Rest of the Prescribed-Shell Guide

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B. Three builders in one worktree, stacked: the two broken `inspect.md` blocks first, then the sixteen guide claims re-derived from the code, then the two callers' association prose.
- [x] **[APPLY]:** Three commits `a1e652f`, `6913dc4`, `7df6488`, merged as `d5cf28b`. Three files, all in the write set.
- [x] **[UNIFY]:** `git diff --stat 804a8ba..d5cf28b`: 3 files, +23/-23. Four guards green on the merged tree (audit-lockins, prescribed-shell-canonicalization, quiet-grep-pipeline-audit, action-shell-blocks). No debug artifacts, no version or changelog touched.

## What

A 69-claim sweep of `skills/do-work/docs/prescribed-shell-primitives.md` against the Go code found 21
claims that do not hold. REQ-596 corrected the eleven inside the two sections it owned. Ten remain, in
sections REQ-596 never opened, at guide lines 30, 32, 80, 92 (two), 106 (two) and 128 (two).

## Why

Same finding class as REQ-555, REQ-595 and REQ-596: the guide is the pointer target from sixteen shipped
files and describes implementations that do not exist. Two of the ten are marked `stale` rather than
`partly-stale`, meaning the sentence is false rather than incomplete.

## Context

Found by the section sweep run for REQ-596. The full report with each claim, the code that contradicts
it, a file:line citation and suggested replacement text is in
`do-work/runs/work-2026-09-05-231943/REQ-596-section-sweep.json`, third element of the array. That file
is the input to this request — do not re-derive the census, but do re-verify each claim against the code
before rewriting it, because a replacement sentence that is also wrong is worse than the one it replaces.

## Detailed Requirements

- Correct the ten claims the sweep names, in the lifecycle-timing, merge-aware-commit-diff,
  commit-file-listing, verified-exact-publication and portfolio-summary sections.
- Two are `stale` rather than `partly-stale` and should be checked first: the claim that the
  diff-tree form is the suite default for path-only consumers, and the claim about how `ln` and `mv`
  treat a directory standing in the destination's place.
- Check every other claim in each section you open, the way REQ-595 checked all fourteen Mechanics cells
  and REQ-596 checked its two sections. The sweep reports 37 claims checked across these sections; a
  claim it marked accurate is evidence, but a claim it did not reach is not.
- No sentence may generalize over two commands that behave differently. That shape has been caught three
  times in this file within the visible history — REQ-555, REQ-595 and REQ-596 — and the sentence
  REQ-596 replaced predates that history.
- **Two shipped callers now contradict the guide and are in this request's write set.**
  `skills/do-work/actions/commit.md:85` still names the `do-work/archive/**/REQ-*.md` glob REQ-596
  replaced, and still says the delegated check "prints one `<owner>\t<path>` row per candidate — a
  `REQ-NNN` id, or `-` for unassociated", with line 87 adding that files coming back `-` remain
  unassociated. `protected-inventory associate` never prints a `-` row; only the separate
  `associate-files.sh` entry point does. Check `../do-work-toolbox/actions/inspect.md` for the same two
  claims. The guide is the canonical home those actions point at, so a caller disagreeing with it is the
  same defect one file over.

## Constraints

- Guide prose only: no behaviour change, no code change, no route-column change, no Mechanics-column
  change, no heading change.
- `_dev/tests/audit-lockins.sh` Finding 7 pins the route column, the orchestration claim and the
  Mechanics column's shell vocabulary. Do not weaken any of the three.
- If a claim cannot be verified against code, say so in the request rather than guessing.

## Dependencies

Depends on REQ-596, which corrected the other eleven findings from the same sweep in the same file.

## Open Questions

None.


---
id: REQ-597
status: claimed
domain: general
created_at: 2026-09-06T05:49:08Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: B
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
write_set: [skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md]
title: 'Correct ten stale claims across the rest of the prescribed-shell guide'
claimed_at: 2026-09-06T06:59:15Z
---

# Correct Ten Stale Claims Across the Rest of the Prescribed-Shell Guide

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Triage

**Route: B** — Explore then build.

**Reasoning:** The ten claims are named with file-and-line citations from a sweep that has already run,
so their existence is settled. What is not settled is their replacements. REQ-596 shipped a replacement
sentence that was false in its own right — it compressed a correct finding into an incorrect one — and
the review caught it. Every replacement here has to be re-derived from the code rather than transcribed
from the sweep's suggestion, and the sections these ten sit in have to be checked whole, because that is
the requirement that turned REQ-595 from one cell into three and REQ-596 from three claims into eleven.

**Planning:** Skipped. Three files, and the edit set is whatever the re-verification confirms.

**The caller drift is the newest part and the least examined.** `actions/commit.md` names a glob the
guide replaced one request ago and describes a `-` row the command never prints. Nobody has checked
`inspect.md` for the same two claims, or the other shipped callers for anything similar.

## Plan

**Planning not required** — Route B: three files, and the work is whatever the re-verification and the
section sweep confirm.

*Skipped by work action*

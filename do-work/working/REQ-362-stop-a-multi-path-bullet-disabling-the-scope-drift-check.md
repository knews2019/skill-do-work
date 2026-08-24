---
id: REQ-362
title: "[impact-rule-change] Stop a multi-path bullet disabling the scope-drift check"
status: claimed
claimed_at: 2026-08-24T18:57:05Z
status_changed_at: 2026-08-24T18:57:05Z
route: B
created_at: 2026-08-24T11:15:00Z
user_request: UR-068
addendum_to: REQ-344
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-24T18:57:05Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work/tools/checks/scope-drift.sh
---

# Stop a Multi-Path Bullet Disabling the Scope-Drift Check

## What

`scope-drift.sh` extracts only the **first** backticked path per Implementation Summary bullet, so a
bullet listing several files hides every path after the first. Hand-formatting a summary silently
disables the measurement that exists to make drift visible.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Found on REQ-344, which touched nine files against a declared write set of two and documented all
seven extensions honestly. Its Implementation Summary packed six paths onto one bullet, three of them
as bare filenames, and the checker reported **two** undeclared files rather than seven:

```
DRIFT: touched but never declared in ## Scope:
  skills/do-work-board/tools/queue-kanban/frontmatter.go
  skills/do-work/actions/capture-reference.md
```

That REQ's own prose asserts the tool "reports seven files touched but never declared". It does not,
and nobody noticed until review — the under-report reads exactly like a clean-ish result.

The same list is what `actions/work.md` Step 9 validates against
`git diff --name-only <pre>..<merge_hash>`, so the same formatting defeats that check too. A checker
whose coverage depends on how a human wrapped a list is a checker that will keep being quietly
disabled.

## Detailed Requirements

- Every backticked repo-relative path on an Implementation Summary bullet is extracted, not just the
  first.
- A bullet naming a file without a path (`capture.md` rather than
  `skills/do-work/actions/capture.md`) either resolves or is reported as unparseable — silence is the
  failure mode this REQ removes.
- Re-running against REQ-344's archived record reports all seven undeclared files.
- Existing single-path bullets behave exactly as before. The suite's current expectations must not
  move.

## Constraints

- `_dev/primes/prime-shell-commands.md` governs. Read it first.
- The Implementation Summary template (`actions/work-reference.md` § Step 6.25) prescribes one
  repo-relative path per bullet. Fixing the extractor does not license abandoning that convention —
  but the extractor must not depend on it, because it is the only thing standing between a
  hand-formatted list and a disabled check.
- Do not make the checker fail on a bullet that legitimately names no file (a "What was done"
  sentence). Distinguish "no paths here" from "paths I could not parse".

## Red-Green Proof

**RED prompt/case:** Run `skills/do-work/tools/checks/scope-drift.sh` against a REQ whose
Implementation Summary lists six files on one bullet. It reports two undeclared files, not seven.
Reproduced on REQ-344, 2026-08-24.

**GREEN when:** that same REQ reports all seven, a bare filename is either resolved or reported as
unparseable, and every existing single-path REQ produces the same output as before.

**Validation:** Inferred during REQ-344's review.

---
*Source: REQ-344 review finding F5 (UR-068).*

## Triage

**Route: B** — The defect and one-file target are clear, but exploration must choose a multi-backtick extractor, define bare-filename resolution versus loud rejection, and preserve every existing single-path fixture before implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

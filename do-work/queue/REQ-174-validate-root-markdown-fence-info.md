---
id: REQ-174
title: Validate root Markdown fence info
status: pending
created_at: 2026-08-11T17:00:04Z
user_request: UR-039
domain: testing
prime_files: [skills/do-work/tools/prime-do-work-update.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-172, REQ-173]
batch: accepted-p2-fixes
write_set: [_dev/tests/shipped-package-reference-contract.sh]
---

# Validate Root Markdown Fence Info

## What

Make root and list fence classification share the CommonMark rule that a backtick-fence info string cannot itself contain a backtick.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The current root fence branch accepts an invalid opener and masks a following link that the repository's pinned Goldmark renderer publishes.

## Detailed Requirements

- Apply one opener/info validity rule to root and list fence branches.
- Reject backticks in backtick-fence info strings.
- Preserve tilde-fence behavior.
- Add the reproduced root-level Goldmark-differential fixture.

## Constraints

- Consolidate the existing list validation rather than adding another independent exception.
- Preserve the classifier behavior earned by REQ-150 and REQ-163.

## Red-Green Proof

**RED prompt/case:** Classify `````lang`invalid\n[live](visible.md)\n````` and compare it with the pinned Goldmark renderer.
**Why RED now:** The shell classifier masks every line and returns no target while Goldmark renders `visible.md` as a link.
**GREEN when:** The classifier returns `visible.md`, agrees with Goldmark, and existing root/list/tilde fence fixtures still pass.
**Validation:** User confirmed

## Full Context

See `do-work/user-requests/UR-039/input.md` and the preceding validated-feedback report.

---
*Source: fix accepted*

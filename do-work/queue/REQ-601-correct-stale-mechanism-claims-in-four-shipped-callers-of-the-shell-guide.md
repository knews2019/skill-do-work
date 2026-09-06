---
id: REQ-601
status: pending
domain: general
created_at: 2026-09-06T07:22:52Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
depends_on: [REQ-597]
related: [REQ-595, REQ-596, REQ-597]
write_set: [skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work/actions/install.md, skills/do-work-toolbox/actions/present-work.md, skills/do-work-board/actions/board.md]
title: 'Correct stale mechanism claims in four shipped callers of the shell guide'
---

# Correct Stale Mechanism Claims in Four Shipped Callers of the Shell Guide

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

REQ-597's caller sweep checked every shipped file that points at
`skills/do-work/docs/prescribed-shell-primitives.md`. Four files outside that request's write set
restate something the guide owns in a way the code contradicts. Each was verified against the code by
running it.

## Detailed Requirements

- `skills/do-work-toolbox/actions/ai-report-reference.md:31` and `:37` say staging is "adjacent to
  the target" / "adjacent to `generated/`". Both commands stage under the system temporary directory
  (`os.MkdirTemp("", …)` at `report_image.go:61` and `:164`). The guide already says so.
- `skills/do-work/actions/install.md:335` names a stray `SKILL.md.download` as a Red Flag. That exact
  name is never created; `os.CreateTemp` appends a random suffix, so the real stray is
  `SKILL.md.download.<random>` and a reader looking for the literal name never finds it. Line 50 has
  the same name.
- `skills/do-work-toolbox/actions/present-work.md:136` says the helper "verifies each output against
  the source separately". No output is ever compared to the source: it is read once and both outputs
  are written from that buffer; the only check is on the existing canonical file before replacement.
- `skills/do-work-board/actions/board.md:87` says a non-git project "skips it silently". It exits zero
  with a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning finding, which the launcher prints.
- `skills/do-work/actions/install.md:246` says the installer uses "the canonical Local Git ignore
  helper". It never calls that command; it implements the same contract inline. Say that.
- Derive every replacement from the code. The sweep's drafts are input, not answers — four of the prior
  sweep's drafts for the guide were false.

## Constraints

- Prose only, no behaviour change, no code change. Fourteen other caller sites were checked and hold;
  do not rewrite them.

## Dependencies

Depends on REQ-597, which corrects the guide these callers point at.

## Open Questions

None.

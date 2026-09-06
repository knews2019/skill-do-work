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
write_set: [skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work/actions/install.md, skills/do-work-toolbox/actions/present-work.md, skills/do-work-board/actions/board.md, skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work-toolbox/actions/architecture-report.md, skills/do-work/actions/work-reference.md]
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
- **One sentence in the guide itself, written by REQ-596 and found wrong by REQ-599's builder.** The
  conflict-resolution bullet in "Protected inventory fallbacks" says "when the timestamps are equal or
  missing the first claim found stands". `requestmodel.ParseTimestamp` returns the zero time for a
  missing `completed_at` and the comparison is `completed.After(current.completed)`, so a `working/`
  REQ with no timestamp — which is every in-flight REQ — loses the path to any `archive/` REQ that has
  one, regardless of walk order. Only when both are missing does the first found stand. Say exactly
  that, and confirm it with a fixture holding one in-flight and one archived REQ claiming the same path.
- Derive every replacement from the code. The sweep's drafts are input, not answers — four of the prior
  sweep's drafts for the guide were false.
- **More instances of the same phantom-script class, found by REQ-597's caller builder after this
  request was captured** (evidence in `do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`):
  `present-work.md:136` ("the helper verifies each output against the source separately", the twin of
  the corrected portfolio sentence) and `:140` ("the compatibility script"); `ai-report-reference.md:31`,
  `:37`, `:47` ("a retained script" for report images — none exists); `architecture-report.md:46`;
  `install.md:50`, `:246`, `:261`, `:335`; `board.md:87`; and `work-reference.md:322` ("hands the child
  the console's own handles" — both child streams get the CLI's stderr). Two files join the write set for
  them. Line numbers are as of commit `d5cf28b`.
- **The tie-break sentence above, corrected again by REQ-599's review:** the first claim found stands
  whenever the two timestamps compare equal, present or missing (the comparison is a strict `After`),
  and `ParseTimestamp`'s error is discarded, so an unparseable `completed_at` counts as missing. Write:
  a REQ with a parseable `completed_at` beats one without, whichever root is read first; when both parse
  equal, or neither parses, the first claim found stands (`working/` before `archive/`, name order within
  a root). The fixture must cover equal-present, both-missing and unparseable.

## Constraints

- Prose only, no behaviour change, no code change. Fourteen other caller sites were checked and hold;
  do not rewrite them.

## Dependencies

Depends on REQ-597, which corrects the guide these callers point at.

## Open Questions

None.

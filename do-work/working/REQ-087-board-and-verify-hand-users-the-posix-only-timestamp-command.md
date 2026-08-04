---
id: REQ-087
title: The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing
status: claimed
claimed_at: 2026-08-04T00:04:10Z
created_at: 2026-08-03T21:45:43Z
user_request: UR-015
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: false
depends_on: [REQ-078]
maintenance: false
addendum_to: REQ-078
review_generated: true
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/web/board.js]
---

# The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing

## What

REQ-078 made `actions/work-reference.md`'s Timestamp rule the only place in `actions/` that spells a
command for obtaining a stamp, and gave that rule a Windows form that actually runs. Three sites in
`tools/queue-kanban/` still hand a user the bare POSIX command:

- `tools/queue-kanban/verify.go:287` — the `Remedy:` string on a future-dated-timestamp finding:
  "re-stamp it with `date -u +%Y-%m-%dT%H:%M:%SZ` (the Timestamp rule in actions/work-reference.md)".
- `tools/queue-kanban/web/board.js:154` — the claim-stopwatch tooltip.
- `tools/queue-kanban/web/board.js:553` — the future-stamp data-warning text.

A Windows user who follows any of the three gets a command that does not exist on their box — the
exact failure REQ-078 fixed one layer up.

## Why

These were left out of REQ-078 on purpose, not by oversight: they are a **different surface with a
different tradeoff**. An action file can say "see the Timestamp rule" because the agent reading it can
open the rule. A board tooltip cannot — the person reading it is looking at a web page, and replacing
a usable command with a file reference makes the tooltip worse. So the fix here is a judgment call
about what a UI should say, not a continuation of REQ-078's sweep.

Low severity: nothing is corrupted, and the finding these strings decorate is itself correct. The
cost is a Windows user copying a dead command out of a UI that looks authoritative.

## Detailed Requirements

1. **Decide what each surface should say, per surface** — the three sites do not have to match. A
   plausible split: `verify.go`'s remedy keeps a command but adds the Windows one (it is CLI output,
   read next to a shell); the two `board.js` strings drop to the shape ("the current UTC instant,
   `YYYY-MM-DDTHH:MM:SSZ`") and cite the rule, since a tooltip's job is to explain the badge, not to
   be a manual. Argue for whatever you pick.
2. **Do not paste the Windows one-liner into all three.** Three copies of a two-branch command in
   display strings is the inline-copy problem REQ-078 just removed, relocated.
3. **If a command survives anywhere, it must agree with the rule** — `ToUniversalTime`, the `\T`/`\Z`
   escapes, the `powershell -NoProfile -Command` wrapper for `cmd`.
4. **Consider whether a contract assertion is warranted** and say why if not. REQ-078's assertion is
   scoped to `actions/` deliberately; widening it to `tools/` would flag `timestamp.go`'s rationale
   comments and the test fixture, which are correct as they are.

## Constraints

- `tools/queue-kanban/` changes are folded into the skill's own version and changelog — no
  independent versioning (the tool's conventions are in `tools/queue-kanban/prime-do-kanban.md`).
- No behaviour change. These are display strings; the findings and badges they annotate stay exactly
  as they are.

## Dependencies

`depends_on: [REQ-078]` — the Windows form these strings would have to agree with ships there.

## Builder Guidance

**Certainty: Firm on the inventory, open on the wording.** All three sites were found by a three-shape
grep during REQ-078's review and are listed above with line numbers. What each should say is the
actual work.

## Full Context

Found by `actions/review-work.md`'s Restatement Sweep during REQ-078's review — see that REQ's
`## Review` → Findings.

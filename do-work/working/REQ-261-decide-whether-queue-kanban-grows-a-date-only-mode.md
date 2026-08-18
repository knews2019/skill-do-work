---
id: REQ-261
title: Delete the date-only tripwire and keep the rule
status: claimed
claimed_at: 2026-08-18T22:59:48Z
route: A
created_at: 2026-08-18T19:30:47Z
status_changed_at: 2026-08-18T21:01:24Z
user_request: UR-055
addendum_to: REQ-253
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work-reference.md
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route A
    - 1-file write set
    - 2 acceptance criteria
    - full-suite verification
---

# Delete the Date-Only Tripwire and Keep the Rule

## What

The Timestamp rule's date-only paragraph ends with "revisit if a second consumer appears". Remove that clause. Keep everything else in the sentence — the shell one-liners, and the reason there is no tool subcommand (adding one would widen the skill's single compiled-dependency exception for something the POSIX floor already covers).

The clause is the only part that does not survive its own argument: it keys on how many consumers exist, and consumer count does not bear on whether a shell one-liner suffices. Leaving it invites a re-litigation the surrounding sentence already settles — a list where a condition belongs (CLAUDE.md → State conditions, not lists).

## Requirements

- The "revisit if a second consumer appears" clause is gone; the rest of the date-only paragraph is unchanged in meaning.
- No date-only subcommand is added to the board tool.
- The paragraph still reads as one coherent sentence after the removal — check the ui-review consumer clause REQ-253 added still sits naturally beside it.
- `bash _dev/tests/maintainer-verify.sh` exits 0 (the Timestamp-rule citation contract counts 54 instant / 17 date-only sites today; a prose-only removal must not move them).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Discovered by REQ-253's builder ([low]; the tripwire sentence was left verbatim as a deliberate maintainer call). Note the board's pinned write-surface count is unaffected either way — `now`-style output is read-only.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-253: the date-only paragraph's "revisit if a second consumer appears" condition is now true. Should I process this as a new task? → Yes, and the answer is decided: delete the tripwire clause and keep the rule
  Recommended: Yes, add to queue (will flip to 'pending') — and the builder decides between the subcommand and a re-stated threshold.
  Also: No, discard it — two consumers on the shell one-liner is still fine.

**Answered [2026-08-18]:** User approved via `do-work clarify` **and settled the underlying question**, after asking where the sentence came from. Provenance established during clarify: the clause arrived with the repository import (recorded author "Claude", root commit `8d5c2ab`) from a pass that restructured the Timestamp rule — it is builder prose, not a maintainer decision. The user's reasoning, which is now the REQ's requirement: the rule itself is sound (`date -u +%F` works on the POSIX floor, and the board tool is the skill's only sanctioned compiled dependency, so putting a date behind a Go binary would widen that exception for nothing), but the tripwire keys on **consumer count**, which has no bearing on that argument — a shell one-liner is no worse at two callers than at one. Delete the clause; keep the rule. Do not add a date-only subcommand.

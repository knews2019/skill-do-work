---
id: REQ-339
title: "[impact-rule-change] Addendum: count every prescribed-shell case the runner reports"
status: pending
status_changed_at: 2026-08-23T22:32:23Z
created_at: 2026-08-23T19:30:00Z
user_request: UR-065
addendum_to: REQ-325
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - _dev/tests/prescribed-shell-harness.sh
---

# Count Every Prescribed-Shell Case the Runner Reports

## What

`prescribed_shell_finish` counts a case file's cases with `grep -cE '^# [a-z0-9][a-z0-9-]*: '`, so a
header with a space or a comma before its colon is invisible. `generate-report-image.sh` reports 7
cases and contains 9: `# generate-report-image caller contract: …` and
`# generate-report-image, interrupted directly: …` are both uncounted. The aggregate figure the
runner prints ("96 named script cases across 17 per-script files") inherits the same undercount.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-234 replaced a hand-maintained literal with a derived count precisely so the figure would stop
being a remembered number reported as a measured one. A regex that silently skips real cases is the
same untruth with a different cause, and it is the more durable kind: nothing fails, the number just
reads low forever. A reader using the count to judge coverage is misled in the direction that
matters.

## Detailed Requirements

- The count matches the number of case blocks the file actually contains.
- Decide and state what a case header *is* — the counting rule lives in a comment beside the count
  (REQ-234's contract), so whatever shape is accepted has to be written down there too.
- Sweep every file under `_dev/tests/prescribed-shell-cases/` for headers in the uncounted shapes
  and confirm the new count agrees with a hand count on at least the two files that change.
- The aggregate figure in `prescribed-shell-scripts-behavior.sh` moves in step (it greps the same
  pattern).

## Constraints

- `_dev/primes/prime-shell-commands.md` governs any shell that ships or gates. Read it first — in
  particular § *Closed Enumerations Go Stale*: prefer widening the condition over enumerating the
  two header spellings that happen to exist today.
- Renaming the offending headers to fit the existing regex is the cheaper fix and is on the table —
  but say why, because it leaves the next author free to write a header the count skips.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-325: the prescribed-shell case count
  skips any header with a space or comma before its colon, so `generate-report-image.sh` reports 7
  of its 9 cases and the runner's aggregate is low by the same amount. Should I process this as a
  new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - [2026-08-23] User approved via clarify: the reported case count must match what the file
    holds — a figure that reads low forever is the same untruth REQ-234 removed when it
    replaced the hand-maintained literal. Nothing put out of scope; renaming the two odd
    headers instead of widening the rule stays on the table as the REQ's Constraints say,
    provided the reason is stated.

## Red-Green Proof

**RED prompt/case:** `grep -cE '^# [a-z0-9][a-z0-9-]*: ' _dev/tests/prescribed-shell-cases/generate-report-image.sh`
prints 7 while `grep -c '^# generate-report-image'` on the same file prints 9.

**GREEN when:** the count reported by that file equals its hand-counted case blocks, and a case
header written in either previously-uncounted shape is counted.

**Validation:** Inferred during REQ-325's implementation — a Discovered Task, not a user request.

---
*Source: Discovered Task, REQ-325 (UR-065).*

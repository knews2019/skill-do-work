---
id: REQ-600
status: claimed
domain: general
created_at: 2026-09-06T06:53:35Z
user_request: UR-105
review_generated: true
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-06T07:26:22Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
maintenance: true
depends_on: [REQ-594]
related: [REQ-593, REQ-594]
title: 'Put the SIGPIPE trap in the prime shell authors read, and fix the one shipped block that carries it'
claimed_at: 2026-09-06T07:26:22Z
status_changed_at: 2026-09-06T12:38:32Z
---

# Put the SIGPIPE Trap Where Shell Authors Read It

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B. Prime read. Sweep every shipped fence first, then one prime section, one block, and the scan wired where the fences are already extracted.
- [x] **[APPLY]:** Five files, commits `22a1ea4` and `cb606f2`, merged as `a25c752`. The three `_dev/tests/` files are the widening the Constraints allowed and D4 records.
- [x] **[UNIFY]:** `git diff --stat 9e00a092..a25c7522`: 5 files, +155/-77. Fast gate on main exit 0; action-shell-blocks 74 fences and 33 shipped shell files; audit 95 tracked shell files, 19+7 shapes; audit-lockins passed. No debug artifacts.

## What

Two requests and about 140 fixed sites later, `_dev/primes/prime-shell-commands.md` still carries no
mention of the defect class. `CLAUDE.md` names that file as what to read before writing or reviewing
shell anywhere it ships, and calls it "the hard-won trap list". The trap is documented in two run records
and two request files, which is not where the next person writing shell will look.

Separately, one shipped markdown block carries the forbidden shape:
`skills/do-work-knowledge/actions/memory-reference.md:88` is inside a ```bash block that agents copy
and run — `ollama list 2>/dev/null | grep -qiE 'embed'`.

## Why

REQ-594's guard covers tracked `*.sh` files. Shipped guidance is where the shape gets copied from, and
a prime is where an author decides how to write the line in the first place. A guard that catches the
defect after it is written is worth less than a prime that stops it being written.

## Context

Both found by the independent three-lens review of REQ-594. The reviewer noted that the guard's
reader-set limitation IS disclosed — in a run record and a request file, not in the prime.

## Detailed Requirements

- Add the class to `_dev/primes/prime-shell-commands.md` as its own section, beside "Unchecked Exit
  Status Reads as Content", which it is a specific and nastier case of. State: the condition (a writer
  piped into an early-leaving reader under `pipefail`), that it is wrong in **both** directions, the
  measured window (silent below roughly 36 KB of producer output, certain above about 200 KB), the fix
  (capture and read as a herestring, asserting the producer's status separately when it can fail), and
  the readers the guard cannot see (`rg -q`, `head`, `sed -n '1p;q'`, `awk '/x/{exit}'`, `read`).
- Point at `_dev/tests/quiet-grep-pipeline-audit.sh` as the guard, and say plainly what it does not
  cover, so the prime is not read as "the guard has this".
- Fix `memory-reference.md:88`. Its producer output is far too small to reach the window in practice,
  so this is about what shipped guidance teaches, not a live failure — say so rather than overstating it.
- Check the other shipped action files for prescribed blocks carrying the same shape in the same pass.
  A markdown scan is not what REQ-594's guard does, and one instance found by hand says nothing about
  the rest.

## Constraints

- Prose and one prescribed block. No change to the guard, whose fixture pins 19 shapes.
- `_dev/tests/action-shell-blocks.sh` already checks prescribed shell blocks; if the fix belongs there
  as an assertion rather than only as prose, say so and do it, but do not weaken what it already pins.

## Open Questions

None.


---
id: REQ-594
status: claimed
domain: testing
created_at: 2026-09-06T03:08:29Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: false
depends_on: [REQ-593]
related: [REQ-593]
write_set: []
title: 'Generalize the SIGPIPE fix: about 130 quiet-grep pipelines remain across the maintainer test tree'
claimed_at: 2026-09-06T03:36:45Z
---

# Generalize the SIGPIPE Fix Across the Maintainer Test Tree

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

REQ-593 removed a defect where a writer piped into `grep -q` under `set -o pipefail` reported the
writer's SIGPIPE death as grep's verdict. It fixed eleven sites and added a source scanner — but the
scanner reads `${BASH_SOURCE[0]}`, its own file and nothing else. The class is closed in **one file
of twenty-three**.

`grep -rc -E '\|[[:space:]]*grep[[:space:]]+-[A-Za-z]*q' _dev/tests/` returns 23 non-empty files and
roughly 130 sites, including nineteen in the gate's own `_dev/tests/contracts/core-checks.sh` and
thirty-one in `_dev/tests/prescribed-shell-cases/qualify.sh`. `core-checks.sh` already disagrees with
itself: one site uses the herestring form, another still pipes.

The consequence is measured, not theoretical. The defect is silent below roughly 36 KB of producer
output and certain above about 200 KB, and it is wrong in **both** directions — a positive matcher
misses a pattern that is present, and a negative matcher fails to flag one it should. Two scripts have
already produced false failures in the gate during a single work run.

## Requirements

- No check anywhere under `_dev/tests/` decides on a quiet grep fed by a pipeline whose writer can die.
- The guard is repository-wide rather than per-file, and it is not defeatable by ordinary shell:
  REQ-593's scanner was evaded by five spellings a reviewer found in minutes — the pipe at end of line
  with the grep on the next, `grep --quiet`, `grep --silent`, `| LC_ALL=C grep -q`, and
  `| command grep -q`.
- Each conversion preserves what the assertion measured. REQ-593 found that capturing a producer's
  output and discarding its exit status can silently narrow an assertion — a truncated archive whose
  partial listing still contains the marker passed where the pipeline form failed.
- A probe fixture carries every evasion spelling, so the widened guard is itself pinned.

## Context

Found during the independent review of REQ-593, which fixed the two named matchers and their nine
siblings in `_dev/tests/update-script-behavior.sh`. This request is the generalization that request's
own Scope deferred, promoted from a note to a request because the review showed the deferral had no
destination.

## Full Context

See the REQ-593 review in `do-work/archive/` (or `do-work/working/` while it is in flight) for the
five evasion spellings with their reproductions, and the site census by file.

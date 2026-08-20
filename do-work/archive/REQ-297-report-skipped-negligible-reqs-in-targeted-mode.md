---
id: REQ-297
title: "Report skipped-negligible REQs in targeted mode too"
status: completed
completed_at: 2026-08-20T13:28:52Z
commit:
created_at: 2026-08-19T16:35:00Z
user_request: UR-060
addendum_to: REQ-290
domain: general
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: false
related: [REQ-290]
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
---

# Report Skipped-Negligible REQs in Targeted Mode Too

## What

`--skip-impact-negligible` has two reporting paths, and a targeted run reaches neither properly.
Close the gap, and define what the skipped count counts.

## Why

REQ-290 added the flag with an explicit requirement: "A run that silently drops REQs reads as 'the
queue is empty' when it is not." Two reporting paths were built for it:

1. A `(K skipped as impact-negligible)` suffix on the queue status summary line — for runs that do
   find work.
2. A seventh Composed Exit Summary section, with its own headline — for runs that find nothing.

Targeted mode uses neither cleanly. It "stop[s] after the last one completes (skip the loop-or-exit
logic in Step 10)" and has its own exit line — `UR-011: no runnable REQs (2 completed, 0 pending).`
— so the per-REQ `REQ-NNN — [title] (impact-negligible)` list never renders and the new headline is
unreachable. `do-work run UR-060 --skip-impact-negligible` where every member is negligible produces
exactly the output the REQ set out to prevent.

Separately, **K's scope is undefined**: the whole queue, or only the resolved token set? In a
targeted run those differ, and a count over the whole queue would be actively misleading.

## Detailed Requirements

- A targeted run that drops REQs names them. Prefer reusing the existing skipped-as-negligible
  section over inventing a targeted-mode-only format — the composed summary's own rule says a
  section applies "if it has at least one REQ", so the cheapest correct fix may be letting the
  targeted exit path render that section rather than adding prose.
- Define K's scope in one clause where the suffix is specified: in targeted mode it counts the
  resolved token set, not the queue.
- The UR-expansion exit line (`UR-NNN: no runnable REQs (…)`) should account for members dropped by
  the flag rather than folding them into a count that reads as "nothing to do here".
- Keep it small. This is a reporting gap in text that already exists, not a new mechanism.

## Acceptance

- `do-work run UR-NNN --skip-impact-negligible` where every member is negligible names the flag and
  lists what it dropped; it never reads as an empty UR.
- The skipped count's scope is stated and matches the run mode.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Full Context

Important finding 5 from REQ-290's review. See `do-work/user-requests/UR-060/input.md`.

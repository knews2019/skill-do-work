---
id: REQ-627
title: 'Show free disk space on the Testing page'
status: pending
created_at: 2026-10-02T18:02:37Z
user_request: UR-132
domain: frontend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-626"]
batch: board-links-and-disk-space
write_set: ["skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/generate.go", "skills/do-work-board/tools/queue-kanban/serve.go", "skills/do-work-board/tools/queue-kanban/web/board-testing.js", "skills/do-work-board/tools/queue-kanban/web/template.html", "skills/do-work-board/tools/queue-kanban/disk_space_test.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/javascript_behavior_*_test.go"]
---

# Show Free Disk Space on the Testing Page

## What
The Testing page always shows one line with the repo root's free and total disk space, for example `disk: 50.2 GiB free of 177.5 GiB`, fed by the measurement REQ-625 (the low-disk-space verify probe) already takes, and coloured by the same thresholds (warning under 10 GiB, critical under 3 GiB).

## Why
REQ-625 shipped in 0.305.60, but a verify finding appears only below the thresholds, so on a healthy machine the free space is shown nowhere. The user wants to see the number before it becomes a finding, on the Testing page.

## Detailed Requirements
- Carry the healthy measurement into the board payload: a small field on `generatedBoardData` (free bytes, total bytes, and the measured directory reduced through `reduceAbsolutePaths`, or a skip reason when the platform cannot measure). Reuse `diskSpaceMeasurer` and `formatGibibytes` from verify.go; do not measure twice if one call can feed both the probe and the payload.
- `serve` recomputes the value on every request, outside the mtime cache, the same way the verify findings are (free space changes while no file changes).
- The static snapshot carries the value measured at generation time, labelled as such (for example "at generation").
- The Testing page renders one always-visible line near the top: free of total in GiB with one decimal, coloured neutral above 10 GiB, warning below 10 GiB, critical below 3 GiB. When the measurement was skipped, the line says so ("disk: not measured on <GOOS>") rather than disappearing.
- No absolute path reaches the page. The existing payload test `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` or a sibling covers the new field.
- The probe's finding behaviour from REQ-625 is unchanged; this REQ adds a readout, not a second threshold.

## Constraints
No new dependency. One syscall per request. Nothing is deleted. Shipped files change, so this is a release. The field is display-only; nothing schedules or gates on it.

## Dependencies
None. Builds on REQ-625, which is archived (commit add2c162); independent of REQ-626 (Give every board page and lens its own URL).

## Builder Guidance
High certainty on placement (Testing page, user decision) and thresholds. Builder latitude on the payload field name and the exact markup. A Go test with a fake measurer plus a JavaScript behaviour test in the existing harness are the expected proof.

## Red-Green Proof
**RED prompt/case:** Open the Testing page on a machine with 50 GiB free: no free-space figure appears anywhere; the generated board payload has no disk-space field.
**Why RED now:** REQ-625 keeps only findings below threshold; the healthy measurement is discarded.
**GREEN when:** The Testing page shows `disk: <free> GiB free of <total> GiB` on a healthy machine (neutral colour); with a fake measurer at 9 GiB the same line turns warning and at 2 GiB critical; with the unsupported sentinel it reads "not measured"; the payload carries the field with no absolute path.
**Validation:** User adjusted (page chosen at capture)

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (7095 tokens, over the 2000 budget; `slugged: partial`). Matched: changing queue-kanban model and UI; the `disk-space-blind-spot` entry is the direct predecessor.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, over budget; `slugged: partial`). Matched: static output.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-132/input.md` for complete verbatim input.

*Source: in one of the pages list the free space too since it's already collected*

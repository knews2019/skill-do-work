---
id: REQ-204
title: Harden ai-report generated-batch lifecycle
status: pending
status_changed_at: 2026-08-17T18:10:49Z
domain: general
created_at: 2026-08-15T19:39:11Z
user_request: UR-042
addendum_to: REQ-198
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: ai-report-generated-batch-lifecycle
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
---

# Review Fix: Harden AI-Report Generated-Batch Lifecycle

## What

Close the complete generated-image batch lifecycle: an interrupted caller must terminate and reap exactly its own helpers, and publication must fail closed if the destination appears at the final boundary. Solve both instances as one rule because both require ownership of the batch from private staging through terminal publication or cleanup.

## Context

REQ-198 fixed the original all-failed directory shape, but review showed that signal cleanup handles files without process ownership and that a check-then-plain-`mv` can nest staging inside a newly appeared destination while returning success.

## Instances

- [ ] Batch interruption: signal and reap recorded helper PIDs (and their owned descendants) before removing exact staging; no optional full-host backend may outlive the caller.
- [ ] Final publication: coordinate destination appearance after the last check and prove the operation returns nonzero, preserves the colliding directory, and leaves no nested/private stage.

## Requirements

- Preserve normal wait-all and per-status freshness behavior.
- On HUP/INT/TERM, terminate and reap exactly the current batch's recorded process tree before staging cleanup.
- Use a portable exclusive/atomic directory publication boundary or a verified rollback that cannot report success after nesting.
- Never delete or overwrite a colliding destination.
- Add exact prescribed-block behavior replays for both adversarial paths.

## Red-Green Proof

**RED prompt/case:** Signal only the batch shell while slow helpers run, then coordinate creation of `generated/` after the final absence check but before publication. Current behavior can leave helpers alive and can return success with staging nested under the colliding directory.
**Why RED now:** File cleanup alone does not own the process tree, and plain `mv` treats an existing destination directory as a container.
**GREEN when:** The signal replay proves no owned process survives and no staging/public path leaks; the collision replay returns nonzero, preserves the destination byte-for-byte, creates no nested stage, and normal all-failed/mixed paths still pass.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] The original empty-directory bug is fixed, but review found two deeper lifecycle edges in the same shell batch. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.

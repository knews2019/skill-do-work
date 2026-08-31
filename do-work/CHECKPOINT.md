# Session Checkpoint

## Completed

- REQ-408 — shared repository-model foundation committed as `ac2e3acd` (metadata `e6488553`); follow-ups REQ-428 and REQ-429 queued.
- REQ-409 — canonical cleanup implementation archived `completed-with-issues`; follow-ups REQ-430 through REQ-433 queued after failed remediation re-review.
- REQ-410 — canonical doctor/forensics implementation archived `completed-with-issues`; follow-ups REQ-434 and REQ-435 queued after failed remediation re-review.

## In Progress (interrupted)

- REQ-426: Preserve setuid, setgid and sticky bits on managed files instead of stripping them — claimed 2026-08-31T10:30:10Z — writer: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2

## Still Queued

- REQ-411 through REQ-420 remain in the ordered UR-081 chain.
- REQ-426 and REQ-428 through REQ-435 are pending; REQ-427 is pending a corrected Go-floor decision after exact Go 1.23 verification failed.

## Session Notes

- REQ-411 was released at a clean pre-plan boundary for the fresh session; no implementation files, branches, or worktrees exist for it.
- REQ-411 now waits for the two repository-model review fixes and the doctor/forensics delegation fix it consumes.
- REQ-430 through REQ-433 are serialized by explicit dependencies; REQ-412 also waits for REQ-433 so state transactions build on the corrected cleanup/archive semantics.
- REQ-427's prior Go 1.23.0 answer rested on a newer-toolchain test and is invalid. Exact Go 1.25 passes; exact Go 1.23 and 1.24 fail on the rooted-filesystem API boundary.

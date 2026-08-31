# Session Checkpoint

## Completed

- REQ-408 — shared repository-model foundation committed as `ac2e3acd` (metadata `e6488553`); follow-ups REQ-428 and REQ-429 queued.
- REQ-409 — canonical cleanup implementation archived `completed-with-issues`; follow-ups REQ-430 through REQ-433 queued after failed remediation re-review.
- REQ-410 — canonical doctor/forensics implementation archived `completed-with-issues`; follow-ups REQ-434 and REQ-435 queued after failed remediation re-review.
- REQ-426 — special mode bits preserved through managed replacement and real install; merge `73bd4a6f`, release `7e0112bc`, metadata `d090af93`.
- REQ-429 — complete normalized schema-field projection merged as `67942dd9`; independent review approved at 100% with no findings.
- REQ-430 — UR closure now depends on terminal member archival; merged as `5f3531d0`, review accepted the core behavior and routed one recovery-command edge to consent-gated REQ-445.
- REQ-434 — unsupported timestamp shapes no longer anchor doctor ordering repairs; merged as `509cbee4`, independently approved at 98% with no Important findings.
- REQ-431 — documentation rewrites now follow their owning moves and compose from current bytes; merged as `3d695fcb`, independently approved at 100% with no findings.
- REQ-432 — consumed scratch can no longer bypass the commit-mode empty-index guard; merged as `ad6c252e`, with adjacent doctor remediation routed to consent-gated REQ-446.
- REQ-435 — doctor-forensics delegation now has complete typed report projection, stable recovery references, and deterministic board-warning mapping; merged as `c1536cbf`, independently approved at 98%.
- REQ-433 — misplaced archived UR items now have independent conflict domains; merged as `f14803a8`, independently approved at 98%.

## In Progress (interrupted)

- REQ-411: Implement dependency-aware queue selection and actionable summaries — claimed 2026-08-31T17:49:48Z — writer: t2s-Virtual-Machine.local:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2

## Still Queued

- REQ-411 through REQ-420 remain in the ordered UR-081 chain.
- REQ-436 through REQ-444 remain queued; REQ-436's default audit decision is recorded and pending.
- REQ-445 awaits user consent through `do-work clarify` before it can enter the runnable queue.
- REQ-446 awaits user consent through `do-work clarify` before it can enter the runnable queue.

## Session Notes

- REQ-411 was released at a clean pre-plan boundary for the fresh session; no implementation files, branches, or worktrees exist for it.
- REQ-411 is ready now that the repository-model and doctor/forensics fixes it consumes are complete.
- REQ-432 and REQ-433 remain serialized by explicit dependencies; REQ-412 also waits for REQ-433 so state transactions build on the corrected cleanup/archive semantics.
- REQ-427 resolved the compatibility floor at exact Go 1.25.0 after exact Go 1.23 and 1.24 failed on the rooted-filesystem API boundary.

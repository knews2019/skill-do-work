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
- REQ-411 — dependency-aware queue selection and actionable summaries merged as `6209227b`, independently approved at 98% after one remediation.
- REQ-436 — atomic replacement and cleanup moves now preserve complete special modes; merged as `f0715c41`, independently approved at 98%.
- REQ-412 — canonical request lifecycle transactions merged as `7fc958be`; independently accepted at 83% after one remediation, with residual findings routed to REQ-447 and REQ-459.

## In Progress


## Still Queued

- REQ-413 through REQ-420 remain in the ordered UR-081 chain.
- REQ-437 through REQ-444 remain queued.
- REQ-447 extends the complete-mode publication audit to the separate queue-kanban module.
- REQ-459 repairs release staging for command-owned calibration changes.
- REQ-445 awaits user consent through `do-work clarify` before it can enter the runnable queue.
- REQ-446 awaits user consent through `do-work clarify` before it can enter the runnable queue.

## Session Notes

- REQ-411 now provides the canonical typed selector used to compute subsequent waves.
- REQ-412 established canonical request-state transactions; downstream UR-081 work can now consume them instead of duplicating lifecycle writes.
- REQ-427 resolved the compatibility floor at exact Go 1.25.0 after exact Go 1.23 and 1.24 failed on the rooted-filesystem API boundary.

## In Progress (interrupted)

  Last known state: Initial builder `a7c975c5` merged as `94560fde`; initial review found nine Important issues. The sole remediation `a43b2587` merged as `82534d36`, passed the full focused/race/full/vet/Go 1.25/Windows/differential/contract/install/update/canonical stack, and has not yet received its required fresh re-review.
  Key files being modified: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`, `internal/gittransaction/git_transaction.go`, `internal/gittransaction/git_transaction_test.go`, command/result registration, and the CLI prime (documented 32-path ceiling).
  Known issues: no post-remediation residual is recorded yet; the fresh review agents were interrupted for handoff before writing `REQ-418-rereview.md`. Resume with a new independent reviewer, then complete/archive/release or route every residual Important finding because the remediation allowance is exhausted.


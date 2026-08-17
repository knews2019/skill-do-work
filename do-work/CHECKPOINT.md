---
session_ended: 2026-08-17T00:45:00Z
last_completed: REQ-209
queue_state: 0 pending, 4 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 3
session_depth: moderate
---

# Session Checkpoint

## Completed This Session

- REQ-208: Deterministic P50 estimator script, reference file, and schema (Route B, 0.194.0; review 96% Pass)
- REQ-210: Recalculate P50 estimates in verify-requests after material repairs (Route A, 0.195.1; review 94% Pass)
- REQ-209: Wire P50 estimation into the work action (Route B, 0.195.0 + 0.195.2 remediation; review 95% Pass)

## In Progress (interrupted)

- REQ-213 — Board surfaces negative claimed→completed duration as a completion anomaly — claimed 2026-08-17T08:15:06Z — writer: vm:/home/user/skill-do-work


- REQ-211 — Calibrate estimator scoring table to archive actuals — claimed 2026-08-17T08:07:47Z — writer: vm:/home/user/skill-do-work

## Still Queued

- REQ-203: Harden presentation target-ID source-seam tests (pending-answers)
- REQ-204: Harden ai-report generated-batch lifecycle (pending-answers)
- REQ-205: Make portfolio publication independent and exact (pending-answers)
- REQ-206: Finish active publication delegation (pending-answers)

## Session Notes

- UR-047 (P50 active-duration estimation) captured, verified, built, reviewed, and archived end-to-end in one session; the feature was dogfooded on its own build (capture-time estimates, a verify-triggered recalculation on REQ-208, estimate lines printed at each claim).
- The estimation subsystem's durable map: `tools/estimate-p50.sh` (deterministic arithmetic + critical-path graph mode), `actions/estimate-reference.md` (lazy-loaded extraction guide), the optional `estimate:` schema block (frozen once execution begins), work.md Step 3.6 (primary wire point), verify-requests Step 7 item 4 (recalc after material repair). Both wirings are contract-pinned in `hardened_check_scripts`.
- `maintainer-verify.sh` in this container: 41 pre-existing environment FAILs (missing `just`, tar/gzip exec, stat-mode probes). The regression gate used all session: FAIL-set diff against the recorded baseline log. Full exit-0 verification needs a machine with those tools.
- ShellCheck 0.11.0 was installed to /usr/local/bin; Go pinned per-invocation via `GOTOOLCHAIN=go1.26.1`.
- AI report published at `ai-reports/2026-08-17_0034_UR-047-p50-active-duration-estimation/` (render-judged light + dark).
- UR-042's REQ-203–206 remain `pending-answers` by design (generation-2 review follow-ups awaiting consent). Run `do-work clarify` to answer them as a batch.

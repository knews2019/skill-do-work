---
id: UR-005
title: Concurrent-orchestrator lock guard for the work action
created_at: 2026-07-01T22:42:23Z
requests: [REQ-018]
word_count: 112
---

# Concurrent-orchestrator lock guard for the work action

## Summary

The user directed the session to "read and run do-work/HANDOFF.md" (2026-07-02). The handoff's continuation plan includes capturing the structural fix for the 2026-07-01 multi-orchestrator collision: a lock guard in the work action so a second `do-work run` on the same tree detects the active orchestrator instead of silently processing the same queue. Running the handoff resolved its open DECISION §1 in favor of capturing.

## Capture-Time Clarifications (user-selected, 2026-07-02)

1. **On detecting an active lock:** Ask to override — prompt the user (proceed anyway / take over the lock / abort). Chosen with the stated caveat that unattended runs need a non-interactive fallback to refuse-and-report.
2. **Stale locks (crashed session):** Heartbeat + age threshold — lock carries a timestamp refreshed at step boundaries; older than a generous threshold ⇒ presumed stale, warn and take over.
3. **RED/GREEN proof:** Confirmed as captured in REQ-018's Red-Green Proof section.

## Full Verbatim Input

User invocation (this session): "ok, read and run do-work/HANDOFF.md"

From `do-work/HANDOFF.md` § Suggested next prompts: "do-work capture-request: add a concurrent-orchestrator lock guard to the work action — the structural fix for the collision hazard"

From `do-work/HANDOFF.md` § DECISIONS pending (user), item 1: "**Queue ownership / multi-orchestrator hazard.** ≥2 Claude sessions ran do-work pipelines on this working tree simultaneously. This time it merged cleanly (the other session committed our builder's work correctly but archived REQ-017 with a hollow paper trail, since repaired). Recommend: ONE designated session per repo for `do-work run`. Optional structural fix worth capturing as a REQ: a lock-file guard in `actions/work.md` Step 2 (claim detects a concurrent orchestrator). User hasn't decided yet."

---
*Captured: 2026-07-01T22:42:23Z*

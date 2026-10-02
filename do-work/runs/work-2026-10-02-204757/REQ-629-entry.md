## 0.305.64 — Resuming a Handoff Keeps Its Claimed Requests (2026-10-02)

A session that resumed a handoff ran `recover`, followed the next step it suggested, `recover --take-over`, and lost its claimed requests: the takeover put each one back in the queue as pending and removed its written sections, even for work that was already merged. The suggested step is now safe.

- When `recover` finds a claimed request, its suggested next step is `advance REQ-NNN`, which changes nothing and names the request's next phase. The finding still says that `recover --take-over REQ-NNN` resets the claim, and that it is only for a claim no live session owns.
- A handoff now writes one `advance REQ-NNN` line per claimed request into its paste block and warns that a takeover resets a claim; the normal resume command never picks up claimed requests on its own.
- The work guide and the crash-recovery reference say the same.


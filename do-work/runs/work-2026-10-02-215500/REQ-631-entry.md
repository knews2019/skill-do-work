## 0.305.66 — Recover's Takeover Message Says Where a Reset Claim Goes (2026-10-02)

Since 0.305.64, `recover` warns that `recover --take-over` resets a claimed request. The warning said the request goes back "as pending", but a request with an unanswered question goes back as `pending-answers`, and a blocked one stays `blocked`.

- The warning now says the takeover returns the request to the queue and strips its routing and orchestrator sections, without naming a single status. Nothing else changes.


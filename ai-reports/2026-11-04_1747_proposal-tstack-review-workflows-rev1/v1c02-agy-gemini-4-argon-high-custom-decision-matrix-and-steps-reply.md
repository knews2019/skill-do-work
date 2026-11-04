Confirmed staging output: [`index-v1c02-agy-gemini-4-argon-high-custom-decision-matrix-and-steps.html`](index-v1c02-agy-gemini-4-argon-high-custom-decision-matrix-and-steps.html)

### Clarity Changes
- **Consolidated Options O1, O2, and O3 into a filterable side-by-side decision matrix (`#rev-1`):** Replaced separate option cards and scattered workflow notes with a seven-row comparison table covering activation scope, reviewer identity, finding ledger and diff binding, round caps, failure/dispute handling, cost, and queue status, with interactive row category filters.
- **Translated the O2 recommendation and `REQ-710` into five ordered implementation slices (`#rev-2`):** Paired an inline SVG control-flow diagram (contrasting the untouched default `do-work run` Step 7 path with the opt-in `AI-Evaluation-Loop` path) with five sequential build steps mapped to local files (`work.md`, `review-work.md`, `validate-feedback.md`) and `REQ-710` requirements.
- **Grouped open lifecycle decisions and the Jester advisory contract (`#rev-3`):** Pulled `Q1` (unresolved handoff vs. holding only the opted-in request), `Q2` (harness provider resolution), `Q3` (PR delivery exclusion), and the `UR-163` / `REQ-710` Jester advisory rules into one scannable reference grid before the folded appendix.
- **Folded upstream `tstack` details, secondary opportunities, corrections, glossary, and source map (`#rev-4`):** Placed deep reference material behind five expandable `<details>` sections with an "Expand/Fold all" control so the primary decision and implementation plan stay compact.

### Corrections to the Original Report
- **`skills/do-work/actions/work.md` line 350 vs. lines 354 and 366–379:** The original report cites `work.md:350` for Step 7 review dispatch, remediation, and `completed-with-issues`. In the current repository, line 350 is the `Already-green repository-gate repair exception` in Step 6.5; `### Step 7: Review` starts at line 354 and the single-remediation / `completed-with-issues` routing lives at lines 366–379.
- **`skills/do-work/actions/work.md` line 519 vs. lines 520–521:** The original report cites `work.md:519` for skipping a failed reviewer and continuing toward archive. In the current Error Handling table, line 519 covers repeated focused-test or repository-gate failures, line 520 covers `Review: Acceptance = Fail`, and line 521 is `Review work agent fails | Skip review, note it in the REQ file, continue to archive — review failure is not a gate`.
- **`skills/do-work/actions/work.md` line 462 vs. line 465:** Opportunity 2 cites `work.md:462` for the commit/delivery boundary; line 462 is Step 8 item 3 (`Terminal judgment`), while `### Step 9: Commit Phase (Git repos only)` starts at line 465.
- **Unverified external claims:** Upstream `t3dotgg/tstack` sources at commit `1085050f3baed1ac724caa4a8166145fb2ac88af` were not re-fetched over the network (external services are prohibited by the task boundary), and runtime cost, latency, and convergence remain unmeasured.

### Open Questions
- **Q1 (Unresolved handoff state):** When an opted-in `AI-Evaluation-Loop` hits its round limit, encounters a persistent dispute, or lacks a required reviewer, should `do-work` hold the claimed request in place while continuing unrelated queue work, or transition it to a dedicated unresolved status?
- **Q2 (Initial provider-selection harnesses):** Which execution harnesses must support multi-provider council resolution in the first slice, and how should partially available provider sets be reported before launch?
- **Q3 (Exact numerical round cap):** `REQ-710` requires a maximum round limit and stopping conditions but leaves the exact numerical default (the proposal suggests 3 rounds, extendable to 5 for consequential open issues) to the builder.

Variation slug: decision matrix and steps

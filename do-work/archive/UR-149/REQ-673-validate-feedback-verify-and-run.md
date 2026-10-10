---
id: REQ-673
title: 'validate-feedback --capture runs verify-requests on the new UR, prints one combined report, and with --run continues into do-work run UR-NNN'
status: cancelled
created_at: 2026-10-09T21:26:55Z
user_request: UR-149
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-672]
related: [REQ-670, REQ-671, REQ-672, REQ-674]
batch: validate-feedback-capture
completed_at: 2026-10-10T12:59:17Z
---
# validate-feedback --capture Verifies the New UR, Then Optionally Runs It
## What
After REQ-672's capture, `validate-feedback --capture` runs `verify-requests` on the new `UR-NNN` without a second prompt and prints one combined report: the triage summary table, the UR and its REQs with titles, and the verify verdict. With `--run`, and only when verification finds no gaps, it continues into `do-work run UR-NNN`. With gaps, it stops before the run and prints the gaps with the exact command to run after fixing them.
## Why
Report item C4 (UR-149 input): `do-work verify-requests` was typed right after a capture roughly 12 to 15 times in 30 days, and "capture and run" / "capture them and run them" was typed in three sessions on 2026-10-08. `next-steps.md:27` and `capture.md:261` only suggest the verify step.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/next-steps.md:27`: after capture, suggest "`do-work verify-requests` before `do-work run`". `skills/do-work/actions/capture.md:261` makes the same suggestion for detailed requests.
- `skills/do-work/actions/verify-requests.md:32` and `:38`: capture QA accepts one `UR-NNN`, so no change in verify-requests is needed.
- `skills/do-work/actions/verify-requests.md` Step 7 (Offer Fixes, around `:155`) asks the user whether to apply fixes after the report.
- `skills/do-work/actions/work.md:103` (the report cites `:105` at 0.305.84): a `UR-NNN` token scopes the run to that UR's REQs and keeps `depends_on` gating; an explicit `REQ-NNN` list would bypass it.
- `skills/do-work/SKILL.md:20`: "Stop after capture unless the same user invocation explicitly requested execution too." `--run` is that explicit request; this REQ relies on that rule and adds no new one.
## Detailed Requirements
1. Add a new Step 8, verify and optional run, that runs only with `--capture`, after REQ-672 captured a UR.
2. Run `skills/do-work/actions/verify-requests.md` with the new `UR-NNN`, with no second user prompt.
3. Print one combined report: the triage summary table, the UR and its REQs with titles, and the verify verdict.
4. Without `--run`: stop after the combined report. No run starts.
5. With `--run` and a verify result with no gaps: continue to `do-work run UR-NNN` (the UR token, so `depends_on` still gates).
6. With `--run` and verify gaps: stop before the run and print the gaps with the exact command to run after fixing them.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No change to verify-requests, to capture's templates, or to how `do-work run` selects work. No new REQ fields, no new statuses.
- Capture ≠ Execute: a run starts only with `--run` (or the "capture and run" phrase) and a clean verify.
- Out of scope: a `--review N` front end; capturing without the flag or phrase.
## Assumptions (recorded at capture, no questions asked)
- In this chain, verify-requests runs through its report and does not enter Step 7 (Offer Fixes) on its own; the gaps are printed and the chain stops, which matches "stop and print the gaps". The builder records in Decisions if verify-requests' contract requires otherwise.
- "No gaps" means the verify report lists no Important, Minor or Ambiguous gap. A score below a threshold with no listed gap counts as no gaps; the gap list decides, not the score.
- The exact command after fixing gaps is `do-work verify-requests UR-NNN` followed by `do-work run UR-NNN`, with the real UR id filled in.
- When REQ-672 captured nothing (empty set) or every finding folded into existing REQs (UR with an empty `requests:` array), verify still runs on the UR's folds; with `--run`, the run is skipped when the UR owns no REQs, and the report says so.
## Dependencies
Depends on REQ-672 (the UR this step verifies and runs).
## Builder Guidance
Certainty is high on behaviour; the report fixes the steps and the stop conditions. Latitude: the combined report's layout.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox validate-feedback --capture` on a paste with 2 Accept, 1 Discuss, 1 Push back items, Discuss answered Accept; then the same with `--capture --run`, once with a clean verify and once with a verify gap.
**Why RED now:** the action ends at the handoff (`validate-feedback.md:113-117`); verify and run are typed by hand.
**GREEN when:** with `--capture`, `verify-requests` runs on the new UR without a second prompt, one combined report is printed, and no run starts. With `--capture --run` and a clean verify, the chain ends in `do-work run UR-NNN`. With a verify gap, it stops before the run and prints the gaps and the exact command.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers status contracts and downstream readers; this REQ chains three actions.
## Full Context
See `do-work/user-requests/UR-149/input.md` for complete verbatim input (sections Request C4, What happened, Where the behaviour lives today, Proposed direction C4, Acceptance check). No queued candidate shares this root cause.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-validate-feedback-capture.md`, Request item C4: "Run `verify-requests` on that UR automatically and print one combined report. With `--run`, and only when verification finds no gaps, continue into `do-work run` on just that UR."*

## Cancelled

- **When:** 2026-10-10T12:59:17Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`

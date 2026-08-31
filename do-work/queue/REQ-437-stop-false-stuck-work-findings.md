---
id: REQ-437
title: 'Stop false stuck-work findings for active and terminal REQs'
status: pending
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
---

# Stop False Stuck-Work Findings for Active and Terminal REQs

## What

Make doctor reserve `STUCK-WORK` for nonterminal working REQs that lack recent activity. A terminal file stranded under `working/` must receive the terminal-location finding without a contradictory ownership warning, and a recently edited old claim must not be called abandoned.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Doctor reports actively edited and terminal working REQs as “stuck”.`
- **Evidence:** `doctor_scan.go:101-103` invokes stuck detection for every working record; `doctor_scan.go:202-227` considers only `claimed_at`. `forensics.md` says a recently modified file is likely active.
- **Origin / earned by:** Commit `210d1459` centralized a longstanding contradiction between the timestamp-only detector and the action's recent-activity rule.
- **Surface-cost:** Terminal filtering is N/A. Modification-time evidence is Earned by the old-claim/recent-edit replay; one snapshot field and focused regressions are cheaper than recurring false ownership warnings.

## Detailed Requirements

- Exclude terminal statuses from `STUCK-WORK` while retaining the appropriate stranded-terminal finding.
- Carry trustworthy file modification-time evidence through the repository snapshot and use the action's recent-activity boundary consistently.
- Preserve existing age severity behavior for genuinely inactive nonterminal claims.
- Keep diagnosis deterministic and do not introduce timestamp recovery or automatic ownership changes.

## Constraints

- This root cause is distinct from REQ-435 (Complete the Doctor-Forensics Delegation Contract); coordinate overlapping files without merging the intents.
- The activity check must be one evidence field and predicate, not a general liveness subsystem.

## Red-Green Proof

**RED prompt/case:** Scan one terminal REQ stranded under `working/` and one nonterminal REQ claimed more than a day ago but modified moments ago.
**Why RED now:** Doctor evaluates every working file and its stuck predicate sees only `claimed_at`.
**GREEN when:** The terminal fixture receives no `STUCK-WORK`, the recently edited fixture receives no `STUCK-WORK`, a genuinely old inactive claim still does, and the terminal-location finding remains intact.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. Use the smallest snapshot evidence and predicate change that closes both named false-positive replays.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 6 from the validated external feedback.*

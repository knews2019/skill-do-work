---
id: REQ-641
title: '[impact-rule-change] Wave-end consistency check: the last review in a fan-out wave also sweeps for elements earlier members redefined'
status: claimed
created_at: 2026-10-07T19:12:33Z
user_request: UR-138
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-639, REQ-640]
batch: coordination-lessons
depends_on: [REQ-640]
write_set: [skills/do-work/actions/review-work.md, skills/do-work/actions/work-reference.md]
claimed_at: 2026-10-07T20:07:41Z
---
# Wave-End Consistency Check
## What
Catch the one fan-out defect no per-REQ review can see: a sibling merged later restates a contract an earlier sibling redefined.
## Why
An earlier member's review ran before the later member merged, so only the last review sees both. The lesson comes from the UR-137 run on 2026-10-07 (fixing accepted consumer-review findings) and Matt Maher's video "The New Way to Work With AI" (2026-10-07).
## Detailed Requirements
1. `skills/do-work/actions/review-work.md` Step 6 Restatement Sweep:
   - Record the result in the appended `## Review` block as one line: `**Restatement sweep:** redefined <elements> | nothing redefined`.
   - Add the line to the "Append to REQ File" template, and point the existing Verification Checklist item at it.
   - Go checks section presence only, so the extra line is safe.
   - New condition in the sweep's trigger (step 1): when the REQ under review is the wave's last successful integration (wave membership from the run manifest; a set-aside member does not skip the check), the trigger set also includes every element the wave's earlier members recorded as redefined, read from their archived `## Review` lines.
   - One-sentence reason in the text: an earlier member's review ran before the later member merged, so only the last review sees both.
   - Findings route exactly as today (impact token, `→ report only` unless `impact-critical`).
2. `skills/do-work/actions/work-reference.md` Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick": one sentence pointing at the wave-end sweep as the semantic half the merge cannot provide.
## Constraints
- Prose-only changes to shipped action files. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- No new test: no existing test exercises review prose.
## Dependencies
Depends on REQ-640 (mid-run messages), which depends on REQ-639 (delegated integration, the coordinator shape). The serial order was requested because all three touch `work.md`, `review-work.md` and the Fan-Out Dispatch section of `work-reference.md`.
## Builder Guidance
Certainty is high: the source names the line format, the trigger condition, the reason sentence and the one work-reference sentence. Latitude is limited to wording. Finding routing does not change.
## Red-Green Proof
**RED prompt/case:** In a two-REQ wave, REQ-1 redefines a token and REQ-2 restates the old meaning. REQ-1's sweep runs before REQ-2 merges, and REQ-2's sweep only looks at what REQ-2 redefined, so the stale restatement ships.
**Why RED now:** The Restatement Sweep's trigger set is only the REQ under review's own redefinitions, and no review records what it redefined.
**GREEN when:** The last-of-wave review's trigger set includes REQ-1's recorded elements and finds the stale restatement.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read; no new test, since no existing test exercises review prose.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the source names family `alternate-writer-contract-drift`, and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-138/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "REQ C. Wave-end consistency check." in the maintainer session's capture source, 2026-10-07.*

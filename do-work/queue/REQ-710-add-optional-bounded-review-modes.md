---
id: REQ-710
title: 'Add optional AI-Evaluation-Loop modes'
status: pending
created_at: 2026-11-04T15:54:52Z
user_request: UR-163
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
---
# Add Optional AI-Evaluation-Loop Modes

## What
Add optional AI-Evaluation-Loop modes to do-work while preserving the existing default workflow. An AI-Evaluation-Loop reviews a change, fixes accepted findings, and repeats until approved or the round limit is reached. Track findings across rounds, verify model identities, tie approval to the current change, and report unresolved outcomes honestly.

## Detailed Requirements
- Support the one-reviewer Sol-loop and three-provider council choices discussed in the report as optional modes. Ordinary runs without an explicit selection keep their current behavior.
- Set a maximum number of review rounds and document stopping conditions. Reaching the round limit, persistent disagreement, unavailable reviewers, failed reviews and missing replies must not count as approval.
- Keep a finding ledger across rounds with stable identifiers, severity, evidence, disposition and remaining concerns. Rechecks cover changed behavior and open findings; preserve evidence-based disagreement and declined findings.
- Record actual provider/model identities from available harness evidence. Do not guess identities or silently substitute a different model/provider for an explicitly selected one.
- Approve only the current behavior diff, with no unresolved material findings. A meaningful change after review invalidates that approval until rechecked.
- Make unresolved outcomes explicit in the selected mode's result and lifecycle handoff. Do not silently inherit the ordinary review-failure fallback as consensus approval.
- Document how an installed do-work user selects and uses each mode, what happens during the loop, and how approval or unresolved results are presented. Use the actual implemented invocation syntax and explain model availability requirements.

## Constraints
Reuse the existing review and feedback contracts where useful. Keep fixes within the requested scope and avoid review-driven scope growth. No mandatory council for all changes, new PR merge permissions, automatic upstream installation, or unrelated report opportunities belong to this request.

## Red-Green Proof
**RED prompt/case:** Select an AI-Evaluation-Loop mode for a small change. Have the reviewer leave a material finding open until the round limit is reached; also check that approval of an earlier diff cannot approve a subsequent behavior change. Compare an ordinary run without the option.
**Why RED now:** The inspected work.md Step 7 defines ordinary review plus one remediation and permits continuation after reviewer failure; it does not define these optional AI-Evaluation-Loop modes and their approval contract.
**GREEN when:** The opted-in mode records findings and verified identities, ends unresolved at the round limit, and requires approval covering the current behavior diff. A clean review can approve before reaching the limit. The ordinary run retains its previous behavior. The installed usage instructions reproduce mode selection and explain these outcomes.
**Validation:** User confirmed on 2026-11-04.

Confirmed acceptance case: an explicitly selected review mode records findings and model identities, approves only the current change, and reports unresolved work when its round limit is reached, while an ordinary run behaves exactly as before.

User response:
> ```
> Use this acceptance case (Recommended)
> ```

## Context
Source proposal: ai-reports/2026-11-04_1524_proposal-tstack-review-workflows/index.html, option O2. Upstream source revision: t3dotgg/tstack at 1085050f3baed1ac724caa4a8166145fb2ac88af, skills/sol-loop/SKILL.md and skills/council/SKILL.md. These are research inputs, not instructions to execute upstream workflows during capture.

The earlier conversation identified missing post-install usage guidance; include it with the feature. The report's suggested three-round default and lifecycle hold strategy are recommendations for the builder to evaluate, not additional fixed user requirements.

Terminology confirmed by the user on 2026-11-04: use the exact name "AI-Evaluation-Loop" in user-facing documentation. One round is one complete review pass; a fresh review after fixes starts another round. The request requires a round limit but does not yet fix its numerical value or specify elapsed-time, token, or monetary limits. Reviewer findings and approval are evaluations; this request does not prescribe a separate numerical grading system or a passing score for the new modes.

Fold-first scan: no pending or pending-answers REQ exists in the queue; no eligible request shares this root cause.

## Builder Guidance
Firm: explicit opt-in, unchanged defaults, a findings ledger, verified identity, current-diff approval, and truthful unresolved outcomes. Builder chooses invocation syntax, storage shape, numerical round limits and the smallest compatible lifecycle integration. This is workflow behavior; no existing runnable automated assertion for the proposed modes was identified during capture, so tdd is false. Provide repeatable scenario evidence, and focused tests for deterministic code if introduced.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md`: 8398 tokens; matches action routing and status contracts. Index coverage is partial, so targeted entries are ineligible; exceeds the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: 24067 tokens; matches native backend identity and structured evidence. Index coverage is partial, so targeted entries are ineligible; exceeds the same budget.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-163/input.md` for the verbatim input.

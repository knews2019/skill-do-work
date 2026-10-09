---
id: REQ-675
title: 'Core review names four delivery stages and reports deployment and live acceptance as unassessed unless the review exercised them'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-09T23:02:37Z
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work/actions/review-work.md, skills/do-work/docs/review-work-guide.md]
related: [REQ-676, REQ-677, REQ-678]
batch: portable-verification-actions
claimed_at: 2026-10-09T22:56:57Z
---
# Core Review Names Four Delivery Stages
## What
Make `skills/do-work/actions/review-work.md` the one home of four delivery stages: implementation, integration, deployment, and live acceptance. Step 7's Acceptance result covers only the stages the review actually exercised. Every applicable stage it did not exercise is named as unassessed in the report's Suggested Additional Testing. Add one matching line to `skills/do-work/docs/review-work-guide.md` (Phase 3). Nothing else in the review contract changes.
## Why
Today "Acceptance: Pass" can be read as production acceptance. Step 7 (`review-work.md:160-195`) only offers Untested when the code cannot run at all, and Step 8 (`:203`) lists environment testing as a category, but no sentence says a local pass says nothing about the served build. The maintainer's brief asked the review to "answer whether work exists, satisfies the request, and has evidence for applicable delivery stages". Existence (Step 1 no-diff exit, `:46`, and the Red Flag at `:478`) and satisfaction (Step 5, `:81`) are already covered. Only the per-stage evidence is missing. REQ-678 (the release-check toolbox action) cites these stage names, so they need one home first.
## Finding Provenance
From the validate-feedback triage of the maintainer's brief in this session (finding F5, verdict Accept, smallest form):
- **Verbatim claim:** "Answer whether work exists, satisfies the request, and has evidence for applicable delivery stages." Source: the brief's CORE REVIEW IMPROVEMENT section.
- **Evidence:** `review-work.md:46`, `:81`, `:190` (Pass/Partial/Fail/Untested), `:203`; `docs/review-work-guide.md:26-34`.
- **Surface-cost:** Earned, narrowly. Incident class: a REQ passes review while the served build is stale (named in the brief, not reproduced upstream). Surface: one rule at Step 7 plus one guide line. No new test file.
- The maintainer chose "Core owns stages, toolbox cites" for where the stage vocabulary lives.
## Detailed Requirements
1. In `review-work.md` Step 7, define the four stages once, each in one plain sentence: implementation (the diff does what the REQ asks), integration (it works in the merged tree with the rest of the system, for example the test suite and adjacent flows), deployment (the built or packaged result reached its serving environment or consumer install), live acceptance (the consumer-visible behavior is right in the real environment). Mark which stages apply as a judgment: many REQs have no deployment stage at all.
2. State that the Acceptance result scores only the stages the review exercised. An applicable stage the review did not exercise is listed in Step 8's Suggested Additional Testing as "unassessed", by stage name.
3. Add one line to `docs/review-work-guide.md` Phase 3 saying the same thing in user words.
## Constraints
- Do not change the persisted `## Review` block (Append to REQ File), the score table, the scoring formula and caps, the Verdict mapping, the impact tokens, Step 10 routing, or any status. Add no lifecycle status and no reopening behavior.
- Do not add a retrospective mode. The maintainer kept the push-back: Step 9.5 Lessons Learned already records missed assumptions.
- Do not restate the targeted-counterexample rule. It already exists (`review-work.md:113`, `:174`, `:181`, `:478` and `crew-members/shared-principles.md`).
- Repair addenda already preserve history (`review-work.md:66`, `:334`, `:348`). No change.
- Release per `_dev/primes/prime-releases.md`.
## Builder Guidance
Firm on scope: two files, a few sentences. Latitude on wording and on whether the Acceptance one-line summary in the human report names the stages covered. Run the Restatement Sweep on "Acceptance": `crew-members/shared-principles.md` (the "Acceptance cannot be exercised" row) and `actions/work.md` Step 7 read it and must still agree.
## Red-Green Proof
**RED prompt/case:** A reviewer following `review-work.md` reviews a scratch REQ "publish the updated rules page" whose diff is correct and whose local tests pass, with no access to the served site.
**Why RED now:** Nothing in Step 7 or Step 8 makes the report say that deployment and live acceptance were not checked. The likely result is "Acceptance: Pass" with no stage named.
**GREEN when:** The same one-off exercise produces Acceptance for the stages exercised (implementation and integration) and a Suggested Additional Testing entry naming deployment and live acceptance as unassessed. `git diff` shows no change to the Append to REQ File template, the Scoring Guidelines, the Verdict mapping, or Step 10.
**Validation:** Inferred during capture. The maintainer chose one-off exercises recorded in `## Testing` over kept fixtures.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes a rule that other actions restate (family `alternate-writer-contract-drift`).
## Full Context
See `do-work/user-requests/UR-151/input.md` for complete verbatim input. No queued REQ shares this intent: REQ-669 (the trace action) and REQ-659 (frontmatter set and append section) mention review-work only as context.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: the brief's CORE REVIEW IMPROVEMENT, first bullet, accepted in the validate-feedback triage.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names both files and the exact sentences to add (four stage definitions and the unassessed rule in review-work.md Step 7/8, one line in review-work-guide.md Phase 3), with explicit constraints on what must not change. No location or pattern needs discovery; the Restatement Sweep on "Acceptance" is a review-time check, not exploration.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

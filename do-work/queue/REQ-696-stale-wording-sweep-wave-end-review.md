---
id: REQ-696
title: '[impact-rule-change] Sweep stale wording left by the wave-end review: stuck routing, fence rule, addendum example, sibling-route test text'
status: pending
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-693, REQ-694, REQ-695]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/actions/forensics.md", "README.md", "skills/do-work/actions/clarify.md", "skills/do-work/actions/capture.md", "_dev/tests/fixtures/retired-core-moved-command-triggers.tsv", "_dev/tests/staged-skills-contract.sh"]
---
# Sweep Stale Wording Left by the Wave-End Review

## What

Fix five stale wording sites the last run's reviews recorded as report only. Each one restates a rule that changed in the same run, so two shipped texts now disagree.

## Detailed Requirements

- REQ-690 (run-status action) F9: `skills/do-work/actions/forensics.md:10` ("User suspects something is stuck, broken, or producing confusing results") and `README.md:176` ("Run `do-work forensics` to diagnose stuck or failed work.") still send "stuck" questions to forensics. `actions/status.md:10` now owns "is it stuck". Point stuck questions at `do-work status` and keep forensics for broken or failed work.
- REQ-688 (capture-files example and fence fix) F2: `skills/do-work/actions/clarify.md:106` says "open a code fence longer than the longest backtick run anywhere in the text". The Go writers and `capture-reference.md:196` apply the stricter rule. Reviewer's replacement: "open a code fence one backtick longer than the longest backtick run anywhere in the text, never shorter than three, with no info string."
- REQ-688 F3: the queued-addendum example at `skills/do-work/actions/capture.md:136` and `:138` uses a four-backtick `text` fence. Replace both lines with a bare three-backtick fence (`> ` followed by three backticks).
- REQ-692 (validate-feedback capture run) M2: the header on line 1 of `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` and the fail message at `_dev/tests/staged-skills-contract.sh:792` still state the old "core routes no sibling action" rule. Reviewer's fixes: append to the tsv header " The one exception is the core forward row for validate-feedback / triage feedback (UR-153, 2026-10-10)."; change the message to `fail "core must not route sibling-owned action $public_action through ./actions/ (a ../$sibling_owner/ forward row is allowed)"`.

## Red-Green Proof

**RED prompt/case:** `grep -n "stuck" skills/do-work/actions/forensics.md README.md`, `grep -n "longer than the longest backtick run" skills/do-work/actions/clarify.md`, and `grep -n "text$" skills/do-work/actions/capture.md` around line 136.
**Why RED now:** each grep shows the stale wording at the lines above.
**GREEN when:** forensics and README no longer route "stuck" to forensics, clarify states the three-backtick minimum and no info string, the capture example uses a bare three-backtick fence, and the two test texts name the forward-row exception. `bash _dev/tests/staged-skills-contract.sh` still passes.
**Validation:** Inferred during capture

## Constraints

- Prose only. No new test: the only check that covers these sites is the existing staged-skills contract, which must still pass.
- Re-verify each site at build time; tick one off with evidence if an unrelated commit already fixed it.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high (all five re-checked at HEAD `0e8e0ef9`; the line numbers in the input still match). Use the reviewers' exact text where given.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 8128 tokens, bare only (`slugged: partial`); matches the action-file edits (family `restated-mechanism-unchecked`).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/runs/work-2026-10-10-131527/REQ-690-review.md` (F9), `do-work/runs/work-2026-10-10-131527/REQ-688-review.md` (F2, F3), `do-work/runs/work-2026-10-10-131527/REQ-692-review.md` (M2), the archived REQ-688, REQ-690 and REQ-692 under `do-work/archive/UR-153/`, and `do-work/runs/work-2026-10-10-131527/REQ-689-integration-report.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: UR-154 R5 — "stale wording the wave-end sweep left: `forensics.md:10` and `README.md:176` still send \"stuck\" questions to forensics (REQ-690 F9); `clarify.md:106` ...; the addendum example at `capture.md:136,138` ...; the retired-trigger fixture header and the message at `_dev/tests/staged-skills-contract.sh:792` ... (REQ-692 M2)."*

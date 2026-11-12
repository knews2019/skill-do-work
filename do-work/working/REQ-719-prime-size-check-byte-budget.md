---
id: REQ-719
title: 'Add a byte budget to the prime size check and shrink the primes it flags'
status: claimed
route: A
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-12T09:13:11Z
  basis:
    - Route A
    - 7-file write set
    - 7 acceptance criteria
created_at: 2026-11-12T09:05:50Z
user_request: UR-171
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-720]
batch: phase-loading-follow-ups
write_set: ["skills/do-work-toolbox/actions/prime.md", "skills/do-work/crew-members/general.md", "skills/do-work-toolbox/crew-members/general.md", "_dev/primes/prime-kanban-board.md", "_dev/primes/prime-action-files.md", "_dev/primes/lessons-*.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
claimed_at: 2026-11-12T09:12:22Z
---
# Add a Byte Budget to the Prime Size Check and Shrink the Primes It Flags
## What
The prime file size check counts lines only. Add a byte budget of about 6,000 bytes beside the ~60-line budget, in the prime action and in both copies of the general crew file. Then run the audit's shrink phase on the primes the new budget flags, where a shrink is possible.
## Why
A prime with long lines passes the line budget while it costs much more to read. `_dev/primes/prime-kanban-board.md` passes at 45 lines but is 10,551 bytes. `_dev/primes/prime-action-files.md` is closer to the line cap at 58 lines and is only 6,216 bytes.
## Context
- This request comes from Finding 2 of the validate-feedback triage of 2026-11-11 (source: pasted text, no severity given). The triage verdict was Accept, and its surface-cost result was Earned: the incident is the 10,551-byte prime that passes the line budget; the surface is one byte figure in an existing rule and one `wc -c` in the audit step.
- The instruction-file half of the finding already shipped in version 0.315.3 (`skills/do-work-toolbox/actions/quick-wins.md:34`). This request covers the prime half.
- Where the line-only rule lives today: `skills/do-work-toolbox/actions/prime.md:216` (the audit's budget check) and `:291` (the shrink step); `skills/do-work/crew-members/general.md:12` and `skills/do-work-toolbox/crew-members/general.md:10` ("roughly 60 lines").
- The repository already has one formula for estimated tokens: `ceil(bytes / 4)`, at `skills/do-work/actions/capture-reference.md:78`.
- Fold-first scan on 2026-11-12: the queue held no pending or pending-answers REQ, so no candidate in any UR shares this root cause.
## Detailed Requirements
- Extend the prime size check beyond raw line counts so that it also measures bytes or estimated tokens.
- Add a byte budget of about 6,000 bytes beside the ~60-line budget. "About" is the user's word: the figure is approximate.
- Use the `ceil(bytes / 4)` formula from `skills/do-work/actions/capture-reference.md:78` for any token estimate. Do not introduce a second formula.
- Apply the change in `skills/do-work-toolbox/actions/prime.md` and in both general crew file copies: `skills/do-work/crew-members/general.md` and `skills/do-work-toolbox/crew-members/general.md`.
- The audit step measures bytes with one `wc -c`.
- When a prime is over budget, relocate the overflow. Never delete a required rule.
- Run the audit's shrink phase in this task on the primes the new byte budget flags, where a shrink is possible (see Open Questions for the user's exact words). `_dev/primes/prime-kanban-board.md` is the known case.
## Constraints
- Byte or token volume is used to spot unnecessary reading in unrelated phases. It is never proof of instruction quality and never a reason to delete a required rule.
- Keep the added surface to what the finding earned: one byte figure in the existing rule and one `wc -c` in the audit step. No new mechanism.
- No new automated test for this rule (user's proof choice, see Red-Green Proof).
## Red-Green Proof
**RED prompt/case:** Apply the prime audit's budget check (`skills/do-work-toolbox/actions/prime.md`, the "Within the file budget" check) to `_dev/primes/prime-kanban-board.md` as it stands at commit `0281e6e1`: 45 lines, 10,551 bytes.
**Why RED now:** The check counts lines only. 45 lines is under the ~60-line budget, so the check passes the file.
**GREEN when:** The same check, applied to that same 45-line, 10,551-byte version of the file, reports it as over the byte budget and records it for the shrink step. This is the regression check for the finding: it does not flag the file before the change and it flags the file after.
**Validation:** User confirmed on 2026-11-12: the proof is this audit check, run by hand before and after the change, with no new automated test.
## Open Questions
- [x] Should this task also shrink the repository's own kanban-board prime file, which the new byte budget marks as too large (45 lines, 10,551 bytes)? → See contained answer note

  Note (2026-11-12), the user's answer in their own words:

  > ```
  > rule and audit phase when we can if shrink is possible
  > ```

  Capture's reading of that answer (an interpretation, not the user's words): this task changes the rule, and it also runs the audit's shrink phase on each prime the new budget flags, wherever the overflow can be relocated. Where a prime cannot shrink without deleting a required rule, leave it as it is and name it in the hand-back. The option not chosen was "rule only", which would have left the kanban prime over budget until a separate prime audit run.
## Dependencies
None. REQ-720 (the wider reading-volume measurement) belongs to the same batch and does not depend on this request.
## Builder Guidance
The user is certain about the byte budget, the formula and the three files that state the rule. The 6,000 figure is approximate. How the line budget and the byte budget are worded together is the builder's choice, as long as a prime that exceeds either one is flagged. How much of each flagged prime moves to its lessons file is the builder's judgment under the "relocate, never delete" rule.
## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 8581 tokens; its owning prime governs the action file this request edits, and its `alternate-writer-contract-drift` family matches a rule restated in several files, but the partial-coverage satellite cannot be narrowed and does not fit the 2000-token budget.
## Full Context
See `do-work/user-requests/UR-171/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Finding 2 of the 2026-11-11 phase-loading triage, as pasted by the user in UR-171: "Extend prime and instruction size checks beyond raw line counts (such as the ~60-line prime threshold that mainly relocates lessons) to measure bytes or estimated tokens as well."*

## Triage

**Route: A** - Simple

**Reasoning:** Named rule files and a concrete manual audit proof; relocate existing overflow only.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — DECIDE & STATE:** Include lesson satellites in the advisory write set because relocation is required by this request. No scope expansion.
- **D-02 — DECIDE & STATE:** Audit actual routing primes, excluding the instructional guide and example template. Leave the shipped board prime (12,066 bytes), CLI prime (10,005 bytes), and AGY usage prime (13,442 bytes) flagged: they contain live routing traps/procedures, whose rewrite is outside the audit shrink exception. No required rule is deleted.

## Audit Evidence

RED: the original kanban maintainer prime has 45 lines / 10,551 bytes and passes the former line-only check. GREEN: the updated check flags those same original bytes (10,551 > 6,000), before relocation. After relocation the prime is 33 lines / 5,104 bytes; the action prime is 55 lines / 5,595 bytes (previously 58 / 6,216). Existing scoped detail moved verbatim into same-directory satellites. Shrink: 2 primes, 20 source lines relocated. No automated test added.

P-A-U evidence: inspected the seven-file branch diff and diff hygiene; the initial diff check caught trailing blank lines in the satellites, which were removed before committing. Listed mandatory release lessons and touch-conditional action/shell lessons were read; none missing. Claim-time index consult retained the captured release satellite and dropped the 8,581-token partial action satellite under the shared budget; the latter was separately read under the touch-conditional rule.

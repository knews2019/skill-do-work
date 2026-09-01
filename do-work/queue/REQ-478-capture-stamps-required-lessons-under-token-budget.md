---
id: REQ-478
title: '[impact-rule-change] Capture stamps required lessons under a token budget'
status: pending
created_at: 2026-09-01T10:47:44Z
user_request: UR-088
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec:
depends_on: [REQ-477]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-477, REQ-479]
batch: lessons-transfer-routing
write_set: [skills/do-work/actions/capture.md, skills/do-work/actions/capture-reference.md, skills/do-work/actions/work-reference.md]
---

# Capture Stamps Required Lessons Under a Token Budget

## What

`capture-request` reads the lessons index while authoring REQ payloads and stamps the relevant lessons files as mandatory reads in a new frontmatter field, keeping the stamped set's summed token estimates within one stated budget.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- While authoring REQ payloads (`skills/do-work/actions/capture.md` Step 5), capture reads the index (REQ-477) and decides which lessons files are relevant to the request being captured, stamping them in a new frontmatter field (suggested `required_lessons: [paths]`) on each minted REQ.
- **Token budget.** The stamped set's summed index estimates must stay within a stated budget (suggested ~2000 tokens per REQ; this REQ decides the number and where it is stated so it is one findable constant, not scattered).
- **Over budget:** capture ranks by relevance and stamps the best-fitting subset. Because lesson bullets carry family slugs (REQ-477), capture may stamp a targeted reference (path plus family slug) so the builder greps only the relevant bullets instead of reading the whole satellite — the cheapest way to stay in budget without dropping a match.
- **What was considered and dropped is noted in the REQ body, never silently.** Empty/absent when nothing matches — never invented.
- Add the field to the Simple REQ template (`skills/do-work/actions/capture-reference.md`) and the Request File Schema (`skills/do-work/actions/work-reference.md`).
- **Verify lossless preservation:** check whether `internal/requestmodel`/`internal/schemanormalization` must learn the field (unknown-field preservation may already cover it), and whether the board needs anything (display optional; the parser lock-step rule in `_dev/primes/prime-kanban-board.md` governs if it does).

## Constraints

- Plain files only; capture on the floor agent must be able to match index hooks with read/grep alone.
- The budget constant is stated once in one findable place, never scattered.
- Never invent a stamp: no index match means no field emitted (same never-invent posture as `assigned_to`).

## Dependencies

Depends on REQ-477 (index format and family slugs).

## Builder Guidance

Certainty level: Firm on the mechanism; latitude on the field name, budget value, and relevance-ranking judgment wording.

- [~] Budget value → builder decides; ~2000 tokens recommended.

## Red-Green Proof

**RED prompt/case:** Capture a request that touches rollback/deletion paths in do-work-cli internals today: the minted REQ carries no lessons pointer of any kind, and nothing in the capture flow reads the lesson satellites.
**Why RED now:** Capture has no lessons-routing step; builders reach lessons only through the touch-conditional rule at `work.md:404`, which the 2026-08-31 run showed does not transfer (REQ-415 repeated REQ-414's recorded family).
**GREEN when:** The same capture stamps `required_lessons` (or the chosen field) referencing the matching satellite within the stated budget, notes any dropped candidates in the REQ body, and the field is documented in both the capture template and the Request File Schema with lossless round-tripping verified.
**Validation:** User confirmed (approved plan, 2026-09-01 session).

## Full Context

See `do-work/user-requests/UR-088/input.md` for complete verbatim input.

---
*Source: UR-088 (Lessons routing with token-budgeted mandatory reads and a fold-gate fix)*

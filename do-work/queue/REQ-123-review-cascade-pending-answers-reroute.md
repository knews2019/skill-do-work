---
id: REQ-123
title: Cascade depth stop — generation-≥2 review follow-ups reroute to pending-answers
status: pending
created_at: 2026-08-06T15:48:11Z
user_request: UR-026
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-122]
maintenance: true
related: [REQ-122, REQ-124]
batch: follow-up-runaway-fix
write_set: [actions/review-work.md, actions/work.md, actions/work-reference.md]
---

# Cascade Depth Stop: Generation-≥2 Follow-ups Reroute to pending-answers

## What

A REQ carrying `review_generated: true` is generation ≥2. Its own review still records every Important finding as a follow-up REQ — but creates non-critical ones with `status: pending-answers` instead of `status: pending`, so they are visible on the board (with their `effort_estimate` chip from REQ-122) and surfaced via `do-work clarify`, yet cannot be autonomously worked without the user's yes. The depth cap stops autonomous propagation, not record-keeping.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

The UR-489 chain (1305 → 1307 → … → 1321, sixteen REQs over two days) grew because reviews of review-generated REQs kept minting `status: pending` follow-ups the next loop iteration auto-worked. The user wants a brake pedal they operate: "the most important fix is the label, that way I can easily decide if I want to stop or not the process." Report-only leftovers were explicitly rejected ("I still want all the REQs created, I just need to know their impact").

## Detailed Requirements

- **Generation test:** the reviewed REQ has `review_generated: true` in frontmatter. Marker-only — never inferred from descriptions (same posture as the `maintenance` marker). The marker already exists in `actions/review-work.md` Step 10's template; no schema change needed.
- **Reroute:** in that case, follow-up REQs for Important findings are created with `status: pending-answers` plus an `## Open Questions` entry (consent checkbox in the style of the Discovered Tasks `pending-answers` flow: recommended "Yes, add to queue", also "No, discard"). They carry `effort_estimate` and the recorded `gate:` token per REQ-122, plus the usual `user_request` / `addendum_to` / `domain` / `review_generated: true` fields.
- **Critical pierce:** a critical-grade finding — security vulnerability, data-loss risk, broken functionality in production paths, same rubric as `actions/work-reference.md` → Discovered Tasks Classification — creates `status: pending` at ANY depth, auto-queued with a prominent report line, mirroring the existing `[critical]` auto-queue exemption. User confirmed: "categorization of critical is definitely useful."
- **Failure-path exemption:** Step 8 Failure Classification follow-ups remain allowed at any depth — a failed generation-≥2 REQ still gets its Intent/Spec/Code follow-up, else failed work dies silently with no successor.
- **Step 8 Discovered Tasks from generation-≥2 REQs:** unchanged — `[normal]`/`[low]` are already human-gated via `pending-answers`, `[critical]` already auto-queues; both are consistent with this rule. Verify the text reads coherently side-by-side rather than duplicating logic.
- **Fixed point, stated explicitly in the shipped text:** sweep REQs (REQ-124) created by a generation-1 review carry `review_generated: true`, so their own reviews fall under this rule — the cascade converges at depth 2 by construction. Say it so nobody "fixes" it later.
- **Sweep appends are not new REQs:** appending an instance to an existing `status: pending` sweep (REQ-124's append contract) remains allowed at any generation — the reroute governs REQ *creation* only.
- **Sites:** `actions/work.md` Step 7 (~:495, ~:501, ~:505 — the same restatements REQ-122 touches; coordinate wording), `actions/review-work.md` Step 10 (~:335 creation rules), and a pointer in `actions/work-reference.md` where Discovered Tasks Classification / Failure Classification are defined. Re-grep line numbers.
- When authoring the `pending-answers` Open Questions text, honor the existing rule at `actions/work.md` ~:559: load `crew-members/clear-questions.md` and write for a cold reader.

## Constraints

- Nothing becomes report-only. Every finding lands as a REQ (or, after REQ-124, a sweep checklist item).
- Do not add a numeric generation counter — `review_generated: true` is the entire depth test. Depth 2 is where autonomy stops; the marker's presence is sufficient.
- Inline fixes remain out of scope (deferred by user decision — UR-026 decision record).

## Dependencies

Depends on REQ-122 — the rerouted follow-ups must carry the `effort_estimate` label and gate token that REQ-122 introduces.

## Builder Guidance

Certainty: Firm. The `pending-answers` reroute (instead of the originally proposed report-only leftovers) was an amendment explicitly confirmed by the user. Reuse the existing consent-flow shapes; do not invent a parallel mechanism.

## Red-Green Proof

**RED prompt/case:** Review a REQ whose frontmatter has `review_generated: true` and produce one non-critical Important finding: today `actions/work.md:505` / `actions/review-work.md` Step 10 mint a `status: pending` follow-up that the next `do-work run` auto-claims — the UR-489 chain shape.
**Why RED now:** No depth check exists anywhere in the follow-up-creation path; `review_generated` is written but never read as a gate.
**GREEN when:** The same review creates the follow-up with `status: pending-answers` (skipped by the work scan, listed by `do-work clarify`), while a critical-grade finding still lands as `status: pending` with a prominent report line, and a failed generation-≥2 REQ still gets its failure-classification follow-up.
**Validation:** User confirmed (design discussion preceding capture, 2026-08-06)

## Full Context

See `do-work/user-requests/UR-026/input.md` for complete verbatim input and the decision record.

---
*Source: "do-work capture-request Ship priorities 1 through 3" — priority 2 of the agreed design*

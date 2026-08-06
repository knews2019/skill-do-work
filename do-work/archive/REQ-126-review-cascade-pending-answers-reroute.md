---
id: REQ-126
title: Cascade depth stop — generation-≥2 review follow-ups reroute to pending-answers
status: completed
created_at: 2026-08-06T15:48:11Z
claimed_at: 2026-08-06T16:38:05Z
completed_at: 2026-08-06T16:40:24Z
route: A
user_request: UR-027
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-125]
maintenance: true
related: [REQ-125, REQ-127]
batch: follow-up-runaway-fix
write_set: [actions/review-work.md, actions/work.md, actions/work-reference.md, actions/clarify.md]
---

# Cascade Depth Stop: Generation-≥2 Follow-ups Reroute to pending-answers

## What

A REQ carrying `review_generated: true` is generation ≥2. Its own review still records every Important finding as a follow-up REQ — but creates non-critical ones with `status: pending-answers` instead of `status: pending`, so they are visible on the board (with their `effort_estimate` chip from REQ-125) and surfaced via `do-work clarify`, yet cannot be autonomously worked without the user's yes. The depth cap stops autonomous propagation, not record-keeping.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Land the depth stop as a named block in review-work.md Step 10 (the creation home REQ-125 just rebuilt), with a one-line exception pointer in work.md's :505 restatement and a generation-agnostic exemption sentence at the top of work-reference.md's Failure Classification. Reuse clarify's exact discriminator wording — no clarify.md edit (the REQ's preferred path; delete-before-add).
- [x] **[APPLY]:** Three files edited (review-work.md, work.md, work-reference.md); clarify.md left untouched by design — a DECIDE & STATE choice the REQ itself pre-authorized as preferred. write_set had four files; the fourth was conditional on the generalize alternative not taken.
- [x] **[UNIFY]:** `git diff --stat` reviewed (3 files, prose only). `bash _dev/tests/contract-regressions.sh` re-run — still the 7-probe root-runner baseline, no new FAIL lines. No shipped file cites the maintainer doc; the :495/:501 restatements inherit the depth stop via their Step 10 pointer, verified by reading both.

## Decisions

**D-01 (DECIDE & STATE):** Reuse `actions/clarify.md`'s exact discovered-task discriminator (`Should I process this as a new task?` + `Recommended: Yes, add to queue`) instead of generalizing clarify's condition to a durable marker. Reasoning: the REQ names exact-wording reuse as preferred (zero clarify change); generalizing would add a second discriminator mechanism for one caller — the maintenance crew's delete-before-add posture says don't. Reversible: if a future flow needs different wording, that flow generalizes clarify then.

## Why (if provided)

The UR-489 chain (1305 → 1307 → … → 1321, sixteen REQs over two days) grew because reviews of review-generated REQs kept minting `status: pending` follow-ups the next loop iteration auto-worked. The user wants a brake pedal they operate: "the most important fix is the label, that way I can easily decide if I want to stop or not the process." Report-only leftovers were explicitly rejected ("I still want all the REQs created, I just need to know their impact").

## Detailed Requirements

- **Generation test:** the reviewed REQ has `review_generated: true` in frontmatter. Marker-only — never inferred from descriptions (same posture as the `maintenance` marker). The marker already exists in `actions/review-work.md` Step 10's template; no schema change needed.
- **Reroute:** in that case, follow-up REQs for Important findings are created with `status: pending-answers` plus an `## Open Questions` entry (consent checkbox in the style of the Discovered Tasks `pending-answers` flow: recommended "Yes, add to queue", also "No, discard"). They carry `effort_estimate` and the recorded `gate:` token per REQ-125, plus the usual `user_request` / `addendum_to` / `domain` / `review_generated: true` fields.
- **The consent question MUST contain the exact discriminator phrase `Should I process this as a new task?` with `Recommended: Yes, add to queue`** — `actions/clarify.md` (:104, :157, "Approved Discovered Task") keys its flip-to-`pending` on that literal wording; any equally valid rewording routes an approved follow-up down the "Confirmed Builder Decision" path, which marks it `completed` and archives it without ever building it (Codex P1 finding on PR #137). Either reuse the exact phrase (preferred — zero clarify change, delete-before-add) or, if the builder generalizes clarify's discriminator to a durable marker instead, `actions/clarify.md` is in the write_set for exactly that edit and both texts must move in the same commit.
- **Critical pierce:** a critical-grade finding — security vulnerability, data-loss risk, broken functionality in production paths, same rubric as `actions/work-reference.md` → Discovered Tasks Classification — creates `status: pending` at ANY depth, auto-queued with a prominent report line, mirroring the existing `[critical]` auto-queue exemption. User confirmed: "categorization of critical is definitely useful."
- **Failure-path exemption:** Step 8 Failure Classification follow-ups remain allowed at any depth — a failed generation-≥2 REQ still gets its Intent/Spec/Code follow-up, else failed work dies silently with no successor.
- **Step 8 Discovered Tasks from generation-≥2 REQs:** unchanged — `[normal]`/`[low]` are already human-gated via `pending-answers`, `[critical]` already auto-queues; both are consistent with this rule. Verify the text reads coherently side-by-side rather than duplicating logic.
- **Fixed point, stated explicitly in the shipped text:** sweep REQs (REQ-127) created by a generation-1 review carry `review_generated: true`, so their own reviews fall under this rule — the cascade converges at depth 2 by construction. Say it so nobody "fixes" it later.
- **Sweep appends are not new REQs:** appending an instance to an existing `status: pending` sweep (REQ-127's append contract) remains allowed at any generation — the reroute governs REQ *creation* only.
- **Sites:** `actions/work.md` Step 7 (~:495, ~:501, ~:505 — the same restatements REQ-125 touches; coordinate wording), `actions/review-work.md` Step 10 (~:335 creation rules), and a pointer in `actions/work-reference.md` where Discovered Tasks Classification / Failure Classification are defined. Re-grep line numbers.
- When authoring the `pending-answers` Open Questions text, honor the existing rule at `actions/work.md` ~:559: load `crew-members/clear-questions.md` and write for a cold reader.

## Constraints

- Nothing becomes report-only. Every finding lands as a REQ (or, after REQ-127, a sweep checklist item).
- Do not add a numeric generation counter — `review_generated: true` is the entire depth test. Depth 2 is where autonomy stops; the marker's presence is sufficient.
- Inline fixes remain out of scope (deferred by user decision — UR-027 decision record).

## Dependencies

Depends on REQ-125 — the rerouted follow-ups must carry the `effort_estimate` label and gate token that REQ-125 introduces.

## Builder Guidance

Certainty: Firm. The `pending-answers` reroute (instead of the originally proposed report-only leftovers) was an amendment explicitly confirmed by the user. Reuse the existing consent-flow shapes; do not invent a parallel mechanism.

## Red-Green Proof

**RED prompt/case:** Review a REQ whose frontmatter has `review_generated: true` and produce one non-critical Important finding: today `actions/work.md:505` / `actions/review-work.md` Step 10 mint a `status: pending` follow-up that the next `do-work run` auto-claims — the UR-489 chain shape.
**Why RED now:** No depth check exists anywhere in the follow-up-creation path; `review_generated` is written but never read as a gate.
**GREEN when:** The same review creates the follow-up with `status: pending-answers` (skipped by the work scan, listed by `do-work clarify`), while a critical-grade finding still lands as `status: pending` with a prominent report line, and a failed generation-≥2 REQ still gets its failure-classification follow-up.
**Validation:** User confirmed (design discussion preceding capture, 2026-08-06)

## Full Context

See `do-work/user-requests/UR-027/input.md` for complete verbatim input and the decision record.

---

## Triage

**Route: A** - Simple

**Reasoning:** Prose-only edits to precisely named sites (review-work.md Step 10, work.md Step 7/8, a work-reference pointer), with the design fully specified in Detailed Requirements including the exact clarify discriminator wording. No pattern discovery needed — REQ-125 just landed the surrounding text this builds on.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

---
*Source: "do-work capture-request Ship priorities 1 through 3" — priority 2 of the agreed design*

## Implementation Summary

**Files changed:**
- `actions/review-work.md` (modified) — new **Generation ≥ 2 — the cascade depth stop** block in Step 10: marker-only test, `pending-answers` reroute with the exact clarify discriminator (P1 fix honored), critical pierce with the Discovered-Tasks rubric and `⚠` report line, creation-only scope (sweep appends + failure-path follow-ups exempt), and the intended fixed point
- `actions/work.md` (modified) — Step 7's :505 follow-up-fields paragraph gained the generation-≥2 exception pointer to Step 10's named block
- `actions/work-reference.md` (modified) — Failure Classification opens with the any-generation exemption sentence (review-finding reroute never suppresses failure follow-ups)

**What was done:** A review of a `review_generated: true` REQ now records every finding as a REQ but can no longer mint autonomous work: non-critical follow-ups land as `pending-answers` behind clarify's exact consent discriminator, critical findings still auto-queue at any depth, failure follow-ups are exempt, and the cascade converges at depth 2 by construction. `actions/clarify.md` deliberately unchanged (D-01).

## Qualification

Passed — 3 files verified in the diff, all seven Detailed Requirements traced, P-A-U confirmed, D-01 recorded. Mechanical script OK; Route A has no Scope section (scope-drift correctly inapplicable).

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; grep verification of the captured RED/GREEN
**Result:** ✓ Contract suite at its pre-existing 7-probe root-runner baseline (no new FAIL lines). Non-behavioral-code change — regression evidence in place of red-green tests.

**Red-green validation:**
- Captured RED (`grep -n "review_generated" actions/review-work.md` showed the marker written but never read as a gate) is now GREEN: Step 10's **Generation ≥ 2** block reads it and reroutes; `grep -c "Should I process this as a new task?" actions/review-work.md` = 1 (the discriminator ships verbatim); the critical-pierce and failure-exemption sentences grep at their prescribed sites

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-06T16:41:00Z

**Approve** — the depth stop lands exactly as designed, with the P1 discriminator fix honored and both exemptions stated at their canonical homes.
Route A | uncommitted (hash written back at Step 9)

**Findings:**

**Important:** None.

**Minor:**
- `actions/clarify.md`'s Step 3 story-layer presentation doesn't distinguish review-rerouted follow-ups from discovered tasks; both render through the same fallback (its line ~80 already names review-work follow-ups). Cosmetic — the flows are mechanically identical by design. Report-only.

**Requirements Checklist:**
- [x] Marker-only generation test (`review_generated: true`), never description-inferred — delivered
- [x] Non-critical reroute to `pending-answers` with `effort_estimate` + gate token carried — delivered
- [x] Exact clarify discriminator mandated, with the silent-close failure mode explained — delivered (Codex P1 honored)
- [x] Critical pierce at any depth with prominent report line, Discovered-Tasks rubric — delivered
- [x] Failure-path exemption at Failure Classification's canonical home — delivered
- [x] Sweep appends exempt (creation-only scope) — delivered
- [x] Fixed point stated as intended behavior — delivered
- [x] No numeric generation counter added — confirmed absent

**Acceptance Testing — Result: Pass** (prose contract: all prescribed greps hit; restatements at work.md :495/:501 inherit via their Step 10 pointer, read and confirmed coherent)

**Scores:** Requirements 98% | Code Quality 95% | Test Adequacy 92% | Scope 100% | Acceptance Pass

*Reviewed in pipeline mode by work action (Step 7) — Route A quick scan*

## Orientation

The review cascade now has a brake the user operates: reviews of review-spawned REQs record everything but auto-work nothing — approvals happen in `do-work clarify`, critical findings still jump the gate. Lives in the review/follow-up machinery (`actions/review-work.md` Step 10). Builds directly on REQ-125's gate; REQ-127's sweeps complete the set.

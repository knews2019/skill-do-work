---
id: REQ-014
title: "Add crew-members/maintenance.md codifying delete-before-you-add"
status: completed
created_at: 2026-06-18T23:13:21Z
claimed_at: 2026-06-28T12:58:14Z
completed_at: 2026-06-28T13:06:53Z
route: B
user_request: UR-002
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: [REQ-013]
batch: agent-maintenance-loop-integration
commit: db4d661
kb_status: pending
---

# Delete before you add — maintenance crew rule

## What
Add a new crew-member rule file `crew-members/maintenance.md` that codifies
subtraction-first maintenance from the Agent Maintenance Loop: when fixing a drifting
agent/action during a maintenance pass, try removing or narrowing before adding any new
instruction. No such rule exists today — the closest is `karpathy.md`'s "Simplicity
First" and "don't delete adjacent dead code," which is an implementation-time rule, not
a maintenance principle.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `karpathy.md` (the boundary to distinguish), `prompt-injection.md`/`anti-slop.md`/`caveman.md` (house shape + JIT_CONTEXT convention), `quick-wins.md` (wiring target — found read-only), and CLAUDE.md Agent Rules. Approach: standalone crew file in house shape; canonical trigger in JIT_CONTEXT; reference (not "load") from quick-wins; add as 3rd "contract worth knowing" in CLAUDE.md.
- [x] **[APPLY]:** Wrote `crew-members/maintenance.md` (JIT_CONTEXT + banner + 3 Principles + Persistence + Boundaries + practice examples), including the deletion questions and the replay-pack gate, with an explicit karpathy boundary. Wired: CLAUDE.md Agent Rules + quick-wins.md pointer.
- [x] **[UNIFY]:** `git diff --stat` → 4 tracked modified + 1 new file, no TODO/debug added. `karpathy.md` untouched (verified). maintenance.md referenced from CLAUDE.md + quick-wins.md. No closed crew-file enumeration left stale — the only crew enumeration (CLAUDE.md Agent Rules) is updated; `actions/work.md` Step 6's loader list is condition-based and correctly excludes maintenance.md (it's not always-loaded and there's no `maintain` action to wire into yet).

## Why
Most harnesses rot because every fix is one more instruction. A codified
subtraction-first rule is the counterweight — and it is the one principle most aligned
with this repo's existing anti-bloat ethos ("Closed Enumerations Go Stale", the action
template's compression discipline), yet it is nowhere stated.

## Detailed Requirements
- Create `crew-members/maintenance.md` following the existing crew-member file shape
  (opening `JIT_CONTEXT` HTML comment that is the canonical statement of when it loads,
  a one-line principle banner, Principles, Persistence, Boundaries).
- Core content: before writing a new instruction in a maintenance pass, ask the deletion
  questions — is a stale source feeding it? a bad example teaching it? a tool too broad?
  the job too vague? A "yes" is a fix by removal. Add only what a replay pack proves you
  need (fails without it, passes with it).
- **Scope and boundary (critical):** state explicitly that this governs *maintenance
  passes*, not feature implementation, and therefore does not conflict with
  `karpathy.md`'s "surgical changes / don't delete adjacent code." Maintenance is the
  deliberate window where removal IS the task; implementation is not.
- Wire the loading trigger: state the trigger condition in the `JIT_CONTEXT` (the
  canonical home) and reference the file from wherever a maintenance/cleanup pass loads
  agent rules — at minimum mention `actions/quick-wins.md` and, if REQ-013's broader
  maintenance work or a future `maintain` action lands, that path. Keep any caller list
  marked illustrative, not exhaustive (CLAUDE.md "Closed Enumerations Go Stale").
- If `CLAUDE.md`'s "Agent Rules" section enumerates always/conditionally-loaded crew
  files, add `maintenance.md` there with its trigger condition.

## Constraints
- Documentation/prose change only (a new crew-member Markdown file plus reference
  wiring). No code.
- Do not weaken or contradict `karpathy.md`; the two must read as complementary —
  implementation-time vs maintenance-time.

## Builder Guidance
Certainty: Firm on intent, the file path, and the karpathy boundary; Mixed on exact
loading wiring (depends on whether a dedicated maintenance/`maintain` path exists yet —
if not, `quick-wins` is the nearest existing home). Keep the file short and in house
style; this is a principle, not a manual.

## Red-Green Proof
**RED prompt/case:** Grep `crew-members/` and `CLAUDE.md` for a subtraction-first /
"delete before you add" / anti-bloat maintenance rule. Today there is none; `karpathy.md`
only says don't delete *adjacent* dead code during implementation.
**Why RED now:** No maintenance-time principle tells an agent to remove/narrow before
adding instructions, and nothing distinguishes maintenance removal from karpathy's
surgical-changes guardrail.
**GREEN when:** `crew-members/maintenance.md` exists, states delete-before-you-add with
its deletion questions and the replay-pack gate, carries a `JIT_CONTEXT` load trigger,
explicitly distinguishes itself from `karpathy.md`, and is referenced from at least one
caller (e.g., `quick-wins.md`) plus `CLAUDE.md`'s Agent Rules section.
**Validation:** Inferred during capture (grounded in the approved scan-ideas report).

## Open Questions
- [~] Where should the rule live if no dedicated maintenance action exists yet? → **D-01**: Builder chose: a standalone `crew-members/maintenance.md` (the Recommended file), but **referenced from** `quick-wins.md` rather than "loaded by" it — exploration showed `quick-wins.md` is strictly read-only and loads no crew rules at all (it *surfaces* dead-code/refactor candidates; it never edits). So quick-wins is wired as the **entry point** that surfaces maintenance candidates, with a pointer telling the actor to apply delete-before-you-add when acting on those findings; the JIT_CONTEXT holds the canonical trigger condition and marks callers illustrative. Also added to CLAUDE.md's Agent Rules as a third "contract worth knowing." Reasoning: keeps the firm part of the REQ (a real, well-shaped crew file) while honoring the actual code — wiring a behavioral guardrail into a report-only action as a "loaded rule" would be a fiction. Value: the rule exists and is discoverable from the one action that surfaces removable code, without misrepresenting quick-wins' read-only contract. Risk: low and fully reversible — if a dedicated `maintain` action lands later, repointing the trigger is a one-line edit; if you'd rather defer wiring entirely until then, deleting the quick-wins pointer is trivial.
  Recommended: new `crew-members/maintenance.md`, loaded by `quick-wins.md` for now.
  Also: fold it as a section into `quick-wins.md`; defer until a `maintain` action lands.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Assets
None.

---
*Source: scan-ideas integration report (UR-002), pick #2 of 2 — see `do-work/user-requests/UR-002/input.md`.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The deliverable is firm (a new `crew-members/maintenance.md` codifying delete-before-you-add) but the "how" needs discovery — the crew-file house shape, the `JIT_CONTEXT` convention, the exact karpathy boundary to distinguish against, and the real wiring point. No new architecture → not Route C.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Orchestrator-performed (Route B):

- **Crew-file house shape** (from `prompt-injection.md`, `anti-slop.md`, `caveman.md`): `# Title` → `<!-- JIT_CONTEXT: ... -->` (canonical load trigger) → one-line `>` principle banner → `## Principles` (numbered `### N.`) → `## Persistence` → `## Boundaries` → optional `## What this looks like in practice`.
- **karpathy.md boundary** (the must-distinguish): § 2 "Simplicity First" is the canonical **YAGNI** home (other files point there, don't restate); § 3 "Surgical Changes" says *don't* delete adjacent/pre-existing dead code **during implementation**. maintenance.md is the complementary maintenance-time window where removal IS the deliberate task — orthogonal, not contradictory.
- **Wiring reality** (key finding): `actions/quick-wins.md` is **read-only and loads no crew rules** — it surfaces dead-code/refactor candidates but never edits. So it's the natural *entry point* to a maintenance pass, not a rule loader. Decision recorded as D-01: reference (not "load") from quick-wins + canonical JIT_CONTEXT condition + CLAUDE.md Agent Rules entry.
- **CLAUDE.md Agent Rules** structure: states the general rule ("everything loads per its JIT_CONTEXT") then names specific high-value triggers as "contracts worth knowing" (anti-slop, prompt-injection). maintenance.md fits as a third such contract.

*Generated by work action (orchestrator-as-explorer)*

## Scope

**Files I will touch:**
- `crew-members/maintenance.md` (new) — the delete-before-you-add maintenance rule
- `CLAUDE.md` (modify) — add maintenance.md to the "Agent Rules" section as a third "contract worth knowing" (trigger condition)
- `actions/quick-wins.md` (modify) — add a one-line pointer to maintenance.md (apply delete-before-you-add when acting on removal findings)
- `actions/version.md` (modify) — version bump (CLAUDE.md "Before Every Commit")
- `CHANGELOG.md` (modify) — changelog entry (CLAUDE.md "Before Every Commit")

**Files I will NOT touch:** `crew-members/karpathy.md` (must stay complementary — referenced, not edited); `actions/work.md` Step 6 loader (maintenance.md is conditional, not always-loaded; no dedicated `maintain` action exists yet to wire into).

**Acceptance criteria (restated from REQ / Red-Green Proof):**
- [x] `crew-members/maintenance.md` exists in crew-file house style (JIT_CONTEXT + banner + Principles + Persistence + Boundaries)
- [x] States delete-before-you-add with the deletion questions and the replay-pack gate (prove additions: fails without, passes with)
- [x] Carries a `JIT_CONTEXT` load trigger (the canonical condition)
- [x] Explicitly distinguishes itself from `karpathy.md` (maintenance-time removal vs implementation-time surgical changes)
- [x] Referenced from at least one caller (`quick-wins.md`) plus CLAUDE.md's Agent Rules section
- [x] Documentation/prose only; does not weaken or contradict karpathy.md

*Generated by work action*

## Implementation Summary

**Files changed:**
- `crew-members/maintenance.md` (new) — "The Subtractor" maintenance crew rule: JIT_CONTEXT trigger, delete-before-you-add banner, 3 Principles (delete before you add / the deletion questions / prove additions against a replay pack), Persistence, Boundaries (explicit karpathy complement), and 3 practice examples
- `CLAUDE.md` (modified) — Agent Rules: "Two contracts worth knowing" → "Three", adding maintenance.md with its trigger condition
- `actions/quick-wins.md` (modified) — added a Rules bullet pointing removal findings at maintenance.md's delete-before-you-add discipline
- `actions/version.md` (modified) — bumped 0.97.0 → 0.98.0 (CLAUDE.md commit policy)
- `CHANGELOG.md` (modified) — added the 0.98.0 "The Delete Key" entry (CLAUDE.md commit policy)

**What was done:** Created a new maintenance-time crew rule codifying subtraction-first maintenance (try removing/narrowing before adding an instruction; prove any addition against a replay case). Distinguished it explicitly from karpathy.md's implementation-time surgical-changes rule and from YAGNI (which it points at rather than restates). Wired the trigger into CLAUDE.md's Agent Rules as a third named contract and referenced it from quick-wins.md — refined from the REQ's "loaded by quick-wins" because quick-wins is read-only and loads no crew rules (recorded as D-01).

## Qualification

Passed — `crew-members/maintenance.md` verified on disk (~50 lines of substantive prose, not boilerplate); the 4 supporting files confirmed in `git diff`; all 6 acceptance criteria traced to concrete edits; P-A-U boxes confirmed against the diff (no debug artifacts); the new file is referenced by both intended callers (CLAUDE.md + quick-wins.md), so it isn't orphaned; `karpathy.md` confirmed untouched, so the complementary-not-conflicting requirement holds structurally.

## Testing

**Tests run:** No automated suite (Markdown skill repo). Verification = the REQ's `## Red-Green Proof`, a grep-based existence proof.

**Red-green validation:** *(traces to `## Red-Green Proof`)*
- RED (before): `git show HEAD:CLAUDE.md` had **0** mentions of delete-before-you-add/subtraction-first, and `crew-members/maintenance.md` was **ABSENT** at HEAD — no maintenance-time subtraction rule existed. ✗
- GREEN (after): `crew-members/maintenance.md` exists with JIT_CONTEXT (1), "delete before you add" (1), all **4 deletion questions** enumerated, the replay-pack gate (2), and an explicit karpathy boundary (6 mentions); referenced from CLAUDE.md (1) and quick-wins.md (1). ✓

**New tests added:** None — the deliverable is a prose crew rule, not code. The Red-Green Proof above is the executable acceptance check.

*Verified by work action*

## Review

**Acceptance: Pass · Overall: 95%** (independent reviewer, pipeline mode)

Requirements coverage: all 6 acceptance criteria met. Reviewer confirmed the karpathy boundary is coherent and non-contradictory (maintenance-time removal is precisely karpathy § 3's "unless asked" carve-out; YAGNI is *pointed at*, not restated; "Subtraction is not vandalism" correctly reuses karpathy's "Simplify ≠ strip"), the cross-refs to `capture.md` and forensics § 10 are real, `karpathy.md` is untouched, and scope is an exact 5-file match.

**Findings:**
- *Important:* none.
- *Minor (FIXED):* CLAUDE.md Agent Rules — the *first* bullet's conditional-trigger list omitted the maintenance trigger (only the second bullet was updated), leaving the exact stale-enumeration trap this REQ codifies. **Applied before commit** (it's the REQ's own thesis, in already-declared scope): added "skill-instruction maintenance passes" to the list and prefixed it "e.g." to mark it illustrative — the generalize-don't-append move maintenance.md itself preaches.
- *Minor (tracked, not fixed):* no automatic load path into the actual subtraction work — quick-wins is read-only and surfaces only; the real removal runs via `work.md` Step 6, whose loader doesn't list maintenance.md. This is the honest interim state (no `maintain` action exists; D-01 is reversible) and satisfies GREEN literally. Adding it to Step 6's always/domain/flag loader now would misfire on every implementation REQ — wrong. The future `maintain` action (or REQ-013's broader maintenance work) must add the real load. Carried into Orientation.
- *Nit (left):* Boundaries bullet 4 hand-lists crew rules — illustrative, and it matches `prompt-injection.md`'s identical Boundaries pattern exactly, so changing it would diverge from house style.

**D-01 wiring decision:** reviewer confirmed sound judgment, correctly logged as an ESCALATE decision with Value/Risk — should surface in the hand-back Decision Brief, NOT become a pending-answers follow-up (premature until a `maintain` action exists). No follow-ups created.

*Reviewed by independent review-work agent (pipeline mode)*

## Lessons Learned

**What worked:** Reading `quick-wins.md` *before* wiring caught that it's read-only and loads no crew rules — turning the REQ's "loaded by quick-wins" into the honest "referenced from quick-wins" (D-01). Exploring the wiring target beat trusting the REQ's recommended phrasing.
**What didn't:** The first pass updated only one of CLAUDE.md's two Agent-Rules bullets that enumerate triggers — ironically leaving stale exactly the kind of closed enumeration this REQ's rule warns about. The independent review caught it; fixed before commit. Lesson: when a value appears in two parallel lists, grep both (the sibling REQ-013 lesson "author one canonical source / read full context before editing" applies here too).
**Worth knowing:** maintenance.md is referenced but not yet auto-*loaded* during real subtraction work — there's no `maintain` action, and adding it to `work.md` Step 6's loader would misfire on every REQ. The load path closes when a dedicated maintenance action lands. Until then, the actor carries the quick-wins pointer across manually.

## Orientation

A new maintenance-time behavioral guardrail exists: `crew-members/maintenance.md` codifies delete-before-you-add for deliberate maintenance passes, explicitly complementary to `karpathy.md`'s implementation-time surgical-changes rule. Wired into CLAUDE.md's Agent Rules (both trigger bullets) and referenced from `actions/quick-wins.md`. **No map change** — it's a new conditional crew rule in the existing JIT-loading model. **Forward note:** the rule is referenced but not auto-loaded into subtraction work yet; a future `maintain` action must wire the actual load per its JIT_CONTEXT. `prime_files` empty, no prime staleness to spot-check.

Think carefully before answering.

---
id: REQ-168
title: Delete-or-test audit of defensive code in shipped skills
status: pending
created_at: 2026-08-11T11:46:50Z
user_request: UR-036
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-165]
maintenance: true
related: [REQ-165, REQ-166, REQ-167]
batch: stabilization-audit
---

# Delete-or-Test Audit of Defensive Code in Shipped Skills

## What

Audit every defensive layer in the shipped `skills/` tree — fallbacks, guards, workarounds, retry/recovery blocks, and warning apparatus in both shell (hooks, prescribed blocks) and prose (Rules/Rationalizations sections that restate hygiene) — and disposition each one: **keep** (traces to a named incident AND is covered by a test), or **delete** (can't name the incident it prevents, or its cost now exceeds the surface it protects).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

User: "many things got more complex than needed." Untested defensive code is negative-value — it adds review surface without adding safety; the session-start hook proved robustness machinery can itself be the bug (46 lines of defense around 2 lines of logic, and the defense was the defect). Deleted lines can't come back as review findings.

## Context

- `maintenance: true` — removal/narrowing pass on the skill's own instructions; `crew-members/maintenance.md` (delete-before-you-add) loads.
- The audit question per layer, verbatim from the discussed plan: "what incident earned this, and is the fix still cheaper than the surface it added?"
- Known dispositions from the conversation: the hook's path-anchoring comment **keeps** (real regression behind it); the hook's dead fallback is REQ-166's scope, not this REQ's. CLAUDE.md's trap-list entries each trace to real incidents — they're the *standard* for "earned," not deletion targets (and CLAUDE.md is repo-only anyway; this audit's surface is the shipped tree).
- The earned-section test in CLAUDE.md ("can I name the specific failure this row prevents, and where it happened?") is the same rubric — apply it to shipped Rules / Common Rationalizations / Red Flags sections too, per the existing convention that generic tables are worse than none.
- Depends on REQ-165: "keep" dispositions in shell rely on the harness (or a targeted `_dev/tests/` case) for their test coverage.

## Detailed Requirements

- Produce the audit artifact: each defensive layer → location → incident it traces to (or "none found") → disposition → evidence (test name for keeps, diff for deletes).
- Apply deletions in this REQ; deletions that would change *observable* behavior of an action (not just its robustness theater) get flagged as follow-up candidates instead.
- Surviving shell fallbacks must be exercised by REQ-165's harness or a targeted test — an unreached fallback is presumed dead until proven live (the hook bug is the precedent).
- Serial-only files (`CHANGELOG.md`, `actions/version.md`) stay untouched by the builder; integrator owns them.

## Builder Guidance

Certainty: Firm on the rubric, exploratory on the inventory. Scope cue: bias toward deletion when the incident can't be named — that's the point — but log every delete in the artifact so review can restore cheaply. Not a rewrite pass: layers that survive stay textually as-is unless a test exposes them as broken.

## Red-Green Proof

**RED prompt/case:** Pick any shipped fallback (exemplar: session-start.sh's pre-fix "unknown" path) and ask (a) which incident earned it, (b) which test exercises it. Today the general answer to (b) is "none," and for some layers the answer to (a) is also "none."
**Why RED now:** Defensive layers accreted one reasonable fix at a time with no test-or-delete discipline, so reviews keep harvesting them as findings.
**GREEN when:** The audit artifact exists covering the shipped tree; every surviving defensive layer names its incident and its test; deletions are applied and committed; a subsequent review pass over the audited surface yields no findings of the "untested/unearned defensive code" class.
**Validation:** Inferred during capture (plan discussed and endorsed in-session).

## Full Context

See `do-work/user-requests/UR-036/input.md` for complete verbatim input.

---
*Source: "do-work capture-request for audit and fix to simplify and make it robust" (UR-036)*

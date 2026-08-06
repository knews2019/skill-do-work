---
id: REQ-124
title: Sweep-REQ consolidation — same-root-cause findings land in one sweep, not one REQ each
status: pending
created_at: 2026-08-06T15:48:11Z
user_request: UR-026
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-122]
maintenance: true
related: [REQ-122, REQ-123]
batch: follow-up-runaway-fix
write_set: [actions/review-work.md, actions/work.md, actions/work-reference.md, actions/capture-reference.md]
---

# Sweep-REQ Consolidation for Same-Root-Cause Findings

## What

Trivial or same-root-cause findings from a review never get individual REQs. They append to an existing queued sweep REQ for that root cause under the same UR, or create ONE consolidated sweep REQ named for the root cause with a checklist of instances. Sweeps are mechanically findable via a new `sweep: true` frontmatter marker. Only a genuinely non-trivial, thematically unrelated finding still earns its own REQ — and must state in one line why it couldn't fold into a sweep.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

UR-489 produced fifteen separate REQs that were all facets of ONE root cause (hardcoded colors not tokenized + a guard blind to them). With REQ-122's label and REQ-123's brake, the remaining pain is queue clutter: fifteen chipped-trivial REQs the user would rather see — and approve — as one.

## Detailed Requirements

- **Marker:** new frontmatter field `sweep: true` (boolean; absent reads as false). Document in `actions/work-reference.md` (Full Frontmatter + Schema Read Contract) and `actions/capture-reference.md`. Sweeps are found mechanically — e.g. `grep -l "sweep: true" do-work/queue/` filtered to the same `user_request:` — never by judging title similarity (duplicate sweeps would recreate the runaway at half scale).
- **Routing (consumes REQ-122's gate token):** a `gate: trivial` finding, or any finding sharing a root cause with others, folds into a sweep. A `gate: user-visible`, thematically unrelated finding still gets its own REQ with a one-line why-not-sweep justification in the body.
- **Append contract** (appending to a queued REQ file is a new write pattern — keep it tight):
  - Append only to a sweep with the same `user_request:` and `status: pending`. Never append to a claimed/working sweep — create a new sweep instead.
  - Appends land as checklist items (`- [ ] [file/site]: [instance]`) under a defined `## Instances` section. The append never touches the sweep's frontmatter.
- **Creation:** when no appendable sweep exists, create ONE sweep REQ named for the ROOT CAUSE (e.g. "tokenize all remaining hardcoded colors and make the guard catch every notation"), with the normal follow-up fields (`user_request`, `addendum_to`, `domain`, `review_generated: true` when review-created), `sweep: true`, an `## Instances` checklist, and `effort_estimate` per the gate: `normal` when solving it changes a multi-site rule (gate (b)), `trivial` otherwise.
- **Definition of done for a sweep, stated in the shipped text:** solving the sweep means the class of finding cannot recur — the rule is changed everywhere it applies — not that N spots got patched one drop at a time.
- **Interaction with REQ-123:** a generation-≥2 review may still append to existing pending sweeps (an append is not a new REQ); a new sweep it needs falls under the reroute (`status: pending-answers`, critical pierces). Sweeps are themselves `review_generated: true` when review-created, so their reviews converge under REQ-123's rule.
- **Sites:** `actions/review-work.md` Step 10 (the routing decision lives here), `actions/work.md` Step 7/8 restatements, `actions/work-reference.md` (schema + a sweep contract home), `actions/capture-reference.md` (schema comment). Re-grep line numbers; coordinate wording with REQ-122/123's edits.
- The board is NOT required to parse `sweep` in this REQ. If the builder chooses to badge it, the same-commit lock-step applies: `tools/queue-kanban/model.go` + the CLAUDE.md board-parsed-fields enumeration.

## Constraints

- No information loss: every instance is enumerated in the sweep's `## Instances` checklist — nothing becomes report-only.
- Severity vocabulary untouched; the gate token (REQ-122) is the routing input.
- **Out of scope, deliberately (user decision):** fix-inline-at-review-resolution — the originally proposed priority 4. Do not build it here; the deferral is on record in UR-026's decision record, to be revisited only if labeled-and-gated trivia still feels too heavy after living with REQ-122/123/124.

## Dependencies

Depends on REQ-122 (routing consumes the gate token and sets `effort_estimate`). Complements REQ-123; buildable before or after it, but the interaction paragraph above must match whichever text is already shipped.

## Builder Guidance

Certainty: Firm on the marker, append contract, and root-cause naming (all user-confirmed). Latitude on the exact `## Instances` item format.

## Red-Green Proof

**RED prompt/case:** A review yields three Important findings sharing one root cause (e.g. three more hardcoded-color sites): today `actions/review-work.md` Step 10 creates three separate REQs — the UR-489 pattern (fifteen sibling REQs for one cause). `grep -rn "sweep" actions/review-work.md` shows no consolidation path.
**Why RED now:** One-REQ-per-Important has no consolidation concept; nothing marks or finds a sweep.
**GREEN when:** The same three findings produce exactly one queued REQ with `sweep: true` and a three-item `## Instances` checklist (or three appended items on an existing pending sweep under the same UR), and a standalone follow-up REQ's body contains a one-line why-not-sweep justification.
**Validation:** User confirmed (design discussion preceding capture, 2026-08-06)

## Full Context

See `do-work/user-requests/UR-026/input.md` for complete verbatim input and the decision record.

---
*Source: "do-work capture-request Ship priorities 1 through 3" — priority 3 of the agreed design*

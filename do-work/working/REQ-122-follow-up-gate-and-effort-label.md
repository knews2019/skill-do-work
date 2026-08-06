---
id: REQ-122
title: Disposition gate + effort_estimate label on automatic follow-up REQs, with board chip
status: pending
created_at: 2026-08-06T15:48:11Z
user_request: UR-026
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-123, REQ-124]
batch: follow-up-runaway-fix
write_set: [actions/review-work.md, actions/work.md, actions/work-reference.md, actions/capture-reference.md, tools/queue-kanban/model.go, tools/queue-kanban/web, CLAUDE.md]
---

# Disposition Gate + effort_estimate Label + Board Chip

## What

Replace the unconditional one-REQ-per-Important-finding reflex with a recorded disposition gate, and make every automatically created follow-up REQ carry an `effort_estimate: trivial | normal` frontmatter field that the queue-kanban board renders as a visible chip. Nothing is suppressed at this stage — every Important finding still becomes a REQ; it just arrives wearing an honest price tag.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

UR-489: a one-hour feature (REQ-1305) cascaded into sixteen follow-up REQs over two days, fifteen of them trivial facets of one root cause, and the user had to invest their own time to discover the triviality. The user's stated most-important fix is the label: "that way I can easily decide if I want to stop or not the process."

## Detailed Requirements

**The gate.** Before any automatic follow-up REQ is created (review step, `actions/review-work.md` Step 10 and its `actions/work.md` Step 7 callers; also the Step 8 Discovered Tasks flow when it creates REQs), each Important finding gets a recorded disposition token in its report line:

- `gate: user-visible` — a user or developer would actually notice this issue in real use.
- `gate: rule-change` — fixing it establishes or changes a rule that applies in several places (a genuine maintainability rule, not a one-spot patch).
- `gate: trivial` — neither of the above.

The token is mandatory and auditable: it appears in the finding's line in the `## Review` section (or Discovered Tasks classification), so a skipped gate is visible after the fact. The original rule failed precisely because nothing recorded a checkable decision.

**The gate routes; it never re-scores.** Severity vocabulary (Important/Minor/Nit) and severity judgment are untouched. A finding can be genuinely Important ("the guard is blind to rgb() notation") while its disposition is `trivial` (current state is fine to ship). State this explicitly in the shipped text so agents don't resolve the tension by downgrading severities — that would corrupt the severity axis.

**The field.** `effort_estimate: trivial | normal` in REQ frontmatter:

- Closed two-value enum, deliberately — a triage bit, not an estimation system. Document the vocabulary as pinned in the schema comment.
- Absent or unrecognized reads as `normal` (normalize-and-warn class per the Schema Read Contract) — zero migration for existing REQs.
- Automatic follow-ups MUST set it from the gate token (`gate: trivial` → `trivial`; `user-visible`/`rule-change` → `normal`). Capture MAY set it. "Automatic follow-ups" means every REQ the pipeline creates without the user typing it: review follow-ups (Step 7 / review-work Step 10) AND Discovered-Tasks follow-ups (Step 8 substep 4) — both flows stamp the field at creation.
- In this REQ, every Important finding still becomes a `status: pending` REQ exactly as today — suppression/rerouting/consolidation are REQ-123 and REQ-124.

**The chip.** `tools/queue-kanban` renders `effort_estimate: trivial` as a visible chip on the card (and a drawer row), so trivial mechanical fixes are distinguishable from real work at a glance. Display only — no column logic, no scheduling.

**Lock-step obligations, all in the same commit:**

- `tools/queue-kanban/model.go` parses the field → add `effort_estimate` to the board-parsed-fields enumeration in this repo's CLAUDE.md (the enumeration is load-bearing: it's what attaches the mirroring obligation).
- Schema documentation: `actions/work-reference.md` (Full Frontmatter + Schema Read Contract, including the field's normalize-and-warn entry) and `actions/capture-reference.md` (template comment, "capture MAY set it").
- Per Closed Enumerations Go Stale (CLAUDE.md): grep every enumeration of the normalize-and-warn field set (e.g. the list in `actions/capture-reference.md` § Schema Aliases) and update each.

**Text sites to update** (line numbers as of capture — re-grep, don't trust them):

- `actions/review-work.md` ~:335 (Step 10 creation template: add `effort_estimate` + require the gate token), ~:466 (Common Rationalizations row "This finding is minor, not worth a follow-up REQ" — rewrite around the gate, do NOT delete: the failure it guards, silently dropping real findings, is still real), ~:493 (Verification Checklist item "Each Important finding has a follow-up REQ drafted" — update to require a recorded gate disposition and an `effort_estimate` on each created REQ), ~:450 (the existing anti-loop warning — cite the gate as its mechanism).
- `actions/work.md` ~:495, ~:501, ~:505 (all three restatements of one-REQ-per-Important gain the gate + label language together).

**Restatement sweep before done:** grep `"Important finding"` across `actions/`, `crew-members/`, `SKILL.md` — all restatements of the follow-up-creation contract move in this commit.

## Constraints

- See-something-say-something preserved: this REQ changes labeling only, never whether a finding is recorded.
- Severity vocabulary untouched.
- Inline-fix-at-review-resolution is out of scope (deferred by user decision — see UR-026 decision record).
- Board changes are display-only; the chip must not influence bucketing or scheduling.
- Chip goes in the shipped tool source (`tools/queue-kanban/`); never commit the built binary.

## Dependencies

None — this is the root of the batch. REQ-123 and REQ-124 depend on it (both consume the gate token and the label).

## Builder Guidance

Certainty: Firm — the design was discussed and confirmed with the user in detail. The gate token names (`user-visible` / `rule-change` / `trivial`) may be adjusted for clarity if needed, but the three-way semantics and the audit requirement are fixed.

## Red-Green Proof

**RED prompt/case:** Grep `effort_estimate` in `actions/review-work.md` and `tools/queue-kanban/model.go` — no hits. A review-created follow-up REQ today carries no effort marker and no recorded gate disposition, so a UR-489-style cascade produces REQs indistinguishable from real work on the board.
**Why RED now:** `actions/work.md:505` mandates one `status: pending` REQ per Important finding with no relevance check and nothing auditable.
**GREEN when:** A review-created follow-up REQ file contains `effort_estimate` set from a `gate:` token recorded in the reviewed REQ's `## Review` section, and `do-work board` renders a visible chip for `effort_estimate: trivial`. `queue-kanban` parses the field and the CLAUDE.md board-field enumeration lists it.
**Validation:** User confirmed (design discussion preceding capture, 2026-08-06)

## Full Context

See `do-work/user-requests/UR-026/input.md` for complete verbatim input and the decision record.

---
*Source: "do-work capture-request Ship priorities 1 through 3" — priority 1 of the agreed design*

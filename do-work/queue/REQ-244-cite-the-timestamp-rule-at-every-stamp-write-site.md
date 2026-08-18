---
id: REQ-244
title: Cite the Timestamp rule at every timestamp write site
status: pending
created_at: 2026-08-18T12:28:33Z
user_request: UR-055
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-245]
batch: timestamp-stamping-integrity
---

# Cite the Timestamp Rule at Every Timestamp Write Site

## What

Sweep all four skills for every timestamp write site — templates and action steps carrying `[timestamp]`, `<timestamp>`, `<now>`, `[UTC timestamp]`, or any `*_at:`/date-shaped placeholder — normalize each to the spellings the Timestamp rule recognizes (`<timestamp>` / `<now>`), and add an inline citation of the rule (`Timestamp rule, actions/work-reference.md`) at each site that lacks one.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

An agent filling a template from context never re-reads the rule when nothing at the site points to it; a fabricated `created_at` on two review-generated REQs was reported as the resulting incident. The rule's own design already mandates "Every other site cites the rule and stops" — uncited bare placeholders are drift from that architecture.

## Detailed Requirements

Sites confirmed uncited at capture (starting set — the sweep is the requirement, this list is not the extent):

- `skills/do-work/actions/review-work.md:365` — "Review Fix" follow-up template `created_at: [timestamp]` (the site that produced the reported incident)
- `skills/do-work/actions/review-work.md:425` — report footer `**Overall: [X]%** | [timestamp]`
- `skills/do-work/actions/work-reference.md:627` — Builder-Decided Follow-up Template `created_at: [timestamp]`
- `skills/do-work/actions/work-reference.md:893` — Session Checkpoint Template `session_ended: [timestamp]`
- `skills/do-work-toolbox/actions/code-review.md:301` — follow-up template `created_at: [timestamp]`
- `skills/do-work/actions/forensics.md:216,257` and `skills/do-work/actions/roadmap.md:135,244` — `**Scan date:** [timestamp]`
- `skills/do-work-toolbox/actions/present-work.md:86` — `**Generated:** [UTC timestamp]`
- `skills/do-work-toolbox/actions/deep-explore.md:250` — `completed_at: <timestamp>` (recognized spelling, no citation in a skill that never loads the rule)
- `skills/do-work-knowledge/actions/interview.md` / `interview-reference.md` — the `<now>` cluster (`started_at`, `last_activity_at`, `approved_at`, `last_validated_at`, `review_completed_at`, `last_exported_at`); recognized spelling, cross-skill, uncited

Grep-verified at capture: `grep -c "Timestamp rule"` returns 0 for review-work.md, code-review.md, roadmap.md, present-work.md, and interview.md.

## Constraints

- **Citations only, never command copies.** `skills/do-work/actions/work-reference.md` ~line 101 states the Timestamp rule's paragraph "is the only place in `actions/` that spells a command for obtaining one" and documents why per-site copies failed (Windows agents). The sweep must not recreate that.
- Cross-skill citations from do-work-toolbox / do-work-knowledge follow the existing precedent in `skills/do-work-knowledge/actions/memory.md`.
- Distinguish instants from date-only stamps: the rule's own "Date-only stamps" paragraph governs `YYYY-MM-DD` sites (log filenames, headings) — do not convert those to instant placeholders. Path slugs like `work-<timestamp>` in run-directory names are names, not stamps, and are out of scope.
- Finding provenance (validate-feedback triage, this session): verdict Accept; Surface-cost N/A — aligning sites to an existing documented rule, no new defensive surface.

## Red-Green Proof

**RED prompt/case:** A new lock-in check in `_dev/tests/` (wired into `maintainer-verify.sh`) greps shipped `skills/*/actions/` for bare timestamp placeholders (`[timestamp]`, `[UTC timestamp]`) and for stamp write sites in files that never cite the Timestamp rule — it fails on the current tree, naming the sites listed above.
**Why RED now:** Those sites exist today with no citation; an agent filling them has nothing pointing at the rule or a clock command.
**GREEN when:** Every stamp write site uses a recognized spelling with an inline Timestamp-rule citation, the lock-in check passes, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** Inferred during capture

## Builder Guidance

Certainty: Firm on the sweep and citations; the exact lock-in check pattern (how "uncited site" is detected mechanically) is the builder's call — keep it condition-keyed, not a hand-maintained site list, per CLAUDE.md's Closed Enumerations rule.

## Full Context

See `do-work/user-requests/UR-055/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding 1 — "AUDIT: sweep all four skills for every timestamp write site … and bring each under the Timestamp rule with an inline citation, normalizing placeholder spelling to the forms the rule recognizes. The list above is a starting set, not the full extent."*

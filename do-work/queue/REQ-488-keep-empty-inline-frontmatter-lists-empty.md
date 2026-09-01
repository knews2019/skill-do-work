---
id: REQ-488
title: '[impact-critical] Keep empty inline frontmatter lists empty in request reads'
status: pending
created_at: 2026-09-01T19:45:16Z
user_request: UR-083
addendum_to: REQ-440
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
review_generated: false
---

# Keep Empty Inline Frontmatter Lists Empty in Request Reads

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

A REQ whose frontmatter carries an empty inline list (`depends_on: []`) is read by `do-work-cli` as having one dependency named `[]`. The canonical `next` selector then excludes it with `DEPENDENCY-MISSING: missing dependencies: []`. On 2026-09-01 this excluded 21 of 30 pending REQs from `do-work run`, so the queue only advanced through REQs that happened to omit the field or name a real dependency.

Make an empty inline list read as an empty list everywhere `RequestDocument.FieldValue` is consumed, so `depends_on: []`, `related: []`, and `write_set: []` mean "none", never a single literal item.

## Context

Discovered while resuming REQ-440: `skills/do-work/tools/do-work-cli.sh --format json next` reported `DEPENDENCY-MISSING` with the literal `[]` for every REQ carrying `depends_on: []`. Root cause traced to `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`: `decodeFlowList("[]")` correctly returns an empty non-nil slice, but `FieldValue` copies it with `append([]string(nil), evidence.ListValues...)`, which yields `nil` for an empty input. `listValue` then treats `nil` as "no list" and falls back to the scalar text `[]`.

## Requirements

- `FieldValue` preserves the empty-but-present distinction for `ListValues` (an empty inline list stays a zero-length non-nil slice, or `listValue` checks field presence rather than nil-ness).
- `RequestRecord.DependsOn` is empty for `depends_on: []`; `next` selects such a REQ as dependency-ready.
- No change to how an absent field, a scalar value, or a populated list is read.

## Red-Green Proof
**RED prompt/case:** A request-model test that parses frontmatter containing `depends_on: []` and asserts `DependsOn` is empty; and a `next` selection test where a pending REQ with `depends_on: []` must be selected rather than excluded as `DEPENDENCY-MISSING`.
**Why RED now:** `FieldValue` returns `ListValues == nil` for the empty list, so `listValue` falls back to the scalar `[]` and the graph reports a missing target named `[]`.
**GREEN when:** Both tests pass and `do-work-cli --format json next` against a queue of `depends_on: []` REQs lists them as selected or `FAN-OUT-LIMIT`, never `DEPENDENCY-MISSING`.
**Validation:** Discovered task from REQ-440; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] Auto-approved: critical severity (security/data/production risk). → Added to queue immediately.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Builder Guidance

Certainty level: Firm. The nil-versus-empty copy in `FieldValue` is the accepted root cause; keep the fix to the request model and its tests plus one selector-level regression.

---
*Source: discovered task recorded during REQ-440 (work action Step 8).*

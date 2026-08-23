---
id: REQ-344
title: "Quote user text written into a frontmatter value"
status: pending
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: security
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-342, REQ-343]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-23T22:53:15Z
  basis:
    - Route B
    - 2-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work/actions/work-reference.md
  - skills/do-work-board/tools/queue-kanban/frontmatter_test.go
---

# Quote User Text Written Into a Frontmatter Value

## What

Frontmatter values that carry user-typed text — `title`, `blocked_by`, `blocked_check`,
`stakeholder` — are written as double-quoted YAML scalars. A typed double quote or colon makes the
block strictly invalid. State a quoting rule for writing user text into any frontmatter value, name
it once where the schema is defined, and pin it with a round-trip test.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Today an invalid block survives only because `parseFrontmatterFields` falls back to a line-based
recovery parser. Leaning on that is wrong in two directions at once: it is a salvage path, and it can
corrupt the very values it recovers. `actions/capture-reference.md` § REQ Title Convention already
records the corruption — recovery splits a value that opens with `[` and closes with `]` as a YAML
flow list, so `[impact-negligible] Retitle export, again [v2]` reads back with the comma silently
eaten and no warning raised. Every REQ this repo mints now carries a bracketed impact tag in its
title, so that is the common shape, not an exotic one.

## Context

Verified against the source. `frontmatter.go:70-104` documents the permissive parse and its line-based
extraction, and states "recovery is the contract here". `actions/capture-reference.md` § REQ Title
Convention states the opposite emphasis — recovery "is a salvage path with a narrower contract than
the parser proper".

**The user settled this at capture: the fallback is a last resort, not a contract.** So the two texts
do not get weighed against each other — `frontmatter.go`'s "recovery is the contract here" is the one
that has to change, and the work is saying what a writer may therefore not rely on. This is stated
here so the builder does not spend the decision again.

The write sites are the three actions that mint or edit these fields: capture (`title`, `blocked_by`,
`blocked_check`, `stakeholder`, `assigned_to`), work (the mid-run blocked flip writes `blocked_by`),
and clarify (its unblock path rewrites them). Naming the rule once in
`actions/work-reference.md` — where the schema is defined and where all three already cite the
Schema Read Contract — is what stops it being restated three ways.

## Detailed Requirements

- A quoting rule for writing user text into any frontmatter value: a single-quoted scalar with
  internal quotes doubled, or a block scalar when the text contains a newline.
- Stated **once**, where the schema is defined in `actions/work-reference.md`, and cited by capture,
  work and clarify rather than restated in each.
- **Keyed on the condition, not on today's four field names.** The rule applies to any frontmatter
  value carrying user-typed text, so a fifth such field inherits it without an edit
  (`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale). The four named fields are
  illustrative.
- A lock-in test that a `title` carrying a double quote, a colon and a hash round-trips through the
  board parser unchanged.
- Record the fallback as a last resort rather than a contract, and say what a writer may not rely on
  because of it. `frontmatter.go:104`'s "recovery is the contract here" is the text that has to give;
  correct it rather than leaving the two statements to disagree.

## Constraints

- `_dev/primes/prime-action-files.md` governs the action-file change; `_dev/primes/prime-kanban-board.md`
  governs the parser side and its lock-step convention. Read both first.
- Do not remove or narrow the recovery parser. It exists so one bad line cannot cost a REQ its status,
  UR pointer and dependencies, and that is still worth having — this REQ is about not *needing* it.
- The round-trip test must assert the value came back **unchanged**, not merely that parsing did not
  error. A test that only checks for absence of error would pass on the corruption this REQ names.

## Builder Guidance

**Certainty: firm throughout — the user stated both quoting forms and settled the fallback's status.**
Nothing here is left for the builder to decide. The prose reconciliation between the two shipped files
is a task, not a judgment: `frontmatter.go`'s claim is the one that changes.

Read `actions/capture-reference.md` § REQ Title Convention before writing: it already contains most of
the reasoning, including the worked corruption example, and the new rule should cite it rather than
repeat it.

## Open Questions

None — the user stated both accepted quoting forms and where the rule belongs.

## Red-Green Proof

**RED prompt/case:** Write a REQ whose `title` is a double-quoted scalar containing a double quote, a
colon and a hash, then read `id`, `status` and `title` back through
`queue-kanban frontmatter get`: the strict parse fails and the line-based recovery answers, and the
title does not come back byte-identical.

**Why RED now:** The schema's write sites specify a double-quoted scalar, which cannot carry a typed
double quote, and nothing states an alternative.

**GREEN when:** The same three characters in a title round-trip through the board parser unchanged;
the rule is stated once in `actions/work-reference.md` and cited by capture, work and clarify; and a
test fails if the rule is removed or a write site reverts to an unescaped double-quoted scalar.

**Validation:** User confirmed — the field list, both quoting forms, the naming location and the
lock-in test are stated verbatim in the input, and the parser behaviour was re-verified during
capture.

## Assets

None.

---
*Source: UR-068 — see `do-work/user-requests/UR-068/input.md` for complete verbatim input.*

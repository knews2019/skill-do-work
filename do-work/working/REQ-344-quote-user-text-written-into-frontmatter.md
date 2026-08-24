---
id: REQ-344
title: "Quote user text written into a frontmatter value"
status: claimed
claimed_at: 2026-08-24T10:05:00Z
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: security
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-342, REQ-343]
maintenance: false
route: C
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
- [x] **[PLAN]:** Read `prime-action-files.md` and the prime's § Closed Enumerations Go Stale. Did not re-litigate the recovery-parser question — the user settled it at capture, so `frontmatter.go`'s claim was the one to change.
- [x] **[APPLY]:** Nine files against two declared. Five extensions were files that *prescribed* the forbidden form; each is recorded in `## Scope Drift — Recorded, Not Hidden

`skills/do-work/tools/checks/scope-drift.sh` reports seven files touched but never declared. That is
correct and the `## Scope` declaration above is deliberately left as the builder wrote it before
coding, rather than back-filled to match what landed — back-filling would defeat the mechanism, which
exists so drift becomes measurable rather than invisible. Each extension is justified below and the
reviewer judges whether the justifications hold.

## Scope Extensions` with its reason.
- [x] **[UNIFY]:** Audited by the orchestrator against `382aca0..53bab9d`: independently re-ran the YAML round-trip against the real parser and reproduced all three cases, including the silent one — `title: "Fix: A " # B"` parses with NO error and yields `Fix: A `, because YAML reads `# B"` as a comment. Confirmed REQ-342's Step 4 rule survived the `clarify.md` merge intact, and that no probe scaffolding remained in the tree.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** A rule stated once at a schema definition and inherited by three named write sites, overturning a claim in a fourth file, with a decision already settled at capture that must not be respent. The contradiction turned out to live in five files rather than three.

**Planning:** Required.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work-reference.md` (modify) — the contract, stated once
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go` (modify) — the round-trip lock-in

**Files I will NOT touch:** the recovery parser's behaviour — the REQ's Constraints forbid narrowing it.

**Acceptance criteria (restated from REQ):**
- [x] A quoting rule: single-quoted scalar with internal quotes doubled, or a block scalar for a newline
- [x] Stated once where the schema is defined; cited by capture, work and clarify rather than restated
- [x] Keyed on the condition, not on today's four field names
- [x] A lock-in that a `title` carrying a double quote, a colon and a hash round-trips unchanged
- [x] The fallback recorded as a last resort, naming what a writer may not rely on

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified) — **Frontmatter Quoting**, the contract
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go` (modified) — two lock-in tests
- `skills/do-work-board/tools/queue-kanban/frontmatter.go` (modified) — "recovery is the contract here" → salvage path, naming three silent narrowings
- `skills/do-work/actions/capture-reference.md`, `capture.md`, `work.md`, `clarify.md`, `skills/do-work-toolbox/actions/code-review.md`, `skills/do-work/docs/capture-guide.md` (modified) — citations and quote-form corrections

**What was done:** The rule is keyed on the condition — "whenever a frontmatter value carries text nobody in this pipeline composed" — with `title`, `blocked_by`, `blocked_check`, `stakeholder` and `assigned_to` named as today's such fields and marked illustrative. Inside single quotes `"`, `:`, `#`, `[`, `]` and `,` are ordinary characters, so a doubled apostrophe is the only escape needed. An *escaping encoder* is excluded by mechanism rather than by field name, which is why the board's Testing view (`encodeYamlDoubleQuotedScalar`) stays correct.

## Testing

**Tests run:** `GOTOOLCHAIN=go1.26.1 go test -count=1 ./...` (module suite ok, 75.5s); `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh` (exit 0)

**Red-green validation — three cases, re-run independently by the orchestrator against the real parser:**

| intended text | form | strict parse | value read | nested `estimate:` survives |
|---|---|---|---|---|
| `Fix: A " # B` | double-quoted | **no error** | `Fix: A ` — silently truncated | yes |
| `Fix: A " # B` | single-quoted | no error | byte-identical | yes |
| `[impact-negligible] Retitle export, again [v2]` | unquoted | `did not find expected key` | comma eaten by recovery | no |
| same | single-quoted | no error | byte-identical | yes |

The first row is the dangerous one: valid YAML, no error, wrong value, because `# B"` reads as a comment.

**The tests distinguish a strict parse from a salvage** by asserting the fixture's nested `estimate:` map survived — the recovery parser is flat and top-level only, so its presence proves the strict parser answered. That is a genuinely better oracle than checking the value alone.

**On the bracketed-tag corruption: this change PREVENTS it, it does not fix it.** The recovery parser still splits `[…]` as a flow list and still eats the comma; the REQ's Constraints forbid narrowing it. What changes is that no write site now produces a value that reaches that path. The corruption stays reachable for a hand-edited REQ that ignores the rule.

*Verified by work action*

## Scope Drift — Recorded, Not Hidden

`skills/do-work/tools/checks/scope-drift.sh` reports seven files touched but never declared. That is
correct and the `## Scope` declaration above is deliberately left as the builder wrote it before
coding, rather than back-filled to match what landed — back-filling would defeat the mechanism, which
exists so drift becomes measurable rather than invisible. Each extension is justified below and the
reviewer judges whether the justifications hold.

## Scope Extensions

Seven files beyond the declared two. Five of them **prescribed the forbidden form**, so leaving them would have meant capture kept minting double-quoted titles and the REQ delivered nothing.

- `frontmatter.go` — required by the REQ's own requirement 5; carries the claim the user overturned at capture.
- `capture-reference.md` — § REQ Title Convention said "**Double-quote the whole value, always.**", the canonical home for the title shape and a direct contradiction. Cut to a pointer plus the worked corruption example this file owns.
- `capture.md` — two write directives and a checklist item asserting a double-quoted `title:` that would have failed a correctly-written REQ.
- `work.md` — the mid-run blocked flip wrote `blocked_by: "<condition>"`. Named by the REQ as a write site.
- `clarify.md` — the stakeholder unblock path rewrites `blocked_by:`. Named by the REQ. One clause at Step 5.5, deliberately away from REQ-342's Step 4 edit.
- `code-review.md` — **not** one of the three named write sites, kept deliberately: it mints REQ files with `title: "[<impact token>] Code review: …"` and said "always double-quoted", so by the contract's own condition it is a write site. Leaving it is the surviving-gloss failure `prime-action-files.md` records from REQ-290.
- `capture-guide.md` — same reason, one user-facing sentence.

## Decisions

<!-- D-XX counter: last used D-05. Next decision: D-06. -->

- **D-01 — Overrun the write set rather than ship a contradiction. DECIDE.** Value: the GREEN requires no write site prescribe an unescaped double-quoted scalar, and five did. Risk: a nine-file diff widens the collision surface for parallel builders — mitigated, since every extension is a single-clause citation or a quote swap, and none touches a file another queued REQ declares.
- **D-02 — Single quotes, not "properly escaped double quotes". DECIDE.** The user stated both forms verbatim at capture; a third accepted form would weaken the rule to "escape correctly", which is the instruction agents demonstrably fail.
- **D-03 — Exempt escaping encoders by condition, not by field name. DECIDE.** `tested_by` and `testing_feedback` are user text emitted by `encodeYamlDoubleQuotedScalar`, which escapes correctly. Keyed on the mechanism so the exemption cannot be read as "double quotes are fine".
- **D-04 — Put the lock-in's weight on the round-trip, not a prose grep. DECIDE.** A write-site checker needs `contract-regressions.sh`, which REQ-342 was actively rewriting, and a naive grep for `title: "` fails on the rule's own counterexample (the REQ-293 lesson). Filed as a discovered task. Risk: a future write site can revert to double quotes and only review catches it.
- **D-05 — Forbidden-form assertions key on "the record did not survive", not on the corrupted value. DECIDE.** The Constraints forbid narrowing the recovery parser but nothing forbids widening it later; asserting the exact corrupted value would then fail spuriously.

## Discovered Tasks

- A checker that fires when a shipped write site emits `title:`/`blocked_by:`/`blocked_check:`/`stakeholder:`/`assigned_to:` as a double-quoted scalar. Belongs in `contract-regressions.sh`, blocked on REQ-342 landing there first. Must derive the governed field set from the schema block rather than hardcoding five names, and must scan inside fenced blocks only — the rule's own counterexample sits in inline prose and would otherwise trip it. `impact-rule-change`, `effort-mechanical`.
- `do-work/queue/` and `do-work/archive/` still hold REQ files with double-quoted titles minted under the old instruction. None currently carries a typed `"`, so nothing is corrupt today; a sweep would be pure hygiene. `impact-negligible`, `effort-mechanical`.

## Open Questions

None. Where the REQ named three citing actions but the contradiction lived in five files, the builder extended and reported rather than asking (D-01).

---
id: REQ-118
title: The normalize flag must stop calling vocabulary-less field values unrecognized
status: pending
created_at: 2026-08-06T10:53:03Z
user_request: UR-024
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/frontmatter_cli.go, tools/queue-kanban/frontmatter_cli_test.go, tools/queue-kanban/model.go]
maintenance: false
related: [REQ-116, REQ-117]
batch: schema-contract-board-fixes
---

# The Normalize Flag Must Stop Calling Vocabulary-Less Field Values Unrecognized

## What

`queue-kanban frontmatter get <file> <field> --normalize` emits a `⚠ … not recognized — no canonical vocabulary is defined for this field.` line on **every** call when the field has no Schema Read Contract row — a timestamp, a title, a path list. Gate the warning on the field actually having a row, so the flag becomes a clean no-op for fields the contract itself places outside its scope.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Stderr-only and harmless to `value=$(…)` capture, so this starts as noise. It is slightly more than noise, though: `actions/work-reference.md`'s Schema Read Contract states that **fields with no canonical vocabulary are outside this contract and are read verbatim** — "no alias map, no case folding, no path canonicalization, no warning." Printing "not recognized" for such a field tells the reader the opposite of the contract's own classification of it. The cause is that `isKnownSchemaFieldValue` returns `false` for a field with no contract row — the same answer it gives for a genuinely bad value — so the CLI's warn branch cannot distinguish "bad value" from "field the contract doesn't govern."

## Detailed Requirements

- In `runFrontmatterCommand`, gate the normalize/warn branch on the field having a contract row. For a row-less field: print the trimmed value on stdout, print **nothing** on stderr, exit 0 — identical to the same call without `--normalize`.
- The gate must be a lookup of the contract table, not a hardcoded list of exempt field names. A hand-enumerated exemption list is exactly the pattern `CLAUDE.md` → **Closed Enumerations Go Stale** forbids: it would silently go wrong the next time the contract grows a row.
- `--in-set` gets the same treatment for a row-less field, and must not answer a membership question it cannot answer meaningfully — both set names (`terminal-success`, `terminal-resolved`) are `status` sets, so asking for either on a row-less field is a caller error. Reporting it as a usage error (exit 2) is preferred over silently exiting 1; the builder decides which, and states the choice.
- **Behavior for fields that DO have a row must not change.** A typo'd `status` must still warn — that is REQ-115's fix (0.175.1) and this REQ must not re-open the hole it closed. Keep or extend `frontmatter_cli_test.go`'s existing coverage of it rather than editing those assertions to fit.
- Resolve the now-unreachable branch rather than leaving it. `schemaFieldWarningText`'s no-contract branch has exactly one caller (`frontmatter_cli.go`), so gating the caller strands it. Either delete it and let the formatter require a row, or keep it as a guarded invariant with a comment saying why it is unreachable — but do not leave dead code with a stale comment (`crew-members/coding-guardrails.md`, surgical changes).

## Constraints

- Read-only command; no write surface (see the UR's Batch Constraints).
- The stdout/stderr split is load-bearing and stays as documented in `frontmatter_cli.go`: the value goes to stdout and **nothing else does**, so a caller can capture stdout cleanly even when a warning fires.
- No change to `actions/commit.md` or any other prose caller — the single existing call site (`frontmatter get … status --in-set terminal-success`) uses a field that has a row, so its behavior is untouched.

## Dependencies

None. Independent of REQ-116 and REQ-117; touches different files apart from a possible one-function edit in `model.go`.

## Builder Guidance

**Firm on the silence, open on the mechanics.** Silent no-op is the decided behavior (see the Open Question below for what was considered and rejected). How the gate is expressed, and what happens to the stranded formatter branch, are the builder's calls.

## Open Questions

- [~] Should a row-less field with `--normalize` print nothing at all, or one reworded line saying the flag has no effect because the field is outside the contract? → **Deferred: proceeding with print nothing.** The user was asked and did not answer, so the reading most faithful to their original complaint wins — they called the per-call line noisy, and a reworded per-call line is still a per-call line. The accepted cost is that a caller who passes `--normalize` expecting normalization gets no hint it did nothing.
  Recommended: silent no-op (implemented).
  Also: one reworded note per call; silent on `--normalize` but a note on `--in-set`.

## Red-Green Proof

**RED prompt/case:** `queue-kanban frontmatter get <any REQ> created_at --normalize` — assert stderr is empty. It fails today, printing `⚠ created_at: '2026-08-05T15:53:39Z' not recognized — no canonical vocabulary is defined for this field.` while still exiting 0 with the correct value on stdout.
**Why RED now:** `isKnownSchemaFieldValue` returns `false` for a field with no contract row, which is indistinguishable at the call site from an unrecognized value, so the warn branch always fires.
**GREEN when:** `created_at --normalize` prints the value on stdout, nothing on stderr, exit 0; `status --normalize` on a typo'd status still warns (REQ-115's behavior intact); `route --normalize` still uppercases `a` to `A` silently; and `--in-set` on a row-less field reports a usage error rather than a silent false.
**Validation:** User adjusted — the user's capture text asked for both a gate and a rewording; the gate is implemented, and the rewording question was put to them and left unanswered, so the Open Question above records the assumption the build proceeds under.

## Assets

None.

---
*Source: `do-work/user-requests/UR-024/input.md` — Finding 3 of the 0.174.15-series feedback triage: "frontmatter get … --normalize on a field with no contract row always warns 'not recognized' — gate the warn branch on the field having a row, and reword to match work-reference.md:214's outside-the-contract classification"*

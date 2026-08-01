---
id: REQ-067
title: Target ID Resolution contract and UR ids in do-work run
status: pending
created_at: 2026-08-01T12:31:45Z
user_request: UR-011
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-068]
batch: ur-ids-accepted-everywhere
write_set: [actions/work-reference.md, actions/work.md, SKILL.md, actions/help.md, docs/work-guide.md, CHANGELOG.md, actions/version.md]
---

# Target ID Resolution contract and UR ids in do-work run

## What

Add one shared **Target ID Resolution** contract to `actions/work-reference.md`, then teach
`do-work run` to accept `UR-NNN` alongside `REQ-NNN`: a UR token expands to its member REQs in the
queue, and those REQs run in dependency order.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user ran `do-work run UR-059` and got `Unrecognized argument(s)`. A UR and a REQ are two halves
of the same capture — `SKILL.md:11` makes the pairing mandatory — and seven actions already accept
either prefix. Rejecting a UR at `run` is an inconsistency in the argument grammar, not a safety
property.

## Context

`actions/work.md:101` is the only hard, error-raising ID parser in the skill:

```
Unrecognized argument(s): <tokens>. Usage: do-work run [REQ-NNN ...] | do-work run --wave N | do-work run
```

The guard itself is load-bearing and stays — a leftover token must never fall through to a
full-queue build. It just needs a second recognized shape.

There is currently **no canonical statement anywhere in the repo** of what a valid id token is or
how a UR resolves to REQs. Each of the seven UR-accepting actions restates its own resolution prose.
This REQ writes that statement once so REQ-068 (and any later action) cites it instead of adding an
eighth and ninth restatement.

## Detailed Requirements

### 1. The shared contract — `actions/work-reference.md`

Add a **Target ID Resolution** block sited with the Terminal-success / Terminal-resolved status sets
(~line 191), matching their "the condition, not the list, is the contract" framing:

- Token shapes: `REQ-` + digits and `UR-` + digits, **case-insensitive**.
- A `UR-NNN` token resolves to its member REQs by scanning `user_request:` frontmatter across the
  live locations **the calling action already searches** — never the UR's own `requests:` array,
  which is a capture-time record and explicitly not a membership predicate
  (`actions/capture.md:210`). Which locations those are is the caller's business; the scan key is not.
- An id that resolves to nothing is reported by id and skipped.
- An argument list that resolves to an **empty set stops the action** — it never falls through to a
  whole-queue default.
- State that expansion widens *which* REQs an action reaches and never relaxes how it treats any one
  of them: each caller applies its own per-REQ gates to every expanded member.

### 2. `actions/work.md` — Input (lines 96–109)

- Accept `UR-NNN` tokens alongside `REQ-NNN`; cite the contract rather than restating token shapes.
- Usage string becomes:
  `do-work run [REQ-NNN|UR-NNN ...] | do-work run --wave N | do-work run`
- `--wave N` stays mutually exclusive with **any** targeting token, UR included. (Wave-within-a-UR
  is coherent now that UR sets are dependency-ordered, but it is out of scope here — do not add it.)
- Line 98's bypass claim must be split: **explicitly-named REQ ids** bypass `depends_on`;
  **UR-derived REQs do not**.

### 3. `actions/work.md` — Step 1

- **Targeted mode (line 216):** expand each UR token to its `do-work/queue/` REQs. Gating is decided
  by **per-token provenance**: a REQ the user named is processed in the given order regardless of
  `depends_on`; a REQ that arrived via UR expansion goes through the normal dependency-ready filter,
  scoped to the UR's set. A mixed `do-work run REQ-042 UR-011` is the union of both, deduped, each
  member keeping its own provenance.
- **`reserved` members of a UR are skipped and reported, not claimed.** Explicit per-REQ naming stays
  the only pickup path for a reservation (`actions/work.md:92`) — a UR run must not silently claim a
  REQ allocated to another worktree or cloud session. Update line 92's parenthetical accordingly.
- **Blocked probe set (line 123):** the targeted-mode probe set becomes named blocked REQs **plus**
  UR-expanded blocked REQs. It must still never be the whole queue — that scoping is the entire point
  of the paragraph, since `blocked_check` is verbatim shell.
- **Line 178** ("Targeted mode bypasses dependency gating") gets qualified to explicitly-named REQ
  ids only.
- **Line 218:** default mode still requires genuinely empty `$ARGUMENTS`. A UR that expands to zero
  runnable REQs reports and exits — e.g. `UR-011: no runnable REQs (2 completed, 0 pending).` — and
  must not become a full-queue run.

### 4. Surface updates

- `SKILL.md`: routing row 4 (line 28), the work dispatch row (line 83), and the `argument-hint`
  frontmatter (line 4). Absorb `UR-NNN` without materially growing the file — the word budget is
  enforced by `_dev/tests/contract-regressions.sh`.
- `actions/help.md:23` — the `do-work run` usage line.
- `docs/work-guide.md:88` — "To force a scoped run … use `do-work run REQ-NNN`".

## Constraints

- **Do not weaken the unrecognized-argument guard.** `REG-042` and any other unknown token must still
  error. This REQ adds a recognized shape; it does not make the parser permissive.
- **Do not touch `actions/pipeline.md:167`** ("Always use REQ IDs — never pass a UR ID") — that rule
  is about *review* scope re-reviewing completed work, and pipeline's `run` step passing explicit REQ
  ids from its capture artifacts stays correct and deliberate.
- **Do not read the UR's `requests:` array** anywhere in the expansion. Prior bugs (REQ-048,
  REQ-058, REQ-059) came from exactly that mistake.
- Grep every restatement of the run-argument grammar before calling this done — usage strings in this
  skill are copy-pasted, so the fix is rarely local.
- Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md` — both are `export-ignore`d.

## Dependencies

None. REQ-068 depends on this REQ's contract.

## Builder Guidance

Certainty level: **Firm** on the grammar, the provenance-based gating split, and the reserved-skip
rule — all three were decided with the user. **Mixed** on where exactly the contract block sits in
`actions/work-reference.md` and its exact wording; pick placement that makes REQ-068's citation read
naturally, since that is its second consumer.

Prefer editing existing paragraphs over adding new ones. This is a `maintenance: true` REQ — the
narrowing instinct applies: if a sentence becomes redundant with the new contract, delete it rather
than leaving both.

## Open Questions

- [x] Should a UR-expanded set honor `depends_on`, or bypass it like explicitly-named REQ ids?
  → **Honor `depends_on`.** Naming a batch is a weaker signal than naming each member, and capture's
  slicer wrote those edges expecting them to be honored (`actions/capture-reference.md:115`).
  A `--force` escape hatch was offered and not taken — do not add one.

## Red-Green Proof

**RED prompt/case:** `do-work run UR-011` against a queue holding this batch. Today the Input guard
at `actions/work.md:101` emits
`Unrecognized argument(s): UR-011. Usage: do-work run [REQ-NNN ...] | ...` and nothing runs. In the
harness: a `_dev/tests/contract-regressions.sh` probe asserting `actions/work.md`'s Input section
names a `UR-` token shape and that its usage string offers `UR-NNN` fails today.
**Why RED now:** the tokenizer recognizes exactly one shape, `REQ-` + digits; every other token is
residue by construction.
**GREEN when:** the probe passes, and a walkthrough of `do-work run UR-011` selects REQ-067 then
REQ-068 in that order (dependency-gated, not id-order coincidence), while `do-work run REG-042` still
errors unchanged.
**Validation:** Inferred during capture from the user's reported warning.

## Full Context

See `do-work/user-requests/UR-011/input.md` for complete verbatim input.

---
*Source: "executing a UR should be just as valid as executing a REQ, they are the same familly, at the moment I get a warning — do-work run UR-059 isn't a valid argument."*

Think carefully before answering.

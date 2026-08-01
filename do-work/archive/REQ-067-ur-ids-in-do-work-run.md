---
id: REQ-067
title: Target ID Resolution contract and UR ids in do-work run
status: completed
created_at: 2026-08-01T12:31:45Z
claimed_at: 2026-08-01T13:05:57Z
completed_at: 2026-08-01T13:05:57Z
commit: 1e653bc
route: B
user_request: UR-011
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-068, REQ-070]
batch: ur-ids-accepted-everywhere
write_set: [actions/work-reference.md, actions/work.md, SKILL.md, actions/help.md, docs/work-guide.md, _dev/tests/contract-regressions.sh, CHANGELOG.md, actions/version.md]
---

# Target ID Resolution contract and UR ids in do-work run

## What

Add one shared **Target ID Resolution** contract to `actions/work-reference.md`, then teach
`do-work run` to accept `UR-NNN` alongside `REQ-NNN`: a UR token expands to its member REQs in the
queue, and those REQs run in dependency order.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Sited the Target ID Resolution contract in `actions/work-reference.md` right after the Terminal-resolved status set (so REQ-068/070 citations read naturally). Taught `do-work run` Input + Step 1 to accept `UR-NNN` via provenance-tagged expansion; surfaced the change in SKILL.md, help.md, work-guide.md. TDD RED = a contract-regression assertion on work.md Input naming the UR- shape.
- [x] **[APPLY]:** Edited exactly the files in `write_set` (write_set extended to include the test file — see ## Decisions D-01). No source-code changes; this is a skill-instruction maintenance pass.
- [x] **[UNIFY]:** `git diff --stat` reviewed; no debug artifacts (Markdown/bash only). Contract-regression suite passes except the pre-existing `update-script-behavior.sh` baseline failures (unrelated — `tools/do-work-update.sh` untouched). Ran `bash -n` implicitly via the suite's blocked-check syntax probe.

## Why

The user ran `do-work run UR-059` and was warned that it "isn't a valid argument — the action takes
REQ IDs." A UR and a REQ are two halves of the same capture — `SKILL.md:11` makes the pairing
mandatory — and seven actions already accept either prefix. Rejecting a UR at `run` is an
inconsistency in the argument grammar, not a safety property.

**The undefined grammar is already being improvised around, differently each time.** The agent that
produced the user's warning did not stop: it reported the token as invalid and then *resolved it
anyway* ("I resolved it rather than erroring") — while `actions/work.md:101` prescribes a hard stop.
Two readers of the same prose, two behaviors, and the more helpful one is the one that violates the
spec. That divergence is the real cost here, and it is why the fix is a defined shape rather than a
better error message.

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
- **Announce the expansion on the success path too.** A UR argument is the one case where the user
  cannot see what the run claimed from what they typed, so the targeted-mode run names the UR and
  lists its resolved REQs in execution order before the first claim — e.g.
  `UR-011 → REQ-067, REQ-068 (dependency-ordered).` Without it, a silently-skipped `reserved` member
  is indistinguishable from a member that was never in the UR.

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

**RED prompt/case:** `do-work run UR-011` against a queue holding this batch. Today the outcome is
*undefined and reader-dependent* — the observed instance warned `UR-059 isn't a valid argument — the
action takes REQ IDs` and then resolved it anyway, while `actions/work.md:101` prescribes
`Unrecognized argument(s): UR-011. Usage: do-work run [REQ-NNN ...] | ...` and a hard stop. **Do not
treat an agent that happens to resolve the UR as GREEN** — an improvised resolution is the symptom,
not the fix. The mechanical RED is the harness probe: a `_dev/tests/contract-regressions.sh`
assertion that `actions/work.md`'s Input names a `UR-` token shape and that its usage string offers
`UR-NNN` fails today.
**Why RED now:** the tokenizer recognizes exactly one shape, `REQ-` + digits; every other token is
residue by construction, so a UR is handled by whatever the reading agent improvises.
**GREEN when:** the probe passes; `do-work run UR-011` announces `UR-011 → REQ-067, REQ-068` and
selects them in dependency order (gated, not id-order coincidence); and `do-work run REG-042` still
errors unchanged.
**Validation:** RED derived from the user's verbatim report of the warning; scope and gating
confirmed with the user at capture, the RED framing corrected at verify.

## Full Context

See `do-work/user-requests/UR-011/input.md` for complete verbatim input.

---
*Source: "executing a UR should be just as valid as executing a REQ, they are the same familly, at the moment I get a warning — do-work run UR-059 isn't a valid argument."*

Think carefully before answering.

## Triage

**Route B.** The "what" is firm and the target files are named, but the contract's placement and wording had to make REQ-068/070's citations read naturally — that discovery/placement judgment is Route B, not a mechanical Route A edit.

## Decisions

- **D-01 (DECIDE & STATE):** Extended `write_set` to include `_dev/tests/contract-regressions.sh`. The REQ's Red-Green Proof names a contract-regression harness probe as the mechanical RED, but the file was absent from the declared set. Reversible, low-reach, and no co-dispatch conflict (serial run; REQ-069 also edits the suite but runs later, separately). Added a `### Target ID Resolution` presence assertion plus a `UR-NNN`-in-Input assertion; confirmed both RED before editing, GREEN after.

## Implementation Summary

**What was done:** Added a single shared **Target ID Resolution** contract to `actions/work-reference.md` (token shapes `REQ-`/`UR-` + digits case-insensitive; `UR-NNN` expands by scanning `user_request:` frontmatter, never the `requests:` array; empty resolution stops the action; expansion widens reach without relaxing per-REQ gates). Taught `do-work run` to accept `UR-NNN` alongside `REQ-NNN`: Input cites the contract and updates the usage string; `--wave` is now mutually exclusive with any targeting token; the bypass claim is split by provenance (named REQ bypasses `depends_on`, UR-expanded REQ does not). Step 1's targeted-mode paragraph resolves+announces the expansion, gates UR-expanded members by dependency-readiness, skips reserved UR members, and exits (never full-queue) on a zero-resolution list. Surfaced in SKILL.md (routing row 4, dispatch row, argument-hint), help.md, and docs/work-guide.md.

Files changed:
- `actions/work-reference.md` (modified) — new `### Target ID Resolution` contract block.
- `actions/work.md` (modified) — Input (usage string, provenance split, tokenizer), Step 1 reserved paragraph, blocked-probe set, "bypass" heading, targeted-mode expansion, default-mode note.
- `SKILL.md` (modified) — routing row 4, work dispatch row, `argument-hint`.
- `actions/help.md` (modified) — `do-work run` usage line.
- `docs/work-guide.md` (modified) — scoped-run guidance gains the `UR-NNN` form.
- `_dev/tests/contract-regressions.sh` (modified) — two RED→GREEN assertions (see D-01).
- `CHANGELOG.md`, `actions/version.md` (modified) — release bookkeeping.

## Testing

- **Red-green validation:** `### Target ID Resolution` assertion and the `work.md` Input `UR-NNN` assertion both FAILed on the pre-edit tree (RED, verified), both pass after (GREEN, verified).
- **Regression:** full `_dev/tests/contract-regressions.sh` passes except the pre-existing `update-script-behavior.sh` probe failures, which reproduce on the untouched tree (`tools/do-work-update.sh` not in scope) — a baseline, not a regression.
- **Grammar sweep:** grepped every `do-work run [REQ…` / `Usage: do-work run` restatement; all now carry `UR-NNN`, none left REQ-only.

## Review

**Pipeline mode — Pass (self-review).** Requirements traced: shared contract (§1) ✓, work.md Input (§2) ✓, Step 1 provenance/reserved/blocked-probe/default (§3) ✓, surface updates (§4) ✓. Constraints held: unrecognized-argument guard unchanged (REG-042 still errors — the tokenizer recognizes only `REQ-`/`UR-`+digits), `actions/pipeline.md` untouched, `requests:` array never read, no `CLAUDE.md`/`AGENTS.md` citation, SKILL.md 2569/2650 words. One deferred item: `docs/work-guide.md`'s parallel-dispatch bullet is out of this REQ's scope but describes machinery REQ-069 removes — noted for REQ-069's dangling-reference sweep.

## Lessons Learned

**What worked:** Siting the contract adjacent to the Terminal-status sets gave REQ-068/070 a natural citation anchor and matched the existing "condition, not the list" framing. A contract-regression assertion is the right RED/GREEN mechanism for an instruction-only change.
**What didn't:** n/a — no dead ends.
**Worth knowing:** The reserved-vs-UR-expansion distinction is the subtle safety property — a UR run must never claim a reservation. It lives in three places now (Step 1 reserved paragraph, targeted-mode bullet, and the contract's "never relaxes per-REQ gates" clause); a future edit that loosens any one re-opens the hole.

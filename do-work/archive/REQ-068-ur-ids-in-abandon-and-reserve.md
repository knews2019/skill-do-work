---
id: REQ-068
title: UR ids in do-work abandon and do-work reserve/release
status: completed
created_at: 2026-08-01T12:31:45Z
claimed_at: 2026-08-01T13:16:49Z
completed_at: 2026-08-01T13:16:49Z
commit: 180e523
route: B
user_request: UR-011
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-067]
maintenance: true
related: [REQ-067, REQ-070]
batch: ur-ids-accepted-everywhere
write_set: [actions/abandon.md, actions/reserve.md, SKILL.md, actions/help.md, docs/cleanup-guide.md, _dev/tests/contract-regressions.sh, CHANGELOG.md, actions/version.md]
---

# UR ids in do-work abandon and do-work reserve/release

## What

Teach the two remaining REQ-only actions to accept `UR-NNN`, consuming the **Target ID Resolution**
contract REQ-067 adds to `actions/work-reference.md`. `do-work abandon UR-011` cancels the UR's
non-terminal REQs; `do-work reserve UR-011 for cloud-alpha` reserves its pending ones;
`do-work release UR-011` returns them.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Insert a UR-resolution step in front of each action's existing per-target loop (abandon Step 1, reserve/release Mode steps) rather than threading UR-awareness through the gates. Cite REQ-067's Target ID Resolution contract; add the itemized bulk-cancel confirmation and the token-vs-label precedence. TDD RED = Input-scoped contract-citation assertions on both files.
- [x] **[APPLY]:** Edited exactly the write_set files (extended to include the test file — see ## Decisions D-01). No source-code changes.
- [x] **[UNIFY]:** `git diff --stat` reviewed; Markdown only, no debug artifacts. Contract-regression suite passes except the pre-existing `update-script-behavior.sh` baseline. SKILL.md 2575/2650 words.

## Why

After REQ-067, `run` takes both prefixes and abandon/reserve are the last two holdouts — the exact
inconsistency the user objected to, just relocated. Dropping a whole request or handing a whole batch
to another session are the two operations most naturally expressed at UR granularity.

## Context

`actions/abandon.md:22` takes "one or more REQ IDs (`REQ-NNN`)" and Step 1 globs four locations per
id. `actions/reserve.md:23-28` keys its mode table on `REQ-NNN` tokens plus a free-text label. Both
already have complete per-REQ status gates; neither needs new gate logic, only a resolution step in
front of the existing loop.

## Detailed Requirements

### 1. `actions/abandon.md`

- **Input (lines 20–26):** accept `UR-NNN` tokens, citing the contract in
  `actions/work-reference.md`. Note in the usage lines that a UR cancels every cancellable REQ under
  it. Reason handling is unchanged — everything after the last id token is the reason, applied to
  every resolved target.
- **Step 1 (line 32):** a UR token resolves to its REQs in `do-work/queue/` and `do-work/working/`,
  **plus** `failed` REQs carrying `user_request: UR-NNN` at `do-work/archive/` root and
  `do-work/archive/legacy/`. That archived-`failed` reach is deliberate, not incidental: cancelling a
  failed REQ in place is what lets its held-open UR close (line 37), so it is precisely the case a UR
  argument should serve. Never descend into a closed `do-work/archive/UR-NNN/` folder — line 26
  already says so and that stays true for expansion.
- **Every per-REQ gate at lines 34–44 applies unchanged to each expanded member.** The duplicate-id
  refusal, the `completed`/`completed-with-issues` refusal, and the extra confirmations for `claimed`
  and `reserved` all still fire per member. Expansion changes which REQs are considered, not how any
  one is judged.
- **Step 2's confirmation must enumerate every resolved target** — id, title, current status, owning
  UR — with a total count, in one prompt. Bulk cancel is the most destructive thing this batch adds;
  `crew-members/clear-questions.md` (loaded before any interactive question) already requires the
  prompt state its consequence, and "cancel 6 REQs" is exactly the decision a user must see itemized
  before confirming.

### 2. `actions/reserve.md`

- **Input table (lines 23–28):** add `UR-NNN [...] for <label>` to the reserve row and `release UR-NNN`
  to the release row. Mixing REQ and UR tokens in one invocation resolves to their union.
- **Expansion feeds the existing loops unchanged:** reserve still only takes `pending` members
  (line 37) and reports/skips everything else with its existing per-status reason; release still only
  touches `reserved` members (line 44).
- **State the token/label precedence explicitly** in or under the mode table: in `release <token>`, a
  `UR-` + digits token resolves as an **id**, so a reservation label literally named `UR-011` is
  unreachable by that route. Without this line the two `release` rows are ambiguous, and a user who
  labelled a session after a UR would silently release the wrong set.
- The label rules are untouched — still YAML-quoted, still never interpolated into a shell command,
  still asked for when absent.

### 3. Surface updates

- `SKILL.md`: the abandon and reserve/release routing rows (lines 59–60) and their dispatch rows.
  Absorb `UR-NNN` without materially growing the file — the word budget is enforced by
  `_dev/tests/contract-regressions.sh`.
- `actions/help.md:25–27` — the abandon / reserve / release usage lines.
- `docs/cleanup-guide.md:13` — `do-work abandon REQ-NNN`.

## Constraints

- **Cite REQ-067's contract; do not restate it.** Avoiding a fourth and fifth independent copy of the
  resolution rule is the reason the contract exists.
- **Do not read the UR's `requests:` array** — scan `user_request:` frontmatter, same as everywhere.
- **Do not add a bulk-confirm bypass.** No `--yes`, no "confirm once for the batch without listing
  members." The itemized prompt is the safety property that makes UR-level cancel acceptable.
- **Do not change any per-REQ status gate.** If a member would be refused when named directly, it is
  refused when reached through a UR.
- Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md` — both are `export-ignore`d.

## Dependencies

`depends_on: [REQ-067]` — this REQ consumes the **Target ID Resolution** contract REQ-067 writes into
`actions/work-reference.md`. It also shares `SKILL.md`, `actions/help.md`, `CHANGELOG.md`, and
`actions/version.md` with REQ-067, so the two must not be co-dispatched.

## Builder Guidance

Certainty level: **Firm** on scope (both actions, confirmed with the user), on the itemized
confirmation, and on the token/label precedence. **Mixed** on how much abandon's Step 1 needs to
change structurally — if resolving ids to a target list in front of the existing per-target loop is a
clean insertion, prefer that over threading UR-awareness through the gates.

This is a `maintenance: true` REQ: prefer narrowing and citing over adding prose.

## Open Questions

- [x] *(batch-level, recorded here and in `do-work/user-requests/UR-011/input.md`)* How wide should
  UR acceptance go — `run` only, or every ID-taking action?
  → **run + abandon + reserve/release.** All three ID-taking actions that reject URs today. A
  middle option (reject, but list the UR's REQ ids in the error) was offered and not taken. At
  verify the user extended the principle symmetrically — REQ-070 covers the inverse gap in
  `actions/roadmap.md`.

## Red-Green Proof

**RED prompt/case:** `do-work abandon UR-011 superseded` and `do-work reserve UR-011 for cloud-alpha`.
Today `actions/abandon.md:32` has no defined handling for a non-`REQ-` token at all — its globs
substitute a REQ *number* into `REQ-NNN-*.md`, so the best case is `UR-011: not found` and the real
case is reader-dependent. `actions/reserve.md`'s mode table matches no row for a `UR-` token either,
so `release UR-011` falls into the free-text-label arm and matches zero reservations. In the harness: a
`_dev/tests/contract-regressions.sh` probe asserting both files' Input sections name a `UR-` token
shape and cite the Target ID Resolution contract fails today.
**Why RED now:** neither action has any notion of a UR; both key entirely on `REQ-NNN` tokens.
**GREEN when:** the probe passes, and a walkthrough shows `do-work abandon UR-011` presenting one
confirmation listing every cancellable member with its status, while a `completed` member is still
refused and a `reserved` member still demands its extra confirmation.
**Validation:** User chose this scope at capture time.

## Full Context

See `do-work/user-requests/UR-011/input.md` for complete verbatim input.

---
*Source: "executing a UR should be just as valid as executing a REQ, they are the same familly" — scope confirmed at capture as run + abandon + reserve/release.*

Think carefully before answering.

## Triage

**Route B.** Firm "what" and named files, but the clean insertion point for UR resolution (in front of each existing per-target loop, gates untouched) and the token/label precedence wording needed discovery — Route B, not a mechanical Route A edit. Depends on REQ-067's contract, which landed first (commit `1e653bc`).

## Decisions

- **D-01 (DECIDE & STATE):** Extended `write_set` to include `_dev/tests/contract-regressions.sh` — same reasoning as REQ-067: the Red-Green Proof names a harness probe, absent from the declared set. Added Input-scoped contract-citation assertions for abandon.md and reserve.md plus a reserve-Input UR-NNN check; scoped to the Input block because abandon.md already writes `archive/UR-NNN/` paths elsewhere (a file-wide grep would pass vacuously). RED confirmed before edits, GREEN after.

## Implementation Summary

**What was done:** Taught `do-work abandon` and `do-work reserve`/`release` to accept `UR-NNN`, consuming REQ-067's Target ID Resolution contract. Both actions resolve tokens in front of their existing per-target loops, so every per-REQ gate (duplicate refusal, `completed`/`completed-with-issues` refusal, `claimed`/`reserved` extra confirmations; reserve's `pending`-only capture; release's `reserved`-only touch) applies unchanged to expanded members. abandon's Step 2 confirmation now enumerates every resolved target with a total count in one prompt, and forbids a `--yes`/per-member bypass. reserve states the token-vs-label precedence explicitly (a `UR-` token resolves as an id in `release <token>`, so a label named `UR-011` is unreachable that way).

Files changed:
- `actions/abandon.md` (modified) — Input accepts `UR-NNN` + cites contract; Step 1 UR-resolution paragraph (queue/working + archived-`failed` reach, never into closed UR folders); Step 2 itemized bulk-cancel confirmation.
- `actions/reserve.md` (modified) — Input mode table gains UR forms; contract citation + token-vs-label precedence; reserve/release resolution steps expand UR before the existing loops.
- `SKILL.md` (modified) — abandon/reserve routing rows + dispatch rows.
- `actions/help.md` (modified) — abandon/reserve/release usage lines note UR.
- `docs/cleanup-guide.md` (modified) — Pass 1 mentions `do-work abandon UR-NNN`.
- `_dev/tests/contract-regressions.sh` (modified) — three RED→GREEN assertions (see D-01).
- `CHANGELOG.md`, `actions/version.md` (modified) — release bookkeeping.

## Testing

- **Red-green validation:** abandon-Input and reserve-Input contract-citation assertions + reserve-Input `UR-NNN` assertion all FAILed pre-edit (RED, verified — abandon.md's pre-existing `archive/UR-NNN/` reference is why the citation check, not a bare `UR-NNN` grep, is the meaningful probe), all pass after (GREEN, verified).
- **Regression:** full suite passes except the pre-existing `update-script-behavior.sh` baseline (untouched files). SKILL.md within the 2650-word router budget.

## Review

**Pipeline mode — Pass (self-review).** Requirements traced: abandon Input+Step1+Step2 (§1) ✓, reserve Input table + precedence + reserve/release loops (§2) ✓, surface updates (§3) ✓. Constraints held: cites the contract without restating it, `requests:` array never read, no bulk-confirm bypass (Step 2 forbids `--yes`/per-member loop), no per-REQ status gate changed (resolution sits in front of the loops), no `CLAUDE.md`/`AGENTS.md` citation.

## Lessons Learned

**What worked:** Resolving tokens in front of the existing per-target loop kept every status gate untouched — the cleanest possible insertion, exactly as the REQ's Builder Guidance predicted.
**What didn't:** A file-wide `UR-NNN` grep is a vacuous test here — abandon.md already contains `archive/UR-NNN/` folder paths, so the assertion passed before any real change. Scoping to the Input block and asserting the *contract citation* is the non-vacuous probe.
**Worth knowing:** The token-vs-label ambiguity in `release <token>` is a genuine footgun — a session labelled after a UR (`reserved_for: "UR-011"`) is now unreachable by `release UR-011` (it resolves as an id). The precedence line under reserve's mode table is the only thing that warns a user before they release the wrong set.

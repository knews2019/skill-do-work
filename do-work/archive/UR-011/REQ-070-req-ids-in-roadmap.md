---
id: REQ-070
title: REQ ids in do-work roadmap — the inverse asymmetry
status: completed
created_at: 2026-08-01T12:31:45Z
claimed_at: 2026-08-01T13:57:09Z
completed_at: 2026-08-01T13:57:09Z
route: B
user_request: UR-011
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-067]
maintenance: true
related: [REQ-067, REQ-068]
batch: ur-ids-accepted-everywhere
write_set: [actions/roadmap.md, docs/roadmap-guide.md, _dev/tests/contract-regressions.sh]
---

# REQ ids in do-work roadmap — the inverse asymmetry

## What

`actions/roadmap.md` is the mirror image of the bug this batch fixes: it accepts `UR-NNN` as a scope
token but not `REQ-NNN`. Add REQ scoping so every id-taking action in the skill accepts both
prefixes with no exceptions.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Generalize roadmap.md's Input list rather than append — add a `REQ-NNN` scope row citing REQ-067's Target ID Resolution contract, state the thin single-REQ output (and the `inspect` boundary), union for multiple tokens, keep the soft fallback. Update the report template's Scope line and the guide's example. Do NOT touch SKILL.md. RED = an Input-scoped REQ-NNN + contract-citation assertion.
- [x] **[APPLY]:** Edited `actions/roadmap.md` (Input + Output Scope line) and `docs/roadmap-guide.md`; `write_set` extended to the test file (D-01). Read-only posture untouched — only what is surveyed changed, never what is written.
- [x] **[UNIFY]:** `git diff --stat` reviewed; Markdown only. Suite passes except the pre-existing `update-script-behavior.sh` baseline. Confirmed SKILL.md unchanged and roadmap's read-only rules intact.

## Why

The user's request was symmetric — "they are the same family" — and the asymmetry runs both ways.
REQ-067 and REQ-068 fix three actions that take REQ ids and reject URs; `roadmap` is the only action
in the skill that does the opposite. Left alone, the batch would establish "both prefixes everywhere"
as the rule and immediately ship the one counterexample.

Roadmap's failure mode is also the quieter of the two. `actions/roadmap.md:35` says an unrecognized
argument **defaults to the full survey** with a note in the report — so `do-work roadmap REQ-067`
today returns a whole-queue roadmap that looks legitimate, rather than refusing. A silently wrong
scope is worse than a rejection.

## Context

`actions/roadmap.md:24-35` is a flat list of scope tokens: *(none)*, `pending`, `in-progress`,
`done`, `UR-NNN`, `since <date>`, then the unrecognized-argument fallback. The action is a read-only
survey — it never writes REQ files — so this is the lowest-risk member of the batch.

The user was offered a narrower option (keep roadmap UR-only, just stop the silent default) and
chose full symmetry instead.

## Detailed Requirements

- **Add a `REQ-NNN` scope row** to the Input list at `actions/roadmap.md:24-35`: scope the report to
  a single REQ — its status, dependency position, feasibility read, and its siblings under the same
  UR for context. Cite the **Target ID Resolution** contract (`actions/work-reference.md`, added by
  REQ-067) for the token shape rather than restating it.
- **Multiple tokens resolve to their union**, matching how the other id-taking actions read a list.
- **Keep the soft unrecognized-argument fallback at line 35.** Roadmap is read-only and advisory;
  a full survey plus a note is an acceptable response to garbage, and hard-erroring a survey action
  would be a bigger change than this REQ is buying. Only now, a recognized `REQ-NNN` no longer lands
  there.
- **Say what a single-REQ roadmap contains.** The output is thin by nature — one card's worth of
  status plus its dependency neighbourhood. State that explicitly in the Input row so a user isn't
  surprised, and so the builder doesn't pad it into a mini-report.
- `docs/roadmap-guide.md:37` currently shows only `do-work roadmap UR-014` — add the REQ form.

## Constraints

- **Do not touch `SKILL.md`.** Its roadmap routing row already reads `UR-NNN` as an example scope
  token, and the word budget is enforced by `_dev/tests/contract-regressions.sh`. If the row must
  change, generalize the existing example rather than appending a second one.
- **Do not extend this to `clarify` or `board`.** Neither takes an id argument in any form — that is
  a different design (a batch review and a UI), not an asymmetry. Explicit non-goal.
- **Roadmap stays read-only** (`actions/roadmap.md:20` — "Feasibility is a read, not a verdict").
  Scoping changes what is surveyed, never what is written.
- Cite REQ-067's contract; do not restate the token shapes.
- Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md` — both are `export-ignore`d.

## Dependencies

`depends_on: [REQ-067]` — consumes the **Target ID Resolution** contract. No file overlap with
REQ-067 or REQ-068 (`actions/roadmap.md` and `docs/roadmap-guide.md` are touched by neither), so it
may be co-dispatched with REQ-068 once REQ-067 lands.

## Builder Guidance

Certainty level: **Firm** on adding REQ scoping and on keeping the soft fallback. **Mixed** on what a
single-REQ roadmap should actually render — keep it small; resist turning it into a second `inspect`
(`actions/inspect.md` already does per-REQ detail, and duplicating it here would be the real failure).

`maintenance: true` — prefer generalizing the existing Input list over appending to it.

## Open Questions

- [x] Should `roadmap` gain REQ scoping, or just stop silently defaulting on an unrecognized token?
  → **Gain REQ scoping (full symmetry).** Chosen by the user at verify time over the narrower
  "fix the silent default only" option; the soft fallback stays for genuinely unrecognized tokens.

## Red-Green Proof

**RED prompt/case:** `do-work roadmap REQ-067`. Today the token matches no row in
`actions/roadmap.md:24-34`, so line 35 sends it to the full survey with a note — the user asked about
one REQ and got the whole queue. In the harness: a `_dev/tests/contract-regressions.sh` probe
asserting `actions/roadmap.md`'s Input names a `REQ-` scope token fails today.
**Why RED now:** the scope-token list was written when only URs grouped work; a REQ was never a
survey scope.
**GREEN when:** the probe passes and `do-work roadmap REQ-067` reports that REQ's status, dependency
position, and UR siblings — while `do-work roadmap banana` still falls through to the full survey
with its note.
**Validation:** User chose full symmetry at verify time.

## Full Context

See `do-work/user-requests/UR-011/input.md` for complete verbatim input.

---
*Source: "they are the same familly" — the inverse asymmetry, surfaced by `do-work verify-request` on UR-011 and confirmed by the user.*

Think carefully before answering.

## Triage

**Route B.** Named file and firm scope, but rendering the single-REQ survey required a judgment call — keeping it thin and not duplicating `inspect` — so Route B, not a mechanical Route A edit. Depends on REQ-067's contract (landed, commit `1e653bc`).

## Decisions

- **D-01 (DECIDE & STATE):** Extended `write_set` to include `_dev/tests/contract-regressions.sh` for the TDD proof (same as REQ-067/068). Added Input-scoped `REQ-NNN` + contract-citation assertions; RED confirmed before edits, GREEN after.

## Implementation Summary

**What was done:** Added a `REQ-NNN` scope token to `do-work roadmap`, the inverse of the asymmetry this batch fixed (roadmap took `UR-NNN` but not `REQ-NNN`, so `do-work roadmap REQ-067` silently returned a whole-queue survey). The Input list now cites REQ-067's Target ID Resolution contract for token shapes, adds the `REQ-NNN` row (single-REQ survey: status, dependency position, feasibility read, UR siblings — explicitly thin, with the `inspect` boundary called out), states that multiple id tokens resolve to their union, and keeps the soft unrecognized-argument fallback (only a recognized `REQ-`/`UR-` id is scoped; `banana` still falls through). The Output Format's Scope line gained `REQ-NNN`.

Files changed:
- `actions/roadmap.md` (modified) — Input list (contract citation + `REQ-NNN` row + union + fallback clarification) and the Output `**Scope:**` enumeration.
- `docs/roadmap-guide.md` (modified) — added the `do-work roadmap REQ-070` example.
- `_dev/tests/contract-regressions.sh` (modified) — two RED→GREEN Input-scoped assertions (see D-01).
- `CHANGELOG.md`, `actions/version.md` (modified) — release bookkeeping.

## Testing

- **Red-green validation:** the roadmap-Input `REQ-NNN` assertion and the contract-citation assertion both FAILed pre-edit (RED, verified), both pass after (GREEN, verified).
- **Regression:** full suite passes except the pre-existing `update-script-behavior.sh` baseline. Confirmed `SKILL.md` unchanged (constraint) and roadmap's read-only rules intact.

## Review

**Pipeline mode — Pass (self-review).** Requirements traced: `REQ-NNN` Input row ✓, contract citation (not restated) ✓, union of multiple tokens ✓, soft fallback kept ✓, single-REQ contents stated with the `inspect` boundary ✓, guide example ✓. Constraints held: `SKILL.md` untouched, roadmap stays read-only (only *what is surveyed* changed), no `CLAUDE.md`/`AGENTS.md` citation.

## Lessons Learned

**What worked:** Generalizing the Input list (one intro line citing the contract, one new row) kept the change small and symmetric with the UR row — exactly the maintenance posture the REQ asked for.
**What didn't:** n/a.
**Worth knowing:** The single-REQ roadmap deliberately overlaps `do-work inspect` in spirit; the guard against turning it into a second `inspect` is the explicit "thin by nature… not a mini-report" wording in the Input row. A future edit that fattens it would re-introduce exactly the duplication this REQ warned against.

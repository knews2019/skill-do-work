---
id: REQ-069
title: REQ ids in do-work roadmap — the inverse asymmetry
status: pending
created_at: 2026-08-01T12:31:45Z
user_request: UR-011
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-067]
maintenance: true
related: [REQ-067, REQ-068]
batch: ur-ids-accepted-everywhere
write_set: [actions/roadmap.md, docs/roadmap-guide.md]
---

# REQ ids in do-work roadmap — the inverse asymmetry

## What

`actions/roadmap.md` is the mirror image of the bug this batch fixes: it accepts `UR-NNN` as a scope
token but not `REQ-NNN`. Add REQ scoping so every id-taking action in the skill accepts both
prefixes with no exceptions.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

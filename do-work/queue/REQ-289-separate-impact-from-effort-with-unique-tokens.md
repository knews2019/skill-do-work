---
id: REQ-289
title: Separate impact from effort, with unique greppable tokens on both axes
status: pending
created_at: 2026-08-19T14:33:51Z
user_request: UR-060
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
related: [REQ-290]
batch: impact-effort-split
write_set:
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/review-work.md
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
- skills/do-work/actions/work.md
- skills/do-work/actions/estimate-reference.md
- skills/do-work-board/tools/queue-kanban/model.go
- skills/do-work-board/tools/queue-kanban/generate.go
- skills/do-work-board/tools/queue-kanban/web/board-cards.js
- skills/do-work-board/tools/queue-kanban/web/board-detail.js
- _dev/tests/contract-regressions.sh
---

# Separate Impact from Effort, With Unique Greppable Tokens on Both Axes

## What

`effort_estimate` has two writers with two different meanings. Capture sets it as a size judgment;
review MUST-stamp it from an impact gate. Split the two axes into two fields, and give every value
on both axes a token that is unique repo-wide and findable by plain-text search.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user wants to know a REQ's **impact** — whether anyone would ever notice the work — so they can
stop implementing REQs whose impact is negligible. That test already runs
(`review-work.md:340-341`), but its verdict is discarded into a field that means size, so the signal
the user wants never reaches them in a usable form.

It also causes damage today. `work.md:263` short-circuits the estimator to its floor value on
`effort_estimate: trivial`. A finding nobody would notice but that takes three hours to fix is
stamped `trivial` by the impact gate and then forecast at five minutes.

## Context

- `capture-reference.md:24` — capture sets `effort_estimate` "when the request is clearly a small
  mechanical fix". A size judgment.
- `review-work.md:357` — automatic follow-up creation MUST stamp it from the `gate:` token.
- `review-work.md:340-341` — the gate's two questions are "Would any user or developer actually
  notice this issue in real use?" and "does fixing it establish or change a rule that applies in
  several places?". Both are impact questions.
- `work-reference.md:674` — writes "Severity and effort are different axes" inside the paragraph
  that wires one to the other.

## Detailed Requirements

### The `impact:` field and its vocabulary

Add `impact:` to the REQ frontmatter schema. Four values, all prefixed so `grep 'impact-'` finds
every use:

| Token | Means | Replaces |
|---|---|---|
| `impact-critical` | Security, data loss, or a broken production path. Pierces the consent gate and auto-queues at any depth. | `[critical]` |
| `impact-user-visible` | A user or developer would notice it in real use. | `gate: user-visible` |
| `impact-rule-change` | Nobody notices it, but fixing it sets or changes a rule that applies in several places. | `gate: rule-change` |
| `impact-negligible` | Neither. The user's stop signal. | `gate: trivial`, `[low]` |

The discovered-task `[normal]` maps to no single token — it splits by the gate's two questions into
`impact-user-visible` or `impact-rule-change`. That is correct; `[normal]` was doing two jobs.

### The `effort_estimate:` rename

`effort-mechanical` (was `trivial`) and `effort-substantive` (was `normal`, and still the
absent-default). Add read-only aliases `trivial` -> `effort-mechanical` and
`normal` -> `effort-substantive` so every existing REQ stays valid unchanged; the alias mechanism
already exists at `capture-reference.md:122` ("aliases are read-only, never propagated on write").
The `effort_estimate` row's current "no aliases — closed two-value enum, deliberately" changes with
it.

**Do not rewrite existing REQ files.** The aliases carry them. New writes emit the new tokens.

### Why the tokens change at all

`trivial` currently matches 104 lines under `skills/` and `normal` matches 520. Neither is findable
by plain-text search, and both are used on two different axes — the conflation in its most literal
form. All six proposed tokens currently match zero files. This is CLAUDE.md's Naming for Reach rule
applied to enum values.

### Rewiring

- `review-work.md:343` records the impact token directly on the finding line. The `gate:` name and
  its separate three-word vocabulary retire, so no translation table sits between the finding line
  and the frontmatter.
- `review-work.md:357` and its follow-up template stamp `impact:`, not `effort_estimate`.
- `review-work.md:352` (sweep REQs) and `review-work.md:539` (Verification Checklist) follow.
- `work-reference.md:674` (Discovered Tasks Classification) stops stamping effort from the impact
  gate. It stamps `impact:` from the gate and judges effort separately. Its "Severity and effort are
  different axes" sentence becomes true instead of self-contradicting. Its restatements at
  `work.md:497`, `:503`, `:507`, and `:564` follow.
- `work.md:263`'s trivial short-circuit stays keyed on `effort_estimate`, which now genuinely means
  effort — it becomes correct rather than accidentally right. `estimate-reference.md:63,69` follows.
- `capture-reference.md:24` and capture's Step 1 assessments gain the impact judgment described
  below.

### Capture behavior — judged by default, asked above a threshold

Every REQ carries an `impact:` value; absent must not be the common case.

- Capture applies the gate's two questions itself and writes a value.
- It asks the user **only** when the two questions disagree or it cannot judge. `clear-questions.md`
  governs that question when it fires.
- This adds no friction to the common capture and leaves the field never silently absent.

### Board

`work-reference.md:137` requires the board's Go parser and the schema line to change in the same
commit. Give `impact:` the same present-value-only normalize-and-warn treatment `effort_estimate`
already has (`model.go:184,194,712-723`, enum table `model.go:997`, warning `model.go:1146`,
export `generate.go:152-158,549-551`), and render it. `board-cards.js:164-188` is the existing chip
precedent. Read `_dev/primes/prime-kanban-board.md` before touching the tool.

## Constraints

- `maintenance: true` — load `crew-members/maintenance.md`. The net move is consolidation: three
  impact vocabularies become one, and one overloaded field becomes two honest ones. Prefer retiring
  to adding wherever the choice appears.
- **Name REQ-228 in the implementation and say why it does not bind.** It recorded "No new
  frontmatter field. Not on REQs, not on URs. `effort_estimate` stays a two-value triage bit." That
  decision was about timeline projection, not about the impact/effort conflation. If this is not
  written down, the next reviewer re-litigates it.
- Do not grow `effort_estimate` toward t-shirt sizes. Its two-value posture is deliberate and
  re-affirmed three times (UR-027, REQ-125, REQ-228); this REQ renames its values, it does not widen
  the enum.
- The board tool's write-surface count in CLAUDE.md is unaffected — nothing here adds a write
  surface.

## Dependencies

REQ-290 depends on this REQ. Nothing depends on REQ-290.

## Red-Green Proof

**RED prompt/case:** Run a `do-work review-work` pass that produces one finding answering "no" to
both gate questions but that would take hours to fix. Then inspect the follow-up REQ it creates and
the estimate the work loop assigns it.

**Why RED now:** The finding is stamped `effort_estimate: trivial` from the impact gate, and
`work.md:263` short-circuits its forecast to the floor value. The REQ claims to be a small
mechanical fix on the strength of an impact judgment. Nothing anywhere records that nobody would
notice it, in a field the user can filter on.

**GREEN when:** The follow-up carries `impact: impact-negligible` and an `effort_estimate` judged as
effort. `grep -rn 'impact-' do-work/queue/` returns every REQ's impact verdict. No file stamps
`effort_estimate` from a gate token. A REQ still carrying the literal `effort_estimate: trivial`
reads as `effort-mechanical` through the alias and remains valid.

Lock-in checks: no action file stamps `effort_estimate` from an impact judgment; the six tokens are
each unique repo-wide; the legacy aliases resolve; the board parser and the schema row are in the
same commit. `bash _dev/tests/maintainer-verify.sh` exits 0.

**Validation:** User confirmed — vocabulary shape, surface, and capture behavior all chosen via the
ask tool during capture.

## Full Context

See `do-work/user-requests/UR-060/input.md` for complete verbatim input.

---
id: REQ-289
title: Separate impact from effort, with unique greppable tokens on both axes
status: completed
completed_at: 2026-08-19T15:48:05Z
commit: 2ea7be5
claimed_at: 2026-08-19T14:41:49Z
created_at: 2026-08-19T14:33:51Z
user_request: UR-060
domain: general
route: C
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
related: [REQ-290]
batch: impact-effort-split
estimate:
  p50_active_minutes: 60
  confidence: low
  calculated_at: 2026-08-19T14:41:49Z
  basis:
    - Route C
    - 11-file write set
    - 4 subsystems involved
    - 4 acceptance criteria
    - persistence changes
    - cross-route regression gates
    - full-suite verification
write_set:
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/review-work.md
- skills/do-work/actions/work.md
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
- skills/do-work/actions/estimate-reference.md
- skills/do-work/docs/review-work-guide.md
- skills/do-work-board/tools/queue-kanban/model.go
- skills/do-work-board/tools/queue-kanban/generate.go
- skills/do-work-board/tools/queue-kanban/timeline.go
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/model_test.go
- skills/do-work-board/tools/queue-kanban/timeline_test.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/web/board-cards.js
- skills/do-work-board/tools/queue-kanban/web/board-detail.js
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-toolbox/actions/maintainability-audit.md
- skills/do-work-toolbox/actions/maintainability-audit-reference.md
- _dev/tests/contract-regressions.sh
kb_status: promoted
kb_entry: REQ-289-separate-impact-from-effort-with-unique.md
---

# Separate Impact from Effort, With Unique Greppable Tokens on Both Axes

## What

`effort_estimate` has two writers with two different meanings. Capture sets it as a size judgment;
review MUST-stamp it from an impact gate. Split the two axes into two fields, and give every value
on both axes a token that is unique repo-wide and findable by plain-text search.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both prime files, the five crew rules, and the approved plan. Approach, in the plan's order: (1) RED — four lock-in checks (A no action file stamps `effort_estimate` from an impact judgment; B the six tokens are defined and the bare axis words never appear un-prefixed; C the legacy `trivial`/`normal` aliases resolve on both the schema row and `model.go`; D the schema row's `impact-*` set and `model.go`'s `canonicalValues` agree both directions) plus six new entries on the retired-vocabulary loop, all confirmed failing first. (2) `work-reference.md` schema row + Schema Read Contract rows and `model.go`/`timeline.go` in one edit, with named constants so no bare `"trivial"` literal survives. (3) retire `gate:` across `review-work.md` — the finding line carries the bare impact token, which IS the frontmatter value. (4) collapse the `[critical]`/`[normal]`/`[low]` ladder and delete the orchestrator-classifies rubric in `work-reference.md`. (5)-(7) `work.md`, `estimate-reference.md`, capture restatements. (8) board render: Go export, card chip, drawer row, CSS, timeline labels. (9) toolbox's third `gate:` consumer. Verification is the plan's Validation block, ending at `bash _dev/tests/maintainer-verify.sh` exit 0.
- [x] **[APPLY]:** Implemented the approved plan's eleven tasks in order, RED first. Scope stayed inside the declared write set except for one file, `skills/do-work-board/tools/queue-kanban/generate_test.go`, which the rename made red and which is recorded as D-05 rather than written silently.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file (23 paths, two of them not this REQ's work: `do-work/CHECKPOINT.md` is the orchestrator's claim record and this REQ file itself). Linters: `gofmt -l .` empty, `go vet ./...` clean, `go test -count=1 ./...` ok, `shellcheck --severity=warning _dev/tests/contract-regressions.sh` exit 0. `git diff | grep -E 'console\.log|debugger|TODO|FIXME|fmt\.Println'` returns nothing. Per-file verification is listed in the Implementation Summary.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** Eleven declared files across four subsystems (core action files, the board's Go parser and its JS renderers, and the contract-regression suite), with a schema change that `work-reference.md:137` requires to land in the same commit as the Go parser. Many interlocking requirements — two vocabularies, an alias mechanism, and rewiring six documented call sites.

**Planning:** Required

## Plan

Route C plan produced by the Plan agent, verified by the orchestrator against the working tree
before acceptance (every cited line re-read; the two "silent breakage" claims reproduced by hand).

**Execution order:** (1) RED — four lock-in checks into `_dev/tests/contract-regressions.sh`;
(2) schema row + Go parser in one edit (`work-reference.md:137,236` + `model.go`), with named
constants so no bare `"trivial"` literal survives; (3) retire `gate:` across `review-work.md`;
(4) collapse the `[critical]`/`[normal]`/`[low]` ladder in `work-reference.md`'s Discovered Tasks
Classification onto the impact vocabulary; (5) `work.md` restatements; (6) `estimate-reference.md`;
(7) capture's impact assessment; (8) board render (`generate.go`, `board-cards.js`,
`board-detail.js`, `board.css`, `board-timeline.js`).

**The net move is subtraction.** One new field (`impact:`, four values) is paid for by retiring two
three-token vocabularies (`gate:` and `[critical]`/`[normal]`/`[low]`), the orchestrator-classifies
rubric at `work-reference.md:668-671`, and two `effort_estimate` write-paths. Every other site is a
one-for-one substitution.

**Plan validation (orchestrator):**
1. *Requirement coverage* — every Detailed Requirement maps to a task: the field and vocabulary
   (T2), the rename and aliases (T2), retiring `gate:` and the severity words (T3, T4), unwiring
   effort from the impact gate (T3 `:357`, T4 `:674`), capture behavior (T7), board (T2e, T8),
   REQ-228 (T2b). No uncovered requirement.
2. *No orphan tasks* — every task traces to a requirement or to a blocker the requirement creates.
3. *Scope sanity* — nine tasks, past the 3-task quality band. Flagged deliberately and **not**
   split: the schema row, the Go parser, and the tests are a single atomic unit
   (`work-reference.md:137` requires the parser and the schema line in one commit, and
   `maintainer-verify.sh` runs `go test ./...`), so a split would ship a knowingly red tree.

*Generated by Plan agent; validated by the orchestrator.*

## Exploration

**No separate Explore agent was spawned** — see D-01. The Plan agent's pass (56 tool uses) already
returned the file-level findings Step 5 asks for, with exact line numbers for every site and the
existing precedents to copy.

**Key files and the precedent each supplies:**
- `model.go:996-1001` — `schemaReadContractFields`, the present-value-only normalize-and-warn table.
  `domain` supplies the alias-map precedent; `effort_estimate` supplies the closed-enum precedent.
- `model.go:712-724`, `:775-777`, `:1146` — the three-site read pattern (normalize, assign
  `Original*`/`*Unrecognized`, register the warning). `:1134-1135` states the rule for adding a field.
- `board-cards.js:164-188` — the chip precedent, including the never-silently-drop leg where an
  unrecognized value chips with an `invalid` flag.
- `contract-regressions.sh:1600-1650` — the retired-vocabulary loop pattern the new checks extend.
- `review-work.md:343` — the bare-id provenance precedent (`(UR-489/UR-027)`) for naming REQ-228 in
  shipped prose without citing a `do-work/archive/` path the consumer will not have.

**Concerns carried into implementation:** the silent-breakage risk at `timeline.go:281,307-310`
(bare `"trivial"` literals that a rename leaves compiling and wrong), and the `impact:` default,
which must be `impact-user-visible` — defaulting to `impact-negligible` would arm REQ-290's filter
to skip every REQ predating the field.

## Scope

**Files I will touch:**

*Core action files*
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/actions/capture-reference.md`
- `skills/do-work/actions/estimate-reference.md`
- `skills/do-work/docs/review-work-guide.md`

*Board tool*
- `skills/do-work-board/tools/queue-kanban/model.go`
- `skills/do-work-board/tools/queue-kanban/generate.go`
- `skills/do-work-board/tools/queue-kanban/timeline.go`
- `skills/do-work-board/tools/queue-kanban/durations.go`
- `skills/do-work-board/tools/queue-kanban/model_test.go`
- `skills/do-work-board/tools/queue-kanban/timeline_test.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (added during implementation — D-05)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js`
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`
- `skills/do-work-board/tools/queue-kanban/web/board.css`

*Toolbox*
- `skills/do-work-toolbox/actions/maintainability-audit.md`
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md`

*Tests*
- `_dev/tests/contract-regressions.sh`

**Integrator-owned, not the builder's** (Step 9 / `CLAUDE.md` § Before Every Commit): `VERSION`,
`skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, and its byte-identical
mirror `skills/do-work/CHANGELOG.md`.

**Acceptance criteria (restated from the REQ's Red-Green Proof):**
1. A follow-up created by review carries `impact:` set to the finding's recorded token, and an
   `effort_estimate` judged as effort — never derived from the impact verdict.
2. `grep -rn 'impact-' do-work/queue/` returns every REQ's impact verdict.
3. No file stamps `effort_estimate` from a gate/impact token.
4. A REQ still carrying the literal `effort_estimate: trivial` reads as `effort-mechanical` through
   the read-only alias and remains valid; no existing REQ file is rewritten.
5. All six tokens are unique repo-wide by plain-text search.
6. The board parser and the schema row change in the same commit, held by an agreement assertion.
7. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Implementation Summary

**Files changed:**

*Tests (written first — RED before any implementation edit)*
- `_dev/tests/contract-regressions.sh` (modified) — four lock-in checks plus six retired-vocabulary entries; +189 lines.

*Core schema and its parser (one atomic unit)*
- `skills/do-work/actions/work-reference.md` (modified) — new `impact:` schema row with the REQ-228 parenthetical; `effort_estimate` renamed and its MUST-stamp-from-the-gate clause deleted; new `impact` Schema Read Contract row; `effort_estimate` row gains the two read-only aliases; Discovered Tasks Classification collapsed onto the four tokens with the orchestrator-classifies rubric deleted.
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — six named constants, the `Impact`/`OriginalImpact`/`ImpactUnrecognized` trio, present-value-only read, `"impact"` contract entry, `effort_estimate` aliases, warning registration.
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified) — the silent-breakage site: both bare `"trivial"` literals and both returned bucket names now reference the constants.
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified) — bucket comment onto the constants.

*Action files*
- `skills/do-work/actions/review-work.md` (modified) — `gate:` retired entirely; the finding line carries the bare impact token, which is the frontmatter value; follow-up template swaps `effort_estimate:` for `impact:`; the critical pierce keys on `impact-critical`.
- `skills/do-work/actions/work.md` (modified) — six restatements; the short-circuit renamed to mechanical-effort.
- `skills/do-work/actions/capture.md` (modified) — Step 1 impact assessment (judge by default, ask only on disagreement), a checklist line, and one plain-English `user-visible` reworded (D-04).
- `skills/do-work/actions/capture-reference.md` (modified) — uncommented `impact:` in the Simple REQ template, size-only `effort_estimate` gloss, `impact` added to the normalize-and-warn list.
- `skills/do-work/actions/estimate-reference.md` (modified) — four renames; `--trivial` kept verbatim with a note that the flag names the estimator's floor mode, not the schema token.
- `skills/do-work/docs/review-work-guide.md` (modified) — shipped prose rewritten onto the impact field.
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modified) — the third `gate:` consumer retired.
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (modified) — same, plus the numeric `Impact: [1-5]` score disambiguated in prose (D-06) and a cross-package citation path corrected (D-08).

*Board render*
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — the three impact fields exported.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified) — impact chip, including the never-silently-drop unrecognized leg.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified) — Impact meta row above Effort estimate.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — `.badge-impact`, neutral for every value (D-07).
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified) — user-facing labels to mechanical/substantive.

*Go tests*
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified) — canonical pair updated, alias-resolution cases added, new `TestImpactFieldFollowsPresentValueOnlyContract`.
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified) — fixtures onto the constants.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — one assertion string; the one write outside the originally declared set (D-05).

**What was done:** Split the two axes `effort_estimate` was carrying into two fields. Added `impact:` with four prefix-unique tokens (`impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible`) and renamed `effort_estimate`'s values to `effort-mechanical`/`effort-substantive`, with read-only `trivial`/`normal` aliases so no existing REQ file needed rewriting. Retired the `gate:` field and its separate three-word vocabulary, along with the `[critical]`/`[normal]`/`[low]` discovered-task severity words, so the token recorded on a review finding is now the same string that lands in the follow-up's frontmatter — no translation table between them. Deleted every site that derived `effort_estimate` from an impact judgment. The board's Go parser, its schema row, and the tests moved in the same commit, held by a both-directions agreement assertion.

**REQ-228 does not bind.** It ruled "No new frontmatter field. Not on REQs, not on URs. `effort_estimate` stays a two-value triage bit" — read in place, bounded by the constraint directly above it, that was about a *derived forecast* not persisting, because persisting it would make the board a fourth write surface. `impact:` is not a forecast and nothing derives it: it is an authored judgment written once at capture or follow-up creation, which the board only reads. REQ-228's other half is honored exactly — this REQ renames `effort_estimate`'s two values and does not widen the enum.

## Qualification

**Passed — 22 project files verified, 7 acceptance criteria traced, P-A-U confirmed.**

Mechanical (`tools/checks/qualify.sh`): OK — every manifest file shows in the diff, the P-A-U box
audit is consistent with the diff, no debug artifacts, and the summary lists project files rather
than only `do-work/` paths.

Judgment checks, run by the orchestrator against git state rather than the builder's report:

- **Substantive (check 2).** Every changed file carries real edits. The central move is a deletion
  and it is present in the diff: `work-reference.md`'s "Automatic follow-up creation MUST stamp it
  from the review gate's recorded disposition" clause is gone, and the orchestrator-classifies
  critical/normal/low rubric is gone. The "Severity and effort are different axes" sentence now
  reads "Impact and effort are different axes" in a paragraph that no longer wires one to the other.
- **Requirements traced (check 3).** All seven acceptance criteria confirmed — see `## Testing`.
- **Flowing (check 6).** Verified end to end through the real board reader rather than by
  inspection: a purpose-built three-REQ fixture was parsed by a freshly built `queue-kanban` and its
  exported `board-data.js` read back. `impact: impact-negligible` round-trips verbatim; a legacy REQ
  carrying `effort_estimate: trivial` resolves to `effort-mechanical` with `originalEffortEstimate`
  preserved; a typo'd `impact: kinda-important` resolves to `impact-user-visible` with
  `impactUnrecognized: true` and the original preserved; an absent `impact` stays absent rather than
  being fabricated. The data path is live, not stubbed.

**The lock-in checks were mutation-tested, not taken on trust.** Each of the four was re-run against
a deliberately reintroduced instance of the exact defect it claims to pin, and each failed with a
message naming the real failure:

| Mutation | Check that caught it |
|---|---|
| Dropped `"trivial": effortMechanical` from `model.go`'s alias map | C — "the board is the reader that makes an archived REQ still resolve" |
| Dropped `impactRuleChange` from `model.go`'s `canonicalValues` | D — named which side held the extra token |
| Appended `— gate: trivial` to `review-work.md` | retired-vocabulary loop |
| Appended a "stamping `effort_estimate` from the impact token" sentence | A — named the offending line |

All three mutated files were restored from byte-verified backups and re-diffed to empty.

## Testing

**Test commands run** (the repo's canonical gate, per `CLAUDE.md` § Verify):

- `bash _dev/tests/maintainer-verify.sh` — **exit 0**, run unpiped. Exit code is the only proof; the
  run was never piped through `tail`, which would have hidden the status.
- `cd skills/do-work-board/tools/queue-kanban && gofmt -l . && go vet ./... && go test -count=1 ./...` — clean, exit 0.
- `bash _dev/tests/contract-regressions.sh` — exit 0.
- `bash _dev/tests/shipped-package-reference-contract.sh` — exit 0.
- `shellcheck --severity=warning _dev/tests/contract-regressions.sh` — exit 0.

**Baseline comparison:** Step 5.75 recorded a green baseline (clean tree outside `do-work/`, full
suite passing), so every check above is a true pass rather than a pre-existing failure excluded from
the gate. No new regressions.

**Red-green validation.** `tdd: true`, and the RED came first: the four lock-in checks were written
and confirmed failing (`contract-regressions.sh` exit 1, 18 FAIL lines) before any implementation
edit. Traced back to the REQ's `## Red-Green Proof`:

- *RED as captured* — a finding answering "no" to both gate questions is stamped `effort_estimate:
  trivial` from an impact judgment, and `work.md`'s short-circuit then forecasts it at the floor.
  Check A named exactly that defect at fourteen sites, `review-work.md:357` and
  `work-reference.md:674` among them — the two the REQ itself cites.
- *GREEN as captured* — a follow-up now carries `impact:` set to the finding's recorded token and an
  `effort_estimate` judged as effort; `grep -rn 'impact-'` returns every verdict; no file stamps
  `effort_estimate` from a gate token; and a REQ still carrying the literal `effort_estimate:
  trivial` reads as `effort-mechanical` through the alias — proven on a live parse, not asserted.

**Acceptance criteria, each with its evidence:**

1. Follow-up carries `impact:` from the recorded token, effort judged separately — `review-work.md`
   template and `work-reference.md` Discovered Tasks diff.
2. `grep -rn 'impact-'` surfaces every verdict — all six tokens resolve to ≥ 9 files each.
3. No file stamps `effort_estimate` from a gate token — Check A, mutation-tested.
4. Legacy REQs valid unchanged — 92 archived REQs still carry the literal token, `git status
   do-work/archive` is empty, and the alias resolves on a live parse.
5. Six tokens unique repo-wide — Check B, both legs; no un-prefixed axis word survives.
6. Parser and schema row in the same commit — Check D, a both-directions agreement assertion.
7. `maintainer-verify.sh` exits 0 — confirmed.

## Review

**Overall: 88%** | 2026-08-19T15:48:05Z

| Dimension | Score |
|-----------|-------|
| Requirements | 97% |
| Code Quality | 88% |
| Test Adequacy | 72% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 Check A scans a strict subset of the files it claims (`docs/`, `crew-members/`, sibling `actions/` excluded) — impact-rule-change → folded into REQ-293
- F2 Check A's proximity rule catches 1 of 6 realistic re-drift phrasings; the verb set is one literal and any period breaks the window — impact-rule-change → folded into REQ-293
- F3 The retired-vocabulary loop pins bold markup, not the token; the backticked form the tree actually carried passes clean — impact-rule-change → folded into REQ-293
- F4 No check holds the `impact-user-visible` default, the one property REQ-290 depends on — impact-user-visible → folded into REQ-293
- F5 The impact chip has no test coverage of any kind; a ten-line precedent sits in the same file — impact-rule-change → folded into REQ-293
- F6 Capture's impact guard is one-directional; three forces push every REQ to `impact-user-visible` — impact-user-visible → REQ-294 created
- F7 Three bare "impact" wordings left in the toolbox audit where D-06's disambiguation was not applied — impact-rule-change → REQ-295 created

**Minor findings:** 4 (report only) — five non-token `impact-*` compounds introduced by the diff (M1); acceptance criterion 2 recorded as traced when it is deferred by design (M2); `board.md`'s parsed-field list not updated for the new chip, widening a gap `effort_estimate` already had (M3); "92 archived REQs" stated as measured where the line-anchored count is 87 (M4).

**Nit findings:** 1 — `estimate-reference.md` added a clarifying sentence where a rename would have deleted the divergence; recorded as a discovered task rather than buried.

**Acceptance:** Pass — `maintainer-verify.sh` exit 0 unpiped, the Go lane clean, and the alias/default/unrecognized paths proven on a live five-fixture parse rather than by inspection.

**Restatement sweep:** applied repo-wide, outside the declared Scope, across all four skill packages plus `decisions/`, `_dev/`, `tools/`, `README.md`, and `CLAUDE.md`. Zero stale restatements of the retired vocabulary survive in any live file; archive and CHANGELOG hits are historical records and correctly left alone. The only sweep hits were F7 and M3.

**REQ-228 argument verified true.** The reviewer read the archived REQ and confirmed the shipped parenthetical describes what the surrounding text actually says: REQ-228's constraint sits directly under "A forecast is a derived display; the moment it persists anywhere it becomes a fourth write surface", and its requirement 9 reads "Nothing is written. The projection is derived at render time." The claim that it does not bind a field nothing derives is sound, not a plausible-sounding rationalization.

**Suggested testing:** 4 items (carried into REQ-293 and REQ-294 acceptance criteria)
**Follow-ups created:** REQ-293 (sweep, `sweep_key: impact-effort-lockin-checks-underpin`), REQ-294, REQ-295; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Mutation-testing the lock-in checks instead of accepting "the suite is green" — three of the four
  proved real on the first try, and the exercise is what made the review's counter-finding legible
  rather than arguable.
- Verifying backwards compatibility through the actual parser on purpose-built fixtures. The claim
  "the aliases carry legacy REQs" is unfalsifiable by grep and trivial to settle with a real parse.
- Extending the write set at Step 5.5 from the *planner's* survey rather than discovering the gaps
  mid-build. Nine of the twenty-one files were added before the builder started; two of them would
  have failed the build and one would have shipped a silent bug.

**What didn't:**
- **The orchestrator's mutation test for Check A was self-confirming and it took the reviewer to see
  it.** The mutation used the verb "stamping" — the single literal the check greps. A green mutation
  test proved only that the check catches its own phrasing. A mutation must be written in words the
  check's author did not choose, or it tests the tester's vocabulary rather than the property.
- The plan's literal wording for Check A was unimplementable and the builder had to redesign it
  mid-build (D-03). Specifying a check by its *regex* rather than by the property it must hold pushed
  the design decision to the wrong moment.

**Worth knowing:**
- `timeline.go` compared `effort_estimate` against a bare `"trivial"` string, and `timeline_test.go`
  built its fixtures from the same literal. A rename would have left both compiling, both passing,
  and every REQ silently bucketed into the substantive median. When renaming an enum, grep the
  *value* across the whole module, not the constant name — the tests are as likely to hold the
  literal as the code is.
- Three separate vocabularies described one axis before this REQ (`gate:`, the `[critical]`/
  `[normal]`/`[low]` ladder, and `effort_estimate`'s double duty). The third consumer
  (`maintainability-audit-reference.md`) already stated this REQ's thesis in its own words — "trivial
  as a gate value routes the review flow and never doubles as a severity" — while carrying the
  conflation. A file arguing against a problem is not evidence it avoided it.
- The `impact:` default must stay `impact-user-visible`. Absence must never be mistakable for the
  user's stop signal, because REQ-290's filter acts on that signal.

## Orientation

Now a REQ says whether anyone would ever notice it, separately from how big it is. `impact:` is a new
frontmatter field with four prefix-unique tokens, living in the request schema
(`skills/do-work/actions/work-reference.md`) and read by the queue pipeline, the review flow, and the
board's parser; `effort_estimate` keeps its two-value shape but now means size only, with read-only
aliases so nothing already written had to change. The `gate:` field and the `[critical]`/`[normal]`/
`[low]` ladder are retired — three vocabularies for one axis became one.

**[MAP CHANGED]** — this renames a concept and adds a schema field, so the queue's data model is not
what it was: readers that knew `effort_estimate` as "the triage bit the review gate stamps" now need
to know it as size, with impact as its own axis. The user-facing payoff is deliberately not here yet:
REQ-290 puts the token in REQ titles and adds `do-work run --skip-impact-negligible`, which is what
turns the field into the stop signal UR-060 asked for.

**Prime staleness spot-check:** `_dev/primes/prime-action-files.md` and
`_dev/primes/prime-kanban-board.md` both re-read; every path each cites still exists, and neither
carries a statement this change falsifies. `prime-kanban-board.md`'s parser-lock-step convention is
reinforced rather than contradicted — this REQ moved the schema row and the Go parser in one commit
and added an agreement assertion that fails whichever side drifts alone.

## Discovered Tasks

- **impact-user-visible** No JavaScript behavior probe covers the card chips. The suite proves the
  `impact` field end to end through the Go parser and the generated `board-data.js`, and it
  parse-checks the assembled client, but nothing executes the chip-rendering branch in
  `web/board-cards.js` — for the new impact chip **or** for the effort chip that has shipped for
  several versions. A typo in either gate would pass the whole suite. `generate_test.go` already has
  the node-probe harness (`runJavaScriptBehaviorProbe`) this would use.
- **impact-negligible** The board's timeline projection still carries the pre-rename vocabulary in
  its internal names: `trivialSamples` / `normalSamples` / `trivialMinutes` / `normalMinutes` JSON
  keys in `generate.go`, and `TrivialMedianMinutes` / `NormalMedianMinutes` in `timeline.go`. Named
  an explicit non-goal by the approved plan (internal payload names with no schema meaning) and left
  alone deliberately; the user-facing labels they feed were renamed. Worth a sweep only if the
  divergence between an internal name and its rendered label starts costing reader time.
- **impact-negligible** `tools/estimate-p50.sh --trivial` and its `- trivial short-circuit` basis
  string still spell the retired word. Pinned by `_dev/tests/p50-estimator-determinism.sh:77-80`,
  named a non-goal by the plan, and now explicitly documented at `estimate-reference.md` as naming
  the estimator's floor mode rather than the schema token — so the divergence is stated rather than
  silent.

**Disposition (Step 8):** The first discovery (no JavaScript behavior probe covers the card chips)
was independently found by the review as F5 and is queued as an instance of the REQ-293 sweep — not
duplicated here. The two `impact-negligible` discoveries were both deliberate non-goals of the
approved plan, so they go through the Open-Questions consent flow as REQ-296
(`status: pending-answers`) rather than being auto-queued.

## Decisions

- **D-01**: No separate Explore agent for Step 5. Reasoning: the Plan agent's pass already returned
  every file, line number, and existing precedent Step 5 exists to discover; a second agent would
  re-find files the plan cites by line. DECIDE & STATE — reversible and low-reach (if implementation
  hits an unmapped area, exploration can still run).
- **D-02**: Extend the write set from 11 files to 21. Reasoning: the planner found five files whose
  omission breaks the build or ships a false statement, and four more carrying drifted comments or
  user-facing labels. Two are hard blockers verified by hand — `timeline.go:281,307-310` compares
  against a bare `"trivial"` literal, so the rename leaves it **compiling and silently wrong** (every
  REQ falls into the substantive median, and `timeline_test.go` keeps passing because its fixtures
  use the same literals); `model_test.go` asserts the retired pair at seven sites, so
  `maintainer-verify.sh`'s `go test ./...` fails until it moves. `maintainability-audit-reference.md`
  is the third consumer of the `gate:` vocabulary — leaving it orphans a retired vocabulary in a
  shipped file, and narrowing the retirement check to make it pass is how the vocabulary survives.
  ESCALATE — the write set is the collision guard for parallel work, so widening it is a scope
  judgment rather than a builder's call.
  **Value:** the REQ can actually be completed; without `timeline.go` it ships a silent projection
  bug, and without `model_test.go` the canonical gate cannot pass.
  **Risk:** a larger diff to review, and the board tool's tests move in the same commit as its
  parser. Low and reversible — REQ-290 is the only other REQ in this UR and it is strictly
  sequential, so no parallel builder can collide with the widened set.
- **D-03**: Lock-in Check A pins the derivation *verb* applied to the field, not bare
  co-occurrence of `effort_estimate` with an impact word. The plan's literal wording ("no line
  contains `effort_estimate` together with any of `impact-`, `gate`, …") is unimplementable as
  written: the new `impact:` schema row has to name `effort_estimate` to say what question it
  answers and to quote REQ-228, and `work.md`'s short-circuit line legitimately says "stamp
  `calculated_at`" in the same sentence. The check therefore fires on `stamp*` within 80
  non-sentence-ending characters of `effort_estimate`, or on the retired word `gate` anywhere on an
  `effort_estimate` line. Verified against the pre-change tree: it names exactly the fourteen defect
  sites and no others. DECIDE & STATE — reversible and low-reach; the pinned property is unchanged,
  only its detection is precise.
- **D-04**: Reworded `capture.md`'s "GREEN outcome in user-visible terms" to "in terms a user or
  developer would observe". `user-visible` is now an impact token, so leaving it as plain English in
  a shipped action file is the exact two-meanings-one-word problem this REQ exists to remove, and it
  was the single false positive standing between the tree and a strict un-prefixed-axis-word check.
  DECIDE & STATE — one clause, in a declared file, meaning unchanged.
- **D-05**: Wrote one file outside the declared write set —
  `skills/do-work-board/tools/queue-kanban/generate_test.go`. Its single assertion pins the timeline
  forecast sentence's literal `"55 normal at 40 min"`, which Task 8's user-facing label rename turns
  red; the file is the same class D-02 already widened the set for (`model_test.go`,
  `timeline_test.go` — tests asserting the retired pair) and was simply missed in the survey. One
  line changed, `normal` → `substantive`. ESCALATE — the write set is the collision guard, so
  widening it is not a builder's call, and it is reported rather than buried.
  **Value:** the alternative was dropping the `board-timeline.js` label rename, which would ship a
  user-facing tooltip still reading "trivial"/"normal" against a schema that no longer has those
  values — a stale restatement of exactly the kind the review flow's Restatement Sweep exists to
  catch.
  **Risk:** near zero and fully reversible — one assertion string, in the board module whose parser
  and tests already move together in this commit, and REQ-290 is strictly sequential so no parallel
  builder can be holding the file.
- **D-06**: `maintainability-audit-reference.md`'s finding template now carries two things spelled
  "impact": the new lowercase `impact:` line (the REQ frontmatter token) and its pre-existing
  `Impact: [1-5]` score, which is what severity derives from. I kept both and disambiguated in prose
  — the field rules now say "the 1-5 `Impact` score" wherever severity is meant, and the token line's
  values are all `impact-` prefixed, so a plain-text search separates them cleanly. Renaming the
  numeric score would cascade through the ranking rule, the severity rule, and the audit action's
  checklist, none of which this REQ is about. DECIDE & STATE, but flagged for the reviewer as the one
  place the new vocabulary sits next to an older use of the same word.
- **D-07**: `.badge-impact` gets the same neutral/muted treatment as `.badge-effort-estimate` and
  `.badge-assigned` for **every** value, `impact-critical` included. The plan permitted the blocked
  accent for `impact-critical` if the reason were stated; the reason runs the other way. The badge
  describes the work's worth, not its pipeline state, and a red mark on a card whose state is
  perfectly healthy is the confusion the accent colours exist to prevent. Stated in the CSS comment
  beside the rule. DECIDE & STATE — reversible, one declaration.
- **D-08**: `maintainability-audit-reference.md`'s new citation of the schema spells the full
  cross-package path (`../../do-work/actions/work-reference.md`). A toolbox file writing
  `actions/work-reference.md` resolves to nothing; `_dev/tests/shipped-package-reference-contract.sh`
  caught it. DECIDE & STATE — the repo's cross-referencing rule, not a choice.

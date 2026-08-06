---
id: REQ-127
title: Sweep-REQ consolidation — same-root-cause findings land in one sweep, not one REQ each
status: completed
created_at: 2026-08-06T15:48:11Z
claimed_at: 2026-08-06T16:41:40Z
completed_at: 2026-08-06T16:43:22Z
route: A
user_request: UR-027
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-125]
maintenance: true
related: [REQ-125, REQ-126]
batch: follow-up-runaway-fix
write_set: [actions/review-work.md, actions/work.md, actions/work-reference.md, actions/capture-reference.md]
---

# Sweep-REQ Consolidation for Same-Root-Cause Findings

## What

Trivial or same-root-cause findings from a review never get individual REQs. They append to an existing queued sweep REQ for that root cause under the same UR, or create ONE consolidated sweep REQ named for the root cause with a checklist of instances. Sweeps are mechanically findable via a new `sweep: true` frontmatter marker. Only a genuinely non-trivial, thematically unrelated finding still earns its own REQ — and must state in one line why it couldn't fold into a sweep.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Land the routing as a **Sweep consolidation** block in review-work.md Step 10 ahead of the creation template (so it composes with REQ-125's gate and sits above REQ-126's depth stop, which already references sweep appends); schema line in work-reference.md's Full Frontmatter next to `effort_estimate`; recognition-only comment in capture-reference.md; one-sentence routing pointer in work.md's :505 restatement.
- [x] **[APPLY]:** Four files edited, all in the declared write_set. Deviation from the REQ's letter recorded as D-01 (no Schema Read Contract table row for `sweep`).
- [x] **[UNIFY]:** `git diff --stat` reviewed (4 files, prose only). Contract suite re-run at the 7-probe baseline. Verified the prescribed lookup runs: `grep -rl "^sweep: true" do-work/queue/` exits cleanly with no matches on the live queue (no directory error — the Codex P2 fix held through to shipped text). No maintainer-doc citations added.

## Decisions

**D-01 (DECIDE & STATE):** `sweep` is documented in Full Frontmatter but NOT added to the Schema Read Contract's normalize-and-warn table, deviating from this REQ's "Full Frontmatter + Schema Read Contract" phrasing. Reasoning: the contract table's rows mirror `tools/queue-kanban/model.go`'s table one-for-one ("Every row mirrors…" — model.go:850), so a prose-only row would break that mirror or force board parsing this REQ explicitly leaves optional; `review_generated: true` — the marker `sweep` most resembles — follows exactly this pattern (schema-documented, table-absent). Reversible: if the board ever badges sweeps, the field joins both tables and the CLAUDE.md enumeration in that commit.

## Why (if provided)

UR-489 produced fifteen separate REQs that were all facets of ONE root cause (hardcoded colors not tokenized + a guard blind to them). With REQ-125's label and REQ-126's brake, the remaining pain is queue clutter: fifteen chipped-trivial REQs the user would rather see — and approve — as one.

## Detailed Requirements

- **Marker:** new frontmatter field `sweep: true` (boolean; absent reads as false). Document in `actions/work-reference.md` (Full Frontmatter + Schema Read Contract) and `actions/capture-reference.md`. Sweeps are found mechanically — e.g. `grep -rl "^sweep: true" do-work/queue/` filtered to the same `user_request:` — never by judging title similarity (duplicate sweeps would recreate the runaway at half scale). (The `-r` and `^`-anchor are load-bearing: `grep -l` against the bare directory exits 2 with "Is a directory" and returns no candidates, silently defeating consolidation — Codex P2 finding on PR #137; per CLAUDE.md's prescribed-commands rule, the shipped command must actually emit what the step consumes.)
- **Routing (consumes REQ-125's gate token):** a `gate: trivial` finding, or any finding sharing a root cause with others, folds into a sweep. A `gate: user-visible`, thematically unrelated finding still gets its own REQ with a one-line why-not-sweep justification in the body.
- **Append contract** (appending to a queued REQ file is a new write pattern — keep it tight):
  - Append only to a sweep with the same `user_request:` and `status: pending`. Never append to a claimed/working sweep — create a new sweep instead.
  - Appends land as checklist items (`- [ ] [file/site]: [instance]`) under a defined `## Instances` section. The append never touches the sweep's frontmatter.
- **Creation:** when no appendable sweep exists, create ONE sweep REQ named for the ROOT CAUSE (e.g. "tokenize all remaining hardcoded colors and make the guard catch every notation"), with the normal follow-up fields (`user_request`, `addendum_to`, `domain`, `review_generated: true` when review-created), `sweep: true`, an `## Instances` checklist, and `effort_estimate` per the gate: `normal` when solving it changes a multi-site rule (gate (b)), `trivial` otherwise.
- **Definition of done for a sweep, stated in the shipped text:** solving the sweep means the class of finding cannot recur — the rule is changed everywhere it applies — not that N spots got patched one drop at a time.
- **Interaction with REQ-126:** a generation-≥2 review may still append to existing pending sweeps (an append is not a new REQ); a new sweep it needs falls under the reroute (`status: pending-answers`, critical pierces). Sweeps are themselves `review_generated: true` when review-created, so their reviews converge under REQ-126's rule.
- **Sites:** `actions/review-work.md` Step 10 (the routing decision lives here), `actions/work.md` Step 7/8 restatements, `actions/work-reference.md` (schema + a sweep contract home), `actions/capture-reference.md` (schema comment). Re-grep line numbers; coordinate wording with REQ-125/123's edits.
- The board is NOT required to parse `sweep` in this REQ. If the builder chooses to badge it, the same-commit lock-step applies: `tools/queue-kanban/model.go` + the CLAUDE.md board-parsed-fields enumeration.

## Constraints

- No information loss: every instance is enumerated in the sweep's `## Instances` checklist — nothing becomes report-only.
- Severity vocabulary untouched; the gate token (REQ-125) is the routing input.
- **Out of scope, deliberately (user decision):** fix-inline-at-review-resolution — the originally proposed priority 4. Do not build it here; the deferral is on record in UR-027's decision record, to be revisited only if labeled-and-gated trivia still feels too heavy after living with REQ-125/123/124.

## Dependencies

Depends on REQ-125 (routing consumes the gate token and sets `effort_estimate`). Complements REQ-126; buildable before or after it, but the interaction paragraph above must match whichever text is already shipped.

## Builder Guidance

Certainty: Firm on the marker, append contract, and root-cause naming (all user-confirmed). Latitude on the exact `## Instances` item format.

## Red-Green Proof

**RED prompt/case:** A review yields three Important findings sharing one root cause (e.g. three more hardcoded-color sites): today `actions/review-work.md` Step 10 creates three separate REQs — the UR-489 pattern (fifteen sibling REQs for one cause). `grep -rn "sweep" actions/review-work.md` shows no consolidation path.
**Why RED now:** One-REQ-per-Important has no consolidation concept; nothing marks or finds a sweep.
**GREEN when:** The same three findings produce exactly one queued REQ with `sweep: true` and a three-item `## Instances` checklist (or three appended items on an existing pending sweep under the same UR), and a standalone follow-up REQ's body contains a one-line why-not-sweep justification.
**Validation:** User confirmed (design discussion preceding capture, 2026-08-06)

## Full Context

See `do-work/user-requests/UR-027/input.md` for complete verbatim input and the decision record.

---
*Source: "do-work capture-request Ship priorities 1 through 3" — priority 3 of the agreed design*

## Implementation Summary

**Files changed:**
- `actions/review-work.md` (modified) — **Sweep consolidation** block in Step 10: mechanical lookup (`grep -rl "^sweep: true"` + UR + pending filters, never title-judging), append contract (`## Instances` checklist lines only, frontmatter untouched, claimed sweeps never appended), root-cause-named creation with `sweep: true` + gate-derived `effort_estimate`, class-cannot-recur definition of done, standalone-REQ justification line, and the generation-≥2 interaction (appends allowed, new sweeps reroute)
- `actions/work.md` (modified) — :505 restatement gained the consolidation routing sentence
- `actions/work-reference.md` (modified) — `sweep:` schema line in Full Frontmatter (marker class, greppable by design, append rules, board/table posture per D-01)
- `actions/capture-reference.md` (modified) — recognition-only `sweep: true` line in the additional-frontmatter list (never emitted by capture)

**What was done:** Same-root-cause and trivial findings now consolidate into one greppably-marked sweep REQ per root cause with an enumerated `## Instances` checklist — no information loss, no per-facet REQs — and only a standalone user-visible finding still earns an individual REQ, with its why-not-sweep line. Inline-fix-at-review-resolution stays deliberately unbuilt per the REQ's Constraints.

## Qualification

Passed — 4 files verified in the diff, all Detailed Requirements traced (D-01 records the one deliberate deviation), P-A-U confirmed. Mechanical script OK; Route A, no Scope section.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; live execution of the prescribed sweep lookup; grep verification of RED/GREEN
**Result:** ✓ Contract suite at its pre-existing 7-probe root-runner baseline. `grep -rl "^sweep: true" do-work/queue/` executes cleanly against the live queue (exit 1, no matches, no "Is a directory" error — the P2 fix survived into shipped text). Non-behavioral-code change — regression evidence in place of red-green tests.

**Red-green validation:**
- Captured RED (`grep -rn "sweep" actions/review-work.md` had no consolidation path) is now GREEN: Step 10 carries the **Sweep consolidation** block; three same-root-cause findings route to one `sweep: true` REQ with a three-item `## Instances` checklist per the prescribed flow

*Verified by work action*

## Review

**Overall: 95%** | 2026-08-06T16:45:00Z

**Approve** — consolidation, marker, append contract, and the depth-stop interaction all land at their canonical homes; the one deviation from the REQ's letter is argued and recorded (D-01).
Route A | uncommitted (hash written back at Step 9)

**Findings:**

**Important:** None.

**Minor:**
- The `## Instances` checklist item format (`- [ ] [file/site]: [instance]`) is defined only by example in Step 10; a builder wanting stricter structure has no template block. Acceptable for a human-checked list. Report-only.
- D-01's contract-table omission deviates from the REQ's written requirement — reviewed and endorsed: the mirror constraint the REQ didn't anticipate makes the letter unsatisfiable without forcing optional board work; the `review_generated` precedent covers it.

**Requirements Checklist:**
- [x] `sweep: true` marker, greppable, schema-documented — delivered (Full Frontmatter; contract table per D-01)
- [x] Mechanical lookup with UR + `status: pending` filters, title-judging forbidden — delivered
- [x] Append contract: `## Instances` lines only, frontmatter untouched, claimed sweeps excluded — delivered
- [x] Root-cause naming with the user's example verbatim — delivered
- [x] Class-cannot-recur definition of done, stated in shipped text — delivered
- [x] Standalone-REQ one-line justification — delivered
- [x] Generation-≥2 interaction (appends allowed, new sweeps reroute, sweeps carry `review_generated`) — delivered
- [x] Inline fixes NOT built — confirmed absent, deferral preserved in Constraints

**Acceptance Testing — Result: Pass** (prescribed lookup executed against the live queue; all prescribed text greps hit)

**Scores:** Requirements 96% | Code Quality 95% | Test Adequacy 92% | Scope 100% | Acceptance Pass

*Reviewed in pipeline mode by work action (Step 7) — Route A quick scan*

## Orientation

Review findings that share a root cause now land as one sweep REQ with a checklist of instances instead of a REQ per facet — solving a sweep means the class can't recur. Lives in the review/follow-up machinery (`actions/review-work.md` Step 10), completing the three-layer UR-027 set: REQ-125's gate prices findings, REQ-126's depth stop gates the cascade, this consolidates the trivia. Closes UR-027.

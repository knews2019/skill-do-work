---
id: REQ-045
title: "Dispatch re-validation completeness: Route A gap, undefined loser, partition visibility, absent-set gloss, missing ratchet"
status: completed
route: C
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T12:50:00Z
completed_at: 2026-07-29T13:04:42Z
commit: 14146c2
user_request: UR-008
addendum_to: REQ-036
domain: general
prime_files: []
tdd: false
depends_on: []
related: [REQ-043, REQ-044]
batch: deep-review-followups
write_set: [actions/work.md, actions/work-reference.md, _dev/tests/contract-regressions.sh]
maintenance: true
---

# Dispatch Re-Validation Completeness (REQ-036 Follow-Up)

## What

REQ-036 added Step 5.5 re-validation of write-set disjointness, but the coverage claim has a route-shaped hole and three coherence gaps around the partition path.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/general.md`, `coding-guardrails.md`, `maintenance.md` (`maintenance: true`); read on-disk `actions/work.md`, `actions/work-reference.md` (post-REQ-043/044), `_dev/tests/contract-regressions.sh`. Approach in `## Plan` — one invariant (route picks the post-dispatch validation point) stated once, every restatement re-pointed at it. `prime_files` is empty; no prime to read.
- [x] **[APPLY]:** Ten edits across exactly the three declared files. No new files, no files outside `write_set`.
- [x] **[UNIFY]:** See `## Implementation Summary` → **UNIFY**.

## Prior Implementation

REQ-036 (archived, commit `4296e11`, v0.145.0): Step 5.5 re-checks disjointness against every other in-flight REQ's current `write_set` when firming from `## Scope` (`actions/work.md` ~:340), partition-survives-mirror stated in both homes, Step 1 gate pointer and Step 4/Step 6 glosses aligned.

## Detailed Requirements

1. **Close the Route A gap.** `actions/work.md`'s dispatch-gate bullet (~:184) claims "Each REQ firms its real set from its `## Scope` at Step 5.5, which re-validates disjointness" — but Route A skips Step 5.5 entirely (stated at ~:329/:338), and the route isn't decided until Step 3, after co-dispatch committed. A co-dispatched Route A REQ's builder writes under an unvalidated capture hint — the unguarded path REQ-036 claims to have closed. Fix: define where a Route A REQ's set is validated (at triage when the route is assigned, or force serial for co-dispatched REQs that triage to Route A), and make the three sentences (~:184, ~:338, ~:340) agree.
2. **Define the serialization "loser."** "the orchestrator serializes the loser" never says which of the overlapping pair loses. State it: the REQ currently at Step 5.5 (the one that just firmed into the overlap) is held; an already-dispatched sibling mid-build is never the loser.
3. **Make dispatch-time partitions visible to the re-check.** Step 1's partition directive hands each builder "its own subset as its declared set" but never writes the partition into the sibling's `write_set` frontmatter — so at a sibling's Step 5.5 the re-check reads the full un-partitioned set and spuriously serializes the partition the orchestrator itself issued. Either the partition lands in each REQ's frontmatter at dispatch, or orchestrator-issued partitions are explicitly exempt from the overlap re-check.
4. **Fix the absent-`write_set` gloss.** Step 6's gloss says an absent/empty set "means no boundary was declared because this REQ was dispatched serially" — but Route A REQs and crash-recovered REQs (whose set recovery cleared) also arrive absent. Reword so a builder doesn't read a recovered REQ's stripped set as full-scope serial freedom without a stated boundary.
5. **Add a ratchet.** REQ-036's deliverable — the Step 5.5 re-validation clause — has no regression guard in `_dev/tests/contract-regressions.sh` and could be reverted silently. Pin it.

## Constraints

- Coordinate item 4 with REQ-044 item 3 (which changes when recovery clears the set) — whichever lands second reconciles the gloss with the surviving clear-semantics.
- Fails-safe is the tiebreak: when in doubt, the resolution serializes rather than co-dispatches.

## Red-Green Proof
**RED prompt/case:** `actions/work.md` ~:184 asserts every co-dispatched REQ re-validates at Step 5.5 while ~:338 says Route A skips Step 5.5 — the two sentences contradict, and no text validates a Route A set post-dispatch.
**Why RED now:** The re-validation contract was written against the Route B/C pipeline shape only.
**GREEN when:** Route A's validation point (or serial-only rule) is stated and the three sentences agree; the loser is defined; a partitioned sibling's Step 5.5 re-check cannot spuriously flag the orchestrator's own partition; the absent-set gloss names all three producers of absence; a contract-regressions ratchet pins the re-validation clause.
**Validation:** User confirmed (approved capture of the reviewed finding set)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** C (Full Pipeline)
**Reasoning:** Five interlocking coherence defects in the parallel-dispatch / Step 5.5 re-validation contract — the concurrency-scheduling surface REQ-036 built. Touches `actions/work.md` (dispatch gate ~:184, Step 5.5, Step 6 gloss), `actions/work-reference.md`, and a ratchet in `_dev/tests/contract-regressions.sh`. Self-referential, fails-safe-critical, and it must reconcile item 4 with REQ-044's just-landed conditional `write_set` clear → Route C.
**Complexity indicators:** 5 requirements, one needing a design decision (Route A validation point vs force-serial — item 1); a cross-REQ reconciliation (item 4 ↔ landed REQ-044 item 3); repo-wide restatement grep; a new ratchet.
**Rigor:** Full adversarial-grade review at Step 7 (main-context, against the full diff + independent restatement sweep + the three-sentence agreement check).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Plan

**Shape of the fix:** one invariant stated once and pointed at from every restatement — *every co-dispatched REQ gets exactly one post-dispatch re-validation, and its route picks which one.* Routes B/C keep Step 5.5; Route A gets Step 3. That single sentence is what makes REQ-036's coverage claim true instead of aspirational, and it is what the three contradicting sentences are rewritten to agree on.

1. **Item 3 first (it is a precondition for item 1's coherence).** Step 1's partition bullet gains a dispatch-time frontmatter write: the narrowed subset lands in the REQ's `write_set` as it is dispatched, narrowing-only. Without this, every downstream re-check reads the un-partitioned set and the new Step 3 check inherits the same spurious-serialization bug Step 5.5 already had. Mirror the new writer into `actions/work-reference.md`'s `write_set` schema line (the canonical "who writes this field" home) and into the Scope-template's partition-survives-the-mirror paragraph.
2. **Item 1.** Rewrite `actions/work.md`'s dispatch-gate bullet to state the route-picks-the-point invariant (and why this gate *cannot* be the only enforcement point: the route is unknown until Step 3). Add a Step 3 paragraph — the Route-A re-validation — whose only disposition is **serialize** (no partition: Route A has no `## Scope` for a partition to be recorded in and no later step that would re-check a narrowed one). Re-point Step 5.5's mirror note and re-validation clause at it. Also fix the Step 4 plan-validation restatement and the Orchestrator Checklist's Step 3 line.
3. **Item 2.** Split Step 5.5's overlap disposition into its own paragraph headed by the rule: *the REQ at this re-check is the loser.* An already-dispatched sibling mid-build is never held. Add the deadlock guard the two-discoverer case needs: a REQ already held for an overlap is not a contender.
4. **Item 4.** Reword Step 6's absent-set gloss. State the *condition* that produces absence (never ran Step 5.5, so nothing mirrored a `## Scope` in) rather than a closed list, then name today's instances — capture seeded nothing, or crash recovery cleared a mirror whose `## Scope` it stripped — and explicitly note REQ-044's other half: a capture-seeded set is *preserved* through recovery, so absence is never the recovery signal for a Scope-less REQ.
5. **Item 5.** Four block-scoped assertions in `_dev/tests/contract-regressions.sh` (Step 5.5 block, Step 3 block, parallel-dispatch block) pinning: the re-validation clause, the loser definition, Route A's validation point, the partition frontmatter write. Block-scoped because the file-wide `write_set`/`pairwise disjoint` assertions already present would be satisfied by a match anywhere in `actions/work.md` — the exact masking REQ-044 hit.

**Verification:** remove the pinned Step 5.5 clause → suite FAILs naming it; restore → passes. Then `bash _dev/tests/contract-regressions.sh` green, plus a re-grep of every swept token to confirm no restatement was left telling the old story.

## Exploration

Restatement sweep over `SKILL.md`, `next-steps.md`, `README.md`, `actions/`, `crew-members/`, `docs/`, `prompts/`, `interviews/`, `specs/`, `tools/`. Tokens grepped: `write_set`, `partition` (case-insensitive), `loser`, `serializes the`, `re-validat|revalidat`, `pairwise disjoint`, `Step 5.5`, `Route A`.

**Sites touched (5 restatements of the phrasings in scope):**

- `actions/work.md:183` — partition directive bullet. Issued the directive but never persisted it (item 3's root).
- `actions/work.md:184` — dispatch-gate bullet. The false coverage claim (item 1's RED).
- `actions/work.md:302` — Step 4 plan-validation disposition note: "plus the Step 5.5 re-validation" — a fourth restatement of the enforcement point the REQ did not list; it inherits item 1's hole verbatim.
- `actions/work.md:338` — Step 5.5 mirror note, the "Route A skips Step 5.5" sentence.
- `actions/work.md:340` — Step 5.5 re-validation clause + the undefined loser + the partition-survives-mirror clause (items 1, 2, 3 all land here).
- `actions/work.md:399` — Step 6 write-boundary bullet, absent-set gloss (item 4).
- `actions/work.md:656` — Orchestrator Checklist Step 3 line. A real read site: an orchestrator working the checklist would never learn the new gate exists.
- `actions/work-reference.md:100` — `write_set` schema line, the canonical "who writes this field" home. Said "`## Scope` is the source and this field is its mirror, never the reverse" — the dispatch-time partition write is a second writer and had to be declared here or it reads as a contract violation.
- `actions/work-reference.md:559` — Scope-template partition-survives-the-mirror paragraph (item 3's second home).

**Sites checked and deliberately left (no contradiction with the new text):**

- `actions/work-reference.md:208` — Crash Recovery substep 1. REQ-044's conditional clear already tells the correct story, including "a Route A REQ that never runs Step 5.5 at all." Item 4's gloss was written to match it, not to change it.
- `actions/work.md:710` — Rules bullet ("an unknown or overlapping set means serialize"). Scheduling-only claim, still true.
- `actions/board.md:92`, `actions/board.md:117`, `tools/queue-kanban/prime-do-kanban.md:51` — the board's display-only overlap annotation. All three already say the co-dispatch decision belongs to the work pipeline's gate and that the gate takes the safety-first reading of an absent set. Unaffected.
- `docs/work-guide.md:87` — user-facing parallel-dispatch summary. Restates only the Step 1 gate and the absent-set-runs-alone rule, both unchanged; it never mentions Step 5.5, the loser, or partitions, so there is nothing to re-point. Outside `write_set` and needs no edit.
- `actions/capture-reference.md:113` — see `## Discovered Tasks`. Its "the work pipeline's Scope step firms it up and overwrites it" is true for Routes B/C only; a Route A REQ keeps capture's value forever. Pre-existing imprecision, not created by this change, and outside the declared `write_set`.

**Three-sentence agreement check (item 1's GREEN condition):** `actions/work.md:184` now states the route-picks-the-validation-point invariant; the Step 5.5 mirror note (~:338) and re-validation clause (~:340) both scope themselves to Routes B/C and point at Step 3 for Route A; the new Step 3 paragraph is the third corner and says the same thing from its own side. No sentence claims Step 5.5 covers Route A.

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — dispatch-gate bullet, partition bullet, new Step 3 Route-A re-validation, Step 5.5 mirror note + re-validation clause + loser, Step 4 plan-validation note, Step 6 absent-set gloss, Orchestrator Checklist
- `actions/work-reference.md` (modify) — `write_set` schema line (partition writer), Scope-template partition-survives-the-mirror paragraph
- `_dev/tests/contract-regressions.sh` (modify) — four block-scoped ratchets

**Files I will NOT touch:** `docs/work-guide.md`, `actions/capture-reference.md`, `actions/board.md`, `tools/queue-kanban/*` (all swept, none contradict — see `## Exploration`); `actions/version.md` and `CHANGELOG.md` (orchestrator's Step 9).

**Acceptance criteria (restated from REQ):**
- [ ] Route A's post-dispatch validation point is stated, and `actions/work.md` ~:184 / ~:338 / ~:340 agree with it and with each other
- [ ] The serialization loser is named: the REQ at the re-check, never an already-dispatched sibling mid-build
- [ ] A partitioned sibling's Step 5.5 re-check cannot spuriously flag the orchestrator's own partition
- [ ] The absent-`write_set` gloss names every producer of absence, consistent with REQ-044's landed conditional clear
- [ ] `_dev/tests/contract-regressions.sh` pins REQ-036's re-validation clause (red-green verified)
- [ ] `bash _dev/tests/contract-regressions.sh` passes

## Implementation Summary

**Files changed:**

- `actions/work.md` (modify) — seven passages: Step 1's partition bullet (dispatch-time frontmatter write, narrowing-only); Step 1's dispatch-gate bullet (the route-picks-the-validation-point invariant, replacing the false blanket claim); a new two-paragraph **Route-A re-validation** in Step 3; Step 4's plan-validation disposition note (re-pointed); Step 5.5's mirror note (points at Step 3 for Route A); Step 5.5's re-validation clause (scoped to Routes B/C, plus the partition-visibility explanation) and a new **The REQ at this re-check is the loser** paragraph; Step 6's absent-`write_set` gloss (item 4); the Orchestrator Checklist's Step 3 line.
- `actions/work-reference.md` (modify) — `write_set` schema line now declares all writers of the field (capture seed, Step 1 partition, Step 5.5 mirror, recovery's conditional clear); the Scope-template partition-survives-the-mirror paragraph now states that the gate persisted the subset at dispatch and why.
- `_dev/tests/contract-regressions.sh` (modify) — four block-scoped assertions plus a header comment recording the defect each one guards.

**What was done:**

The five defects shared one root: REQ-036 wrote the re-validation contract against the Route B/C pipeline shape and asserted it covered everything. The fix states the covering invariant once — *every co-dispatched REQ gets exactly one post-dispatch re-validation, and its route picks which one* — gives Route A a real validation point at Step 3 (serialize-only, since it has no `## Scope` to hold a partition and no later step to re-check one), names the loser as the REQ at the re-check, persists dispatch-time partitions into frontmatter so a sibling's re-check compares against the subset the gate actually issued, replaces the absent-set gloss with a condition-first statement that reconciles with REQ-044's landed conditional clear, and pins all four contracts against silent reversion.

Two clauses were added beyond the five requirements because the new rules created the cases: **(a)** the deadlock guard — once "the discoverer holds" is the rule, two REQs can each discover the same overlap, so an already-held REQ is explicitly not a contender; **(b)** narrowing-only on the dispatch-time write, which is what keeps the new second writer from being able to contradict `## Scope`.

**UNIFY** — `git diff --stat`: `_dev/tests/contract-regressions.sh` +34, `actions/work-reference.md` +2/-2, `actions/work.md` +15/-7 (3 files, 51 insertions, 9 deletions). Files verified individually:

- `actions/work.md` — read all seven changed passages in full step context (Step 1 gate block, Step 3, Step 5.5, Step 6 bullet, checklist). Confirmed the three item-1 sentences agree and no passage still claims Step 5.5 covers Route A. Re-grepped 11 edit markers: each present exactly once.
- `actions/work-reference.md` — reviewed both hunks against the diff. Crash Recovery substep 1 (REQ-044) deliberately untouched; item 4's gloss was written to match it.
- `_dev/tests/contract-regressions.sh` — `bash -n` clean. Red-green verified per assertion: removing the Step 5.5 re-validation clause, the loser definition, the Step 3 validation point, or the Step 1 partition write each FAILs with its own message; restoring passes. (The first RED probe initially showed no failure — a flawed probe, not a flawed ratchet: the pinned phrase now occurs in both the Step 3 and Step 5.5 blocks and a non-global substitution hit only the Step 3 copy, which the Step-5.5-scoped assertion correctly ignores. Re-run with a global substitution, it fires.) `actions/work.md` confirmed byte-identical to its pre-probe baseline afterwards.
- No debug artifacts: diff grepped for `console.log`/`TODO`/`FIXME`/`XXX`/`DEBUG` and the probe's `REDACTED-FOR-RED-PROOF` marker — none present.
- `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` (exit 0).
- No linter beyond this suite applies to Markdown prose; the one shell file was checked with `bash -n`.

## Decisions

**D-01 — Route A's validation point is Step 3, and serializing is its only disposition there (item 1).** ESCALATE. The alternatives were (a) validate at triage, (b) force every co-dispatched REQ that triages to Route A to serial unconditionally, (c) give Route A a Step 5.5. Chose (a) with serialize-only: it keeps Route A's cheapness (a co-dispatched Route A REQ that genuinely does not overlap still runs concurrently, which (b) would forfeit) while making the fails-safe outcome the only one available where verification is impossible. Rejected partitioning at Step 3 specifically: Route A has no `## Scope` for a partition to be recorded in and no later step that would re-check a narrowed set, so a Route A partition is a boundary nothing can verify — precisely the unguarded shape this REQ exists to close. Rejected (c) as the more expensive fix: Step 5.5 exists to declare scope for work whose scope is unclear, and Route A is defined by having clear scope; adding the step would contradict the route's reason to exist. **Value:** the coverage claim becomes true, and the invariant is now one sentence a reader can check against three sites. **Risk:** a co-dispatched Route A REQ can now be held at triage where it previously ran unguarded — a throughput cost, not a correctness one, and reversible by narrowing the check.
**D-02 — Dispatch-time partitions land in `write_set` frontmatter rather than being exempted from the re-check (item 3).** ESCALATE. The exemption alternative is one sentence and needs no write, but it makes the re-check knowingly blind: it would have to trust an out-of-band directive the field contradicts, and every future reader of `write_set` inherits a field that does not describe reality. Persisting the subset makes the re-check see the truth, so no exemption is needed and the board's overlap badge stops showing a conflict the orchestrator already resolved. **Value:** one representation of the boundary instead of two, and the re-check needs no special case. **Risk:** a second writer to a field documented as a one-directional mirror — mitigated by making it narrowing-only and by declaring it in the schema line, which is where a future reader looks for who writes the field.
**D-03 — Added the deadlock guard ("a REQ already held for an overlap is not a contender").** DECIDE & STATE. Defining the loser as the discoverer (item 2) creates a two-discoverer case: a Route B/C REQ held at its Step 5.5 leaves the overlap in its mirrored `write_set`, and the Route A sibling then discovers the same overlap at Step 3. Without this clause both hold and neither resumes. One sentence in each of the two re-check homes, and still fails-safe (the held REQ is not writing).
**D-04 — Stated item 4's producers as a condition, not a closed list.** DECIDE & STATE. The gloss now leads with the condition that produces absence — *this REQ never ran Step 5.5, so nothing mirrored a `## Scope` into the field* — and names Route A as today's instance. A hand-enumerated list of producers is exactly the pattern that goes stale as routes or recovery semantics change (it is how the defect arose: REQ-036's list said "serial"). This also corrected the enumeration itself: at builder-dispatch time a *serial* Route B/C run still writes and mirrors a `## Scope`, so serial-ness never produces absence there — not running Step 5.5 does.

## Discovered Tasks

- `actions/capture-reference.md:113` says capture's `write_set` value is one "the work pipeline's Scope step … firms it up and overwrites it." True for Routes B and C only — a Route A REQ never runs Step 5.5, so it keeps capture's value for the whole run (and, after this REQ, that value is what Step 3 re-validates). Pre-existing imprecision, outside this REQ's declared `write_set`; a one-clause qualification would close it.

## Review

**Acceptance: Pass — overall ~95%.** Main-context adversarial review against the full 3-file diff + independent restatement sweep + the three-sentence agreement check.

**Requirements (all 5 met):**
1. Route A gap closed by one invariant — *every co-dispatched REQ gets exactly one post-dispatch re-validation, its route picks which* (B/C → Step 5.5, Route A → Step 3, serialize-only). The three previously-contradicting sentences (dispatch gate, Route-A-skips-5.5, Step 5.5 clause) now all point at it; Orchestrator Checklist Step 3 updated.
2. Loser defined: the REQ at the re-check is held; a dispatched sibling mid-build is never held — plus a deadlock guard (a held REQ is not a contender) for the two-discoverer case.
3. Partitions land in `write_set` frontmatter at dispatch (narrowing-only), so a sibling's re-check compares against the real subset; schema line updated to declare all writers.
4. Absent-set gloss reworded condition-first ("never ran Step 5.5 ⇒ Route A today"), correctly reconciled with REQ-044's landed conditional clear (a capture-seeded set is preserved through recovery, so a no-Scope REQ is never absent); warns against reading a recovered-and-stripped set as full-scope freedom.
5. Ratchet added (block-scoped assertions pinning the Step 5.5 clause, the loser, the Step 3 point, and the Step 1 partition write); red-green verified per assertion; contract-regressions passes.

**Coherence/sweep:** restatement sweep across `actions/` clean — the invariant is stated once and referenced everywhere; no passage still claims Step 5.5 covers Route A. D-04 correctly applies the repo's "Closed Enumerations Go Stale" principle.

No Important/Critical findings. One Discovered Task queued as REQ-057 (`pending-answers`): a pre-existing `capture-reference.md` imprecision, outside this REQ's write_set.

## Lessons Learned
**What worked:** Collapsing five defects to one root cause (REQ-036 wrote the contract against the B/C shape only) and stating a single covering invariant made the three-sentence agreement mechanical rather than a hunt.
**What didn't:** The first red-green probe on the new ratchet showed a false green — the pinned phrase now lives in both the Step 3 and Step 5.5 blocks, and a non-global substitution hit only one copy while the block-scoped assertion correctly ignored it. Fixed by a global substitution. Lesson: when a phrase is intentionally restated in two blocks, a ratchet's red-probe must perturb the specific block the assertion scopes to.
**Worth knowing:** "Serial" never produces an absent `write_set` at build time — a serial B/C REQ still mirrors a `## Scope`; only *not running Step 5.5* (i.e. Route A) does. The REQ-036 defect originated in exactly that stale enumeration.

## Orientation
The parallel-dispatch write-set contract now has full route coverage: Route A co-dispatched REQs re-validate (serialize-only) at Step 3, Routes B/C at Step 5.5, both against partitions that are now persisted into `write_set` at dispatch. Lives in `actions/work.md` (Step 1 gate, Step 3, Step 5.5, Step 6 gloss) + `actions/work-reference.md` (write_set schema line, Scope template) with a contract-regression ratchet. No map change — closes the route-shaped hole in the REQ-036 subsystem.

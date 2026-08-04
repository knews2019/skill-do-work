---
id: REQ-036
title: Re-validate write-set disjointness when Step 5.5 firms the sets
status: completed
claimed_at: 2026-07-29T07:54:52Z
completed_at: 2026-07-29T08:07:08Z
commit: 4296e11
route: B
kb_status: promoted
kb_entry: REQ-036-re-validate-write-set-disjointness-when.md
created_at: 2026-07-28T21:57:21Z
user_request: UR-007
addendum_to: REQ-032
depends_on: [REQ-033]
related: [REQ-032, REQ-033, REQ-035]
batch: parallel-dispatch
domain: general
prime_files: []
tdd: false
review_generated: true
write_set:
  - actions/work.md
  - actions/work-reference.md
maintenance: false
---

# Re-validate write-set disjointness when Step 5.5 firms the sets

## What

REQ-032's dispatch gate decides co-dispatch on capture-seeded `write_set`s, which `actions/capture-reference.md` itself calls "a hint, never a commitment" — and Step 5.5 then overwrites the field from the `## Scope` list with no conflict re-check. Add the missing re-validation: on a concurrent dispatch, the Step 5.5 mirror runs the same no-concurrent-claimant check Step 6 already requires for mid-build extensions, and serializes/partitions when it fails.

## Why (if provided)

Review of REQ-032 (confirmed by 2 independent adversarial verifiers per finding): REQ-A seeds [src/a.ts], REQ-B seeds [src/b.ts] → co-dispatched; each Scope then declares src/shared.ts; both mirror, both write it, no rule fires. The same field has two write paths — the Step 6 mid-build extension is guarded, the Step 5.5 firming is not. The mirror also overwrites a dispatch-time partition directive (the narrowed set handed to a builder), erasing the partition.

## Detailed Requirements

- One clause in Step 5.5 (plus a pointer in the Step 1 gate): when REQs were co-dispatched, the mirror write re-checks pairwise disjointness against every other in-flight REQ's current `write_set` before replacing the field; on overlap, the orchestrator serializes the loser (or issues/refreshes an explicit partition directive) before its builder starts.
- State how a partition directive survives the mirror: the narrowed set handed at dispatch is what Step 5.5 writes (the partition is the Scope), not the un-partitioned Scope list.
- Gloss the absent-set case in the Step 6 write-boundary bullet for builders (absent = no declared boundary because the REQ was dispatched serially — not "write nothing"; split finding of REQ-032's review).
- One sentence reconciling Step 4's plan-validation file-conflict warning with the gate ("warnings, not blockers" is the plan validator's disposition; dispatch-time serialization is the enforcement point — split finding).
- Floor agents: zero behavior change; all new text stays inside the optional parallel-dispatch path.

## Constraints

- `actions/work.md` + `actions/work-reference.md` only; keep `_dev/tests/contract-regressions.sh` green.

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" is fully specified (5 named edits in the REQ's Detailed Requirements); what needs discovery is the exact current text/anchors of the four edit sites, which shifted after REQ-035 rewrote the adjacent gate/Step-2/Step-8 text. No architectural decision to plan — the requirements are the design. Exploration (folded into the orchestrator's scoping, since the current files are already in context post-REQ-035) locates the anchors.

**Planning:** Not required (Route B)

## Plan

**Planning not required** - Route B: Exploration-guided implementation.

*Skipped by work action*

## Exploration

Anchors verified against disk at claim time (post-REQ-035, commit `fd56267`):

- **Step 5.5 mirror paragraph** — `actions/work.md:337`, begins "**Mirror the file list into `write_set`.**" This is where the co-dispatch re-validation clause attaches: the mirror currently "replac[es] whatever capture seeded" with no conflict re-check.
- **Step 1 parallel-dispatch gate bullets** — `actions/work.md:182` (absent/empty ⇒ overlaps ⇒ serialize) and `:183` (overlapping pairs serialize or take an explicit partition directive at dispatch). The Step-1 pointer to Step 5.5's re-validation attaches here.
- **Step 6 write-boundary bullet** — `actions/work.md:396`, "**Write only inside the declared `write_set`:**". This is where the absent-set gloss for builders attaches (absent = dispatched serially, no declared boundary — not "write nothing").
- **Step 4 plan-validation file-conflict warning** — `actions/work.md:301`, item 4 "**File conflicts:**" ("warnings, not blockers" is the plan validator's disposition). The reconciliation sentence attaches here or at the gate.
- **The "hint, never a commitment" characterization** the REQ cites lives at `actions/capture-reference.md:113` (out of scope — referenced, not edited).
- **Scope Declaration Template** — `actions/work-reference.md` "Scope Declaration Template (Step 5.5)" carries the one-directional mirror note; the partition-survives-mirror rule may attach there so the template and Step 5.5 tell one story.

No `## Scope` source exists yet for the re-validation (it is new text); no parser/schema shape change (model.go untouched).

*Exploration folded into orchestrator scoping (Route B)*

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — Step 5.5 mirror paragraph (`:337`): add the co-dispatch re-validation clause + partition-survives-mirror rule; Step 1 gate (`:182`–`:183`): add the pointer that Step 5.5 is the dispatch-time enforcement point; Step 6 write-boundary bullet (`:396`): gloss the absent-set case for builders; Step 4 plan-validation item 4 (`:301`): reconcile "warnings, not blockers" with dispatch-time serialization
- `actions/work-reference.md` (modify) — Scope Declaration Template (Step 5.5): state how a dispatch-time partition directive survives the mirror (the narrowed set handed at dispatch is what Step 5.5 writes, not the un-partitioned Scope list), keeping the template and Step 5.5 consistent

**Files I will NOT touch:** `_dev/tests/contract-regressions.sh` (keep green; no new ratchet — the existing `pairwise disjoint` gate ratchet already anchors this surface, and the REQ scopes to the two action files); `actions/capture-reference.md` (the "hint, never a commitment" line is referenced, not edited); `tools/queue-kanban/model.go` (no schema shape change)

**Acceptance criteria (restated from REQ):**
- [ ] One clause in Step 5.5 (plus a Step 1 gate pointer): on a concurrent dispatch, the mirror write re-checks pairwise disjointness against every other in-flight REQ's current `write_set` before replacing the field; on overlap, the orchestrator serializes the loser (or issues/refreshes an explicit partition directive) before its builder starts
- [ ] State how a partition directive survives the mirror: the narrowed set handed at dispatch is what Step 5.5 writes (the partition is the Scope), not the un-partitioned Scope list
- [ ] Gloss the absent-set case in the Step 6 write-boundary bullet: absent = no declared boundary because the REQ was dispatched serially — not "write nothing"
- [ ] One sentence reconciling Step 4's plan-validation file-conflict warning with the gate: "warnings, not blockers" is the plan validator's disposition; dispatch-time serialization is the enforcement point
- [ ] Floor agents: zero behavior change; all new text stays inside the optional parallel-dispatch path
- [ ] `bash _dev/tests/contract-regressions.sh` stays green

*Scope declared by work action (orchestrator, session do-work-20260729T065754Z-5724)*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]** Five additive edits, all scoped to the optional parallel-dispatch / co-dispatch path so a serial (floor) run reads exactly as before. `actions/work.md`: (1) Step 5.5 mirror — add a co-dispatch re-validation paragraph (re-check pairwise disjointness of the firmed `## Scope` list against every other in-flight REQ's current `write_set` before replacing the field; serialize/partition the loser before its builder starts; state that a dispatch-time partition directive survives the mirror), keeping the existing one-directional Scope→write_set paragraph intact; (2) Step 1 gate — add a bullet that the disjointness decided on capture-seeded sets is re-validated at Step 5.5 (this gate = initial scheduling decision, Step 5.5 = dispatch-time enforcement); (3) Step 6 write-boundary bullet — gloss the absent-set case (serial dispatch ⇒ no declared boundary, full-scope freedom, not "write nothing"); (4) Step 4 plan-validation item 4 — one sentence: the file-conflict flag is a warning, enforcement is Step 1 + Step 5.5. `actions/work-reference.md`: (5) Scope Declaration Template — partition-survives-mirror rule so template and Step 5.5 tell one story.
- [x] **[APPLY]** Edits applied, in declared scope: `actions/work.md` (Step 5.5 re-validation paragraph; Step 1 gate pointer bullet; Step 6 absent-set gloss; Step 4 item 4 sentence) and `actions/work-reference.md` (Scope Declaration Template partition-survives-mirror note).
- [x] **[UNIFY]** `git diff --stat` = `actions/work.md` (+7/-2) and `actions/work-reference.md` (+2) only; `git status --porcelain -uall` shows no stray untracked files outside `do-work/`; `bash _dev/tests/contract-regressions.sh` exits 0 (all ratchets green, none added); no debug artifacts. Files: `actions/work.md`, `actions/work-reference.md` (both modified, in declared scope).

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified)
- `actions/work-reference.md` (modified)

**What was done:** Closed the unguarded second write path to `write_set`. `actions/work.md` Step 5.5 gained a co-dispatch re-validation paragraph — when more than one REQ is in flight, the mirror re-checks the firmed `## Scope` list for pairwise disjointness against every other in-flight REQ's current `write_set` before replacing the field (the same no-concurrent-claimant check Step 6 runs for mid-build extensions), serializing or partitioning the loser before its builder starts, and stating that a dispatch-time partition directive survives the mirror. The Step 1 gate gained a pointer naming Step 5.5 as the dispatch-time enforcement point (the gate itself being the initial scheduling decision on capture-seeded hints); the Step 6 write-boundary bullet now glosses the absent-set case (serial dispatch ⇒ no declared boundary, not "write nothing"); and Step 4's plan-validation file-conflict item now reconciles its flag as a warning whose real enforcement is Step 1 + Step 5.5. `actions/work-reference.md`'s Scope Declaration Template gained the matching partition-survives-mirror note. Every added clause is gated to the parallel-dispatch path; the serial/floor path is behaviorally unchanged (the absent-set gloss only widens/clarifies a serial builder's freedom). No ratchet added, no schema shape change.

*Summary written by work action (orchestrator)*

## Qualification

**Passed.** `qualify.sh` OK (2 files present + in diff, P-A-U all `[x]`, no debug artifacts), `scope-drift.sh` clean (2 declared = 2 touched). Orchestrator judgment: the diff adds substantive spec paragraphs (not whitespace); all 4 acceptance criteria trace to diff changes verified by read — Step 5.5 re-validation clause + partition-survives-mirror (work.md:337 + work-reference.md template), Step 1 enforcement-point pointer, Step 6 absent-set gloss, Step 4 warning-vs-enforcement reconciliation; the serial/floor path is behaviorally unchanged (every clause gated to parallel dispatch; the absent-set gloss only widens the serial builder's freedom). Check 6 (data flows) N/A — instruction prose.

*Verified by work action (orchestrator)*

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Contract regression checks passed (exit 0) — all ratchets green, none added. Matches the Step 5.75 baseline; no regression.

**Red-green validation:** N/A — non-behavioral Markdown instruction change confined to the optional parallel-dispatch path. Proof is the green regression suite plus the direct-read verification that the serial path is untouched.

*Verified by work action*

## Review

**Overall: 98%** | 2026-07-29T08:05:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor, 1 nit — **both remediated inline** (see `## Remediation`). Full adversarial rigor (reviewer + independent contradiction-hunter); no Important findings surfaced, so no refuter round was needed. The hunter's design-soundness check confirmed the new rule is sound under one orchestrator's sequential Step 5.5: whichever co-dispatched REQ firms its Scope *last* reads the other's already-firmed `write_set` and catches a firming-introduced overlap before its builder starts (no double-write), and no third unguarded write path to `write_set` exists (capture seeds pre-dispatch; crash-recovery clears; Step 5.5 mirror now guarded; Step 6 extension already guarded; board never writes it).
**Acceptance:** Pass — `contract-regressions.sh` green; all 5 edit sites present, cross-references resolve, serial/floor path unchanged.
**Follow-ups created:** None.

*Reviewed by review-work action (pipeline mode, full adversarial rigor per the session's calibration)*

## Remediation

The review found no Important defects, but flagged two coherence gaps in REQ-036's *own* newly-added sentences — precisely the "two-places framing clash" this coherence REQ exists to eliminate. Both were fixed inline (same file, in scope, trivial), rather than shipping self-inconsistent instructions plus a report:

- **Minor — `actions/work.md`** — the Step 1 gate pointer said "This gate is the initial *scheduling* decision; Step 5.5 is the **dispatch-time enforcement point**" (demoting Step 1 to non-enforcement), while the Step 4 reconciliation sentence credited Step 1 *as* an enforcement point. Reworded the Step 1 pointer to "Both steps enforce, at different times: this gate serializes or partitions overlapping pairs on the capture-seeded hint, and Step 5.5 **re-enforces** after each REQ firms its real set" — so both sentences now tell one story.
- **Nit — `actions/work.md`** — the Step 5.5 clause labelled the referenced Step 6 check "the same **no-concurrent-claimant** check," which conflates with the orchestrator lock's `claimed_reqs` *claimant* bookkeeping (a collision REQ-035 just made more prominent). Relabelled to "the same **disjointness re-check**" — the check is about a file/write-set claim, not a lock claim.

**Re-verification:** `contract-regressions.sh` exit 0; `qualify.sh`/`scope-drift.sh` clean (still the 2 declared files); grep confirms both enforcement sentences now agree and the `no-concurrent-claimant` label is gone.

*Remediated by work action (orchestrator)*

## Lessons Learned

**What worked:** Gating every new clause to the parallel-dispatch path (explicit "a serial run skips this entirely" hedges) made the serial/floor invariant checkable by grep, not just by argument. The adversarial hunter's ordering analysis — proving the *last* REQ to firm its Scope always catches a firming-introduced overlap — is the kind of soundness check a requirements-walk alone would skip.

**What didn't:** Writing two cross-referencing sentences in *different* steps (the Step 1 gate pointer and the Step 4 reconciliation) produced a micro-incoherence about whether Step 1 "enforces" — a coherence REQ briefly re-introducing the exact defect class it targets. Sentences that answer the same reader question ("which step enforces?") must be drafted side-by-side, not one-per-step.

**Worth knowing:** After REQ-035, "claimant" is a loaded word — it now means a lock `claimed_reqs` holder. Write-set/file-overlap checks should say "disjointness" or "no other REQ claims this file," never "claimant," to avoid conflation.

## Orientation

Step 5.5's `write_set` mirror now **re-validates disjointness on co-dispatch** before replacing the field, closing the second unguarded write path that REQ-032's gate left open (capture-seeded disjointness could be silently broken when Scope firmed the set), and a dispatch-time partition directive now survives the mirror. Lives in the work pipeline's Step 1 gate / Step 4 plan-validation / Step 5.5 mirror / Step 6 write-boundary bullet (`actions/work.md`) and the Scope Declaration Template (`actions/work-reference.md`). A leaf extension of the existing optional parallel-dispatch path — no new module, no data-flow or contract change, so no `[MAP CHANGED]`. Builds on REQ-035's `claimed_reqs`; the sibling REQ-037 (worktree merge placement) is independent of it.

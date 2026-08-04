---
id: REQ-102
title: Scope work.md Step 10 preserve rules to every non-own entry, and pin both label-destruction paths
status: completed
created_at: 2026-08-04T20:08:59Z
claimed_at: 2026-08-04T20:10:49Z
completed_at: 2026-08-04T20:15:26Z
route: A
kb_status: pending
user_request: UR-018
addendum_to: REQ-094
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
review_generated: true
related: [REQ-094]
batch: parallel-building
write_set: [actions/work.md, _dev/tests/contract-regressions.sh]
---

# Scope Step 10 Preserve Rules to Every Non-Own Entry

## What

REQ-094's review (Important finding) found that `actions/work.md`'s two Step 10 preserve rules — the wholesale-rewrite clause (~line 637) and the session-start delete clause (~line 647) — are scoped to entries "carrying another checkout's `writer:` label", which silently excludes the **label-less report-only** case. A label-less entry in a clean, committed checkpoint survives Crash Recovery's report-only branch (`actions/work-reference.md`, label-less legacy bullet) but then satisfies "no entry carrying another checkout's label remains", so the session-start delete removes it — and the next run classifies that `working/` REQ as "not named there" and ages it into the three-hour takeover ladder, which is exactly what the report-only branch refused.

## Detailed Requirements

- Change both `actions/work.md` clauses to scope preservation to **every entry this checkout did not write** (own-label entries are the only removable/enrichable ones), matching `actions/work-reference.md`'s canonical Step 10 clause ("enriches only this checkout's own entries, and carries every other one through verbatim").
- Add a contract assertion to `_dev/tests/contract-regressions.sh` pinning **both label-destruction paths**: the Step 10 preserve-foreign clause and the session-start scoped delete (REQ-094's review noted neither is pinned — a later "simplify the checkpoint rewrite" pass would reopen the hole with the suite green). Follow the file's existing assertion idioms.
- Suite must stay green (`bash _dev/tests/contract-regressions.sh` exit 0).

## Red-Green Proof

**RED prompt/case:** A clean committed checkpoint holds one label-less legacy entry; following `actions/work.md`'s session-start step verbatim authorizes deleting the file (no *labeled* foreign entry remains), destroying the entry Crash Recovery had classified report-only. Also: grep the two clauses — neither the preserve rule nor the scoped delete is pinned by any suite assertion.
**Why RED now:** The Step 10 echoes paraphrased the canonical condition instead of quoting it.
**GREEN when:** Both clauses read "every entry this checkout did not write" (or equivalent quoting the canonical condition); a new suite assertion fails if either clause loses its preserve language; suite green.
**Validation:** Review-generated (REQ-094 review, Important #1).

---
*Source: REQ-094 review, Important finding #1*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the exact files and clauses; the fix is quoting the canonical condition at two echo sites plus one new suite assertion in an established idiom. Review-generated with a precise failure narrative.

**Planning:** Not required

---

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified) — rescoped both Step 10 label-destruction clauses (2 lines).
- `_dev/tests/contract-regressions.sh` (modified) — one sed-extracted block + two assertions pinning both clauses.

**Requirements traced:**
- [x] Wholesale-rewrite preserve clause (Step 10, Session Checkpoint) now scopes preservation to "every `## In Progress (interrupted)` entry this checkout did not write", naming both qualifying cases (foreign `writer:` label, no label at all); enrichment stays scoped to this checkout's own labeled entries.
- [x] Session-start delete clause (step 3) now gates file deletion on "no entry this checkout did not write remains", with the label-less case named as a claim of unknown origin that recovery already refused to touch.
- [x] Both label-destruction paths pinned in `_dev/tests/contract-regressions.sh` (own comment block explaining the regression: a later "simplify the checkpoint rewrite" pass reopening the label-less hole with the suite green).
- [x] Assertions proven non-vacuous — both FAIL against `git show HEAD:actions/work.md` (full suite exits 1 with exactly the two REQ-102 FAIL lines).
- [x] `bash _dev/tests/contract-regressions.sh` exits 0.

## Decisions

- **D-01 (DECIDE & STATE):** Pinned phrase is `did not write` (as "entry this checkout did not write through verbatim" / "no entry this checkout did not write remains") rather than a quote of `actions/work-reference.md`'s "enriches only this checkout's own entries". The preserve rule's load-bearing half is the *carry-through* condition, not the enrichment condition, and both echo sites state the carry-through — pinning the enrichment half would leave the delete clause unpinned.
- **D-02 (DECIDE & STATE):** Two assertions over one shared block (`sed -n '/^#### Session Checkpoint/,/^## Clarify Questions/p'`) instead of one combined assertion. Both clauses live inside that single range, but each path fails for its own reason and deserves its own FAIL message naming which clause drifted — matching the file's existing one-requirement-per-assertion idiom (the REQ-077/REQ-094 groups).
- **D-03 (DECIDE & STATE):** Did *not* add an `assert_block_not_contains` for the old "carrying another checkout's `writer:` label" phrasing. That wording legitimately survives in the same block at session-start step 2 (which classifies entries and correctly distinguishes unlabeled from foreign-labeled), so a negative assertion would either fire on correct text or need a fragile longer anchor. The two positive assertions already fail if either clause loses its scoping.
- **D-04 (SILENT→stated):** Kept both edits as in-place rewordings of the existing sentences; no paragraph restructuring, per the REQ's surgical constraint.

## Discovered Tasks

- **Third echo of the same labeled-only scoping, outside this REQ's write boundary.** `actions/work-reference.md`'s Session Checkpoint Template carries the hole in its inline template comment (~line 806): `any entry carrying another checkout's label is copied through verbatim`. The canonical prose 15 lines below it (~821) is correct (`carries every other one through verbatim`), so the file contradicts itself — an agent copying the template comment reproduces exactly the label-less bug REQ-102 just closed at the two `actions/work.md` echo sites. One-line reword; not touched here because `actions/work-reference.md` is not in this REQ's `write_set`.
- **Checked and clean (no action needed):** `actions/cleanup.md` Pass 0 (line 40) and `actions/forensics.md` Check 1 (line 39) state their rule positively — *remove this checkout's own-label entry only* — so a label-less entry is out of scope there by construction. They mention the foreign label only as an illustration, not as the scope test, so they do not carry the hole.

*Orchestrator disposition: the third-echo find is routed into queued REQ-096's addendum (that REQ owns `actions/work-reference.md` and already carries a Step-10-template reword item) rather than a new pending-answers REQ.*

## Qualification

Passed — 2 files verified in working diff, both clauses match the canonical carry-through condition, assertions confirmed present and non-vacuous (builder's pre-edit run: suite exit 1 with exactly the two new FAIL lines; orchestrator re-ran post-edit suite: exit 0).

## Testing

- `bash _dev/tests/contract-regressions.sh` → exit 0 (orchestrator-run).
- Red-green: RED proven against `git show HEAD:actions/work.md` (two FAIL lines, nothing else moved); GREEN on working tree.

## Review

**Overall: 96%** | 2026-08-04 | Route A quick scan (orchestrator)

Both echo sites now quote the canonical condition ("entry this checkout did not write") and enumerate the two qualifying cases; session-start step 2's classification prose correctly retains the labeled/unlabeled distinction (D-03's reasoning for skipping a negative assertion is sound — the old phrase legitimately survives there as classification, not scope). Assertions follow the file's one-requirement-per-assertion idiom with a regression-explaining comment. **Findings:** none blocking; 1 discovered (third echo in work-reference.md's template comment) routed to REQ-096. **Acceptance: Pass.**

*Reviewed by orchestrator (Route A quick scan)*

## Lessons Learned

**What worked:** Proving assertion non-vacuity by running the full suite against `git show HEAD:<file>` and diffing the FAIL set — exactly two new lines, nothing else moved.
**Worth knowing:** This closed the *second and third* copies of a scoping condition whose first copy was canonical — and turned up a fourth (the template comment in work-reference.md, routed to REQ-096). Echo sites that paraphrase a canonical condition drift; echoes should quote it.

## Orientation

Step 10's checkpoint rewrite and session-start delete now preserve every in-progress entry this checkout did not write — label-less legacy entries included, closing the takeover-ladder re-entry hole from REQ-094's review. Contract suite pins both destruction paths.

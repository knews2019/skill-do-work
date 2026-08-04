---
id: REQ-102
title: Scope work.md Step 10 preserve rules to every non-own entry, and pin both label-destruction paths
status: pending
created_at: 2026-08-04T20:08:59Z
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

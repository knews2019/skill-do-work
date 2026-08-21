---
id: UR-019
title: Triaged downstream sync-review findings (0.170.1 → 0.174.3)
created_at: 2026-08-05T09:43:47Z
requests: [REQ-105, REQ-106, REQ-107]
word_count: 340
---

# Triaged Downstream Sync-Review Findings (0.170.1 → 0.174.3)

## Summary

A consumer repo pulled the skill's 0.170.1 → 0.174.3 diff and reviewed it (Go tests pass; findings are doc-level). The user routed the findings through `do-work validate-feedback` with instructions to verify each against this repo before queueing. All three findings verified as real and were accepted; the "context, not a REQ" item (missing `adr-018-*` / archived REQ-095 downstream) resolved as a consumer-side partial sync — both files exist here — so no fourth REQ.

One reviewer sub-claim was corrected during triage: the 0.172.0 CHANGELOG entry does **not** say "seeded by capture" — that claim lives only in `actions/work-reference.md` (the `assigned_to` schema line). REQ-105 is scoped accordingly.

## Extracted Requests

| REQ | Finding | Verdict at triage |
|-----|---------|-------------------|
| REQ-105 | `assigned_to` capture seeding documented but not implemented | Accept — add earmark detection + seeding to capture (reviewer-preferred direction) |
| REQ-106 | Auto-wave readiness predicate contradicts targeting-token provenance | Accept — add provenance carve-out to the Auto-wave list, mirroring the serial rule |
| REQ-107 | board.js assigned-badge comment over-claims ("nothing trims" above `truncateBadgeText`) | Accept — cosmetic comment reword |

## Batch Constraints

- Findings originate from a third-party review of the consumer's copy; each was independently verified against this repo during triage (validate-feedback) before capture.
- The schema lock-step rule applies to REQ-105 and REQ-106: `actions/work-reference.md`'s `assigned_to` schema line, `tools/queue-kanban/model.go` parser comments, and capture instructions must agree in the same commit.
- REQ-105 and REQ-106 both touch skill instruction files but are additive/corrective, not removal passes — `maintenance: false` on all three.

## Full Verbatim Input

do-work validate-feedback: do-work capture — three findings from a downstream review of the 0.170.1 → 0.174.3 sync
(consumer repo pulled the skill and reviewed the diff; Go tests pass, findings are
doc-level). Please verify each against this repo before queueing — the reviewer only
saw the consumer copy.

1. `assigned_to` capture seeding is documented but not implemented.
   `actions/work-reference.md` (Request File Schema, `assigned_to` line) and the
   0.172.0 CHANGELOG entry both say the field is "seeded by capture when the user
   earmarks work", but `actions/capture.md` and `actions/capture-reference.md`
   contain zero mentions of `assigned_to` (grep confirms). An agent following the
   capture action can never seed it, so the advertised behavior doesn't exist.
   Fix either direction: add earmark detection + seeding instructions to capture
   (preferred — e.g. "leave this for cloud-alpha" in the request text), or drop the
   "seeded by capture" claim from the schema line and changelog. Keep the schema
   line and any parser comments in lock-step per the existing lock-step rule.

2. Auto-wave readiness predicate contradicts targeting-token provenance.
   `actions/work.md` says `--fan-out` "composes with everything that selects a set,
   `--wave N` and targeting tokens included", and elsewhere that an explicitly-named
   `REQ-NNN` bypasses `depends_on`. But `actions/work-reference.md` → Fan-Out
   Dispatch → Auto-wave defines the ready set with condition 2 "dependency-ready"
   stated unconditionally. Reading only the reference, an explicitly-named but
   dependency-blocked REQ is excluded from the wave; reading work.md, it's included.
   Decide which is intended and make both files say it — likely a provenance
   carve-out sentence in the Auto-wave list mirroring the serial scan's rule.

3. Cosmetic: `tools/queue-kanban/web/board.js`, assigned badge block — the comment
   says "nothing here folds, trims, or rewrites the value" directly above
   `truncateBadgeText(request.assignedTo, 18)`. The tooltip does carry the full
   value, so the behavior is fine; the comment over-claims. Reword to distinguish
   display truncation from value normalization.

Context, not a REQ: the consumer sync arrived without `decisions/records/adr-018-*`
and `do-work/archive/UR-018/REQ-095-two-clone-acceptance-run.md`, both of which the
CHANGELOG/work-reference cite. If they exist here, it's a consumer-side partial
sync (ignore); if they don't, that's a fourth finding — the 0.174.2 entry announces
an ADR that was never committed.

---
*Captured: 2026-08-05T09:43:47Z*

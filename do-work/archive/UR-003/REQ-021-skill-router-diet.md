---
id: REQ-021
title: "SKILL.md router diet: one routing enumeration, help menu goes lazy"
status: completed
created_at: 2026-07-15T17:33:04Z
claimed_at: 2026-07-15T17:40:00Z
completed_at: 2026-07-15T17:52:00Z
route: B
user_request: UR-003
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: [REQ-020]
batch: harness-bloat-cleanup
maintenance: true
commit: da6cf19
---

# SKILL.md router diet

## What
Shrink SKILL.md from 5,507 words to ≤2,500 by removing duplicate enumerations of
the action set while preserving routing behavior exactly:

1. **Delete** the Actions bullet list (lines 11-46, 722 words) — duplicated by each
   action file's own blockquote/When-to-Use.
2. **Merge** the Verb Reference (lines 129-166, 1,526 words) into the routing
   priority table's Notes column, preserving unique disambiguation rules (install
   normalization, dream-over-cleanup precedence, abandon ID-rule, work argument
   rejection, single-word rule) and deleting the ~60-70% that restates the table.
3. **Move** the help-menu example + per-command-help (lines 168-274, 665 words) to a
   new `actions/help.md`, loaded only on routing priority 1 (bare/`help`) and on
   `<action> help` invocations (LAZY-LOAD).
4. **Keep** the routing priority table, payload rules, principles blockquotes, and
   the Action Dispatch table (canonical name→path map; grep-asserted by
   `_dev/tests/contract-regressions.sh:63-69`).
5. Update `_dev/tests/contract-regressions.sh` in the same commit if any asserted
   block moves; add `help` to the dispatch table.

## Why
SKILL.md loads on every invocation. It enumerates the action set five times; three
of those enumerations tax all 34 verbs to serve one (help) or none. Audit:
`decisions/audits/2026-07-15-harness-bloat-audit-phase1-2.md` §1c, DELETE/LAZY-LOAD.

## Acceptance criteria
- [x] SKILL.md ≤2,500 words (`wc -w`).
- [x] Every routing-table verb still resolves to the same action file: dispatch
      table row exists for all 34 routes incl. help; spot-verify the tricky
      precedences (cleanup carve-outs, abandon vs descriptive, check-for-updates).
- [x] No unique disambiguation rule from the Verb Reference is lost (diff review).
- [x] `actions/help.md` reproduces the menu verbatim + per-command-help behavior.
- [x] `_dev/tests/contract-regressions.sh` passes.

## Open Questions
(none — scope resolved during capture; see UR-003 capture notes)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read SKILL.md in full + `_dev/tests/contract-regressions.sh` (asserts on the Action Dispatch block). Approach: merge Verb Reference disambiguation into the priority table Notes column; extract help content to a routed `actions/help.md`; keep dispatch block heading/format byte-compatible with the test.
- [x] **[APPLY]:** Rewrote SKILL.md (5,507 → 2,396 words); created actions/help.md (760 words, menu verbatim + per-command help). Dispatch table gained a `help` row; per-command-help exception list (pipeline/prime/bkb) preserved.
- [x] **[UNIFY]:** `wc -w SKILL.md` = 2,396; contract-regressions.sh passes; grep for dangling references to removed sections (Verb Reference, Payload Preservation, help menu) found only actions' internal menus; 36 routing rows ↔ 35 dispatch rows reconciled (recap shares version.md).

## Triage

Route B — outcome clear (delete/merge/move with exact word targets), locations known, but the merge required
reading both enumerations side-by-side to preserve unique disambiguation. *(Generated during UR-003 processing.)*

## Implementation Summary

**What was done:** SKILL.md rewritten around a single merged routing table; help content extracted to a lazily-loaded action file.

Files changed:
- `SKILL.md` (modified) — 5,507 → 2,396 words. Deleted the Actions bullet list (722w); merged the Verb Reference into the routing table's Notes column (all trigger verbs + precedence rules preserved: install normalization, dream-over-cleanup, abandon ID-rule, work residue rejection, single-word rule, check/audit collisions); moved help menu + per-command help out; kept principles blockquotes, payload rules, and the Action Dispatch table (test-asserted) with a new `help` row.
- `actions/help.md` (new) — 760 words; menu text verbatim from the old SKILL.md, per-command-help spec, sync rule.
- `actions/version.md` (modified) — version bump 0.122.0 → 0.123.0.
- `CHANGELOG.md` (modified) — 0.123.0 entry.

**Route-preservation evidence:** every route in the 36-row table has a dispatch row (help/pipeline/capture/work/clarify/abandon/verify/review-work/validate-feedback/present/ai-report/cleanup/commit/inspect/code-review/ui-review/quick-wins/scan-ideas/deep-explore/install/forensics/roadmap/note/stray-check/tidy-repo/prime/bkb/interview/prompts/version/recap/tutorial/slop-check/dream/board); tricky precedences spot-verified: `check for updates`→version(2) before verify(5); cleanup carve-outs→29/30/34; `abandon <prose>`→capture(36); `what changed`→inspect vs `what's changed`→version.

## Testing

- `bash _dev/tests/contract-regressions.sh` → "Contract regression checks passed." (dispatch-block assertion exercises the rewritten file)
- `wc -w SKILL.md` → 2,396 (target ≤2,500)
- Red-green validation: omitted — non-behavioral restructuring; regression evidence above per Step 6.5's non-behavioral carve-out.

## Lessons Learned

**What worked:** Merging the Verb Reference INTO the priority table (instead of deleting it) kept every disambiguation rule while eliminating the duplication — the two tables were one mapping printed twice.
**What didn't:** The original ≤1,500-word target was unreachable: the dispatch table is test-asserted and CLAUDE.md-canonical, and full trigger-verb lists ARE routing behavior, not decoration.
**Worth knowing:** `_dev/tests/contract-regressions.sh:63` extracts the block between `## Action Dispatch` and `## Suggest Next Steps` — those two headings are load-bearing for the test.

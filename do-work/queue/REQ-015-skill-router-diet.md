---
id: REQ-015
title: "SKILL.md router diet: one routing enumeration, help menu goes lazy"
status: pending
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: [REQ-020]
batch: harness-bloat-cleanup
maintenance: true
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
- [ ] SKILL.md ≤2,500 words (`wc -w`).
- [ ] Every routing-table verb still resolves to the same action file: dispatch
      table row exists for all 34 routes incl. help; spot-verify the tricky
      precedences (cleanup carve-outs, abandon vs descriptive, check-for-updates).
- [ ] No unique disambiguation rule from the Verb Reference is lost (diff review).
- [ ] `actions/help.md` reproduces the menu verbatim + per-command-help behavior.
- [ ] `_dev/tests/contract-regressions.sh` passes.

## Open Questions
(none — scope resolved during capture; see UR-003 capture notes)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:**
- [ ] **[APPLY]:**
- [ ] **[UNIFY]:**

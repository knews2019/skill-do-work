---
id: UR-021
title: Act on the census's two durable findings
created_at: 2026-08-05T15:53:39Z
requests: [REQ-111, REQ-112]
word_count: 4
---

# Act on the Census's Two Durable Findings

## Summary

Captured from `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` after the user asked what the audit was for and why it should merge. The answer agreed: the 169-row table is perishable (415 line-number citations into actively-edited files), but two findings are about the *absence* of a mechanism rather than about line numbers, so they don't rot. These are those two, captured so the audit's value survives the table going stale.

## Extracted Requests

| REQ | Finding | Census §  |
|-----|---------|-----------|
| REQ-111 | 7 of the 9 Schema Read Contract enum fields have no normalizer anywhere in the repo | §2 (`work-reference.md` Schema Read Contract row), §3 closing paragraph |
| REQ-112 | `frontmatter.go` has no CLI surface, so no prose step can call it — all ~95 prose frontmatter reads are hand reimplementations by construction | §1 structural fact 1, §4 candidate 1 |

## Batch Constraints

- **REQ-112 depends on REQ-111.** The subcommand's `--normalize` flag is only correct once the missing normalizers exist; shipping the surface first would expose a flag that silently no-ops on 7 of 9 fields.
- **Both must clear the compiled-tooling exception.** `actions/board.md` is the only capability allowed to *need* a compiler (ADR-016; `CLAUDE.md` → Shipped Tooling). Neither REQ may make a shell-floor action depend on Go — the accelerator shape (`next-req` / `next-version` / `now`) is the only permitted form: gated on an already-built binary, with the prose procedure documented as the floor.
- Neither REQ removes or narrows the skill's own instructions, so both carry `maintenance: false`.

## Full Verbatim Input

do 1-3

(Referring to the recommendation: (1) revert the two version bumps and both changelog entries, (2) capture 2 REQs for the findings that survive, (3) let the audit merge as evidence behind those REQs.)

---
*Captured: 2026-08-05T15:53:39Z*

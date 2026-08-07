---
id: UR-022
title: Move the census's residual value into the queue and stop maintaining the table
created_at: 2026-08-05T19:32:52Z
requests: [REQ-114]
word_count: 12
---

# Move the Census's Residual Value Into the Queue

## Full Verbatim Input

also you said you'll move this to a REQ decisions/audits/2026-08-05-shell-logic-in-prose-census.md

## Context

The user is holding me to an earlier framing: that the census's value is its findings, and findings belong in the queue rather than in a standalone audit document. Two of them already shipped (REQ-111, REQ-112). What remained was the audit's 169-row table plus three uncaptured extraction candidates.

The decay the audit warned about then arrived within hours: merging `origin/main` (which added ~25 lines to `actions/work-reference.md`) shifted every citation past that insertion point. Verified concretely — the audit claimed L790/L794–809/L821; the cited text now lives at L815/L822–826/L846. The table had begun to lie about the exact thing it existed to record.

## Batch Constraints

- The audit's path must survive: `REQ-111` and `REQ-112` are archived and immutable, and both carry `*Source:*` lines citing `decisions/audits/2026-08-05-shell-logic-in-prose-census.md`. Deleting the file outright would dangle two references in the permanent intent trail.
- Candidates must be restated **without** line-number citations, or they inherit the decay that motivated this request.

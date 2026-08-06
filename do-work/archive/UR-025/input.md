---
id: UR-025
title: Two Codex review findings on PR #133
created_at: 2026-08-06T11:26:17Z
requests: [REQ-119, REQ-120]
word_count: 96
---

# Two Codex Review Findings on PR #133

## Summary

Codex reviewed PR #133 (REQ-116/117/118) and raised two findings, both verified against the code before capture. One is a behaviour gap this branch created by fixing `domain` and not `route`; the other is a rule violation in shipped files that also pre-dates this branch.

## Extracted Requests

| REQ | Title | Codex severity |
|---|---|---|
| REQ-119 | An off-vocabulary `route` warns like `domain` does | P2 |
| REQ-120 | Shipped files stop citing the export-ignored maintainer doc | P1 |

## Batch Constraints

- Both are read-side only, and neither adds a write surface to the tool.
- REQ-119 must not blank an unrecognized route. REQ-116 deliberately reports it case-folded rather than substituting route's empty-string default, because the raw letter is the evidence re-triage reads. Adding the warning must not change that.
- REQ-120's check is `_dev/tests/contract-regressions.sh`'s maintainer-document probe. Its per-file allowlist exists for mentions of a *consumer project's* CLAUDE.md — none of the four hits qualify, so the fix is to restate inline, not to widen the allowlist.

## Full Verbatim Input

Codex review comment 1 (P2, `tools/queue-kanban/model.go:653`) — "Warn when a route remains outside the enum": When a ticket contains an unrecognized value such as `route: z`, this path changes it to `Z` but records no unrecognized flag and appends no board warning. That leaves the board out of lock-step with `actions/work-reference.md:191-198`, which requires every off-vocabulary enum read to emit a warning rather than merely display the bad value. Carry the recognition result and surface it through `board.Warnings`, as the adjacent domain handling does.

Codex review comment 2 (P1, `tools/queue-kanban/model.go:975`) — "Remove maintainer-doc references from shipped files": This new comment, together with the new `tools/queue-kanban/prime-do-kanban.md` REQ-116 lesson, cites the repository's export-ignored `CLAUDE.md` from shipped paths. I ran `_dev/tests/contract-regressions.sh`; its maintainer-document check now fails on both additions, and consumers receive dangling references. Restate the relevant rules inline or point to a shipped source instead.

---
*Captured: 2026-08-06T11:26:17Z*
*Provenance: both findings are third-party review output from the `chatgpt-codex-connector` bot on PR #133, verified against the code during capture. No instruction-like content was detected in either comment body. Codex's second comment states the check "now fails on both additions"; verification found it was **already** failing on `main` for two other citations in `tools/queue-kanban/frontmatter_cli.go` (from REQ-112), so the check cannot go green by fixing only this branch's two — REQ-120 covers all four.*

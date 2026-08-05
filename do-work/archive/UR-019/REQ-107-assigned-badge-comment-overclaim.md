---
id: REQ-107
title: Board assigned-badge comment over-claims — "nothing trims" sits directly above truncateBadgeText
status: completed
created_at: 2026-08-05T09:43:47Z
claimed_at: 2026-08-05T10:39:08Z
completed_at: 2026-08-05T10:40:04Z
route: A
user_request: UR-019
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set: [tools/queue-kanban/web/board.js]
related: [REQ-097, REQ-105, REQ-106]
batch: sync-review-0174
---

# Board Assigned-Badge Comment Over-Claims

## What

In `tools/queue-kanban/web/board.js`, the assigned-badge comment says the value is "Rendered verbatim — session names have no canonical vocabulary, so nothing here folds, trims, or rewrites the value" — directly above `truncateBadgeText(request.assignedTo, 18)`, which trims the badge text for display. The behavior is fine (the tooltip and drawer row carry the full value); the comment over-claims. Reword it to distinguish value normalization (none — verbatim-read class) from display truncation (badge clipped at 18 chars, full value in the tooltip and drawer).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reword only the six-line comment above the assigned badge: keep the never-normalized claim (true, and in lock-step with the schema's verbatim-read class), scope the no-rewriting claim to the *value*, and state explicitly that the badge's visible text is truncated for layout with the full value in the tooltip and drawer row. No code change.
- [x] **[APPLY]:** Comment replaced in `tools/queue-kanban/web/board.js`; no executable line touched.
- [x] **[UNIFY]:** `git diff --stat` shows only `tools/queue-kanban/web/board.js` beyond do-work/ bookkeeping, and the hunk is comment-only. `go build` and `go test ./...` pass in `tools/queue-kanban/`. No debug artifacts.

## Context

Cosmetic, comment-only — no behavior change, no schema change, so the `model.go`/work-reference lock-step is untouched. Found by a downstream consumer's review of the 0.170.1 → 0.174.3 sync; verified here at triage (comment at `tools/queue-kanban/web/board.js:533-538`, truncation call at `:542`, full-value tooltip at `:546`, drawer row at `:1722-1723`).

## Red-Green Proof
**RED prompt/case:** Reading `board.js:535` ("nothing here folds, trims, or rewrites the value") against `:542` (`truncateBadgeText(request.assignedTo, 18)`) — the comment contradicts the very next statement.
**Why RED now:** The comment conflates "no value normalization" (true) with "no display trimming" (false).
**GREEN when:** The comment states both halves accurately: the *value* is never normalized (verbatim-read class), while the *badge display* is truncated at 18 chars with the full value preserved in the tooltip and drawer row.
**Validation:** Inferred during capture (triage-verified against the repo)

## Full Context
See `do-work/user-requests/UR-019/input.md` for complete verbatim input.

## Triage

**Route: A (Direct to Builder)** — single named file, comment-only reword, behavior explicitly out of scope.

## Implementation Summary

**What was done:** Reworded the assigned-badge comment in `tools/queue-kanban/web/board.js` so it no longer claims "nothing here folds, trims, or rewrites the value" directly above a `truncateBadgeText` call. The comment now separates the two facts: the *value* is never normalized (no case folding, no alias maps — matching the schema's verbatim-read class), while the badge's *visible text* is truncated for layout, with the full value carried by the title tooltip and the drawer row.

**Files changed:**
- `tools/queue-kanban/web/board.js` (modified) — assigned-badge comment block only; no executable change

**Lock-step verification:** The comment's normalization claim still matches `actions/work-reference.md`'s verbatim-read contract and `model.go`'s display-only parse — neither needed an edit (no semantic change).

**Tests:** `go build` + `go test ./...` pass in `tools/queue-kanban/`.

## Qualification

Passed — 1 file verified (`tools/checks/qualify.sh` OK), requirement traced (the one requirement is the reword; the diff hunk is exactly that comment block), P-A-U confirmed against the diff.

## Testing

- `go test ./...` in `tools/queue-kanban/` — ok (2.8s). Comment-only change; regression evidence in lieu of red-green (non-behavioral, per the REQ's own framing).
- **Red-green validation (doc-level):** RED: comment said "nothing here folds, trims, or rewrites the value" three lines above `truncateBadgeText(request.assignedTo, 18)`. GREEN: comment now says only the badge's visible text is truncated and the full value survives in tooltip + drawer — verified against the code paths it describes (`makeBadge` call, `.title` assignment, drawer's `appendMetaRow("Assigned to", ...)`).

## Review

**Acceptance: Pass** (Route A quick scan, calibrated depth). The reworded comment's every claim was checked against the code it annotates: truncation at 18 chars (badge only), tooltip set from the untruncated value, drawer row untruncated, no normalization anywhere in the parse or render path. Wording stays consistent with the schema's verbatim-read language without restating it at length. Scope: exactly the declared file, exactly the comment. No findings.

## Orientation

The board's assigned-badge comment now tells the truth about truncation: value verbatim, display clipped, full value one hover away. Lives in the queue-kanban board frontend; no map change.

---
*Source: downstream sync-review finding 3, verified at triage — see UR-019*

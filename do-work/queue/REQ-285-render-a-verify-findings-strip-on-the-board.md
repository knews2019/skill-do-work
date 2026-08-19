---
id: REQ-285
title: Render a verify-findings strip on the board
status: pending
created_at: 2026-08-19T13:47:06Z
user_request: UR-058
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-284]
maintenance: false
related: [REQ-284, REQ-286]
batch: verify-findings-on-board
write_set: [skills/do-work-board/tools/queue-kanban/web/template.html, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work-board/tools/queue-kanban/web/board.css]
---

# Render a Verify-Findings Strip on the Board

## What

Add a second always-visible strip to the board client, modeled on the existing completion-anomalies
strip, that renders `verifyFindings` and `verifySkipped` from the board payload.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Seeing the board should mean seeing the problems. Today the thirteen categories the data-warnings banner
does not cover are invisible to anyone who does not run `verify` from a shell.

## Context

`web/template.html:137-146` and `web/board-cards.js:412-431` are the completion-anomalies strip: a titled
section with a count, a hint line, and per-item cards, sitting outside the view panels so it survives
every view switch, the recently-done window, and the shared filters. That is the shape to copy, and its
visual weight is already calibrated.

## Detailed Requirements

- One new strip: title, count, one card per finding, the category as the badge, the remedy as the muted
  second line.
- Skipped probes render in a collapsed footer on the same strip. They must render — a skipped probe shown
  as nothing reads as "checked and clean".
- Reuse the `.board-anomalies-*` rules rather than writing a parallel palette.
- Hidden when there are no findings and no skips, exactly as the anomalies strip hides itself.
- Exempt from the shared filters and from the 24h/48h/7d window, for the same reason the anomalies strip
  is: a finding must not be hideable by a filter combination.
- Nothing in this REQ parses or re-derives a finding. The client renders the producer's list blindly;
  category suppression already happened in Go (REQ-284).

## Constraints

- Framework-free client, consistent with the rest of `web/`.
- Read-only. No write surface, no repair affordance, no link that mutates queue state.

## Dependencies

REQ-284 — the `verifyFindings` / `verifySkipped` payload fields must exist first.

## Builder Guidance

Firm on placement and shape, latitude on the exact wording of the title and hint line. Keep it visually
subordinate to the anomalies strip when both are showing; anomalies are broken bookkeeping, findings are
process drift.

## Red-Green Proof

**RED prompt/case:** Generate a board from a fixture queue carrying one REQ claimed more than 3h ago and
one `worktree-agent-*` leftover branch, open the generated `index.html`, and look at it.

**Why RED now:** the board renders neither. `board.Warnings` carries no claim-age or worktree entry, and
there is no strip that could show one.

**GREEN when:** the rendered page shows a findings strip with a count of 2, one card reading the
`claim-needs-attention` detail with its remedy underneath, one card for the worktree leftover, and the
strip stays on screen after switching views and after applying a filter that hides the claimed column.

**Validation:** User confirmed (accepted as F1, F7, F8 in the `do-work validate-feedback` triage).

## Full Context

See `do-work/user-requests/UR-058/input.md` for the complete verbatim input and the triage verdicts.

---
*Source: upstream suggestion for `knews2019/skill-do-work`, observed against v0.212.25 — "Suggested shape" item 3.*

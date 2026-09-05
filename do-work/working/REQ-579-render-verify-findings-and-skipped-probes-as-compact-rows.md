---
id: REQ-579
title: 'Render verify findings and skipped probes as compact rows in one list'
status: claimed
created_at: 2026-09-05T00:19:58Z
user_request: UR-118
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-580, REQ-578]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
claimed_at: 2026-09-05T12:21:55Z
---

# Render Verify Findings and Skipped Probes as Compact Rows in One List

## What

The Verify Findings strip (`#board-findings`, added by REQ-285) renders each finding as a bordered card in a multi-column grid and lists skipped probes as bullets inside a collapsed disclosure below. Replace both with one flat list of compact rows: one row per finding and one row per skipped probe, no cards, no columns, no separate disclosure. A finding and a skipped probe are the same kind of thing to the reader ("verify has something to tell you"), so they get one visual language.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user's words: "for a warning I like the small lines, the big boxes are using up too much space for something not that important". The cards also borrow the REQ-card shape, so they read as work items, which REQ-482 (stack verify-finding cards full width, now cancelled in favour of this REQ) had already noticed. The contradiction being removed is two visual languages for one idea: a bordered card for a finding, a bullet in a disclosure for a skipped probe.

## Detailed Requirements

- **D1, one list of rows.** The header line stays: "Verify findings", the finding count, and the hint. Add the skipped count to the header when it is non-zero (for example "2 findings · 2 probes not checked"). Below the header, one plain list. Each finding row carries: a small uppercase category chip (as today's `.board-finding-category`), the detail sentence, and the remedy in muted text on the same line after the detail. Long text wraps. No border, no card background, at most a thin divider between rows. Each skipped-probe row carries a "not checked" chip and the probe text in muted weight, in the same list, after the findings. Delete the `<details>` disclosure and the `board-anomalies-cards` grid from the strip's markup; the `.board-finding*` and `.board-findings-skipped*` classes are replaced by the row classes, not kept beside them.
- **D2, two weights, from the producer only.** A finding with `fixable: true` renders with a muted "cleanup can fix" tag and lighter text, because a command resolves it. Every other finding renders at normal weight. A skipped probe renders muted. Do not add any severity, colour scale or ordering in the client that the payload does not carry; `fixable` is the only weight (the same rule as today's comment in `renderVerifyFindingsStrip`: never inferred here, the producer sets the flag).
- **D3, grouped by subject.** Add a `Subject` string to `VerifyFinding` in verify.go and a `subject` field to the board payload in generate.go. The producer sets it to the thing the finding is about: the worktree name for the worktree probes, the REQ id for REQ-level probes, the changelog path or version for the release probes. A probe with no natural subject leaves it empty. The client groups rows with the same non-empty subject together and prints the subject once as a small heading line above its rows; rows with an empty subject render as single rows in producer order, after the grouped ones. Grouping is by exact string match on the payload field, never by parsing the detail text.
- **D5, hide rules unchanged.** The strip stays hidden when there are no findings and no skipped probes. REQ-578 (hide the strip on the Activity view) keeps working against the same `#board-findings` element and `hidden` attribute; do not restructure the strip in a way that breaks it.
- The strip is not a REQ card. Do not reuse `board-anomalies-cards` or any `.board-request*` class for finding rows.
- Board changes follow `_dev/primes/prime-kanban-board.md` (versioning, parser lock-step, build outputs). The generated payload gains one optional field, so the snapshot fixtures under the queue-kanban tests need regenerating per that prime.

## Constraints

- REQ-580 (stop probing committed queue state for an undetermined worktree) edits the same probe in verify.go and its test. Both REQs may add a field or a branch around `appendWorktreeFindings`; declare the overlap, do not serialise on it.
- The completion-anomalies strip above the findings strip is not part of this request and keeps its cards.

## Builder Guidance

The user is certain about the outcome (rows, not cards) and approved D1 to D3 and D5 as written. Latitude: spacing, divider style, chip typography and the exact header wording are the builder's call, as long as the strip reads as a list of warnings and not as a set of cards. Load `crew-members/anti-slop.md` before styling; the result is a human-facing artifact.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane (`javascript_behavior_c_test.go` pattern), render with two verify findings sharing `subject: "worktree-agent-REQ-506-focused-evidence"` (one `fixable: true`, one not) and one skipped probe in `boardData`. Assert: `#board-findings` contains no `.board-finding` card and no `details` element; it contains exactly three row elements in one list; the two findings sit under one subject heading carrying that worktree name; the fixable row and the skipped row carry the muted class and the other finding row does not. In Go, assert that `generatedVerifyFinding` serialises a `subject` field for a worktree leftover finding.
**Why RED now:** `renderVerifyFindingsStrip` in board-cards.js builds one `.board-finding` card per finding inside the `board-anomalies-cards` grid and puts skipped probes in a `<details>` list; `VerifyFinding` has no subject and the payload carries none.
**GREEN when:** The assertions above pass, the existing strip hide-when-empty test still passes, and the board screenshot shows the strip as a flat list of rows with the skipped probes in it.
**Validation:** User confirmed: "for a warning I like the small lines" and "ok, do it" to the D1 to D5 proposal; screenshot 1 shows the RED state.

## Assets

- `do-work/user-requests/UR-118/assets/REQ-579-screenshot-1-verify-findings-two-cards-and-skipped-probes.png`: queue-kanban board, light theme, Board view, 2026-09-05. Strip header "VERIFY FINDINGS 2". Two white bordered cards side by side, both labelled WORKTREE-MERGE-STATE-UNDETERMINED, for `worktree-agent-REQ-506-focused-evidence` and `worktree-agent-REQ-577-launcher-fixture`, each with the detail "git could not say whether this is merged ... inspect it by hand; cleanup Pass 5 deletes nothing it cannot establish a merge target for". Below them an expanded disclosure "2 probe(s) could not run — unverified, not clean" with two bullets, each "committed-queue-state probe for <worktree>: `git diff main...<branch> -- do-work/` failed (no such branch, or unrelated histories)". The cards use roughly a third of the strip's width each and most of its height; the two bullets take two lines.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views" and "static output". Over the budget on its own.

*Source: "for a warning I like the small lines, the big boxes are using up too much space for something not that important" / "how to make it good, beautiful and non-contradictory?" / "ok, do it"*

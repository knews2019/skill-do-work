---
id: REQ-388
title: '[impact-rule-change] Settle the last two drawer/clipboard divergences: fence info strings and ids inside paths'
status: pending
created_at: 2026-08-26T23:08:00Z
status_changed_at: 2026-08-26T23:08:00Z
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-386]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/citations.go
  - skills/do-work-board/tools/queue-kanban/citations_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
---

# Settle The Last Two Drawer/Clipboard Divergences

## What

Two places where the drawer's glossary and the paste's appendix list different ids for the same body.
Decide which surface is right in each case and make them agree.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-383's stated rule is that the drawer and the paste say the same thing about the same body. After
its review fixes, two cases still break it in opposite directions.

**D-A — a fence info string.** `collectMentionSurfaces` registers a fenced block's `Info` segment as a
surface, so a resolved id in an info string (` ```yaml REQ-1679-template `) is emitted and lands in
the paste's appendix. goldmark renders the info string only as `class="language-…"`, never as a text
node, so the drawer's glossary has no matching entry. **The paste says more than the drawer.**

**D-B — an id nested in a repo-relative path.** The file-path alternative claims the whole run and Go
never re-scans inside it. `board-detail.js` deliberately rewinds `lastIndex` by one on a skipped match
so a REQ id inside a skipped path still links, so the drawer glossaries an id the appendix omits.
**The drawer says more than the paste.** Two real instances: REQ-112 and REQ-239. It disappears in
serve mode, where the path becomes a file link instead.

## Context

**D-B is a deliberate, documented difference, not an accident.** The clipboard's own comment in
REQ-379 said so: *"A URL or a repo-relative path is one opaque run. The drawer resumes INSIDE it so a
nested id can still link; a clipboard payload must not, or an expansion lands mid-path and the paste
carries a path that no longer names a file."* That reasoning is about the EXPANSION, which would
corrupt the path — it says nothing about the appendix line, which is additive and cannot corrupt
anything. So the likely resolution is asymmetric: never expand inside a path on either surface, but
let both record the glossary entry.

**D-A is a genuine question.** An id in an info string is arguably an illustration (the reader never
sees it) or arguably a real reference (the author typed it deliberately).

## Detailed Requirements

- **Record a `## Decisions` entry for each**, naming which surface changes and why. These are rule
  changes, not bug fixes, and the next person needs the reasoning more than the diff.
- **Whatever is decided, one test must drive BOTH surfaces over the same body** and compare their two
  reference lists directly. Two tests asserting two expectations is how these diverged.
- **An expansion must still never land inside a path** on either surface — that part of REQ-379's
  reasoning stands regardless of what the appendix does.

## Constraints

- **No new board write surface.**
- Do not widen into REQ-385, REQ-386 or REQ-387.
- **`depends_on: [REQ-386]`** because both edit `citations.go` and `citations_test.go`. The edge
  serializes a shared file; it is not a need for REQ-386's output. `write_set` gates nothing — only
  `depends_on` does — and this batch has already shipped that mistake twice.

## Dependencies

`depends_on: [REQ-386]`, which shares `citations.go` and `citations_test.go` with this REQ.
Transitively this also orders it after REQ-381 and REQ-385 — REQ-385 matters, because the two also
share `web/board-detail.js`.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Red-Green Proof

**RED prompt/case:** Open REQ-112 in a STATIC board, read the drawer's Referenced tickets list, then
Copy it and read the appendix — the drawer lists REQ-110 and the appendix does not. Second case: a
body containing ` ```yaml REQ-1679-template ` — the appendix lists REQ-1679 and the drawer's glossary
does not.

**Why RED now:** the two surfaces derive their reference lists from different inputs (rendered DOM vs
source bytes) with no test comparing them.

**GREEN when:** for both bodies the drawer glossary and the paste appendix list the same ids, and one
test asserts that equality rather than two tests asserting two expectations.

**Validation:** Both reproduced by adversarial review of REQ-383; the two real instances of D-B were
enumerated across all 453 documents.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, findings S3, S4 and C3.*

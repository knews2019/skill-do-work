---
id: REQ-054
title: "Write docs/board-guide.md — user-facing guide to the Kanban board"
status: completed
route: B
claimed_at: 2026-07-29T13:55:51Z
commit: 50caea3
domain: general
tdd: false
maintenance: false
prime_files: []
created_at: 2026-07-29T09:32:07Z
user_request: UR-007
addendum_to: REQ-034
depends_on: []
write_set: ["docs/board-guide.md"]
---

# Write docs/board-guide.md — user-facing guide to the Kanban board

## What

Create `docs/board-guide.md`, a short human-facing tour of the board: modes (serve / static / summary), columns, badges (including the write-set `overlaps` badge and glob dialect caveats), the card drawer, and the Testing view. Follow the structure and tone of the existing `docs/*-guide.md` files.

## Why

Board features are currently documented only in agent-facing `actions/board.md`. The builder recommended against a guide (ongoing drift surface), but the user overrode via `do-work clarify` on 2026-07-29: a linkable feature tour is wanted while the feature set is fresh — and `docs/` already carries a per-action guide convention this slots into. (Follow-up from REQ-034, surfaced in REQ-041.)

## Constraints

- Keep it short — the drift-surface concern was real; cover what a user sees, not internals. Load `crew-members/anti-slop.md` before writing (human-facing artifact).
- Don't duplicate `actions/board.md`'s agent instructions; link roles are: guide = human tour, action file = agent contract.

## Acceptance

- `docs/board-guide.md` exists, covers modes/columns/badges/Testing view, and reads like the sibling guides.

## Triage

**Route:** B
**Reasoning:** New human-facing doc (`docs/board-guide.md`); survey the board's modes/columns/badges/Testing view and follow the sibling `docs/*-guide.md` structure. Human-facing artifact -> load `crew-members/anti-slop.md`. Route B.
**Rigor:** Standard main-context review (part of the parallel disjoint-write_set batch 051/052/054/057/058; single-file, no spec-cluster overlap).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Scope

Files I will touch: `docs/board-guide.md`

## Plan-Act-Unify

- [x] **PLAN** — Survey `actions/board.md` for the feature facts (modes, columns, badge set, drawer rows, Testing view), read the shipped frontend (`tools/queue-kanban/web/template.html`, `web/board.js`) for the labels a user actually sees, and match the structure/tone of `docs/roadmap-guide.md` + `docs/cleanup-guide.md` (H1 + one-line summary → what you see → key rules → Usage).
- [x] **ACT** — Wrote `docs/board-guide.md`: Modes table, Board view (columns + Notes/anomalies strips + toolbar), Badges table, a `Reading the overlaps badge` subsection carrying the four glob caveats, Card drawer, Testing view, Usage.
- [x] **UNIFY** — `docs/board-guide.md` created (new file, 788 words, 68 lines) and nothing else touched. Anti-slop self-check: no throat-clearing intro, no marketing tone, no self-grading; every section is something a user sees in the UI, with internals (Go source, `annotateWriteSetOverlap`, API endpoints, `--repo-root`) deliberately left to `actions/board.md`. Compressed twice — the drawer's exhaustive frontmatter-row enumeration was collapsed to a labelled summary because it was the file's highest drift surface for no reader gain. No agent instructions from `actions/board.md` are duplicated: the guide describes behavior in second person, the action file keeps the build/dispatch contract.

## Implementation Summary

**Files changed:**
- `docs/board-guide.md` (new) — human-facing board tour (modes, columns, badges, Testing view) following the sibling `docs/*-guide.md` shape

**What was done**

Wrote the human-facing board tour the REQ asked for, following the sibling `docs/*-guide.md` shape: H1 + a one-line read of what the board is and its read-only posture, a Go-toolchain note as a blockquote, then a Modes table (serve / static / summary, each with its default and its override flag), the Board view's four columns and the Ready/Waiting split inside Pending, the always-visible Notes and Completion-anomalies strips, the toolbar (filters, recently-done window, Columns/By-UR lens, Board/Calendar/Testing switch), a Badges reference table, the card drawer, the Testing view's four columns and per-card actions, and a Usage block closing on the `just-kanban` shortcut.

The `overlaps` badge got its own subsection because it is the one badge a user can misread as a safety guarantee. It states the badge's role (scheduling heads-up; `do-work run`'s gate makes the co-dispatch call) and then its four under-report modes: no-badge ≠ safe when no `write_set` was declared, `*` never crossing `/` and `**` not being recursive, a malformed pattern matching nothing except its own identical twin, and a directory entry never badging a file inside it.

Boundary held throughout: feature facts sourced from `actions/board.md` and the shipped frontend labels, but none of the action file's agent contract (build steps, `--repo-root`, `.git/info/exclude`, `/api/testing/*` internals, the Go module layout) restated — guide is the human tour, action file is the agent contract.

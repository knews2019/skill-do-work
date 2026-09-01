---
source_type: req_lesson
req_id: REQ-097
req_path: do-work/archive/UR-018/REQ-097-assigned-to-advisory-field.md
date: 2026-08-04
domain: backend
module: actions
tags: [backend, assigned, advisory, field, schema]
---

# Lessons from REQ-097: assigned_to advisory field — schema line, scan skip-and-report, board parse (lock-step)

## What the REQ was about

Add the cooperative claim marker the whole model rests on: a single advisory frontmatter field, `assigned_to: "<session-name>"`, on a pending REQ. The default work-loop scan skips-and-reports such REQs; explicitly targeting one (`do-work run REQ-NNN`) overrides and clears the field. The board parses it display-only. **No verb, no status, no staleness clock, no release command** — the 0.163.0 forbidden-token ratchet stays fully intact.

## Solution summary

Added the advisory cooperative claim marker end to end — schema, one courtesy reader in the work loop, and display-only board rendering — with the parser and the schema line in the same commit per the lock-step rule. No verb, no status, no timestamp, no clock.

## What worked

- Writing the three Go tests first genuinely paid, because the RED was a *compile* failure — which is the strongest possible RED for a new field, and it proves the test is actually reaching the code under test rather than passing vacuously.
- Making the verbatim fixture hostile (`Cloud-Alpha_2` instead of the plan's `cloud-alpha`). The plan's value passes with or without normalization, so it would not have tested the contract it was written for.

## What didn't work

- The first board smoke reported `assignedTo` absent from the payload and looked like a threading bug. It was a **stale compiled binary** — the batch handdown warns about exactly this, and the warning was still not enough to stop it happening. The tell was `grep` finding the JS identifier but never the data value.
- The second smoke then grepped `index.html`, which never carries the data: `generate` writes the payload to a sibling **`board-data.js`** (plus a lazy `board-markdown.js`). Two false negatives in a row, neither of them a code defect.

## Worth knowing

- `generate --out <dir>` produces three files. **`board-data.js` is the payload** — assert against that, never `index.html`.
- There is exactly **one** payload construction site (`generate.go`'s `generatedRequest` literal), shared by `serve` and `generate`, so a new display field needs one copy line, not two. Confirm by grepping the `WriteSet:` copy before assuming otherwise.
- `omitempty` on the payload string is load-bearing, not tidiness: it is what preserves the *absence reads as unknown* convention the board applies to `write_set`.
- The `.badge-*` rules use `--ink-soft` / `--surface-3` / `--line-firm`; `--text-muted`, `--surface-raised` and `--border-subtle` do not exist in `board.css`. A wrong variable name fails silently — the badge just renders with the base `.badge` styling.

## Back-reference

See `do-work/archive/UR-018/REQ-097-assigned-to-advisory-field.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `13328a8`.

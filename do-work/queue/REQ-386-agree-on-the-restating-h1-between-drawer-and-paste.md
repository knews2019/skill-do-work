---
id: REQ-386
title: '[impact-user-visible] Make the drawer and the paste agree about a body H1 that restates the title'
status: pending
created_at: 2026-08-26T23:06:00Z
status_changed_at: 2026-08-26T23:06:00Z
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-381]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/citations.go
  - skills/do-work-board/tools/queue-kanban/citations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
---

# Make The Drawer And The Paste Agree About A Body H1 That Restates The Title

## What

The drawer deletes a body's opening H1 when it restates the frontmatter title, then decides which
mention expands. The clipboard keeps that H1 and counts it as the first prose mention. Pick one rule
and apply it to both.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`linkifyDetailBody` (`web/board-detail.js`) removes the restating H1 before linkifying, so the drawer
expands the first mention in the BODY. The Go walk sees the whole file and expands the first mention
in the H1. Three real records hit this today:

| Record | Its restating H1 |
|---|---|
| REQ-041 | Confirm: three small board/pipeline hardening follow-ups from REQ-034 |
| REQ-042 | Confirm: three worktree-mode evidence-path hardening follow-ups from REQ-037's review |
| REQ-085 | Run REQ-073's Live Two-Builder Acceptance Test and Record What It Found |

Two consequences, and the second is the one that matters:

- The two surfaces annotate **different occurrences** of the same id. The drawer expands the prose
  mention and the paste leaves it bare, expanding the heading instead.
- **The paste stops round-tripping.** Once the H1 reads `… from REQ-034 (Capture-time decomposition
  nudge…)` it no longer equals the record's `title:` field, so on save-back the drawer's exact-match
  H1 removal stops firing and the reader sees the title twice.

## Context

The H1-stripping rule is old and deliberate: REQ/UR bodies conventionally open with an H1 restating
the frontmatter title, and showing it inside a drawer that already displays the title is duplication.
`copyTextWithHeading` implements the same de-duplication on the clipboard's rendered-text FALLBACK
path — so the clipboard already knows the rule, and only its primary path ignores it.

Not a regression: REQ-379's scanner annotated the H1 too. REQ-383 moved the decision into Go and kept
the behaviour.

## Detailed Requirements

- **Decide which surface is right and say why in a `## Decisions` entry.** The two candidate rules:
  suppress annotation inside a restating H1 (the paste keeps a heading identical to its title, and
  save-back stays clean), or stop stripping the H1 in the drawer (both show it, no rule needed). The
  first preserves the round trip and is the smaller change; the second removes a rule instead of
  adding one, which this repo generally prefers. Argue it rather than assuming.
- **Whichever is chosen, the H1 comparison must be the SAME comparison on both sides.** The drawer
  uses `normalizeHeadingText` (collapse whitespace, trim, lowercase). Go needs the same normalization
  or the two disagree on a heading that differs only in spacing.
- **The round trip is the acceptance test**: copy one of the three records above, save the paste as a
  file, rebuild the board from it, and confirm the drawer still strips its H1.

## Constraints

- **No new board write surface**, and no change to what the frontmatter fence carries.
- Do not widen into REQ-385, REQ-387 or REQ-388.

## Dependencies

`depends_on: [REQ-381]`, which shares `citations.go`, `citations_test.go` and `generate.go` with this
REQ. Transitively this also orders it after REQ-385.

An earlier draft of this section claimed a "disjoint write set from REQ-385", which this REQ's own
frontmatter contradicted on two files. That was the third write-set/`depends_on` mismatch in this
batch, and it is why the `ungated-write-set-overlap` probe now exists.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Red-Green Proof

**RED prompt/case:** Open REQ-041 in the drawer and read where the title appears; then Copy it and
read the paste. The drawer expands the prose mention on the later line; the paste expands the one in
the H1 and leaves the prose mention bare.

**Why RED now:** `linkifyDetailBody` removes the H1 at `web/board-detail.js:208` before the mention
walk; the Go walk never sees that removal.

**GREEN when:** both surfaces annotate the same occurrence for all three records, and a pasted copy of
REQ-041 saved back to disk still has its H1 stripped by the drawer.

**Validation:** Reproduced by adversarial review of REQ-383 against the generated board in headless
Chromium; the three affected records were then enumerated across the whole tree.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, finding C2.*

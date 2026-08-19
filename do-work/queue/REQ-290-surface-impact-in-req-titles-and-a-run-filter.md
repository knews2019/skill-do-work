---
id: REQ-290
title: Surface impact in REQ titles and add a run filter that skips negligible work
status: pending
created_at: 2026-08-19T14:33:51Z
user_request: UR-060
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-289]
maintenance: false
related: [REQ-289]
batch: impact-effort-split
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
- skills/do-work/actions/review-work.md
---

# Surface Impact in REQ Titles and Add a Run Filter That Skips Negligible Work

## What

Make the `impact:` field REQ-289 introduces actually usable for the decision the user wants to make:
put the token in the REQ title so it is searchable today, and give `do-work run` a flag that skips
negligible-impact work.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A field nobody can filter on does not let the user stop work. Today `do-work run` has no size or
impact filter of any kind — `work.md:103-113` accepts only targeting tokens, `--fan-out`, and
`--wave`, and Step 1's selection scan reads status, dependency readiness, claim state, `assigned_to`,
and wave depth. Nothing else.

## Context

The board's search box already matches on `request.title` (`web/board-filters.js:20-31`), so a token
in the title is filterable with zero board work. A dedicated board filter control would need Go, JS,
and CSS changes; the title gets there for free.

There is no REQ title format rule anywhere in the repo today. `work-reference.md:119` and
`capture-reference.md:14` carry only the placeholder "Brief descriptive title". Prefix conventions
already exist in practice but are unenforced: `capture-reference.md:166` ("Addendum: ...") and
`review-work.md:362` ("Review fix: ...").

## Detailed Requirements

### The token goes in the title, not the filename

- The `impact:` token appears in the REQ's `title:` frontmatter value.
- It does **not** go into the filename slug. `REQ-NNN-slug.md` stays as it is. Impact can be revised
  at clarify, and a filename-borne token would mean renaming files mid-pipeline for a metadata edit.
  The board searches titles, not filenames, so nothing is lost.
- Give the convention a home in `capture-reference.md` alongside the templates, and follow the
  existing prefix precedents rather than inventing a third shape.
- `capture.md` and `review-work.md` emit titles conforming to it.
- `impact:` frontmatter stays the source of truth. The title is a mirror; when they disagree, the
  field wins.

### The run filter

- `do-work run --skip-impact-negligible` omits `impact: impact-negligible` REQs from Step 1's
  selection.
- One boolean flag, not a general `--impact <token>` selector. YAGNI until a second use appears.
- `work.md:107` rejects unrecognized arguments, so the usage string at `:110` changes in the same
  edit.
- The flag changes *which* REQs are selected. State how it composes with `--wave` and with explicit
  `REQ-NNN` targeting: explicit targeting names a REQ deliberately and must not be silently skipped.
- Report what was skipped. A run that silently drops REQs reads as "the queue is empty" when it is
  not.

## Constraints

- Blocked on REQ-289 — there is no field to filter on until it lands. `depends_on` enforces this.
- Write-set overlap with REQ-289 on `work.md`, `capture.md`, `capture-reference.md`, and
  `review-work.md` is expected and safe because the two are strictly sequential. Do not run them in
  parallel.
- No new board code. If a board filter control is wanted later it is a separate REQ.

## Dependencies

REQ-289 must complete first.

## Red-Green Proof

**RED prompt/case:** With a queue holding both negligible-impact and user-visible REQs, try to run
only the ones worth doing, and try to find the negligible ones on the board.

**Why RED now:** `do-work run` has no flag for it — `work.md:110`'s usage string offers targeting
tokens, `--fan-out`, and `--wave`, and unrecognized arguments are rejected. The board's search box
matches titles, and no title carries the impact token, so searching finds nothing.

**GREEN when:** `do-work run --skip-impact-negligible` runs the queue and omits exactly the
`impact-negligible` REQs, reporting how many it skipped. Typing `impact-negligible` into the board's
search box lists exactly those REQs, with no board code changed.

**Validation:** User confirmed — "Field + title prefix + run filter" chosen via the ask tool, with
the stated purpose "if I ever want to stop spawning and processing them I can".

## Full Context

See `do-work/user-requests/UR-060/input.md` for complete verbatim input.

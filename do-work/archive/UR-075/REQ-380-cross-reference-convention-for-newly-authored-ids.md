---
id: REQ-380
title: '[impact-rule-change] Cross-Reference Convention for newly authored REQ and UR ids'
status: completed
created_at: 2026-08-26T13:02:24Z
user_request: UR-075
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379]
batch: ticket-id-autocomplete
write_set:
- skills/do-work/actions/capture-reference.md
claimed_at: '2026-08-27T11:27:28Z'
status_changed_at: '2026-08-27T11:27:28Z'
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T11:27:28Z'
completed_at: '2026-08-27T11:33:05Z'
commit: 253b294
kb_status: skipped
---

# Cross-Reference Convention For Newly Authored REQ And UR Ids

## What

One canonical section in `skills/do-work/actions/capture-reference.md` saying that any flow writing
a REQ or UR id into prose names it — `REQ-1679 (Admin can delete a card)` on first mention, bare
afterwards — while frontmatter fields keep bare ids.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add one short condition-keyed convention beside the canonical title convention; no checker or changes to other actions.
- [x] **[APPLY]:** Added the first-mention title form, bare frontmatter exception, and prohibition on retrofitting existing references.
- [x] **[UNIFY]:** Reviewed capture-reference.md diff for all requirements and confirmed no other authoring action or checker changed. Canonical gate recorded below.

## Why

REQ-378 and REQ-379 patch the symptom at render and copy time. The cause is that the reference was
written as a bare number in the first place, which means it stays cryptic everywhere the board is
not: in a plain editor, in `git grep` output, in a pasted diff, in a review thread. The user's
framing is the general rule, not a board rule — "we shouldn't assume that people know what those
numbers mean."

## Context

`capture-reference.md` already carries two conventions in exactly this shape, both keyed on a
condition so their applicability lists cannot go stale:

- **REQ Title Convention** (line 9): "The condition is the rule — **any flow that mints a REQ
  carrying an `impact:` value follows this section** — so a new one inherits it without this list
  being re-counted."
- the **fold-before-mint contract** (line 24), with the same construction.

The convention this REQ writes down is also already in practice: that same file's own worked
example at line 221 reads
`- REQ-018 (durations-measured-face-constants-lack-provenance) — the second finding, that the face
numbers carry no provenance`. So this names an existing habit rather than inventing one, which is
why it is one section and not a campaign.

## Detailed Requirements

- **One new section in `capture-reference.md`**, written in that file's established
  condition-keyed form. Content:

  > **Cross-Reference Convention.** Any flow that writes a REQ or UR id into prose names it:
  > `REQ-1679 (Admin can delete a card)` on first mention in a document, bare afterwards. The
  > condition is the rule — capture, review findings, remediation notes, discovered tasks and
  > handoffs all inherit it without this list being re-counted. Frontmatter fields (`depends_on`,
  > `related`, `addendum_to`, `user_request`) stay bare ids: they are parsed, not read. Never
  > rewrite an existing file's references to add titles.

- **First mention only**, matching REQ-378's display rule — the two must not disagree about what
  the convention is.
- **Frontmatter stays bare.** State it explicitly; it is the part someone will get wrong.
- **No other action file is edited.** The flows that mint REQ prose already cite
  `capture-reference.md` as the canonical home for its conventions, so a second copy anywhere else
  is drift, not reinforcement.
- **Never retro-fit.** The convention governs newly authored references; existing REQ files are
  left exactly as they are.

## Constraints

- **Do not add a checker for this.** A prose convention enforced by a grep would fire on the
  section defining it — the failure `_dev/primes/prime-kanban-board.md` records from REQ-293, where
  "a prose-grep guard will fail on its own contract being stated". If a lock-in is wanted later it
  is its own REQ with that trap solved first.
- **One section, no list.** Do not enumerate the flows that inherit the rule; key it on the
  condition, per root `CLAUDE.md` § "State conditions, not lists".
- Keep it short. This is a convention statement, not a guide.

## Dependencies

None. Independent of REQ-378 and REQ-379 — it changes an action file, they change the board client,
and no file is shared.

## Builder Guidance

**Certainty level: Firm.** The wording above is the deliverable; adjust it for fit with the
surrounding sections' voice, not for scope.

Read `_dev/primes/prime-action-files.md` first — it owns the template, the earned-sections rule and
the cross-referencing conventions for anything under `skills/*/actions/`.

Resist growing this. The temptation is to also touch `review-work.md`, `work-reference.md` and the
handoff actions "for consistency"; that is exactly the drift the canonical-home pattern exists to
prevent.

## Red-Green Proof

**RED prompt/case:** Search `skills/do-work/actions/capture-reference.md` for a stated rule about
how a REQ or UR id is written into prose. There is none — the id-in-parens form appears only
incidentally, inside one worked example of the Folded Requests section, where a reader has no
reason to read it as a rule that applies anywhere else.

**Why RED now:** No section of any action file states the convention, so every flow that writes a
cross-reference decides the form for itself, and bare ids are the common outcome.

**GREEN when:** `capture-reference.md` carries a named Cross-Reference Convention section stating
the first-mention-with-title form, the frontmatter exception, and the never-retro-fit rule, keyed on
a condition rather than a list of flows — and no other action file gained a competing copy.

**Validation:** User confirmed — "Board + author-side convention" was selected from two presented
options during planning, with the alternative being board-only.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26 — "we shouldn't assume that people know what those numbers mean."*

## Triage

**Route: A** — One canonical authoring convention; no runtime change or new checker.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture-reference.md` (modified) — one canonical Cross-Reference Convention section.

**What was done:** Newly authored prose names the referenced ticket on first mention; frontmatter and existing references remain unchanged.

## Orientation

REQ and UR references in newly authored prose now carry a title on first mention, so they remain understandable outside the board.

## Qualification

Passed: mechanical qualify.sh exit 0; the source diff is one substantive section satisfying all five required behaviors. No data flow changed. P-A-U confirmed.

## Review

Acceptance: Pass. Overall: 100%. Independent Route A review checked the original UR, complete REQ, actual source diff, and a narrow restatement sweep. No findings or follow-ups.

## Testing

- `bash _dev/tests/maintainer-verify.sh`: exit 0 on the final source/release state. Contract suite, Go vet/tests, and strict JavaScript lane passed. The default gate explicitly skipped its optional browser lane; this REQ changes prose only.
- `git diff --check`: exit 0.
- Independent source review confirmed the new section and all constraints. No new checker was added, as requested.
- Optional Chrome 151 full-browser baseline was stopped after an unrelated drawer dump-dom test stalled; not counted as a passing browser run.

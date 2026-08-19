---
id: REQ-288
title: Fix the three unfiled contradictions in clarify's Step 4
status: pending
created_at: 2026-08-19T14:33:51Z
user_request: UR-059
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/clarify.md
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/review-work.md
- _dev/tests/contract-regressions.sh
---

# Fix the Three Unfiled Contradictions in clarify's Step 4

## What

`skills/do-work/actions/clarify.md` Step 4 carries three shipped contradictions. Each is a pair of
statements that cannot both be followed, and none is filed anywhere. Fix all three in one pass.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Clarify is the highest-volume writer of user answers in the suite. Two of these three defects
destroy work: K3 can archive a REQ that still holds an unanswered question, and K4 can archive an
approved follow-up `completed` without ever building it. `review-work.md:409` already documents K4's
failure mode and defends against it with nothing but a literal string comparison.

## Context

Found by a contradiction sweep run during a `do-work validate-feedback` session on REQ-276. Six
contradictions surfaced; REQ-276 and REQ-270 already cover two of them, and the sixth was deferred
by the user. These three had no owner.

Two candidates were checked and cleared during the same sweep, so do not "fix" them here:
`write_set` being both display-only and a builder write boundary is consistent (normative for the
builder at `work.md:392`, inert for the dispatcher, and both docs say so); and the discard path does
not contradict Principle 8, which carves out consent gates and still writes a `## Cancelled` block.

## Detailed Requirements

### K2 — the durable record is defined two ways

`crew-members/clear-questions.md:41` (Principle 8) defines the durable record as
`- [x] [question] -> [answer]` **plus a dated note carrying the reasoning, including anything the
answer put out of scope**. `clarify.md:96` — the block that names itself the canonical entry point —
says the answer line alone **is** the durable record, and Step 4's branches at `:104-110` instruct
nothing beyond the answer text. `work.md:248` sides with clear-questions.

- Conform `clarify.md` Step 4 to Principle 8: the answer write includes the dated reasoning note,
  and that note carries anything the answer put out of scope.
- The note is free-form. Do **not** add a templated `## Decision Record` section to
  `actions/capture-reference.md` — the shape stays flexible.
- `clarify.md:96`'s claim about what the durable record *is* must stop being narrower than the file
  it cites.

### K2's date constraint — cite the rule, never copy a command

The note carries a date, and ungoverned prose dates get fabricated. This is the failure UR-055 was
captured about; its fix was scoped to `*_at` instants only.

- `work-reference.md:101` states that the Timestamp rule is the only place in `actions/` that spells
  a clock command, and `_dev/tests/` fails the build if any action file reintroduces a copy. K2's
  note **cites** the rule. It never spells a command.
- The date-only paragraph at `work-reference.md:113` already carries the `YYYY-MM-DD` shape and its
  command, but scopes itself to an enumerated list of two consumers (the `do-work-knowledge`
  memory-logs mirror and the `do-work-toolbox` ui-review report header). REQ prose notes are not on
  that list. **Rewrite `:113` to key on the condition** — any UTC calendar date written into a
  durable record — instead of naming consumers. CLAUDE.md's "State conditions, not lists" and
  `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale both require this direction.
- The condition-keyed paragraph also governs the `**Answered [YYYY-MM-DD]:**` line clarify already
  writes into every REQ it touches, which today has no shipped format spec at all.

### K3 — per-question verbs set whole-file state

`clarify.md:102` introduces Step 4's branches with "For each question, the user can:", then `:107`
confirm sets `status: completed`, `:109` skip leaves the REQ `pending-answers`, and `:110` discard
sets `status: cancelled` and archives. `clarify.md:98` explicitly supports multi-question REQs
("A REQ with 6 questions needs 2 prompts"). Discard Q2 and skip Q3 and `:109`/`:110` cannot both be
obeyed.

- Per-question verbs record **only** the answer text and its note. They set no REQ-level status.
- Step 5 computes the REQ's status once, from the whole set of per-question outcomes, extending the
  aggregation it already performs at `:114` for the `pending` flip.
- Any remaining `- [ ]` keeps the REQ `pending-answers`. A REQ is never archived while a question is
  open.
- The Verification Checklist at the end of `clarify.md` gains the aggregate case.

### K4 — a confirm can archive work that still needs building

`clarify.md:107` routes "*All other REQs*" to `completed` and archives them on a confirm.
`review-work.md:409` names this exact failure and guards it only by requiring the literal string
`Should I process this as a new task?`.

- Narrow the "*All other REQs*" clause so it covers builder-decision follow-ups only; everything
  else falls through to the answer path.
- Add a frontmatter marker that the emitters stamp at REQ creation, and route the branch on the
  marker instead of on question prose. `review-work.md` stamps it where it creates follow-ups.
- The marker gets a Schema Read Contract row in `work-reference.md` like every other read field.
- After this lands, a reworded consent question can no longer misroute a follow-up.

## Constraints

- All three defects write `skills/do-work/actions/clarify.md`. They must be fixed by one builder in
  one pass — splitting them puts concurrent writers on one file
  (`actions/capture-reference.md:116`).
- `maintenance: true` — this is instruction maintenance on the skill's own action files. Load
  `crew-members/maintenance.md` and prefer narrowing and deleting over adding.
- Do not restate any Timestamp-rule command anywhere. Citations only.

## Dependencies

**REQ-262 (`pending`) interacts with K2's date fix and must not be worked in parallel with it.**
REQ-262 asks whether three prompt-kit templates "join the date-only paragraph's consumer list (UTC,
cited) or are declared template-content-out-of-scope". K2 replaces that consumer list with a
condition, which removes REQ-262's first option as stated. The two REQs write different files
(`skills/do-work-knowledge/prompts/**` vs `skills/do-work/actions/work-reference.md`), so there is no
file collision — but whoever works REQ-262 after this lands must re-read the rewritten paragraph and
answer against the condition rather than against a list. Flag this to the user at hand-back rather
than editing REQ-262 from here.

## Red-Green Proof

**RED prompt/case:**
1. A REQ in `do-work/queue/` holding two open questions. Discard the first and skip the second in one
   `do-work clarify` run.
2. A `review-work`-created follow-up REQ whose consent question is reworded — same meaning, different
   words from `Should I process this as a new task?`. Confirm the recommended choice.
3. A `do-work clarify` answer that changes what gets built and puts something out of scope.

**Why RED now:**
1. `clarify.md:109` and `:110` give contradictory instructions; the REQ is archived `cancelled` while
   still carrying an unanswered `- [ ]`, or left in the queue already marked terminal.
2. `clarify.md:107`'s "*All other REQs*" branch archives it `completed` and it is never built —
   the failure `review-work.md:409` predicts.
3. Only the answer text is written. The reasoning and the out-of-scope item exist nowhere, and the
   next reader re-derives the decision from `Recommended:`.

**GREEN when:**
1. The REQ stays `pending-answers` in `do-work/queue/`, holding the discarded answer and the still-open
   question. Nothing is archived.
2. The follow-up lands `pending` and the work loop picks it up. Routing no longer depends on the
   question's wording.
3. The REQ carries the answer line **and** a dated note with the reasoning and the out-of-scope item,
   with the date obtained per the Timestamp rule's condition-keyed date-only paragraph.

Each case gets a lock-in check in `_dev/tests/contract-regressions.sh` naming the defect it pins.
`bash _dev/tests/maintainer-verify.sh` exits 0.

**Validation:** User confirmed (K3 and K4 chosen via the ask tool; K2's date constraint supplied by
the user in the same exchange).

## Full Context

See `do-work/user-requests/UR-059/input.md` for complete verbatim input.

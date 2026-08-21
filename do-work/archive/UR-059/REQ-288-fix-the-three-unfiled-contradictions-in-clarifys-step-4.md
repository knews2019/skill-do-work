---
id: REQ-288
title: Fix the three unfiled contradictions in clarify's Step 4
status: completed
created_at: 2026-08-19T14:33:51Z
claimed_at: 2026-08-21T00:34:42Z
completed_at: 2026-08-21T00:57:58Z
kb_status: pending
commit: c25ee71
route: C
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
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---

## Triage

**Route: C** - Complex

**Reasoning:** Three separate contradictions in one file, one of which needs a frontmatter marker with a Schema Read Contract row, plus a paragraph rewrite in a second file, a defense rewrite in a third, and lock-in checks in a fourth. Two of the three destroy work when they fire, so the ordering and the interactions had to be worked out before editing.

**Planning:** Required

## Plan

### The shape of the fix, defect by defect

**K2 — durable record.** Widen `clarify.md`'s named-entry-point block to Principle 8's definition (answer line **plus** dated reasoning note including what the answer put out of scope), and make Step 4's branches instruct it. The date **cites** the Timestamp rule; no clock command enters this file.

**K2's date constraint.** Rewrite `work-reference.md`'s date-only paragraph to key on the condition — *any UTC calendar date written into a durable record* — instead of enumerating two consumers. REQ prose notes are then governed by the condition rather than needing to be added to a list.

**K3 — per-question verbs.** Strip every REQ-level status write out of Step 4's branches. Rename Step 5 to say what it now does and have it compute status once from the whole outcome set, with "any remaining `- [ ]` wins" first so a REQ can never be archived holding an open question.

**K4 — marker routing.** **The marker already exists.** `builder_decided: true` is stamped by the Builder-Decided Follow-up Template (`work-reference.md` → Step 8) and by nothing else; `review_generated: true` follow-ups never carry it. So K4 needs no new field — it needs the `completed` branch keyed on the existing marker, plus the Schema Read Contract row the REQ correctly says any read field must have. `crew-members/maintenance.md` § 1 is explicit that subtraction and reuse come before addition, and this is the cheaper fix by a wide margin: a new field would need a second emitter, a second contract row, and a migration story for every REQ already carrying `builder_decided`.

Then delete `review-work.md`'s literal-string MUST, which exists only to defend the prose-keyed routing K4 removes. That is a subtraction, not a rewrite.

### Ordering

K2's date fix must land before or with K2's note requirement, since the note's date depends on the rewritten paragraph. K3 must land before K4, because K4's branch lives in the Step 5 that K3 rewrites. Everything is one file-set and one pass, as the REQ's Constraints require.

### Plan validation

- **Requirement coverage:** K2 → clarify entry point + Step 4 + `work-reference.md` paragraph; K3 → Step 4 branches + Step 5 + checklist; K4 → Step 5 branch + contract row + `review-work.md` deletion; all three → `contract-regressions.sh`.
- **No orphan tasks:** the `review-work.md` edit is the only task not literally in the requirement list, and the REQ names its line as the thing defending the defect.
- **Scope sanity:** 4 tasks. Under the 5-task flag.

*Planned inline by the orchestrator*

## Exploration

- `clarify.md:96` is the named-entry-point block; `:102-110` the branches; `:114` the old Step 5.
- **`builder_decided: true` already exists and is already the right discriminator.** Emitted at `work-reference.md:655` (Builder-Decided Follow-up Template, Step 8). Read today by `verify-requests.md:201`, `forensics.md:50`, and `clarify.md:100/104/108`. **It has no Schema Read Contract row** — the one thing K4 genuinely needed to add.
- `review-work.md:364-368` shows review follow-ups carry `review_generated: true` and **not** `builder_decided`, which is what makes the marker a correct discriminator rather than a coincidence.
- `review-work.md:413` is the literal-string MUST that K4 obsoletes.
- `contract-regressions.sh:59` `assert_block_contains` / `:70` `assert_block_not_contains`, with blocks extracted by `sed -n '/start/,/end/p'`. Existing clarify blocks are extracted the same way at `:3875`.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/actions/clarify.md` (modify) — K2 entry point + Step 4 note requirement; K3 branch destatusing + Step 5 aggregation; K4 marker branch; discovered-task lead de-prosed; Verification Checklist
- `skills/do-work/actions/work-reference.md` (modify) — date-only paragraph condition-keyed; `builder_decided` Schema Read Contract row
- `skills/do-work/actions/review-work.md` (modify) — literal-string MUST replaced by the marker rule
- `_dev/tests/contract-regressions.sh` (modify) — eleven lock-in checks across the three defects

**Files I will NOT touch:**
- `actions/capture-reference.md` — the REQ explicitly forbids a templated `## Decision Record`; the note stays free-form.
- `do-work/queue/REQ-262-*` — the REQ says flag the interaction at hand-back rather than editing it. It had already been completed earlier in this run; see D-04.

**Acceptance criteria (restated from the REQ):**
1. Discard one question and skip another: the REQ stays `pending-answers` in the queue, nothing archived.
2. A reworded consent question on a review follow-up still lands `pending` and gets built.
3. An answer that changes what gets built records the answer **and** a dated note with reasoning and the out-of-scope item.
4. Each case gets a lock-in check naming the defect it pins.
5. No Timestamp-rule command is restated anywhere; citations only.
6. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Decisions

- **D-01** (ESCALATE): K4 reuses the existing `builder_decided: true` marker instead of adding a new frontmatter field, which is what the REQ's text asked for. Reasoning: the field already exists, is already stamped by exactly the right emitter and by no other, and is already read by three actions — it is the discriminator K4 describes, already in the tree. `crew-members/maintenance.md` § 1 and § 2 both say to look for the existing cause before adding, and a new marker would need a second emitter, a second contract row, and a story for every REQ already carrying the old one. What the REQ genuinely identified as missing — a Schema Read Contract row — was added. Value: the fix is a routing change plus one table row rather than a new schema field with a migration. Risk: `builder_decided` now carries routing weight it did not before, so a REQ that wrongly carries it would be archived `completed`. Mitigated because only one template emits it, and the contract row now names every read site.
- **D-02** (DECIDE & STATE): Deleted `review-work.md`'s "MUST contain the exact discriminator phrase" requirement rather than keeping it as belt-and-braces. Reasoning: it exists only to defend prose-keyed routing. Kept as a live MUST it would be a rule with no consequence — the failure it names can no longer happen — and a rule nobody's behavior depends on is exactly the bloat `maintenance.md` § 1 says to remove. The house shape and the phrase survive as guidance, not as a gate.
- **D-03** (DECIDE & STATE): Step 5 was renamed from "Activate answered REQs" to "Resolve each REQ's status from its whole question set". Reasoning: it no longer only activates — it can now also produce `cancelled` and `completed`, the two dispositions K3 moved out of Step 4. Grepped for the old name across `skills/` and `_dev/`: zero references, so nothing broke.
- **D-04** (ESCALATE — flagged for the user): The REQ's Dependencies section says whoever works REQ-262 after K2 lands must answer against the rewritten condition rather than the old consumer list, and asks that this be flagged at hand-back rather than edited from here. **REQ-262 was completed earlier in this same run**, before REQ-288 was claimed, so the order is the reverse of what the REQ anticipated. Checked directly: REQ-262 chose the *out-of-scope* option, wrote it as a condition-keyed carve-out, and did not add anything to the consumer list — so it is compatible with K2's rewrite, and the carve-out survives it verbatim (verified by grep after the edit). No action needed, and nothing about REQ-262 was edited. Flagged here because the REQ asked for it to be surfaced.

## Implementation Summary

**What was done:** Fixed all three contradictions in one pass. The durable record now includes its dated reasoning note and cites the Timestamp rule for the date; that rule's date-only paragraph was rewritten to key on the condition instead of two named consumers. Per-question verbs no longer set REQ-level status — Step 5 computes it once from every question's outcome, with any remaining open question holding the REQ in the queue. The `completed` fast path routes on the existing `builder_decided: true` marker, which gained the Schema Read Contract row it never had, and the literal-string defense that existed only to protect the old prose-keyed routing was deleted.

**Files changed:**
- `skills/do-work/actions/clarify.md` (modified) — named-entry-point block redefined as answer + dated note; Step 4 branches record outcomes only and instruct the note; Step 5 renamed and rewritten to aggregate; discovered-task lead de-prosed; six Verification Checklist lines replacing two.
- `skills/do-work/actions/work-reference.md` (modified) — date-only paragraph condition-keyed; `builder_decided` Schema Read Contract row added.
- `skills/do-work/actions/review-work.md` (modified) — the exact-phrase MUST replaced by the marker rule.
- `_dev/tests/contract-regressions.sh` (modified) — eleven lock-in checks, each naming its defect.

**Tests touched:** eleven new assertions in `contract-regressions.sh`. No existing assertion changed meaning.

## Qualification

Passed — 4 files verified, 6 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `git diff --stat` reviewed per file. `shellcheck` clean on the edited test (it is in `maintainer-verify`'s lint set). No debug artifacts. `maintainer-verify.sh` exits 0.
- **Substantive:** the clarify diff removes two whole-file status writes and adds a five-branch aggregation step; not whitespace.
- **Requirements traced:** AC1 → K3's "any remaining `- [ ]` wins" plus its checklist line and lock-in check; AC2 → K4's marker branch plus `review-work.md`'s rewritten bullet; AC3 → K2's entry-point redefinition and Step 4 note instruction; AC4 → eleven checks, each with a defect-naming message; AC5 → `grep -cE "date -u|Get-Date" clarify.md` returns 0, and a lock-in check pins that; AC6 → verify exits 0.
- **Flowing:** instruction prose — no data path. The behavior it governs is an agent following it, which is why the lock-in checks assert the instructions rather than a runtime.

## Testing

- `bash _dev/tests/contract-regressions.sh` — passed, including all eleven new checks.
- `bash _dev/tests/maintainer-verify.sh` — exit 0.

**Red-green validation** — the REQ's captured RED is three scenarios in an action file, so each is pinned by the lock-in check that would fail if the instruction regressed, and each check was **mutation-tested** by reintroducing the defect:

| Captured RED | Mutation applied | Result |
|---|---|---|
| K2: only the answer text is written; reasoning and out-of-scope item exist nowhere | "together with a dated note" → "WITHOUT a dated note" | **FAIL**, naming K2's durable-record defect |
| K2 date: an action file spells a clock command | "cite it, never spell a clock command here" → "obtain it with `date -u +%F`" | **FAIL** ×2 — both the cite-the-rule check and the no-command check |
| K2 date: the paragraph enumerates consumers again | "any UTC calendar date written into a durable record" → "the two sites listed below" | **FAIL**, naming the stale-enumeration defect |
| K3: a per-question verb sets whole-file state | "never sets the REQ's status and never archives" → "sets the status directly" | **FAIL**, naming K3 |
| K4: routing keyed on question prose | marker branch → `"Should I process this as a new task?"` branch, and "Never infer this branch from question prose" → the opposite | **FAIL** ×2 on both K4 checks (verified in isolation) |

Every mutation was reverted and the suite re-run clean. No check passed vacuously.

**AC5 checked directly, not only by the lock-in check:** `grep -cE "date -u|Get-Date" skills/do-work/actions/clarify.md` → `0`.

## Review

**Overall: 94%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| K2: Step 4 conforms to Principle 8 — answer plus dated reasoning note with out-of-scope items | ✅ |
| K2: no templated `## Decision Record` added; note stays free-form | ✅ nothing added to `capture-reference.md` |
| K2: `clarify.md:96`'s claim stops being narrower than the file it cites | ✅ redefined to both halves |
| K2 date: cite the rule, never spell a command | ✅ grep returns 0, and a check pins it |
| K2 date: `work-reference.md` paragraph keyed on the condition | ✅ |
| K3: per-question verbs set no REQ-level status | ✅ |
| K3: Step 5 computes status once from the whole set | ✅ |
| K3: any remaining `- [ ]` keeps the REQ `pending-answers`; never archived | ✅ stated first in Step 5 |
| K3: Verification Checklist gains the aggregate case | ✅ six lines replacing two |
| K4: "*All other REQs*" narrowed to builder-decision follow-ups only | ✅ |
| K4: routing on a frontmatter marker, not question prose | ✅ on the existing `builder_decided` — D-01 |
| K4: the marker gets a Schema Read Contract row | ✅ |
| Each case gets a lock-in check naming its defect | ✅ eleven, all mutation-tested |

### Findings

**Important — none.**

**Minor:**

- **M1:** D-01 diverges from the REQ's literal instruction to *add* a marker, reusing one instead. Anyone reviewing against the REQ's words rather than its purpose would flag it; the argument is that the marker it describes already exists with the right emitter and the right absence on review follow-ups, and `maintenance.md` puts reuse ahead of addition. The genuinely missing half — the contract row — was added.
- **M2:** Step 5 now carries five outcomes where it carried one. It is longer than the step it replaced, in a file whose whole defect class was over-loaded instructions. The length buys a single place where status is decided, which is exactly what K3 asked for, but this step is now the one most worth watching for accretion.

**Nit:**

- **N1:** The eleven lock-in checks assert on prose fragments, so an innocuous rewording of a correct instruction can fail them. That is the standing trade in `contract-regressions.sh` and matches every existing check in it; the messages name the defect rather than the wording, so a failure tells the next editor what the rule is for.

### Restatement Sweep

Redefined elements: what "the durable record" is; what a per-question verb may write; what routes a confirm to `completed`; what the date-only paragraph governs.

- `crew-members/clear-questions.md` Principle 8 — the source K2 conforms *to*; unchanged and now agreed with rather than contradicted. Checked its text still describes both halves: it does.
- `actions/work.md` Step 3.5 — the other caller of the Canonical answered-question format. Re-read: it already required the reasoning ("Append the reasoning too, including anything the answer put out of scope"), so it was on Principle 8's side all along and needs no edit. That is the asymmetry K2 was filing.
- `actions/verify-requests.md:201`, `actions/forensics.md:50` — the other `builder_decided` readers. Both read it as a property of the REQ, neither as a routing key; unaffected, and both now covered by the new contract row.
- `actions/review-work.md:413` — the literal-string defense; **deleted and replaced** (D-02).
- `actions/work-reference.md:712` — the Discovered Tasks consent-question template still contains the phrase `Should I process this as a new task?`. **Left as is deliberately:** it is now house shape rather than a routing key, exactly as `review-work.md`'s rewritten bullet says, and the phrase is good question-writing.
- Step 5's old name "Activate answered REQs" — grepped across `skills/` and `_dev/`: zero references (D-03).
- `actions/clarify.md`'s own Verification Checklist — the in-file restatement of every rule above; rewritten in the same commit.

No stale restatement remains.

### Acceptance Testing

The deliverable is instruction text, so acceptance is: would an agent following the new text produce the captured GREEN? Each of the three scenarios was traced through the rewritten steps by reading them as a builder would, and each is additionally pinned by a lock-in check that was proved to fail when the defect is reintroduced. The mutation table in `## Testing` is the evidence: eleven checks, six mutations, every one caught.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope Discipline | 90% |
| Risk | Low |
| Acceptance | Pass |

Scope Discipline 90% for D-01's divergence from the literal instruction and D-02's deletion, both argued. Risk Low rather than None: K3 and K4 change what clarify does to real REQ files, and clarify is the suite's highest-volume writer of user answers — but both changes are strictly more conservative than what they replace (fewer archives, fewer terminal flips, no prose-dependent routing).

### Follow-up REQs Created

None. D-04's interaction with REQ-262 was checked and needs no work.

## Lessons Learned

**What worked:** Checking whether the marker K4 asked for already existed before adding one. It did — `builder_decided: true`, stamped by exactly the right emitter and absent from exactly the REQs that must not take that branch — so the fix collapsed from "new field, second emitter, contract row, migration" to "route on this, add the missing row". `maintenance.md`'s deletion questions are what prompted the check; the REQ itself said "add a marker" and would have been followed literally otherwise.

**What didn't:** The mutation runs kept hitting the command timeout, because `contract-regressions.sh` runs the whole prescribed-shell suite before its own assertions and takes minutes. Two mutations were lost to timeouts before switching to extracting the block and running the two `grep` patterns directly against an original and a mutated copy — same assertion, seconds instead of minutes. For a prose-assertion check, testing the pattern against two copies of the file is equivalent to running the suite and is the right loop to iterate in.

**Worth knowing:** K4's failure mode is worth remembering as a shape, not just as a bug. `review-work.md` had *documented* it — it knew the routing was prose-keyed, knew a rewording would misroute an approved follow-up into being archived unbuilt, and defended it by requiring an exact sentence in every consent question it wrote. A rule that requires every future author to reproduce a magic string is a defense with a per-use failure probability. The marker moves the decision to something the emitter sets once and nobody retypes.

## Orientation

Clarify no longer contradicts itself in the step that writes user answers. A per-question verb records only that question's outcome, so discarding one question while skipping another is now expressible and a REQ holding an open question can no longer be archived; the REQ's status is decided once, at the end, from every answer together. The confirm-to-`completed` fast path routes on the `builder_decided` marker rather than on a question's wording, so rephrasing a consent question can no longer archive a follow-up that still needed building. Answers now carry their reasoning and what the answer put out of scope, with the date governed by a rule that keys on the condition instead of a list of two consumers. Lives in the core queue actions (`skills/do-work/actions/clarify.md`, with contract text in `work-reference.md` and `review-work.md`), indexed by `_dev/primes/prime-action-files.md`.

**[MAP CHANGED]** — `builder_decided: true` is now a routing key with a Schema Read Contract row, not just a property. The Timestamp rule's date-only paragraph governs by condition rather than by an enumerated consumer list, so every future dated record is in scope without being added anywhere. Clarify's Step 5 was renamed and is now where every terminal disposition is decided.

Prime staleness spot-check: `_dev/primes/prime-action-files.md` — referenced paths still resolve; its Cross-Referencing and Lessons sections are unaffected by this change.

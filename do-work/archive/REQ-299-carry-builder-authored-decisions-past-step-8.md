---
id: REQ-299
title: "[impact-rule-change] Review fix: carry builder-authored sections past Step 8, starting with ## Decisions"
status: completed
created_at: 2026-08-19T20:03:19Z
status_changed_at: 2026-08-20T08:22:51Z
claimed_at: 2026-08-21T02:52:01Z
completed_at: 2026-08-21T03:20:43Z
kb_status: pending
commit:
route: C
user_request: UR-055
addendum_to: REQ-270
domain: general
review_generated: true
sweep: true
sweep_key: builder-authored-sections-unread-outside-step-8
impact: impact-rule-change
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-21T02:54:38Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/review-work.md
- _dev/tests/contract-regressions.sh
---

# Review Fix: Carry Builder-Authored Sections Past Step 8, Starting with `## Decisions`

## What

REQ-270 closed the case where a worktree builder's `## Discovered Tasks` never reached
Step 8, and keyed its new rule on the condition — but scoped that condition's home to
**Step 8's substeps** (`actions/work.md` → *Where a builder-authored section is read from*,
which opens "Some substeps **below**"). Its independent review then found a second
builder-authored section with the identical defect, read from **outside** Step 8, where the
new rule structurally cannot reach:

`## Decisions` is written by the builder (`actions/work.md` Step 6, "Log Decisions as
D-XX" — unqualified, so a worktree builder is told to write a file it may not write), and
read at two sites: `actions/review-work.md` Step 4's traceability check ("If a
`## Decisions` section exists in the REQ: verify that significant implementation
choices … are documented") and the end-of-run Decision Brief's HANDLED block
(`actions/work-reference.md` → **Decision Brief (hand-back format)**).

**The failure is silent in both directions.** Review's traceability check finds no section
and reports clean rather than flagging an absence it cannot distinguish from "the builder
made no decisions"; the Decision Brief renders an empty HANDLED list, so a DECIDE & STATE
choice the builder actually made never reaches the user. Under fan-out that is every
builder's decisions, every run.

**Done means the class cannot recur:** the rule lives somewhere every reader of a
builder-authored section inherits it, not in Step 8's preamble, and the sections Step 6
tells the builder to author are the same set the hand-back contract names.

## Instances

- [x] `## Decisions` — builder-authored at `actions/work.md` Step 6; read by
  `actions/review-work.md` Step 4 and by the Decision Brief's HANDLED block. The instance
  that fired this REQ.
- [x] `actions/work.md` Step 6's `## Decisions` instruction itself — it must say where the
  section goes when the builder may not write the REQ, exactly as REQ-270 fixed the
  `## Discovered Tasks` bullet beside it.
- [x] The rule's home — REQ-270 put it under Step 8 because that was the only reader then
  known. Move or restate it so a reader outside Step 8 inherits it, and mark the reader
  list illustrative rather than closed.

## Context

Found during the independent review of REQ-270 (finding 2, Important, gate:
`impact-rule-change` — it changes where a rule lives and which readers inherit it, across
several sites). REQ-270's other Important finding was a stale restatement in
`crew-members/general.md` that actively defeated its own fix; that one was repaired inside
REQ-270 rather than deferred, because the fix did not work end to end without it. This one
is genuinely a different scope: it needs a reader outside Step 8 and a check that did not
exist.

Created `pending-answers` per the generation-≥2 cascade depth stop: REQ-270 carries
`review_generated: true`.

## Requirements

- Every reader of a builder-authored REQ section — at Step 8 or anywhere else — reads the
  hand-back when the REQ file does not carry the section, and the rule states that as its
  condition rather than listing the readers it knows about today.
- `actions/work.md` Step 6's `## Decisions` instruction names the hand-back for a builder
  that may not write the main tree.
- The Decision Brief and review-work's traceability check distinguish "no section anywhere"
  from "the builder recorded nothing", and say which.
- **Every change holds in a consumer install**, where the suite is vendored under
  `.claude/skills/` and only `do-work/` is at the project root — see **Consumer-Install
  Constraint** below. Proven against a consumer-shaped fixture, not by reading.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** The property check REQ-270's reviewer specified, which pins the
property rather than any wording: extract every `## <Name>` section that `actions/work.md`
Step 6 instructs the **builder** to author, and assert each one is named in
`actions/work-reference.md`'s per-builder-output hand-back contents. It fails today, on
`## Decisions`.
**Why RED now:** `work-reference.md`'s hand-back row names `## Discovered Tasks` only —
REQ-270 added it — while Step 6 tells the builder to author `## Decisions` too.
**GREEN when:** The same check passes, and the full suite still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure
Ratchet (Step 6.5)**.

## Open Questions

- [x] REQ-270 made Step 8 read a worktree builder's `## Discovered Tasks` from its hand-back
  when the REQ file cannot carry it. Its review found the same gap for `## Decisions` — the
  numbered record of the judgment calls a builder made — which is read by the code review and
  by the end-of-run report you see, both outside Step 8, so under parallel building those
  decisions silently never reach you. Closing it means moving where that rule lives so
  readers outside Step 8 inherit it, and adding a check that every section the builder is
  told to write is one the hand-back is told to carry. Should I process this as a new task?
  → **Yes, add to queue — the full scope**, with one added constraint: it must work from the
  perspective of the installed skill in another repo, not only in this maintainer checkout.
  See **Consumer-Install Constraint** below, which the answer added to the Requirements.

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the technical fix is clear, but it decides how much of
  the pipeline's instruction set reorganizes around parallel building — REQ-270 deliberately
  drew the boundary at Step 8, and widening it touches the review action and the hand-back
  format you read at the end of every run. It is also only worth paying for if you intend to
  keep using worktree fan-out; in a purely serial run nothing here ever fires.

**Answered [2026-08-20]:** User approved via `do-work clarify` at full scope, after asking
how the fix would work and being shown the two reader sites verbatim. The user added the
consumer-install constraint in the same answer; it is a Requirement, not a nice-to-have.

## Consumer-Install Constraint (added by the user at clarify, 2026-08-20)

**Everything this REQ changes must hold in a consumer install, where the suite is vendored
under `.claude/skills/` and only `do-work/` sits at the project root — not just in this
maintainer checkout, where the two coincide.** The user asked for this explicitly. It is the
same class of defect REQ-282 shipped a fix for the day before: a path that resolves in the
maintainer layout and silently resolves to nothing in a consumer's.

Concretely, before this REQ is done:

- **The hand-back path must resolve for an installed skill.** The rule points a reader at
  `do-work/runs/work-<stamp>/REQ-NNN-handback.md`. That is relative to the **project root**,
  where `do-work/` lives in both layouts, so it should hold — **verify it rather than assume
  it**, and state in the rule which root the path is relative to, since a builder resolving
  it against the vendored skill directory finds nothing.
- **No shipped instruction may cite a maintainer-only path.** The property check lives in
  `_dev/tests/`, which is export-ignored and never installed; the *rule* it enforces ships.
  Nothing in `work.md`, `work-reference.md`, or `review-work.md` may reference `_dev/`, the
  check, or this repo's layout. `_dev/tests/shipped-package-reference-contract.sh` already
  enforces that class — make sure it still passes and that it actually covers the new text.
- **Cross-package citations stay resolvable from the citing file's own directory.** The
  three touched action files are all in the core package, so no `../` hop should be needed;
  if one appears, it must be correct at the installed depth, not the repo depth.
- **Verification is an execution, not a reading.** Build a consumer-shaped fixture — suite
  vendored under `.claude/skills/`, `do-work/` at the root — and confirm the instructions a
  builder and an orchestrator would follow resolve there. REQ-282's review used exactly this
  fixture shape and it is cheap to rebuild.

## Plan

**Approach.** The rule has to leave Step 8. `actions/work-reference.md` is the one file
every reader already loads and cites — `work.md` Step 8, `review-work.md` Step 4, and the
Decision Brief all live there or point there — so the rule moves into it as a top-level
section, `## Reading a Builder-Authored Section (any step)`, keyed on the condition
("whenever you read a `##` section the builder authors") with the reader list marked
illustrative. Step 8's two paragraphs collapse to a pointer.

**The check must not carry a list.** The property is "every `##` section Step 6 tells the
builder to author is named in the hand-back contents". Extracting "tells the builder to
author" mechanically needs a signal in the text. Rather than hand-maintain an exclusion
list of the sections Step 6 merely *reads* (`## Scope`, `## Red-Green Proof`,
`## Implementation Summary`), the check classifies **every** `##` token mentioned in Step
6's builder-instruction block: each is either routed to the hand-back (and then must appear
in `work-reference.md`'s per-builder-output row, both directions) or carries an explicit
`not yours to write` disclaimer in every bullet that mentions it. A section mentioned in
Step 6 that classifies neither way fails. A new section added later cannot pass silently.

**Steps.**

1. RED first: add the property check to `_dev/tests/contract-regressions.sh` and confirm it
   fails on `## Decisions` against the untouched tree.
2. `work-reference.md`: new `## Reading a Builder-Authored Section (any step)` carrying the
   moved rule, the absence-vs-silence distinction, the illustrative reader list, and the
   root the hand-back path resolves against (project root, where `do-work/` sits in both a
   maintainer checkout and a consumer install).
3. `work.md` Step 8: replace the two paragraphs with a pointer to that section.
4. `work.md` Step 6: add the worktree routing clause to the `## Decisions` bullet; add the
   `not yours to write` disclaimer to the three read-only section mentions.
5. `work-reference.md` per-builder-output row: name both routed sections, drop "today",
   cite the new section instead of Step 8.
6. `review-work.md` Step 4 and the Decision Brief HANDLED bullet: read per the new section
   and distinguish "no section anywhere" from "the builder recorded nothing".
7. Consumer-install proof: install the suite into a throwaway project with
   `tools/install-do-work-suite.sh`, then resolve every citation the three touched files
   make from their installed locations, and confirm the hand-back path resolves against the
   project root there.
8. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modify) — the property check and the rule-home pins
- `skills/do-work/actions/work-reference.md` (modify) — the rule's new home, hand-back row, Decision Brief
- `skills/do-work/actions/work.md` (modify) — Step 8 pointer, Step 6 section classification
- `skills/do-work/actions/review-work.md` (modify) — Step 4 traceability bullet

**Files I will NOT touch:** `_dev/tests/shipped-package-reference-contract.sh` — its coverage boundary is a separate defect and goes to Discovered Tasks; `skills/do-work/crew-members/general.md`, already repaired inside REQ-270.

**Acceptance criteria (restated from REQ):**
- [x] Every reader of a builder-authored section reads the hand-back when the REQ lacks it, stated as a condition rather than a reader list
- [x] Step 6's `## Decisions` instruction names the hand-back for a builder that may not write the main tree
- [x] The Decision Brief and review-work's traceability check distinguish "no section anywhere" from "the builder recorded nothing"
- [x] Every change holds in a consumer install, proven against a consumer-shaped fixture
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]** — rule moves to `actions/work-reference.md`; the check classifies every `##` section Step 6 names rather than carrying a list
- [x] **[APPLY]** — four files, all inside the declared scope
- [x] **[UNIFY]** — `git diff --stat` reviewed (4 files, +173/-10); `bash -n` on the changed shell; no debug artifacts; each changed file re-read at its edit sites

## Implementation Summary

The rule left Step 8. `actions/work-reference.md` gained
`## Reading a Builder-Authored Section (any step)` — the moved rule, keyed on the condition
("whenever you read a `##` section the builder authors"), with the reader list named as
illustrative, the hand-back path stated as project-root-relative, and the
absence-vs-silence paragraph. `actions/work.md` Step 8's two paragraphs collapsed to one
pointer at it.

`actions/work.md` Step 6 now classifies every `##` section it names. The `## Decisions`
bullet carries the worktree routing clause REQ-270 gave `## Discovered Tasks`; the three
sections the builder only reads (`## Red-Green Proof`, `## Scope`, `## Implementation
Summary`) each say they are **not yours to write**. `actions/work-reference.md`'s
per-builder-output row names both routed sections and states that the two sets are one.

Both readers outside Step 8 inherit the rule: `actions/review-work.md` Step 4's
traceability bullet reads `## Decisions` per the shared section and reports which absence it
found; the Decision Brief's HANDLED bullet does the same and renders
`• REQ-NNN: decisions not recovered — hand-back unread` rather than omitting the block.

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified) — new rule section; hand-back contents row; Decision Brief HANDLED bullet
- `skills/do-work/actions/work.md` (modified) — Step 8 pointer; Step 6 section classification across four bullets
- `skills/do-work/actions/review-work.md` (modified) — Step 4 traceability bullet
- `_dev/tests/contract-regressions.sh` (modified) — the property check plus eight pins on the rule's home and its two outside readers

## Decisions

- **D-01 — The property check classifies every `##` section Step 6 names, per section
  rather than per mention.** DECIDE & STATE. Extracting "instructs the builder to author"
  mechanically needed a signal in the text. The alternative was a hand-maintained exclusion
  list of the sections Step 6 merely reads, which goes stale the moment Step 6 grows
  (`CLAUDE.md` → Closed Enumerations Go Stale). Instead each section is routed to the
  hand-back or marked `not yours to write`, and a section classified by neither fails.
  Per-section rather than per-mention because several bullets name the same section and one
  clear statement of who writes it is enough for a reader; the cost is that a future
  authoring bullet for a section disclaimed elsewhere would not fail here, which is a
  visible self-contradiction rather than the silent loss this REQ exists to close.
- **D-02 — The rule's new home is `actions/work-reference.md`, not a crew-member file.**
  DECIDE & STATE. Its readers are the orchestrator and the reviewer, not the builder, and
  all three sites already load and cite `work-reference.md`. A crew-member file loads
  during implementation, which is the one moment the rule does not apply.
- **D-03 — The Decision Brief renders an explicit not-recovered line instead of omitting
  the HANDLED block.** DECIDE & STATE. Omission is what made the original failure silent;
  an empty block and an unread hand-back had the same rendering. Requirement 3 asked the
  two to be distinguishable, and the only place a user sees it is the brief.

## Discovered Tasks

- **impact-rule-change** `_dev/tests/shipped-package-reference-contract.sh` resolves
  **cross-package** citations only. A dangling same-package citation ships silently: a
  planted `actions/no-such-file.md`, `../docs/no-such-file.md`, and
  `crew-members/no-such-file.md` in `actions/work-reference.md` each left the contract at
  PASS, while `../../do-work-board/actions/no-such-file.md` failed it. Found while
  confirming the contract covers this REQ's new text, as its Consumer-Install Constraint
  required. A sweep of the installed suite found no live dangling same-package citation
  today — the seven candidates are all project-owned paths or the deliberately hostile
  `prompts/init.md` example — so this is coverage, not breakage. A fix needs the
  same-package resolver arm plus a policy for those two false-positive classes.

## Testing

**Tests run:** `QUEUE_KANBAN_BROWSER=<chromium> bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 — `Maintainer verification passed.`

**Red-green validation:**
- `_dev/tests/contract-regressions.sh` REQ-299 property check: ✗ before the fix
  (`routed by Step 6 but not carried by the hand-back: ['## Decisions']`) → ✓ after. The
  captured RED was reproduced on the untouched tree first: Step 6 named
  `## Decisions`, `## Discovered Tasks`, `## Implementation Summary`, `## Red-Green Proof`,
  `## Scope`, while the hand-back row named `## Discovered Tasks` alone.

**Mutation-tested — every new check was reverted individually and confirmed to fail:**
- M1 strip the routing clause from Step 6's `## Decisions` bullet → set mismatch, caught
- M2 drop `## Decisions` from the hand-back contents row → set mismatch, caught
- M3 remove the `not yours to write` mark from `## Scope` → unclassified section, caught
- M4 delete the `## Reading a Builder-Authored Section (any step)` section → caught
- M5 close the reader list ("These are the readers") → caught
- M6 replace "relative to the project root" → caught
- M7 drop the absence-vs-silence paragraph → caught
- M8 have Step 8 restate the rule instead of pointing at it → caught
- M9 remove review-work's citation of the rule → caught
- M10/M11 collapse the two absences into one rendering, in review-work and the Decision
  Brief → caught in both files

**Consumer-install proof (execution, not reading):** the suite was installed into a
throwaway git project with `tools/install-do-work-suite.sh`, producing the consumer layout
(`.claude/skills/do-work*/`, no `do-work/` under the skill root). From the three touched
files at their **installed** paths, 57 same-package and cross-package citations were
resolved against the installed tree — 0 unresolved — and none of the three cites `_dev/`.
The rule's hand-back path is project-root relative and the installed tree confirms
`do-work/` is a project-root sibling of `.claude/`, never a child of the vendored skill.

**New tests added:**
- `_dev/tests/contract-regressions.sh` — REQ-299 property check (Step 6's routed sections
  ≡ the hand-back contents, plus the classify-every-section guard and a vacuity guard) and
  eight pins on the rule's home, its condition-keyed reader list, its stated path root, and
  its two readers outside Step 8.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

## Review — 2026-08-21T03:19:19Z

**Overall: 94%**

| Dimension | Score |
|---|---|
| Requirements Compliance | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope Discipline | 100% |
| Risk | None |

**Acceptance: Pass.**

**Requirements Compliance.** All five requirements are delivered and each has evidence.
The condition-keyed rule, the Step 6 routing clause, the two-absence distinction in both
readers, the consumer-install execution, and the green gate.

**Finding-Closure Ratchet.** This REQ carries `review_generated: true`. Its captured GREEN
names the REQ-299 property check in `_dev/tests/contract-regressions.sh`; that check was
confirmed to fail on the untouched tree with
`routed by Step 6 but not carried by the hand-back: ['## Decisions']` and to pass after.
Closure evidence matches the named check.

**Findings**

- **F1 — Important, fixed in-REQ.** The first draft of `review-work.md`'s traceability
  bullet nested bold inside bold (`**… **Reading a Builder-Authored Section (any step)****`),
  which renders as literal asterisks. Rewritten as a lead-in plus a plain citation.
- **F2 — Minor, fixed in-REQ.** The `## Red-Green Proof` disclaimer first attributed the
  section to the orchestrator. It arrives with the REQ from capture or verify-requests;
  reworded to "it arrived with the REQ".
- **F3 — Minor, accepted.** The two-absence pin greps for a phrase
  (`could not be read` / `hand-back unread`) rather than the property. A rewording that
  keeps the distinction but drops both phrases would fail the pin spuriously. The failure it
  guards is deletion of the distinction, which is what the mutation test exercised; the
  false-positive cost is one grep to re-word.
- **F4 — Minor, accepted.** The property check's section-name pattern is
  `\`(## [A-Z][A-Za-z -]*)\``, so a future section whose name carries a digit or an
  ampersand would be invisible to it. Every section in the pipeline today is alphabetic; a
  wider pattern would start matching inline code that is not a section name.

**Code Quality (90%).** The check is 3 short passes over one extracted block and one table
row, with named intermediates and error text that says what to do. The 10-point deduction is
F3 and F4 together: two places where a phrase or a character class stands in for the
property.

## Lessons Learned

- **A rule scoped to a step cannot be inherited by readers at other steps.** REQ-270 wrote
  "Some substeps below" and closed one instance; the second instance was already in the file
  and outside that scope. When a rule's condition is "whenever anyone does X", its home has
  to be a place every doer of X reads, and its opening sentence has to name the condition
  rather than the location.
- **A check that needs a list can often be turned into a check that needs a mark.** The
  property "every section Step 6 tells the builder to author" resisted mechanical extraction
  until the requirement flipped: instead of the test knowing which sections those are, every
  section mention in the block states who writes it. The list moved into the prose, where it
  is co-located with what it describes and cannot be forgotten by a test author.
- **A guard passing is not a guard covering.** `shipped-package-reference-contract.sh`
  passed on planted dangling citations of three different same-package shapes. The only way
  to learn that was to plant them. Confirm a guard covers new text by breaking the text,
  never by reading the guard.
- **Mutation-test the pin, not just the fix.** M11 first came back green because the
  mutation removed one of two phrases the file carried. A mutation that does not actually
  remove the property proves nothing — check the mutation applied to everything the check
  can see.

## Orientation

**What changed in the map.** The rule for reading a section the builder authored is no
longer part of Step 8. It is a named section of `actions/work-reference.md` —
**Reading a Builder-Authored Section (any step)** — that any step, action, or report can
cite, and `actions/work.md` Step 8, `actions/review-work.md` Step 4, and the Decision
Brief all now cite it rather than restating it.

**What this makes true.** Under worktree fan-out, a builder's numbered decisions reach the
end-of-run brief and the code review instead of vanishing with its worktree. The set of
sections a builder is told to author and the set the hand-back is told to carry are now one
set, enforced mechanically, so the next section added to Step 6 cannot repeat this.

**Subsystem:** the do-work pipeline's instruction set — `actions/work.md`,
`actions/work-reference.md`, `actions/review-work.md`. Prime:
`_dev/primes/prime-action-files.md`.

**Discovered-task follow-up:** the shipped-reference-contract coverage gap was queued as REQ-312 (`pending-answers`, `impact-rule-change`).

---
id: REQ-270
title: Carry a worktree builder's hand-back sections into Step 8
status: completed
created_at: 2026-08-18T21:45:21Z
claimed_at: 2026-08-19T19:48:27Z
completed_at: 2026-08-19T20:07:36Z
commit:
status_changed_at: 2026-08-18T22:20:09Z
user_request: UR-055
addendum_to: REQ-259
domain: general
review_generated: true
sweep: true
sweep_key: worktree-handback-sections-unread-by-step-8
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: true
route: B
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-19T19:49:00Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
- skills/do-work/crew-members/general.md
- skills/do-work-toolbox/crew-members/general.md
---

# Carry a Worktree Builder's Hand-Back Sections into Step 8

## What

`skills/do-work/actions/work.md` Step 8 substep 4 reads `## Discovered Tasks` **from the REQ file**, describing it as "appended by the implementation agent as a separate section". In worktree dispatch mode the implementation agent **cannot write the REQ file** — the REQ lives in the main tree, which is exactly what `State stays home` forbids a builder to touch — so the section is never there, Step 8 finds nothing, and every out-of-scope find the builder recorded is silently dropped.

This is not hypothetical and it is not a builder error: it fired on REQ-259, whose builder correctly routed three Discovered Tasks to its hand-back per the dispatch brief. The section was absent from the REQ until the orchestrator noticed and transcribed it by hand. Nothing in the pipeline would have reported the loss.

**State the condition, not the instance.** `## Discovered Tasks` is the one that fired; the defect is the general shape — *any* Step 8 substep that expects to read a builder-authored section from the REQ file is silently disarmed under fan-out. Sweep the substeps for that shape rather than patching the one known case.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-259's independent review, Important finding 2 (gate: rule-change). Created `pending-answers` per the generation-≥2 cascade depth stop, since REQ-259 is itself `review_generated: true`; the reviewer assessed it as Important rather than critical-grade, so it does not pierce to `pending` on its own.

Worth knowing while it sits in the queue: **the loss is silent and unbounded in the past.** Every fan-out run before this fix depended on the orchestrator happening to transcribe by hand. This session's wave-1 REQs were transcribed manually once the review surfaced it, so nothing is lost here — but that is a person catching it, not the pipeline.

## Requirements

- Step 8's builder-authored-section reads work in worktree dispatch mode: the section is taken from the hand-back when the REQ file does not carry it.
- The rule is keyed on the **condition** (a Step 8 substep reading a section the builder authors) and marks any list of affected substeps illustrative, so a substep added later inherits the behavior instead of waiting to be remembered.
- The failure mode is loud rather than silent: a run that finds neither the REQ section nor a readable hand-back says so.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [x] REQ-259's review found that a worktree builder's `## Discovered Tasks` never reach Step 8, because Step 8 reads that section from the REQ file and a worktree builder may not write it — so every out-of-scope find is silently dropped unless the orchestrator transcribes it by hand. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — rely on the orchestrator transcribing hand-back sections as part of integration, and say so in the dispatch contract instead of changing Step 8.

**Answered [2026-08-18]:** User approved via `do-work clarify`, presented as a live defect in the pipeline itself that fired during this very run. Approved keyed on the general shape (any Step 8 substep reading a builder-authored section from the REQ file), not on `## Discovered Tasks` alone.

---

## Triage

**Route: B** - Medium

**Reasoning:** The failing substep is named exactly, but the REQ requires sweeping Step 8's substeps for the general shape rather than patching the one known case, and the hand-back's contract lives in a second file — the "where" needs discovery before anything is written.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Swept all eight Step 8 substeps for the shape (a substep reading a section the *builder*
authors from the REQ file). One fires today:**

| Substep | Reads | Author | Disarmed under fan-out? |
|---|---|---|---|
| 1 update frontmatter | frontmatter | orchestrator | no |
| 2 verify `## Implementation Summary` | that section | **orchestrator** (Step 6.25 writes it from the builder's report) | no |
| 3 builder-decided questions | `## Open Questions` `- [~]` items | **orchestrator** (Step 3.5, pre-dispatch) | no |
| 4 queue Discovered Tasks | `## Discovered Tasks` | **builder** | **yes — this is the instance** |
| 5 cycle detection | `addendum_to` frontmatter | capture | no |
| 6 archive move | file location | orchestrator | no |
| 7 deferred prime-link writes | Step 7.5's in-memory list | orchestrator | no |
| 7.5 calibration line | `estimate:` frontmatter | Step 3.6 | no |
| 8 worktree cleanup | git state | git | no |

So substep 4 is the only live case, which is exactly why the rule must be keyed on the
condition — a one-instance patch here would be indistinguishable from a fix and would
leave the next builder-authored section to be discovered the same way.

**A second half of the same defect, upstream in the same file:** `work.md` Step 6's
builder-instruction bullet tells the builder "Out-of-scope finds go to
`## Discovered Tasks`" — a section in the REQ file, which the very next bullet forbids a
worktree builder to write. The builder is being told to write somewhere it may not. Fixing
only the reader would leave the instruction contradicting itself.

**Where the section can come from instead:** `work-reference.md` → Worktree Dispatch Mode
already establishes `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md` as the one
main-tree path a builder may write (*Sole integrator*), and the Fan-Out run-directory table
already lists it as the per-builder output. Its stated contents are "branch, file manifest,
integration seams" — builder-authored REQ sections are not named, which is why REQ-259's
builder routing them there was correct behavior that nothing downstream read. For a remote
builder the same content travels on the branch instead (same section, `git show` after the
merge).

**Maintenance posture (`maintenance: true`):** subtract first. The fix removes a false
provenance claim ("appended by the implementation agent") rather than adding a caveat
beside it, and extends one existing table cell rather than adding a paragraph. One new
paragraph is earned — the condition itself has no home today.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — the condition-keyed reading rule at Step 8, the corrected substep 4, and Step 6's builder instruction
- `skills/do-work/actions/work-reference.md` (modify) — name builder-authored sections in the hand-back's stated contents
- `skills/do-work/crew-members/general.md` (modify) — added at review (D-03): its always-loaded Discovered-Tasks Contract contradicted the fix
- `skills/do-work-toolbox/crew-members/general.md` (modify) — byte-mirror of the above

**Files I will NOT touch:** `crew-members/background-agents.md` (owns the durability pattern, not this pipeline's read sites), `actions/review-work.md`.

**Acceptance criteria (restated from REQ):**
- [x] Step 8's builder-authored-section reads take the section from the hand-back when the REQ file does not carry it
- [x] The rule is keyed on the condition and marks any list of affected substeps illustrative
- [x] A run that finds neither the REQ section nor a readable hand-back says so
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## Decisions

- **D-01**: No new lock-in test in this REQ — but not for the reason first recorded. The
  original claim, that every mechanical check here would pin a spelling, was too strong,
  and the review named a concrete one that does not: extract every `## <Name>` section
  Step 6 instructs the **builder** to author, and assert each is named in
  `work-reference.md`'s hand-back contents. That check fails today on `## Decisions`, which
  is finding 2's whole subject — so it cannot be added here without also fixing a section
  outside this REQ's scope. It is carried into **REQ-299** as that REQ's Red-Green Proof,
  which is its correct home. DECIDE & STATE.
- **D-02**: Fixed Step 6's builder-instruction bullet as well as Step 8's reader. It told
  the builder to write `## Discovered Tasks` into the REQ file while the bullet directly
  above forbids a worktree builder to write the main tree — the same defect facing the
  other way. Repairing only the reader would leave the instruction self-contradicting and
  the builder still guessing, which is what REQ-259's builder had to do. Inside the
  declared write set. DECIDE & STATE.

- **D-03 (recorded at review; a scope judgment, so a Decision rather than a silent write)**:
  Extended the write set to both `crew-members/general.md` mirrors. The review found its
  always-loaded Discovered-Tasks Contract still telling every builder to append the section
  to the REQ file, with no worktree carve-out. `review-work.md` routes a stale restatement
  to a follow-up REQ, and I deliberately did not follow that route: this restatement does
  not merely go stale, it **defeats this REQ's own fix**. A worktree builder obeying it
  writes the REQ inside its own checkout, and the hand-back merge's queue guard drops
  exactly that commit — silently destroying the content, which is REQ-259's original bug
  reproduced through the one path this REQ claims to close. Shipping the reader-side fix
  while the always-loaded writer-side instruction contradicts it would have been a fix that
  does not work. Value: the fix holds end to end. Risk: low — one sentence, both mirrors
  kept byte-identical, nothing outside the worktree case changes. The second Important
  finding (`## Decisions`, read outside Step 8) is genuinely different scope and went to
  REQ-299.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/crew-members/general.md` (modified)
- `skills/do-work-toolbox/crew-members/general.md` (modified)

**What was done:** Step 8 gained one condition-keyed paragraph — *Where a builder-authored
section is read from* — stating that any substep reading a section the builder authors
reads the REQ file first and this REQ's hand-back second, with `## Discovered Tasks` named
as the only such substep today and explicitly illustrative. A second paragraph makes the
absent case loud: in worktree dispatch mode, a section in neither place **and** a missing or
unreadable hand-back is reported rather than read as "the builder found nothing". Substep 4
lost the provenance claim that was false under fan-out ("appended by the implementation
agent") and now points at the rule. Step 6's builder instruction says where the section goes
when the builder may not write the REQ. `work-reference.md`'s run-directory table now names
builder-authored REQ sections in the hand-back's stated contents, so the brief that quotes
that row tells the builder to put them there. After review, both `crew-members/general.md`
mirrors — the always-loaded Discovered-Tasks Contract every builder reads before the
route-specific instructions — say the same thing, so the writer-side instruction no longer
contradicts the reader-side fix.

## Qualification

Passed — 4 files verified in the diff, 4 requirements traced, P-A-U confirmed. Scope-drift
check clean against the extended Scope declaration (D-03). Instructions-only change, so
judgment check 6 (data flows) does not apply.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Passing (exit 0; 76 named script cases, contract regressions, both Go modules)

No red-green validation: the change has no executable surface. See D-01 for why the one
property check the review named belongs to REQ-299 rather than here.

**New tests added:** none — D-01.
**Existing tests updated:** none.

*Verified by work action*

## Review

**Reviewer:** independent agent, orchestrated mode against the working-tree diff.

| Dimension | Score |
|---|---|
| Requirements Compliance | 60% |
| Code Quality (prose) | 75% |
| Scope Discipline | 70% |
| Test Adequacy | N/A (instructions-only) |
| Risk | None |
| **Acceptance** | **Partial** |
| **Overall** | **58%** |

The reviewer independently re-ran the Step 8 substep sweep and confirmed the table: substeps
2 and 3 really are orchestrator-authored, so substep 4 was the only live case *inside Step
8*. It then found two Important things the sweep had not reached.

**F1 — a stale restatement that defeated the fix. Fixed here (D-03), not deferred.**
`crew-members/general.md`'s Discovered-Tasks Contract — always loaded at Step 6, ahead of
the route-specific instructions this REQ corrected — still said "Append a
`## Discovered Tasks` section to the REQ" with no worktree carve-out. A builder obeying it
writes the REQ inside its own checkout; the hand-back merge's queue guard then drops that
commit as committed queue state, destroying the content. That is REQ-259's bug reproduced
through the exact path this REQ closes, so the reader-side fix alone would not have worked.
Both mirrors repaired.

**F2 — the same defect one section over. Went to REQ-299.** `## Decisions` is
builder-authored (Step 6) and read at two sites *outside* Step 8 — review-work's
traceability check and the Decision Brief's HANDLED block — where this REQ's rule, whose
home is Step 8's preamble, structurally cannot reach. Silent in both directions: the review
reports clean because it cannot tell "no section" from "no decisions", and the user's
end-of-run report renders an empty HANDLED list. The reviewer's judgment that Step 8 was
the wrong home for the condition is right, and REQ-299 carries it with the property check
as its Red-Green Proof.

**Minor, fixed:** the new paragraph's central sentence was stated unconditionally while the
paragraph beside it gates on serial mode; a floor agent in serial mode would have gone
looking for a hand-back path that does not exist. Now gated explicitly.

**Minor, not fixed (noted):** `_dev/primes/prime-kanban-board.md:28`'s REQ-252 lesson
("capture Discovered Tasks in the REQ itself — hand-back prose is one slip from
evaporating") now reads as the opposite of this REQ's stance. It is an immutable historical
lesson entry and correct to leave; worth knowing it exists. `do-work/RESTART-PROMPT.md:61`'s
manual-transcription workaround goes stale on archive, but it is session scratch, not
shipped, and the workaround it describes remains harmless.

## Lessons Learned

**What worked:** Sweeping every substep in a table rather than checking the one the REQ
named. Eight rows, seven of them "no", and writing down *why* each was no is what made the
sweep auditable — the reviewer re-derived it and agreed, which a prose claim would not have
allowed.

**What didn't:** I swept the substeps and stopped at the file boundary. The condition was
"a reader of a builder-authored section", and I searched for readers only inside Step 8 —
so I never grepped for other places that read those sections, and never grepped for other
places that *tell the builder to write* them. Both misses were one `grep -rn '## Discovered
Tasks' skills/` away. **When a REQ says "state the condition, not the instance", the sweep
has to be for the condition across the repo, not for the instance across one step.**

**Worth knowing:** a fix in an action file can be silently defeated by an always-loaded
crew-member file, because crew rules load *before* the route-specific instructions and the
builder reads them as the contract. Any change to what a builder writes has to check
`crew-members/general.md` first — it is the instruction with the widest reach and the
lowest visibility, and it is mirrored byte-for-byte into `do-work-toolbox`, so both copies
move together or they drift.

## Orientation

A worktree builder's out-of-scope findings now survive the trip home: Step 8 reads
`## Discovered Tasks` from the hand-back when the REQ file cannot carry it, the builder is
told where to put it in both the action file and the always-loaded crew contract, and a run
that can find neither the section nor a readable hand-back says so instead of assuming the
builder found nothing. Lives in the work-loop orchestration instructions
(`_dev/primes/prime-action-files.md`). No map change — one new rule inside Step 8, no new
file, no new mechanism. Prime spot-check: `prime-action-files.md`'s referenced paths all
still exist.

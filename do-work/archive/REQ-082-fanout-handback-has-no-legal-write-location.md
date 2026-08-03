---
id: REQ-082
title: The fan-out hand-back file has no legal write location
status: completed
claimed_at: 2026-08-03T22:14:08Z
completed_at: 2026-08-03T22:17:23Z
commit: 1cff0a7
kb_status: pending
route: B
created_at: 2026-08-03T17:09:21Z
user_request: UR-016
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-073, REQ-084]
batch: audit-remediation-external
addendum_to: REQ-073
write_set: [actions/work-reference.md, _dev/tests/contract-regressions.sh]
---

# The Fan-Out Hand-Back File Has No Legal Write Location

## What

Fan-Out Dispatch makes a per-builder output file mandatory — `REQ-NNN-handback.md` inside
`do-work/runs/work-<timestamp>/` — and `crew-members/background-agents.md` is explicit that the
sub-agent writes that file itself. Worktree Dispatch Mode is equally explicit that a builder **never
writes the main tree**, and `do-work/` exists in the main tree only. There is no location satisfying
both rules, so the mandatory hand-back has no legal execution.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty. Loaded `crew-members/general.md`, `coding-guardrails.md`,
  `maintenance.md` (`maintenance: true`), `testing.md` (`tdd: true`). The direction was settled by the
  user at capture (Open Questions, resolved), so the plan is placement and wording only: the exception
  goes in the sentence it modifies (*Sole integrator*), the table row and the brief-delivery trap
  *reference* it rather than restating it, and *State stays home*'s three-item list is replaced — not
  extended — with its condition. Budget: one sentence for the exception, per the Constraints.
- [x] **[APPLY]:** Two files, both declared. Four prose edits and four assertions; nothing else touched.
- [x] **[UNIFY]:** `git diff --stat` → `actions/work-reference.md` 4 hunks, +4/−4 lines net of
  rewrites; `_dev/tests/contract-regressions.sh` +33. Verified: `shellcheck` clean, `bash -n` parses,
  suite exits 0. Read each hunk against the Constraints — sole-integrator is not weakened anywhere
  else (the builder still never touches the queue, a status, an archive move, `CHECKPOINT.md`,
  `actions/version.md`, `CHANGELOG.md`, a seam, or a sibling's anything), the survivable-not-prevented
  framing at the Fan-Out closing line is untouched, and the forbidden-token sweep is still green (no
  lock, heartbeat, registry, or liveness probe entered). No debug artifacts.

## Why

An agent reaching this contradiction has three moves and two of them corrupt the run. It can write the
main tree (violating sole-integrator, the rule that makes worktree isolation mean anything); it can
write the worktree's own `do-work/runs/…`, where the file lands in the builder's branch and gets swept
into the merge as committed scratch while the orchestrator reads nothing; or it can skip the file and
return findings in conversation, which is the exact failure the durability pattern exists to prevent.
Nothing in the shipped prose tells it which.

This is the output-direction twin of a trap REQ-073 already closed for inputs. Its requirement 8 fixed
brief *delivery* — "the brief must reach the builder as prompt content or an absolute main-tree path.
A repo-relative path resolves inside the worktree, against its own stale tracked copy of `do-work/`" —
and the return path has the identical resolution problem with no equivalent sentence.

## Context

The three statements, verbatim:

- `actions/work-reference.md:313` (Fan-Out Dispatch, the mandatory run-directory table):
  > | per-builder output | `REQ-NNN-handback.md` — branch, file manifest, integration seams |

  introduced by "**The run directory is mandatory here, not optional.**" (`:307`).

- `crew-members/background-agents.md:44-48`:
  > **Each sub-agent writes its own findings file; returns only a one-line status.** Give every
  > sub-agent an output path inside the run directory … The agent writes its *full* findings to that
  > file and returns **only a one-line status** to the orchestrator — never the full findings inline.

- `actions/work-reference.md:275` (Sole integrator):
  > The builder never writes the main tree or its branch.

  reinforced by `:273` (State stays home): "`do-work/` — the queue, `working/`, `CHECKPOINT.md` —
  exists in the main tree only."

**Two secondary problems in the same neighbourhood.**

1. *State stays home* enumerates three things (`the queue, working/, CHECKPOINT.md`). `do-work/runs/`
   did not exist when that sentence was written and is not in the list, so a reader can argue the run
   directory is out of scope — which is exactly the closed-enumeration failure `CLAUDE.md` warns
   about. The sentence needs its condition stated, not a fourth item appended.
2. `crew-members/background-agents.md:33-41` makes the run directory an ordinary **committable** path
   under `do-work/`, deliberately, so a mid-run commit carries it. Combined with a builder writing
   into it from inside a worktree, that turns run scratch into branch content — and REQ-084's new
   merge-base probe would then correctly flag the builder for writing queue state. The two REQs must
   not disagree about whether that file is a violation.

**Why no check caught it.** The contradiction is between two files that no assertion compares, and
`_dev/tests/contract-regressions.sh` has no assertion touching the hand-back file at all. It is also
unreachable in practice today: everything since REQ-073 shipped has been built serially (see
REQ-085), so no run has ever produced a hand-back file.

## Detailed Requirements

1. **State the one exception explicitly in `actions/work-reference.md`**, at *Sole integrator* — the
   sentence the exception modifies — and reference it from the Fan-Out Dispatch table row rather than
   restating it. A builder may write **exactly one** path: its own
   `do-work/runs/work-<timestamp>/REQ-NNN-handback.md`.
2. **The path reaches the builder the same way its brief does** — as an absolute main-tree path, never
   repo-relative. Requirement 8's existing trap sentence covers the mechanism; say that it applies in
   both directions instead of writing a second copy of the reasoning.
3. **The builder never stages, commits, or merges that file.** It is a main-tree working file owned by
   the orchestrator's run directory, not branch content. State this as part of the exception, because
   the natural reading of "you may write it" includes committing it.
4. **The exception is bounded to that one filename.** Not "the run directory", not "files under
   `do-work/runs/`" — one path per builder, derived from its own REQ id. A builder writing a sibling's
   hand-back file, the manifest, or anything else remains a sole-integrator violation.
5. **Restate *State stays home* as a condition, not a list.** Replace the three-item enumeration with
   the rule it is trying to express — every path under `do-work/` is main-tree-only and
   orchestrator-owned — and mark any examples as illustrative. Per `CLAUDE.md` → Closed Enumerations
   Go Stale, grep for other enumerations of the same set and generalize each.
6. **`manifest.md` stays the orchestrator's.** `background-agents.md:51-56` has the orchestrator
   maintain it per wave; make sure the amended prose cannot be read as licensing a builder to update
   it. This is the natural over-reach from requirement 1.
7. **Reconcile with REQ-084 in whichever order they land.** REQ-084 adds a probe for a builder branch
   carrying `do-work/` changes. The hand-back file must not trip it — which it cannot, if requirements
   2 and 3 hold, because the file is written to the main tree and never committed. Whichever REQ lands
   second must confirm that in its Qualification rather than assuming it.
8. **Add a contract assertion pinning the exception.** The failure mode is a later maintenance pass
   reading "the builder never writes the main tree" as absolute and deleting the carve-out as
   redundant, which silently restores the contradiction. Assert that the exception exists and that it
   names a single path.

## Constraints

- **`maintenance: true`.** The candidate fix narrows two shipped instructions (`Sole integrator`'s
  absolute prohibition; *State stays home*'s enumeration), so `crew-members/maintenance.md`'s
  delete-before-you-add rule governs. Requirement 5 is a replacement, not an addition, and
  requirement 1 should cost one sentence — if the exception needs a paragraph, the shape is wrong.
- **Do not weaken sole-integrator anywhere else.** The builder still never touches the queue, a
  status, an archive move, `CHECKPOINT.md`, `actions/version.md`, `CHANGELOG.md`, an integration seam,
  or a sibling's anything. This REQ opens one file, by name.
- **Do not describe the durability pattern as preventing failures.** `background-agents.md:11-14` and
  Fan-Out Dispatch's closing line (`actions/work-reference.md:317`) both require the
  survivable-not-prevented framing to be carried, not softened.
- **No new durable coordination state.** The forbidden-token sweep
  (`_dev/tests/contract-regressions.sh:132-137`) must stay green: no lock, heartbeat, claim registry,
  or liveness probe enters via the hand-back path.
- The run directory's lifecycle (created before any spawn, deleted once consumed) is
  `background-agents.md`'s, and is cited rather than restated.

## Dependencies

`addendum_to: REQ-073`, which introduced the mandatory hand-back. `related: REQ-084` for requirement 7.
No `depends_on`: buildable immediately in either order, and REQ-085's live run is the thing that would
have surfaced this, not a prerequisite for fixing it.

## Builder Guidance

**Certainty: Firm on the diagnosis and on the chosen direction; open on wording and placement.** The
direction was decided by the user at capture time — see the resolved question below — so do not
re-open it or re-derive the trade-off.

Prefer the shortest wording that survives a maintenance pass. This section of
`actions/work-reference.md` is already the file's largest, and REQ-073's Builder Guidance rule still
applies: anything restating what Worktree Dispatch Mode already says gets cut rather than written.

## Open Questions

- [x] The hand-back file must be written somewhere, and both candidate homes are currently forbidden.
  Should the builder be granted a narrow main-tree write, or should the file be dropped in favour of
  the builder reporting its manifest in its reply?
  → **Grant the narrow write.** Resolved by the user at capture time (2026-08-03, ask-tool prompt
  during `do-work capture-request`). The builder may write exactly one path — its own
  `REQ-NNN-handback.md`, by absolute main-tree path — and never commits it. Reasoning: the file exists
  *because* the transcript is not durable (`crew-members/background-agents.md` § Why This Matters), so
  dropping it to preserve an absolute prohibition trades a real recovery property for a tidier rule.
  **Out of scope as a result:** the return-the-manifest-in-the-reply shape, and any broadening of the
  exception past that one filename.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Red-Green Proof

**RED prompt/case:** A contract-suite assertion (`_dev/tests/contract-regressions.sh`) that fails
while `actions/work-reference.md` mandates a per-builder hand-back file without naming a legal write
location for it. Concretely: assert that the *Sole integrator* paragraph contains a bounded exception
naming `handback`, and that *State stays home* no longer expresses its scope as a three-item list.

**Why RED now:** `grep -c "handback" actions/work-reference.md` finds the mandate at `:313` and
nothing at `:275`; `_dev/tests/contract-regressions.sh` has no assertion mentioning the hand-back file.

**GREEN when:** The assertions pass; the full contract suite stays green; and a reader can answer
"where does the builder write its hand-back file, and who commits it?" from `actions/work-reference.md`
alone, without opening `crew-members/background-agents.md`.

**Manual proof, for the part an assertion cannot reach:** the amended prose must let a human dispatch
two builders and receive two hand-back files without either builder writing anything the orchestrator
did not authorize. That check is REQ-085's run, not this REQ's — note the dependency in the review
rather than claiming it here.

**Validation:** User confirmed — the direction was chosen from an explicit two-option prompt at capture
time; the contradiction itself was verified by reading all three cited statements against the tree.

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the
external audit, and the batch constraints.

---
*Source: external audit finding F2, second claim (P1) — "background-agent rules require builders to
write into the main-tree run directory, while worktree rules forbid builders from writing the main
tree" — accepted by `do-work validate-feedback` triage as the sharpest finding of the six.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The contradiction, its three cited statements, and the chosen resolution were all
supplied — the direction came from a user answer at capture, which Builder Guidance forbids
re-opening. What needed discovery was placement: which sentence carries the exception, and whether
the closed enumeration in requirement 5 has siblings elsewhere.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**All three statements re-read against the tree** (line numbers moved: REQ-077's new section pushed
Worktree Dispatch Mode down). *State stays home* is `:287`, *Sole integrator* `:289`, the Fan-Out
table row `:327`, the brief-delivery trap `:333`.

**Requirement 5's sibling-enumeration grep came back empty.** Searching `actions/`, `crew-members/`
and `docs/` for any other place expressing "what lives under `do-work/`" as a list found exactly one
site — the sentence being replaced. The rule is stated once, which is why generalizing it in place is
the whole fix rather than the first of several.

**Where the exception belongs, and where it does not.** The Constraints budget one sentence and
REQ-073's rule ("anything restating what Worktree Dispatch Mode already says gets cut") both point the
same way: the carve-out goes *inside* the `Sole integrator` sentence it modifies, and the two places
that would otherwise duplicate it — the Fan-Out table row and the brief-delivery trap — get a
back-reference and a clause respectively. The GREEN condition is the test of whether that was enough:
a reader must be able to answer "where does the builder write it, and who commits it?" from
`actions/work-reference.md` alone.

**`crew-members/background-agents.md` deliberately not touched.** It says "give every sub-agent an
output path inside the run directory" without requiring it be absolute — which is correct *there*,
because that file governs background fan-outs generally, most of which have no worktree and no
resolution hazard. The absolute-path requirement is worktree-specific, so it belongs in the
worktree-specific section, and adding it to the generic file would be the restatement the guidance
says to cut. Recorded so the omission reads as a decision rather than a miss.

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — the bounded exception at *Sole integrator*; *State stays
  home* restated as a condition; the Fan-Out table's output row and `manifest.md` row; the
  brief-delivery trap extended to the return direction
- `_dev/tests/contract-regressions.sh` (modify) — four assertions pinning the exception's existence,
  its no-commit clause, the condition form, and the absence of the old enumeration

**Files I will NOT touch:** `crew-members/background-agents.md` (reasoning above — the run directory's
lifecycle and the generic findings-file pattern stay its own), and every other statement of
sole-integrator.

**Acceptance criteria (restated from REQ):**
- [ ] The exception is stated at *Sole integrator*, referenced (not restated) from the table row
- [ ] The path reaches the builder absolute, and the trap sentence says it applies both directions
- [ ] The builder never stages, commits, or merges the file
- [ ] The exception is bounded to one filename, derived from the builder's own REQ id
- [ ] *State stays home* is a condition, not a list; sibling enumerations generalized
- [ ] `manifest.md` cannot be read as a builder's to write
- [ ] REQ-084 reconciliation confirmed rather than assumed
- [ ] A contract assertion pins the exception and its single-path bound

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** `Sole integrator` now reads "The builder never writes the main tree or its branch,
**with exactly one exception: its own `do-work/runs/work-<timestamp>/REQ-NNN-handback.md`**" —
reached by the absolute main-tree path the orchestrator hands it, **never staged, committed, or
merged**, and explicitly bounded: a sibling's hand-back, `manifest.md`, anything else under
`do-work/runs/`, and every other main-tree path remain violations. One sentence, as the Constraints
budgeted, carrying the *why* (the file exists because the transcript is not durable) so a maintenance
pass cannot read it as redundant.

`State stays home` was **replaced, not extended**: "Every path under `do-work/` exists in the main tree
only and is the orchestrator's — the queue, `working/`, `CHECKPOINT.md` and `runs/` are examples of the
rule, not its extent." The old three-item list predated `do-work/runs/` entirely, which is what let a
reader argue the run directory was out of scope.

The Fan-Out table's output row now points at the exception rather than restating it; its `manifest.md`
row says outright that the file is the orchestrator's and never a builder's (requirement 6's
over-reach, closed). The brief-delivery trap sentence gained the return direction, where the same
repo-relative resolution is *worse*: the write succeeds, lands in the builder's branch, and the
orchestrator reads nothing.

Four assertions: the exception exists, it says never-committed, the condition form is present, and the
old enumeration is absent. The last is a negative on purpose — the enumeration is the thing that goes
stale, so re-adding it must fail even if the condition sentence survives alongside it.

## Qualification

**Passed** — 2 files verified; 8 requirements traced.

- r1 → the exception at `Sole integrator`; the table row references it. r2 → the trap sentence's new
  both-directions clause. r3 → "never staged, committed, or merged", asserted. r4 → the bound is
  spelled out by exclusion (sibling / `manifest.md` / anything else under `runs/`). r5 → the condition
  replaces the list; the sibling grep found no other enumeration. r6 → `manifest.md` row amended.
  r7 → see below. r8 → four assertions, each observed failing.
- **Requirement 7 — REQ-084 reconciliation, confirmed rather than assumed.** REQ-084 has not been
  built yet, so this REQ lands first. The confirmation available now is structural: the hand-back file
  is written **to the main tree** (never the worktree) and is **never committed**, so it cannot appear
  in a builder branch's `do-work/` diff and cannot trip a merge-base probe for a builder carrying
  queue-state changes, however that probe is written. REQ-084 lands second and its Qualification must
  restate this against its actual implementation — flagged in its own terms below so it is not lost.
- **P-A-U audit:** three boxes with evidence; diff matches. **Contamination check:** REQ-077, 078 and
  079 all touched `actions/work-reference.md` and `_dev/tests/contract-regressions.sh`. Expected —
  four REQs in one batch amending the same two contract files — and verified by reading: this REQ's
  hunks are in Worktree Dispatch Mode and a new assertion block, none of which any prior REQ touched.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Passing (exit 0)

**Red-green validation:** traced to `## Red-Green Proof`.

- *RED (as captured):* at `HEAD`, `grep -c handback actions/work-reference.md` finds the mandate in the
  Fan-Out table and nothing at `Sole integrator`, and no assertion in the suite mentions the hand-back
  file. ✗ before → ✓ after.
- *GREEN condition, observed:* deleting the carve-out from `Sole integrator` fails the suite with
  **two** named failures (the exception's absence and the no-commit clause's absence). Restored → 0.
- *The negative assertion earns its place:* restoring the three-item `State stays home` enumeration
  fails the suite with **two** more (the condition form missing, and the stale enumeration present).
  Either one alone would have been satisfiable by a sentence that says both things — the pair is what
  makes re-adding the list a failure. Restored → 0.
- *GREEN condition, the readability half:* `actions/work-reference.md` now answers "where does the
  builder write its hand-back file, and who commits it?" in one sentence, without opening
  `crew-members/background-agents.md`. Checked by reading the `Sole integrator` paragraph cold.
- *Constraints re-verified as negatives:* the forbidden-token sweep is green (no lock/heartbeat/
  registry/liveness token entered), and the survivable-not-prevented line at the Fan-Out section's
  close is byte-identical.

**Not tested, and cannot be here:** the manual proof — two builders dispatched, two hand-back files
received, neither builder writing anything unauthorized. That is REQ-085's live run. The REQ says to
note the dependency rather than claim it, and this is that note.

**New tests added:** 4 assertions in `_dev/tests/contract-regressions.sh`.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

## Decisions

- **D-01 — The exception lives inside the `Sole integrator` sentence, not as its own paragraph.**
  DECIDE & STATE. The Constraints said one sentence and that "if the exception needs a paragraph, the
  shape is wrong." Putting it in the sentence it modifies also defends it: a maintenance pass reading
  `Sole integrator` cannot see the absolute prohibition without also seeing the carve-out and its
  reason. A separate paragraph would be deletable on its own, which is failure mode requirement 8
  names. Reversible.
- **D-02 — The bound is stated by exclusion, not just by naming the one path.** DECIDE & STATE.
  "Exactly one exception: its own `…/REQ-NNN-handback.md`" is already bounded, but requirement 4 is
  specifically about the over-reads — "the run directory", "files under `do-work/runs/`", a sibling's
  file, `manifest.md`. Naming those as still-violations costs half a sentence and closes each one, and
  they are exactly what an agent under time pressure would generalize to. Reversible.
- **D-03 — `crew-members/background-agents.md` left untouched.** DECIDE & STATE. Its "an output path
  inside the run directory" is correct for the file's actual scope: background fan-outs generally, most
  of which have no worktree and therefore no path-resolution hazard. The absolute-path requirement is
  worktree-specific and belongs in the worktree-specific section, where it now sits next to the
  identical rule for the inbound brief. Adding it to the generic file would be a restatement REQ-073's
  guidance says to cut, and would put a worktree rule in a file loaded for non-worktree runs.
  Reversible if a fan-out ever ships that reads only that file.

## Lessons Learned

**What worked:**
- **Fixing the enumeration with a *negative* assertion alongside the positive one.** A single
  "contains the condition" assertion is satisfiable by a sentence that states the condition *and* keeps
  the stale list beside it — which is the likeliest accidental regression, since adding is easier than
  replacing. The pair makes re-adding the list a failure on its own.
- **Putting the carve-out inside the sentence it modifies.** It is the cheapest possible defence
  against the specific failure requirement 8 describes, and it cost nothing.

**What didn't:** nothing failed. The sibling-enumeration grep requirement 5 mandated came back empty,
which is worth recording as a *result* — the previous two REQs in this batch both had their site
inventories turn out to be floors, and this one genuinely did not.

**Worth knowing:**
- **A contradiction between two files is invisible to a per-file assertion suite, and this one was
  also unreachable in practice.** Nothing has fan-out dispatched since REQ-073 shipped, so no run has
  ever produced a hand-back file — the contract has been broken for its entire life without a single
  observable symptom. Contract defects in unexercised paths are found by reading, not by running,
  which is an argument for REQ-085 actually happening.
- **The dangerous direction of a path-resolution bug is the one where the write succeeds.** Inbound,
  a repo-relative brief path yields nothing or a stale snapshot — visible. Outbound, the same mistake
  writes a real file into the builder's branch, gets swept into the merge as committed scratch, and
  the orchestrator quietly reads nothing. Same mechanism, and the return direction is worse, which is
  why the trap sentence now says so rather than being left as an exercise.

## Orientation

The fan-out hand-back file now has one legal place to be written. Worktree Dispatch Mode granted
builders a mandatory output file and simultaneously forbade the only tree it could live in; a builder
hitting that had three moves and two of them corrupted the run silently. `Sole integrator` now names
the single path a builder may write — its own `REQ-NNN-handback.md`, absolute, never committed — and
everything else stays forbidden by name. Lives in the fan-out contract
(`actions/work-reference.md` → Worktree Dispatch Mode).

[MAP CHANGED] Sole-integrator is no longer an absolute prohibition; it is a prohibition with exactly
one named exception. And `do-work/`'s main-tree-only rule is now stated as a condition covering every
path under it, rather than a list of the three directories that existed when it was written — so
`do-work/runs/` is covered, and so is the next one.

`prime_files` is empty, so no prime staleness spot-check applies.

## Review

**Overall: 96%** | 2026-08-03T22:14:08Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Findings:** 0 important, 3 minor
**Acceptance:** Partial — every assertable condition passes and each assertion was observed failing,
but the REQ's own **manual proof** (two builders, two hand-back files, nothing unauthorized written)
belongs to REQ-085 and has not run. The REQ instructs that this be noted rather than claimed.
**Suggested testing:** 2 items
**Follow-ups created:** None

### Findings

- **[Minor] Requirement 7's reconciliation is structural, not observed.** REQ-084 has not been built,
  so "the hand-back file cannot trip the builder-wrote-`do-work/` probe" rests on the file being
  main-tree and uncommitted — true by construction, but no probe exists yet to run it against. Flagged
  into REQ-084's own Qualification, which the REQ asks the second-lander to do.
- **[Minor] The prose is enforced; the behaviour is not.** Four assertions pin what
  `actions/work-reference.md` *says*. Nothing stops a real builder writing a sibling's hand-back file
  — that would be caught, if at all, by REQ-084's probe (for committed writes) or by nothing at all
  (for uncommitted main-tree writes). Inherent to a prose-driven pipeline; recorded so a green suite
  is not mistaken for enforcement.
- **[Minor] `crew-members/background-agents.md` still says "an output path inside the run directory"
  without requiring absolute.** Correct for its generic scope (D-03), but a fan-out orchestrator who
  reads only that file could hand a relative path. The mitigation is that the worktree section — which
  a worktree fan-out must read — now states it in both directions.

### Requirements Checklist

| # | Requirement | Status |
|---|---|---|
| 1 | Exception at *Sole integrator*, referenced from the table | Delivered — D-01 |
| 2 | Absolute path, both directions | Delivered — trap sentence extended |
| 3 | Never staged, committed, or merged | Delivered — and asserted |
| 4 | Bounded to one filename | Delivered — bound stated by exclusion, D-02 |
| 5 | *State stays home* as a condition; siblings generalized | Delivered — replaced, not extended; sibling grep empty |
| 6 | `manifest.md` stays the orchestrator's | Delivered — table row amended |
| 7 | Reconcile with REQ-084 | Delivered as far as possible — structural confirmation, flagged for REQ-084 |
| 8 | Contract assertion pinning the exception | Delivered — 4 assertions, all observed failing |

### Acceptance Testing

Suite exits 0. Deleting the carve-out produces two named failures; restoring the stale enumeration
produces two more; both reverted cleanly. The Constraints were re-verified as negatives: the
forbidden-token sweep is green and the survivable-not-prevented framing is byte-identical.
`shellcheck` clean, `bash -n` parses, `qualify.sh` and `scope-drift.sh` OK.

### Suggested Additional Testing

- **REQ-085's live two-builder run** is the outstanding proof, and this REQ is one of the reasons it
  matters: the contract it fixes has never been exercised.
- **A builder handed a repo-relative hand-back path on purpose**, to confirm the failure is as
  described (a real file in the wrong tree, silently) rather than an error. Worth doing once inside
  REQ-085's run rather than as a separate exercise.

*Reviewed by review-work action (pipeline mode, in-session)*

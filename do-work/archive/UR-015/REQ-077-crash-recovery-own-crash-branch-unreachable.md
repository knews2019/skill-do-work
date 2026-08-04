---
id: REQ-077
title: Crash recovery's own-crash branch is unreachable, and its retired premise survives in the same file
status: completed
claimed_at: 2026-08-03T21:29:26Z
completed_at: 2026-08-03T21:44:08Z
commit: 598ef35
kb_status: pending
created_at: 2026-08-03T16:53:42Z
user_request: UR-015
domain: general
prime_files: []
tdd: true
route: C
depends_on: []
maintenance: true
addendum_to: REQ-071
write_set: [actions/work.md, actions/work-reference.md, _dev/tests/contract-regressions.sh]
---

# Crash recovery's own-crash branch is unreachable, and its retired premise survives in the same file

## What

REQ-071 (v0.164.0) gated crash recovery on the checkpoint's `## In Progress (interrupted)` record.
That record has exactly one write site — `actions/work.md:627`, Step 10, **session end** — so on a
hard crash it does not exist, every abandoned REQ classifies as a foreign claim, and the REQ leaves
the pipeline permanently. Separately, the premise REQ-071 retired still sits ~110 lines below the
paragraph REQ-071 rewrote, in the same file, and the guard meant to pin it absent matches a different
wording.

Two findings, one REQ: same file, same premise, one sweep.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty. Loaded `crew-members/general.md`, `coding-guardrails.md`,
  `maintenance.md` (`maintenance: true`), `testing.md` (`tdd: true`). Approach in `## Plan` below:
  the fix's added half (a claim-time write) goes in a named reference procedure that Step 2 points
  at; the fix's subtractive half retires the premise's three fingerprints and *merges* two narrow
  guards into one generalized sweep rather than adding a third. Per `maintenance.md` § 2, the
  drift's cause is a **stale source** (a sentence describing machinery that no longer matches the
  contract) and a **guard too narrow** (a haystack that cannot see the line it pins) — both fixed by
  removal/widening; only requirement 1 is a genuine addition, and it has a replay case (§ 3) in the
  Red-Green Proof.
- [x] **[APPLY]:** Three files, all declared in `## Scope`, no others. No new machinery: the
  in-progress record already existed and already had a consumer; only its write site moved.
- [x] **[UNIFY]:** `git diff --stat` → 3 implementation files (`actions/work.md` 15±,
  `actions/work-reference.md` 22±, `_dev/tests/contract-regressions.sh` 67±) plus `do-work/` state.
  Verified per file: **`actions/work.md`** — read all five hunks (Step 2 substeps 1–3, Step 8
  substep 6, the failure path, the blocked flip, Step 10's session-start note, both checklist
  lines); confirmed no surviving premise fingerprint by grep. **`actions/work-reference.md`** — read
  all four hunks; confirmed the six REQ-071 phrases the suite pins are byte-identical and that the
  new section sits *outside* the `crash_recovery_block` sed range so no existing assertion's haystack
  changed. **`_dev/tests/contract-regressions.sh`** — `bash -n` parses, `shellcheck` clean (no new
  warnings), suite exits 0. No debug artifacts in the diff (`console.log`/`debugger`/`TODO`/`FIXME`
  grep over added lines: none).

## Why

REQ-071 traded a **lossy but visible** failure (recovery strips uncommitted `## Plan`/`## Exploration`)
for a **silent and permanent** one (the REQ is never worked again). That trade was never stated as
such. The 0.164.0 entry discloses the mechanism in a bullet — "No checkpoint at all is treated as
ambiguous, not as permission — which is the common case, since a hard crash never gets to write one" —
but frames it as a safety property, not as the automatic recovery path becoming unreachable.

The premise remnant matters for the same reason REQ-075 existed: a reader reasoning forward from
`actions/work.md:224` concludes the opposite of the contract stated at `:113`.

## Context

**Finding F1 — the own-crash branch is unreachable.** Full trace, verified against the tree:

1. Hard crash → no `do-work/CHECKPOINT.md` (only Step 10 writes it; `grep -n "CHECKPOINT.md" actions/*.md`
   confirms no other write site).
2. Next run: `actions/work-reference.md:224` requires the REQ be *named* in `## In Progress (interrupted)`.
   Absent checkpoint → **foreign claim** → left byte-identical, not moved.
3. Under 3h → reported, **no takeover offered at all**. Unattended → outcome (b), never taken over.
4. `actions/cleanup.md` Pass 0 step 5 sweeps `do-work/working/` only for **terminal** statuses; a
   `claimed` file is untouched there too.
5. Step 1's selection scan globs `do-work/queue/` — the file is in `working/`, so it is never selected.

Net: a 30-second crash silently removes the REQ from the pipeline, leaving one warning line per
subsequent run, until a human runs `actions/forensics.md` Check 1's manual reset.

**Finding F3 — the retired premise survives in the file REQ-071 edited.**

- `actions/work.md:113` (REQ-071 rewrote): "**A claimed REQ is not automatically this session's to reclaim**"
- `actions/work.md:224` (Step 2, untouched since REQ-069 / `76cdf39`): "every `working/` file is this
  session's to recover under the exclusive-session model, so no lock or claim record is needed"

The guard at `_dev/tests/contract-regressions.sh:364` pins the string
`Every .working/. file is this session's own leftover` — capital E, "own leftover". Line 224 reads
lowercase "every … is this session's **to recover**". Different fingerprint, suite green, premise
alive — and it is the sentence that justifies why Step 2 writes no claim record, which is exactly what
F1's fix changes.

This is REQ-075's own lesson ("a retired premise leaves two fingerprints — the thing it said and the
thing it was called") going unapplied to REQ-071, by a REQ that shipped later the same day.

## Detailed Requirements

1. **Write `## In Progress (interrupted)` at claim time**, in `actions/work.md` Step 2, alongside the
   `status: claimed` / `claimed_at` frontmatter write. This is the whole fix for F1: it makes the
   own-crash branch reachable in the ordinary case while leaving REQ-071's protection intact, because
   a claim this session never made still never appears there.
2. **The record must be a list, not a single id.** Fan-out dispatch (REQ-073) claims N REQs
   concurrently under one owner; a singular record would classify N−1 of them as foreign claims after
   a crash. Verify the Session Checkpoint Template in `actions/work-reference.md:762` accepts a list,
   and amend it if it does not.
3. **The claim-time write must not become a second checkpoint owner.** Step 10 still writes the full
   checkpoint. Step 2 appends or refreshes only the in-progress record — specify which, and specify
   what happens when no `CHECKPOINT.md` exists yet at claim time (create it with just that section is
   the expected answer; state it rather than leaving it inferred).
4. **Remove the record when the REQ leaves `working/`** — Step 8's archive move — or a completed REQ
   stays listed as in-progress and the next run reports a contradiction against
   `actions/work-reference.md:224`'s "finding one there is a contradiction to report" rule.
5. **Rewrite `actions/work.md:224`** so Step 2's no-claim-record sentence agrees with `:113`. Note the
   sentence's *stated reason* changes with requirement 1: after this REQ there **is** a claim record
   (the checkpoint's in-progress list), so the correct wording is not a reworded version of the old
   claim but a statement of what that record is and is not (it is recovery's classification input, not
   a lock; it coordinates nothing and no second owner reads it).
6. **Broaden the guard at `_dev/tests/contract-regressions.sh:364`** from the literal REQ-069 sentence
   to a pattern that catches both fingerprints of this premise. Per the Closed Enumerations rule, state
   the trigger condition in the assertion's comment and treat the matched wordings as illustrative.
7. **Disclose the F1 regression in the changelog entry**, not only the fix. A reader of 0.164.0 could
   not tell that automatic recovery had stopped working; the entry for this REQ should say so plainly.
8. **Do not weaken REQ-071.** The three-hour threshold still only gates the *offer*, a human still
   authorizes every takeover, the unattended path still never strips, and an absent checkpoint is
   still ambiguous rather than permission. Those are load-bearing and stay.

## Constraints

- `actions/work.md` and `actions/work-reference.md` are the same pair REQ-071 and REQ-074 edited.
  Re-read `actions/work-reference.md`'s Crash Recovery (Step 1) in full before touching substep 1 —
  REQ-074's `status_changed_at` stamp and the `claimed_at`-read-before-substep-1 ordering trap both
  live there and must survive unchanged.
- The claim-time write is one small file write per claim. Do not let it grow into a lock, a heartbeat,
  or a liveness check — that is precisely the machinery REQ-069 deleted at 0.161.0 and REQ-073
  declined to revive. It is a classification *input*, nothing more.
- Prescribed-command hygiene applies: any command this REQ adds must actually emit what the following
  step consumes, and `do-work/CHECKPOINT.md` must be reached by a deterministic path, not a variable
  inherited from an earlier block.

## Dependencies

`addendum_to: REQ-071` — this amends REQ-071's classification gate. Interacts with REQ-074 (its stamp
sits in the substeps this REQ makes reachable) and REQ-073 (its N-concurrent-claims case drives
requirement 2). No `depends_on`: buildable immediately.

## Builder Guidance

**Certainty: Firm on the diagnosis, firm on the fix's shape, open on placement details.** The trace in
Context was verified end-to-end against the tree; do not re-derive it from scratch, but do re-verify
the two line numbers, since REQ-077's own edits will move them.

The open latitude is *where exactly* the claim-time write goes (a Step 2 substep vs. a named procedure
in `actions/work-reference.md` that Step 2 points at) and how the removal in requirement 4 is worded.
Prefer the reference-file procedure if the prose exceeds a couple of sentences — `actions/work.md` is
the dispatcher and already long.

Keep the diff surgical. This is a `maintenance: true` REQ: `crew-members/maintenance.md`'s
delete-before-you-add rule applies, and requirement 5 is a rewrite, not an addition.

## Red-Green Proof

**RED case:** With `do-work/working/REQ-999-x.md` present at `status: claimed`, no `do-work/CHECKPOINT.md`,
and `claimed_at` two minutes old, run Step 1's crash recovery as written. Today it reports a foreign
claim, offers nothing, and leaves the REQ stranded outside `do-work/queue/` forever.

Second RED, mechanical: `grep -n "this session's to recover" actions/work.md` returns line 224 while
`_dev/tests/contract-regressions.sh` exits 0.

**Why RED now:** The classification input is written only at session end, so the branch that consumes
it cannot be reached by the event it exists to handle. And the guard matches a wording the file no
longer contains, so it green-lights the wording the file does contain.

**GREEN when:** (1) The same scenario classifies as an own crash, recovers via substeps 1–3, stamps
`status_changed_at` per REQ-074, and returns the REQ to `do-work/queue/`. (2) `actions/work.md:224`'s
successor agrees with `:113`. (3) The broadened assertion fails when either fingerprint is
reintroduced — prove it by reverting the sentence and watching the suite name the file. (4) A
two-REQ in-progress list survives a crash with both REQs recovered, not one.

**Validation:** Inferred during an adversarial audit; remediation plan reviewed and approved by the
user before capture.

## Full Context

See `do-work/user-requests/UR-015/input.md` for the audit's provenance and the findings it cleared.

---

## Addendum (2026-08-03)

An **external audit**, triaged separately via `do-work validate-feedback` and captured as UR-016,
reached finding F3 independently: "Step 1 says foreign claims must remain byte-identical, but Step 2
still says every file in `working/` belongs to this session and will be reclaimed. The regression
assertion only scans Step 1, so the suite remains green." The user's instruction was to fold that
audit's evidence into this REQ rather than duplicate it, so F3 stays this REQ's and the external
finding adds one thing to it.

**The addition — requirement 6 must widen the guard's *scope*, not only its pattern.** The assertion at
`_dev/tests/contract-regressions.sh:361-364` reads its haystack from a block, not the file:

```bash
work_step_one_block="$(sed -n '/^### Step 1: Find Next Request/,/^### Step 2\.0/p' "$repo_root/actions/work.md")"
assert_block_not_contains "$work_step_one_block" "Every .working/. file is this session's own leftover" ...
```

`actions/work.md`'s Step 1 spans lines 114–210 (`### Step 1: Find Next Request` to
`### Step 2.0: Pre-Claim Archive Collision Check`). **The stale sentence is at `:224`, inside
`### Step 2: Claim the Request` — outside that block entirely.** So a broadened *pattern* alone still
returns green: the assertion never looks at the line it is meant to catch. Requirement 6 must widen the
haystack to the whole file (or to a block that contains Step 2) as well as generalizing the wording, and
the requirement's own proof step — "prove it by reverting the sentence and watching the suite name the
file" — is what will surface this if it is missed.

**Also worth recording:** there is a *second* REQ-071 guard with the same shape, at
`_dev/tests/contract-regressions.sh:313-316` — `assert_file_not_contains "actions/work-reference.md"
'no other live session whose in-flight claim a recovery could disturb'`. It is scoped to
`actions/work-reference.md` by argument, so it cannot see `actions/work.md` at all. Both guards pin the
same retired premise in one file each; treat them as one set when broadening, per `CLAUDE.md` → Closed
Enumerations Go Stale, rather than fixing the one this REQ names and leaving its twin narrow.

No contradiction with anything above: the external finding is a subset of F3, and its contribution is
scope, not diagnosis.

---

## Triage

**Route: C** - Complex

**Reasoning:** Eight numbered requirements across three files, one of which changes a stated
contract (crash recovery's classification input) that two guards and four prose sites restate. The
addendum adds a guard-*scope* finding on top of the guard-*pattern* one, so the fix is a sweep, not
an edit.

**Planning:** Required

## Plan

**Approach — one premise, four prose sites, two guards.**

1. **`actions/work-reference.md` — new named procedure `## In-Progress Record (Step 2)`**, placed
   between `## Composed Exit Summary (Step 1)` and `## Triage Section Template (Step 3)` so it sits
   in step order *and* stays outside the `sed`-delimited `crash_recovery_block` the suite reads
   (`/^## Crash Recovery (Step 1)/,/^## Worktree Dispatch Mode/`). Per Builder Guidance, the prose
   exceeds a couple of sentences, so it goes in the reference file and `actions/work.md` Step 2
   points at it. The procedure states: the record is a **list**, one entry per claimed REQ; create
   `do-work/CHECKPOINT.md` with only the heading and this section when absent; append, never rewrite;
   remove the entry when the REQ leaves `working/`; deterministic path, no inherited shell variable.
   An explicit *is-not* paragraph keeps it from growing into a lock/heartbeat.

2. **`actions/work-reference.md` — reconcile the two Crash Recovery sentences the fix falsifies.**
   The foreign-claim bullet (`:225`) currently says a hard crash usually leaves no checkpoint "— this
   is the ordinary path, not a corner case," which stops being true the moment Step 2 writes at claim
   time. The inputs paragraph (`:227`) says both inputs "already exist for other reasons," which also
   stops being true. Both get rewritten; `absent checkpoint is ambiguous` and the other five pinned
   phrases stay verbatim (requirement 8).

3. **`actions/work-reference.md` — Session Checkpoint Template (requirement 2).** The template is
   already a bullet list but shows one entry and never says the list is per-REQ. Add a second
   placeholder entry plus a sentence tying it to Step 2's appends and to fan-out's N concurrent
   claims.

4. **`actions/work.md` Step 2 — the claim-time write (requirement 1) and the rewrite (requirement 5).**
   New substep 3 records the claim. `:224`'s no-claim-record sentence is *replaced*, not reworded:
   after this REQ a claim record exists, so the correct successor states what it is (recovery's
   classification input) and what it is not (a lock — grants nothing, coordinates nothing, no second
   owner reads it).

5. **`actions/work.md` — the same premise's third fingerprint at Step 8 substep 6** ("no lock or
   claim record is updated") and the removal sites (requirement 4): archive move on success, archive
   move on failure, and the mid-run blocked flip's move back to `do-work/queue/`. Step 10's
   session-start note and the Orchestrator Checklist get one clause each so the checklist still
   describes the pipeline.

6. **`_dev/tests/contract-regressions.sh` — one guard set, both files, both fingerprints
   (requirement 6 + addendum).** The two REQ-071 guards are merged into a loop over
   `{actions/work.md, actions/work-reference.md}` × fingerprints, using `assert_file_not_contains`
   (whole file) rather than a `sed` block — the addendum's finding is that the old haystack could
   never see the live restatement. The comment states the *trigger condition* and marks the wordings
   illustrative, per the Closed Enumerations rule. Two positive assertions pin the claim-time write
   and the record-is-a-list rule so a later edit cannot delete the fix silently.

**Plan validation:** 8 requirements → 6 tasks, every requirement mapped (r1→4, r2→3, r3→1, r4→5,
r5→4, r6→6, r7→changelog at Step 9, r8→verified as a negative in 2 and 6). No orphan tasks. Six
tasks is above the 3-task quality line — flagged per Step 4, but they are one sweep of one premise
across its known sites, and splitting would ship a half-swept premise, which is the failure mode the
REQ was written about.

*Generated in-session (no subagent dispatch this run)*

## Exploration

**Verified line numbers** (the REQ warned they would move; these are pre-edit):

- `actions/work.md:224` — Step 2 substep 1, the "every `working/` file is this session's to recover
  … so no lock or claim record is needed" sentence. Confirmed present.
- `actions/work.md:549` — **third fingerprint, not named in the REQ**: Step 8 substep 6, "no lock or
  claim record is updated (the exclusive-session model keeps none)." Same premise, same file, missed
  by both the REQ's F3 trace and the external audit's.
- `actions/work.md:638` — Step 10's session-start note, already correct ("a `working/` REQ named
  there is this session's own to recover") but silent on who writes the record.
- `actions/work-reference.md:224-227` — the classification bullets and the "two inputs … already
  exist for other reasons" paragraph.
- `actions/work-reference.md:762` — the Session Checkpoint Template's in-progress section.
- `_dev/tests/contract-regressions.sh:313-316` — the twin guard, scoped to `actions/work-reference.md`
  by argument.
- `_dev/tests/contract-regressions.sh:360-365` — `work_step_one_block` is
  `sed -n '/^### Step 1: Find Next Request/,/^### Step 2\.0/p'`, i.e. lines 114–210. The live
  restatement is at 224. **Addendum confirmed against the tree: the assertion cannot see the line it
  exists to catch.**

**Downstream reader check.** `tools/queue-kanban/verify.go:227` (`appendCheckpointGhostFindings`)
flags any REQ id `CHECKPOINT.md` names that exists in neither `queue/`, `working/`, nor `archive/`.
A claim-time entry names a REQ that is in `working/` by construction, so the new write cannot
manufacture a ghost — and requirement 4's removal is what keeps it that way after archival.

**Assertion helpers available:** `assert_contains`, `assert_block_contains`,
`assert_block_not_contains`, `assert_file_missing`, `assert_file_not_contains` (the last is
`grep -Eiq`, i.e. case-insensitive — so one pattern catches the capital-E and lowercase
fingerprints without alternation).

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — Step 2 claim-time record + rewritten premise sentence; Step 8
  substep 6 + failure path + blocked-flip removal clauses; Step 10 session-start note; Orchestrator
  Checklist
- `actions/work-reference.md` (modify) — new `## In-Progress Record (Step 2)`; Crash Recovery bullet
  and inputs paragraph; Session Checkpoint Template
- `_dev/tests/contract-regressions.sh` (modify) — merged premise guard set (both files, both
  fingerprints, whole-file scope) + two positive assertions

**Files I will NOT touch:** `actions/forensics.md` (Check 1's manual reset stays the escape hatch and
is unaffected), `actions/cleanup.md` (Pass 0 sweeps terminal statuses only — out of scope, and the
REQ's F1 trace step 4 is an observation, not a requirement), `tools/queue-kanban/verify.go` (the
ghost check already behaves correctly under the new write). `CHANGELOG.md` and `actions/version.md`
are Step 9 lifecycle files, not implementation scope.

**Acceptance criteria (restated from REQ):**
- [ ] Step 2 writes the `## In Progress (interrupted)` record at claim time
- [ ] The record is specified as a list, and the Session Checkpoint Template accepts one
- [ ] Step 2 appends/refreshes only that section; absent-checkpoint behaviour is stated, not inferred
- [ ] The entry is removed when the REQ leaves `working/`
- [ ] `actions/work.md`'s Step 2 sentence agrees with Step 1's contract
- [ ] The guard catches both fingerprints **and** looks at the lines they live on
- [ ] The changelog discloses the 0.164.0 regression, not only the fix
- [ ] REQ-071's threshold, human-authorization, unattended, and ambiguity rules are unchanged

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified)
- `actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Moved the `## In Progress (interrupted)` record's write site from session end
(Step 10) to claim time (Step 2), which is what makes Crash Recovery's own-crash branch reachable by
the event it handles. The mechanics live in a new named procedure,
`actions/work-reference.md` → **In-Progress Record (Step 2)**, placed in step order between the
Composed Exit Summary and the Triage template — deliberately outside the `sed`-delimited
`crash_recovery_block` the suite reads, so no existing assertion's haystack changed. The procedure
specifies the record as a list (one entry per claimed REQ), append-never-rewrite, create-with-only-
that-section when `do-work/CHECKPOINT.md` is absent, removal on every departure from `working/`, and
a deterministic literal path; a dedicated paragraph states what it is *not* (no exclusivity, no
coordination, no second reader) so it cannot accrete into the lock REQ-069 deleted.

`actions/work.md` Step 2 gained substep 3 (the write, pointing at the procedure) and lost the
premise sentence in substep 1. The removal side is wired into all three departures: Step 8 substep 6
(archive on success), the failure path's move to `archive/` root, and the mid-run blocked flip's move
back to `do-work/queue/`. Step 10's session-start note and two Orchestrator Checklist lines were
reconciled.

Two prose sites in `actions/work-reference.md` that the fix falsified were rewritten: the
foreign-claim bullet's "a hard crash usually leaves none — this is the ordinary path" tail, and the
inputs paragraph's "both inputs already exist for other reasons / adds no durable state" claim. The
six phrases the suite pins from REQ-071 and REQ-074 are byte-identical. The Session Checkpoint
Template gained a second placeholder entry and a paragraph stating Step 10 enriches Step 2's list
rather than authoring it.

In the suite, the two narrow REQ-071 guards were **merged into one generalized sweep** — a nested
loop over `{actions/work.md, actions/work-reference.md}` × three fingerprints, using
`assert_file_not_contains` (whole file) instead of a `sed` block. That closes the addendum's finding:
the predecessor read only Step 1 while the live restatement sat in Step 2, so a broadened pattern
alone would still have passed. Four positive assertions pin the fix's four numbered requirements
(claim-time write, record-is-a-list, never-a-lock, removal-on-archive). Net assertion change is
+4 −0 with two narrow ones generalized into one set, not three literals appended.

## Qualification

**Passed** — 3 files verified on disk and in the diff; 8 requirements traced.

- **Files exist / show in diff:** `qualify.sh` PASS on all three (no `(new)` files, so no
  placeholder or dead-code checks apply).
- **Substantive:** every hunk is contract prose or an assertion; no whitespace-only or import-shuffle
  changes.
- **Requirements traced:** r1 → `actions/work.md` Step 2 substep 3; r2 → Session Checkpoint Template
  + the procedure's "The record is a list" bullet; r3 → the procedure's absent-checkpoint bullet
  ("only the heading and this one section", "Step 10 still owns the full checkpoint"); r4 → Step 8
  substep 6 + failure path + blocked flip; r5 → Step 2 substep 1 rewrite (verified by grep that no
  fingerprint survives); r6 → the merged guard sweep (both files, whole-file scope, illustrative
  wordings, trigger condition in the comment); r7 → Step 9 changelog (below); r8 → verified as a
  negative — all six REQ-071/REQ-074 pinned phrases still assert green and the threshold /
  human-authorization / unattended / ambiguity prose is untouched.
- **P-A-U audit:** all three boxes carry evidence, and the diff matches (`shellcheck` clean,
  `bash -n` parses, zero debug artifacts).
- **Flowing:** N/A — no data paths; the "wiring" analogue is that the new section has a live
  consumer, which the four positive assertions and the Step 2/Step 8 pointers establish.
- **Contamination check:** first REQ of the session; no previous Implementation Summary to compare.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing (exit 0, including the four new assertions)

**Red-green validation:** traced to `## Red-Green Proof`.

- *RED 1 (mechanical, as captured):* against `HEAD`'s tree,
  `grep -n "this session's to recover" actions/work.md` returns line 224 **and** the HEAD suite exits
  **0** — reproduced by running the HEAD copy of the suite against the HEAD copy of `actions/work.md`.
  ✗ before → the new suite run against that same stale tree exits 1 with 4 FAILs, two of them naming
  `actions/work.md` for the premise fingerprints ✓ after.
- *RED 2 (both fingerprints):* the single generalized pattern
  `every .working/. file is this session's` matches the REQ-069 wording
  ("Every \`working/\` file is this session's own leftover") **and** the live one
  ("every \`working/\` file is this session's to recover") — `grep -Ei` makes the capital-E variant
  free. ✗ before (the old literal matched only the first) → ✓ after.
- *RED 3 (guard scope — the addendum's finding):* the two FAILs above are raised on lines 224 and 549,
  both outside the old `work_step_one_block` range (114–210). The old assertion could not see either.
- *RED 4 (positive assertions are live):* rewording "**The record is a list.**" in the new procedure
  makes the suite fail and name `actions/work-reference.md`; restoring it returns exit 0. Deleting
  the Step 2 write makes the suite name `actions/work.md`.
- *Negative control:* after the fix, all three fingerprints grep-miss both files.

**Behavioural probe (GREEN condition 4 — a two-REQ list survives a crash with *both* recovered):**
run over a synthetic `do-work/` tree in a scratch directory, so the live queue was never touched.
Two `working/` REQs at `status: claimed` (REQ-901, REQ-902) and a checkpoint naming both in
`## In Progress (interrupted)`:

- Applying the classification rule to each id: both are named **and** present in `working/` → both
  classify as **own crash**, neither as a foreign claim. Under a singular record REQ-901 would have
  been foreign — this is the case requirement 2 exists for.
- `queue-kanban verify --repo-root <scratch>` reports no checkpoint-ghost finding for either id, so
  the claim-time write cannot manufacture the ghost that `verify.go`'s
  `appendCheckpointGhostFindings` looks for.
- Departure (requirement 4): removing REQ-901's entry alone leaves the section heading present and
  REQ-902's entry byte-identical — the fan-out case where one builder finishes before its sibling.

**Also dogfooded live:** this session applied the new procedure to its own claim — Step 10's
session-start rule deleted the previous session's `CHECKPOINT.md` after recovery found nothing to
recover, and Step 2 recreated it containing only the `# Session Checkpoint` heading and this REQ's
in-progress entry. That is the create-when-absent branch of requirement 3, executed rather than
reasoned about.

**New tests added:** 4 assertions in `_dev/tests/contract-regressions.sh` (claim-time write,
record-is-a-list, never-a-lock, removal-on-archive).

**Existing tests updated (cross-REQ impact):** the two REQ-071 premise guards
(`_dev/tests/contract-regressions.sh`, formerly at `:313-316` and `:361-364`) were merged into one
generalized sweep — intentional widening, not a relaxation: every string the old pair caught is still
caught, over a strictly larger haystack.

*Verified by work action*

## Decisions

- **D-01 — Fixed the Architecture diagram's checkpoint gloss inline rather than routing it to a
  follow-up.** DECIDE & STATE. `actions/work-reference.md:12` described `CHECKPOINT.md` as "resume
  context from previous session," which already contradicted REQ-071's contract and contradicts this
  REQ's more sharply. The Restatement Sweep's default is to route an out-of-scope stale restatement
  to a follow-up, but this one sits **inside a declared scope file**, in the same premise class the
  REQ exists to sweep — leaving it would repeat verbatim the failure the REQ was written about
  (REQ-071 fixed the paragraph and left the gloss 110 lines away). One line, zero behavioural reach,
  trivially reversible.

- **D-02 — Stated the removal rule as a trigger condition and marked the mover list illustrative,
  rather than enumerating the three pipeline movers requirement 4 names.** DECIDE & STATE. The first
  draft enumerated Step 8-success, Step 8-failure, and the blocked flip — and the sweep immediately
  found two more movers outside that list (`actions/cleanup.md` Pass 0 step 5,
  `actions/forensics.md` Check 1). A closed list in a rule that means "whenever X happens" is the
  exact defect class this repo already has a standing rule against, so the enumeration was replaced
  with the condition plus an explicitly-illustrative list. Reversible; strictly more correct.

- **D-03 — One entry per REQ id, refreshed in place on a re-claim, rather than a pure append.**
  DECIDE & STATE. Requirement 3 said "appends or refreshes … specify which" and append alone is the
  simpler answer, but a REQ *can* be claimed twice in one session: the mid-run blocked flip returns
  it to `do-work/queue/`, and Step 1's `blocked_check` re-probe can unblock and re-claim it in the
  same run. With a pure append, a single missed removal leaves a duplicate that survives the next
  removal and becomes a permanent phantom claim — a stranded-claim report that no departure can ever
  clear. Keyed-by-id refresh costs one clause and makes the record self-healing. Reversible.

## Lessons Learned

**What worked:**
- **Running the old suite against the old tree** was the cheapest possible proof of the regression:
  `OLD-SUITE-EXIT=0` with the stale sentence sitting at line 224 is one line of evidence that no
  amount of reading the assertion would have produced as convincingly.
- **Placing the new reference section outside the `sed`-delimited block an existing assertion reads.**
  Checked before writing, not after: `crash_recovery_block` runs
  `/^## Crash Recovery (Step 1)/,/^## Worktree Dispatch Mode/`, so a section inserted between those
  two headings would have silently changed six existing assertions' haystack.

**What didn't:**
- **The first draft of the removal rule enumerated three movers and was wrong within the hour** — the
  Restatement Sweep found two more in files the REQ never mentioned. Writing a closed list *while
  fixing a REQ about closed lists going stale* is how strong the pull toward enumeration is.
- **A single-file `git archive` reproduction of HEAD failed twice** before working: `_dev/`,
  `CLAUDE.md`, and `.gitattributes` are all `export-ignore`d, so a tarball of HEAD has no test suite
  and the suite's own export-ignore assertions fail against it. Reproducing a suite regression needs
  `git show HEAD:<path>` into the live repo, not `git archive`.

**Worth knowing:**
- `assert_file_not_contains` is `grep -Eiq` — **case-insensitive**. One pattern therefore covers a
  capital-E and lowercase fingerprint for free, which is why the merged sweep needed three patterns
  and not six.
- The premise had **three** fingerprints in the tree, not the two the REQ diagnosed: the third
  (`actions/work.md:549`, Step 8 substep 6, "no lock or claim record is updated") was in neither the
  REQ's F3 trace nor the external audit's. Both audits found the sentence that *justified* the
  missing claim record; neither found the one that *asserted* it again 325 lines later. When
  sweeping a premise, grep the claim it makes, not the wording you found it in.
- The record now has a live machine consumer worth remembering: `tools/queue-kanban/verify.go`'s
  `appendCheckpointGhostFindings` flags any id the checkpoint names that exists in none of
  `queue/`/`working/`/`archive/`. A claim-time entry always names a `working/` REQ, so the write is
  ghost-safe by construction — but a future change that writes an id before its file exists would
  trip it.

## Orientation

Crash recovery can now actually recover a crash. The checkpoint's in-progress list — recovery's
classification input since 0.164.0 — is written when a REQ is *claimed* rather than only at session
end, so a run that dies mid-REQ leaves the record its successor needs, and the REQ returns to the
queue instead of stranding in `do-work/working/` until a human runs forensics. Lives in the work
pipeline's claim/recover path (`actions/work.md` Steps 2/8/10, `actions/work-reference.md` → Crash
Recovery and the new In-Progress Record procedure).

[MAP CHANGED] The pipeline now keeps a **claim record** where it previously kept none. It is
deliberately not a lock — no exclusivity, no coordination, no second reader, and a paragraph in the
canonical home says so — but "the skill keeps no claim record at all" is no longer a true sentence
about this system, and two contract assertions now enforce that it is never written again.

`prime_files` is empty, so no prime staleness spot-check applies.

## Review

**Overall: 96%** | 2026-08-03T21:29:26Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 1 important, 2 minor
**Acceptance:** Pass — the RED scenario now classifies as an own crash; all four GREEN conditions
demonstrated, three mechanically and one over a synthetic tree.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-086

### Findings

- **[Important] The removal rule has two more movers than the pipeline's three, in files this REQ
  never touched.** Found by the Restatement Sweep: `actions/cleanup.md` Pass 0 step 5 and
  `actions/forensics.md` Check 1 both relocate a REQ out of `do-work/working/` and say nothing about
  the in-progress entry, so each leaves a stale one behind — which makes the next run report a
  contradiction on a REQ that finished normally, training readers to ignore the one warning that
  distinguishes a real stranded claim. Partly mitigated in this REQ (D-02 restated the rule as a
  trigger condition and names both movers in the canonical procedure's illustrative list), so what
  remains is that each site's own reader never sees it. Routed to **REQ-086**.
- **[Minor] `docs/work-guide.md:66` and `:119` still describe the checkpoint as a session-end
  artifact.** User-facing, and it is the reading that made the 0.164.0 regression invisible. Also
  routed to REQ-086 rather than fixed here — a different file and a different audience, and the
  guide deserves a sentence written for it rather than a contract clause pasted in.
- **[Minor] The fix is prose, and prose has no runtime.** Nothing mechanically forces an orchestrator
  to write the record at claim time; the four new assertions pin the *instruction*, not the
  behaviour. That is inherent to a prose-driven pipeline and not a defect of this REQ — noted so a
  later reader does not mistake a green suite for a verified claim-time write. The live dogfood in
  `## Testing` is the closest thing to behavioural evidence available.

### Requirements Checklist

| # | Requirement | Status |
|---|---|---|
| 1 | Write the record at claim time (Step 2) | Delivered |
| 2 | Record is a list; template accepts one | Delivered |
| 3 | Not a second checkpoint owner; absent-file behaviour stated | Delivered (append/refresh both specified — D-03) |
| 4 | Remove on departure from `working/` | Delivered (as a trigger condition — D-02; three sites' own prose → REQ-086) |
| 5 | Step 2's premise sentence rewritten to agree with Step 1 | Delivered |
| 6 | Guard catches both fingerprints **and** sees the lines they live on | Delivered (three fingerprints, both files, whole-file scope) |
| 7 | Changelog discloses the regression, not just the fix | Delivered at Step 9 |
| 8 | REQ-071 not weakened | Delivered — all six pinned phrases byte-identical, suite green |

### Acceptance Testing

- `bash _dev/tests/contract-regressions.sh` → exit 0 (three new assertion groups, all live).
- HEAD suite vs HEAD tree → exit 0 with the stale sentence present (regression reproduced); new
  suite vs that same tree → exit 1, 4 FAILs, two naming `actions/work.md` for the premise.
- `tools/checks/qualify.sh` → OK. `tools/checks/scope-drift.sh` → OK (no drift either direction).
- `queue-kanban verify --repo-root .` → no findings.
- Synthetic two-REQ tree: both claims classify as own crash; no checkpoint-ghost finding; removing
  one entry leaves the section and its sibling intact.

### Suggested Additional Testing

- **A real interrupted run.** Kill a `do-work run` mid-REQ and re-run it, confirming the REQ returns
  to `do-work/queue/` with `status_changed_at` stamped and its generated sections stripped. Nothing
  here simulates a genuine process death.
- **Fan-out.** Two builders claimed concurrently, then a crash — the case requirement 2 exists for,
  and the one the synthetic probe only approximates. This is also REQ-085's territory.
- **A consumer repo on an older tarball.** The new Step 2 substep points at a reference section that
  only exists after this release; confirm a mid-version consumer degrades to "no such section" rather
  than stalling.

### Scores (on the record — not the headline)

Requirements 100 · Code Quality 95 · Test Adequacy 92 · Scope 100 → 96.75, reported as 96%. Test
Adequacy is held under 95 because the strongest evidence is a suite that pins instructions and a
synthetic-tree probe applying the rule by hand; there is no executable that exercises the claim-time
write end-to-end, and there cannot be one in a prose pipeline.

*Reviewed by review-work action (pipeline mode, in-session)*

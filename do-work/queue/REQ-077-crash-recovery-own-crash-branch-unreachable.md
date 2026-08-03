---
id: REQ-077
title: Crash recovery's own-crash branch is unreachable, and its retired premise survives in the same file
status: pending
created_at: 2026-08-03T16:53:42Z
user_request: UR-015
domain: general
prime_files: []
tdd: true
depends_on: []
maintenance: true
addendum_to: REQ-071
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

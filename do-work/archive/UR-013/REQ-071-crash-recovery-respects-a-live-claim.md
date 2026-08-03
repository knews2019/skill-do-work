---
id: REQ-071
title: Crash recovery must respect a live claim before stripping and re-queueing
status: completed
claimed_at: 2026-08-03T14:37:19Z
completed_at: 2026-08-03T14:48:01Z
commit: 5c39899
route: B
kb_status: pending
created_at: 2026-08-03T11:41:15Z
user_request: UR-013
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-072, REQ-073]
batch: parallel-builds
write_set: [actions/work-reference.md, actions/work.md, _dev/tests/contract-regressions.sh, docs/work-guide.md]
---

# Crash Recovery Must Respect a Live Claim Before Stripping and Re-Queueing

## What

Crash recovery is currently unconditional and destructive. `actions/work-reference.md:220-223`
tells the pipeline that every `REQ-*.md` in `do-work/working/` is its own interrupted leftover, so it
resets the frontmatter, **strips thirteen generated sections**, and moves the file back to
`do-work/queue/`. Make it respect a live claim instead: recover only what this session can show is
its own, and ask a human before taking over anything else.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** `prime_files` is empty; loaded `crew-members/general.md`, `coding-guardrails.md`, and `clear-questions.md` (this REQ specifies an interactive prompt). Approach: insert a **classification gate** ahead of `actions/work-reference.md`'s existing recovery substeps rather than rewriting them — substeps 1–3 stay byte-for-byte, they just stop being unconditional. The discriminator is the checkpoint's `## In Progress (interrupted)` record and nothing else; the age/eligibility rules and the human-authorization rationale live with the threshold in the same subsection. `actions/work.md` Step 1 gets the ordering (checkpoint first, as recovery's input) plus a one-paragraph summary that points at the reference for mechanics. Assertions go in `_dev/tests/contract-regressions.sh` next to the exclusive-session block, pinning the removed premise as absent and each replacement rule as present, plus one real ordering check.
- [x] **[APPLY]:** Assertions written first and confirmed RED (9 of 10 failing on the untouched tree), then the prose. Scope held to the four declared files; the two stale restatements found inside `actions/work.md` (Orchestrator Checklist, Verification Checklist) were fixed in place rather than filed, since both are inside the declared file and both are made false by this same change.
- [x] **[UNIFY]:** `git diff --stat` → 4 project files (`_dev/tests/contract-regressions.sh` +66, `actions/work-reference.md` +32/-1, `actions/work.md` +14/-12, `docs/work-guide.md` +2). Verified each: **contract-regressions.sh** — `bash -n` clean, `shellcheck` clean, new block sits after the three-attempt counter and uses only existing helpers plus one inline ordering check; **work-reference.md** — classification gate added above substeps 1–3, substeps themselves unchanged (diff shows only the surrounding paragraphs and the closing sentence), worktree sweep untouched and still reachable on every path; **work.md** — Step 1 ordering + summary, Step 10 checkpoint note, two checklist lines; **work-guide.md** — one paragraph in `## Checkpoints`, no other line touched. Debug-artifact scan over the diff (`console.log`/`debugger`/`TODO`/`FIXME`/`XXX`): none. Suite green.

## Why

Two independent reasons, and the first one alone justifies the change:

1. **It is wrong for the ordinary single-session crash.** Crash and restart thirty seconds later and
   recovery discards a finished `## Exploration`, a finished `## Plan`, and a declared `## Scope`.
   The pipeline only commits at Step 9, so those sections are almost always uncommitted and
   unrecoverable from git. The trail they form is what `SKILL.md:13` calls the skill's primary value.
2. It is also the single most destructive thing that happens when a second session is started against
   the same checkout — but that is a beneficiary, not the justification. **Write and review this REQ as
   "recovery is too aggressive."**

## Context

- `actions/work.md:116-118` — Step 1 opens with the exclusive-session assumption and then runs
  recovery **before** the CHECKPOINT read and before the queue glob.
- `actions/work-reference.md:220` — the rationale sentence: "every `working/` file is this session's
  own leftover to recover; there is no other live session whose in-flight claim a recovery could
  disturb, so recovery no longer consults any lock."
- `actions/work-reference.md:223` — substep 1 **removes** `claimed_at`. It must be read before it is
  discarded.
- `actions/work.md:225` — `claimed_at` is already written at claim time as a UTC ISO-8601 instant via
  `date -u +%Y-%m-%dT%H:%M:%SZ`, and already notes that a future-dated stamp "freezes the board's
  claim stopwatch."
- `actions/forensics.md:32,39` — Check 1 (Stuck Work) already reads `claimed_at` and already
  prescribes the manual reset. Reuse that definition; do not re-derive it.
- `tools/queue-kanban/future_timestamp_test.go:14-17` — existing skew precedent: 2-minute allowance,
  and "unparseable is not future."
- Nothing here trips the removed-machinery guard: `_dev/tests/contract-regressions.sh:132-137` bans
  exactly `Concurrent-Orchestrator Lock Guard`, `coexisting_sessions`, `claimed_reqs`, `heartbeat_at`,
  and `orchestrator-lock\.json`.

## Detailed Requirements

1. **Discriminate own-crash from foreign claim using the checkpoint file.** A claimed REQ named by
   `do-work/CHECKPOINT.md` is this session's own crash and recovers exactly as it does today. A
   claimed REQ *not* named there is left untouched — no frontmatter reset, no section stripping, no
   move — reported, and offered for takeover only when stale.
2. **Flip the order of the two things Step 1 does.** `actions/work.md:116-118` must read
   `do-work/CHECKPOINT.md` *before* Crash Recovery, since the checkpoint is now recovery's input.
3. **No checkpoint at all is ambiguous — ask, never strip.** Absence of a checkpoint must not be read
   as permission to recover.
4. **Three hours reports; it never authorizes.** A large REQ with review loops can legitimately exceed
   three hours, so the threshold is not a liveness test — it only bounds how long a dead claim goes
   unnoticed. The decision to take over always comes from a human. State this rationale where the
   threshold is defined, so a later edit cannot "simplify" it into an automatic takeover.
5. **Guard the timestamp toward asking.** An unparseable or future-dated `claimed_at` yields a
   negative or meaningless age; treat it as *immediately* eligible for the takeover prompt rather than
   never eligible, or a REQ becomes permanently protected. Follow the existing 2-minute skew
   allowance.
6. **Unattended runs must not block.** With no human to answer, a foreign claim is left alone and
   reported, and the run continues to the next queue item. Never stall the loop on the prompt, and
   never resolve it by stripping.
7. **Read `claimed_at` before substep 1 discards it** (`actions/work-reference.md:223`), the same
   ordering trap that already applies to the `## Scope` / `write_set` decision in that substep.
8. **Update the rationale sentence** at `actions/work-reference.md:220`. Its "there is no other live
   session whose in-flight claim a recovery could disturb" clause is exactly the premise this REQ
   removes.

## Constraints

- **No new durable state.** No lock file, no marker, no liveness counter, no new frontmatter field.
  This REQ reads `claimed_at` and `do-work/CHECKPOINT.md`, both of which already exist.
- The exclusive-session invariant at `actions/work-reference.md:53` is **REQ-073's** to reword. Do not
  touch it here; if this REQ needs the invariant to read differently, say so in `## Lessons Learned`
  rather than editing it.
- Interactive prompts go through `crew-members/clear-questions.md` — one decision per question, and
  the options must state their consequence.

## Dependencies

None. Deliberately first of the batch: smallest, independently justified, and it reduces the blast
radius of REQ-073.

## Builder Guidance

**Certainty level: Firm** on all six, but the provenance differs and that matters if one has to be
traded off. **User-stated, verbatim:** leave a claimed ticket claimed, the three-hour threshold, and
asking before takeover. **Analysis-derived and approved via the plan:** the checkpoint discriminator,
the timestamp guard, and the non-blocking unattended path. If any requirement has to give, it comes
from the second group — never the first.

Latitude on: where the threshold constant is stated and how it is worded; the exact prompt text; and
how the report line is formatted. Keep it simple — this is prose in two action files plus assertions,
not machinery.

## Red-Green Proof

**RED prompt/case:** Two probes, both runnable today.

1. Contract suite: an assertion that recovery is conditional on a checkpoint match — e.g. that
   `actions/work-reference.md` no longer claims "there is no other live session whose in-flight claim
   a recovery could disturb", and that `actions/work.md` Step 1 reads CHECKPOINT before recovery.
   Fails on the current tree.
2. Behavioural: hand-create `do-work/working/REQ-999-probe.md` with `status: claimed`, a
   `claimed_at` ten minutes in the past, and a `## Plan` section, alongside a `do-work/CHECKPOINT.md`
   naming a *different* REQ. Start the pipeline.

**Why RED now:** Probe 2 today strips `## Plan` and moves the file to `do-work/queue/`, because
`actions/work-reference.md:220` instructs recovery to treat every `working/` file as its own.

**GREEN when:** Probe 1's assertions pass. Probe 2 leaves `do-work/working/REQ-999-probe.md`
byte-identical, reports it as a foreign claim, and — being under three hours old — does not even offer
takeover. Re-stamping `claimed_at` to four hours ago offers takeover and still strips nothing until a
human says yes. Re-stamping it to a future instant also offers takeover rather than protecting the
file forever.

**Validation:** User adjusted — the user proposed the claim-respecting behaviour and the three-hour
ask ("basically when a ticket is claimed, leave it claimed, if it's older then 3h then ask if it
should be taken over"). The checkpoint discriminator, the timestamp guard, and the non-blocking
unattended path were added during capture and approved in the plan.

## Full Context

See `do-work/user-requests/UR-013/input.md` for complete verbatim input.

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" is fully specified with file-and-line anchors, but the change is prose rewriting across two coupled instruction files plus new suite assertions — the existing recovery wording, the forensics Check 1 reset definition it must reuse, and the assertion helpers all need reading before a line is written. Route A would skip that read.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The destructive procedure (the thing being narrowed):** `actions/work-reference.md:218-233` → `## Crash Recovery (Step 1)`. Its opening rationale (`:220`) asserts "every `working/` file is this session's own leftover to recover; there is no other live session whose in-flight claim a recovery could disturb" — the premise requirement 8 removes. Substep 1 (`:223`) resets `status`, **removes `claimed_at`**, and carries the existing `## Scope`-conditional `write_set` rule; substep 2 (`:224`) strips thirteen sections; substep 3 moves the file to `do-work/queue/`. The worktree sweep (`:227-231`) is name-based and independent of the per-file loop — it must stay reachable no matter how the classification lands, or a leftover branch survives every run.

**The ordering to flip:** `actions/work.md:116-118`. Step 1's opening says "go straight to the `do-work/CHECKPOINT.md` check and Crash Recovery below" — one clause naming both, with no ordering. The checkpoint read itself is specified remotely, at Step 10's **On session start (Step 1 addition)** note (`actions/work.md:625-630`), whose substep 3 already defers deletion until recovery is finished — that deferral is what makes the checkpoint readable as recovery's input, so nothing there needs changing.

**Reusable definitions found (do not re-derive):**
- `actions/forensics.md:39` — Check 1's suggested remediation is the canonical manual reset (reset `status: pending`, stamp `status_changed_at`, remove `claimed_at`/`route`, strip sections, move back to the queue). Point at it. Its 1-hour/24-hour bands are *reporting* thresholds for a read-only diagnostic, a different purpose from this REQ's takeover threshold — worth saying so, so the two aren't read as drift.
- `actions/forensics.md:158` and `tools/queue-kanban/future_timestamp_test.go:14-19` — the repo-wide skew precedent: ~2 minutes allowed, and "unparseable is not future." The board's version deliberately does **not** flag an unparseable stamp; this REQ's guard must invert that for its own purpose (requirement 5), which is a real difference to state rather than a contradiction to hide.
- `actions/work-reference.md:91` — the Timestamp rule; `claimed_at` is a UTC ISO-8601 instant, so an age is a plain subtraction.

**Checkpoint shape (the discriminator's raw material):** `actions/work-reference.md:694-728`. Frontmatter carries `last_completed: REQ-NNN`; the body has `## Completed This Session`, `## In Progress (interrupted)`, `## Still Queued`. Only `## In Progress (interrupted)` records an in-flight claim — the other three record REQs that should *not* be sitting in `working/` at all. A checkpoint is written at Step 10 (session end), so a hard crash typically leaves none: the no-checkpoint case is the common case, not the corner case, which is why requirement 3's "ask, never strip" carries most of this REQ's weight.

**Suite mechanics:** `_dev/tests/contract-regressions.sh` — helpers `assert_contains` (grep -Eq), `assert_file_not_contains` (grep -Eiq, prints hits), `assert_block_contains`/`assert_block_not_contains` over a `sed`-extracted block, plus inline `grep -c`-style counters for exactly-once rules. Ordering assertions have no helper yet — the two existing `grep -roh | wc -l` counters are the precedent for writing an inline check. The forbidden-token sweep (`:132-137`) bans `Concurrent-Orchestrator Lock Guard`, `coexisting_sessions`, `claimed_reqs`, `heartbeat_at`, `orchestrator-lock\.json`; none of this REQ's vocabulary collides.

**No prime files.** `prime_files: []` — the orchestrator's own instruction files have no prime.

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — rewrite `## Crash Recovery (Step 1)`: drop the removed premise, add the own-crash/foreign-claim classification, the takeover threshold with its non-authorizing rationale, the timestamp guard, the unattended path, and the read-`claimed_at`-first ordering note
- `actions/work.md` (modify) — Step 1: state the checkpoint read as recovery's input and order it first; replace the "every `working/` file is this session's own leftover" sentence with the claim-respecting summary
- `_dev/tests/contract-regressions.sh` (modify) — add the RED assertions: removed premise absent, classification present, checkpoint-before-recovery ordering, threshold rationale present, timestamp guard present
- `docs/work-guide.md` (modify) — one sentence in `## Checkpoints`: the checkpoint is now also what a later session's recovery reads (see `## Decisions` D-01)

**Files I will NOT touch:**
- `actions/work-reference.md:53` `## Execution Model — Exclusive Session` — REQ-073's to reword (this REQ's Constraints)
- `actions/forensics.md` — Check 1 is the reused definition, cited not edited
- `tools/queue-kanban/` — no schema field added or changed
- `actions/cleanup.md` — Pass 0/Pass 5 already handle their own cases

**Acceptance criteria (restated from REQ):**
- [ ] A claimed `working/` REQ named by `CHECKPOINT.md`'s in-progress record recovers exactly as today
- [ ] A claimed `working/` REQ not named there is left byte-identical — no reset, no stripping, no move — and reported
- [ ] No checkpoint at all does not authorize recovery
- [ ] The three-hour threshold reports and gates the *offer*; only a human authorizes takeover, and that rationale is stated where the threshold is
- [ ] Unparseable, future-dated, or absent `claimed_at` makes a REQ immediately eligible for the offer, never permanently protected; 2-minute skew allowance honored
- [ ] An unattended run leaves foreign claims alone, reports them, and continues — never stalls, never strips
- [ ] `claimed_at` is read before substep 1 discards it
- [ ] The `actions/work-reference.md:220` rationale sentence no longer claims no live claim can be disturbed

## Decisions

- **D-01 — The checkpoint discriminator is the `## In Progress (interrupted)` record only, not any mention of the id.** DECIDE & STATE. The checkpoint names REQs in four places (`last_completed`, `## Completed This Session`, `## In Progress (interrupted)`, `## Still Queued`) and requirement 1 says "named by `do-work/CHECKPOINT.md`" without picking one. Only the in-progress record documents an in-flight claim; the other three describe REQs that should not be in `working/` at all, so a match there is a contradiction. Treating a contradiction as authorization to strip is the exact failure mode this REQ removes, so the narrow reading is the safe one. Reversible — widening it later is a one-line edit.
- **D-02 — Extended Scope to `docs/work-guide.md` (one paragraph in `## Checkpoints`).** DECIDE & STATE, declared in `## Scope` before coding rather than discovered after. The guide told users the checkpoint exists "so the next session can resume cleanly"; after this change it also decides what a later session may reset, and a user who deletes it as stale bookkeeping would silently turn every own-crash recovery into a hands-off report. That is a user-visible consequence of this REQ, not adjacent polish. One paragraph, no other line touched.
- **D-03 — Fixed two stale restatements inside `actions/work.md` rather than filing them as Discovered Tasks.** DECIDE & STATE. The Orchestrator Checklist said Step 1 "recover[s] every working/ file" and the Verification Checklist required that no REQ files remain in `do-work/working/` after the loop — both are made false by this change, and the second would have a reviewer or a later run treat a deliberately-preserved foreign claim as a defect. Both live inside a file this REQ already declares, so the surgical-changes rule points at fixing them here; leaving them would ship a self-contradicting action file.
- **D-04 — Did not load `crew-members/maintenance.md`.** DECIDE & STATE. This REQ narrows an instruction in the skill's own files, which looks like maintenance-crew territory, but `maintenance: false` and `actions/work.md` Step 6 substep 5a is marker-only by design ("do not infer it from the description"). Honoring the marker over the smell is the point of the rule. The delete-before-you-add posture was not needed anyway: this change adds a gate rather than removing a capability, and the one removal it does make (the superseded premise sentence) is required by requirement 8.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified)
- `actions/work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `docs/work-guide.md` (modified)

**What was done:** Put a classification gate in front of `## Crash Recovery (Step 1)`'s existing substeps instead of rewriting them. A `working/` REQ named in the checkpoint's `## In Progress (interrupted)` record is an **own crash** and recovers exactly as before; anything else — including the common case of no checkpoint at all — is a **foreign claim** that is left byte-identical and reported. A foreign claim is offered for takeover only past a three-hour age, or immediately when `claimed_at` is unparseable, future-dated (2-minute skew allowance), or absent; the offer is a `clear-questions.md`-shaped two-option prompt, and with no human to answer the outcome is "leave it, report it, continue." The threshold's non-authorizing rationale is stated where the threshold is, so a later edit cannot collapse it into an automatic takeover. `claimed_at` is read during classification, before substep 1 discards it. The premise that licensed unconditional recovery ("no other live session whose in-flight claim a recovery could disturb") is gone from both files. `actions/work.md` Step 1 now reads `do-work/CHECKPOINT.md` first and says why — it is recovery's input — with the Step 10 session-start note, the Orchestrator Checklist, and the Verification Checklist reconciled to match. Ten new contract-suite assertions pin the removed premise as absent, each replacement rule as present, and the checkpoint-read-before-recovery ordering.

## Qualification

Passed — 4 files verified in the diff, 8 acceptance criteria traced, P-A-U confirmed against the actual diff.

- **Files exist / show in diff:** all four declared files appear in `git diff --stat`; no `(new)` files, so no dead-code or placeholder question arises.
- **Substantive:** `work-reference.md` +32/-1 is the classification gate, threshold rationale, prompt, and unattended path — not whitespace. `contract-regressions.sh` +66 is ten assertions plus one inline ordering check, all reached (suite output shows them running: RED before, green after). `work.md` +14/-12 is four distinct edits. `work-guide.md` +2 is the user-facing consequence.
- **Requirements traced:** (1) classification → `work-reference.md` own-crash/foreign-claim bullets; (2) ordering → `work.md` Step 1 opening + the suite's line-number ordering check; (3) no-checkpoint → "absent checkpoint is ambiguous, not permission"; (4) threshold rationale → the "bounds how long a dead claim goes unnoticed; it never authorizes anything" paragraph; (5) timestamp guard → the "unparseable, future-dated, or absent" bullet with the 2-minute allowance; (6) unattended → "With no human to answer … never resolve a missing answer by stripping"; (7) read-before-discard → "Read `claimed_at` while classifying, not afterwards: substep 1 removes it"; (8) premise removed → `assert_file_not_contains` on the sentence, plus the Step 1 rewrite.
- **Flowing (not hollow):** the deliverable is prescribed behavior, so the probe below executes it rather than trusting the prose — all five fixture cases resolve to a decision with the right side effects.
- **Scope drift:** none. Four files touched, four declared, `write_set` mirrored from `## Scope` before implementation.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing (suite exits 0; `record-commit-hash`/`blanked-req-scan` and `update-script-behavior` sub-probes both pass)

**Red-green validation:**
- `_dev/tests/contract-regressions.sh` REQ-071 block (10 assertions): ✗ 9 failing before implementation → ✓ after. The tenth (the checkpoint-before-recovery *ordering* check) passed on the untouched tree by accident of the old wording and is a forward lock, not a RED — the RED for requirement 2 is the companion `Crash Recovery's input` assertion, which did fail.
- Behavioural probe (REQ's Red-Green Proof, probe 2) over a synthetic `do-work/` tree in a scratch directory, so the live queue was never touched. A `working/REQ-999-probe.md` with `status: claimed`, a `## Plan` section, and a checkpoint naming a *different* REQ in `## In Progress (interrupted)`:
  - `claimed_at` 10 minutes old → report only, no offer; file byte-identical, `## Plan` intact, not moved ✓
  - `claimed_at` 4 hours old → takeover offered; with no answer, still byte-identical ✓
  - `claimed_at` 3 hours in the future → takeover offered (not protected forever) ✓
  - `claimed_at` unparseable → takeover offered ✓
  - control: checkpoint naming REQ-999 in `## In Progress (interrupted)` → own crash, substeps 1–3 run, `## Plan` stripped, `status: pending`, `claimed_at`/`route` removed, file moved to `queue/` ✓
  - control: no checkpoint at all → foreign claim, left intact ✓

**Existing tests updated (cross-REQ impact):** none — the new assertions are additive; no prior assertion changed.

**Pre-flight baseline:** clean (working tree clean outside `do-work/`, suite green before any edit). The `_dev/tests/update-script-behavior.sh` failure recorded in the previous session's checkpoint did **not** reproduce — it passes on this tree, so nothing was excluded from this REQ's gate.

## Discovered Tasks

- **[low]** `actions/forensics.md` Check 1's manual reset stamps `status_changed_at` with the current UTC instant (so the board's state timer knows when the reset happened); Crash Recovery substep 1 resets `status` to `pending`/`pending-answers` and stamps nothing. A recovered REQ therefore loses its flip instant and the board dates it from `created_at`/file-mtime. Pre-existing (predates this REQ) and not in this REQ's requirements — surfaced by the restatement sweep while checking that the two reset procedures agree.

## Review

**Overall: 94%**

**Pipeline mode — Approve (self-review).** Acceptance: **Pass**.

**Requirements Compliance: 100%** — all eight Detailed Requirements delivered and traced in `## Qualification`. Both certainty groups honored: the three user-stated items (leave a claimed ticket claimed; three-hour threshold; ask before takeover) are verbatim in the shipped prose, and the three analysis-derived items (checkpoint discriminator, timestamp guard, non-blocking unattended path) landed without trading against them. The UR's batch constraints hold: no durable coordination state (the two inputs, `CHECKPOINT.md` and `claimed_at`, both pre-exist), no new action file, no SKILL.md routing row, and `actions/work-reference.md:53`'s exclusive-session invariant untouched as REQ-073's to reword.

**Code Quality: 95%** — the classification gate is inserted above substeps 1–3 rather than rewriting them, so the diff on the destructive path is zero and the change is legible as "when do these run" instead of "what do these do." Pointer discipline held: the manual reset points at `actions/forensics.md` Check 1 rather than being restated, and the 1h/24h reporting bands are explicitly labelled a different purpose so they don't read as a second copy of the threshold. Minor: the section grew by ~330 words in a file that is already the skill's largest reference — justified here (six of the eight requirements are rules that only exist if written down), but it is the kind of growth that compounds.

**Test Adequacy: 90%** — ten assertions pin each replacement rule plus the removed premise, and one is a genuine ordering check rather than a presence grep. Red-green evidence is real: 9 of 10 failed before the prose and pass after, and the behavioural probe exercised all four timestamp cases plus both controls against a synthetic tree. Deducted for the known ceiling of grep-based assertions on prose — they pin that a *phrase* is present, not that an agent following the prose reaches the right decision. The behavioural probe is what covers that gap, and it is a one-off walkthrough rather than a committed test.

**Scope Discipline: 100%** — four files touched, four declared in `## Scope` before implementation, `write_set` mirrored from it. `tools/checks/scope-drift.sh` reports no drift in either direction. The one scope extension (`docs/work-guide.md`) was declared up front and logged as D-02, not discovered at review.

**Risk: None.** No executable behavior changed — the deliverable is prescribed orchestrator behavior. The change is strictly less destructive than what it replaces, so the failure direction is "a dead claim sits in `working/` until a human notices," reported on every run, rather than "uncommitted work is destroyed." No security, performance, or data-integrity surface.

**Restatement Sweep — run, two findings, both fixed in place.** The diff redefines two things other text restates: (a) *what recovery does to a `working/` file*, and (b) *what `do-work/CHECKPOINT.md` is for*. Swept every statement and consumer of both across `actions/`, `docs/`, `crew-members/`, `tools/`, `SKILL.md`, `next-steps.md`, and `_dev/`:
- `actions/work.md` Orchestrator Checklist — "crash recovery (recover every working/ file)" was the removed premise restated as a checklist item. Fixed.
- `actions/work.md` Verification Checklist — "No REQ files remain in `do-work/working/` after the work loop ends" became false for a deliberately-preserved foreign claim, and would have a later reviewer read correct behavior as a defect. Fixed.
- Verified still-accurate, no change needed: `actions/cleanup.md:316` (scoped to *terminal-status* REQs, which a claimed foreign claim is not), `actions/cleanup.md` Pass 5 and `docs/cleanup-guide.md` (worktree leftovers, and the sweep stays unconditional — it walks names, not the recovery loop), `actions/pipeline.md:12` (describes CHECKPOINT.md's macro-vs-micro split, not recovery), `actions/forensics.md` Check 1 (now the cited home for the manual reset; unchanged and still correct).
- One divergence found and **not** fixed: forensics' manual reset stamps `status_changed_at` and Crash Recovery substep 1 does not. Pre-existing, outside this REQ's requirements, filed as `## Discovered Tasks` `[low]` rather than swept in.

**Coding-guardrails spot-check:** Think Before Coding — four `## Decisions` entries, all DECIDE & STATE with reasoning; the one genuine ambiguity (which checkpoint record counts) was resolved toward the safe reading and logged as D-01. Simplicity First — no machinery added; the gate is prose plus one inline shell check reusing existing helpers. Surgical Changes — the two restatement fixes are one line each and both trace to this REQ. Goal-Driven Execution — RED confirmed before GREEN, plus an executed behavioural probe rather than "the suite passes."

**Findings**
- **Minor:** the takeover prompt is illustrative prose, not a template with a defined shape, so two agents may word it differently. Acceptable — the REQ granted latitude on exact prompt text, and `crew-members/clear-questions.md` is cited as the governing contract.
- **Nit:** `actions/work-reference.md`'s Crash Recovery section is now the longest subsection between the schema and Worktree Dispatch Mode. If it grows again, the threshold/prompt mechanics are the natural split point.

**Suggested additional testing**
- The first real crash after this ships is the acceptance test that matters: confirm the run reports the claimed REQ instead of silently re-queueing it, and that the report is legible enough to act on without opening the reference file.
- A run where the checkpoint *does* name the in-progress REQ (a graceful stop, then resume) should still recover automatically — the probe covered it synthetically, not through the live pipeline.

**Follow-up REQs created:** none. No Important findings; the one divergence found is a `[low]` Discovered Task.

## Lessons Learned

**What worked:**
- **Gating the destructive procedure instead of rewriting it.** Substeps 1–3 are unchanged byte-for-byte; only their precondition moved. That kept the diff on the dangerous path at zero, made the review one question ("when do these run?") instead of two, and means anything that already trusted the substeps still can.
- **Writing the assertions before the prose, on a prose-only REQ.** Committing to the exact phrases first (`absent checkpoint is ambiguous`, `never authorizes`, `substep 1 removes it`) forced each requirement to become one findable sentence rather than an idea diffused across a paragraph — the assertion is what stops requirement 4's rationale from being "simplified" away later, which is the whole point of that requirement.

**What didn't:**
- **The first instinct — "recover it if this session claimed it" — has nothing to read.** The skill keeps no session identity anywhere (that was the point of REQ-069), so ownership is unknowable in principle. The checkpoint is not a session id; it is the only durable record that *some* session was mid-REQ. Accepting that the discriminator is weaker than ownership is what made the design work: it can only be trusted in one direction, so the no-match side must be non-destructive.
- **Claiming the manual reset was "defined once" in forensics.** It isn't — substeps 1–3 define the automatic version right below the sentence making the claim. Caught in the restatement sweep and reworded to a plain pointer. Asserting single-sourcing while sitting next to a second copy is worse than not claiming it.

**Worth knowing:**
- **A hard crash usually leaves no checkpoint at all** — Step 10 writes it at session end. So the common real-world case is the *foreign claim* branch, not the own-crash branch, and the practical effect of this REQ is that a crashed REQ now waits for a human instead of being auto-recovered. That is the intended trade (the REQ's `## Why` values the uncommitted trail above the convenience), but anyone measuring "how often does recovery fire" will see it drop to near zero, and that is not a bug.
- **A stale checkpoint is now load-bearing.** Deleting `do-work/CHECKPOINT.md` as tidy-up turns every subsequent own-crash recovery into a hands-off report. Step 10's substep 3 already defers deletion until recovery finishes, which is what keeps that safe — do not "simplify" that deferral either.
- **The exclusive-session invariant now reads slightly oddly next to this.** `actions/work-reference.md:53` still says the pipeline supports "one active REQ, one coder context," while recovery now spends prose on a claim it cannot account for. Not a contradiction — the invariant is a *product contract*, and this is what the pipeline does when reality violates it — but REQ-073 is rewording that invariant, and it is worth checking there that the two paragraphs read as one position. Flagged here rather than edited, per this REQ's Constraints.

## Orientation

Now a crashed or foreign-claimed REQ keeps its plan, exploration, and scope instead of being stripped and re-queued — the work pipeline's Step 1 recovery gate (`actions/work.md` + `actions/work-reference.md`). Recovery went from unconditional to conditional on the checkpoint, with a human-authorized takeover past three hours.

`[MAP CHANGED]` — `do-work/CHECKPOINT.md` gained a second, load-bearing role: it was resume context, and is now also the input that decides which `working/` files may be recovered. Any reader that treated the checkpoint as disposable bookkeeping is now wrong.

No `prime_files` declared (`prime_files: []`) — the orchestrator's own instruction files have no prime, so no prime staleness spot-check applied and no prime-link write was deferred.

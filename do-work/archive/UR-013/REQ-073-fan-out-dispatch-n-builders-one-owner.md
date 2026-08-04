---
id: REQ-073
title: "Fan-out dispatch: N concurrent builders under one queue owner"
status: completed
claimed_at: 2026-08-03T15:09:50Z
completed_at: 2026-08-03T15:18:56Z
commit: 9ba2cda
route: C
kb_status: promoted
kb_entry: REQ-073-fan-out-dispatch-n-concurrent-builders-u.md
created_at: 2026-08-03T11:41:15Z
user_request: UR-013
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
addendum_to: REQ-069
maintenance: false
related: [REQ-071, REQ-072]
batch: parallel-builds
write_set: [actions/work-reference.md, _dev/tests/contract-regressions.sh, CLAUDE.md, docs/work-guide.md, actions/work.md]
---

# Fan-Out Dispatch: N Concurrent Builders Under One Queue Owner

## What

Raise Worktree Dispatch Mode from one builder to several concurrent builders under a single queue
owner. The exclusive-session contract is about **who owns `do-work/`**, not about how many builds run,
so a builder that owns nothing shared needs no coordination at all.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** `prime_files` is empty. Loaded `crew-members/general.md`, `coding-guardrails.md`, and `crew-members/background-agents.md` (requirement 6 adopts its pattern, so it had to be read rather than cited from memory). Approach written into `## Plan` above: reword the invariant, delete the two capping sentences, add one Fan-Out Dispatch subsection *inside* Worktree Dispatch Mode, repoint the exactly-once assertion, reconcile three downstream restatements. The measured word budget (158/200 in the Execution Model section) was established before drafting, since requirement 5's new wording has to fit inside it.
- [x] **[APPLY]:** Assertions first, confirmed RED (11 failures). Then the invariant, the cap, the section, and the three reconciliations. Two unplanned discoveries handled in place: the assertion counter aborted the whole suite silently on a no-match grep (fixed with `|| true`, since a missing invariant must read as a FAIL line, not a crash), and the restatement sweep found a family of files justifying "nothing schedules on `write_set`" with a reason this REQ falsifies — the three in declared files were fixed, the five outside routed to REQ-075.
- [x] **[UNIFY]:** `git diff --stat` → 5 project files (`_dev/tests/contract-regressions.sh` +73/-6, `actions/work-reference.md` +29/-4, `CLAUDE.md` +4/-2, `actions/work.md` +4/-2, `docs/work-guide.md` +2/-1). Verified each: **work-reference.md** — invariant reworded and re-measured at **198/200 words**; both capping phrases gone (`grep -c` = 0 for each); Fan-Out Dispatch added inside Worktree Dispatch Mode with every existing guarantee left byte-identical (sole integrator, state stays home, merge-never-rebase, four-step hand-back, the merge range and its cumulative-remediation rule, post-merge verification, both cleanup paths, Naming and the operative-name rule); **contract-regressions.sh** — `bash -n` clean, `shellcheck` clean, forbidden-token and reservation sweeps still green, SKILL.md at **2,588/2,650 words** (unchanged — no routing row added); **CLAUDE.md** — the Before-Every-Commit scope clause and the `write_set` reason; **actions/work.md** — two `write_set` reason lines; **docs/work-guide.md** — the one bullet. Debug-artifact scan over the diff: none. No new action file, no SKILL.md row, no new frontmatter field, no lock/heartbeat/registry vocabulary.

## Why

The user wants several REQs built at once without the checking overhead. Almost all of the machinery
for it already ships — `actions/work-reference.md:235` already specifies sole-integrator, "state stays
home", the four-step hand-back, the per-REQ `<pre>..<merge_hash>` range, name-based crash sweeps, and
cleanup, and the pipeline is worktree-aware throughout (`actions/work.md:377, 388, 424, 536, 572, 581,
694`). **Only the builder count is capped**, by two sentences. This REQ removes that cap and states
the boundary that makes it safe without adding any coordination.

## Context

This REQ is an addendum to **REQ-069** (`do-work/archive/UR-012/`), which adopted the exclusive-session
model at v0.161.0 and deleted ~6,500 words of orchestrator-lock, heartbeat, `claimed_reqs` and
co-dispatch re-validation machinery. That machinery all existed to police **two queue owners**, which
stays banned. This REQ reopens only the cheap half: several builders under one owner. None of the
deleted machinery returns.

- `actions/work-reference.md:237` — "The single active builder…" and "…only one builder is ever in
  flight." The two sentences to change.
- `actions/work-reference.md:53` — `## Execution Model — Exclusive Session`, asserting "one active
  session, one active REQ, one coder context."
- `_dev/tests/contract-regressions.sh:184` — asserts the exact string `one active REQ, one coder
  context` appears **exactly once** across `actions/`.
- `_dev/tests/contract-regressions.sh:241` — `router_word_budget=2650`; `SKILL.md` is at 2,588.
- `crew-members/background-agents.md:3` — loads for "any action that fans work out to background or
  parallel sub-agents", with `work (multi-REQ)` a named caller. Its run-directory pattern therefore
  becomes mandatory here, not optional.
- `actions/board.md:92` — the `overlaps` badge is advisory; **absence reads as unknown, not safe**, and
  it misses glob-vs-glob, `**`, and directory entries.
- `actions/work.md:694` — the ritual already belongs to the integrating commit, not the builder's.

## Detailed Requirements

1. **Lift the builder cap** at `actions/work-reference.md:235-237`: one integrator, several concurrent
   builders. Keep every existing guarantee — sole integrator, state stays home, merge never rebase,
   the four-step hand-back, post-merge verification, and the by-operative-name cleanup.
2. **Add a Fan-Out Dispatch subsection** covering: the human picks which REQs run together; declared
   `write_set` overlaps are surfaced as **advisory input to that pick, never a gate**; and the real
   non-interference proof is `git merge --no-ff --no-commit` refusing to merge.
3. **State the merge gate's honest limit.** Git detects conflicts by line proximity, not meaning. Two
   REQs each adding an entry to a shared registry merge cleanly and can still be wrong. The existing
   integration-seam rule is what covers that, and it only works with a single integrator.
4. **Add a Serial-Only list** naming what never parallelises: queue transitions (claim, status flips,
   archive moves), REQ id allocation, and `actions/version.md` + `CHANGELOG.md` — one changelog entry
   per REQ, written by the owner at merge time.
5. **Reword the invariant** at `actions/work-reference.md:53` from "one active REQ, one coder context"
   to a one-queue-owner formulation, and keep the two-queue-owners ban explicit so none of the deleted
   machinery has a reason to come back. Update the exactly-once assertion at
   `_dev/tests/contract-regressions.sh:184` to the new wording.
6. **Adopt the fan-out durability pattern** from `crew-members/background-agents.md` — a run directory
   created before any spawn, one input and one output file per builder, a manifest, bounded waves, and
   synthesis from files rather than from conversation:

   | Guardrail slot | Fan-out use |
   | --- | --- |
   | run directory (`:26`) | `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/` |
   | per-builder input | `REQ-NNN-brief.md` — REQ body, worktree path, branch name, never-touch list, hand-back format |
   | per-builder output (`:42`) | `REQ-NNN-handback.md` — branch, file manifest, integration seams |
   | `manifest.md` (`:48`) | REQ id → builder, `worktree-agent-REQ-NNN-<suffix>`, handback file, landed status |
   | bounded waves (`:48`) | builders per wave, sized to the harness concurrency limit |

7. **Dispatch mechanism is deliberately unspecified.** Because the owner synthesizes from files and
   never from conversation, a spawned subagent and a human-driven session are indistinguishable to it.
   Do **not** document two separate routes.
8. **Brief delivery path trap.** The brief must reach the builder as prompt content or an **absolute
   main-tree path**. A repo-relative path resolves against the worktree's own stale tracked copy of
   `do-work/`.
9. **State that integration is serial.** Implementation parallelises; merge → qualify → review →
   changelog → archive runs one REQ at a time, and each merge invalidates the previous post-merge
   verification (`actions/work.md:536`). Set the expectation rather than pretending otherwise.
10. **Reconcile `CLAUDE.md`'s "Before Every Commit" section.** It mandates a version bump plus changelog
    entry for *every* commit and auto-loads in any session — including a builder in a worktree —
    pushing builders straight into the serial-only files. `actions/work.md:694` already scopes the
    ritual to the integrating commit. **Fix the instruction, not each brief:** a rule that must be
    overridden by every brief will meet a brief that forgets.
11. **Update the user-facing statement** at `docs/work-guide.md:89` ("One REQ at a time"), which becomes
    false for builders while staying true for queue ownership.

## Constraints

- **No new durable coordination state**, and none of REQ-069's deleted machinery: no lock, no
  heartbeat, no claim registry, no liveness probe, no takeover logic. The forbidden-token sweep at
  `_dev/tests/contract-regressions.sh:132-137` must stay green.
- **No new action file and no new SKILL.md routing row.** Budget headroom is 62 words.
- **`write_set` stays display-only.** Nothing may schedule, gate, or dispatch on it. `tools/queue-kanban/`
  column logic is untouched.
- Carry over `crew-members/background-agents.md:11-14`'s own ceiling note: the durability pattern makes
  fan-out failures **survivable, not prevented**. Do not describe it as a fix.
- Worktrees live outside the repo working tree (`actions/work-reference.md:235`) — a nested second
  checkout is a documented corruption path.
- **A worktree per builder is mandatory, not optional.** The original request offered "with our
  without a new workspace, is a valid variation"; the shared-tree variation was **ruled out** during
  capture and the user did not contest it. Reason: sharing one working tree means every test run,
  qualification check, and review diff in both sessions reads a tree containing the other builder's
  unfinished edits, so the evidence steps stop meaning anything and nothing downstream can tell. The
  staging race is the lesser problem. Keep this rationale in the shipped prose — without it, a future
  reader will re-offer the shared tree as a simplification.
- **"Without letting the other session know that it is running in parallel"** — the user's phrasing,
  and it is satisfied structurally rather than by suppressing information: a builder owns no queue
  state, so there is nothing for it to know. Never implement this as hiding state from a builder that
  could otherwise read it.

## Dependencies

Functionally independent of REQ-071 and REQ-072, but **run it last**: it is the largest contract change,
REQ-071 reduces its blast radius, and REQ-072's `verify` covers the worktree state it creates.

## Builder Guidance

**Certainty level: Firm** on the one-owner boundary, the serial-only list, the advisory-only status of
`write_set`, the mandatory fan-out run directory, and the deliberate silence on dispatch mechanism.

**Exploratory** on the invariant's new wording — it must survive the exactly-once assertion, stay
inside the ≤200-word Execution Model cap REQ-069 established, and read as a boundary rather than a
description. Expect to iterate on that sentence.

This is an instruction-editing REQ on the skill's own files. Prefer **rewording and deleting** over
adding: if the fan-out section grows past what the existing Worktree Dispatch Mode already says, most
of it is redundant and should be cut rather than written. That guidance governs the **shipped prose**,
not this REQ's requirement count — see `## Open Questions`.

## Open Questions

- [x] Eleven detailed requirements against a one-sentence ask carrying a lightness cue ("without the
  anxiety and checking overhead") — keep them all, or trim to the load-bearing ones? → **Keep all
  eleven.** Resolved by the user at verify time (`do-work verify-requests`, 2026-08-03). Each
  requirement traces to a specific failure surfaced during the conversation, and the builder should
  inherit that reasoning rather than rediscover it. The lightness cue was about **runtime** overhead —
  no locks, no heartbeats, no coordination checks — not about how thoroughly the request is specified.

## Red-Green Proof

**RED prompt/case:** Contract-suite assertions, which is the harness REQ-069 itself used:

1. `grep -c "The single active builder" actions/work-reference.md` is 1 today; must be 0.
2. `grep -c "only one builder is ever in flight" actions/work-reference.md` is 1 today; must be 0.
3. A grep for the Fan-Out Dispatch heading and for the Serial-Only list naming `CHANGELOG.md` returns
   nothing today; must match after.
4. The exactly-once invariant assertion at `_dev/tests/contract-regressions.sh:184` still passes
   against the **new** wording (and the old string is gone from `actions/`).
5. The forbidden-token sweep at `:132-137` and the router budget at `:241` both stay green.

**Why RED now:** `actions/work-reference.md:237` states the cap in both sentences, no fan-out section
exists, and the invariant still reads "one active REQ, one coder context."

**GREEN when:** All five assertions pass, and a live two-builder run on two non-overlapping REQs shows:
both branches merge cleanly; each REQ gets its own changelog entry with a strictly increasing version;
`do-work/working/` never holds a file the owner did not put there; `git worktree list` and
`git branch --list 'worktree-agent-*'` are empty after both archives; and the run directory is deleted.
A deliberately overlapping pair must **fail** at `git merge --no-ff --no-commit` rather than merging
silently.

**Validation:** User confirmed — the user asked for "parallel build, but without the anxiety and
checking overhead" and chose "Keep builders out of it" from an explicit option prompt, fixing the shape
as builders that never enter Step 1 and own nothing.

## Full Context

See `do-work/user-requests/UR-013/input.md` for complete verbatim input.

---

## Triage

**Route: C** - Complex

**Reasoning:** Eleven requirements that reword the skill's central operating invariant, add a subsection to its largest reference file, retire a contract-suite assertion's pinned string, and reconcile three downstream restatements — including one in an auto-loading maintainer doc. The invariant's new wording is flagged Exploratory in the REQ's own Builder Guidance and has to survive an exactly-once assertion inside a ≤200-word cap. Planning first.

**Planning:** Required

## Plan

**The core insight the REQ hands me: almost nothing needs building.** Worktree Dispatch Mode already specifies per-builder isolation, the four-step hand-back, per-REQ `<operative_name>` and `<pre>..<merge_hash>` ranges, sole-integrator, name-collision variants, post-merge verification, and name-based sweeps. All of that is already *per REQ*, so it is already N-safe. Exactly two sentences cap the count, and one invariant describes ownership in a way that conflates owner with builder. So this REQ is mostly **rewording and reconciling**, with one new subsection — and per the Builder Guidance, anything in that subsection restating what Worktree Dispatch Mode already says gets cut instead of written.

**Work, in order:**

1. **The invariant** (`actions/work-reference.md` § Execution Model — Exclusive Session). Replace "one active `do-work` session, one active REQ, one coder context" with a one-queue-owner formulation that (a) names *ownership* as the boundary, (b) says builders are not owners and any number may run, (c) keeps the two-owners ban explicit so none of REQ-069's deleted machinery has a reason to return. Constraints: the section is at **158 words of its ≤200 cap** (measured), so the rewrite has ~40 words of headroom; the new phrase has to be a distinctive greppable compound for the exactly-once assertion. Chosen token: **`one queue owner per checkout`**. Flagged Exploratory in the REQ — expect to iterate on this sentence.
2. **Lift the cap** (§ Worktree Dispatch Mode, opening paragraph). "The single active builder" → each builder; delete the trailing "even though only one builder is ever in flight" clause outright rather than qualifying it.
3. **Add § Fan-Out Dispatch** inside Worktree Dispatch Mode — not a sibling section, because every guarantee it depends on lives in that section and a sibling would invite a second copy of them. Contents, one bullet each: the human picks; `write_set` is advisory input to the pick and never a gate; the merge is the non-interference proof; the merge gate's honest limit (line proximity, not meaning) and why the integration-seam rule is what covers it; integration is serial and each merge invalidates the previous post-merge verification; a worktree per builder is mandatory with its ruled-out-alternative rationale kept in shipped prose; the mandatory run-directory mapping table; the brief-delivery path trap; and the deliberate silence on dispatch mechanism.
4. **Retire the pinned string** in `_dev/tests/contract-regressions.sh` — the exactly-once counter greps the literal old wording, so step 1 turns it red. Repoint it at the new token, keep it exactly-once, and add assertions for the fan-out section, the serial-only list naming `CHANGELOG.md`, and the absence of both capping sentences.
5. **Reconcile the three downstream restatements:** `CLAUDE.md` § Before Every Commit (scope the ritual to the integrating commit — it auto-loads in any session, a worktree builder included, and requirement 10 is explicit that the *instruction* gets fixed rather than every brief), `docs/work-guide.md`'s "One REQ at a time" bullet (false for builders, true for ownership), and a restatement sweep for anything else asserting one-builder or one-REQ-at-a-time.

**What this REQ must not do:** add an action file or a SKILL.md routing row (62 words of budget headroom); reintroduce any lock, heartbeat, claim registry, liveness probe, or takeover logic (the forbidden-token sweep must stay green); schedule anything on `write_set`; touch `tools/queue-kanban/` column logic.

*Plan validation:* all 11 requirements map — 1→step 2, 2→step 3, 3→step 3, 4→step 3, 5→steps 1+4, 6→step 3's table, 7→step 3, 8→step 3, 9→step 3, 10→step 5, 11→step 5. No orphan tasks. Five tasks, above the 3-task comfort line; they are one contract change split by file, and shipping any subset leaves the tree self-contradicting (a lifted cap with an un-reworded invariant, or a new wording with a red assertion).

## Exploration

- `actions/work-reference.md` § **Execution Model — Exclusive Session** — 158 words across three paragraphs (invariant, Current-REQ relevance, Three-attempt stop). Only the first is in scope; the other two are about coder behavior and are count-neutral.
- `actions/work-reference.md` § **Worktree Dispatch Mode** — the two capping sentences are both in the **opening paragraph**: "The single active builder runs in its own git worktree…" and "…worth having even though only one builder is ever in flight." The REQ's line references (`:235-237`) drifted by ~28 lines when REQ-071 grew the Crash Recovery section earlier in this same run; located by string, not line.
- Everything downstream in that section is **already per-REQ**: `<operative_name>` (per REQ, held from dispatch), the four-step hand-back (run per branch), `<pre>..<merge_hash>` (explicitly per REQ, cumulative across that REQ's re-merges), post-merge verification, and the two cleanup paths. Nothing in it assumes a single builder — which is why requirement 1 is a two-sentence edit rather than a rewrite.
- `_dev/tests/contract-regressions.sh` — `exclusive_invariant_count` greps the literal `one active REQ, one coder context` and requires exactly 1 across `actions/`. Step 1 makes it 0, so this assertion is the REQ's built-in RED. The neighbouring `three_attempt_count` counter is the pattern to copy for the replacement.
- `_dev/tests/contract-regressions.sh` forbidden-token sweep — bans `Concurrent-Orchestrator Lock Guard`, `coexisting_sessions`, `claimed_reqs`, `heartbeat_at`, `orchestrator-lock\.json`; and the reservation sweep bans `status: reserved`, `reserved_for`, `reserved_at`, `do-work reserve`, `actions/reserve\.md`. None of this REQ's vocabulary collides — "fan-out", "queue owner", "builder", "integration seam" are all clear.
- `crew-members/background-agents.md` — JIT_CONTEXT already names `work (multi-REQ)` as a caller, so the run-directory pattern is *already* mandated for a multi-builder work run; requirement 6 is making that concrete, not new. Its own ceiling note (`survivable, not prevented`) is at the top of the file and must be carried, not paraphrased into a promise.
- `actions/board.md:92` — the `overlaps` badge's advisory status and its miss-classes are already documented there, including "absence reads as unknown, not safe." Point at it; do not restate the glob dialect.
- `CLAUDE.md:26` § Before Every Commit — items 1 and 2 are the version bump and changelog entry. This file auto-loads in any session rooted here, including a builder's worktree, and `actions/work.md`'s Rules already scope the ritual to the integrating commit in worktree mode. So the gap is only in CLAUDE.md's own wording.
- `docs/work-guide.md:91` — "**One REQ at a time.** The loop finishes a REQ before starting the next. This pipeline runs one active REQ in one coder context…" — restates the old invariant nearly verbatim, so it is the one user-facing doc that must move with it.

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — reword the Execution Model invariant; lift the two capping sentences; add § Fan-Out Dispatch
- `_dev/tests/contract-regressions.sh` (modify) — repoint the exactly-once invariant assertion; add fan-out, serial-only, and cap-removal assertions
- `CLAUDE.md` (modify) — scope § Before Every Commit to the integrating commit
- `docs/work-guide.md` (modify) — the "One REQ at a time" bullet
- `actions/work.md` (modify) — only if the sweep finds a one-builder restatement there (its worktree-mode Rules bullet is the candidate)

**Files I will NOT touch:**
- `tools/queue-kanban/` — no column logic, no scheduling, no schema field; `write_set` stays display-only
- `SKILL.md` — no routing row, no new action
- `crew-members/background-agents.md` — the pattern is cited, not edited
- `actions/cleanup.md`, `actions/capture.md` — Pass 5 and REQ allocation already behave correctly under N builders
- `actions/board.md` — the `overlaps` advisory is already correct; cited, not restated

**Acceptance criteria (restated from REQ):**
- [ ] `grep -c "The single active builder"` and `grep -c "only one builder is ever in flight"` in `actions/work-reference.md` are both 0
- [ ] A Fan-Out Dispatch section exists, and a Serial-Only list naming `CHANGELOG.md` exists
- [ ] Every existing worktree guarantee survives: sole integrator, state stays home, merge never rebase, four-step hand-back, post-merge verification, by-operative-name cleanup
- [ ] `write_set` overlaps are advisory input to the human's pick, never a gate; nothing schedules on them
- [ ] The merge gate's honest limit is stated, with the integration-seam rule named as what covers it
- [ ] The invariant reads as one-queue-owner, keeps the two-owners ban explicit, stays inside the ≤200-word Execution Model cap, and passes the exactly-once assertion under its new wording
- [ ] The fan-out run-directory pattern is adopted from `crew-members/background-agents.md` with its ceiling note carried, not softened
- [ ] Dispatch mechanism is left unspecified — no two documented routes
- [ ] The brief-delivery path trap is stated
- [ ] Integration is stated as serial, including that each merge invalidates the previous post-merge verification
- [ ] A worktree per builder is mandatory, with the shared-tree rationale kept in shipped prose
- [ ] `CLAUDE.md` § Before Every Commit no longer pushes a worktree builder into the serial-only files
- [ ] `docs/work-guide.md`'s "One REQ at a time" is true for queue ownership and no longer false for builders
- [ ] Forbidden-token sweep green; SKILL.md router budget green; no new action file or routing row

## Decisions

- **D-01 — Fan-Out Dispatch is a subsection *inside* Worktree Dispatch Mode, not a sibling section.** DECIDE & STATE. Every guarantee fan-out relies on — sole integrator, state stays home, merge-never-rebase, the four-step hand-back, the per-REQ merge range, post-merge verification, by-operative-name cleanup — already lives in that section. A sibling section would have had to restate or cross-reference all of them, and a restatement of seven contracts is seven future drift sites. Placing it inside let the whole thing be written as "what fan-out *adds*", which is also what kept it short.
- **D-02 — The invariant's new token is `one queue owner per checkout`, and the section heading stays "Exclusive Session".** DECIDE & STATE, with the wording flagged Exploratory by the REQ. Renaming the heading would break the `^## Execution Model — Exclusive Session` assertion plus every pointer to it in `actions/work.md`, `actions/cleanup.md`, `actions/board.md`, `actions/capture-reference.md`, and two guides — a large sweep with no requirement behind it. It also isn't needed: the *session* is still exclusive; it is the *builds* that are not. The new paragraph says so outright ("That session is what 'exclusive' names here") so a reader hitting the heading after the rewrite isn't left to reconcile it alone.
- **D-03 — Deleted the "worth having even though only one builder is ever in flight" clause outright instead of qualifying it.** DECIDE & STATE, following the REQ's prefer-deleting guidance. The clause existed only to pre-empt "why isolate a single builder?" — a question fan-out answers by construction. What replaced it is a *load-bearing* sentence rather than a longer hedge: everything in the section is written per REQ and therefore already holds for N builders, which is the fact that makes requirement 1 a two-sentence edit rather than a rewrite.
- **D-04 — Fixed a latent `set -euo pipefail` bug in the assertion I was adding.** DECIDE & STATE. `grep -roh <pattern> | wc -l` aborts the entire suite when the pattern does not match: grep exits 1, `pipefail` propagates it, `set -e` kills the script — so the suite exited 1 with **no FAIL line at all**, which reads as a crash rather than as the missing invariant. Wrapped the grep in `{ … || true; }` and said why in a comment. The pre-existing `three_attempt_count` counter has the same shape and is only safe because its string currently matches; noted as a Discovered Task rather than changed, since it is outside this REQ.
- **D-05 — Did not load `crew-members/maintenance.md`.** DECIDE & STATE, same call as REQ-071's. This is a removal-and-narrowing pass on the skill's own instructions and *looks* exactly like maintenance-crew territory, but `maintenance: false` and `actions/work.md` Step 6 substep 5a is marker-only by design. The REQ's own Builder Guidance supplied the same posture ("prefer rewording and deleting over adding"), so honoring the marker cost nothing. Worth noting that this is the second REQ in this batch where the marker and the smell disagreed — see `## Discovered Tasks`.
- **D-06 — Split the restatement-sweep findings by declared scope rather than fixing all eight.** DECIDE & STATE, per `actions/review-work.md` Step 6's rule that an out-of-scope stale restatement is a finding routed to a follow-up, not the builder's scope drift. Three sites in declared files (`actions/work.md` ×2, `CLAUDE.md`) were fixed here; five outside (`actions/board.md` ×2, `docs/board-guide.md`, `tools/queue-kanban/prime-do-kanban.md`, `actions/capture-reference.md`) went to **REQ-075** with the corrected wording and an assertion to stop the sixth copy. Fixing all eight would have quietly doubled this REQ's file count on the largest contract change of the batch.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `CLAUDE.md` (modified)
- `actions/work.md` (modified)
- `docs/work-guide.md` (modified)

**What was done:** Raised worktree dispatch from one builder to N under a single queue owner, by rewording rather than building. The Execution Model invariant became **one queue owner per checkout** — the session that claims, flips status, and archives — with builders explicitly not owners, any number allowed to build concurrently, and the two-queue-owners ban kept explicit so none of REQ-069's deleted machinery has a reason to return; the paragraph came in at 198 of its 200-word cap. Both capping sentences in Worktree Dispatch Mode's opening are gone, replaced by the fact that makes the lift safe: everything in that section is already written per REQ (one `<operative_name>`, one hand-back, one `<pre>..<merge_hash>`, one cleanup — each per REQ), so it already holds for any builder count. A new **Fan-Out Dispatch** subsection adds only what fan-out genuinely introduces: a human picks the set; `write_set` is advisory input to that pick and never a gate; the merge is the non-interference proof, with its honest limit stated (git detects conflicts by line proximity, not meaning, so two REQs appending to a shared registry merge cleanly and can still be jointly wrong — the integration-seam rule covers that, and only under one integrator); integration is serial and each merge invalidates the previous post-merge verification; a worktree per builder is mandatory, with the ruled-out shared-tree rationale kept in shipped prose; a serial-only list naming queue transitions, REQ id allocation, and `actions/version.md` + `CHANGELOG.md`; the mandatory run-directory mapping from `crew-members/background-agents.md` with its survivable-not-prevented ceiling carried; the brief-delivery path trap; and deliberate silence on dispatch mechanism. The contract suite's exactly-once invariant assertion was repointed at the new token and hardened against a silent-abort bug, plus eight new assertions. `CLAUDE.md`'s Before-Every-Commit ritual is now scoped to the integrating commit so a worktree builder skips it, and `docs/work-guide.md`'s "One REQ at a time" bullet is true for ownership and no longer false for builders.

## Qualification

Passed — 5 files verified in the diff, 14 acceptance criteria traced, P-A-U confirmed against the actual diff.

- **Files exist / show in diff:** all 5 present, no `(new)` files. Mechanical checks: both retired cap phrases `grep -c` to 0; the new invariant token appears exactly once across `actions/`; the Execution Model section measures 198 words; `SKILL.md` is unchanged at 2,588 words, so no routing row crept in.
- **Substantive:** `actions/work-reference.md` +29/-4 is the invariant, the opening paragraph, and the new subsection — not reflow. `contract-regressions.sh` +73/-6 is nine assertions plus the `|| true` hardening, all of which were observed failing before the prose and passing after.
- **Requirements traced:** 1→the opening-paragraph rewrite (cap phrases gone, guarantees named as surviving); 2→Fan-Out Dispatch's first two bullets; 3→"Its limit is honest: git detects conflicts by line proximity, not meaning" plus the integration-seam pointer; 4→the Serial-only paragraph naming all three classes; 5→the invariant plus the repointed exactly-once assertion; 6→the guardrail-slot table with the ceiling note; 7→"Dispatch mechanism is deliberately unspecified"; 8→"prompt content or an absolute main-tree path"; 9→"Integration is serial", including the invalidation of the previous post-merge verification; 10→`CLAUDE.md`'s scope clause; 11→`docs/work-guide.md`'s bullet.
- **Constraints held:** no new durable coordination state (no lock, heartbeat, registry, liveness probe, or takeover logic — the forbidden-token and reservation sweeps stay green); no new action file and no SKILL.md row; `write_set` stays display-only and `tools/queue-kanban/` is untouched; the background-agents ceiling note is carried, not softened into a promise; worktrees-outside-the-repo is left as the existing section states it; the shared-tree rationale and the "nothing for a builder to know" framing are both in shipped prose rather than implemented as hiding state.
- **Flowing (not hollow):** the deliverable is prescribed behavior; every asserted rule is a findable sentence, and the nine assertions are what stop them being reworded away.
- **Scope drift:** none. Five files touched, five declared (`actions/work.md` was declared conditionally — "only if the sweep finds a one-builder restatement there" — and the sweep did). `tools/checks/scope-drift.sh` clean.
- **Contamination check (Step 10):** REQ-072 touched 14 files. Three overlap this REQ — `actions/work.md`, `CLAUDE.md`, `_dev/tests/contract-regressions.sh` — all declared in this REQ's capture-time `write_set`, all in different sections (Step 9's accelerator vs. the Rules `write_set` bullet; the Shipped-Tooling write-surface bullet vs. Before-Every-Commit's scope clause; separate assertion blocks). Expected overlap for three REQs in one UR editing one skill, not contamination.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (plus `go test ./...` in `tools/queue-kanban/` as a regression check — this REQ touches no Go, and the tool must not have moved)
**Result:** ✓ All passing (contract suite green including the `record-commit-hash`, `blanked-req-scan`, and `update-script-behavior` sub-probes; Go suite ok)

**Red-green validation:**
- 11 new/repointed contract assertions: ✗ before → ✓ after. The RED output named every one: the ownership invariant missing (found 0), the retired invariant wording still present, both cap phrases still present, and six Fan-Out Dispatch rules absent.
- The RED run also exposed a real defect in the assertion itself: the first attempt exited 1 with **no output at all**, because `grep -roh <pattern> | wc -l` under `set -euo pipefail` aborts the script when the pattern does not match. That is the exact condition the assertion exists to detect, so the assertion would have been useless in the one case it matters — fixed before the prose was written (`## Decisions` D-04). A test whose failure mode is a silent crash is worse than no test.
- Mechanical GREEN checks, run directly rather than trusted from the suite: `grep -c "The single active builder"` → 0; `grep -c "only one builder is ever in flight"` → 0; `grep -roh 'one queue owner per checkout' actions/ | wc -l` → 1; Execution Model section → 198 words (≤200); `wc -w SKILL.md` → 2,588 (≤2,650).

**Not run — the REQ's live two-builder acceptance case.** The `## Red-Green Proof`'s GREEN condition includes a live run of two builders on two non-overlapping REQs (both branches merging cleanly, one changelog entry each with strictly increasing versions, `working/` holding only what the owner put there, `git worktree list` and `git branch --list 'worktree-agent-*'` empty afterwards, the run directory deleted) and a deliberately-overlapping pair failing at `git merge --no-ff --no-commit`. That is a **harness-capability test, not a contract test**: it requires dispatching real concurrent builders into real worktrees, which this session did not do (it ran the whole batch serially, one REQ at a time, as the pipeline's default). Reported as untested rather than claimed — see `## Review` → Suggested additional testing, where it is the first item.

**Existing tests updated (cross-REQ impact):** one, intentional. The `exclusive_invariant_count` assertion (added by REQ-069) pinned the literal string `one active REQ, one coder context`. Requirement 5 replaces that wording, so the assertion was **repointed** at `one queue owner per checkout` and kept exactly-once, with a companion assertion that the retired string is gone from `actions/`, `docs/`, and `SKILL.md`. REQ-069's intent — the invariant is stated once, and a second copy is drift — is preserved exactly; only the string it pins changed.

**Pre-flight baseline:** clean (working tree clean outside `do-work/`, contract suite green before any edit).

## Discovered Tasks

- **[low]** `_dev/tests/contract-regressions.sh`'s `three_attempt_count` counter has the same `grep | wc -l` shape that silently aborted the suite under `set -euo pipefail` (see `## Decisions` D-04). It is safe only because its pattern currently matches; the day someone reworders the three-attempt rule, the suite will exit 1 with no FAIL line and read as a crash. One `{ … || true; }` wrap. Not fixed here — outside this REQ's requirements, and it is the kind of change that wants its own RED (delete the string, watch the suite go silent).
- **[low]** Two REQs in this batch (REQ-071, REQ-073) were removal/narrowing passes on the skill's own instructions with `maintenance: false`, so `crew-members/maintenance.md` never loaded (marker-only by design, and correctly honored both times). Either capture's maintenance assessment is under-firing on this REQ shape, or the marker's definition is narrower than the crew file's JIT_CONTEXT — worth one look at `actions/capture.md`'s assessment step to see which. Deliberately not a heuristic change: `actions/work.md` Step 6 5a explains why inferring the marker from the description is worse than missing it.

## Review

**Overall: 89%**

**Pipeline mode — Approve with follow-ups (self-review).** Acceptance: **Partial** — the contract change is complete and mechanically pinned, but the REQ's live two-builder GREEN condition was not run (this session built the batch serially). One follow-up REQ created.

**Requirements Compliance: 100%** — all eleven delivered and traced in `## Qualification`. The Firm items are all present and unqualified: the one-owner boundary, the serial-only list, `write_set` advisory-only, the mandatory run directory, the deliberate silence on dispatch mechanism. The Exploratory item (the invariant's wording) landed at 198/200 words, passes the exactly-once assertion under its new token, and reads as a boundary rather than a description.

**Code Quality: 92%** — the change is mostly subtraction and repointing, which is what the REQ asked for. The opening paragraph's replacement sentence carries real weight instead of hedging: "everything in this section is written per REQ and therefore already holds for any number of concurrent builders" is *why* eleven requirements cost 29 net lines in the reference file. Fan-Out Dispatch cites `actions/board.md` for the badge's miss-classes and `crew-members/background-agents.md` for the durability ceiling rather than restating either. Deducted for the section's growth in a file that is already the skill's largest reference, and because the guardrail-slot table is the one place a reader could mistake a mapping for a spec — the run-directory pattern's actual rules stay in the crew file, which the table's intro says but a skimmer might miss.

**Test Adequacy: 85%** — nine assertions pin each rule, the RED was total and legible, and the mechanical GREEN checks (both cap-phrase counts, the exactly-once count, the 198-word measurement, the SKILL.md budget) were run directly rather than inferred from a green suite. Two honest deductions. First: **the live two-builder run did not happen**, and that is the only test that would exercise what this REQ actually enables — concurrent branches merging cleanly, one changelog entry each, an empty `worktree-agent-*` list afterwards, and a deliberately-overlapping pair failing at the merge. Grep assertions prove the *prose* says the right thing; they cannot prove two builders compose. Second: the assertions are string-pinned, so a faithful rewording that preserves meaning would break them — the accepted cost of ratcheting prose, and the reason each assertion's message explains the rule rather than just naming it.

**Scope Discipline: 100%** — five files touched, five declared, `tools/checks/scope-drift.sh` clean in both directions. `actions/work.md` was declared *conditionally* ("only if the sweep finds a one-builder restatement there"), and the sweep did — so the conditional resolved rather than being stretched. The five out-of-scope sweep findings were routed to a follow-up instead of widening the diff, which is the rule rather than my preference.

**Risk: Low.** No executable behavior changed and no state was added. The real risk is *social*: this REQ makes a capability legible that was previously capped, so a user may fan out before their harness can honestly support it. Three things bound that — worktree dispatch remains "optional, advanced harnesses only" with a silent-degrade precondition; the merge is named as the only mechanical proof of non-interference, with its limit stated; and integration stays serial, so the failure mode of an over-ambitious fan-out is a conflict at merge time, which stops the run, rather than a silent bad integration. Nothing in the change makes a second *queue owner* any more likely to work — that ban is now stated more explicitly than before.

**Restatement Sweep — run, eight findings, three fixed here and five routed.** The diff redefines two things other text restates: *the execution invariant* and, downstream of it, *why `write_set` schedules nothing*.
- `_dev/tests/contract-regressions.sh` — the exactly-once assertion pinned the old invariant string verbatim. This is the sweep working as designed: the suite went red the moment the invariant changed, which is exactly what a ratchet is for. Repointed, kept exactly-once, plus a companion assertion that the retired string is gone from `actions/`, `docs/`, and `SKILL.md`.
- `docs/work-guide.md` — "One REQ at a time … runs one active REQ in one coder context" restated the old invariant nearly verbatim in the one user-facing doc for this action. Rewritten to separate ownership from build count.
- `CLAUDE.md` § Before Every Commit — requirement 10's finding, confirmed by the sweep rather than taken on faith: the file auto-loads in any session rooted here, so a builder in a worktree of this repo reads "bump the version and write a changelog entry" and walks into two serial-only files. Scoped to the integrating commit, with the reason stated (a rule every brief must override will meet a brief that forgets).
- **Five sites routed to REQ-075** (`status: pending`, Important): `actions/board.md` ×2, `docs/board-guide.md`, `tools/queue-kanban/prime-do-kanban.md`, `actions/capture-reference.md`. Each says some version of "under the exclusive-session model one REQ runs at a time, so the badge schedules nothing." The **conclusion survives; the premise does not** — and that is the dangerous shape: a reader who accepts the reasoning and learns the premise is false concludes `write_set` should now gate, which the contract forbids outright. Rated Important for that reason, not for tidiness. REQ-075 carries the corrected wording, points at the canonical home, and adds an assertion so a sixth copy cannot land quietly.
- Verified still-accurate, no change needed: `actions/work.md` Step 1, Step 2, and Step 8's "exclusive-session model" mentions (all about *ownership* — still true), `actions/cleanup.md`'s Pass 0 safety note (same), `actions/help.md`'s "one REQ at a time" run summary (accurate for the default serial loop), `actions/capture.md`'s compound-input rationalization row (about splitting REQs, unrelated).

**Coding-guardrails spot-check:** Think Before Coding — six `## Decisions`; the heading-rename question and the sweep-scope question were both surfaced with reasoning rather than resolved silently. Simplicity First — eleven requirements, 29 net lines in the reference file, zero new machinery; the "prefer deleting" guidance was followed to the point of removing a clause rather than qualifying it (D-03). Surgical Changes — every changed line traces to a requirement or to a sweep finding inside a declared file. Goal-Driven Execution — partial, and stated as such: the contract is verified, the capability is not.

**Findings**
- **Important (routed to REQ-075):** five files justify `write_set`'s display-only status with a premise this REQ falsifies. Detail above.
- **Minor:** the section heading still reads "Execution Model — Exclusive Session" while the content now permits concurrent builders. Judged accurate rather than stale — the *session* is exclusive, the builds are not, and the paragraph says so explicitly — but a renaming sweep is the honest alternative if the distinction proves confusing in practice. Deliberately not done: it would touch two guides and five actions for no requirement (D-02).
- **Minor:** the contract-suite assertions are string-pinned, so a meaning-preserving rewording of any Fan-Out Dispatch rule breaks the suite. Inherent to prose ratchets; mitigated by each message explaining the rule so whoever hits it can tell a rewording from a regression.
- **Nit:** `crew-members/background-agents.md`'s ceiling note is carried as "verbatim in spirit" rather than quoted. Intentional — an exact quote of a whole paragraph is the restatement problem this REQ spent its sweep on — but it is the one place the wording drifted from its source on purpose.

**Suggested additional testing**
- **The live two-builder run, first priority** — it is the REQ's own GREEN condition and the only test of what actually changed. Two non-overlapping REQs, two worktrees, two branches: confirm both merge cleanly, each gets its own changelog entry with a strictly increasing version, `do-work/working/` never holds a file the owner did not put there, `git worktree list` and `git branch --list 'worktree-agent-*'` are empty after both archives, and the run directory is deleted. Then the negative case: a deliberately overlapping pair must **fail** at `git merge --no-ff --no-commit` rather than merging silently.
- **The line-proximity limit, deliberately provoked** — two REQs each appending an entry to the same registry or list. They will merge cleanly; the point is to confirm the integration-seam rule and post-merge verification catch the joint wrongness that git cannot see. This is the failure mode the prose warns about, and it has never been exercised.
- **A worktree builder in this repo, reading `CLAUDE.md`** — confirm the new scope clause actually stops it from bumping the version, since that is a behavioral claim about an auto-loading instruction and grep cannot verify it.

**Follow-up REQs created:** REQ-075 (`status: pending`, `review_generated: true`, `addendum_to: REQ-073`, `maintenance: true`) — the five stale `write_set` justifications, with the corrected wording and an assertion against a sixth copy.

## Lessons Learned

**What worked:**
- **Reading the target section for what was already per-REQ before planning any addition.** Eleven requirements looked like a large build; the section turned out to already specify `<operative_name>`, the hand-back, the merge range, verification, and cleanup *per REQ*, so the "cap" was literally two sentences. That observation became the replacement sentence and is why this landed in 29 net lines. The generalizable move: before adding a concurrency story, check whether the existing prose is already written per-unit — if it is, the cap is a claim, not an architecture.
- **Measuring the word budget before drafting.** The Execution Model cap is ≤200 and the section was at 158, so the invariant rewrite had ~40 words. Knowing that first produced a boundary sentence; discovering it after would have produced a good paragraph that then had to be mangled. Two trim passes were still needed (201 → 201 → 198), which is a fair sign the cap is doing real work.

**What didn't:**
- **The first version of my own assertion was worse than no assertion.** `grep -roh <pattern> | wc -l` under `set -euo pipefail` aborts the suite when the pattern is absent — so the check for "the invariant exists" exited 1 with **no output whatsoever** in precisely the case it was written to catch. It read as a crash, not a finding. Any counter-style assertion in a `pipefail` script needs `{ … || true; }`, and the tell is that the assertion's own failure mode is silence.
- **Writing the decision down is not the same as updating the sweep's input** — the same lesson REQ-072 hit from the other side. Here the conditional Scope entry ("`actions/work.md` — only if the sweep finds a restatement") worked, because the condition was written into the declaration where `scope-drift.sh` could see the file listed. A conditional declared *only* in prose would have failed the check.
- **My first sweep pass looked for the invariant's string and stopped.** The eight-site `write_set` family only surfaced on a second pass that searched for *consumers of the reasoning* rather than copies of the wording. A sweep that greps the changed token finds restatements; a sweep that asks "what argued from this premise?" finds the dangerous ones. The five in REQ-075 share no phrase with anything I edited.

**Worth knowing:**
- **The premise, not the conclusion, is what goes stale.** Every one of the eight sweep findings kept a true conclusion ("nothing schedules on `write_set`") attached to a premise that had just become false ("one REQ runs at a time"). That shape is more dangerous than a plainly wrong sentence, because a careful reader *reasons forward* from the dead premise and lands on the opposite of the contract. When a change falsifies a premise, grep for what was justified by it, not just for restatements of it.
- **`CLAUDE.md` auto-loads in a worktree of this repo.** That is why requirement 10 insisted the *instruction* be fixed rather than each builder's brief. It also means this file's rules reach agents nobody briefed — worth remembering before adding any imperative here.
- **The one-owner/N-builders split is now the load-bearing distinction, and the heading doesn't say it.** "Execution Model — Exclusive Session" is still accurate (the session is exclusive; the builds are not), and the paragraph makes that explicit, but the mismatch is the first thing a fresh reader will trip on. If it comes up twice, rename the heading and sweep its seven pointers — the cost is known and small, and I chose not to spend it without a requirement.
- **REQ-071's takeover prompt and this REQ's ownership boundary now read as one position** — which the batch's REQ-071 flagged as needing a check. Recovery spends prose on a claim it cannot account for, and the invariant now says why that is coherent rather than contradictory: a foreign claim is evidence the *ownership* rule was violated, and the pipeline's answer is to report and ask rather than to coordinate. No edit was needed; the reworded invariant absorbed it.

## Orientation

Now several REQs can be built at once — each builder in its own git worktree on its own branch, under a single queue owner who merges, verifies, and archives them one at a time. Lives in `actions/work-reference.md`'s Worktree Dispatch Mode (new **Fan-Out Dispatch** subsection) and its Execution Model invariant; no new machinery, no coordination state, no new action.

`[MAP CHANGED]` — the skill's central operating invariant was reworded from **one active REQ, one coder context** to **one queue owner per checkout**. That is a renamed concept, not a tightened rule: the boundary is now *who owns `do-work/`* rather than *how many builds run*, and everything that reasoned from the old premise had to be re-derived (three files here, five more in REQ-075). Two queue owners on one checkout stays banned and is now stated more explicitly than before, so none of REQ-069's deleted lock/heartbeat/registry machinery gains a reason to return.

Why it matters: this is the batch's requirement the user actually asked for — "parallel build, but without the anxiety and checking overhead" — and it is delivered structurally rather than by adding checks. A builder owns nothing shared, so there is nothing for it to coordinate and nothing for the owner to poll; the only mechanical gate is a merge that either applies or refuses.

No `prime_files` declared (`prime_files: []`) — the orchestrator's own instruction files have no prime, so no prime staleness spot-check applied and no prime-link write was deferred. (`tools/queue-kanban/prime-do-kanban.md` is touched by the follow-up REQ-075, not by this one.)

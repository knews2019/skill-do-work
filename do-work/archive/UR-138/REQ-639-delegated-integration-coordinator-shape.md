---
id: REQ-639
title: '[impact-rule-change] Delegated integration: in fan-out mode the main session may hand each REQ''s integration to one agent at a time and stay a coordinator'
status: completed
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  basis:
  - Route B
  - 3-file write set
  - 5 acceptance criteria
  - cross-route regression gates
  calculated_at: 2026-10-07T19:16:34Z
created_at: 2026-10-07T19:12:33Z
user_request: UR-138
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-640, REQ-641]
batch: coordination-lessons
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/docs/work-guide.md]
claimed_at: 2026-10-07T19:14:46Z
dispatch_at: 2026-10-07T19:22:07Z
builder_handback_at: 2026-10-07T19:26:06Z
integration_at: 2026-10-07T19:27:52Z
review_at: 2026-10-07T19:32:07Z
kb_status: pending
commit: c40c6460c3b9900a2c42d912b2ab17a7169260a0
heavy_verified_at: 2026-10-07T19:38:00Z
heavy_verified_revision: c40c6460c3b9900a2c42d912b2ab17a7169260a0
completed_at: 2026-10-07T19:38:48Z
release_at: 2026-10-07T19:38:48Z
---
# Delegated Integration: The Coordinator Shape
## What
In fan-out mode the main session may hand each REQ's integration to one agent at a time, in series, and stay a coordinator that only writes pre-dispatch judgment, routes messages, and checkpoints. There is no flag: who plays the orchestrator role for a span is a placement, like Step 7's review agent ("Spawn an agent with actions/review-work.md ... Or read it and follow it in the current session"), and the dispatch mechanism stays unspecified.
## Why
The integration span is long and context-heavy, and running it in the main session blocks the conversation the user steers from. The shape was observed in the UR-137 run on 2026-10-07 (fixing accepted consumer-review findings), where the main session coordinated, builders ran in parallel worktrees, and one background integrator per REQ ran in series. The lesson also draws on Matt Maher's video "The New Way to Work With AI" (2026-10-07).
## Decision (from the source)
The coordinator writes the pre-dispatch sections (Triage, Open Questions, Plan, lessons consult, Scope, pre-flight) before each dispatch, as the fan-out contract already requires ("briefs written before any spawn"). Writing Scope afterwards from the hand-back makes scope-drift detection vacuous.
## Verified Facts (from the source)
- Plain `recover` preserves a live claim and reports RECOVERY-TAKEOVER-AVAILABLE with next_argv `advance REQ-NNN`, which continues the claim (work-reference.md Crash Recovery; recovery_commands.go:88-104).
- `recover --assume-sole-authority` and `--take-over` reset every sibling claim.
- Targeted `advance` on a working REQ classifies by section presence, and every section step in work.md is idempotent, so an agent entering late continues at the first missing section.
- `advance --checkpoint` preserves only foreign or unlabelled records, so an integrator running it could drop same-writer sibling claims.
- work.md:284 forbids passing `dispatch_at` read back from the file to record-timing-event.
## Detailed Requirements
1. `skills/do-work/actions/work.md` Step 6: add a new condition before "Spawn a general-purpose agent". If this REQ's row in the run manifest records a landed hand-back and the hand-back file exists, consume it: do not dispatch, take the dispatch instant from that manifest row for record-timing-event, stamp `builder_handback_at` only if absent, and go straight to the hand-back merge, which already proves the branch state. Never re-dispatch a builder for a landed hand-back. State it as a condition so it also serves a fresh session after a crash (background-agents.md: agents whose findings file exists are done).
2. `work.md` Step 6 dispatch paragraph: the manifest row for a REQ carries the held dispatch instant, so an integrator can record the builder-work timing event without reading `dispatch_at` back from the file.
3. `work.md` Step 10, one sentence: under delegated integration the integrator never runs `advance --checkpoint` or the loop, because the checkpoint preserves only foreign records and would drop same-writer sibling claims; the coordinator runs it, then cleanup, after the last integrator returns.
4. `skills/do-work/actions/work-reference.md`, Worktree Dispatch Mode, Fan-Out Dispatch: a new paragraph "Delegated integration — the coordinator shape" stating, in this order:
   - (a) After the wave's builders are dispatched, the orchestrator may hand each REQ's integration (hand-back merge through Step 9) to one agent at a time, in series, in the same checkout. Release and changelog stay serial-only (cite the existing Serial-only paragraph).
   - (b) The integrator enters through plain `do-work run REQ-NNN`: recover reports the live claim and `advance REQ-NNN` continues it; the existing section steps skip what is already written; Step 6's landed-hand-back condition takes it to the merge. Never run-with-recovery, `--assume-sole-authority` or `--take-over` for an integrator: those reset every sibling claim and requeue the whole wave.
   - (c) One writer under the project root at a time: while an integrator runs, the coordinator writes nothing under the project root (no REQ sections, no manifest edits, no captures); otherwise this is the "two sessions in one working tree" case the Execution Model leaves unspecified. The coordinator's turn to write is the gap between integrators; its live notes go to its own session scratch until then.
   - (d) Why: the integration span is long and context-heavy, and running it in the main session blocks the conversation the user steers from.
   - (e) The integrator brief, `REQ-NNN-integrate.md` in the run directory, written by the coordinator in the gap before that integrator starts: REQ id, run directory, hand-back path, operative name, wave membership, any mid-run addendum text, and the instruction not to run Step 10.
   - Add one row to the guardrail-slot table: `per-integrator input | REQ-NNN-integrate.md — written by the coordinator between integrators, never by a builder`.
   - Keep "Dispatch mechanism is deliberately unspecified" as is. The new paragraph must not say subagent versus session.
5. `skills/do-work/docs/work-guide.md`, "Building several REQs at once": two sentences saying the session that runs `--fan-out` can stay a coordinator by handing each REQ's integration to one agent at a time, and that the saving is still build-phase only.
## Constraints
- Prose-only changes to shipped action files plus one user-guide paragraph. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- No flag: the integrator is a placement, and the dispatch mechanism stays unspecified.
## Dependencies
None. First in the chain: REQ-640 (mid-run messages) depends on this REQ, and REQ-641 (wave-end consistency check) depends on REQ-640. The serial order was requested because all three touch `work.md`, `review-work.md` and the Fan-Out Dispatch section of `work-reference.md`.
## Builder Guidance
Certainty is high: the source names every file, the exact insertion points, the order of the new paragraph's parts, and the table row text. Latitude is limited to wording. Do not add a flag, a new action file, or a subagent-versus-session statement. REQ-640 later adds a mid-run addendum path that the integrator brief in requirement 4(e) already names ("any mid-run addendum text"); write the brief field now as stated.
## Red-Green Proof
**RED prompt/case:** Today `work.md` Step 6 unconditionally spawns a builder, so an integrator entering a REQ with a landed hand-back would build it twice.
**Why RED now:** No condition in Step 6 checks the run manifest for a landed hand-back before dispatch, and `work-reference.md` names no integrator placement.
**GREEN when:** Step 6 names the landed-hand-back condition, and `work-reference.md` names the integrator placement, the one-writer rule, the never-rwr rule, and the brief.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read plus `_dev/tests/shipped-package-reference-contract.sh` and `_dev/tests/contract-regressions.sh`.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the source names family `alternate-writer-contract-drift`, and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-138/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (ticked by the orchestrator from the builder's hand-back) Read the brief, the REQ, prime-action-files.md, the six `alternate-writer-contract-drift` lesson bullets, and the general, coding-guardrails, shared-principles and communication-style crew rules. Approach: five prose insertions at the REQ's named anchors, 4(b) worded per the coordinator note, then a restatement sweep across work.md, work-reference.md, work-guide.md and background-agents.md. Checked recovery_commands.go before wording 4(b).
- [x] **[APPLY]:** (ticked by the orchestrator from the builder's hand-back) Edits stayed inside the three declared files, committed as a3ddde35. Sweep-driven changes inside scope: the manifest.md table row, two checklist lines, and the hand-back merge step 0 stage set (D-04). The orchestrator's one-sentence timing fix 578776b2 is in work.md too (D-08).
- [x] **[UNIFY]:** (ticked by the orchestrator from the builder's hand-back) git diff --stat 0efba253..a3ddde35: 3 files, 12 insertions, 7 deletions; git diff --check clean; shipped-package-reference-contract.sh PASS 1s, contract-regressions.sh PASS 21s. Files checked by cold read: work.md Step 6 landed paragraph, dispatch paragraph, Step 10, checklist; work-reference.md new paragraph, guardrail table, step 0; work-guide.md paragraph. No debug artifacts.
*Source: "REQ A. Delegated integration: the coordinator shape." in the maintainer session's capture source, 2026-10-07.*

## Triage

**Route: B** - Medium

**Reasoning:** The REQ fixes the outcome, the three files, the insertion points and the order of the new paragraph, so no Plan agent is needed. The rule it adds has reach the REQ does not list: Step 6, Step 10, Fan-Out Dispatch, the background-agents durability pattern and the user guide restate who integrates and who runs the checkpoint, so exploration confirms the exact paragraphs and every restatement before Scope is fixed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Exploration ran in the orchestrator session (targeted greps plus full reads of work.md, the Worktree Dispatch Mode section of work-reference.md, the work-guide fan-out paragraphs and background-agents.md), against `0efba253`.

- **Required lessons consult:** the index matches two satellites, `_dev/primes/lessons-action-files.md` (5879 tokens) and `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (18680 tokens). Both are `slugged: partial`, so no targeted form is legal and both stay dropped for budget; `required_lessons` stays absent. The six `alternate-writer-contract-drift` bullets of the action-files satellite (lines 52-60, REQ-477/498/513/461/531/566) were read anyway, because they are this REQ's failure shape: sweep every restatement by ownership condition, not only the cited lines. `lessons-releases.md` does not match: this REQ changes no release rule or shipped-package reference contract.
- **Requirement 1 (landed hand-back condition):** goes in `skills/do-work/actions/work.md` Step 6 between the **Overlapping parallel writers** paragraph (line 276) and "Spawn a **general-purpose agent**" (line 278). The hand-back merge it jumps to is the paragraph at line 303. The crash-recovery anchor it cites is `crew-members/background-agents.md` lines 128-131 ("Agents whose findings file already exists are done; do not re-run them").
- **Requirement 2 (held dispatch instant in the manifest row):** the dispatch paragraph is work.md line 284; its last sentence already forbids passing `dispatch_at` read back from the file.
- **Requirement 3 (checkpoint stays the coordinator's):** Step 10 is work.md line 466 ("The checkpoint command is the sole session-end writer and preserves every live foreign or unlabelled in-progress record").
- **Requirement 4 (coordinator-shape paragraph):** Fan-Out Dispatch in `skills/do-work/actions/work-reference.md` runs lines 449-484. The **Serial-only** paragraph (line 468) is the citation target for 4(a); the natural slot for the new paragraph is after it and before **The run directory is mandatory here** (line 470). The guardrail-slot table is lines 472-478; "Dispatch mechanism is deliberately unspecified" (line 484) stays as is. 4(c) cites the Execution Model sentence at line 17 ("so do two sessions in one working tree").
- **Requirement 5:** `skills/do-work/docs/work-guide.md` "Building several REQs at once" is lines 132-134; the second paragraph already says "Integration stays serial regardless, so the time you save is in the build phase".
- **Restatements found by the sweep (the REQ does not list them):**
  - work-reference.md line 477, guardrail table row `manifest.md` ("REQ id → builder, `<operative_name>`, handback file, landed status"). Requirement 2 makes the row carry the held dispatch instant, so this row must name it too; it is inside the declared file.
  - work.md line 538 (Rules: "the orchestrator is the sole integrator") and line 297 (builder bullet) stay true when the integrator plays the orchestrator role for its span; candidates for no change, builder judges.
  - work.md Orchestrator Checklist line 487 (Step 6 "spawn agent") and line 496 (Step 10 "advance --checkpoint + cleanup") restate the two changed steps; the builder judges whether each needs the new condition named.
  - work.md line 37 ("integration stays serial", "the dispatch mechanism stays unspecified") and work-reference.md line 457 (**Integration is serial**) agree with integrators in series; no change expected.
  - work-guide.md line 138 (**What isn't specified**: "two `do-work` sessions in the same working tree") is consistent only if the new guide sentences keep the one-writer rule visible; line 76 (`advance --checkpoint` is the sole checkpoint writer) and line 165 (continue a working claim with `advance REQ-NNN`; `recover --take-over` would reset it) agree with 4(b) and requirement 3.
  - `crew-members/background-agents.md` line 181 ("The integrator attempts every completed hand-back and integrates one branch at a time") and lines 128-131 agree; no edit expected there, so it is not in Scope.
- **Test pins:** no test under `_dev/tests` pins the prose being changed. `_dev/tests/shipped-package-reference-contract.sh` resolves `file` → **Heading** citations, so every heading in work-reference.md must keep its text, and any new citation must name a real heading or a bold label the test accepts.

*Generated by work action*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — Step 6 landed-hand-back condition, the manifest row carrying the held dispatch instant, one Step 10 sentence; checklist lines only if the sweep requires them
- `skills/do-work/actions/work-reference.md` (modify) — Fan-Out Dispatch paragraph on delegated integration, one new guardrail-table row, the manifest.md row names the held dispatch instant
- `skills/do-work/docs/work-guide.md` (modify) — two sentences in Building several REQs at once

**Files I will NOT touch:** crew-members/background-agents.md (swept, consistent), review-work.md (REQ-640 and REQ-641 own it), any Go source or test, SKILL.md, any heading text in work-reference.md. Written by finalization in the main tree, not on the branch: every release path, _dev/primes/lessons-action-files.md and do-work/lessons-index.md.

**Acceptance criteria (restated from REQ):**
- [ ] work.md Step 6 has a condition before the builder spawn: a manifest row recording a landed hand-back plus an existing hand-back file means consume it, never dispatch; take the dispatch instant from the manifest row for record-timing-event; stamp builder_handback_at only if absent; go straight to the hand-back merge. It is stated as a condition, so it also serves a fresh session after a crash
- [ ] work.md Step 6 dispatch paragraph says the REQ's manifest row carries the held dispatch instant, so an integrator records builder-work timing without reading dispatch_at back from the file
- [ ] work.md Step 10 has one sentence: under delegated integration the integrator never runs advance --checkpoint or the loop; the coordinator runs it, then cleanup, after the last integrator returns
- [ ] work-reference.md Fan-Out Dispatch has a paragraph "Delegated integration — the coordinator shape" stating (a) series handoff with release and changelog serial-only, (b) entry through plain do-work run REQ-NNN and never run-with-recovery, --assume-sole-authority or --take-over, (c) one writer under the project root at a time, (d) the why, (e) the REQ-NNN-integrate.md brief and its fields, in that order
- [ ] The guardrail-slot table has the row: per-integrator input | REQ-NNN-integrate.md — written by the coordinator between integrators, never by a builder
- [ ] "Dispatch mechanism is deliberately unspecified" is unchanged; the new paragraph names no subagent-versus-session distinction and adds no flag
- [ ] work-guide.md Building several REQs at once has two sentences: the --fan-out session can stay a coordinator by handing each REQ's integration to one agent at a time, and the saving is still build-phase only
- [ ] Every work-reference.md heading is unchanged; shipped-package-reference-contract.sh and contract-regressions.sh pass
- [ ] Every restatement of a changed rule across work.md, work-reference.md, docs/work-guide.md and crew-members/background-agents.md is swept

## Pre-Flight

**Git:** ✓ Integration tip 0efba253 on `main` (this REQ's claim commit); the only dirt is this REQ's own trail (working REQ, untracked run directory `do-work/runs/work-2026-10-07-191445/`) and `do-work/working/baseline.json`, which this pre-flight rewrote
**Tests baseline:** ✓ focused contract tests green (`do-work/runs/work-2026-10-07-191445/REQ-639-probe.sh`, launched by advance: shipped-package-reference-contract.sh 1s, contract-regressions.sh 22s)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 0efba253 from the detached checkout `.git/work-run-2026-10-07-191445/drain-head` (19:18:33Z to 19:20:52Z, 137s, load under 4, no other gate running; queue-kanban-fast-tests and do-work-cli-fast-tests executed, slowest files strict_behavior_regression_test.go 20.14s and finalization_recovery_test.go 21.24s < 30s); preflight and green-gate records satisfied
**Dependencies:** ✓ prose-only change; no toolchain or module dependency

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** work.md Step 6 gains the **Landed hand-back — consume it, never re-dispatch.** condition before the builder spawn, and the dispatch paragraph now writes the held dispatch instant into the REQ's run-manifest row so a session that did not dispatch the builder records timing from that row. Step 10 says a delegated integrator stops after Step 9 and never runs `advance --checkpoint`, because the checkpoint would drop its siblings' same-writer claims; the coordinator runs it after the last integrator. work-reference.md Fan-Out Dispatch gains the **Delegated integration — the coordinator shape.** paragraph after **Serial-only** (series handoff, entry by `advance REQ-NNN` without `recover`, the reset verbs forbidden, one writer under the project root, the why, the `REQ-NNN-integrate.md` brief), the guardrail table gains a `per-integrator input` row, the `manifest.md` row names the held dispatch instant, and hand-back step 0 stages `REQ-NNN-integrate.md`. Two checklist lines and two sentences in the user guide's "Building several REQs at once" follow. After integration the orchestrator narrowed the landed-hand-back timing sentence: the builder-work event is recorded only when the hand-back has just landed, because `record-timing-event` times from the instant to now (D-08). After review the orchestrator corrected the Step 10 reason, because `advance --checkpoint` keeps every claim, and keyed the landed condition on the hand-back file alone, as `crew-members/background-agents.md` does (review F1 and F2, D-09 and D-10). Merge range 14360235..c40c6460 (builder commit a3ddde35, orchestrator commits 578776b2 and 357047d7, merges 65ad25f0, dfd3872c and c40c6460).

## Qualification

**Diff range:** 14360235..dfd3872c for the qualify record (builder commit a3ddde35, orchestrator commit 578776b2, merges 65ad25f0 and dfd3872c); cumulative range after the review fix 14360235..c40c6460.
**Gate records:** qualify satisfied; scope-drift satisfied (the three changed files equal the declared Scope).
**Warnings judged:** none raised.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. 1: the **Landed hand-back** condition sits before "Spawn a general-purpose agent" and serves a crash-recovered session too. 2: the dispatch paragraph writes the held instant into the manifest row. 3: Step 10 says a delegated integrator never runs the checkpoint or the loop, and the coordinator does. 4: the work-reference.md paragraph states (a) series handoff with release and changelog serial-only, (b) entry by the read-only `advance REQ-NNN` and the three forbidden reset verbs (builder D-01 and D-02 changed the wording, not the prohibition), (c) one writer under the project root, (d) the why, (e) the `REQ-NNN-integrate.md` brief and its fields, in that order; the guardrail row and the manifest row are present. 5: the user guide has the two sentences. "Dispatch mechanism is deliberately unspecified" is unchanged, no flag was added, and no heading changed.
**Orchestrator fix before review:** the first merge told any late session to record builder-work timing from the manifest instant. `record-timing-event` has no end flag (timing_commands.go parses no `--ended-at`) and times to now, so a crash-recovered session would charge the builder with the whole gap. 578776b2 makes the event conditional on a hand-back that has just landed (D-08).
**After review:** 357047d7 fixes review F1 to F3 in work.md only; merged with the same pre as c40c6460. The cumulative range has the same three files, all in Scope.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY cross-checked against git diff --stat 14360235..c40c6460 (3 files, all in Scope).
**Live data flow:** prose-only change. The readers are agents running `do-work run`, and every new citation names a real heading or bold label (shipped-package-reference-contract.sh PASS on the fixed branch).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at c40c6460 (detached checkout `.git/work-run-2026-10-07-191445/drain-head`)
**Result:** ✓ All passing. Exit 0 on the first run at this revision, 19:33:48Z to 19:35:59Z, gate wall 131s, load under 3 and no other gate running. Stage queue-kanban-fast-tests executed (423 tests, wall 41s, slowest file strict_behavior_regression_test.go 20.54s < 30s). Stage do-work-cli-fast-tests executed (876 tests, wall 60s, slowest file internal/finalization/finalization_recovery_test.go 20.58s < 30s). Green-gate record satisfied by advance. The two orchestrator fixes (578776b2 and 357047d7) were merged before this gate ran, so no gate ran at 65ad25f0 or dfd3872c.

**Focused tests:** `do-work/runs/work-2026-10-07-191445/REQ-639-probe.sh` (shipped-package-reference-contract.sh and contract-regressions.sh), launched by advance, exit 0; probe record satisfied. The first launch exited 1 on the orchestrator's own uncommitted lesson bullet, whose archive link cannot resolve before finalization; the bullet was set aside until Step 8 and the rerun passed. The builder ran both tests on its branch (PASS 1s and 21s), and the reviewer ran them on the integrated tree at dfd3872c (PASS).

**Red-green validation:** not applicable. This is a prose-only change to shipped action files and the user guide, with no behavior a test pins (Exploration: no test under `_dev/tests` pins the changed prose).

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 14360235..c40c6460
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md and skills/do-work/docs/work-guide.md matched subtree skills

*Verified by work action*

## Review

**Overall: 90%** | 2026-10-07T19:32:07Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 80% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `skills/do-work/actions/work.md:468` (Step 10) gives a false reason: "the checkpoint preserves only foreign or unlabelled records, so an integrator's checkpoint would drop the same-writer claims of its siblings". The real `advance --checkpoint` (`tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go`, `stripCheckpointSummaries`) keeps the whole `## In Progress (interrupted)` section byte-for-byte and drops no claim of any writer; it only rewrites `session_ended`/`queue_state` and strips the retired summary sections. The rule itself is still right (Step 10 also loops into fresh selection and would claim new work, and the checkpoint stamps `session_ended` while the coordinator's session is live), but the stated mechanism is the source's wrong "verified fact", and `work-reference.md:470` defers to it ("`actions/work.md` Step 10 says why"). A maintainer who checks the code could delete the rule as unfounded. — impact-rule-change → report only
- `skills/do-work/actions/work.md:278` (**Landed hand-back — consume it, never re-dispatch.**) requires both a manifest row that records a landed hand-back and the file on disk, yet claims to serve "a fresh session after a crash" by citing `crew-members/background-agents.md`, whose recovery step (lines 128-131) says the opposite: do not trust the manifest's per-row label, because a crashed orchestrator may never have updated it; the file's existence alone means done. If the coordinator crashes after the hand-back file lands but before it flips the row, the condition is false and Step 6 dispatches a second builder for work already on a branch, the exact case the paragraph names. Keying the condition on the hand-back file (with the manifest row as the source of the dispatch instant only) would match the cited contract. — impact-rule-change → report only

**Minor findings:**
- `skills/do-work/actions/work.md:278`: "only when the hand-back has just landed" has no threshold, and its last clause ("and it skips the event and says so instead") has an unclear subject. Integrators run in series, so every integrator after the first arrives after earlier integrations finished; by this rule the builder-work event is skipped for most delegated REQs, and the held instant that `work.md:286` and the `manifest.md` row (`work-reference.md:480`) now carry has almost no consumer. The orchestrator's narrowing (578776b2) is correct against `record-timing-event`, which has no end flag (`lifecycletiming/timing_commands.go`); the gap is who records the event instead. The coordinator could record it at landing, since the timing log lives in the git common directory, not the working tree (`lifecycletiming/lifecycle_timing.go:540-542`). — impact-rule-change → report only
- `skills/do-work/crew-members/background-agents.md:52` (and its mirrors under `skills/do-work-toolbox/` and `skills/do-work-knowledge/`) lists the manifest's columns without the held dispatch instant. The list reads as a minimum for the generic pattern, so it is not wrong, only incomplete (the builder already raised it as a discovered task). — impact-negligible → report only
- Nit: `skills/do-work/actions/work.md:468`: requirement 3 asked for one Step 10 sentence; the addition is two sentences. Content is in scope. — impact-negligible → report only

**Acceptance:** Pass — every Detailed Requirement and acceptance criterion traces to the merge range 14360235..dfd3872c; `shipped-package-reference-contract.sh` PASS (1s) and `contract-regressions.sh` PASS (21s) on the integrated tree; no `work-reference.md` heading changed; `git diff --check` clean; "Dispatch mechanism is deliberately unspecified" untouched and the new paragraph names no subagent-versus-session split and no flag. Requirement 1 is delivered with the justified timing narrowing (D-08); requirement 4(b) is delivered with the documented deviation D-01 (entry by `advance REQ-NNN`, no `recover`), which the code confirms: a mid-wave plain `recover` runs finalization discovery, and any dirty `do-work/` path not owned by a claimed working REQ (untracked briefs, hand-backs, `baseline.json`) yields `FINALIZATION-DISCOVERY-AMBIGUOUS` (`finalization/finalization_discovery.go`, `coherentClaimOnlyTopology` then `ambiguousSharedRemainder`). D-02 is confirmed too: `--take-over REQ-NNN` filters to that one claim and `--assume-sole-authority` resets every claim (`lifecycleadvance/recovery_commands.go`). Restatement sweep: redefined the Step 6 dispatch entry, the integrator entry verbs, the `manifest.md` columns, the hand-back step 0 stage set and who runs Step 10; swept `work.md` Step 1 and Rules, `work-reference.md` lines 17, 310, 430, 457 and the guardrail table, `docs/work-guide.md` lines 130, 138 and 165, `crew-members/background-agents.md` lines 52, 128-131 and 181, and `actions/run-with-recovery.md`. Stale or contradicted items are the findings above.
**Suggested testing:** 2 items
**Follow-ups created:** None (5 findings report only)

### Domain Review

General domain: no debug artifacts, no prime file touched, and every new bold-label citation (**Delegated integration — the coordinator shape**, **Landed hand-back — consume it, never re-dispatch.**, **Hand-back merge**, **Serial-only**, **Execution Model — Claim Anywhere, One Releaser**) resolves to text in the target file. P-A-U boxes are all ticked.

*Reviewed by review-work action*

**Post-review fix check (c40c6460):** Read `git diff dfd3872c..c40c6460 -- skills/` (only `skills/do-work/actions/work.md`, +2 -2). F1 resolved: Step 10 now says the integrator's loop would claim work nobody handed it and its checkpoint would stamp `session_ended` for a session the coordinator still runs. Both are true of the code (`checkpoint_commands.go` sets `session_ended`; queue-mode `advance` claims). The `work-reference.md:470` pointer ("Step 10 says why") now resolves to a correct reason. F2 resolved: the condition now depends only on the hand-back file that the manifest row names, whatever the row's status says, which matches `background-agents.md:128-131`. The row still exists in the crash case because it is written at dispatch. F3 partly resolved: the wording is now clear sentences with a clear subject. "Just landed" still has no threshold, and no rule says who records the event when the integrator skips it, so most integrators after the first in a series still skip the builder-work event. That remains a report-only gap (impact-rule-change → report only). New issues introduced: none. The test `shipped-package-reference-contract.sh` was not re-run for this two-line change, and no citation text changed.

**Orchestrator disposition:** F1 and F2 were fixed before the gate ran, as orchestrator commit 357047d7 on the builder branch, re-merged with the same pre as c40c6460 (D-09 and D-10). F3's unclear clause was rewritten in the same commit. Its remaining gap (who records builder-work timing when the integrator skips it) stays report only and is listed in Discovered Tasks. F4 and F5 stay report only.

## Decisions

Builder decisions D-01 to D-07 are recorded in full in the hand-back, `do-work/runs/work-2026-10-07-191445/REQ-639-handback.md` § Decisions. Summary:
- **D-01 (DECIDE & STATE):** 4(b) entry is the read-only `advance REQ-NNN`, not `recover`, because a mid-wave `recover` refuses on the wave's own untracked run directory (confirmed by review in finalization_discovery.go).
- **D-02 (DECIDE & STATE):** `--take-over` is described by what it resets (one claim and its sections), not "every sibling claim" as captured; the prohibition is unchanged (confirmed by review in recovery_commands.go).
- **D-03 (DECIDE & STATE):** the Step 6 and Step 10 checklist lines name the new conditions.
- **D-04 (DECIDE & STATE):** hand-back merge step 0 stages `REQ-NNN-integrate.md`, so the new brief is not left as unattributable dirt.
- **D-05 (DECIDE & STATE):** "the orchestrator is the sole integrator" and related sentences are unchanged, because an integrator plays the orchestrator role for its span.
- **D-06 (DECIDE & STATE):** the user guide's **What isn't specified** stays unchanged; the new guide sentence keeps the one-writer rule visible.
- **D-07 (DECIDE & STATE):** the manifest-row instant is conditioned on "When the run has a manifest".

Orchestrator decisions:
- **D-08 (DECIDE & STATE): builder-work timing only when the hand-back has just landed.** `record-timing-event` has no end flag and times from `--started-at` to now (timing_commands.go). As merged first, the landed-hand-back rule would have charged a late or post-crash session's whole gap to the builder. Fixed as 578776b2 before review. Reversible prose.
- **D-09 (DECIDE & STATE): Step 10 states the true reason (review F1).** The captured "verified fact" that `advance --checkpoint` keeps only foreign or unlabelled records is wrong: `stripCheckpointSummaries` keeps the whole In Progress section. The integrator still must not run Step 10, because its loop would claim work nobody handed it and its checkpoint would stamp `session_ended` for the coordinator's live session. Fixed as 357047d7.
- **D-10 (DECIDE & STATE): the landed condition is keyed on the hand-back file (review F2).** `crew-members/background-agents.md` says not to trust a manifest row's status after a crash; the file alone means done. The row still supplies the path and the dispatch instant. Fixed as 357047d7.
- **D-11 (DECIDE & STATE): the lesson goes under a new family, `restated-mechanism-unchecked`.** The builder proposed `alternate-writer-contract-drift`. That family is about writers left behind by a contract change; the costly failure here was different: two captured claims about what commands do were wrong, and one shipped through the merge.
- **D-12 (DECIDE & STATE): no builder-work timing event for this REQ.** The integrator started after the hand-back landed (19:26:06Z) and read its recipes first, so a start-to-now event would have charged the integrator's reading to the builder. This follows D-08's own rule.

## Discovered Tasks

- impact-rule-change: Under delegated integration, "only when the hand-back has just landed" leaves most integrators after the first in a series skipping the builder-work event, and no rule says who records it instead. The coordinator could record it at landing, because the timing stream lives in the Git common directory, not the working tree (review F3). → report only
- impact-negligible: The landed-hand-back condition does not say what to do when the manifest row has no dispatch instant, for example a manifest written before this release (builder). → report only
- impact-negligible: `crew-members/background-agents.md` step 3 and its two mirrors list the manifest's columns without the held dispatch instant; it reads as a minimum for the generic pattern (builder, review F4). → report only
- impact-negligible: Step 10's delegated-integration addition is two sentences where requirement 3 asked for one (review F5). → report only

## Lessons Learned

**What worked:** Consuming a landed hand-back without re-dispatch, in this REQ's own run: the integrator entered with `advance REQ-639`, the pre-dispatch sections were already written, and every later step ran from the merge. Running the review in the background right after qualify meant its two Important findings were fixed before the gate, so the gate and the heavy lane ran once.
**What didn't:** Trusting the REQ's "Verified Facts". Two of five were wrong (`--take-over` scope and what `advance --checkpoint` keeps); the builder checked the first against the code and missed the second, so a false Step 10 reason survived to review. Writing the lesson bullet into the satellite before archive also failed the focused contract test, because its link points at the archive path that only finalization creates.
**Worth knowing:** `record-timing-event` has no end flag; any rule that records a span later than its end overstates it. `advance --checkpoint` drops no claims; it rewrites `session_ended` and `queue_state` and strips the retired summary sections.

## Orientation

Now a `do-work run --fan-out` session can stay a coordinator and hand each REQ's integration to one agent at a time, and any session that finds a landed hand-back consumes it instead of re-building; this lives in the work pipeline's dispatch and integration rules (`actions/work.md` Step 6 and Step 10, `actions/work-reference.md` Fan-Out Dispatch). `_dev/primes/prime-action-files.md` is not stale: no path it names moved or was removed.

## Heavy Verification Plan

- Base revision: 14360235fae98656fbef2dd1929991a2337b7fab
- Target revision: c40c6460c3b9900a2c42d912b2ab17a7169260a0
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/work-reference.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills; skills/do-work/docs/work-guide.md matched subtree skills

## Heavy Verification Result

- Target revision: c40c6460c3b9900a2c42d912b2ab17a7169260a0
- Execution revision: c40c6460c3b9900a2c42d912b2ab17a7169260a0 (detached checkout `.git/work-run-2026-10-07-191445/drain-head`, QUEUE_KANBAN_BROWSER set)
- staged-skills: exit 0, executed (no prior evidence reused), 41s, no HEAVY-RUN-LANE-SKIPPED finding (19:36:06Z to 19:36:48Z)

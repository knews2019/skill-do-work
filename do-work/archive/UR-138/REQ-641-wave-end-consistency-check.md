---
id: REQ-641
title: '[impact-rule-change] Wave-end consistency check: the last review in a fan-out wave also sweeps for elements earlier members redefined'
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
  calculated_at: 2026-10-07T20:09:55Z
created_at: 2026-10-07T19:12:33Z
user_request: UR-138
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-639, REQ-640]
batch: coordination-lessons
depends_on: [REQ-640]
write_set: [skills/do-work/actions/review-work.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md, skills/do-work/docs/work-guide.md]
claimed_at: 2026-10-07T20:07:41Z
dispatch_at: 2026-10-07T20:15:07Z
builder_handback_at: 2026-10-07T20:19:29Z
integration_at: 2026-10-07T20:20:58Z
review_at: 2026-10-07T20:30:05Z
kb_status: pending
commit: 78037a3e350e19d871ad5890a2aaf0efdbdef776
heavy_verified_at: 2026-10-07T20:30:49Z
heavy_verified_revision: 78037a3e350e19d871ad5890a2aaf0efdbdef776
completed_at: 2026-10-07T20:31:13Z
release_at: 2026-10-07T20:31:13Z
---
# Wave-End Consistency Check
## What
Catch the one fan-out defect no per-REQ review can see: a sibling merged later restates a contract an earlier sibling redefined.
## Why
An earlier member's review ran before the later member merged, so only the last review sees both. The lesson comes from the UR-137 run on 2026-10-07 (fixing accepted consumer-review findings) and Matt Maher's video "The New Way to Work With AI" (2026-10-07).
## Detailed Requirements
1. `skills/do-work/actions/review-work.md` Step 6 Restatement Sweep:
   - Record the result in the appended `## Review` block as one line: `**Restatement sweep:** redefined <elements> | nothing redefined`.
   - Add the line to the "Append to REQ File" template, and point the existing Verification Checklist item at it.
   - Go checks section presence only, so the extra line is safe.
   - New condition in the sweep's trigger (step 1): when the REQ under review is the wave's last successful integration (wave membership from the run manifest; a set-aside member does not skip the check), the trigger set also includes every element the wave's earlier members recorded as redefined, read from their archived `## Review` lines.
   - One-sentence reason in the text: an earlier member's review ran before the later member merged, so only the last review sees both.
   - Findings route exactly as today (impact token, `→ report only` unless `impact-critical`).
2. `skills/do-work/actions/work-reference.md` Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick": one sentence pointing at the wave-end sweep as the semantic half the merge cannot provide.
## Constraints
- Prose-only changes to shipped action files. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- No new test: no existing test exercises review prose.
## Dependencies
Depends on REQ-640 (mid-run messages), which depends on REQ-639 (delegated integration, the coordinator shape). The serial order was requested because all three touch `work.md`, `review-work.md` and the Fan-Out Dispatch section of `work-reference.md`.
## Builder Guidance
Certainty is high: the source names the line format, the trigger condition, the reason sentence and the one work-reference sentence. Latitude is limited to wording. Finding routing does not change.
## Red-Green Proof
**RED prompt/case:** In a two-REQ wave, REQ-1 redefines a token and REQ-2 restates the old meaning. REQ-1's sweep runs before REQ-2 merges, and REQ-2's sweep only looks at what REQ-2 redefined, so the stale restatement ships.
**Why RED now:** The Restatement Sweep's trigger set is only the REQ under review's own redefinitions, and no review records what it redefined.
**GREEN when:** The last-of-wave review's trigger set includes REQ-1's recorded elements and finds the stale restatement.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read; no new test, since no existing test exercises review prose.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 6141 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the source names family `alternate-writer-contract-drift`, and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-138/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the brief, the REQ, prime-action-files.md, the alternate-writer-contract-drift and restated-mechanism-unchecked lesson bullets, work-reference.md delegated integration and the manifest row, and work.md Step 7. Approach: widen Restatement Sweep step 1 with a wave-end condition, widen the step 4 skip to the whole trigger set, add the record line to the Append to REQ File template, repoint the checklist, align work.md Step 7, add one sentence to the work-reference.md non-interference bullet. Routing untouched. (Ticked by the integrator from the builder hand-back.)
- [x] **[APPLY]:** 39b21c30 first build; 0b306f69 applies the coordinator answers (the orchestrator decides "last" and passes the fact and sibling ids; a missing line is named unread for the sweep; archive location cited via Folder Structure); a09d4e87 adds the work-guide.md sentence (D-09). (Ticked by the integrator from the builder hand-back.)
- [x] **[UNIFY]:** git diff --stat 2291812d..a09d4e87: review-work.md, work-reference.md, work.md, docs/work-guide.md; git diff --check clean; shipped-package-reference-contract.sh PASS and contract-regressions.sh pass on the branch; restatement grep accounting in the hand-back. (Ticked by the integrator from the builder hand-back.)
*Source: "REQ C. Wave-end consistency check." in the maintainer session's capture source, 2026-10-07.*

## Triage

**Route: B** - Medium

**Reasoning:** The REQ fixes the line format, the trigger condition, the reason sentence and the one work-reference sentence, so no Plan agent is needed. The rule it changes, the Restatement Sweep's trigger, is restated outside the two named files (work.md Step 7 states the trigger as "this REQ's diff"), and the new condition reads inputs review-work.md never names today (the run manifest, earlier members' archived reviews), so exploration confirms each anchor and every restatement before Scope is fixed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Exploration ran in the orchestrator session (targeted greps, full reads of the cited paragraphs, and a read of the Go code behind the REQ's one mechanism claim), against `2291812d`.

- **Required lessons consult:** the index matches the same two satellites as at capture. `_dev/primes/lessons-action-files.md` is now 6141 tokens and `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` 18680; both are `slugged: partial`, so no targeted form is legal and both stay dropped for budget (drop section refreshed); `required_lessons` stays absent. Read anyway: the `alternate-writer-contract-drift` bullets (lines 52-60 and REQ-640's line 62: grep the restated phrase, not only the cited line) and REQ-639's `restated-mechanism-unchecked` bullet (line 61: read the command before prose restates what it does).
- **Mechanism claim, checked:** "Go checks section presence only." True. The only Go reader of the review block is `internal/lifecycleadvance/advance_commands.go` lines 139-250 (`hasSection(sections, "Review")`, presence and order). `finalization_discovery.go` line 522 matches only `## Review Fold — REQ-NNN` headings. No Go, shell, JS or Python reader parses `Overall:`, the findings lines, or any other line inside `## Review`. An extra bold line is safe.
- **Insertion points:**
  - Requirement 1, record line: `skills/do-work/actions/review-work.md` → **Append to REQ File** (line 380), template fenced at lines 384-408. The new `**Restatement sweep:**` line belongs with the other bold record lines (lines 402-405, between Minor findings and Follow-ups created).
  - Requirement 1, checklist: **Verification Checklist** (line 484), item at line 489 ("Restatement Sweep applied … or recorded as \"nothing redefined, sweep N/A\""). Its quoted N/A wording must become the template line's `nothing redefined` form.
  - Requirement 1, trigger: Step 6 **Restatement Sweep** paragraph (heading line 130), step 1 at line 134, step 2 line 135, step 3 line 136 (routing, unchanged), step 4 line 137 ("Skip it when nothing was redefined"; must not let a last-of-wave review skip when its own diff redefines nothing but earlier members did), Origin line 139.
  - Requirement 2: `skills/do-work/actions/work-reference.md` Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick" at line 456. It already says git detects conflicts "by line proximity, not meaning" and points at the integration-seam rule; the new sentence names the wave-end sweep as the semantic half.
- **Inputs the new condition reads that review-work.md never names today:** review-work.md has no mention of a run manifest or run directory. Wave membership lives in the run manifest's member rows (`work-reference.md` line 480 manifest row, `crew-members/background-agents.md` step 3 "Write a manifest per wave") and, under delegated integration, in the integrator brief `REQ-NNN-integrate.md` ("the wave membership", line 470). Earlier members are archived before the next member integrates (Fan-Out Dispatch "Integration is serial", line 457: merge → qualify → test → review → changelog → archive one REQ at a time), so their `## Review` blocks are in archived REQ files, flat (`do-work/archive/REQ-NNN-*.md`) or under `do-work/archive/UR-NNN/` when the UR closed. A set-aside member was never archived and has no line to read.
- **Restatement found outside the REQ's write set (Scope grows):** `skills/do-work/actions/work.md` Step 7 line 368, **Restatement sweep (MUST)**, restates the trigger as "If this REQ's diff redefines something other text restates". Under the new condition a last-of-wave review sweeps even when its own diff redefines nothing, so this line goes stale. Step 7 **How to run it** (line 370) lists what the review agent receives (review-work.md, REQ path, domain file, merge range) and gives it no run manifest; the reviewer of the last member cannot apply the condition without it. Both are in work.md, which joins Scope.
- **Restatements checked and consistent (no change expected):**
  - `work-reference.md` line 462 (Auto-wave: "The merge is the non-interference proof, not the pick (above)"): about what the pick claims; it points back at line 456, which will carry the new sentence.
  - `work-reference.md` line 470 (Delegated integration): the integrator brief already carries the wave membership; the integrator passes it to its reviewer through Step 7.
  - `docs/work-guide.md` line 134 ("collisions surface when the branches merge"): user-facing and still true for textual collisions; a sentence on meaning collisions is optional (builder judges; a candidate discovered task).
  - `docs/review-work-guide.md` (Phase 2 table, Follow-ups line 55): never describes the sweep's trigger; routing is unchanged.
  - `crew-members/background-agents.md`: nothing on review or the sweep.
  - `review-work.md` line 471 (Common Rationalizations, out-of-scope stale restatement): routing unchanged, consistent.
- **Test pins:** no test parses review prose. `_dev/tests/shipped-package-reference-contract.sh` resolves `file` → **Heading** citations, so every work-reference.md heading keeps its text and any new citation names a real heading or bold label.

*Generated by work action*

## Scope

**Files I will touch:**
- `skills/do-work/actions/review-work.md` (modify) — Restatement Sweep step 1 wave-end trigger condition with its one-sentence reason, step 4 skip rule kept consistent; the record line in the Append to REQ File template; the Verification Checklist item points at that line
- `skills/do-work/actions/work-reference.md` (modify) — one sentence in the Fan-Out Dispatch bullet on the non-interference proof naming the wave-end sweep as the semantic half the merge cannot provide
- `skills/do-work/actions/work.md` (modify) — Step 7 Restatement sweep (MUST) restates the widened trigger; Step 7 How to run it passes the run manifest to the reviewer in a fan-out run

- `skills/do-work/docs/work-guide.md` (modify) — one user-facing sentence that the last review in a wave also catches meaning collisions; scope extended by the coordinator at hand-back from the builder's discovered task (recorded as a D-XX in the hand-back)

**Files I will NOT touch:** docs/work-guide.md and docs/review-work-guide.md (swept, consistent; a user-facing sentence is a discovered-task candidate), crew-members/background-agents.md, any Go source or test, SKILL.md, any heading text in work-reference.md. Written by finalization in the main tree, not on the branch: every release path, _dev/primes/lessons-action-files.md and do-work/lessons-index.md.

**Acceptance criteria (restated from REQ):**
- [ ] The appended `## Review` block records one line: `**Restatement sweep:** redefined <elements> | nothing redefined`, present in the Append to REQ File template
- [ ] The Verification Checklist's sweep item points at that line
- [ ] Restatement Sweep step 1: when the REQ under review is the wave's last successful integration (wave membership from the run manifest; a set-aside member does not skip the check), the trigger set also includes every element earlier members recorded as redefined, read from their archived `## Review` lines
- [ ] The text gives the one-sentence reason: an earlier member's review ran before the later member merged, so only the last review sees both
- [ ] Findings route exactly as today (impact token, `→ report only` unless `impact-critical`)
- [ ] work-reference.md "The non-interference proof is the merge, not the pick" bullet has one sentence pointing at the wave-end sweep as the semantic half the merge cannot provide
- [ ] Every restatement of the sweep trigger agrees (work.md Step 7); every work-reference.md heading is unchanged; shipped-package-reference-contract.sh and contract-regressions.sh pass

## Pre-Flight

**Git:** ✓ Integration tip 2291812d on `main` (this REQ's claim commit); the only dirt is this REQ's own trail (working REQ, run manifest row, untracked `do-work/runs/work-2026-10-07-191445/REQ-641-brief.md` and `REQ-641-probe.sh`) and `do-work/working/baseline.json`, which this pre-flight rewrote
**Tests baseline:** ✓ focused contract tests green (`do-work/runs/work-2026-10-07-191445/REQ-641-probe.sh`, launched by advance: shipped-package-reference-contract.sh 1s, contract-regressions.sh 21s)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 2291812d from the detached checkout `.git/work-run-2026-10-07-191445/drain-head` (20:11:17Z to 20:13:25Z, gate wall 127s, load under 4, no other gate running; queue-kanban-fast-tests and do-work-cli-fast-tests executed, slowest files strict_behavior_regression_test.go 20.45s and finalization_recovery_test.go 20.39s < 30s); preflight and green-gate records satisfied
**Dependencies:** ✓ prose-only change; no toolchain or module dependency

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** review-work.md Step 6 **Restatement Sweep** gains a wave-end condition under step 1: when the orchestrator passes that this REQ is its wave's last successful integration, with the earlier members' REQ ids, the trigger set also includes every element those members recorded on the **Restatement sweep:** line of their archived `## Review` (found via work-reference.md **Folder Structure**). It gives the reason (an earlier member's review ran before the later member merged, so only the last review sees both), and a member with no such line is named unread for the sweep, never re-derived from its diff. Step 4 skips only when the whole trigger set is empty, and the sweep's label no longer restates the old trigger. The Append to REQ File template gains the `**Restatement sweep:** redefined <elements> | nothing redefined` line, holding this diff's own elements only, and the Verification Checklist item points at it. work.md Step 7 **Restatement sweep (MUST)** names the wave-end case, and **How to run it** has the orchestrator decide "last" from the run manifest (every other member finalized or set aside) and pass the fact and the ids; a run with no manifest has no wave. work-reference.md's "The non-interference proof is the merge, not the pick" bullet gains one sentence naming the wave-end sweep as the semantic half no merge provides. work-guide.md gains one user-facing sentence (D-09, coordinator Scope extension). Finding routing is unchanged, no heading changed, no integration seams. Merge range 992f28a9..f7416a80 (builder commits 39b21c30, 0b306f69 and a09d4e87, merge f7416a80).

## Qualification

**Diff range:** 992f28a9..f7416a80 for the qualify record (builder commits 39b21c30, 0b306f69 and a09d4e87, merge f7416a80); cumulative range after the review fixes 992f28a9..78037a3e.
**Gate records:** qualify satisfied; scope-drift satisfied (the four changed files equal the declared Scope, including docs/work-guide.md, which the coordinator added at hand-back, D-09).
**Warnings judged:** none raised.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. 1: the template gains `**Restatement sweep:** redefined [...] | nothing redefined`, and the Verification Checklist item points at it; Restatement Sweep step 1 gains the wave-end condition (last successful integration, earlier members' recorded elements read from their archived `## Review`), with the reason sentence; routing (step 3) is unchanged. A set-aside member does not skip the check, because work.md Step 7 counts set-aside members as done when it decides "last". 2: the work-reference.md non-interference bullet gains one sentence naming the wave-end sweep as the semantic half. Constraint sweep: work.md Step 7's trigger restatement agrees; no heading changed. The builder went beyond the REQ in three recorded places: the sweep's label parenthetical (D-06), the step 4 empty-set skip (D-07) and the user guide sentence (D-09).
**After review:** f5f95c65 fixes review F1 (work.md) and F3 (work-reference.md Folder Structure), merged with the same pre as 78037a3e. The cumulative range has the same four files, all in Scope.
**P-A-U honesty:** the three boxes were ticked by the integrator from the builder's hand-back; APPLY cross-checked against git diff --stat 992f28a9..78037a3e (4 files, all in Scope).
**Live data flow:** prose-only change. The readers are agents running `do-work run` and `review-work`; every new citation names a real heading or bold label (shipped-package-reference-contract.sh PASS on the fixed branch).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 78037a3e (detached checkout `.git/work-run-2026-10-07-191445/drain-head`)
**Result:** ✓ All passing. Exit 0 on the first run at this revision, 20:26:44Z to 20:28:55Z, gate wall 131s, load under 3 and no other gate running. Stage queue-kanban-fast-tests executed (422 tests, wall 41s, slowest file strict_behavior_regression_test.go 21.05s < 30s). Stage do-work-cli-fast-tests executed (876 tests, wall 62s, slowest file internal/finalization/finalization_recovery_test.go 20.70s < 30s). Green-gate record satisfied by advance. The review fixes (f5f95c65) were merged before this gate ran, so no gate ran at f7416a80.

**Focused tests:** `do-work/runs/work-2026-10-07-191445/REQ-641-probe.sh` (shipped-package-reference-contract.sh and contract-regressions.sh), launched by advance, exit 0; probe, scope-drift, run-blocked-check and green-gate records satisfied. The builder ran both tests on its branch (PASS 1s and 20s), and the reviewer ran them on the integrated tree at f7416a80 (PASS 1s and 21s).

**Red-green validation:** not applicable as a test. This is a prose-only change to shipped action files and the user guide, and no test parses review prose. The captured GREEN condition was checked live: this REQ is the last successful integration of its wave, and its own review applied the new wave-end condition (see `## Review`).

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 992f28a9..78037a3e
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — the four changed files matched subtree skills

## Review

**Overall: 93%** after one re-check (first pass 92%) | 2026-10-07T20:30:05Z

| Dimension | Score |
|-----------|-------|
| Requirements | 93% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Full review: `do-work/runs/work-2026-10-07-191445/REQ-641-review.md` (first pass on f7416a80, re-check on f7416a80..78037a3e).

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `skills/do-work/actions/work.md:370` decides "last" from the run manifest showing every other member "finalized or set aside", but the manifest contract (`work-reference.md:480`, `crew-members/background-agents.md` step 3 and Manifest Format) records only landed status and no step writes finalized or set-aside to a row; an orchestrator that follows the contract never passes the fact, or passes it on "landed" too early — impact-rule-change → fixed (f5f95c65, re-merged 78037a3e: "last" keyed on the archive and the orchestrator's set-aside knowledge; the manifest names membership only)

**Minor findings:**
- F2 `work.md:370` / `review-work.md:135`: if the member that integrates last in order is set aside before its review, the previous member's review already ran without the fact and no review receives it, so the check is skipped with no record; a set-aside member that never reached review is named "unread for the sweep" — impact-rule-change → report only
- F3 `review-work.md:135` cites `work-reference.md` → **Folder Structure** for the archive location, but that tree shows flat archived REQs only as "Legacy REQs (no UR)"; members of an open UR are archived flat (REQ-639, REQ-640), so a reader trusting the tree looks in `archive/UR-NNN/` and wrongly names them unread — impact-rule-change → fixed (f5f95c65, re-merged 78037a3e)
- F4 `skills/do-work/actions/sample-archived-req.md` `## Review` example (called "a complete example" by `work.md:541`) has no `**Restatement sweep:**` line — impact-negligible → report only
- F6 (re-check) `work.md:370` "the orchestrator knows which it set aside": under delegated integration the integrator did not set earlier members aside, a set-aside REQ keeps its claim in `working/` like an in-flight one, and the integrator brief (`work-reference.md:470`) does not list set-aside members, so the integrator cannot pass the fact — impact-rule-change → report only
- F5 (Nit) `review-work.md:135` "a member whose review has no such line": a pre-template review can hold the phrase as inline prose (REQ-639's Acceptance line), which a grep matches; read here as unread — impact-negligible → report only
**Acceptance:** Pass — wave-end condition applied to this review as written (both earlier members unread, nothing re-derived); `shipped-package-reference-contract.sh` PASS and `contract-regressions.sh` pass at f7416a80; no heading changed in work-reference.md, review-work.md, work.md or docs/work-guide.md; `git diff --check` clean; re-check at 78037a3e: `shipped-package-reference-contract.sh` PASS, F1 and F3 resolved, no restatement broken
**Restatement sweep:** redefined the Restatement Sweep trigger set (wave-end inheritance, step 4 empty-set skip, widened label), the `**Restatement sweep:**` record line of the `## Review` template, the Verification Checklist sweep item, the wave-end fact work.md Step 7 passes the review agent (re-check: decided from the archive, not manifest status), and the work-reference.md Folder Structure archive bullet (re-check: open-UR REQs sit flat); unread for the sweep: REQ-639, REQ-640
**Suggested testing:** 2 items
**Follow-ups created:** None (4 findings report only, 2 fixed)

### Domain Review

General domain: no debug artifacts, no prime file touched, every new citation resolves; P-A-U boxes all ticked.

*Reviewed by review-work action*

## Decisions

Builder decisions D-01 to D-07 and D-09 are recorded in full in the hand-back, `do-work/runs/work-2026-10-07-191445/REQ-641-handback.md` § Decisions and § Addendum. Summary:
- **D-01 (DECIDE & STATE):** the orchestrator, not the reviewer, decides "last" and passes the fact plus the earlier members' REQ ids, like the merge range. Value: the reviewer parses no manifest. Risk: an orchestrator that forgets skips the check silently; work.md Step 7's MUST paragraph names the case.
- **D-02 (DECIDE & STATE):** a run with no manifest has no wave.
- **D-03 (DECIDE & STATE):** an earlier member with no record line is named "unread for the sweep", never re-derived from its diff.
- **D-04 (DECIDE & STATE):** the archive location is cited as work-reference.md → **Folder Structure**, not restated.
- **D-05 (DECIDE & STATE):** the record line holds this diff's own elements only, never inherited ones.
- **D-06 (DECIDE & STATE):** the sweep's label parenthetical widened so it does not restate the old trigger.
- **D-07 (DECIDE & STATE):** step 4 skips only when the whole trigger set is empty.
- **D-09 (DECIDE & STATE):** the coordinator extended Scope and `write_set` with `skills/do-work/docs/work-guide.md` at hand-back; one user-facing sentence.

Integrator decisions:
- **D-10 (DECIDE & STATE): review F1 and F3 fixed before the gate (f5f95c65, re-merged as 78037a3e).** F1: no step writes "finalized or set aside" into a manifest row, so work.md now takes membership from the manifest and done-ness from the archive and the orchestrator's own set-aside knowledge. F3: the cited **Folder Structure** listed flat archive files as legacy only, while REQs of a still-open UR archive flat; its tree comment and `archive/` bullet now say so (work-reference.md is in the write set). Value: the check fires on evidence that exists, and the cited location is right. Risk: none; the gate had not run, so the fix cost one re-merge.
- **D-11 (DECIDE & STATE): review F2 and F6 stay report only, recorded as one discovered task.** Both are the set-aside edge of the wave-end decision (the last member in order set aside before review; an integrator that cannot tell set-aside from not-yet-integrated). A fix needs the coordinator's brief contract to carry set-aside members, which is a design choice beyond this REQ's line and trigger. Value: the shipped rule is right for the normal case now. Risk: in those edge cases the check is skipped with no record.
- **D-12 (DECIDE & STATE): review F4 and F5 stay report only.** F4: `sample-archived-req.md` is already behind the template (it uses `**Findings:**`), so adding one line there would not make it a faithful copy. F5: affects only reviews written before this release.
- **D-13 (DECIDE & STATE): the wave-end acceptance test.** This REQ is the last successful integration of the UR-138 wave after REQ-639 (delegated integration, the coordinator shape) and REQ-640 (mid-run messages). The integrator passed the reviewer the fact and both ids; the reviewer read both archived reviews, found no `**Restatement sweep:**` line in either (REQ-639's Acceptance line holds only an inline prose phrase), named both unread for the sweep and re-derived nothing. That is the outcome the new text prescribes.
- **D-14 (DECIDE & STATE): no builder-work timing event.** The hand-back landed at 20:19:29Z, before this integrator started, and the recorder times from start to now (work.md Step 6, **Landed hand-back**).
- **D-15 (DECIDE & STATE): the lesson stays in family `alternate-writer-contract-drift`, as the builder proposed,** extended with review F1: a rule that tells an agent to read a state must name a writer of that state.

## Discovered Tasks

- impact-rule-change: the wave-end decision has two set-aside gaps (review F2 and F6). If the member that integrates last in order is set aside before its review, no review receives the wave-end fact. Under delegated integration the integrator cannot tell a set-aside member from one not yet integrated, because the `REQ-NNN-integrate.md` brief (`actions/work-reference.md` → **Delegated integration — the coordinator shape**) lists the wave membership but not the set-aside members. Listing set-aside members in the brief, and passing the fact to the last reviewed member when a later one is set aside, would close both. → report only
- impact-negligible: `skills/do-work/actions/sample-archived-req.md` `## Review` example lacks the `**Restatement sweep:**` line and already uses an older findings shape (review F4). → report only

## Lessons Learned

**What worked:** Applying the new rule to its own review was a real test, not a formality. It ran on the live chain (REQ-639 and REQ-640 archived, REQ-641 last) and showed both the normal path (both members unread, nothing re-derived) and two input gaps (F1, F3). Launching the reviewer right after qualify again put its Important finding before the gate, so the fix cost one re-merge and no extra gate.
**What didn't:** The first build keyed "last" on a manifest state ("finalized or set aside") that nothing writes; only this run's coordinator wrote it as free text, so it looked true here. The cited **Folder Structure** described flat archives as legacy only, which this very wave contradicted.
**Worth knowing:** When a rule tells an agent to read a state, name the writer of that state, and check the cited home describes the case you are in. A durable per-REQ fact that a later reviewer inherits belongs on a fixed line in the review block, decided by the orchestrator, not inferred by the reviewer.

## Orientation

Now the review of the last REQ to integrate in a fan-out wave also checks whether a later REQ restated something an earlier one redefined. Each review records what its own change redefined on one `**Restatement sweep:**` line, and the orchestrator passes the last reviewer the earlier members' ids, decided from the archive. This lives in `actions/review-work.md` Step 6 → **Restatement Sweep** and the Append to REQ File template, with the orchestrator side in `actions/work.md` Step 7. `_dev/primes/prime-action-files.md` is not stale: no path it names moved or was removed.

## Heavy Verification Plan

- Base revision: 992f28a907ccd860544f2faaae002cc2572a0f23
- Target revision: 78037a3e350e19d871ad5890a2aaf0efdbdef776
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/actions/work-reference.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills; skills/do-work/docs/work-guide.md matched subtree skills

## Heavy Verification Result

- Target revision: 78037a3e350e19d871ad5890a2aaf0efdbdef776
- Execution revision: 78037a3e350e19d871ad5890a2aaf0efdbdef776 (detached checkout `.git/work-run-2026-10-07-191445/drain-head`, QUEUE_KANBAN_BROWSER set)
- staged-skills: exit 0, executed, 40s, no HEAVY-RUN-LANE-SKIPPED finding (20:28:55Z to 20:29:36Z)


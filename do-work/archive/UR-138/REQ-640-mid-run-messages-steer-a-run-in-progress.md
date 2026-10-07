---
id: REQ-640
title: '[impact-rule-change] Mid-run messages: route a user message that arrives during a run to the right REQ in the user''s own words without stopping the run'
status: completed
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  basis:
  - Route B
  - 3-file write set
  - 4 acceptance criteria
  - cross-route regression gates
  calculated_at: 2026-10-07T19:41:13Z
created_at: 2026-10-07T19:12:33Z
user_request: UR-138
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-639, REQ-641]
batch: coordination-lessons
depends_on: [REQ-639]
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, skills/do-work/actions/capture.md, skills/do-work/actions/work-reference.md, skills/do-work/docs/work-guide.md]
claimed_at: 2026-10-07T19:39:49Z
dispatch_at: 2026-10-07T19:48:43Z
builder_handback_at: 2026-10-07T19:53:16Z
integration_at: 2026-10-07T19:54:35Z
review_at: 2026-10-07T20:03:18Z
kb_status: pending
commit: 0982069ec58c4ce3850ac398825fb2e028b3fd20
heavy_verified_at: 2026-10-07T20:05:46Z
heavy_verified_revision: 0982069ec58c4ce3850ac398825fb2e028b3fd20
completed_at: 2026-10-07T20:06:23Z
release_at: 2026-10-07T20:06:23Z
---
# Mid-Run Messages: Steer a Run While It Is in Progress
## What
One rule for a user message that arrives while a run is in progress, so a change of emphasis reaches the right REQ in the user's own words without stopping the run. Disk is the durable home; forwarding to a running builder is a courtesy.
## Why
Today a mid-run steering message has no path: capture turns it into a next-run addendum REQ and the in-flight builder never sees it. The lesson comes from the UR-137 run on 2026-10-07 (fixing accepted consumer-review findings) and Matt Maher's video "The New Way to Work With AI" (2026-10-07).
## Verified Facts (from the source)
- The takeover reset strips only generated headings (requeststate/state_apply.go:981).
- The advance classifier ignores unknown sections.
- Go reads only the `addendum_to` frontmatter field.
- So a `## Addendum (mid-run)` body section on a working REQ is safe and survives.
## Detailed Requirements
1. `skills/do-work/actions/work.md`: new subsection `### Mid-Run Messages (any step)` directly after Step 3.5, beside its existing "question escalated mid-run and answered" paragraph. Condition-keyed, not a list of message types:
   - A question about progress or state is answered from disk (REQ files, run manifest, hand-backs), never from memory of the conversation, and the run does not stop.
   - A change that affects an in-flight REQ whose builder has not handed back: the durable record is `## Addendum (mid-run)` on the working REQ, holding the user's words under `actions/clarify.md` Step 4's Outside-text containment plus one framing line (extends, narrows, or corrects which requirement). In the serial loop the orchestrator writes it now. Under delegated integration (REQ-639, the coordinator shape) the coordinator puts the text in that REQ's integrator brief and the integrator writes the section before the merge (one writer under the project root). Forward it to the builder where the harness can message a running agent; otherwise say in the progress output that the builder will not see it and review will judge it. The section is user intent, so the takeover reset keeps it.
   - A change that affects an in-flight REQ whose hand-back already landed: capture's existing in-flight path (new REQ with `addendum_to`), reported as queued for the next loop.
   - A change to a queued REQ, or new work: capture's existing paths, run in the coordinator's writing gap.
   - A change of emphasis for the hand-back only (what to show first): write it to the run manifest, not memory; the Decision Brief reads it.
   - Two rules across all branches: forward the user's words, never a paraphrase, because meaning is lost in restatement; route the message to every REQ it affects, not only the one the user named.
2. `skills/do-work/actions/review-work.md` Step 5 item 1: extract requirements from What/Detailed Requirements, any `## Addendum` section, and the UR.
3. `skills/do-work/actions/capture.md` Step 2 table, `do-work/working/` row: keep "NEVER modify" for capture; add one pointer clause: a message that arrives during a live run in this session is routed by `actions/work.md` -> Mid-Run Messages first.
4. `work.md` Orchestrator Checklist: one line "Mid-run message: route per Mid-Run Messages; never stop the run for it."
## Constraints
- Prose-only changes to shipped action files. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- The subsection is keyed on conditions, not a closed list of message types.
## Dependencies
Depends on REQ-639 (delegated integration, the coordinator shape): the delegated-integration branch of the new subsection uses REQ-639's integrator brief and one-writer rule. REQ-641 (wave-end consistency check) depends on this REQ. The serial order was requested because all three touch `work.md`, `review-work.md` and the Fan-Out Dispatch section of `work-reference.md`.
## Builder Guidance
Certainty is high: the source names every insertion point and the branch conditions. Latitude is limited to wording. Do not turn the subsection into a list of message types.
## Red-Green Proof
**RED prompt/case:** During a run, the user says "the builder should surface waiting-on-me items, and put decisions first in the update". Today this has no path: capture turns it into a next-run addendum REQ and the in-flight builder never sees it.
**Why RED now:** `work.md` has no mid-run message rule, and review does not read `## Addendum` sections when extracting requirements.
**GREEN when:** The new subsection names the in-flight path, and review reads the addendum.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read plus the citation contract test (`_dev/tests/shipped-package-reference-contract.sh`).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 6007 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the source names family `alternate-writer-contract-drift`, and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-138/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (ticked by the orchestrator from the builder's hand-back) Read prime-action-files.md, the alternate-writer-contract-drift and restated-mechanism-unchecked bullets of lessons-action-files.md, general.md, coding-guardrails.md, shared-principles.md and communication-style.md. Re-checked the Verified Facts in Go: the takeover reset keeps `## Addendum`, and `advanceSections` refuses a duplicated section name, so one section holds every entry. Approach: one condition-keyed subsection, one checklist line, the addendum added to every requirement-source reader, a capture pointer in the house arrow form, and the Decision Brief and manifest row naming the emphasis note.
- [x] **[APPLY]:** (ticked by the orchestrator from the builder's hand-back) Edits as planned in the four write-set files (d7e2e186), then the work-guide paragraph after the coordinator extended Scope (92f662d7). Nothing outside the write set.
- [x] **[UNIFY]:** (ticked by the orchestrator from the builder's hand-back) `git diff --stat` on the branch: capture.md, review-work.md, work-reference.md, work.md and work-guide.md. `git diff --check` clean. shipped-package-reference-contract.sh PASS 1s, contract-regressions.sh PASS 21s. Every new citation resolves, no debug artifacts.
*Source: "REQ B. Mid-run messages: steer a run while it is in progress." in the maintainer session's capture source, 2026-10-07.*

## Triage

**Route: B** - Medium

**Reasoning:** The REQ fixes the outcome, the three files, the insertion points and the branch conditions, so no Plan agent is needed. The rule it adds has reach the REQ does not list: who writes into a working REQ mid-run, what the integrator brief carries, and where review reads requirements from are restated in work.md, work-reference.md, capture.md, clarify.md and review-work.md, so exploration confirms each anchor and every restatement before Scope is fixed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Exploration ran in the orchestrator session (targeted greps, full reads of the cited paragraphs, and reads of the Go code behind each Verified Fact), against `31a6e35f`.

- **Required lessons consult:** the index matches two satellites, `_dev/primes/lessons-action-files.md` (now 6007 tokens) and `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (18680 tokens). Both are `slugged: partial`, so no targeted form is legal and both stay dropped for budget; `required_lessons` stays absent. Read anyway: the six `alternate-writer-contract-drift` bullets (lines 52-60) and the new REQ-639 bullet `restated-mechanism-unchecked` (line 61): a REQ's Verified Facts are claims, so each one below was checked against the code.
- **Verified Facts, checked:**
  - Takeover reset (`internal/requeststate/state_apply.go` line 981, `generatedRecoveryHeading`): strips only Triage, Exploration, Plan, Scope, Pre-Flight, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation, Decisions, Discovered Tasks and Timing. `## Addendum (mid-run)` survives. True.
  - Advance classifier (`internal/lifecycleadvance/advance_commands.go` lines 139-259) tests only named lifecycle sections. True, with one catch: `advanceSections` (line 356) refuses any section name that appears twice ("duplicate lifecycle section"). A second mid-run message for the same REQ must append inside the existing `## Addendum (mid-run)` section, never add a second heading with the same name.
  - Go reads only the `addendum_to` frontmatter field (`requestmodel.AddendumTo`); no Go code reads a `## Addendum` body section, and no validator refuses an unknown body section. True.
  - Outside-text containment (`actions/clarify.md` lines 101-106): a body passage goes in a blockquote whose lines open a code fence longer than the longest backtick run. `VisibleSections` skips fenced lines, so a `## ` line inside the user's words cannot open a section.
- **Claim in the REQ that is not yet true:** "write it to the run manifest; the Decision Brief reads it". Today `actions/work-reference.md` → **Decision Brief (hand-back format)** (lines 846-868) names no run-manifest source, and the guardrail table row for `manifest.md` (line 480) lists only REQ id, builder, operative name, hand-back file, landed status and held dispatch instant. For the new branch to work, both must name the hand-back emphasis note. A run manifest exists only when the run has a run directory (fan-out, or a background builder per `crew-members/background-agents.md`); a serial in-session run may have none. The builder states what the rule does then.
- **Insertion points:**
  - Requirement 1: `skills/do-work/actions/work.md` Step 3.5 runs lines 137-163; the "question escalated to the user mid-run and answered" paragraph is line 155. The new `### Mid-Run Messages (any step)` goes after line 163 and before `### Mechanical Evidence-Gate Loop` (line 165).
  - Requirement 2: `skills/do-work/actions/review-work.md` Step 5 item 1 is line 87. Step 2 line 55 ("What was requested — the What/Detailed Requirements sections") restates the same source list and needs the same addendum clause.
  - Requirement 3: `skills/do-work/actions/capture.md` Step 2 table, `do-work/working/` row, line 124. The **Immutability Rule** paragraph (line 63) restates the same rule ("If someone wants to add to an in-flight ... request, create a new addendum REQ") and a capture reader in a live run lands there first; builder judges a pointer there too.
  - Requirement 4: `work.md` Orchestrator Checklist lines 476-500; the Step 3.5 line is 483.
- **Restatements checked and consistent (no change expected):**
  - `work-reference.md` line 51: `working/` is "Immutable to all actions except the work pipeline", which covers an orchestrator or integrator writing the addendum.
  - `work-reference.md` line 470, **Delegated integration — the coordinator shape**: the brief already carries "any mid-run addendum text", and the one-writer rule agrees with the coordinator putting text in the brief during the gap. Seam for the builder: a message that arrives while that REQ's own integrator is already running cannot go in its brief; state where it goes (forward to the running integrator where the harness can message it, otherwise the coordinator's session scratch until the gap) without breaking the one-writer rule.
  - `work.md` line 135 (addendum REQs at Triage), line 519-527 (Progress Reporting, where "the builder will not see it" is said), `docs/work-guide.md` lines 72 and 101, `docs/capture-guide.md` line 76 ("never modifies files in working/"; still true for capture), `crew-members/background-agents.md` (nothing on mid-run messages).
- **Test pins:** no test pins the prose being changed. `_dev/tests/shipped-package-reference-contract.sh` resolves `file` → **Heading** citations, so capture.md's pointer must cite `actions/work.md` → **Mid-Run Messages (any step)** in the house arrow form with the real heading text, and every work-reference.md heading keeps its text.

*Generated by work action*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — new Mid-Run Messages (any step) subsection after Step 3.5; one Orchestrator Checklist line
- `skills/do-work/actions/review-work.md` (modify) — Step 5 item 1 and Step 2 requirement sources include any Addendum section
- `skills/do-work/actions/capture.md` (modify) — working row pointer clause in the Step 2 table; Immutability Rule pointer only if the builder judges it needed
- `skills/do-work/actions/work-reference.md` (modify) — Decision Brief names the run manifest's hand-back emphasis note as a source; the manifest.md guardrail row names that note
- `skills/do-work/docs/work-guide.md` (modify) — one user-facing sentence on messages sent during a run; scope extended by the coordinator at hand-back after the builder reported it as a discovered task (the builder records the extension as a D-XX in its hand-back)

**Files I will NOT touch:** docs/capture-guide.md (swept, consistent), crew-members/background-agents.md, actions/clarify.md (cited, not changed), any Go source or test, SKILL.md, any heading text in work-reference.md. Written by finalization in the main tree, not on the branch: every release path, _dev/primes/lessons-action-files.md and do-work/lessons-index.md.

**Acceptance criteria (restated from REQ):**
- [ ] work.md has `### Mid-Run Messages (any step)` directly after Step 3.5, keyed on conditions, not a list of message types
- [ ] A progress or state question is answered from disk (REQ files, run manifest, hand-backs), never from conversation memory, and the run does not stop
- [ ] A change to an in-flight REQ whose builder has not handed back is recorded as `## Addendum (mid-run)` on the working REQ: the user's words under Outside-text containment plus one framing line; written now in the serial loop, or carried in the integrator brief and written by the integrator before the merge under delegated integration; forwarded to a running builder where the harness allows, otherwise the progress output says review will judge it; the takeover reset keeps it
- [ ] A change to an in-flight REQ whose hand-back already landed uses capture's in-flight path (new REQ with addendum_to), reported as queued for the next loop
- [ ] A change to a queued REQ, or new work, uses capture's existing paths in the coordinator's writing gap
- [ ] A change of emphasis for the hand-back only is written to the run manifest, not memory, and the Decision Brief reads it
- [ ] Across all branches: forward the user's words, never a paraphrase; route to every affected REQ, not only the named one
- [ ] review-work.md Step 5 item 1 extracts requirements from What/Detailed Requirements, any `## Addendum` section, and the UR
- [ ] capture.md Step 2 table `do-work/working/` row keeps NEVER modify and adds the pointer to work.md Mid-Run Messages for a message during a live run in this session
- [ ] work.md Orchestrator Checklist has: Mid-run message: route per Mid-Run Messages; never stop the run for it
- [ ] Every work-reference.md heading is unchanged; shipped-package-reference-contract.sh and contract-regressions.sh pass

## Pre-Flight

**Git:** ✓ Integration tip 31a6e35f on `main` (this REQ's claim commit); the only dirt is this REQ's own trail (working REQ, untracked `do-work/runs/work-2026-10-07-191445/REQ-640-probe.sh`) and `do-work/working/baseline.json`, which this pre-flight rewrote
**Tests baseline:** ✓ focused contract tests green (`do-work/runs/work-2026-10-07-191445/REQ-640-probe.sh`, launched by advance: shipped-package-reference-contract.sh 1s, contract-regressions.sh 21s)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 31a6e35f from the detached checkout `.git/work-run-2026-10-07-191445/drain-head` (19:43:58Z to 19:46:09Z, gate wall 131s, load under 4, no other gate running; queue-kanban-fast-tests and do-work-cli-fast-tests executed, slowest files strict_behavior_regression_test.go 20.57s and finalization_recovery_test.go 20.97s < 30s); preflight and green-gate records satisfied
**Dependencies:** ✓ prose-only change; no toolchain or module dependency

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/capture.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** work.md gains `### Mid-Run Messages (any step)` between Step 3.5 and the Mechanical Evidence-Gate Loop. It is keyed on what the message changes: a progress question is answered from disk; a change to a REQ whose builder has not handed back becomes one `## Addendum (mid-run)` section on the working REQ (written now in the serial loop, or by the integrator from its brief under delegated integration); a change to handed-back work, queued work or new work goes through capture; an emphasis-only change is the hand-back emphasis note in the run manifest. The user's words are carried verbatim and routed to every affected REQ. A second message for one REQ appends an entry inside the same section, because advance refuses a duplicated section name. The Orchestrator Checklist gains one line, and Route C plan validation reads any `## Addendum` section too. review-work.md Step 2 and Step 5 item 1 extract requirements from any `## Addendum` section. capture.md's `do-work/working/` row keeps NEVER modify and points to the new subsection. work-reference.md's Decision Brief reads the emphasis note first, and the `manifest.md` guardrail row names it. work-guide.md tells the user they can keep talking during a run. No heading changed, no integration seams. After review the orchestrator made three one-line fixes on the builder branch (9079ab68, review F1 to F3): the emphasis-note branch waits for the coordinator's writing gap under delegated integration, the not-yet-handed-back branch is keyed on the session being the only writer under the project root, and the Decision Brief note may reorder the sections as well as their contents (D-10). The review re-check then found that work.md Progress Reporting still stated a fixed section order; 5c68b501 names the emphasis note as its one exception (review F5, D-12). Merge range 2ab092c0..0982069e (builder commits d7e2e186 and 92f662d7, orchestrator commits 9079ab68 and 5c68b501, merges b36f7c65, 33c41c50 and 0982069e).

## Qualification

**Diff range:** 2ab092c0..b36f7c65 for the qualify record (builder commits d7e2e186 and 92f662d7, merge b36f7c65); cumulative range after the review fixes 2ab092c0..0982069e.
**Gate records:** qualify satisfied; scope-drift satisfied (the five changed files equal the declared Scope, including work-guide.md, which the coordinator added at hand-back, D-09).
**Warnings judged:** none raised.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. 1: `### Mid-Run Messages (any step)` sits right after Step 3.5 and before the Mechanical Evidence-Gate Loop, keyed on what the message changes; it carries the two cross-branch rules (the user's own words, every affected REQ), the disk-answered progress branch, the `## Addendum (mid-run)` branch with the framing line, Outside-text containment, the serial and delegated writers, forwarding or the progress-line fallback and the takeover-reset note, the capture branch for handed-back work, the capture branch for queued or new work, and the run-manifest emphasis note that the Decision Brief reads. 2: review-work.md Step 5 item 1 (and Step 2, D-06) reads any `## Addendum` section. 3: capture.md's `do-work/working/` row keeps NEVER modify and points to the new subsection. 4: the checklist line is present verbatim. No work-reference.md heading changed. The builder went beyond the REQ in three places, each recorded: one entry per message inside one section (D-03, advance refuses a duplicated section name), Route C plan validation reads addenda (D-05), and the user guide paragraph (D-09).
**After review:** 9079ab68 fixes review F1 to F3 in work.md and work-reference.md only, merged with the same pre as 33c41c50; 5c68b501 fixes re-check finding F5 in work.md only, merged as 0982069e. The cumulative range has the same five files, all in Scope.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY cross-checked against git diff --stat 2ab092c0..0982069e (5 files, all in Scope).
**Live data flow:** prose-only change. The readers are agents running `do-work run`, `review-work` and `capture`, and every new citation names a real heading or bold label (shipped-package-reference-contract.sh PASS on the fixed branch).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 33c41c50 (detached checkout `.git/work-run-2026-10-07-191445/drain-head`)
**Result:** ✓ All passing. Exit 0 on the first run at this revision, 19:59:15Z to 20:01:33Z, gate wall 138s, load under 3 and no other gate running. Stage queue-kanban-fast-tests executed (423 tests, wall 43s, slowest file strict_behavior_regression_test.go 20.78s < 30s). Stage do-work-cli-fast-tests executed (876 tests, wall 62s, slowest file internal/finalization/finalization_recovery_test.go 21.00s < 30s). Green-gate record satisfied by advance. The review fixes (9079ab68) were merged before this gate ran, so no gate ran at b36f7c65. The heavy lane also ran green at 33c41c50 (staged-skills exit 0, executed, 41s, 20:01:42Z to 20:02:23Z).

**Repository gate after the re-check fix:** the review re-check found F5 after the advance test gate was satisfied, so 5c68b501 was merged as 0982069e and the gate ran directly at that revision from the same detached checkout (advance refuses gate input once the REQ is past testing): exit 0, 20:02:43Z to 20:04:59Z, wall 136s; queue-kanban-fast-tests 423 tests, wall 42s, slowest file 21.07s < 30s; do-work-cli-fast-tests 876 tests, wall 59s, slowest file 20.46s < 30s. The focused contract tests passed inside it (shipped-package-reference-contract.sh PASS).

**Focused tests:** `do-work/runs/work-2026-10-07-191445/REQ-640-probe.sh` (shipped-package-reference-contract.sh and contract-regressions.sh), launched by advance, exit 0; probe record satisfied. The builder ran both tests on its branch (PASS 1s and 21s), and the reviewer ran them on the integrated tree at b36f7c65 (PASS).

**Red-green validation:** not applicable as a test. This is a prose-only change to shipped action files and the user guide, with no behavior a test pins (Exploration: no test under `_dev/tests` pins the changed prose). The captured GREEN condition is checked by reading: the new subsection names the in-flight path, and review-work.md reads `## Addendum` sections.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 2ab092c0..0982069e (same single lane as the 2ab092c0..33c41c50 plan)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — the five changed files matched subtree skills

*Verified by work action*

## Review

**Overall: 94%** after two re-checks (first pass 91%) | 2026-10-07T20:03:18Z

| Dimension | Score |
|-----------|-------|
| Requirements | 97% |
| Code Quality | 94% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

Full review: `do-work/runs/work-2026-10-07-191445/REQ-640-review.md` (requirements checklist, Go claim checks, restatement sweep, re-checks).

**Mechanism claims checked against Go:** all true. The takeover reset keeps `## Addendum (mid-run)` (`generatedRecoveryHeading`, state_apply.go:981). `advance` refuses a duplicated section name (`advanceSections`, advance_commands.go:344-360); capture's dated `## Addendum (<date>)` has a different name, so both can coexist. Outside-text containment keeps the user's words from opening a section, because of the `> ` line prefix.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `skills/do-work/actions/work.md` emphasis-note branch: no writing-gap clause, so under delegated integration it told the coordinator to edit `manifest.md` while an integrator runs, which the coordinator shape forbids — impact-rule-change → fixed before the gate (9079ab68, D-10)

**Minor findings:**
- F2 `work.md` not-yet-handed-back branch named only the serial loop and delegated integration, not a self-integrating fan-out run — impact-rule-change → fixed (9079ab68)
- F3 `work-reference.md` Decision Brief: the note could only reorder items inside each section, so the REQ's own example ("decisions first") could not be honored — impact-rule-change → fixed (9079ab68)
- F4 `work-reference.md` Crash Recovery and `docs/work-guide.md` recovery paragraph say a takeover strips "orchestrator sections", though the reset keeps the orchestrator-written addendum — impact-rule-change → report only (D-11)
- F5 (re-check) `work.md` Progress Reporting still stated a fixed section order after the F3 fix — impact-rule-change → fixed (5c68b501, D-12); the second re-check at 0982069e found F5 resolved and nothing new

**Acceptance:** Pass — every Detailed Requirement and acceptance criterion traces to the diff; shipped-package-reference-contract.sh and contract-regressions.sh PASS on the integrated tree; no work-reference.md heading changed; `git diff --check` clean
**Suggested testing:** 0 items (prose-only)
**Follow-ups created:** None (no finding is impact-critical)

*Reviewed by review-work action*

## Decisions

Builder decisions D-01 to D-09 are recorded in full in the hand-back, `do-work/runs/work-2026-10-07-191445/REQ-640-handback.md` § Decisions and § Addendum. Summary:
- **D-01 (DECIDE & STATE):** the emphasis branch is keyed on "the run has a run directory"; without one, the session that renders the Decision Brief holds the note itself and creates no file.
- **D-02 (DECIDE & STATE):** a change to a REQ whose integrator is already running takes the handed-back path (capture with `addendum_to`), run by the coordinator in its next writing gap.
- **D-03 (DECIDE & STATE):** one `## Addendum (mid-run)` section per REQ with timestamped entries inside, because `advance` refuses a duplicated section name.
- **D-04 (DECIDE & STATE):** under delegated integration the coordinator keeps the words in session scratch until its writing gap, then puts them in the integrator brief.
- **D-05 (DECIDE & STATE):** Route C plan validation also reads any `## Addendum` section (restates the same source list).
- **D-06 (DECIDE & STATE):** review-work.md Step 2 changed the same way as Step 5 item 1.
- **D-07 (DECIDE & STATE):** capture.md's **Immutability Rule** stays unchanged; the pointer sits only in the Step 2 table.
- **D-08 (DECIDE & STATE):** the checklist line sits right after the Step 3.5 line.
- **D-09 (DECIDE & STATE):** the coordinator extended Scope and `write_set` with `skills/do-work/docs/work-guide.md` at hand-back; one user-facing paragraph.

Orchestrator decisions:
- **D-10 (DECIDE & STATE): review F1 to F3 fixed before the gate (9079ab68).** Each was one sentence in a Scope file and the gate had not run, so the fix cost one re-merge. F1 restores the one-writer rule on the fourth branch; F2 keys the write-now case on the condition (sole writer) instead of naming two run shapes; F3 lets the emphasis note reorder sections, because the REQ's own RED example asks for "decisions first".
- **D-11 (DECIDE & STATE): review F4 stays report only.** The new subsection already says the addendum is user intent, not a generated section, so the takeover keeps it; the same "orchestrator sections" wording also sits in `actions/restart-with-parallel-handoff.md`, outside the write set, and changing two of three copies would make them disagree. Reversible prose.
- **D-12 (DECIDE & STATE): review re-check F5 fixed after the gate (5c68b501).** The F3 fix made work.md Progress Reporting contradict the Decision Brief in the same file. One clause; cost a re-merge and a direct gate and heavy re-run at 0982069e.
- **D-13 (DECIDE & STATE): no builder-work timing event for this REQ.** The hand-back landed at 19:53:16Z and the integrator read its recipes first, so a start-to-now event would charge the integrator's reading to the builder (work.md Step 6, **Landed hand-back**).
- **D-14 (DECIDE & STATE): the lesson stays in family `alternate-writer-contract-drift`, as the builder proposed.** A new requirement source is a reader-side contract change, and the failure mode (restated source lists and a rule repeated on sibling branches) is that family's.

## Discovered Tasks

- impact-rule-change: `skills/do-work/actions/work-reference.md` Crash Recovery, `skills/do-work/docs/work-guide.md` recovery paragraph and `skills/do-work/actions/restart-with-parallel-handoff.md` say a takeover strips the claim's "orchestrator sections"; the reset strips only generated sections, so the orchestrator-written `## Addendum (mid-run)` survives. Saying "generated sections" in all three would remove the doubt (review F4). → report only

## Lessons Learned

**What worked:** Consuming the landed hand-back with no re-dispatch, and launching the reviewer right after qualify: its Important finding (the one-writer clause missing from one branch) was fixed before the gate. Checking the REQ's Verified Facts against Go first (the REQ-639 lesson) found the duplicated-section refusal that shaped the one-section-with-entries rule.
**What didn't:** A fix to one restated rule can break another restatement in the same file. The F3 fix let the emphasis note reorder the Decision Brief's sections, and work.md Progress Reporting still stated the fixed order; the re-check caught it, and it cost one more merge and a second gate run.
**Worth knowing:** A new requirement source is a reader-side contract change: grep the source-list phrase ("What/Detailed Requirements"), not only the line the REQ names. When several sibling branches each repeat one clause (here, "in the coordinator's next writing gap"), read every branch for it; the one without it is where the rule breaks.

## Orientation

Now a user can keep talking to a running `do-work run`: a progress question is answered from disk, a change to a REQ still being built reaches it as a `## Addendum (mid-run)` section in the user's own words that review judges the build against, later changes go through capture, and a "show me X first" note orders the end-of-run Decision Brief. This lives in the work pipeline's step rules (`actions/work.md` → **Mid-Run Messages (any step)**), with readers in review, capture and the Decision Brief. `_dev/primes/prime-action-files.md` is not stale: no path it names moved or was removed.

## Heavy Verification Plan

- Base revision: 2ab092c004cc2294cfa9149a0f24171b684b367e
- Target revision: 0982069ec58c4ce3850ac398825fb2e028b3fd20
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/capture.md matched subtree skills; skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/actions/work-reference.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills; skills/do-work/docs/work-guide.md matched subtree skills

## Heavy Verification Result

- Target revision: 0982069ec58c4ce3850ac398825fb2e028b3fd20
- Execution revision: 0982069ec58c4ce3850ac398825fb2e028b3fd20 (detached checkout `.git/work-run-2026-10-07-191445/drain-head`, QUEUE_KANBAN_BROWSER set)
- staged-skills: exit 0, executed, 32s, no HEAVY-RUN-LANE-SKIPPED finding (20:04:59Z to 20:05:31Z)

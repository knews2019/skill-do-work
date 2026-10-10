---
id: REQ-675
title: 'Core review names four delivery stages and reports deployment and live acceptance as unassessed unless the review exercised them'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-09T23:02:37Z
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work/actions/review-work.md, skills/do-work/docs/review-work-guide.md]
related: [REQ-676, REQ-677, REQ-678]
batch: portable-verification-actions
claimed_at: 2026-10-09T22:56:57Z
dispatch_at: 2026-10-09T23:08:40Z
builder_handback_at: 2026-10-09T23:12:36Z
integration_at: 2026-10-09T23:12:53Z
review_at: 2026-10-09T23:24:45Z
kb_status: pending
commit: fec8d220748dda570dcefbfa49d924e40be3b55c
heavy_verified_at: 2026-10-09T23:25:09Z
heavy_verified_revision: fec8d220748dda570dcefbfa49d924e40be3b55c
completed_at: 2026-10-09T23:25:40Z
release_at: 2026-10-09T23:25:40Z
---
# Core Review Names Four Delivery Stages
## What
Make `skills/do-work/actions/review-work.md` the one home of four delivery stages: implementation, integration, deployment, and live acceptance. Step 7's Acceptance result covers only the stages the review actually exercised. Every applicable stage it did not exercise is named as unassessed in the report's Suggested Additional Testing. Add one matching line to `skills/do-work/docs/review-work-guide.md` (Phase 3). Nothing else in the review contract changes.
## Why
Today "Acceptance: Pass" can be read as production acceptance. Step 7 (`review-work.md:160-195`) only offers Untested when the code cannot run at all, and Step 8 (`:203`) lists environment testing as a category, but no sentence says a local pass says nothing about the served build. The maintainer's brief asked the review to "answer whether work exists, satisfies the request, and has evidence for applicable delivery stages". Existence (Step 1 no-diff exit, `:46`, and the Red Flag at `:478`) and satisfaction (Step 5, `:81`) are already covered. Only the per-stage evidence is missing. REQ-678 (the release-check toolbox action) cites these stage names, so they need one home first.
## Finding Provenance
From the validate-feedback triage of the maintainer's brief in this session (finding F5, verdict Accept, smallest form):
- **Verbatim claim:** "Answer whether work exists, satisfies the request, and has evidence for applicable delivery stages." Source: the brief's CORE REVIEW IMPROVEMENT section.
- **Evidence:** `review-work.md:46`, `:81`, `:190` (Pass/Partial/Fail/Untested), `:203`; `docs/review-work-guide.md:26-34`.
- **Surface-cost:** Earned, narrowly. Incident class: a REQ passes review while the served build is stale (named in the brief, not reproduced upstream). Surface: one rule at Step 7 plus one guide line. No new test file.
- The maintainer chose "Core owns stages, toolbox cites" for where the stage vocabulary lives.
## Detailed Requirements
1. In `review-work.md` Step 7, define the four stages once, each in one plain sentence: implementation (the diff does what the REQ asks), integration (it works in the merged tree with the rest of the system, for example the test suite and adjacent flows), deployment (the built or packaged result reached its serving environment or consumer install), live acceptance (the consumer-visible behavior is right in the real environment). Mark which stages apply as a judgment: many REQs have no deployment stage at all.
2. State that the Acceptance result scores only the stages the review exercised. An applicable stage the review did not exercise is listed in Step 8's Suggested Additional Testing as "unassessed", by stage name.
3. Add one line to `docs/review-work-guide.md` Phase 3 saying the same thing in user words.
## Constraints
- Do not change the persisted `## Review` block (Append to REQ File), the score table, the scoring formula and caps, the Verdict mapping, the impact tokens, Step 10 routing, or any status. Add no lifecycle status and no reopening behavior.
- Do not add a retrospective mode. The maintainer kept the push-back: Step 9.5 Lessons Learned already records missed assumptions.
- Do not restate the targeted-counterexample rule. It already exists (`review-work.md:113`, `:174`, `:181`, `:478` and `crew-members/shared-principles.md`).
- Repair addenda already preserve history (`review-work.md:66`, `:334`, `:348`). No change.
- Release per `_dev/primes/prime-releases.md`.
## Builder Guidance
Firm on scope: two files, a few sentences. Latitude on wording and on whether the Acceptance one-line summary in the human report names the stages covered. Run the Restatement Sweep on "Acceptance": `crew-members/shared-principles.md` (the "Acceptance cannot be exercised" row) and `actions/work.md` Step 7 read it and must still agree.
## Red-Green Proof
**RED prompt/case:** A reviewer following `review-work.md` reviews a scratch REQ "publish the updated rules page" whose diff is correct and whose local tests pass, with no access to the served site.
**Why RED now:** Nothing in Step 7 or Step 8 makes the report say that deployment and live acceptance were not checked. The likely result is "Acceptance: Pass" with no stage named.
**GREEN when:** The same one-off exercise produces Acceptance for the stages exercised (implementation and integration) and a Suggested Additional Testing entry naming deployment and live acceptance as unassessed. `git diff` shows no change to the Append to REQ File template, the Scoring Guidelines, the Verdict mapping, or Step 10.
**Validation:** Inferred during capture. The maintainer chose one-off exercises recorded in `## Testing` over kept fixtures.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes a rule that other actions restate (family `alternate-writer-contract-drift`).
## Full Context
See `do-work/user-requests/UR-151/input.md` for complete verbatim input. No queued REQ shares this intent: REQ-669 (the trace action) and REQ-659 (frontmatter set and append section) mention review-work only as context.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the brief, REQ, UR-151, general.md, coding-guardrails.md, shared-principles.md, communication-style.md, anti-slop.md, prime-action-files.md, prime-releases.md, lessons-releases.md (whole) and the alternate-writer-contract-drift bullets of lessons-action-files.md. Plan: define the stages once in Step 7 between "What NOT to do" and "If you can't run the code"; put the "scores only exercised stages" rule after the Score block without touching the Pass/Partial/Fail/Untested bullets; put the "unassessed" entry rule in Step 8 as a category so the rule lives in one place; one paragraph in guide Phase 3. (from the builder hand-back)
- [x] **[APPLY]:** Done as planned. +9 lines in review-work.md, +2 in the guide. No other file touched. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat e313e870`: review-work.md +9, review-work-guide.md +2, 2 files, 11 insertions, 0 deletions. REQ-675-probe.sh exit 0; `_dev/tests/shipped-package-reference-contract.sh` exit 0; `git diff --check` exit 0; `_dev/tests/contract-regressions.sh` exit 1 only on the quiet-grep pipeline audit of the REQ-676 and REQ-677 run probes (run artifacts, not this diff; fixed by the coordinator in c4dda51e). Two hunks in review-work.md, both inside Steps 7 and 8; Append to REQ File template, score table, Scoring Guidelines, Verdict mapping and Step 10 untouched. Both changed files read in full around the edits; no em-dashes in the added prose. (from the builder hand-back)
*Source: the brief's CORE REVIEW IMPROVEMENT, first bullet, accepted in the validate-feedback triage.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names both files and the exact sentences to add (four stage definitions and the unassessed rule in review-work.md Step 7/8, one line in review-work-guide.md Phase 3), with explicit constraints on what must not change. No location or pattern needs discovery; the Restatement Sweep on "Acceptance" is a review-time check, not exploration.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/docs/review-work-guide.md` (modified)

**What was done:** Step 7 of review-work.md now defines the four delivery stages once (implementation, integration, deployment, live acceptance), says which apply is a judgment, and adds one paragraph after the Score block: the Acceptance result scores only the stages the review exercised and its one-line summary names them. Step 8 gains a first category, "Unassessed delivery stages", required whenever an applicable stage was not exercised. The guide's Phase 3 gains one paragraph saying the same in user words. Merged as f7dad3b1 over 623f00d5; no integration seams. Review finding F2 was fixed on the builder branch (67e0ed86: the Step 7 sentence now says the covered stages go on the persisted `**Acceptance:**` line too) and re-merged as fec8d220; the cumulative range is 623f00d5..fec8d220.

## Decisions

(from the builder hand-back)

- D-01 DECIDE & STATE: the "unassessed" entry rule lives in Step 8 as the first category (required whenever such a stage exists), and Step 7 only points to it. One home for the rule, and it sits where the reviewer writes the entries.
- D-02 DECIDE & STATE: the Pass/Partial/Fail/Untested bullets are unchanged. The scoping sentence follows the Score block instead, so the score values every reader routes on keep their exact wording.
- D-03 DECIDE & STATE: used the latitude on the one-line summary: Step 7 tells the reviewer to name the stages covered there. The Append template itself is unchanged.
- D-04 DECIDE & STATE: the new Step 8 bullet uses a colon, not the em-dash its siblings use, per the no-em-dash rule. Small style mismatch inside the list, accepted.
- D-05 DECIDE & STATE: the existing Step 8 category "Integration testing" (third-party services, APIs) shares a word with the new "Integration" stage. Left as is: renaming it is outside the write boundary's intent. Noted for the reviewer.

## Discovered Tasks

(from the builder hand-back)

- **impact-negligible** (the builder wrote `impact-important`, which is not an impact token; restamped because the problem is already fixed) `_dev/tests/contract-regressions.sh` failed on base e313e870 because `do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh` and `REQ-677-probe.sh` decided on a quiet grep fed from a pipeline (quiet-grep pipeline audit). Resolved by the coordinator in c4dda51e (the probes now read the argument hint through a here-string), before this REQ's merge; the repository gate at f7dad3b1 is the evidence. → report only

## Qualification

**Gate records:** `advance REQ-675 --diff-range 623f00d5..f7dad3b1` returned the `qualify` gate as `satisfied` with provenance `merged_range` and no findings (no debug artifact, no output primitive, P-A-U boxes ticked).

**Requirement trace against `git diff 623f00d5..f7dad3b1` (2 files, +11/-0):**
1. Four stages defined once, one plain sentence each, in Step 7 (`review-work.md` new lines 188-193): implementation, integration, deployment and live acceptance, worded as the REQ asks; "Which ones apply is a judgment, and many REQs have no deployment stage at all" is the applicability rule. Met.
2. Acceptance scores only exercised stages (new paragraph after the Score block, line 202) and Step 8 lists each applicable unexercised stage as unassessed, by stage name, as a required first category (line 210). Met.
3. One matching paragraph in `review-work-guide.md` Phase 3 (line 36). Met; it is two sentences, which fits "one line" in intent.

**Constraints:** the three review-work.md hunks sit at lines 188, 202 and 210, all inside Step 7 (starts at 160) and Step 8 (starts at 204), before Scoring Guidelines (221). The Append to REQ File template (390), the score table, Verdict mapping, impact tokens, Step 10 (334) and every status are untouched. No retrospective mode, no restated counterexample rule. Added lines carry no em-dash.

**Scope:** Route A, `write_set` names exactly the two touched files. Declared equals touched; no `do-work/` path in the range.

**Re-merge (review finding F2):** the cumulative range 623f00d5..fec8d220 adds one sentence inside the same Step 7 paragraph (line 202). Still 2 files, +11/-0; the queue guard printed nothing; `advance` had already moved to the review phase, so qualification of the one-sentence delta is this judgment: no debug artifact, no em-dash, inside Step 7, template untouched.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at f7dad3b1 and again at fec8d220 after the F2 re-merge (direct, unpiped), then at f7dad3b1 `advance REQ-675 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-09-225703/REQ-675-probe.sh`.
**Result:** ✓ Repository gate passed on the first run at both merges (f7dad3b1: exit 0, gate wall 133 s, under load from two sibling builders; fec8d220: exit 0, gate wall 144 s). At fec8d220 the GREEN probe was re-run directly: exit 0, 1 s. `advance` recorded `green-gate` satisfied (`advance_executed`) and ran the GREEN probe: exit 0 (`BLOCKED-PROBE-SUCCEEDED`), which checks that Steps 7-8 name implementation, integration, deployment, live acceptance and unassessed, that guide Phase 3 says unassessed, and that `shipped-package-reference-contract.sh` holds.

**Red-green validation:** non-test behavioral exercise (the REQ chose a one-off exercise over a kept fixture), scenario from `## Red-Green Proof`: a REQ "publish the updated rules page", correct diff, local tests pass, no served site.
- Builder (from the hand-back, by reading the text, no reviewer agent spawned): RED at base e313e870 produced `**Acceptance:** Pass — rules page updated, tests pass.` with no stage named and no required deployment or live line. GREEN at f12beed1 produced `**Acceptance:** Pass — implementation and integration: the diff updates the rules page as asked and the local test suite passes.` plus Suggested Additional Testing lines `Deployment: unassessed. ...` and `Live acceptance: unassessed. ...`.
- Re-run at the merge: the GREEN probe above (exit 0 at both merges). The independent reviewer re-applied the edited Steps 7 and 8 to the same scenario at f7dad3b1: GREEN. It produced `**Acceptance:** Pass — implementation and integration: the diff updates the rules page content as the REQ asks and the local test suite passes; deployment and live acceptance unassessed (no access to the served site).` and Suggested Additional Testing lines `Deployment: unassessed. ...` and `Live acceptance: unassessed. Open the served rules page and confirm the new rules show.` The base text yielded a bare `Pass` with no stage lines (RED). The delta re-review re-checked the F2 sentence at fec8d220 (see `## Review`).
- Not exercised: a reviewer agent run on a real scratch REQ file end to end; the change is prose only.

**New tests added:** none (one-off exercise per the REQ; run probe `do-work/runs/work-2026-10-09-225703/REQ-675-probe.sh` is a run artifact, not a kept test).

**Heavy verification plan:**
- Range: 623f00d5710f47a160a963f0d53f23abe204d1c7..fec8d220748dda570dcefbfa49d924e40be3b55c (first planned at f7dad3b1 with the same single lane)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`: skills/do-work/actions/review-work.md and skills/do-work/docs/review-work-guide.md matched subtree skills

*Verified by work action*

## Review

**Overall: 95%** | 2026-10-09T23:24:45Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1. `crew-members/shared-principles.md:15` "Acceptance cannot be exercised → Record Untested" does not cover the partial case (implementation and integration exercised, deployment or live acceptance not), so a reviewer can score Untested instead of a scoped Pass. — impact-rule-change → report only
- F2. Resolved in 67e0ed86 (merged fec8d220): `actions/review-work.md:202` now names the `**Acceptance:**` line of the Append to REQ File block as where the covered stages go. — resolved
- F3. `actions/sample-archived-req.md:99` example `**Acceptance:**` line names no stage. The sample was already behind the current template. — impact-negligible → report only
- F4 (Nit). `actions/review-work.md:460` Route A "Suggested testing is usually empty or 1 item" versus the now-required unassessed-stage lines. "Usually" keeps them compatible. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration: every requirement is in the diff, the F2 fix names the persisted `**Acceptance:**` line, the GREEN probe re-run at fec8d220 exits 0, and the behavioral exercise on the "publish the updated rules page" scenario is GREEN; deployment and live acceptance unassessed (release happens at finalization, no consumer install checked).
**Restatement sweep:** redefined the Acceptance result's scope (only exercised delivery stages, named in the one-line summary and on the persisted `**Acceptance:**` line) and Step 8's required unassessed-stage category; swept shared-principles.md:14-15, work.md:358-368 and :518, the review-work.md Append template, Scoring Guidelines, Verdict mapping, Step 10, Calibrating Review Depth and Verification Checklist, review-work-guide.md, sample-archived-req.md, advance_commands_test.go:327, completed-work-presentation-reference.md:39 and every "one-line summary" hit in skills/; stale: shared-principles.md:15 (F1), sample-archived-req.md:99 (F3), review-work.md:460 (F4)
**Suggested testing:** 2 items
**Follow-ups created:** None (3 open findings report only, F2 resolved)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Putting the scoping sentence after the Score block, not inside the Pass/Partial/Fail/Untested bullets, kept the values that `work.md` routes on byte-identical, so the restatement sweep found no routing reader to change.
**What didn't:** The first build said "name the stages in the result's one-line summary" without naming which line. Review found three candidates, and the persisted `**Acceptance:**` line is filled from a template placeholder 200 lines later that gives no cue, so the durable record could still carry a bare Pass (finding F2, fixed in 67e0ed86).
**Worth knowing:** A rule about what a review writes must name the persisted line it lands on, not only the human report. `crew-members/shared-principles.md:15` ("Acceptance cannot be exercised → Record Untested") still does not cover the partial case where only deployment or live acceptance could not be checked (finding F1, report only). REQ-678 (release-check) cites the four stage names from the **Delivery stages** block in Step 7.

## Orientation

Now a review says which delivery stages its Acceptance result covers and lists the applicable unchecked ones (deployment, live acceptance) as unassessed; lives in the review contract, `skills/do-work/actions/review-work.md` Steps 7 and 8 (prime: `_dev/primes/prime-action-files.md`). [MAP CHANGED] new shared vocabulary: the four delivery stages now have one home that REQ-678 will cite.

## Heavy Verification Plan

- Base: 623f00d5710f47a160a963f0d53f23abe204d1c7
- Target: fec8d220748dda570dcefbfa49d924e40be3b55c
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. Reasons: skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/docs/review-work-guide.md matched subtree skills.

## Heavy Verification Result

- Target revision: fec8d220748dda570dcefbfa49d924e40be3b55c
- Execution revision: fec8d220748dda570dcefbfa49d924e40be3b55c (detached checkout `.git/work-run-work-2026-10-09-225703/drain-head-REQ-675`, removed after the run)
- staged-skills: exit 0, executed (fingerprint_mismatch), 35 s. An earlier run at the first merge f7dad3b1 also passed (exit 0, executed, 43 s) before review finding F2 moved the target.

## Timing

Observed 2026-10-09T23:08:40Z to 2026-10-09T23:24:25Z: 15m 45s total, 14m 42s attributed across 8 events, 1m 03s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 7m 05s | 4 |
| builder-work | 3m 57s | 1 |
| review | 3m 20s | 1 |
| handback-merge | 20s | 2 |

Slowest stage: builder-work / builder worktree build, 3m 57s, outcome success.
Slowest command: verification-gate / repository gate and probe, 2m 51s, exit 0, bash (2 argv tokens).

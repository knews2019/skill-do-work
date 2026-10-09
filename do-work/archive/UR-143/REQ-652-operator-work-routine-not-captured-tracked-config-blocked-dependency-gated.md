---
id: REQ-652
title: '[impact-rule-change] Routine operator work is never a REQ; a tracked operator configuration is a blocked, dependency-gated REQ that clarify completes'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-09T16:09:30Z
created_at: 2026-10-09T16:06:43Z
user_request: UR-143
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-650, REQ-651, REQ-653]
batch: october-review-triage
write_set: ["skills/do-work/actions/capture.md", "skills/do-work/actions/clarify.md", "skills/do-work/actions/work.md", "skills/do-work/actions/work-reference.md", "skills/do-work/docs/work-guide.md"]
claimed_at: 2026-10-09T16:08:10Z
dispatch_at: 2026-10-09T16:11:57Z
builder_handback_at: 2026-10-09T16:16:57Z
integration_at: 2026-10-09T16:17:15Z
review_at: 2026-10-09T16:28:49Z
kb_status: pending
commit: e8b3687cd31c0f6bdd9e9841f099fa22156bb7b1
heavy_verified_at: 2026-10-09T16:29:10Z
heavy_verified_revision: e8b3687cd31c0f6bdd9e9841f099fa22156bb7b1
completed_at: 2026-10-09T16:29:40Z
release_at: 2026-10-09T16:29:40Z
---
# Routine Operator Work Is Never a REQ; a Tracked Operator Configuration Is a Blocked, Dependency-Gated REQ That Clarify Completes
## What
Capture, clarify and the run orchestrator get one small rule for work the operator does with their own hands. A routine operator act after the code ships (deploy, publish, approve, verify on live hosts) is never captured; the user does it on their own time, and only the AI-buildable parts become REQs. When the request itself asks to track a special operator configuration, capture that one REQ as `status: blocked` with `blocked_by` naming the operator, `blocked_at`, and `depends_on` naming the AI-buildable REQs from the same request, with the checklist in its body. The board already shows such a REQ under Pending → Waiting while a dependency is unmet and under Needs input · Blocked once all are met. Clarify's blocked prompt gains a fourth option, "Done — I did it myself", which records the receipts and flips the REQ to `completed`, so it shows under Done and `cleanup` archives it. A run that discovers mid-flight that only a routine operator act remains completes the REQ with what the AI built and names the operator step in the hand-back; it does not flip to `blocked`, does not fail, and does not recommend abandon.
## Why
Upstream proposal "operator work is not a REQ" (pasted 2026-10-09, verified against 0.305.84 as F11): one consumer REQ, "install, enable and verify CMS publishing in production", touched 53 commits over six days with three block/clarify/earmark round trips and two multi-hour sittings, and the user cancelled it: "I don't need a REQ for it, I'll publish on my own, these REQs are for things that still need to implemented by ai-coders." The proposal asked for "no REQ, a report-only line, and an abandon recommendation". The maintainer rejected that shape and said: "operator action needs to show up in the Needs input · Blocked column, when the code is ready to be released (but not before, before it should have dependencies that might block it, in which case it should be in pending column)"; "when the operator section is done, the REQ should be moved to the DONE column"; and "don't open operator tasks everytime, I'll deploy on my own time, that is nothing special, unless special configuration needs to be done (again this is prose that mostly this orchestrator skill should not concern itself too much)".
## Verified Facts (from triage)
- `skills/do-work/actions/capture.md:107` (Earmark assessment) routes "I need to be at the keyboard for this one" to a `blocked` capture; `:108` (External-condition assessment) defines `blocked` as a condition the work cannot start without, with examples that are all preconditions to AI work. Neither says what to do when the operator act is the whole remaining work.
- `skills/do-work/actions/work.md:523` (Error Handling row added by REQ-648) flips a queued REQ the orchestrator wants to park on the operator to `blocked`. `skills/do-work/actions/work-reference.md:788` (Failure Classification, Environment row) holds the blocked-flip test with no branch for "the missing thing is the whole remaining work".
- `skills/do-work-board/tools/queue-kanban/model.go:1723` places `status: blocked` with unmet `depends_on` under Pending → Waiting; `:1735` places `blocked` with none under Needs input · Blocked; `completed` goes to Done. `model_test.go:638-647` already pins blocked-with-unmet-deps → Waiting, so no board change is needed.
- `skills/do-work/actions/clarify.md:179-190`: the plain-blocked prompt offers "1. Yes — unblock it  2. Not yet — leave it  3. Abandon this REQ"; "Yes" runs the `unblock` transaction, which returns the REQ to `pending` and the next run claims it. `clarify.md:267` already flips a confirmed `builder_decided: true` follow-up to `completed` in place, and `actions/cleanup.md` Pass 0 archives terminal REQs left in the queue.
- The core contract in `skills/do-work/SKILL.md` ("Capture does not execute": a UR plus one or more REQs) and `capture.md` Philosophy are untouched by this rule: a request whose only content is a routine operator act still gets its UR, with zero REQs only when nothing AI-buildable is in it, which `capture.md:16` already allows as the fold-only shape; say so in the capture summary.
## Detailed Requirements
1. `actions/capture.md:108`, External-condition assessment, two sentences keyed on the condition, not a list: a routine operator act after the code ships (deploy, publish, approve, verify on live hosts are examples, not a list to match) is never captured as a REQ; the user does it on their own time, and only the AI-buildable parts become REQs, with the capture summary naming the operator step left to the user. When the request itself asks to track a special operator configuration, capture that one REQ as `status: blocked`, `blocked_by: '<the operator, in the user's words>'`, `blocked_at`, and `depends_on` listing every AI-buildable REQ minted from the same request, with the checklist or runbook in the REQ body; it waits under Pending → Waiting until those complete, then surfaces under Needs input · Blocked.
2. `actions/capture.md:107`, Earmark assessment: the "at the keyboard" sentence points at the new rule instead of flatly capturing `blocked`.
3. `actions/clarify.md` Step 5.5, plain-blocked prompt: add `4. Done — I did it myself`. On that answer clarify appends a `## Operator receipts` section holding the user's words verbatim (apply Step 4's Outside-text containment), sets `status: completed` and `status_changed_at: <now>` in place, and does not run `unblock`. The summary in Step 6 lists each REQ completed this way. Precedent to cite: the confirmed `builder_decided: true` path at `clarify.md:267`, and `actions/cleanup.md` Pass 0 for the archive move. No new CLI verb.
4. `actions/work.md:523` row and `actions/work-reference.md:788` Environment row, one sentence each: when a run finds that the remaining work is a routine operator act, the REQ completes with what the AI built and the hand-back names the operator step; it is not flipped to `blocked`, not failed, and abandon is not recommended. The blocked flip stays for a precondition to AI work and for a tracked special configuration.
5. `docs/work-guide.md`, the Earmarking or blocked paragraph: one sentence on the shape and where it shows on the board (Pending → Waiting, then Needs input · Blocked, then Done).
6. Acceptance greps: `grep -n "own time" skills/do-work/actions/capture.md` finds the routine-act sentence; `grep -n "Done — I did it myself" skills/do-work/actions/clarify.md` finds the option; `grep -n "operator" skills/do-work/actions/work.md skills/do-work/actions/work-reference.md` finds the two inverse sentences.
7. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` green. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Prose only. No new field, status, flag, transaction, probe or board change. `blocked`, `depends_on`, `unblock` and the default scan keep their meaning.
- Keep it short. The maintainer said the orchestrator skill "should not concern itself too much" with operator work: two sentences in capture, one option in clarify, one sentence each in the two run rows, one in the guide. Delete before you add where an existing sentence already says part of it.
- Record the challenge in the Decisions of this REQ: the proposal's "no REQ, report-only line, abandon" shape was rejected by the maintainer in favour of the blocked + `depends_on` placement the board already has and a clarify completion path.
- Read `_dev/primes/prime-action-files.md` before editing; cross-references use that prime's spelling.
- Batch constraint: this REQ is its own release and its own commit; REQ-653 depends on it and must not be folded into it.
## Dependencies
None upstream. REQ-653 (fan-out prose to a companion reference) depends on this REQ because both edit `actions/work.md` and `actions/work-reference.md`.
## Builder Guidance
Certainty is high on the shape (the maintainer stated each column and the Done move). Latitude: exact wording, and whether the clarify Done path writes the receipts section before or after the status line.
## Red-Green Proof
**RED prompt/case:** A capture whose request ends "then deploy it to production and verify on the live host" and a clarify session on a `blocked` REQ whose operator work is finished.
**Why RED now:** Capture mints a `blocked` operator REQ (capture.md:107-108) with no `depends_on`; clarify offers only unblock, leave, or abandon, so the finished REQ goes back to `pending`, the run claims it, finds nothing to build, and flips it to `blocked` again (work.md:523).
**GREEN when:** Capture writes no REQ for the routine deploy and says so in its summary; a request that asks to track a special configuration yields one `blocked` REQ with `depends_on` on the AI-buildable REQs; clarify's prompt has the Done option and leaves a `completed` REQ with `## Operator receipts`, which the board shows under Done and `cleanup` archives; the acceptance greps in requirement 6 all hit; both contract tests are green.
**Validation:** User adjusted. The maintainer replaced the proposal's shape with the blocked, dependency-gated, Done-column shape in two messages, quoted verbatim in the UR.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (6855 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` (a rule added in one action stales sibling restatements) is the risk when capture, clarify, work and the guide each carry part of this rule.
## Full Context
See `do-work/user-requests/UR-143/input.md` for complete verbatim input (the full proposal, the maintainer's answers and refinements, and the triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read brief, REQ, UR-143, crew rules, both primes and every `alternate-writer-contract-drift` lesson. Edit only the five named files, delete before adding: rewrite the existing "at the keyboard" clause in capture and the operator sentence in the guide rather than add beside them. Key the routine-act rule on its condition, with example verbs marked as examples. Swept `skills/` for other restatements of the operator-blocked rule (lesson REQ-566/REQ-648); found the `assigned_to` schema line in work-reference and the board action's Needs-input paragraph, both still accurate (see D-04).
- [x] **[APPLY]:** Five prose edits as planned, one commit. No Go, field, status, flag, transaction, probe or board change. All existing bold labels and headings unchanged.
- [x] **[UNIFY]:** `git diff --stat HEAD~1`: 5 files changed, 8 insertions(+), 6 deletions(-). `git diff --check` clean. `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 (2 s). `bash _dev/tests/contract-regressions.sh` exit 0 (30 s). Each file read back: capture (both assessments, cross-ref to clarify Step 5.5), clarify (prompt block, Done bullet with cross-ref `actions/cleanup.md` → **Pass 0: Sweep Finished Queue Items**, Step 6), work.md row 523, work-reference row 788, work-guide line 119. No debug artifacts.
*Source: upstream proposal "operator work is not a REQ" (verified by `do-work-toolbox validate-feedback` on 2026-10-09 as F11), reshaped by the maintainer: "operator action needs to show up in the Needs input · Blocked column, when the code is ready to be released …", "when the operator section is done, the REQ should be moved to the DONE column", "don't open operator tasks everytime, I'll deploy on my own time … unless special configuration needs to be done".*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the five files and the exact lines, and states each sentence to add; prose only, no Go, no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** Capture's External-condition assessment now says a routine operator act after the code ships is never a REQ (the capture summary names the step) and a tracked special operator configuration is one `blocked` REQ with `depends_on` on the AI-buildable REQs; the Earmark assessment's "at the keyboard" clause points at that rule. Clarify's plain-blocked prompt gains `4. Done — I did it myself`, which writes `## Operator receipts` and flips the REQ to `completed` in place without `unblock`, and Step 6 reports such REQs; the two run rows (work.md Error Handling, work-reference.md Environment) complete a REQ whose remaining work is a routine operator act instead of flipping it to `blocked`; the work guide's earmark paragraph states the three shapes and the board path. Merged from the builder branch in 92a6951c, then re-merged in e8b3687c after the review fix 0a1a8da0 (the Done flip stamps `completed_at`, not `status_changed_at`); cumulative range fa288233..e8b3687c, 5 files, +8/-6.

## Decisions

*(from the builder hand-back, `do-work/runs/work-2026-10-09-160815/REQ-652-handback.md`, verbatim)*

- **D-01 (DECIDE & STATE, recorded challenge):** The upstream proposal's "no REQ, report-only line, recommend abandon" shape was rejected by the maintainer in favour of the blocked + `depends_on` placement the board already has (model.go Waiting / Needs input · Blocked) and a clarify completion path into Done. Built the maintainer's shape only.
- **D-02 (DECIDE & STATE):** A request whose only content is a routine operator act writes no UR and no REQ ("none means no UR and no REQ"). The REQ's Verified Facts say `capture.md:16` already allows a zero-REQ UR, but that exception covers only a capture whose every request was folded; widening it would change the core contract the REQ says stays untouched. Value: no orphan UR, no change to verify-requests linkage. Risk: the verbatim input of an operator-only request is not stored; reversible by one clause if the maintainer wants the UR kept.
- **D-03 (DECIDE & STATE):** The clarify precedent is cited as "the confirmed `builder_decided: true` path ... in Step 5" (the operative bullet is at line 135; line 267 is its checklist line). No Verification Checklist line added: the existing checklist has no blocked-prompt items to sit beside, so a line would be ceremony.
- **D-04 (DECIDE & STATE):** Left unchanged: the `assigned_to` schema line in work-reference ("work that waits on the user as operator is `status: blocked`") and the board action's Needs input paragraph. Both stay true under the new rule (a routine act no longer waits as a REQ; a tracked configuration still is `blocked`), and both are outside the write set.
- **D-05 (DECIDE & STATE):** work-guide keeps a short sentence for the "at the keyboard" precondition case so the rewrite did not drop it.

- **D-02 coordinator ruling:** accepted as built. A request holding only a routine operator act is not queue intent, so capture declines it the way it already declines a duplicate (capture.md Step 2 "If same: tell user, skip" writes nothing either); the user's words stay in the conversation. The fold-only UR exception is not widened.
- **D-06 (DECIDE & STATE, integrator, from review finding F1):** clarify's Done flip stamps `completed_at: <now>`, not the `status_changed_at` that Detailed Requirement 3 named. The schema (`actions/work-reference.md` → Request File Schema) requires `completed_at` on every terminal flip, and the board resolves a terminal card's Done placement only from `completed_at` or a commit hash; with `status_changed_at` alone the card lands under Completion anomalies, which misses the maintainer's "moved to the DONE column". Fixed on the builder branch (0a1a8da0) and re-merged (e8b3687c).

## Discovered Tasks

*(from the builder hand-back, verbatim)*

- impact-negligible: REQ-652 Verified Facts misstate `capture.md:16` as allowing a zero-REQ UR for non-fold captures; only matters if a later REQ relies on that claim → report only

## Qualification

**Gate records:** `advance --diff-range fa288233..92a6951c` returned the `qualify` gate `state: satisfied`, `provenance: merged_range`, no findings (no debug artifact, no output primitive, P-A-U boxes ticked with the builder's text).

**Requirement trace against the merged diff** (`git diff fa288233..92a6951c --stat`: 5 files, +8/-6):
1. capture.md External-condition assessment: the routine-act sentence ("examples, not a list", "on their own time", capture summary names the step) and the tracked-configuration sentence (`blocked`, `blocked_by`, `blocked_at`, `depends_on` on every AI-buildable REQ, checklist in body, Pending → Waiting then Needs input · Blocked, clarify completes it) are both present. Met.
2. capture.md Earmark assessment: the "at the keyboard" clause now points at the operator rule instead of flatly capturing `blocked`. Met.
3. clarify.md Step 5.5: option `4. Done — I did it myself`; the Done bullet writes `## Operator receipts` under Step 4's Outside-text containment (heading exists, clarify.md:101), sets `completed` + `status_changed_at` in place, runs no `unblock`, cites the Step 5 `builder_decided: true` path (clarify.md:135) and `actions/cleanup.md` → **Pass 0: Sweep Finished Queue Items** (cleanup.md:35, exists). Step 6 lists REQs completed by the operator. Met.
4. work.md row 523 and work-reference.md row 788: each gains one sentence completing the REQ for a routine operator act, with the blocked flip kept for a precondition and a tracked configuration (work-reference wording; work.md keeps its existing flip sentence for the park case). Met.
5. work-guide.md earmark paragraph: states the three shapes and the board path Pending → Waiting, Needs input · Blocked, Done. Met.
6. Acceptance greps: rerun by the probe at the test gate.
7. Contract tests: the builder ran both green; the repository gate below reruns them on the merged tree.

**Scope (Route A, judged against `write_set`):** touched files equal the five `write_set` entries exactly; no `do-work/` path inside the range (queue guard empty before the merge); no Go, field, status, flag or transaction change. D-02 deviates from the REQ's Verified Facts (zero-REQ UR) by writing no UR; recorded as a decision and ruled on by the coordinator, so not drift.

**Review-fix re-merge:** the fix 0a1a8da0 changes one line of `skills/do-work/actions/clarify.md`, inside `write_set`; the cumulative range fa288233..e8b3687c still touches exactly the five `write_set` files (+8/-6). `advance` did not re-run the qualify gate (the REQ was already past it, `ADVANCE-GATE-INPUT-IRRELEVANT`); the delta was read by hand: no debug artifact, the new cross-reference `actions/work-reference.md` → Request File Schema lands on the `## Request File Schema — Full Frontmatter` heading. Requirement 3 now stamps `completed_at` instead of `status_changed_at` (D-06).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 92a6951c (repository gate, run directly, unpiped); then `advance REQ-652 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-09-160815/REQ-652-probe.sh` (the four acceptance greps plus `_dev/tests/shipped-package-reference-contract.sh`).
**Result:** ✓ All passing. Repository gate exit 0, wall 137 s (includes `contract-regressions.sh` and the reference-contract test; both Go modules re-executed on fingerprint mismatch and passed). Probe status 0 (`BLOCKED-PROBE-SUCCEEDED`), `green-gate` satisfied.

**Rerun at the review-fix merge e8b3687c:** repository gate exit 0, wall 151 s; probe status 0 (run directly; `advance` was past the test gate).

**Regression evidence (non-behavioral prose change, `tdd: false`):** the Red-Green Proof is a prose read. RED: before the merge, capture.md routed "at the keyboard" straight to `blocked` and had no routine-act rule, clarify offered only unblock/leave/abandon, and both run rows flipped to `blocked`. GREEN: the four acceptance greps (requirement 6) hit on the merged tree, and the merged diff carries each sentence the GREEN condition names (traced in `## Qualification`).

**New tests added:** none (prose-only change; the probe is run evidence, not a committed test).

**Heavy verification plan:**
- Range: fa288233683ec423c33a6055daafd118f9195573..e8b3687cd31c0f6bdd9e9841f099fa22156bb7b1 (replanned at the review-fix merge; first plan at 92a6951c selected the same lane)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed paths matched subtree `skills`

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-09T16:28:49Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 clarify.md:194 Done flip stamped `status_changed_at` instead of `completed_at`, so the board flagged a completion anomaly and never showed the REQ under Done. RESOLVED by 0a1a8da0 (now stamps `completed_at`, citing Request File Schema). — impact-user-visible → report only

**Minor findings:** F2 clarify.md:194 misdescribed the Step 5 `builder_decided: true` path as "the same in-place flip", RESOLVED by 0a1a8da0 (now "Like the confirmed `builder_decided: true` path in Step 5, it completes a REQ that needs no build") — impact-negligible → report only; F3 a Done-completed REQ shows doctor `STRANDED-TERMINAL-REQUEST` until cleanup Pass 0, which is expected, open — impact-negligible → report only; F4 D-02 adds a second unlisted exception to capture.md:10 / SKILL.md:20 "always a UR" next to the existing duplicate skip, accepted, open — impact-negligible → report only
**Acceptance:** Pass — capture, run rows and guide match board placement, and the clarify Done path stamps `completed_at` and lands under Done (model.go code-path trace at e8b3687c)
**Restatement sweep:** redefined the operator-blocked rule (routine act never captured, tracked configuration blocked with depends_on, clarify Done → completed, run completes instead of flipping); swept skills/: work-reference.md:114 `assigned_to` line agrees, do-work-board/actions/board.md:92 Needs input paragraph agrees, board-cards.js:240 tooltip agrees, work-guide.md:102 and capture-guide.md:63 not stale; the one disagreement (clarify Done stamping vs work-reference.md:172-191) was resolved by 0a1a8da0
**Suggested testing:** 1 items
**Re-review:** first pass at 92a6951c scored 82% (Partial) on F1; fixed on the builder branch in 0a1a8da0, re-merged e8b3687c, gate re-run green, delta re-reviewed to the score above. Full report: `do-work/runs/work-2026-10-09-160815/REQ-652-review.md`.
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Sweeping `skills/` for restatements of the operator-blocked rule before editing; the `assigned_to` schema line, the board's Needs input paragraph and the badge tooltip all still agreed (D-04), so the change stayed at five files.
**What didn't:** The REQ's Detailed Requirement 3 named `status_changed_at` for clarify's Done flip, and the builder followed it. That stamp is for flips without a dedicated stamp; a terminal flip needs `completed_at`, and without it (or a commit hash) the board files the card under Completion anomalies instead of Done. Review caught it (F1); fixed in 0a1a8da0.
**Worth knowing:** Any new path that sets a terminal status outside the finalization transaction must stamp `completed_at` itself (`actions/work-reference.md` → Request File Schema, the STAMPING RULE on `completed_at`). Check the field the downstream reader resolves from, not the field the REQ text names.

## Orientation

Now capture keeps routine operator work out of the queue, a tracked operator configuration waits under Pending → Waiting behind its code, and clarify's "Done — I did it myself" completes it; lives in the capture / clarify / work action prose (`_dev/primes/prime-action-files.md`). Leaf change, no map change.

## Heavy Verification Plan

- Base: fa288233683ec423c33a6055daafd118f9195573
- Target: e8b3687cd31c0f6bdd9e9841f099fa22156bb7b1
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed paths matched subtree `skills`

## Heavy Verification Result

- Target revision: e8b3687cd31c0f6bdd9e9841f099fa22156bb7b1
- Execution revision: e8b3687cd31c0f6bdd9e9841f099fa22156bb7b1 (detached checkout at the merge)
- staged-skills: exit 0, executed (fingerprint_mismatch), wall 48 s

## Timing

Observed 2026-10-09T16:11:57Z to 2026-10-09T16:29:11Z: 17m 14s total, 19m 52s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 7m 44s | 1 |
| verification-gate | 6m 50s | 3 |
| builder-work | 5m 00s | 1 |
| handback-merge | 18s | 1 |

Slowest stage: review / review, fix F1, re-merge, re-review, 7m 44s, outcome success.
Slowest command: verification-gate / repository gate (137 s) plus probe and heavy plan, 3m 05s, exit 0, .
